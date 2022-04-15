package goaviatrix

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

type TransitGatewayPeering struct {
	TransitGatewayName1                 string `form:"gateway1,omitempty" json:"gateway_1,omitempty"`
	TransitGatewayName2                 string `form:"gateway2,omitempty" json:"gateway_2,omitempty"`
	Gateway1ExcludedCIDRs               string `form:"source_filter_cidrs,omitempty"`
	Gateway2ExcludedCIDRs               string `form:"destination_filter_cidrs,omitempty"`
	Gateway1ExcludedTGWConnections      string `form:"source_exclude_connections,omitempty"`
	Gateway2ExcludedTGWConnections      string `form:"destination_exclude_connections,omitempty"`
	PrivateIPPeering                    bool   `form:"private_ip_peering,omitempty"`
	InsaneModeOverInternet              bool   `form:"insane_mode_over_internet,omitempty"`
	TunnelCount                         int    `form:"tunnel_count,omitempty"`
	Gateway1ExcludedCIDRsSlice          []string
	Gateway2ExcludedCIDRsSlice          []string
	Gateway1ExcludedTGWConnectionsSlice []string
	Gateway2ExcludedTGWConnectionsSlice []string
	PrependAsPath1                      string
	PrependAsPath2                      string
	CID                                 string `form:"CID,omitempty"`
	Action                              string `form:"action,omitempty"`
	SingleTunnel                        bool   `form:"single_tunnel"`
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
	PrivateNetworkPeering  bool                        `json:"private_network_peering"`
	Tunnels                []TunnelsDetail             `json:"tunnels"`
	InsaneModeOverInternet bool                        `json:"insane_mode_over_internet"`
	TunnelCount            int                         `json:"tunnel_count"`
}

type TransitGatewayPeeringDetail struct {
	ExcludedCIDRs          []string `json:"exclude_filter_list"`
	ExcludedTGWConnections []string `json:"exclude_connections"`
	ConnBGPPrependAsPath   string   `json:"conn_bgp_prepend_as_path"`
}

type TunnelsDetail struct {
	LicenseId [][]string `json:"license_id"`
}

func (c *Client) CreateTransitGatewayPeering(transitGatewayPeering *TransitGatewayPeering) error {
	transitGatewayPeering.CID = c.CID
	transitGatewayPeering.Action = "create_inter_transit_gateway_peering"
	return c.PostAPI(transitGatewayPeering.Action, transitGatewayPeering, BasicCheck)
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

	transitGatewayPeering.Gateway1ExcludedCIDRsSlice = data.Results.Site1.ExcludedCIDRs
	transitGatewayPeering.Gateway1ExcludedTGWConnectionsSlice = data.Results.Site1.ExcludedTGWConnections
	transitGatewayPeering.Gateway2ExcludedCIDRsSlice = data.Results.Site2.ExcludedCIDRs
	transitGatewayPeering.Gateway2ExcludedTGWConnectionsSlice = data.Results.Site2.ExcludedTGWConnections
	transitGatewayPeering.PrivateIPPeering = data.Results.PrivateNetworkPeering
	transitGatewayPeering.InsaneModeOverInternet = data.Results.InsaneModeOverInternet
	transitGatewayPeering.TunnelCount = data.Results.TunnelCount
	transitGatewayPeering.PrependAsPath1 = data.Results.Site1.ConnBGPPrependAsPath
	transitGatewayPeering.PrependAsPath2 = data.Results.Site2.ConnBGPPrependAsPath

	if len(data.Results.Tunnels[0].LicenseId) == 1 {
		transitGatewayPeering.SingleTunnel = true
	}

	return transitGatewayPeering, nil
}

func (c *Client) UpdateTransitGatewayPeering(transitGatewayPeering *TransitGatewayPeering) error {
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
