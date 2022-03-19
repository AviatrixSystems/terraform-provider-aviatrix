package aviatrix

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixCloudnEdgeGateway() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixCloudnEdgeGatewayCreate,
		ReadWithoutTimeout:   resourceAviatrixCloudnEdgeGatewayRead,
		UpdateWithoutTimeout: resourceAviatrixCloudnEdgeGatewayUpdate,
		DeleteWithoutTimeout: resourceAviatrixCloudnEdgeGatewayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Edge gateway name.",
			},
			"management_connection_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Management connection type. Valid values: 'DHCP' and 'Static'.",
			},
			"wan_interface_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "WAN interface IP.",
			},
			"wan_default_gateway": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "WAN default gateway.",
			},
			"lan_interface_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "LAN interface IP.",
			},
			"over_private_network": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
				Description: "Management over private network.",
			},
			"management_interface_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Management interface IP.",
			},
			"default_gateway_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Default gateway IP.",
			},
			"dns_server": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "DNS server.",
			},
			"secondary_dns_server": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Secondary DNS server.",
			},
			"image_download_path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old != ""
				},
				Description: "The location where the Edge gateway image will be stored.",
			},
			"local_as_number": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Changes the Aviatrix CloudN ASN number before you setup Aviatrix Transit Gateway connection configurations.",
				ValidateFunc: goaviatrix.ValidateASN,
			},
			"prepend_as_path": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "AS path prepend.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: goaviatrix.ValidateASN,
				},
				MaxItems: 25,
			},
		},
	}
}

func marshalCloudnEdgeGatewayInput(d *schema.ResourceData) *goaviatrix.CloudnEdgeGateway {
	cloudnEdgeGateway := &goaviatrix.CloudnEdgeGateway{
		GatewayName:              d.Get("gw_name").(string),
		ManagementConnectionType: d.Get("management_connection_type").(string),
		OverPrivateNetwork:       d.Get("over_private_network").(bool),
		WanInterfaceIp:           d.Get("wan_interface_ip").(string),
		WanDefaultGateway:        d.Get("wan_default_gateway").(string),
		LanInterfaceIp:           d.Get("lan_interface_ip").(string),
		ManagementInterfaceIp:    d.Get("management_interface_ip").(string),
		DefaultGatewayIp:         d.Get("default_gateway_ip").(string),
		DnsServer:                d.Get("dns_server").(string),
		SecondaryDnsServer:       d.Get("secondary_dns_server").(string),
		ImageDownloadPath:        d.Get("image_download_path").(string),
	}

	return cloudnEdgeGateway
}

func resourceAviatrixCloudnEdgeGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	cloudnEdgeGateway := marshalCloudnEdgeGatewayInput(d)

	if cloudnEdgeGateway.ManagementConnectionType == "DHCP" && (cloudnEdgeGateway.ManagementInterfaceIp != "" || cloudnEdgeGateway.DefaultGatewayIp != "" ||
		cloudnEdgeGateway.DnsServer != "" || cloudnEdgeGateway.SecondaryDnsServer != "") {
		return diag.Errorf("'management_interface_ip', 'default_gateway_ip', 'dns_server' and 'secondary_dns_server' are only valid when 'management_connection_type' is Static")
	}

	d.SetId(cloudnEdgeGateway.GatewayName)
	flag := false
	defer resourceAviatrixCloudnEdgeGatewayReadIfRequired(ctx, d, meta, &flag)

	if err := client.CreateCloundEdgeGateway(ctx, cloudnEdgeGateway); err != nil {
		return diag.Errorf("could not create cloudn edge gateway: %v", err)
	}

	gateway := &goaviatrix.TransitVpc{
		GwName: cloudnEdgeGateway.GatewayName,
	}

	_, prependAsPathOk := d.GetOk("prepend_as_path")
	if _, ok := d.GetOk("local_as_number"); ok {
		localASNumber := d.Get("local_as_number").(string)
		err := client.SetLocalASNumber(gateway, localASNumber)
		if err != nil {
			return diag.Errorf("failed to create cloudn edge gateway: could not set local_as_number: %v", err)
		}

		if prependAsPathOk {
			var prependASPath []string
			for _, v := range d.Get("prepend_as_path").([]interface{}) {
				prependASPath = append(prependASPath, v.(string))
			}

			err := client.SetPrependASPath(gateway, prependASPath)
			if err != nil {
				return diag.Errorf("failed to create cloudn edge gateway: could not set prepend_as_path: %v", err)
			}
		}
	} else if prependAsPathOk {
		return diag.Errorf("failed to create cloudn edge gateway: prepend_as_path must be empty when local_as_number has not been set")
	}

	return resourceAviatrixCloudnEdgeGatewayReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixCloudnEdgeGatewayReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixCloudnEdgeGatewayRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixCloudnEdgeGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Get("gw_name").(string) == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no gw_name received. Import Id is %s", id)
		d.Set("gw_name", id)
		d.SetId(id)
	}

	cloudnEdgeGateway, err := client.GetCloudnEdgeGateway(ctx, d.Get("gw_name").(string))
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read CloudN Edge Gateway: %v", err)
	}

	d.Set("gw_name", cloudnEdgeGateway.GatewayName)
	d.Set("over_private_network", cloudnEdgeGateway.OverPrivateNetwork)
	d.Set("wan_interface_ip", cloudnEdgeGateway.WanInterfaceIp)
	d.Set("wan_default_gateway", cloudnEdgeGateway.WanDefaultGateway)
	d.Set("lan_interface_ip", cloudnEdgeGateway.LanInterfaceIp)
	d.Set("management_interface_ip", cloudnEdgeGateway.ManagementInterfaceIp)
	d.Set("default_gateway_ip", cloudnEdgeGateway.DefaultGatewayIp)
	d.Set("dns_server", cloudnEdgeGateway.DnsServer)
	d.Set("secondary_dns_server", cloudnEdgeGateway.SecondaryDnsServer)

	if cloudnEdgeGateway.Dhcp {
		d.Set("management_connection_type", "DHCP")
	} else {
		d.Set("management_connection_type", "Static")
	}

	gateway := &goaviatrix.TransitVpc{
		GwName: d.Get("gw_name").(string),
	}
	transitGatewayAdvancedConfig, err := client.GetTransitGatewayAdvancedConfig(gateway)
	if err != nil {
		return diag.Errorf("failed to read cloudn edge gateway transit gateway advanced config: %v", err)
	}
	if transitGatewayAdvancedConfig.LocalASNumber != "" {
		d.Set("local_as_number", transitGatewayAdvancedConfig.LocalASNumber)
		d.Set("prepend_as_path", transitGatewayAdvancedConfig.PrependASPath)
	}

	d.SetId(cloudnEdgeGateway.GatewayName)
	return nil
}

func resourceAviatrixCloudnEdgeGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	d.Partial(true)
	gateway := &goaviatrix.TransitVpc{
		GwName: d.Get("gw_name").(string),
	}

	if d.HasChanges("local_as_number", "prepend_as_path") {
		var prependASPath []string
		for _, v := range d.Get("prepend_as_path").([]interface{}) {
			prependASPath = append(prependASPath, v.(string))
		}

		if (d.HasChange("local_as_number") && d.HasChange("prepend_as_path")) || len(prependASPath) == 0 {
			// prependASPath must be deleted from the controller before local_as_number can be changed
			// Handle the case where prependASPath is empty here so that the API is not called twice
			err := client.SetPrependASPath(gateway, nil)
			if err != nil {
				return diag.Errorf("failed to delete prepend_as_path during Aviatrix CloudN Registration update: %v", err)
			}
		}

		if d.HasChange("local_as_number") {
			localASNumber := d.Get("local_as_number").(string)
			err := client.SetLocalASNumber(gateway, localASNumber)
			if err != nil {
				return diag.Errorf("failed to update Aviatrix CloudN Registration: could not set local_as_number: %v", err)
			}
		}

		if d.HasChange("prepend_as_path") && len(prependASPath) > 0 {
			err := client.SetPrependASPath(gateway, prependASPath)
			if err != nil {
				return diag.Errorf("failed to update Aviatrix CloudN Registration prepend_as_path: %v", err)
			}
		}
	}
	d.Partial(false)

	return resourceAviatrixCloudnEdgeGatewayRead(ctx, d, meta)
}

func resourceAviatrixCloudnEdgeGatewayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	gwName := d.Get("gw_name").(string)
	imageDownloadPath := d.Get("image_download_path").(string)

	err := client.DeleteCloudnEdgeGateway(ctx, gwName)
	if err != nil {
		return diag.Errorf("could not delete cloudn edge gateway: %v", err)
	}

	err = os.Remove(imageDownloadPath + "/" + gwName + ".iso")
	if err != nil {
		log.Printf("[WARN] could not remove the image file: %v", err)
	}

	return nil
}
