package goaviatrix

import (
	"context"
	"strings"
)

// SmartGatewayGroupStatus is one group in get_smart_gw_fabric_status.
type SmartGatewayGroupStatus struct {
	GatewayGroupName            string `json:"gateway_group_name"`
	SmartGatewayUnderlayEnabled bool   `json:"smart_gw_underlay_enabled"`
	SmartGatewayResolverEnabled bool   `json:"smart_gw_resolver_enabled"`
	// IsCloudSpokeWithoutBgpLu is a cloud spoke without BGP-LU, where both flags do nothing.
	IsCloudSpokeWithoutBgpLu bool `json:"is_csp_spoke"`
}

// Gateway group types that use the Smart Gateway flags; mirrors GwGroupType in common/gwgroup.py.
const (
	GatewayGroupTypeSpoke       = "SPOKE"
	GatewayGroupTypeTransit     = "TRANSIT"
	GatewayGroupTypeEdgeSpoke   = "EDGESPOKE"
	GatewayGroupTypeEdgeTransit = "EDGETRANSIT"
)

// normalizeGatewayGroupType turns "GwGroupType.SPOKE" or "spoke" into "SPOKE".
func normalizeGatewayGroupType(gatewayType string) string {
	s := strings.TrimSpace(gatewayType)
	if i := strings.LastIndex(s, "."); i >= 0 {
		s = s[i+1:]
	}
	return strings.ToUpper(s)
}

// IsSmartGatewayEligibleType reports whether a gw_type uses the Smart Gateway flags.
func IsSmartGatewayEligibleType(gatewayType string) bool {
	switch normalizeGatewayGroupType(gatewayType) {
	case GatewayGroupTypeSpoke, GatewayGroupTypeTransit, GatewayGroupTypeEdgeSpoke, GatewayGroupTypeEdgeTransit:
		return true
	}
	return false
}

// GetSmartGatewayFabricStatus returns the underlay and resolver flags of every gateway group.
func (c *Client) GetSmartGatewayFabricStatus(ctx context.Context) ([]SmartGatewayGroupStatus, error) {
	action := "get_smart_gw_fabric_status"
	form := map[string]string{
		"CID":    c.CID,
		"action": action,
	}
	var resp struct {
		Results []SmartGatewayGroupStatus `json:"results"`
	}
	if err := c.GetAPIContext(ctx, &resp, action, form, BasicCheck); err != nil {
		return nil, err
	}
	return resp.Results, nil
}

// ListGatewayGroupTypes maps gateway group name to gw_type.
func (c *Client) ListGatewayGroupTypes(ctx context.Context) (map[string]string, error) {
	action := "list_gateway_groups"
	form := map[string]string{
		"CID":    c.CID,
		"action": action,
	}
	var resp struct {
		Results []struct {
			Name        string `json:"name"`
			GatewayType string `json:"gw_type"`
		} `json:"results"`
	}
	if err := c.GetAPIContext(ctx, &resp, action, form, BasicCheck); err != nil {
		return nil, err
	}
	types := make(map[string]string, len(resp.Results))
	for _, g := range resp.Results {
		if g.Name == "" {
			continue
		}
		types[g.Name] = g.GatewayType
	}
	return types, nil
}

// GetSmartGatewayResolverStatus reports whether the route resolver is on for one gateway group.
func (c *Client) GetSmartGatewayResolverStatus(ctx context.Context, group string) (bool, error) {
	action := "get_smart_gw_resolver_status"
	form := map[string]string{
		"CID":                c.CID,
		"action":             action,
		"gateway_group_name": group,
	}
	var resp struct {
		Results string `json:"results"`
	}
	if err := c.GetAPIContext(ctx, &resp, action, form, BasicCheck); err != nil {
		return false, err
	}
	return strings.Contains(resp.Results, "resolver=enabled") && !strings.Contains(resp.Results, "[cloud spoke"), nil
}

// SetSmartGatewayUnderlay turns the underlay on or off for one gateway group.
func (c *Client) SetSmartGatewayUnderlay(ctx context.Context, group string, enable bool) error {
	action := "disable_smart_gw_underlay"
	if enable {
		action = "enable_smart_gw_underlay"
	}
	return c.setSmartGatewayFlag(ctx, action, group)
}

// SetSmartGatewayResolver turns the route resolver on or off for one gateway group.
func (c *Client) SetSmartGatewayResolver(ctx context.Context, group string, enable bool) error {
	action := "disable_smart_gw_resolver"
	if enable {
		action = "enable_smart_gw_resolver"
	}
	return c.setSmartGatewayFlag(ctx, action, group)
}

// setSmartGatewayFlag posts one per-group Smart Gateway toggle.
func (c *Client) setSmartGatewayFlag(ctx context.Context, action, group string) error {
	form := map[string]string{
		"CID":                c.CID,
		"action":             action,
		"gateway_group_name": group,
	}
	return c.PostAPIContext2(ctx, nil, action, form, BasicCheck)
}
