package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
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
	LearnedCidrsApproval string `form:"learned_cidrs_approval,omitempty"`
	EnableAcceleration   string `form:"enable_global_acceleration"`
}

type AwsTgwVpnConnEdit struct {
	TgwName              string          `json:"tgw_name,omitempty"`
	RouteDomainName      string          `json:"associated_route_domain_name,omitempty"`
	ConnName             string          `json:"vpc_name,omitempty"`
	PublicIP             string          `json:"public_ip,omitempty"`
	OnpremASNRaw         json.RawMessage `json:"aws_side_asn,omitempty"`
	OnpremASN            string
	RemoteCIDR           []string              `json:"remote_cidrs,omitempty"`
	VpnID                string                `json:"vpc_id,omitempty"`
	InsideIpCIDRTun1     string                `json:"inside_ip_cidr_tun_1,omitempty"`
	InsideIpCIDRTun2     string                `json:"inside_ip_cidr_tun_2,omitempty"`
	PreSharedKeyTun1     string                `json:"pre_shared_key_tun_1,omitempty"`
	PreSharedKeyTun2     string                `json:"pre_shared_key_tun_2,omitempty"`
	LearnedCidrsApproval string                `json:"learned_cidrs_approval,omitempty"`
	EnableAcceleration   bool                  `json:"enable_acceleration"`
	VpnTunnelData        map[string]TunnelData `json:"vpn_tunnel_data"`
}

type TunnelData struct {
	Status               string `json:"status"`
	RouteCount           int    `json:"route_count"`
	VpnOutsideAddress    string `json:"vpn_outside_address"`
	VpnInsideAddress     string `json:"vpn_inside_address"`
	TgwAsn               string `json:"tgw_asn"`
	StatusMessage        string `json:"status_message"`
	LastStatusChangeTime string `json:"last_status_change"`
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
	var data AwsTgwVpnConnCreateResp
	form := map[string]string{
		"CID":                        c.CID,
		"action":                     "attach_edge_vpn_to_tgw",
		"tgw_name":                   awsTgwVpnConn.TgwName,
		"route_domain_name":          awsTgwVpnConn.RouteDomainName,
		"connection_name":            awsTgwVpnConn.ConnName,
		"public_ip":                  awsTgwVpnConn.PublicIP,
		"onprem_asn":                 awsTgwVpnConn.OnpremASN,
		"remote_cidr":                awsTgwVpnConn.RemoteCIDR,
		"learned_cidrs_approval":     awsTgwVpnConn.LearnedCidrsApproval,
		"enable_global_acceleration": awsTgwVpnConn.EnableAcceleration,
	}
	if awsTgwVpnConn.InsideIpCIDRTun1 != "" {
		form["inside_ip_cidr_tun_1"] = awsTgwVpnConn.InsideIpCIDRTun1
	}
	if awsTgwVpnConn.InsideIpCIDRTun2 != "" {
		form["inside_ip_cidr_tun_2"] = awsTgwVpnConn.InsideIpCIDRTun2
	}
	if awsTgwVpnConn.PreSharedKeyTun1 != "" {
		form["pre_shared_key_tun_1"] = awsTgwVpnConn.PreSharedKeyTun1
	}
	if awsTgwVpnConn.PreSharedKeyTun2 != "" {
		form["pre_shared_key_tun_2"] = awsTgwVpnConn.PreSharedKeyTun2
	}
	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return "", err
	}
	if data.Results.VpnID == "" {
		return "", errors.New("could not get vpn_id information")
	}
	return data.Results.VpnID, nil
}

func (c *Client) GetAwsTgwVpnConn(awsTgwVpnConn *AwsTgwVpnConn) (*AwsTgwVpnConn, error) {
	var data AwsTgwVpnConnResp
	form := map[string]string{
		"CID":           c.CID,
		"action":        "list_all_tgw_attachments",
		"tgw_name":      awsTgwVpnConn.TgwName,
		"resource_type": "vpn",
	}
	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	allAwsTgwVpnConn := data.Results
	for i := range allAwsTgwVpnConn {
		if allAwsTgwVpnConn[i].TgwName == awsTgwVpnConn.TgwName && allAwsTgwVpnConn[i].VpnID == awsTgwVpnConn.VpnID {
			awsTgwVpnConn.RouteDomainName = allAwsTgwVpnConn[i].RouteDomainName
			awsTgwVpnConn.ConnName = allAwsTgwVpnConn[i].ConnName
			awsTgwVpnConn.PublicIP = allAwsTgwVpnConn[i].PublicIP
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
			awsTgwVpnConn.EnableAcceleration = "no"
			if allAwsTgwVpnConn[i].EnableAcceleration {
				awsTgwVpnConn.EnableAcceleration = "yes"
			}

			// aws_side_asn can return as either string or int from API
			var asnString string
			if len(allAwsTgwVpnConn[i].OnpremASNRaw) != 0 {
				// First try as string
				err = json.Unmarshal(allAwsTgwVpnConn[i].OnpremASNRaw, &asnString)
				if err != nil {
					// String failed, must be int
					var asnInt int
					err = json.Unmarshal(allAwsTgwVpnConn[i].OnpremASNRaw, &asnInt)
					if err != nil {
						return nil, fmt.Errorf("json decode list_all_tgw_attachments aws_side_asn field failed: aws_side_asn = %s: %v", string(allAwsTgwVpnConn[i].OnpremASNRaw), err)
					}
					asnString = strconv.Itoa(asnInt)
				}
			}
			awsTgwVpnConn.OnpremASN = asnString

			log.Debugf("Found AwsTgwVpnConn: %#v", awsTgwVpnConn)

			return awsTgwVpnConn, nil
		}
	}

	return nil, ErrNotFound
}

func (c *Client) DeleteAwsTgwVpnConn(awsTgwVpnConn *AwsTgwVpnConn) error {
	awsTgwVpnConn.CID = c.CID
	awsTgwVpnConn.Action = "detach_vpn_from_tgw"
	return c.PostAPI(awsTgwVpnConn.Action, awsTgwVpnConn, BasicCheck)
}

func (c *Client) EnableVpnConnectionLearnedCidrsApproval(awsTgwVpnConn *AwsTgwVpnConn) error {
	form := map[string]string{
		"CID":                    c.CID,
		"action":                 "enable_learned_cidrs_approval",
		"tgw_name":               awsTgwVpnConn.TgwName,
		"attachment_name":        awsTgwVpnConn.VpnID,
		"learned_cidrs_approval": awsTgwVpnConn.LearnedCidrsApproval,
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) DisableVpnConnectionLearnedCidrsApproval(awsTgwVpnConn *AwsTgwVpnConn) error {
	form := map[string]string{
		"CID":                    c.CID,
		"action":                 "disable_learned_cidrs_approval",
		"tgw_name":               awsTgwVpnConn.TgwName,
		"attachment_name":        awsTgwVpnConn.VpnID,
		"learned_cidrs_approval": awsTgwVpnConn.LearnedCidrsApproval,
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) GetAwsTgwVpnTunnelData(awsTgwVpnConn *AwsTgwVpnConn) (*AwsTgwVpnConnEdit, error) {
	params := map[string]string{
		"action":          "list_attachment_route_table_details",
		"CID":             c.CID,
		"tgw_name":        awsTgwVpnConn.TgwName,
		"attachment_name": awsTgwVpnConn.VpnID,
	}

	type Resp struct {
		Return  bool              `json:"return"`
		Results AwsTgwVpnConnEdit `json:"results"`
		Reason  string            `json:"reason"`
	}

	var data Resp

	err := c.GetAPI(&data, params["action"], params, BasicCheck)
	if err != nil {
		return nil, err
	}

	return &data.Results, nil
}
