package goaviatrix

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"
)

type FirewallInstance struct {
	CID                 string `form:"CID,omitempty"`
	Action              string `form:"action,omitempty"`
	VpcID               string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	GwName              string `form:"gw_name,omitempty" json:"gw_name,omitempty"`
	FirewallName        string `form:"firewall_name,omitempty" json:"firewall_name,omitempty"`
	FirewallImage       string `form:"firewall_image,omitempty" json:"firewall_image,omitempty"`
	EgressSubnet        string `form:"egress_subnet,omitempty" json:"egress_subnet,omitempty"`
	ManagementSubnet    string `form:"management_subnet,omitempty" json:"management_subnet,omitempty"`
	KeyName             string `form:"key_name,omitempty" json:"key_name,omitempty"`
	IamRole             string `form:"iam_role,omitempty" json:"iam_role,omitempty"`
	BootstrapBucketName string `form:"bootstrap_bucket_name,omitempty" json:"bootstrap_bucket_name,omitempty"`
	InstanceID          string `form:"firewall_id,omitempty" json:"instance_id,omitempty"`
	Enabled             bool   `form:"enabled,omitempty" json:"enabled,omitempty"`
}

type FirewallInstanceResp struct {
	Return  bool             `json:"return"`
	Results FirewallInstance `json:"results"`
	Reason  string           `json:"reason"`
}

type FirewallInstanceCreateResp struct {
	Return  bool   `json:"return"`
	Results string `json:"results"`
	Reason  string `json:"reason"`
}

func (c *Client) CreateFirewallInstance(firewallInstance *FirewallInstance) (string, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return "", errors.New(("url Parsing failed for add_firewall_instance: ") + err.Error())
	}
	getInstanceById := url.Values{}
	getInstanceById.Add("CID", c.CID)
	getInstanceById.Add("action", "add_firewall_instance")
	getInstanceById.Add("gw_name", firewallInstance.GwName)
	getInstanceById.Add("firewall_name", firewallInstance.FirewallName)
	getInstanceById.Add("firewall_image", firewallInstance.FirewallImage)
	getInstanceById.Add("egress_subnet", firewallInstance.EgressSubnet)
	getInstanceById.Add("management_subnet", firewallInstance.ManagementSubnet)
	getInstanceById.Add("key_name", firewallInstance.InstanceID)
	getInstanceById.Add("iam_role", firewallInstance.InstanceID)
	getInstanceById.Add("bootstrap_bucket_name", firewallInstance.InstanceID)

	Url.RawQuery = getInstanceById.Encode()
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

	index := strings.Index(data.Results, "-")
	if index != -1 {
		return data.Results[index-1 : index+18], nil
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
