package goaviatrix

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
)

type CloudnEdgeGateway struct {
	Action                   string `form:"action,omitempty"`
	CID                      string `form:"CID,omitempty"`
	Type                     string `form:"type,omitempty"`
	Caag                     bool   `form:"caag,omitempty"`
	GatewayName              string `form:"gateway_name,omitempty" json:"name"`
	ManagementConnectionType string
	OverPrivateNetwork       bool   `form:"mgmt_over_private_network,omitempty" json:"mgmt_over_private_network"`
	WanInterfaceIp           string `form:"wan_ip,omitempty" json:"wan_ip"`
	WanDefaultGateway        string `form:"wan_default_gateway,omitempty" json:"wan_default_gateway"`
	LanInterfaceIp           string `form:"lan_ip,omitempty" json:"lan_ip"`
	ManagementInterfaceIp    string `form:"mgmt_ip,omitempty" json:"mgmt_ip"`
	DefaultGatewayIp         string `form:"mgmt_default_gateway,omitempty" json:"mgmt_default_gateway"`
	DnsServer                string `form:"dns_server_ip,omitempty" json:"dns_server_ip"`
	SecondaryDnsServer       string `form:"dns_server_ip_secondary,omitempty" json:"dns_server_ip_secondary"`
	Dhcp                     bool   `form:"dhcp,omitempty" json:"dhcp"`
	Hpe                      bool   `form:"hpe,omitempty"`
	ImageDownloadPath        string
}

func (c *Client) CreateCloundEdgeGateway(ctx context.Context, cloudnEdgeGateway *CloudnEdgeGateway) error {
	cloudnEdgeGateway.Action = "create_edge_gateway"
	cloudnEdgeGateway.CID = c.CID
	cloudnEdgeGateway.Type = "caag"
	cloudnEdgeGateway.Caag = true
	cloudnEdgeGateway.Hpe = false

	if cloudnEdgeGateway.ManagementConnectionType == "DHCP" {
		cloudnEdgeGateway.Dhcp = true
	}

	resp, err := c.PostAPIDownloadContext(ctx, cloudnEdgeGateway.Action, cloudnEdgeGateway, BasicCheck)
	if err != nil {
		return err
	}

	outFile, err := os.Create(cloudnEdgeGateway.ImageDownloadPath + "/" + cloudnEdgeGateway.GatewayName + ".iso")
	if err != nil {
		return err
	}

	_, err = io.Copy(outFile, resp)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetCloudnEdgeGateway(ctx context.Context, gwName string) (*CloudnEdgeGateway, error) {
	form := map[string]string{
		"action":      "get_cloudwan_device_details",
		"CID":         c.CID,
		"device_name": gwName,
	}

	type Resp struct {
		Return  bool              `json:"return"`
		Results CloudnEdgeGateway `json:"results"`
		Reason  string            `json:"reason"`
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

func (c *Client) DeleteCloudnEdgeGateway(ctx context.Context, gwName string) error {
	form := map[string]string{
		"action": "delete_edge_gateway",
		"CID":    c.CID,
		"name":   gwName,
		"caag":   "true",
	}

	return c.PostAPIContext(ctx, form["action"], form, BasicCheck)
}
