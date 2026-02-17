package goaviatrix

import (
	"context"
)

type PortRange struct {
	Lo int `json:"lo"`
	Hi int `json:"hi"`
}

type TCPolicy struct {
	UUID          string      `json:"uuid,omitempty"`
	Name          string      `json:"name,omitempty"`
	SrcSgs        []string    `json:"src_sgs,omitempty"`
	DstSgs        []string    `json:"dst_sgs,omitempty"`
	PortRanges    []PortRange `json:"port_ranges,omitempty"`
	Protocol      string      `json:"protocol,omitempty"`
	LinkHierarchy string      `json:"link_hierarchy,omitempty"`
	SlaClass      string      `json:"sla_class,omitempty"`
	Logging       bool        `json:"logging,omitempty"`
	RouteType     string      `json:"route_type,omitempty"`
}

type PolicyList struct {
	Policies []TCPolicy `json:"policies"`
}

type TrafficClassifierResp struct {
	TrafficClassifier []PolicyList `json:"traffic_classifier_policies"`
}

func (c *Client) CreateTrafficClassifier(ctx context.Context, policyList *PolicyList) error {
	endpoint := "ipsla/traffic-classifier"

	return c.PutAPIContext25(ctx, endpoint, policyList)
}

func (c *Client) GetTrafficClassifier(ctx context.Context) (*[]PolicyList, error) {
	endpoint := "ipsla/traffic-classifier"

	var data TrafficClassifierResp
	err := c.GetAPIContext25(ctx, &data, endpoint, nil)
	if err != nil {
		return nil, err
	}

	if len(data.TrafficClassifier) == 0 {
		return nil, ErrNotFound
	}

	return &data.TrafficClassifier, nil
}

func (c *Client) DeleteTrafficClassifier(ctx context.Context) error {
	endpoint := "ipsla/traffic-classifier"
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}
