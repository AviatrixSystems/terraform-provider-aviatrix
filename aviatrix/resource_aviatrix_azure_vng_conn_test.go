package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixAzureVngConn_basic(t *testing.T) {
	if os.Getenv("SKIP_AZURE_VNG_CONN") == "yes" {
		t.Skip("Skipping azure vng conn test as SKIP_AZURE_VNG_CONN is set")
	}

	connectionName := acctest.RandString(5)
	resourceName := "aviatrix_azure_vng_conn.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAzureVngConnDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureVngConnBasic(connectionName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAzureVngConnExists(resourceName, connectionName),
					resource.TestCheckResourceAttr(resourceName, "primary_gateway_name", "test-tgw-azure"),
					resource.TestCheckResourceAttr(resourceName, "connection_name", connectionName),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AZURE_VNG_VNET_ID")),
					resource.TestCheckResourceAttr(resourceName, "vng_name", os.Getenv("AZURE_VNG")),
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

func testAccAzureVngConnBasic(connectionName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name        = "test-acc-azure"
	cloud_type          = 8
	arm_subscription_id = "%s"
	arm_directory_id    = "%s"
	arm_application_id  = "%s"
	arm_application_key = "%s"
}
resource "aviatrix_transit_gateway" "test" {
	cloud_type     = 8
	account_name   = aviatrix_account.test.account_name
	gw_name        = "test-tgw-azure"
	vpc_id         = "%s"
	vpc_reg        = "%s"
	gw_size        = "Standard_B2ms"
	subnet         = "%s"
	ha_subnet      = "%[7]s"
	ha_gw_size     = "Standard_B2ms"
	connected_transit = true
	enable_active_mesh = true
	enable_transit_firenet = true
}
resource "aviatrix_azure_vng_conn" "test" {
	primary_gateway_name = aviatrix_transit_gateway.test.gw_name
	connection_name      = "%s"
}
	`, os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_DIRECTORY_ID"),
		os.Getenv("ARM_APPLICATION_ID"), os.Getenv("ARM_APPLICATION_KEY"),
		os.Getenv("AZURE_VNG_VNET_ID"), os.Getenv("AZURE_REGION"),
		os.Getenv("AZURE_VNG_SUBNET"), connectionName)
}

func testAccCheckAzureVngConnExists(resourceName string, connectionName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("azure vng conn not found: %s", resourceName)
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		resp, err := client.GetAzureVngConnStatus(connectionName)
		if err == nil && !resp.Attached {
			return fmt.Errorf("azure vng conn %s not attached to the gateway", connectionName)
		} else if err != nil {
			return fmt.Errorf("azure vng conn %s not found", connectionName)
		}

		return nil
	}
}

func testAccCheckAzureVngConnDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_azure_vng_conn" {
			continue
		}

		connectionName := rs.Primary.Attributes["connection_name"]

		_, err := client.GetAzureVngConnStatus(connectionName)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("azure_vpn_conn still exists")
		}
	}

	return nil
}
