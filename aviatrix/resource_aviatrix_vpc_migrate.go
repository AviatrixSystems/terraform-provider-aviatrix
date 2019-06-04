package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/terraform"
)

func resourceAviatrixVpcMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found AVIATRIX Vpc State v0; migrating to v1")
		return migrateVpcStateV0toV1(is)
	default:
		return is, fmt.Errorf("unexpected schema version: %d", v)
	}
}

func migrateVpcStateV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() || is.Attributes == nil {
		log.Println("[DEBUG] Empty Vpc State; nothing to migrate.")
		return is, nil
	}
	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)

	is.Attributes["aviatrix_firenet_vpc"] = "false"

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}
