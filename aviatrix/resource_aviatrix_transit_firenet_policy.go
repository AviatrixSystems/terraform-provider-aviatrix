package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixTransitFireNetPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixTransitFireNetPolicyCreate,
		Read:   resourceAviatrixTransitFireNetPolicyRead,
		Delete: resourceAviatrixTransitFireNetPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},

		Schema: map[string]*schema.Schema{
			"transit_firenet_gateway_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the transit firenet gateway.",
			},
			"inspected_resource_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the resource to be added to transit firenet policy.",
			},
		},
	}
}

func resourceAviatrixTransitFireNetPolicyCreate(d *schema.ResourceData, meta any) error {
	client := mustClient(meta)

	transitFireNetPolicy := &goaviatrix.TransitFireNetPolicy{
		TransitFireNetGatewayName: getString(d, "transit_firenet_gateway_name"),
		InspectedResourceName:     getString(d, "inspected_resource_name"),
	}

	log.Printf("[INFO] Creating Aviatrix transit firenet policy: %#v", transitFireNetPolicy)

	d.SetId(transitFireNetPolicy.TransitFireNetGatewayName + "~" + transitFireNetPolicy.InspectedResourceName)
	flag := false
	defer func() { _ = resourceAviatrixTransitFireNetPolicyReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	err := client.CreateTransitFireNetPolicy(transitFireNetPolicy)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix transit firenet Policy: %w", err)
	}

	return resourceAviatrixTransitFireNetPolicyReadIfRequired(d, meta, &flag)
}

func resourceAviatrixTransitFireNetPolicyReadIfRequired(d *schema.ResourceData, meta any, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixTransitFireNetPolicyRead(d, meta)
	}
	return nil
}

func resourceAviatrixTransitFireNetPolicyRead(d *schema.ResourceData, meta any) error {
	client := mustClient(meta)

	transitFireNetGatewayName := getString(d, "transit_firenet_gateway_name")
	inspectedResourceName := getString(d, "inspected_resource_name")

	if transitFireNetGatewayName == "" || inspectedResourceName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no transit firenet name or inspected resource name received. Import Id is %s", id)
		// Split only on the first "~" since "inspected_resource_name" can itself contain "~~"
		// (e.g. "SPOKE_SUBNET_GROUP:<gw_name>~~<subnet_group_name>"), and a naive Split on every
		// "~" would truncate/mangle the inspected resource name and cause import to fail with
		// "Cannot import non-existent remote object".
		parts := strings.SplitN(id, "~", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid ID, expected ID transit_firenet_gateway_name~inspected_resource_name, instead got %s", id)
		}
		mustSet(d, "transit_firenet_gateway_name", parts[0])
		mustSet(d, "inspected_resource_name", parts[1])
		d.SetId(id)
	}

	transitFireNetPolicy := &goaviatrix.TransitFireNetPolicy{
		TransitFireNetGatewayName: getString(d, "transit_firenet_gateway_name"),
		InspectedResourceName:     getString(d, "inspected_resource_name"),
	}

	err := client.GetTransitFireNetPolicy(transitFireNetPolicy)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix transit gateway policy: %w", err)
	}

	d.SetId(transitFireNetPolicy.TransitFireNetGatewayName + "~" + transitFireNetPolicy.InspectedResourceName)
	return nil
}

func resourceAviatrixTransitFireNetPolicyDelete(d *schema.ResourceData, meta any) error {
	client := mustClient(meta)

	transitFireNetPolicy := &goaviatrix.TransitFireNetPolicy{
		TransitFireNetGatewayName: getString(d, "transit_firenet_gateway_name"),
		InspectedResourceName:     getString(d, "inspected_resource_name"),
	}

	log.Printf("[INFO] Deleting Aviatrix transit firenet policy %#v", transitFireNetPolicy)

	err := client.DeleteTransitFireNetPolicy(transitFireNetPolicy)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix transit firenet policy: %w", err)
	}

	return nil
}
