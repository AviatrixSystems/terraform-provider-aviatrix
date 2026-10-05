package goaviatrix

import "context"

type DistributedFirewallingIntraVpc struct {
	VpcId       string `json:"vpc_id"`
	AccountName string `json:"account_name"`
	Region      string `json:"region,omitempty"`
}

type DistributedFirewallingIntraVpcList struct {
	VPCs []DistributedFirewallingIntraVpc `json:"vpcs"`
}

func (c *Client) CreateDistributedFirewallingIntraVpc(ctx context.Context, vpcList *DistributedFirewallingIntraVpcList) error {
	endpoint := "microseg/intra-vpc"
	return c.PutAPIContext25(ctx, endpoint, vpcList)
}

func (c *Client) GetDistributedFirewallingIntraVpc(ctx context.Context) (*DistributedFirewallingIntraVpcList, error) {
	endpoint := "microseg/intra-vpc"

	var vpcList DistributedFirewallingIntraVpcList
	err := c.GetAPIContext25(ctx, &vpcList, endpoint, nil)
	if err != nil {
		return nil, err
	} else if len(vpcList.VPCs) == 0 {
		return nil, ErrNotFound
	}

	return &vpcList, nil
}

func (c *Client) DeleteDistributedFirewallingIntraVpc(ctx context.Context) error {
	endpoint := "microseg/intra-vpc"
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}
