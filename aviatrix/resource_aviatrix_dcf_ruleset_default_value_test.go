package aviatrix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// TestDcfRuleset_OptionalFieldsHaveDefaults verifies that all optional string fields
// in the rules schema have Default: "" set.
//
// WHY THIS MATTERS:
// -----------------
// The "rules" attribute uses schema.TypeSet, which identifies elements by hash.
// When ANY field differs between state and config, Terraform sees it as a different
// element (remove old + add new).
//
// PROBLEM SCENARIO (without Default):
// -----------------------------------
// 1. User creates rules without specifying optional field "some_field"
// 2. Backend doesn't return "some_field" (field missing from API response)
// 3. Read function sets some_field = "" in state
// 4. User updates only rule1's name
// 5. Without Default: "" in schema:
//   - Config interprets unspecified field as nil
//   - State has "" (from Read)
//   - nil != "" → Hash mismatch for ALL rules → ALL rules show in diff (bad!)
//
// SOLUTION:
// ---------
// Add Default: "" to all optional string fields. Then:
//   - Config uses "" for unspecified field (from Default)
//   - State has "" (from Read)
//   - "" == "" → Only actually changed rules show in diff (good!)
//
// This test ensures no one forgets to add Default: "" when adding new optional fields.
func TestDcfRuleset_OptionalFieldsHaveDefaults(t *testing.T) {
	resourceSchema := resourceAviatrixDCFRuleset()
	rulesSchema := resourceSchema.Schema["rules"]
	ruleSchema, ok := rulesSchema.Elem.(*schema.Resource)
	assert.True(t, ok)
	assert.NotNil(t, ruleSchema)

	for fieldName, fieldSchema := range ruleSchema.Schema {
		// Only check optional string fields (Required fields are always specified by user)
		if fieldSchema.Type != schema.TypeString {
			continue
		}
		if fieldSchema.Required {
			continue
		}
		if fieldSchema.Computed && !fieldSchema.Optional {
			// Computed-only fields are set by the provider, not the user
			continue
		}

		// Optional string field - must have Default: "" to avoid spurious diffs
		assert.NotNilf(t, fieldSchema.Default, "Field %s is missing Default: \"\"", fieldName)
	}
}
