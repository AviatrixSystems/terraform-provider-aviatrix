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
			"enable_max_performance": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
				Description: "Indicates whether the maximum amount of HPE tunnels will be created. " +
					"Only valid when transit and spoke gateways are each launched in Insane Mode and in the same cloud type. " +
					"Available as of provider version R2.22.3+.",
			},
		},
	}
}

func resourceAviatrixSpokeTransitAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	attachment := marshalSpokeTransitAttachmentInput(d)

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

	transitGatewayPeering := &goaviatrix.TransitGatewayPeering{
		TransitGatewayName1: spokeGwName,
		TransitGatewayName2: transitGwName,
	}

	transitGatewayPeering, err = client.GetTransitGatewayPeeringDetails(transitGatewayPeering)
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}

	d.Set("enable_max_performance", !transitGatewayPeering.NoMaxPerformance)

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
		SpokeGwName:      d.Get("spoke_gw_name").(string),
		TransitGwName:    d.Get("transit_gw_name").(string),
		RouteTables:      strings.Join(getStringSet(d, "route_tables"), ","),
		NoMaxPerformance: !d.Get("enable_max_performance").(bool),
	}

	return spokeTransitAttachment
}
