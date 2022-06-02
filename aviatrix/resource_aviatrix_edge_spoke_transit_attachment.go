package aviatrix

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
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
				Default:     false,
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
				ForceNew:    true,
				Default:     0,
				Description: "Insane mode tunnel number.",
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
	}

	return edgeSpokeTransitAttachment
}

func resourceAviatrixEdgeSpokeTransitAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	attachment := marshalEdgeSpokeTransitAttachmentInput(d)

	if attachment.EnableInsaneMode && attachment.InsaneModeTunnelNumber == 0 {
		diag.Errorf("'insane_mode_tunnel_number' must be set when insane mode is enabled")
	}

	if !attachment.EnableInsaneMode && attachment.InsaneModeTunnelNumber != 0 {
		diag.Errorf("'insane_mode_tunnel_number' is only valid when insane mode is enabled")
	}

	d.SetId(attachment.SpokeGwName + "~" + attachment.TransitGwName)
	flag := false
	defer resourceAviatrixEdgeSpokeTransitAttachmentReadIfRequired(ctx, d, meta, &flag)

	try, maxTries, backoff := 0, 8, 1000*time.Millisecond
	for {
		try++
		err := client.CreateSpokeTransitAttachment(attachment)
		if err != nil {
			if strings.Contains(err.Error(), "is not up") {
				if try == maxTries {
					return diag.Errorf("could not attach Edge as a Spoke: %s to transit %s: %v", attachment.SpokeGwName, attachment.TransitGwName, err)
				}
				time.Sleep(backoff)
				// Double the backoff time after each failed try
				backoff *= 2
				continue
			}
			return diag.Errorf("could not attach Edge as a Spoke: %s to transit %s: %v", attachment.SpokeGwName, attachment.TransitGwName, err)
		}
		break
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
	} else {
		d.Set("insane_mode_tunnel_number", 0)
	}

	if len(attachment.SpokePrependAsPath) != 0 {
		err = d.Set("spoke_prepend_as_path", attachment.SpokePrependAsPath)
		if err != nil {
			return diag.Errorf("could not set spoke_prepend_as_path: %v", err)
		}
	} else {
		d.Set("spoke_prepend_as_path", nil)
	}

	if len(attachment.TransitPrependAsPath) != 0 {
		err = d.Set("transit_prepend_as_path", attachment.TransitPrependAsPath)
		if err != nil {
			return diag.Errorf("could not set transit_prepend_as_path: %v", err)
		}
	} else {
		d.Set("transit_prepend_as_path", nil)
	}

	d.SetId(spokeGwName + "~" + transitGwName)
	return nil
}

func resourceAviatrixEdgeSpokeTransitAttachmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

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
