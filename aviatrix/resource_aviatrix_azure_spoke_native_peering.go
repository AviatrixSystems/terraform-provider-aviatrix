package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixAzureSpokeNativePeering() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixAzureSpokeNativePeeringCreate,
		Read:   resourceAviatrixAzureSpokeNativePeeringRead,
		Update: resourceAviatrixAzureSpokeNativePeeringUpdate,
		Delete: resourceAviatrixAzureSpokeNativePeeringDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
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
			"private_route_table_config": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Private route table configuration.",
			},
		},
	}
}

func resourceAviatrixAzureSpokeNativePeeringCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	azureSpokeNativePeering := &goaviatrix.AzureSpokeNativePeering{
		TransitGatewayName: getString(d, "transit_gateway_name"),
		SpokeAccountName:   getString(d, "spoke_account_name"),
		SpokeRegion:        getString(d, "spoke_region"),
		SpokeVpcID:         getString(d, "spoke_vpc_id"),
	}

	log.Printf("[INFO] Creating Aviatrix Azure spoke native peering: %#v", azureSpokeNativePeering)

	d.SetId(azureSpokeNativePeering.TransitGatewayName + "~" + azureSpokeNativePeering.SpokeAccountName + "~" + azureSpokeNativePeering.SpokeVpcID)
	flag := false
	defer func() { _ = resourceAviatrixAzureSpokeNativePeeringReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	err := client.CreateAzureSpokeNativePeering(azureSpokeNativePeering)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Azure spoke native peering: %w", err)
	}

	// Configure private route table config if provided
	routeTables := getStringSet(d, "private_route_table_config")
	if len(routeTables) > 0 {
		// Gateway name for azure native spoke is: account_name:vpc_id (with dots replaced by dashes)
		gatewayName := azureSpokeNativePeering.SpokeAccountName + ":" + strings.ReplaceAll(azureSpokeNativePeering.SpokeVpcID, ".", "-")
		gw := &goaviatrix.Gateway{
			GwName: gatewayName,
		}
		err := client.EditPrivateRouteTableConfig(gw, routeTables)
		if err != nil {
			return fmt.Errorf("could not edit private route table config: %w", err)
		}
	}
	return resourceAviatrixAzureSpokeNativePeeringReadIfRequired(d, meta, &flag)
}

func resourceAviatrixAzureSpokeNativePeeringReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixAzureSpokeNativePeeringRead(d, meta)
	}
	return nil
}

func resourceAviatrixAzureSpokeNativePeeringRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	transitGatewayName := getString(d, "transit_gateway_name")
	spokeAccountName := getString(d, "spoke_account_name")
	spokeVpcID := getString(d, "spoke_vpc_id")

	if transitGatewayName == "" || spokeAccountName == "" || spokeVpcID == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no transit gateway name, or spoke account name, or spoke vpc id received. Import Id is %s", id)
		mustSet(d, "transit_gateway_name", strings.Split(id, "~")[0])
		mustSet(d, "spoke_account_name", strings.Split(id, "~")[1])
		mustSet(d, "spoke_vpc_id", strings.Split(id, "~")[2])
		d.SetId(id)
	}

	azureSpokeNativePeering := &goaviatrix.AzureSpokeNativePeering{
		TransitGatewayName: getString(d, "transit_gateway_name"),
		SpokeAccountName:   getString(d, "spoke_account_name"),
		SpokeVpcID:         getString(d, "spoke_vpc_id"),
	}

	azureSpokeNativePeering, err := client.GetAzureSpokeNativePeering(azureSpokeNativePeering)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix azure spoke native peering: %w", err)
	}
	mustSet(d, "transit_gateway_name", azureSpokeNativePeering.TransitGatewayName)
	mustSet(d, "spoke_account_name", azureSpokeNativePeering.SpokeAccountName)
	mustSet(d, "spoke_region", azureSpokeNativePeering.SpokeRegion)
	mustSet(d, "spoke_vpc_id", azureSpokeNativePeering.SpokeVpcID)
	mustSet(d, "private_route_table_config", azureSpokeNativePeering.PrivateRouteTableConfig)

	d.SetId(azureSpokeNativePeering.TransitGatewayName + "~" + azureSpokeNativePeering.SpokeAccountName + "~" + azureSpokeNativePeering.SpokeVpcID)
	return nil
}

func resourceAviatrixAzureSpokeNativePeeringUpdate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	if d.HasChange("private_route_table_config") {
		log.Printf("[INFO] resourceAviatrixAzureSpokeNativePeeringUpdate has changed private_route_table_config")

		spokeAccountName := getString(d, "spoke_account_name")
		spokeVpcID := getString(d, "spoke_vpc_id")

		routeTables := getStringSet(d, "private_route_table_config")
		// Gateway name for azure native spoke is: account_name:vpc_id (with dots replaced by dashes)
		gatewayName := spokeAccountName + ":" + strings.ReplaceAll(spokeVpcID, ".", "-")
		gateway := &goaviatrix.Gateway{
			GwName: gatewayName,
		}
		err := client.EditPrivateRouteTableConfig(gateway, routeTables)
		if err != nil {
			return fmt.Errorf("could not edit private route table config: %w", err)
		}
	}

	return resourceAviatrixAzureSpokeNativePeeringRead(d, meta)
}

func resourceAviatrixAzureSpokeNativePeeringDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	azureSpokeNativePeering := &goaviatrix.AzureSpokeNativePeering{
		TransitGatewayName: getString(d, "transit_gateway_name"),
		SpokeAccountName:   getString(d, "spoke_account_name"),
		SpokeVpcID:         getString(d, "spoke_vpc_id"),
	}

	log.Printf("[INFO] Deleting Aviatrix Azure spoke native peering: %#v", azureSpokeNativePeering)

	err := client.DeleteAzureSpokeNativePeering(azureSpokeNativePeering)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Azure spoke native peering: %w", err)
	}

	return nil
}
