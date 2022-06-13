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

func prePrivateModeMulticloudEndpointCheck(t *testing.T, msgCommon string) {
	if os.Getenv("AWS_VPC_ID") == "" {
		t.Fatal(fmt.Sprintf("%s must be set for Private Mode multicloud endpoint tests. %s", "AWS_VPC_ID", msgCommon))
	}
}

func TestAccAviatrixPrivateModeMulticloudEndpoint_basic(t *testing.T) {
	rName := acctest.RandString(5)

	skipAcc := os.Getenv("SKIP_PRIVATE_MODE_MULTICLOUD_ENDPOINT")
	if skipAcc == "yes" {
		t.Skip("Skipping Controller Private Mode multicloud endpoint tests as SKIP_PRIVATE_MODE_MULTICLOUD_ENDPOINT is set")
	}
	msgCommon := "Set SKIP_PRIVATE_MODE_MULTICLOUD_ENDPOINT to yes to skip Controller Private Mode load balancer tests"
	resourceName := "aviatrix_private_mode_multicloud_endpoint.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, msgCommon)
			prePrivateModeCheck(t, msgCommon)
			prePrivateModeMulticloudEndpointCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccAviatrixPrivateModeMulticloudEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAviatrixPrivateModeMulticloudEndpointBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccAviatrixPrivateModeMulticloudEndpointExists(resourceName),
					testAccAviatrixPrivateModeLbExists("aviatrix_private_mode_lb.test"),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "region", os.Getenv("AWS_REGION")),
					resource.TestCheckResourceAttr(resourceName, "controller_lb_vpc_id", os.Getenv("CONTROLLER_VPC_ID")),
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

func testAccAviatrixPrivateModeMulticloudEndpointBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name       = "tfa-%[1]s"
	cloud_type         = 1
	aws_account_number = "%[2]s"
	aws_iam            = false
	aws_access_key     = "%[3]s"
	aws_secret_key     = "%[4]s"
}

resource "aviatrix_controller_private_mode_config" "test" {
	enable_private_mode = true
}

resource "aviatrix_private_mode_lb" "test" {
	account_name = aviatrix_account.test_account.account_name
	vpc_id       = "%[5]s"
	region       = "%[6]s"
	lb_type      = "controller"

	depends_on = [aviatrix_controller_private_mode_config.test]
}

resource "aviatrix_private_mode_multicloud_endpoint" "test" {
	account_name         = "tfa-%[1]s"
	vpc_id               = "%[7]s"
	region               = "%[6]s"
	controller_lb_vpc_id = aviatrix_private_mode_lb.test.vpc_id
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), os.Getenv("CONTROLLER_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_VPC_ID"))
}

func testAccAviatrixPrivateModeMulticloudEndpointExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("private mode load multicloud endpoint Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no private mode multicloud endpoint ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		vpcId := rs.Primary.ID
		_, err := client.GetPrivateModeMulticloudEndpoint(context.Background(), vpcId)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccAviatrixPrivateModeMulticloudEndpointDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aviatrix_private_mode_lb" {
			vpcId := rs.Primary.ID
			_, err := client.GetPrivateModeLoadBalancer(context.Background(), vpcId)
			if err != nil {
				if err == goaviatrix.ErrNotFound {
					continue
				}
				return fmt.Errorf("failed to destroy Private Mode load balancer: %s", err)
			}
			return fmt.Errorf("failed to destroy Private Mode load balancer")
		} else if rs.Type == "aviatrix_private_mode_multicloud_load_balancer" {
			vpcId := rs.Primary.ID
			_, err := client.GetPrivateModeMulticloudEndpoint(context.Background(), vpcId)

			if err != nil {
				if err == goaviatrix.ErrNotFound {
					continue
				}
				return fmt.Errorf("failed to destroy Private Mode multicloud endpoint: %s", err)
			}
			return fmt.Errorf("failed to destroy Private Mode multicloud endpoint")
		}
		if rs.Type != "aviatrix_private_mode_lb" {
			continue
		}

	}

	return nil
}
