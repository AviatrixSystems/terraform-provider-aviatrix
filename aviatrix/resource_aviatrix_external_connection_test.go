package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

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
						"bgp_remote_as_number": "65002",
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
		bgp_remote_as_number = 65002
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
