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

func preARMPeerCheck(t *testing.T, msgCommon string) (string, string, string, string) {
	vNet1 := os.Getenv("ARM_VNET_ID")
	if vNet1 == "" {
		t.Fatal("Environment variable ARM_VNET_ID is not set" + msgCommon)
	}
	vNet2 := os.Getenv("ARM_VNET_ID2")
	if vNet2 == "" {
		t.Fatal("Environment variable ARM_VNET_ID2 is not set" + msgCommon)
	}

	region1 := os.Getenv("ARM_REGION")
	if region1 == "" {
		t.Fatal("Environment variable ARM_REGION is not set" + msgCommon)
	}
	region2 := os.Getenv("ARM_REGION2")
	if region2 == "" {
		t.Fatal("Environment variable ARM_REGION2 is not set" + msgCommon)
	}
	return vNet1, vNet2, region1, region2
}

func TestAccAviatrixARMPeer_basic(t *testing.T) {
	var armPeer goaviatrix.ARMPeer
	rInt := acctest.RandInt()
	resourceName := "aviatrix_arm_peer.test_arm_peer"

	skipAcc := os.Getenv("SKIP_ARM_PEER")
	if skipAcc == "yes" {
		t.Skip("Skipping aviatrix ARM peering test as SKIP_ARM_PEER is set")
	}
	msgCommon := ". Set SKIP_ARM_PEER to yes to skip AWS peer tests"

	preAccountCheck(t, msgCommon)

	vNet1, vNet2, region1, region2 := preARMPeerCheck(t, msgCommon)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckARMPeerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccARMPeerConfigBasic(rInt, vNet1, vNet2, region1, region2),
				Check: resource.ComposeTestCheckFunc(
					tesAccCheckARMPeerExists(resourceName, &armPeer),
					resource.TestCheckResourceAttr(
						resourceName, "vnet_name_resource_group1", vNet1),
					resource.TestCheckResourceAttr(
						resourceName, "vnet_name_resource_group2", vNet2),
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

func testAccARMPeerConfigBasic(rInt int, vNet1 string, vNet2 string, region1 string, region2 string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name        = "tf-testing-%d"
	cloud_type          = 8
    arm_subscription_id = "%s"
    arm_directory_id    = "%s"
    arm_application_id  = "%s"
    arm_application_key = "%s"
}

resource "aviatrix_arm_peer" "test_arm_peer" {
	account_name1             = "${aviatrix_account.test_account.account_name}"
	account_name2             = "${aviatrix_account.test_account.account_name}"
	vnet_name_resource_group1 = "%s"
	vnet_name_resource_group2 = "%s"
	vnet_reg1                 = "%s"
	vnet_reg2                 = "%s"
}
	`, rInt, os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_DIRECTORY_ID"),
		os.Getenv("ARM_APPLICATION_ID"), os.Getenv("ARM_APPLICATION_KEY"),
		vNet1, vNet2, region1, region2)
}

func tesAccCheckARMPeerExists(n string, armPeer *goaviatrix.ARMPeer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("armPeer Not Created: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ARMPeer ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundPeer := &goaviatrix.ARMPeer{
			VNet1: rs.Primary.Attributes["vnet_name_resource_group1"],
			VNet2: rs.Primary.Attributes["vnet_name_resource_group2"],
		}

		_, err := client.GetARMPeer(foundPeer)
		if err != nil {
			return err
		}
		if foundPeer.VNet1 != rs.Primary.Attributes["vnet_name_resource_group1"] {
			return fmt.Errorf("vnet_name_resource_group1 Not found in created attributes")
		}
		if foundPeer.VNet2 != rs.Primary.Attributes["vnet_name_resource_group2"] {
			return fmt.Errorf("vnet_name_resource_group2 Not found in created attributes")
		}

		*armPeer = *foundPeer
		return nil
	}
}

func testAccCheckARMPeerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_arm_peer" {
			continue
		}
		foundPeer := &goaviatrix.ARMPeer{
			VNet1: rs.Primary.Attributes["vnet_name_resource_group1"],
			VNet2: rs.Primary.Attributes["vnet_name_resource_group2"],
		}
		_, err := client.GetARMPeer(foundPeer)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("armPeer still exists")
		}
	}
	return nil
}
