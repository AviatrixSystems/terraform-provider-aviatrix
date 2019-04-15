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

func TestAccAviatrixTransitGw_basic(t *testing.T) {
	var gateway goaviatrix.Gateway
	rName := fmt.Sprintf("%s", acctest.RandString(5))
	resourceName := "aviatrix_transit_vpc.test_transit_vpc"

	msgCommon := ". Set SKIP_TRANSIT to yes to skip Transit Gateway tests"

	skipGw := os.Getenv("SKIP_TRANSIT")
	if skipGw == "yes" {
		t.Skip("Skipping Transit gateway test as SKIP_TRANSIT is set")
	}

	preGatewayCheck(t, msgCommon)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitGwDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGwConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGwExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "vpc_size", "t2.micro"),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-%s",
						rName)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AWS_VPC_ID")),
					resource.TestCheckResourceAttr(resourceName, "subnet", os.Getenv("AWS_VPC_NET")),
					resource.TestCheckResourceAttr(resourceName, "vpc_reg", os.Getenv("AWS_REGION")),
					resource.TestCheckResourceAttr(resourceName, "vpc_reg", os.Getenv("AWS_REGION")),
					resource.TestCheckResourceAttr(resourceName, "tag_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tag_list.0", "k1:v1"),
					resource.TestCheckResourceAttr(resourceName, "tag_list.1", "k2:v2"),
				),
			},
		},
	})
}

func testAccTransitGwConfigBasic(rName string) string {
	return fmt.Sprintf(`

resource "aviatrix_account" "test" {
  account_name = "tfa-%s"
  cloud_type = 1
  aws_account_number = "%s"
  aws_iam = "false"
  aws_access_key = "%s"
  aws_secret_key = "%s"
}

resource "aviatrix_transit_vpc" "test_transit_vpc" {
  cloud_type = 1
  account_name = "${aviatrix_account.test.account_name}"
  gw_name = "tfg-%[1]s"
  vpc_id = "%[5]s"
  vpc_reg = "%[6]s"
  vpc_size = "t2.micro"
  subnet = "%[7]s"
  tag_list = ["k1:v1","k2:v2"]
}

	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_VPC_NET"))
}

func testAccCheckTransitGwExists(n string, gateway *goaviatrix.Gateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("transit gateway Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no transit gateway ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundGateway := &goaviatrix.Gateway{
			GwName:      rs.Primary.Attributes["gw_name"],
			AccountName: rs.Primary.Attributes["account_name"],
		}

		_, err := client.GetGateway(foundGateway)

		if err != nil {
			return err
		}

		if foundGateway.GwName != rs.Primary.ID {
			return fmt.Errorf("transit gateway not found")
		}

		*gateway = *foundGateway

		return nil
	}
}

func testAccCheckTransitGwDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_transit_vpc" {
			continue
		}
		foundGateway := &goaviatrix.Gateway{
			GwName:      rs.Primary.Attributes["gw_name"],
			AccountName: rs.Primary.Attributes["account_name"],
		}
		_, err := client.GetGateway(foundGateway)

		if err == nil {
			return fmt.Errorf("transit gateway still exists")
		}
	}
	return nil
}
