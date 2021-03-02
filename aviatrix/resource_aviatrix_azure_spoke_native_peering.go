package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixAzureSpokeNativePeering() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixAzureSpokeNativePeeringCreate,
		Read:   resourceAviatrixAzureSpokeNativePeeringRead,
		Delete: resourceAviatrixAzureSpokeNativePeeringDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"transit_gateway_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of an azure transit gateway with transit firenet enabled.",
			},
			"spoke_account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "An Aviatrix account that corresponds to a subscription in Azure.",
			},
			"spoke_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Spoke VNet region.",
			},
			"spoke_vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Combination of the Spoke VNet name and resource group.",
			},
		},
	}
}

func resourceAviatrixAzureSpokeNativePeeringCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	azureSpokeNativePeering := &goaviatrix.AzureSpokeNativePeering{
		TransitGatewayName: d.Get("transit_gateway_name").(string),
		SpokeAccountName:   d.Get("spoke_account_name").(string),
		SpokeRegion:        d.Get("spoke_region").(string),
		SpokeVpcID:         d.Get("spoke_vpc_id").(string),
	}

	log.Printf("[INFO] Creating Aviatrix Azure spoke native peering: %#v", azureSpokeNativePeering)

	err := client.CreateAzureSpokeNativePeering(azureSpokeNativePeering)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Azure spoke native peering: %s", err)
	}

	d.SetId(azureSpokeNativePeering.TransitGatewayName + "~" + azureSpokeNativePeering.SpokeAccountName + "~" + azureSpokeNativePeering.SpokeVpcID)
	return resourceAviatrixAzureSpokeNativePeeringRead(d, meta)
}

func resourceAviatrixAzureSpokeNativePeeringRead(d *schema.ResourceData, meta interface{}) error {
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

	azureSpokeNativePeering := &goaviatrix.AzureSpokeNativePeering{
		TransitGatewayName: d.Get("transit_gateway_name").(string),
		SpokeAccountName:   d.Get("spoke_account_name").(string),
		SpokeVpcID:         d.Get("spoke_vpc_id").(string),
	}

	azureSpokeNativePeering, err := client.GetAzureSpokeNativePeering(azureSpokeNativePeering)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix azure spoke native peering: %s", err)
	}

	d.Set("transit_gateway_name", azureSpokeNativePeering.TransitGatewayName)
	d.Set("spoke_account_name", azureSpokeNativePeering.SpokeAccountName)
	d.Set("spoke_region", azureSpokeNativePeering.SpokeRegion)
	d.Set("spoke_vpc_id", azureSpokeNativePeering.SpokeVpcID)

	d.SetId(azureSpokeNativePeering.TransitGatewayName + "~" + azureSpokeNativePeering.SpokeAccountName + "~" + azureSpokeNativePeering.SpokeVpcID)
	return nil
}

func resourceAviatrixAzureSpokeNativePeeringDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	azureSpokeNativePeering := &goaviatrix.AzureSpokeNativePeering{
		TransitGatewayName: d.Get("transit_gateway_name").(string),
		SpokeAccountName:   d.Get("spoke_account_name").(string),
		SpokeVpcID:         d.Get("spoke_vpc_id").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix Azure spoke native peering: %#v", azureSpokeNativePeering)

	err := client.DeleteAzureSpokeNativePeering(azureSpokeNativePeering)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Azure spoke native peering: %s", err)
	}

	return nil
}
