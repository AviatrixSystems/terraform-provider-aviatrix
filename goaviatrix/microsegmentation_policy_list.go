package goaviatrix

import "context"

type MicrosegPortRange struct {
	Hi int `json:"hi"`
	Lo int `json:"lo"`
}

type MicrosegPolicy struct {
	Name          string                       `json:"name"`
	Action        string                       `json:"action"`
	Logging       bool                         `json:"logging,omitempty"`
	DstAppDomains []string                     `json:"dst_ads"`
	SrcAppDomains []string            `json:"src_ads"`
	PortRanges    []MicrosegPortRange `json:"port_ranges,omitempty"`
	Priority      int                 `json:"priority"`
	Protocol      string                       `json:"protocol"`
	Watch         bool                         `json:"watch,omitempty"`
	UUID          string                       `json:"uuid,omitempty"`
}

type MicrosegPolicyList struct {
	Policies []MicrosegPolicy `json:"policies"`
}

func (c *Client) CreateMicrosegPolicyList(ctx context.Context, policyList *MicrosegPolicyList) error {
	endpoint := "microseg/policy-list"
	return c.PutAPIContext25(ctx, endpoint, policyList)
}

func (c *Client) GetMicrosegPolicyList(ctx context.Context) (*MicrosegPolicyList, error) {
	endpoint := "microseg/policy-list"

	var policyList MicrosegPolicyList
	err := c.GetAPIContext25(ctx, &policyList, endpoint, nil)
	if err != nil {
		return nil, err
	} else if len(policyList.Policies) == 0 {
		return nil, ErrNotFound
	}

	return &policyList, nil
}

func (c *Client) UpdateMicrosegPolicyList(ctx context.Context, policyList *MicrosegPolicyList) error {
	endpoint := "microseg/policy-list"
	return c.PutAPIContext25(ctx, endpoint, policyList)
}

func (c *Client) DeleteMicrosegPolicyList(ctx context.Context) error {
	endpoint := "microseg/policy-list"
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}
