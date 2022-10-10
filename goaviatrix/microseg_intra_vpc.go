package goaviatrix

import "context"

type MicrosegIntraVpc struct {
	VpcId       string `json:"vpc_id"`
	AccountName string `json:"account_name"`
	Region      string `json:"region,omitempty"`
}

type MicrosegIntraVpcList struct {
	VPCs []MicrosegIntraVpc `json:"vpcs"`
}

func (c *Client) CreateMicrosegIntraVpc(ctx context.Context, vpcList *MicrosegIntraVpcList) error {
	endpoint := "microseg/intra-vpc"
	return c.PutAPIContext25(ctx, endpoint, vpcList)
}

func (c *Client) GetMicrosegIntraVpc(ctx context.Context) (*MicrosegIntraVpcList, error) {
	endpoint := "microseg/intra-vpc"

	var vpcList MicrosegIntraVpcList
	err := c.GetAPIContext25(ctx, &vpcList, endpoint, nil)
	if err != nil {
		return nil, err
	} else if len(vpcList.VPCs) == 0 {
		return nil, ErrNotFound
	}

	return &vpcList, nil
}

func (c *Client) DeleteMicrosegIntraVpc(ctx context.Context) error {
	endpoint := "microseg/intra-vpc"
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}
