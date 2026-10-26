package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixGeoVPN() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixGeoVPNCreate,
		Read:   resourceAviatrixGeoVPNRead,
		Update: resourceAviatrixGeoVPNUpdate,
		Delete: resourceAviatrixGeoVPNDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},

		Schema: map[string]*schema.Schema{
			"cloud_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{goaviatrix.AWS}),
				Description:  "Type of cloud service provider, requires an integer value. Currently only AWS(1) is supported.",
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
	client := mustClient(meta)

	geoVPN := &goaviatrix.GeoVPN{
		CloudType:   getInt(d, "cloud_type"),
		AccountName: getString(d, "account_name"),
		ServiceName: getString(d, "service_name"),
		DomainName:  getString(d, "domain_name"),
	}

	log.Printf("[INFO] Enabling Aviatrix Geo VPN: %#v", geoVPN)

	elbDNSNames := make([]string, 0)
	for _, elbDNSName := range getList(d, "elb_dns_names") {
		elbDNSNames = append(elbDNSNames, mustString(elbDNSName))
	}

	if len(elbDNSNames) == 0 {
		return fmt.Errorf("please specify 'elb_dns_names' to enable Aviatrix Geo VPN")
	}

	d.SetId(geoVPN.ServiceName + "~" + geoVPN.DomainName)
	flag := false
	defer func() { _ = resourceAviatrixGeoVPNReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	geoVPN.ElbDNSName = elbDNSNames[0]
	err := client.EnableGeoVPN(context.Background(), geoVPN)
	if err != nil {
		return fmt.Errorf("failed to enable Aviatrix Geo VPN due to: %w", err)
	}

	for i := 1; i < len(elbDNSNames); i++ {
		geoVPN.ElbDNSName = elbDNSNames[i]
		err := client.AddElbToGeoVPN(geoVPN)
		if err != nil {
			return fmt.Errorf("failed to add elb: %s to Aviatrix Geo VPN due to: %w", geoVPN.ElbDNSName, err)
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
	client := mustClient(meta)

	domainName := getString(d, "domain_name")
	serviceName := getString(d, "service_name")
	if domainName == "" || serviceName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no domain name or service name received. Import id is %s", id)
		mustSet(d, "cloud_type", goaviatrix.AWS)
		mustSet(d, "service_name", strings.Split(id, "~")[0])
		mustSet(d, "domain_name", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	geoVPN := &goaviatrix.GeoVPN{
		CloudType:   getInt(d, "cloud_type"),
		ServiceName: getString(d, "service_name"),
		DomainName:  getString(d, "domain_name"),
	}

	geoVPNDetail, err := client.GetGeoVPNInfo(geoVPN)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("unable to read Aviatrix Geo VPN Info due to %w", err)
	}
	mustSet(d, "cloud_type", geoVPNDetail.CloudType)
	mustSet(d, "account_name", geoVPNDetail.AccountName)
	mustSet(d, "service_name", geoVPNDetail.ServiceName)
	mustSet(d, "domain_name", geoVPNDetail.DomainName)
	if err := d.Set("elb_dns_names", strings.Split(geoVPNDetail.ElbDNSName, ",")); err != nil {
		log.Printf("[WARN] Error setting 'elb_dns_names' for (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceAviatrixGeoVPNUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Aviatrix Geo VPN")

	client := mustClient(meta)

	geoVPN := &goaviatrix.GeoVPN{
		CloudType:   getInt(d, "cloud_type"),
		ServiceName: getString(d, "service_name"),
		DomainName:  getString(d, "domain_name"),
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
		oldString := mustSlice(oldElb)
		newString := mustSlice(newElb)
		oldElbList := goaviatrix.ExpandStringList(oldString)
		newElbList := goaviatrix.ExpandStringList(newString)
		toDeleteElbs = goaviatrix.Difference(oldElbList, newElbList)
		toAddElbs = goaviatrix.Difference(newElbList, oldElbList)

		for i := 0; i < len(toDeleteElbs); i++ {
			geoVPN.ElbDNSName = toDeleteElbs[i]
			err := client.DeleteElbFromGeoVPN(geoVPN)
			if err != nil {
				return fmt.Errorf("failed to delete ELB: %s from Aviatrix Geo VPN due to: %w", toDeleteElbs[i], err)
			}
		}

		for i := 0; i < len(toAddElbs); i++ {
			geoVPN.ElbDNSName = toAddElbs[i]
			err := client.AddElbToGeoVPN(geoVPN)
			if err != nil {
				return fmt.Errorf("failed to add ELB: %s to Aviatrix Geo VPN due to: %w", toAddElbs[i], err)
			}
		}
	}

	d.Partial(false)
	return resourceAviatrixGeoVPNRead(d, meta)
}

func resourceAviatrixGeoVPNDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	geoVPN := &goaviatrix.GeoVPN{
		CloudType: getInt(d, "cloud_type"),
	}

	log.Printf("[INFO] Disabling Aviatrix Geo VPN: %#v", geoVPN)

	err := client.DisableGeoVPN(geoVPN)
	if err != nil {
		return fmt.Errorf("failed to disable Aviatrix Geo VPN due to: %w", err)
	}

	return nil
}
