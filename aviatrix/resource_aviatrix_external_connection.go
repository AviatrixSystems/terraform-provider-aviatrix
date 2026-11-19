package aviatrix

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"net/netip"
	"slices"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixExternalConnection() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixExternalConnectionCreate,
		ReadWithoutTimeout:   resourceAviatrixExternalConnectionRead,
		UpdateWithoutTimeout: resourceAviatrixExternalConnectionUpdate,
		DeleteWithoutTimeout: resourceAviatrixExternalConnectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: resourceAviatrixExternalConnectionCustomizeDiff,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "External connection name.",
			},
			"group_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "UUID of the gateway group this connection belongs to.",
			},
			"routing_protocol": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Routing protocol for the connection (e.g. bgp).",
			},
			"tunnel_protocol": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Tunnel protocol (e.g. LAN, IPsec, GRE).",
			},
			"tunnels": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				Description: "Set of tunnel specifications.",
				Set:         tunnelSetHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"local_gw_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Local Aviatrix gateway name.",
						},
						"remote_device_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Remote device IPv4 address.",
						},
						"local_device_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "Local device IPv4 address. If not provided, assigned by the server.",
						},
						"assigned_local_device_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Local device IPv4 address in use, as reported by the server.",
						},
						"bgp_remote_as_number": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "Per-tunnel BGP remote ASN override.",
						},
						"local_device_ipv6": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "Local device IPv6 address (with prefix).",
						},
						"remote_device_ipv6": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "Remote device IPv6 address (with prefix).",
						},
						"bgp_md5_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Sensitive:   true,
							Description: "BGP MD5 authentication key.",
						},
						"direct_connect": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether this tunnel uses direct-connect underlay.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tunnel status (UP, DOWN, UNKNOWN).",
						},
					},
				},
			},
			"bgp_local_as_number": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "BGP local ASN. Defaults to the gateway group's local ASN.",
			},
			"bgp_remote_as_number": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Connection-level BGP remote ASN. When omitted, the tunnels' common remote ASN is reported.",
			},
			"peer_vnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Peer VNet ID (applicable when peering to a CSP VNet).",
			},
			"ipv6_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Enable IPv6. If true, each tunnel must include IPv6 fields.",
			},
			"enable_bgp_multihop": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    false,
				Description: "Enable BGP multihop.",
			},
			"jumbo_frame": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    false,
				Description: "Enable jumbo frames.",
			},
			"enable_bfd": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    false,
				Description: "Enable BFD.",
			},
			"bfd_transmit_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    false,
				Description: "BFD transmit interval in milliseconds.",
			},
			"bfd_receive_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    false,
				Description: "BFD receive interval in milliseconds.",
			},
			"bfd_detect_multiplier": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    false,
				Description: "BFD detect multiplier.",
			},
			"edge_underlay": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Whether the connection uses Edge underlay.",
			},
			"bgp_lan_activemesh": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Whether the connection uses BGP LAN ActiveMesh.",
			},
			"conn_learned_cidrs_approval": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    false,
				Description: "Whether the connection learned CIDRs are approved.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server-assigned UUID of the external connection.",
			},
		},
	}
}

func marshalExternalConnectionInput(d *schema.ResourceData) *goaviatrix.ExternalConnection {
	conn := &goaviatrix.ExternalConnection{
		Name:                     getString(d, "name"),
		GroupUUID:                getString(d, "group_uuid"),
		RoutingProtocol:          getString(d, "routing_protocol"),
		TunnelProtocol:           getString(d, "tunnel_protocol"),
		BgpLocalAsNum:            getInt(d, "bgp_local_as_number"),
		BgpRemoteAsNum:           getInt(d, "bgp_remote_as_number"),
		PeerVnetID:               getString(d, "peer_vnet_id"),
		IPv6Enabled:              getBool(d, "ipv6_enabled"),
		EnableBgpMultihop:        getBool(d, "enable_bgp_multihop"),
		JumboFrame:               getBool(d, "jumbo_frame"),
		EnableBfd:                getBool(d, "enable_bfd"),
		BfdTxInterval:            getInt(d, "bfd_transmit_interval"),
		BfdRxInterval:            getInt(d, "bfd_receive_interval"),
		BfdDetectMult:            getInt(d, "bfd_detect_multiplier"),
		EdgeUnderlay:             getBool(d, "edge_underlay"),
		BgpLanActivemesh:         getBool(d, "bgp_lan_activemesh"),
		ConnLearnedCidrsApproval: getBool(d, "conn_learned_cidrs_approval"),
	}

	conn.Tunnels = marshalTunnelSpecs(d)
	return conn
}

// tunnelSetHash must cover every user-settable tunnel field: the SDK plans no
// change for a set whose element hashes are unchanged, without comparing fields.
// Server-reported fields (status, assigned_local_device_ip) are left out so that
// they never cause a diff.
func tunnelSetHash(v any) int {
	m, ok := v.(map[string]any)
	if !ok {
		return 0
	}
	str := func(k string) string {
		s, _ := m[k].(string)
		return s
	}
	asn, _ := m["bgp_remote_as_number"].(int)
	directConnect, _ := m["direct_connect"].(bool)
	return schema.HashString(fmt.Sprintf("%s-%s-%s-%d-%s-%s-%s-%t",
		str("local_gw_name"), str("remote_device_ip"), str("local_device_ip"), asn,
		str("local_device_ipv6"), str("remote_device_ipv6"), str("bgp_md5_key"), directConnect))
}

func resourceAviatrixExternalConnectionCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ any) error {
	// A flag not yet known at plan time is assumed to allow the configuration.
	ipv6Enabled := !d.NewValueKnown("ipv6_enabled") || getBool(d, "ipv6_enabled")
	activeMesh := !d.NewValueKnown("bgp_lan_activemesh") || getBool(d, "bgp_lan_activemesh")
	pc := tunnelPlanContextFrom(d)
	if !d.NewValueKnown("edge_underlay") {
		pc.perTunnelASN = true
	}

	tunnels := expandTunnelList(getSet(d, "tunnels").List())
	allKnown := true
	seen := make(map[tunnelKey]bool)
	gateways := make(map[string]bool)
	remoteASNs := make(map[int]bool)
	if connASN := getInt(d, "bgp_remote_as_number"); connASN != 0 && configSets(d, "bgp_remote_as_number") {
		remoteASNs[connASN] = true
	}
	for _, t := range tunnels {
		if !ipv6Enabled && (t.LocalDeviceIPv6 != "" || t.RemoteDeviceIPv6 != "") {
			return fmt.Errorf("tunnel (local_gw_name=%s, remote_device_ip=%s): local_device_ipv6 and remote_device_ipv6 require ipv6_enabled = true",
				t.LocalGwName, t.RemoteDeviceIP)
		}
		if t.BgpRemoteAsNum != 0 {
			remoteASNs[t.BgpRemoteAsNum] = true
		}
		// Values not yet known at plan time read as empty.
		if t.LocalGwName == "" || t.RemoteDeviceIP == "" {
			allKnown = false
			continue
		}
		key := tunnelKeyFrom(t)
		if seen[key] {
			return fmt.Errorf("duplicate tunnel (local_gw_name=%s, remote_device_ip=%s)", t.LocalGwName, t.RemoteDeviceIP)
		}
		seen[key] = true
		if !activeMesh && gateways[t.LocalGwName] {
			return fmt.Errorf("gateway %s has more than one tunnel; multiple tunnels per gateway require bgp_lan_activemesh = true", t.LocalGwName)
		}
		gateways[t.LocalGwName] = true
	}
	if !pc.perTunnelASN && len(remoteASNs) > 1 {
		return errors.New("all tunnels must use the same bgp_remote_as_number, matching the connection-level bgp_remote_as_number if set; " +
			"tunnels may differ only when edge_underlay is true and the connection-level bgp_remote_as_number is not set")
	}
	if md5KeysMixed(d) {
		return errors.New("bgp_md5_key must be set on all tunnels or on none")
	}

	if d.Id() == "" || !allKnown || !d.HasChange("tunnels") {
		return nil
	}
	oldRaw, _ := d.GetChange("tunnels")
	if _, replace := planTunnelChanges(expandTunnelList(mustSchemaSet(oldRaw).List()), tunnels, pc); replace {
		return forceNewTunnels(d)
	}
	return nil
}

// forceNewTunnels plans the replacement of the connection. ForceNew on the tunnels
// set itself only takes effect when the number of tunnels changes, so it is set on
// an attribute of a tunnel being added, which every replacing change has.
func forceNewTunnels(d *schema.ResourceDiff) error {
	oldRaw, newRaw := d.GetChange("tunnels")
	added := mustSchemaSet(newRaw).Difference(mustSchemaSet(oldRaw)).List()
	if len(added) == 0 {
		return nil
	}
	return d.ForceNew(fmt.Sprintf("tunnels.%d.local_gw_name", tunnelSetHash(added[0])))
}

// md5KeysMixed reports whether some tunnels set bgp_md5_key and others do not.
// The server checks this only within one create or add-tunnel request, so a mixed
// connection would fail when it is next created, such as on a replacement. It
// reads the raw config so that a key not yet known at plan time is not taken as
// unset.
func md5KeysMixed(d *schema.ResourceDiff) bool {
	rc := d.GetRawConfig()
	if rc.IsNull() || !rc.IsKnown() {
		return false
	}
	tunnels := rc.GetAttr("tunnels")
	if tunnels.IsNull() || !tunnels.IsKnown() {
		return false
	}
	var withKey, withoutKey bool
	for it := tunnels.ElementIterator(); it.Next(); {
		_, tunnel := it.Element()
		if !tunnel.IsKnown() {
			return false
		}
		key := tunnel.GetAttr("bgp_md5_key")
		switch {
		case !key.IsKnown():
			return false
		case key.IsNull() || key.AsString() == "":
			withoutKey = true
		default:
			withKey = true
		}
	}
	return withKey && withoutKey
}

// tunnelChangeGetter is implemented by both schema.ResourceData and schema.ResourceDiff.
type tunnelChangeGetter interface {
	Getter
	GetChange(key string) (any, any)
	GetRawConfig() cty.Value
}

func configSets(d tunnelChangeGetter, key string) bool {
	rc := d.GetRawConfig()
	if rc.IsNull() || !rc.IsKnown() {
		return false
	}
	return !rc.GetAttr(key).IsNull()
}

// tunnelPlanContext holds the connection state that decides how a tunnel change
// is applied, read from the prior state so that plan and apply agree.
type tunnelPlanContext struct {
	// assignedIPs maps each existing tunnel to the local IP the controller is
	// using, so a configured local_device_ip equal to it is not a change.
	assignedIPs map[tunnelKey]string
	// connRemoteASN is the connection-level remote ASN; a tunnel ASN equal to it
	// is not a change.
	connRemoteASN int
	// perTunnelASN is set when tunnels may use different remote ASNs: on an Edge
	// underlay connection without a connection-level ASN.
	perTunnelASN bool
}

func tunnelPlanContextFrom(d tunnelChangeGetter) tunnelPlanContext {
	oldRaw, _ := d.GetChange("tunnels")
	assignedIPs := make(map[tunnelKey]string)
	if oldSet, ok := oldRaw.(*schema.Set); ok {
		for _, t := range oldSet.List() {
			tm := mustMap(t)
			key := tunnelKey{localGwName: mustString(tm["local_gw_name"]), remoteDevice: mustString(tm["remote_device_ip"])}
			assignedIPs[key], _ = tm["assigned_local_device_ip"].(string)
		}
	}
	oldASN, _ := d.GetChange("bgp_remote_as_number")
	connRemoteASN, _ := oldASN.(int)
	return tunnelPlanContext{
		assignedIPs:   assignedIPs,
		connRemoteASN: connRemoteASN,
		perTunnelASN:  getBool(d, "edge_underlay") && !configSets(d, "bgp_remote_as_number"),
	}
}

type tunnelOpKind int

const (
	tunnelOpAdd tunnelOpKind = iota
	tunnelOpDelete
	tunnelOpSetMd5Key
)

type tunnelOp struct {
	kind   tunnelOpKind
	tunnel goaviatrix.ExternalConnectionTunnel
}

// planTunnelChanges returns the tunnel API calls that turn oldTunnels into
// newTunnels, or replace = true when the change can only be made by replacing the
// connection. The server patches only bgp_md5_key; other tunnel changes are a
// delete followed by an add, because without ActiveMesh a gateway holds one tunnel
// of a connection. Tunnels on new gateways are added first and removed tunnels
// deleted last, so that the connection loses at most one tunnel at a time; the
// server refuses to delete its last tunnel. The remote ASN and direct_connect can
// only change by replacing the connection, except that the ASN of one tunnel may
// change on an Edge underlay connection without a connection-level ASN.
func planTunnelChanges(oldTunnels, newTunnels []goaviatrix.ExternalConnectionTunnel, pc tunnelPlanContext) (ops []tunnelOp, replace bool) {
	sortTunnels(oldTunnels)
	sortTunnels(newTunnels)

	if !pc.perTunnelASN && pc.connRemoteASN != 0 {
		for _, t := range newTunnels {
			if t.BgpRemoteAsNum != 0 && t.BgpRemoteAsNum != pc.connRemoteASN {
				return nil, true
			}
		}
	}

	oldByKey := make(map[tunnelKey]goaviatrix.ExternalConnectionTunnel, len(oldTunnels))
	for _, t := range oldTunnels {
		oldByKey[tunnelKeyFrom(t)] = t
	}
	newKeys := make(map[tunnelKey]bool, len(newTunnels))

	type replacement struct {
		old, new goaviatrix.ExternalConnectionTunnel
	}
	var replacements []replacement
	var added, md5Changes []goaviatrix.ExternalConnectionTunnel
	for _, newT := range newTunnels {
		key := tunnelKeyFrom(newT)
		newKeys[key] = true
		oldT, exists := oldByKey[key]
		switch {
		case !exists:
			added = append(added, newT)
		case oldT.DirectConnect != newT.DirectConnect:
			return nil, true
		case tunnelNeedsRecreate(oldT, newT, pc.assignedIPs[key], pc.connRemoteASN):
			replacements = append(replacements, replacement{old: oldT, new: newT})
		case oldT.BgpMd5Key != newT.BgpMd5Key:
			md5Changes = append(md5Changes, newT)
		}
	}
	var removed []goaviatrix.ExternalConnectionTunnel
	for _, oldT := range oldTunnels {
		if !newKeys[tunnelKeyFrom(oldT)] {
			removed = append(removed, oldT)
		}
	}

	var onNewGateways []goaviatrix.ExternalConnectionTunnel
	for _, newT := range added {
		i := slices.IndexFunc(removed, func(t goaviatrix.ExternalConnectionTunnel) bool {
			return t.LocalGwName == newT.LocalGwName
		})
		if i < 0 {
			onNewGateways = append(onNewGateways, newT)
			continue
		}
		if removed[i].DirectConnect != newT.DirectConnect {
			return nil, true
		}
		replacements = append(replacements, replacement{old: removed[i], new: newT})
		removed = slices.Delete(removed, i, i+1)
	}

	for _, t := range onNewGateways {
		ops = append(ops, tunnelOp{kind: tunnelOpAdd, tunnel: t})
	}
	for _, r := range replacements {
		ops = append(ops, tunnelOp{kind: tunnelOpDelete, tunnel: r.old}, tunnelOp{kind: tunnelOpAdd, tunnel: r.new})
	}
	for _, t := range removed {
		ops = append(ops, tunnelOp{kind: tunnelOpDelete, tunnel: t})
	}
	for _, t := range md5Changes {
		ops = append(ops, tunnelOp{kind: tunnelOpSetMd5Key, tunnel: t})
	}

	count := len(oldTunnels)
	for _, op := range ops {
		switch op.kind {
		case tunnelOpAdd:
			count++
		case tunnelOpDelete:
			// The controller refuses to delete a connection's last tunnel, so a
			// sequence that needs to only works by replacing the connection.
			if count < 2 {
				return nil, true
			}
			count--
		}
	}
	return ops, false
}

func sortTunnels(tunnels []goaviatrix.ExternalConnectionTunnel) {
	slices.SortFunc(tunnels, func(a, b goaviatrix.ExternalConnectionTunnel) int {
		return cmp.Or(cmp.Compare(a.LocalGwName, b.LocalGwName), cmp.Compare(a.RemoteDeviceIP, b.RemoteDeviceIP))
	})
}

// localIPAddr returns the IP of a local_device_ip, which on GCP is given as
// "vpc:subnet:ip" while the server reports only the IP.
func localIPAddr(v string) string {
	return v[strings.LastIndex(v, ":")+1:]
}

// ipv6Addr returns the canonical address of an IPv6 value. Tunnel IPv6 fields are
// configured with a prefix length, but the server reports the bare address.
func ipv6Addr(v string) string {
	addr, _, _ := strings.Cut(v, "/")
	if a, err := netip.ParseAddr(addr); err == nil {
		return a.String()
	}
	return addr
}

// tunnelStateFromServer builds the state view of a tunnel read from the server.
// Where the configured value means the same as the reported one, it is kept so
// that the set hash still matches the configuration: an empty local_device_ip or
// zero bgp_remote_as_number lets the server choose, and the server omits the
// tunnel ASN when a connection-level ASN is set. bgp_md5_key is never returned;
// the server reports a set key as "***", which is kept when no key is configured
// so that a key set outside Terraform shows as a change.
func tunnelStateFromServer(prior goaviatrix.ExternalConnectionTunnel, hasPrior bool, server goaviatrix.ExternalConnectionTunnel) goaviatrix.ExternalConnectionTunnel {
	t := server
	if !hasPrior {
		return t
	}
	if prior.LocalDeviceIP == "" || localIPAddr(prior.LocalDeviceIP) == server.LocalDeviceIP {
		t.LocalDeviceIP = prior.LocalDeviceIP
	}
	if prior.BgpRemoteAsNum == 0 || server.BgpRemoteAsNum == 0 {
		t.BgpRemoteAsNum = prior.BgpRemoteAsNum
	}
	if ipv6Addr(prior.LocalDeviceIPv6) == ipv6Addr(server.LocalDeviceIPv6) {
		t.LocalDeviceIPv6 = prior.LocalDeviceIPv6
	}
	if ipv6Addr(prior.RemoteDeviceIPv6) == ipv6Addr(server.RemoteDeviceIPv6) {
		t.RemoteDeviceIPv6 = prior.RemoteDeviceIPv6
	}
	if server.BgpMd5Key != "" && prior.BgpMd5Key != "" {
		t.BgpMd5Key = prior.BgpMd5Key
	}
	return t
}

func marshalTunnelSpecs(d *schema.ResourceData) []goaviatrix.ExternalConnectionTunnel {
	tunnelList := getSet(d, "tunnels").List()
	tunnels := make([]goaviatrix.ExternalConnectionTunnel, 0, len(tunnelList))

	for _, t := range tunnelList {
		tm := mustMap(t)
		tunnel := goaviatrix.ExternalConnectionTunnel{
			LocalGwName:    mustString(tm["local_gw_name"]),
			RemoteDeviceIP: mustString(tm["remote_device_ip"]),
			DirectConnect:  mustBool(tm["direct_connect"]),
		}
		if v, ok := tm["local_device_ip"]; ok && v != "" {
			tunnel.LocalDeviceIP = mustString(v)
		}
		if v, ok := tm["bgp_remote_as_number"]; ok && v != 0 {
			tunnel.BgpRemoteAsNum = mustInt(v)
		}
		if v, ok := tm["local_device_ipv6"]; ok && v != "" {
			tunnel.LocalDeviceIPv6 = mustString(v)
		}
		if v, ok := tm["remote_device_ipv6"]; ok && v != "" {
			tunnel.RemoteDeviceIPv6 = mustString(v)
		}
		if v, ok := tm["bgp_md5_key"]; ok && v != "" {
			tunnel.BgpMd5Key = mustString(v)
		}
		tunnels = append(tunnels, tunnel)
	}

	return tunnels
}

type tunnelKey struct {
	localGwName  string
	remoteDevice string
}

func tunnelKeyFrom(t goaviatrix.ExternalConnectionTunnel) tunnelKey {
	return tunnelKey{localGwName: t.LocalGwName, remoteDevice: t.RemoteDeviceIP}
}

func resourceAviatrixExternalConnectionCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	conn := marshalExternalConnectionInput(d)

	uuid, err := client.CreateExternalConnection(ctx, conn)
	if err != nil {
		return diag.Errorf("failed to create external connection: %s", err)
	}

	d.SetId(uuid)
	return resourceAviatrixExternalConnectionRead(ctx, d, meta)
}

func resourceAviatrixExternalConnectionRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	uuid := d.Id()

	conn, err := client.GetExternalConnection(ctx, uuid)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read external connection: %s", err)
	}

	mustSet(d, "uuid", conn.UUID)
	mustSet(d, "name", conn.Name)
	mustSet(d, "group_uuid", conn.GroupUUID)
	mustSet(d, "routing_protocol", conn.RoutingProtocol)
	mustSet(d, "tunnel_protocol", conn.TunnelProtocol)
	mustSet(d, "bgp_local_as_number", conn.BgpLocalAsNum)
	mustSet(d, "bgp_remote_as_number", conn.BgpRemoteAsNum)
	mustSet(d, "peer_vnet_id", conn.PeerVnetID)
	mustSet(d, "ipv6_enabled", conn.IPv6Enabled)
	mustSet(d, "enable_bgp_multihop", conn.EnableBgpMultihop)
	mustSet(d, "jumbo_frame", conn.JumboFrame)
	mustSet(d, "enable_bfd", conn.EnableBfd)
	mustSet(d, "bfd_transmit_interval", conn.BfdTxInterval)
	mustSet(d, "bfd_receive_interval", conn.BfdRxInterval)
	mustSet(d, "bfd_detect_multiplier", conn.BfdDetectMult)
	mustSet(d, "edge_underlay", conn.EdgeUnderlay)
	mustSet(d, "bgp_lan_activemesh", conn.BgpLanActivemesh)
	mustSet(d, "conn_learned_cidrs_approval", conn.ConnLearnedCidrsApproval)

	prior := make(map[tunnelKey]goaviatrix.ExternalConnectionTunnel)
	if oldTunnels, ok := d.GetOk("tunnels"); ok {
		for _, t := range expandTunnelList(mustSchemaSet(oldTunnels).List()) {
			prior[tunnelKeyFrom(t)] = t
		}
	}

	tunnels := make([]map[string]any, 0, len(conn.Tunnels))
	for _, server := range conn.Tunnels {
		p, hasPrior := prior[tunnelKeyFrom(server)]
		t := tunnelStateFromServer(p, hasPrior, server)
		tunnels = append(tunnels, map[string]any{
			"local_gw_name":            t.LocalGwName,
			"remote_device_ip":         t.RemoteDeviceIP,
			"local_device_ip":          t.LocalDeviceIP,
			"assigned_local_device_ip": server.LocalDeviceIP,
			"bgp_remote_as_number":     t.BgpRemoteAsNum,
			"local_device_ipv6":        t.LocalDeviceIPv6,
			"remote_device_ipv6":       t.RemoteDeviceIPv6,
			"bgp_md5_key":              t.BgpMd5Key,
			"direct_connect":           t.DirectConnect,
			"status":                   t.Status,
		})
	}
	mustSet(d, "tunnels", tunnels)

	return nil
}

func resourceAviatrixExternalConnectionUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	uuid := d.Id()

	if d.HasChanges("enable_bgp_multihop", "jumbo_frame", "enable_bfd",
		"bfd_transmit_interval", "bfd_receive_interval", "bfd_detect_multiplier",
		"conn_learned_cidrs_approval") {

		update := &goaviatrix.ExternalConnectionUpdate{}
		if d.HasChange("enable_bgp_multihop") {
			v := getBool(d, "enable_bgp_multihop")
			update.EnableBgpMultihop = &v
		}
		if d.HasChange("jumbo_frame") {
			v := getBool(d, "jumbo_frame")
			update.JumboFrame = &v
		}
		if d.HasChange("enable_bfd") {
			v := getBool(d, "enable_bfd")
			update.EnableBfd = &v
		}
		if d.HasChange("bfd_transmit_interval") {
			v := getInt(d, "bfd_transmit_interval")
			update.BfdTxInterval = &v
		}
		if d.HasChange("bfd_receive_interval") {
			v := getInt(d, "bfd_receive_interval")
			update.BfdRxInterval = &v
		}
		if d.HasChange("bfd_detect_multiplier") {
			v := getInt(d, "bfd_detect_multiplier")
			update.BfdDetectMult = &v
		}
		if d.HasChange("conn_learned_cidrs_approval") {
			v := getBool(d, "conn_learned_cidrs_approval")
			update.ConnLearnedCidrsApproval = &v
		}

		_, err := client.UpdateExternalConnection(ctx, uuid, update)
		if err != nil {
			d.Partial(true)
			return diag.Errorf("failed to update external connection: %s", err)
		}
	}

	if d.HasChange("tunnels") {
		oldRaw, newRaw := d.GetChange("tunnels")
		ops, replace := planTunnelChanges(
			expandTunnelList(mustSchemaSet(oldRaw).List()),
			expandTunnelList(mustSchemaSet(newRaw).List()),
			tunnelPlanContextFrom(d))
		if replace {
			d.Partial(true)
			return diag.Errorf("the tunnel changes can only be applied by replacing the external connection; " +
				"run terraform plan again to plan the replacement")
		}
		for _, op := range ops {
			if err := applyTunnelOp(ctx, client, uuid, op); err != nil {
				// Keep the prior state; the next refresh reports what the server has.
				d.Partial(true)
				return diag.FromErr(err)
			}
		}
	}

	return resourceAviatrixExternalConnectionRead(ctx, d, meta)
}

func applyTunnelOp(ctx context.Context, client *goaviatrix.Client, uuid string, op tunnelOp) error {
	t := op.tunnel
	var err error
	var action string
	switch op.kind {
	case tunnelOpAdd:
		action = "create"
		err = client.CreateExternalConnectionTunnel(ctx, uuid, &t)
	case tunnelOpDelete:
		action = "delete"
		err = client.DeleteExternalConnectionTunnel(ctx, uuid, t.LocalGwName, t.RemoteDeviceIP)
	case tunnelOpSetMd5Key:
		action = "update bgp_md5_key of"
		err = client.UpdateExternalConnectionTunnel(ctx, uuid, &goaviatrix.ExternalConnectionTunnelUpdate{
			LocalGwName:    t.LocalGwName,
			RemoteDeviceIP: t.RemoteDeviceIP,
			BgpMd5Key:      &t.BgpMd5Key,
		})
	}
	if err != nil {
		return fmt.Errorf("failed to %s tunnel (local_gw_name=%s, remote_device_ip=%s): %w", action, t.LocalGwName, t.RemoteDeviceIP, err)
	}
	return nil
}

func resourceAviatrixExternalConnectionDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := mustClient(meta)

	uuid := d.Id()
	err := client.DeleteExternalConnection(ctx, uuid)
	if err != nil {
		return diag.Errorf("failed to delete external connection: %v", err)
	}

	return nil
}

func expandTunnelList(raw []any) []goaviatrix.ExternalConnectionTunnel {
	tunnels := make([]goaviatrix.ExternalConnectionTunnel, 0, len(raw))
	for _, t := range raw {
		tm := mustMap(t)
		tunnel := goaviatrix.ExternalConnectionTunnel{
			LocalGwName:      mustString(tm["local_gw_name"]),
			RemoteDeviceIP:   mustString(tm["remote_device_ip"]),
			LocalDeviceIP:    mustString(tm["local_device_ip"]),
			BgpRemoteAsNum:   mustInt(tm["bgp_remote_as_number"]),
			LocalDeviceIPv6:  mustString(tm["local_device_ipv6"]),
			RemoteDeviceIPv6: mustString(tm["remote_device_ipv6"]),
			BgpMd5Key:        mustString(tm["bgp_md5_key"]),
			DirectConnect:    mustBool(tm["direct_connect"]),
		}
		tunnels = append(tunnels, tunnel)
	}
	return tunnels
}

// tunnelNeedsRecreate returns true when newT differs from the tunnel running on the
// server in a field that can only be changed by delete+recreate. Values are compared
// by meaning, not spelling, so a configuration that already matches the server is
// applied without touching the tunnel. bgp_md5_key is excluded because it can be
// updated in-place via PATCH, and direct_connect because it requires replacing the
// connection.
func tunnelNeedsRecreate(oldT, newT goaviatrix.ExternalConnectionTunnel, oldAssignedIP string, connRemoteASN int) bool {
	return localIPChanged(oldT.LocalDeviceIP, oldAssignedIP, newT.LocalDeviceIP) ||
		remoteASNChanged(oldT.BgpRemoteAsNum, newT.BgpRemoteAsNum, connRemoteASN) ||
		ipv6Changed(oldT.LocalDeviceIPv6, newT.LocalDeviceIPv6) ||
		ipv6Changed(oldT.RemoteDeviceIPv6, newT.RemoteDeviceIPv6)
}

// ipv6Changed compares IPv6 values by address and prefix length. The server
// reports the bare address, so an old value read without a prefix, such as after
// import, is compared by address only.
func ipv6Changed(oldV, newV string) bool {
	if ipv6Addr(oldV) != ipv6Addr(newV) {
		return true
	}
	_, oldBits, oldHasPrefix := strings.Cut(oldV, "/")
	_, newBits, _ := strings.Cut(newV, "/")
	return oldHasPrefix && oldBits != newBits
}

// localIPChanged treats an empty new local_device_ip as "server chooses", which
// keeps whatever IP the tunnel already has.
func localIPChanged(oldIP, oldAssignedIP, newIP string) bool {
	if newIP == "" || newIP == oldIP {
		return false
	}
	if strings.Contains(oldIP, ":") && strings.Contains(newIP, ":") {
		return true
	}
	running := oldAssignedIP
	if running == "" {
		running = localIPAddr(oldIP)
	}
	return localIPAddr(newIP) != running
}

// remoteASNChanged treats a zero new ASN as "server chooses". A zero old ASN means
// the tunnel runs with the connection-level ASN.
func remoteASNChanged(oldASN, newASN, connRemoteASN int) bool {
	if newASN == 0 || newASN == oldASN {
		return false
	}
	running := oldASN
	if running == 0 {
		running = connRemoteASN
	}
	return newASN != running
}
