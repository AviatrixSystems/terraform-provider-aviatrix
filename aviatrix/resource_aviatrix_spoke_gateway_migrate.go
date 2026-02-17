package aviatrix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixSpokeGatewayResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"manage_transit_gateway_attachment": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAviatrixSpokeGatewayStateUpgradeV0(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if _, ok := rawState["manage_transit_gateway_attachment"]; !ok {
		rawState["manage_transit_gateway_attachment"] = "true"
	}

	return rawState, nil
}

func resourceAviatrixSpokeGatewayResourceV1() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"manage_ha_gateway": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAviatrixSpokeGatewayStateUpgradeV1(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if _, ok := rawState["manage_ha_gateway"]; !ok {
		rawState["manage_ha_gateway"] = "true"
	}

	return rawState, nil
}
