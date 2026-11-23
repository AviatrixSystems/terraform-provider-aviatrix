package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

// TestDCFPolicyScopeSchema verifies the resource's public schema contract.
func TestDCFPolicyScopeSchema(t *testing.T) {
	resource := resourceAviatrixDCFPolicyScope()

	require.NotNil(t, resource.Importer)
	assert.True(t, resource.Schema["name"].Required)
	assert.True(t, resource.Schema["uuid"].Computed)
	assert.True(t, resource.Schema["included"].Required)
	assert.Equal(t, 1, resource.Schema["included"].MaxItems)
	assert.True(t, resource.Schema["excluded"].Optional)
	assert.True(t, resource.Schema["excluded"].Computed)

	included, ok := resource.Schema["included"].Elem.(*schema.Resource)
	require.True(t, ok, "included Elem should be *schema.Resource")
	assert.ElementsMatch(t, dcfPolicyScopeSelectorKeys, selectorSchemaKeys(included.Schema))
	assert.Equal(t, 1, included.Schema["preset"].MaxItems)
}

// TestDCFPolicyScopeSchemaSelectorRequirements verifies custom validation for selector blocks.
func TestDCFPolicyScopeSchemaSelectorRequirements(t *testing.T) {
	resource := resourceAviatrixDCFPolicyScope()

	validConfig := terraform.NewResourceConfigRaw(map[string]any{
		"name": "all",
		"included": []any{map[string]any{
			"preset": []any{"all"},
		}},
	})
	_, err := resource.Diff(context.Background(), nil, validConfig, nil)
	require.NoError(t, err)

	emptyConfig := terraform.NewResourceConfigRaw(map[string]any{
		"name":     "empty",
		"included": []any{map[string]any{}},
	})
	_, err = resource.Diff(context.Background(), nil, emptyConfig, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "included")

	unknownConfig := terraform.NewResourceConfigRaw(map[string]any{
		"name": "unknown",
		"included": []any{map[string]any{
			"preset": []any{"all"},
		}},
	})
	unknownConfig.ComputedKeys = []string{"included"}
	_, err = resource.Diff(context.Background(), nil, unknownConfig, nil)
	require.NoError(t, err)
}

// TestExpandDCFPolicyScopeSelectors verifies expansion of populated selector sets.
func TestExpandDCFPolicyScopeSelectors(t *testing.T) {
	value := []any{map[string]any{
		"csps":    schema.NewSet(schema.HashString, []any{"azure", "aws"}),
		"vpc_ids": schema.NewSet(schema.HashString, []any{"vpc-1"}),
		"preset":  schema.NewSet(schema.HashString, nil),
	}}

	selectors := expandDCFPolicyScopeSelectors(value)

	require.Len(t, selectors.Selectors, 2)
	assert.ElementsMatch(t, []string{"azure", "aws"}, selectors.Selectors[0]["csps"])
	assert.Equal(t, []string{"vpc-1"}, selectors.Selectors[1]["vpc_ids"])
}

// TestFlattenDCFPolicyScopeSelectorsMergesCanonicalKeys verifies repeated API keys produce stable state.
func TestFlattenDCFPolicyScopeSelectorsMergesCanonicalKeys(t *testing.T) {
	flattened := flattenDCFPolicyScopeSelectors(goaviatrix.DCFPolicyScopeSelectors{
		Selectors: []goaviatrix.DCFPolicyScopeSelector{
			{"vpc_ids": {"vpc-1"}},
			{"csps": {"aws"}},
			{"vpc_ids": {"vpc-2"}},
		},
	})

	require.Len(t, flattened, 1)
	block, ok := flattened[0].(map[string]any)
	require.True(t, ok, "flattened block should be map[string]any")
	assert.Equal(t, []string{"aws"}, block["csps"])
	assert.Equal(t, []string{"vpc-1", "vpc-2"}, block["vpc_ids"])
}

// TestDCFPolicyScopePresetValidation verifies the preset enum contract.
func TestDCFPolicyScopePresetValidation(t *testing.T) {
	resource := resourceAviatrixDCFPolicyScope()
	selectors, ok := resource.Schema["included"].Elem.(*schema.Resource)
	require.True(t, ok, "included Elem should be *schema.Resource")
	preset, ok := selectors.Schema["preset"].Elem.(*schema.Schema)
	require.True(t, ok, "preset Elem should be *schema.Schema")

	_, validErrors := preset.ValidateFunc("all", "included.0.preset")
	assert.Empty(t, validErrors)
	_, invalidErrors := preset.ValidateFunc("everything", "included.0.preset")
	assert.NotEmpty(t, invalidErrors)
}

// TestDCFPolicyScopeCSPValidation verifies csps only accepts controller-supported cloud providers.
func TestDCFPolicyScopeCSPValidation(t *testing.T) {
	resource := resourceAviatrixDCFPolicyScope()
	selectors, ok := resource.Schema["included"].Elem.(*schema.Resource)
	require.True(t, ok, "included Elem should be *schema.Resource")
	csps, ok := selectors.Schema["csps"].Elem.(*schema.Schema)
	require.True(t, ok, "csps Elem should be *schema.Schema")

	for _, value := range dcfPolicyScopeCSPValues {
		_, errs := csps.ValidateFunc(value, "included.0.csps")
		assert.Empty(t, errs, value)
	}
	for _, value := range []string{"aws", "Azure", "AVX-TEST", "gcp-gov"} {
		_, errs := csps.ValidateFunc(value, "included.0.csps")
		assert.NotEmpty(t, errs, value)
	}
}

// TestMarshalDCFPolicyScopeInputPreservesExcludedOmission verifies omitted exclusions remain absent.
func TestMarshalDCFPolicyScopeInputPreservesExcludedOmission(t *testing.T) {
	resourceSchema := resourceAviatrixDCFPolicyScope().Schema
	data := schema.TestResourceDataRaw(t, resourceSchema, map[string]any{
		"name": "all",
		"included": []any{map[string]any{
			"preset": []any{"all"},
		}},
	})

	input := marshalDCFPolicyScopeInput(data)

	assert.Nil(t, input.Excluded)
}

// TestMarshalDCFPolicyScopeInputIncludesConfiguredExcluded verifies explicit exclusions reach the API request.
func TestMarshalDCFPolicyScopeInputIncludesConfiguredExcluded(t *testing.T) {
	resourceSchema := resourceAviatrixDCFPolicyScope().Schema
	data := schema.TestResourceDataRaw(t, resourceSchema, map[string]any{
		"name": "all-except-maintenance",
		"included": []any{map[string]any{
			"preset": []any{"all"},
		}},
		"excluded": []any{map[string]any{
			"gateway_group_names": []any{"maintenance"},
		}},
	})

	input := marshalDCFPolicyScopeInput(data)

	require.NotNil(t, input.Excluded)
	require.Len(t, input.Excluded.Selectors, 1)
	assert.Equal(t, []string{"maintenance"}, input.Excluded.Selectors[0]["gateway_group_names"])
}

// TestDCFPolicyScopeExcludedOmitted verifies when an excluded block is treated as absent from configuration.
func TestDCFPolicyScopeExcludedOmitted(t *testing.T) {
	excludedType := cty.List(cty.EmptyObject)
	rawWith := func(excluded cty.Value) cty.Value {
		return cty.ObjectVal(map[string]cty.Value{"excluded": excluded})
	}

	tests := []struct {
		name      string
		rawConfig cty.Value
		want      bool
	}{
		{name: "null", rawConfig: rawWith(cty.NullVal(excludedType)), want: true},
		{name: "empty list", rawConfig: rawWith(cty.ListValEmpty(cty.EmptyObject)), want: true},
		{name: "configured", rawConfig: rawWith(cty.ListVal([]cty.Value{cty.EmptyObjectVal})), want: false},
		{name: "unknown", rawConfig: rawWith(cty.UnknownVal(excludedType)), want: false},
		{name: "no raw config", rawConfig: cty.NullVal(cty.EmptyObject), want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, dcfPolicyScopeExcludedOmitted(tt.rawConfig))
		})
	}
}

// selectorSchemaKeys returns all keys from a selector schema.
func selectorSchemaKeys(selectorSchema map[string]*schema.Schema) []string {
	keys := make([]string, 0, len(selectorSchema))
	for key := range selectorSchema {
		keys = append(keys, key)
	}
	return keys
}

// TestAccAviatrixDCFPolicyScope_basic covers create, update, import, and destroy against a controller.
func TestAccAviatrixDCFPolicyScope_basic(t *testing.T) {
	const resourceName = "aviatrix_dcf_policy_scope.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"aviatrix": func() (*schema.Provider, error) { return testAccProviderVersionValidation, nil },
		},
		CheckDestroy: func(state *terraform.State) error {
			client := mustClient(testAccProviderVersionValidation.Meta())
			for _, stateResource := range state.RootModule().Resources {
				if stateResource.Type != "aviatrix_dcf_policy_scope" {
					continue
				}
				_, err := client.GetDCFPolicyScope(context.Background(), stateResource.Primary.ID)
				if errors.Is(err, goaviatrix.ErrNotFound) {
					continue
				}
				if err != nil {
					return fmt.Errorf("checking destroyed DCF policy scope: %w", err)
				}
				return fmt.Errorf("DCF policy scope %q still exists", stateResource.Primary.ID)
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: `
resource "aviatrix_dcf_policy_scope" "test" {
  name = "terraform-policy-scope"
  included {
    preset = ["all"]
  }
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "terraform-policy-scope"),
					resource.TestCheckResourceAttr(resourceName, "included.0.preset.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "excluded.0.preset.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
				),
			},
			{
				Config: `
resource "aviatrix_dcf_policy_scope" "test" {
  name = "terraform-policy-scope-updated"
  included {
    vpc_ids = ["vpc-terraform-test"]
  }
  excluded {
    preset = ["none"]
  }
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "terraform-policy-scope-updated"),
					resource.TestCheckResourceAttr(resourceName, "included.0.vpc_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "excluded.0.preset.#", "1"),
				),
			},
			{
				Config: `
resource "aviatrix_dcf_policy_scope" "test" {
  name = "terraform-policy-scope-updated"
  included {
    preset = ["all"]
  }
  excluded {
    vpc_ids = ["vpc-terraform-test"]
  }
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "excluded.0.vpc_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "excluded.0.preset.#", "0"),
				),
			},
			{
				// Removing excluded must reset it to the controller default instead of keeping prior state.
				Config: `
resource "aviatrix_dcf_policy_scope" "test" {
  name = "terraform-policy-scope-updated"
  included {
    preset = ["all"]
  }
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "excluded.0.vpc_ids.#", "0"),
					resource.TestCheckTypeSetElemAttr(resourceName, "excluded.0.preset.*", "none"),
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
