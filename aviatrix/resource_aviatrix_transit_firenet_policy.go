package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixTransitFireNetInspection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixTransitFireNetInspectionCreate,
		Read:   resourceAviatrixTransitFireNetInspectionRead,
		Delete: resourceAviatrixTransitFireNetInspectionDelete,
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
				Description: "Name of the resource to be added to transit firenet inspection.",
			},
		},
	}
}

func resourceAviatrixTransitFireNetInspectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	transitFireNetInspection := &goaviatrix.TransitFireNetInspection{
		TransitFireNetGatewayName: d.Get("transit_firenet_gateway_name").(string),
		InspectedResourceName:     d.Get("inspected_resource_name").(string),
	}

	log.Printf("[INFO] Creating Aviatrix transit firenet inspection: %#v", transitFireNetInspection)

	err := client.CreateTransitFireNetInspection(transitFireNetInspection)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix transit firenet inspection: %s", err)
	}

	d.SetId(transitFireNetInspection.TransitFireNetGatewayName + "~" + transitFireNetInspection.InspectedResourceName)
	return resourceAviatrixTransitFireNetInspectionRead(d, meta)
}

func resourceAviatrixTransitFireNetInspectionRead(d *schema.ResourceData, meta interface{}) error {
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

	transitFireNetInspection := &goaviatrix.TransitFireNetInspection{
		TransitFireNetGatewayName: d.Get("transit_firenet_gateway_name").(string),
		InspectedResourceName:     d.Get("inspected_resource_name").(string),
	}

	err := client.GetTransitFireNetInspection(transitFireNetInspection)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix transit gateway inspection: %s", err)
	}

	d.SetId(transitFireNetInspection.TransitFireNetGatewayName + "~" + transitFireNetInspection.InspectedResourceName)
	return nil
}

func resourceAviatrixTransitFireNetInspectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	transitFireNetInspection := &goaviatrix.TransitFireNetInspection{
		TransitFireNetGatewayName: d.Get("transit_firenet_gateway_name").(string),
		InspectedResourceName:     d.Get("inspected_resource_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix transit firenet inspection %#v", transitFireNetInspection)

	err := client.DeleteTransitFireNetInspection(transitFireNetInspection)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix transit firenet inspection: %s", err)
	}

	return nil
}
