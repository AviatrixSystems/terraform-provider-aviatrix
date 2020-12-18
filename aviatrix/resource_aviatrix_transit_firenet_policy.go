package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixTransitFireNetPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixTransitFireNetPolicyCreate,
		Read:   resourceAviatrixTransitFireNetPolicyRead,
		Delete: resourceAviatrixTransitFireNetPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceAviatrixTransitFireNetPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	transitFireNetPolicy := &goaviatrix.TransitFireNetPolicy{
		TransitFireNetGatewayName: d.Get("transit_firenet_gateway_name").(string),
		InspectedResourceName:     d.Get("inspected_resource_name").(string),
	}

	log.Printf("[INFO] Creating Aviatrix transit firenet policy: %#v", transitFireNetPolicy)

	err := client.CreateTransitFireNetPolicy(transitFireNetPolicy)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix transit firenet Policy: %s", err)
	}

	d.SetId(transitFireNetPolicy.TransitFireNetGatewayName + "~" + transitFireNetPolicy.InspectedResourceName)
	return resourceAviatrixTransitFireNetPolicyRead(d, meta)
}

func resourceAviatrixTransitFireNetPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	transitFireNetGatewayName := d.Get("transit_firenet_gateway_name").(string)
	inspectedResourceName := d.Get("inspected_resource_name").(string)

	if transitFireNetGatewayName == "" || inspectedResourceName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no transit firenet name or inspected resource name received. Import Id is %s", id)
		d.Set("transit_firenet_gateway_name", strings.Split(id, "~")[0])
		d.Set("inspected_resource_name", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	transitFireNetPolicy := &goaviatrix.TransitFireNetPolicy{
		TransitFireNetGatewayName: d.Get("transit_firenet_gateway_name").(string),
		InspectedResourceName:     d.Get("inspected_resource_name").(string),
	}

	err := client.GetTransitFireNetPolicy(transitFireNetPolicy)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix transit gateway policy: %s", err)
	}

	d.SetId(transitFireNetPolicy.TransitFireNetGatewayName + "~" + transitFireNetPolicy.InspectedResourceName)
	return nil
}

func resourceAviatrixTransitFireNetPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	transitFireNetPolicy := &goaviatrix.TransitFireNetPolicy{
		TransitFireNetGatewayName: d.Get("transit_firenet_gateway_name").(string),
		InspectedResourceName:     d.Get("inspected_resource_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix transit firenet policy %#v", transitFireNetPolicy)

	err := client.DeleteTransitFireNetPolicy(transitFireNetPolicy)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix transit firenet policy: %s", err)
	}

	return nil
}
