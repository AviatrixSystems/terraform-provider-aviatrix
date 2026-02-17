package aviatrix

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixAWSTgw_basic(t *testing.T) {
	var awsTgw goaviatrix.AWSTgw

	rName := acctest.RandString(5)
	resourceName := "aviatrix_aws_tgw.aws_tgw_test"

	skipAcc := os.Getenv("SKIP_AWS_TGW")
	if skipAcc == "yes" {
		t.Skip("Skipping AWS TGW test as SKIP_AWS_TGW is set")
	}
	msg := ". Set SKIP_AWS_TGW to yes to skip AWS TGW  tests"

	awsSideAsNumber := "64512"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, msg)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSTgwDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSTgwConfigBasic(rName, awsSideAsNumber),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSTgwExists(resourceName, &awsTgw),
					resource.TestCheckResourceAttr(resourceName, "tgw_name", fmt.Sprintf("tft-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfaa-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "region", os.Getenv("AWS_REGION")),
					resource.TestCheckResourceAttr(resourceName, "aws_side_as_number", awsSideAsNumber),
				),
			},
		},
	})
}

func testAccAWSTgwConfigBasic(rName string, awsSideAsNumber string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_aws_tgw" "aws_tgw_test" {
	account_name       = aviatrix_account.test_account.account_name
	aws_side_as_number = "%s"
	region             = "%s"
	tgw_name           = "tft-%s"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		awsSideAsNumber, os.Getenv("AWS_REGION"), rName)
}

func testAccCheckAWSTgwExists(n string, awsTgw *goaviatrix.AWSTgw) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("AWS TGW Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no AWS TGW ID is set")
		}

		client := mustClient(testAccProvider.Meta())

		foundAwsTgw := &goaviatrix.AWSTgw{
			Name: rs.Primary.Attributes["tgw_name"],
		}

		foundAwsTgw2, err := client.GetAWSTgw(foundAwsTgw)
		if err != nil {
			return err
		}
		if foundAwsTgw2.Name != rs.Primary.ID {
			return fmt.Errorf("AWS TGW not found")
		}

		*awsTgw = *foundAwsTgw
		return nil
	}
}

func testAccCheckAWSTgwDestroy(s *terraform.State) error {
	client := mustClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_aws_tgw" {
			continue
		}

		foundAWSTgw := &goaviatrix.AWSTgw{
			Name: rs.Primary.Attributes["tgw_name"],
		}

		_, err := client.GetAWSTgw(foundAWSTgw)
		if err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return nil
			}
			return fmt.Errorf("AWS TGW still exists: %w", err)
		}

		return fmt.Errorf("AWS TGW still exists")
	}

	return nil
}
