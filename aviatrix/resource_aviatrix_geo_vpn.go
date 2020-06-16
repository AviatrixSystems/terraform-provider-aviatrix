package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixGeoVPN() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixGeoVPNCreate,
		Read:   resourceAviatrixGeoVPNRead,
		Update: resourceAviatrixGeoVPNUpdate,
		Delete: resourceAviatrixGeoVPNDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cloud_type": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Type of cloud service provider, requires an integer value. Currently only AWS(1) is supported.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "This parameter represents the name of a Cloud-Account in Aviatrix controller.",
			},
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The hostname that users will connect to. A DNS record will be created for this name in the specified domain name.",
			},
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The hosted domain name. It must be hosted by AWS Route53 or Azure DNS in the selected account.",
			},
			"elb_dns_names": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "List of ELB names to attach to this Geo VPN name.",
			},
		},
	}
}

func resourceAviatrixGeoVPNCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	geoVPN := &goaviatrix.GeoVPN{
		CloudType:   d.Get("cloud_type").(int),
		AccountName: d.Get("account_name").(string),
		ServiceName: d.Get("service_name").(string),
		DomainName:  d.Get("domain_name").(string),
	}

	log.Printf("[INFO] Enabling Aviatrix Geo VPN: %#v", geoVPN)

	elbDNSNames := make([]string, 0)
	for _, elbDNSName := range d.Get("elb_dns_names").([]interface{}) {
		elbDNSNames = append(elbDNSNames, elbDNSName.(string))
	}

	if len(elbDNSNames) == 0 {
		return fmt.Errorf("please specify 'elb_dns_names' to enable Aviatrix Geo VPN")
	}

	geoVPN.ElbDNSName = elbDNSNames[0]
	err := client.EnableGeoVPN(geoVPN)
	if err != nil {
		return fmt.Errorf("failed to enable Aviatrix Geo VPN due to: %s", err)
	}

	d.SetId(geoVPN.ServiceName + "~" + geoVPN.DomainName)

	flag := false
	defer resourceAviatrixGeoVPNReadIfRequired(d, meta, &flag)

	for i := 1; i < len(elbDNSNames); i++ {
		geoVPN.ElbDNSName = elbDNSNames[i]
		err := client.AddElbToGeoVPN(geoVPN)
		if err != nil {
			return fmt.Errorf("failed to add elb: %s to Aviatrix Geo VPN due to: %s", geoVPN.ElbDNSName, err)
		}
	}

	return resourceAviatrixGeoVPNReadIfRequired(d, meta, &flag)
}

func resourceAviatrixGeoVPNReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixGeoVPNRead(d, meta)
	}
	return nil
}

func resourceAviatrixGeoVPNRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	domainName := d.Get("domain_name").(string)
	serviceName := d.Get("service_name").(string)
	if domainName == "" || serviceName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no domain name or service name received. Import id is %s", id)
		d.Set("cloud_type", 1)
		d.Set("service_name", strings.Split(id, "~")[0])
		d.Set("domain_name", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	geoVPN := &goaviatrix.GeoVPN{
		CloudType:   d.Get("cloud_type").(int),
		ServiceName: d.Get("service_name").(string),
		DomainName:  d.Get("domain_name").(string),
	}

	geoVPNDetail, err := client.GetGeoVPNInfo(geoVPN)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("unable to read Aviatrix Geo VPN Info due to %v", err)
	}

	d.Set("cloud_type", geoVPNDetail.CloudType)
	d.Set("account_name", geoVPNDetail.AccountName)
	d.Set("service_name", geoVPNDetail.ServiceName)
	d.Set("domain_name", geoVPNDetail.DomainName)
	if err := d.Set("elb_dns_names", geoVPNDetail.ElbDNSNames); err != nil {
		log.Printf("[WARN] Error setting 'elb_dns_names' for (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceAviatrixGeoVPNUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Aviatrix Geo VPN")

	client := meta.(*goaviatrix.Client)

	geoVPN := &goaviatrix.GeoVPN{
		CloudType:   d.Get("cloud_type").(int),
		ServiceName: d.Get("service_name").(string),
		DomainName:  d.Get("domain_name").(string),
	}

	var toDeleteElbs []string
	var toAddElbs []string

	d.Partial(true)
	if d.HasChange("elb_dns_names") {
		oldElb, newElb := d.GetChange("elb_dns_names")
		if oldElb == nil {
			oldElb = new([]interface{})
		}
		if newElb == nil {
			newElb = new([]interface{})
		}
		oldString := oldElb.([]interface{})
		newString := newElb.([]interface{})
		oldElbList := goaviatrix.ExpandStringList(oldString)
		newElbList := goaviatrix.ExpandStringList(newString)
		toDeleteElbs = goaviatrix.Difference(oldElbList, newElbList)
		toAddElbs = goaviatrix.Difference(newElbList, oldElbList)

		for i := 0; i < len(toDeleteElbs); i++ {
			geoVPN.ElbDNSName = toDeleteElbs[i]
			err := client.DeleteElbFromGeoVPN(geoVPN)
			if err != nil {
				return fmt.Errorf("failed to delete ELB: %s from Aviatrix Geo VPN due to: %s", toDeleteElbs[i], err)
			}
		}

		for i := 0; i < len(toAddElbs); i++ {
			geoVPN.ElbDNSName = toAddElbs[i]
			err := client.AddElbToGeoVPN(geoVPN)
			if err != nil {
				return fmt.Errorf("failed to add ELB: %s to Aviatrix Geo VPN due to: %s", toAddElbs[i], err)
			}
		}
		d.SetPartial("elb_dns_names")
	}

	return resourceAviatrixGeoVPNRead(d, meta)
}

func resourceAviatrixGeoVPNDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	geoVPN := &goaviatrix.GeoVPN{
		CloudType: d.Get("cloud_type").(int),
	}

	log.Printf("[INFO] Disabling Aviatrix Geo VPN: %#v", geoVPN)

	err := client.DisableGeoVPN(geoVPN)
	if err != nil {
		return fmt.Errorf("failed to disable Aviatrix Geo VPN due to: %s", err)
	}

	return nil
}
