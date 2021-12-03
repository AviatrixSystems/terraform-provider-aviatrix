package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixSpokeExternalDeviceConn_basic(t *testing.T) {
	var externalDeviceConn goaviatrix.ExternalDeviceConn

	rName := acctest.RandString(5)
	resourceName := "aviatrix_spoke_external_device_conn.test"

	skipAcc := os.Getenv("SKIP_SPOKE_EXTERNAL_DEVICE_CONN")
	if skipAcc == "yes" {
		t.Skip("Skipping spoke external device connection tests as 'SKIP_SPOKE_EXTERNAL_DEVICE_CONN' is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, ". Set 'SKIP_SPOKE_EXTERNAL_DEVICE_CONN' to 'yes' to skip Site2Cloud spoke external device connection tests")
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpokeExternalDeviceConnDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpokeExternalDeviceConnConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeExternalDeviceConnExists(resourceName, &externalDeviceConn),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AWS_VPC_ID")),
					resource.TestCheckResourceAttr(resourceName, "connection_name", rName),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "connection_type", "bgp"),
					resource.TestCheckResourceAttr(resourceName, "bgp_local_as_num", "123"),
					resource.TestCheckResourceAttr(resourceName, "bgp_remote_as_num", "345"),
					resource.TestCheckResourceAttr(resourceName, "remote_gateway_ip", "172.12.13.14"),
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

func testAccSpokeExternalDeviceConnConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_spoke_gateway" "test" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	gw_name      = "tfg-%s"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t2.micro"
	subnet       = "%s"
    enable_bgp   = true
}
resource "aviatrix_spoke_external_device_conn" "test" {
	vpc_id            = aviatrix_spoke_gateway.test.vpc_id
	connection_name   = "%s"
	gw_name           = aviatrix_spoke_gateway.test.gw_name
	connection_type   = "bgp"
	bgp_local_as_num  = "123"
	bgp_remote_as_num = "345"
	remote_gateway_ip = "172.12.13.14"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"), rName)
}

func testAccCheckSpokeExternalDeviceConnExists(n string, externalDeviceConn *goaviatrix.ExternalDeviceConn) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("spoke external device connection Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no spoke external device connection ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundExternalDeviceConn := &goaviatrix.ExternalDeviceConn{
			VpcID:          rs.Primary.Attributes["vpc_id"],
			ConnectionName: rs.Primary.Attributes["connection_name"],
		}
		foundExternalDeviceConn2, err := client.GetExternalDeviceConnDetail(foundExternalDeviceConn)
		if err != nil {
			return err
		}
		if foundExternalDeviceConn2.ConnectionName+"~"+foundExternalDeviceConn2.VpcID != rs.Primary.ID {
			return fmt.Errorf("spoke external device connection not found")
		}

		*externalDeviceConn = *foundExternalDeviceConn2
		return nil
	}
}

func testAccCheckSpokeExternalDeviceConnDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_spoke_external_device_conn" {
			continue
		}

		foundExternalDeviceConn := &goaviatrix.ExternalDeviceConn{
			VpcID:          rs.Primary.Attributes["vpc_id"],
			ConnectionName: rs.Primary.Attributes["connection_name"],
		}

		_, err := client.GetExternalDeviceConnDetail(foundExternalDeviceConn)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("site2cloud still exists %s", err.Error())
		}
	}

	return nil
}
