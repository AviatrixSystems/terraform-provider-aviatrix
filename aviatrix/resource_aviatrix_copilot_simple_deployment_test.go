package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixCopilotSimpleDeployment_basic(t *testing.T) {
	if os.Getenv("SKIP_COPILOT_SIMPLE_DEPLOYMENT") == "yes" {
		t.Skip("Skipping Copilot Simple Deployment test as SKIP_COPILOT_SIMPLE_DEPLOYMENT is set")
	}

	resourceName := "aviatrix_copilot_simple_deployment.test"
	rName := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCopilotSimpleDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCopilotSimpleDeploymentBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCopilotSimpleDeploymentExists(resourceName),
				),
			},
		},
	})
}

func testAccCopilotSimpleDeploymentBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name 	   = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = "false"
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_copilot_simple_deployment" "test" {
	cloud_type                          = 1
	account_name                        = aviatrix_account.test.account_name
	region                              = "%s"
	vpc_id                              = "%s"
	subnet                              = "%s"
	controller_service_account_username = "%s"
	controller_service_account_password = "%s"
}
 `, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_REGION2"), os.Getenv("AWS_VPC_ID2"), os.Getenv("AWS_SUBNET2"),
		os.Getenv("AVIATRIX_USERNAME"), os.Getenv("AVIATRIX_PASSWORD"))
}

func testAccCheckCopilotSimpleDeploymentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("copilot not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no copilot id is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		_, err := client.GetCopilotAssociationStatus(context.Background())
		if err != nil {
			if err == goaviatrix.ErrNotFound {
				return fmt.Errorf("could not find copilot")
			}
			return err
		}

		return nil
	}
}

func testAccCheckCopilotSimpleDeploymentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_copilot_simple_deployment" {
			continue
		}

		_, err := client.GetCopilotAssociationStatus(context.Background())
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("copilot still exists")
		}
	}

	return nil
}
