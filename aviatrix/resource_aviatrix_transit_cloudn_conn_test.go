package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testTransitCloudnConnPreCheck(t *testing.T) {
	for _, key := range []string{"AWS_VPC_ID", "AWS_ACCOUNT_NUMBER", "AWS_ACCESS_KEY", "AWS_SECRET_KEY", "AWS_REGION", "AWS_SUBNET", "CLOUDN_IP", "CLOUDN_NEIGHBOR_IP"} {
		if os.Getenv(key) == "" {
			t.Fatalf("%s must be set for Transit CloudN Connection acceptance tests. Set SKIP_TRANSIT_CLOUDN_CONN to yes to skip acceptance tests", key)
		}
	}
}

func TestAccAviatrixTransitCloudnConn_basic(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_transit_cloudn_conn.test"

	skipAcc := os.Getenv("SKIP_TRANSIT_CLOUDN_CONN")
	if skipAcc == "yes" {
		t.Skip("Skipping transit CloudN Connection tests as 'SKIP_TRANSIT_CLOUDN_CONN' is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, ". Set 'SKIP_TRANSIT_CLOUDN_CONN' to 'yes' to skip Site2Cloud transit CloudN Conn tests")
			testTransitCloudnConnPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitCloudnConnDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitCloudnConnBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitCloudnConnExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AWS_VPC_ID")),
					resource.TestCheckResourceAttr(resourceName, "connection_name", rName),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "bgp_local_as_num", "65003"),
					resource.TestCheckResourceAttr(resourceName, "cloudn_as_num", "65004"),
					resource.TestCheckResourceAttr(resourceName, "cloudn_remote_ip", os.Getenv("CLOUDN_IP")),
					resource.TestCheckResourceAttr(resourceName, "cloudn_neighbor_ip", os.Getenv("CLOUDN_NEIGHBOR_IP")),
					resource.TestCheckResourceAttr(resourceName, "cloudn_neighbor_as_num", "65005"),
					resource.TestCheckResourceAttr(resourceName, "enable_ha", "false"),
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

func testAccTransitCloudnConnBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%[1]s"
	cloud_type         = 1
	aws_account_number = "%[2]s"
	aws_iam            = false
	aws_access_key     = "%[3]s"
	aws_secret_key     = "%[4]s"
}
resource "aviatrix_transit_gateway" "test" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	gw_name      = "tfg-%[1]s"
	vpc_id       = "%[5]s"
	vpc_reg      = "%[6]s"
	gw_size      = "t2.micro"
	subnet       = "%[7]s"
}
resource "aviatrix_transit_cloudn_conn" "test" {
	vpc_id                 = aviatrix_transit_gateway.test.vpc_id
	connection_name        = "%[1]s"
	gw_name                = aviatrix_transit_gateway.test.gw_name
	bgp_local_as_num       = "65003"
	cloudn_remote_ip       = "%[8]s"
	cloudn_as_num          = "65004"
	cloudn_neighbor_ip     = "%[9]s"
	cloudn_neighbor_as_num = "65005"
	enable_ha              = false
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"),
		os.Getenv("CLOUDN_IP"), os.Getenv("CLOUDN_NEIGHBOR_IP"))
}

func testAccCheckTransitCloudnConnExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("transit cloudn connection Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no transit cloudn connection ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		transitCloudnConn := &goaviatrix.TransitCloudnConn{
			VpcID:          rs.Primary.Attributes["vpc_id"],
			ConnectionName: rs.Primary.Attributes["connection_name"],
		}

		cloudnConn, err := client.GetTransitCloudnConn(context.Background(), transitCloudnConn)
		if err != nil {
			return err
		}
		if fmt.Sprintf("%s~%s", cloudnConn.ConnectionName, cloudnConn.VpcID) != rs.Primary.ID {
			return fmt.Errorf("transit cloudn connection not found")
		}

		return nil
	}
}

func testAccCheckTransitCloudnConnDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_transit_cloudn_conn" {
			continue
		}

		transitCloudnConn := &goaviatrix.TransitCloudnConn{
			VpcID:          rs.Primary.Attributes["vpc_id"],
			ConnectionName: rs.Primary.Attributes["connection_name"],
		}

		_, err := client.GetTransitCloudnConn(context.Background(), transitCloudnConn)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("site2cloud transit cloudn connection still exists %s", err.Error())
		}
	}

	return nil
}
