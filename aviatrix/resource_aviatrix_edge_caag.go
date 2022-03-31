package aviatrix

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixEdgeCaag() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixEdgeCaagCreate,
		ReadWithoutTimeout:   resourceAviatrixEdgeCaagRead,
		UpdateWithoutTimeout: resourceAviatrixEdgeCaagUpdate,
		DeleteWithoutTimeout: resourceAviatrixEdgeCaagDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Edge as a CaaG name.",
			},
			"management_interface_config": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Management interface configuration. Valid values: 'DHCP' and 'Static'.",
				ValidateFunc: validation.StringInSlice([]string{"DHCP", "Static"}, false),
			},
			"wan_interface_ip_prefix": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "WAN interface IP / prefix.",
			},
			"wan_default_gateway_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "WAN default gateway IP.",
				ValidateFunc: validation.IsIPAddress,
			},
			"lan_interface_ip_prefix": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "LAN interface IP / prefix.",
			},
			"enable_over_private_network": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Enable management over private network.",
			},
			"management_interface_ip_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Management interface IP / prefix.",
			},
			"management_default_gateway_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "Management default gateway IP.",
				ValidateFunc: validation.IsIPAddress,
			},
			"dns_server_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "DNS server IP.",
				ValidateFunc: validation.IsIPAddress,
			},
			"secondary_dns_server_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "Secondary DNS server IP.",
				ValidateFunc: validation.IsIPAddress,
			},
			"image_download_path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old != ""
				},
				Description: "The location where the Edge as a CaaG image will be stored.",
			},
			"local_as_number": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Local AS number.",
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
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of Edge as a CaaG.",
			},
		},
	}
}

func marshalEdgeCaagInput(d *schema.ResourceData) *goaviatrix.EdgeCaag {
	edgeCaag := &goaviatrix.EdgeCaag{
		Name:                        d.Get("name").(string),
		ManagementInterfaceConfig:   d.Get("management_interface_config").(string),
		EnableOverPrivateNetwork:    d.Get("enable_over_private_network").(bool),
		WanInterfaceIpPrefix:        d.Get("wan_interface_ip_prefix").(string),
		WanDefaultGatewayIp:         d.Get("wan_default_gateway_ip").(string),
		LanInterfaceIpPrefix:        d.Get("lan_interface_ip_prefix").(string),
		ManagementInterfaceIpPrefix: d.Get("management_interface_ip_prefix").(string),
		ManagementDefaultGatewayIp:  d.Get("management_default_gateway_ip").(string),
		DnsServerIp:                 d.Get("dns_server_ip").(string),
		SecondaryDnsServerIp:        d.Get("secondary_dns_server_ip").(string),
		ImageDownloadPath:           d.Get("image_download_path").(string),
	}

	return edgeCaag
}

func resourceAviatrixEdgeCaagCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	edgeCaag := marshalEdgeCaagInput(d)

	if edgeCaag.ManagementInterfaceConfig == "DHCP" && (edgeCaag.ManagementInterfaceIpPrefix != "" || edgeCaag.ManagementDefaultGatewayIp != "" ||
		edgeCaag.DnsServerIp != "" || edgeCaag.SecondaryDnsServerIp != "") {
		return diag.Errorf("'management_interface_ip', 'management_default_gateway_ip', 'dns_server_ip' and 'secondary_dns_server_ip' are only valid when 'management_interface_config' is Static")
	}

	if edgeCaag.ManagementInterfaceConfig == "Static" && (edgeCaag.ManagementInterfaceIpPrefix == "" || edgeCaag.ManagementDefaultGatewayIp == "" ||
		edgeCaag.DnsServerIp == "" || edgeCaag.SecondaryDnsServerIp == "") {
		return diag.Errorf("'management_interface_ip', 'management_default_gateway_ip', 'dns_server_ip' and 'secondary_dns_server_ip' are required when 'management_interface_config' is Static")
	}

	d.SetId(edgeCaag.Name)
	flag := false
	defer resourceAviatrixEdgeCaagReadIfRequired(ctx, d, meta, &flag)

	if err := client.CreateEdgeCaag(ctx, edgeCaag); err != nil {
		return diag.Errorf("could not create Edge as a CaaG: %v", err)
	}

	gateway := &goaviatrix.TransitVpc{
		GwName: edgeCaag.Name,
	}

	_, prependAsPathOk := d.GetOk("prepend_as_path")
	if _, ok := d.GetOk("local_as_number"); ok {
		localASNumber := d.Get("local_as_number").(string)
		err := client.SetLocalASNumber(gateway, localASNumber)
		if err != nil {
			return diag.Errorf("could not create Edge as a CaaG: could not set local_as_number: %v", err)
		}

		if prependAsPathOk {
			var prependASPath []string
			for _, v := range d.Get("prepend_as_path").([]interface{}) {
				prependASPath = append(prependASPath, v.(string))
			}

			err := client.SetPrependASPath(gateway, prependASPath)
			if err != nil {
				return diag.Errorf("could not create Edge as a CaaG: could not set prepend_as_path: %v", err)
			}
		}
	} else if prependAsPathOk {
		return diag.Errorf("could not create Edge as a CaaG: prepend_as_path must be empty when local_as_number has not been set")
	}

	return resourceAviatrixEdgeCaagReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixEdgeCaagReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixEdgeCaagRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixEdgeCaagRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Get("name").(string) == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no name received. Import Id is %s", id)
		d.Set("name", id)
		d.SetId(id)
	}

	edgeCaag, err := client.GetEdgeCaag(ctx, d.Get("name").(string))
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read Edge as a CaaG: %v", err)
	}

	d.Set("name", edgeCaag.Name)
	d.Set("enable_over_private_network", edgeCaag.EnableOverPrivateNetwork)
	d.Set("wan_interface_ip_prefix", edgeCaag.WanInterfaceIpPrefix)
	d.Set("wan_default_gateway_ip", edgeCaag.WanDefaultGatewayIp)
	d.Set("lan_interface_ip_prefix", edgeCaag.LanInterfaceIpPrefix)
	d.Set("management_default_gateway_ip", edgeCaag.ManagementDefaultGatewayIp)
	d.Set("dns_server_ip", edgeCaag.DnsServerIp)
	d.Set("secondary_dns_server_ip", edgeCaag.SecondaryDnsServerIp)
	d.Set("state", edgeCaag.State)

	if edgeCaag.Dhcp {
		d.Set("management_interface_config", "DHCP")
	} else {
		d.Set("management_interface_config", "Static")
		d.Set("management_interface_ip_prefix", edgeCaag.ManagementInterfaceIpPrefix)
	}

	gateway := &goaviatrix.TransitVpc{
		GwName: d.Get("name").(string),
	}
	transitGatewayAdvancedConfig, err := client.GetTransitGatewayAdvancedConfig(gateway)
	if err != nil {
		return diag.Errorf("could not read Edge as a CaaG transit gateway advanced config: %v", err)
	}

	d.Set("local_as_number", transitGatewayAdvancedConfig.LocalASNumber)
	d.Set("prepend_as_path", transitGatewayAdvancedConfig.PrependASPath)

	d.SetId(edgeCaag.Name)
	return nil
}

func resourceAviatrixEdgeCaagUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	d.Partial(true)
	gateway := &goaviatrix.TransitVpc{
		GwName: d.Get("name").(string),
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
				return diag.Errorf("could not delete prepend_as_path during Edge as a CaaG update: %v", err)
			}
		}

		if d.HasChange("local_as_number") {
			localASNumber := d.Get("local_as_number").(string)
			err := client.SetLocalASNumber(gateway, localASNumber)
			if err != nil {
				return diag.Errorf("could not update Edge as a CaaG: could not set local_as_number: %v", err)
			}
		}

		if d.HasChange("prepend_as_path") && len(prependASPath) > 0 {
			err := client.SetPrependASPath(gateway, prependASPath)
			if err != nil {
				return diag.Errorf("could not update Edge as a CaaG prepend_as_path: %v", err)
			}
		}
	}
	d.Partial(false)

	return resourceAviatrixEdgeCaagRead(ctx, d, meta)
}

func resourceAviatrixEdgeCaagDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	name := d.Get("name").(string)
	state := d.Get("state").(string)
	imageDownloadPath := d.Get("image_download_path").(string)

	err := client.DeleteEdgeCaag(ctx, name, state)
	if err != nil {
		return diag.Errorf("could not delete Edge as a CaaG: %v", err)
	}

	err = os.Remove(imageDownloadPath + "/" + name + ".iso")
	if err != nil {
		log.Printf("[WARN] could not remove the image file: %v", err)
	}

	return nil
}
