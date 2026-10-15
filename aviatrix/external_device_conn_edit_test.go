package aviatrix

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

// fakeEdcDiff stands in for *schema.ResourceData. A key missing from want is
// unchanged, which is how a real diff reports every attribute the plan does not
// touch, so a test case only has to name the attributes it changes.
type fakeEdcDiff struct {
	state map[string]string
	want  map[string]string
	// boolState and boolWant hold the TypeBool attributes. They are separate
	// from the string maps because getBool panics on anything but a bool, so a
	// bool attribute cannot be spelled as "true" here and still be read the way
	// the production code reads it.
	boolState map[string]bool
	boolWant  map[string]bool
}

func (f fakeEdcDiff) isBool(key string) bool {
	if _, ok := f.boolState[key]; ok {
		return true
	}
	_, ok := f.boolWant[key]
	return ok
}

func (f fakeEdcDiff) old(key string) string { return f.state[key] }

func (f fakeEdcDiff) newV(key string) string {
	if v, ok := f.want[key]; ok {
		return v
	}
	return f.state[key]
}

func (f fakeEdcDiff) oldBool(key string) bool { return f.boolState[key] }

func (f fakeEdcDiff) newBool(key string) bool {
	if v, ok := f.boolWant[key]; ok {
		return v
	}
	return f.boolState[key]
}

func (f fakeEdcDiff) Get(key string) any {
	if f.isBool(key) {
		return f.newBool(key)
	}
	return f.newV(key)
}

func (f fakeEdcDiff) GetChange(key string) (any, any) {
	if f.isBool(key) {
		return f.oldBool(key), f.newBool(key)
	}
	return f.old(key), f.newV(key)
}

func (f fakeEdcDiff) HasChange(key string) bool {
	if f.isBool(key) {
		return f.oldBool(key) != f.newBool(key)
	}
	return f.old(key) != f.newV(key)
}

func (f fakeEdcDiff) HasChanges(keys ...string) bool {
	for _, k := range keys {
		if f.HasChange(k) {
			return true
		}
	}
	return false
}

// haState is a two-tunnel HA connection: one tunnel per side.
func haState() map[string]string {
	return map[string]string{
		"connection_type":           "bgp",
		"tunnel_protocol":           "IPsec",
		"local_tunnel_cidr":         "169.254.0.1/30",
		"remote_tunnel_cidr":        "169.254.0.2/30",
		"remote_gateway_ip":         "203.0.113.1",
		"bgp_remote_as_num":         "65002",
		"pre_shared_key":            "primary-key",
		"backup_local_tunnel_cidr":  "169.254.1.1/30",
		"backup_remote_tunnel_cidr": "169.254.1.2/30",
		"backup_remote_gateway_ip":  "203.0.113.2",
		"backup_bgp_remote_as_num":  "65002",
		"backup_pre_shared_key":     "backup-key",
	}
}

// activemeshState is the four-tunnel 2x2 layout: two tunnels per side, which
// terraform stores as two-element comma-separated lists.
func activemeshState() map[string]string {
	return map[string]string{
		"connection_type":           "bgp",
		"tunnel_protocol":           "IPsec",
		"local_tunnel_cidr":         "169.254.0.1/30,169.254.2.1/30",
		"remote_tunnel_cidr":        "169.254.0.2/30,169.254.2.2/30",
		"remote_gateway_ip":         "203.0.113.1",
		"bgp_remote_as_num":         "65002",
		"backup_local_tunnel_cidr":  "169.254.1.1/30,169.254.3.1/30",
		"backup_remote_tunnel_cidr": "169.254.1.2/30,169.254.3.2/30",
		"backup_remote_gateway_ip":  "203.0.113.2",
		"backup_bgp_remote_as_num":  "65002",
	}
}

func TestPlanExternalDeviceConnEdits(t *testing.T) {
	tests := []struct {
		name    string
		state   map[string]string
		want    map[string]string
		expect  []edcTunnelEdit
		errWant string
	}{
		{
			name:  "no change plans nothing",
			state: haState(),
		},
		{
			name:  "remote ASN change edits only its own side",
			state: haState(),
			want:  map[string]string{"bgp_remote_as_num": "65003"},
			expect: []edcTunnelEdit{{
				side:          "primary",
				matchLocalIP:  "169.254.0.1/30",
				matchRemoteIP: "169.254.0.2/30",
				edit:          goaviatrix.EditExternalDeviceConnTunnel{BgpRemoteAsNumber: "65003"},
			}},
		},
		{
			name:  "PSK change only touches the side it belongs to",
			state: haState(),
			want:  map[string]string{"backup_pre_shared_key": "new-backup-key"},
			expect: []edcTunnelEdit{{
				side:          "backup",
				matchLocalIP:  "169.254.1.1/30",
				matchRemoteIP: "169.254.1.2/30",
				edit:          goaviatrix.EditExternalDeviceConnTunnel{PreSharedKey: "new-backup-key"},
			}},
		},
		{
			name:  "both inner addresses change together",
			state: haState(),
			want: map[string]string{
				"local_tunnel_cidr":  "169.254.9.1/30",
				"remote_tunnel_cidr": "169.254.9.2/30",
			},
			expect: []edcTunnelEdit{{
				side:          "primary",
				matchLocalIP:  "169.254.0.1/30",
				matchRemoteIP: "169.254.0.2/30",
				edit: goaviatrix.EditExternalDeviceConnTunnel{
					BgpLocalIP:  "169.254.9.1/30",
					BgpRemoteIP: "169.254.9.2/30",
				},
			}},
		},
		{
			name:  "a one-sided address change sends only that side",
			state: haState(),
			want:  map[string]string{"remote_tunnel_cidr": "169.254.9.2/30"},
			expect: []edcTunnelEdit{{
				side:          "primary",
				matchLocalIP:  "169.254.0.1/30",
				matchRemoteIP: "169.254.0.2/30",
				edit:          goaviatrix.EditExternalDeviceConnTunnel{BgpRemoteIP: "169.254.9.2/30"},
			}},
		},
		{
			name:  "remote gateway IP change becomes a remote_device_ip edit",
			state: haState(),
			want:  map[string]string{"remote_gateway_ip": "203.0.113.9"},
			expect: []edcTunnelEdit{{
				side:          "primary",
				matchLocalIP:  "169.254.0.1/30",
				matchRemoteIP: "169.254.0.2/30",
				edit:          goaviatrix.EditExternalDeviceConnTunnel{RemoteDeviceIP: "203.0.113.9"},
			}},
		},
		{
			name:  "only the position that changed is edited",
			state: activemeshState(),
			want:  map[string]string{"remote_tunnel_cidr": "169.254.0.2/30,169.254.9.2/30"},
			expect: []edcTunnelEdit{{
				side:          "primary",
				position:      1,
				matchLocalIP:  "169.254.2.1/30",
				matchRemoteIP: "169.254.2.2/30",
				edit:          goaviatrix.EditExternalDeviceConnTunnel{BgpRemoteIP: "169.254.9.2/30"},
			}},
		},
		{
			name:  "a per-side value reaches every tunnel of a four-tunnel side",
			state: activemeshState(),
			want:  map[string]string{"bgp_remote_as_num": "65003"},
			expect: []edcTunnelEdit{
				{
					side:          "primary",
					position:      0,
					matchLocalIP:  "169.254.0.1/30",
					matchRemoteIP: "169.254.0.2/30",
					edit:          goaviatrix.EditExternalDeviceConnTunnel{BgpRemoteAsNumber: "65003"},
				},
				{
					side:          "primary",
					position:      1,
					matchLocalIP:  "169.254.2.1/30",
					matchRemoteIP: "169.254.2.2/30",
					edit:          goaviatrix.EditExternalDeviceConnTunnel{BgpRemoteAsNumber: "65003"},
				},
			},
		},
		{
			name:  "spacing in the configured list is not sent to the controller",
			state: activemeshState(),
			want:  map[string]string{"remote_tunnel_cidr": "169.254.0.2/30, 169.254.9.2/30"},
			expect: []edcTunnelEdit{{
				side:          "primary",
				position:      1,
				matchLocalIP:  "169.254.2.1/30",
				matchRemoteIP: "169.254.2.2/30",
				edit:          goaviatrix.EditExternalDeviceConnTunnel{BgpRemoteIP: "169.254.9.2/30"},
			}},
		},
		{
			name:    "adding a tunnel is refused rather than mapped positionally",
			state:   haState(),
			want:    map[string]string{"remote_tunnel_cidr": "169.254.0.2/30,169.254.9.2/30"},
			errWant: "cannot add or remove tunnels",
		},
		{
			name:    "clearing the remote ASN is refused",
			state:   haState(),
			want:    map[string]string{"bgp_remote_as_num": ""},
			errWant: `cannot clear "bgp_remote_as_num"`,
		},
		{
			name:    "clearing the PSK is refused",
			state:   haState(),
			want:    map[string]string{"pre_shared_key": ""},
			errWant: `cannot clear "pre_shared_key"`,
		},
		{
			name:    "clearing the remote gateway IP is refused",
			state:   haState(),
			want:    map[string]string{"remote_gateway_ip": ""},
			errWant: `cannot clear "remote_gateway_ip"`,
		},
		{
			name: "editing a side that has no tunnel is refused",
			state: map[string]string{
				"connection_type":    "bgp",
				"tunnel_protocol":    "IPsec",
				"local_tunnel_cidr":  "169.254.0.1/30",
				"remote_tunnel_cidr": "169.254.0.2/30",
			},
			want:    map[string]string{"backup_bgp_remote_as_num": "65003"},
			errWant: "records no tunnel there",
		},
		{
			// An exact compare would plan an edit here, and the diff would never
			// clear: the tunnel would be re-established on every apply.
			name:  "dropping the prefix length is not a change",
			state: haState(),
			want: map[string]string{
				"local_tunnel_cidr":  "169.254.0.1",
				"remote_tunnel_cidr": "169.254.0.2",
			},
		},
		{
			name: "adding the prefix length to a bare address is not a change either",
			state: map[string]string{
				"connection_type":    "bgp",
				"tunnel_protocol":    "IPsec",
				"local_tunnel_cidr":  "169.254.0.1",
				"remote_tunnel_cidr": "169.254.0.2",
			},
			want: map[string]string{
				"local_tunnel_cidr":  "169.254.0.1/30",
				"remote_tunnel_cidr": "169.254.0.2/30",
			},
		},
		{
			name:    "swapping two addresses of one side is refused",
			state:   activemeshState(),
			want:    map[string]string{"local_tunnel_cidr": "169.254.2.1/30,169.254.0.1/30"},
			errWant: "cannot move 169.254.2.1/30 onto another tunnel of the primary side",
		},
		{
			name:    "a swap spelled without prefix lengths is still refused",
			state:   activemeshState(),
			want:    map[string]string{"remote_tunnel_cidr": "169.254.2.2,169.254.0.2"},
			errWant: "cannot move 169.254.2.2 onto another tunnel of the primary side",
		},
		{
			// A partial move collides exactly as a swap does: tunnel 0 would
			// take an address tunnel 1 still holds.
			name:    "a move onto an address another tunnel still holds is refused",
			state:   activemeshState(),
			want:    map[string]string{"local_tunnel_cidr": "169.254.2.1/30,169.254.9.1/30"},
			errWant: "cannot move 169.254.2.1/30 onto another tunnel of the primary side",
		},
		{
			name:  "moving both addresses somewhere nobody holds is planned",
			state: activemeshState(),
			want:  map[string]string{"local_tunnel_cidr": "169.254.5.1/30,169.254.9.1/30"},
			expect: []edcTunnelEdit{
				{
					side: "primary", position: 0,
					matchLocalIP:  "169.254.0.1/30",
					matchRemoteIP: "169.254.0.2/30",
					edit:          goaviatrix.EditExternalDeviceConnTunnel{BgpLocalIP: "169.254.5.1/30"},
				},
				{
					side: "primary", position: 1,
					matchLocalIP:  "169.254.2.1/30",
					matchRemoteIP: "169.254.2.2/30",
					edit:          goaviatrix.EditExternalDeviceConnTunnel{BgpLocalIP: "169.254.9.1/30"},
				},
			},
		},
		{
			name:    "clearing a tunnel address is refused",
			state:   activemeshState(),
			want:    map[string]string{"local_tunnel_cidr": ""},
			errWant: `cannot clear "local_tunnel_cidr"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := fakeEdcDiff{state: test.state, want: test.want}
			got, err := planExternalDeviceConnEdits(d, "gw1", "conn1")
			if test.errWant != "" {
				if err == nil {
					t.Fatalf("expected an error containing %q, got plan %+v", test.errWant, got)
				}
				if !strings.Contains(err.Error(), test.errWant) {
					t.Fatalf("expected an error containing %q, got %q", test.errWant, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(test.expect) {
				t.Fatalf("expected %d edit(s), got %d: %+v", len(test.expect), len(got), got)
			}
			for i, want := range test.expect {
				// Every edit carries the same gateway and connection, so the
				// cases above leave them out.
				want.edit.GwName = "gw1"
				want.edit.ConnectionName = "conn1"
				if got[i] != want {
					t.Errorf("edit %d:\n got %+v\nwant %+v", i, got[i], want)
				}
			}
		})
	}
}

// TestExternalDeviceConnTunnel covers the whole-tunnel form, which exists so a
// remote_device_ip edit can be skipped when the tunnel already has the value.
func TestExternalDeviceConnTunnel(t *testing.T) {
	tunnels := []goaviatrix.TunnelInfo{
		{Name: "conn1-1", BgpLocalIP: "169.254.0.1/30", BgpRemoteIP: "169.254.0.2/30", PeerIP: "203.0.113.1"},
		{Name: "conn1-2", BgpLocalIP: "169.254.1.1/30", BgpRemoteIP: "169.254.1.2/30", PeerIP: "203.0.113.2"},
	}

	got, err := goaviatrix.ExternalDeviceConnTunnel(tunnels, "", "169.254.1.2/30")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "conn1-2" || got.PeerIP != "203.0.113.2" {
		t.Errorf("expected tunnel conn1-2 with peer IP 203.0.113.2, got %+v", got)
	}

	// An error must not come back as a usable zero tunnel, or a caller would
	// read an empty PeerIP as "the tunnel has no peer IP" and skip an edit.
	if _, err := goaviatrix.ExternalDeviceConnTunnel(tunnels, "", "169.254.9.2/30"); err == nil {
		t.Error("expected an error for an address no tunnel has")
	}
}

func TestSameBgpIP(t *testing.T) {
	tests := []struct {
		a, b   string
		expect bool
	}{
		{a: "169.254.0.1/30", b: "169.254.0.1/30", expect: true},
		{a: "169.254.0.1/30", b: "169.254.0.1", expect: true},
		{a: " 169.254.0.1 ", b: "169.254.0.1/30", expect: true},
		{a: "169.254.0.1/30", b: "169.254.0.2/30"},
		// Two unset addresses are not "the same": equal would match an
		// unreported tunnel, or skip an edit nobody knows is needed.
		{a: "", b: ""},
		{a: "", b: "169.254.0.1/30"},
		{a: "/30", b: "/30"},
		// The outer peer IP never carries a prefix length.
		{a: "203.0.113.1", b: "203.0.113.1", expect: true},
		{a: "203.0.113.1", b: "203.0.113.2"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%q vs %q", test.a, test.b), func(t *testing.T) {
			if got := goaviatrix.SameBgpIP(test.a, test.b); got != test.expect {
				t.Errorf("expected SameBgpIP(%q, %q) = %t, got %t", test.a, test.b, test.expect, got)
			}
		})
	}
}

func TestExternalDeviceConnTunnelName(t *testing.T) {
	tunnels := []goaviatrix.TunnelInfo{
		{Name: "conn1-1", BgpLocalIP: "169.254.0.1/30", BgpRemoteIP: "169.254.0.2/30"},
		{Name: "conn1-2", BgpLocalIP: "169.254.1.1/30", BgpRemoteIP: "169.254.1.2/30"},
	}

	tests := []struct {
		name     string
		tunnels  []goaviatrix.TunnelInfo
		localIP  string
		remoteIP string
		want     string
		errWant  string
	}{
		{
			name:     "matches on the remote address",
			tunnels:  tunnels,
			localIP:  "169.254.1.1/30",
			remoteIP: "169.254.1.2/30",
			want:     "conn1-2",
		},
		{
			name:    "falls back to the local address",
			tunnels: tunnels,
			localIP: "169.254.0.1/30",
			want:    "conn1-1",
		},
		{
			name:     "a missing prefix length still matches",
			tunnels:  tunnels,
			remoteIP: " 169.254.0.2 ",
			want:     "conn1-1",
		},
		{
			name:     "an unknown address names what could have been selected",
			tunnels:  tunnels,
			remoteIP: "169.254.9.2/30",
			errWant:  "conn1-1 (local 169.254.0.1/30, remote 169.254.0.2/30)",
		},
		{
			name: "a controller that reports no per-tunnel IPs says so",
			tunnels: []goaviatrix.TunnelInfo{
				{Name: "conn1-1"}, {Name: "conn1-2"},
			},
			remoteIP: "169.254.0.2/30",
			errWant:  "does not report per-tunnel BGP IPs",
		},
		{
			name: "duplicate addresses are refused rather than guessed",
			tunnels: []goaviatrix.TunnelInfo{
				{Name: "conn1-1", BgpRemoteIP: "169.254.0.2/30"},
				{Name: "conn1-2", BgpRemoteIP: "169.254.0.2/30"},
			},
			remoteIP: "169.254.0.2/30",
			errWant:  "2 tunnels share the remote BGP IP",
		},
		{
			name:     "a connection with no tunnels is an error, not a nil name",
			tunnels:  nil,
			remoteIP: "169.254.0.2/30",
			errWant:  "no tunnels to edit",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := goaviatrix.ExternalDeviceConnTunnelName(test.tunnels, test.localIP, test.remoteIP)
			if test.errWant != "" {
				if err == nil {
					t.Fatalf("expected an error containing %q, got %q", test.errWant, got)
				}
				if !strings.Contains(err.Error(), test.errWant) {
					t.Fatalf("expected an error containing %q, got %q", test.errWant, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != test.want {
				t.Errorf("expected tunnel %q, got %q", test.want, got)
			}
		})
	}
}

// assertEdcEditableAttrsNotForceNew checks that a resource's schema really lets
// the in-place edit path run. A ForceNew left on any editable attribute would
// make terraform replace the connection instead of calling
// updateExternalDeviceConnInPlace, so the feature would quietly not exist - and
// nothing else in the provider would notice.
//
// Shared by the transit and spoke resources, which are the two that wire the
// edit path in.
func assertEdcEditableAttrsNotForceNew(t *testing.T, schemaMap map[string]*schema.Schema, resourceName string) {
	t.Helper()

	for _, attr := range append(slices.Clone(edcEditableAttrs), edcBgpAdminShutdownAttr) {
		field, ok := schemaMap[attr]
		if !ok {
			t.Errorf("%s: %q should exist in the schema", resourceName, attr)
			continue
		}
		if field.ForceNew {
			t.Errorf("%s: %q is ForceNew, so changing it would replace the connection rather than "+
				"be edited in place", resourceName, attr)
		}
	}

	// edcEditableConn reads these two with Get rather than GetChange, which is
	// only correct while they force replacement: a ForceNew attribute's current
	// value is also its old value.
	for _, attr := range []string{"connection_type", "tunnel_protocol"} {
		field, ok := schemaMap[attr]
		if !ok {
			t.Errorf("%s: %q should exist in the schema", resourceName, attr)
			continue
		}
		if !field.ForceNew {
			t.Errorf("%s: %q has to stay ForceNew: edcEditableConn treats its current value as its "+
				"old value when deciding whether the connection is editable", resourceName, attr)
		}
	}
}

// testAccCheckEdcIDUnchanged records the resource's ID the first time it runs and
// fails on any later run if the ID differs. That is how an acceptance test tells
// an in-place edit apart from a replacement: terraform reports both as a
// successful apply, but a replacement produces a new ID.
func testAccCheckEdcIDUnchanged(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("%s not found in state", n)
		}
		if *id == "" {
			*id = rs.Primary.ID
			return nil
		}
		if rs.Primary.ID != *id {
			return fmt.Errorf("%s was replaced: its ID went from %s to %s, but these attributes are "+
				"meant to be edited in place", n, *id, rs.Primary.ID)
		}
		return nil
	}
}

func TestValueForPosition(t *testing.T) {
	tests := []struct {
		name    string
		list    []string
		i, n    int
		want    string
		errWant string
	}{
		{name: "an unchanged attribute contributes nothing", list: nil, i: 0, n: 2},
		{name: "one value per tunnel is positional", list: []string{"a", "b"}, i: 1, n: 2, want: "b"},
		{name: "a single value covers every tunnel", list: []string{"a"}, i: 1, n: 2, want: "a"},
		{
			name:    "any other length is refused",
			list:    []string{"a", "b", "c"},
			i:       0,
			n:       2,
			errWant: "has 3 values but the connection has 2 tunnel(s)",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := valueForPosition(test.list, test.i, test.n, "remote_gateway_ip")
			if test.errWant != "" {
				if err == nil || !strings.Contains(err.Error(), test.errWant) {
					t.Fatalf("expected an error containing %q, got %v", test.errWant, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != test.want {
				t.Errorf("expected %q, got %q", test.want, got)
			}
		})
	}
}

func TestBgpAdminShutdownArg(t *testing.T) {
	if got := goaviatrix.BgpAdminShutdownArg(true); got != "yes" {
		t.Errorf("expected a shutdown to be sent as %q, got %q", "yes", got)
	}
	// Not "" and not omitted: the edit action distinguishes "bring the
	// neighbors back up" from "say nothing about the admin state", and the
	// form tag carries omitempty, so an empty string would be dropped from
	// the request and the lift would silently not happen.
	if got := goaviatrix.BgpAdminShutdownArg(false); got != "no" {
		t.Errorf("expected a lifted shutdown to be sent as %q, got %q", "no", got)
	}
}

func TestPlanExternalDeviceConnUpdate(t *testing.T) {
	tests := []struct {
		name      string
		want      map[string]string
		boolState map[string]bool
		boolWant  map[string]bool
		steps     []edcUpdateStep
	}{
		{
			name:  "an unchanged connection needs no steps",
			steps: nil,
		},
		{
			name:  "a per-tunnel change alone is one step",
			want:  map[string]string{"bgp_remote_as_num": "65100"},
			steps: []edcUpdateStep{edcStepTunnelEdits},
		},
		{
			name:      "a shutdown alone is one step",
			boolState: map[string]bool{"bgp_admin_shutdown": false},
			boolWant:  map[string]bool{"bgp_admin_shutdown": true},
			steps:     []edcUpdateStep{edcStepShutdown},
		},
		{
			name:      "lifting a shutdown alone is one step",
			boolState: map[string]bool{"bgp_admin_shutdown": true},
			boolWant:  map[string]bool{"bgp_admin_shutdown": false},
			steps:     []edcUpdateStep{edcStepNoShutdown},
		},
		{
			// The session goes down before the tunnel is re-keyed, rather than
			// being re-keyed while up and then shut.
			name:      "a shutdown is applied before the tunnel edits it accompanies",
			want:      map[string]string{"pre_shared_key": "rotated"},
			boolState: map[string]bool{"bgp_admin_shutdown": false},
			boolWant:  map[string]bool{"bgp_admin_shutdown": true},
			steps:     []edcUpdateStep{edcStepShutdown, edcStepTunnelEdits},
		},
		{
			// The mirror image: the session comes back up only once it has the
			// addresses the configuration asks for.
			name:      "lifting a shutdown is applied after the tunnel edits it accompanies",
			want:      map[string]string{"remote_tunnel_cidr": "169.254.9.2/30"},
			boolState: map[string]bool{"bgp_admin_shutdown": true},
			boolWant:  map[string]bool{"bgp_admin_shutdown": false},
			steps:     []edcUpdateStep{edcStepTunnelEdits, edcStepNoShutdown},
		},
		{
			name:      "an unchanged shutdown contributes no step of its own",
			want:      map[string]string{"remote_gateway_ip": "203.0.113.9"},
			boolState: map[string]bool{"bgp_admin_shutdown": true},
			boolWant:  map[string]bool{"bgp_admin_shutdown": true},
			steps:     []edcUpdateStep{edcStepTunnelEdits},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := fakeEdcDiff{
				state:     haState(),
				want:      test.want,
				boolState: test.boolState,
				boolWant:  test.boolWant,
			}

			got := planExternalDeviceConnUpdate(d)
			if len(got) != len(test.steps) {
				t.Fatalf("expected steps %v, got %v", test.steps, got)
			}
			for i := range got {
				if got[i] != test.steps[i] {
					t.Fatalf("expected steps %v, got %v", test.steps, got)
				}
			}
		})
	}
}

func TestEdcEditableConn(t *testing.T) {
	tests := []struct {
		name           string
		connType       string
		tunnelProtocol string
		errWant        string
	}{
		{
			name:           "bgp over ipsec is editable",
			connType:       "bgp",
			tunnelProtocol: "IPsec",
		},
		{
			// The controller compares case-insensitively and configurations are
			// written both ways, so neither spelling may force a replacement.
			name:           "the tunnel protocol is matched case-insensitively",
			connType:       "bgp",
			tunnelProtocol: "ipsec",
		},
		{
			name:           "a static connection is refused",
			connType:       "static",
			tunnelProtocol: "IPsec",
			errWant:        "only a BGP connection over IPsec can be edited",
		},
		{
			name:           "a LAN connection is refused",
			connType:       "bgp",
			tunnelProtocol: "LAN",
			errWant:        "only a BGP connection over IPsec can be edited",
		},
		{
			name:           "a GRE connection is refused",
			connType:       "bgp",
			tunnelProtocol: "GRE",
			errWant:        "only a BGP connection over IPsec can be edited",
		},
		{
			// Read normalises an empty protocol to "IPsec", so refusing it would
			// refuse to edit real IPsec connections.
			name:           "an empty tunnel protocol counts as IPsec",
			connType:       "bgp",
			tunnelProtocol: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := haState()
			state["connection_type"] = test.connType
			state["tunnel_protocol"] = test.tunnelProtocol
			d := fakeEdcDiff{state: state}

			err := edcEditableConn(d, "conn-1", []string{"bgp_admin_shutdown", "pre_shared_key"})
			if test.errWant == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), test.errWant) {
				t.Fatalf("expected an error containing %q, got %v", test.errWant, err)
			}
			// The refusal has to name the attributes that provoked it, so that
			// an operator can tell which part of the configuration to revert.
			for _, attr := range []string{"bgp_admin_shutdown", "pre_shared_key"} {
				if !strings.Contains(err.Error(), attr) {
					t.Errorf("expected the error to name %q, got %v", attr, err)
				}
			}
		})
	}
}

// TestEdcReplacementAttrs pins the rule that replaced the unconditional
// ForceNew. Getting it wrong fails quietly: it plans an in-place update the
// controller refuses, leaving a config unappliable without -replace.
func TestEdcReplacementAttrs(t *testing.T) {
	tests := []struct {
		name           string
		connType       string
		tunnelProtocol string
		want           map[string]string
		expect         []string
	}{
		{
			name:           "an editable connection replaces nothing",
			connType:       "bgp",
			tunnelProtocol: "IPsec",
			want:           map[string]string{"pre_shared_key": "rotated", "bgp_remote_as_num": "65003"},
		},
		{
			name:           "a LAN connection still replaces on every changed attribute",
			connType:       "bgp",
			tunnelProtocol: "LAN",
			want:           map[string]string{"bgp_remote_as_num": "65003", "backup_pre_shared_key": "rotated"},
			expect:         []string{"bgp_remote_as_num", "backup_pre_shared_key"},
		},
		{
			name:           "a static connection still replaces",
			connType:       "static",
			tunnelProtocol: "IPsec",
			want:           map[string]string{"remote_gateway_ip": "203.0.113.9"},
			expect:         []string{"remote_gateway_ip"},
		},
		{
			name:           "a GRE connection still replaces",
			connType:       "bgp",
			tunnelProtocol: "GRE",
			want:           map[string]string{"local_tunnel_cidr": "169.254.9.1/30"},
			expect:         []string{"local_tunnel_cidr"},
		},
		{
			// An unrelated edit to a LAN connection must not start replacing it.
			name:           "an unrelated change to a LAN connection replaces nothing",
			connType:       "bgp",
			tunnelProtocol: "LAN",
			want:           map[string]string{"connection_bgp_send_communities": "all"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := haState()
			state["connection_type"] = test.connType
			state["tunnel_protocol"] = test.tunnelProtocol

			got := edcReplacementAttrs(fakeEdcDiff{state: state, want: test.want})
			if !slices.Equal(got, test.expect) {
				t.Errorf("expected %v to require replacement, got %v", test.expect, got)
			}
		})
	}
}

func TestEdcCollidingMove(t *testing.T) {
	tests := []struct {
		name   string
		have   []string
		want   []string
		expect string // the colliding address, empty when there is none
	}{
		{
			name:   "a swap collides",
			have:   []string{"169.254.0.1/30", "169.254.2.1/30"},
			want:   []string{"169.254.2.1/30", "169.254.0.1/30"},
			expect: "169.254.2.1/30",
		},
		{
			name:   "a swap spelled without prefix lengths collides",
			have:   []string{"169.254.0.1/30", "169.254.2.1/30"},
			want:   []string{"169.254.2.1", "169.254.0.1"},
			expect: "169.254.2.1",
		},
		{
			// A partial move is the same hazard as a swap: position 0 takes an
			// address position 1 still holds, so the two overlap until the
			// second request lands.
			name:   "a partial move collides",
			have:   []string{"169.254.0.1/30", "169.254.2.1/30"},
			want:   []string{"169.254.2.1/30", "169.254.9.1/30"},
			expect: "169.254.2.1/30",
		},
		{
			name:   "a repeated value held elsewhere collides",
			have:   []string{"169.254.0.1/30", "169.254.0.1/30", "169.254.2.1/30"},
			want:   []string{"169.254.2.1/30", "169.254.2.1/30", "169.254.0.1/30"},
			expect: "169.254.2.1/30",
		},
		{
			name: "the same order does not collide",
			have: []string{"169.254.0.1/30", "169.254.2.1/30"},
			want: []string{"169.254.0.1/30", "169.254.2.1/30"},
		},
		{
			name: "addresses nobody holds do not collide",
			have: []string{"169.254.0.1/30", "169.254.2.1/30"},
			want: []string{"169.254.5.1/30", "169.254.9.1/30"},
		},
		{
			name: "one changed address nobody holds does not collide",
			have: []string{"169.254.0.1/30", "169.254.2.1/30"},
			want: []string{"169.254.0.1/30", "169.254.9.1/30"},
		},
		{
			name: "a single-tunnel side cannot collide",
			have: []string{"169.254.0.1/30"},
			want: []string{"169.254.0.1/30"},
		},
		{
			name: "a list of a different length is left to the length check",
			have: []string{"169.254.0.1/30", "169.254.2.1/30"},
			want: []string{"169.254.2.1/30"},
		},
		{
			name: "an empty new list is left to the clear check",
			have: []string{"169.254.0.1/30", "169.254.2.1/30"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			addr, ok := edcCollidingMove(test.have, test.want)
			if ok != (test.expect != "") {
				t.Errorf("expected edcCollidingMove(%v, %v) to report %t, got %t",
					test.have, test.want, test.expect != "", ok)
			}
			if addr != test.expect {
				t.Errorf("expected the colliding address %q, got %q", test.expect, addr)
			}
		})
	}
}

func TestEdcChangedAttrs(t *testing.T) {
	d := fakeEdcDiff{
		state:     haState(),
		want:      map[string]string{"pre_shared_key": "rotated", "backup_remote_gateway_ip": "203.0.113.9"},
		boolState: map[string]bool{"bgp_admin_shutdown": false},
		boolWant:  map[string]bool{"bgp_admin_shutdown": true},
	}

	// Reported in the order of the list passed in, not in diff order, so the
	// message reads the same way every time.
	got := edcChangedAttrs(d, append(slices.Clone(edcEditableAttrs), edcBgpAdminShutdownAttr)...)
	want := []string{"pre_shared_key", "backup_remote_gateway_ip", "bgp_admin_shutdown"}
	if !slices.Equal(got, want) {
		t.Errorf("expected changed attributes %v, got %v", want, got)
	}
}

// TestEdcCustomizeDiffBgpOverLan pins the two refusals a BGP-over-LAN
// connection must produce, through the resource's CustomizeDiff: a plan-time
// error for bgp_admin_shutdown, which replacement cannot express either, and a
// forced replacement for bgp_remote_as_num, which would otherwise promise an
// in-place edit the controller refuses.
//
// TestEdcEditableConn covers the protocol check itself; this covers the wiring
// around it, including the create-time branch fakeEdcDiff cannot reach.
func TestEdcCustomizeDiffBgpOverLan(t *testing.T) {
	lanConn := map[string]any{
		"vpc_id":            "vpc-1",
		"gw_name":           "transit1",
		"connection_name":   "bgpolan-edit-negative",
		"connection_type":   "bgp",
		"tunnel_protocol":   "LAN",
		"bgp_local_as_num":  "65100",
		"bgp_remote_as_num": "65201",
		"remote_lan_ip":     "10.60.0.21",
	}

	// Lets a case name only what it changes.
	with := func(base map[string]any, changes map[string]any) map[string]any {
		out := make(map[string]any, len(base)+len(changes))
		for k, v := range base {
			out[k] = v
		}
		for k, v := range changes {
			out[k] = v
		}
		return out
	}

	// An already-applied connection. State() is nil without an ID, and
	// TestResourceDataRaw never sets one.
	priorState := func(t *testing.T, res *schema.Resource, raw map[string]any) *terraform.InstanceState {
		t.Helper()
		d := schema.TestResourceDataRaw(t, res.Schema, raw)
		d.SetId("bgpolan-edit-negative~vpc-1")
		state := d.State()
		if state == nil {
			t.Fatal("expected a non-nil state from TestResourceDataRaw once an ID is set")
		}
		return state
	}

	t.Run("bgp_admin_shutdown on an existing LAN connection is a plan-time error", func(t *testing.T) {
		res := resourceAviatrixTransitExternalDeviceConn()
		state := priorState(t, res, lanConn)
		config := terraform.NewResourceConfigRaw(with(lanConn, map[string]any{"bgp_admin_shutdown": true}))

		_, err := res.Diff(context.Background(), state, config, nil)
		if err == nil {
			t.Fatal("expected the diff to be refused, got no error")
		}
		for _, want := range []string{
			"only a BGP connection over IPsec can be edited",
			edcBgpAdminShutdownAttr,
			"bgpolan-edit-negative",
			"LAN",
		} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("expected the error to mention %q, got %v", want, err)
			}
		}
	})

	t.Run("bgp_admin_shutdown at create time is refused before anything exists", func(t *testing.T) {
		res := resourceAviatrixTransitExternalDeviceConn()
		config := terraform.NewResourceConfigRaw(with(lanConn, map[string]any{"bgp_admin_shutdown": true}))

		// Nil prior state is a create; catching it here avoids a created-then-
		// tainted connection.
		_, err := res.Diff(context.Background(), nil, config, nil)
		if err == nil {
			t.Fatal("expected a create carrying bgp_admin_shutdown over LAN to be refused, got no error")
		}
		if !strings.Contains(err.Error(), "only a BGP connection over IPsec can be edited") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("bgp_remote_as_num on a LAN connection forces a replacement", func(t *testing.T) {
		res := resourceAviatrixTransitExternalDeviceConn()
		state := priorState(t, res, lanConn)
		config := terraform.NewResourceConfigRaw(with(lanConn, map[string]any{"bgp_remote_as_num": "65202"}))

		diff, err := res.Diff(context.Background(), state, config, nil)
		if err != nil {
			t.Fatalf("a replacement is a legitimate plan, so Diff must not error: %v", err)
		}
		attr, ok := diff.Attributes["bgp_remote_as_num"]
		if !ok {
			t.Fatalf("expected a diff entry for bgp_remote_as_num, attributes: %#v", diff.Attributes)
		}
		if !attr.RequiresNew {
			t.Errorf("expected bgp_remote_as_num to require a new connection on a LAN connection, got %#v", attr)
		}
	})

	t.Run("the same change over IPsec is applied in place", func(t *testing.T) {
		res := resourceAviatrixTransitExternalDeviceConn()
		ipsecConn := with(lanConn, map[string]any{
			"tunnel_protocol":    "IPsec",
			"remote_lan_ip":      "",
			"remote_gateway_ip":  "172.12.13.14",
			"local_tunnel_cidr":  "169.254.10.1/30",
			"remote_tunnel_cidr": "169.254.10.2/30",
		})
		state := priorState(t, res, ipsecConn)
		config := terraform.NewResourceConfigRaw(with(ipsecConn, map[string]any{"bgp_remote_as_num": "65202"}))

		diff, err := res.Diff(context.Background(), state, config, nil)
		if err != nil {
			t.Fatalf("Diff returned an error: %v", err)
		}
		attr, ok := diff.Attributes["bgp_remote_as_num"]
		if !ok {
			t.Fatalf("expected a diff entry for bgp_remote_as_num, attributes: %#v", diff.Attributes)
		}
		if attr.RequiresNew {
			t.Error("bgp_remote_as_num must be editable in place on a BGP-over-IPsec connection")
		}
	})
}
