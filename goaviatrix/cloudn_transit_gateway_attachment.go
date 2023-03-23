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
	CloudnLanInterfaceNeighborIP     string `json:"cloudn_neighbor_ip"`
	CloudnLanInterfaceNeighborBgpAsn string `json:"cloudn_neighbor_as_number"`
	CloudnNeighbor                   string `form:"cloudn_neighbor"`
	EnableOverPrivateNetwork         bool   `form:"direct_connect" json:"direct_connect_primary"`
	EnableJumboFrame                 bool   `json:"jumbo_frame"`
	EnableDeadPeerDetection          bool
	DpdConfig                        string   `json:"dpd_config"`
	RoutingProtocol                  string   `form:"routing_protocol"`
	Action                           string   `form:"action"`
	CID                              string   `form:"CID"`
	EnableLearnedCidrsApproval       bool     `form:"connection_learned_cidrs_approval"`
	EnableLearnedCidrsApprovalValue  string   `json:"conn_learned_cidrs_approval"`
	ApprovedCidrs                    []string `json:"conn_approved_learned_cidrs"`
	PrependAsPath                    string   `json:"conn_bgp_prepend_as_path"`
	Async                            bool     `form:"async,omitempty"`
}

func (c *Client) CreateCloudnTransitGatewayAttachment(ctx context.Context, attachment *CloudnTransitGatewayAttachment) error {
	attachment.Action = "attach_cloudwan_device_to_transit_gateway"
	attachment.CID = c.CID
	attachment.RoutingProtocol = "bgp"
	attachment.CloudnNeighbor = `"[{"ip_addr": ` + attachment.CloudnLanInterfaceNeighborIP + `, "as_num": ` + attachment.CloudnLanInterfaceNeighborBgpAsn + `}]"`
	attachment.Async = true
	return c.PostAsyncAPIContext(ctx, attachment.Action, attachment, BasicCheck)
}

func (c *Client) GetCloudnTransitGatewayAttachment(ctx context.Context, connName string) (*CloudnTransitGatewayAttachment, error) {
	deviceName, err := c.GetDeviceName(connName)
	if err != nil {
		if err == ErrNotFound {
			return nil, err
		}
		return nil, fmt.Errorf("could not get cloudn transit gateway attachment device name: %w", err)
	}

	vpcID, err := c.GetDeviceAttachmentVpcID(connName)
	if err != nil {
		if err == ErrNotFound {
			return nil, err
		}
		return nil, fmt.Errorf("could not get cloudn transit gateway attachment VPC id: %w", err)
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
	data.Results.Connections.EnableLearnedCidrsApproval = data.Results.Connections.EnableLearnedCidrsApprovalValue == "yes"

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

func (c *Client) EditCloudnTransitGatewayAttachmentASPathPrepend(ctx context.Context, attachment *CloudnTransitGatewayAttachment, prependASPath []string) error {
	action := "edit_transit_connection_as_path_prepend"
	return c.PostAPIContext(ctx, action, struct {
		CID            string `form:"CID"`
		Action         string `form:"action"`
		GatewayName    string `form:"gateway_name"`
		ConnectionName string `form:"connection_name"`
		PrependASPath  string `form:"connection_as_path_prepend"`
	}{
		CID:            c.CID,
		Action:         action,
		GatewayName:    attachment.TransitGatewayName,
		ConnectionName: attachment.ConnectionName,
		PrependASPath:  strings.Join(prependASPath, ","),
	}, BasicCheck)
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
		"async":           "true",
	}

	return c.PostAsyncAPI(form["action"], form, BasicCheck)
}
