package goaviatrix

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	log "github.com/sirupsen/logrus"
)

type SpokeTransitAttachment struct {
	Action                       string `form:"action,omitempty" json:"action,omitempty"`
	CID                          string `form:"CID,omitempty" json:"CID,omitempty"`
	SpokeGwName                  string `form:"spoke_gw,omitempty" json:"spoke_gw,omitempty"`
	TransitGwName                string `form:"transit_gw,omitempty" json:"transit_gw,omitempty"`
	RouteTables                  string `form:"route_table_list,omitempty" json:"route_table_list,omitempty"`
	SpokeBgpEnabled              bool
	SpokePrependAsPath           []string
	TransitPrependAsPath         []string
	EnableOverPrivateNetwork     bool   `form:"over_private_network,omitempty" json:"over_private_network,omitempty"`
	EnableJumboFrame             bool   `form:"jumbo_frame,omitempty" json:"jumbo_frame,omitempty"`
	EnableInsaneMode             bool   `form:"insane_mode,omitempty" json:"insane_mode,omitempty"`
	InsaneModeTunnelNumber       int    `form:"tunnel_count,omitempty" json:"tunnel_count,omitempty"`
	NoMaxPerformance             bool   `form:"no_max_performance,omitempty" json:"no_max_performance,omitempty"`
	EdgeWanInterfaces            string `form:"edge_wan_interfaces,omitempty" json:"edge_wan_interfaces,omitempty"`
	EdgeWanInterfacesResp        []string
	DstWanInterfaces             string   `form:"dst_wan_interfaces,omitempty" json:"dst_wan_interfaces,omitempty"`
	SpokeGatewayLogicalIfNames   []string `form:"spoke_gw_logical_ifnames,omitempty" json:"spoke_gw_logical_ifnames,omitempty"`
	TransitGatewayLogicalIfNames []string `form:"transit_gw_logical_ifnames,omitempty" json:"transit_gw_logical_ifnames,omitempty"`
	DisableActivemesh            bool     `form:"disable_activemesh,omitempty" json:"disable_activemesh,omitempty"`
	EnableFirenetForEdge         bool     `form:"enable_firenet_for_edge,omitempty" json:"enable_firenet_for_edge"`
	EnableAzAffinity             bool     `form:"enable_az_affinity,omitempty" json:"enable_az_affinity,omitempty"`
	SpokeToTransitFilterCIDRs    string   `form:"spoke_to_transit_filter_cidrs,omitempty" json:"spoke_to_transit_filter_cidrs,omitempty"`
	TransitToSpokeFilterCIDRs    string   `form:"transit_to_spoke_filter_cidrs,omitempty" json:"transit_to_spoke_filter_cidrs,omitempty"`
	// Read-back (response-only) values parsed from the peering-details API so
	// Terraform can detect out-of-band changes to the advertisement filters.
	SpokeToTransitFilterCIDRsSlice []string
	TransitToSpokeFilterCIDRsSlice []string
}

type EdgeSpokeTransitAttachmentResp struct {
	Return  bool                              `json:"return"`
	Results EdgeSpokeTransitAttachmentResults `json:"results"`
	Reason  string                            `json:"reason"`
}

type EdgeSpokeTransitAttachmentResults struct {
	Site1                        SiteDetail `json:"site_1"`
	Site2                        SiteDetail `json:"site_2"`
	EnableOverPrivateNetwork     bool       `json:"private_network_peering"`
	EnableJumboFrame             bool       `json:"jumbo_frame"`
	EnableInsaneMode             bool       `json:"insane_mode"`
	InsaneModeTunnelNumber       int        `json:"insane_mode_tunnel_count"`
	EdgeWanInterfaces            []string   `json:"src_wan_interfaces"`
	SpokeGatewayLogicalIfNames   []string   `json:"src_gw_logical_ifnames"`
	TransitGatewayLogicalIfNames []string   `json:"dst_wan_interfaces"`
	DisableActivemesh            bool       `json:"disable_activemesh,omitempty"`
	EnableFirenetForEdge         bool       `json:"enable_firenet_for_edge,omitempty"`
}

type SiteDetail struct {
	ConnBgpPrependAsPath string   `json:"conn_bgp_prepend_as_path"`
	FilterCIDRs          []string `json:"filter_list"`
}

func (c *Client) CreateSpokeTransitAttachment(ctx context.Context, spokeTransitAttachment *SpokeTransitAttachment) error {
	action := "attach_spoke_to_transit_gw"
	spokeTransitAttachment.CID = c.CID
	spokeTransitAttachment.Action = action
	return c.PostAPIContext2(ctx, nil, action, spokeTransitAttachment, BasicCheck)
}

// EditSpokeTransitAttachmentFilters updates the spoke-to-transit and
// transit-to-spoke route advertisement filter CIDRs on an existing
// spoke-transit attachment without recreating the attachment.
func (c *Client) EditSpokeTransitAttachmentFilters(ctx context.Context, spokeTransitAttachment *SpokeTransitAttachment) error {
	action := "edit_attach_spoke_to_transit_gw"
	spokeTransitAttachment.CID = c.CID
	spokeTransitAttachment.Action = action
	return c.PostAPIContext2(ctx, nil, action, spokeTransitAttachment, BasicCheck)
}

// GetSpokeTransitAttachmentFilterCIDRs reads the advertisement filter cidrs for
// a (non-edge) spoke-transit attachment from the peering-details API and maps
// them onto the SpokeTransitAttachment model, so drift detection lives on the
// attachment model rather than on the transit-peering struct. site_1 is the
// spoke side (spoke->transit), site_2 the transit side (transit->spoke).
func (c *Client) GetSpokeTransitAttachmentFilterCIDRs(spokeTransitAttachment *SpokeTransitAttachment) error {
	form := map[string]string{
		"action":   "get_inter_transit_gateway_peering_details",
		"CID":      c.CID,
		"gateway1": spokeTransitAttachment.SpokeGwName,
		"gateway2": spokeTransitAttachment.TransitGwName,
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
	if err := c.GetAPI(&data, form["action"], form, check); err != nil {
		return err
	}
	spokeTransitAttachment.SpokeToTransitFilterCIDRsSlice = data.Results.Site1.FilterCIDRs
	spokeTransitAttachment.TransitToSpokeFilterCIDRsSlice = data.Results.Site2.FilterCIDRs
	return nil
}

func (c *Client) GetSpokeTransitAttachment(spokeTransitAttachment *SpokeTransitAttachment) (*SpokeTransitAttachment, error) {
	return c.GetSpokeTransitAttachmentContext(context.Background(), spokeTransitAttachment)
}

func (c *Client) GetSpokeTransitAttachmentContext(ctx context.Context, spokeTransitAttachment *SpokeTransitAttachment) (*SpokeTransitAttachment, error) {
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

	err := c.GetAPIContext(ctx, &data, form["action"], form, checkFunc)
	if err != nil {
		return nil, err
	}

	// take care TransitGwName can be either the gateway name or the group name
	transitGrpName := spokeTransitAttachment.TransitGwName
	transitGw, err := c.GetGateway(&Gateway{GwName: spokeTransitAttachment.TransitGwName})
	if err == nil {
		transitGrpName = transitGw.GroupName
	}

	if data.Results.GwName == spokeTransitAttachment.SpokeGwName {
		if data.Results.TransitGwName == spokeTransitAttachment.TransitGwName ||
			data.Results.EgressTransitGwName == spokeTransitAttachment.TransitGwName ||
			data.Results.TransitGwName == transitGrpName || data.Results.EgressTransitGwName == transitGrpName {
			spokeTransitAttachment.RouteTables = strings.Join(data.Results.RouteTables, ",")
			spokeTransitAttachment.SpokeBgpEnabled = data.Results.BgpEnabled
			return spokeTransitAttachment, nil
		}
	}

	log.Errorf("Couldn't find spoke transit attachment %s to transit %s", spokeTransitAttachment.SpokeGwName, transitGrpName)
	// ErrNotFound lets Terraform Read clear state on refresh when the spoke exists
	// but is not attached to this transit (e.g. attach failed or was never applied).
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
	spokeTransitAttachment.DisableActivemesh = data.Results.DisableActivemesh
	spokeTransitAttachment.EnableFirenetForEdge = data.Results.EnableFirenetForEdge

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

	// set the spoke gateway logical interface names
	if data.Results.SpokeGatewayLogicalIfNames != nil {
		spokeTransitAttachment.SpokeGatewayLogicalIfNames = data.Results.SpokeGatewayLogicalIfNames
	}

	// set the transit gateway logical interface names
	if data.Results.TransitGatewayLogicalIfNames != nil {
		spokeTransitAttachment.TransitGatewayLogicalIfNames = data.Results.TransitGatewayLogicalIfNames
	}

	// advertisement filter cidrs, read back so Terraform can detect drift.
	// site_1 is the spoke side (spoke->transit), site_2 the transit side.
	spokeTransitAttachment.SpokeToTransitFilterCIDRsSlice = data.Results.Site1.FilterCIDRs
	spokeTransitAttachment.TransitToSpokeFilterCIDRsSlice = data.Results.Site2.FilterCIDRs

	return spokeTransitAttachment, nil
}

func DiffSuppressFuncEdgeSpokeTransitAttachmentEdgeWanInterfaces(k, old, new string, d *schema.ResourceData) bool {
	o, n := d.GetChange("edge_wan_interfaces")

	oSet, ok := o.(*schema.Set)
	if !ok {
		return false
	}
	nSet, ok := n.(*schema.Set)
	if !ok {
		return false
	}

	edgeWanInterfacesOld := ExpandStringList(oSet.List())
	edgeWanInterfacesNew := ExpandStringList(nSet.List())

	defaultWanInterfaces := getStringSet(d, "default_edge_wan_interfaces")

	if (len(edgeWanInterfacesOld) != 0 && Equivalent(edgeWanInterfacesOld, defaultWanInterfaces) && len(edgeWanInterfacesNew) == 0) ||
		(len(edgeWanInterfacesOld) == 0 && len(edgeWanInterfacesNew) != 0 && Equivalent(edgeWanInterfacesNew, defaultWanInterfaces)) {
		return true
	}

	return false
}
