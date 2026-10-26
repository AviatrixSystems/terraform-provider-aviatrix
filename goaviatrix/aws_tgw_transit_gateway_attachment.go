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
	Return  bool             `json:"return"`
	Results []AttachmentInfo `json:"results"`
	Reason  string           `json:"reason"`
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
	AccessFromEdge          []string        `json:"access_from_edge"`
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
		"async":             "true",
	}
	return c.PostAsyncAPI(form["action"], form, BasicCheck)
}

func (c *Client) GetAwsTgwTransitGwAttachment(awsTgwTransitGwAttachment *AwsTgwTransitGwAttachment) (*AwsTgwTransitGwAttachment, error) {
	tgwVpcAttachment := &AwsTgwVpcAttachment{
		TgwName: awsTgwTransitGwAttachment.TgwName,
		VpcID:   awsTgwTransitGwAttachment.VpcID,
	}
	tgwAttachmentInfo, err := c.GetAwsTgwAttachmentInfo(tgwVpcAttachment)
	if err != nil {
		return nil, err
	}
	if tgwAttachmentInfo.GwName != "" {
		awsTgwTransitGwAttachment.TransitGatewayName = tgwAttachmentInfo.GwName
		awsTgwTransitGwAttachment.Region = tgwAttachmentInfo.Region
		awsTgwTransitGwAttachment.VpcAccountName = tgwAttachmentInfo.AccountName
		return awsTgwTransitGwAttachment, nil
	}
	return nil, ErrNotFound
}

func (c *Client) DeleteAwsTgwTransitGwAttachment(awsTgwTransitGwAttachment *AwsTgwTransitGwAttachment) error {
	form := map[string]string{
		"CID":      c.CID,
		"action":   "detach_vpc_from_tgw",
		"tgw_name": awsTgwTransitGwAttachment.TgwName,
		"vpc_name": awsTgwTransitGwAttachment.VpcID,
		"async":    "true",
	}
	check := func(action, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "is not attached to") {
				return nil
			}
			return fmt.Errorf("rest API %s %s failed: %s", action, method, reason)
		}
		return nil
	}
	return c.PostAsyncAPI(form["action"], form, check)
}
