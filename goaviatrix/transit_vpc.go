package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Gateway simple struct to hold gateway details
type TransitVpc struct {
	AccountName                  string `form:"account_name,omitempty" json:"account_name,omitempty"`
	Action                       string `form:"action,omitempty"`
	CID                          string `form:"CID,omitempty"`
	CloudType                    int    `form:"cloud_type,omitempty" json:"cloud_type,omitempty"`
	DnsServer                    string `form:"dns_server,omitempty" json:"dns_server,omitempty"`
	GwName                       string `form:"gw_name,omitempty" json:"vpc_name,omitempty"`
	GwSize                       string `form:"gw_size,omitempty"`
	VpcID                        string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	VNetNameResourceGroup        string `form:"vnet_and_resource_group_names,omitempty" json:"vpc_id,omitempty"`
	Subnet                       string `form:"public_subnet,omitempty" json:"vpc_net,omitempty"`
	HASubnet                     string `form:"ha_subnet,omitempty"`
	HAZone                       string `form:"new_zone,omitempty"`
	HASubnetGCP                  string `form:"new_subnet,omitempty"`
	PeeringHASubnet              string `json:"public_subnet,omitempty"`
	VpcRegion                    string `form:"region,omitempty" json:"vpc_region,omitempty"`
	VpcSize                      string `form:"gw_size,omitempty" json:"gw_size,omitempty"`
	EnableNAT                    string `form:"nat_enabled,omitempty" json:"enable_nat,omitempty"`
	SingleAzHa                   string `form:"single_az_ha,omitempty"`
	EnableVpcDnsServer           string `json:"use_vpc_dns,omitempty"`
	TagList                      string `form:"tags,omitempty"`
	EnableHybridConnection       bool   `form:"enable_hybrid_connection" json:"tgw_enabled,omitempty"`
	ConnectedTransit             string `form:"connected_transit" json:"connected_transit,omitempty"`
	InsaneMode                   string `form:"insane_mode,omitempty"`
	ReuseEip                     string `form:"reuse_eip,omitempty"`
	AllocateNewEipRead           bool   `json:"newly_allocated_eip,omitempty"`
	Eip                          string `form:"eip,omitempty"`
	EnableActiveMesh             string `form:"enable_activemesh,omitempty" json:"enable_activemesh,omitempty"`
	Zone                         string `form:"zone,omitempty" json:"zone,omitempty"`
	EnableAdvertiseTransitCidr   bool
	BgpManualSpokeAdvertiseCidrs string `form:"bgp_manual_spoke,omitempty"`
	EnableTransitFireNet         string `form:"enable_transit_firenet,omitempty"`
	LanVpcID                     string `form:"lan_vpc_id,omitempty"`
	LanPrivateSubnet             string `form:"lan_private_subnet,omitempty"`
	LearnedCidrsApproval         string `form:"learned_cidrs_approval,omitempty"`
	EncVolume                    string `form:"enc_volume,omitempty"`
	BgpOverLan                   string `form:"bgp_over_lan,omitempty"`
	EnablePrivateOob             string `form:"private_oob,omitempty"`
	OobManagementSubnet          string `form:"oob_mgmt_subnet,omitempty"`
	HAOobManagementSubnet        string
}

type TransitGatewayAdvancedConfig struct {
	BgpPollingTime                    string
	PrependASPath                     []string
	LocalASNumber                     string
	BgpEcmpEnabled                    bool
	ActiveStandbyEnabled              bool
	ActiveStandbyConnections          []StandbyConnection
	LearnedCIDRsApprovalMode          string
	ConnectionLearnedCIDRApprovalInfo []LearnedCIDRApprovalInfo
	TunnelAddrLocal                   string
	TunnelAddrLocalBackup             string
	PeerVnetId                        []string
	BgpHoldTime                       int
}

type StandbyConnection struct {
	ConnectionName    string
	ActiveGatewayType string
}

type TransitGatewayAdvancedConfigRespResult struct {
	BgpPollingTime                    int                       `json:"bgp_polling_time"`
	PrependASPath                     string                    `json:"bgp_prepend_as_path"`
	LocalASNumber                     string                    `json:"local_asn_num"`
	BgpEcmpEnabled                    string                    `json:"bgp_ecmp"`
	ActiveStandby                     string                    `json:"active-standby"`
	ActiveStandbyStatus               map[string]string         `json:"active_standby_status"`
	LearnedCIDRsApprovalMode          string                    `json:"learned_cidrs_approval_mode"`
	ConnectionLearnedCIDRApprovalInfo []LearnedCIDRApprovalInfo `json:"connection_learned_cidrs_approval_info"`
	TunnelAddrLocal                   string                    `json:"tunnel_addr_local"`
	TunnelAddrLocalBackup             string                    `json:"tunnel_addr_local_backup"`
	PeerVnetId                        []string                  `json:"peer_vnet_id"`
	BgpHoldTime                       int                       `json:"bgp_hold_time"`
}

type LearnedCIDRApprovalInfo struct {
	ConnName        string `json:"conn_name"`
	EnabledApproval string `json:"conn_learned_cidrs_approval"`
}

type TransitGatewayAdvancedConfigResp struct {
	Return  bool                                   `json:"return"`
	Results TransitGatewayAdvancedConfigRespResult `json:"results"`
	Reason  string                                 `json:"reason"`
}

type TransitGwFireNetInterfaces struct {
	VpcName                  string `json:"vpc_name"`
	VpcRegion                string `json:"vpc_region"`
	TransitVpc               string `json:"transit_vpc"`
	FireNetInterfacesEnabled bool   `json:"dmz_enabled"`
	Name                     string `json:"name"`
}

type TransitGwFireNetInterfacesResp struct {
	Return  bool                       `json:"return"`
	Results TransitGwFireNetInterfaces `json:"results"`
	Reason  string                     `json:"reason"`
}

func (c *Client) LaunchTransitVpc(gateway *TransitVpc) error {
	gateway.CID = c.CID
	gateway.Action = "create_transit_gw"
	resp, err := c.Post(c.baseURL, gateway)
	if err != nil {
		return errors.New("HTTP Post create_transit_gw failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode create_transit_gw failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API create_transit_gw Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) EnableHaTransitGateway(gateway *TransitVpc) error {
	gateway.CID = c.CID
	gateway.Action = "create_peering_ha_gateway"
	resp, err := c.Post(c.baseURL, gateway)
	if err != nil {
		return errors.New("HTTP Post create_peering_ha_gateway failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode create_peering_ha_gateway failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API create_peering_ha_gateway Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) EnableHaTransitVpc(gateway *TransitVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New("url Parsing failed for enable_transit_ha " + err.Error())
	}
	enableTransitHa := url.Values{}
	enableTransitHa.Add("CID", c.CID)
	enableTransitHa.Add("action", "enable_transit_ha")
	enableTransitHa.Add("gw_name", gateway.GwName)
	enableTransitHa.Add("eip", gateway.Eip)

	if gateway.CloudType == GCP {
		enableTransitHa.Add("new_zone", gateway.HAZone)
	} else {
		enableTransitHa.Add("public_subnet", gateway.HASubnet)
		enableTransitHa.Add("oob_mgmt_subnet", gateway.HAOobManagementSubnet)
	}

	Url.RawQuery = enableTransitHa.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get enable_transit_ha failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode enable_transit_ha failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		if strings.Contains(data.Reason, "HA GW already exists") {
			log.Infof("HA is already enabled %s", data.Reason)
			return nil
		}
		log.Errorf("Enabling HA failed with error %s", data.Reason)

		return errors.New("Rest API enable_transit_ha Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) AttachTransitGWForHybrid(gateway *TransitVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for enable_transit_gateway_interface_to_aws_tgw ") + err.Error())
	}
	enableTransitGatewayInterfaceToAwsTgw := url.Values{}
	enableTransitGatewayInterfaceToAwsTgw.Add("CID", c.CID)
	enableTransitGatewayInterfaceToAwsTgw.Add("action", "enable_transit_gateway_interface_to_aws_tgw")
	enableTransitGatewayInterfaceToAwsTgw.Add("gateway_name", gateway.GwName)
	Url.RawQuery = enableTransitGatewayInterfaceToAwsTgw.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get enable_transit_gateway_interface_to_aws_tgw failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		if strings.Contains(err.Error(), "already enabled tgw interface") {
			return nil
		}
		return errors.New("Json Decode enable_transit_gateway_interface_to_aws_tgw failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API enable_transit_gateway_interface_to_aws_tgw Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DetachTransitGWForHybrid(gateway *TransitVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for disable_transit_gateway_interface_to_aws_tgw") + err.Error())
	}
	disableTransitGatewayInterfaceToAwsTgw := url.Values{}
	disableTransitGatewayInterfaceToAwsTgw.Add("CID", c.CID)
	disableTransitGatewayInterfaceToAwsTgw.Add("action", "disable_transit_gateway_interface_to_aws_tgw")
	disableTransitGatewayInterfaceToAwsTgw.Add("gateway_name", gateway.GwName)
	Url.RawQuery = disableTransitGatewayInterfaceToAwsTgw.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get disable_transit_gateway_interface_to_aws_tgw failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode disable_transit_gateway_interface_to_aws_tgw failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API disable_transit_gateway_interface_to_aws_tgw Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) EnableConnectedTransit(gateway *TransitVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for enable_connected_transit_on_gateway") + err.Error())
	}
	enableConnectedTransitOnGateway := url.Values{}
	enableConnectedTransitOnGateway.Add("CID", c.CID)
	enableConnectedTransitOnGateway.Add("action", "enable_connected_transit_on_gateway")
	enableConnectedTransitOnGateway.Add("gateway_name", gateway.GwName)
	Url.RawQuery = enableConnectedTransitOnGateway.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get enable_connected_transit_on_gateway failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode enable_connected_transit_on_gateway failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API enable_connected_transit_on_gateway Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DisableConnectedTransit(gateway *TransitVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for disable_connected_transit_on_gateway") + err.Error())
	}
	disableConnectedTransitOnGateway := url.Values{}
	disableConnectedTransitOnGateway.Add("CID", c.CID)
	disableConnectedTransitOnGateway.Add("action", "disable_connected_transit_on_gateway")
	disableConnectedTransitOnGateway.Add("gateway_name", gateway.GwName)
	Url.RawQuery = disableConnectedTransitOnGateway.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get disable_connected_transit_on_gateway failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode disable_connected_transit_on_gateway failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API disable_connected_transit_on_gateway Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) EnableGatewayFireNetInterfaces(gateway *TransitVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for enable_gateway_firenet_interfaces") + err.Error())
	}
	disableConnectedTransitOnGateway := url.Values{}
	disableConnectedTransitOnGateway.Add("CID", c.CID)
	disableConnectedTransitOnGateway.Add("action", "enable_gateway_firenet_interfaces")
	disableConnectedTransitOnGateway.Add("gateway_name", gateway.GwName)
	Url.RawQuery = disableConnectedTransitOnGateway.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Post enable_gateway_firenet_interfaces failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode enable_gateway_firenet_interfaces failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API enable_gateway_firenet_interfaces Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DisableGatewayFireNetInterfaces(gateway *TransitVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for disable_gateway_firenet_interfaces") + err.Error())
	}
	disableConnectedTransitOnGateway := url.Values{}
	disableConnectedTransitOnGateway.Add("CID", c.CID)
	disableConnectedTransitOnGateway.Add("action", "disable_gateway_firenet_interfaces")
	disableConnectedTransitOnGateway.Add("gateway", gateway.GwName)
	Url.RawQuery = disableConnectedTransitOnGateway.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Post disable_gateway_firenet_interfaces failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode disable_gateway_firenet_interfaces failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API disable_gateway_firenet_interfaces Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) EnableGatewayFireNetInterfacesWithGWLB(gateway *TransitVpc) error {
	data := map[string]string{
		"action":       "enable_gateway_firenet_interfaces",
		"CID":          c.CID,
		"gateway_name": gateway.GwName,
		"mode":         "gwlb",
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) EnableAdvertiseTransitCidr(transitGw *TransitVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for enable_advertise_transit_cidr") + err.Error())
	}
	enableAdvertiseTransitCidr := url.Values{}
	enableAdvertiseTransitCidr.Add("CID", c.CID)
	enableAdvertiseTransitCidr.Add("action", "enable_advertise_transit_cidr")
	enableAdvertiseTransitCidr.Add("gateway_name", transitGw.GwName)
	if transitGw.EnableAdvertiseTransitCidr {
		enableAdvertiseTransitCidr.Add("advertise_transit_cidr", "yes")
	} else {
		enableAdvertiseTransitCidr.Add("advertise_transit_cidr", "no")
	}

	Url.RawQuery = enableAdvertiseTransitCidr.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get enable_advertise_transit_cidr failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode enable_advertise_transit_cidr failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API enable_advertise_transit_cidr Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DisableAdvertiseTransitCidr(transitGw *TransitVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for disable_advertise_transit_cidr") + err.Error())
	}
	enableAdvertiseTransitCidr := url.Values{}
	enableAdvertiseTransitCidr.Add("CID", c.CID)
	enableAdvertiseTransitCidr.Add("action", "disable_advertise_transit_cidr")
	enableAdvertiseTransitCidr.Add("gateway_name", transitGw.GwName)
	enableAdvertiseTransitCidr.Add("advertise_transit_cidr", "no")

	Url.RawQuery = enableAdvertiseTransitCidr.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get disable_advertise_transit_cidr failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode disable_advertise_transit_cidr failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API disable_advertise_transit_cidr Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) SetBgpManualSpokeAdvertisedNetworks(transitGw *TransitVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for edit_aviatrix_transit_advanced_config") + err.Error())
	}
	editTransitAdvancedConfig := url.Values{}
	editTransitAdvancedConfig.Add("CID", c.CID)
	editTransitAdvancedConfig.Add("action", "edit_aviatrix_transit_advanced_config")
	editTransitAdvancedConfig.Add("subaction", "bgp_manual_spoke")
	editTransitAdvancedConfig.Add("gateway_name", transitGw.GwName)
	editTransitAdvancedConfig.Add("bgp_manual_spoke_advertise_cidrs", transitGw.BgpManualSpokeAdvertiseCidrs)
	Url.RawQuery = editTransitAdvancedConfig.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get edit_aviatrix_transit_advanced_config failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode edit_aviatrix_transit_advanced_config failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API edit_aviatrix_transit_advanced_config Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) EnableTransitLearnedCidrsApproval(gateway *TransitVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'enable_transit_learned_cidrs_approval': ") + err.Error())
	}
	enableTransitLearnedCidrsApproval := url.Values{}
	enableTransitLearnedCidrsApproval.Add("CID", c.CID)
	enableTransitLearnedCidrsApproval.Add("action", "enable_transit_learned_cidrs_approval")
	enableTransitLearnedCidrsApproval.Add("gateway_name", gateway.GwName)
	Url.RawQuery = enableTransitLearnedCidrsApproval.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get 'enable_transit_learned_cidrs_approval' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'enable_transit_learned_cidrs_approval' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'enable_transit_learned_cidrs_approval' Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DisableTransitLearnedCidrsApproval(gateway *TransitVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'disable_transit_learned_cidrs_approval': ") + err.Error())
	}
	disableTransitLearnedCidrsApproval := url.Values{}
	disableTransitLearnedCidrsApproval.Add("CID", c.CID)
	disableTransitLearnedCidrsApproval.Add("action", "disable_transit_learned_cidrs_approval")
	disableTransitLearnedCidrsApproval.Add("gateway_name", gateway.GwName)
	Url.RawQuery = disableTransitLearnedCidrsApproval.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get 'disable_transit_learned_cidrs_approval' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'disable_transit_learned_cidrs_approval' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'disable_transit_learned_cidrs_approval' Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) SetBgpPollingTime(transitGateway *TransitVpc, newPollingTime string) error {
	action := "change_bgp_polling_time"
	return c.PostAPI(action, struct {
		CID         string `form:"CID"`
		Action      string `form:"action"`
		GatewayName string `form:"gateway_name"`
		PollingTime string `form:"bgp_polling_time"`
	}{
		CID:         c.CID,
		Action:      action,
		GatewayName: transitGateway.GwName,
		PollingTime: newPollingTime,
	}, BasicCheck)
}

func (c *Client) SetPrependASPath(transitGateway *TransitVpc, prependASPath []string) error {
	action, subaction := "edit_aviatrix_transit_advanced_config", "prepend_as_path"
	return c.PostAPI(action+"/"+subaction, struct {
		CID           string `form:"CID"`
		Action        string `form:"action"`
		Subaction     string `form:"subaction"`
		GatewayName   string `form:"gateway_name"`
		PrependASPath string `form:"bgp_prepend_as_path"`
	}{
		CID:           c.CID,
		Action:        action,
		Subaction:     subaction,
		GatewayName:   transitGateway.GwName,
		PrependASPath: strings.Join(prependASPath, " "),
	}, BasicCheck)
}

func (c *Client) SetLocalASNumber(transitGateway *TransitVpc, localASNumber string) error {
	action := "edit_transit_local_as_number"
	return c.PostAPI(action, struct {
		CID           string `form:"CID"`
		Action        string `form:"action"`
		GatewayName   string `form:"gateway_name"`
		LocalASNumber string `form:"local_as_num"`
	}{
		CID:           c.CID,
		Action:        action,
		GatewayName:   transitGateway.GwName,
		LocalASNumber: localASNumber,
	}, func(action, reason string, ret bool) error {
		if !ret {
			// Tried to set ASN to the same value, don't fail
			if strings.Contains(reason, "No change on transit gateway") {
				return nil
			}
			return fmt.Errorf("rest API %s Post failed: %s", action, reason)
		}
		return nil
	})
}

func (c *Client) SetBgpEcmp(transitGateway *TransitVpc, enabled bool) error {
	action := "enable_bgp_ecmp"
	if !enabled {
		action = "disable_bgp_ecmp"
	}
	return c.PostAPI(action, struct {
		CID         string `form:"CID"`
		Action      string `form:"action"`
		GatewayName string `form:"gateway_name"`
	}{
		CID:         c.CID,
		Action:      action,
		GatewayName: transitGateway.GwName,
	}, BasicCheck)
}

func (c *Client) GetTransitGatewayAdvancedConfig(transitGateway *TransitVpc) (*TransitGatewayAdvancedConfig, error) {
	action := "list_aviatrix_transit_advanced_config"
	resp, err := c.Post(c.baseURL, struct {
		CID         string `form:"CID"`
		Action      string `form:"action"`
		GatewayName string `form:"transit_gateway_name"`
	}{
		CID:         c.CID,
		Action:      action,
		GatewayName: transitGateway.GwName,
	})
	if err != nil {
		return nil, fmt.Errorf("HTTP POST %s failed: %v", action, err)
	}

	var data TransitGatewayAdvancedConfigResp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body %s failed: %v", action, err)
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return nil, fmt.Errorf("json Decode %s failed: %v\n Body: %s", action, err, b.String())
	}

	if !data.Return {
		return nil, fmt.Errorf("rest API %s Post failed: %s", action, data.Reason)
	}

	prependASPathStrings := strings.Split(data.Results.PrependASPath, " ")
	var filteredStrings []string
	for _, v := range prependASPathStrings {
		if v != "" {
			filteredStrings = append(filteredStrings, v)
		}
	}

	var standbyConnections []StandbyConnection
	for k, v := range data.Results.ActiveStandbyStatus {
		gwType := "Primary"
		if strings.HasSuffix(v, "-hagw") {
			gwType = "HA"
		}

		standbyConnections = append(standbyConnections, StandbyConnection{
			ConnectionName:    k,
			ActiveGatewayType: gwType,
		})
	}

	return &TransitGatewayAdvancedConfig{
		BgpPollingTime:                    strconv.Itoa(data.Results.BgpPollingTime),
		PrependASPath:                     filteredStrings,
		LocalASNumber:                     data.Results.LocalASNumber,
		BgpEcmpEnabled:                    data.Results.BgpEcmpEnabled == "yes",
		ActiveStandbyEnabled:              data.Results.ActiveStandby == "yes",
		ActiveStandbyConnections:          standbyConnections,
		LearnedCIDRsApprovalMode:          data.Results.LearnedCIDRsApprovalMode,
		ConnectionLearnedCIDRApprovalInfo: data.Results.ConnectionLearnedCIDRApprovalInfo,
		TunnelAddrLocal:                   data.Results.TunnelAddrLocal,
		TunnelAddrLocalBackup:             data.Results.TunnelAddrLocalBackup,
		PeerVnetId:                        data.Results.PeerVnetId,
		BgpHoldTime:                       data.Results.BgpHoldTime,
	}, nil
}

func (c *Client) SetTransitLearnedCIDRsApprovalMode(gw *TransitVpc, mode string) error {
	data := map[string]string{
		"action":       "set_transit_learned_cidrs_approval_mode",
		"CID":          c.CID,
		"gateway_name": gw.GwName,
		"mode":         mode,
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) EnableTransitConnectionLearnedCIDRApproval(gwName, connName string) error {
	data := map[string]string{
		"action":          "enable_transit_connection_learned_cidrs_approval",
		"CID":             c.CID,
		"gateway_name":    gwName,
		"connection_name": connName,
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) DisableTransitConnectionLearnedCIDRApproval(gwName, connName string) error {
	data := map[string]string{
		"action":          "disable_transit_connection_learned_cidrs_approval",
		"CID":             c.CID,
		"gateway_name":    gwName,
		"connection_name": connName,
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) EditTransitConnectionBGPManualAdvertiseCIDRs(gwName, connName string, cidrs []string) error {
	data := map[string]string{
		"action":                                "edit_transit_connection_bgp_manual_advertise_cidrs",
		"CID":                                   c.CID,
		"gateway_name":                          gwName,
		"connection_name":                       connName,
		"connection_bgp_manual_advertise_cidrs": strings.Join(cidrs, ","),
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) ChangeBgpHoldTime(gwName string, holdTime int) error {
	data := map[string]string{
		"action":        "change_bgp_hold_time",
		"gateway_name":  gwName,
		"bgp_hold_time": strconv.Itoa(holdTime),
		"CID":           c.CID,
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}
