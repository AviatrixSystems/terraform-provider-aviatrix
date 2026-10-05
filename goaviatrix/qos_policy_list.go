package goaviatrix

import (
	"context"
)

type QosPolicy struct {
	UUID         string   `json:"uuid,omitempty"`
	Name         string   `json:"name,omitempty"`
	DscpValues   []string `json:"dscps,omitempty"`
	QosClassUuid string   `json:"qos_class,omitempty"`
}

type QosPolicyList struct {
	Policies []QosPolicy `json:"policies"`
}

type QosPolicyListResp struct {
	QosPolicy []QosPolicyList `json:"qos_policies"`
}

func (c *Client) UpdateQosPolicyList(ctx context.Context, qosPolicyList *QosPolicyList) error {
	endpoint := "qos/policy"

	return c.PutAPIContext25(ctx, endpoint, qosPolicyList)
}

func (c *Client) GetQosPolicyList(ctx context.Context) (*[]QosPolicyList, error) {
	endpoint := "qos/policy"

	var data QosPolicyListResp
	err := c.GetAPIContext25(ctx, &data, endpoint, nil)
	if err != nil {
		return nil, err
	}

	if len(data.QosPolicy) == 0 {
		return nil, ErrNotFound
	}

	return &data.QosPolicy, nil
}

func (c *Client) DeleteQosPolicyList(ctx context.Context) error {
	endpoint := "qos/policy"
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}
