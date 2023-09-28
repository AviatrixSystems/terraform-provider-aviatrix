package goaviatrix

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

type SpokeTransitAttachment struct {
	Action                   string `form:"action,omitempty"`
	CID                      string `form:"CID,omitempty"`
	SpokeGwName              string `form:"spoke_gw,omitempty"`
	TransitGwName            string `form:"transit_gw,omitempty"`
	RouteTables              string `form:"route_table_list,omitempty"`
	SpokeBgpEnabled          bool
	SpokePrependAsPath       []string
	TransitPrependAsPath     []string
	EnableOverPrivateNetwork bool   `form:"over_private_network,omitempty"`
	EnableJumboFrame         bool   `form:"jumbo_frame,omitempty"`
	EnableInsaneMode         bool   `form:"insane_mode,omitempty"`
	InsaneModeTunnelNumber   int    `form:"tunnel_count,omitempty"`
	NoMaxPerformance         bool   `form:"no_max_performance,omitempty"`
	EdgeWanInterfaces        string `form:"edge_wan_interfaces,omitempty"`
	EdgeWanInterfacesResp    []string
}

type EdgeSpokeTransitAttachmentResp struct {
	Return  bool                              `json:"return"`
	Results EdgeSpokeTransitAttachmentResults `json:"results"`
	Reason  string                            `json:"reason"`
}

type EdgeSpokeTransitAttachmentResults struct {
	Site1                    SiteDetail `json:"site_1"`
	Site2                    SiteDetail `json:"site_2"`
	EnableOverPrivateNetwork bool       `json:"private_network_peering"`
	EnableJumboFrame         bool       `json:"jumbo_frame"`
	EnableInsaneMode         bool       `json:"insane_mode"`
	InsaneModeTunnelNumber   int        `json:"insane_mode_tunnel_count"`
	EdgeWanInterfaces        []string   `json:"src_wan_interfaces"`
}

type SiteDetail struct {
	ConnBgpPrependAsPath string `json:"conn_bgp_prepend_as_path"`
}

func (c *Client) CreateSpokeTransitAttachment(spokeTransitAttachment *SpokeTransitAttachment) error {
	action := "attach_spoke_to_transit_gw"
	spokeTransitAttachment.CID = c.CID
	spokeTransitAttachment.Action = action
	return c.PostAPI(action, spokeTransitAttachment, BasicCheck)
}

func (c *Client) GetSpokeTransitAttachment(spokeTransitAttachment *SpokeTransitAttachment) (*SpokeTransitAttachment, error) {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "get_gateway_info",
		"gateway_name": spokeTransitAttachment.SpokeGwName,
	}

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "does not exist") {
				log.Errorf("Couldn't find Spoke Transit Attachment: %s", reason)
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	var data GatewayDetailApiResp

	err := c.GetAPI(&data, form["action"], form, checkFunc)
	if err != nil {
		return nil, err
	}

	if data.Results.GwName == spokeTransitAttachment.SpokeGwName {
		if data.Results.TransitGwName == spokeTransitAttachment.TransitGwName || data.Results.EgressTransitGwName == spokeTransitAttachment.TransitGwName {
			spokeTransitAttachment.RouteTables = strings.Join(data.Results.RouteTables, ",")
			spokeTransitAttachment.SpokeBgpEnabled = data.Results.BgpEnabled
			return spokeTransitAttachment, nil
		}
	}

	log.Errorf("Couldn't find Aviatrix gateway %s", spokeTransitAttachment.SpokeGwName)
	return nil, ErrNotFound
}

func (c *Client) DeleteSpokeTransitAttachment(spokeTransitAttachment *SpokeTransitAttachment) error {
	action := "detach_spoke_from_transit_gw"
	spokeTransitAttachment.CID = c.CID
	spokeTransitAttachment.Action = action
	return c.PostAPI(action, spokeTransitAttachment, BasicCheck)
}

func (c *Client) GetEdgeSpokeTransitAttachment(ctx context.Context, spokeTransitAttachment *SpokeTransitAttachment) (*SpokeTransitAttachment, error) {
	form := map[string]string{
		"CID":      c.CID,
		"action":   "show_multi_cloud_transit_peering_details",
		"gateway1": spokeTransitAttachment.SpokeGwName,
		"gateway2": spokeTransitAttachment.TransitGwName,
	}

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "does not exist") {
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	var data EdgeSpokeTransitAttachmentResp

	err := c.GetAPIContext(ctx, &data, form["action"], form, checkFunc)
	if err != nil {
		return nil, err
	}

	spokeTransitAttachment.EnableOverPrivateNetwork = data.Results.EnableOverPrivateNetwork
	spokeTransitAttachment.EnableJumboFrame = data.Results.EnableJumboFrame
	spokeTransitAttachment.EnableInsaneMode = data.Results.EnableInsaneMode
	spokeTransitAttachment.InsaneModeTunnelNumber = data.Results.InsaneModeTunnelNumber
	spokeTransitAttachment.EdgeWanInterfacesResp = data.Results.EdgeWanInterfaces

	if data.Results.Site1.ConnBgpPrependAsPath != "" {
		var prependAsPath []string
		for _, str := range strings.Split(data.Results.Site1.ConnBgpPrependAsPath, " ") {
			prependAsPath = append(prependAsPath, strings.TrimSpace(str))
		}
		spokeTransitAttachment.SpokePrependAsPath = prependAsPath
	}

	if data.Results.Site2.ConnBgpPrependAsPath != "" {
		var prependAsPath []string
		for _, str := range strings.Split(data.Results.Site2.ConnBgpPrependAsPath, " ") {
			prependAsPath = append(prependAsPath, strings.TrimSpace(str))
		}
		spokeTransitAttachment.TransitPrependAsPath = prependAsPath
	}

	return spokeTransitAttachment, nil
}
