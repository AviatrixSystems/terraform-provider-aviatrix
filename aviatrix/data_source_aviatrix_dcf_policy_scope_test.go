package aviatrix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// TestDataSourceDCFPolicyScopeSchema verifies name is the only lookup argument.
func TestDataSourceDCFPolicyScopeSchema(t *testing.T) {
	dataSource := dataSourceAviatrixDCFPolicyScope()

	assert.True(t, dataSource.Schema["name"].Required)
	assert.True(t, dataSource.Schema["uuid"].Computed)
	assert.False(t, dataSource.Schema["uuid"].Optional)
	assert.True(t, dataSource.Schema["included"].Computed)
	assert.True(t, dataSource.Schema["excluded"].Computed)
}

// TestAccDataSourceAviatrixDCFPolicyScope_basic looks up the Global Scope and a Terraform-managed scope.
func TestAccDataSourceAviatrixDCFPolicyScope_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"aviatrix": func() (*schema.Provider, error) { return testAccProviderVersionValidation, nil },
		},
		Steps: []resource.TestStep{
			{
				Config: `
resource "aviatrix_dcf_policy_scope" "test" {
  name = "terraform-data-source-scope"
  included {
    vpc_ids = ["vpc-terraform-test"]
  }
}

data "aviatrix_dcf_policy_scope" "global" {
  name = "Global Scope"
}

data "aviatrix_dcf_policy_scope" "managed" {
  name = aviatrix_dcf_policy_scope.test.name
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aviatrix_dcf_policy_scope.global", "uuid", dcfDefaultGlobalPolicyScopeID),
					resource.TestCheckTypeSetElemAttr("data.aviatrix_dcf_policy_scope.global", "included.0.preset.*", "all"),
					resource.TestCheckTypeSetElemAttr("data.aviatrix_dcf_policy_scope.global", "excluded.0.preset.*", "none"),
					resource.TestCheckResourceAttrPair("data.aviatrix_dcf_policy_scope.managed", "uuid", "aviatrix_dcf_policy_scope.test", "uuid"),
					resource.TestCheckTypeSetElemAttr("data.aviatrix_dcf_policy_scope.managed", "included.0.vpc_ids.*", "vpc-terraform-test"),
				),
			},
		},
	})
}
