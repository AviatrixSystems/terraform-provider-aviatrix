package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
)

type AwsTgwTransitGwAttachment struct {
	Action             string `form:"action,omitempty"`
	CID                string `form:"CID,omitempty"`
	TgwName            string `form:"tgw_name"`
	Region             string `form:"region"`
	SecurityDomainName string `form:"security_domain_name"`
	VpcAccountName     string `form:"vpc_account_name"`
	VpcID              string `form:"vpc_id"`
	TransitGatewayName string
}

type TgwAttachmentResp struct {
	Return  bool            `json:"return"`
	Results AttachmentsList `json:"results"`
	Reason  string          `json:"reason"`
}

type AttachmentsList struct {
	Attachments map[string]AttachmentInfo `json:"attachments"`
}

type AttachmentInfo struct {
	VpcID       string `json:"vpc_id"`
	GwName      string `json:"avx_gw_name"`
	AccountName string `json:"acct_name"`
	Region      string `json:"region"`
}

func (c *Client) CreateAwsTgwTransitGwAttachment(awsTgwTransitGwAttachment *AwsTgwTransitGwAttachment) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for attach_vpc_to_tgw") + err.Error())
	}
	attachVpcFromTgw := url.Values{}
	attachVpcFromTgw.Add("CID", c.CID)
	attachVpcFromTgw.Add("action", "attach_vpc_to_tgw")
	attachVpcFromTgw.Add("region", awsTgwTransitGwAttachment.Region)
	attachVpcFromTgw.Add("vpc_account_name", awsTgwTransitGwAttachment.VpcAccountName)
	attachVpcFromTgw.Add("vpc_name", awsTgwTransitGwAttachment.VpcID)
	attachVpcFromTgw.Add("tgw_name", awsTgwTransitGwAttachment.TgwName)
	attachVpcFromTgw.Add("route_domain_name", "Aviatrix_Edge_Domain")
	attachVpcFromTgw.Add("gateway_name", awsTgwTransitGwAttachment.TransitGatewayName)
	Url.RawQuery = attachVpcFromTgw.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get attach_vpc_to_tgw failed: " + err.Error())
	}

	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode attach_vpc_to_tgw failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API attach_vpc_to_tgw Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetAwsTgwTransitGwAttachment(awsTgwTransitGwAttachment *AwsTgwTransitGwAttachment) (*AwsTgwTransitGwAttachment, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for 'list_tgw_details': ") + err.Error())
	}
	listTgwDetails := url.Values{}
	listTgwDetails.Add("CID", c.CID)
	listTgwDetails.Add("action", "list_tgw_details")
	listTgwDetails.Add("tgw_name", awsTgwTransitGwAttachment.TgwName)
	Url.RawQuery = listTgwDetails.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get 'list_tgw_details' failed: " + err.Error())
	}

	var data TgwAttachmentResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode 'list_tgw_details' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		if strings.Contains(data.Reason, "does not exist") {
			return nil, ErrNotFound
		}
		return nil, errors.New("Rest API 'list_tgw_details' Get failed: " + data.Reason)
	}
	if _, ok := data.Results.Attachments[awsTgwTransitGwAttachment.VpcID]; ok {
		if data.Results.Attachments[awsTgwTransitGwAttachment.VpcID].GwName != "" {
			awsTgwTransitGwAttachment.TransitGatewayName = data.Results.Attachments[awsTgwTransitGwAttachment.VpcID].GwName
			awsTgwTransitGwAttachment.Region = data.Results.Attachments[awsTgwTransitGwAttachment.VpcID].Region
			awsTgwTransitGwAttachment.VpcAccountName = data.Results.Attachments[awsTgwTransitGwAttachment.VpcID].AccountName
			return awsTgwTransitGwAttachment, nil
		}
		return nil, ErrNotFound
	}

	return nil, ErrNotFound
}

func (c *Client) DeleteAwsTgwTransitGwAttachment(awsTgwTransitGwAttachment *AwsTgwTransitGwAttachment) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'detach_vpc_from_tgw': ") + err.Error())
	}
	detachVpcFromTgw := url.Values{}
	detachVpcFromTgw.Add("CID", c.CID)
	detachVpcFromTgw.Add("action", "detach_vpc_from_tgw")
	detachVpcFromTgw.Add("tgw_name", awsTgwTransitGwAttachment.TgwName)
	detachVpcFromTgw.Add("vpc_name", awsTgwTransitGwAttachment.VpcID)
	Url.RawQuery = detachVpcFromTgw.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get 'detach_vpc_from_tgw' failed: " + err.Error())
	}

	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'detach_vpc_from_tgw' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		if strings.Contains(data.Reason, "is not attached to") {
			return nil
		}
		return errors.New("Rest API 'detach_vpc_from_tgw' Get failed: " + data.Reason)
	}

	return nil
}
