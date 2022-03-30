package goaviatrix

import "context"

type MicrosegmentationPortRange struct {
	Hi int `json:"hi"`
	Lo int `json:"lo"`
}

type MicrosegmentationPolicy struct {
	Name          string                       `json:"name"`
	Action        string                       `json:"action"`
	Logging       bool                         `json:"logging,omitempty"`
	DstAppDomains []string                     `json:"dst_ads"`
	SrcAppDomains []string                     `json:"src_ads"`
	PortRanges    []MicrosegmentationPortRange `json:"port_ranges,omitempty"`
	Priority      int                          `json:"priority"`
	Protocol      string                       `json:"protocol"`
	Watch         bool                         `json:"watch,omitempty"`
	UUID          string                       `json:"uuid,omitempty"`
}

type MicrosegmentationPolicyList struct {
	Policies []MicrosegmentationPolicy `json:"policies"`
}

func (c *Client) CreateMicrosegmentationPolicyList(ctx context.Context, policyList *MicrosegmentationPolicyList) error {
	endpoint := "microseg/policy-list"
	return c.PutAPIContext25(ctx, endpoint, policyList)
}

func (c *Client) GetMicrosegmentationPolicyList(ctx context.Context) (*MicrosegmentationPolicyList, error) {
	endpoint := "microseg/policy-list"

	var policyList MicrosegmentationPolicyList
	err := c.GetAPIContext25(ctx, &policyList, endpoint, nil)
	if err != nil {
		return nil, err
	} else if len(policyList.Policies) == 0 {
		return nil, ErrNotFound
	}

	return &policyList, nil
}

func (c *Client) UpdateMicrosegmentationPolicyList(ctx context.Context, policyList *MicrosegmentationPolicyList) error {
	endpoint := "microseg/policy-list"
	return c.PutAPIContext25(ctx, endpoint, policyList)
}

func (c *Client) DeleteMicrosegmentationPolicyList(ctx context.Context) error {
	endpoint := "microseg/policy-list"
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}
