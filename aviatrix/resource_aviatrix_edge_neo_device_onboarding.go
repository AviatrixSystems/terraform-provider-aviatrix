package aviatrix

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixEdgeNEODeviceOnboarding() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixEdgeNEODeviceOnboardingCreate,
		ReadWithoutTimeout:   resourceAviatrixEdgeNEODeviceOnboardingRead,
		DeleteWithoutTimeout: resourceAviatrixEdgeNEODeviceOnboardingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Edge CSP account name.",
			},
			"device_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Edge CSP name.",
			},
			"serial_number": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: ".",
			},
			"hardware_model": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: ".",
			},
			"device_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Edge NEO device ID.",
			},
			"network": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Description: "Network configurations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interface_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Interface name.",
						},
						"enable_dhcp": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Enable DHCP.",
						},
						"gateway_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Gateway IP.",
						},
						"subnet_cidr": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Subnet CIDR.",
						},
						"dns_server_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "DNS server IP.",
						},
						"proxy_server_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Proxy server IP.",
						},
						"tags": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tags.",
						},
					},
				},
			},
		},
	}
}

func marshalEdgeNEODeviceOnboardingInput(d *schema.ResourceData) *goaviatrix.EdgeNEODevice {
	edgeNEODevice := &goaviatrix.EdgeNEODevice{
		AccountName:   d.Get("account_name").(string),
		DeviceName:    d.Get("device_name").(string),
		SerialNumber:  d.Get("serial_number").(string),
		HardwareModel: d.Get("hardware_model").(string),
	}

	network := d.Get("network").(*schema.Set).List()
	for _, nw0 := range network {
		nw1 := nw0.(map[string]interface{})

		nw2 := &goaviatrix.EdgeNEODeviceNetwork{
			InterfaceName: nw1["interface_name"].(string),
			EnableDhcp:    nw1["enable_dhcp"].(bool),
			GatewayIp:     nw1["gateway_ip"].(string),
			SubnetCidr:    nw1["subnet_cidr"].(string),
			DnsServerIp:   nw1["dns_server_ip"].(string),
			ProxyServerIp: nw1["proxy_server_ip"].(string),
			Tags:          nw1["tags"].(string),
		}

		edgeNEODevice.Network = append(edgeNEODevice.Network, nw2)
	}

	return edgeNEODevice
}

func resourceAviatrixEdgeNEODeviceOnboardingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	edgeNEODevice := marshalEdgeNEODeviceOnboardingInput(d)

	flag := false
	defer resourceAviatrixEdgeNEODeviceOnboardingReadIfRequired(ctx, d, meta, &flag)

	if err := client.OnboardEdgeNEODevice(ctx, edgeNEODevice); err != nil {
		return diag.Errorf("could not onboard Edge NEO device: %v", err)
	}

	d.SetId(edgeNEODevice.AccountName + "~" + edgeNEODevice.DeviceName)
	return resourceAviatrixEdgeNEODeviceOnboardingReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixEdgeNEODeviceOnboardingReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixEdgeNEODeviceOnboardingRead(ctx, d, meta)
	}
	return nil
}

//func resourceAviatrixEdgeNEODeviceOnboardingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
//	return nil
//}

func resourceAviatrixEdgeNEODeviceOnboardingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	accountName := d.Get("account_name").(string)
	deviceName := d.Get("device_name").(string)
	if accountName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no account name received. Import Id is %s", id)
		parts := strings.Split(id, "~")
		if len(parts) != 32 {
			return diag.Errorf("Invalid Import ID received, ID must be in the format account_name~device_name")
		}
		accountName = parts[0]
		deviceName = parts[1]
		d.SetId(id)
	}

	edgeNEODeviceResp, err := client.GetEdgeNEODevice(ctx, accountName, deviceName)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read Edge NEO device: %v", err)
	}

	d.Set("account_name", accountName)
	d.Set("device_name", edgeNEODeviceResp.DeviceName)
	d.Set("device_id", edgeNEODeviceResp.DeviceId)
	d.Set("serial_number", edgeNEODeviceResp.SerialNumber)
	d.Set("hardware_model", edgeNEODeviceResp.HardwareModel)

	var network []map[string]interface{}
	for _, nw0 := range edgeNEODeviceResp.Network {
		nw1 := make(map[string]interface{})
		nw1["interface_name"] = nw0.InterfaceName
		nw1["enable_dhcp"] = nw0.EnableDhcp
		nw1["gateway_ip"] = nw0.GatewayIp
		nw1["subnet_cidr"] = nw0.SubnetCidr
		nw1["dns_server_ip"] = nw0.DnsServerIp
		nw1["proxy_server_ip"] = nw0.ProxyServerIp
		nw1["tags"] = nw0.Tags

		network = append(network, nw1)
	}

	if err = d.Set("network", network); err != nil {
		return diag.Errorf("failed to set network: %s\n", err)
	}

	d.SetId(accountName + "~" + deviceName)
	return nil
}

func resourceAviatrixEdgeNEODeviceOnboardingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	accountName := d.Get("account_name").(string)
	serialNumber := d.Get("serial_number").(string)

	err := client.DeleteEdgeNEODevice(ctx, accountName, serialNumber)
	if err != nil {
		return diag.Errorf("could not delete Edge NEO device: %v", err)
	}

	return nil
}
