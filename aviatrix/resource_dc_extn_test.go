package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func preDCextnCheck(t *testing.T, msgEnd string) {
	if os.Getenv("AWS_REGION") == "" {
		t.Fatal("AWS_REGION must be set for acceptance tests." + msgEnd)
	}

	if os.Getenv("DCX_SUBNET") == "" {
		t.Fatal("DCX_SUBNET must be set for acceptance tests." + msgEnd)
	}
}

func TestDCX_basic(t *testing.T) {
	var dcx goaviatrix.DCExtn
	rName := fmt.Sprintf("tf-tst-%s", acctest.RandString(5))

	skipAcc := os.Getenv("SKIP_DCX")
	if skipAcc == "yes" {
		t.Skip("Skipping DCX test as SKIP_DCX is set")
	}

	msg := ". Set SKIP_DCX to yes to skip data centre extension tests"
	preAccountCheck(t, msg)

	preDCextnCheck(t, msg)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDCXDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDCXConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCXExists("aviatrix_dc_extn.test_dcx", &dcx),
					resource.TestCheckResourceAttr(
						"aviatrix_dc_extn.test_dcx", "account_name", "aws"),
					resource.TestCheckResourceAttr(
						"aviatrix_dc_extn.test_dcx", "gw_name", rName),
					resource.TestCheckResourceAttr(
						"aviatrix_dc_extn.test_dcx", "tunnel_type", "udp"),
					resource.TestCheckResourceAttr(
						"aviatrix_dc_extn.test_dcx", "subnet_cidr",
						os.Getenv("DCX_SUBNET")),
				),
			},
		},
	})
}

func testAccDCXConfigBasic(rName string) string {
	return fmt.Sprintf(`

resource "aviatrix_account" "foo" {
  account_name = "%s"
  cloud_type = 1
  aws_account_number = "%s"
  aws_iam = "false"
  aws_access_key = "%s"
  aws_secret_key = "%s"
}

resource "aviatrix_dc_extn" "test_dcx" {
  cloud_type = "1"
  account_name = "aws"
  gw_name = "%s"
  vpc_reg = "%s"
  gw_size = "t2.micro"
  subnet_cidr = "%s"
  internet_access = "no"
  public_subnet = "no"
  tunnel_type = "udp"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"),
		os.Getenv("AWS_SECRET_KEY"), rName, os.Getenv("AWS_REGION"), os.Getenv("DCX_SUBNET"))
}

func testAccCheckDCXExists(n string, _ *goaviatrix.DCExtn) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("DCX Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no DCX ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundGw := &goaviatrix.Gateway{
			AccountName: rs.Primary.Attributes["account_name"],
			GwName:      rs.Primary.Attributes["gw_name"],
		}

		_, err := client.GetGateway(foundGw)

		if err != nil {
			return err
		}

		if foundGw.GwName != rs.Primary.ID {
			return fmt.Errorf("DCX not found")
		}

		//*dcx = *foundGw

		return nil
	}
}

func testAccCheckDCXDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_dc_extn" {
			continue
		}
		foundGw := &goaviatrix.Gateway{
			AccountName: rs.Primary.Attributes["account_name"],
			GwName:      rs.Primary.Attributes["gw_name"],
		}
		_, err := client.GetGateway(foundGw)

		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("DCX still exists")
		}
	}
	return nil
}
