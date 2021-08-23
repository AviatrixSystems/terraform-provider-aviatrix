package goaviatrix

import (
	"fmt"
	"strconv"
	"strings"
)

type DeviceTransitGatewayAttachment struct {
	DeviceName              string `form:"device_name,omitempty"`
	TransitGatewayName      string `form:"transit_gw,omitempty"`
	ConnectionName          string `form:"connection_name,omitempty"`
	RoutingProtocol         string `form:"routing_protocol,omitempty"`
	TransitGatewayBgpAsn    string `form:"bgp_local_as_number,omitempty"`
	DeviceBgpAsn            string `form:"external_device_as_number,omitempty"`
	Phase1Authentication    string `form:"phase1_authentication,omitempty"`
	Phase1DHGroups          string `form:"phase1_dh_groups,omitempty"`
	Phase1Encryption        string `form:"phase1_encryption,omitempty"`
	Phase2Authentication    string `form:"phase2_authentication,omitempty"`
	Phase2DHGroups          string `form:"phase2_dh_groups,omitempty"`
	Phase2Encryption        string `form:"phase2_encryption,omitempty"`
	EnableGlobalAccelerator string `form:"enable_global_accelerator,omitempty"`
	PreSharedKey            string `form:"pre_shared_key,omitempty"`
	LocalTunnelIP           string `form:"local_tunnel_ip,omitempty"`
	RemoteTunnelIP          string `form:"remote_tunnel_ip,omitempty"`
	VpcID                   string
	EventTriggeredHA        bool
	ManualBGPCidrs          []string
	Action                  string `form:"action"`
	CID                     string `form:"CID"`
}

func (c *Client) CreateDeviceTransitGatewayAttachment(attachment *DeviceTransitGatewayAttachment) error {
	attachment.Action = "attach_cloudwan_device_to_transit_gateway"
	attachment.CID = c.CID
	return c.PostAPI(attachment.Action, attachment, BasicCheck)
}

func (c *Client) GetDeviceTransitGatewayAttachment(attachment *DeviceTransitGatewayAttachment) (*DeviceTransitGatewayAttachment, error) {
	deviceName, err := c.GetDeviceName(attachment.ConnectionName)
	if err != nil {
		if err == ErrNotFound {
			return nil, err
		}
		return nil, fmt.Errorf("could not get device name: %v", err)
	}

	vpcID, err := c.GetDeviceAttachmentVpcID(attachment.ConnectionName)
	if err != nil {
		if err == ErrNotFound {
			return nil, err
		}
		return nil, fmt.Errorf("could not get device attachment VPC id: %v", err)
	}

	form := map[string]string{
		"CID":       c.CID,
		"action":    "get_site2cloud_conn_detail",
		"vpc_id":    vpcID,
		"conn_name": attachment.ConnectionName,
	}

	var data Site2CloudConnDetailResp

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "does not exist") {
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	err = c.GetAPI(&data, form["action"], form, checkFunc)
	if err != nil {
		return nil, err
	}

	return &DeviceTransitGatewayAttachment{
		DeviceName:              deviceName,
		TransitGatewayName:      data.Results.Connections.GwName,
		ConnectionName:          attachment.ConnectionName,
		TransitGatewayBgpAsn:    data.Results.Connections.BgpLocalASN,
		DeviceBgpAsn:            data.Results.Connections.BgpRemoteASN,
		Phase1Authentication:    data.Results.Connections.Algorithm.Phase1Auth[0],
		Phase1DHGroups:          data.Results.Connections.Algorithm.Phase1DhGroups[0],
		Phase1Encryption:        data.Results.Connections.Algorithm.Phase1Encrption[0],
		Phase2Authentication:    data.Results.Connections.Algorithm.Phase2Auth[0],
		Phase2DHGroups:          data.Results.Connections.Algorithm.Phase2DhGroups[0],
		Phase2Encryption:        data.Results.Connections.Algorithm.Phase2Encrption[0],
		EnableGlobalAccelerator: strconv.FormatBool(data.Results.Connections.EnableGlobalAccelerator),
		LocalTunnelIP:           data.Results.Connections.BgpLocalIP,
		RemoteTunnelIP:          data.Results.Connections.BgpRemoteIP,
		ManualBGPCidrs:          data.Results.Connections.ManualBGPCidrs,
		EventTriggeredHA:        data.Results.Connections.EventTriggeredHA == "enabled",
		VpcID:                   vpcID,
	}, nil
}

func (c *Client) GetDeviceAttachmentVpcID(connectionName string) (string, error) {
	form := map[string]string{
		"CID":    c.CID,
		"action": "list_cloudwan_attachments",
	}

	type CloudWanAttachments struct {
		VpcID string `json:"vpc_id"`
		Name  string `json:"name"`
	}

	type Resp struct {
		Return  bool                  `json:"return"`
		Results []CloudWanAttachments `json:"results"`
		Reason  string                `json:"reason"`
	}

	var data Resp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return "", err
	}

	for _, attachment := range data.Results {
		if attachment.Name == connectionName {
			return attachment.VpcID, nil
		}
	}

	return "", ErrNotFound
}

func (c *Client) DeleteDeviceAttachment(connectionName string) error {
	vpcID, err := c.GetDeviceAttachmentVpcID(connectionName)
	if err != nil {
		return fmt.Errorf("could not get device attachment VPC id: %v", err)
	}

	form := map[string]string{
		"CID":             c.CID,
		"action":          "detach_cloudwan_device",
		"vpc_id":          vpcID,
		"connection_name": connectionName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}
