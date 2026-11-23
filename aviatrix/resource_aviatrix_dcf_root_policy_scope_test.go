package aviatrix

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

// TestDCFRootPolicyScopeSchema verifies the singleton resource schema.
func TestDCFRootPolicyScopeSchema(t *testing.T) {
	rootScope := resourceAviatrixDCFRootPolicyScope()

	require.NotNil(t, rootScope.Importer)
	scopeUUIDs := rootScope.Schema["scope_uuids"]
	assert.True(t, scopeUUIDs.Required)
	assert.Equal(t, schema.TypeSet, scopeUUIDs.Type)
	assert.Equal(t, 1, scopeUUIDs.MinItems)

	elem, ok := scopeUUIDs.Elem.(*schema.Schema)
	require.True(t, ok, "scope_uuids Elem should be *schema.Schema")
	validator := elem.ValidateFunc
	_, validErrors := validator("123e4567-e89b-12d3-a456-426614174000", "scope_uuids")
	assert.Empty(t, validErrors)
	_, invalidErrors := validator("not-a-uuid", "scope_uuids")
	assert.NotEmpty(t, invalidErrors)
}

// TestDCFRootPolicyScopeReadRejectsWrongImportID verifies the semantic singleton import ID.
func TestDCFRootPolicyScopeReadRejectsWrongImportID(t *testing.T) {
	data := schema.TestResourceDataRaw(t, resourceAviatrixDCFRootPolicyScope().Schema, nil)
	data.SetId("wrong")

	diagnostics := resourceAviatrixDCFRootPolicyScopeRead(context.Background(), data, nil)

	require.Len(t, diagnostics, 1)
	assert.Contains(t, diagnostics[0].Summary, dcfRootPolicyScopeID)
}

// TestDCFRootPolicyScopeDeleteResetsGlobalScope verifies destroy restores the required default scope.
func TestDCFRootPolicyScopeDeleteResetsGlobalScope(t *testing.T) {
	var body goaviatrix.DCFPolicyScopeConfigPatch
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "/v2.5/api/microseg/policy-list3/"+dcfSystemRootPolicyBlockUUID, r.URL.Path)
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		_, _ = w.Write([]byte(`{
			"name":"System Root Policy Block",
			"uuid":"` + dcfSystemRootPolicyBlockUUID + `",
			"system_resource":true,
			"sub_policies":[],
			"scope_config":{"scope_uuids":["` + dcfDefaultGlobalPolicyScopeID + `"]}
		}`))
	}))
	t.Cleanup(server.Close)
	client := &goaviatrix.Client{
		HTTPClient:   server.Client(),
		CID:          "test-cid",
		ControllerIP: strings.TrimPrefix(server.URL, "https://"),
	}
	data := schema.TestResourceDataRaw(t, resourceAviatrixDCFRootPolicyScope().Schema, nil)
	data.SetId(dcfRootPolicyScopeID)

	diagnostics := resourceAviatrixDCFRootPolicyScopeDelete(context.Background(), data, client)

	assert.Empty(t, diagnostics)
	assert.Equal(t, []string{dcfDefaultGlobalPolicyScopeID}, body.ScopeConfig.ScopeUUIDs)
	assert.Empty(t, data.Id())
}

// TestAccAviatrixDCFRootPolicyScope covers PATCH create/update, no-diff refresh, import, and destroy reset.
func TestAccAviatrixDCFRootPolicyScope(t *testing.T) {
	const resourceName = "aviatrix_dcf_root_policy_scope.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"aviatrix": func() (*schema.Provider, error) { return testAccProviderVersionValidation, nil },
		},
		CheckDestroy: testAccCheckDCFRootPolicyScopeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDCFRootPolicyScopeConfig(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "scope_uuids.#", "1"),
				),
			},
			{
				Config: testAccDCFRootPolicyScopeConfig(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "scope_uuids.#", "2"),
				),
			},
			{
				Config:             testAccDCFRootPolicyScopeConfig(true),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// testAccDCFRootPolicyScopeConfig returns acceptance configuration with one or two root scopes.
func testAccDCFRootPolicyScopeConfig(includeSecond bool) string {
	rootScopes := "aviatrix_dcf_policy_scope.first.uuid"
	if includeSecond {
		rootScopes += ", aviatrix_dcf_policy_scope.second.uuid"
	}
	return fmt.Sprintf(`
resource "aviatrix_dcf_policy_scope" "first" {
  name = "terraform-root-policy-scope-first"
  included {
    preset = ["all"]
  }
}

resource "aviatrix_dcf_policy_scope" "second" {
  name = "terraform-root-policy-scope-second"
  included {
    preset = ["none"]
  }
}

resource "aviatrix_dcf_root_policy_scope" "test" {
  scope_uuids = [%s]
}
`, rootScopes)
}

// testAccCheckDCFRootPolicyScopeDestroy verifies destroy restores the Global Scope.
func testAccCheckDCFRootPolicyScopeDestroy(_ *terraform.State) error {
	client := mustClient(testAccProviderVersionValidation.Meta())
	policyBlock, err := client.GetDCFPolicyBlock(context.Background(), dcfSystemRootPolicyBlockUUID)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("system root policy block was removed")
		}
		return fmt.Errorf("reading system root policy block after destroy: %w", err)
	}
	if policyBlock.ScopeConfig == nil ||
		len(policyBlock.ScopeConfig.ScopeUUIDs) != 1 ||
		policyBlock.ScopeConfig.ScopeUUIDs[0] != dcfDefaultGlobalPolicyScopeID {
		return fmt.Errorf("system root policy block was not reset to the Global Scope")
	}
	return nil
}
