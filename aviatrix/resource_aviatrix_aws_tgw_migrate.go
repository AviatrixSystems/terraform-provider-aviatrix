package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func resourceAviatrixAWSTgwMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found AVIATRIX AWS TGW State v0; migrating to v1")
		return migrateAWSTgwStateV0toV1(is)
	default:
		return is, fmt.Errorf("unexpected schema version: %d", v)
	}
}

func migrateAWSTgwStateV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() || is.Attributes == nil {
		log.Println("[DEBUG] Empty AWS TGW State; nothing to migrate.")
		return is, nil
	}
	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)

	is.Attributes["manage_vpc_attachment"] = "true"

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}

func resourceAviatrixAWSTgwResourceV1() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"manage_transit_gateway_attachment": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAviatrixAWSTgwStateUpgradeV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if _, ok := rawState["manage_transit_gateway_attachment"]; !ok {
		rawState["manage_transit_gateway_attachment"] = true
	}
	return rawState, nil
}
