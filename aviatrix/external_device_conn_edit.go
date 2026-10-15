package aviatrix

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

// edcEditableAttrs are the external device connection attributes the in-place
// edit API can change. Every one of them used to be ForceNew.
//
// Only BGP-over-IPsec connections are editable this way, so a change to any of
// these on a static, LAN or GRE connection is still an error rather than a
// silent no-op - the connection has to be recreated in that case, and saying so
// is better than appearing to apply an edit the controller would refuse.
var edcEditableAttrs = []string{
	"bgp_remote_as_num",
	"local_tunnel_cidr",
	"remote_tunnel_cidr",
	"remote_gateway_ip",
	"pre_shared_key",
	"backup_bgp_remote_as_num",
	"backup_local_tunnel_cidr",
	"backup_remote_tunnel_cidr",
	"backup_remote_gateway_ip",
	"backup_pre_shared_key",
}

// edcBgpAdminShutdownAttr is the connection-wide BGP admin state.
//
// Deliberately not a member of edcEditableAttrs: everything in that list
// addresses one tunnel and is sent with a tunnel selector, whereas this covers
// every tunnel of the connection and the controller rejects a request that
// carries both. It therefore travels in its own request, and the two are only
// sequenced together - see updateExternalDeviceConnInPlace.
const edcBgpAdminShutdownAttr = "bgp_admin_shutdown"

// edcSide names the attributes describing one side of a connection. Terraform
// models HA as a primary set plus a parallel backup_* set, whereas the edit API
// addresses one tunnel at a time by name; a side is the bridge between the two.
type edcSide struct {
	label      string
	asn        string
	localCidr  string
	remoteCidr string
	peerIP     string
	psk        string
}

var edcSides = []edcSide{
	{
		label:      "primary",
		asn:        "bgp_remote_as_num",
		localCidr:  "local_tunnel_cidr",
		remoteCidr: "remote_tunnel_cidr",
		peerIP:     "remote_gateway_ip",
		psk:        "pre_shared_key",
	},
	{
		label:      "backup",
		asn:        "backup_bgp_remote_as_num",
		localCidr:  "backup_local_tunnel_cidr",
		remoteCidr: "backup_remote_tunnel_cidr",
		peerIP:     "backup_remote_gateway_ip",
		psk:        "backup_pre_shared_key",
	},
}

// edcDiff is the part of *schema.ResourceData the planner reads. Narrowing it to
// an interface keeps the planning rules testable directly, rather than only
// through a hand-built terraform.InstanceDiff.
type edcDiff interface {
	Get(key string) any
	GetChange(key string) (any, any)
	HasChange(key string) bool
	HasChanges(keys ...string) bool
}

// edcTunnelEdit is one planned edit, still keyed by the addresses the tunnel has
// before the update runs rather than by its name. The name is resolved as late
// as possible, because a remote_device_ip edit re-keys the tunnel and so changes
// the name of every tunnel it touches.
type edcTunnelEdit struct {
	side          string
	position      int
	matchLocalIP  string
	matchRemoteIP string
	edit          goaviatrix.EditExternalDeviceConnTunnel
}

// valueForPosition picks the value a list supplies for tunnel i of n.
//
// A list as long as the side is positional. A single value applies to every
// tunnel on the side, which is how the per-side attributes that are not lists
// at all - the ASN and the PSK - reach each tunnel. Any other length is a
// mistake worth reporting, since guessing would edit the wrong tunnel.
func valueForPosition(list []string, i, n int, attr string) (string, error) {
	switch {
	case len(list) == 0:
		return "", nil
	case len(list) == n:
		return list[i], nil
	case len(list) == 1:
		return list[0], nil
	default:
		return "", fmt.Errorf("%q has %d values but the connection has %d tunnel(s) on this side; "+
			"specify either one value per tunnel, separated by \",\", or a single value for all of them",
			attr, len(list), n)
	}
}

// edcCollidingMove returns the first new address that another tunnel of the
// same side currently holds.
//
// Refused rather than applied because the requests go out one per tunnel, in
// order: writing that address leaves the connection with two tunnels holding it
// until the request that moves the other one off lands. If anything interrupts
// the apply in between, the connection is left in that state, and the next
// apply cannot repair it - ExternalDeviceConnTunnelName resolves a tunnel by the
// address it holds and refuses a duplicate ("2 tunnels share ..."), so every
// subsequent edit of that side fails until the addresses are fixed by hand.
//
// A full swap ([A,B] -> [B,A]) and a partial move ([A,B] -> [B,C]) are both
// caught; moving to addresses nobody holds is not.
func edcCollidingMove(have, want []string) (string, bool) {
	if len(have) < 2 || len(have) != len(want) {
		return "", false
	}
	for i := range want {
		if goaviatrix.SameBgpIP(have[i], want[i]) {
			continue
		}
		for j := range have {
			if j != i && goaviatrix.SameBgpIP(have[j], want[i]) {
				return want[i], true
			}
		}
	}
	return "", false
}

// planExternalDeviceConnEdits turns the diff into one edit per tunnel that has
// at least one changed field. Nothing is sent to the controller here.
func planExternalDeviceConnEdits(d edcDiff, gwName, connName string) ([]edcTunnelEdit, error) {
	var plan []edcTunnelEdit

	for _, side := range edcSides {
		oldRemoteRaw, newRemoteRaw := d.GetChange(side.remoteCidr)
		oldLocalRaw, newLocalRaw := d.GetChange(side.localCidr)
		oldRemote := goaviatrix.SplitTrimmed(mustString(oldRemoteRaw))
		newRemote := goaviatrix.SplitTrimmed(mustString(newRemoteRaw))
		oldLocal := goaviatrix.SplitTrimmed(mustString(oldLocalRaw))
		newLocal := goaviatrix.SplitTrimmed(mustString(newLocalRaw))

		// The tunnel count comes from the old state, which is what the tunnels
		// currently in place are described by. A new list of a different length
		// would be asking to add or remove tunnels, which this API cannot do.
		// Either list covers the whole side, so take whichever state recorded.
		n := cmp.Or(len(oldRemote), len(oldLocal))
		if n == 0 {
			// No tunnels recorded on this side; nothing addressable here. A
			// change to a backup_* attribute of a connection that has no backup
			// side lands here and is caught below.
			if d.HasChanges(side.asn, side.localCidr, side.remoteCidr, side.peerIP, side.psk) {
				return nil, fmt.Errorf("cannot edit the %s side of connection %q: its state records no "+
					"tunnel there, so there is nothing to edit", side.label, connName)
			}
			continue
		}
		for _, l := range []struct {
			list []string
			attr string
			old  []string
		}{{newRemote, side.remoteCidr, oldRemote}, {newLocal, side.localCidr, oldLocal}} {
			// An emptied list would otherwise skip the length check below, plan
			// no field edit at all, and leave an apply that reports success
			// while Read restores the old address - a diff that never clears.
			if d.HasChange(l.attr) && len(l.list) == 0 && len(l.old) != 0 {
				return nil, fmt.Errorf("cannot clear %q: an in-place edit can change a tunnel address but "+
					"cannot remove one. Replace the connection instead", l.attr)
			}
			if len(l.list) != 0 && len(l.list) != n {
				return nil, fmt.Errorf("%q has %d values but the %s side of connection %q has %d tunnel(s); "+
					"an in-place edit can change a tunnel's addresses but cannot add or remove tunnels",
					l.attr, len(l.list), side.label, connName, n)
			}
			if addr, ok := edcCollidingMove(l.old, l.list); ok {
				return nil, fmt.Errorf("cannot move %s onto another tunnel of the %s side of connection %q: "+
					"%q assigns it to a tunnel other than the one holding it now. The tunnels are edited one "+
					"request at a time, so part way through the apply two of them would hold this address, and "+
					"an apply interrupted there leaves a connection no later edit can address; change the "+
					"addresses one at a time, or replace the connection",
					addr, side.label, connName, l.attr)
			}
		}

		asn, psk, peerIPs, err := edcChangedSideValues(d, side)
		if err != nil {
			return nil, err
		}

		for i := range n {
			item := edcTunnelEdit{
				side:     side.label,
				position: i,
				edit: goaviatrix.EditExternalDeviceConnTunnel{
					GwName:            gwName,
					ConnectionName:    connName,
					BgpRemoteAsNumber: asn,
					PreSharedKey:      psk,
				},
			}
			if i < len(oldRemote) {
				item.matchRemoteIP = oldRemote[i]
			}
			if i < len(oldLocal) {
				item.matchLocalIP = oldLocal[i]
			}
			if item.matchRemoteIP == "" && item.matchLocalIP == "" {
				return nil, fmt.Errorf("cannot edit tunnel %d of the %s side of connection %q: its state "+
					"records neither a local nor a remote tunnel address, so the tunnel cannot be identified",
					i+1, side.label, connName)
			}

			// Only send an address that actually changed. A one-sided change is
			// legal: the controller validates the supplied address against the
			// tunnel's existing partner address. Compared with SameBgpIP, not
			// ==, so a config spelled without the /30 is not a change against
			// the prefixed form the controller reports - an exact compare there
			// re-sends the edit, and bounces the session, on every apply.
			if d.HasChange(side.remoteCidr) && i < len(newRemote) &&
				!goaviatrix.SameBgpIP(newRemote[i], item.matchRemoteIP) {
				item.edit.BgpRemoteIP = newRemote[i]
			}
			if d.HasChange(side.localCidr) && i < len(newLocal) &&
				!goaviatrix.SameBgpIP(newLocal[i], item.matchLocalIP) {
				item.edit.BgpLocalIP = newLocal[i]
			}
			if peerIP, err := valueForPosition(peerIPs, i, n, side.peerIP); err != nil {
				return nil, err
			} else {
				item.edit.RemoteDeviceIP = peerIP
			}

			if !item.edit.IsEmpty() {
				plan = append(plan, item)
			}
		}
	}

	return plan, nil
}

// edcChangedSideValues returns the new ASN, PSK and outer peer IPs for a side,
// each empty when unchanged. The edit API treats an omitted field as "leave
// alone", so an unchanged field is simply not sent.
func edcChangedSideValues(d edcDiff, side edcSide) (asn, psk string, peerIPs []string, err error) {
	if d.HasChange(side.asn) {
		asn = strings.TrimSpace(getString(d, side.asn))
		if asn == "" {
			return "", "", nil, fmt.Errorf("cannot clear %q: an in-place edit can change the remote ASN "+
				"of a BGP connection but cannot remove it", side.asn)
		}
	}
	if d.HasChange(side.psk) {
		psk = getString(d, side.psk)
		if psk == "" {
			// An empty PSK at create time means "generate one", but on an edit it
			// would ask the controller to set an empty key. Recreating the
			// connection is the only way back to a generated PSK.
			return "", "", nil, fmt.Errorf("cannot clear %q: an in-place edit can change the pre-shared key "+
				"but cannot ask for a new one to be generated", side.psk)
		}
	}
	if d.HasChange(side.peerIP) {
		peerIPs = goaviatrix.SplitTrimmed(getString(d, side.peerIP))
		if len(peerIPs) == 0 {
			return "", "", nil, fmt.Errorf("cannot clear %q: an in-place edit can change the remote "+
				"gateway IP but cannot remove it", side.peerIP)
		}
	}
	return asn, psk, peerIPs, nil
}

// edcEditableConn reports whether the edit API can touch this connection at
// all, and explains the refusal when it cannot.
//
// The edit action is limited to BGP over IPsec: the controller refuses LAN and
// GRE connections outright, and refuses the BGP fields on a connection whose
// routing is static. Saying so here is better than letting the request go and
// reporting whatever the controller said, because the remedy is different - the
// connection has to be replaced rather than edited - and because the diff that
// got us here would otherwise look like an edit terraform could have made.
//
// connection_type and tunnel_protocol are themselves ForceNew, so their current
// values are also their old values and this needs no GetChange.
func edcEditableConn(d edcDiff, connName string, changed []string) error {
	if edcConnIsEditable(d) {
		return nil
	}
	return fmt.Errorf("cannot update %s in place: only a BGP connection over IPsec can be edited, "+
		"and connection %q is %q over %q. Changing these attributes on this connection requires "+
		"replacing it", strings.Join(changed, ", "), connName,
		getString(d, "connection_type"), getString(d, "tunnel_protocol"))
}

// edcConnIsEditable is edcEditableConn without the explanation.
//
// An empty tunnel_protocol counts as IPsec: Read normalises the controller's
// empty value to "IPsec", so a state written before that reads back empty on a
// connection that is in fact IPsec.
func edcConnIsEditable(d edcDiff) bool {
	tunnelProtocol := getString(d, "tunnel_protocol")
	return getString(d, "connection_type") == "bgp" &&
		(tunnelProtocol == "" || strings.EqualFold(tunnelProtocol, "IPsec"))
}

// edcReplacementAttrs names the changed attributes that still have to replace
// the connection; empty on a BGP-over-IPsec one, where all of them are editable.
//
// Since they stopped being unconditionally ForceNew, without this a LAN
// connection whose bgp_remote_as_num changed would plan an in-place update, be
// refused at apply time, and leave the configuration unappliable without
// -replace.
func edcReplacementAttrs(d edcDiff) []string {
	if edcConnIsEditable(d) {
		return nil
	}
	return edcChangedAttrs(d, edcEditableAttrs...)
}

// edcCustomizeDiff decides at plan time which changes still replace the
// connection, and rejects the ones that can neither be applied nor replaced.
// Shared by the transit and spoke resources.
//
// bgp_admin_shutdown cannot be handled by replacement - create has no such
// argument - so on a connection that cannot be edited it is a plan-time error.
// That also keeps an unsupported configuration from creating the connection and
// only then failing, which taints the resource.
func edcCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ any) error {
	// Being created: nothing to replace, and the per-tunnel attributes are just
	// part of the create.
	if d.Id() == "" {
		if getBool(d, edcBgpAdminShutdownAttr) {
			return edcEditableConn(d, getString(d, "connection_name"), []string{edcBgpAdminShutdownAttr})
		}
		return nil
	}

	if d.HasChange(edcBgpAdminShutdownAttr) {
		if err := edcEditableConn(d, getString(d, "connection_name"), []string{edcBgpAdminShutdownAttr}); err != nil {
			return err
		}
	}

	for _, attr := range edcReplacementAttrs(d) {
		if err := d.ForceNew(attr); err != nil {
			return fmt.Errorf("could not mark %q as requiring connection %q to be replaced: %w",
				attr, getString(d, "connection_name"), err)
		}
	}

	return nil
}

// edcChangedAttrs names the attributes of attrs the diff changes, for an error
// message that tells the operator which part of their configuration is at
// fault.
func edcChangedAttrs(d edcDiff, attrs ...string) []string {
	var changed []string
	for _, attr := range attrs {
		if d.HasChange(attr) {
			changed = append(changed, attr)
		}
	}
	return changed
}

// applyEdcBgpAdminShutdown shuts the connection's BGP neighbors down, or brings
// them back up.
func applyEdcBgpAdminShutdown(client *goaviatrix.Client, gwName, connName string, shutdown bool) error {
	log.Printf("[INFO] Setting bgp_admin_shutdown=%t on external device conn %s", shutdown, connName)
	edit := &goaviatrix.EditExternalDeviceConnBgpAdminShutdown{
		GwName:           gwName,
		ConnectionName:   connName,
		BgpAdminShutdown: goaviatrix.BgpAdminShutdownArg(shutdown),
	}
	if err := client.EditExternalDeviceConnBgpAdminShutdown(edit); err != nil {
		verb := "bring the BGP neighbors of"
		if shutdown {
			verb = "administratively shut down the BGP neighbors of"
		}
		return fmt.Errorf("failed to %s external device conn %s: %w", verb, connName, err)
	}
	return nil
}

// edcUpdateStep is one stage of an in-place update. An update runs its steps in
// the order planExternalDeviceConnUpdate returns them.
type edcUpdateStep int

const (
	// edcStepShutdown sends the connection-wide BGP admin shutdown.
	edcStepShutdown edcUpdateStep = iota
	// edcStepTunnelEdits sends the per-tunnel field edits, one request each.
	edcStepTunnelEdits
	// edcStepNoShutdown lifts the connection-wide BGP admin shutdown.
	edcStepNoShutdown
)

func (s edcUpdateStep) String() string {
	switch s {
	case edcStepShutdown:
		return "bgp_admin_shutdown=true"
	case edcStepTunnelEdits:
		return "per-tunnel edits"
	case edcStepNoShutdown:
		return "bgp_admin_shutdown=false"
	}
	return fmt.Sprintf("edcUpdateStep(%d)", int(s))
}

// planExternalDeviceConnUpdate returns the steps an in-place update has to run,
// in the order it has to run them. Nothing is sent to the controller here.
//
// The BGP admin state and the per-tunnel fields are ordered around each other
// rather than merely both applied, because the order is what makes a combined
// change safe. A change that shuts BGP down goes first, so that the tunnel is
// re-addressed or re-keyed while the session is already down; a change that
// lifts a shutdown goes last, so the session only comes back up once the new
// configuration is in place. Either way the session is down across the
// disruptive part of the apply instead of flapping on a half-applied
// configuration, which is the same ordering the gateway uses when it emits the
// neighbor deltas.
//
// The two shutdown steps are mutually exclusive: both derive from the one
// bgp_admin_shutdown attribute, so a plan holds at most one of them.
func planExternalDeviceConnUpdate(d edcDiff) []edcUpdateStep {
	var steps []edcUpdateStep

	shutdownChanged := d.HasChange(edcBgpAdminShutdownAttr)
	if shutdownChanged && getBool(d, edcBgpAdminShutdownAttr) {
		steps = append(steps, edcStepShutdown)
	}
	if d.HasChanges(edcEditableAttrs...) {
		steps = append(steps, edcStepTunnelEdits)
	}
	if shutdownChanged && !getBool(d, edcBgpAdminShutdownAttr) {
		steps = append(steps, edcStepNoShutdown)
	}

	return steps
}

// updateExternalDeviceConnInPlace applies every in-place edit the diff calls
// for: the connection-wide BGP admin state and the per-tunnel fields.
//
// Called from Update on the transit and spoke external device connection
// resources, in place of the ForceNew these attributes used to carry. The
// ordering rule lives in planExternalDeviceConnUpdate.
func updateExternalDeviceConnInPlace(d *schema.ResourceData, client *goaviatrix.Client, gwName, connName, vpcID string) error {
	steps := planExternalDeviceConnUpdate(d)
	if len(steps) == 0 {
		return nil
	}

	changed := edcChangedAttrs(d, append(slices.Clone(edcEditableAttrs), edcBgpAdminShutdownAttr)...)
	if err := edcEditableConn(d, connName, changed); err != nil {
		return err
	}

	for _, step := range steps {
		var err error
		switch step {
		case edcStepShutdown:
			err = applyEdcBgpAdminShutdown(client, gwName, connName, true)
		case edcStepTunnelEdits:
			err = editExternalDeviceConnTunnels(d, client, gwName, connName, vpcID)
		case edcStepNoShutdown:
			err = applyEdcBgpAdminShutdown(client, gwName, connName, false)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// createExternalDeviceConnBgpAdminShutdown realises a bgp_admin_shutdown that
// was asked for in a configuration being created for the first time.
//
// The create action has no such argument - the controller always writes the
// connection with its BGP neighbors up - so the only way to honour the
// configuration is a second call once the connection exists.
//
// A failure here is expensive: d.SetId has already run, so Create fails with an
// ID set, terraform taints the object and the next apply *replaces* the
// connection rather than retrying the shutdown. Recovery is `terraform untaint`
// then an apply through Update. The error is still returned - an apply must not
// report a shutdown that did not happen - and edcCustomizeDiff keeps the common
// cause, an unsupported connection type, from reaching here.
//
// Does nothing when the configuration does not ask for a shutdown, which is
// the overwhelmingly common case: there is nothing to send, since created-with-
// BGP-up is already the requested state.
func createExternalDeviceConnBgpAdminShutdown(d *schema.ResourceData, client *goaviatrix.Client, gwName, connName string) error {
	if !getBool(d, edcBgpAdminShutdownAttr) {
		return nil
	}
	if err := edcEditableConn(d, connName, []string{edcBgpAdminShutdownAttr}); err != nil {
		return err
	}
	return applyEdcBgpAdminShutdown(client, gwName, connName, true)
}

// editExternalDeviceConnTunnels applies the per-tunnel in-place edits the diff
// calls for.
//
// Called from updateExternalDeviceConnInPlace, which has already established
// that the connection is editable and sequenced any BGP admin shutdown around
// this. The edge_spoke resource has no use for either: its tunnel_protocol only
// accepts LAN, which the edit API refuses.
//
// This is not atomic, at two levels. The controller applies the fields of one
// request one at a time and leaves the earlier ones written if a later one
// fails, and this function issues one request per tunnel. A failure part way
// through therefore leaves some tunnels edited; for every attribute Read brings
// back, the next plan shows the remainder as still needing to change, which is
// the intended way to recover.
//
// pre_shared_key is the sharp exception: helper/schema persists the new state
// alongside the error and Read never reads a PSK back, so a failure after the
// first tunnel is re-keyed leaves state holding the new key for both while the
// controller keeps the old one on the tunnel it did not reach. The next plan is
// a clean no-op over a diverged connection. Recovery is setting the keys again
// or replacing the connection; the docs say so.
func editExternalDeviceConnTunnels(d *schema.ResourceData, client *goaviatrix.Client, gwName, connName, vpcID string) error {
	if !d.HasChanges(edcEditableAttrs...) {
		return nil
	}

	plan, err := planExternalDeviceConnEdits(d, gwName, connName)
	if err != nil {
		return err
	}
	if len(plan) == 0 {
		return nil
	}

	tunnels, err := client.GetExternalDeviceConnTunnels(connName, vpcID)
	if err != nil {
		return fmt.Errorf("could not read the tunnels of external device conn %q in order to edit them: %w", connName, err)
	}

	// The BGP fields, the ASN and the PSK first, then every remote_device_ip.
	// A remote_device_ip edit re-keys the tunnel, so doing those last means
	// every earlier request can use names resolved from one read.
	for _, item := range plan {
		edit := item.edit
		edit.RemoteDeviceIP = ""
		if edit.IsEmpty() {
			continue
		}
		name, err := goaviatrix.ExternalDeviceConnTunnelName(tunnels, item.matchLocalIP, item.matchRemoteIP)
		if err != nil {
			return fmt.Errorf("could not find the %s tunnel %d of connection %q to edit: %w",
				item.side, item.position+1, connName, err)
		}
		edit.TunnelName = name
		log.Printf("[INFO] Editing tunnel %s of external device conn %s", name, connName)
		if err := client.EditExternalDeviceConnTunnel(&edit); err != nil {
			return fmt.Errorf("failed to edit tunnel %s of external device conn %s: %w", name, connName, err)
		}
	}

	for _, item := range plan {
		if item.edit.RemoteDeviceIP == "" {
			continue
		}
		// Re-read: an earlier remote_device_ip edit renames the tunnel it
		// touched and, under single_ip_ha, the tunnel in the other slot too.
		tunnels, err = client.GetExternalDeviceConnTunnels(connName, vpcID)
		if err != nil {
			return fmt.Errorf("could not re-read the tunnels of external device conn %q after editing them: %w", connName, err)
		}
		// Identify by the addresses the tunnel has now, which are the new ones
		// where the first pass changed them.
		matchLocal, matchRemote := item.matchLocalIP, item.matchRemoteIP
		if item.edit.BgpLocalIP != "" {
			matchLocal = item.edit.BgpLocalIP
		}
		if item.edit.BgpRemoteIP != "" {
			matchRemote = item.edit.BgpRemoteIP
		}
		tunnel, err := goaviatrix.ExternalDeviceConnTunnel(tunnels, matchLocal, matchRemote)
		if err != nil {
			return fmt.Errorf("could not find the %s tunnel %d of connection %q to change its remote device IP: %w",
				item.side, item.position+1, connName, err)
		}
		// Skip a write of what the tunnel already has: one remote_gateway_ip
		// covers the whole side and single_ip_ha mirrors it to the partner, so
		// a later tunnel often has it already. Sending it would re-establish
		// the tunnel for no change.
		name := tunnel.Name
		if goaviatrix.SameBgpIP(tunnel.PeerIP, item.edit.RemoteDeviceIP) {
			log.Printf("[DEBUG] Skipping remote device IP update for tunnel %s of external device conn %s: "+
				"already set to %s", name, connName, item.edit.RemoteDeviceIP)
			continue
		}
		edit := goaviatrix.EditExternalDeviceConnTunnel{
			GwName:         item.edit.GwName,
			ConnectionName: item.edit.ConnectionName,
			TunnelName:     name,
			RemoteDeviceIP: item.edit.RemoteDeviceIP,
		}
		log.Printf("[INFO] Updating remote device IP of tunnel %s of external device conn %s to %s",
			name, connName, item.edit.RemoteDeviceIP)
		if err := client.EditExternalDeviceConnTunnel(&edit); err != nil {
			return fmt.Errorf("failed to edit the remote device IP of tunnel %s of external device conn %s: %w",
				name, connName, err)
		}
	}

	return nil
}
