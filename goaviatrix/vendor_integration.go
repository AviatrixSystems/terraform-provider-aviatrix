package goaviatrix

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
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
	ApiToken     string `form:"api_token,omitempty"`
	RouteTable   string `form:"route_table,omitempty"`
	PublicIP     string `form:"public_ip,omitempty"`
	Save         bool
	Synchronize  bool `form:"sync,omitempty"`
}

type FirewallManager struct {
	CID           string
	Action        string
	VpcID         string
	GatewayName   string
	VendorType    string
	PublicIP      string
	Username      string
	Password      string
	Template      string
	TemplateStack string
	RouteTable    string
	Save          bool
	Synchronize   bool
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
	editFireNetFirewallVendorInfo.Add("api_token", vendorInfo.ApiToken)
	editFireNetFirewallVendorInfo.Add("route_table", vendorInfo.RouteTable)
	editFireNetFirewallVendorInfo.Add("public_ip", vendorInfo.PublicIP)
	Url.RawQuery = editFireNetFirewallVendorInfo.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get edit_firenet_firewall_vendor_info failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode edit_firenet_firewall_vendor_info failed: " + err.Error() + "\n Body: " + bodyString)
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
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode show_firenet_firewall_vendor_config failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API show_firenet_firewall_vendor_config Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) EditFireNetFirewallManagerVendorInfo(ctx context.Context, firewallManager *FirewallManager) error {
	params := map[string]string{
		"action":          "edit_firenet_firewall_manager_vendor_info",
		"CID":             c.CID,
		"vpc_id":          firewallManager.VpcID,
		"gw_name":         firewallManager.GatewayName,
		"firewall_vendor": firewallManager.VendorType,
		"public_ip":       firewallManager.PublicIP,
		"user":            firewallManager.Username,
		"password":        firewallManager.Password,
		"template":        firewallManager.Template,
		"template_stack":  firewallManager.TemplateStack,
		"route_table":     firewallManager.RouteTable,
	}

	return c.PostAPIContext(ctx, params["action"], params, BasicCheck)
}

func (c *Client) SyncFireNetFirewallManagerVendorConfig(ctx context.Context, firewallManager *FirewallManager) error {
	params := map[string]string{
		"action":  "show_firenet_firewall_vendor_config",
		"CID":     c.CID,
		"vpc_id":  firewallManager.VpcID,
		"gw_name": firewallManager.GatewayName,
		"sync":    "true",
	}

	return c.PostAPIContext(ctx, params["action"], params, BasicCheck)
}
