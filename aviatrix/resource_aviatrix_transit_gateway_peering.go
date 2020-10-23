package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixTransitGatewayPeering() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixTransitGatewayPeeringCreate,
		Read:   resourceAviatrixTransitGatewayPeeringRead,
		Update: resourceAviatrixTransitGatewayPeeringUpdate,
		Delete: resourceAviatrixTransitGatewayPeeringDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"transit_gateway_name1": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The first transit gateway name to make a peer pair.",
			},
			"transit_gateway_name2": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The second transit gateway name to make a peer pair.",
			},
			"gateway1_excluded_network_cidrs": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of excluded network CIDRs for the first transit gateway.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"gateway1_excluded_tgw_connections": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of excluded TGW connections for the first transit gateway.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"gateway2_excluded_network_cidrs": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of excluded network CIDRs for the second transit gateway.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"gateway2_excluded_tgw_connections": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of excluded TGW connections for the second transit gateway.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enable_peering_over_private_network": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "(Optional) Enable peering over private network. Insane mode is required on both transit gateways. Available as of provider version R2.17.1",
			},
		},
	}
}

func resourceAviatrixTransitGatewayPeeringCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	var gw1Cidrs []string
	for _, cidr := range d.Get("gateway1_excluded_network_cidrs").([]interface{}) {
		gw1Cidrs = append(gw1Cidrs, cidr.(string))
	}
	var gw2Cidrs []string
	for _, cidr := range d.Get("gateway2_excluded_network_cidrs").([]interface{}) {
		gw2Cidrs = append(gw2Cidrs, cidr.(string))
	}

	var gw1Tgws []string
	for _, tgw := range d.Get("gateway1_excluded_tgw_connections").([]interface{}) {
		gw1Tgws = append(gw1Tgws, tgw.(string))
	}
	var gw2Tgws []string
	for _, tgw := range d.Get("gateway2_excluded_tgw_connections").([]interface{}) {
		gw2Tgws = append(gw2Tgws, tgw.(string))
	}

	transitGatewayPeering := &goaviatrix.TransitGatewayPeering{
		TransitGatewayName1:            d.Get("transit_gateway_name1").(string),
		TransitGatewayName2:            d.Get("transit_gateway_name2").(string),
		Gateway1ExcludedCIDRs:          strings.Join(gw1Cidrs, ","),
		Gateway2ExcludedCIDRs:          strings.Join(gw2Cidrs, ","),
		Gateway1ExcludedTGWConnections: strings.Join(gw1Tgws, ","),
		Gateway2ExcludedTGWConnections: strings.Join(gw2Tgws, ","),
		PrivateIPPeering:               d.Get("enable_peering_over_private_network").(bool),
	}

	log.Printf("[INFO] Creating Aviatrix Transit Gateway peering: %#v", transitGatewayPeering)

	err := client.CreateTransitGatewayPeering(transitGatewayPeering)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Transit Gateway peering: %s", err)
	}

	d.SetId(transitGatewayPeering.TransitGatewayName1 + "~" + transitGatewayPeering.TransitGatewayName2)
	return resourceAviatrixTransitGatewayPeeringRead(d, meta)
}

func resourceAviatrixTransitGatewayPeeringRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	transitGwName1 := d.Get("transit_gateway_name1").(string)
	transitGwName2 := d.Get("transit_gateway_name2").(string)

	if transitGwName1 == "" || transitGwName2 == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no transit gateway names received. Import Id is %s", id)
		d.Set("transit_gateway_name1", strings.Split(id, "~")[0])
		d.Set("transit_gateway_name2", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	transitGatewayPeering := &goaviatrix.TransitGatewayPeering{
		TransitGatewayName1: d.Get("transit_gateway_name1").(string),
		TransitGatewayName2: d.Get("transit_gateway_name2").(string),
	}

	err := client.GetTransitGatewayPeering(transitGatewayPeering)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix Transit Gateway peering: %s", err)
	}

	transitGatewayPeering, err = client.GetTransitGatewayPeeringDetails(transitGatewayPeering)
	if err != nil {
		return fmt.Errorf("could not get transit peering details: %v", err)
	}

	if err := d.Set("gateway1_excluded_network_cidrs", transitGatewayPeering.Gateway1ExcludedCIDRsSlice); err != nil {
		return err
	}
	if err := d.Set("gateway2_excluded_network_cidrs", transitGatewayPeering.Gateway2ExcludedCIDRsSlice); err != nil {
		return err
	}
	if err := d.Set("gateway1_excluded_tgw_connections", transitGatewayPeering.Gateway1ExcludedTGWConnectionsSlice); err != nil {
		return err
	}
	if err := d.Set("gateway2_excluded_tgw_connections", transitGatewayPeering.Gateway2ExcludedTGWConnectionsSlice); err != nil {
		return err
	}

	d.Set("enable_peering_over_private_network", transitGatewayPeering.PrivateIPPeering)

	d.SetId(transitGatewayPeering.TransitGatewayName1 + "~" + transitGatewayPeering.TransitGatewayName2)
	return nil
}

func resourceAviatrixTransitGatewayPeeringUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	var gw1Cidrs []string
	for _, cidr := range d.Get("gateway1_excluded_network_cidrs").([]interface{}) {
		gw1Cidrs = append(gw1Cidrs, cidr.(string))
	}
	var gw2Cidrs []string
	for _, cidr := range d.Get("gateway2_excluded_network_cidrs").([]interface{}) {
		gw2Cidrs = append(gw2Cidrs, cidr.(string))
	}

	var gw1Tgws []string
	for _, tgw := range d.Get("gateway1_excluded_tgw_connections").([]interface{}) {
		gw1Tgws = append(gw1Tgws, tgw.(string))
	}
	var gw2Tgws []string
	for _, tgw := range d.Get("gateway2_excluded_tgw_connections").([]interface{}) {
		gw2Tgws = append(gw2Tgws, tgw.(string))
	}

	transitGatewayPeering := &goaviatrix.TransitGatewayPeering{
		TransitGatewayName1:            d.Get("transit_gateway_name1").(string),
		TransitGatewayName2:            d.Get("transit_gateway_name2").(string),
		Gateway1ExcludedCIDRs:          strings.Join(gw1Cidrs, ","),
		Gateway2ExcludedCIDRs:          strings.Join(gw2Cidrs, ","),
		Gateway1ExcludedTGWConnections: strings.Join(gw1Tgws, ","),
		Gateway2ExcludedTGWConnections: strings.Join(gw2Tgws, ","),
	}

	log.Printf("[INFO] Updating Aviatrix Transit Gateway peering: %#v", transitGatewayPeering)

	err := client.UpdateTransitGatewayPeering(transitGatewayPeering)
	if err != nil {
		return fmt.Errorf("failed to update Aviatrix Transit Gateway peering: %s", err)
	}

	d.SetId(transitGatewayPeering.TransitGatewayName1 + "~" + transitGatewayPeering.TransitGatewayName2)
	return resourceAviatrixTransitGatewayPeeringRead(d, meta)
}

func resourceAviatrixTransitGatewayPeeringDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	transitGatewayPeering := &goaviatrix.TransitGatewayPeering{
		TransitGatewayName1: d.Get("transit_gateway_name1").(string),
		TransitGatewayName2: d.Get("transit_gateway_name2").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix Transit Gateway peering: %#v", transitGatewayPeering)

	err := client.DeleteTransitGatewayPeering(transitGatewayPeering)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Transit Gateway peering: %s", err)
	}

	return nil
}
