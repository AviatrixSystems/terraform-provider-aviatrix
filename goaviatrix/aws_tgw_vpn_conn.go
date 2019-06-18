package goaviatrix

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"strings"
)

// VGWConn simple struct to hold VGW Connection details
type AwsTgwVpnConn struct {
	Action          string `form:"action,omitempty"`
	TgwName         string `form:"tgw_name,omitempty"`
	RouteDomainName string `form:"route_domain_name,omitempty"`
	CID             string `form:"CID,omitempty"`
	ConnName        string `form:"connection_name,omitempty"`
	PublicIP        string `form:"public_ip,omitempty"`
	OnpremASN       string `form:"onprem_asn,omitempty"`
	RemoteCIDR      string `form:"remote_cidr,omitempty"`
	VpnID           string `form:"vpn_id,omitempty"`
}

type AwsTgwVpnConnEdit struct {
	TgwName         string   `json:"tgw_name,omitempty"`
	RouteDomainName string   `json:"associated_route_domain_name,omitempty"`
	ConnName        string   `json:"vpc_name,omitempty"`
	PublicIP        string   `json:"public_ip,omitempty"`
	OnpremASN       string   `json:"aws_side_asn,omitempty"`
	RemoteCIDR      []string `json:"remote_cidrs,omitempty"`
	VpnID           string   `json:"vpc_id,omitempty"`
}

type AwsTgwVpnConnCreateResp struct {
	Return  bool   `json:"return"`
	Results string `json:"results"`
	Reason  string `json:"reason"`
}

type AwsTgwVpnConnResp struct {
	Return  bool                `json:"return"`
	Results []AwsTgwVpnConnEdit `json:"results"`
	Reason  string              `json:"reason"`
}

type GetAwsTgwVpnConnVpnIdResp struct {
	Return  bool     `json:"return"`
	Results []string `json:"results"`
	Reason  string   `json:"reason"`
}

func (c *Client) CreateAwsTgwVpnConn(awsTgwVpnConn *AwsTgwVpnConn) (string, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return "", errors.New(("url Parsing failed for connect_transit_gw_to_vgw") + err.Error())
	}
	attachEdgeVpnToTgw := url.Values{}
	attachEdgeVpnToTgw.Add("CID", c.CID)
	attachEdgeVpnToTgw.Add("action", "attach_edge_vpn_to_tgw")
	attachEdgeVpnToTgw.Add("tgw_name", awsTgwVpnConn.TgwName)
	attachEdgeVpnToTgw.Add("route_domain_name", awsTgwVpnConn.RouteDomainName)
	attachEdgeVpnToTgw.Add("connection_name", awsTgwVpnConn.ConnName)
	attachEdgeVpnToTgw.Add("public_ip", awsTgwVpnConn.PublicIP)
	attachEdgeVpnToTgw.Add("onprem_asn", awsTgwVpnConn.OnpremASN)
	attachEdgeVpnToTgw.Add("remote_cidr", awsTgwVpnConn.RemoteCIDR)

	Url.RawQuery = attachEdgeVpnToTgw.Encode()

	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return "", errors.New("HTTP Get attach_edge_vpn_to_tgw failed: " + err.Error())
	}

	var data AwsTgwVpnConnCreateResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", errors.New("Json Decode attach_edge_vpn_to_tgw failed: " + err.Error())
	}

	if !data.Return {
		return "", errors.New("Rest API attach_edge_vpn_to_tgw Get failed: " + data.Reason)
	}

	if !strings.Contains(data.Results, "vpn-") {
		return "", errors.New("cannot get vpn_id from result text")
	}

	tempStr := strings.Split(data.Results, "vpn-")[1]
	vpnID := "vpn-" + tempStr[0:17]

	return vpnID, nil
}

func (c *Client) GetAwsTgwVpnConn(awsTgwVpnConn *AwsTgwVpnConn) (*AwsTgwVpnConn, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for list_all_tgw_attachments") + err.Error())
	}
	listAllTgwAttachments := url.Values{}
	listAllTgwAttachments.Add("CID", c.CID)
	listAllTgwAttachments.Add("action", "list_all_tgw_attachments")
	listAllTgwAttachments.Add("tgw_name", awsTgwVpnConn.TgwName)
	listAllTgwAttachments.Add("resource_type", "vpn")

	Url.RawQuery = listAllTgwAttachments.Encode()

	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get list_all_tgw_attachments failed: " + err.Error())
	}

	var data AwsTgwVpnConnResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_all_tgw_attachments failed: " + err.Error())
	}

	if !data.Return {
		return nil, errors.New("Rest API list_all_tgw_attachments Get failed: " + data.Reason)
	}

	allAwsTgwVpnConn := data.Results
	for i := range allAwsTgwVpnConn {
		if allAwsTgwVpnConn[i].TgwName == awsTgwVpnConn.TgwName && allAwsTgwVpnConn[i].VpnID == awsTgwVpnConn.VpnID {
			awsTgwVpnConn.RouteDomainName = allAwsTgwVpnConn[i].RouteDomainName
			awsTgwVpnConn.ConnName = allAwsTgwVpnConn[i].ConnName
			awsTgwVpnConn.PublicIP = allAwsTgwVpnConn[i].PublicIP
			awsTgwVpnConn.OnpremASN = allAwsTgwVpnConn[i].OnpremASN
			awsTgwVpnConn.RemoteCIDR = strings.Join(allAwsTgwVpnConn[i].RemoteCIDR, ",")

			log.Printf("[DEBUG] Found AwsTgwVpnConn: %#v", awsTgwVpnConn)

			return awsTgwVpnConn, nil
		}
	}

	return nil, ErrNotFound
}

func (c *Client) DeleteAwsTgwVpnConn(awsTgwVpnConn *AwsTgwVpnConn) error {
	awsTgwVpnConn.CID = c.CID
	awsTgwVpnConn.Action = "detach_vpn_from_tgw"
	resp, err := c.Post(c.baseURL, awsTgwVpnConn)
	if err != nil {
		return errors.New("HTTP Post detach_vpn_from_tgw failed: " + err.Error())
	}

	var data APIResp

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode detach_vpn_from_tgw failed: " + err.Error())
	}

	if !data.Return {
		return errors.New("Rest API detach_vpn_from_tgw Post failed: " + data.Reason)
	}

	return nil
}
