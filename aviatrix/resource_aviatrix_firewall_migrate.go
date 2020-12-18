package aviatrix

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixFirewallResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"manage_firewall_policies": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAviatrixFirewallStateUpgradeV0(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if _, ok := rawState["manage_firewall_policies"]; !ok {
		rawState["manage_firewall_policies"] = true
	}
	return rawState, nil
}
