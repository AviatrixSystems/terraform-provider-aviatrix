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

func TestAccAviatrixSymmetricRoutingConfig_basic(t *testing.T) {
	rName := acctest.RandString(5)

	skipAcc := os.Getenv("SKIP_SYMMETRIC_ROUTING_CONFIG")
	if skipAcc == "yes" {
		t.Skip("Skipping Symmetric Routing Config test as SKIP_SYMMETRIC_ROUTING_CONFIG is set")
	}

	resourceName := "aviatrix_symmetric_routing_config.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSymmetricRoutingConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSymmetricRoutingConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSymmetricRoutingConfigExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_symmetric_routing", "true"),
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

func testAccSymmetricRoutingConfigBasic(rName string) string {
	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t2.micro"
	}
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_spoke_gateway" "test_spoke_gateway" {
	cloud_type     = 1
	account_name   = aviatrix_account.test_acc_aws.account_name
	gw_name        = "tfg-aws-%[1]s"
	vpc_id         = "%[5]s"
	vpc_reg        = "%[6]s"
	gw_size        = "%[7]s"
	subnet         = "%[8]s"
	single_ip_snat = false
}
resource "aviatrix_symmetric_routing_config" "test" {
	gw_name                  = aviatrix_spoke_gateway.test_spoke_gateway.gw_name
	enable_symmetric_routing = true
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID4"), os.Getenv("AWS_REGION"), awsGwSize, os.Getenv("AWS_SUBNET4"))
}

func testAccCheckSymmetricRoutingConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("symmetric routing config ID not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no symmetric routing config ID is set")
		}

		if rs.Primary.Attributes["gw_name"] != rs.Primary.ID {
			return fmt.Errorf("symmetric routing config ID not found")
		}

		return nil
	}
}

func testAccCheckSymmetricRoutingConfigDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_symmetric_routing_config" {
			continue
		}

		_, err := client.GetSymmetricRoutingStatus(context.Background(), rs.Primary.Attributes["gw_name"])
		if err == nil {
			return fmt.Errorf("symmetric routing config still exists")
		}
	}

	return nil
}
