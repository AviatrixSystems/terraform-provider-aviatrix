package aviatrix

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixDistributedFirewallingIntraVpc_basic(t *testing.T) {
	skipAcc := os.Getenv("SKIP_DISTRIBUTED_FIREWALLING_INTRA_VPC")
	if skipAcc == "yes" {
		t.Skip("Skipping Distributed-firewalling Intra VPC test as SKIP_DISTRIBUTED_FIREWALLING_INTRA_VPC is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_distributed_firewalling_intra_vpc.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProvidersVersionValidation,
		CheckDestroy: testAccDistributedFirewallingIntraVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDistributedFirewallingIntraVpcBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDistributedFirewallingIntraVpcExists(resourceName),
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

func testAccDistributedFirewallingIntraVpcBasic(rName string) string {
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
resource "aviatrix_spoke_gateway" "test"{
  cloud_type   = aviatrix_account.test.cloud_type
  account_name = aviatrix_account.test.account_name
  gw_name      = "azure-spoke-0"
  vpc_id       = aviatrix_vpc.test.vpc_id
  vpc_reg      = aviatrix_vpc.test.region
  gw_size      = "Standard_D3_v2"
  subnet       = aviatrix_vpc.test.public_subnets[0].cidr
  enable_bgp   = true
}
resource "aviatrix_spoke_gateway" "test1"{
  cloud_type   = aviatrix_account.test.cloud_type
  account_name = aviatrix_account.test.account_name
  gw_name      = "azure-spoke-1"
  vpc_id       = aviatrix_vpc.test1.vpc_id
  vpc_reg      = aviatrix_vpc.test1.region
  gw_size      = "Standard_D3_v2"
  subnet       = aviatrix_vpc.test1.public_subnets[0].cidr
  enable_bgp   = true
}
resource "aviatrix_distributed_firewalling_intra_vpc" "test"{
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

  depends_on = [
    aviatrix_spoke_gateway.test,
    aviatrix_spoke_gateway.test1
  ]
}
	`, rName, os.Getenv("azure_subscription_id"), os.Getenv("azure_tenant_id"),
		os.Getenv("azure_client_id"), os.Getenv("azure_client_secret"))
}

func testAccCheckDistributedFirewallingIntraVpcExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no Distributed-firewalling Intra VPC resource found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Distributed-firewalling Intra VPC ID is set")
		}

		client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

		_, err := client.GetDistributedFirewallingIntraVpc(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get Distributed-firewalling Policy List status: %v", err)
		}

		if strings.Replace(client.ControllerIP, ".", "-", -1) != rs.Primary.ID {
			return fmt.Errorf("distributed-firewalling policy list ID not found")
		}

		return nil
	}
}

func testAccDistributedFirewallingIntraVpcDestroy(s *terraform.State) error {
	client := testAccProviderVersionValidation.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_distributed_firewalling_intra_vpc" {
			continue
		}

		_, err := client.GetDistributedFirewallingIntraVpc(context.Background())
		if err == nil || err != goaviatrix.ErrNotFound {
			return fmt.Errorf("distributed-firewalling intra vpc configured when it should be destroyed")
		}
	}

	return nil
}
