package aviatrix

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixSpokeTransitAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSpokeTransitAttachmentCreate,
		Read:   resourceAviatrixSpokeTransitAttachmentRead,
		Update: resourceAviatrixSpokeTransitAttachmentUpdate,
		Delete: resourceAviatrixSpokeTransitAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"spoke_gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the spoke gateway to attach to transit network.",
			},
			"transit_gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the transit gateway to attach the spoke gateway to.",
			},
			"route_tables": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				ForceNew:    true,
				Description: "Learned routes will be propagated to these route tables.",
			},
			"spoke_prepend_as_path": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "AS Path Prepend customized by specifying AS PATH for a BGP connection. Applies on spoke gateway.",
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
			"spoke_bgp_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the spoke gateway is BGP enabled or not.",
			},
		},
	}
}

func resourceAviatrixSpokeTransitAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	attachment := marshalSpokeTransitAttachmentInput(d)

	spoke, err := client.GetGateway(&goaviatrix.Gateway{GwName: attachment.SpokeGwName})
	if err != nil {
		return fmt.Errorf("could not find spoke gateway: %s", err)
	}

	if !spoke.EnableBgp && (len(attachment.SpokePrependAsPath) != 0 || len(attachment.TransitPrependAsPath) != 0) {
		return fmt.Errorf("'spoke_prepend_as_path' and 'transit_prepend_as_path' are only valid for BGP enabled spoke gateway")
	}

	d.SetId(attachment.SpokeGwName + "~" + attachment.TransitGwName)
	flag := false
	defer resourceAviatrixSpokeTransitAttachmentReadIfRequired(d, meta, &flag)

	try, maxTries, backoff := 0, 8, 1000*time.Millisecond
	for {
		try++
		err := client.CreateSpokeTransitAttachment(attachment)
		if err != nil {
			if strings.Contains(err.Error(), "is not up") {
				if try == maxTries {
					return fmt.Errorf("could not attach spoke: %s to transit %s: %v", attachment.SpokeGwName, attachment.TransitGwName, err)
				}
				time.Sleep(backoff)
				// Double the backoff time after each failed try
				backoff *= 2
				continue
			}
			return fmt.Errorf("failed to attach spoke: %s to transit %s: %v", attachment.SpokeGwName, attachment.TransitGwName, err)
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
			return fmt.Errorf("could not set spoke_prepend_as_path: %v", err)
		}
	}

	if len(attachment.TransitPrependAsPath) != 0 {
		transGwPeering := &goaviatrix.TransitGatewayPeering{
			TransitGatewayName1: attachment.TransitGwName,
			TransitGatewayName2: attachment.SpokeGwName,
		}

		err := client.EditTransitConnectionASPathPrepend(transGwPeering, attachment.TransitPrependAsPath)
		if err != nil {
			return fmt.Errorf("could not set transit_prepend_as_path: %v", err)
		}
	}

	return resourceAviatrixSpokeTransitAttachmentReadIfRequired(d, meta, &flag)
}

func resourceAviatrixSpokeTransitAttachmentReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixSpokeTransitAttachmentRead(d, meta)
	}
	return nil
}

func resourceAviatrixSpokeTransitAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	spokeGwName := d.Get("spoke_gw_name").(string)
	transitGwName := d.Get("transit_gw_name").(string)
	if spokeGwName == "" || transitGwName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no spoke_gw_name or transit_gw_name received. Import Id is %s", id)
		d.SetId(id)
		parts := strings.Split(id, "~")
		if len(parts) != 2 {
			return fmt.Errorf("import id is invalid, expecting spoke_gw_name~transit_gw_name, but received: %s", id)
		}
		spokeGwName = parts[0]
		transitGwName = parts[1]
	}

	spokeTransitAttachment := &goaviatrix.SpokeTransitAttachment{
		SpokeGwName:   spokeGwName,
		TransitGwName: transitGwName,
	}

	attachment, err := client.GetSpokeTransitAttachment(spokeTransitAttachment)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not find spoke_transit_attachment: %v", err)
	}

	d.Set("spoke_gw_name", attachment.SpokeGwName)
	d.Set("transit_gw_name", attachment.TransitGwName)
	d.Set("spoke_bgp_enabled", attachment.SpokeBgpEnabled)

	if attachment.RouteTables != "" {
		var routeTables []string
		for _, routeTable := range strings.Split(attachment.RouteTables, ",") {
			routeTables = append(routeTables, strings.Split(routeTable, "~~")[0])
		}

		err = d.Set("route_tables", routeTables)
		if err != nil {
			return fmt.Errorf("could not set route_tables: %v", err)
		}
	}

	if attachment.SpokeBgpEnabled {
		transitGatewayPeering := &goaviatrix.TransitGatewayPeering{
			TransitGatewayName1: spokeGwName,
			TransitGatewayName2: transitGwName,
		}

		transitGatewayPeering, err = client.GetTransitGatewayPeeringDetails(transitGatewayPeering)
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}

		if transitGatewayPeering.PrependAsPath1 != "" {
			var prependAsPath []string
			for _, str := range strings.Split(transitGatewayPeering.PrependAsPath1, " ") {
				prependAsPath = append(prependAsPath, strings.TrimSpace(str))
			}

			err = d.Set("spoke_prepend_as_path", prependAsPath)
			if err != nil {
				return fmt.Errorf("could not set spoke_prepend_as_path: %v", err)
			}
		}

		if transitGatewayPeering.PrependAsPath2 != "" {
			var prependAsPath []string
			for _, str := range strings.Split(transitGatewayPeering.PrependAsPath2, " ") {
				prependAsPath = append(prependAsPath, strings.TrimSpace(str))
			}

			err = d.Set("transit_prepend_as_path", prependAsPath)
			if err != nil {
				return fmt.Errorf("could not set transit_prepend_as_path: %v", err)
			}
		}
	}

	d.SetId(attachment.SpokeGwName + "~" + attachment.TransitGwName)
	return nil
}

func resourceAviatrixSpokeTransitAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	if !d.Get("spoke_bgp_enabled").(bool) && d.HasChanges("spoke_prepend_as_path", "transit_prepend_as_path") {
		return fmt.Errorf("'spoke_prepend_as_path' and 'transit_prepend_as_path' are only valid for BGP enabled spoke gateway")
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
			return fmt.Errorf("could not update spoke_prepend_as_path: %v", err)
		}
	}

	if d.HasChange("transit_prepend_as_path") {
		transitGwPeering := &goaviatrix.TransitGatewayPeering{
			TransitGatewayName1: transitGwName,
			TransitGatewayName2: spokeGwName,
		}

		err := client.EditTransitConnectionASPathPrepend(transitGwPeering, getStringList(d, "transit_prepend_as_path"))
		if err != nil {
			return fmt.Errorf("could not update transit_prepend_as_path: %v", err)
		}
	}

	d.Partial(false)
	d.SetId(spokeGwName + "~" + transitGwName)
	return resourceAviatrixSpokeTransitAttachmentRead(d, meta)
}

func resourceAviatrixSpokeTransitAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	spokeTransitAttachment := &goaviatrix.SpokeTransitAttachment{
		SpokeGwName:   d.Get("spoke_gw_name").(string),
		TransitGwName: d.Get("transit_gw_name").(string),
	}
	if err := client.DeleteSpokeTransitAttachment(spokeTransitAttachment); err != nil {
		return fmt.Errorf("could not detach spoke: %s from transit %s: %v", spokeTransitAttachment.SpokeGwName,
			spokeTransitAttachment.TransitGwName, err)
	}

	return nil
}

func marshalSpokeTransitAttachmentInput(d *schema.ResourceData) *goaviatrix.SpokeTransitAttachment {
	spokeTransitAttachment := &goaviatrix.SpokeTransitAttachment{
		SpokeGwName:          d.Get("spoke_gw_name").(string),
		TransitGwName:        d.Get("transit_gw_name").(string),
		RouteTables:          strings.Join(getStringSet(d, "route_tables"), ","),
		SpokePrependAsPath:   getStringList(d, "spoke_prepend_as_path"),
		TransitPrependAsPath: getStringList(d, "transit_prepend_as_path"),
	}

	return spokeTransitAttachment
}
