package goaviatrix

import (
	"context"
	"reflect"
)

type ControllerMetadata struct {
	Region     string `json:"region"`
	VpcId      string `json:"vpc_id"`
	InstanceId string `json:"instance_id"`
	CloudType  string `json:"cloud_type"`
}

func (c *Client) GetControllerMetadata(ctx context.Context) (*ControllerMetadata, error) {
	endpoint := "get-controller-metadata"

	var data ControllerMetadata
	err := c.GetAPIContext25(ctx, &data, endpoint, nil)
	if err != nil {
		return nil, err
	}

	v := reflect.ValueOf(data)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Interface() != "" {
			return &data, nil
		}
	}

	return nil, ErrNotFound
}
