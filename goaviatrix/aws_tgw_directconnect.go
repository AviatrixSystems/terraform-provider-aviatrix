package goaviatrix

import (
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"
)

type AwsTgwDirectConnect struct {
	CID                      string `form:"CID,omitempty"`
	Action                   string `form:"action,omitempty"`
	TgwName                  string `form:"tgw_name,omitempty"`
	DirectConnectAccountName string `form:"directconnect_account_name,omitempty"`
	DxGatewayID              string `form:"directconnect_gateway_id,omitempty"`
	DxGatewayName            string `form:"directconnect_gateway_name,omitempty"`
	SecurityDomainName       string `form:"route_domain_name,omitempty"`
	AllowedPrefix            string `form:"allowed_prefix,omitempty"`
	DirectConnectID          string `form:"directconnect_id, omitempty"`
	LearnedCidrsApproval     string `form:"learned_cidrs_approval,omitempty"`
	Async                    bool   `form:"async,omitempty"`
}

type AwsTgwDirectConnEdit struct {
	TgwName                  string   `json:"tgw_name,omitempty"`
	DirectConnectAccountName string   `json:"acct_name,omitempty"`
	DxGatewayID              string   `json:"name,omitempty"`
	SecurityDomainName       string   `json:"associated_route_domain_name,omitempty"`
	AllowedPrefix            []string `json:"allowed_prefix,omitempty"`
	LearnedCidrsApproval     string   `json:"learned_cidrs_approval,omitempty"`
}

type AwsTgwDirectConnResp struct {
	Return  bool                   `json:"return"`
	Results []AwsTgwDirectConnEdit `json:"results"`
	Reason  string                 `json:"reason"`
}

func (c *Client) CreateAwsTgwDirectConnect(awsTgwDirectConnect *AwsTgwDirectConnect) error {
	awsTgwDirectConnect.CID = c.CID
	awsTgwDirectConnect.Action = "attach_direct_connect_to_tgw"
	awsTgwDirectConnect.Async = true
	return c.PostAsyncAPI(awsTgwDirectConnect.Action, awsTgwDirectConnect, BasicCheck)
}

func (c *Client) GetAwsTgwDirectConnect(awsTgwDirectConnect *AwsTgwDirectConnect) (*AwsTgwDirectConnect, error) {
	var data AwsTgwDirectConnResp
	form := map[string]string{
		"CID":      c.CID,
		"action":   "list_all_tgw_attachments",
		"tgw_name": awsTgwDirectConnect.TgwName,
	}
	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, errors.New("HTTP Get list_all_tgw_attachments failed: " + err.Error())
	}
	allAwsTgwDirectConn := data.Results
	for i := range allAwsTgwDirectConn {
		if allAwsTgwDirectConn[i].TgwName == awsTgwDirectConnect.TgwName && allAwsTgwDirectConn[i].DxGatewayID == awsTgwDirectConnect.DxGatewayID {
			awsTgwDirectConnect.DirectConnectAccountName = allAwsTgwDirectConn[i].DirectConnectAccountName
			awsTgwDirectConnect.SecurityDomainName = allAwsTgwDirectConn[i].SecurityDomainName
			awsTgwDirectConnect.AllowedPrefix = strings.Join(allAwsTgwDirectConn[i].AllowedPrefix, ",")
			awsTgwDirectConnect.LearnedCidrsApproval = allAwsTgwDirectConn[i].LearnedCidrsApproval
			log.Debugf("Found Aws Tgw Direct Conn: %#v", awsTgwDirectConnect)
			return awsTgwDirectConnect, nil
		}
	}
	return nil, ErrNotFound
}

func (c *Client) UpdateDirectConnAllowedPrefix(awsTgwDirectConnect *AwsTgwDirectConnect) error {
	awsTgwDirectConnect.CID = c.CID
	awsTgwDirectConnect.Action = "update_tgw_directconnect_allowed_prefix"
	return c.PostAPI(awsTgwDirectConnect.Action, awsTgwDirectConnect, BasicCheck)
}

func (c *Client) DeleteAwsTgwDirectConnect(awsTgwDirectConnect *AwsTgwDirectConnect) error {
	awsTgwDirectConnect.CID = c.CID
	awsTgwDirectConnect.Action = "detach_directconnect_from_tgw"
	return c.PostAPI(awsTgwDirectConnect.Action, awsTgwDirectConnect, BasicCheck)
}

func (c *Client) EnableDirectConnectLearnedCidrsApproval(awsTgwDirectConnect *AwsTgwDirectConnect) error {
	form := map[string]string{
		"CID":                    c.CID,
		"action":                 "enable_learned_cidrs_approval",
		"tgw_name":               awsTgwDirectConnect.TgwName,
		"attachment_name":        awsTgwDirectConnect.DxGatewayName,
		"learned_cidrs_approval": awsTgwDirectConnect.LearnedCidrsApproval,
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) DisableDirectConnectLearnedCidrsApproval(awsTgwDirectConnect *AwsTgwDirectConnect) error {
	form := map[string]string{
		"CID":                    c.CID,
		"action":                 "disable_learned_cidrs_approval",
		"tgw_name":               awsTgwDirectConnect.TgwName,
		"attachment_name":        awsTgwDirectConnect.DxGatewayName,
		"learned_cidrs_approval": awsTgwDirectConnect.LearnedCidrsApproval,
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}
