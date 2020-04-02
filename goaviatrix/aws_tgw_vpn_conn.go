package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

// VGWConn simple struct to hold VGW Connection details
type AwsTgwVpnConn struct {
	Action               string `form:"action,omitempty"`
	TgwName              string `form:"tgw_name,omitempty"`
	RouteDomainName      string `form:"route_domain_name,omitempty"`
	CID                  string `form:"CID,omitempty"`
	ConnName             string `form:"connection_name,omitempty"`
	PublicIP             string `form:"public_ip,omitempty"`
	OnpremASN            string `form:"onprem_asn,omitempty"`
	RemoteCIDR           string `form:"remote_cidr,omitempty"`
	VpnID                string `form:"vpn_id,omitempty"`
	InsideIpCIDRTun1     string `form:"inside_ip_cidr_tun_1,omitempty"`
	InsideIpCIDRTun2     string `form:"inside_ip_cidr_tun_2,omitempty"`
	PreSharedKeyTun1     string `form:"pre_shared_key_tun_1,omitempty"`
	PreSharedKeyTun2     string `form:"pre_shared_key_tun_2,omitempty"`
	LearnedCidrsApproval string `form:"learned_cidrs_approval,omitempty""`
}

type AwsTgwVpnConnEdit struct {
	TgwName              string   `json:"tgw_name,omitempty"`
	RouteDomainName      string   `json:"associated_route_domain_name,omitempty"`
	ConnName             string   `json:"vpc_name,omitempty"`
	PublicIP             string   `json:"public_ip,omitempty"`
	OnpremASN            string   `json:"aws_side_asn,omitempty"`
	RemoteCIDR           []string `json:"remote_cidrs,omitempty"`
	VpnID                string   `json:"vpc_id,omitempty"`
	InsideIpCIDRTun1     string   `json:"inside_ip_cidr_tun_1,omitempty"`
	InsideIpCIDRTun2     string   `json:"inside_ip_cidr_tun_2,omitempty"`
	PreSharedKeyTun1     string   `json:"pre_shared_key_tun_1,omitempty"`
	PreSharedKeyTun2     string   `json:"pre_shared_key_tun_2,omitempty"`
	LearnedCidrsApproval string   `json:"learned_cidrs_approval,omitempty"`
}

type AwsTgwVpnConnCreateResp struct {
	Return  bool    `json:"return"`
	Results vpnInfo `json:"results"`
	Reason  string  `json:"reason"`
}

type vpnInfo struct {
	Text  string `json:"text,omitempty"`
	VpnID string `json:"vpn_id,omitempty"`
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
	if awsTgwVpnConn.InsideIpCIDRTun1 != "" {
		attachEdgeVpnToTgw.Add("inside_ip_cidr_tun_1", awsTgwVpnConn.InsideIpCIDRTun1)
	}
	if awsTgwVpnConn.InsideIpCIDRTun2 != "" {
		attachEdgeVpnToTgw.Add("inside_ip_cidr_tun_2", awsTgwVpnConn.InsideIpCIDRTun2)
	}
	if awsTgwVpnConn.PreSharedKeyTun1 != "" {
		attachEdgeVpnToTgw.Add("pre_shared_key_tun_1", awsTgwVpnConn.PreSharedKeyTun1)
	}
	if awsTgwVpnConn.PreSharedKeyTun2 != "" {
		attachEdgeVpnToTgw.Add("pre_shared_key_tun_2", awsTgwVpnConn.PreSharedKeyTun2)
	}
	attachEdgeVpnToTgw.Add("learned_cidrs_approval", awsTgwVpnConn.LearnedCidrsApproval)

	Url.RawQuery = attachEdgeVpnToTgw.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return "", errors.New("HTTP Get attach_edge_vpn_to_tgw failed: " + err.Error())
	}

	var data AwsTgwVpnConnCreateResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return "", errors.New("Json Decode attach_edge_vpn_to_tgw failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return "", errors.New("Rest API attach_edge_vpn_to_tgw Get failed: " + data.Reason)
	}
	if data.Results.VpnID == "" {
		return "", errors.New("could not get vpn_id information")
	}
	return data.Results.VpnID, nil
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

	allAwsTgwVpnConn := data.Results
	for i := range allAwsTgwVpnConn {
		if allAwsTgwVpnConn[i].TgwName == awsTgwVpnConn.TgwName && allAwsTgwVpnConn[i].VpnID == awsTgwVpnConn.VpnID {
			awsTgwVpnConn.RouteDomainName = allAwsTgwVpnConn[i].RouteDomainName
			awsTgwVpnConn.ConnName = allAwsTgwVpnConn[i].ConnName
			awsTgwVpnConn.PublicIP = allAwsTgwVpnConn[i].PublicIP
			awsTgwVpnConn.OnpremASN = allAwsTgwVpnConn[i].OnpremASN
			awsTgwVpnConn.RemoteCIDR = strings.Join(allAwsTgwVpnConn[i].RemoteCIDR, ",")
			if allAwsTgwVpnConn[i].InsideIpCIDRTun1 != "" {
				awsTgwVpnConn.InsideIpCIDRTun1 = allAwsTgwVpnConn[i].InsideIpCIDRTun1
			}
			if allAwsTgwVpnConn[i].PreSharedKeyTun1 != "" {
				awsTgwVpnConn.PreSharedKeyTun1 = allAwsTgwVpnConn[i].PreSharedKeyTun1
			}
			if allAwsTgwVpnConn[i].InsideIpCIDRTun2 != "" {
				awsTgwVpnConn.InsideIpCIDRTun2 = allAwsTgwVpnConn[i].InsideIpCIDRTun2
			}
			if allAwsTgwVpnConn[i].PreSharedKeyTun2 != "" {
				awsTgwVpnConn.PreSharedKeyTun2 = allAwsTgwVpnConn[i].PreSharedKeyTun2
			}
			awsTgwVpnConn.LearnedCidrsApproval = allAwsTgwVpnConn[i].LearnedCidrsApproval

			log.Debugf("Found AwsTgwVpnConn: %#v", awsTgwVpnConn)

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
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode detach_vpn_from_tgw failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API detach_vpn_from_tgw Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) EnableVpnConnectionLearnedCidrsApproval(awsTgwVpnConn *AwsTgwVpnConn) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'enable_learned_cidrs_approval': ") + err.Error())
	}
	enableLearnedCidrsApproval := url.Values{}
	enableLearnedCidrsApproval.Add("CID", c.CID)
	enableLearnedCidrsApproval.Add("action", "enable_learned_cidrs_approval")
	enableLearnedCidrsApproval.Add("tgw_name", awsTgwVpnConn.TgwName)
	enableLearnedCidrsApproval.Add("attachment_name", awsTgwVpnConn.VpnID)
	enableLearnedCidrsApproval.Add("learned_cidrs_approval", awsTgwVpnConn.LearnedCidrsApproval)
	Url.RawQuery = enableLearnedCidrsApproval.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get 'enable_learned_cidrs_approval' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'enable_learned_cidrs_approval' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'enable_learned_cidrs_approval' Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DisableVpnConnectionLearnedCidrsApproval(awsTgwVpnConn *AwsTgwVpnConn) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'disable_learned_cidrs_approval': ") + err.Error())
	}
	enableLearnedCidrsApproval := url.Values{}
	enableLearnedCidrsApproval.Add("CID", c.CID)
	enableLearnedCidrsApproval.Add("action", "disable_learned_cidrs_approval")
	enableLearnedCidrsApproval.Add("tgw_name", awsTgwVpnConn.TgwName)
	enableLearnedCidrsApproval.Add("attachment_name", awsTgwVpnConn.VpnID)
	enableLearnedCidrsApproval.Add("learned_cidrs_approval", awsTgwVpnConn.LearnedCidrsApproval)
	Url.RawQuery = enableLearnedCidrsApproval.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get 'disable_learned_cidrs_approval' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'disable_learned_cidrs_approval' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'disable_learned_cidrs_approval' Get failed: " + data.Reason)
	}
	return nil
}
