package aviatrix

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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

func resourceAviatrixFQDNResourceV1() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"manage_domain_names": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAviatrixFQDNStateUpgradeV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if _, ok := rawState["manage_domain_names"]; !ok {
		rawState["manage_domain_names"] = "true"
	}
	return rawState, nil
}