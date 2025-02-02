package aviatrix

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixEdgeSpokeTransitAttachment() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixEdgeSpokeTransitAttachmentCreate,
		ReadWithoutTimeout:   resourceAviatrixEdgeSpokeTransitAttachmentRead,
		UpdateWithoutTimeout: resourceAviatrixEdgeSpokeTransitAttachmentUpdate,
		DeleteWithoutTimeout: resourceAviatrixEdgeSpokeTransitAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"spoke_gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Edge as a Spoke to attach to the transit network.",
			},
			"transit_gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the transit gateway to attach the Edge as a Spoke to.",
			},
			"enable_over_private_network": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     true,
				Description: "Enable over private network.",
			},
			"enable_jumbo_frame": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Enable jumbo frame.",
			},
			"enable_insane_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Enable jumbo frame.",
			},
			"insane_mode_tunnel_number": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Insane mode tunnel number. Valid range for HPE over private network: 0-49. Valid range for HPE over internet: 2-20.",
			},
			"spoke_prepend_as_path": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "AS Path Prepend customized by specifying AS PATH for a BGP connection. Applies on Edge as a Spoke.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: goaviatrix.ValidateASN,
				},
				MaxItems: 25,
			},
			"transit_prepend_as_path": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "AS Path Prepend customized by specifying AS PATH for a BGP connection. Applies on transit gateway.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: goaviatrix.ValidateASN,
				},
				MaxItems: 25,
			},
			"number_of_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Number of retries.",
			},
			"retry_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "Retry interval in seconds.",
			},
			"edge_wan_interfaces": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Description: "Set of Edge WAN interfaces.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: goaviatrix.DiffSuppressFuncEdgeSpokeTransitAttachmentEdgeWanInterfaces,
			},
			"default_edge_wan_interfaces": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Default Edge WAN interfaces.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"spoke_gateway_logical_ifnames": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Spoke gateway logical interface names for edge gateways, where the peering originates. Required for Megaport cloud type.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"transit_gateway_logical_ifnames": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Transit gateway logical interface names for edge gateways, where the peering terminates. Required for all edge gateways.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func marshalEdgeSpokeTransitAttachmentInput(d *schema.ResourceData) *goaviatrix.SpokeTransitAttachment {
	edgeSpokeTransitAttachment := &goaviatrix.SpokeTransitAttachment{
		SpokeGwName:              d.Get("spoke_gw_name").(string),
		TransitGwName:            d.Get("transit_gw_name").(string),
		EnableOverPrivateNetwork: d.Get("enable_over_private_network").(bool),
		EnableJumboFrame:         d.Get("enable_jumbo_frame").(bool),
		EnableInsaneMode:         d.Get("enable_insane_mode").(bool),
		InsaneModeTunnelNumber:   d.Get("insane_mode_tunnel_number").(int),
		SpokePrependAsPath:       getStringList(d, "spoke_prepend_as_path"),
		TransitPrependAsPath:     getStringList(d, "transit_prepend_as_path"),
		EdgeWanInterfaces:        strings.Join(getStringSet(d, "edge_wan_interfaces"), ","),
	}

	return edgeSpokeTransitAttachment
}

func resourceAviatrixEdgeSpokeTransitAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	attachment := marshalEdgeSpokeTransitAttachmentInput(d)

	// get edge spoke logical interface names
	edgeSpoke, err := client.GetEdgeSpoke(ctx, attachment.SpokeGwName)
	if err != nil {
		return diag.Errorf("failed to get spoke gateway details: %v", err)
	}
	if edgeSpoke.CloudType == goaviatrix.EDGEMEGAPORT {
		attachment.SpokeGatewayLogicalIfNames = getStringList(d, "spoke_gateway_logical_ifnames")
	}

	// get edge transit logical interface names
	transitGatewayDetails, err := getGatewayDetails(client, attachment.TransitGwName)
	if err != nil {
		return diag.Errorf("could not get transit gateway details for %s: %v", attachment.TransitGwName, err)
	}
	transitCloudType := transitGatewayDetails.CloudType
	// get the dst wan interfaces for Equinix & AEP EAT gateway
	if goaviatrix.IsCloudType(transitCloudType, goaviatrix.EDGEEQUINIX|goaviatrix.EDGENEO) {
		attachment.DstWanInterfaces, err = getDstWanInterfaces(attachment.TransitGatewayLogicalIfNames, transitGatewayDetails)
		if err != nil {
			return diag.Errorf("could not get dst wan interfaces for transit gateway: %v", err)
		}
	} else if goaviatrix.IsCloudType(transitCloudType, goaviatrix.EDGEMEGAPORT) {
		// get the logical if names for megaport
		attachment.TransitGatewayLogicalIfNames = getStringList(d, "transit_gateway_logical_ifnames")
	}

	if attachment.EnableInsaneMode {
		if attachment.EnableOverPrivateNetwork && (attachment.InsaneModeTunnelNumber < 0 || attachment.InsaneModeTunnelNumber > 49) {
			return diag.Errorf("valid range for HPE over private network: 0-49")
		}
		if !attachment.EnableOverPrivateNetwork && (attachment.InsaneModeTunnelNumber < 2 || attachment.InsaneModeTunnelNumber > 20) {
			return diag.Errorf("valid range for HPE over internet: 2-20")
		}
	}

	d.SetId(attachment.SpokeGwName + "~" + attachment.TransitGwName)
	flag := false
	defer resourceAviatrixEdgeSpokeTransitAttachmentReadIfRequired(ctx, d, meta, &flag)

	numberOfRetries := d.Get("number_of_retries").(int)
	retryInterval := d.Get("retry_interval").(int)

	for i := 0; ; i++ {
		err := client.CreateSpokeTransitAttachment(context.Background(), attachment)
		if err != nil {
			if !strings.Contains(err.Error(), "not ready") && !strings.Contains(err.Error(), "not up") &&
				!strings.Contains(err.Error(), "try again") {
				return diag.Errorf("could not attach Edge as a Spoke: %s to transit %s: %v", attachment.SpokeGwName, attachment.TransitGwName, err)
			}
		} else {
			break
		}
		if i < numberOfRetries {
			time.Sleep(time.Duration(retryInterval) * time.Second)
		} else {
			d.SetId("")
			return diag.Errorf("could not attach Edge as a Spoke: %s to transit %s: %v", attachment.SpokeGwName, attachment.TransitGwName, err)
		}
	}

	if len(attachment.SpokePrependAsPath) != 0 {
		transGwPeering := &goaviatrix.TransitGatewayPeering{
			TransitGatewayName1: attachment.SpokeGwName,
			TransitGatewayName2: attachment.TransitGwName,
		}

		err := client.EditTransitConnectionASPathPrepend(transGwPeering, attachment.SpokePrependAsPath)
		if err != nil {
			return diag.Errorf("could not set spoke_prepend_as_path: %v", err)
		}
	}

	if len(attachment.TransitPrependAsPath) != 0 {
		transGwPeering := &goaviatrix.TransitGatewayPeering{
			TransitGatewayName1: attachment.TransitGwName,
			TransitGatewayName2: attachment.SpokeGwName,
		}

		err := client.EditTransitConnectionASPathPrepend(transGwPeering, attachment.TransitPrependAsPath)
		if err != nil {
			return diag.Errorf("could not set transit_prepend_as_path: %v", err)
		}
	}

	return resourceAviatrixEdgeSpokeTransitAttachmentReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixEdgeSpokeTransitAttachmentReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixEdgeSpokeTransitAttachmentRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixEdgeSpokeTransitAttachmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	spokeGwName := d.Get("spoke_gw_name").(string)
	transitGwName := d.Get("transit_gw_name").(string)
	if spokeGwName == "" || transitGwName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no spoke_gw_name or transit_gw_name received. Import Id is %s", id)
		d.SetId(id)
		parts := strings.Split(id, "~")
		if len(parts) != 2 {
			return diag.Errorf("import id is invalid, expecting spoke_gw_name~transit_gw_name, but received: %s", id)
		}
		d.Set("spoke_gw_name", parts[0])
		d.Set("transit_gw_name", parts[1])
		spokeGwName = parts[0]
		transitGwName = parts[1]
	}

	spokeTransitAttachment := &goaviatrix.SpokeTransitAttachment{
		SpokeGwName:   spokeGwName,
		TransitGwName: transitGwName,
	}

	attachment, err := client.GetEdgeSpokeTransitAttachment(ctx, spokeTransitAttachment)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not find Edge as a Spoke transit attachment: %v", err)
	}

	d.Set("enable_over_private_network", attachment.EnableOverPrivateNetwork)
	d.Set("enable_jumbo_frame", attachment.EnableJumboFrame)
	d.Set("enable_insane_mode", attachment.EnableInsaneMode)
	if attachment.EnableInsaneMode {
		d.Set("insane_mode_tunnel_number", attachment.InsaneModeTunnelNumber)
	}

	if len(attachment.SpokePrependAsPath) != 0 {
		err = d.Set("spoke_prepend_as_path", attachment.SpokePrependAsPath)
		if err != nil {
			return diag.Errorf("could not set spoke_prepend_as_path: %v", err)
		}
	}

	if len(attachment.TransitPrependAsPath) != 0 {
		err = d.Set("transit_prepend_as_path", attachment.TransitPrependAsPath)
		if err != nil {
			return diag.Errorf("could not set transit_prepend_as_path: %v", err)
		}
	}

	edgeSpoke, err := client.GetEdgeSpoke(ctx, spokeGwName)
	if err != nil {
		return diag.Errorf("couldn't get wan interfaces for edge gateway %s: %s", spokeGwName, err)
	}

	if edgeSpoke.CloudType == goaviatrix.EDGEMEGAPORT {
		if len(attachment.SpokeGatewayLogicalIfNames) > 0 {
			_ = d.Set("spoke_gateway_logical_ifnames", attachment.SpokeGatewayLogicalIfNames)
		}
	} else {
		var defaultWanInterfaces []string
		for _, if0 := range edgeSpoke.InterfaceList {
			if if0.Type == "WAN" {
				defaultWanInterfaces = append(defaultWanInterfaces, if0.IfName)
			}
		}

		edgeWanInterfacesInput := getStringSet(d, "edge_wan_interfaces")

		if len(attachment.EdgeWanInterfacesResp) == 0 || (len(edgeWanInterfacesInput) == 0 && goaviatrix.Equivalent(attachment.EdgeWanInterfacesResp, defaultWanInterfaces)) ||
			(len(edgeWanInterfacesInput) != 0 && goaviatrix.Equivalent(edgeWanInterfacesInput, defaultWanInterfaces)) {
			_ = d.Set("edge_wan_interfaces", nil)
		} else {
			_ = d.Set("edge_wan_interfaces", attachment.EdgeWanInterfacesResp)
		}

		if len(defaultWanInterfaces) != 0 && len(attachment.EdgeWanInterfacesResp) != 0 {
			_ = d.Set("default_edge_wan_interfaces", defaultWanInterfaces)
		}
	}

	transitGateway, err := getGatewayDetails(client, transitGwName)
	if err != nil {
		return diag.Errorf("could not get transit gateway details for %s: %v", transitGwName, err)
	}

	if goaviatrix.IsCloudType(transitGateway.CloudType, goaviatrix.EdgeRelatedCloudTypes) {
		if len(attachment.TransitGatewayLogicalIfNames) > 0 {
			_ = d.Set("transit_gateway_logical_ifnames", attachment.TransitGatewayLogicalIfNames)
		}
	}

	d.SetId(spokeGwName + "~" + transitGwName)
	return nil
}

func resourceAviatrixEdgeSpokeTransitAttachmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	enableInsaneMode := d.Get("enable_insane_mode").(bool)
	enableOverPrivateNetwork := d.Get("enable_over_private_network").(bool)
	insaneModeTunnelNumber := d.Get("insane_mode_tunnel_number").(int)

	if enableInsaneMode {
		if enableOverPrivateNetwork && (insaneModeTunnelNumber < 0 || insaneModeTunnelNumber > 49) {
			return diag.Errorf("valid range for HPE over private network: 0-49")
		}
		if !enableOverPrivateNetwork && (insaneModeTunnelNumber < 2 || insaneModeTunnelNumber > 20) {
			return diag.Errorf("valid range for HPE over internet: 2-20")
		}
	}

	d.Partial(true)

	spokeGwName := d.Get("spoke_gw_name").(string)
	transitGwName := d.Get("transit_gw_name").(string)

	if d.HasChange("spoke_prepend_as_path") {
		transitGatewayPeering := &goaviatrix.TransitGatewayPeering{
			TransitGatewayName1: spokeGwName,
			TransitGatewayName2: transitGwName,
		}

		err := client.EditTransitConnectionASPathPrepend(transitGatewayPeering, getStringList(d, "spoke_prepend_as_path"))
		if err != nil {
			return diag.Errorf("could not update spoke_prepend_as_path: %v", err)
		}

	}

	if d.HasChange("transit_prepend_as_path") {
		transitGwPeering := &goaviatrix.TransitGatewayPeering{
			TransitGatewayName1: transitGwName,
			TransitGatewayName2: spokeGwName,
		}

		err := client.EditTransitConnectionASPathPrepend(transitGwPeering, getStringList(d, "transit_prepend_as_path"))
		if err != nil {
			return diag.Errorf("could not update transit_prepend_as_path: %v", err)
		}

	}

	if d.HasChange("insane_mode_tunnel_number") {
		transitGatewayPeering := &goaviatrix.TransitGatewayPeering{
			TransitGatewayName1: spokeGwName,
			TransitGatewayName2: transitGwName,
			TunnelCount:         insaneModeTunnelNumber,
		}

		err := client.UpdateTransitGatewayPeering(transitGatewayPeering)
		if err != nil {
			return diag.Errorf("could not update insane_mode_tunnel_number for edge spoke transit attachment: %v : %v", spokeGwName+"~"+transitGwName, err)
		}
	}

	d.Partial(false)
	d.SetId(spokeGwName + "~" + transitGwName)
	return resourceAviatrixEdgeSpokeTransitAttachmentRead(ctx, d, meta)
}

func resourceAviatrixEdgeSpokeTransitAttachmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	attachment := &goaviatrix.SpokeTransitAttachment{
		SpokeGwName:   d.Get("spoke_gw_name").(string),
		TransitGwName: d.Get("transit_gw_name").(string),
	}

	if err := client.DeleteSpokeTransitAttachment(attachment); err != nil {
		return diag.Errorf("could not detach Edge as a Spoke: %s from transit %s: %v", attachment.SpokeGwName, attachment.TransitGwName, err)
	}

	return nil
}

func getDstWanInterfaces(
	logicalIfNames []string,
	gatewayDetails *goaviatrix.Gateway,
) (string, error) {
	reversedInterfaceNames := ReverseIfnameTranslation(gatewayDetails.IfNamesTranslation)
	dstWanInterfaceStr, err := SetWanInterfaces(convertToInterfaceSlice(logicalIfNames), reversedInterfaceNames)
	if err != nil {
		return "", err
	}
	return dstWanInterfaceStr, nil
}
