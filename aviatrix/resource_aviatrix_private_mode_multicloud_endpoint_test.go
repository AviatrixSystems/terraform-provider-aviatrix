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
	var lb goaviatrix.PrivateModeLbRead
	msgCommon := "Set SKIP_PRIVATE_MODE_MULTICLOUD_ENDPOINT to yes to skip Controller Private Mode load balancer tests"
	resourceName := "aviatrix_private_mode_multicloud_endpoint.test"

	awsVpcId := os.Getenv("AWS_VPC_ID")
	awsRegion := os.Getenv("AWS_REGION")
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
				Config: testAccAviatrixPrivateModeLbBasic(rName),
				Check:  testAccAviatrixPrivateModeLbExists("aviatrix_private_mode_lb.test", &lb),
			},
			{
				Config: testAccAviatrixPrivateModeMulticloudEndpointBasic(rName, awsVpcId, awsRegion, lb),
				Check: resource.ComposeTestCheckFunc(
					testAccAviatrixPrivateModeMulticloudEndpointExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "region", os.Getenv("AWS_REGION")),
					resource.TestCheckResourceAttr(resourceName, "controller_lb_vpc_id", lb.VpcId),
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

func testAccAviatrixPrivateModeMulticloudEndpointBasic(rName, vpcId, vpcRegion string, lb goaviatrix.PrivateModeLbRead) string {
	return fmt.Sprintf(`
resource "aviatrix_private_mode_multicloud_endpoint" "test" {
	account_name = "tfa-%s"
	vpc_id = "%s"
	region = "%s"
	controller_lb_vpc_id = "%s"
}
	`, rName, vpcId, vpcRegion, lb.VpcId)
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
