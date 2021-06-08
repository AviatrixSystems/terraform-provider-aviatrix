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
		UpdateWithoutTimeout: resourceAviatrixCloudnTransitGatewayAttachmentUpdate,
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
			"enable_jumbo_frame": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable jumbo frame.",
			},
			"enable_dead_peer_detection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable dead peer detection.",
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
		EnableJumboFrame:                 d.Get("enable_jumbo_frame").(bool),
		EnableDeadPeerDetection:          d.Get("enable_dead_peer_detection").(bool),
	}
}

func resourceAviatrixCloudnTransitGatewayAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	attachment := marshalCloudnTransitGatewayAttachmentInput(d)

	err := client.CreateCloudnTransitGatewayAttachment(ctx, attachment)
	if err != nil {
		return diag.Errorf("could not create cloudn transit gateway attachment: %v", err)
	}
	d.SetId(attachment.ConnectionName)

	var vpcID string
	if attachment.EnableJumboFrame || !attachment.EnableDeadPeerDetection {
		vpcID, err = client.GetDeviceAttachmentVpcID(attachment.ConnectionName)
		if err != nil {
			return diag.Errorf("could not get cloudn transit gateway attachment VPC id after creating: %v", err)
		}
	}

	if attachment.EnableJumboFrame {
		err = client.EnableJumboFrameOnConnectionToCloudn(ctx, attachment.ConnectionName, vpcID)
		if err != nil {
			return diag.Errorf("could not enable jumbo frame after creating cloudn transit gateway attachment: %v", err)
		}
	}

	if !attachment.EnableDeadPeerDetection {
		s2c := goaviatrix.Site2Cloud{
			VpcID:      vpcID,
			TunnelName: attachment.ConnectionName,
		}
		err = client.DisableDeadPeerDetection(&s2c)
		if err != nil {
			return diag.Errorf("could not enable dead peer detection after creating cloudn transit gateway attachment: %v", err)
		}
	}

	return resourceAviatrixCloudnTransitGatewayAttachmentRead(ctx, d, meta)
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
	d.Set("enable_jumbo_frame", attachment.EnableJumboFrame)
	d.Set("enable_dead_peer_detection", attachment.EnableDeadPeerDetection)

	d.SetId(attachment.ConnectionName)
	return nil
}

func resourceAviatrixCloudnTransitGatewayAttachmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	attachment := marshalCloudnTransitGatewayAttachmentInput(d)

	var vpcID string
	if d.HasChanges("enable_jumbo_frame", "enable_dead_peer_detection") {
		var err error
		vpcID, err = client.GetDeviceAttachmentVpcID(attachment.ConnectionName)
		if err != nil {
			return diag.Errorf("could not get cloudn transit gateway attachment VPC id during update: %v", err)
		}
	}

	if d.HasChange("enable_jumbo_frame") {
		if attachment.EnableJumboFrame {
			err := client.EnableJumboFrameOnConnectionToCloudn(ctx, attachment.ConnectionName, vpcID)
			if err != nil {
				return diag.Errorf("could not enable jumbo frame during cloudn transit gateway attachment update: %v", err)
			}
		} else {
			err := client.DisableJumboFrameOnConnectionToCloudn(ctx, attachment.ConnectionName, vpcID)
			if err != nil {
				return diag.Errorf("could not disable jumbo frame during cloudn transit gateway attachment update: %v", err)
			}
		}
	}

	if d.HasChange("enable_dead_peer_detection") {
		s2c := goaviatrix.Site2Cloud{
			VpcID:      vpcID,
			TunnelName: attachment.ConnectionName,
		}
		if attachment.EnableDeadPeerDetection {
			err := client.EnableDeadPeerDetection(&s2c)
			if err != nil {
				return diag.Errorf("could not enable dead peer detection during cloudn transit gateway attachment update: %v", err)
			}
		} else {
			err := client.DisableDeadPeerDetection(&s2c)
			if err != nil {
				return diag.Errorf("could not disable dead peer detection during cloudn transit gateway attachment update: %v", err)
			}
		}
	}

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
