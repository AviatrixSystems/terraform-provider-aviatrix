package aviatrix

import (
	"context"
	"errors"
	"fmt"

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
							Computed:    true,
							Description: "Local device IPv4 address. If not provided, assigned by the server.",
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
				ForceNew:    true,
				Description: "BGP local ASN.",
			},
			"bgp_remote_as_number": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Connection-level BGP remote ASN.",
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

func tunnelSetHash(v any) int {
	m, ok := v.(map[string]any)
	if !ok {
		return 0
	}
	return schema.HashString(fmt.Sprintf("%s-%s", m["local_gw_name"], m["remote_device_ip"]))
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

	// Build a lookup of bgp_md5_key from local state since the server never returns it.
	stateMd5 := make(map[tunnelKey]string)
	if oldTunnels, ok := d.GetOk("tunnels"); ok {
		for _, t := range mustSchemaSet(oldTunnels).List() {
			tm := mustMap(t)
			key := tunnelKey{localGwName: mustString(tm["local_gw_name"]), remoteDevice: mustString(tm["remote_device_ip"])}
			stateMd5[key] = mustString(tm["bgp_md5_key"])
		}
	}

	tunnels := make([]map[string]any, 0, len(conn.Tunnels))
	for _, t := range conn.Tunnels {
		key := tunnelKeyFrom(t)
		tunnels = append(tunnels, map[string]any{
			"local_gw_name":        t.LocalGwName,
			"remote_device_ip":     t.RemoteDeviceIP,
			"local_device_ip":      t.LocalDeviceIP,
			"bgp_remote_as_number": t.BgpRemoteAsNum,
			"local_device_ipv6":    t.LocalDeviceIPv6,
			"remote_device_ipv6":   t.RemoteDeviceIPv6,
			"bgp_md5_key":          stateMd5[key],
			"direct_connect":       t.DirectConnect,
			"status":               t.Status,
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
			return diag.Errorf("failed to update external connection: %s", err)
		}
	}

	if d.HasChange("tunnels") {
		oldRaw, newRaw := d.GetChange("tunnels")
		oldTunnels := expandTunnelList(mustSchemaSet(oldRaw).List())
		newTunnels := expandTunnelList(mustSchemaSet(newRaw).List())

		oldMap := make(map[tunnelKey]goaviatrix.ExternalConnectionTunnel)
		for _, t := range oldTunnels {
			oldMap[tunnelKeyFrom(t)] = t
		}

		newMap := make(map[tunnelKey]goaviatrix.ExternalConnectionTunnel)
		for _, t := range newTunnels {
			newMap[tunnelKeyFrom(t)] = t
		}

		// Delete tunnels that were removed or structurally changed
		for key, oldT := range oldMap {
			if newT, exists := newMap[key]; !exists || tunnelStructurallyChanged(oldT, newT) {
				err := client.DeleteExternalConnectionTunnel(ctx, uuid, oldT.LocalGwName, oldT.RemoteDeviceIP)
				if err != nil {
					return diag.Errorf("failed to delete tunnel (local_gw_name=%s, remote_device_ip=%s): %s",
						oldT.LocalGwName, oldT.RemoteDeviceIP, err)
				}
			}
		}

		// Create tunnels that were added or structurally changed
		for key, newT := range newMap {
			if oldT, exists := oldMap[key]; !exists || tunnelStructurallyChanged(oldT, newT) {
				err := client.CreateExternalConnectionTunnel(ctx, uuid, &newT)
				if err != nil {
					return diag.Errorf("failed to create tunnel (local_gw_name=%s, remote_device_ip=%s): %s",
						newT.LocalGwName, newT.RemoteDeviceIP, err)
				}
			}
		}

		// PATCH tunnels where only bgp_md5_key changed (avoids unnecessary delete+recreate)
		for key, newT := range newMap {
			if oldT, exists := oldMap[key]; exists && !tunnelStructurallyChanged(oldT, newT) && oldT.BgpMd5Key != newT.BgpMd5Key {
				err := client.UpdateExternalConnectionTunnel(ctx, uuid, &goaviatrix.ExternalConnectionTunnelUpdate{
					LocalGwName:    newT.LocalGwName,
					RemoteDeviceIP: newT.RemoteDeviceIP,
					BgpMd5Key:      &newT.BgpMd5Key,
				})
				if err != nil {
					return diag.Errorf("failed to update tunnel bgp_md5_key (local_gw_name=%s, remote_device_ip=%s): %s",
						newT.LocalGwName, newT.RemoteDeviceIP, err)
				}
			}
		}
	}

	return resourceAviatrixExternalConnectionRead(ctx, d, meta)
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

// tunnelStructurallyChanged returns true when fields that require delete+recreate have changed.
// bgp_md5_key is excluded because it can be updated in-place via PATCH.
func tunnelStructurallyChanged(a, b goaviatrix.ExternalConnectionTunnel) bool {
	return a.LocalDeviceIP != b.LocalDeviceIP ||
		a.BgpRemoteAsNum != b.BgpRemoteAsNum ||
		a.LocalDeviceIPv6 != b.LocalDeviceIPv6 ||
		a.RemoteDeviceIPv6 != b.RemoteDeviceIPv6 ||
		a.DirectConnect != b.DirectConnect
}
