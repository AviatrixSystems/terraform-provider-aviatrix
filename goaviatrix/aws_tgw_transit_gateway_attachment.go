package goaviatrix

import (
	"encoding/json"
	"fmt"
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
	VpcID                   string `json:"vpc_id"`
	GwName                  string `json:"avx_gw_name"`
	AccountName             string `json:"acct_name"`
	Region                  string `json:"region"`
	SecurityDomainName      string `json:"associated_route_domain_name"`
	VpcName                 string `json:"vpc_name"`
	AwsSideAsn              string
	AwsSideAsnRaw           json.RawMessage `json:"aws_side_asn"`
	EnableGlobalAccelerator bool            `json:"enable_acceleration"`
}

func (c *Client) CreateAwsTgwTransitGwAttachment(awsTgwTransitGwAttachment *AwsTgwTransitGwAttachment) error {
	form := map[string]string{
		"CID":               c.CID,
		"action":            "attach_vpc_to_tgw",
		"region":            awsTgwTransitGwAttachment.Region,
		"vpc_account_name":  awsTgwTransitGwAttachment.VpcAccountName,
		"vpc_name":          awsTgwTransitGwAttachment.VpcID,
		"tgw_name":          awsTgwTransitGwAttachment.TgwName,
		"route_domain_name": "Aviatrix_Edge_Domain",
		"gateway_name":      awsTgwTransitGwAttachment.TransitGatewayName,
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) GetAwsTgwTransitGwAttachment(awsTgwTransitGwAttachment *AwsTgwTransitGwAttachment) (*AwsTgwTransitGwAttachment, error) {
	var data TgwAttachmentResp
	form := map[string]string{
		"CID":      c.CID,
		"action":   "list_tgw_details",
		"tgw_name": awsTgwTransitGwAttachment.TgwName,
	}
	check := func(action, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "does not exist") {
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s Post failed: %s", action, reason)
		}
		return nil
	}
	err := c.GetAPI(&data, form["action"], form, check)
	if err != nil {
		return nil, err
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
	form := map[string]string{
		"CID":      c.CID,
		"action":   "detach_vpc_from_tgw",
		"tgw_name": awsTgwTransitGwAttachment.TgwName,
		"vpc_name": awsTgwTransitGwAttachment.VpcID,
	}
	check := func(action, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "is not attached to") {
				return nil
			}
			return fmt.Errorf("rest API %s Post failed: %s", action, reason)
		}
		return nil
	}
	return c.PostAPI(form["action"], form, check)
}
