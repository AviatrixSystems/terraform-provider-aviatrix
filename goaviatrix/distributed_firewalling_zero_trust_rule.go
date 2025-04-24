package goaviatrix

import (
	"context"
	"fmt"
	"strconv"
)

type DistributedFirewallingZeroTrustRule struct {
	Action  string `json:"action"`
	Logging bool   `json:"logging"`
	UUID    string `json:"uuid,omitempty"`
}

func (c *Client) UpdateDistributedFirewallingZeroTrust(ctx context.Context, request *DistributedFirewallingZeroTrustRule) error {
	endpoint := "microseg/edit-zero-trust"
	return c.PutAPIContext25(ctx, endpoint, request)
}

func (c *Client) GetDistributedFirewallingZeroTrustRule(ctx context.Context) (map[string]string, error) {

	postRuleListUUID := "defa11a1-3000-5000-0000-000000000000"
	zeroTrustRuleUUID := "defa11a1-0000-0000-0000-000000000000"

	endpoint := "/microseg/policy-list3/" + postRuleListUUID

	var response struct {
		DcfPolicies struct {
			Policies []DistributedFirewallingZeroTrustRule `json:"policies"`
		} `json:"dcf_policies"`
	}

	// Fetch the data from the API
	err := c.GetAPIContext25(ctx, &response, endpoint, nil)
	if err != nil {
		return nil, err
	}

	// Iterate through the policies to find the rule with the matching UUID
	for _, policy := range response.DcfPolicies.Policies {
		if policy.UUID == zeroTrustRuleUUID {
			return map[string]string{
				"Action":  policy.Action,
				"Logging": strconv.FormatBool(policy.Logging),
				"UUID":    policy.UUID,
			}, nil
		}
	}

	return nil, fmt.Errorf("zero trust rule with UUID %s not found", zeroTrustRuleUUID)
}
