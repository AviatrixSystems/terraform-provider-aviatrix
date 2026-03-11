package aviatrix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAllResources_OptionalFieldsInSetHaveDefaults verifies that all optional
// fields in TypeSet blocks with nested Resource schemas have Default: ""
// set across ALL provider resources.
//
// WHY THIS MATTERS:
// -----------------
// TypeSet identifies elements by hash. When ANY field differs between state and config,
// Terraform sees it as a different element (remove old + add new).
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
// Add Default: "" to all optional fields. Then:
//   - Config uses "" for unspecified field (from Default)
//   - State has "" (from Read)
//   - "" == "" → Only actually changed rules show in diff (good!)
//
// This test ensures no one forgets to add Default: "" when adding new optional fields.
func TestAllResources_OptionalFieldsInSetHaveDefaults(t *testing.T) {
	provider := Provider()

	for resourceName, resource := range provider.ResourcesMap {
		require.NotEmpty(t, resource.Schema)
		checkResourceSchema(t, resourceName, "", resource.Schema)
	}
}

// checkResourceSchema recursively walks schema attributes, finding TypeSet blocks
// with nested Resource schemas and verifying optional fields have Default set.
func checkResourceSchema(t *testing.T, resourceName, path string, schemas map[string]*schema.Schema) {
	for attrName, attrSchema := range schemas {
		attrPath := path
		if attrPath != "" {
			attrPath += "."
		}
		attrPath += attrName

		if attrSchema.Type != schema.TypeSet {
			continue
		}
		elemResource, ok := attrSchema.Elem.(*schema.Resource)
		if !ok {
			continue
		}
		require.NotNil(t, elemResource)
		require.NotEmpty(t, elemResource.Schema)

		for fieldName, fieldSchema := range elemResource.Schema {
			fieldPath := attrPath + "." + fieldName
			// Check optional fields
			switch fieldSchema.Type {
			case schema.TypeSet:
				if nestedResource, ok := fieldSchema.Elem.(*schema.Resource); ok {
					require.NotEmpty(t, nestedResource.Schema, fieldPath)
					checkResourceSchema(t, resourceName, fieldPath, nestedResource.Schema)
				}
			case schema.TypeMap, schema.TypeList:
				// For TypeMap, this terraform plugin sdk only supports map[string]string, so we do not need to check the default value
				// For TypeList, terraform compares based on index, and not the has so we don't need this check
				continue
			case schema.TypeString, schema.TypeInt, schema.TypeFloat, schema.TypeBool:
				if fieldSchema.Required {
					continue
				}
				if fieldSchema.Computed && !fieldSchema.Optional {
					continue
				}
				assert.NotNilf(t, fieldSchema.Default, "%s.%s: field %q is optional but missing Default, please add a Default value to the field", resourceName, fieldPath, fieldName)
				continue
			default:
				t.Errorf("unsupported field type: %s", fieldSchema.Type)
			}
		}
	}
}
