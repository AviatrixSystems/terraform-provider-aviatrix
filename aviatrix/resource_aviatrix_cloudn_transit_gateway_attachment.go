package aviatrix

import (
	"context"
	"log"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixCloudnTransitGatewayAttachment() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixCloudnTransitGatewayAttachmentCreate,
		ReadWithoutTimeout:   resourceAviatrixCloudnTransitGatewayAttachmentRead,
		DeleteWithoutTimeout: resourceAviatrixCloudnTransitGatewayAttachmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"device_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Device name.",
			},
			"transit_gateway_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Transit Gateway name.",
			},
			"connection_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Connection name.",
			},
			"transit_gateway_bgp_asn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Transit Gateway BGP AS Number.",
			},
			"cloudn_bgp_asn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "CloudN BGP AS Number.",
			},
			"cloudn_lan_interface_neighbor_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "CloudN LAN Interface Neighbor's IP.",
			},
			"cloudn_lan_interface_neighbor_bgp_asn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "CloudN LAN Interface Neighbor's BGP AS Number.",
			},
			"enable_over_private_network": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
				Description: "Enable over private network.",
			},
		},
	}
}

func marshalCloudnTransitGatewayAttachmentInput(d *schema.ResourceData) *goaviatrix.CloudnTransitGatewayAttachment {
	return &goaviatrix.CloudnTransitGatewayAttachment{
		DeviceName:                       d.Get("device_name").(string),
		TransitGatewayName:               d.Get("transit_gateway_name").(string),
		ConnectionName:                   d.Get("connection_name").(string),
		TransitGatewayBgpAsn:             d.Get("transit_gateway_bgp_asn").(string),
		CloudnBgpAsn:                     d.Get("cloudn_bgp_asn").(string),
		CloudnLanInterfaceNeighborIP:     d.Get("cloudn_lan_interface_neighbor_ip").(string),
		CloudnLanInterfaceNeighborBgpAsn: d.Get("cloudn_lan_interface_neighbor_bgp_asn").(string),
		EnableOverPrivateNetwork:         d.Get("enable_over_private_network").(bool),
	}
}

func resourceAviatrixCloudnTransitGatewayAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	attachment := marshalCloudnTransitGatewayAttachmentInput(d)

	flag := false
	defer resourceAviatrixCloudnTransitGatewayAttachmentReadIfRequired(ctx, d, meta, &flag)

	err := client.CreateCloudnTransitGatewayAttachment(ctx, attachment)
	if err != nil {
		return diag.Errorf("could not create cloudn transit gateway attachment: %v", err)
	}

	d.SetId(attachment.ConnectionName)
	return resourceAviatrixCloudnTransitGatewayAttachmentReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixCloudnTransitGatewayAttachmentReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixCloudnTransitGatewayAttachmentRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixCloudnTransitGatewayAttachmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	connName := d.Get("connection_name").(string)
	if connName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)
		connName = id
		d.Set("connection_name", connName)
		d.SetId(connName)
	}

	attachment := marshalCloudnTransitGatewayAttachmentInput(d)

	attachment, err := client.GetCloudnTransitGatewayAttachment(ctx, connName)
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("could not get cloudn transit gateway attachment: %v", err)
	}

	d.Set("device_name", attachment.DeviceName)
	d.Set("transit_gateway_name", attachment.TransitGatewayName)
	d.Set("connection_name", attachment.ConnectionName)
	d.Set("transit_gateway_bgp_asn", attachment.TransitGatewayBgpAsn)
	d.Set("cloudn_bgp_asn", attachment.CloudnBgpAsn)
	d.Set("cloudn_lan_interface_neighbor_ip", attachment.CloudnLanInterfaceNeighborIP)
	d.Set("cloudn_lan_interface_neighbor_bgp_asn", attachment.CloudnLanInterfaceNeighborBgpAsn)
	d.Set("enable_over_private_network", attachment.EnableOverPrivateNetwork)

	d.SetId(attachment.ConnectionName)
	return nil
}

func resourceAviatrixCloudnTransitGatewayAttachmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	attachment := marshalCloudnTransitGatewayAttachmentInput(d)

	err := client.DeleteDeviceAttachment(attachment.ConnectionName)
	if err != nil {
		return diag.Errorf("could not delete cloudn transit gateway attachment: %v", err)
	}

	return nil
}
