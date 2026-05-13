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

func TestAccAviatrixEdgeSpokeExternalDeviceConn_basic(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "aviatrix_edge_spoke_external_device_conn.test"

	skipAcc := os.Getenv("SKIP_EDGE_SPOKE_EXTERNAL_DEVICE_CONN")
	if skipAcc == "yes" {
		t.Skip("Skipping Edge as a Spoke external device connection tests as 'SKIP_EDGE_SPOKE_EXTERNAL_DEVICE_CONN' is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preEdgeSpokeExternalDeviceConnCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEdgeSpokeExternalDeviceConnDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEdgeSpokeExternalDeviceConnConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEdgeSpokeExternalDeviceConnExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site_id", os.Getenv("EDGE_SPOKE_SITE_ID")),
					resource.TestCheckResourceAttr(resourceName, "connection_name", rName),
					resource.TestCheckResourceAttr(resourceName, "gw_name", os.Getenv("EDGE_SPOKE_NAME")),
					resource.TestCheckResourceAttr(resourceName, "bgp_local_as_num", "65001"),
					resource.TestCheckResourceAttr(resourceName, "bgp_remote_as_num", "65002"),
					resource.TestCheckResourceAttr(resourceName, "local_lan_ip", "1.2.3.4"),
					resource.TestCheckResourceAttr(resourceName, "remote_lan_ip", "5.6.7.8"),
					resource.TestCheckResourceAttr(resourceName, "connection_type", "bgp"),
					resource.TestCheckResourceAttr(resourceName, "enable_bgp_lan_activemesh", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_jumbo_frame", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

	skipBgpBfd := os.Getenv("SKIP_BGP_BFD_EDGE_SPOKE_EXTERNAL_DEVICE_CONN")
	if skipBgpBfd != "yes" {
		resourceNameBfd := "aviatrix_edge_spoke_external_device_conn.test-bfd"
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preEdgeSpokeExternalDeviceConnCheck(t)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckEdgeSpokeExternalDeviceConnDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccEdgeSpokeExternalDeviceConnConfigBgpBfd(),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckEdgeSpokeExternalDeviceConnExists(resourceNameBfd),
						resource.TestCheckResourceAttr(resourceNameBfd, "site_id", os.Getenv("EDGE_SPOKE_SITE_ID")),
						resource.TestCheckResourceAttr(resourceNameBfd, "connection_name", "cloudn-bgp-lan"),
						resource.TestCheckResourceAttr(resourceNameBfd, "gw_name", os.Getenv("EDGE_SPOKE_NAME")),
						resource.TestCheckResourceAttr(resourceNameBfd, "bgp_local_as_num", "65182"),
						resource.TestCheckResourceAttr(resourceNameBfd, "bgp_remote_as_num", "65220"),
						resource.TestCheckResourceAttr(resourceNameBfd, "local_lan_ip", "10.220.86.182"),
						resource.TestCheckResourceAttr(resourceNameBfd, "remote_lan_ip", "10.220.86.100"),
						resource.TestCheckResourceAttr(resourceNameBfd, "enable_bfd", "true"),
						resource.TestCheckResourceAttr(resourceNameBfd, "bgp_bfd.0.transmit_interval", "400"),
						resource.TestCheckResourceAttr(resourceNameBfd, "bgp_bfd.0.receive_interval", "400"),
						resource.TestCheckResourceAttr(resourceNameBfd, "bgp_bfd.0.multiplier", "5"),
					),
				},
				{
					ResourceName:      resourceNameBfd,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	} else {
		t.Skip("Skipping BGP BFD Edge as a Spoke external device connection tests as 'SKIP_BGP_BFD_EDGE_SPOKE_EXTERNAL_DEVICE_CONN' is set")
	}
}

func preEdgeSpokeExternalDeviceConnCheck(t *testing.T) {
	if os.Getenv("EDGE_SPOKE_NAME") == "" {
		t.Fatal("Environment variable EDGE_SPOKE_NAME is not set")
	}

	if os.Getenv("EDGE_SPOKE_SITE_ID") == "" {
		t.Fatal("Environment variable EDGE_SPOKE_SITE_ID is not set")
	}
}

func testAccEdgeSpokeExternalDeviceConnConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_edge_spoke_external_device_conn" "test" {
	site_id           = "%s"
	connection_name   = "%s"
	gw_name           = "%s"
	bgp_local_as_num  = "65001"
	bgp_remote_as_num = "65002"
	local_lan_ip      = "1.2.3.4"
	remote_lan_ip     = "5.6.7.8"
	connection_type   = "bgp"
	enable_bgp_lan_activemesh = true
	enable_jumbo_frame = true
}
	`, os.Getenv("EDGE_SPOKE_SITE_ID"), rName, os.Getenv("EDGE_SPOKE_NAME"))
}

func testAccEdgeSpokeExternalDeviceConnConfigBgpBfd() string {
	return fmt.Sprintf(`
resource "aviatrix_edge_spoke_external_device_conn" "test-bfd" {
	site_id           = "%s"
	connection_name   = "cloudn-bgp-lan"
	gw_name           = "%s"
	bgp_local_as_num  = "65182"
	bgp_remote_as_num = "65220"
	local_lan_ip      = "10.220.86.182"
	remote_lan_ip     = "10.220.86.100"
	enable_bfd        = true
	bgp_bfd {
		transmit_interval = 400
		receive_interval = 400
		multiplier = 5
	}
}
	`, os.Getenv("EDGE_SPOKE_SITE_ID"), os.Getenv("EDGE_SPOKE_NAME"))
}

func testAccCheckEdgeSpokeExternalDeviceConnExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("edge as a spoke external device conn not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no edge as a spoke external device conn ID is set")
		}

		client := mustClient(testAccProvider.Meta())

		externalDeviceConn := &goaviatrix.ExternalDeviceConn{
			VpcID:          rs.Primary.Attributes["site_id"],
			ConnectionName: rs.Primary.Attributes["connection_name"],
			GwName:         rs.Primary.Attributes["gw_name"],
		}

		localGateway, err := getGatewayDetails(client, externalDeviceConn.GwName)
		if err != nil {
			return fmt.Errorf("could not get local gateway details: %w", err)
		}

		conn, err := client.GetExternalDeviceConnDetail(externalDeviceConn, localGateway)
		if err != nil {
			return err
		}

		expectedID := conn.ConnectionName + "~" + conn.VpcID + "~" + conn.GwName
		if expectedID != rs.Primary.ID {
			return fmt.Errorf("edge as a spoke external device conn ID mismatch: expected %s, got %s", expectedID, rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckEdgeSpokeExternalDeviceConnDestroy(s *terraform.State) error {
	client := mustClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_edge_spoke_external_device_conn" {
			continue
		}

		externalDeviceConn := &goaviatrix.ExternalDeviceConn{
			VpcID:          rs.Primary.Attributes["site_id"],
			ConnectionName: rs.Primary.Attributes["connection_name"],
			GwName:         rs.Primary.Attributes["gw_name"],
		}

		localGateway, err := getGatewayDetails(client, externalDeviceConn.GwName)
		if err != nil {
			return fmt.Errorf("could not get local gateway details: %w", err)
		}

		_, err = client.GetExternalDeviceConnDetail(externalDeviceConn, localGateway)
		if !errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("edge as a spoke external device conn still exists %s", err.Error())
		}
	}

	return nil
}

func TestEdgeSpokeExternalDeviceConnSchema_RemoteLanIPv6Fields(t *testing.T) {
	resource := resourceAviatrixEdgeSpokeExternalDeviceConn()
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
	assert.Contains(t, backupRemoteLanIPv6Field.Description, "Backup Remote LAN IPv6 address")
}

func TestEdgeSpokeExternalDeviceConnSchema_RemoteLanIPv6FieldsReference(t *testing.T) {
	// Test that remote_lan_ipv6_ip follows the same pattern as remote_lan_ip
	resource := resourceAviatrixEdgeSpokeExternalDeviceConn()
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

func TestEdgeSpokeExternalDeviceConnSchema_EnableIPv6ForceNew(t *testing.T) {
	resource := resourceAviatrixEdgeSpokeExternalDeviceConn()
	schemaMap := resource.Schema

	// Test enable_ipv6 field exists and has correct properties
	enableIPv6Field, ok := schemaMap["enable_ipv6"]
	assert.True(t, ok, "enable_ipv6 field should exist in schema")
	assert.Equal(t, schema.TypeBool, enableIPv6Field.Type)
	assert.True(t, enableIPv6Field.Optional, "enable_ipv6 should be optional")
	assert.True(t, enableIPv6Field.ForceNew, "enable_ipv6 should be ForceNew")
	assert.Contains(t, enableIPv6Field.Description, "Enable IPv6")
}

func TestBuildEdgeSpokeExternalDeviceConnForHa(t *testing.T) {
	resource := resourceAviatrixEdgeSpokeExternalDeviceConn()

	tests := []struct {
		name        string
		config      map[string]interface{}
		haGwName    string
		expectError bool
		errorSubstr string
		validate    func(t *testing.T, conn *goaviatrix.ExternalDeviceConn)
	}{
		{
			name: "Valid HA config builds conn successfully",
			config: map[string]interface{}{
				"site_id":                  "site-123",
				"connection_name":          "conn-test",
				"gw_name":                  "edge-gw",
				"bgp_local_as_num":         "65001",
				"backup_bgp_remote_as_num": "65002",
				"backup_local_lan_ip":      "10.1.1.3",
				"backup_remote_lan_ip":     "10.1.1.4",
				"connection_type":          "bgp",
				"tunnel_protocol":          "LAN",
				"enable_edge_underlay":     false,
			},
			haGwName:    "edge-gw-hagw",
			expectError: false,
			validate: func(t *testing.T, conn *goaviatrix.ExternalDeviceConn) {
				assert.Equal(t, "site-123", conn.VpcID)
				assert.Equal(t, "conn-test", conn.ConnectionName)
				assert.Equal(t, "edge-gw-hagw", conn.GwName)
				assert.Equal(t, "10.1.1.3", conn.LocalLanIP)
				assert.Equal(t, "10.1.1.4", conn.RemoteLanIP)
				assert.Equal(t, 65001, conn.BgpLocalAsNum)
				assert.Equal(t, 65002, conn.BgpRemoteAsNum)
			},
		},
		{
			name: "Valid HA config with IPv6",
			config: map[string]interface{}{
				"site_id":                   "site-123",
				"connection_name":           "conn-test",
				"gw_name":                   "edge-gw",
				"bgp_local_as_num":          "65001",
				"backup_bgp_remote_as_num":  "65002",
				"backup_local_lan_ip":       "10.1.1.3",
				"backup_remote_lan_ip":      "10.1.1.4",
				"backup_remote_lan_ipv6_ip": "2001:db8::2",
				"connection_type":           "bgp",
				"tunnel_protocol":           "LAN",
				"enable_edge_underlay":      false,
				"enable_ipv6":               true,
			},
			haGwName:    "edge-gw-hagw",
			expectError: false,
			validate: func(t *testing.T, conn *goaviatrix.ExternalDeviceConn) {
				assert.True(t, conn.EnableIpv6)
				assert.Equal(t, "2001:db8::2", conn.RemoteLanIPv6)
			},
		},
		{
			// backup_bgp_remote_as_num is now validated before backup_remote_lan_ip
			name: "Fails when backup_bgp_remote_as_num is missing",
			config: map[string]interface{}{
				"site_id":                  "site-123",
				"connection_name":          "conn-test",
				"gw_name":                  "edge-gw",
				"bgp_local_as_num":         "65001",
				"backup_bgp_remote_as_num": "",
				"backup_local_lan_ip":      "10.1.1.3",
				"backup_remote_lan_ip":     "10.1.1.4",
				"connection_type":          "bgp",
				"tunnel_protocol":          "LAN",
			},
			haGwName:    "edge-gw-hagw",
			expectError: true,
			errorSubstr: "backup_bgp_remote_as_num",
		},
		{
			name: "Fails when backup_bgp_remote_as_num is not a valid integer",
			config: map[string]interface{}{
				"site_id":                  "site-123",
				"connection_name":          "conn-test",
				"gw_name":                  "edge-gw",
				"bgp_local_as_num":         "65001",
				"backup_bgp_remote_as_num": "not-a-number",
				"backup_local_lan_ip":      "10.1.1.3",
				"backup_remote_lan_ip":     "10.1.1.4",
				"connection_type":          "bgp",
				"tunnel_protocol":          "LAN",
			},
			haGwName:    "edge-gw-hagw",
			expectError: true,
			errorSubstr: "backup_bgp_remote_as_num",
		},
		{
			name: "Fails when backup_remote_lan_ip is missing",
			config: map[string]interface{}{
				"site_id":                  "site-123",
				"connection_name":          "conn-test",
				"gw_name":                  "edge-gw",
				"bgp_local_as_num":         "65001",
				"backup_bgp_remote_as_num": "65002",
				"backup_local_lan_ip":      "10.1.1.3",
				"backup_remote_lan_ip":     "",
				"connection_type":          "bgp",
				"tunnel_protocol":          "LAN",
			},
			haGwName:    "edge-gw-hagw",
			expectError: true,
			errorSubstr: "backup_remote_lan_ip",
		},
		{
			name: "Fails when bgp_local_as_num is not a valid integer",
			config: map[string]interface{}{
				"site_id":                  "site-123",
				"connection_name":          "conn-test",
				"gw_name":                  "edge-gw",
				"bgp_local_as_num":         "not-a-number",
				"backup_bgp_remote_as_num": "65002",
				"backup_local_lan_ip":      "10.1.1.3",
				"backup_remote_lan_ip":     "10.1.1.4",
				"connection_type":          "bgp",
				"tunnel_protocol":          "LAN",
			},
			haGwName:    "edge-gw-hagw",
			expectError: true,
			errorSubstr: "bgp_local_as_num",
		},
		{
			name: "Succeeds when bgp_local_as_num is empty",
			config: map[string]interface{}{
				"site_id":                  "site-123",
				"connection_name":          "conn-test",
				"gw_name":                  "edge-gw",
				"bgp_local_as_num":         "",
				"backup_bgp_remote_as_num": "65002",
				"backup_local_lan_ip":      "10.1.1.3",
				"backup_remote_lan_ip":     "10.1.1.4",
				"connection_type":          "bgp",
				"tunnel_protocol":          "LAN",
			},
			haGwName:    "edge-gw-hagw",
			expectError: false,
			validate: func(t *testing.T, conn *goaviatrix.ExternalDeviceConn) {
				assert.Equal(t, 0, conn.BgpLocalAsNum)
			},
		},
		{
			name: "Fails when enable_ipv6 but backup_remote_lan_ipv6_ip missing",
			config: map[string]interface{}{
				"site_id":                  "site-123",
				"connection_name":          "conn-test",
				"gw_name":                  "edge-gw",
				"bgp_local_as_num":         "65001",
				"backup_bgp_remote_as_num": "65002",
				"backup_local_lan_ip":      "10.1.1.3",
				"backup_remote_lan_ip":     "10.1.1.4",
				"connection_type":          "bgp",
				"tunnel_protocol":          "LAN",
				"enable_ipv6":              true,
			},
			haGwName:    "edge-gw-hagw",
			expectError: true,
			errorSubstr: "backup_remote_lan_ipv6_ip",
		},
		{
			name: "Includes backup_bgp_md5_key when set",
			config: map[string]interface{}{
				"site_id":                  "site-123",
				"connection_name":          "conn-test",
				"gw_name":                  "edge-gw",
				"bgp_local_as_num":         "65001",
				"backup_bgp_remote_as_num": "65002",
				"backup_local_lan_ip":      "10.1.1.3",
				"backup_remote_lan_ip":     "10.1.1.4",
				"backup_bgp_md5_key":       "secret-key",
				"connection_type":          "bgp",
				"tunnel_protocol":          "LAN",
			},
			haGwName:    "edge-gw-hagw",
			expectError: false,
			validate: func(t *testing.T, conn *goaviatrix.ExternalDeviceConn) {
				assert.Equal(t, "secret-key", conn.BgpMd5Key)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resource.Schema, tt.config)
			conn, err := buildEdgeSpokeExternalDeviceConnForHa(d, tt.haGwName)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorSubstr != "" {
					assert.Contains(t, err.Error(), tt.errorSubstr)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, conn)
			if tt.validate != nil {
				tt.validate(t, conn)
			}
		})
	}
}

func TestBuildEdgeSpokeExternalDeviceConnForDisableHa(t *testing.T) {
	resource := resourceAviatrixEdgeSpokeExternalDeviceConn()

	tests := []struct {
		name        string
		config      map[string]interface{}
		haGwName    string
		expectError bool
		errorSubstr string
		validate    func(t *testing.T, conn *goaviatrix.ExternalDeviceConn)
	}{
		{
			name: "Valid disable HA config - all backup fields empty",
			config: map[string]interface{}{
				"site_id":              "site-123",
				"connection_name":      "conn-test",
				"gw_name":              "edge-gw",
				"connection_type":      "bgp",
				"tunnel_protocol":      "LAN",
				"enable_edge_underlay": false,
			},
			haGwName:    "edge-gw-hagw",
			expectError: false,
			validate: func(t *testing.T, conn *goaviatrix.ExternalDeviceConn) {
				assert.Equal(t, "site-123", conn.VpcID)
				assert.Equal(t, "conn-test", conn.ConnectionName)
				assert.Equal(t, "edge-gw-hagw", conn.GwName)
				assert.False(t, conn.EnableEdgeUnderlay)
			},
		},
		{
			name: "Valid disable HA with edge underlay",
			config: map[string]interface{}{
				"site_id":              "site-123",
				"connection_name":      "conn-test",
				"gw_name":              "edge-gw",
				"enable_edge_underlay": true,
			},
			haGwName:    "edge-gw-hagw",
			expectError: false,
			validate: func(t *testing.T, conn *goaviatrix.ExternalDeviceConn) {
				assert.True(t, conn.EnableEdgeUnderlay)
			},
		},
		{
			name: "Fails when backup_remote_lan_ip is set",
			config: map[string]interface{}{
				"site_id":              "site-123",
				"connection_name":      "conn-test",
				"gw_name":              "edge-gw",
				"backup_remote_lan_ip": "10.1.1.4",
				"enable_edge_underlay": false,
			},
			haGwName:    "edge-gw-hagw",
			expectError: true,
			errorSubstr: "backup_remote_lan_ip",
		},
		{
			name: "Fails when backup_local_lan_ip is set",
			config: map[string]interface{}{
				"site_id":              "site-123",
				"connection_name":      "conn-test",
				"gw_name":              "edge-gw",
				"backup_local_lan_ip":  "10.1.1.3",
				"enable_edge_underlay": false,
			},
			haGwName:    "edge-gw-hagw",
			expectError: true,
			errorSubstr: "backup_remote_lan_ip",
		},
		{
			name: "Fails when backup_bgp_remote_as_num is set",
			config: map[string]interface{}{
				"site_id":                  "site-123",
				"connection_name":          "conn-test",
				"gw_name":                  "edge-gw",
				"backup_bgp_remote_as_num": "65002",
				"enable_edge_underlay":     false,
			},
			haGwName:    "edge-gw-hagw",
			expectError: true,
			errorSubstr: "backup_bgp_remote_as_num",
		},
		{
			name: "Fails when backup_remote_lan_ipv6_ip is set",
			config: map[string]interface{}{
				"site_id":                   "site-123",
				"connection_name":           "conn-test",
				"gw_name":                   "edge-gw",
				"backup_remote_lan_ipv6_ip": "2001:db8::2",
				"enable_edge_underlay":      false,
			},
			haGwName:    "edge-gw-hagw",
			expectError: true,
			errorSubstr: "backup_remote_lan_ipv6_ip",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resource.Schema, tt.config)
			conn, err := buildEdgeSpokeExternalDeviceConnForDisableHa(d, tt.haGwName)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorSubstr != "" {
					assert.Contains(t, err.Error(), tt.errorSubstr)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, conn)
			if tt.validate != nil {
				tt.validate(t, conn)
			}
		})
	}
}
