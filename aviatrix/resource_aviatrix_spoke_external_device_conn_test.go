package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixSpokeExternalDeviceConn_basic(t *testing.T) {
	var externalDeviceConn goaviatrix.ExternalDeviceConn

	rName := acctest.RandString(5)
	resourceName := "aviatrix_spoke_external_device_conn.test"

	skipAcc := os.Getenv("SKIP_SPOKE_EXTERNAL_DEVICE_CONN")
	if skipAcc == "yes" {
		t.Skip("Skipping spoke external device connection tests as 'SKIP_SPOKE_EXTERNAL_DEVICE_CONN' is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, ". Set 'SKIP_SPOKE_EXTERNAL_DEVICE_CONN' to 'yes' to skip Site2Cloud spoke external device connection tests")
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpokeExternalDeviceConnDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpokeExternalDeviceConnConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpokeExternalDeviceConnExists(resourceName, &externalDeviceConn),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AWS_VPC_ID")),
					resource.TestCheckResourceAttr(resourceName, "connection_name", rName),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "connection_type", "bgp"),
					resource.TestCheckResourceAttr(resourceName, "bgp_local_as_num", "123"),
					resource.TestCheckResourceAttr(resourceName, "bgp_remote_as_num", "345"),
					resource.TestCheckResourceAttr(resourceName, "remote_gateway_ip", "172.12.13.14"),
					resource.TestCheckResourceAttr(resourceName, "enable_bfd", "true"),
					resource.TestCheckResourceAttr(resourceName, "bgp_bfd.0.transmit_interval", "400"),
					resource.TestCheckResourceAttr(resourceName, "bgp_bfd.0.receive_interval", "400"),
					resource.TestCheckResourceAttr(resourceName, "bgp_bfd.0.multiplier", "5"),
					resource.TestCheckResourceAttr(resourceName, "connection_bgp_send_communities", "444:444"),
					resource.TestCheckResourceAttr(resourceName, "connection_bgp_send_communities_additive", "true"),
					resource.TestCheckResourceAttr(resourceName, "connection_bgp_send_communities_block", "false"),
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

func testAccSpokeExternalDeviceConnConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_spoke_gateway" "test" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	gw_name      = "tfg-%s"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t2.micro"
	subnet       = "%s"
    enable_bgp   = true
}
resource "aviatrix_spoke_external_device_conn" "test" {
	vpc_id            = aviatrix_spoke_gateway.test.vpc_id
	connection_name   = "%s"
	gw_name           = aviatrix_spoke_gateway.test.gw_name
	connection_type   = "bgp"
	bgp_local_as_num  = "123"
	bgp_remote_as_num = "345"
	remote_gateway_ip = "172.12.13.14"
	enable_bfd = true
	bgp_bfd {
		transmit_interval = 400
		receive_interval = 400
		multiplier = 5
	}
	connection_bgp_send_communities           = "444:444"
	connection_bgp_send_communities_additive  = true
	connection_bgp_send_communities_block     = false
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"), rName)
}

func testAccCheckSpokeExternalDeviceConnExists(n string, externalDeviceConn *goaviatrix.ExternalDeviceConn) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("spoke external device connection Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no spoke external device connection ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundExternalDeviceConn := &goaviatrix.ExternalDeviceConn{
			VpcID:          rs.Primary.Attributes["vpc_id"],
			ConnectionName: rs.Primary.Attributes["connection_name"],
			GwName:         rs.Primary.Attributes["gw_name"],
		}
		localGateway, err := getGatewayDetails(client, foundExternalDeviceConn.GwName)
		if err != nil {
			return fmt.Errorf("could not get local gateway details: %w", err)
		}
		foundExternalDeviceConn2, err := client.GetExternalDeviceConnDetail(foundExternalDeviceConn, localGateway)
		if err != nil {
			return err
		}
		if foundExternalDeviceConn2.ConnectionName+"~"+foundExternalDeviceConn2.VpcID != rs.Primary.ID {
			return fmt.Errorf("spoke external device connection not found")
		}

		*externalDeviceConn = *foundExternalDeviceConn2
		return nil
	}
}

func testAccCheckSpokeExternalDeviceConnDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_spoke_external_device_conn" {
			continue
		}

		foundExternalDeviceConn := &goaviatrix.ExternalDeviceConn{
			VpcID:          rs.Primary.Attributes["vpc_id"],
			ConnectionName: rs.Primary.Attributes["connection_name"],
			GwName:         rs.Primary.Attributes["gw_name"],
		}
		localGateway, err := getGatewayDetails(client, foundExternalDeviceConn.GwName)
		if err != nil {
			return fmt.Errorf("could not get local gateway details: %w", err)
		}

		_, err = client.GetExternalDeviceConnDetail(foundExternalDeviceConn, localGateway)
		if err != goaviatrix.ErrNotFound {
			return fmt.Errorf("site2cloud still exists: %w", err)
		}
	}

	return nil
}

func TestSpokeExternalDeviceConnSchema_RemoteLanIPv6Fields(t *testing.T) {
	resource := resourceAviatrixSpokeExternalDeviceConn()
	schemaMap := resource.Schema

	// Test remote_lan_ipv6_ip field exists and has correct properties
	remoteLanIPv6Field, ok := schemaMap["remote_lan_ipv6_ip"]
	assert.True(t, ok, "remote_lan_ipv6_ip field should exist in schema")
	assert.Equal(t, schema.TypeString, remoteLanIPv6Field.Type)
	assert.True(t, remoteLanIPv6Field.Optional, "remote_lan_ipv6_ip should be optional")
	assert.True(t, remoteLanIPv6Field.Computed, "remote_lan_ipv6_ip should be computed")
	assert.True(t, remoteLanIPv6Field.ForceNew, "remote_lan_ipv6_ip should be ForceNew")
	assert.Contains(t, remoteLanIPv6Field.Description, "Remote LAN IPv6 address")

	// Test backup_remote_lan_ipv6_ip field exists and has correct properties
	backupRemoteLanIPv6Field, ok := schemaMap["backup_remote_lan_ipv6_ip"]
	assert.True(t, ok, "backup_remote_lan_ipv6_ip field should exist in schema")
	assert.Equal(t, schema.TypeString, backupRemoteLanIPv6Field.Type)
	assert.True(t, backupRemoteLanIPv6Field.Optional, "backup_remote_lan_ipv6_ip should be optional")
	assert.True(t, backupRemoteLanIPv6Field.Computed, "backup_remote_lan_ipv6_ip should be computed")
	assert.True(t, backupRemoteLanIPv6Field.ForceNew, "backup_remote_lan_ipv6_ip should be ForceNew")
	assert.Contains(t, backupRemoteLanIPv6Field.Description, "Backup Remote LAN IPv6 address")
}

func TestSpokeExternalDeviceConnSchema_RemoteLanIPv6FieldsReference(t *testing.T) {
	// Test that remote_lan_ipv6_ip follows the same pattern as remote_lan_ip
	resource := resourceAviatrixSpokeExternalDeviceConn()
	schemaMap := resource.Schema

	remoteLanIPField, ok := schemaMap["remote_lan_ip"]
	assert.True(t, ok, "remote_lan_ip field should exist for reference")

	remoteLanIPv6Field, ok := schemaMap["remote_lan_ipv6_ip"]
	assert.True(t, ok, "remote_lan_ipv6_ip field should exist")

	// Both should be TypeString
	assert.Equal(t, remoteLanIPField.Type, remoteLanIPv6Field.Type, "remote_lan_ipv6_ip should have same type as remote_lan_ip")
	assert.Equal(t, remoteLanIPField.ForceNew, remoteLanIPv6Field.ForceNew, "remote_lan_ipv6_ip should have same ForceNew as remote_lan_ip")

	// Test backup fields follow same pattern
	backupRemoteLanIPField, ok := schemaMap["backup_remote_lan_ip"]
	assert.True(t, ok, "backup_remote_lan_ip field should exist for reference")

	backupRemoteLanIPv6Field, ok := schemaMap["backup_remote_lan_ipv6_ip"]
	assert.True(t, ok, "backup_remote_lan_ipv6_ip field should exist")

	assert.Equal(t, backupRemoteLanIPField.Type, backupRemoteLanIPv6Field.Type, "backup_remote_lan_ipv6_ip should have same type as backup_remote_lan_ip")
	assert.Equal(t, backupRemoteLanIPField.ForceNew, backupRemoteLanIPv6Field.ForceNew, "backup_remote_lan_ipv6_ip should have same ForceNew as backup_remote_lan_ip")
}
