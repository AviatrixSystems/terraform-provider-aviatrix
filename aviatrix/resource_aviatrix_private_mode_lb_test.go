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

func prePrivateModeCheck(t *testing.T, msgEnd string) {
	for _, key := range []string{"CONTROLLER_VPC_ID", "AWS_REGION"} {
		if os.Getenv(key) == "" {
			t.Fatal(fmt.Sprintf("%s must be set for Private Mode tests using load balancers. %s", key, msgEnd))
		}
	}
}
func TestAccAviatrixPrivateModeLb_basic(t *testing.T) {
	rName := acctest.RandString(5)

	skipAcc := os.Getenv("SKIP_PRIVATE_MODE_LB")
	if skipAcc == "yes" {
		t.Skip("Skipping Controller Private Mode load balancer tests as SKIP_PRIVATE_MODE_LB is set")
	}
	msgCommon := "Set SKIP_PRIVATE_MODE_LB to yes to skip Controller Private Mode load balancer tests"
	resourceName := "aviatrix_private_mode_lb.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, msgCommon)
			prePrivateModeCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccAviatrixPrivateModeLbDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAviatrixPrivateModeLbBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccAviatrixPrivateModeLbExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "region", os.Getenv("AWS_REGION")),
					resource.TestCheckResourceAttr(resourceName, "lb_type", "controller"),
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

func testAccAviatrixPrivateModeLbBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_controller_private_mode_config" "test" {
	enable_private_mode = true
}

resource "aviatrix_private_mode_lb" "test" {
	account_name = aviatrix_account.test_account.account_name
	vpc_id       = "%s"
	region       = "%s"
	lb_type      = "controller"

	depends_on = [aviatrix_controller_private_mode_config.test]
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), os.Getenv("CONTROLLER_VPC_ID"), os.Getenv("AWS_REGION"))
}

func testAccAviatrixPrivateModeLbExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("private mode load balancer Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no private mode load balancer ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		vpcId := rs.Primary.ID
		_, err := client.GetPrivateModeLoadBalancer(context.Background(), vpcId)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccAviatrixPrivateModeLbDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_private_mode_lb" {
			continue
		}

		vpcId := rs.Primary.ID
		_, err := client.GetPrivateModeLoadBalancer(context.Background(), vpcId)
		if err != nil {
			if err == goaviatrix.ErrNotFound {
				return nil
			}
			return fmt.Errorf("error getting Private Mode load balancer after destroy: %s", err)
		}
		return fmt.Errorf("failed to destroy Private Mode load balancer")
	}

	return nil
}
