package aviatrix

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/terraform"
)

func resourceAviatrixFQDNMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found AVIATRIX FQDN State v0; migrating to v1")
		return migrateFQDNStateV0toV1(is)
	default:
		return is, fmt.Errorf("unexpected schema version: %d", v)
	}
}

func migrateFQDNStateV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() || is.Attributes == nil {
		log.Println("[DEBUG] Empty FQDN State; nothing to migrate.")
		return is, nil
	}
	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)

	is.Attributes["gw_filter_tag_list.#"] = is.Attributes["gw_list.#"]
	delete(is.Attributes, "gw_list.#")
	num := 0
	for k, v := range is.Attributes {
		if strings.HasPrefix(k, "gw_list.") {
			is.Attributes["gw_filter_tag_list."+strconv.Itoa(num)+".gw_name"] = v
			is.Attributes["gw_filter_tag_list."+strconv.Itoa(num)+".source_ip_list.#"] = "0"
			num += 1
			delete(is.Attributes, k)
		}
	}

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}
