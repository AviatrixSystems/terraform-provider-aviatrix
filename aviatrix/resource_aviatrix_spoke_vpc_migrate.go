package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func resourceSpokeVpcMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found AVIATRIX Spoke Vpc State v0; migrating to v1")
		return migrateSpokeVpcStateV0toV1(is)
	default:
		return is, fmt.Errorf("unexpected schema version: %d", v)
	}
}

func migrateSpokeVpcStateV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() || is.Attributes == nil {
		log.Println("[DEBUG] Empty Spoke Vpc State; nothing to migrate.")
		return is, nil
	}
	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)

	if is.Attributes["vpc_id"] == "" && is.Attributes["vnet_and_resource_group_names"] != "" {
		is.Attributes["vpc_id"] = is.Attributes["vnet_and_resource_group_names"]
		delete(is.Attributes, "vnet_and_resource_group_names")
	}

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}
