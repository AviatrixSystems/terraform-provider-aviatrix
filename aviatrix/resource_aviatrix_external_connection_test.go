package aviatrix

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"slices"
	"strconv"
	"testing"

	ctyjson "github.com/hashicorp/go-cty/cty/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestAccAviatrixExternalConnection_basic(t *testing.T) {
	if os.Getenv("SKIP_EXTERNAL_CONNECTION") == "yes" {
		t.Skip("Skipping external connection test as SKIP_EXTERNAL_CONNECTION is set")
	}

	resourceName := "aviatrix_external_connection.test"
	connName := "ext-conn-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckExternalConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccExternalConnectionBasic(connName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExternalConnectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", connName),
					resource.TestCheckResourceAttr(resourceName, "routing_protocol", "bgp"),
					resource.TestCheckResourceAttr(resourceName, "tunnel_protocol", "LAN"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "tunnels.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "tunnels.*", map[string]string{
						"local_gw_name":    "gw-1",
						"remote_device_ip": "10.0.4.10",
						"direct_connect":   "false",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// Import has no configuration, so it records the server-assigned local_device_ip.
				ImportStateVerifyIgnore: []string{"tunnels"},
			},
		},
	})
}

func TestAccAviatrixExternalConnection_fullOptions(t *testing.T) {
	if os.Getenv("SKIP_EXTERNAL_CONNECTION") == "yes" {
		t.Skip("Skipping external connection test as SKIP_EXTERNAL_CONNECTION is set")
	}

	resourceName := "aviatrix_external_connection.test_full"
	connName := "ext-conn-full-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckExternalConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccExternalConnectionFull(connName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExternalConnectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", connName),
					resource.TestCheckResourceAttr(resourceName, "routing_protocol", "bgp"),
					resource.TestCheckResourceAttr(resourceName, "tunnel_protocol", "LAN"),
					resource.TestCheckResourceAttr(resourceName, "bgp_local_as_number", "65000"),
					resource.TestCheckResourceAttr(resourceName, "bgp_remote_as_number", "65001"),
					resource.TestCheckResourceAttr(resourceName, "enable_bgp_multihop", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_bfd", "true"),
					resource.TestCheckResourceAttr(resourceName, "bfd_transmit_interval", "300"),
					resource.TestCheckResourceAttr(resourceName, "bfd_receive_interval", "300"),
					resource.TestCheckResourceAttr(resourceName, "bfd_detect_multiplier", "3"),
					resource.TestCheckResourceAttr(resourceName, "bgp_lan_activemesh", "true"),
					resource.TestCheckResourceAttr(resourceName, "tunnels.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "tunnels.*", map[string]string{
						"local_gw_name":        "gw-1",
						"remote_device_ip":     "10.0.4.10",
						"bgp_remote_as_number": "65001",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "tunnels.*", map[string]string{
						"local_gw_name":    "gw-2",
						"remote_device_ip": "10.0.4.20",
						"direct_connect":   "true",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// Import has no configuration, so it records the server-assigned local_device_ip.
				ImportStateVerifyIgnore: []string{"tunnels"},
			},
		},
	})
}

func TestAccAviatrixExternalConnection_tunnelUpdate(t *testing.T) {
	if os.Getenv("SKIP_EXTERNAL_CONNECTION") == "yes" {
		t.Skip("Skipping external connection test as SKIP_EXTERNAL_CONNECTION is set")
	}

	resourceName := "aviatrix_external_connection.test_update"
	connName := "ext-conn-upd-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckExternalConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccExternalConnectionOneTunnel(connName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExternalConnectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tunnels.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "tunnels.*", map[string]string{
						"local_gw_name":    "gw-1",
						"remote_device_ip": "10.0.4.10",
					}),
				),
			},
			{
				Config: testAccExternalConnectionTwoTunnels(connName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExternalConnectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tunnels.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "tunnels.*", map[string]string{
						"local_gw_name":    "gw-1",
						"remote_device_ip": "10.0.4.10",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "tunnels.*", map[string]string{
						"local_gw_name":    "gw-2",
						"remote_device_ip": "10.0.4.20",
					}),
				),
			},
			{
				Config: testAccExternalConnectionOneTunnel(connName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExternalConnectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tunnels.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "tunnels.*", map[string]string{
						"local_gw_name":    "gw-1",
						"remote_device_ip": "10.0.4.10",
					}),
				),
			},
		},
	})
}

func testAccExternalConnectionBasic(name string) string {
	return fmt.Sprintf(`
resource "aviatrix_external_connection" "test" {
	name             = "%s"
	group_uuid       = "test-group-uuid"
	routing_protocol = "bgp"
	tunnel_protocol  = "LAN"

	tunnels {
		local_gw_name    = "gw-1"
		remote_device_ip = "10.0.4.10"
	}
}
`, name)
}

func testAccExternalConnectionFull(name string) string {
	return fmt.Sprintf(`
resource "aviatrix_external_connection" "test_full" {
	name               = "%s"
	group_uuid         = "test-group-uuid"
	routing_protocol   = "bgp"
	tunnel_protocol    = "LAN"
	bgp_local_as_number  = 65000
	bgp_remote_as_number = 65001
	enable_bgp_multihop  = true
	enable_bfd           = true
	bfd_transmit_interval = 300
	bfd_receive_interval  = 300
	bfd_detect_multiplier = 3
	bgp_lan_activemesh   = true

	tunnels {
		local_gw_name        = "gw-1"
		remote_device_ip     = "10.0.4.10"
		bgp_remote_as_number = 65001
	}

	tunnels {
		local_gw_name    = "gw-2"
		remote_device_ip = "10.0.4.20"
		direct_connect   = true
	}
}
`, name)
}

func testAccExternalConnectionOneTunnel(name string) string {
	return fmt.Sprintf(`
resource "aviatrix_external_connection" "test_update" {
	name             = "%s"
	group_uuid       = "test-group-uuid"
	routing_protocol = "bgp"
	tunnel_protocol  = "LAN"

	tunnels {
		local_gw_name    = "gw-1"
		remote_device_ip = "10.0.4.10"
	}
}
`, name)
}

func testAccExternalConnectionTwoTunnels(name string) string {
	return fmt.Sprintf(`
resource "aviatrix_external_connection" "test_update" {
	name             = "%s"
	group_uuid       = "test-group-uuid"
	routing_protocol = "bgp"
	tunnel_protocol  = "LAN"

	tunnels {
		local_gw_name    = "gw-1"
		remote_device_ip = "10.0.4.10"
	}

	tunnels {
		local_gw_name    = "gw-2"
		remote_device_ip = "10.0.4.20"
	}
}
`, name)
}

func TestAccAviatrixExternalConnection_connUpdate(t *testing.T) {
	if os.Getenv("SKIP_EXTERNAL_CONNECTION") == "yes" {
		t.Skip("Skipping external connection test as SKIP_EXTERNAL_CONNECTION is set")
	}

	resourceName := "aviatrix_external_connection.test_conn_upd"
	connName := "ext-conn-cupd-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckExternalConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccExternalConnectionMutableFields(connName, true, false, false, false, 0, 0, 0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExternalConnectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_bgp_multihop", "true"),
					resource.TestCheckResourceAttr(resourceName, "jumbo_frame", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_bfd", "false"),
					resource.TestCheckResourceAttr(resourceName, "conn_learned_cidrs_approval", "false"),
				),
			},
			{
				Config: testAccExternalConnectionMutableFields(connName, false, true, true, true, 500, 500, 5),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExternalConnectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_bgp_multihop", "false"),
					resource.TestCheckResourceAttr(resourceName, "jumbo_frame", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_bfd", "true"),
					resource.TestCheckResourceAttr(resourceName, "bfd_transmit_interval", "500"),
					resource.TestCheckResourceAttr(resourceName, "bfd_receive_interval", "500"),
					resource.TestCheckResourceAttr(resourceName, "bfd_detect_multiplier", "5"),
					resource.TestCheckResourceAttr(resourceName, "conn_learned_cidrs_approval", "true"),
				),
			},
			{
				Config: testAccExternalConnectionMutableFields(connName, true, false, false, false, 0, 0, 0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExternalConnectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_bgp_multihop", "true"),
					resource.TestCheckResourceAttr(resourceName, "jumbo_frame", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_bfd", "false"),
					resource.TestCheckResourceAttr(resourceName, "conn_learned_cidrs_approval", "false"),
				),
			},
		},
	})
}

func TestAccAviatrixExternalConnection_tunnelMd5Update(t *testing.T) {
	if os.Getenv("SKIP_EXTERNAL_CONNECTION") == "yes" {
		t.Skip("Skipping external connection test as SKIP_EXTERNAL_CONNECTION is set")
	}

	resourceName := "aviatrix_external_connection.test_md5"
	connName := "ext-conn-md5-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckExternalConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccExternalConnectionTunnelMd5(connName, "secret-key-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExternalConnectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tunnels.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "tunnels.*", map[string]string{
						"bgp_md5_key": "secret-key-1",
					}),
				),
			},
			{
				Config: testAccExternalConnectionTunnelMd5(connName, "secret-key-2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExternalConnectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tunnels.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "tunnels.*", map[string]string{
						"bgp_md5_key": "secret-key-2",
					}),
				),
			},
		},
	})
}

func testAccExternalConnectionMutableFields(name string, bgpMultihop, jumboFrame, enableBfd, cidrsApproval bool, bfdTx, bfdRx, bfdMult int) string {
	bfdBlock := ""
	if enableBfd {
		bfdBlock = fmt.Sprintf(`
	bfd_transmit_interval = %d
	bfd_receive_interval  = %d
	bfd_detect_multiplier = %d`, bfdTx, bfdRx, bfdMult)
	}

	return fmt.Sprintf(`
resource "aviatrix_external_connection" "test_conn_upd" {
	name                       = "%s"
	group_uuid                 = "test-group-uuid"
	routing_protocol           = "bgp"
	tunnel_protocol            = "LAN"
	enable_bgp_multihop        = %t
	jumbo_frame                = %t
	enable_bfd                 = %t%s
	conn_learned_cidrs_approval = %t

	tunnels {
		local_gw_name    = "gw-1"
		remote_device_ip = "10.0.4.10"
	}
}
`, name, bgpMultihop, jumboFrame, enableBfd, bfdBlock, cidrsApproval)
}

func testAccExternalConnectionTunnelMd5(name, md5Key string) string {
	return fmt.Sprintf(`
resource "aviatrix_external_connection" "test_md5" {
	name             = "%s"
	group_uuid       = "test-group-uuid"
	routing_protocol = "bgp"
	tunnel_protocol  = "LAN"

	tunnels {
		local_gw_name    = "gw-1"
		remote_device_ip = "10.0.4.10"
		bgp_md5_key      = "%s"
	}
}
`, name, md5Key)
}

// The controller fills in both ASNs when the config leaves them out, so an
// omitted ASN must not plan a replacement of a live connection.
func TestExternalConnectionAsnDiff(t *testing.T) {
	tunnelHash := tunnelSetHash(map[string]any{
		"local_gw_name":    "transit-gw",
		"remote_device_ip": "10.1.0.91",
	})
	tunnelPrefix := fmt.Sprintf("tunnels.%d.", tunnelHash)
	state := &terraform.InstanceState{
		ID: "8aa0ac233f4e42a4897c6030f38dd09a",
		Attributes: map[string]string{
			"id":                              "8aa0ac233f4e42a4897c6030f38dd09a",
			"uuid":                            "8aa0ac233f4e42a4897c6030f38dd09a",
			"name":                            "conn",
			"group_uuid":                      "25649ff81aad431e8911671714c5922b",
			"routing_protocol":                "bgp",
			"tunnel_protocol":                 "LAN",
			"bgp_local_as_number":             "65001",
			"bgp_remote_as_number":            "65101",
			"peer_vnet_id":                    "",
			"ipv6_enabled":                    "false",
			"edge_underlay":                   "false",
			"bgp_lan_activemesh":              "false",
			"enable_bgp_multihop":             "true",
			"jumbo_frame":                     "false",
			"enable_bfd":                      "false",
			"conn_learned_cidrs_approval":     "false",
			"tunnels.#":                       "1",
			tunnelPrefix + "local_gw_name":    "transit-gw",
			tunnelPrefix + "remote_device_ip": "10.1.0.91",
		},
	}

	tests := []struct {
		name          string
		asns          map[string]any
		wantChanged   []string
		wantReplacing bool
	}{
		{name: "omitted", asns: map[string]any{}},
		{
			name: "same values",
			asns: map[string]any{"bgp_local_as_number": 65001, "bgp_remote_as_number": 65101},
		},
		{
			name:          "remote changed",
			asns:          map[string]any{"bgp_remote_as_number": 65102},
			wantChanged:   []string{"bgp_remote_as_number"},
			wantReplacing: true,
		},
		{
			name:          "local changed",
			asns:          map[string]any{"bgp_local_as_number": 65002},
			wantChanged:   []string{"bgp_local_as_number"},
			wantReplacing: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := map[string]any{
				"name":             "conn",
				"group_uuid":       "25649ff81aad431e8911671714c5922b",
				"routing_protocol": "bgp",
				"tunnel_protocol":  "LAN",
				"tunnels": []any{map[string]any{
					"local_gw_name":    "transit-gw",
					"remote_device_ip": "10.1.0.91",
				}},
			}
			for k, v := range tt.asns {
				raw[k] = v
			}

			diff, err := resourceAviatrixExternalConnection().Diff(
				context.Background(), state, terraform.NewResourceConfigRaw(raw), nil)
			require.NoError(t, err)

			for _, attr := range []string{"bgp_local_as_number", "bgp_remote_as_number"} {
				var attrDiff *terraform.ResourceAttrDiff
				if diff != nil {
					attrDiff = diff.Attributes[attr]
				}
				// On replacement an omitted computed ASN is "known after apply", not a change.
				changed := attrDiff != nil && !attrDiff.NewComputed && attrDiff.Old != attrDiff.New
				assert.Equal(t, slices.Contains(tt.wantChanged, attr), changed, "%s diff: %+v", attr, attrDiff)
			}
			assert.Equal(t, tt.wantReplacing, diff != nil && diff.RequiresNew())
		})
	}
}

// externalConnectionDiff plans configTunnels against a state holding stateTunnels,
// on a connection with IPv6 enabled.
func externalConnectionDiff(t *testing.T, stateTunnels, configTunnels []map[string]any) (*terraform.InstanceDiff, error) {
	t.Helper()
	attrs := map[string]string{
		"id":                          "8aa0ac233f4e42a4897c6030f38dd09a",
		"uuid":                        "8aa0ac233f4e42a4897c6030f38dd09a",
		"name":                        "conn",
		"group_uuid":                  "25649ff81aad431e8911671714c5922b",
		"routing_protocol":            "bgp",
		"tunnel_protocol":             "LAN",
		"bgp_local_as_number":         "65001",
		"bgp_remote_as_number":        "65101",
		"peer_vnet_id":                "",
		"ipv6_enabled":                "true",
		"edge_underlay":               "false",
		"bgp_lan_activemesh":          "false",
		"enable_bgp_multihop":         "true",
		"jumbo_frame":                 "false",
		"enable_bfd":                  "false",
		"conn_learned_cidrs_approval": "false",
		"tunnels.#":                   strconv.Itoa(len(stateTunnels)),
	}
	for _, tunnel := range stateTunnels {
		prefix := fmt.Sprintf("tunnels.%d.", tunnelSetHash(tunnel))
		for k, v := range tunnel {
			attrs[prefix+k] = fmt.Sprint(v)
		}
	}
	configList := make([]any, 0, len(configTunnels))
	for _, tunnel := range configTunnels {
		configList = append(configList, tunnel)
	}
	raw := map[string]any{
		"name":                "conn",
		"group_uuid":          "25649ff81aad431e8911671714c5922b",
		"routing_protocol":    "bgp",
		"tunnel_protocol":     "LAN",
		"bgp_local_as_number": 65001,
		"ipv6_enabled":        true,
		"tunnels":             configList,
	}
	state := withRawConfig(t, &terraform.InstanceState{ID: attrs["id"], Attributes: attrs}, raw)
	return resourceAviatrixExternalConnection().Diff(context.Background(), state, terraform.NewResourceConfigRaw(raw), nil)
}

// unknownValue is the value the SDK reads as unknown in a raw test config.
const unknownValue = "74D93920-ED26-11E3-AC10-0800200C9A66"

// withRawConfig sets the raw config on state, which Terraform always sends with a
// plan request but the SDK's Diff does not derive from the ResourceConfig.
func withRawConfig(t *testing.T, state *terraform.InstanceState, raw map[string]any) *terraform.InstanceState {
	t.Helper()
	b, err := json.Marshal(raw)
	require.NoError(t, err)
	state.RawConfig, err = ctyjson.Unmarshal(b, resourceAviatrixExternalConnection().CoreConfigSchema().ImpliedType())
	require.NoError(t, err)
	return state
}

func TestExternalConnectionTunnelDiff(t *testing.T) {
	stateTunnel := func(gw, remoteIP, assignedIP, localIPv6, remoteIPv6 string) map[string]any {
		return map[string]any{
			"local_gw_name":            gw,
			"remote_device_ip":         remoteIP,
			"local_device_ip":          "",
			"assigned_local_device_ip": assignedIP,
			"bgp_remote_as_number":     0,
			"local_device_ipv6":        localIPv6,
			"remote_device_ipv6":       remoteIPv6,
			"bgp_md5_key":              "",
			"direct_connect":           false,
			"status":                   "UP",
		}
	}
	stateTunnels := []map[string]any{
		stateTunnel("transit-gw", "10.1.0.91", "10.1.0.10", "fd00::1/126", "fd00::2/126"),
		stateTunnel("transit-gw-ha", "10.1.1.91", "10.1.1.10", "fd00:1::1/126", "fd00:1::2/126"),
	}
	configTunnel := func(gw, remoteIP, localIPv6, remoteIPv6 string, override map[string]any) map[string]any {
		tunnel := map[string]any{
			"local_gw_name":      gw,
			"remote_device_ip":   remoteIP,
			"local_device_ipv6":  localIPv6,
			"remote_device_ipv6": remoteIPv6,
		}
		maps.Copy(tunnel, override)
		return tunnel
	}

	tests := []struct {
		name        string
		config      map[string]any
		haConfig    map[string]any
		oneTunnel   bool
		wantDiff    bool
		wantReplace bool
	}{
		{name: "unchanged"},
		{
			name:     "bgp_md5_key added",
			config:   map[string]any{"bgp_md5_key": "aviatrixTest23"},
			haConfig: map[string]any{"bgp_md5_key": "aviatrixTest24"},
			wantDiff: true,
		},
		{name: "local_device_ip set", config: map[string]any{"local_device_ip": "10.1.0.11"}, wantDiff: true},
		{name: "bgp_remote_as_number set to the connection ASN", config: map[string]any{"bgp_remote_as_number": 65101}, wantDiff: true},
		{name: "local_device_ipv6 changed", config: map[string]any{"local_device_ipv6": "fd00::5/126"}, wantDiff: true},
		{name: "remote_device_ipv6 changed", config: map[string]any{"remote_device_ipv6": "fd00::6/126"}, wantDiff: true},
		{
			name:        "bgp_remote_as_number changed on every tunnel",
			config:      map[string]any{"bgp_remote_as_number": 65102},
			haConfig:    map[string]any{"bgp_remote_as_number": 65102},
			wantDiff:    true,
			wantReplace: true,
		},
		{name: "direct_connect set", config: map[string]any{"direct_connect": true}, wantDiff: true, wantReplace: true},
		{
			name:        "local_device_ip set on the only tunnel",
			config:      map[string]any{"local_device_ip": "10.1.0.11"},
			oneTunnel:   true,
			wantDiff:    true,
			wantReplace: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := stateTunnels
			config := []map[string]any{
				configTunnel("transit-gw", "10.1.0.91", "fd00::1/126", "fd00::2/126", tt.config),
				configTunnel("transit-gw-ha", "10.1.1.91", "fd00:1::1/126", "fd00:1::2/126", tt.haConfig),
			}
			if tt.oneTunnel {
				state, config = state[:1], config[:1]
			}
			diff, err := externalConnectionDiff(t, state, config)
			require.NoError(t, err)
			assert.Equal(t, tt.wantDiff, diff != nil && len(diff.Attributes) > 0, "diff: %+v", diff)
			assert.Equal(t, tt.wantReplace, diff != nil && diff.RequiresNew())
		})
	}
}

func TestTunnelStateFromServer(t *testing.T) {
	server := goaviatrix.ExternalConnectionTunnel{
		LocalGwName:      "transit-gw",
		RemoteDeviceIP:   "10.1.0.91",
		LocalDeviceIP:    "10.1.0.10",
		BgpRemoteAsNum:   65101,
		LocalDeviceIPv6:  "fd00::1",
		RemoteDeviceIPv6: "fd00::2",
		BgpMd5Key:        "***",
		Status:           "UP",
	}

	tests := []struct {
		name     string
		prior    *goaviatrix.ExternalConnectionTunnel
		override func(*goaviatrix.ExternalConnectionTunnel)
		want     func(*goaviatrix.ExternalConnectionTunnel)
	}{
		{
			name: "import takes server values, including the masked key",
			want: func(*goaviatrix.ExternalConnectionTunnel) {},
		},
		{
			name:  "key set outside Terraform is reported",
			prior: &goaviatrix.ExternalConnectionTunnel{LocalDeviceIP: "10.1.0.10", BgpRemoteAsNum: 65101, LocalDeviceIPv6: "fd00::1", RemoteDeviceIPv6: "fd00::2"},
			want:  func(*goaviatrix.ExternalConnectionTunnel) {},
		},
		{
			name:  "omitted local_device_ip and ASN stay omitted",
			prior: &goaviatrix.ExternalConnectionTunnel{LocalDeviceIPv6: "fd00::1", RemoteDeviceIPv6: "fd00::2", BgpMd5Key: "k"},
			want: func(w *goaviatrix.ExternalConnectionTunnel) {
				w.LocalDeviceIP, w.BgpRemoteAsNum, w.BgpMd5Key = "", 0, "k"
			},
		},
		{
			name:  "configured spelling is kept",
			prior: &goaviatrix.ExternalConnectionTunnel{LocalDeviceIP: "vpc:subnet:10.1.0.10", BgpRemoteAsNum: 65101, LocalDeviceIPv6: "fd00:0::1/126", RemoteDeviceIPv6: "fd00::2/126"},
			override: func(s *goaviatrix.ExternalConnectionTunnel) {
				s.BgpMd5Key = ""
			},
			want: func(w *goaviatrix.ExternalConnectionTunnel) {
				w.LocalDeviceIP, w.LocalDeviceIPv6, w.RemoteDeviceIPv6, w.BgpMd5Key = "vpc:subnet:10.1.0.10", "fd00:0::1/126", "fd00::2/126", ""
			},
		},
		{
			name:     "tunnel ASN omitted by server under a connection-level ASN",
			prior:    &goaviatrix.ExternalConnectionTunnel{LocalDeviceIP: "10.1.0.10", BgpRemoteAsNum: 65101, LocalDeviceIPv6: "fd00::1", RemoteDeviceIPv6: "fd00::2", BgpMd5Key: "k"},
			override: func(s *goaviatrix.ExternalConnectionTunnel) { s.BgpRemoteAsNum = 0 },
			want:     func(w *goaviatrix.ExternalConnectionTunnel) { w.BgpRemoteAsNum, w.BgpMd5Key = 65101, "k" },
		},
		{
			name:  "out-of-band changes are reported",
			prior: &goaviatrix.ExternalConnectionTunnel{LocalDeviceIP: "10.1.0.20", BgpRemoteAsNum: 65102, LocalDeviceIPv6: "fd00::5/126", RemoteDeviceIPv6: "fd00::6/126", BgpMd5Key: "k", DirectConnect: true},
			override: func(s *goaviatrix.ExternalConnectionTunnel) {
				s.BgpMd5Key = ""
			},
			want: func(w *goaviatrix.ExternalConnectionTunnel) { w.BgpMd5Key = "" },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := server
			if tt.override != nil {
				tt.override(&s)
			}
			var prior goaviatrix.ExternalConnectionTunnel
			if tt.prior != nil {
				prior = *tt.prior
			}
			want := s
			tt.want(&want)
			assert.Equal(t, want, tunnelStateFromServer(prior, tt.prior != nil, s))
		})
	}
}

func TestTunnelNeedsRecreate(t *testing.T) {
	running := goaviatrix.ExternalConnectionTunnel{
		LocalGwName:      "transit-gw",
		RemoteDeviceIP:   "10.1.0.91",
		LocalDeviceIPv6:  "fd00::1/126",
		RemoteDeviceIPv6: "fd00::2/126",
	}
	const assignedIP = "10.1.0.10"
	const connASN = 65101

	tests := []struct {
		name         string
		old          func(*goaviatrix.ExternalConnectionTunnel)
		new          func(*goaviatrix.ExternalConnectionTunnel)
		wantRecreate bool
	}{
		{name: "local_device_ip set to the assigned IP", new: func(n *goaviatrix.ExternalConnectionTunnel) { n.LocalDeviceIP = assignedIP }},
		{name: "local_device_ip set to another IP", new: func(n *goaviatrix.ExternalConnectionTunnel) { n.LocalDeviceIP = "10.1.0.11" }, wantRecreate: true},
		{
			name: "local_device_ip removed from config",
			old:  func(o *goaviatrix.ExternalConnectionTunnel) { o.LocalDeviceIP = assignedIP },
		},
		{
			name: "GCP vpc changed",
			old:  func(o *goaviatrix.ExternalConnectionTunnel) { o.LocalDeviceIP = "vpc-a:subnet:10.1.0.10" },
			new:  func(n *goaviatrix.ExternalConnectionTunnel) { n.LocalDeviceIP = "vpc-b:subnet:10.1.0.10" }, wantRecreate: true,
		},
		{name: "tunnel ASN set to the connection-level ASN", new: func(n *goaviatrix.ExternalConnectionTunnel) { n.BgpRemoteAsNum = connASN }},
		{name: "tunnel ASN set to another ASN", new: func(n *goaviatrix.ExternalConnectionTunnel) { n.BgpRemoteAsNum = 65102 }, wantRecreate: true},
		{
			name: "tunnel ASN removed from config",
			old:  func(o *goaviatrix.ExternalConnectionTunnel) { o.BgpRemoteAsNum = 65099 },
		},
		{name: "IPv6 respelled", new: func(n *goaviatrix.ExternalConnectionTunnel) { n.LocalDeviceIPv6 = "fd00:0:0::1/126" }},
		{name: "IPv6 changed", new: func(n *goaviatrix.ExternalConnectionTunnel) { n.RemoteDeviceIPv6 = "fd00::6/126" }, wantRecreate: true},
		{name: "IPv6 prefix length changed", new: func(n *goaviatrix.ExternalConnectionTunnel) { n.RemoteDeviceIPv6 = "fd00::2/64" }, wantRecreate: true},
		{name: "IPv6 prefix length removed", new: func(n *goaviatrix.ExternalConnectionTunnel) { n.LocalDeviceIPv6 = "fd00::1" }, wantRecreate: true},
		{
			name: "IPv6 prefix length added to an address read from the server",
			old:  func(o *goaviatrix.ExternalConnectionTunnel) { o.LocalDeviceIPv6 = "fd00::1" },
		},
		{name: "only bgp_md5_key changed", new: func(n *goaviatrix.ExternalConnectionTunnel) { n.BgpMd5Key = "k" }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldT, newT := running, running
			if tt.old != nil {
				tt.old(&oldT)
			}
			if tt.new != nil {
				tt.new(&newT)
			}
			assert.Equal(t, tt.wantRecreate, tunnelNeedsRecreate(oldT, newT, assignedIP, connASN))
		})
	}
}

func TestExternalConnectionCustomizeDiff(t *testing.T) {
	tunnel := func(gw, remoteIP string) map[string]any {
		return map[string]any{"local_gw_name": gw, "remote_device_ip": remoteIP}
	}
	withASN := func(tunnel map[string]any, asn int) map[string]any {
		tunnel["bgp_remote_as_number"] = asn
		return tunnel
	}
	withMd5 := func(tunnel map[string]any, key string) map[string]any {
		tunnel["bgp_md5_key"] = key
		return tunnel
	}
	tests := []struct {
		name    string
		conn    map[string]any
		tunnels []any
		wantErr string
	}{
		{name: "valid", tunnels: []any{tunnel("gw-1", "10.1.0.91"), tunnel("gw-2", "10.1.0.92")}},
		{
			name:    "IPv6 fields without ipv6_enabled",
			tunnels: []any{map[string]any{"local_gw_name": "gw-1", "remote_device_ip": "10.1.0.91", "remote_device_ipv6": "fd00::2/126"}},
			wantErr: "require ipv6_enabled = true",
		},
		{
			name: "duplicate tunnel",
			tunnels: []any{
				tunnel("gw-1", "10.1.0.91"),
				map[string]any{"local_gw_name": "gw-1", "remote_device_ip": "10.1.0.91", "bgp_md5_key": "k"},
			},
			wantErr: "duplicate tunnel",
		},
		{
			name:    "two tunnels on one gateway without ActiveMesh",
			tunnels: []any{tunnel("gw-1", "10.1.0.91"), tunnel("gw-1", "10.1.0.92")},
			wantErr: "multiple tunnels per gateway require bgp_lan_activemesh = true",
		},
		{
			name:    "two tunnels on one gateway with ActiveMesh",
			conn:    map[string]any{"bgp_lan_activemesh": true},
			tunnels: []any{tunnel("gw-1", "10.1.0.91"), tunnel("gw-1", "10.1.0.92")},
		},
		{
			name:    "tunnel ASN differs from the connection ASN",
			conn:    map[string]any{"bgp_remote_as_number": 65101},
			tunnels: []any{withASN(tunnel("gw-1", "10.1.0.91"), 65102)},
			wantErr: "same bgp_remote_as_number",
		},
		{
			name:    "tunnel ASNs differ",
			tunnels: []any{withASN(tunnel("gw-1", "10.1.0.91"), 65101), withASN(tunnel("gw-2", "10.1.0.92"), 65102)},
			wantErr: "same bgp_remote_as_number",
		},
		{
			name:    "tunnel ASNs differ on an edge underlay",
			conn:    map[string]any{"edge_underlay": true},
			tunnels: []any{withASN(tunnel("gw-1", "10.1.0.91"), 65101), withASN(tunnel("gw-2", "10.1.0.92"), 65102)},
		},
		{
			name:    "two tunnels on one gateway with ActiveMesh not yet known",
			conn:    map[string]any{"bgp_lan_activemesh": unknownValue},
			tunnels: []any{tunnel("gw-1", "10.1.0.91"), tunnel("gw-1", "10.1.0.92")},
		},
		{
			name:    "tunnel ASNs differ with edge_underlay not yet known",
			conn:    map[string]any{"edge_underlay": unknownValue},
			tunnels: []any{withASN(tunnel("gw-1", "10.1.0.91"), 65101), withASN(tunnel("gw-2", "10.1.0.92"), 65102)},
		},
		{
			name:    "bgp_md5_key on all tunnels",
			tunnels: []any{withMd5(tunnel("gw-1", "10.1.0.91"), "k1"), withMd5(tunnel("gw-2", "10.1.0.92"), "k2")},
		},
		{
			name:    "bgp_md5_key on some tunnels",
			tunnels: []any{withMd5(tunnel("gw-1", "10.1.0.91"), "k1"), tunnel("gw-2", "10.1.0.92")},
			wantErr: "bgp_md5_key must be set on all tunnels or on none",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := map[string]any{
				"name":             "conn",
				"group_uuid":       "25649ff81aad431e8911671714c5922b",
				"routing_protocol": "bgp",
				"tunnel_protocol":  "LAN",
				"ipv6_enabled":     false,
				"tunnels":          tt.tunnels,
			}
			maps.Copy(raw, tt.conn)
			state := &terraform.InstanceState{}
			// The raw config cannot hold an unknown value, so only the
			// ResourceConfig carries it.
			if !slices.Contains(slices.Collect(maps.Values(tt.conn)), any(unknownValue)) {
				state = withRawConfig(t, state, raw)
			}
			_, err := resourceAviatrixExternalConnection().Diff(context.Background(),
				state, terraform.NewResourceConfigRaw(raw), nil)
			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}

func TestPlanTunnelChanges(t *testing.T) {
	tunnel := func(gw, remoteIP string, opts ...func(*goaviatrix.ExternalConnectionTunnel)) goaviatrix.ExternalConnectionTunnel {
		t := goaviatrix.ExternalConnectionTunnel{LocalGwName: gw, RemoteDeviceIP: remoteIP}
		for _, opt := range opts {
			opt(&t)
		}
		return t
	}
	localIP := func(ip string) func(*goaviatrix.ExternalConnectionTunnel) {
		return func(t *goaviatrix.ExternalConnectionTunnel) { t.LocalDeviceIP = ip }
	}
	asn := func(asn int) func(*goaviatrix.ExternalConnectionTunnel) {
		return func(t *goaviatrix.ExternalConnectionTunnel) { t.BgpRemoteAsNum = asn }
	}
	md5 := func(key string) func(*goaviatrix.ExternalConnectionTunnel) {
		return func(t *goaviatrix.ExternalConnectionTunnel) { t.BgpMd5Key = key }
	}
	directConnect := func(t *goaviatrix.ExternalConnectionTunnel) { t.DirectConnect = true }
	three := []goaviatrix.ExternalConnectionTunnel{
		tunnel("gw-1", "10.0.1.5"), tunnel("gw-2", "10.0.2.5"), tunnel("gw-3", "10.0.3.5"),
	}
	assigned := map[tunnelKey]string{
		{"gw-1", "10.0.1.5"}: "10.0.1.10",
		{"gw-2", "10.0.2.5"}: "10.0.2.10",
		{"gw-3", "10.0.3.5"}: "10.0.3.10",
	}

	tests := []struct {
		name         string
		old          []goaviatrix.ExternalConnectionTunnel
		new          []goaviatrix.ExternalConnectionTunnel
		perTunnelASN bool
		wantOps      []string
		wantReplace  bool
	}{
		{name: "unchanged", old: three, new: three},
		{
			name:    "bgp_md5_key changed",
			old:     three,
			new:     []goaviatrix.ExternalConnectionTunnel{tunnel("gw-1", "10.0.1.5", md5("k")), three[1], three[2]},
			wantOps: []string{"md5 gw-1/10.0.1.5"},
		},
		{
			name:    "local_device_ip changed",
			old:     three,
			new:     []goaviatrix.ExternalConnectionTunnel{three[0], tunnel("gw-2", "10.0.2.5", localIP("10.0.2.11")), three[2]},
			wantOps: []string{"delete gw-2/10.0.2.5", "add gw-2/10.0.2.5"},
		},
		{
			name: "local_device_ip changed on every tunnel",
			old:  three,
			new: []goaviatrix.ExternalConnectionTunnel{
				tunnel("gw-1", "10.0.1.5", localIP("10.0.1.11")),
				tunnel("gw-2", "10.0.2.5", localIP("10.0.2.11")),
				tunnel("gw-3", "10.0.3.5", localIP("10.0.3.11")),
			},
			wantOps: []string{
				"delete gw-1/10.0.1.5", "add gw-1/10.0.1.5",
				"delete gw-2/10.0.2.5", "add gw-2/10.0.2.5",
				"delete gw-3/10.0.3.5", "add gw-3/10.0.3.5",
			},
		},
		{
			name:    "local_device_ip set to the assigned IP",
			old:     three,
			new:     []goaviatrix.ExternalConnectionTunnel{tunnel("gw-1", "10.0.1.5", localIP("10.0.1.10")), three[1], three[2]},
			wantOps: nil,
		},
		{
			name:    "remote_device_ip changed",
			old:     three,
			new:     []goaviatrix.ExternalConnectionTunnel{tunnel("gw-1", "10.0.1.6"), three[1], three[2]},
			wantOps: []string{"delete gw-1/10.0.1.5", "add gw-1/10.0.1.6"},
		},
		{
			name:    "tunnel moved to another gateway",
			old:     three,
			new:     []goaviatrix.ExternalConnectionTunnel{three[0], three[1], tunnel("gw-4", "10.0.3.5")},
			wantOps: []string{"add gw-4/10.0.3.5", "delete gw-3/10.0.3.5"},
		},
		{
			name:    "only tunnel moved to another gateway",
			old:     three[:1],
			new:     []goaviatrix.ExternalConnectionTunnel{tunnel("gw-4", "10.0.1.5")},
			wantOps: []string{"add gw-4/10.0.1.5", "delete gw-1/10.0.1.5"},
		},
		{name: "tunnel added", old: three[:2], new: three, wantOps: []string{"add gw-3/10.0.3.5"}},
		{name: "tunnel removed", old: three, new: three[:2], wantOps: []string{"delete gw-3/10.0.3.5"}},
		{
			name:        "local_device_ip changed on the only tunnel",
			old:         three[:1],
			new:         []goaviatrix.ExternalConnectionTunnel{tunnel("gw-1", "10.0.1.5", localIP("10.0.1.11"))},
			wantReplace: true,
		},
		{
			name: "bgp_remote_as_number changed on every tunnel",
			old:  three,
			new: []goaviatrix.ExternalConnectionTunnel{
				tunnel("gw-1", "10.0.1.5", asn(65102)), tunnel("gw-2", "10.0.2.5", asn(65102)), tunnel("gw-3", "10.0.3.5", asn(65102)),
			},
			wantReplace: true,
		},
		{
			name:    "bgp_remote_as_number set to the connection ASN",
			old:     three,
			new:     []goaviatrix.ExternalConnectionTunnel{tunnel("gw-1", "10.0.1.5", asn(65101)), three[1], three[2]},
			wantOps: nil,
		},
		{
			name:         "per-tunnel bgp_remote_as_number changed on an edge underlay",
			old:          three,
			new:          []goaviatrix.ExternalConnectionTunnel{tunnel("gw-1", "10.0.1.5", asn(65102)), three[1], three[2]},
			perTunnelASN: true,
			wantOps:      []string{"delete gw-1/10.0.1.5", "add gw-1/10.0.1.5"},
		},
		{
			name:        "direct_connect changed",
			old:         three,
			new:         []goaviatrix.ExternalConnectionTunnel{tunnel("gw-1", "10.0.1.5", directConnect), three[1], three[2]},
			wantReplace: true,
		},
		{
			name:        "direct_connect changed with remote_device_ip",
			old:         three,
			new:         []goaviatrix.ExternalConnectionTunnel{tunnel("gw-1", "10.0.1.6", directConnect), three[1], three[2]},
			wantReplace: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pc := tunnelPlanContext{assignedIPs: assigned, connRemoteASN: 65101, perTunnelASN: tt.perTunnelASN}
			ops, replace := planTunnelChanges(slices.Clone(tt.old), slices.Clone(tt.new), pc)
			require.Equal(t, tt.wantReplace, replace)
			if replace {
				return
			}
			kinds := map[tunnelOpKind]string{tunnelOpAdd: "add", tunnelOpDelete: "delete", tunnelOpSetMd5Key: "md5"}
			var got []string
			for _, op := range ops {
				got = append(got, fmt.Sprintf("%s %s/%s", kinds[op.kind], op.tunnel.LocalGwName, op.tunnel.RemoteDeviceIP))
			}
			assert.Equal(t, tt.wantOps, got)
		})
	}
}

func testAccCheckExternalConnectionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("external connection not found: %s", resourceName)
		}

		client := mustClient(testAccProvider.Meta())

		_, err := client.GetExternalConnection(context.Background(), rs.Primary.ID)
		if errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("external connection not found")
		}

		return err
	}
}

func testAccCheckExternalConnectionDestroy(s *terraform.State) error {
	client := mustClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_external_connection" {
			continue
		}

		_, err := client.GetExternalConnection(context.Background(), rs.Primary.ID)
		if !errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("external connection still exists")
		}
	}

	return nil
}
