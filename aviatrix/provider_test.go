package aviatrix

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

var testAccProvidersVersionValidation map[string]*schema.Provider
var testAccProviderVersionValidation *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"aviatrix": testAccProvider,
	}

	testAccProviderVersionValidation = Provider()
	testAccProviderVersionValidation.ConfigureFunc = aviatrixConfigureWithoutVersionValidation
	testAccProvidersVersionValidation = map[string]*schema.Provider{
		"aviatrix": testAccProviderVersionValidation,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("AVIATRIX_CONTROLLER_IP"); v == "" {
		t.Fatal("AVIATRIX_CONTROLLER_IP must be set for acceptance tests.")
	}
	if v := os.Getenv("AVIATRIX_USERNAME"); v == "" {
		t.Fatal("AVIATRIX_USERNAME must be set for acceptance tests.")
	}
	if v := os.Getenv("AVIATRIX_PASSWORD"); v == "" {
		t.Fatal("AVIATRIX_PASSWORD must be set for acceptance tests.")
	}
}

func TestCIDTimeout(t *testing.T) {
	if os.Getenv("SKIP_CID_EXPIRY") == "yes" {
		t.Skip("Skipping CID expiry retry test as SKIP_CID_EXPIRY is set")
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProvidersVersionValidation,
		Steps: []resource.TestStep{
			{
				Config: `
resource "aviatrix_rbac_group" "a" {
	group_name = "test-group"
}
`,
				Check: testCID,
			},
		},
	})
}

func testCID(s *terraform.State) error {
	client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

	log.Printf("found CID: %s", client.CID)
	time.Sleep(time.Hour + 30*time.Minute)

	group := &goaviatrix.RbacGroup{
		GroupName: "test-group",
	}
	rGroup, err := client.GetPermissionGroup(group)
	if err != nil {
		return fmt.Errorf("CID test failed with: %v", err)
	}

	fmt.Printf("Found user: %v", rGroup.GroupName)
	return nil
}
