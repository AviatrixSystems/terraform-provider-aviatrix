package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixAwsTgwVpnConn() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixAwsTgwVpnConnCreate,
		Read:   resourceAviatrixAwsTgwVpnConnRead,
		Delete: resourceAviatrixAwsTgwVpnConnDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"tgw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "This parameter represents the name of an AWS TGW.",
			},
			"route_domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of a route domain, to which the vpn will be attached.",
			},
			"connection_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unique name of the connection.",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Public IP address. Example: '40.0.0.0'.",
			},
			"remote_as_number": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "AWS side as a number. Integer between 1-65535. Example: '12'.",
			},
			"remote_cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Remote CIDRs joined as a string with ','.",
			},
			"vpn_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the vpn connection.",
			},
		},
	}
}

func resourceAviatrixAwsTgwVpnConnCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	awsTgwVpnConn := &goaviatrix.AwsTgwVpnConn{
		TgwName:         d.Get("tgw_name").(string),
		ConnName:        d.Get("connection_name").(string),
		PublicIP:        d.Get("public_ip").(string),
		RouteDomainName: d.Get("route_domain_name").(string),
	}

	if awsTgwVpnConn.RouteDomainName != "Default_Domain" {
		return fmt.Errorf("invalid 'route_domain_name'. only 'Default_Domain' is supported")
	}

	remoteAsn := d.Get("remote_as_number").(string)
	remoteCIDR := d.Get("remote_cidr").(string)

	if remoteAsn != "" && remoteCIDR != "" {
		return fmt.Errorf("remote_asn and remote_cidr cannot be set at the same time. Please choose: " +
			"remote_asw or remote_cidr")
	}

	if remoteAsn != "" {
		awsTgwVpnConn.OnpremASN = remoteAsn
	}
	if remoteCIDR != "" {
		awsTgwVpnConn.RemoteCIDR = remoteCIDR
	}

	log.Printf("[INFO] Creating Aviatrix AWS TGW VPN Connection: %#v", awsTgwVpnConn)

	err := client.CreateAwsTgwVpnConn(awsTgwVpnConn)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix AWS TGW VPN Connection: %s", err)
	}

	d.SetId(awsTgwVpnConn.TgwName + "~" + awsTgwVpnConn.ConnName)

	return resourceAviatrixAwsTgwVpnConnRead(d, meta)
}

func resourceAviatrixAwsTgwVpnConnRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	tgwName := d.Get("tgw_name").(string)
	connName := d.Get("connection_name").(string)

	if tgwName == "" || connName == "" {
		id := d.Id()

		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)

		if !strings.Contains(id, "~") {
			log.Printf("[DEBUG] Import Id: %s is invalid", id)
		}

		d.Set("tgw_name", strings.Split(id, "~")[0])
		d.Set("connection_name", strings.Split(id, "~")[1])

		d.SetId(id)
	}

	awsTgwVpnConn := &goaviatrix.AwsTgwVpnConn{
		TgwName:  d.Get("tgw_name").(string),
		ConnName: d.Get("connection_name").(string),
	}

	vpnConn, err := client.GetAwsTgwVpnConn(awsTgwVpnConn)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix Aws Tgw VPN Connection: %s", err)
	}
	log.Printf("[INFO] Found Aviatrix Aws Tgw VPN Connection: %#v", vpnConn)

	d.Set("tgw_name", vpnConn.TgwName)
	d.Set("route_domain_name", vpnConn.RouteDomainName)
	d.Set("connection_name", vpnConn.ConnName)
	d.Set("public_ip", vpnConn.PublicIP)
	d.Set("remote_as_number", vpnConn.OnpremASN)
	d.Set("remote_cidr", vpnConn.RemoteCIDR)
	d.Set("vpn_id", vpnConn.VpnID)

	d.SetId(vpnConn.TgwName + "~" + vpnConn.ConnName)
	return nil
}

func resourceAviatrixAwsTgwVpnConnDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	awsTgwVpnConn := &goaviatrix.AwsTgwVpnConn{
		TgwName:  d.Get("tgw_name").(string),
		VpnID:    d.Get("vpn_id").(string),
		ConnName: d.Get("connection_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix aws_tgw_vpn_conn: %#v", awsTgwVpnConn)

	err := client.DeleteAwsTgwVpnConn(awsTgwVpnConn)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return nil
		}
		return fmt.Errorf("failed to delete Aviatrix AwsTgwVpnConn: %s", err)
	}

	return nil
}
