package goaviatrix

import (
	"context"
)

type DistributedFirewallingDeploymentPolicy struct {
	Providers   []string `json:"providers"`
	SetDefaults bool     `json:"set_defaults,omitempty"`
}

func (c *Client) GetDistributedFirewallingDeploymentPolicy(ctx context.Context) (*DistributedFirewallingDeploymentPolicy, error) {
	endpoint := "microseg/deploy-policy"
	var deploymentPolicy DistributedFirewallingDeploymentPolicy

	err := c.GetAPIContext25(ctx, &deploymentPolicy, endpoint, nil)
	if err != nil {
		return nil, err
	}

	return &deploymentPolicy, nil
}

func (c *Client) CreateDistributedFirewallingDeploymentPolicy(ctx context.Context, deploymentPolicy *DistributedFirewallingDeploymentPolicy) error {
	endpoint := "microseg/deploy-policy"
	deploymentPolicyresp := &DistributedFirewallingDeploymentPolicy{}
	return c.PostAPIContext25(ctx, &deploymentPolicyresp, endpoint, deploymentPolicy)
}
