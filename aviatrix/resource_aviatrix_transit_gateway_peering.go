package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
			"prepend_as_path1": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "AS Path Prepend customized by specifying AS PATH for a BGP connection. Applies on transit_gateway_name1.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: goaviatrix.ValidateASN,
				},
				MaxItems: 25,
			},
			"prepend_as_path2": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "AS Path Prepend customized by specifying AS PATH for a BGP connection. Applies on transit_gateway_name2.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: goaviatrix.ValidateASN,
				},
				MaxItems: 25,
			},
			"enable_peering_over_private_network": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "(Optional) Enable peering over private network. Insane mode is required on both transit gateways. Available as of provider version R2.17.1",
			},
			"enable_single_tunnel_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Enable peering with Single-Tunnel mode.",
			},
			"enable_insane_mode_encryption_over_internet": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Enable Insane Mode Encryption over Internet. Type: Boolean. Default: false.",
			},
			"tunnel_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(2, 20),
				Description:  "Number of public tunnels. Only valid with 'enable_insane_mode_encryption_over_internet'. Type: Integer. Valid Range: 2-20.",
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
		InsaneModeOverInternet:         d.Get("enable_insane_mode_encryption_over_internet").(bool),
	}

	if transitGatewayPeering.PrivateIPPeering && transitGatewayPeering.InsaneModeOverInternet {
		return fmt.Errorf("enable_peering_over_private_network conflicts with enable_insane_mode_encryption_over_internet")
	}

	if d.Get("enable_single_tunnel_mode").(bool) {
		if transitGatewayPeering.PrivateIPPeering {
			transitGatewayPeering.SingleTunnel = true
		} else {
			return fmt.Errorf("enable_single_tunnel_mode is only valid when enable_peering_over_private_network is set to true")
		}
	}

	tunnelCount := d.Get("tunnel_count").(int)
	if tunnelCount != 0 {
		if transitGatewayPeering.InsaneModeOverInternet {
			transitGatewayPeering.TunnelCount = tunnelCount
		} else {
			return fmt.Errorf("tunnel_count is only valid when enable_insane_mode_encryption_over_internet is set to true")
		}
	} else {
		if transitGatewayPeering.InsaneModeOverInternet {
			return fmt.Errorf("enable_insane_mode_encryption_over_internet being set to true requires valid tunnel_count")
		}
	}

	log.Printf("[INFO] Creating Aviatrix Transit Gateway peering: %#v", transitGatewayPeering)

	d.SetId(transitGatewayPeering.TransitGatewayName1 + "~" + transitGatewayPeering.TransitGatewayName2)
	flag := false
	defer resourceAviatrixTransitGatewayPeeringReadIfRequired(d, meta, &flag)

	err := client.CreateTransitGatewayPeering(transitGatewayPeering)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Transit Gateway peering: %s", err)
	}

	if _, ok := d.GetOk("prepend_as_path1"); ok {
		var prependASPath []string
		for _, v := range d.Get("prepend_as_path1").([]interface{}) {
			prependASPath = append(prependASPath, v.(string))
		}
		transGwPeering := &goaviatrix.TransitGatewayPeering{
			TransitGatewayName1: d.Get("transit_gateway_name1").(string),
			TransitGatewayName2: d.Get("transit_gateway_name2").(string),
		}

		err = client.EditTransitConnectionASPathPrepend(transGwPeering, prependASPath)
		if err != nil {
			return fmt.Errorf("could not set prepend_as_path1: %v", err)
		}
	}

	if _, ok := d.GetOk("prepend_as_path2"); ok {
		var prependASPath []string
		for _, v := range d.Get("prepend_as_path2").([]interface{}) {
			prependASPath = append(prependASPath, v.(string))
		}
		transGwPeering := &goaviatrix.TransitGatewayPeering{
			TransitGatewayName1: d.Get("transit_gateway_name2").(string),
			TransitGatewayName2: d.Get("transit_gateway_name1").(string),
		}

		err = client.EditTransitConnectionASPathPrepend(transGwPeering, prependASPath)
		if err != nil {
			return fmt.Errorf("could not set prepend_as_path2: %v", err)
		}
	}

	return resourceAviatrixTransitGatewayPeeringReadIfRequired(d, meta, &flag)
}

func resourceAviatrixTransitGatewayPeeringReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixTransitGatewayPeeringRead(d, meta)
	}
	return nil
}

func resourceAviatrixTransitGatewayPeeringRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	transitGwName1 := d.Get("transit_gateway_name1").(string)
	transitGwName2 := d.Get("transit_gateway_name2").(string)

	if transitGwName1 == "" || transitGwName2 == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no transit gateway names received. Import Id is %s", id)
		parts := strings.Split(id, "~")
		if len(parts) != 2 {
			return fmt.Errorf("invalid import id expected transit_gateway_name1~transit_gateway_name2")
		}
		d.Set("transit_gateway_name1", parts[0])
		d.Set("transit_gateway_name2", parts[1])
		d.SetId(id)
	}

	transitGatewayPeering := &goaviatrix.TransitGatewayPeering{
		TransitGatewayName1: d.Get("transit_gateway_name1").(string),
		TransitGatewayName2: d.Get("transit_gateway_name2").(string),
	}

	transitGatewayPeering, err := client.GetTransitGatewayPeeringDetails(transitGatewayPeering)
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not get transit peering details: %v", err)
	}

	gw1CidrsFromConfig := getStringList(d, "gateway1_excluded_network_cidrs")
	err = setConfigValueIfEquivalent(d, "gateway1_excluded_network_cidrs", gw1CidrsFromConfig, transitGatewayPeering.Gateway1ExcludedCIDRsSlice)
	if err != nil {
		return fmt.Errorf("could not write gateway1_excluded_network_cidrs to state: %v", err)
	}
	gw2CidrsFromConfig := getStringList(d, "gateway2_excluded_network_cidrs")
	err = setConfigValueIfEquivalent(d, "gateway2_excluded_network_cidrs", gw2CidrsFromConfig, transitGatewayPeering.Gateway2ExcludedCIDRsSlice)
	if err != nil {
		return fmt.Errorf("could not write gateway2_excluded_network_cidrs to state: %v", err)
	}
	gw1TgwsFromConfig := getStringList(d, "gateway1_excluded_tgw_connections")
	err = setConfigValueIfEquivalent(d, "gateway1_excluded_tgw_connections", gw1TgwsFromConfig, transitGatewayPeering.Gateway1ExcludedTGWConnectionsSlice)
	if err != nil {
		return fmt.Errorf("could not write gateway1_excluded_tgw_connections to state: %v", err)
	}
	gw2TgwsFromConfig := getStringList(d, "gateway2_excluded_tgw_connections")
	err = setConfigValueIfEquivalent(d, "gateway2_excluded_tgw_connections", gw2TgwsFromConfig, transitGatewayPeering.Gateway2ExcludedTGWConnectionsSlice)
	if err != nil {
		return fmt.Errorf("could not write gateway2_excluded_tgw_connections to state: %v", err)
	}

	if transitGatewayPeering.PrependAsPath1 != "" {
		var prependAsPath []string
		for _, str := range strings.Split(transitGatewayPeering.PrependAsPath1, " ") {
			prependAsPath = append(prependAsPath, strings.TrimSpace(str))
		}

		err = d.Set("prepend_as_path1", prependAsPath)
		if err != nil {
			return fmt.Errorf("could not set prepend_as_path1: %v", err)
		}
	}
	if transitGatewayPeering.PrependAsPath2 != "" {
		var prependAsPath []string
		for _, str := range strings.Split(transitGatewayPeering.PrependAsPath2, " ") {
			prependAsPath = append(prependAsPath, strings.TrimSpace(str))
		}

		err = d.Set("prepend_as_path2", prependAsPath)
		if err != nil {
			return fmt.Errorf("could not set prepend_as_path2: %v", err)
		}
	}

	d.Set("enable_peering_over_private_network", transitGatewayPeering.PrivateIPPeering)
	d.Set("enable_single_tunnel_mode", transitGatewayPeering.SingleTunnel)
	d.Set("enable_insane_mode_encryption_over_internet", transitGatewayPeering.InsaneModeOverInternet)
	if transitGatewayPeering.InsaneModeOverInternet {
		d.Set("tunnel_count", transitGatewayPeering.TunnelCount)
	}

	d.SetId(transitGatewayPeering.TransitGatewayName1 + "~" + transitGatewayPeering.TransitGatewayName2)
	return nil
}

func resourceAviatrixTransitGatewayPeeringUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	d.Partial(true)

	transitGatewayPeering := &goaviatrix.TransitGatewayPeering{
		TransitGatewayName1: d.Get("transit_gateway_name1").(string),
		TransitGatewayName2: d.Get("transit_gateway_name2").(string),
	}
	if d.HasChange("gateway1_excluded_network_cidrs") || d.HasChange("gateway2_excluded_network_cidrs") ||
		d.HasChange("gateway1_excluded_tgw_connections") || d.HasChange("gateway2_excluded_tgw_connections") {
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

		transitGatewayPeering.Gateway1ExcludedCIDRs = strings.Join(gw1Cidrs, ",")
		transitGatewayPeering.Gateway2ExcludedCIDRs = strings.Join(gw2Cidrs, ",")
		transitGatewayPeering.Gateway1ExcludedTGWConnections = strings.Join(gw1Tgws, ",")
		transitGatewayPeering.Gateway2ExcludedTGWConnections = strings.Join(gw2Tgws, ",")

		log.Printf("[INFO] Updating Aviatrix Transit Gateway peering: %#v", transitGatewayPeering)
		err := client.UpdateTransitGatewayPeering(transitGatewayPeering)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix Transit Gateway peering: %s", err)
		}
	}

	if d.HasChange("prepend_as_path1") {
		var prependASPath []string
		for _, v := range d.Get("prepend_as_path1").([]interface{}) {
			prependASPath = append(prependASPath, v.(string))
		}

		err := client.EditTransitConnectionASPathPrepend(transitGatewayPeering, prependASPath)
		if err != nil {
			return fmt.Errorf("could not update prepend_as_path1: %v", err)
		}

	}

	if d.HasChange("prepend_as_path2") {
		var prependASPath []string
		for _, v := range d.Get("prepend_as_path2").([]interface{}) {
			prependASPath = append(prependASPath, v.(string))
		}
		transitGwPeering := &goaviatrix.TransitGatewayPeering{
			TransitGatewayName1: d.Get("transit_gateway_name2").(string),
			TransitGatewayName2: d.Get("transit_gateway_name1").(string),
		}
		err := client.EditTransitConnectionASPathPrepend(transitGwPeering, prependASPath)
		if err != nil {
			return fmt.Errorf("could not update prepend_as_path2: %v", err)
		}

	}

	d.Partial(false)
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
