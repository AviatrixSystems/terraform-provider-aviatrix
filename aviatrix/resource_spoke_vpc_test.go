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

func TestAccAviatrixSpokeGw_basic(t *testing.T) {
	var gateway goaviatrix.Gateway
	rName := fmt.Sprintf("%s", acctest.RandString(5))
	resourceName := "aviatrix_spoke_vpc.test_spoke_vpc"

	msgCommon := ". Set SKIP_SPOKE to yes to skip Spoke Gateway tests"

	skipGw := os.Getenv("SKIP_SPOKE")
	if skipGw == "yes" {
		t.Skip("Skipping Spoke Gateway test as SKIP_SPOKE is set")
	}

	preGatewayCheck(t, msgCommon)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpokeGwDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpokeGwConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeGwExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "vpc_size", "t2.micro"),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-%s",
						rName)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AWS_VPC_ID")),
					resource.TestCheckResourceAttr(resourceName, "subnet", os.Getenv("AWS_VPC_NET")),
					resource.TestCheckResourceAttr(resourceName, "vpc_reg", os.Getenv("AWS_REGION")),
					resource.TestCheckResourceAttr(resourceName, "vpc_reg", os.Getenv("AWS_REGION")),
					resource.TestCheckResourceAttr(resourceName, "enable_nat", "no"),
					resource.TestCheckResourceAttr(resourceName, "tag_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tag_list.0", "k1:v1"),
					resource.TestCheckResourceAttr(resourceName, "tag_list.1", "k2:v2"),
				),
			},
		},
	})
}

func testAccSpokeGwConfigBasic(rName string) string {
	return fmt.Sprintf(`

resource "aviatrix_account" "test" {
  account_name = "tfa-%s"
  cloud_type = 1
  aws_account_number = "%s"
  aws_iam = "false"
  aws_access_key = "%s"
  aws_secret_key = "%s"
}

resource "aviatrix_spoke_vpc" "test_spoke_vpc" {
  cloud_type = 1
  account_name = "${aviatrix_account.test.account_name}"
  gw_name = "tfg-%[1]s"
  vpc_id = "%[5]s"
  vpc_reg = "%[6]s"
  vpc_size = "t2.micro"
  subnet = "%[7]s"
  enable_nat = "no"
  tag_list = ["k1:v1","k2:v2"]
}

	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_VPC_NET"))
}

func testAccCheckSpokeGwExists(n string, gateway *goaviatrix.Gateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("spoke gateway Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no spoke gateway ID is set")
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
			return fmt.Errorf("spoke gateway not found")
		}

		*gateway = *foundGateway

		return nil
	}
}

func testAccCheckSpokeGwDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_spoke_vpc" {
			continue
		}
		foundGateway := &goaviatrix.Gateway{
			GwName:      rs.Primary.Attributes["gw_name"],
			AccountName: rs.Primary.Attributes["account_name"],
		}
		_, err := client.GetGateway(foundGateway)

		if err == nil {
			return fmt.Errorf("spoke gateway still exists")
		}
	}
	return nil
}
