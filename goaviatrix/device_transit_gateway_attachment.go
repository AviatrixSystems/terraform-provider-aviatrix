package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
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
	ManualBGPCidrs          []string
	Action                  string `form:"action"`
	CID                     string `form:"CID"`
}

func (c *Client) CreateDeviceTransitGatewayAttachment(attachment *DeviceTransitGatewayAttachment) error {
	attachment.Action = "attach_cloudwan_device_to_transit_gateway"
	attachment.CID = c.CID
	resp, err := c.Post(c.baseURL, attachment)
	if err != nil {
		return errors.New("HTTP Post attach_cloudwan_device_to_transit_gateway failed: " + err.Error())
	}

	type Resp struct {
		Return  bool   `json:"return,omitempty"`
		Results string `json:"results,omitempty"`
		Reason  string `json:"reason,omitempty"`
	}
	var data Resp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return errors.New("Reading response body attach_cloudwan_device_to_transit_gateway failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return errors.New("Json Decode attach_cloudwan_device_to_transit_gateway failed: " + err.Error() + "\n Body: " + b.String())
	}
	if !data.Return {
		return errors.New("Rest API attach_cloudwan_device_to_transit_gateway Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetDeviceTransitGatewayAttachment(attachment *DeviceTransitGatewayAttachment) (*DeviceTransitGatewayAttachment, error) {
	deviceName, err := c.GetDeviceName(attachment.ConnectionName)
	if err != nil {
		return nil, fmt.Errorf("could not get device name: %v", err)
	}

	vpcID, err := c.GetDeviceAttachmentVpcID(attachment.ConnectionName)
	if err != nil {
		return nil, fmt.Errorf("could not get device attachment VPC id: %v", err)
	}

	resp, err := c.Post(c.baseURL, struct {
		CID            string `form:"CID"`
		Action         string `form:"action"`
		ConnectionName string `form:"conn_name"`
		VpcID          string `form:"vpc_id"`
	}{
		CID:            c.CID,
		Action:         "get_site2cloud_conn_detail",
		ConnectionName: attachment.ConnectionName,
		VpcID:          vpcID,
	})
	if err != nil {
		return nil, errors.New("HTTP POST get_site2cloud_conn_detail failed: " + err.Error())
	}

	var data Site2CloudConnDetailResp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return nil, errors.New("Reading response body get_site2cloud_conn_detail failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return nil, errors.New("Json Decode get_site2cloud_conn_detail failed: " + err.Error() +
			"\n Body: " + b.String())
	}

	if !data.Return {
		if strings.Contains(data.Reason, "does not exist") {
			return nil, ErrNotFound
		}
		return nil, errors.New("Rest API get_site2cloud_conn_detail Post failed: " + data.Reason)
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
	}, nil
}

func (c *Client) GetDeviceAttachmentVpcID(connectionName string) (string, error) {
	resp, err := c.Post(c.baseURL, struct {
		CID    string `form:"CID"`
		Action string `form:"action"`
	}{
		CID:    c.CID,
		Action: "list_cloudwan_attachments",
	})
	if err != nil {
		return "", errors.New("HTTP POST list_cloudwan_attachments failed: " + err.Error())
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
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return "", errors.New("Reading response body list_cloudwan_attachments failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return "", errors.New("Json Decode list_cloudwan_attachments failed: " + err.Error() +
			"\n Body: " + b.String())
	}
	if !data.Return {
		return "", errors.New("Rest API list_cloudwan_attachments Post failed: " + data.Reason)
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

	resp, err := c.Post(c.baseURL, struct {
		CID            string `form:"CID"`
		Action         string `form:"action"`
		VpcID          string `form:"vpc_id"`
		ConnectionName string `form:"connection_name"`
	}{
		CID:            c.CID,
		Action:         "detach_cloudwan_device",
		VpcID:          vpcID,
		ConnectionName: connectionName,
	})
	if err != nil {
		return errors.New("HTTP POST detach_cloudwan_device failed: " + err.Error())
	}

	type Resp struct {
		Return  bool   `json:"return"`
		Results string `json:"results"`
		Reason  string `json:"reason"`
	}
	var data Resp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return errors.New("Reading response body detach_cloudwan_device failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return errors.New("Json Decode detach_cloudwan_device failed: " + err.Error() +
			"\n Body: " + b.String())
	}
	if !data.Return {
		return errors.New("Rest API detach_cloudwan_device Post failed: " + data.Reason)
	}

	return nil
}
