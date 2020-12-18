package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixSpokeTransitAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSpokeTransitAttachmentCreate,
		Read:   resourceAviatrixSpokeTransitAttachmentRead,
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
		},
	}
}

func resourceAviatrixSpokeTransitAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	attachment := marshalSpokeTransitAttachmentInput(d)
	if err := client.CreateSpokeTransitAttachment(attachment); err != nil {
		return fmt.Errorf("could not attach spoke: %s to transit %s: %v", attachment.SpokeGwName, attachment.TransitGwName, err)
	}

	d.SetId(attachment.SpokeGwName + "~" + attachment.TransitGwName)
	return resourceAviatrixSpokeTransitAttachmentRead(d, meta)
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

	d.SetId(attachment.SpokeGwName + "~" + attachment.TransitGwName)
	return nil
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
		SpokeGwName:   d.Get("spoke_gw_name").(string),
		TransitGwName: d.Get("transit_gw_name").(string),
	}

	var routeTables []string
	for _, v := range d.Get("route_tables").(*schema.Set).List() {
		routeTables = append(routeTables, v.(string))
	}
	if len(routeTables) != 0 {
		spokeTransitAttachment.RouteTables = strings.Join(routeTables, ",")
	}

	return spokeTransitAttachment
}
