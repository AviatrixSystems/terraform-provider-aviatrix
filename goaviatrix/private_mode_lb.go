package goaviatrix

import (
	"context"
	"fmt"
	"strings"
)

type PrivateModeLb struct {
	CID                   string `json:"CID"`
	Action                string `json:"action"`
	AccountName           string `json:"account_name"`
	VpcId                 string `json:"vpc_id"`
	Region                string `json:"region"`
	LbType                string `json:"lb_type"`
	MulticloudAccessVpcId string `json:"endpoint_vpc_id,omitempty"`
	EdgeVpc               bool   `json:"edge_vpc,omitempty"`
	Proxies               []PrivateModeMulticloudProxy
}

type PrivateModeLbRead struct {
	AccountName           string `json:"account_name"`
	VpcId                 string `json:"vpc_id"`
	Region                string `json:"region"`
	LbType                string `json:"lb_type"`
	MulticloudAccessVpcId string `json:"multicloud_access_vpc_id,omitempty"`
	EdgeVpc               bool   `json:"edge_vpc,omitempty"`
	Proxies               []PrivateModeMulticloudProxy
}

type PrivateModeMulticloudProxy struct {
	InstanceId string
	ProxyType  string
}

func privateModeLoadBalancerCheckFunc(action, method, reason string, ret bool) error {
	if !ret {
		if strings.Contains(reason, "Private Mode is not enabled") {
			return ErrNotFound
		}
		return fmt.Errorf("rest API %s %s failed: %s", action, method, reason)
	}
	return nil
}

func (c *Client) CreatePrivateModeControllerLoadBalancer(ctx context.Context, privateModeLb *PrivateModeLb) error {
	privateModeLb.CID = c.CID
	privateModeLb.Action = "create_private_mode_controller_load_balancer"
	return c.PostAPIContext2(ctx, nil, privateModeLb.Action, privateModeLb, BasicCheck)
}

func (c *Client) CreatePrivateModeMulticloudLoadBalancer(ctx context.Context, privateModeLb *PrivateModeLb) error {
	privateModeLb.CID = c.CID
	privateModeLb.Action = "create_private_mode_multicloud_load_balancer"
	return c.PostAPIContext2(ctx, nil, privateModeLb.Action, privateModeLb, BasicCheck)
}

func (c *Client) UpdatePrivateModeMulticloudProxies(ctx context.Context, privateModeLb *PrivateModeLb) error {
	action := "update_private_mode_multicloud_proxies"
	form := map[string]interface{}{
		"CID":           c.CID,
		"action":        action,
		"account_name":  privateModeLb.AccountName,
		"instance_info": privateModeLb.Proxies,
	}

	return c.PostAPIContext2(ctx, nil, action, form, BasicCheck)
}

func (c *Client) GetPrivateModeLoadBalancer(ctx context.Context, loadBalancerVpcId string) (*PrivateModeLbRead, error) {
	action := "get_private_mode_load_balancer_detail"
	form := map[string]string{
		"CID":                  c.CID,
		"action":               action,
		"load_balancer_vpc_id": loadBalancerVpcId,
	}

	type PrivateModeLoadBalancerResp struct {
		Results PrivateModeLbRead
	}

	var resp PrivateModeLoadBalancerResp
	err := c.PostAPIContext2(ctx, &resp, action, form, privateModeLoadBalancerCheckFunc)
	if err != nil {
		return nil, err
	}

	return &resp.Results, nil
}

func (c *Client) DeletePrivateModeLoadBalancer(ctx context.Context, loadBalancerVpcId string) error {
	action := "delete_private_mode_load_balancer"
	form := map[string]string{
		"CID":    c.CID,
		"action": action,
		"vpc_id": loadBalancerVpcId,
	}

	return c.PostAPIContext2(ctx, nil, action, form, BasicCheck)
}
