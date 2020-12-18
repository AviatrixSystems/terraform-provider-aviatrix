package aviatrix

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func preAzureSpokeNativePeeringCheck(t *testing.T, msgCommon string) {
	vNet1 := os.Getenv("AZURE_VNET_ID")
	if vNet1 == "" {
		t.Fatal("Environment variable 'AZURE_VNET_ID' is not set" + msgCommon)
	}
	vNet2 := os.Getenv("AZURE_VNET_ID2")
	if vNet2 == "" {
		t.Fatal("Environment variable 'AZURE_VNET_ID2' is not set" + msgCommon)
	}
	region1 := os.Getenv("AZURE_REGION")
	if region1 == "" {
		t.Fatal("Environment variable 'AZURE_REGION' is not set" + msgCommon)
	}
	region2 := os.Getenv("AZURE_REGION2")
	if region2 == "" {
		t.Fatal("Environment variable 'AZURE_REGION2' is not set" + msgCommon)
	}
	gwSize := os.Getenv("AZURE_GW_SIZE")
	if gwSize == "" {
		t.Fatal("Environment variable 'AZURE_GW_SIZE' is not set" + msgCommon)
	}
	subnet := os.Getenv("AZURE_SUBNET")
	if subnet == "" {
		t.Fatal("Environment variable 'AZURE_SUBNET' is not set" + msgCommon)
	}
}

func TestAccAviatrixAzureSpokeNativePeering_basic(t *testing.T) {
	var azureSpokeNativePeering goaviatrix.AzureSpokeNativePeering

	rName := fmt.Sprintf("%s", acctest.RandString(5))
	resourceName := "aviatrix_azure_spoke_native_peering.test"

	skipAcc := os.Getenv("SKIP_AZURE_SPOKE_NATIVE_PEERING")
	if skipAcc == "yes" {
		t.Skip("Skipping Aviatrix Azure spoke native peering tests as SKIP_AZURE_SPOKE_NATIVE_PEERING is set")
	}

	msgCommon := ". Set SKIP_AZURE_SPOKE_NATIVE_PEERING to yes to skip Azure spoke native peering tests"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, msgCommon)
			preAzureSpokeNativePeeringCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAzureSpokeNativePeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureSpokeNativePeeringConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					tesAccCheckAzureSpokeNativePeeringExists(resourceName, &azureSpokeNativePeering),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_name", fmt.Sprintf("tfg-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "spoke_account_name", fmt.Sprintf("tf-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "spoke_region", os.Getenv("AZURE_REGION2")),
					resource.TestCheckResourceAttr(resourceName, "spoke_vpc_id", os.Getenv("AZURE_VNET_ID2")),
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

func testAccAzureSpokeNativePeeringConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name        = "tf-%s"
	cloud_type          = 8
	arm_subscription_id = "%s"
	arm_directory_id    = "%s"
	arm_application_id  = "%s"
	arm_application_key = "%s"
}
resource "aviatrix_transit_gateway" "test" {
	cloud_type   = aviatrix_account.test_account.cloud_type
	account_name = aviatrix_account.test_account.account_name
	gw_name      = "tfg-%s"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "%s"
	subnet       = "%s"
}
resource "aviatrix_azure_spoke_native_peering" "test" {
	transit_gateway_name = aviatrix_transit_gateway.test.gw_name
	spoke_account_name   = aviatrix_transit_gateway.test.account_name
	spoke_region         = "%s"
	spoke_vpc_id         = "%s"
}
	`, rName, os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_DIRECTORY_ID"),
		os.Getenv("ARM_APPLICATION_ID"), os.Getenv("ARM_APPLICATION_KEY"), rName,
		os.Getenv("AZURE_VNET_ID"), os.Getenv("AZURE_REGION"),
		os.Getenv("AZURE_GW_SIZE"), os.Getenv("AZURE_SUBNET"),
		os.Getenv("AZURE_REGION2"), os.Getenv("AZURE_VNET_ID2"))
}

func tesAccCheckAzureSpokeNativePeeringExists(n string, azureSpokeNativePeering *goaviatrix.AzureSpokeNativePeering) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("azure spoke native peering Not Created: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Azure spoke native peering ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundPeering := &goaviatrix.AzureSpokeNativePeering{
			TransitGatewayName: rs.Primary.Attributes["transit_gateway_name"],
			SpokeAccountName:   rs.Primary.Attributes["spoke_account_name"],
			SpokeVpcID:         rs.Primary.Attributes["spoke_vpc_id"],
		}

		foundPeering2, err := client.GetAzureSpokeNativePeering(foundPeering)
		if err != nil {
			if err == goaviatrix.ErrNotFound {
				return fmt.Errorf("azure spoke native peering does not exist")
			}
			return err
		}
		if foundPeering2.TransitGatewayName != rs.Primary.Attributes["transit_gateway_name"] {
			return fmt.Errorf("'transit_gateway_name' Not found in created attributes")
		}
		if foundPeering2.SpokeAccountName != rs.Primary.Attributes["spoke_account_name"] {
			return fmt.Errorf("'spoke_account_name' Not found in created attributes")
		}
		if foundPeering2.SpokeVpcID != rs.Primary.Attributes["spoke_vpc_id"] {
			return fmt.Errorf("'spoke_vpc_id' Not found in created attributes")
		}

		*azureSpokeNativePeering = *foundPeering2
		return nil
	}
}

func testAccCheckAzureSpokeNativePeeringDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_azure_spoke_native_peering" {
			continue
		}

		foundPeering := &goaviatrix.AzureSpokeNativePeering{
			TransitGatewayName: rs.Primary.Attributes["transit_gateway_name"],
			SpokeAccountName:   rs.Primary.Attributes["spoke_account_name"],
			SpokeVpcID:         rs.Primary.Attributes["spoke_vpc_id"],
		}

		_, err := client.GetAzureSpokeNativePeering(foundPeering)
		if err != goaviatrix.ErrNotFound {
			if strings.Contains(err.Error(), "not found in controller database") {
				return nil
			}
			return fmt.Errorf("azure spoke native peering still exists")
		}
	}

	return nil
}
