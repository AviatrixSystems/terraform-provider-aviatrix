package goaviatrix

import "context"

type DistributedFirewallingPortRange struct {
	Hi int `json:"hi,omitempty"`
	Lo int `json:"lo"`
}

type DistributedFirewallingPolicy struct {
	Name            string                              `json:"name"`
	Action          string                              `json:"action"`
	Logging         bool                                `json:"logging,omitempty"`
	DstSmartGroups  []string                            `json:"dst_ads"`
	SrcSmartGroups  []string                            `json:"src_ads"`
	PortRanges      []DistributedFirewallingPortRange   `json:"port_ranges,omitempty"`
	Priority        int                                 `json:"priority"`
	Protocol        string                              `json:"protocol"`
	Watch           bool                                `json:"watch,omitempty"`
	InspectionRules []DistributedFirewallingInspectRule `json:"nested_rules,omitempty"`
	UUID            string                              `json:"uuid,omitempty"`
	SystemResource  bool                                `json:"system_resource,omitempty"`
}

type DistributedFirewallingInspectRule struct {
	Name           string                            `json:"name"`
	Action         string                            `json:"action"`
	Logging        bool                              `json:"logging,omitempty"`
	DstSmartGroups []string                          `json:"dst_ads"`
	SrcSmartGroups []string                          `json:"src_ads"`
	PortRanges     []DistributedFirewallingPortRange `json:"port_ranges,omitempty"`
	Priority       int                               `json:"priority"`
	Protocol       string                            `json:"protocol"`
	Watch          bool                              `json:"watch,omitempty"`
	UUID           string                            `json:"uuid,omitempty"`
}

type DistributedFirewallingPolicyList struct {
	Policies []DistributedFirewallingPolicy `json:"policies"`
}

func (c *Client) CreateDistributedFirewallingPolicyList(ctx context.Context, policyList *DistributedFirewallingPolicyList) error {
	endpoint := "microseg/policy-list"
	return c.PutAPIContext25(ctx, endpoint, policyList)
}

func (c *Client) GetDistributedFirewallingPolicyList(ctx context.Context) (*DistributedFirewallingPolicyList, error) {
	endpoint := "microseg/policy-list"

	var policyList DistributedFirewallingPolicyList
	err := c.GetAPIContext25(ctx, &policyList, endpoint, nil)
	if err != nil {
		return nil, err
	} else if len(policyList.Policies) == 0 {
		return nil, ErrNotFound
	}

	return &policyList, nil
}

func (c *Client) UpdateDistributedFirewallingPolicyList(ctx context.Context, policyList *DistributedFirewallingPolicyList) error {
	endpoint := "microseg/policy-list"
	return c.PutAPIContext25(ctx, endpoint, policyList)
}

func (c *Client) DeleteDistributedFirewallingPolicyList(ctx context.Context) error {
	endpoint := "microseg/policy-list"
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}
