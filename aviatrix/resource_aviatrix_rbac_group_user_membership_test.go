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

func TestAccAviatrixRbacGroupUserMembership_basic(t *testing.T) {
	var gotUsers []string

	rName := acctest.RandString(5)

	skipAcc := os.Getenv("SKIP_RBAC_GROUP_USER_MEMBERSHIP")
	if skipAcc == "yes" {
		t.Skip("Skipping rbac group user membership tests as SKIP_RBAC_GROUP_USER_MEMBERSHIP is set")
	}

	resourceName := "aviatrix_rbac_group_user_membership.test"
	msgCommon := ". Set SKIP_RBAC_GROUP_USER_MEMBERSHIP to 'yes' to skip rbac group user membership tests"

	usersStep1 := []string{
		fmt.Sprintf("tf-user-a-%s", rName),
		fmt.Sprintf("tf-user-b-%s", rName),
		fmt.Sprintf("tf-user-c-%s", rName),
	}
	usersStep2 := []string{
		fmt.Sprintf("tf-user-a-%s", rName),
		fmt.Sprintf("tf-user-c-%s", rName),
		fmt.Sprintf("tf-user-d-%s", rName),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRbacGroupUserMembershipDestroy(usersStep2),
		Steps: []resource.TestStep{
			{
				Config: testAccRbacGroupUserMembershipConfig(rName, usersStep1, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRbacGroupUserMembershipExists(resourceName, &gotUsers),
					resource.TestCheckResourceAttr(resourceName, "group_name", fmt.Sprintf("tf-%s", rName)),
					testCheckStringSet(resourceName, "user_names", usersStep1),
				),
			},
			{
				Config: testAccRbacGroupUserMembershipConfig(rName, usersStep2, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRbacGroupUserMembershipExists(resourceName, &gotUsers),
					testCheckStringSet(resourceName, "user_names", usersStep2),
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

func testAccRbacGroupUserMembershipConfig(rName string, users []string, removeOnDestroy bool) string {
	group := fmt.Sprintf("tf-%s", rName)

	var userBlocks strings.Builder
	for i, u := range users {
		fmt.Fprintf(&userBlocks, `
resource "aviatrix_account_user" "u%d" {
  username = "%s"
  email    = "%s@xyz.com"
  password = "Password-1234"
}
`, i+1, u, u)
	}

	var userRefs strings.Builder
	for i := range users {
		fmt.Fprintf(&userRefs, "    aviatrix_account_user.u%d.username,\n", i+1)
	}

	return fmt.Sprintf(`
resource "aviatrix_rbac_group" "test" {
  group_name = "%s"
}

%s

resource "aviatrix_rbac_group_user_membership" "test" {
  group_name = aviatrix_rbac_group.test.group_name
  user_names = [
%s  ]
  remove_users_on_destroy = %t
}
`, group, userBlocks.String(), userRefs.String(), removeOnDestroy)
}

func testAccCheckRbacGroupUserMembershipExists(n string, got *[]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("rbac group user membership not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no rbac group user membership ID set")
		}

		client := mustClient(testAccProvider.Meta())
		group := rs.Primary.Attributes["group_name"]
		if group == "" {
			group = rs.Primary.ID
		}

		current, err := client.ListRbacGroupUsers(group)
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

// For destroy, because the resource may only remove users if remove_users_on_destroy=true,
// we verify that none of the test users remain members. If the group is gone, that's fine.
func testAccCheckRbacGroupUserMembershipDestroy(lastAppliedUsers []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := mustClient(testAccProvider.Meta())

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aviatrix_rbac_group_user_membership" {
				continue
			}
			group := rs.Primary.Attributes["group_name"]
			if group == "" {
				group = rs.Primary.ID
			}

			current, err := client.ListRbacGroupUsers(group)
			if err != nil {
				// Group missing is acceptable (treated as removed from state)
				if errors.Is(err, goaviatrix.ErrNotFound) || strings.Contains(strings.ToLower(err.Error()), "not found") {
					return nil
				}
				return err
			}

			// Ensure none of the last-applied users are still members
			for _, u := range lastAppliedUsers {
				if slices.Contains(current, u) {
					return fmt.Errorf("rbac group %q still contains user %q after destroy", group, u)
				}
			}
		}
		return nil
	}
}
