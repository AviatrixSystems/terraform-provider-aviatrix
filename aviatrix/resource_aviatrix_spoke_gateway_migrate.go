package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func resourceAviatrixSpokeGatewayMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found AVIATRIX Spoke Gateway State v0; migrating to v1")
		return migrateSpokeGatewayStateV0toV1(is)
	default:
		return is, fmt.Errorf("unexpected schema version: %d", v)
	}
}

func migrateSpokeGatewayStateV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() || is.Attributes == nil {
		log.Println("[DEBUG] Empty Spoke Gateway State; nothing to migrate.")
		return is, nil
	}
	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)

	if (is.Attributes["enable_snat"] == "true" || is.Attributes["enable_snat"] == "True") && is.Attributes["snat_mode"] == "primary" {
		is.Attributes["single_ip_snat"] = "true"
	}

	delete(is.Attributes, "snat_mode")
	delete(is.Attributes, "snat_policy")
	delete(is.Attributes, "dnat_policy")

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}
