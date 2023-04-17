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
		UpdateWithoutTimeout: resourceAviatrixEdgeNEODeviceOnboardingUpdate,
		DeleteWithoutTimeout: resourceAviatrixEdgeNEODeviceOnboardingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Edge NEO account name.",
			},
			"device_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Device name.",
			},
			"serial_number": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Serial number.",
			},
			"hardware_model": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Hardware Model.",
			},
			"device_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Device ID.",
			},
			"network": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Network configurations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interface_name": {
							Type:        schema.TypeString,
							Required:    true,
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
						"ipv4_cidr": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IPV4 CIDR.",
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
	for _, network0 := range network {
		network1 := network0.(map[string]interface{})

		network2 := &goaviatrix.EdgeNEODeviceNetwork{
			InterfaceName: network1["interface_name"].(string),
			EnableDhcp:    network1["enable_dhcp"].(bool),
			GatewayIp:     network1["gateway_ip"].(string),
			Ipv4Cidr:      network1["ipv4_cidr"].(string),
			DnsServerIp:   network1["dns_server_ip"].(string),
			ProxyServerIp: network1["proxy_server_ip"].(string),
		}

		edgeNEODevice.Network = append(edgeNEODevice.Network, network2)
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
	for _, network0 := range edgeNEODeviceResp.Network {
		network1 := make(map[string]interface{})
		network1["interface_name"] = network0.InterfaceName
		network1["enable_dhcp"] = network0.EnableDhcp
		network1["gateway_ip"] = network0.GatewayIp
		network1["ipv4_cidr"] = network0.Ipv4Cidr
		network1["dns_server_ip"] = network0.DnsServerIp
		network1["proxy_server_ip"] = network0.ProxyServerIp

		network = append(network, network1)
	}

	if err = d.Set("network", network); err != nil {
		return diag.Errorf("failed to set network: %s\n", err)
	}

	d.SetId(accountName + "~" + deviceName)
	return nil
}

func resourceAviatrixEdgeNEODeviceOnboardingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	edgeNEODevice := marshalEdgeNEODeviceOnboardingInput(d)

	if d.HasChanges("account_name", "device_name", "serial_number", "hardware_model") {
		return diag.Errorf("account_name, device_name, serial_number and hardware_model are not allowed to be updated")
	}

	d.Partial(true)

	if d.HasChange("network") {
		if err := client.OnboardEdgeNEODevice(ctx, edgeNEODevice); err != nil {
			return diag.Errorf("could not update network configurations during Edge NEO device update: %v", err)
		}
	}

	d.Partial(false)
	return resourceAviatrixEdgeNEODeviceOnboardingRead(ctx, d, meta)
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
