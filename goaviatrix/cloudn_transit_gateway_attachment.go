package goaviatrix

import (
	"context"
	"fmt"
	"strings"
)

type CloudnTransitGatewayAttachment struct {
	DeviceName                       string `form:"device_name"`
	TransitGatewayName               string `form:"transit_gw" json:"gw_name"`
	ConnectionName                   string `form:"connection_name"`
	TransitGatewayBgpAsn             string `form:"bgp_local_as_number" json:"bgp_local_asn_number"`
	CloudnBgpAsn                     string `form:"external_device_as_number" json:"bgp_remote_asn_number"`
	CloudnLanInterfaceNeighborIP     string `form:"cloudn_neighbor_ip" json:"cloudn_neighbor_ip"`
	CloudnLanInterfaceNeighborBgpAsn string `form:"cloudn_neighbor_as_number" json:"cloudn_neighbor_as_number"`
	EnableOverPrivateNetwork         bool   `form:"direct_connect" json:"direct_connect_primary"`
	EnableJumboFrame                 bool   `json:"jumbo_frame"`
	EnableDeadPeerDetection          bool
	DpdConfig                        string `json:"dpd_config"`
	RoutingProtocol                  string `form:"routing_protocol"`
	Action                           string `form:"action"`
	CID                              string `form:"CID"`
}

func (c *Client) CreateCloudnTransitGatewayAttachment(ctx context.Context, attachment *CloudnTransitGatewayAttachment) error {
	attachment.Action = "attach_cloudwan_device_to_transit_gateway"
	attachment.CID = c.CID
	attachment.RoutingProtocol = "bgp"
	return c.PostAPIContext(ctx, attachment.Action, attachment, BasicCheck)
}

func (c *Client) GetCloudnTransitGatewayAttachment(ctx context.Context, connName string) (*CloudnTransitGatewayAttachment, error) {
	deviceName, err := c.GetDeviceName(connName)
	if err != nil {
		return nil, fmt.Errorf("could not get cloudn transit gateway attachment device name: %v", err)
	}

	vpcID, err := c.GetDeviceAttachmentVpcID(connName)
	if err != nil {
		return nil, fmt.Errorf("could not get cloudn transit gateway attachment VPC id: %v", err)
	}

	type site2cloudResp struct {
		Connections CloudnTransitGatewayAttachment
	}

	type resp struct {
		APIResp
		Results site2cloudResp
	}

	form := map[string]string{
		"action":    "get_site2cloud_conn_detail",
		"CID":       c.CID,
		"conn_name": connName,
		"vpc_id":    vpcID,
	}

	check := func(action, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "Connection does not exist") {
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s %s failed: %s", action, method, reason)
		}
		return nil
	}
	var data resp
	err = c.GetAPIContext(ctx, &data, form["action"], form, check)
	if err != nil {
		return nil, err
	}

	data.Results.Connections.ConnectionName = connName
	data.Results.Connections.DeviceName = deviceName
	data.Results.Connections.EnableDeadPeerDetection = data.Results.Connections.DpdConfig == "enable"
	return &data.Results.Connections, nil
}

func (c *Client) EnableJumboFrameOnConnectionToCloudn(ctx context.Context, connName, vpcID string) error {
	form := map[string]string{
		"action":          "enable_jumbo_frame_on_connection_to_cloudn",
		"CID":             c.CID,
		"connection_name": connName,
		"vpc_id":          vpcID,
	}
	return c.PostAPIContext(ctx, form["action"], form, BasicCheck)
}

func (c *Client) DisableJumboFrameOnConnectionToCloudn(ctx context.Context, connName, vpcID string) error {
	form := map[string]string{
		"action":          "disable_jumbo_frame_on_connection_to_cloudn",
		"CID":             c.CID,
		"connection_name": connName,
		"vpc_id":          vpcID,
	}
	return c.PostAPIContext(ctx, form["action"], form, BasicCheck)
}
