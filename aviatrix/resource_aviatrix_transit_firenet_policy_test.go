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

func TestAccAviatrixTransitFireNetPolicy_basic(t *testing.T) {
	var transitFireNetPolicy goaviatrix.TransitFireNetPolicy

	rName := fmt.Sprintf("%s", acctest.RandString(5))

	skipAcc := os.Getenv("SKIP_TRANSIT_FIRENET_POLICY")
	if skipAcc == "yes" {
		t.Skip("Skipping transit firenet policy tests as 'SKIP_TRANSIT_FIRENET_POLICY' is set")
	}

	skipAccAWS := os.Getenv("SKIP_TRANSIT_FIRENET_POLICY_AWS")
	skipAccARM := os.Getenv("SKIP_TRANSIT_FIRENET_POLICY_ARM")
	if skipAcc == "yes" && skipAccAWS == "yes" && skipAccARM == "yes" {
		t.Skip("Skipping transit firenet policy tests as 'SKIP_TRANSIT_FIRENET_POLICY_AWS' and 'SKIP_TRANSIT_FIRENET_POLICY_ARM' are all set")
	}

	if skipAccAWS != "yes" {
		resourceName := "aviatrix_transit_firenet_policy.test"
		msgCommonAws := ". Set 'SKIP_TRANSIT_FIRENET_POLICY_AWS' to 'yes' to skip transit firenet policy tests in AWS"
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheck(t, msgCommonAws)
				preGateway2Check(t, msgCommonAws)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckTransitFireNetPolicyDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccTransitFireNetPolicyConfigBasicAWS(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckTransitFireNetPolicyExists(resourceName, &transitFireNetPolicy),
						resource.TestCheckResourceAttr(resourceName, "transit_firenet_gateway_name", fmt.Sprintf("tfg-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "inspected_resource_name", fmt.Sprintf("SPOKE:tfg-aws-%s", rName)),
					),
				},
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	} else {
		t.Log("Skipping transit firenet policy tests in AWS as 'SKIP_TRANSIT_FIRENET_POLICY_AWS' is set")
	}

	if skipAccARM != "yes" {
		resourceName := "aviatrix_transit_firenet_policy.test"
		msgCommonARM := ". Set 'SKIP_TRANSIT_FIRENET_POLICY_ARM' to 'yes' to skip transit firenet policy tests in AZURE"
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheckARM(t, msgCommonARM)
				preGateway2CheckARM(t, msgCommonARM)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckTransitFireNetPolicyDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccTransitFireNetPolicyConfigBasicARM(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckTransitFireNetPolicyExists(resourceName, &transitFireNetPolicy),
						resource.TestCheckResourceAttr(resourceName, "transit_firenet_gateway_name", fmt.Sprintf("tfg-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "inspected_resource_name", fmt.Sprintf("SPOKE:tfg-arm-%s", rName)),
					),
				},
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	} else {
		t.Log("Skipping transit firenet policy tests in AZURE as 'SKIP_TRANSIT_FIRENET_POLICY_ARM' is set")
	}
}

func testAccTransitFireNetPolicyConfigBasicAWS(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_transit_gateway" "test_transit_gw" {
	cloud_type             = 1
	account_name           = aviatrix_account.test_account.account_name
	gw_name                = "tfg-%s"
	vpc_id                 = "%s"
	vpc_reg                = "%s"
	gw_size                = "c5.xlarge"
	subnet                 = "%s"
	enable_active_mesh     = true
	connected_transit      = true 
	enable_transit_firenet = true
}
resource "aviatrix_spoke_gateway" "test_spoke_gateway" {
	cloud_type         = 1
	account_name       = aviatrix_account.test_account.account_name
	gw_name            = "tfg-aws-%s"
	vpc_id             = "%s"
	vpc_reg            = "%s"
	gw_size            = "t2.micro"
	subnet             = "%s"
	enable_active_mesh = true
	transit_gw         = aviatrix_transit_gateway.test_transit_gw.gw_name
}
resource "aviatrix_transit_firenet_policy" "test" {
	transit_firenet_gateway_name = aviatrix_transit_gateway.test_transit_gw.gw_name
	inspected_resource_name      = join(":", ["SPOKE", aviatrix_spoke_gateway.test_spoke_gateway.gw_name])
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"), rName,
		os.Getenv("AWS_VPC_ID2"), os.Getenv("AWS_REGION2"), os.Getenv("AWS_SUBNET2"))
}

func testAccTransitFireNetPolicyConfigBasicARM(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_arm" {
	account_name        = "tfa-arm-%s"
	cloud_type          = 8
	arm_subscription_id = "%s"
	arm_directory_id    = "%s"
	arm_application_id  = "%s"
	arm_application_key = "%s"
}
resource "aviatrix_transit_gateway" "test_arm" {
	cloud_type             = 8
	account_name           = aviatrix_account.test_acc_arm.account_name
	gw_name                = "tfg-%s"
	vpc_id                 = "%s"
	vpc_reg                = "%s"
	gw_size                = "%s"
	subnet                 = "%s"
	enable_active_mesh     = true
	connected_transit      = true 
	enable_transit_firenet = true
}
resource "aviatrix_spoke_gateway" "test_spoke_gateway" {
	cloud_type         = 8
	account_name       = aviatrix_account.test_acc_arm.account_name
	gw_name            = "tfg-arm-%s"
	vpc_id             = "%s"
	vpc_reg            = "%s"
	gw_size            = "%s"
	subnet             = "%s"
	enable_active_mesh = true
	transit_gw         = aviatrix_transit_gateway.test_arm.gw_name
}
resource "aviatrix_transit_firenet_policy" "test" {
	transit_firenet_gateway_name = aviatrix_transit_gateway.test_arm.gw_name
	inspected_resource_name      = join(":", ["SPOKE", aviatrix_spoke_gateway.test_spoke_gateway.gw_name])
}
	`, rName, os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_DIRECTORY_ID"), os.Getenv("ARM_APPLICATION_ID"),
		os.Getenv("ARM_APPLICATION_KEY"), rName, os.Getenv("ARM_VNET_ID"), os.Getenv("ARM_REGION"),
		os.Getenv("ARM_GW_SIZE"), os.Getenv("ARM_SUBNET"), rName,
		os.Getenv("ARM_VNET_ID2"), os.Getenv("ARM_REGION2"),
		os.Getenv("ARM_GW_SIZE"), os.Getenv("ARM_SUBNET2"))
}

func testAccCheckTransitFireNetPolicyExists(n string, transitFireNetPolicy *goaviatrix.TransitFireNetPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("transit firenet policy Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no transit firenet policy ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundTransitFireNetPolicy := &goaviatrix.TransitFireNetPolicy{
			TransitFireNetGatewayName: rs.Primary.Attributes["transit_firenet_gateway_name"],
			InspectedResourceName:     rs.Primary.Attributes["inspected_resource_name"],
		}

		err := client.GetTransitFireNetPolicy(foundTransitFireNetPolicy)
		if err != nil {
			return fmt.Errorf("transit firenet policy not found")
		}

		*transitFireNetPolicy = *foundTransitFireNetPolicy
		return nil
	}
}

func testAccCheckTransitFireNetPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_transit_firenet_policy" {
			continue
		}

		foundTransitFireNetPolicy := &goaviatrix.TransitFireNetPolicy{
			TransitFireNetGatewayName: rs.Primary.Attributes["transit_firenet_gateway_name"],
			InspectedResourceName:     rs.Primary.Attributes["inspected_resource_name"],
		}

		err := client.GetTransitFireNetPolicy(foundTransitFireNetPolicy)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("transit firenet policy still exists")
		}
	}

	return nil
}
