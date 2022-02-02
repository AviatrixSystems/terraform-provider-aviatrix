package goaviatrix

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type SpokeGatewaySubnetGroup struct {
	CID             string
	GatewayName     string
	SubnetGroupName string
	SubnetList      []string
}

type SpokeGatewaySubnetGroupResp struct {
	SubnetGroupName string   `json:"subnet_group_name"`
	SubnetList      []string `json:"subnet_list"`
}

type ListSpokeGatewaySubnetsResp struct {
	SubnetListAll []string `json:"subnet_list_all"`
}

func (c *Client) AddSpokeGatewaySubnetGroup(ctx context.Context, spokeGatewaySubnetGroup *SpokeGatewaySubnetGroup) error {
	form := map[string]string{
		"action":            "add_spoke_gateway_subnet_group",
		"CID":               c.CID,
		"gateway_name":      spokeGatewaySubnetGroup.GatewayName,
		"subnet_group_name": spokeGatewaySubnetGroup.SubnetGroupName,
	}

	return c.PostAPIContext(ctx, form["action"], form, BasicCheck)
}

func (c *Client) GetSpokeGatewaySubnetGroup(ctx context.Context, spokeGatewaySubnetGroup *SpokeGatewaySubnetGroup) error {
	form := map[string]string{
		"action":            "get_spoke_gateway_subnet_group",
		"CID":               c.CID,
		"gateway_name":      spokeGatewaySubnetGroup.GatewayName,
		"subnet_group_name": spokeGatewaySubnetGroup.SubnetGroupName,
	}

	type Resp struct {
		Return  bool                        `json:"return"`
		Results SpokeGatewaySubnetGroupResp `json:"results"`
		Reason  string                      `json:"reason"`
	}

	var data Resp

	check := func(action, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "Can not find subnet group") || strings.Contains(reason, "not found") {
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s %s failed: %s", action, method, reason)
		}
		return nil
	}

	err := c.GetAPIContext(ctx, &data, form["action"], form, check)
	if err != nil {
		return err
	}

	spokeGatewaySubnetGroup.SubnetList = data.Results.SubnetList

	return nil
}

func (c *Client) UpdateSpokeGatewaySubnetGroup(ctx context.Context, spokeGatewaySubnetGroup *SpokeGatewaySubnetGroup) error {
	subnetsJson, _ := json.Marshal(spokeGatewaySubnetGroup.SubnetList)

	form := map[string]string{
		"action":            "update_spoke_gateway_subnet_group",
		"CID":               c.CID,
		"gateway_name":      spokeGatewaySubnetGroup.GatewayName,
		"subnet_group_name": spokeGatewaySubnetGroup.SubnetGroupName,
		"subnet_list":       string(subnetsJson),
	}

	return c.PostAPIContext(ctx, form["action"], form, BasicCheck)
}

func (c *Client) DeleteSpokeGatewaySubnetGroup(ctx context.Context, spokeGatewaySubnetGroup *SpokeGatewaySubnetGroup) error {
	form := map[string]string{
		"action":            "delete_spoke_gateway_subnet_group",
		"CID":               c.CID,
		"gateway_name":      spokeGatewaySubnetGroup.GatewayName,
		"subnet_group_name": spokeGatewaySubnetGroup.SubnetGroupName,
	}

	return c.PostAPIContext(ctx, form["action"], form, BasicCheck)
}

func (c *Client) GetSubnetsForInspection(gatewayName string) ([]string, error) {
	form := map[string]string{
		"action":       "list_spoke_gateway_subnets",
		"CID":          c.CID,
		"gateway_name": gatewayName,
	}

	type Resp struct {
		Return  bool                        `json:"return,omitempty"`
		Results ListSpokeGatewaySubnetsResp `json:"results,omitempty"`
		Reason  string                      `json:"reason,omitempty"`
	}

	var data Resp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	return data.Results.SubnetListAll, nil
}
