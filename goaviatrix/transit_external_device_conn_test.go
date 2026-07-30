package goaviatrix

import (
	"testing"

	"github.com/ajg/form"
	"github.com/stretchr/testify/assert"
)

// TestExternalDeviceConnBfdFormFieldNames guards the wire field names the
// connect_transit_gw_to_external_device create request sends. The controller
// argparse args are enable_bfd / bfd_transmit_interval / bfd_receive_interval /
// bfd_detect_multiplier; if the form tags drift from those, the controller
// silently drops BFD and it is never configured (AVX-79351).
func TestExternalDeviceConnBfdFormFieldNames(t *testing.T) {
	conn := &ExternalDeviceConn{
		EnableBfd:     true,
		BfdTxIntv:     300,
		BfdRxIntv:     300,
		BfdMultiplier: 3,
	}
	vals, err := form.EncodeToValues(conn)
	assert.NoError(t, err)

	assert.Equal(t, "true", vals.Get("enable_bfd"))
	assert.Equal(t, "300", vals.Get("bfd_transmit_interval"))
	assert.Equal(t, "300", vals.Get("bfd_receive_interval"))
	assert.Equal(t, "3", vals.Get("bfd_detect_multiplier"))

	// The pre-AVX-79351 names must no longer be emitted.
	assert.Empty(t, vals.Get("bgp_bfd_enabled"))
	assert.Empty(t, vals.Get("bfd_tx_intv"))
	assert.Empty(t, vals.Get("bfd_rx_intv"))
	assert.Empty(t, vals.Get("bfd_multiplier"))
}

// twoTunnels returns a 2-element Tunnels slice; only the length is consulted by
// the HA-classification helpers under test.
func twoTunnels() []TunnelInfo {
	return []TunnelInfo{{}, {}}
}

func edgeTransitGateway() *Gateway {
	return &Gateway{EdgeGateway: true, TransitVpc: "yes"}
}

// TestConfigureHAForTwoDevices covers the local-gateway-HA-enabled path
// (populateNonLANTunnelInfo -> configureLocalGatewayHAEnabled ->
// configureHAForTwoDevices), where two tunnels point at two different remote
// IPs. The wire shape is identical for a genuine remote-HA connection and for a
// single non-HA connection carrying two tunnels, so prior state disambiguates.
func TestConfigureHAForTwoDevices(t *testing.T) {
	const (
		ip1            = "48.194.108.37"
		ip2            = "20.81.31.86"
		localTunnel1   = "169.254.100.2/30"
		localTunnel2   = "169.254.100.6/30"
		remoteTunnel1  = "169.254.100.1/30"
		remoteTunnel2  = "169.254.100.5/30"
		backupRemoteAS = 65001
	)
	detail := EditExternalDeviceConnDetail{
		RemoteGatewayIP:        ip1 + "," + ip2,
		LocalTunnelCidr:        localTunnel1,
		BackupLocalTunnelCidr:  localTunnel2,
		RemoteTunnelCidr:       remoteTunnel1,
		BackupRemoteTunnelCidr: remoteTunnel2,
		Tunnels:                twoTunnels(),
	}
	remoteIP := []string{ip1, ip2}

	t.Run("edge transit, single non-HA multi-tunnel connection (AVX-78064)", func(t *testing.T) {
		conn := &ExternalDeviceConn{}
		configureHAForTwoDevices(conn, detail, edgeTransitGateway(), remoteIP, backupRemoteAS, false)

		assert.Equal(t, Disabled, conn.HAEnabled)
		assert.Equal(t, ip1+","+ip2, conn.RemoteGatewayIP)
		assert.Equal(t, localTunnel1+","+localTunnel2, conn.LocalTunnelCidr)
		assert.Equal(t, remoteTunnel1+","+remoteTunnel2, conn.RemoteTunnelCidr)
		assert.Empty(t, conn.BackupRemoteGatewayIP)
		assert.Empty(t, conn.BackupLocalTunnelCidr)
		assert.Empty(t, conn.BackupRemoteTunnelCidr)
	})

	t.Run("edge transit, genuine remote HA (AVX-68132 regression guard)", func(t *testing.T) {
		conn := &ExternalDeviceConn{}
		configureHAForTwoDevices(conn, detail, edgeTransitGateway(), remoteIP, backupRemoteAS, true)

		assert.Equal(t, Enabled, conn.HAEnabled)
		assert.Equal(t, ip2, conn.BackupRemoteGatewayIP)
		assert.Equal(t, localTunnel1, conn.LocalTunnelCidr)
		assert.Equal(t, localTunnel2, conn.BackupLocalTunnelCidr)
		assert.Equal(t, remoteTunnel1, conn.RemoteTunnelCidr)
		assert.Equal(t, remoteTunnel2, conn.BackupRemoteTunnelCidr)
		assert.Equal(t, backupRemoteAS, conn.BackupBgpRemoteAsNum)
	})

	t.Run("non-edge gateway always combines into one non-HA connection", func(t *testing.T) {
		// prior HA true must not promote a non-edge gateway to HA.
		conn := &ExternalDeviceConn{}
		configureHAForTwoDevices(conn, detail, &Gateway{}, remoteIP, backupRemoteAS, true)

		assert.Equal(t, Disabled, conn.HAEnabled)
		assert.Equal(t, ip1+","+ip2, conn.RemoteGatewayIP)
		assert.Equal(t, localTunnel1+","+localTunnel2, conn.LocalTunnelCidr)
		assert.Equal(t, remoteTunnel1+","+remoteTunnel2, conn.RemoteTunnelCidr)
		assert.Empty(t, conn.BackupRemoteGatewayIP)
	})
}

// TestConfigureLocalGatewayHADisabled covers the local-gateway-HA-disabled path
// (populateNonLANTunnelInfo -> configureLocalGatewayHADisabled).
func TestConfigureLocalGatewayHADisabled(t *testing.T) {
	const (
		ip1            = "48.194.108.37"
		ip2            = "20.81.31.86"
		localTunnel1   = "169.254.100.2/30"
		localTunnel2   = "169.254.100.6/30"
		remoteTunnel1  = "169.254.100.1/30"
		remoteTunnel2  = "169.254.100.5/30"
		backupRemoteAS = 65001
	)

	t.Run("two tunnels, two different IPs, single non-HA connection (AVX-78064)", func(t *testing.T) {
		detail := EditExternalDeviceConnDetail{
			RemoteGatewayIP:        ip1 + "," + ip2,
			LocalTunnelCidr:        localTunnel1,
			BackupLocalTunnelCidr:  localTunnel2,
			RemoteTunnelCidr:       remoteTunnel1,
			BackupRemoteTunnelCidr: remoteTunnel2,
			Tunnels:                twoTunnels(),
		}
		// populateBasicConnectionInfo sets RemoteGatewayIP to the primary IP
		// only; the helper must restore the full comma-joined value.
		conn := &ExternalDeviceConn{RemoteGatewayIP: ip1}
		configureLocalGatewayHADisabled(conn, detail, backupRemoteAS, false)

		assert.Equal(t, Disabled, conn.HAEnabled)
		assert.Equal(t, ip1+","+ip2, conn.RemoteGatewayIP)
		assert.Equal(t, localTunnel1+","+localTunnel2, conn.LocalTunnelCidr)
		assert.Equal(t, remoteTunnel1+","+remoteTunnel2, conn.RemoteTunnelCidr)
		assert.Empty(t, conn.BackupRemoteGatewayIP)
		assert.Empty(t, conn.BackupLocalTunnelCidr)
		assert.Empty(t, conn.BackupRemoteTunnelCidr)
	})

	t.Run("two tunnels, two different IPs, genuine remote HA", func(t *testing.T) {
		detail := EditExternalDeviceConnDetail{
			RemoteGatewayIP:        ip1 + "," + ip2,
			LocalTunnelCidr:        localTunnel1,
			BackupLocalTunnelCidr:  localTunnel2,
			RemoteTunnelCidr:       remoteTunnel1,
			BackupRemoteTunnelCidr: remoteTunnel2,
			Tunnels:                twoTunnels(),
		}
		conn := &ExternalDeviceConn{}
		configureLocalGatewayHADisabled(conn, detail, backupRemoteAS, true)

		assert.Equal(t, Enabled, conn.HAEnabled)
		assert.Equal(t, ip2, conn.BackupRemoteGatewayIP)
		assert.Equal(t, localTunnel1, conn.LocalTunnelCidr)
		assert.Equal(t, localTunnel2, conn.BackupLocalTunnelCidr)
		assert.Equal(t, remoteTunnel1, conn.RemoteTunnelCidr)
		assert.Equal(t, remoteTunnel2, conn.BackupRemoteTunnelCidr)
		assert.Equal(t, backupRemoteAS, conn.BackupBgpRemoteAsNum)
	})

	t.Run("two tunnels, same IP, single device regardless of prior state", func(t *testing.T) {
		detail := EditExternalDeviceConnDetail{
			RemoteGatewayIP:        ip1 + "," + ip1,
			LocalTunnelCidr:        localTunnel1,
			BackupLocalTunnelCidr:  localTunnel2,
			RemoteTunnelCidr:       remoteTunnel1,
			BackupRemoteTunnelCidr: remoteTunnel2,
			Tunnels:                twoTunnels(),
		}
		for _, priorHA := range []bool{false, true} {
			conn := &ExternalDeviceConn{}
			configureLocalGatewayHADisabled(conn, detail, backupRemoteAS, priorHA)

			assert.Equal(t, Disabled, conn.HAEnabled)
			assert.Equal(t, localTunnel1+","+localTunnel2, conn.LocalTunnelCidr)
			assert.Equal(t, remoteTunnel1+","+remoteTunnel2, conn.RemoteTunnelCidr)
			assert.Empty(t, conn.BackupRemoteGatewayIP)
		}
	})

	t.Run("single tunnel is never HA", func(t *testing.T) {
		detail := EditExternalDeviceConnDetail{
			RemoteGatewayIP:  ip1,
			LocalTunnelCidr:  localTunnel1,
			RemoteTunnelCidr: remoteTunnel1,
			Tunnels:          []TunnelInfo{{}},
		}
		conn := &ExternalDeviceConn{}
		configureLocalGatewayHADisabled(conn, detail, backupRemoteAS, false)

		assert.Equal(t, Disabled, conn.HAEnabled)
		assert.Equal(t, localTunnel1, conn.LocalTunnelCidr)
		assert.Equal(t, remoteTunnel1, conn.RemoteTunnelCidr)
		assert.Empty(t, conn.BackupRemoteGatewayIP)
	})
}

// TestConfigureLocalGatewayHAEnabled covers the local-gateway-HA-enabled path
// (populateNonLANTunnelInfo -> configureLocalGatewayHAEnabled).
func TestConfigureLocalGatewayHAEnabled(t *testing.T) {
	const (
		ip1            = "48.194.108.37"
		ip2            = "20.81.31.86"
		localTunnel1   = "169.254.201.1/30"
		localTunnel2   = "169.254.202.1/30"
		remoteTunnel1  = "169.254.201.2/30"
		remoteTunnel2  = "169.254.202.2/30"
		backupRemoteAS = 65001
	)

	t.Run("two tunnels, single remote IP, no remote HA (AVX-79490)", func(t *testing.T) {
		detail := EditExternalDeviceConnDetail{
			RemoteGatewayIP:        ip1,
			LocalTunnelCidr:        localTunnel1,
			BackupLocalTunnelCidr:  localTunnel2,
			RemoteTunnelCidr:       remoteTunnel1,
			BackupRemoteTunnelCidr: remoteTunnel2,
			Tunnels:                twoTunnels(),
		}
		conn := &ExternalDeviceConn{}
		configureLocalGatewayHAEnabled(conn, detail, &Gateway{}, backupRemoteAS, false)

		assert.Equal(t, Disabled, conn.HAEnabled)
		assert.Equal(t, localTunnel1+","+localTunnel2, conn.LocalTunnelCidr)
		assert.Equal(t, remoteTunnel1+","+remoteTunnel2, conn.RemoteTunnelCidr)
		assert.Empty(t, conn.BackupRemoteGatewayIP)
		assert.Empty(t, conn.BackupLocalTunnelCidr)
		assert.Empty(t, conn.BackupRemoteTunnelCidr)
	})

	t.Run("two tunnels, duplicated remote IP, no remote HA", func(t *testing.T) {
		detail := EditExternalDeviceConnDetail{
			RemoteGatewayIP:        ip1 + "," + ip1,
			LocalTunnelCidr:        localTunnel1,
			BackupLocalTunnelCidr:  localTunnel2,
			RemoteTunnelCidr:       remoteTunnel1,
			BackupRemoteTunnelCidr: remoteTunnel2,
			Tunnels:                twoTunnels(),
		}
		conn := &ExternalDeviceConn{}
		configureLocalGatewayHAEnabled(conn, detail, &Gateway{}, backupRemoteAS, false)

		assert.Equal(t, Disabled, conn.HAEnabled)
		assert.Equal(t, localTunnel1+","+localTunnel2, conn.LocalTunnelCidr)
		assert.Equal(t, remoteTunnel1+","+remoteTunnel2, conn.RemoteTunnelCidr)
		assert.Empty(t, conn.BackupRemoteGatewayIP)
		assert.Empty(t, conn.BackupLocalTunnelCidr)
		assert.Empty(t, conn.BackupRemoteTunnelCidr)
	})

	t.Run("two tunnels, two different IPs, remote HA with edge transit", func(t *testing.T) {
		detail := EditExternalDeviceConnDetail{
			RemoteGatewayIP:        ip1 + "," + ip2,
			LocalTunnelCidr:        localTunnel1,
			BackupLocalTunnelCidr:  localTunnel2,
			RemoteTunnelCidr:       remoteTunnel1,
			BackupRemoteTunnelCidr: remoteTunnel2,
			Tunnels:                twoTunnels(),
		}
		conn := &ExternalDeviceConn{}
		configureLocalGatewayHAEnabled(conn, detail, edgeTransitGateway(), backupRemoteAS, true)

		assert.Equal(t, Enabled, conn.HAEnabled)
		assert.Equal(t, ip2, conn.BackupRemoteGatewayIP)
		assert.Equal(t, localTunnel1, conn.LocalTunnelCidr)
		assert.Equal(t, localTunnel2, conn.BackupLocalTunnelCidr)
		assert.Equal(t, remoteTunnel1, conn.RemoteTunnelCidr)
		assert.Equal(t, remoteTunnel2, conn.BackupRemoteTunnelCidr)
		assert.Equal(t, backupRemoteAS, conn.BackupBgpRemoteAsNum)
	})

	t.Run("four tunnels, activemesh enabled, remote HA", func(t *testing.T) {
		detail := EditExternalDeviceConnDetail{
			RemoteGatewayIP:        ip1 + "," + ip2,
			LocalTunnelCidr:        localTunnel1,
			BackupLocalTunnelCidr:  localTunnel2,
			RemoteTunnelCidr:       remoteTunnel1,
			BackupRemoteTunnelCidr: remoteTunnel2,
			Tunnels:                []TunnelInfo{{}, {}, {}, {}},
		}
		conn := &ExternalDeviceConn{}
		configureLocalGatewayHAEnabled(conn, detail, &Gateway{}, backupRemoteAS, true)

		assert.Equal(t, Enabled, conn.HAEnabled)
		assert.Equal(t, ip2, conn.BackupRemoteGatewayIP)
		assert.Equal(t, localTunnel1, conn.LocalTunnelCidr)
		assert.Equal(t, localTunnel2, conn.BackupLocalTunnelCidr)
		assert.Equal(t, remoteTunnel1, conn.RemoteTunnelCidr)
		assert.Equal(t, remoteTunnel2, conn.BackupRemoteTunnelCidr)
		assert.Equal(t, backupRemoteAS, conn.BackupBgpRemoteAsNum)
	})
}
