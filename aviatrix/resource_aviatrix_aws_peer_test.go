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

func preAWSPeerCheck(t *testing.T, msgCommon string) (string, string, string, string) {
	vpcID1 := os.Getenv("AWS_VPC_ID")
	if vpcID1 == "" {
		t.Fatal("Environment variable AWS_VPC_ID is not set" + msgCommon)
	}
	vpcID2 := os.Getenv("AWS_VPC_ID2")
	if vpcID2 == "" {
		t.Fatal("Environment variable AWS_VPC_ID2 is not set" + msgCommon)
	}

	region1 := os.Getenv("AWS_REGION")
	if region1 == "" {
		t.Fatal("Environment variable AWS_REGION is not set" + msgCommon)
	}
	region2 := os.Getenv("AWS_REGION2")
	if region2 == "" {
		t.Fatal("Environment variable AWS_REGION2 is not set" + msgCommon)
	}
	return vpcID1, vpcID2, region1, region2
}

func TestAccAviatrixAWSPeer_basic(t *testing.T) {
	var awsPeer goaviatrix.AWSPeer
	rInt := acctest.RandInt()
	resourceName := "aviatrix_aws_peer.test_aws_peer"

	skipAcc := os.Getenv("SKIP_AWS_PEER")
	if skipAcc == "yes" {
		t.Skip("Skipping aviatrix AWS peering test as SKIP_AWS_PEER is set")
	}
	msgCommon := ". Set SKIP_AWS_PEER to yes to skip AWS peer tests"

	preAccountCheck(t, msgCommon)

	vpcID1, vpcID2, region1, region2 := preAWSPeerCheck(t, msgCommon)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSPeerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSPeerConfigBasic(rInt, vpcID1, vpcID2, region1, region2),
				Check: resource.ComposeTestCheckFunc(
					tesAccCheckAWSPeerExists(resourceName, &awsPeer),
					resource.TestCheckResourceAttr(
						resourceName, "vpc_id1", vpcID1),
					resource.TestCheckResourceAttr(
						resourceName, "vpc_id2", vpcID2),
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

func testAccAWSPeerConfigBasic(rInt int, vpcID1 string, vpcID2 string, region1 string, region2 string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name       = "tf-testing-%d"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = "false"
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_aws_peer" "test_aws_peer" {
	account_name1 = "${aviatrix_account.test_account.account_name}"
	account_name2 = "${aviatrix_account.test_account.account_name}"
	vpc_id1       = "%s"
	vpc_id2       = "%s"
	vpc_reg1      = "%s"
	vpc_reg2      = "%s"
}
	`, rInt, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		vpcID1, vpcID2, region1, region2)
}

func tesAccCheckAWSPeerExists(n string, awsPeer *goaviatrix.AWSPeer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("awsPeer Not Created: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no AWSPeer ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundPeer := &goaviatrix.AWSPeer{
			VpcID1: rs.Primary.Attributes["vpc_id1"],
			VpcID2: rs.Primary.Attributes["vpc_id2"],
		}

		_, err := client.GetAWSPeer(foundPeer)
		if err != nil {
			return err
		}
		if foundPeer.VpcID1 != rs.Primary.Attributes["vpc_id1"] {
			return fmt.Errorf("vpc_id1 Not found in created attributes")
		}
		if foundPeer.VpcID2 != rs.Primary.Attributes["vpc_id2"] {
			return fmt.Errorf("vpc_id2 Not found in created attributes")
		}
		*awsPeer = *foundPeer

		return nil
	}
}

func testAccCheckAWSPeerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_aws_peer" {
			continue
		}
		foundPeer := &goaviatrix.AWSPeer{
			VpcID1: rs.Primary.Attributes["vpc_id1"],
			VpcID2: rs.Primary.Attributes["vpc_id2"],
		}
		_, err := client.GetAWSPeer(foundPeer)

		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("awsPeer still exists")
		}
	}
	return nil
}
