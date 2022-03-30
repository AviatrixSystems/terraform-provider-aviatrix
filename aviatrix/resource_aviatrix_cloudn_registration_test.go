package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAviatrixCloudnRegistration_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_CLOUDN_REGISTRATION")
	if skipAcc == "yes" {
		t.Skip("Skipping Aviatrix CloudN Registration test as SKIP_CLOUDN_REGISTRATION is set")
	}

	rName := fmt.Sprintf("cloudn-%s", acctest.RandString(5))
	resourceName := "aviatrix_cloudn_registration.test_cloudn_registration"
	localASNumber := "65707"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccAviatrixCloudnRegistrationPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccCheckCloudnRegistrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudnRegistrationBasic(rName, localASNumber),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudnRegistrationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address", os.Getenv("CLOUDN_IP")),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "local_as_number", localASNumber),
					resource.TestCheckNoResourceAttr(resourceName, "prepend_as_path"),
				),
			},
			{
				Config: testAccCloudnRegistrationBasicUpdated(rName, localASNumber),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudnRegistrationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address", os.Getenv("CLOUDN_IP")),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "local_as_number", localASNumber),
					resource.TestCheckResourceAttr(resourceName, "prepend_as_path.0", localASNumber),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"username",
					"password",
				},
			},
		},
	})
}

func testAccCloudnRegistrationBasic(rName, localASNumber string) string {
	return fmt.Sprintf(`
resource "aviatrix_cloudn_registration" "test_cloudn_registration" {
	name			= "%[1]s"
	address			= "%[2]s"
	username		= "%[3]s"
	password		= "%[4]s"
	local_as_number		= "%[5]s"
}
`, rName, os.Getenv("CLOUDN_IP"), os.Getenv("CLOUDN_USERNAME"), os.Getenv("CLOUDN_PASSWORD"),
		localASNumber)
}

func testAccCloudnRegistrationBasicUpdated(rName, localASNumber string) string {
	return fmt.Sprintf(`
resource "aviatrix_cloudn_registration" "test_cloudn_registration" {
	name			= "%[1]s"
	address			= "%[2]s"
	username		= "%[3]s"
	password		= "%[4]s"
	local_as_number		= "%[5]s"
	prepend_as_path		= ["%[5]s"]
}
`, rName, os.Getenv("CLOUDN_IP"), os.Getenv("CLOUDN_USERNAME"), os.Getenv("CLOUDN_PASSWORD"),
		localASNumber)
}

func testAccCheckCloudnRegistrationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("aviatrix_cloudn_registration resource not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("aviatrix_cloudn_registration ID is empty")
		}

		client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

		cloudnRegistration := &goaviatrix.CloudnRegistration{
			Name: rs.Primary.Attributes["name"],
		}

		_, err := client.GetCloudnRegistration(context.Background(), cloudnRegistration)
		if err != nil {
			return err
		}
		if cloudnRegistration.Name != rs.Primary.ID {
			return fmt.Errorf("aviatrix_cloudn_registration not found")
		}

		return nil
	}
}

func testAccCheckCloudnRegistrationDestroy(s *terraform.State) error {
	client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_cloudn_registration" {
			continue
		}
		cloudnRegistration := &goaviatrix.CloudnRegistration{
			Name: rs.Primary.Attributes["name"],
		}
		_, err := client.GetCloudnRegistration(context.Background(), cloudnRegistration)
		if err == nil {
			return fmt.Errorf("aviatrix_cloudn_registration still exists")
		}
	}

	return nil
}

func testAccAviatrixCloudnRegistrationPreCheck(t *testing.T) {
	requiredEnv := []string{
		"CLOUDN_IP",
		"CLOUDN_USERNAME",
		"CLOUDN_PASSWORD",
	}

	for _, v := range requiredEnv {
		if os.Getenv(v) == "" {
			t.Fatalf("%s must be set for aviatrix_cloudn_registration acceptance test", v)
		}
	}
}
