package aviatrix

import (
	"context"
	"log"
	"strings"

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
			"enable_learned_cidrs_approval": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Switch to enable/disable encrypted transit approval for cloudn transit gateway attachment. Valid values: true, false.",
			},
			"approved_cidrs": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "Set of approved cidrs. Requires 'enable_learned_cidrs_approval' to be true. Type: Set(String).",
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
		EnableLearnedCidrsApproval:       d.Get("enable_learned_cidrs_approval").(bool),
	}
}

func resourceAviatrixCloudnTransitGatewayAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	attachment := marshalCloudnTransitGatewayAttachmentInput(d)

	enableLearnedCIDRApproval := d.Get("enable_learned_cidrs_approval").(bool)
	approvedCidrs := getStringSet(d, "approved_cidrs")
	if !enableLearnedCIDRApproval && len(approvedCidrs) > 0 {
		return diag.Errorf("error creating cloudn transit gateway attachment: 'approved_cidrs' must be empty if 'enable_learned_cidrs_approval' is false")
	}

	flag := false
	defer resourceAviatrixCloudnTransitGatewayAttachmentReadIfRequired(ctx, d, meta, &flag)

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

	if enableLearnedCIDRApproval {
		err = client.EnableTransitConnectionLearnedCIDRApproval(attachment.TransitGatewayName, attachment.ConnectionName)
		if err != nil {
			return diag.Errorf("could not enable learned cidr approval: %v", err)
		}
		if len(approvedCidrs) > 0 {
			err = client.UpdateTransitConnectionPendingApprovedCidrs(attachment.TransitGatewayName, attachment.ConnectionName, approvedCidrs)
			if err != nil {
				return diag.Errorf("could not update cloudn transit gateway attachment approved cidrs after creation: %v", err)
			}
		}
	}

	if _, ok := d.GetOk("prepend_as_path"); ok {
		var prependASPath []string
		for _, v := range d.Get("prepend_as_path").([]interface{}) {
			prependASPath = append(prependASPath, v.(string))
		}

		gateway := &goaviatrix.TransitVpc{
			GwName: attachment.TransitGatewayName,
		}
		err := client.SetPrependASPath(gateway, prependASPath)
		if err != nil {
			return diag.Errorf("could not update cloudn transit gateway attachment prepend_as_path after creation: %v", err)
		}
	}

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
	d.Set("enable_jumbo_frame", attachment.EnableJumboFrame)
	d.Set("enable_dead_peer_detection", attachment.EnableDeadPeerDetection)
	d.Set("enable_learned_cidrs_approval", attachment.EnableLearnedCidrsApproval)
	err = d.Set("approved_cidrs", attachment.ApprovedCidrs)
	if err != nil {
		return diag.Errorf("failed to set approved_cidrs for cloudn_transit_gateway_attachment on read: %v", err)
	}

	if attachment.PrependAsPath != "" {
		var prependAsPath []string
		for _, str := range strings.Split(attachment.PrependAsPath, " ") {
			prependAsPath = append(prependAsPath, strings.TrimSpace(str))
		}

		err = d.Set("prepend_as_path", prependAsPath)
		if err != nil {
			return diag.Errorf("could not set value for prepend_as_path: %v", err)
		}
	}

	d.SetId(attachment.ConnectionName)
	return nil
}

func resourceAviatrixCloudnTransitGatewayAttachmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)
	d.Partial(true)

	attachment := marshalCloudnTransitGatewayAttachmentInput(d)

	approvedCidrs := getStringSet(d, "approved_cidrs")
	enableLearnedCIDRApproval := d.Get("enable_learned_cidrs_approval").(bool)
	// `approved_cidrs` is Optional/Computed so Terraform will not detect a change even if `approved_cidrs` is removed from Terraform configuration.
	// Only prevent from setting `enable_learned_cidrs_approval = false` when Terraform detects a change to `approved_cidrs`.
	if !enableLearnedCIDRApproval && len(approvedCidrs) > 0 && d.HasChange("approved_cidrs") {
		return diag.Errorf("updating cloudn transit gateway attachment: 'approved_cidrs' must be empty if 'enable_learned_cidrs_approval' is false")
	}

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

	if d.HasChange("enable_learned_cidrs_approval") {
		if attachment.EnableLearnedCidrsApproval {
			err := client.EnableTransitConnectionLearnedCIDRApproval(attachment.TransitGatewayName, attachment.ConnectionName)
			if err != nil {
				return diag.Errorf("could not enable learned CIDRs approval during cloudn transit gateway attachment update: %v", err)
			}
		} else {
			err := client.DisableTransitConnectionLearnedCIDRApproval(attachment.TransitGatewayName, attachment.ConnectionName)
			if err != nil {
				return diag.Errorf("could not disable learned CIDRs approval during cloudn transit gateway attachment update: %v", err)
			}
		}
	}

	if d.HasChange("approved_cidrs") {
		err := client.UpdateTransitConnectionPendingApprovedCidrs(attachment.TransitGatewayName, attachment.ConnectionName, approvedCidrs)
		if err != nil {
			return diag.Errorf("could not update cloudn transit gateway attachment approved cidrs: %v", err)
		}
	}

	if d.HasChange("prepend_as_path") {
		var prependASPath []string
		for _, v := range d.Get("prepend_as_path").([]interface{}) {
			prependASPath = append(prependASPath, v.(string))
		}

		gateway := &goaviatrix.TransitVpc{
			GwName: attachment.TransitGatewayName,
		}
		err := client.SetPrependASPath(gateway, prependASPath)
		if err != nil {
			return diag.Errorf("could not update cloudn transit gateway attachment prepend_as_path: %v", err)
		}
	}

	d.Partial(false)
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
