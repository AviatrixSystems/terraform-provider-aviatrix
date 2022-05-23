package goaviatrix

import "context"

type PrivateModeMulticloudEndpoint struct {
	CID               string
	Action            string
	AccountName       string
	VpcId             string
	Region            string
	ControllerLbVpcId string
}

func (c *Client) CreatePrivateModeMulticloudEndpoint(ctx context.Context, privateModeMulticloudEndpoint *PrivateModeMulticloudEndpoint) error {
	privateModeMulticloudEndpoint.CID = c.CID
	privateModeMulticloudEndpoint.Action = "create_private_mode_multicloud_endpoint"
	return c.PostAPIContext2(ctx, nil, privateModeMulticloudEndpoint.Action, privateModeMulticloudEndpoint, BasicCheck)
}

func (c *Client) GetPrivateModeMulticloudEndpoint(ctx context.Context, vpcId string) (*PrivateModeMulticloudEndpoint, error) {
	action := "list_private_mode_proxies"
	form := map[string]string{
		"CID":    c.CID,
		"action": action,
	}

	type PrivateModeMulticloudEndpointResp struct {
		Result []PrivateModeMulticloudEndpoint
	}

	var resp PrivateModeMulticloudEndpointResp
	err := c.GetAPIContext(ctx, &resp, action, form, BasicCheck)
	if err != nil {
		return nil, err
	}

	for _, endpoint := range resp.Result {
		if endpoint.VpcId == vpcId {
			return &endpoint, nil
		}
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
