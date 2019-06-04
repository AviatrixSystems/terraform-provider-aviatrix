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

func TestAccAviatrixVpc_basic(t *testing.T) {
	var vpc goaviatrix.Vpc
	rName := fmt.Sprintf("%s", acctest.RandString(5))
	resourceName := "aviatrix_vpc.test_vpc"

	skipAcc := os.Getenv("SKIP_VPC")
	if skipAcc == "yes" {
		t.Skip("Skipping VPC test as SKIP_VPC is set")
	}
	msgCommon := ". Set SKIP_VPC to yes to skip VPC tests"
	preAccountCheck(t, msgCommon)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists("aviatrix_vpc.test_vpc", &vpc),
					resource.TestCheckResourceAttr(
						resourceName, "name", fmt.Sprintf("tfg-%s", rName)),
					resource.TestCheckResourceAttr(
						resourceName, "account_name", fmt.Sprintf("tfa-%s", rName)),
					resource.TestCheckResourceAttr(
						resourceName, "cloud_type", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "aviatrix_transit_vpc", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "region", os.Getenv("AWS_REGION")),
					resource.TestCheckResourceAttr(
						resourceName, "cidr", "10.0.0.0/16"),
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

func testAccVpcConfigBasic(rName string) string {
	return fmt.Sprintf(`

resource "aviatrix_account" "test_account" {
    account_name = "tfa-%s"
    cloud_type = 1
    aws_account_number = "%s"
    aws_iam = "false"
    aws_access_key = "%s"
    aws_secret_key = "%s"
}

resource "aviatrix_vpc" "test_vpc" {
	cloud_type = 1
	account_name = "${aviatrix_account.test_account.account_name}"
	name = "tfg-%s"
	region = "%s"
	cidr = "10.0.0.0/16"
}
`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_REGION"))
}

func testAccCheckVpcExists(n string, vpc *goaviatrix.Vpc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("VPC Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no VPC ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundVpc := &goaviatrix.Vpc{
			Name: rs.Primary.Attributes["name"],
		}

		foundVpc2, err := client.GetVpc(foundVpc)
		if err != nil {
			return err
		}
		if foundVpc2.Name != rs.Primary.ID {
			return fmt.Errorf("VPC not found")
		}
		*vpc = *foundVpc2

		return nil
	}
}

func testAccCheckVpcDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_vpc" {
			continue
		}
		foundVpc := &goaviatrix.Vpc{
			Name: rs.Primary.Attributes["name"],
		}
		_, err := client.GetVpc(foundVpc)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("VPC still exists")
		}
	}
	return nil
}
