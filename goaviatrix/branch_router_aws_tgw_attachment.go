package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
)

type BranchRouterAwsTgwAttachment struct {
	ConnectionName          string `form:"connection_name"`
	BranchName              string `form:"branch_name"`
	AwsTgwName              string `form:"tgw_name"`
	BranchRouterAsn         string `form:"external_device_as_number"`
	SecurityDomainName      string `form:"route_domain_name"`
	EnableGlobalAccelerator string `form:"enable_global_accelerator"`
	Action                  string `form:"action"`
	CID                     string `form:"CID"`
}

func (b *BranchRouterAwsTgwAttachment) ID() string {
	return b.ConnectionName + "~" + b.BranchName + "~" + b.AwsTgwName
}

func (c *Client) CreateBranchRouterAwsTgwAttachment(brata *BranchRouterAwsTgwAttachment) error {
	brata.Action = "attach_cloudwan_branch_to_aws_tgw"
	brata.CID = c.CID
	resp, err := c.Post(c.baseURL, brata)
	if err != nil {
		return errors.New("HTTP Post attach_cloudwan_branch_to_aws_tgw failed: " + err.Error())
	}

	var data APIResp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return errors.New("Reading response body attach_cloudwan_branch_to_aws_tgw failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return errors.New("Json Decode attach_cloudwan_branch_to_aws_tgw failed: " + err.Error() + "\n Body: " + b.String())
	}
	if !data.Return {
		return errors.New("Rest API attach_cloudwan_branch_to_aws_tgw Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetBranchRouterAwsTgwAttachment(brata *BranchRouterAwsTgwAttachment) (*BranchRouterAwsTgwAttachment, error) {
	brata.CID = c.CID
	brata.Action = "list_tgw_details"

	resp, err := c.Post(c.baseURL, brata)
	if err != nil {
		return nil, errors.New("HTTP Post list_tgw_details failed: " + err.Error())
	}

	var data TgwAttachmentResp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return nil, errors.New("Reading response body list_tgw_details failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_tgw_details failed: " + err.Error() + "\n Body: " + b.String())
	}
	if !data.Return {
		return nil, errors.New("Rest API list_tgw_details Post failed: " + data.Reason)
	}

	var branchRouterAttachment AttachmentInfo
	var found bool
	for _, attachment := range data.Results.Attachments {
		if attachment.VpcName == brata.ConnectionName {
			branchRouterAttachment = attachment
			found = true
		}
	}
	if !found {
		return nil, ErrNotFound
	}

	return &BranchRouterAwsTgwAttachment{
		ConnectionName:          brata.ConnectionName,
		BranchName:              brata.BranchName,
		AwsTgwName:              brata.AwsTgwName,
		BranchRouterAsn:         branchRouterAttachment.AwsSideAsn,
		SecurityDomainName:      branchRouterAttachment.SecurityDomainName,
		EnableGlobalAccelerator: strconv.FormatBool(branchRouterAttachment.EnableGlobalAccelerator),
	}, nil
}
