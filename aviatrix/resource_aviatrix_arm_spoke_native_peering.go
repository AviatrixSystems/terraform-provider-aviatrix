package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixARMSpokeNativePeering() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixARMSpokeNativePeeringCreate,
		Read:   resourceAviatrixARMSpokeNativePeeringRead,
		Delete: resourceAviatrixARMSpokeNativePeeringDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"transit_gateway_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "",
			},
			"spoke_account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "",
			},
			"spoke_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "",
			},
			"spoke_vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "",
			},
		},
	}
}

func resourceAviatrixARMSpokeNativePeeringCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	armSpokeNativePeering := &goaviatrix.ArmSpokeNativePeering{
		TransitGatewayName: d.Get("transit_gateway_name").(string),
		SpokeAccountName:   d.Get("spoke_account_name").(string),
		SpokeRegion:        d.Get("spoke_region").(string),
		SpokeVpcID:         d.Get("spoke_vpc_id").(string),
	}

	log.Printf("[INFO] Creating Aviatrix arm spoke native peering: %#v", armSpokeNativePeering)

	err := client.CreateArmSpokeNativePeering(armSpokeNativePeering)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix arm spoke native peering: %s", err)
	}

	d.SetId(armSpokeNativePeering.TransitGatewayName + "~" + armSpokeNativePeering.SpokeAccountName + "~" + armSpokeNativePeering.SpokeVpcID)
	return resourceAviatrixARMSpokeNativePeeringRead(d, meta)
}

func resourceAviatrixARMSpokeNativePeeringRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	transitGatewayName := d.Get("transit_gateway_name").(string)
	spokeAccountName := d.Get("spoke_account_name").(string)
	spokeVpcID := d.Get("spoke_vpc_id").(string)

	if transitGatewayName == "" || spokeAccountName == "" || spokeVpcID == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no transit gateway name, or spoke account name, or spoke vpc id received. Import Id is %s", id)
		d.Set("transit_gateway_name", strings.Split(id, "~")[0])
		d.Set("spoke_account_name", strings.Split(id, "~")[1])
		d.Set("spoke_vpc_id", strings.Split(id, "~")[2])
		d.SetId(id)
	}

	armSpokeNativePeering := &goaviatrix.ArmSpokeNativePeering{
		TransitGatewayName: d.Get("transit_gateway_name").(string),
		SpokeAccountName:   d.Get("spoke_account_name").(string),
		SpokeVpcID:         d.Get("spoke_vpc_id").(string),
	}

	armSpokeNativePeering, err := client.GetArmSpokeNativePeering(armSpokeNativePeering)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix arm spoke native peering: %s", err)
	}

	d.Set("transit_gateway_name", armSpokeNativePeering.TransitGatewayName)
	d.Set("spoke_account_name", armSpokeNativePeering.SpokeAccountName)
	d.Set("spoke_region", armSpokeNativePeering.SpokeRegion)
	d.Set("spoke_vpc_id", armSpokeNativePeering.SpokeVpcID)

	d.SetId(armSpokeNativePeering.TransitGatewayName + "~" + armSpokeNativePeering.SpokeAccountName + "~" + armSpokeNativePeering.SpokeVpcID)
	return nil
}

func resourceAviatrixARMSpokeNativePeeringDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	armSpokeNativePeering := &goaviatrix.ArmSpokeNativePeering{
		TransitGatewayName: d.Get("transit_gateway_name").(string),
		SpokeAccountName:   d.Get("spoke_account_name").(string),
		SpokeVpcID:         d.Get("spoke_vpc_id").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix arm spoke native peering: %#v", armSpokeNativePeering)

	err := client.DeleteArmSpokeNativePeering(armSpokeNativePeering)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix arm spoke native peering: %s", err)
	}

	return nil
}
