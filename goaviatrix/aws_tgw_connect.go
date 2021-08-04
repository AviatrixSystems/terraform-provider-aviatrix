package goaviatrix

import (
	"context"
	"fmt"
	"strings"
)

type AwsTgwConnect struct {
	Action                  string `form:"action"`
	CID                     string `form:"CID"`
	TgwName                 string `form:"tgw_name" json:"tgw_name"`
	ConnectionName          string `form:"connection_name" json:"connection_name"`
	TransportAttachmentID   string `form:"transport_vpc_id" json:"transport_attachment_id"`
	TransportAttachmentName string `json:"transport_attachment_name"`
	TransportVpcName        string `json:"transport_vpc_name"`
	SecurityDomainName      string `form:"security_domain_name" json:"route_domain_name"`
	ConnectAttachmentID     string `form:"connect_attachment_id" json:"connect_attachment_id"`
}

func (a *AwsTgwConnect) ID() string {
	return fmt.Sprintf("%s~~%s", a.TgwName, a.ConnectionName)
}

type AwsTgwConnectPeer struct {
	Action              string   `form:"action"`
	CID                 string   `form:"CID"`
	TgwName             string   `form:"tgw_name" json:"tgw_name"`
	ConnectAttachmentID string   `form:"connect_attachment_id" json:"connect_attachment_id"`
	ConnectPeerName     string   `form:"connect_peer_name" json:"connect_peer_name"`
	ConnectPeerID       string   `form:"connect_peer_id" json:"connect_peer_id"`
	PeerGreAddress      string   `form:"peer_gre_address" json:"peer_gre_address"`
	InsideIPCidrs       []string `json:"inside_ip_cidr"`
	InsideIPCidrsString string   `form:"bgp_inside_cidrs"`
	PeerASNumber        string   `form:"peer_as_number" json:"peer_asn"`
	TgwGreAddress       string   `form:"tgw_gre_address" json:"tgw_gre_address"`
	ConnectionName      string   `form:"connection_name" json:"connection_name"`
}

func (a *AwsTgwConnectPeer) ID() string {
	return fmt.Sprintf("%s~~%s~~%s", a.TgwName, a.ConnectionName, a.ConnectPeerName)
}

func (c *Client) AttachTGWConnectToTGW(ctx context.Context, connect *AwsTgwConnect) error {
	connect.Action = "attach_tgw_connect_to_tgw"
	connect.CID = c.CID
	return c.PostAPIContext(ctx, connect.Action, connect, BasicCheck)
}

func (c *Client) DetachTGWConnectFromTGW(ctx context.Context, connect *AwsTgwConnect) error {
	connect.Action = "detach_tgw_connect_from_tgw"
	connect.CID = c.CID
	return c.PostAPIContext(ctx, connect.Action, connect, BasicCheck)
}

func (c *Client) GetTGWConnect(ctx context.Context, connect *AwsTgwConnect) (*AwsTgwConnect, error) {
	form := map[string]string{
		"action":          "get_tgw_connect_by_connection_name",
		"CID":             c.CID,
		"connection_name": connect.ConnectionName,
		"tgw_name":        connect.TgwName,
	}

	check := func(action, method, reason string, ret bool) error {
		if !ret {
			// 'AVXERR-TGW-0072': 'TGW Connect {conn_name} does not exist.'
			if strings.Contains(reason, "AVXERR-TGW-0072") {
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s %s failed: %s", action, method, reason)
		}
		return nil
	}

	var data struct {
		Results AwsTgwConnect
	}
	err := c.GetAPIContext(ctx, &data, form["action"], form, check)
	if err != nil {
		return nil, err
	}
	return &data.Results, nil
}

func (c *Client) CreateTGWConnectPeer(ctx context.Context, peer *AwsTgwConnectPeer) error {
	peer.Action = "create_tgw_connect_peer"
	peer.CID = c.CID
	peer.InsideIPCidrsString = strings.Join(peer.InsideIPCidrs, ",")
	return c.PostAPIContext(ctx, peer.Action, peer, BasicCheck)
}

func (c *Client) DeleteTGWConnectPeer(ctx context.Context, peer *AwsTgwConnectPeer) error {
	peer.Action = "delete_tgw_connect_peer"
	peer.CID = c.CID
	return c.PostAPIContext(ctx, peer.Action, peer, BasicCheck)
}

func (c *Client) GetTGWConnectPeer(ctx context.Context, peer *AwsTgwConnectPeer) (*AwsTgwConnectPeer, error) {
	form := map[string]string{
		"action":            "get_tgw_connect_peer_by_connect_peer_name",
		"CID":               c.CID,
		"connection_name":   peer.ConnectionName,
		"tgw_name":          peer.TgwName,
		"connect_peer_name": peer.ConnectPeerName,
	}

	check := func(action, method, reason string, ret bool) error {
		if !ret {
			// 'AVXERR-TGW-0074': 'TGW Connect peer {conn_name}/{connect_peer_name} does not exist.',
			if strings.Contains(reason, "AVXERR-TGW-0074") {
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s %s failed: %s", action, method, reason)
		}
		return nil
	}

	var data struct {
		Results AwsTgwConnectPeer
	}
	err := c.GetAPIContext(ctx, &data, form["action"], form, check)
	if err != nil {
		return nil, err
	}
	return &data.Results, nil
}
