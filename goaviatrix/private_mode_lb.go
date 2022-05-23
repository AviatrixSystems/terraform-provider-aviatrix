package goaviatrix

import "context"

type PrivateModeLb struct {
	CID                   string
	Action                string
	AccountName           string
	VpcId                 string
	Region                string
	CloudType             int
	LbType                string
	MulticloudAccessVpcId string
	Proxies               []PrivateModeMulticloudProxy
}

type PrivateModeMulticloudProxy struct {
	InstanceId string
	ProxyType  string
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

func (c *Client) GetPrivateModeLoadBalancer(ctx context.Context, loadBalancerVpcId string) (*PrivateModeLb, error) {
	action := "get_private_mode_load_balancer_detail"
	form := map[string]string{
		"CID":                  c.CID,
		"action":               action,
		"load_balancer_vpc_id": loadBalancerVpcId,
	}

	type PrivateModeLoadBalancerResp struct {
		Result PrivateModeLb
	}

	var resp PrivateModeLoadBalancerResp
	err := c.GetAPIContext(ctx, &resp, action, form, BasicCheck)
	if err != nil {
		return nil, err
	}

	return &resp.Result, nil
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
