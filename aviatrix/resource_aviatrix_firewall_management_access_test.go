package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixFirewallManagementAccess_basic(t *testing.T) {
	var firewallManagementAccess goaviatrix.FirewallManagementAccess

	rName := fmt.Sprintf("%s", acctest.RandString(5))

	skipAcc := os.Getenv("SKIP_FIREWALL_MANAGEMENT_ACCESS")
	if skipAcc == "yes" {
		t.Skip("Skipping firewall management access tests as 'SKIP_FIREWALL_MANAGEMENT_ACCESS' is set")
	}

	skipAccAWS := os.Getenv("SKIP_FIREWALL_MANAGEMENT_ACCESS_AWS")
	skipAccAZURE := os.Getenv("SKIP_FIREWALL_MANAGEMENT_ACCESS_AZURE")
	if skipAcc == "yes" && skipAccAWS == "yes" && skipAccAZURE == "yes" {
		t.Skip("Skipping firewall management access tests as 'SKIP_FIREWALL_MANAGEMENT_ACCESS_AWS' and 'SKIP_FIREWALL_MANAGEMENT_ACCESS_AZURE' are all set")
	}

	if skipAccAWS != "yes" {
		resourceName := "aviatrix_firewall_management_access.test"
		msgCommonAws := ". Set 'SKIP_FIREWALL_MANAGEMENT_ACCESS_AWS' to 'yes' to skip firewall management access tests in AWS"
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheck(t, msgCommonAws)
				preGateway2Check(t, msgCommonAws)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckFirewallManagementAccessDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccFirewallManagementAccessConfigBasicAWS(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckFirewallManagementAccessExists(resourceName, &firewallManagementAccess),
						resource.TestCheckResourceAttr(resourceName, "transit_firenet_gateway_name", fmt.Sprintf("tfg-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "management_access_resource_name", fmt.Sprintf("SPOKE:tfg-aws-%s", rName)),
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

	if skipAccAZURE != "yes" {
		resourceName := "aviatrix_firewall_management_access.test"
		msgCommonAZURE := ". Set 'SKIP_FIREWALL_MANAGEMENT_ACCESS_AZURE' to 'yes' to skip firewall management access tests in AZURE"
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheckAZURE(t, msgCommonAZURE)
				preGateway2CheckAZURE(t, msgCommonAZURE)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckFirewallManagementAccessDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccFirewallManagementAccessConfigBasicAZURE(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckFirewallManagementAccessExists(resourceName, &firewallManagementAccess),
						resource.TestCheckResourceAttr(resourceName, "transit_firenet_gateway_name", fmt.Sprintf("tfg-%s", rName)),
						resource.TestCheckResourceAttr(resourceName, "management_access_resource_name", fmt.Sprintf("SPOKE:tfg-azure-%s", rName)),
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
		t.Log("Skipping transit firenet policy tests in AZURE as 'SKIP_TRANSIT_FIRENET_POLICY_AZURE' is set")
	}
}

func testAccFirewallManagementAccessConfigBasicAWS(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_transit_gateway" "test_transit_aws" {
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
resource "aviatrix_spoke_gateway" "test_spoke_aws" {
	cloud_type         = 1
	account_name       = aviatrix_account.test_account.account_name
	gw_name            = "tfg-aws-%s"
	vpc_id             = "%s"
	vpc_reg            = "%s"
	gw_size            = "t2.micro"
	subnet             = "%s"
	enable_active_mesh = true
	transit_gw         = aviatrix_transit_gateway.test_transit_aws.gw_name
}
resource "aviatrix_firewall_management_access" "test" {
	transit_firenet_gateway_name    = aviatrix_transit_gateway.test_transit_aws.gw_name
	management_access_resource_name = join(":", ["SPOKE", aviatrix_spoke_gateway.test_spoke_aws.gw_name])
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"), rName,
		os.Getenv("AWS_VPC_ID2"), os.Getenv("AWS_REGION2"), os.Getenv("AWS_SUBNET2"))
}

func testAccFirewallManagementAccessConfigBasicAZURE(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_azure" {
	account_name        = "tfa-%s"
	cloud_type          = 8
	arm_subscription_id = "%s"
	arm_directory_id    = "%s"
	arm_application_id  = "%s"
	arm_application_key = "%s"
}
resource "aviatrix_transit_gateway" "test_transit_azure" {
	cloud_type             = 8
	account_name           = aviatrix_account.test_acc_azure.account_name
	gw_name                = "tfg-%s"
	vpc_id                 = "%s"
	vpc_reg                = "%s"
	gw_size                = "%s"
	subnet                 = "%s"
	enable_active_mesh     = true
	connected_transit      = true 
	enable_transit_firenet = true
}
resource "aviatrix_spoke_gateway" "test_spoke_azure" {
	cloud_type         = 8
	account_name       = aviatrix_account.test_acc_azure.account_name
	gw_name            = "tfg-azure-%s"
	vpc_id             = "%s"
	vpc_reg            = "%s"
	gw_size            = "%s"
	subnet             = "%s"
	enable_active_mesh = true
	transit_gw         = aviatrix_transit_gateway.test_transit_azure.gw_name
}
resource "aviatrix_firewall_management_access" "test" {
	transit_firenet_gateway_name    = aviatrix_transit_gateway.test_transit_azure.gw_name
	management_access_resource_name = join(":", ["SPOKE", aviatrix_spoke_gateway.test_spoke_azure.gw_name])
}
	`, rName, os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_DIRECTORY_ID"), os.Getenv("ARM_APPLICATION_ID"),
		os.Getenv("ARM_APPLICATION_KEY"), rName, os.Getenv("AZURE_VNET_ID"), os.Getenv("AZURE_REGION"),
		os.Getenv("AZURE_GW_SIZE"), os.Getenv("AZURE_SUBNET"), rName,
		os.Getenv("AZURE_VNET_ID2"), os.Getenv("AZURE_REGION2"),
		os.Getenv("AZURE_GW_SIZE"), os.Getenv("AZURE_SUBNET2"))
}

func testAccCheckFirewallManagementAccessExists(n string, firewallManagementAccess *goaviatrix.FirewallManagementAccess) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("firewall management access Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no firewall management access ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundFirewallManagementAccess := &goaviatrix.FirewallManagementAccess{
			TransitFireNetGatewayName: rs.Primary.Attributes["transit_firenet_gateway_name"],
		}

		foundFirewallManagementAccess2, err := client.GetFirewallManagementAccess(foundFirewallManagementAccess)
		if err != nil {
			return fmt.Errorf("firewall management access not found")
		}

		*firewallManagementAccess = *foundFirewallManagementAccess2
		return nil
	}
}

func testAccCheckFirewallManagementAccessDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_firewall_management_access" {
			continue
		}

		foundFirewallManagementAccess := &goaviatrix.FirewallManagementAccess{
			TransitFireNetGatewayName: rs.Primary.Attributes["transit_firenet_gateway_name"],
		}

		_, err := client.GetFirewallManagementAccess(foundFirewallManagementAccess)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("firewall management access still exists")
		}
	}

	return nil
}
