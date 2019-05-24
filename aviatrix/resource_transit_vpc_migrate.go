package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/terraform"
)

func resourceTransitVpcMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found AVIATRIX Transit Vpc State v0; migrating to v1")
		return migrateTransitVpcStateV0toV1(is)
	default:
		return is, fmt.Errorf("unexpected schema version: %d", v)
	}
}

func migrateTransitVpcStateV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() || is.Attributes == nil {
		log.Println("[DEBUG] Empty Transit Vpc State; nothing to migrate.")
		return is, nil
	}
	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)

	is.Attributes["enable_firenet_interfaces"] = "false"

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}
