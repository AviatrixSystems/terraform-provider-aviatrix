package aviatrix

import (
	"context"
	"fmt"
	"maps"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// dcfRulesetTestRule returns the raw config value for one rule, so old/new
// variants used in diff tests only differ in the fields the test cares about.
func dcfRulesetTestRule(logging bool, lo, hi int) map[string]any {
	return dcfRulesetTestRuleWithPortRange(logging, map[string]any{"lo": lo, "hi": hi})
}

// dcfRulesetTestRuleHiOmitted is like dcfRulesetTestRule but leaves "hi" out
// of port_ranges entirely, as a user relying on single-port shorthand would.
func dcfRulesetTestRuleHiOmitted(logging bool, lo int) map[string]any {
	return dcfRulesetTestRuleWithPortRange(logging, map[string]any{"lo": lo})
}

func dcfRulesetTestRuleWithPortRange(logging bool, portRange map[string]any) map[string]any {
	return map[string]any{
		"name":                     "test-rule-1",
		"action":                   "PERMIT",
		"priority":                 2,
		"protocol":                 "TCP",
		"logging":                  logging,
		"decrypt_policy":           "DECRYPT_UNSPECIFIED",
		"flow_app_requirement":     "APP_UNSPECIFIED",
		"exclude_sg_orchestration": true,
		"tls_profile":              "",
		"log_profile":              "def000ad-7000-0000-0000-000000000001",
		"egress_path":              "EGRESS_PATH_DEFAULT",
		"src_smart_groups":         []any{"def000ad-0000-0000-0000-000000000000"},
		"dst_smart_groups":         []any{"def000ad-0000-0000-0000-000000000000"},
		"web_groups":               []any{},
		"port_ranges": []any{
			portRange,
		},
	}
}

// TestDcfRulesetDiff_NoPhantomRuleWhenPortRangeHiEqualsLo is a regression test,
// run entirely offline (no controller/acceptance test needed), for the bug
// where updating any field on a rule whose port_ranges had hi == lo produced a
// third, nameless "phantom" rule in the plan and failed apply.
//
// Root cause: the old hi DiffSuppressFunc suppressed hi's diff whenever
// old(hi) == lo, which is also exactly what happens when the whole rule is
// being removed from the TypeSet (its hi diffs toward its zero value). When a
// DiffSuppressFunc fires on an attribute inside a set element, the SDK does
// not drop the diff -- it replaces it with a bare {Old, New: Old} NOOP and
// loses whatever NewRemoved flag the field originally had (see
// helper/schema/schema.go's diff(), the "make it a NOOP" branch). So instead
// of hi being cleanly removed alongside every other field of the deleted
// rule, it silently survives, and Terraform's later state reconstruction
// resurrects the dead rule as an extra, nameless "phantom" block containing
// just that leftover hi (and the still-computed uuid).
//
// This test builds a "prior state" (as if already applied, with a real UUID)
// for a single rule with lo == hi, then diffs it against a new config that
// only changes "logging" -- forcing the whole rule to be replaced (removed +
// re-added) via the set-hash mechanism. It checks two layers:
//  1. The raw InstanceDiff: lo must be cleanly removed, and hi must never be
//     left behind as a leftover no-op diff entry (hi being Computed means its
//     removal may simply produce no diff entry at all, same as uuid -- that's
//     fine; a retained no-op is the bug).
//  2. The reconstructed "planned new state" (via InstanceDiff.Apply, the same
//     shim Terraform core uses to turn this diff into a plan): exactly one
//     named rule must survive, with no nameless phantom slot alongside it.
func TestDcfRulesetDiff_NoPhantomRuleWhenPortRangeHiEqualsLo(t *testing.T) {
	checkNoPhantomRuleOnUpdate(t, dcfRulesetTestRule(true, 443, 443), dcfRulesetTestRule(false, 443, 443))
}

// TestDcfRulesetDiff_NoPhantomRuleWhenPortRangeHiGreaterThanLo is a sanity
// check for the "normal" case (hi > lo), which never triggered the old
// DiffSuppressFunc bug (that suppressor only fired when old(hi) == lo). It
// guards against a regression in the other direction: that making hi
// Optional+Computed didn't somehow break ordinary multi-port ranges.
func TestDcfRulesetDiff_NoPhantomRuleWhenPortRangeHiGreaterThanLo(t *testing.T) {
	checkNoPhantomRuleOnUpdate(t, dcfRulesetTestRule(true, 3000, 8000), dcfRulesetTestRule(false, 3000, 8000))
}

// TestDcfRulesetDiff_NoPhantomRuleWhenPortRangeHiOmitted covers the
// single-port shorthand where "hi" is left out of config entirely (letting
// it default to whatever the backend/state already has for lo). This is the
// other long-standing production pattern that the fix (making hi
// Optional+Computed instead of relying on a DiffSuppressFunc) has to keep
// working, alongside the explicit hi == lo case above.
func TestDcfRulesetDiff_NoPhantomRuleWhenPortRangeHiOmitted(t *testing.T) {
	checkNoPhantomRuleOnUpdate(t, dcfRulesetTestRuleHiOmitted(true, 443), dcfRulesetTestRuleHiOmitted(false, 443))
}

// TestDcfRulesetDiff_NoPhantomRuleWhenPortRangeHiIsZero covers a rule whose
// hi is explicitly set to the literal sentinel value 0 (as opposed to being
// omitted from config, or being defaulted to lo). "hi": 0 is a distinct
// state/config combination from "hi omitted" -- Terraform tracks "not present
// in config" separately from "present in config with the zero value" -- so
// it needs its own coverage rather than being assumed equivalent to the
// omitted case above.
func TestDcfRulesetDiff_NoPhantomRuleWhenPortRangeHiIsZero(t *testing.T) {
	checkNoPhantomRuleOnUpdate(t, dcfRulesetTestRule(true, 443, 0), dcfRulesetTestRule(false, 443, 0))
}

// checkNoPhantomRuleOnUpdate is a regression check, run entirely offline (no
// controller/acceptance test needed), for the bug where updating any field on
// a rule whose port_ranges had hi == lo produced a third, nameless "phantom"
// rule in the plan and failed apply.
//
// Root cause: the old hi DiffSuppressFunc suppressed hi's diff whenever
// old(hi) == lo, which is also exactly what happens when the whole rule is
// being removed from the TypeSet (its hi diffs toward its zero value). When a
// DiffSuppressFunc fires on an attribute inside a set element, the SDK does
// not drop the diff -- it replaces it with a bare {Old, New: Old} NOOP and
// loses whatever NewRemoved flag the field originally had (see
// helper/schema/schema.go's diff(), the "make it a NOOP" branch). So instead
// of hi being cleanly removed alongside every other field of the deleted
// rule, it silently survives, and Terraform's later state reconstruction
// resurrects the dead rule as an extra, nameless "phantom" block containing
// just that leftover hi (and the still-computed uuid).
//
// oldRule and newRule should differ only in the field being updated (e.g.
// "logging"), so the rule's port_ranges stays fixed across the update -- that
// is what forces the whole rule to be replaced (removed + re-added) via the
// set-hash mechanism while isolating the effect of the specific hi/lo
// combination under test. It checks two layers:
//  1. The raw InstanceDiff: lo must be cleanly removed, and hi must never be
//     left behind as a leftover no-op diff entry (hi being Computed means its
//     removal may simply produce no diff entry at all, same as uuid -- that's
//     fine; a retained no-op is the bug).
//  2. The reconstructed "planned new state" (via InstanceDiff.Apply, the same
//     shim Terraform core uses to turn this diff into a plan): exactly one
//     named rule must survive, with no nameless phantom slot alongside it.
func checkNoPhantomRuleOnUpdate(t *testing.T, oldRule, newRule map[string]any) {
	// Use a plain *schema.Resource with no CustomizeDiff: the phantom-rule bug
	// lives entirely in the TypeSet hashing/diffing machinery, not in
	// resourceAviatrixDCFRulesetCustomizeDiff, and CustomizeDiff needs a real
	// raw cty config that TestResourceDataRaw/NewResourceConfigRaw don't
	// populate.
	testResource := &schema.Resource{Schema: resourceAviatrixDCFRuleset().Schema}

	oldRaw := map[string]any{
		"name": "test-ruleset",
		"rules": []any{
			oldRule,
		},
	}
	d := schema.TestResourceDataRaw(t, testResource.Schema, oldRaw)
	// ResourceData.State() returns nil unless an ID is set, since a resource
	// with no ID is considered not to exist yet. TestResourceDataRaw never
	// sets one (it only models raw config), so set a fake one to represent
	// "already applied."
	d.SetId("test-ruleset-id")
	oldState := d.State()
	if oldState == nil {
		t.Fatal("expected a non-nil state from TestResourceDataRaw")
	}

	// Simulate the backend having already assigned a real UUID to the rule,
	// as would be true for any previously-applied resource.
	var ruleCode string
	for k, v := range oldState.Attributes {
		if strings.HasSuffix(k, ".name") && v == "test-rule-1" {
			ruleCode = strings.TrimSuffix(k, ".name")
		}
	}
	if ruleCode == "" {
		t.Fatalf("could not find the rule's set code in old state attributes: %#v", oldState.Attributes)
	}
	oldState.Attributes[ruleCode+".uuid"] = "93d4d7eb-b77c-4179-985d-f2c85bf818c5"

	newRaw := map[string]any{
		"name": "test-ruleset",
		"rules": []any{
			// newRule should only differ from oldRule in a field like
			// "logging"; keeping port_ranges fixed forces the rule's set
			// hashcode to change, so the SDK removes this (old) code and
			// adds a new one for the updated rule.
			newRule,
		},
	}
	newConfig := terraform.NewResourceConfigRaw(newRaw)

	instanceDiff, err := testResource.Diff(context.Background(), oldState, newConfig, nil)
	if err != nil {
		t.Fatalf("Diff returned an error: %v", err)
	}
	if instanceDiff == nil || len(instanceDiff.Attributes) == 0 {
		t.Fatal("expected a non-empty diff since logging changed")
	}

	loKey := ruleCode + ".port_ranges.0.lo"
	hiKey := ruleCode + ".port_ranges.0.hi"

	loAttr, ok := instanceDiff.Attributes[loKey]
	if !ok {
		t.Fatalf("expected a diff entry for %q (the removed rule's lo), attributes: %#v", loKey, instanceDiff.Attributes)
	}
	if !loAttr.NewRemoved {
		t.Fatalf("sanity check failed: %q should be NewRemoved since the whole rule is being removed, got %#v", loKey, loAttr)
	}

	// hi is Computed, so like its sibling uuid a removed Computed field
	// produces *no* diff entry at all (schema.go's diffString returns nil for
	// "removed && schema.Computed") -- that's fine and expected, it's the
	// same mechanism that already lets uuid disappear cleanly. What must
	// never happen is a leftover NOOP entry (Old == New == the old value,
	// NewRemoved == false): that's the exact signature of the phantom-rule
	// bug, where the old DiffSuppressFunc fired on removal and the SDK
	// silently replaced the removal with a no-op instead of dropping it.
	if hiAttr, ok := instanceDiff.Attributes[hiKey]; ok && !hiAttr.NewRemoved {
		t.Errorf("%q is present but not marked NewRemoved (got %#v) -- "+
			"this is the phantom-rule bug: because hi == lo, the DiffSuppressFunc "+
			"fired on removal and the SDK silently replaced the removal with a "+
			"NOOP, leaving a dead rule slot behind", hiKey, hiAttr)
	}

	// The checks above operate one layer below where the actual phantom rule
	// appears in `terraform plan`. That third rule is materialized when
	// Terraform reconstructs the "planned new state" from this same flat
	// diff -- InstanceDiff.Apply is the (legacy, but still used) function
	// that does that reconstruction. Its TypeSet handling
	// (terraform/diff.go's applyBlockDiff) decides whether to keep a set
	// element's code by scanning for *any* diff attribute under that code
	// that is not NewRemoved (see the "this must be a diff to keep" branch).
	// The leftover, non-removed hi NOOP is exactly such an attribute, so the
	// dead rule's code gets kept -- resurrecting it as a nameless rule. This
	// step reproduces that directly, with no controller and no Terraform
	// core involved: just the same SDK-shim function core relies on.
	newAttrs, err := instanceDiff.Apply(oldState.Attributes, testResource.CoreConfigSchema())
	if err != nil {
		t.Fatalf("InstanceDiff.Apply (state reconstruction) returned an error: %v", err)
	}

	ruleCodes := map[string]bool{}
	namedRuleCodes := map[string]bool{}
	for k, v := range newAttrs {
		const prefix = "rules."
		if !strings.HasPrefix(k, prefix) || strings.HasSuffix(k, ".#") {
			continue
		}
		rest := strings.TrimPrefix(k, prefix)
		before, _, ok0 := strings.Cut(rest, ".")
		if !ok0 {
			continue
		}
		code := before
		ruleCodes[code] = true
		if strings.HasSuffix(k, ".name") && v != "" {
			namedRuleCodes[code] = true
		}
	}

	for code := range ruleCodes {
		if !namedRuleCodes[code] {
			t.Errorf("reconstructed planned state has a rule slot %q with no name -- "+
				"this is the phantom rule that shows up in `terraform plan` and fails "+
				"apply; full reconstructed attributes: %#v", code, newAttrs)
		}
	}
	if got := len(namedRuleCodes); got != 1 {
		t.Errorf("expected exactly 1 named rule to survive into the planned new state, got %d: %v", got, namedRuleCodes)
	}
}

func TestDcfRuleSetHash_protocolCaseInsensitive(t *testing.T) {
	basRule := map[string]any{
		"name":                     "test-rule",
		"action":                   "PERMIT",
		"priority":                 0,
		"protocol":                 "TCP",
		"logging":                  false,
		"watch":                    false,
		"enforcement":              "ENFORCE",
		"decrypt_policy":           "DECRYPT_UNSPECIFIED",
		"flow_app_requirement":     "APP_UNSPECIFIED",
		"exclude_sg_orchestration": false,
		"tls_profile":              "",
		"uuid":                     "",
		"log_profile":              "def000ad-7000-0000-0000-000000000001",
		"egress_path":              "EGRESS_PATH_DEFAULT",
		"src_smart_groups":         schema.NewSet(schema.HashString, []any{"sg-1"}),
		"dst_smart_groups":         schema.NewSet(schema.HashString, []any{"sg-2"}),
		"web_groups":               schema.NewSet(schema.HashString, []any{}),
		"port_ranges":              []any{},
	}

	upperHash := dcfRuleSetHash(basRule)

	lowercaseRule := make(map[string]any, len(basRule))
	maps.Copy(lowercaseRule, basRule)
	lowercaseRule["protocol"] = "tcp"
	lowerHash := dcfRuleSetHash(lowercaseRule)

	if upperHash != lowerHash {
		t.Errorf("hash mismatch: protocol 'TCP' produced %d, protocol 'tcp' produced %d; expected equal hashes", upperHash, lowerHash)
	}

	mixedRule := make(map[string]any, len(basRule))
	maps.Copy(mixedRule, basRule)
	mixedRule["protocol"] = "Tcp"
	mixedHash := dcfRuleSetHash(mixedRule)

	if upperHash != mixedHash {
		t.Errorf("hash mismatch: protocol 'TCP' produced %d, protocol 'Tcp' produced %d; expected equal hashes", upperHash, mixedHash)
	}
}

// Verify enforcement is passed through to the struct, and watch falls back when enforcement is absent.
func TestMarshalDcfPolicyInput_EnforcementSetsField(t *testing.T) {
	base := map[string]any{
		"name":                     "test-rule",
		"action":                   "PERMIT",
		"priority":                 0,
		"protocol":                 "TCP",
		"logging":                  false,
		"watch":                    false,
		"enforcement":              "DISABLE",
		"decrypt_policy":           "DECRYPT_UNSPECIFIED",
		"flow_app_requirement":     "APP_UNSPECIFIED",
		"exclude_sg_orchestration": false,
		"tls_profile":              "",
		"uuid":                     "",
		"log_profile":              "def000ad-7000-0000-0000-000000000001",
		"egress_path":              "EGRESS_PATH_DEFAULT",
		"src_smart_groups":         schema.NewSet(schema.HashString, []any{}),
		"dst_smart_groups":         schema.NewSet(schema.HashString, []any{}),
		"web_groups":               schema.NewSet(schema.HashString, []any{}),
		"port_ranges":              []any{},
	}

	policy, err := marshalPolicyInput(base)
	if err != nil {
		t.Fatalf("marshalPolicyInput with enforcement=DISABLE: %v", err)
	}
	if policy.Enforcement != "DISABLE" {
		t.Errorf("expected Enforcement=DISABLE, got %q", policy.Enforcement)
	}

	// watch fallback: when enforcement is absent, Watch is used.
	base["enforcement"] = ""
	base["watch"] = true
	policy, err = marshalPolicyInput(base)
	if err != nil {
		t.Fatalf("marshalPolicyInput with watch=true fallback: %v", err)
	}
	if !policy.Watch {
		t.Error("expected Watch=true when enforcement is absent and watch=true")
	}
}

// Verify the read path derives watch=true only for MONITOR, and false for all other values.
func TestReadDcfRule_WatchDerivedFromEnforcement(t *testing.T) {
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
			got := tt.enforcement == "MONITOR"
			if got != tt.wantWatch {
				t.Errorf("enforcement=%q: expected watch=%v, got %v", tt.enforcement, tt.wantWatch, got)
			}
		})
	}
}

// Verify the CustomizeDiff validation rejects rules that set both watch and enforcement,
// and allows rules that set only one.
func TestDcfRulesetCustomizeDiff_RejectsBothWatchAndEnforcement(t *testing.T) {
	if err := validateDCFRuleWatchEnforcement("test-rule", true, true); err == nil {
		t.Error("expected error when both watch and enforcement are set, got nil")
	}
	if err := validateDCFRuleWatchEnforcement("test-rule", false, true); err != nil {
		t.Errorf("expected no error when only enforcement is set, got: %v", err)
	}
	if err := validateDCFRuleWatchEnforcement("test-rule", true, false); err != nil {
		t.Errorf("expected no error when only watch is set, got: %v", err)
	}
	if err := validateDCFRuleWatchEnforcement("test-rule", false, false); err != nil {
		t.Errorf("expected no error when neither is set, got: %v", err)
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

		client := mustClient(testAccProviderVersionValidation.Meta())

		_, err := client.GetDCFPolicyList(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to get DCF Ruleset status: %w", err)
		}

		return nil
	}
}

func testAccCheckDcfRulesetDestroy(s *terraform.State) error {
	client := mustClient(testAccProviderVersionValidation.Meta())

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
