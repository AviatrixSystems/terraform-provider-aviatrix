package goaviatrix

import (
	"context"
)

const (
	FirewallManagerConfigModeDefault = "DEFAULT"
	FirewallManagerConfigModeAdvance = "ADVANCE"
)

type VendorInfo struct {
	CID            string `form:"CID,omitempty"`
	Action         string `form:"action,omitempty"`
	VpcID          string `form:"vpc_id,omitempty"`
	InstanceID     string `form:"firewall_id,omitempty"`
	FirewallName   string `form:"firewall_name,omitempty"`
	VendorType     string `form:"firewall_vendor,omitempty"`
	Username       string `form:"user,omitempty"`
	Password       string `form:"password,omitempty"`
	ApiToken       string `form:"api_token,omitempty"`
	PrivateKeyFile string `form:"private_key_file,omitempty"`
	RouteTable     string `form:"route_table,omitempty"`
	PublicIP       string `form:"public_ip,omitempty"`
	Save           bool
	Synchronize    bool `form:"sync,omitempty"`
}
type FirewallTemplateConfig struct {
	Template      string `json:"template,omitempty"`
	TemplateStack string `json:"template_stack,omitempty"`
	RouteTable    string `json:"route_table,omitempty"`
}
type FirewallManager struct {
	CID                    string                            `json:"CID,omitempty"`
	Action                 string                            `json:"action,omitempty"`
	VpcID                  string                            `json:"vpc_id,omitempty"`
	GatewayName            string                            `json:"gw_name,omitempty"`
	VendorType             string                            `json:"firewall_vendor,omitempty"`
	PublicIP               string                            `json:"public_ip,omitempty"`
	Username               string                            `json:"user,omitempty"`
	Password               string                            `json:"password,omitempty"`
	Template               string                            `json:"template,omitempty"`
	TemplateStack          string                            `json:"template_stack,omitempty"`
	RouteTable             string                            `json:"route_table,omitempty"`
	FirewallTemplateConfig map[string]FirewallTemplateConfig `json:"firewall_template_config,omitempty"`
	ConfigMode             string                            `json:"config_mode,omitempty"`
	Save                   bool                              `json:"save,omitempty"`
	Synchronize            bool                              `json:"sync,omitempty"`
}

func (c *Client) EditFireNetFirewallVendorInfo(vendorInfo *VendorInfo) error {
	form := map[string]string{
		"CID":             c.CID,
		"action":          "edit_firenet_firewall_vendor_info",
		"vpc_id":          vendorInfo.VpcID,
		"firewall_id":     vendorInfo.InstanceID,
		"firewall_name":   vendorInfo.FirewallName,
		"firewall_vendor": vendorInfo.VendorType,
		"user":            vendorInfo.Username,
		"password":        vendorInfo.Password,
		"api_token":       vendorInfo.ApiToken,
		"route_table":     vendorInfo.RouteTable,
		"public_ip":       vendorInfo.PublicIP,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) ShowFireNetFirewallVendorConfig(vendorInfo *VendorInfo) error {
	form := map[string]string{
		"CID":         c.CID,
		"action":      "show_firenet_firewall_vendor_config",
		"vpc_id":      vendorInfo.VpcID,
		"firewall_id": vendorInfo.InstanceID,
		"sync":        "true",
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) EditFireNetFirewallManagerVendorInfo(ctx context.Context, firewallManager *FirewallManager) error {
	firewallManager.CID = c.CID
	firewallManager.Action = "edit_firenet_firewall_manager_vendor_info"

	return c.PostAPIContext2(ctx, nil, firewallManager.Action, firewallManager, BasicCheck)
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

func (c *Client) EditFireNetFirewallVendorInfoWithPrivateKey(vendorInfo *VendorInfo) error {
	params := map[string]string{
		"action":          "edit_firenet_firewall_vendor_info",
		"CID":             c.CID,
		"vpc_id":          vendorInfo.VpcID,
		"firewall_id":     vendorInfo.InstanceID,
		"firewall_name":   vendorInfo.FirewallName,
		"public_ip":       vendorInfo.PublicIP,
		"firewall_vendor": vendorInfo.VendorType,
		"user":            vendorInfo.Username,
	}

	var files []File

	key := File{
		ParamName:      "private_key_file",
		UseFileContent: true,
		FileName:       "key.pem", // fake name for key
		FileContent:    vendorInfo.PrivateKeyFile,
	}
	files = append(files, key)

	return c.PostFileAPI(params, files, BasicCheck)
}
