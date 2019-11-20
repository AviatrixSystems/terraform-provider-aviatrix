package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"strings"
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
}

type AwsTgwDirectConnEdit struct {
	TgwName                  string   `json:"tgw_name,omitempty"`
	DirectConnectAccountName string   `json:"acct_name,omitempty"`
	DxGatewayID              string   `json:"name,omitempty"`
	SecurityDomainName       string   `json:"associated_route_domain_name,omitempty"`
	AllowedPrefix            []string `json:"allowed_prefix,omitempty"`
}

type AwsTgwDirectConnResp struct {
	Return  bool                   `json:"return"`
	Results []AwsTgwDirectConnEdit `json:"results"`
	Reason  string                 `json:"reason"`
}

func (c *Client) CreateAwsTgwDirectConnect(awsTgwDirectConnect *AwsTgwDirectConnect) error {
	awsTgwDirectConnect.CID = c.CID
	awsTgwDirectConnect.Action = "attach_direct_connect_to_tgw"
	resp, err := c.Post(c.baseURL, awsTgwDirectConnect)
	if err != nil {
		return errors.New("HTTP Post attach_direct_connect_to_tgw failed: " + err.Error())
	}

	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode attach_direct_connect_to_tgw failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API attach_direct_connect_to_tgw Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetAwsTgwDirectConnect(awsTgwDirectConnect *AwsTgwDirectConnect) (*AwsTgwDirectConnect, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for list_all_tgw_attachments") + err.Error())
	}
	listAllTgwAttachments := url.Values{}
	listAllTgwAttachments.Add("CID", c.CID)
	listAllTgwAttachments.Add("action", "list_all_tgw_attachments")
	listAllTgwAttachments.Add("tgw_name", awsTgwDirectConnect.TgwName)
	Url.RawQuery = listAllTgwAttachments.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get list_all_tgw_attachments failed: " + err.Error())
	}

	var data AwsTgwDirectConnResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_all_tgw_attachments failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return nil, errors.New("Rest API list_all_tgw_attachments Get failed: " + data.Reason)
	}
	allAwsTgwDirectConn := data.Results
	for i := range allAwsTgwDirectConn {
		if allAwsTgwDirectConn[i].TgwName == awsTgwDirectConnect.TgwName && allAwsTgwDirectConn[i].DxGatewayID == awsTgwDirectConnect.DxGatewayID {
			awsTgwDirectConnect.DirectConnectAccountName = allAwsTgwDirectConn[i].DirectConnectAccountName
			awsTgwDirectConnect.SecurityDomainName = allAwsTgwDirectConn[i].SecurityDomainName
			awsTgwDirectConnect.AllowedPrefix = strings.Join(allAwsTgwDirectConn[i].AllowedPrefix, ",")
			log.Printf("[DEBUG] Found Aws Tgw Direct Conn: %#v", awsTgwDirectConnect)
			return awsTgwDirectConnect, nil
		}
	}
	return nil, ErrNotFound
}

func (c *Client) UpdateDirectConnAllowedPrefix(awsTgwDirectConnect *AwsTgwDirectConnect) error {
	awsTgwDirectConnect.CID = c.CID
	awsTgwDirectConnect.Action = "update_tgw_directconnect_allowed_prefix"
	resp, err := c.Post(c.baseURL, awsTgwDirectConnect)
	if err != nil {
		return errors.New("HTTP Post update_tgw_directconnect_allowed_prefix failed: " + err.Error())
	}

	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode update_tgw_directconnect_allowed_prefix failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API update_tgw_directconnect_allowed_prefix Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DeleteAwsTgwDirectConnect(awsTgwDirectConnect *AwsTgwDirectConnect) error {
	awsTgwDirectConnect.CID = c.CID
	awsTgwDirectConnect.Action = "detach_directconnect_from_tgw"
	resp, err := c.Post(c.baseURL, awsTgwDirectConnect)
	if err != nil {
		return errors.New("HTTP Post detach_directconnect_from_tgw failed: " + err.Error())
	}

	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode detach_directconnect_from_tgw failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API detach_directconnect_from_tgw Post failed: " + data.Reason)
	}
	return nil
}
