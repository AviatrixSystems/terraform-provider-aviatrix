package goaviatrix

import (
	"encoding/json"
	"errors"
	"net/url"
)

type VendorInfo struct {
	CID          string `form:"CID,omitempty"`
	Action       string `form:"action,omitempty"`
	VpcID        string `form:"vpc_id,omitempty"`
	InstanceID   string `form:"firewall_id,omitempty"`
	FirewallName string `form:"firewall_name,omitempty"`
	VendorType   string `form:"firewall_vendor,omitempty"`
	Username     string `form:"user,omitempty"`
	Password     string `form:"password,omitempty"`
	RouteTable   string `form:"route_table,omitempty"`
	PublicIP     string `form:"public_ip,omitempty"`
	Save         bool
	Synchronize  bool `form:"sync,omitempty"`
}

func (c *Client) EditFireNetFirewallVendorInfo(vendorInfo *VendorInfo) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for edit_firenet_firewall_vendor_info") + err.Error())
	}

	editFireNetFirewallVendorInfo := url.Values{}
	editFireNetFirewallVendorInfo.Add("CID", c.CID)
	editFireNetFirewallVendorInfo.Add("action", "edit_firenet_firewall_vendor_info")
	editFireNetFirewallVendorInfo.Add("vpc_id", vendorInfo.VpcID)
	editFireNetFirewallVendorInfo.Add("firewall_id", vendorInfo.InstanceID)
	editFireNetFirewallVendorInfo.Add("firewall_name", vendorInfo.FirewallName)
	editFireNetFirewallVendorInfo.Add("firewall_vendor", vendorInfo.VendorType)
	editFireNetFirewallVendorInfo.Add("user", vendorInfo.Username)
	editFireNetFirewallVendorInfo.Add("password", vendorInfo.Password)
	editFireNetFirewallVendorInfo.Add("route_table", vendorInfo.RouteTable)
	editFireNetFirewallVendorInfo.Add("public_ip", vendorInfo.PublicIP)

	Url.RawQuery = editFireNetFirewallVendorInfo.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get edit_firenet_firewall_vendor_info failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode edit_firenet_firewall_vendor_info failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API edit_firenet_firewall_vendor_info Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) ShowFireNetFirewallVendorConfig(vendorInfo *VendorInfo) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for show_firenet_firewall_vendor_config") + err.Error())
	}

	showFireNetFirewallVendorConfig := url.Values{}
	showFireNetFirewallVendorConfig.Add("CID", c.CID)
	showFireNetFirewallVendorConfig.Add("action", "show_firenet_firewall_vendor_config")
	showFireNetFirewallVendorConfig.Add("vpc_id", vendorInfo.VpcID)
	showFireNetFirewallVendorConfig.Add("firewall_id", vendorInfo.InstanceID)
	showFireNetFirewallVendorConfig.Add("sync", "true")

	Url.RawQuery = showFireNetFirewallVendorConfig.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get show_firenet_firewall_vendor_config failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode show_firenet_firewall_vendor_config failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API show_firenet_firewall_vendor_config Get failed: " + data.Reason)
	}
	return nil
}
