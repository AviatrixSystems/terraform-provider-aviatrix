package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type BranchRouterVirtualWanAttachment struct {
	ConnectionName  string `form:"connection_name"`
	BranchName      string `form:"branch_name"`
	AccountName     string `form:"account_name"`
	ResourceGroup   string `form:"arm_resource_group"`
	HubName         string `form:"virtual_wan_hub_name"`
	BranchRouterAsn string `form:"external_device_as_number"`
	Action          string `form:"action"`
	CID             string `form:"CID"`
}

func (c *Client) CreateBranchRouterVirtualWanAttachment(attachment *BranchRouterVirtualWanAttachment) error {
	attachment.Action = "attach_cloudwan_branch_to_virtual_wan"
	attachment.CID = c.CID
	resp, err := c.Post(c.baseURL, attachment)
	if err != nil {
		return errors.New("HTTP Post attach_cloudwan_branch_to_virtual_wan failed: " + err.Error())
	}

	var data APIResp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return errors.New("Reading response body attach_cloudwan_branch_to_virtual_wan failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return errors.New("Json Decode attach_cloudwan_branch_to_virtual_wan failed: " + err.Error() + "\n Body: " + b.String())
	}
	if !data.Return {
		return errors.New("Rest API attach_cloudwan_branch_to_virtual_wan Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetBranchRouterVirtualWanAttachment(attachment *BranchRouterVirtualWanAttachment) (*BranchRouterVirtualWanAttachment, error) {
	branchName, err := c.GetBranchRouterName(attachment.ConnectionName)
	if err != nil {
		return nil, fmt.Errorf("could not get branch name: %v", err)
	}

	vpcID, err := c.GetBranchRouterAttachmentVpcID(attachment.ConnectionName)
	if err != nil {
		return nil, fmt.Errorf("could not get branch router attachment VPC id: %v", err)
	}

	resp, err := c.Post(c.baseURL, struct {
		CID            string `form:"CID"`
		Action         string `form:"action"`
		ConnectionName string `form:"conn_name"`
		VpcID          string `form:"vpc_id"`
	}{
		CID:            c.CID,
		Action:         "get_site2cloud_conn_detail",
		ConnectionName: attachment.ConnectionName,
		VpcID:          vpcID,
	})
	if err != nil {
		return nil, errors.New("HTTP POST get_site2cloud_conn_detail failed: " + err.Error())
	}

	var data Site2CloudConnDetailResp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return nil, errors.New("Reading response body get_site2cloud_conn_detail failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return nil, errors.New("Json Decode get_site2cloud_conn_detail failed: " + err.Error() +
			"\n Body: " + b.String())
	}

	if !data.Return {
		if strings.Contains(data.Reason, "does not exist") {
			return nil, ErrNotFound
		}
		return nil, errors.New("Rest API get_site2cloud_conn_detail Post failed: " + data.Reason)
	}

	return &BranchRouterVirtualWanAttachment{
		ConnectionName:  attachment.ConnectionName,
		BranchName:      branchName,
		AccountName:     data.Results.Connections.AzureAccountName,
		ResourceGroup:   data.Results.Connections.AzureResourceGroup,
		HubName:         data.Results.Connections.AzureVhubName,
		BranchRouterAsn: data.Results.Connections.BgpRemoteASN,
	}, nil
}
