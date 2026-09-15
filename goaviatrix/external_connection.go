package goaviatrix

import (
	"context"
	"fmt"
	"net/url"
)

type ExternalConnectionTunnel struct {
	LocalGwName      string `json:"local_gw_name"`
	RemoteDeviceIP   string `json:"remote_device_ip"`
	LocalDeviceIP    string `json:"local_device_ip,omitempty"`
	BgpRemoteAsNum   int    `json:"bgp_remote_as_number,omitempty"`
	LocalDeviceIPv6  string `json:"local_device_ipv6,omitempty"`
	RemoteDeviceIPv6 string `json:"remote_device_ipv6,omitempty"`
	BgpMd5Key        string `json:"bgp_md5_key,omitempty"`
	DirectConnect    bool   `json:"direct_connect"`
	Status           string `json:"status,omitempty"`
}

type ExternalConnection struct {
	UUID                     string                     `json:"uuid,omitempty"`
	Name                     string                     `json:"name"`
	GroupUUID                string                     `json:"group_uuid"`
	RoutingProtocol          string                     `json:"routing_protocol"`
	TunnelProtocol           string                     `json:"tunnel_protocol"`
	Tunnels                  []ExternalConnectionTunnel `json:"tunnels"`
	BgpLocalAsNum            int                        `json:"bgp_local_as_number,omitempty"`
	BgpRemoteAsNum           int                        `json:"bgp_remote_as_number,omitempty"`
	PeerVnetID               string                     `json:"peer_vnet_id,omitempty"`
	IPv6Enabled              bool                       `json:"ipv6_enabled"`
	EnableBgpMultihop        bool                       `json:"enable_bgp_multihop"`
	JumboFrame               bool                       `json:"jumbo_frame"`
	EnableBfd                bool                       `json:"enable_bfd"`
	BfdTxInterval            int                        `json:"bfd_transmit_interval,omitempty"`
	BfdRxInterval            int                        `json:"bfd_receive_interval,omitempty"`
	BfdDetectMult            int                        `json:"bfd_detect_multiplier,omitempty"`
	EdgeUnderlay             bool                       `json:"edge_underlay"`
	BgpLanActivemesh         bool                       `json:"bgp_lan_activemesh"`
	ConnLearnedCidrsApproval bool                       `json:"conn_learned_cidrs_approval"`
}

func (c *Client) CreateExternalConnection(ctx context.Context, conn *ExternalConnection) (string, error) {
	type resp struct {
		UUID string `json:"uuid"`
	}

	var data resp
	err := c.PostAPIContext25(ctx, &data, "external-connections", conn)
	if err != nil {
		return "", err
	}

	return data.UUID, nil
}

func (c *Client) GetExternalConnection(ctx context.Context, uuid string) (*ExternalConnection, error) {
	endpoint := fmt.Sprintf("external-connections/%s", uuid)

	var data ExternalConnection
	err := c.GetAPIContext25(ctx, &data, endpoint, nil)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

type ExternalConnectionUpdate struct {
	EnableBgpMultihop        *bool `json:"enable_bgp_multihop,omitempty"`
	JumboFrame               *bool `json:"jumbo_frame,omitempty"`
	EnableBfd                *bool `json:"enable_bfd,omitempty"`
	BfdTxInterval            *int  `json:"bfd_transmit_interval,omitempty"`
	BfdRxInterval            *int  `json:"bfd_receive_interval,omitempty"`
	BfdDetectMult            *int  `json:"bfd_detect_multiplier,omitempty"`
	ConnLearnedCidrsApproval *bool `json:"conn_learned_cidrs_approval,omitempty"`
}

type ExternalConnectionTunnelUpdate struct {
	LocalGwName    string  `json:"local_gw_name"`
	RemoteDeviceIP string  `json:"remote_device_ip"`
	BgpMd5Key      *string `json:"bgp_md5_key"` // bgp_md5_key accept empty string as clear the key
}

func (c *Client) UpdateExternalConnection(ctx context.Context, uuid string, update *ExternalConnectionUpdate) (*ExternalConnection, error) {
	endpoint := fmt.Sprintf("external-connections/%s", uuid)

	var data ExternalConnection
	err := c.PatchAPIContext25(ctx, &data, endpoint, update)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *Client) DeleteExternalConnection(ctx context.Context, uuid string) error {
	endpoint := fmt.Sprintf("external-connections/%s", uuid)
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}

func (c *Client) CreateExternalConnectionTunnel(ctx context.Context, connUUID string, tunnel *ExternalConnectionTunnel) error {
	endpoint := fmt.Sprintf("external-connections/%s/tunnels", connUUID)
	return c.PostAPIContext25(ctx, nil, endpoint, tunnel)
}

func (c *Client) UpdateExternalConnectionTunnel(ctx context.Context, connUUID string, update *ExternalConnectionTunnelUpdate) error {
	endpoint := fmt.Sprintf("external-connections/%s/tunnels", connUUID)
	return c.PatchAPIContext25(ctx, nil, endpoint, update)
}

func (c *Client) DeleteExternalConnectionTunnel(ctx context.Context, connUUID, localGwName, remoteDeviceIP string) error {
	endpoint := fmt.Sprintf("external-connections/%s/tunnels?local_gw_name=%s&remote_device_ip=%s",
		connUUID,
		url.QueryEscape(localGwName),
		url.QueryEscape(remoteDeviceIP),
	)
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}
