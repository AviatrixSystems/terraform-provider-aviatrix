package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func preAzurePeerCheck(t *testing.T, msgCommon string) {
	vNet1 := os.Getenv("AZURE_VNET_ID")
	if vNet1 == "" {
		t.Fatal("Environment variable AZURE_VNET_ID is not set" + msgCommon)
	}
	vNet2 := os.Getenv("AZURE_VNET_ID2")
	if vNet2 == "" {
		t.Fatal("Environment variable AZURE_VNET_ID2 is not set" + msgCommon)
	}

	region1 := os.Getenv("AZURE_REGION")
	if region1 == "" {
		t.Fatal("Environment variable AZURE_REGION is not set" + msgCommon)
	}
	region2 := os.Getenv("AZURE_REGION2")
	if region2 == "" {
		t.Fatal("Environment variable AZURE_REGION2 is not set" + msgCommon)
	}
}

func TestAccAviatrixAzurePeer_basic(t *testing.T) {
	var azurePeer goaviatrix.AzurePeer
	vNet1 := os.Getenv("AZURE_VNET_ID")
	vNet2 := os.Getenv("AZURE_VNET_ID2")
	region1 := os.Getenv("AZURE_REGION")
	region2 := os.Getenv("AZURE_REGION2")

	rInt := acctest.RandInt()
	resourceName := "aviatrix_azure_peer.test_azure_peer"

	skipAcc := os.Getenv("SKIP_AZURE_PEER")
	if skipAcc == "yes" {
		t.Skip("Skipping Aviatrix Azure peering tests as SKIP_AZURE_PEER is set")
	}
	msgCommon := ". Set SKIP_AZURE_PEER to yes to skip Azure peer tests"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, msgCommon)
			preAzurePeerCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAzurePeerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzurePeerConfigBasic(rInt, vNet1, vNet2, region1, region2),
				Check: resource.ComposeTestCheckFunc(
					tesAccCheckAzurePeerExists(resourceName, &azurePeer),
					resource.TestCheckResourceAttr(resourceName, "vnet_name_resource_group1", vNet1),
					resource.TestCheckResourceAttr(resourceName, "vnet_name_resource_group2", vNet2),
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

func testAccAzurePeerConfigBasic(rInt int, vNet1 string, vNet2 string, region1 string, region2 string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name        = "tf-testing-%d"
	cloud_type          = 8
	arm_subscription_id = "%s"
	arm_directory_id    = "%s"
	arm_application_id  = "%s"
	arm_application_key = "%s"
}
resource "aviatrix_azure_peer" "test_azure_peer" {
	account_name1             = aviatrix_account.test_account.account_name
	account_name2             = aviatrix_account.test_account.account_name
	vnet_name_resource_group1 = "%s"
	vnet_name_resource_group2 = "%s"
	vnet_reg1                 = "%s"
	vnet_reg2                 = "%s"
}
	`, rInt, os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_DIRECTORY_ID"),
		os.Getenv("ARM_APPLICATION_ID"), os.Getenv("ARM_APPLICATION_KEY"),
		vNet1, vNet2, region1, region2)

}

func tesAccCheckAzurePeerExists(n string, azurePeer *goaviatrix.AzurePeer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("azure peer Not Created: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Azure peer ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundPeer := &goaviatrix.AzurePeer{
			VNet1: rs.Primary.Attributes["vnet_name_resource_group1"],
			VNet2: rs.Primary.Attributes["vnet_name_resource_group2"],
		}

		_, err := client.GetAzurePeer(foundPeer)
		if err != nil {
			return err
		}
		if foundPeer.VNet1 != rs.Primary.Attributes["vnet_name_resource_group1"] {
			return fmt.Errorf("vnet_name_resource_group1 Not found in created attributes")
		}
		if foundPeer.VNet2 != rs.Primary.Attributes["vnet_name_resource_group2"] {
			return fmt.Errorf("vnet_name_resource_group2 Not found in created attributes")
		}

		*azurePeer = *foundPeer
		return nil
	}
}

func testAccCheckAzurePeerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_azure_peer" {
			continue
		}

		foundPeer := &goaviatrix.AzurePeer{
			VNet1: rs.Primary.Attributes["vnet_name_resource_group1"],
			VNet2: rs.Primary.Attributes["vnet_name_resource_group2"],
		}

		_, err := client.GetAzurePeer(foundPeer)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("azure peer still exists")
		}
	}

	return nil
}
