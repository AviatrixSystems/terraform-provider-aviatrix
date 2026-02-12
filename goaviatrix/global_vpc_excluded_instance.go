package goaviatrix

import (
	"context"
	"fmt"
	"strings"
)

type GlobalVpcExcludedInstance struct {
	UUID         string `json:"id"`
	AccountName  string `json:"account"`
	InstanceName string `json:"instance_name"`
	Region       string `json:"region"`
}

type GlobalVpcExcludedInstanceResp struct {
	Return  bool                      `json:"return"`
	Results GlobalVpcExcludedInstance `json:"results"`
	Reason  string                    `json:"reason"`
}

func (c *Client) CreateGlobalVpcExcludedInstance(ctx context.Context, globalVpcExcludedInstance *GlobalVpcExcludedInstance) (string, error) {
	endpoint := "globalvpc/exclude_list"

	var data GlobalVpcExcludedInstanceResp
	err := c.PostAPIContext25(ctx, &data, endpoint, globalVpcExcludedInstance)
	if err != nil {
		return "", err
	}

	return data.Results.UUID, nil
}

func (c *Client) GetGlobalVpcExcludedInstance(ctx context.Context, uuid string) (*GlobalVpcExcludedInstance, error) {
	endpoint := fmt.Sprintf("globalvpc/exclude_list/%s", uuid)

	var data GlobalVpcExcludedInstanceResp
	err := c.GetAPIContext25(ctx, &data, endpoint, nil)
	if err != nil {
		if strings.Contains(data.Reason, "does not exist") {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &data.Results, nil
}

func (c *Client) UpdateGlobalVpcExcludedInstance(ctx context.Context, globalVpcExcludedInstance *GlobalVpcExcludedInstance, uuid string) error {
	endpoint := fmt.Sprintf("globalvpc/exclude_list/%s", uuid)
	return c.PutAPIContext25(ctx, endpoint, globalVpcExcludedInstance)
}

func (c *Client) DeleteGlobalVpcExcludedInstance(ctx context.Context, uuid string) error {
	endpoint := fmt.Sprintf("globalvpc/exclude_list/%s", uuid)
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}
