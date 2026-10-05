package aviatrix

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func requireEnv(t *testing.T, key string) string {
	t.Helper()
	v := os.Getenv(key)
	if v == "" {
		t.Skipf("skipping: %s not set", key)
	}
	return v
}

func TestAccRbacGroupAccessAccountMembership(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("skipping acceptance test; TF_ACC not set")
	}

	if os.Getenv("SKIP_RBAC_GROUP_ACCESS_ACCOUNT_MEMBERSHIP") == "yes" {
		t.Skip("skipping by env flag")
	}

	// Provider-level precheck (controller creds, etc.)
	testAccPreCheck(t)

	// Pull cloud creds from environment (never commit these)
	azSub := requireEnv(t, "AZURE_SUBSCRIPTION_ID")
	azTenant := requireEnv(t, "AZURE_TENANT_ID")
	azAppID := requireEnv(t, "AZURE_APP_ID")
	azAppSecret := requireEnv(t, "AZURE_APP_SECRET")

	awsAcct := requireEnv(t, "AWS_ACCOUNT_NUMBER")
	awsAK := requireEnv(t, "AWS_ACCESS_KEY_ID")
	awsSK := requireEnv(t, "AWS_SECRET_ACCESS_KEY")

	rName := acctest.RandString(6)
	resourceName := "aviatrix_rbac_group_access_account_membership.test"

	// Attach both; Step 2: keep only Azure
	cfgStep1 := testAccCfg(rName,
		azSub, azTenant, azAppID, azAppSecret,
		awsAcct, awsAK, awsSK,
		[]string{"azure1", "aws1"}, // both
		true,
	)
	cfgStep2 := testAccCfg(rName,
		azSub, azTenant, azAppID, azAppSecret,
		awsAcct, awsAK, awsSK,
		[]string{"azure1"}, // drop AWS
		true,
	)

	var got []string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() {}, // Prechecks done above
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRbacGroupAccessAccountMembershipDestroy([]string{"azure1"}),
		Steps: []resource.TestStep{
			{
				Config: cfgStep1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRbacGroupAccessAccountMembershipExists(resourceName, &got),
					resource.TestCheckResourceAttr(resourceName, "group_name", fmt.Sprintf("tf-%s", rName)),
					testCheckStringSet(resourceName, "access_account_names", []string{"azure1", "aws1"}),
				),
			},
			{
				Config: cfgStep2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRbacGroupAccessAccountMembershipExists(resourceName, &got),
					testCheckStringSet(resourceName, "access_account_names", []string{"azure1"}),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_access_accounts_on_destroy"},
			},
		},
	})
}

// Renders HCL for one RBAC group, two cloud accounts (Azure+AWS), and the membership resource.
// All secrets are passed via args (from env), not hardcoded.
func testAccCfg(
	rName, azSub, azTenant, azAppID, azAppSecret,
	awsAcct, awsAK, awsSK string,
	memberNames []string,
	removeOnDestroy bool,
) string {
	var memberRefs strings.Builder
	for _, name := range memberNames {
		switch name {
		case "azure1":
			fmt.Fprintf(&memberRefs, "    aviatrix_account.azure_account.account_name,\n")
		case "aws1":
			fmt.Fprintf(&memberRefs, "    aviatrix_account.aws_account.account_name,\n")
		default:
			panic("unknown member name " + name)
		}
	}

	return fmt.Sprintf(`
resource "aviatrix_rbac_group" "test" {
  group_name = "tf-%[1]s"
}

resource "aviatrix_account" "azure_account" {
  account_name        = "azure1"
  cloud_type          = 8
  arm_subscription_id = "%[2]s"
  arm_directory_id    = "%[3]s"
  arm_application_id  = "%[4]s"
  arm_application_key = "%[5]s"
}

resource "aviatrix_account" "aws_account" {
  account_name       = "aws1"
  cloud_type         = 1
  aws_account_number = "%[6]s"
  aws_iam            = false
  aws_access_key     = "%[7]s"
  aws_secret_key     = "%[8]s"
}

resource "aviatrix_rbac_group_access_account_membership" "test" {
  group_name = aviatrix_rbac_group.test.group_name
  access_account_names = [
%[9]s  ]
  remove_access_accounts_on_destroy = %[10]t
}
`, rName, azSub, azTenant, azAppID, azAppSecret, awsAcct, awsAK, awsSK, memberRefs.String(), removeOnDestroy)
}

// Reads current membership from backend and records it.
func testAccCheckRbacGroupAccessAccountMembershipExists(n string, got *[]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("rbac group access account membership not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no rbac group access account membership ID set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)
		group := rs.Primary.Attributes["group_name"]
		if group == "" {
			group = rs.Primary.ID
		}

		current, err := client.ListRbacGroupAccessAccounts(group)
		if err != nil {
			if errors.Is(err, goaviatrix.ErrNotFound) {
				return fmt.Errorf("rbac group %q not found in backend", group)
			}
			return err
		}
		sort.Strings(current)
		*got = current
		return nil
	}
}

// On destroy, ensure none of the applied access accounts remain members.
// If the group itself is gone, that's acceptable.
func testAccCheckRbacGroupAccessAccountMembershipDestroy(lastApplied []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*goaviatrix.Client)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aviatrix_rbac_group_access_account_membership" {
				continue
			}
			group := rs.Primary.Attributes["group_name"]
			if group == "" {
				group = rs.Primary.ID
			}

			current, err := client.ListRbacGroupAccessAccounts(group)
			if err != nil {
				// Group missing is acceptable
				lo := strings.ToLower(err.Error())
				if errors.Is(err, goaviatrix.ErrNotFound) || strings.Contains(lo, "not found") {
					return nil
				}
				return err
			}

			for _, a := range lastApplied {
				if slices.Contains(current, a) {
					return fmt.Errorf("rbac group %q still contains access account %q after destroy", group, a)
				}
			}
		}
		return nil
	}
}
