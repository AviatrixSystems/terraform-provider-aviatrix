package goaviatrix

import (
	"fmt"
	"strings"
)

// EditExternalDeviceConnTunnel is one in-place edit of one tunnel of an
// external device connection, sent to the v2 edit_external_device_conn action.
//
// Not to be confused with EditExternalDeviceConnDetail, which despite its name
// is the *read* struct for get_site2cloud_conn_detail.
//
// Only BGP-over-IPsec connections can be edited this way: the controller
// refuses LAN and GRE outright, and refuses the BGP fields on a connection
// whose routing is static.
//
// Every field applies to the single tunnel named by TunnelName, which is
// mandatory and has no default. There is deliberately no backup_* form of any
// field - nothing records that a tunnel is the backup, "primary" is just
// position 0 of the connection's tunnel list, so a backup_* spelling would only
// be a second, worse selector. Use ExternalDeviceConnTunnelName to turn a
// position in the connection-level CIDR lists into the name to send here.
//
// Two connection-level invariants make a write carry to the tunnel in the other
// slot: same_backup_psk mirrors PreSharedKey, and single_ip_ha mirrors
// RemoteDeviceIP. Editing either tunnel of such a pair therefore changes both.
type EditExternalDeviceConnTunnel struct {
	Action         string `form:"action,omitempty"`
	CID            string `form:"CID,omitempty"`
	GwName         string `form:"gateway_name,omitempty"`
	ConnectionName string `form:"conn_name,omitempty"`
	TunnelName     string `form:"tunnel_name,omitempty"`
	// BgpRemoteAsNumber is a string rather than an int so that "unset" and 0
	// stay distinguishable; the schema attribute it comes from is a string too.
	BgpRemoteAsNumber string `form:"bgp_remote_as_number,omitempty"`
	// BgpLocalIP and BgpRemoteIP are host addresses in /30 CIDR notation, e.g.
	// "169.254.0.1/30" - the same spelling the read API reports them in. The
	// controller validates the effective local+remote pair even when only one
	// side is supplied, pairing the supplied one with the tunnel's current
	// value, so sending a single side is legal.
	BgpLocalIP     string `form:"bgp_local_ip,omitempty"`
	BgpRemoteIP    string `form:"bgp_remote_ip,omitempty"`
	RemoteDeviceIP string `form:"remote_device_ip,omitempty"`
	PreSharedKey   string `form:"pre_shared_key,omitempty"`
}

// IsEmpty reports whether this edit would change nothing, i.e. carries only
// selection and no field to write.
func (e *EditExternalDeviceConnTunnel) IsEmpty() bool {
	return e.BgpRemoteAsNumber == "" && e.BgpLocalIP == "" && e.BgpRemoteIP == "" &&
		e.RemoteDeviceIP == "" && e.PreSharedKey == ""
}

// EditExternalDeviceConnTunnel applies one in-place edit to one tunnel.
//
// The action is not atomic across fields: the controller validates and applies
// them one at a time, so a request that fails part way through leaves the
// earlier fields written. Callers must not treat a failure as a no-op.
func (c *Client) EditExternalDeviceConnTunnel(edit *EditExternalDeviceConnTunnel) error {
	edit.CID = c.CID
	edit.Action = "edit_external_device_conn"

	// PostAPI2, not PostAPI: this action is only reachable under /v2/api.
	return c.PostAPI2(edit.Action, edit, BasicCheck)
}

// EditExternalDeviceConnBgpAdminShutdown shuts the BGP neighbors of a whole
// external device connection down, or brings them back up, through the same
// edit_external_device_conn action.
//
// This is a separate request type from EditExternalDeviceConnTunnel, rather
// than one more field on it, because the controller refuses to see the two in
// one request: bgp_admin_shutdown covers every tunnel of the connection and
// takes no tunnel selector, so a request carrying it alongside any per-tunnel
// field - or alongside a tunnel_name - is rejected outright rather than partly
// applied. Keeping the two shapes apart in the type system means a caller
// cannot build a request the controller will refuse.
//
// Submitting it alone is also what keeps it reaching the handler at all. The
// controller dispatches this action through a long if/elif chain in which
// several earlier branches fire on conn_name plus one other field; every one of
// those is conjunctive, so a request carrying only the connection, the gateway
// and this field falls through to the edit handler. A request that bundled
// another field could be claimed by an earlier branch and silently do something
// else.
type EditExternalDeviceConnBgpAdminShutdown struct {
	Action         string `form:"action,omitempty"`
	CID            string `form:"CID,omitempty"`
	GwName         string `form:"gateway_name,omitempty"`
	ConnectionName string `form:"conn_name,omitempty"`
	// BgpAdminShutdown is "yes" or "no" rather than a bool: the controller
	// declares the argument with choices=["yes","no"] and treats anything else
	// as a refusal. Note that the read side reports the same state as a JSON
	// boolean, so the two directions do not share a spelling - use
	// BgpAdminShutdownArg to cross over.
	BgpAdminShutdown string `form:"bgp_admin_shutdown,omitempty"`
}

// BgpAdminShutdownArg renders the boolean the read API reports, and terraform
// stores, as the "yes"/"no" the edit action accepts.
func BgpAdminShutdownArg(shutdown bool) string {
	if shutdown {
		return "yes"
	}
	return "no"
}

// EditExternalDeviceConnBgpAdminShutdown applies a connection-wide BGP admin
// shutdown, or lifts one.
//
// One call covers every tunnel of the connection, including tunnels added while
// it is shut down. The IPsec tunnels are not touched: the neighbors stay
// configured and the gateway holds the sessions down with
// "neighbor <ip> shutdown", which is why the connection keeps reporting a BGP
// status while shut down.
func (c *Client) EditExternalDeviceConnBgpAdminShutdown(edit *EditExternalDeviceConnBgpAdminShutdown) error {
	edit.CID = c.CID
	edit.Action = "edit_external_device_conn"

	// PostAPI2, not PostAPI: this action is only reachable under /v2/api.
	return c.PostAPI2(edit.Action, edit, BasicCheck)
}

// GetExternalDeviceConnTunnels returns the connection's tunnels, each carrying
// the Name that EditExternalDeviceConnTunnel selects by.
//
// The order is the controller's own tunnel order, which is not the order of the
// connection-level CIDR lists - see ExternalDeviceConnTunnelName.
//
// The per-tunnel BgpLocalIP / BgpRemoteIP that ExternalDeviceConnTunnelName
// matches on come from Tunnel.show_details on the controller, which only began
// reporting them in AVX-80972. Against an older controller every tunnel comes
// back with both empty and no tunnel can be named; ExternalDeviceConnTunnelName
// detects that and says so rather than reporting an unknown address.
func (c *Client) GetExternalDeviceConnTunnels(connName, vpcID string) ([]TunnelInfo, error) {
	params := map[string]string{
		"CID":       c.CID,
		"action":    "get_site2cloud_conn_detail",
		"conn_name": connName,
		"vpc_id":    vpcID,
	}
	checkFunc := func(action, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "does not exist") {
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s %s failed: %s", action, method, reason)
		}
		return nil
	}

	var data ExternalDeviceConnDetailResp
	if err := c.GetAPI(&data, params["action"], params, checkFunc); err != nil {
		return nil, err
	}
	if len(data.Results.Connections.ConnectionName) == 0 {
		return nil, ErrNotFound
	}

	return data.Results.Connections.Tunnels, nil
}

// ExternalDeviceConnTunnelName returns the name of the tunnel whose current
// inner addresses are localIP and remoteIP.
//
// Matching is by value rather than by position on purpose. The connection-level
// lists terraform stores are ordered by mirror slot and not by tunnel index: on
// a four-tunnel activemesh connection local_tunnel_cidr holds tunnels[0] and
// tunnels[3] while backup_local_tunnel_cidr holds tunnels[2] and tunnels[1].
// Reproducing that table here would mean keeping a copy of a layout that exists
// for backwards compatibility, so the position is resolved through the value it
// holds instead.
//
// remoteIP is tried first because it is the FRR neighbor address and so is
// unique across the gateway; localIP is a fallback for the case where the
// caller only knows its own side.
func ExternalDeviceConnTunnelName(tunnels []TunnelInfo, localIP, remoteIP string) (string, error) {
	tunnel, err := ExternalDeviceConnTunnel(tunnels, localIP, remoteIP)
	if err != nil {
		return "", err
	}
	return tunnel.Name, nil
}

// ExternalDeviceConnTunnel is the whole-tunnel form of
// ExternalDeviceConnTunnelName, for a caller that needs the tunnel's current
// PeerIP to tell a real edit from a write of what is already there.
func ExternalDeviceConnTunnel(tunnels []TunnelInfo, localIP, remoteIP string) (TunnelInfo, error) {
	if len(tunnels) == 0 {
		return TunnelInfo{}, fmt.Errorf("connection has no tunnels to edit")
	}

	for _, key := range []struct {
		want  string
		field func(TunnelInfo) string
		label string
	}{
		{remoteIP, func(t TunnelInfo) string { return t.BgpRemoteIP }, "remote"},
		{localIP, func(t TunnelInfo) string { return t.BgpLocalIP }, "local"},
	} {
		if strings.TrimSpace(key.want) == "" {
			continue
		}
		var matches []TunnelInfo
		for _, t := range tunnels {
			if SameBgpIP(key.field(t), key.want) {
				matches = append(matches, t)
			}
		}
		if len(matches) == 1 {
			return matches[0], nil
		}
		if len(matches) > 1 {
			return TunnelInfo{}, fmt.Errorf("%d tunnels share the %s BGP IP %q; refusing to guess which one to edit",
				len(matches), key.label, strings.TrimSpace(key.want))
		}
	}

	// Neither address matched. An empty BgpRemoteIP on every tunnel means the
	// controller is not reporting per-tunnel inner IPs, which is a different
	// problem from a genuinely unknown address, so say which one it is.
	reporting := false
	for _, t := range tunnels {
		if t.BgpRemoteIP != "" || t.BgpLocalIP != "" {
			reporting = true
			break
		}
	}
	if !reporting {
		return TunnelInfo{}, fmt.Errorf("connection does not report per-tunnel BGP IPs, so the tunnel to edit " +
			"cannot be identified; in-place edit needs a controller that reports them")
	}

	return TunnelInfo{}, fmt.Errorf("no tunnel has BGP IPs local %q / remote %q; the connection's tunnels are %s",
		strings.TrimSpace(localIP), strings.TrimSpace(remoteIP), describeTunnelBgpIPs(tunnels))
}

// describeTunnelBgpIPs lists the tunnels by name and inner addresses, so that an
// error tells the operator what could have been selected.
func describeTunnelBgpIPs(tunnels []TunnelInfo) string {
	parts := make([]string, 0, len(tunnels))
	for _, t := range tunnels {
		parts = append(parts, fmt.Sprintf("%s (local %s, remote %s)", t.Name, t.BgpLocalIP, t.BgpRemoteIP))
	}
	return strings.Join(parts, ", ")
}
