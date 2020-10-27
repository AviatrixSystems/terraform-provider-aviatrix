package aviatrix

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

func resourceAviatrixSpokeGatewayStateUpgradeV0(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	rawState["manage_transit_gateway_attachment"] = "true"

	return rawState, nil
}
