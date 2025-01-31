package goaviatrix

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

type TransitGatewayPeering struct {
	TransitGatewayName1                 string `form:"gateway1,omitempty" json:"gateway1,omitempty"`
	TransitGatewayName2                 string `form:"gateway2,omitempty" json:"gateway2,omitempty"`
	Gateway1ExcludedCIDRs               string `form:"src_filter_list,omitempty" json:"src_filter_list,omitempty"`
	Gateway2ExcludedCIDRs               string `form:"dst_filter_list,omitempty" json:"dst_filter_list,omitempty"`
	Gateway1ExcludedTGWConnections      string `form:"source_exclude_connections,omitempty" json:"source_exclude_connections,omitempty"`
	Gateway2ExcludedTGWConnections      string `form:"destination_exclude_connections,omitempty" json:"destination_exclude_connections,omitempty"`
	PrivateIPPeering                    string `form:"private_ip_peering,omitempty" json:"private_ip_peering,omitempty"`
	InsaneModeOverInternet              bool   `form:"insane_mode_over_internet,omitempty" json:"insane_mode_over_internet,omitempty"`
	InsaneModeTunnelCount               int    `json:"insane_mode_tunnel_count"`
	TunnelCount                         int    `form:"tunnel_count,omitempty" json:"tunnel_count,omitempty"`
	Gateway1ExcludedCIDRsSlice          []string
	Gateway2ExcludedCIDRsSlice          []string
	Gateway1ExcludedTGWConnectionsSlice []string
	Gateway2ExcludedTGWConnectionsSlice []string
	PrependAsPath1                      string
	PrependAsPath2                      string
	CID                                 string   `form:"CID,omitempty" json:"CID,omitempty"`
	Action                              string   `form:"action,omitempty" json:"action,omitempty"`
	SingleTunnel                        string   `form:"single_tunnel,omitempty" json:"single_tunnel,omitempty"`
	NoMaxPerformance                    bool     `form:"no_max_performance,omitempty" json:"no_max_performance,omitempty"`
	EnableOverPrivateNetwork            bool     `form:"over_private_network,omitempty" json:"over_private_network,omitempty"`
	EnableJumboFrame                    bool     `form:"jumbo_frame,omitempty" json:"jumbo_frame,omitempty"`
	EnableInsaneMode                    bool     `form:"insane_mode,omitempty" json:"insane_mode,omitempty"`
	SrcWanInterfaces                    string   `form:"src_wan_interfaces,omitempty" json:"src_wan_interfaces,omitempty"`
	DstWanInterfaces                    string   `form:"dst_wan_interfaces,omitempty" json:"dst_wan_interfaces,omitempty"`
	Gateway1LogicalIfNames              []string `form:"gateway1_logical_ifnames,omitempty" json:"gateway1_logical_ifnames,omitempty"`
	Gateway2LogicalIfNames              []string `form:"gateway2_logical_ifnames,omitempty" json:"gateway2_logical_ifnames,omitempty"`
}

type TransitGatewayPeeringEdit struct {
	TransitGatewayName1            string `form:"gateway1,omitempty" json:"gateway_1,omitempty"`
	TransitGatewayName2            string `form:"gateway2,omitempty" json:"gateway_2,omitempty"`
	Gateway1ExcludedCIDRs          string `form:"src_filter_list,omitempty"`
	Gateway2ExcludedCIDRs          string `form:"dst_filter_list,omitempty"`
	Gateway1ExcludedTGWConnections string `form:"source_exclude_connections,omitempty"`
	Gateway2ExcludedTGWConnections string `form:"destination_exclude_connections,omitempty"`
	PrivateIPPeering               string `form:"private_ip_peering,omitempty"`
	InsaneModeOverInternet         bool   `form:"insane_mode_over_internet,omitempty"`
	InsaneModeTunnelCount          int    `json:"insane_mode_tunnel_count"`
	TunnelCount                    int    `form:"tunnel_count"`
	CID                            string `form:"CID,omitempty"`
	Action                         string `form:"action,omitempty"`
	SingleTunnel                   string `form:"single_tunnel,omitempty"`
	NoMaxPerformance               bool   `form:"no_max_performance,omitempty"`
}

type TransitGatewayPeeringAPIResp struct {
	Return  bool                      `json:"return"`
	Results [][]TransitGatewayPeering `json:"results"`
	Reason  string                    `json:"reason"`
}

type TransitGatewayPeeringDetailsAPIResp struct {
	Return  bool                                `json:"return"`
	Results TransitGatewayPeeringDetailsResults `json:"results"`
	Reason  string                              `json:"reason"`
}

type TransitGatewayPeeringDetailsResults struct {
	Site1                  TransitGatewayPeeringDetail `json:"site_1"`
	Site2                  TransitGatewayPeeringDetail `json:"site_2"`
	PrivateNetworkPeering  bool                        `json:"private_network_peering"` // over private network
	Tunnels                []TunnelsDetail             `json:"tunnels"`
	InsaneModeOverInternet bool                        `json:"insane_mode_over_internet"`
	InsaneModeTunnelCount  int                         `json:"insane_mode_tunnel_count"`
	TunnelCount            int                         `json:"tunnel_count"`
	EnableJumboFrame       bool                        `json:"jumbo_frame"` // jumbo frame
	NoMaxPerformance       bool                        `json:"no_max_performance"`
	Gateway1LogicalIfNames []string                    `json:"gateway1_logical_ifnames"`
	Gateway2LogicalIfNames []string                    `json:"gateway2_logical_ifnames"`
}

type TransitGatewayPeeringDetail struct {
	ExcludedCIDRs          []string `json:"exclude_filter_list"`
	ExcludedTGWConnections []string `json:"exclude_connections"`
	ConnBGPPrependAsPath   string   `json:"conn_bgp_prepend_as_path"`
}

type TunnelsDetail struct {
	LicenseId      [][]string `json:"license_id"`
	SubTunnelCount int        `json:"sub_tunnel_count"`
}

func (c *Client) CreateTransitGatewayPeering(ctx context.Context, transitGatewayPeering *TransitGatewayPeering) error {
	transitGatewayPeering.CID = c.CID
	transitGatewayPeering.Action = "create_inter_transit_gateway_peering"
	return c.PostAPIContext2(ctx, nil, transitGatewayPeering.Action, transitGatewayPeering, BasicCheck)
	// return c.PostAPI(transitGatewayPeering.Action, transitGatewayPeering, BasicCheck)
}

func (c *Client) GetTransitGatewayPeering(transitGatewayPeering *TransitGatewayPeering) error {
	form := map[string]string{
		"CID":    c.CID,
		"action": "list_inter_transit_gateway_peering",
	}

	var data TransitGatewayPeeringAPIResp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return err
	}

	if len(data.Results) == 0 {
		log.Errorf("Transit gateway peering with gateways %s and %s not found",
			transitGatewayPeering.TransitGatewayName1, transitGatewayPeering.TransitGatewayName2)
		return ErrNotFound
	}
	peeringList := data.Results
	for i := range peeringList {
		for j := range peeringList[i] {
			if peeringList[i][j].TransitGatewayName1 == transitGatewayPeering.TransitGatewayName1 &&
				peeringList[i][j].TransitGatewayName2 == transitGatewayPeering.TransitGatewayName2 ||
				peeringList[i][j].TransitGatewayName1 == transitGatewayPeering.TransitGatewayName2 &&
					peeringList[i][j].TransitGatewayName2 == transitGatewayPeering.TransitGatewayName1 {
				log.Debugf("Found %s<->%s transit gateway peering: %#v",
					transitGatewayPeering.TransitGatewayName1,
					transitGatewayPeering.TransitGatewayName2, peeringList[i][j])
				return nil
			}
		}
	}
	return ErrNotFound
}

func (c *Client) GetTransitGatewayPeeringDetails(transitGatewayPeering *TransitGatewayPeering) (*TransitGatewayPeering, error) {
	form := map[string]string{
		"action":   "get_inter_transit_gateway_peering_details",
		"CID":      c.CID,
		"gateway1": transitGatewayPeering.TransitGatewayName1,
		"gateway2": transitGatewayPeering.TransitGatewayName2,
	}
	check := func(action, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "does not exist") || strings.Contains(reason, "not found") {
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s %s failed: %s", action, method, reason)
		}
		return nil
	}
	var data TransitGatewayPeeringDetailsAPIResp
	err := c.GetAPI(&data, form["action"], form, check)
	if err != nil {
		return nil, err
	}

	if data.Results.Site1.ExcludedCIDRs != nil {
		transitGatewayPeering.Gateway1ExcludedCIDRsSlice = data.Results.Site1.ExcludedCIDRs
	}
	if data.Results.Site1.ExcludedTGWConnections != nil {
		transitGatewayPeering.Gateway1ExcludedTGWConnectionsSlice = data.Results.Site1.ExcludedTGWConnections
	}
	if data.Results.Site2.ExcludedCIDRs != nil {
		transitGatewayPeering.Gateway2ExcludedCIDRsSlice = data.Results.Site2.ExcludedCIDRs
	}
	if data.Results.Site2.ExcludedTGWConnections != nil {
		transitGatewayPeering.Gateway2ExcludedTGWConnectionsSlice = data.Results.Site2.ExcludedTGWConnections
	}

	if data.Results.PrivateNetworkPeering {
		transitGatewayPeering.PrivateIPPeering = "yes"
	} else {
		transitGatewayPeering.PrivateIPPeering = "no"
	}
	transitGatewayPeering.InsaneModeOverInternet = data.Results.InsaneModeOverInternet
	transitGatewayPeering.TunnelCount = data.Results.TunnelCount
	transitGatewayPeering.PrependAsPath1 = data.Results.Site1.ConnBGPPrependAsPath
	transitGatewayPeering.PrependAsPath2 = data.Results.Site2.ConnBGPPrependAsPath
	transitGatewayPeering.NoMaxPerformance = data.Results.NoMaxPerformance
	transitGatewayPeering.InsaneModeTunnelCount = data.Results.InsaneModeTunnelCount
	if len(data.Results.Tunnels) >= 1 && data.Results.Tunnels[0].SubTunnelCount == 1 {
		transitGatewayPeering.SingleTunnel = "yes"
	} else {
		transitGatewayPeering.SingleTunnel = "no"
	}

	return transitGatewayPeering, nil
}

func (c *Client) UpdateTransitGatewayPeering(transitGatewayPeering *TransitGatewayPeering) error {
	transitGatewayPeering.CID = c.CID
	transitGatewayPeering.Action = "edit_inter_transit_gateway_peering"

	return c.PostAPI(transitGatewayPeering.Action, transitGatewayPeering, BasicCheck)
}

func (c *Client) UpdateTransitGatewayPeeringTunnelCount(transitGatewayPeering *TransitGatewayPeeringEdit) error {
	transitGatewayPeering.CID = c.CID
	transitGatewayPeering.Action = "edit_inter_transit_gateway_peering"

	return c.PostAPI(transitGatewayPeering.Action, transitGatewayPeering, BasicCheck)
}

func (c *Client) DeleteTransitGatewayPeering(transitGatewayPeering *TransitGatewayPeering) error {
	form := map[string]string{
		"CID":      c.CID,
		"action":   "delete_inter_transit_gateway_peering",
		"gateway1": transitGatewayPeering.TransitGatewayName1,
		"gateway2": transitGatewayPeering.TransitGatewayName2,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) EditTransitConnectionASPathPrepend(transitGatewayPeering *TransitGatewayPeering, prependASPath []string) error {
	action := "edit_transit_connection_as_path_prepend"
	return c.PostAPI(action, struct {
		CID            string `form:"CID"`
		Action         string `form:"action"`
		GatewayName    string `form:"gateway_name"`
		ConnectionName string `form:"connection_name"`
		PrependASPath  string `form:"connection_as_path_prepend"`
	}{
		CID:            c.CID,
		Action:         action,
		GatewayName:    transitGatewayPeering.TransitGatewayName1,
		ConnectionName: transitGatewayPeering.TransitGatewayName2 + "-peering",
		PrependASPath:  strings.Join(prependASPath, ","),
	}, BasicCheck)
}
