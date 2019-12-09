package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func resourceTransitVpcMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found AVIATRIX Transit Vpc State v0; migrating to v1")
		return migrateTransitVpcStateV0toV1(is)
	case 1:
		log.Println("[INFO] Found AVIATRIX Transit Vpc State v1; migrating to v2")
		return migrateTransitVpcStateV1toV2(is)
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

func migrateTransitVpcStateV1toV2(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() || is.Attributes == nil {
		log.Println("[DEBUG] Empty Transit Vpc State; nothing to migrate.")
		return is, nil
	}
	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)

	if is.Attributes["vpc_id"] == "" && is.Attributes["vnet_name_resource_group"] != "" {
		is.Attributes["vpc_id"] = is.Attributes["vnet_name_resource_group"]
		delete(is.Attributes, "vnet_name_resource_group")
	}

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}
