package goaviatrix

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
)

type EdgeCaag struct {
	Action                      string `form:"action,omitempty"`
	CID                         string `form:"CID,omitempty"`
	Type                        string `form:"type,omitempty"`
	Caag                        bool   `form:"caag,omitempty"`
	Name                        string `form:"gateway_name,omitempty" json:"name"`
	ManagementInterfaceConfig   string
	ManagementEgressIpPrefix    string `form:"mgmt_egress_ip" json:"mgmt_egress_ip"`
	EnableOverPrivateNetwork    bool   `form:"mgmt_over_private_network,omitempty" json:"mgmt_over_private_network"`
	WanInterfaceIpPrefix        string `form:"wan_ip,omitempty" json:"wan_ip"`
	WanDefaultGatewayIp         string `form:"wan_default_gateway,omitempty" json:"wan_default_gateway"`
	LanInterfaceIpPrefix        string `form:"lan_ip,omitempty" json:"lan_ip"`
	ManagementInterfaceIpPrefix string `form:"mgmt_ip,omitempty" json:"mgmt_ip"`
	ManagementDefaultGatewayIp  string `form:"mgmt_default_gateway,omitempty" json:"mgmt_default_gateway"`
	DnsServerIp                 string `form:"dns_server_ip,omitempty" json:"dns_server_ip"`
	SecondaryDnsServerIp        string `form:"dns_server_ip_secondary,omitempty" json:"dns_server_ip_secondary"`
	Dhcp                        bool   `form:"dhcp,omitempty" json:"dhcp"`
	Hpe                         bool   `form:"hpe,omitempty"`
	ZtpFileType                 string `form:"ztp_file_type,omitempty"`
	ZtpFileDownloadPath         string
	State                       string `json:"state"`
}

func (c *Client) CreateEdgeCaag(ctx context.Context, edgeCaag *EdgeCaag) error {
	edgeCaag.Action = "create_edge_gateway"
	edgeCaag.CID = c.CID
	edgeCaag.Type = "caag"
	edgeCaag.Caag = true
	edgeCaag.Hpe = false

	if edgeCaag.ManagementInterfaceConfig == "DHCP" {
		edgeCaag.Dhcp = true
	}

	resp, err := c.PostAPIDownloadContext(ctx, edgeCaag.Action, edgeCaag, BasicCheck)
	if err != nil {
		return err
	}

	var fileName string
	if edgeCaag.ZtpFileType == "iso" {
		fileName = edgeCaag.ZtpFileDownloadPath + "/" + edgeCaag.Name + ".iso"
	} else {
		fileName = edgeCaag.ZtpFileDownloadPath + "/" + edgeCaag.Name + "-cloud-init.txt"
	}

	outFile, err := os.Create(fileName)
	if err != nil {
		return err
	}

	_, err = io.Copy(outFile, resp)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetEdgeCaag(ctx context.Context, name string) (*EdgeCaag, error) {
	form := map[string]string{
		"action":      "get_cloudwan_device_details",
		"CID":         c.CID,
		"device_name": name,
	}

	type Resp struct {
		Return  bool     `json:"return"`
		Results EdgeCaag `json:"results"`
		Reason  string   `json:"reason"`
	}

	var data Resp
	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "does not exist") {
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	err := c.GetAPIContext(ctx, &data, form["action"], form, checkFunc)
	if err != nil {
		return nil, err
	}

	return &data.Results, nil
}

func (c *Client) UpdateEdgeCaag(ctx context.Context, edgeCaag *EdgeCaag) error {
	form := map[string]string{
		"action":         "update_edge_gateway",
		"CID":            c.CID,
		"gateway_name":   edgeCaag.Name,
		"mgmt_egress_ip": edgeCaag.ManagementEgressIpPrefix,
	}

	return c.PostAPIContext(ctx, form["action"], form, BasicCheck)
}

func (c *Client) DeleteEdgeCaag(ctx context.Context, name string, state string) error {
	form := map[string]string{
		"CID": c.CID,
	}

	if state == "check" || state == "waiting" {
		form["action"] = "reset_managed_cloudn_to_factory_state"
		form["device_name"] = name
	} else {
		form["action"] = "delete_edge_gateway"
		form["name"] = name
		form["caag"] = "true"
	}

	return c.PostAPIContext(ctx, form["action"], form, BasicCheck)
}
