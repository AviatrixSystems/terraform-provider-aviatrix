package aviatrix

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestMarshalDCFPolicyBlockInputScopes(t *testing.T) {
	const scopeUUID = "123e4567-e89b-12d3-a456-426614174000"
	resourceSchema := resourceAviatrixDCFPolicyGroup().Schema

	withScopes := schema.TestResourceDataRaw(t, resourceSchema, map[string]any{
		"name": "group",
		"scope_config": []any{map[string]any{
			"scope_uuids": []any{scopeUUID},
		}},
		"ruleset_reference": []any{map[string]any{
			"priority":    1,
			"target_uuid": "ruleset-uuid",
		}},
	})
	policyBlock, err := marshalDCFPolicyBlockInput(withScopes)
	if err != nil {
		t.Fatalf("marshal policy group with scopes: %v", err)
	}
	if policyBlock.ScopeConfig == nil || len(policyBlock.ScopeConfig.ScopeUUIDs) != 1 ||
		policyBlock.ScopeConfig.ScopeUUIDs[0] != scopeUUID {
		t.Fatalf("unexpected policy group scope config: %#v", policyBlock.ScopeConfig)
	}

	withoutScopes := schema.TestResourceDataRaw(t, resourceSchema, map[string]any{
		"name": "group",
		"ruleset_reference": []any{map[string]any{
			"priority":    1,
			"target_uuid": "ruleset-uuid",
		}},
	})
	policyBlock, err = marshalDCFPolicyBlockInput(withoutScopes)
	if err != nil {
		t.Fatalf("marshal policy group without scopes: %v", err)
	}
	if policyBlock.ScopeConfig != nil {
		t.Fatalf("scope_config must be omitted when scope_uuids is unset: %#v", policyBlock.ScopeConfig)
	}
}

func TestAccAviatrixDcfPolicyGroup_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DCF_POLICY_GROUP")
	if skipAcc == "yes" {
		t.Skip("Skipping DCF Policy Group test as SKIP_DCF_POLICY_GROUP is set")
	}
	resourceName := "aviatrix_dcf_policy_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckDcfPolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDcfPolicyGroupBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcfPolicyGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-dcf-policy-group"),
					resource.TestCheckResourceAttr(resourceName, "scope_config.0.scope_uuids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ruleset_reference.0.priority", "0"),
					resource.TestCheckResourceAttr(resourceName, "policy_group_reference.0.priority", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckDcfPolicyGroupBasic() string {
	return `resource "aviatrix_dcf_policy_scope" "test" {
	name = "test-policy-scope"
	included {
		preset = ["all"]
	}
}

resource "aviatrix_smart_group" "ad1" {
	name = "test-smart_group-1"
	selector {
		match_expressions {
			cidr = "10.0.0.0/16"
		}
	}
}

resource "aviatrix_smart_group" "ad2" {
	name = "test-smart-group-2"
	selector {
		match_expressions {
			cidr = "11.0.0.0/16"
		}
	}
}

resource "aviatrix_dcf_ruleset" "test_list" {
	name = "test-dcf-ruleset"
	scope_config {
		scope_uuids = [aviatrix_dcf_policy_scope.test.uuid]
	}
	rules {
		name             = "test-distributed-firewalling-rule"
		action           = "PERMIT"
		logging          = true
		priority         = 0
		protocol         = "TCP"
		src_smart_groups = [
		  aviatrix_smart_group.ad1.uuid
		]
		dst_smart_groups = [
		  aviatrix_smart_group.ad2.uuid
		]

		port_ranges {
		  hi = 10
		  lo = 1
		}
  }
}

resource "aviatrix_dcf_policy_group" "nested_group" {
	name = "test-nested-dcf-policy-group"
	ruleset_reference {
		priority = 0
		target_uuid = aviatrix_dcf_ruleset.test_list.id
	}
}

resource "aviatrix_dcf_policy_group" "test" {
	name = "test-dcf-policy-group"
	scope_config {
		scope_uuids = [aviatrix_dcf_policy_scope.test.uuid]
	}
	ruleset_reference {
		priority = 0
		target_uuid = aviatrix_dcf_ruleset.test_list.id
	}
	policy_group_reference {
		priority = 1
		target_uuid = aviatrix_dcf_policy_group.nested_group.id
	}
}
`
}

// TestAccAviatrixDcfPolicyGroup_removeScope verifies removing scope_config clears the policy group scopes.
// Requires a controller that clears scope_config when a policy group PUT omits it.
func TestAccAviatrixDcfPolicyGroup_removeScope(t *testing.T) {
	if os.Getenv("SKIP_DCF_POLICY_GROUP") == "yes" {
		t.Skip("Skipping DCF Policy Group test as SKIP_DCF_POLICY_GROUP is set")
	}
	resourceName := "aviatrix_dcf_policy_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"aviatrix": func() (*schema.Provider, error) { return testAccProviderVersionValidation, nil },
		},
		CheckDestroy: testAccCheckDcfPolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcfPolicyGroupWithScope(true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcfPolicyGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "scope_config.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "scope_config.0.scope_uuids.*", "aviatrix_dcf_policy_scope.test", "uuid"),
				),
			},
			{
				Config: testAccDcfPolicyGroupWithScope(false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcfPolicyGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "scope_config.#", "0"),
				),
			},
		},
	})
}

func testAccDcfPolicyGroupWithScope(withScope bool) string {
	scopeConfig := ""
	if withScope {
		scopeConfig = `
	scope_config {
		scope_uuids = [aviatrix_dcf_policy_scope.test.uuid]
	}`
	}
	return fmt.Sprintf(`
resource "aviatrix_dcf_policy_scope" "test" {
	name = "test-policy-group-remove-scope"
	included {
		preset = ["all"]
	}
}

resource "aviatrix_dcf_ruleset" "test_list" {
	name = "test-dcf-policy-group-remove-scope-ruleset"
}

resource "aviatrix_dcf_policy_group" "test" {
	name = "test-dcf-policy-group-remove-scope"%s
	ruleset_reference {
		priority    = 0
		target_uuid = aviatrix_dcf_ruleset.test_list.id
	}
}
`, scopeConfig)
}

func testAccCheckDcfPolicyGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no DCF Policy Group resource found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no DCF Policy Group ID is set")
		}

		client := mustClient(testAccProviderVersionValidation.Meta())

		_, err := client.GetDCFPolicyBlock(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to get DCF Policy Group status: %w", err)
		}

		return nil
	}
}

func testAccCheckDcfPolicyGroupDestroy(s *terraform.State) error {
	client := mustClient(testAccProviderVersionValidation.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_dcf_policy_group" {
			continue
		}

		_, err := client.GetDCFPolicyBlock(context.Background(), rs.Primary.ID)
		if err == nil || !strings.Contains(err.Error(), "does not exist") {
			return fmt.Errorf("dcf policy group configured when it should be destroyed %w", err)
		}
	}

	return nil
}
