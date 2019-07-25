package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixVPNUserAccelerator_basic(t *testing.T) {
	rName := fmt.Sprintf("%s", acctest.RandString(5))
	resourceName := "aviatrix_vpn_user_accelerator.test_elb"

	skipXlr := os.Getenv("SKIP_VPN_USER_ACCELERATOR")
	if skipXlr == "yes" {
		t.Skip("SKipping VPN User Accelerator test as SKIP_VPN_USER_ACCELERATOR is set")
	}
	msgCommon := ". Set SKIP_VPN_USER_ACCELERATOR to skip VPN User Accelerator tests"

	preGatewayCheck(t, msgCommon)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVPNUserAcceleratorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNUserAcceleratorConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPNUserAcceleratorExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "elb_name", fmt.Sprintf("tflb-%s", rName)),
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

func testAccVPNUserAcceleratorConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name 	   = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key 	   = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_gateway" "test_gw" {
	cloud_type   = 1
	account_name = aviatrix_account.test_account.account_name
	gw_name      = "tfg-%[1]s"
	vpc_id       = "%[5]s"
	vpc_reg      = "%[6]s"
	gw_size      = "t2.micro"
	subnet       = "%[7]s"
	vpn_access   = true
	vpn_cidr     = "192.168.43.0/24"
	max_vpn_conn = "100"
	enable_elb   = true
	elb_name     = "tflb-%[1]s"
}
resource "aviatrix_vpn_user_accelerator" "test_elb" {
	elb_name = aviatrix_gateway.test_gw.elb_name
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"))
}

func testAccCheckVPNUserAcceleratorExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("vpn user accelerator not found : %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no vpn user accelerator ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		elbList, err := client.GetVpnUserAccelerator()
		if err != nil {
			return err
		}
		if !goaviatrix.Contains(elbList, rs.Primary.ID) {
			return fmt.Errorf("vpn user accelerator ID not found")
		}

		return nil
	}
}

func testAccCheckVPNUserAcceleratorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_vpn_user_accelerator" {
			continue
		}

		elbList, err := client.GetVpnUserAccelerator()
		if err != nil {
			return fmt.Errorf("error retrieving vpn user accelerator: %s", err)
		}
		if goaviatrix.Contains(elbList, rs.Primary.ID) {
			return fmt.Errorf("vpn user accelerator still exists")
		}
	}

	return nil
}
