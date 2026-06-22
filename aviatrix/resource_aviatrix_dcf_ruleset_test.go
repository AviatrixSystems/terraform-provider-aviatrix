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

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestDcfRuleSetHash_protocolCaseInsensitive(t *testing.T) {
	basRule := map[string]interface{}{
		"name":                     "test-rule",
		"action":                   "PERMIT",
		"priority":                 0,
		"protocol":                 "TCP",
		"logging":                  false,
		"watch":                    false,
		"decrypt_policy":           "DECRYPT_UNSPECIFIED",
		"flow_app_requirement":     "APP_UNSPECIFIED",
		"exclude_sg_orchestration": false,
		"tls_profile":              "",
		"uuid":                     "",
		"log_profile":              "def000ad-7000-0000-0000-000000000001",
		"src_smart_groups":         schema.NewSet(schema.HashString, []interface{}{"sg-1"}),
		"dst_smart_groups":         schema.NewSet(schema.HashString, []interface{}{"sg-2"}),
		"web_groups":               schema.NewSet(schema.HashString, []interface{}{}),
		"port_ranges":              []interface{}{},
	}

	upperHash := dcfRuleSetHash(basRule)

	lowercaseRule := make(map[string]interface{}, len(basRule))
	for k, v := range basRule {
		lowercaseRule[k] = v
	}
	lowercaseRule["protocol"] = "tcp"
	lowerHash := dcfRuleSetHash(lowercaseRule)

	if upperHash != lowerHash {
		t.Errorf("hash mismatch: protocol 'TCP' produced %d, protocol 'tcp' produced %d; expected equal hashes", upperHash, lowerHash)
	}

	mixedRule := make(map[string]interface{}, len(basRule))
	for k, v := range basRule {
		mixedRule[k] = v
	}
	mixedRule["protocol"] = "Tcp"
	mixedHash := dcfRuleSetHash(mixedRule)

	if upperHash != mixedHash {
		t.Errorf("hash mismatch: protocol 'TCP' produced %d, protocol 'Tcp' produced %d; expected equal hashes", upperHash, mixedHash)
	}
}

func TestAccAviatrixDcfRuleset_protocolCaseNoDiff(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DCF_RULESET")
	if skipAcc == "yes" {
		t.Skip("Skipping DCF Ruleset test as SKIP_DCF_RULESET is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckDcfRulesetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDcfRulesetWithProtocol("TCP"),
			},
			{
				Config:   testAccCheckDcfRulesetWithProtocol("tcp"),
				PlanOnly: true,
			},
		},
	})
}

func testAccCheckDcfRulesetWithProtocol(protocol string) string {
	return fmt.Sprintf(`
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

resource "aviatrix_dcf_ruleset" "test" {
	name = "test-dcf-ruleset"
	rules {
		name             = "test-distributed-firewalling-rule"
		action           = "PERMIT"
		logging          = true
		priority         = 0
		protocol         = %q
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
`, protocol)
}

func TestAccAviatriDcfRuleset_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DCF_RULESET")
	if skipAcc == "yes" {
		t.Skip("Skipping DCF Ruleset test as SKIP_DCF_RULESET is set")
	}
	resourceName := "aviatrix_dcf_ruleset.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckDcfRulesetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDcfRulesetBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcfRulesetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-dcf-ruleset"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.name", "test-distributed-firewalling-rule"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.action", "PERMIT"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.priority", "0"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.src_smart_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.dst_smart_groups.#", "1"),
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

func testAccCheckDcfRulesetBasic() string {
	return `
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

resource "aviatrix_dcf_ruleset" "test" {
	name = "test-dcf-ruleset"
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
`
}

func TestAccAviatrixDcfRuleset_uuidPreservedOnUpdate(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DCF_RULESET")
	if skipAcc == "yes" {
		t.Skip("Skipping DCF Ruleset test as SKIP_DCF_RULESET is set")
	}
	resourceName := "aviatrix_dcf_ruleset.test"

	var originalUUID string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckDcfRulesetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcfRulesetForUUIDPreservation(true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcfRulesetExists(resourceName),
					testAccCaptureDcfRuleUUID(resourceName, &originalUUID),
				),
			},
			{
				Config: testAccDcfRulesetForUUIDPreservation(false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcfRulesetExists(resourceName),
					testAccCheckDcfRuleUUIDUnchanged(resourceName, &originalUUID),
				),
			},
		},
	})
}

func testAccDcfRulesetForUUIDPreservation(logging bool) string {
	return fmt.Sprintf(`
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

resource "aviatrix_dcf_ruleset" "test" {
	name = "test-dcf-ruleset"
	rules {
		name             = "test-distributed-firewalling-rule"
		action           = "PERMIT"
		logging          = %t
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
`, logging)
}

func findRuleUUID(rs *terraform.ResourceState) (string, error) {
	for key, value := range rs.Primary.Attributes {
		if strings.HasSuffix(key, ".uuid") && strings.HasPrefix(key, "rules.") && value != "" {
			return value, nil
		}
	}
	return "", fmt.Errorf("no rule UUID found in state attributes")
}

func testAccCaptureDcfRuleUUID(resourceName string, dest *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		uuid, err := findRuleUUID(rs)
		if err != nil {
			return fmt.Errorf("after initial creation: %w", err)
		}
		*dest = uuid
		return nil
	}
}

func testAccCheckDcfRuleUUIDUnchanged(resourceName string, originalUUID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		uuid, err := findRuleUUID(rs)
		if err != nil {
			return fmt.Errorf("after update: %w", err)
		}
		if uuid != *originalUUID {
			return fmt.Errorf("rule UUID changed after attribute update: was %q, now %q", *originalUUID, uuid)
		}
		return nil
	}
}

func testAccCheckDcfRulesetExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no DCF Ruleset resource found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no DCF Ruleset ID is set")
		}

		client, ok := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)
		if !ok {
			return fmt.Errorf("failed to assert Meta as *goaviatrix.Client")
		}

		_, err := client.GetDCFPolicyList(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to get DCF Ruleset status: %w", err)
		}

		return nil
	}
}

func testAccCheckDcfRulesetDestroy(s *terraform.State) error {
	client, ok := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)
	if !ok {
		return fmt.Errorf("failed to assert Meta as *goaviatrix.Client")
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_smart_group" {
			continue
		}

		_, err := client.GetDCFPolicyList(context.Background(), rs.Primary.ID)
		if err == nil || !strings.Contains(err.Error(), "does not exist") {
			return fmt.Errorf("dcf ruleset configured when it should be destroyed %w", err)
		}
	}

	return nil
}
