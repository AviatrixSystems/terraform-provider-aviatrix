package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixGlobalVpcExcludedInstance_basic(t *testing.T) {
	if os.Getenv("SKIP_GLOBAL_VPC_EXCLUDED_INSTANCE") == "yes" {
		t.Skip("Skipping global vpc excluded instance test as SKIP_GLOBAL_VPC_EXCLUDED_INSTANCE is set")
	}

	resourceName := "aviatrix_global_vpc_excluded_instance.test"
	rName := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGlobalVpcExcludedInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGlobalVpcExcludedInstanceBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGlobalVpcExcludedInstanceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-gcp-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "instance_name", fmt.Sprintf("tfg-gcp-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "region", os.Getenv("GCP_ZONE")),
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

func testAccGlobalVpcExcludedInstanceBasic(rName string) string {
	gcpGwSize := os.Getenv("GCP_GW_SIZE")
	if gcpGwSize == "" {
		gcpGwSize = "n1-standard-1"
	}
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_gcp" {
	account_name                        = "tfa-gcp-%s"
	cloud_type                          = 4
	gcloud_project_id                   = "%s"
	gcloud_project_credentials_filepath = "%s"
}
resource "aviatrix_transit_gateway" "test_transit_gateway_gcp" {
	cloud_type   = 4
	account_name = aviatrix_account.test_acc_gcp.account_name
	gw_name      = "tfg-gcp-%[1]s"
	vpc_id       = "%[4]s"
	vpc_reg      = "%[5]s"
	gw_size      = "%[6]s"
	subnet       = "%[7]s"
}
resource "aviatrix_global_vpc_excluded_instance" "test" {
	account_name  = aviatrix_account.test_acc_gcp.account_name
	instance_name = aviatrix_transit_gateway.test_transit_gateway_gcp.gw_name
	region        = aviatrix_transit_gateway.test_transit_gateway_gcp.vpc_reg
}
	`, rName, os.Getenv("GCP_ID"), os.Getenv("GCP_CREDENTIALS_FILEPATH"),
		os.Getenv("GCP_VPC_ID"), os.Getenv("GCP_ZONE"), gcpGwSize, os.Getenv("GCP_SUBNET"))
}

func testAccCheckGlobalVpcExcludedInstanceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("global vpc excluded instance not found: %s", resourceName)
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		_, err := client.GetGlobalVpcExcludedInstance(context.Background(), rs.Primary.Attributes["uuid"])
		if err == goaviatrix.ErrNotFound {
			return fmt.Errorf("global vpc excluded instance not found")
		}

		return nil
	}
}

func testAccCheckGlobalVpcExcludedInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_global_vpc_excluded_instance" {
			continue
		}

		_, err := client.GetGlobalVpcExcludedInstance(context.Background(), rs.Primary.Attributes["uuid"])
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("global vpc excluded instance still exists")
		}
	}

	return nil
}
