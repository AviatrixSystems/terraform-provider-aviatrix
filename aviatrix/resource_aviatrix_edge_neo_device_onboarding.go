package aviatrix

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

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
						"dns_server_ips": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Set of DNS server IPs.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"proxy_server_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Proxy server IP.",
						},
					},
				},
			},
			"download_config_file": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Set to true to download the Edge NEO static config file.",
			},
			"config_file_download_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The location where the config file will be stored.",
			},
		},
	}
}

func marshalEdgeNEODeviceOnboardingInput(d *schema.ResourceData) *goaviatrix.EdgeNEODevice {
	edgeNEODevice := &goaviatrix.EdgeNEODevice{
		AccountName:            d.Get("account_name").(string),
		DeviceName:             d.Get("device_name").(string),
		SerialNumber:           d.Get("serial_number").(string),
		HardwareModel:          d.Get("hardware_model").(string),
		DownloadConfigFile:     d.Get("download_config_file").(bool),
		ConfigFileDownloadPath: d.Get("config_file_download_path").(string),
	}

	network := d.Get("network").(*schema.Set).List()
	for _, network0 := range network {
		network1 := network0.(map[string]interface{})

		network2 := &goaviatrix.EdgeNEODeviceNetwork{
			InterfaceName: network1["interface_name"].(string),
			EnableDhcp:    network1["enable_dhcp"].(bool),
			GatewayIp:     network1["gateway_ip"].(string),
			Ipv4Cidr:      network1["ipv4_cidr"].(string),
			ProxyServerIp: network1["proxy_server_ip"].(string),
		}

		for _, dnsServerIp := range network1["dns_server_ips"].(*schema.Set).List() {
			network2.DnsServerIps = append(network2.DnsServerIps, dnsServerIp.(string))
		}

		edgeNEODevice.Network = append(edgeNEODevice.Network, network2)
	}

	return edgeNEODevice
}

func resourceAviatrixEdgeNEODeviceOnboardingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	edgeNEODevice := marshalEdgeNEODeviceOnboardingInput(d)

	if edgeNEODevice.DownloadConfigFile && edgeNEODevice.ConfigFileDownloadPath == "" {
		diag.Errorf("config_file_download_path is required when download_config_file is true")
	}
	if !edgeNEODevice.DownloadConfigFile && edgeNEODevice.ConfigFileDownloadPath != "" {
		diag.Errorf("config_file_download_path must be empty when download_config_file is false")
	}

	flag := false
	defer resourceAviatrixEdgeNEODeviceOnboardingReadIfRequired(ctx, d, meta, &flag)

	if err := client.OnboardEdgeNEODevice(ctx, edgeNEODevice); err != nil {
		return diag.Errorf("could not onboard Edge NEO device: %v", err)
	}

	if edgeNEODevice.DownloadConfigFile {
		if err := client.DownloadEdgeNEOConfigFile(ctx, edgeNEODevice); err != nil {
			return diag.Errorf("could not download Edge NEO static config file: %v", err)
		}
	}

	numberOfRetries := 30
	retryInterval := 20
	for i := 0; ; i++ {
		edgeNEODeviceResp, err := client.GetEdgeNEODevice(ctx, edgeNEODevice.AccountName, edgeNEODevice.DeviceName)
		if err != nil {
			return diag.Errorf("could not read Edge NEO device during onboarding: %v", err)
		}

		if edgeNEODeviceResp.ConnectionStatus == "connected" {
			break
		}

		if i < numberOfRetries {
			time.Sleep(time.Duration(retryInterval) * time.Second)
		} else {
			return diag.Errorf("Edge NEO device connection status could not become connected after 10 minutes")
		}
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

	var edgeNEODeviceResp *goaviatrix.EdgeNEODeviceResp
	var err error

	for i := 0; ; i++ {
		edgeNEODeviceResp, err = client.GetEdgeNEODevice(ctx, accountName, deviceName)

		if err != nil {
			if err == goaviatrix.ErrNotFound {
				d.SetId("")
				return nil
			} else if !strings.Contains(err.Error(), "Failed to list devices: can not access dinfo") {
				return diag.Errorf("could not read Edge NEO device due to: %v", err)
			}
		} else {
			break
		}

		if i < 5 {
			time.Sleep(time.Duration(5) * time.Second)
		} else {
			d.SetId("")
			return diag.Errorf("could not read Edge NEO device after 5 retries: %s", err)
		}
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
		network1["dns_server_ips"] = network0.DnsServerIps
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

	if edgeNEODevice.DownloadConfigFile && edgeNEODevice.ConfigFileDownloadPath == "" {
		diag.Errorf("config_file_download_path is required when download_config_file is true")
	}
	if !edgeNEODevice.DownloadConfigFile && edgeNEODevice.ConfigFileDownloadPath != "" {
		diag.Errorf("config_file_download_path must be empty when download_config_file is false")
	}

	if d.HasChanges("account_name", "serial_number", "hardware_model") {
		return diag.Errorf("account_name, serial_number and hardware_model are not allowed to be updated")
	}

	d.Partial(true)

	if d.HasChanges("device_name", "network") {
		if err := client.UpdateEdgeNEODevice(ctx, edgeNEODevice); err != nil {
			return diag.Errorf("could not update device_name and network configurations during Edge NEO device update: %v", err)
		}

		if edgeNEODevice.DownloadConfigFile {
			if err := client.DownloadEdgeNEOConfigFile(ctx, edgeNEODevice); err != nil {
				return diag.Errorf("could not download Edge NEO static config file: %v", err)
			}
		}
	}

	if d.HasChange("download_config_file") {
		if edgeNEODevice.DownloadConfigFile {
			if err := client.DownloadEdgeNEOConfigFile(ctx, edgeNEODevice); err != nil {
				return diag.Errorf("could not download Edge NEO static config file: %v", err)
			}
		} else {
			oldConfigFileDownloadPath, _ := d.GetChange("config_file_download_path")
			fileName := oldConfigFileDownloadPath.(string) + edgeNEODevice.SerialNumber + "-bootstrap-config.img"
			err := os.Remove(fileName)
			if err != nil {
				log.Printf("[WARN] could not remove the config file: %v", err)
			}
		}
	}

	if d.HasChange("config_file_download_path") {
		if err := client.DownloadEdgeNEOConfigFile(ctx, edgeNEODevice); err != nil {
			return diag.Errorf("could not download Edge NEO static config file: %v", err)
		}
	}

	d.Partial(false)
	return resourceAviatrixEdgeNEODeviceOnboardingRead(ctx, d, meta)
}

func resourceAviatrixEdgeNEODeviceOnboardingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	edgeNEODevice := marshalEdgeNEODeviceOnboardingInput(d)

	err := client.DeleteEdgeNEODevice(ctx, edgeNEODevice.AccountName, edgeNEODevice.SerialNumber)
	if err != nil {
		return diag.Errorf("could not delete Edge NEO device: %v", err)
	}

	if edgeNEODevice.DownloadConfigFile {
		fileName := edgeNEODevice.ConfigFileDownloadPath + edgeNEODevice.SerialNumber + "-bootstrap-config.img"
		err = os.Remove(fileName)
		if err != nil {
			log.Printf("[WARN] could not remove the config file: %v", err)
		}
	}

	return nil
}
