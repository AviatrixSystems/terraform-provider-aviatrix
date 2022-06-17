package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixCopilotSecurityGroupManagementConfig_basic(t *testing.T) {
	if os.Getenv("SKIP_COPILOT_SECURITY_GROUP_MANAGEMENT_CONFIG") == "yes" {
		t.Skip("Skipping copilot security group management config test as SKIP_COPILOT_SECURITY_GROUP_MANAGEMENT_CONFIG is set")
	}

	resourceName := "aviatrix_copilot_security_group_management_config.test"
	rName := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCopilotSecurityGroupManagementConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCopilotSecurityGroupManagementConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCopilotSecurityGroupManagementConfigExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-aws-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AWS_VPC_ID")),
					resource.TestCheckResourceAttr(resourceName, "region", os.Getenv("AWS_REGION")),
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

func testAccCopilotSecurityGroupManagementConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_transit_gateway" "test_transit_gateway_aws" {
	cloud_type   = 1
	account_name = aviatrix_account.test_acc_aws.account_name
	gw_name      = "tfg-aws-%[1]s"
	vpc_id       = "%[5]s"
	vpc_reg      = "%[6]s"
	gw_size      = "t2.micro"
	subnet       = "%[7]s"
}
resource "aviatrix_copilot_security_group_management_config" "test" {
	cloud_type   = 1
	account_name = aviatrix_account.test_acc_aws.account_name
	region       = aviatrix_transit_gateway.test_transit_gateway_aws.vpc_reg
	vpc_id       = aviatrix_transit_gateway.test_transit_gateway_aws.vpc_id
	instance_id  = aviatrix_transit_gateway.test_transit_gateway_aws.cloud_instance_id
}
`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"))
}

func testAccCheckCopilotSecurityGroupManagementConfigExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("could not find copilot security group management config: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("copilot security group management config id is not set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		copilotSecurityGroupManagementConfig, err := client.GetCopilotSecurityGroupManagementConfig(context.Background())
		if err != nil {
			return err
		}
		if copilotSecurityGroupManagementConfig.InstanceID != rs.Primary.ID {
			return fmt.Errorf("could not find copilot security group management id")
		}
		return nil
	}
}

func testAccCheckCopilotSecurityGroupManagementConfigDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_copilot_security_group_management_config" {
			continue
		}

		_, err := client.GetCopilotSecurityGroupManagementConfig(context.Background())
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("copilot security group management is still enabled")
		}
	}
	return nil
}
