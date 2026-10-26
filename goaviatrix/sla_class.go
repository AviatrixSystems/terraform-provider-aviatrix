package goaviatrix

import (
	"context"
	"fmt"
)

type SLAClass struct {
	UUID           string  `json:"uuid"`
	Name           string  `json:"name"`
	Latency        int     `json:"latency_ms"`
	Jitter         int     `json:"jitter_ms"`
	PacketDropRate float64 `json:"packet_drop_rate"`
}

type SLAClassResp struct {
	SLAClass []SLAClass `json:"sla_class"`
}

func (c *Client) CreateSLAClass(ctx context.Context, slaClass *SLAClass) (string, error) {
	endpoint := "ipsla/sla"

	type resp struct {
		UUID string `json:"uuid"`
	}

	var data resp
	err := c.PostAPIContext25(ctx, &data, endpoint, slaClass)
	if err != nil {
		return "", err
	}

	return data.UUID, nil
}

func (c *Client) GetSLAClass(ctx context.Context, uuid string) (*SLAClass, error) {
	endpoint := fmt.Sprintf("ipsla/sla/%s", uuid)

	var data SLAClassResp
	err := c.GetAPIContext25(ctx, &data, endpoint, nil)
	if err != nil {
		return nil, err
	}

	for _, slaClassResult := range data.SLAClass {
		if slaClassResult.UUID == uuid {
			return &slaClassResult, nil
		}
	}

	return nil, ErrNotFound
}

func (c *Client) UpdateSLAClass(ctx context.Context, slaClass *SLAClass, uuid string) error {
	endpoint := fmt.Sprintf("ipsla/sla/%s", uuid)
	return c.PutAPIContext25(ctx, endpoint, slaClass)
}

func (c *Client) DeleteSLAClass(ctx context.Context, uuid string) error {
	endpoint := fmt.Sprintf("ipsla/sla/%s", uuid)
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}
