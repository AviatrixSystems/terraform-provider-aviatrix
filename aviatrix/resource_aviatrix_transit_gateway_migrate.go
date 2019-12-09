package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func resourceAviatrixTransitGatewayMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found AVIATRIX Transit Gateway State v0; migrating to v1")
		return migrateTransitGatewayStateV0toV1(is)
	default:
		return is, fmt.Errorf("unexpected schema version: %d", v)
	}
}

func migrateTransitGatewayStateV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() || is.Attributes == nil {
		log.Println("[DEBUG] Empty Transit Gateway State; nothing to migrate.")
		return is, nil
	}
	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)

	if is.Attributes["enable_firenet_interfaces"] == "true" {
		is.Attributes["enable_firenet"] = "true"
	} else {
		is.Attributes["enable_firenet"] = "false"
	}
	delete(is.Attributes, "enable_firenet_interfaces")

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}
