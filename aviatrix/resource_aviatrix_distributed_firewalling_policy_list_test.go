package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixDistributedFirewallingPolicyList_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DISTRIBUTED_FIREWALLING_POLICY_LIST")
	if skipAcc == "yes" {
		t.Skip("Skipping Distributed Firewalling Policy List test as SKIP_DISTRIBUTED_FIREWALLING_POLICY_LIST is set")
	}
	resourceName := "aviatrix_distributed_firewalling_policy_list.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccDistributedFirewallingPolicyListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDistributedFirewallingPolicyListBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDistributedFirewallingPolicyListExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "policies.0.name", "test-distributed-firewalling-policy"),
					resource.TestCheckResourceAttr(resourceName, "policies.0.action", "PERMIT"),
					resource.TestCheckResourceAttr(resourceName, "policies.0.priority", "0"),
					resource.TestCheckResourceAttr(resourceName, "policies.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "policies.0.src_smart_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policies.0.dst_smart_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policies.0.egress_path", "EGRESS_PATH_DEFAULT"),
				),
			},
			{
				Config: testAccDistributedFirewallingPolicyListLocalEgress(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDistributedFirewallingPolicyListExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "policies.0.egress_path", "EGRESS_PATH_LOCAL"),
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

func testAccDistributedFirewallingPolicyListBasic() string {
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

resource "aviatrix_distributed_firewalling_policy_list" "test" {
	policies {
		name             = "test-distributed-firewalling-policy"
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

func testAccDistributedFirewallingPolicyListLocalEgress() string {
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

resource "aviatrix_distributed_firewalling_policy_list" "test" {
	policies {
		name             = "test-distributed-firewalling-policy"
		action           = "PERMIT"
		logging          = true
		priority         = 0
		protocol         = "TCP"
		egress_path      = "EGRESS_PATH_LOCAL"
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

func testAccCheckDistributedFirewallingPolicyListExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no Distributed-firewalling Policy List resource found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Distributed-firewalling Policy List ID is set")
		}

		client := mustClient(testAccProviderVersionValidation.Meta())

		_, err := client.GetDistributedFirewallingPolicyList(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get Distributed-firewalling Policy List status: %w", err)
		}

		if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
			return fmt.Errorf("distributed-firewalling policy list ID not found")
		}

		return nil
	}
}

func testAccDistributedFirewallingPolicyListDestroy(s *terraform.State) error {
	client := mustClient(testAccProviderVersionValidation.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_smart_group" {
			continue
		}

		_, err := client.GetDistributedFirewallingPolicyList(context.Background())
		if err == nil || !errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("distributed-firewalling policy list configured when it should be destroyed")
		}
	}

	return nil
}

// Verify the enforcement field is Optional+Computed with a validator, and watch is Optional+Computed+Deprecated with no Default.
func TestDistributedFirewallingPolicyListSchema_EnforcementAndWatch(t *testing.T) {
	res := resourceAviatrixDistributedFirewallingPolicyList()
	policyElem, ok := res.Schema["policies"].Elem.(*schema.Resource)
	if !ok {
		t.Fatal("expected policies Elem to be *schema.Resource")
	}

	enforcementField := policyElem.Schema["enforcement"]
	assert.Equal(t, schema.TypeString, enforcementField.Type)
	assert.True(t, enforcementField.Optional)
	assert.True(t, enforcementField.Computed)
	assert.NotNil(t, enforcementField.ValidateFunc)

	watchField := policyElem.Schema["watch"]
	assert.Equal(t, schema.TypeBool, watchField.Type)
	assert.True(t, watchField.Optional)
	assert.True(t, watchField.Computed)
	assert.False(t, watchField.Optional && watchField.Default != nil, "watch must not have a Default when Computed")
	assert.NotEmpty(t, watchField.Deprecated, "watch must be marked as deprecated")
}

// Verify enforcement is resolved from the user's config: an explicit enforcement
// wins, otherwise watch maps to MONITOR/ENFORCE, otherwise ENFORCE. watch is
// always sent consistent with enforcement.
func TestResolveDistributedFirewallingPolicyEnforcement(t *testing.T) {
	nullStr, nullBool := cty.NullVal(cty.String), cty.NullVal(cty.Bool)
	tests := []struct {
		name        string
		enforcement cty.Value
		watch       cty.Value
		want        string
	}{
		{"enforcement=MONITOR", cty.StringVal("MONITOR"), nullBool, "MONITOR"},
		{"enforcement=ENFORCE", cty.StringVal("ENFORCE"), nullBool, "ENFORCE"},
		{"enforcement=DISABLE", cty.StringVal("DISABLE"), nullBool, "DISABLE"},
		{"watch=true", nullStr, cty.True, "MONITOR"},
		{"watch=false", nullStr, cty.False, "ENFORCE"},
		{"neither set", nullStr, nullBool, "ENFORCE"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := cty.ObjectVal(map[string]cty.Value{"enforcement": tt.enforcement, "watch": tt.watch})
			enforcement, watch := resolveDistributedFirewallingPolicyEnforcement(raw)
			assert.Equal(t, tt.want, enforcement)
			assert.Equal(t, tt.want == "MONITOR", watch)
		})
	}
}

// Verify http_filters UUIDs are passed through to the policy request, and omitted when unset.
func TestMarshalDistributedFirewallingPolicyListInput_HTTPFilters(t *testing.T) {
	const httpFilterUUID = "0d1e2f3a-4b5c-6d7e-8f90-a1b2c3d4e5f6"
	tests := []struct {
		name        string
		httpFilters []any
		want        []string
	}{
		{name: "unset", httpFilters: nil, want: nil},
		{name: "single filter", httpFilters: []any{httpFilterUUID}, want: []string{httpFilterUUID}},
	}

	res := resourceAviatrixDistributedFirewallingPolicyList()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy := map[string]any{
				"name":             "test-pol",
				"action":           "PERMIT",
				"src_smart_groups": []any{},
				"dst_smart_groups": []any{},
				"protocol":         "TCP",
			}
			if tt.httpFilters != nil {
				policy["http_filters"] = tt.httpFilters
			}
			d := schema.TestResourceDataRaw(t, res.Schema, map[string]any{"policies": []any{policy}})

			policyList, err := marshalDistributedFirewallingPolicyListInput(d)
			assert.NoError(t, err)
			assert.Len(t, policyList.Policies, 1)
			assert.Equal(t, tt.want, policyList.Policies[0].HTTPFilters)
		})
	}
}

// Verify the read path derives watch=true only for enforcement=monitor, and false for all other values.
func TestReadDistributedFirewallingPolicy_WatchDerivedFromEnforcement(t *testing.T) {
	tests := []struct {
		enforcement string
		wantWatch   bool
	}{
		{"MONITOR", true},
		{"ENFORCE", false},
		{"DISABLE", false},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.enforcement, func(t *testing.T) {
			policy := goaviatrix.DistributedFirewallingPolicy{
				Enforcement: tt.enforcement,
			}
			got := policy.Enforcement == "MONITOR"
			assert.Equal(t, tt.wantWatch, got)
		})
	}
}

// Verify the CustomizeDiff validation rejects policies that set both watch and enforcement,
// and allows policies that set only one.
func TestDistributedFirewallingPolicyListCustomizeDiff_RejectsBothWatchAndEnforcement(t *testing.T) {
	if err := validateDistributedFirewallingPolicyWatchEnforcement(0, true, true); err == nil {
		t.Error("expected error when both watch and enforcement are set, got nil")
	}
	if err := validateDistributedFirewallingPolicyWatchEnforcement(0, false, true); err != nil {
		t.Errorf("expected no error when only enforcement is set, got: %v", err)
	}
	if err := validateDistributedFirewallingPolicyWatchEnforcement(0, true, false); err != nil {
		t.Errorf("expected no error when only watch is set, got: %v", err)
	}
	if err := validateDistributedFirewallingPolicyWatchEnforcement(0, false, false); err != nil {
		t.Errorf("expected no error when neither is set, got: %v", err)
	}
}
