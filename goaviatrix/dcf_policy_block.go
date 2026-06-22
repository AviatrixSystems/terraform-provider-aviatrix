package goaviatrix

import (
	"context"
	"fmt"
)

type DCFPolicyBlock struct {
	Name           string                 `json:"name"`
	SubPolicies    []DCFSubPolicy         `json:"sub_policies"`
	SystemResource bool                   `json:"system_resource,omitempty"`
	UUID           string                 `json:"uuid,omitempty"`
	AttachTo       string                 `json:"attach_to,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

type DCFSubPolicy struct {
	Block           string           `json:"block,omitempty"`
	List            string           `json:"list,omitempty"`
	Name            string           `json:"name"`
	AttachmentPoint *AttachmentPoint `json:"attachment_point,omitempty"`
	Priority        int              `json:"priority"`
}

type AttachmentPoint struct {
	Name       string `json:"name,omitempty"`
	TargetUUID string `json:"target_uuid,omitempty"`
	UUID       string `json:"uuid,omitempty"`
}

func (c *Client) CreateDCFPolicyBlock(ctx context.Context, policyBlock *DCFPolicyBlock) (string, error) {
	endpoint := "microseg/policy-list3"

	for _, sp := range policyBlock.SubPolicies {
		if sp.AttachmentPoint != nil && sp.AttachmentPoint.Name != "" {
			fmt.Printf("Processing subpolicy: %s\n", sp.AttachmentPoint.Name)
			attachmentPoint, err := c.GetDCFAttachmentPoint(ctx, sp.AttachmentPoint.Name)
			if err != nil {
				return "", err
			}
			sp.AttachmentPoint.UUID = attachmentPoint.AttachmentPointID
		}
	}

	policyBlock.Metadata = map[string]interface{}{
		"terraform": map[string]string{
			"resource_type": "terraform-policy-block",
		},
	}

	var data DCFPolicyBlock
	err := c.PostAPIContext25(ctx, &data, endpoint, policyBlock)
	if err != nil {
		return "", err
	}

	return data.UUID, nil
}

func (c *Client) GetDCFPolicyBlock(ctx context.Context, uuid string) (*DCFPolicyBlock, error) {
	endpoint := fmt.Sprintf("microseg/policy-list3/%s", uuid)

	var policyBlock DCFPolicyBlock
	err := c.GetAPIContext25(ctx, &policyBlock, endpoint, nil)
	if err != nil {
		return nil, err
	}

	return &policyBlock, nil
}

func (c *Client) UpdateDCFPolicyBlock(ctx context.Context, policyBlock *DCFPolicyBlock) error {
	endpoint := fmt.Sprintf("microseg/policy-list3/%s", policyBlock.UUID)
	for _, sp := range policyBlock.SubPolicies {
		if sp.AttachmentPoint != nil && sp.AttachmentPoint.Name != "" {
			attachmentPoint, err := c.GetDCFAttachmentPoint(ctx, sp.AttachmentPoint.Name)
			if err != nil {
				return err
			}
			sp.AttachmentPoint.UUID = attachmentPoint.AttachmentPointID
		}
	}
	return c.PutAPIContext25(ctx, endpoint, policyBlock)
}

func (c *Client) DeleteDCFPolicyBlock(ctx context.Context, uuid string) error {
	endpoint := fmt.Sprintf("microseg/policy-list3/%s", uuid)
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}
