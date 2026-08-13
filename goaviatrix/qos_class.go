package goaviatrix

import (
	"context"
	"fmt"
)

type QosClass struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Priority int    `json:"priority"`
}

type QosClassResp struct {
	QosClass []QosClass `json:"qos_class"`
}

func (c *Client) CreateQosClass(ctx context.Context, qosClass *QosClass) (string, error) {
	endpoint := "qos/class"

	type resp struct {
		UUID string `json:"uuid"`
	}

	var data resp
	err := c.PostAPIContext25(ctx, &data, endpoint, qosClass)
	if err != nil {
		return "", err
	}

	return data.UUID, nil
}

func (c *Client) GetQosClass(ctx context.Context, uuid string) (*QosClass, error) {
	endpoint := "qos/class"

	var data QosClassResp
	err := c.GetAPIContext25(ctx, &data, endpoint, nil)
	if err != nil {
		return nil, err
	}

	for _, qosClassResult := range data.QosClass {
		if qosClassResult.UUID == uuid {
			return &qosClassResult, nil
		}
	}

	return nil, ErrNotFound
}

func (c *Client) UpdateQosClass(ctx context.Context, qosClass *QosClass, uuid string) error {
	endpoint := fmt.Sprintf("qos/class/%s", uuid)
	return c.PutAPIContext25(ctx, endpoint, qosClass)
}

func (c *Client) DeleteQosClass(ctx context.Context, uuid string) error {
	endpoint := fmt.Sprintf("qos/class/%s", uuid)
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}
