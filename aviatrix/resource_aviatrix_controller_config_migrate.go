package aviatrix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixControllerConfigResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"sg_management_account_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_management": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAviatrixControllerConfigStateUpgradeV0(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if _, ok := rawState["sg_management_account_name"]; !ok {
		delete(rawState, "sg_management_account_name")
	}

	if _, ok := rawState["security_group_management"]; !ok {
		delete(rawState, "security_group_management")
	}

	return rawState, nil
}
