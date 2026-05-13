package goaviatrix

import (
	"context"
	"fmt"
)

type DCFPolicyList struct {
	AttachTo       string                 `json:"attach_to,omitempty"`
	Name           string                 `json:"name"`
	Policies       []DCFPolicy            `json:"policies"`
	SystemResource bool                   `json:"system_resource,omitempty"`
	UUID           string                 `json:"uuid,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

type DCFPolicy struct {
	Action                 string                 `json:"action"`
	DecryptPolicy          string                 `json:"decrypt_policy,omitempty"`
	DstSmartGroups         []string               `json:"dst_ads"`
	ExcludeSgOrchestration bool                   `json:"exclude_sg_orchestration,omitempty"`
	FlowAppRequirement     string                 `json:"flow_app_requirement,omitempty"`
	Logging                bool                   `json:"logging,omitempty"`
	Name                   string                 `json:"name"`
	PortRanges             []DCFPortRange         `json:"port_ranges,omitempty"`
	Priority               int                    `json:"priority"`
	Protocol               string                 `json:"protocol"`
	SrcSmartGroups         []string               `json:"src_ads"`
	SystemResource         bool                   `json:"system_resource,omitempty"`
	TLSProfile             string                 `json:"tls_profile,omitempty"`
	UUID                   string                 `json:"uuid,omitempty"`
	Watch                  bool                   `json:"watch,omitempty"`
	WebGroups              []string               `json:"web_filters,omitempty"`
	LogProfile             string                 `json:"log_profile,omitempty"`
	Metadata               map[string]interface{} `json:"metadata,omitempty"`
}

type DCFPortRange struct {
	Hi int `json:"hi,omitempty"`
	Lo int `json:"lo"`
}

func (c *Client) CreateDCFPolicyList(ctx context.Context, policyList *DCFPolicyList) (string, error) {
	endpoint := "microseg/policy-list3"
	policyList.Metadata = map[string]interface{}{
		"terraform": map[string]string{
			"resource_type": "terraform-policy-list",
		},
	}
	for i := range policyList.Policies {
		policyList.Policies[i].Metadata = map[string]interface{}{
			"terraform": map[string]string{
				"resource_type": "terraform-policy",
			},
		}
	}
	var data DCFPolicyList
	err := c.PostAPIContext25(ctx, &data, endpoint, policyList)
	if err != nil {
		return "", err
	}

	return data.UUID, nil
}

func (c *Client) GetDCFPolicyList(ctx context.Context, uuid string) (*DCFPolicyList, error) {
	endpoint := fmt.Sprintf("microseg/policy-list3/%s", uuid)

	var policyList DCFPolicyList
	err := c.GetAPIContext25(ctx, &policyList, endpoint, nil)
	if err != nil {
		return nil, err
	}

	return &policyList, nil
}

func (c *Client) UpdateDCFPolicyList(ctx context.Context, policyList *DCFPolicyList) error {
	endpoint := fmt.Sprintf("microseg/policy-list3/%s", policyList.UUID)
	return c.PutAPIContext25(ctx, endpoint, policyList)
}

func (c *Client) DeleteDCFPolicyList(ctx context.Context, uuid string) error {
	endpoint := fmt.Sprintf("microseg/policy-list3/%s", uuid)
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}
