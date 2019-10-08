package goaviatrix

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
)

type FirewallInstance struct {
	CID                 string `form:"CID,omitempty"`
	Action              string `form:"action,omitempty"`
	VpcID               string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	GwName              string `form:"gw_name,omitempty" json:"gw_name,omitempty"`
	FirewallName        string `form:"firewall_name,omitempty" json:"firewall_name,omitempty"`
	FirewallImage       string `form:"firewall_image,omitempty" json:"firewall_image,omitempty"`
	FirewallSize        string `form:"firewall_size,omitempty" json:"firewall_size,omitempty"`
	EgressSubnet        string `form:"egress_subnet,omitempty" json:"egress_subnet,omitempty"`
	ManagementSubnet    string `form:"management_subnet,omitempty" json:"management_subnet,omitempty"`
	KeyName             string `form:"key_name,omitempty" json:"key_name,omitempty"`
	IamRole             string `form:"iam_role,omitempty" json:"iam_role,omitempty"`
	BootstrapBucketName string `form:"bootstrap_bucket_name,omitempty" json:"bootstrap_bucket_name,omitempty"`
	InstanceID          string `form:"firewall_id,omitempty" json:"instance_id,omitempty"`
	Attached            bool
	LanInterface        string `form:"lan_interface,omitempty" json:"lan_interface_id,omitempty"`
	ManagementInterface string `form:"management_interface,omitempty" json:"management_interface_id,omitempty"`
	EgressInterface     string `form:"egress_interface,omitempty" json:"egress_interface_id,omitempty"`
	ManagementPublicIP  string `json:"management_public_ip,omitempty"`
}

type FirewallInstanceResp struct {
	Return  bool             `json:"return"`
	Results FirewallInstance `json:"results"`
	Reason  string           `json:"reason"`
}

type FirewallInstanceCreateResp struct {
	Return  bool                         `json:"return"`
	Results FirewallInstanceCreateResult `json:"results"`
	Reason  string                       `json:"reason"`
}

type FirewallInstanceCreateResult struct {
	Text       string `json:"text"`
	FirewallID string `json:"firewall_id"`
}

func (c *Client) CreateFirewallInstance(firewallInstance *FirewallInstance) (string, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return "", errors.New(("url Parsing failed for add_firewall_instance: ") + err.Error())
	}
	addFirewallInstance := url.Values{}
	addFirewallInstance.Add("CID", c.CID)
	addFirewallInstance.Add("action", "add_firewall_instance")
	addFirewallInstance.Add("gw_name", firewallInstance.GwName)
	addFirewallInstance.Add("firewall_name", firewallInstance.FirewallName)
	addFirewallInstance.Add("firewall_image", firewallInstance.FirewallImage)
	addFirewallInstance.Add("firewall_size", firewallInstance.FirewallSize)
	addFirewallInstance.Add("egress_subnet", firewallInstance.EgressSubnet)
	addFirewallInstance.Add("management_subnet", firewallInstance.ManagementSubnet)
	addFirewallInstance.Add("key_name", firewallInstance.InstanceID)
	addFirewallInstance.Add("iam_role", firewallInstance.InstanceID)
	addFirewallInstance.Add("bootstrap_bucket_name", firewallInstance.InstanceID)
	addFirewallInstance.Add("no_associate", strconv.FormatBool(true))

	Url.RawQuery = addFirewallInstance.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return "", errors.New("HTTP Get add_firewall_instance failed: " + err.Error())
	}

	var data FirewallInstanceCreateResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", errors.New("Json Decode add_firewall_instance failed: " + err.Error())
	}
	if !data.Return {
		return "", errors.New("Rest API add_firewall_instance Get failed: " + data.Reason)
	}
	if data.Results.FirewallID != "" {
		return data.Results.FirewallID, nil
	}
	return "", ErrNotFound
}

func (c *Client) GetFirewallInstance(firewallInstance *FirewallInstance) (*FirewallInstance, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for get_instance_by_id: ") + err.Error())
	}
	getInstanceById := url.Values{}
	getInstanceById.Add("CID", c.CID)
	getInstanceById.Add("action", "get_instance_by_id")
	getInstanceById.Add("instance_id", firewallInstance.InstanceID)

	Url.RawQuery = getInstanceById.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get get_instance_by_id failed: " + err.Error())
	}

	var data FirewallInstanceResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.New("Json Decode get_instance_by_id failed: " + err.Error())
	}
	if !data.Return {
		return nil, errors.New("Rest API get_instance_by_id Get failed: " + data.Reason)
	}
	if data.Results.InstanceID == firewallInstance.InstanceID {
		return &data.Results, nil
	}
	return nil, ErrNotFound
}

func (c *Client) DeleteFirewallInstance(firewallInstance *FirewallInstance) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for delete_firenet_firewall_instance ") + err.Error())
	}
	deleteFirenetFirewallInstance := url.Values{}
	deleteFirenetFirewallInstance.Add("CID", c.CID)
	deleteFirenetFirewallInstance.Add("action", "delete_firenet_firewall_instance")
	deleteFirenetFirewallInstance.Add("vpc_id", firewallInstance.VpcID)
	deleteFirenetFirewallInstance.Add("firewall_id", firewallInstance.InstanceID)

	Url.RawQuery = deleteFirenetFirewallInstance.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get delete_firenet_firewall_instance failed: " + err.Error())
	}

	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode delete_firenet_firewall_instance failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API delete_firenet_firewall_instance Get failed: " + data.Reason)
	}
	return nil
}
