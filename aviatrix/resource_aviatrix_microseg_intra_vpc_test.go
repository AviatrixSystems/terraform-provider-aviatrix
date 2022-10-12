package aviatrix

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixMicrosegIntraVpc_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_MICROSEG_INTRA_VPC")
	if skipAcc == "yes" {
		t.Skip("Skipping Microseg Intra VPC test as SKIP_MICROSEG_INTRA_VPC is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_microseg_intra_vpc.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccMicrosegIntraVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMicrosegIntraVpcBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMicrosegIntraVpcExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "vpcs.0.account_name", fmt.Sprintf("tfa-azure-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "vpcs.0.region", "Central US"),
					resource.TestCheckResourceAttr(resourceName, "vpcs.1.region", "Central US"),
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

func testAccMicrosegIntraVpcBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
  account_name        = "tfa-azure-%[1]s"
  cloud_type          = 8
  arm_subscription_id = "%[2]s"
  arm_directory_id    = "%[3]s"
  arm_application_id  = "%[4]s"
  arm_application_key = "%[5]s"
}
resource "aviatrix_vpc" "test" {
  cloud_type   = aviatrix_account.test.cloud_type
  account_name = aviatrix_account.test.account_name
  region       = "Central US"
  name         = "azure-vpc-0-%[1]s"
  cidr         = "15.0.0.0/20"
}
resource "aviatrix_vpc" "test1" {
  cloud_type   = aviatrix_account.test.cloud_type
  account_name = aviatrix_account.test.account_name
  region       = "Central US"
  name         = "azure-vpc-1-%[1]s"
  cidr         = "16.0.0.0/20"
}
resource "aviatrix_microseg_intra_vpc" "test"{
  vpcs {
    account_name = aviatrix_vpc.test.account_name
    vpc_id       = aviatrix_vpc.test.vpc_id
    region       = aviatrix_vpc.test.region
  }

  vpcs {
    account_name = aviatrix_vpc.test1.account_name
    vpc_id       = aviatrix_vpc.test1.vpc_id
    region       = aviatrix_vpc.test1.region
  }
}
	`, rName, os.Getenv("azure_subscription_id"), os.Getenv("azure_tenant_id"),
		os.Getenv("azure_client_id"), os.Getenv("azure_client_secret"))
}

func testAccCheckMicrosegIntraVpcExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no Micro-segmentation Intra VPC resource found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Micro-segmentation Intra VPC ID is set")
		}

		client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

		_, err := client.GetMicrosegIntraVpc(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get Micro-segmentation Policy List status: %v", err)
		}

		if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
			return fmt.Errorf("micro-segmentation policy list ID not found")
		}

		return nil
	}
}

func testAccMicrosegIntraVpcDestroy(s *terraform.State) error {
	client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_microseg_intra_vpc" {
			continue
		}

		_, err := client.GetMicrosegIntraVpc(context.Background())
		if err == nil || err != goaviatrix.ErrNotFound {
			return fmt.Errorf("micro-segmentation intra vpc configured when it should be destroyed")
		}
	}

	return nil
}
