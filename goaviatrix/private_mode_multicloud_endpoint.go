package goaviatrix

import (
	"context"
)

type PrivateModeMulticloudEndpoint struct {
	CID               string `json:"CID"`
	Action            string `json:"action"`
	AccountName       string `json:"account_name"`
	VpcId             string `json:"endpoint_vpc_id"`
	Region            string `json:"region"`
	ControllerLbVpcId string `json:"load_balancer_vpc_id"`
}

type PrivateModeMultiCloudEndpointRead struct {
	AccountName       string `json:"account_name"`
	VpcId             string `json:"endpoint_vpc_id"`
	Region            string `json:"region"`
	ControllerLbVpcId string `json:"endpoint_lb_vpc_id"`
}

func (c *Client) CreatePrivateModeMulticloudEndpoint(ctx context.Context, privateModeMulticloudEndpoint *PrivateModeMulticloudEndpoint) error {
	privateModeMulticloudEndpoint.CID = c.CID
	privateModeMulticloudEndpoint.Action = "create_private_mode_multicloud_endpoint"
	return c.PostAPIContext2(ctx, nil, privateModeMulticloudEndpoint.Action, privateModeMulticloudEndpoint, BasicCheck)
}

func (c *Client) GetPrivateModeMulticloudEndpoint(ctx context.Context, vpcId string) (*PrivateModeMultiCloudEndpointRead, error) {
	action := "list_private_mode_multicloud_endpoints"
	form := map[string]string{
		"CID":    c.CID,
		"action": action,
	}

	type PrivateModeMulticloudEndpointResp struct {
		Results map[string]PrivateModeMultiCloudEndpointRead
	}

	var resp PrivateModeMulticloudEndpointResp
	err := c.GetAPIContext(ctx, &resp, action, form, BasicCheck)
	if err != nil {
		return nil, err
	}

	if endpoint, ok := resp.Results[vpcId]; ok {
		return &endpoint, nil
	}
	return nil, ErrNotFound
}

func (c *Client) DeletePrivateModeMulticloudEndpoint(ctx context.Context, vpcId string) error {
	action := "delete_private_mode_multicloud_endpoint"
	form := map[string]string{
		"CID":    c.CID,
		"action": action,
		"vpc_id": vpcId,
	}

	return c.PostAPIContext2(ctx, nil, action, form, BasicCheck)
}
