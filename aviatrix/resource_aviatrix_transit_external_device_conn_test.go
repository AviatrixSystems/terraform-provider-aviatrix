package aviatrix

import (
	"errors"
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

func TestAccAviatrixTransitExternalDeviceConn_basic(t *testing.T) {
	var externalDeviceConn goaviatrix.ExternalDeviceConn

	rName := acctest.RandString(5)
	resourceName := "aviatrix_transit_external_device_conn.test"

	skipAcc := os.Getenv("SKIP_TRANSIT_EXTERNAL_DEVICE_CONN")
	if skipAcc == "yes" {
		t.Skip("Skipping transit external device connection tests as 'SKIP_TRANSIT_EXTERNAL_DEVICE_CONN' is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, ". Set 'SKIP_TRANSIT_EXTERNAL_DEVICE_CONN' to 'yes' to skip Site2Cloud transit external device connection tests")
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitExternalDeviceConnDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitExternalDeviceConnConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitExternalDeviceConnExists(resourceName, &externalDeviceConn),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AWS_VPC_ID")),
					resource.TestCheckResourceAttr(resourceName, "connection_name", rName),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "connection_type", "bgp"),
					resource.TestCheckResourceAttr(resourceName, "bgp_local_as_num", "123"),
					resource.TestCheckResourceAttr(resourceName, "bgp_remote_as_num", "345"),
					resource.TestCheckResourceAttr(resourceName, "remote_gateway_ip", "172.12.13.14"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

	skipBgpBfd := os.Getenv("SKIP_BGP_BFD_DEVICE_CONN")
	if skipBgpBfd != "yes" {
		resourceNameBfd := "aviatrix_transit_external_device_conn.transit-1-transit-2"
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheck(t, ". Set 'SKIP_BGP_BFD_DEVICE_CONN' to 'yes' to skip Site2Cloud transit external device connection tests")
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckTransitExternalDeviceConnDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccTransitExternalDeviceConnConfigBgpBfd(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckTransitExternalDeviceConnExists(resourceNameBfd, &externalDeviceConn),
						resource.TestCheckResourceAttr(resourceNameBfd, "vpc_id", os.Getenv("AWS_VPC_ID")),
						resource.TestCheckResourceAttr(resourceNameBfd, "connection_name", "transit-1-transit-2"),
						resource.TestCheckResourceAttr(resourceNameBfd, "gw_name", fmt.Sprintf("%s-aws-transit-1", rName)),
						resource.TestCheckResourceAttr(resourceNameBfd, "connection_type", "bgp"),
						resource.TestCheckResourceAttr(resourceNameBfd, "tunnel_protocol", "IPsec"),
						resource.TestCheckResourceAttr(resourceNameBfd, "bgp_local_as_num", "65075"),
						resource.TestCheckResourceAttr(resourceNameBfd, "bgp_remote_as_num", "65076"),
						resource.TestCheckResourceAttr(resourceNameBfd, "enable_bfd", "true"),
						resource.TestCheckResourceAttr(resourceNameBfd, "bgp_bfd.0.transmit_interval", "400"),
						resource.TestCheckResourceAttr(resourceNameBfd, "bgp_bfd.0.receive_interval", "400"),
						resource.TestCheckResourceAttr(resourceNameBfd, "bgp_bfd.0.multiplier", "5"),
					),
				},
				{
					Config: testAccTransitExternalDeviceConnConfigBgpBfdUpdated(),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckTransitExternalDeviceConnExists(resourceNameBfd, &externalDeviceConn),
						resource.TestCheckResourceAttr(resourceNameBfd, "bgp_bfd.0.transmit_interval", "500"),
						resource.TestCheckResourceAttr(resourceNameBfd, "bgp_bfd.0.receive_interval", "500"),
						resource.TestCheckResourceAttr(resourceNameBfd, "bgp_bfd.0.multiplier", "7"),
					),
				},
				{
					Config: testAccTransitExternalDeviceConnConfigBfdDisabled(),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckTransitExternalDeviceConnExists(resourceNameBfd, &externalDeviceConn),
						resource.TestCheckResourceAttr(resourceNameBfd, "enable_bfd", "false"),
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
		t.Skip("Skipping transit external device connection tests for bgp bfd config as 'SKIP_BGP_BFD_DEVICE_CONN' is set")
	}
}

func TestAccAviatrixTransitExternalDeviceConn_proxyId(t *testing.T) {
	var externalDeviceConn goaviatrix.ExternalDeviceConn

	rName := acctest.RandString(5)
	resourceName := "aviatrix_transit_external_device_conn.test_proxy_id"

	skipAcc := os.Getenv("SKIP_TRANSIT_EXTERNAL_DEVICE_CONN")
	if skipAcc == "yes" {
		t.Skip("Skipping transit external device connection tests as 'SKIP_TRANSIT_EXTERNAL_DEVICE_CONN' is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, ". Set 'SKIP_TRANSIT_EXTERNAL_DEVICE_CONN' to 'yes' to skip Site2Cloud transit external device connection tests")
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitExternalDeviceConnDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitExternalDeviceConnConfigProxyId(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitExternalDeviceConnExists(resourceName, &externalDeviceConn),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AWS_VPC_ID")),
					resource.TestCheckResourceAttr(resourceName, "connection_name", fmt.Sprintf("%s-proxy-id", rName)),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "connection_type", "static"),
					resource.TestCheckResourceAttr(resourceName, "remote_gateway_ip", "172.12.13.15"),
					resource.TestCheckResourceAttr(resourceName, "remote_subnet", "192.168.1.0/24,192.168.2.0/24"),
					resource.TestCheckResourceAttr(resourceName, "local_subnet", "10.0.0.0/16,10.1.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "proxy_id_enabled", "true"),
				),
			},
			{
				Config: testAccTransitExternalDeviceConnConfigProxyIdUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitExternalDeviceConnExists(resourceName, &externalDeviceConn),
					resource.TestCheckResourceAttr(resourceName, "local_subnet", "10.0.0.0/16,10.1.0.0/16,10.2.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "proxy_id_enabled", "false"),
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

func testAccTransitExternalDeviceConnConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_transit_gateway" "test" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	gw_name      = "tfg-%s"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t2.micro"
	subnet       = "%s"
}
resource "aviatrix_transit_external_device_conn" "test" {
	vpc_id            = aviatrix_transit_gateway.test.vpc_id
	connection_name   = "%s"
	gw_name           = aviatrix_transit_gateway.test.gw_name
	connection_type   = "bgp"
	bgp_local_as_num  = "123"
	bgp_remote_as_num = "345"
	remote_gateway_ip = "172.12.13.14"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"), rName)
}

func TestTransitExternalDeviceConnSchema_RemoteLanIPv6Fields(t *testing.T) {
	resource := resourceAviatrixTransitExternalDeviceConn()
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

func TestTransitExternalDeviceConnSchema_RemoteLanIPv6FieldsReference(t *testing.T) {
	// Test that remote_lan_ipv6_ip follows the same pattern as remote_lan_ip
	resource := resourceAviatrixTransitExternalDeviceConn()
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

func TestTransitExternalDeviceConnSchema_EnableIPv6ForceNew(t *testing.T) {
	resource := resourceAviatrixTransitExternalDeviceConn()
	schemaMap := resource.Schema

	// Test enable_ipv6 field exists and has correct properties
	enableIPv6Field, ok := schemaMap["enable_ipv6"]
	assert.True(t, ok, "enable_ipv6 field should exist in schema")
	assert.Equal(t, schema.TypeBool, enableIPv6Field.Type)
	assert.True(t, enableIPv6Field.Optional, "enable_ipv6 should be optional")
	assert.True(t, enableIPv6Field.ForceNew, "enable_ipv6 should be ForceNew")
	assert.Contains(t, enableIPv6Field.Description, "Enable IPv6")
}

func TestTransitExternalDeviceConnSchema_LocalSubnetField(t *testing.T) {
	resource := resourceAviatrixTransitExternalDeviceConn()
	schemaMap := resource.Schema

	// Test local_subnet field exists and has correct properties
	localSubnetField, ok := schemaMap["local_subnet"]
	assert.True(t, ok, "local_subnet field should exist in schema")
	assert.Equal(t, schema.TypeString, localSubnetField.Type, "local_subnet should be TypeString")
	assert.True(t, localSubnetField.Optional, "local_subnet should be optional")
	assert.NotNil(t, localSubnetField.DiffSuppressFunc, "local_subnet should have DiffSuppressFunc")
	assert.Contains(t, localSubnetField.Description, "Local CIDRs", "local_subnet description should mention Local CIDRs")
}

func TestTransitExternalDeviceConnSchema_ProxyIdEnabledField(t *testing.T) {
	resource := resourceAviatrixTransitExternalDeviceConn()
	schemaMap := resource.Schema

	// Test proxy_id_enabled field exists and has correct properties
	proxyIdEnabledField, ok := schemaMap["proxy_id_enabled"]
	assert.True(t, ok, "proxy_id_enabled field should exist in schema")
	assert.Equal(t, schema.TypeBool, proxyIdEnabledField.Type, "proxy_id_enabled should be TypeBool")
	assert.True(t, proxyIdEnabledField.Optional, "proxy_id_enabled should be optional")
	assert.Equal(t, false, proxyIdEnabledField.Default, "proxy_id_enabled should have default value false")
	assert.Contains(t, proxyIdEnabledField.Description, "proxy ID", "proxy_id_enabled description should mention proxy ID")
}

func TestTransitExternalDeviceConnSchema_LocalSubnetAndProxyIdComparison(t *testing.T) {
	// Test that local_subnet and proxy_id_enabled follow similar patterns to remote_subnet
	resource := resourceAviatrixTransitExternalDeviceConn()
	schemaMap := resource.Schema

	remoteSubnetField, ok := schemaMap["remote_subnet"]
	assert.True(t, ok, "remote_subnet field should exist for reference")

	localSubnetField, ok := schemaMap["local_subnet"]
	assert.True(t, ok, "local_subnet field should exist")

	// Both subnet fields should be TypeString and Optional
	assert.Equal(t, remoteSubnetField.Type, localSubnetField.Type, "local_subnet should have same type as remote_subnet")
	assert.Equal(t, remoteSubnetField.Optional, localSubnetField.Optional, "local_subnet should have same Optional as remote_subnet")

	// Both should have DiffSuppressFunc (for handling comma-separated CIDRs with spaces)
	assert.NotNil(t, remoteSubnetField.DiffSuppressFunc, "remote_subnet should have DiffSuppressFunc")
	assert.NotNil(t, localSubnetField.DiffSuppressFunc, "local_subnet should have DiffSuppressFunc")

	// Test proxy_id_enabled exists and is boolean
	proxyIdEnabledField, ok := schemaMap["proxy_id_enabled"]
	assert.True(t, ok, "proxy_id_enabled field should exist")
	assert.Equal(t, schema.TypeBool, proxyIdEnabledField.Type, "proxy_id_enabled should be TypeBool")
}

func TestTransitExternalDeviceConnSchema_EditableAttrsAreNotForceNew(t *testing.T) {
	assertEdcEditableAttrsNotForceNew(t, resourceAviatrixTransitExternalDeviceConn().Schema,
		"aviatrix_transit_external_device_conn")
}

// TestAccAviatrixTransitExternalDeviceConn_editInPlace changes every attribute
// the edit API can change and asserts the connection was edited rather than
// destroyed and recreated. Before this feature each of these was ForceNew, so
// the same plan replaced the connection and dropped the BGP session.
func TestAccAviatrixTransitExternalDeviceConn_editInPlace(t *testing.T) {
	var externalDeviceConn goaviatrix.ExternalDeviceConn

	if os.Getenv("SKIP_TRANSIT_EXTERNAL_DEVICE_CONN_EDIT") == "yes" {
		t.Skip("Skipping the transit external device connection edit test as " +
			"'SKIP_TRANSIT_EXTERNAL_DEVICE_CONN_EDIT' is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_transit_external_device_conn.test"
	var connID string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, ". Set 'SKIP_TRANSIT_EXTERNAL_DEVICE_CONN_EDIT' to 'yes' to skip the "+
				"Site2Cloud transit external device connection edit test")
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"aviatrix": func() (*schema.Provider, error) { return Provider(), nil },
		},
		CheckDestroy: testAccCheckTransitExternalDeviceConnDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitExternalDeviceConnConfigEditable(rName, "345",
					"169.254.10.1/30", "169.254.10.2/30", "172.12.13.14", "psk-before"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitExternalDeviceConnExists(resourceName, &externalDeviceConn),
					testAccCheckEdcIDUnchanged(resourceName, &connID),
					resource.TestCheckResourceAttr(resourceName, "bgp_remote_as_num", "345"),
					resource.TestCheckResourceAttr(resourceName, "remote_tunnel_cidr", "169.254.10.2/30"),
					resource.TestCheckResourceAttr(resourceName, "remote_gateway_ip", "172.12.13.14"),
				),
			},
			{
				// The remote ASN, both inner addresses, the pre-shared key and the
				// outer peer IP all change at once, which is also the case that
				// exercises ordering: the remote_device_ip edit re-keys the tunnel,
				// so it has to run after the address edits.
				Config: testAccTransitExternalDeviceConnConfigEditable(rName, "346",
					"169.254.11.1/30", "169.254.11.2/30", "172.12.13.15", "psk-after"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitExternalDeviceConnExists(resourceName, &externalDeviceConn),
					testAccCheckEdcIDUnchanged(resourceName, &connID),
					resource.TestCheckResourceAttr(resourceName, "bgp_remote_as_num", "346"),
					resource.TestCheckResourceAttr(resourceName, "local_tunnel_cidr", "169.254.11.1/30"),
					resource.TestCheckResourceAttr(resourceName, "remote_tunnel_cidr", "169.254.11.2/30"),
					resource.TestCheckResourceAttr(resourceName, "remote_gateway_ip", "172.12.13.15"),
				),
			},
		},
	})
}

// testAccTransitExternalDeviceConnConfigEditable takes the editable values as
// arguments so both steps of the edit test share one config, which keeps the
// diff between them to exactly the attributes under test.
func testAccTransitExternalDeviceConnConfigEditable(rName, remoteASN, localCidr, remoteCidr, remoteGatewayIP, psk string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_transit_gateway" "test" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	gw_name      = "tfg-%s"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t2.micro"
	subnet       = "%s"
}
resource "aviatrix_transit_external_device_conn" "test" {
	vpc_id             = aviatrix_transit_gateway.test.vpc_id
	connection_name    = "%s"
	gw_name            = aviatrix_transit_gateway.test.gw_name
	connection_type    = "bgp"
	tunnel_protocol    = "IPsec"
	bgp_local_as_num   = "123"
	bgp_remote_as_num  = "%s"
	local_tunnel_cidr  = "%s"
	remote_tunnel_cidr = "%s"
	remote_gateway_ip  = "%s"
	pre_shared_key     = "%s"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"), rName,
		remoteASN, localCidr, remoteCidr, remoteGatewayIP, psk)
}

// TestAccAviatrixTransitExternalDeviceConn_bgpAdminShutdown toggles the
// connection-wide BGP admin state and asserts the connection is edited rather
// than replaced each time.
//
// The first step is the one worth having: it creates the connection with
// bgp_admin_shutdown already true, which the create API cannot express. Getting
// true back out of a plan that reports no further changes is what demonstrates
// the post-create edit ran - if it had been skipped, the connection would exist
// with BGP up and the next plan would want to change it.
func TestAccAviatrixTransitExternalDeviceConn_bgpAdminShutdown(t *testing.T) {
	var externalDeviceConn goaviatrix.ExternalDeviceConn

	if os.Getenv("SKIP_TRANSIT_EXTERNAL_DEVICE_CONN_EDIT") == "yes" {
		t.Skip("Skipping the transit external device connection BGP admin shutdown test as " +
			"'SKIP_TRANSIT_EXTERNAL_DEVICE_CONN_EDIT' is set")
	}

	rName := acctest.RandString(5)
	resourceName := "aviatrix_transit_external_device_conn.test"
	var connID string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, ". Set 'SKIP_TRANSIT_EXTERNAL_DEVICE_CONN_EDIT' to 'yes' to skip the "+
				"Site2Cloud transit external device connection BGP admin shutdown test")
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"aviatrix": func() (*schema.Provider, error) { return Provider(), nil },
		},
		CheckDestroy: testAccCheckTransitExternalDeviceConnDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitExternalDeviceConnConfigBgpAdminShutdown(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitExternalDeviceConnExists(resourceName, &externalDeviceConn),
					testAccCheckEdcIDUnchanged(resourceName, &connID),
					resource.TestCheckResourceAttr(resourceName, "bgp_admin_shutdown", "true"),
				),
			},
			{
				// Lifting it reaches the other branch of the ordering rule.
				Config: testAccTransitExternalDeviceConnConfigBgpAdminShutdown(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitExternalDeviceConnExists(resourceName, &externalDeviceConn),
					testAccCheckEdcIDUnchanged(resourceName, &connID),
					resource.TestCheckResourceAttr(resourceName, "bgp_admin_shutdown", "false"),
				),
			},
			{
				// And back down, so the test covers a shutdown applied to a
				// connection that has been up as well as to a fresh one.
				Config: testAccTransitExternalDeviceConnConfigBgpAdminShutdown(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitExternalDeviceConnExists(resourceName, &externalDeviceConn),
					testAccCheckEdcIDUnchanged(resourceName, &connID),
					resource.TestCheckResourceAttr(resourceName, "bgp_admin_shutdown", "true"),
				),
			},
		},
	})
}

// testAccTransitExternalDeviceConnConfigBgpAdminShutdown varies only
// bgp_admin_shutdown, so that each step's diff is the one attribute under test
// and no tunnel edit rides along with it.
func testAccTransitExternalDeviceConnConfigBgpAdminShutdown(rName string, shutdown bool) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_transit_gateway" "test" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	gw_name      = "tfg-%s"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t2.micro"
	subnet       = "%s"
}
resource "aviatrix_transit_external_device_conn" "test" {
	vpc_id             = aviatrix_transit_gateway.test.vpc_id
	connection_name    = "%s"
	gw_name            = aviatrix_transit_gateway.test.gw_name
	connection_type    = "bgp"
	tunnel_protocol    = "IPsec"
	bgp_local_as_num   = "123"
	bgp_remote_as_num  = "345"
	local_tunnel_cidr  = "169.254.10.1/30"
	remote_tunnel_cidr = "169.254.10.2/30"
	remote_gateway_ip  = "172.12.13.14"
	pre_shared_key     = "psk-shutdown"
	bgp_admin_shutdown = %t
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"), rName, shutdown)
}

func testAccTransitExternalDeviceConnConfigBgpBfd(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_transit_gateway" "transit_gateway_1" {
	single_az_ha = true
	gw_name = "%s-transit_gateway_1"
	vpc_id = "%s"
	cloud_type   = 1
	vpc_reg = "%s"
    gw_size = "t3.medium"
    account_name = "%s"
    subnet = "%s"
    connected_transit = true
    bgp_manual_spoke_advertise_cidrs = "192.0.2.0/24"
    bgp_polling_time = "10"
    local_as_number = "65075"
}
resource "aviatrix_transit_gateway" "transit_gateway_2" {
    single_az_ha = true
    gw_name = "%s-aws-transit-2"
    vpc_id = "%s"
    cloud_type = 1
    vpc_reg = "%s"
    gw_size = "t3.medium"
    account_name = "%s"
    subnet = "%s"
    connected_transit = true
    bgp_manual_spoke_advertise_cidrs = "198.51.100.0/24"
    bgp_polling_time = "10"
    local_as_number = "65076"
}
resource "aviatrix_transit_external_device_conn" "transit-1-transit-2" {
	vpc_id            = aviatrix_transit_gateway.transit_gateway_1.vpc_id
	connection_name   = "transit-1-transit-2"
	gw_name           = aviatrix_transit_gateway.transit_gateway_1.gw_name
	connection_type   = "bgp"
	tunnel_protocol    = "IPsec"
	bgp_local_as_num  = aviatrix_transit_gateway.transit_gateway_1.local_as_number
	bgp_remote_as_num = aviatrix_transit_gateway.transit_gateway_2.local_as_number
	remote_gateway_ip = aviatrix_transit_gateway.transit_gateway_2.eip
	local_tunnel_cidr  = "169.254.10.1/30"
	remote_tunnel_cidr = "169.254.10.2/30"
	pre_shared_key = "psk12"
	enable_bfd = true
	bgp_bfd {
	  transmit_interval = 400
	  receive_interval = 400
	  multiplier = 5
	}
  }
  resource "aviatrix_transit_external_device_conn" "transit-2-transit-1" {
	vpc_id            = avaitrix_transit_gateway.transit_gateway_2.vpc_id
	connection_name   = "transit-2-transit-1"
	gw_name           = aviatrix_transit_gateway.transit_gateway_2.gw_name
	connection_type   = "bgp"
	tunnel_protocol    = "IPsec"
	bgp_local_as_num  = aviatrix_transit_gateway.transit_gateway_2.local_as_number
	bgp_remote_as_num = aviatrix_transit_gateway.transit_gateway_1.local_as_number
	remote_gateway_ip = "3.140.172.87"
	local_tunnel_cidr  = "169.254.10.2/30"
	remote_tunnel_cidr = "169.254.10.1/30"
	pre_shared_key = "psk12"
  }
	`, rName, os.Getenv("AWS_VPC_ID_1"), os.Getenv("AWS_REGION"), os.Getenv("AWS_ACCOUNT_NAME"), os.Getenv("AWS_SUBNET_1"),
		rName, os.Getenv("AWS_VPC_ID_2"), os.Getenv("AWS_REGION"), os.Getenv("AWS_ACCOUNT_NAME"), os.Getenv("AWS_SUBNET_2"))
}

func testAccTransitExternalDeviceConnConfigBgpBfdUpdated() string {
	return `
	resource "aviatrix_transit_external_device_conn" "transit-1-transit-2" {
		vpc_id            = aviatrix_transit_gateway.transit_gateway_1.vpc_id
		connection_name   = "transit-1-transit-2"
		gw_name           = aviatrix_transit_gateway.transit_gateway_1.gw_name
		connection_type   = "bgp"
		tunnel_protocol    = "IPsec"
		bgp_local_as_num  = aviatrix_transit_gateway.transit_gateway_1.local_as_number
		bgp_remote_as_num = aviatrix_transit_gateway.transit_gateway_2.local_as_number
		remote_gateway_ip = aviatrix_transit_gateway.transit_gateway_2.eip
		local_tunnel_cidr  = "169.254.10.1/30"
		remote_tunnel_cidr = "169.254.10.2/30"
		pre_shared_key = "psk12"
		enable_bfd = true
		bgp_bfd {
		  transmit_interval = 500
		  receive_interval = 500
		  multiplier = 7
		}
	}`
}

func testAccTransitExternalDeviceConnConfigBfdDisabled() string {
	return `
	resource "aviatrix_transit_external_device_conn" "transit-1-transit-2" {
		vpc_id            = aviatrix_transit_gateway.transit_gateway_1.vpc_id
		connection_name   = "transit-1-transit-2"
		gw_name           = aviatrix_transit_gateway.transit_gateway_1.gw_name
		connection_type   = "bgp"
		tunnel_protocol    = "IPsec"
		bgp_local_as_num  = aviatrix_transit_gateway.transit_gateway_1.local_as_number
		bgp_remote_as_num = aviatrix_transit_gateway.transit_gateway_2.local_as_number
		remote_gateway_ip = aviatrix_transit_gateway.transit_gateway_2.eip
		local_tunnel_cidr  = "169.254.10.1/30"
		remote_tunnel_cidr = "169.254.10.2/30"
		pre_shared_key = "psk12"
		enable_bfd = false
	}`
}

func testAccCheckTransitExternalDeviceConnExists(n string, externalDeviceConn *goaviatrix.ExternalDeviceConn) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("transit external device connection Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no transit external device connection ID is set")
		}

		client := mustClient(testAccProvider.Meta())

		foundExternalDeviceConn := &goaviatrix.ExternalDeviceConn{
			VpcID:          rs.Primary.Attributes["vpc_id"],
			ConnectionName: rs.Primary.Attributes["connection_name"],
			GwName:         rs.Primary.Attributes["gw_name"],
		}
		localGateway, err := getGatewayDetails(client, foundExternalDeviceConn.GwName)
		if err != nil {
			return fmt.Errorf("could not get local gateway details: %w", err)
		}
		foundExternalDeviceConn2, err := client.GetExternalDeviceConnDetail(foundExternalDeviceConn, localGateway, false)
		if err != nil {
			return err
		}
		if foundExternalDeviceConn2.ConnectionName+"~"+foundExternalDeviceConn2.VpcID != rs.Primary.ID {
			return fmt.Errorf("transit external device connection not found")
		}

		*externalDeviceConn = *foundExternalDeviceConn2
		return nil
	}
}

func testAccCheckTransitExternalDeviceConnDestroy(s *terraform.State) error {
	client := mustClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_transit_external_device_conn" {
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

		_, err = client.GetExternalDeviceConnDetail(foundExternalDeviceConn, localGateway, false)
		if !errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("site2cloud still exists %s", err.Error())
		}
	}

	return nil
}

func TestAccAviatrixEdgeTransitExternalDeviceConn(t *testing.T) {
	var externalDeviceConn goaviatrix.ExternalDeviceConn

	rName := "aviatrix_transit_external_device_conn.eat-2-bgpoipsec-1"

	skipAcc := os.Getenv("SKIP_EDGE_TRANSIT_EXTERNAL_DEVICE_CONN")
	if skipAcc == "yes" {
		t.Skip("Skipping edge transit external device connection tests as 'SKIP_EDGE_TRANSIT_EXTERNAL_DEVICE_CONN' is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, ". Set 'SKIP_EDGE_TRANSIT_EXTERNAL_DEVICE_CONN' to 'yes' to skip Edge transit external device connection tests")
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitExternalDeviceConnDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEdgeTransitExternalDeviceConnConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitExternalDeviceConnExists(rName, &externalDeviceConn),
					resource.TestCheckResourceAttr(rName, "connection_name", "eat-2-bgpoipsec-1"),
					resource.TestCheckResourceAttr(rName, "gw_name", "e2e-edge-transit-2"),
					resource.TestCheckResourceAttr(rName, "enable_jumbo_frame", "true"),
					resource.TestCheckResourceAttr(rName, "tunnel_src_ip", "192.168.20.117,192.168.23.16"),
					resource.TestCheckResourceAttr(rName, "connection_bgp_send_communities", "444:444"),
					resource.TestCheckResourceAttr(rName, "connection_bgp_send_communities_additive", "true"),
					resource.TestCheckResourceAttr(rName, "connection_bgp_send_communities_block", "false"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEdgeTransitExternalDeviceConnConfig() string {
	return `
	resource "aviatrix_transit_external_device_conn" "eat-2-bgpoipsec-1" {
	vpc_id                    = "tsite-2"
	connection_name           = "eat-2-bgpoipsec-1"
	gw_name                   = "e2e-edge-transit-2"
	connection_type           = "bgp"
	tunnel_protocol           = "IPsec"
	bgp_local_as_num          = "65402"
	bgp_remote_as_num         = "65151"
	remote_gateway_ip         = "3.140.40.45"
	pre_shared_key         = "aviatrix,aviatrix"
	direct_connect            = true
	enable_jumbo_frame            = true
	disable_activemesh = false
	ha_enabled                = true
	local_tunnel_cidr         = "169.254.22.54/30, 169.254.238.22/30"
	remote_tunnel_cidr        = "169.254.22.53/30, 169.254.238.21/30"
	backup_local_tunnel_cidr  = "169.254.33.94/30, 169.254.165.254/30"
	backup_remote_tunnel_cidr = "169.254.33.93/30, 169.254.165.253/30"
	backup_bgp_remote_as_num  = "65151"
	backup_remote_gateway_ip  = "18.223.219.22"
	backup_direct_connect     = true
	backup_pre_shared_key         = "aviatrix,aviatrix"
	tunnel_src_ip  = "192.168.20.117,192.168.23.16"
	connection_bgp_send_communities           = "444:444"
	connection_bgp_send_communities_additive  = true
	connection_bgp_send_communities_block     = false
	}`
}

func TestAccAviatrixTransitExternalDeviceConn_disableActivemeshNotSetInState(t *testing.T) {
	var externalDeviceConn goaviatrix.ExternalDeviceConn

	rName := acctest.RandString(5)
	resourceName := "aviatrix_transit_external_device_conn.test_disable_activemesh_not_set"

	skipAcc := os.Getenv("SKIP_TRANSIT_EXTERNAL_DEVICE_CONN")
	if skipAcc == "yes" {
		t.Skip("Skipping transit external device connection tests as 'SKIP_TRANSIT_EXTERNAL_DEVICE_CONN' is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGatewayCheck(t, ". Set 'SKIP_TRANSIT_EXTERNAL_DEVICE_CONN' to 'yes' to skip Site2Cloud transit external device connection tests")
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitExternalDeviceConnDestroy,
		Steps: []resource.TestStep{
			{
				// Step 1: Create without disable_activemesh (should use default value but not set in state)
				Config: testAccTransitExternalDeviceConnConfigDisableActivemeshNotSet(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitExternalDeviceConnExists(resourceName, &externalDeviceConn),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AWS_VPC_ID")),
					resource.TestCheckResourceAttr(resourceName, "connection_name", rName),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "connection_type", "bgp"),
					resource.TestCheckResourceAttr(resourceName, "bgp_local_as_num", "123"),
					resource.TestCheckResourceAttr(resourceName, "bgp_remote_as_num", "345"),
					resource.TestCheckResourceAttr(resourceName, "remote_gateway_ip", "172.12.13.14"),
					// Verify that disable_activemesh is NOT set in state even though default is false
					resource.TestCheckNoResourceAttr(resourceName, "disable_activemesh"),
				),
			},
			{
				// Step 2: Update to explicitly set disable_activemesh = false (should now be in state)
				Config: testAccTransitExternalDeviceConnConfigDisableActivemeshExplicitFalse(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitExternalDeviceConnExists(resourceName, &externalDeviceConn),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AWS_VPC_ID")),
					resource.TestCheckResourceAttr(resourceName, "connection_name", rName),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "connection_type", "bgp"),
					resource.TestCheckResourceAttr(resourceName, "bgp_local_as_num", "123"),
					resource.TestCheckResourceAttr(resourceName, "bgp_remote_as_num", "345"),
					resource.TestCheckResourceAttr(resourceName, "remote_gateway_ip", "172.12.13.14"),
					// Verify that disable_activemesh IS now set in state because user explicitly set it
					resource.TestCheckResourceAttr(resourceName, "disable_activemesh", "false"),
				),
			},
			{
				// Step 3: Update to explicitly set disable_activemesh = true (should still be in state)
				Config: testAccTransitExternalDeviceConnConfigDisableActivemeshExplicitTrue(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitExternalDeviceConnExists(resourceName, &externalDeviceConn),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AWS_VPC_ID")),
					resource.TestCheckResourceAttr(resourceName, "connection_name", rName),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "connection_type", "bgp"),
					resource.TestCheckResourceAttr(resourceName, "bgp_local_as_num", "123"),
					resource.TestCheckResourceAttr(resourceName, "bgp_remote_as_num", "345"),
					resource.TestCheckResourceAttr(resourceName, "remote_gateway_ip", "172.12.13.14"),
					// Verify that disable_activemesh IS set in state because user explicitly set it
					resource.TestCheckResourceAttr(resourceName, "disable_activemesh", "true"),
				),
			},
			{
				// Step 4: Remove disable_activemesh again (should remove from state)
				Config: testAccTransitExternalDeviceConnConfigDisableActivemeshNotSet(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitExternalDeviceConnExists(resourceName, &externalDeviceConn),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AWS_VPC_ID")),
					resource.TestCheckResourceAttr(resourceName, "connection_name", rName),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "connection_type", "bgp"),
					resource.TestCheckResourceAttr(resourceName, "bgp_local_as_num", "123"),
					resource.TestCheckResourceAttr(resourceName, "bgp_remote_as_num", "345"),
					resource.TestCheckResourceAttr(resourceName, "remote_gateway_ip", "172.12.13.14"),
					// Verify that disable_activemesh is NOT set in state again
					resource.TestCheckNoResourceAttr(resourceName, "disable_activemesh"),
				),
			},
		},
	})
}

func testAccTransitExternalDeviceConnConfigDisableActivemeshNotSet(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_transit_gateway" "test" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	gw_name      = "tfg-%s"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t3.medium"
	subnet       = "%s"
}

resource "aviatrix_transit_external_device_conn" "test_disable_activemesh_not_set" {
	vpc_id              = aviatrix_transit_gateway.test.vpc_id
	connection_name     = "%s"
	gw_name             = aviatrix_transit_gateway.test.gw_name
	connection_type     = "bgp"
	bgp_local_as_num    = "123"
	bgp_remote_as_num   = "345"
	remote_gateway_ip   = "172.12.13.14"
	# disable_activemesh not set - should use default value but not appear in state
}
`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"), rName)
}

func testAccTransitExternalDeviceConnConfigDisableActivemeshExplicitFalse(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_transit_gateway" "test" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	gw_name      = "tfg-%s"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t3.medium"
	subnet       = "%s"
}

resource "aviatrix_transit_external_device_conn" "test_disable_activemesh_not_set" {
	vpc_id              = aviatrix_transit_gateway.test.vpc_id
	connection_name     = "%s"
	gw_name             = aviatrix_transit_gateway.test.gw_name
	connection_type     = "bgp"
	bgp_local_as_num    = "123"
	bgp_remote_as_num   = "345"
	remote_gateway_ip   = "172.12.13.14"
	disable_activemesh  = false  # Explicitly set to false - should appear in state
}
`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"), rName)
}

func testAccTransitExternalDeviceConnConfigDisableActivemeshExplicitTrue(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}

resource "aviatrix_transit_gateway" "test" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	gw_name      = "tfg-%s"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t3.medium"
	subnet       = "%s"
}

resource "aviatrix_transit_external_device_conn" "test_disable_activemesh_not_set" {
	vpc_id              = aviatrix_transit_gateway.test.vpc_id
	connection_name     = "%s"
	gw_name             = aviatrix_transit_gateway.test.gw_name
	connection_type     = "bgp"
	bgp_local_as_num    = "123"
	bgp_remote_as_num   = "345"
	remote_gateway_ip   = "172.12.13.14"
	disable_activemesh  = true  # Explicitly set to true - should appear in state
}
`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"), rName)
}

func testAccTransitExternalDeviceConnConfigProxyId(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_transit_gateway" "test" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	gw_name      = "tfg-%s"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t2.micro"
	subnet       = "%s"
}
resource "aviatrix_transit_external_device_conn" "test_proxy_id" {
	vpc_id            = aviatrix_transit_gateway.test.vpc_id
	connection_name   = "%s-proxy-id"
	gw_name           = aviatrix_transit_gateway.test.gw_name
	connection_type   = "static"
	remote_gateway_ip = "172.12.13.15"
	remote_subnet     = "192.168.1.0/24,192.168.2.0/24"
	local_subnet      = "10.0.0.0/16,10.1.0.0/16"
	proxy_id_enabled  = true
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"), rName)
}

func testAccTransitExternalDeviceConnConfigProxyIdUpdated(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_transit_gateway" "test" {
	cloud_type   = 1
	account_name = aviatrix_account.test.account_name
	gw_name      = "tfg-%s"
	vpc_id       = "%s"
	vpc_reg      = "%s"
	gw_size      = "t2.micro"
	subnet       = "%s"
}
resource "aviatrix_transit_external_device_conn" "test_proxy_id" {
	vpc_id            = aviatrix_transit_gateway.test.vpc_id
	connection_name   = "%s-proxy-id"
	gw_name           = aviatrix_transit_gateway.test.gw_name
	connection_type   = "static"
	remote_gateway_ip = "172.12.13.15"
	remote_subnet     = "192.168.1.0/24,192.168.2.0/24"
	local_subnet      = "10.0.0.0/16,10.1.0.0/16,10.2.0.0/16"
	proxy_id_enabled  = false
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		rName, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"), rName)
}

func TestParseExternalDeviceConnImportID(t *testing.T) {
	tests := []struct {
		name             string
		id               string
		expectedConnName string
		expectedVpcID    string
		expectErr        bool
	}{
		{
			name:             "valid id",
			id:               "conn-1~vpc-abc123",
			expectedConnName: "conn-1",
			expectedVpcID:    "vpc-abc123",
		},
		{
			name:             "vpc id containing tilde is preserved",
			id:               "conn-1~vpc~with~tilde",
			expectedConnName: "conn-1",
			expectedVpcID:    "vpc~with~tilde",
		},
		{
			name:      "missing separator",
			id:        "conn-1",
			expectErr: true,
		},
		{
			name:      "empty id",
			id:        "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			connName, vpcID, err := parseExternalDeviceConnImportID(tt.id)
			if tt.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedConnName, connName)
			assert.Equal(t, tt.expectedVpcID, vpcID)
		})
	}
}

// externalDeviceConnTestSchema returns the minimal schema needed to build a
// *schema.ResourceData for exercising handleExternalDeviceConnImport.
func externalDeviceConnTestSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"connection_name": {Type: schema.TypeString, Optional: true},
		"vpc_id":          {Type: schema.TypeString, Optional: true},
		"gw_name":         {Type: schema.TypeString, Optional: true},
	}
}

func TestHandleExternalDeviceConnImport(t *testing.T) {
	t.Run("import with malformed id returns error", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, externalDeviceConnTestSchema(), map[string]any{})
		d.SetId("no-separator")

		err := handleExternalDeviceConnImport(d, nil)

		assert.Error(t, err)
	})

	t.Run("import with gw_name already set does not probe", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, externalDeviceConnTestSchema(), map[string]any{
			"gw_name": "gw-preset",
		})
		d.SetId("conn-1~vpc-abc123")

		// gw_name is already known, so the client probe is skipped and nil is
		// safe.
		err := handleExternalDeviceConnImport(d, nil)

		assert.NoError(t, err)
		assert.Equal(t, "conn-1", getString(d, "connection_name"))
		assert.Equal(t, "vpc-abc123", getString(d, "vpc_id"))
		assert.Equal(t, "gw-preset", getString(d, "gw_name"))
		assert.Equal(t, "conn-1~vpc-abc123", d.Id())
	})
}
