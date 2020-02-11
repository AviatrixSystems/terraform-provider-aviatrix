package aviatrix

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"connection_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "dynamic",
				ForceNew:    true,
				Description: "Connection type. Valid values: 'dynamic', 'static'. 'dynamic' stands for a BGP VPN connection; 'static' stands for a static VPN connection. Default value: 'dynamic'.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "dynamic" && v != "static" {
						errs = append(errs, fmt.Errorf("%q must be either 'dynamic' or 'static', got: %s", key, val))
					}
					return
				},
			},
			"remote_as_number": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "AWS side as a number. Integer between 1-65535. Example: '12'. Required for a dynamic VPN connection.",
			},
			"remote_cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Remote CIDRs joined as a string with ','. Required for a static VPN connection.",
			},
			"inside_ip_cidr_tun_1": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Inside IP CIDR for Tunnel 1. A /30 CIDR in 169.254.0.0/16.",
			},
			"pre_shared_key_tun_1": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				ForceNew:  true,
				Description: "Pre-Shared Key for Tunnel 1. A 8-64 character string with alphanumeric, " +
					"underscore(_) and dot(.). It cannot start with 0",
			},
			"inside_ip_cidr_tun_2": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Inside IP CIDR for Tunnel 2. A /30 CIDR in 169.254.0.0/16.",
			},
			"pre_shared_key_tun_2": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				ForceNew:  true,
				Description: "Pre-Shared Key for Tunnel 2. A 8-64 character string with alphanumeric, " +
					"underscore(_) and dot(.). It cannot start with 0",
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
		TgwName:          d.Get("tgw_name").(string),
		ConnName:         d.Get("connection_name").(string),
		PublicIP:         d.Get("public_ip").(string),
		RouteDomainName:  d.Get("route_domain_name").(string),
		InsideIpCIDRTun1: d.Get("inside_ip_cidr_tun_1").(string),
		InsideIpCIDRTun2: d.Get("inside_ip_cidr_tun_2").(string),
		PreSharedKeyTun1: d.Get("pre_shared_key_tun_1").(string),
		PreSharedKeyTun2: d.Get("pre_shared_key_tun_2").(string),
	}

	connectionType := d.Get("connection_type").(string)
	remoteAsn := d.Get("remote_as_number").(string)
	remoteCIDR := d.Get("remote_cidr").(string)
	if connectionType == "dynamic" && remoteAsn == "" {
		return fmt.Errorf("please specify 'remote_as_number' to create a BGP VPN connection")
	} else if connectionType == "dynamic" && remoteCIDR != "" {
		return fmt.Errorf("please set 'remote_cidr' as empty since it is only requried for a static VPN connection")
	} else if connectionType == "static" && remoteCIDR == "" {
		return fmt.Errorf("please specify 'remote_cidr' to create a static VPN connection")
	} else if connectionType == "static" && remoteAsn != "" {
		return fmt.Errorf("please set 'remote_as_number' as empty since it is only requried for a BGP VPN connection")
	}

	if remoteAsn != "" {
		awsTgwVpnConn.OnpremASN = remoteAsn
	}
	if remoteCIDR != "" {
		awsTgwVpnConn.RemoteCIDR = remoteCIDR
	}

	log.Printf("[INFO] Creating Aviatrix AWS TGW VPN Connection: %#v", awsTgwVpnConn)

	vpnID, err := client.CreateAwsTgwVpnConn(awsTgwVpnConn)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix AWS TGW VPN Connection: %s", err)
	}

	d.SetId(awsTgwVpnConn.TgwName + "~" + vpnID)
	return resourceAviatrixAwsTgwVpnConnRead(d, meta)
}

func resourceAviatrixAwsTgwVpnConnRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	tgwName := d.Get("tgw_name").(string)
	vpnID := d.Get("vpn_id").(string)

	if tgwName == "" || vpnID == "" {
		id := d.Id()

		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)

		if !strings.Contains(id, "~") {
			log.Printf("[DEBUG] Import Id: %s is invalid", id)
		}

		d.Set("tgw_name", strings.Split(id, "~")[0])
		d.Set("vpn_id", strings.Split(id, "~")[1])

		d.SetId(id)
	}

	awsTgwVpnConn := &goaviatrix.AwsTgwVpnConn{
		TgwName: d.Get("tgw_name").(string),
		VpnID:   d.Get("vpn_id").(string),
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
	if vpnConn.OnpremASN != "" && vpnConn.RemoteCIDR == "" {
		d.Set("connection_type", "dynamic")
	} else if vpnConn.OnpremASN == "" && vpnConn.RemoteCIDR != "" {
		d.Set("connection_type", "static")
	}
	d.Set("vpn_id", vpnConn.VpnID)
	d.Set("inside_ip_cidr_tun_1", vpnConn.InsideIpCIDRTun1)
	d.Set("inside_ip_cidr_tun_2", vpnConn.InsideIpCIDRTun2)

	d.SetId(vpnConn.TgwName + "~" + vpnConn.VpnID)
	return nil
}

func resourceAviatrixAwsTgwVpnConnDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	awsTgwVpnConn := &goaviatrix.AwsTgwVpnConn{
		TgwName: d.Get("tgw_name").(string),
		VpnID:   d.Get("vpn_id").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix aws_tgw_vpn_conn: %#v", awsTgwVpnConn)

	err := client.DeleteAwsTgwVpnConn(awsTgwVpnConn)

	time.Sleep(40 * time.Second)

	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return nil
		}

		return fmt.Errorf("failed to delete Aviatrix AwsTgwVpnConn: %s", err)
	}

	return nil
}
