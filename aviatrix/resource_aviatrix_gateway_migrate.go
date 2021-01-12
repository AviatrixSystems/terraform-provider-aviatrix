package aviatrix

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixGatewayResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"monitor_exclude_list": {
				Type: schema.TypeString,
			},
		},
	}
}

func resourceAviatrixGatewayStateUpgradeV0(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if v, ok := rawState["monitor_exclude_list"]; ok {
		excludedInstances, ok := v.(string)
		if !ok {
			rawState["monitor_exclude_list"] = []string{}
		} else {
			rawState["monitor_exclude_list"] = strings.Split(excludedInstances, ",")
		}
	}
	return rawState, nil
}
