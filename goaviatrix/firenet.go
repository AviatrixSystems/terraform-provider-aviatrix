package goaviatrix

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"
)

type FireNet struct {
	CID              string `form:"CID,omitempty"`
	Action           string `form:"action,omitempty"`
	VpcID            string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	GwName           string `form:"gw_name,omitempty" json:"gw_name,omitempty"`
	FirewallInstance []FirewallInstance
	FirewallEgress   bool `form:"firewall_egress,omitempty" json:"firewall_egress,omitempty"`
	Inspection       bool `form:"inspection,omitempty" json:"inspection,omitempty"`
}

type FireNetDetail struct {
	Region           string                 `json:"region,omitempty"`
	VpcID            string                 `json:"vpc_id,omitempty"`
	FirewallInstance []FirewallInstanceInfo `json:"firewall,omitempty"`
	Gateway          []GatewayInfo          `json:"gateway,omitempty"`
	FirewallEgress   string                 `json:"firewall_egress,omitempty"`
	Inspection       string                 `json:"inspection,omitempty"`
}

type GetFireNetResp struct {
	Return  bool          `json:"return"`
	Results FireNetDetail `json:"results"`
	Reason  string        `json:"reason"`
}

type ListFireNetResp struct {
	Return  bool              `json:"return"`
	Results FirewallInterface `json:"results"`
	Reason  string            `json:"reason"`
}

type FirewallInterface struct {
	Instances  []string            `json:"instances"`
	Interfaces map[string][]string `json:"interfaces"`
}

type GatewayInfo struct {
	DomainName string `json:"domain_name"`
	HaStatus   string `json:"ha_status"`
	GwName     string `json:"name"`
	TgwID      string `json:"tgw_id"`
}

type FirewallInstanceInfo struct {
	Enabled             bool   `json:"enabled"`
	GwName              string `json:"gateway"`
	InstanceID          string `json:"id"`
	FirewallName        string `json:"name"`
	LanInterface        string `json:"lan_interface_id,omitempty"`
	ManagementInterface string `json:"management_interface_id,omitempty"`
	EgressInterface     string `json:"egress_interface_id,omitempty"`
}

func (c *Client) CreateFireNet(fireNet *FireNet) error {
	return nil
}

func (c *Client) GetFireNet(fireNet *FireNet) (*FireNetDetail, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for show_firenet_detail: ") + err.Error())
	}
	showFireNetDetail := url.Values{}
	showFireNetDetail.Add("CID", c.CID)
	showFireNetDetail.Add("action", "show_firenet_detail")
	showFireNetDetail.Add("vpc_id", fireNet.VpcID)

	Url.RawQuery = showFireNetDetail.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get show_firenet_detail failed: " + err.Error())
	}
	var data GetFireNetResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.New("Json Decode show_firenet_detail failed: " + err.Error())
	}
	if !data.Return {
		return nil, errors.New("Rest API show_firenet_detail Get failed: " + data.Reason)
	}
	if strings.Split(data.Results.VpcID, "~~")[0] == fireNet.VpcID {
		return &data.Results, nil
	}
	return nil, ErrNotFound
}

func (c *Client) AssociateFirewallWithFireNet(firewallInstance *FirewallInstance) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for associate_firewall_with_firenet: ") + err.Error())
	}

	associateFirewallWithFireNet := url.Values{}
	associateFirewallWithFireNet.Add("CID", c.CID)
	associateFirewallWithFireNet.Add("action", "associate_firewall_with_firenet")
	associateFirewallWithFireNet.Add("vpc_id", firewallInstance.VpcID)
	associateFirewallWithFireNet.Add("gateway_name", firewallInstance.GwName)
	associateFirewallWithFireNet.Add("firewall_id", firewallInstance.InstanceID)
	associateFirewallWithFireNet.Add("firewall_name", firewallInstance.FirewallName)
	associateFirewallWithFireNet.Add("lan_interface", firewallInstance.LanInterface)
	associateFirewallWithFireNet.Add("management_interface", firewallInstance.ManagementInterface)
	associateFirewallWithFireNet.Add("egress_interface", firewallInstance.EgressInterface)

	Url.RawQuery = associateFirewallWithFireNet.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get associate_firewall_with_firenet failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode associate_firewall_with_firenet failed: " + err.Error())
	}
	if !data.Return {
		if strings.Contains(data.Reason, "already associated") {
			return nil
		}
		return errors.New("Rest API associate_firewall_with_firenet Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DisassociateFirewallFromFireNet(firewallInstance *FirewallInstance) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for disassociate_firewall_with_firenet: ") + err.Error())
	}
	disassociateFirewallWithFireNet := url.Values{}
	disassociateFirewallWithFireNet.Add("CID", c.CID)
	disassociateFirewallWithFireNet.Add("action", "disassociate_firewall_with_firenet")
	disassociateFirewallWithFireNet.Add("vpc_id", firewallInstance.VpcID)
	disassociateFirewallWithFireNet.Add("firewall_id", firewallInstance.InstanceID)

	Url.RawQuery = disassociateFirewallWithFireNet.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get disassociate_firewall_with_firenet failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode disassociate_firewall_with_firenet failed: " + err.Error())
	}
	if !data.Return {
		if strings.Contains(data.Reason, "not found") {
			return nil
		}
		return errors.New("Rest API disassociate_firewall_with_firenet Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) AttachFirewallToFireNet(firewallInstance *FirewallInstance) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for attach_firewall_to_firenet: ") + err.Error())
	}

	attachFirewallToFireNet := url.Values{}
	attachFirewallToFireNet.Add("CID", c.CID)
	attachFirewallToFireNet.Add("action", "attach_firewall_to_firenet")
	attachFirewallToFireNet.Add("vpc_id", firewallInstance.VpcID)
	attachFirewallToFireNet.Add("firewall_id", firewallInstance.InstanceID)

	Url.RawQuery = attachFirewallToFireNet.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get attach_firewall_to_firenet failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode attach_firewall_to_firenet failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API attach_firewall_to_firenet Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DetachFirewallFromFireNet(firewallInstance *FirewallInstance) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for detach_firewall_from_firenet: ") + err.Error())
	}

	detachFirewallFromFireNet := url.Values{}
	detachFirewallFromFireNet.Add("CID", c.CID)
	detachFirewallFromFireNet.Add("action", "detach_firewall_from_firenet")
	detachFirewallFromFireNet.Add("vpc_id", firewallInstance.VpcID)
	detachFirewallFromFireNet.Add("firewall_id", firewallInstance.InstanceID)

	Url.RawQuery = detachFirewallFromFireNet.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get detach_firewall_from_firenet failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode detach_firewall_from_firenet failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API detach_firewall_from_firenet Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) ConnectFireNetWithTgw(awsTgw *AWSTgw, vpcSolo VPCSolo, SecurityDomainName string) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for connect_firenet_with_tgw") + err.Error())
	}
	connectFireNetWithTgw := url.Values{}
	connectFireNetWithTgw.Add("CID", c.CID)
	connectFireNetWithTgw.Add("action", "connect_firenet_with_tgw")
	connectFireNetWithTgw.Add("vpc_id", vpcSolo.VpcID)
	connectFireNetWithTgw.Add("tgw_name", awsTgw.Name)
	connectFireNetWithTgw.Add("domain_name", SecurityDomainName)
	Url.RawQuery = connectFireNetWithTgw.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get connect_firenet_with_tgw failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode connect_firenet_with_tgw failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API connect_firenet_with_tgw Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DisconnectFireNetFromTgw(awsTgw *AWSTgw, vpcID string) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for disconnect_firenet_with_tgw") + err.Error())
	}
	disconnectFireNetWithTgw := url.Values{}
	disconnectFireNetWithTgw.Add("CID", c.CID)
	disconnectFireNetWithTgw.Add("action", "disconnect_firenet_with_tgw")
	disconnectFireNetWithTgw.Add("vpc_id", vpcID)
	Url.RawQuery = disconnectFireNetWithTgw.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get disconnect_firenet_with_tgw failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode disconnect_firenet_with_tgw failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API disconnect_firenet_with_tgw Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) EditFireNetInspection(fireNet *FireNet) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for edit_firenet") + err.Error())
	}
	editFireNet := url.Values{}
	editFireNet.Add("CID", c.CID)
	editFireNet.Add("action", "edit_firenet")
	editFireNet.Add("vpc_id", fireNet.VpcID)
	if fireNet.Inspection {
		editFireNet.Add("inspection", "true")
	} else {
		editFireNet.Add("inspection", "false")
	}

	Url.RawQuery = editFireNet.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get edit_firenet failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode edit_firenet failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API edit_firenet Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) EditFireNetEgress(fireNet *FireNet) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for edit_firenet") + err.Error())
	}
	editFireNet := url.Values{}
	editFireNet.Add("CID", c.CID)
	editFireNet.Add("action", "edit_firenet")
	editFireNet.Add("vpc_id", fireNet.VpcID)
	if fireNet.FirewallEgress {
		editFireNet.Add("firewall_egress", "true")
	} else {
		editFireNet.Add("firewall_egress", "false")
	}

	Url.RawQuery = editFireNet.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get edit_firenet failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode edit_firenet failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API edit_firenet Get failed: " + data.Reason)
	}
	return nil
}
