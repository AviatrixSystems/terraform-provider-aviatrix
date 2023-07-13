package goaviatrix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CloudnNeighborAsNumberWrapper string

func (w *CloudnNeighborAsNumberWrapper) UnmarshalJSON(data []byte) (err error) {
	if port, err := strconv.Atoi(string(data)); err == nil {
		str := strconv.Itoa(port)
		*w = CloudnNeighborAsNumberWrapper(str)
		return nil
	}
	var str string
	err = myUnmarshal(data, &str)
	if err != nil {
		return err
	}
	return myUnmarshal([]byte(str), w)
}

type CloudnTransitGatewayAttachment struct {
	DeviceName                       string `form:"device_name"`
	TransitGatewayName               string `form:"transit_gw" json:"gw_name"`
	ConnectionName                   string `form:"connection_name"`
	TransitGatewayBgpAsn             string `form:"bgp_local_as_number" json:"bgp_local_asn_number"`
	CloudnBgpAsn                     string `form:"external_device_as_number" json:"bgp_remote_asn_number"`
	CloudnLanInterfaceNeighborIP     string `json:"cloudn_neighbor_ip"`
	CloudnLanInterfaceNeighborBgpAsn int    `json:"cloudn_neighbor_as_number"`
	CloudnNeighbor                   string `form:"cloudn_neighbor"`
	EnableOverPrivateNetwork         bool   `form:"direct_connect" json:"direct_connect_primary"`
	EnableJumboFrame                 bool   `json:"jumbo_frame"`
	EnableDeadPeerDetection          bool
	DpdConfig                        string   `json:"dpd_config"`
	RoutingProtocol                  string   `form:"routing_protocol"`
	Action                           string   `form:"action"`
	CID                              string   `form:"CID"`
	EnableLearnedCidrsApproval       string   `form:"conn_learned_cidrs_approval" json:"conn_learned_cidrs_approval"`
	ApprovedCidrs                    []string `json:"conn_approved_learned_cidrs"`
	PrependAsPath                    string   `json:"conn_bgp_prepend_as_path"`
	Async                            bool     `form:"async,omitempty"`
}

type CloudnTransitGatewayAttachmentResp struct {
	DeviceName                       string
	TransitGatewayName               string
	ConnectionName                   string
	TransitGatewayBgpAsn             string                        `json:"bgp_local_asn_number"`
	CloudnBgpAsn                     string                        `json:"bgp_remote_asn_number"`
	CloudnLanInterfaceNeighborIP     string                        `json:"cloudn_neighbor_ip"`
	CloudnLanInterfaceNeighborBgpAsn CloudnNeighborAsNumberWrapper `json:"cloudn_neighbor_as_number"`
	EnableOverPrivateNetwork         bool                          `json:"direct_connect_primary"`
	EnableJumboFrame                 bool                          `json:"jumbo_frame"`
	EnableDeadPeerDetection          bool
	DpdConfig                        string   `json:"dpd_config"`
	EnableLearnedCidrsApproval       string   `json:"conn_learned_cidrs_approval"`
	ApprovedCidrs                    []string `json:"conn_approved_learned_cidrs"`
	PrependAsPath                    string   `json:"conn_bgp_prepend_as_path"`
}

func (c *Client) CreateCloudnTransitGatewayAttachment(ctx context.Context, attachment *CloudnTransitGatewayAttachment) error {
	attachment.Action = "attach_cloudwan_device_to_transit_gateway"
	attachment.CID = c.CID
	attachment.RoutingProtocol = "bgp"
	attachment.Async = true

	type CloudnNeighbor struct {
		IpAddr string `json:"ip_addr"`
		AsNum  int    `json:"as_num"`
	}

	cloudnNeighbor := CloudnNeighbor{
		IpAddr: attachment.CloudnLanInterfaceNeighborIP,
		AsNum:  attachment.CloudnLanInterfaceNeighborBgpAsn,
	}

	cloudnNeighborJson, _ := json.Marshal(cloudnNeighbor)
	attachment.CloudnNeighbor = "[" + string(cloudnNeighborJson) + "]"

	return c.PostAsyncAPIContext(ctx, attachment.Action, attachment, BasicCheck)
}

func (c *Client) GetCloudnTransitGatewayAttachment(ctx context.Context, connName string) (*CloudnTransitGatewayAttachmentResp, error) {
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
		Connections CloudnTransitGatewayAttachmentResp
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
	err = c.GetAPIContextCloudnTransitGatewayAttachment(ctx, &data, form["action"], form, check)
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

func (c *Client) GetAPIContextCloudnTransitGatewayAttachment(ctx context.Context, v interface{}, action string, d map[string]string, checkFunc CheckAPIResponseFunc) error {
	Url, err := c.urlEncode(d)
	if err != nil {
		return fmt.Errorf("could not url encode values for action %q: %v", action, err)
	}

	try, maxTries, backoff := 0, 5, 500*time.Millisecond
	var resp *http.Response
	for {
		try++
		resp, err = c.GetContext(ctx, Url, nil)
		if err == nil {
			break
		}

		if try == maxTries {
			return fmt.Errorf("HTTP Get %s failed: %v", action, err)
		}
		time.Sleep(backoff)
		// Double the backoff time after each failed try
		backoff *= 2
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	var data APIResp
	if err := json.NewDecoder(strings.NewReader(bodyString)).Decode(&data); err != nil {
		return fmt.Errorf("Json Decode into standard format failed: %v\n Body: %s", err, bodyString)
	}
	if err := checkFunc(action, "Get", data.Reason, data.Return); err != nil {
		return err
	}

	err = myUnmarshal(buf.Bytes(), &v)
	if err != nil {
		return fmt.Errorf("json unmarshal failed: %v\n Body: %s", err, buf.String())
	}

	return nil
}
