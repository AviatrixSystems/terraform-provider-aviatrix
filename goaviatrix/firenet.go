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
}

type FireNetDetail struct {
	Region           string                 `json:"region,omitempty"`
	VpcID            string                 `json:"vpc_id,omitempty"`
	FirewallInstance []FirewallInstanceInfo `json:"firewall,omitempty"`
	Gateway          []GatewayInfo          `json:"gateway,omitempty"`
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
	Enabled      bool   `json:"enabled"`
	GwName       string `json:"gateway"`
	InstanceID   string `json:"id"`
	FirewallName string `json:"name"`
}

func (c *Client) CreateFireNet(fireNet *FireNet) error {
	return nil
}

func (c *Client) GetFireNet(fireNet *FireNet) (*FireNetDetail, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for show_firenet_detail: ") + err.Error())
	}
	getFireNetDetail := url.Values{}
	getFireNetDetail.Add("CID", c.CID)
	getFireNetDetail.Add("action", "show_firenet_detail")
	getFireNetDetail.Add("vpc_id", fireNet.VpcID)

	Url.RawQuery = getFireNetDetail.Encode()
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

	interfaces, err := c.GetInterfaceInformation(firewallInstance)
	if err != nil || len(interfaces) != 3 {
		return errors.New(("failed to read interface information: ") + err.Error())
	}

	associateFirewallWithFireNet := url.Values{}
	associateFirewallWithFireNet.Add("CID", c.CID)
	associateFirewallWithFireNet.Add("action", "associate_firewall_with_firenet")
	associateFirewallWithFireNet.Add("vpc_id", firewallInstance.VpcID)
	associateFirewallWithFireNet.Add("gateway_name", firewallInstance.GwName)
	associateFirewallWithFireNet.Add("firewall_id", firewallInstance.InstanceID)
	associateFirewallWithFireNet.Add("firewall_name", firewallInstance.FirewallName)
	associateFirewallWithFireNet.Add("lan_interface", interfaces[0])
	associateFirewallWithFireNet.Add("egress_interface", interfaces[1])
	associateFirewallWithFireNet.Add("management_interface", interfaces[2])

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

func (c *Client) GetInterfaceInformation(firewallInstance *FirewallInstance) ([]string, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for list_firenet: ") + err.Error())
	}
	listFireNet := url.Values{}
	listFireNet.Add("CID", c.CID)
	listFireNet.Add("action", "list_firenet")
	listFireNet.Add("subaction", "instance")
	listFireNet.Add("vpc_id", firewallInstance.VpcID)
	listFireNet.Add("gateway_name", firewallInstance.GwName)

	Url.RawQuery = listFireNet.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get list_firenet failed: " + err.Error())
	}
	var data ListFireNetResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_firenet failed: " + err.Error())
	}
	if !data.Return {
		return nil, errors.New("Rest API list_firenet Get failed: " + data.Reason)
	}
	for _, instance := range data.Results.Instances {
		if strings.Contains(instance, firewallInstance.InstanceID) {
			var temp [3]string
			var result []string
			for _, interfaceStr := range data.Results.Interfaces[instance] {
				if strings.Contains(interfaceStr, "lan") {
					temp[0] = interfaceStr
				} else if strings.Contains(interfaceStr, "egress") {
					temp[1] = interfaceStr
				} else if strings.Contains(interfaceStr, "management") {
					temp[2] = interfaceStr
				}
			}
			for i := 0; i < 3; i++ {
				if temp[i] != "" {
					result = append(result, temp[i])
				}
			}
			return result, nil
		}
	}
	return nil, ErrNotFound
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
