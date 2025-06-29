package goaviatrix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Gateway simple struct to hold gateway details
type TransitVpc struct {
	AccountName                  string `form:"account_name,omitempty" json:"account_name,omitempty"`
	Action                       string `form:"action,omitempty" json:"action,omitempty"`
	CID                          string `form:"CID,omitempty" json:"CID,omitempty"`
	CloudType                    int    `form:"cloud_type,omitempty" json:"cloud_type,omitempty"`
	DnsServer                    string `form:"dns_server,omitempty" json:"dns_server,omitempty"`
	GwName                       string `form:"gw_name,omitempty" json:"name,omitempty"`
	GwSize                       string `form:"gw_size,omitempty"`
	VpcID                        string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	VNetNameResourceGroup        string `form:"vnet_and_resource_group_names,omitempty"`
	Subnet                       string `form:"gw_subnet,omitempty" json:"gw_subnet,omitempty"`
	HASubnet                     string `form:"ha_subnet,omitempty"`
	HAZone                       string `form:"new_zone,omitempty"`
	HASubnetGCP                  string `form:"new_subnet,omitempty"`
	PeeringHASubnet              string `json:"public_subnet,omitempty"`
	VpcRegion                    string `form:"vpc_region,omitempty" json:"vpc_region,omitempty"`
	VpcSize                      string `form:"gw_size,omitempty" json:"gw_size,omitempty"`
	EnableNAT                    string `form:"enable_nat,omitempty" json:"enable_nat,omitempty"`
	SingleAzHa                   string `form:"single_az_ha,omitempty"`
	EnableVpcDnsServer           string `json:"use_vpc_dns,omitempty"`
	TagJson                      string `form:"json_tags,omitempty"`
	EnableHybridConnection       bool   `form:"enable_hybrid_connection" json:"tgw_enabled,omitempty"`
	ConnectedTransit             string `form:"connected_transit" json:"connected_transit,omitempty"`
	InsaneMode                   string `form:"insane_mode,omitempty"`
	ReuseEip                     string `form:"reuse_eip,omitempty"`
	AllocateNewEipRead           bool   `json:"newly_allocated_eip,omitempty"`
	Eip                          string `form:"eip,omitempty"`
	Zone                         string `form:"zone,omitempty" json:"zone,omitempty"`
	EnableAdvertiseTransitCidr   bool
	BgpManualSpokeAdvertiseCidrs string `form:"bgp_manual_spoke,omitempty"`
	EnableTransitFireNet         bool   `form:"firenet,omitempty"`
	LanVpcID                     string `form:"lan_vpc,omitempty"`
	LanPrivateSubnet             string `form:"lan_subnet,omitempty"`
	LearnedCidrsApproval         string `form:"learned_cidrs_approval,omitempty"`
	EncVolume                    string `form:"enc_volume,omitempty"`
	CustomerManagedKeys          string `form:"cmk,omitempty"`
	BgpOverLan                   bool   `form:"bgp_lan,omitempty"`
	EnablePrivateOob             string `form:"private_oob,omitempty"`
	OobManagementSubnet          string `form:"oob_mgmt_subnet,omitempty"`
	HAOobManagementSubnet        string
	EnableSummarizeCidrToTgw     bool
	AvailabilityDomain           string              `form:"availability_domain,omitempty"`
	FaultDomain                  string              `form:"fault_domain,omitempty"`
	EnableSpotInstance           bool                `form:"spot_instance,omitempty"`
	SpotPrice                    string              `form:"spot_price,omitempty"`
	DeleteSpot                   bool                `form:"delete_spot,omitempty"`
	ApprovedLearnedCidrs         []string            `form:"approved_learned_cidrs"`
	BgpLanVpcID                  string              `form:"bgp_lan_vpc"`
	BgpLanSpecifySubnet          string              `form:"bgp_lan_subnet"`
	Async                        bool                `form:"async,omitempty"`
	BgpLanInterfacesCount        int                 `form:"bgp_lan_intf_count,omitempty"`
	LbVpcID                      string              `form:"lb_vpc_id,omitempty"`
	Transit                      bool                `form:"transit,omitempty"`
	DeviceID                     string              `form:"device_id,omitempty"`
	SiteID                       string              `form:"site_id,omitempty"`
	Interfaces                   string              `json:"interfaces,omitempty"`
	InterfaceMapping             string              `json:"interface_mapping,omitempty"`
	EipMap                       string              `json:"eip_map,omitempty"`
	LogicalEipMap                map[string][]EipMap `json:"logical_intf_eip_map,omitempty"`
	ZtpFileDownloadPath          string              `json:"-"`
	ManagementEgressIPPrefix     string              `json:"mgmt_egress_ip,omitempty"`
	JumboFrame                   bool                `json:"jumbo_frame,omitempty"`
}

type TransitGatewayAdvancedConfig struct {
	BgpPollingTime                    int
	BgpBfdPollingTime                 int
	PrependASPath                     []string
	LocalASNumber                     string
	BgpEcmpEnabled                    bool
	ActiveStandbyEnabled              bool
	ActiveStandbyConnections          []StandbyConnection
	LearnedCIDRsApprovalMode          string
	ConnectionLearnedCIDRApprovalInfo []LearnedCIDRApprovalInfo
	TunnelAddrLocal                   string
	TunnelAddrLocalBackup             string
	PeerVnetID                        []string
	BgpHoldTime                       int
	EnableSummarizeCidrToTgw          bool
	ApprovedLearnedCidrs              []string
}

type StandbyConnection struct {
	ConnectionName    string
	ActiveGatewayType string
}

type TransitGatewayAdvancedConfigRespResult struct {
	BgpPollingTime                    int                       `json:"bgp_polling_time"`
	BgpBfdPollingTime                 int                       `json:"bgp_neighbor_status_polling_time"`
	PrependASPath                     string                    `json:"bgp_prepend_as_path"`
	LocalASNumber                     string                    `json:"local_asn_num"`
	BgpEcmpEnabled                    string                    `json:"bgp_ecmp"`
	ActiveStandby                     string                    `json:"active-standby"`
	ActiveStandbyStatus               map[string]string         `json:"active_standby_status"`
	LearnedCIDRsApprovalMode          string                    `json:"learned_cidrs_approval_mode"`
	ConnectionLearnedCIDRApprovalInfo []LearnedCIDRApprovalInfo `json:"connection_learned_cidrs_approval_info"`
	TunnelAddrLocal                   string                    `json:"tunnel_addr_local"`
	TunnelAddrLocalBackup             string                    `json:"tunnel_addr_local_backup"`
	PeerVnetID                        []string                  `json:"peer_vnet_id"`
	BgpHoldTime                       int                       `json:"bgp_hold_time"`
	EnableSummarizeCidrToTgw          string                    `json:"summarize_cidr_to_tgw"`
	ApprovedLearnedCidrs              []string                  `json:"approved_learned_cidrs"`
}

type LearnedCIDRApprovalInfo struct {
	ConnName             string   `json:"conn_name"`
	EnabledApproval      string   `json:"conn_learned_cidrs_approval"`
	ApprovedLearnedCidrs []string `json:"conn_approved_learned_cidrs"`
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

type EdgeTransitInterface struct {
	Name           string   `json:"ifname"`
	Type           string   `json:"type"`
	Index          int      `json:"index,omitempty"`
	PublicIp       string   `json:"public_ip,omitempty"`
	Dhcp           bool     `json:"dhcp,omitempty"`
	IpAddress      string   `json:"ipaddr,omitempty"`
	GatewayIp      string   `json:"gateway_ip,omitempty"`
	SecondaryCIDRs []string `json:"secondary_private_cidr_list,omitempty"`
	LogicalIfName  string   `json:"logical_ifname,omitempty"`
}

type EipMap struct {
	PrivateIP string `json:"private_ip"`
	PublicIP  string `json:"public_ip"`
}

type TransitGatewayBgpLanIpInfoResp struct {
	Return  bool                                 `json:"return"`
	Results TransitGatewayBgpLanIpInfoRespResult `json:"results"`
	Reason  string                               `json:"reason"`
}

type TransitGatewayBgpLanIpInfoRespResult struct {
	BgpLanIpList        []string `json:"gce_bgp_lan_all_intf_tuple_list"`
	HaBgpLanIpList      []string `json:"gce_bgp_lan_all_intf_ha_tuple_list"`
	AzureBgpLanIpList   []string `json:"arm_bgp_lan_all_intf_ip_list"`
	AzureHaBgpLanIpList []string `json:"arm_bgp_lan_all_intf_ha_ip_list"`
}

type TransitGatewayBgpLanIpInfo struct {
	BgpLanIpList        []string
	HaBgpLanIpList      []string
	AzureBgpLanIpList   []string
	AzureHaBgpLanIpList []string
}

func (c *Client) LaunchTransitVpc(gateway *TransitVpc) error {
	gateway.CID = c.CID
	gateway.Action = "create_multicloud_primary_gateway"
	var data CreateEdgeEquinixResp
	err := c.PostAPIWithResponse(&data, gateway.Action, gateway, BasicCheck)
	if err != nil {
		return err
	}
	// create the ZTP file for Equinix and Megaport edge transit gateway
	if IsCloudType(gateway.CloudType, EDGEEQUINIX|EDGEMEGAPORT|EDGESELFMANAGED) {
		fileName := getFileName(gateway.ZtpFileDownloadPath, gateway.GwName, gateway.VpcID)
		fileContent, err := processZtpFileContent(data.Result)
		if err != nil {
			return err
		}
		err = createZtpFile(fileName, fileContent)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) EnableHaTransitGateway(gateway *TransitVpc) error {
	gateway.CID = c.CID
	gateway.Action = "create_peering_ha_gateway"

	return c.PostAPI(gateway.Action, gateway, BasicCheck)
}

func (c *Client) EnableHaTransitVpc(gateway *TransitVpc) error {
	form := map[string]string{
		"CID":     c.CID,
		"action":  "create_multicloud_ha_gateway",
		"gw_name": gateway.GwName,
		"eip":     gateway.Eip,
	}

	if gateway.CloudType == GCP {
		form["new_zone"] = gateway.HAZone
	} else {
		form["public_subnet"] = gateway.HASubnet
		form["oob_mgmt_subnet"] = gateway.HAOobManagementSubnet
	}

	checkFunc := func(action, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "HA GW already exists") {
				return nil
			}
			return fmt.Errorf("rest API %s %s failed: %s", action, method, reason)
		}
		return nil
	}

	return c.PostAPI(form["action"], form, checkFunc)
}

func (c *Client) AttachTransitGWForHybrid(gateway *TransitVpc) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "enable_transit_gateway_interface_to_aws_tgw",
		"gateway_name": gateway.GwName,
	}

	checkFunc := func(action, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "already enabled tgw interface") {
				return nil
			}
			return fmt.Errorf("rest API %s %s failed: %s", action, method, reason)
		}
		return nil
	}

	return c.PostAPI(form["action"], form, checkFunc)
}

func (c *Client) UpdateEdgeGateway(gateway *TransitVpc) error {
	form := map[string]interface{}{
		"CID":          c.CID,
		"action":       "update_edge_gateway",
		"gateway_name": gateway.GwName,
	}

	if gateway.Interfaces != "" {
		form["interfaces"] = gateway.Interfaces
	}

	if gateway.EipMap != "" {
		form["eip_map"] = gateway.EipMap
	}

	if len(gateway.LogicalEipMap) > 0 {
		eipMapJSON, err := json.Marshal(gateway.LogicalEipMap)
		if err != nil {
			return fmt.Errorf("failed to marshal eip_map to JSON: %w", err)
		}
		eipMapJSONObj := bytes.NewBuffer(eipMapJSON)
		form["logical_intf_eip_map"] = eipMapJSONObj
	}

	if gateway.ManagementEgressIPPrefix != "" {
		form["mgmt_egress_ip"] = gateway.ManagementEgressIPPrefix
	}

	action, ok := form["action"].(string)
	if !ok {
		return fmt.Errorf("form[action] is not a string, got type %T", form["action"])
	}
	log.Printf("Formm details: %v", form)
	return c.PostAPI(action, form, BasicCheck)
}

func (c *Client) UpdateEdgeGatewayV2(ctx context.Context, gateway *TransitVpc) error {
	gateway.CID = c.CID
	gateway.Action = "update_edge_gateway"
	return c.PostAPIContext2(ctx, nil, gateway.Action, gateway, BasicCheck)
}

func (c *Client) DeleteEdgeGateway(gateway *Gateway) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "delete_multicloud_gateway",
		"gateway_name": gateway.GwName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) DetachTransitGWForHybrid(gateway *TransitVpc) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "disable_transit_gateway_interface_to_aws_tgw",
		"gateway_name": gateway.GwName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) EnableConnectedTransit(gateway *TransitVpc) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "enable_connected_transit_on_gateway",
		"gateway_name": gateway.GwName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) DisableConnectedTransit(gateway *TransitVpc) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "disable_connected_transit_on_gateway",
		"gateway_name": gateway.GwName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) EnableGatewayFireNetInterfaces(gateway *TransitVpc) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "enable_gateway_firenet_interfaces",
		"gateway_name": gateway.GwName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) DisableGatewayFireNetInterfaces(gateway *TransitVpc) error {
	form := map[string]string{
		"CID":     c.CID,
		"action":  "disable_gateway_firenet_interfaces",
		"gateway": gateway.GwName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
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
	form := map[string]string{
		"CID":          c.CID,
		"action":       "enable_advertise_transit_cidr",
		"gateway_name": transitGw.GwName,
	}

	if transitGw.EnableAdvertiseTransitCidr {
		form["advertise_transit_cidr"] = "yes"
	} else {
		form["advertise_transit_cidr"] = "no"
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) DisableAdvertiseTransitCidr(transitGw *TransitVpc) error {
	form := map[string]string{
		"CID":                    c.CID,
		"action":                 "disable_advertise_transit_cidr",
		"gateway_name":           transitGw.GwName,
		"advertise_transit_cidr": "no",
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) SetBgpManualSpokeAdvertisedNetworks(transitGw *TransitVpc) error {
	form := map[string]string{
		"CID":                              c.CID,
		"action":                           "edit_aviatrix_transit_advanced_config",
		"subaction":                        "bgp_manual_spoke",
		"gateway_name":                     transitGw.GwName,
		"bgp_manual_spoke_advertise_cidrs": transitGw.BgpManualSpokeAdvertiseCidrs,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) EnableTransitLearnedCidrsApproval(gateway *TransitVpc) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "enable_bgp_gateway_cidr_approval",
		"gateway_name": gateway.GwName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) DisableTransitLearnedCidrsApproval(gateway *TransitVpc) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "disable_bgp_gateway_cidr_approval",
		"gateway_name": gateway.GwName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) UpdateTransitPendingApprovedCidrs(gateway *TransitVpc) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "set_bgp_gateway_approved_cidr_rules",
		"gateway_name": gateway.GwName,
		"cidr_rules":   strings.Join(gateway.ApprovedLearnedCidrs, ","),
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) SetBgpPollingTime(transitGateway *TransitVpc, newPollingTime int) error {
	action := "change_bgp_polling_time"
	return c.PostAPI(action, struct {
		CID         string `form:"CID"`
		Action      string `form:"action"`
		GatewayName string `form:"gateway_name"`
		PollingTime int    `form:"bgp_polling_time"`
	}{
		CID:         c.CID,
		Action:      action,
		GatewayName: transitGateway.GwName,
		PollingTime: newPollingTime,
	}, BasicCheck)
}

func (c *Client) SetBgpBfdPollingTime(transitGateway *TransitVpc, newPollingTime int) error {
	action := "change_bgp_neighbor_status_polling_time"
	return c.PostAPI(action, struct {
		CID         string `form:"CID"`
		Action      string `form:"action"`
		GatewayName string `form:"gateway_name"`
		PollingTime int    `form:"bgp_neighbor_status_polling_time"`
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
	}, func(action, method, reason string, ret bool) error {
		if !ret {
			// Tried to set ASN to the same value, don't fail
			if strings.Contains(reason, "No change on transit gateway") {
				return nil
			}
			return fmt.Errorf("rest API %s %s failed: %s", action, method, reason)
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
	form := map[string]string{
		"CID":          c.CID,
		"action":       "list_aviatrix_transit_advanced_config",
		"gateway_name": transitGateway.GwName,
	}

	var data TransitGatewayAdvancedConfigResp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
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
		BgpPollingTime:                    data.Results.BgpPollingTime,
		BgpBfdPollingTime:                 data.Results.BgpBfdPollingTime,
		PrependASPath:                     filteredStrings,
		LocalASNumber:                     data.Results.LocalASNumber,
		BgpEcmpEnabled:                    data.Results.BgpEcmpEnabled == "yes",
		ActiveStandbyEnabled:              data.Results.ActiveStandby == "yes",
		ActiveStandbyConnections:          standbyConnections,
		LearnedCIDRsApprovalMode:          data.Results.LearnedCIDRsApprovalMode,
		ConnectionLearnedCIDRApprovalInfo: data.Results.ConnectionLearnedCIDRApprovalInfo,
		TunnelAddrLocal:                   data.Results.TunnelAddrLocal,
		TunnelAddrLocalBackup:             data.Results.TunnelAddrLocalBackup,
		PeerVnetID:                        data.Results.PeerVnetID,
		BgpHoldTime:                       data.Results.BgpHoldTime,
		EnableSummarizeCidrToTgw:          data.Results.EnableSummarizeCidrToTgw == "yes",
		ApprovedLearnedCidrs:              data.Results.ApprovedLearnedCidrs,
	}, nil
}

func (c *Client) SetTransitLearnedCIDRsApprovalMode(gw *TransitVpc, mode string) error {
	data := map[string]string{
		"action":       "set_bgp_gateway_cidr_approval_mode",
		"CID":          c.CID,
		"gateway_name": gw.GwName,
		"mode":         mode,
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) EnableTransitConnectionLearnedCIDRApproval(gwName, connName string) error {
	data := map[string]string{
		"action":          "enable_bgp_connection_cidr_approval",
		"CID":             c.CID,
		"gateway_name":    gwName,
		"connection_name": connName,
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) DisableTransitConnectionLearnedCIDRApproval(gwName, connName string) error {
	data := map[string]string{
		"action":          "disable_bgp_connection_cidr_approval",
		"CID":             c.CID,
		"gateway_name":    gwName,
		"connection_name": connName,
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) UpdateTransitConnectionPendingApprovedCidrs(gwName, connName string, approvedCidrs []string) error {
	data := map[string]string{
		"action":          "set_bgp_connection_approved_cidr_rules",
		"CID":             c.CID,
		"gateway_name":    gwName,
		"connection_name": connName,
		"cidr_rules":      strings.Join(approvedCidrs, ","),
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

func (c *Client) EnableSummarizeCidrToTgw(gwName string) error {
	data := map[string]string{
		"action":       "enable_transit_summarize_cidr_to_tgw",
		"gateway_name": gwName,
		"CID":          c.CID,
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) DisableSummarizeCidrToTgw(gwName string) error {
	data := map[string]string{
		"action":       "disable_transit_summarize_cidr_to_tgw",
		"gateway_name": gwName,
		"CID":          c.CID,
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) EnableMultitierTransit(gwName string) error {
	data := map[string]string{
		"action":       "enable_multitier_transit",
		"gateway_name": gwName,
		"CID":          c.CID,
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) DisableMultitierTransit(gwName string) error {
	data := map[string]string{
		"action":       "disable_multitier_transit",
		"gateway_name": gwName,
		"CID":          c.CID,
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) EditTransitConnectionRemoteSubnet(vpcId, connName, remoteSubnet string) error {
	data := map[string]string{
		"action":      "edit_site2cloud_conn",
		"CID":         c.CID,
		"vpc_id":      vpcId,
		"conn_name":   connName,
		"remote_cidr": remoteSubnet,
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) GetBgpLanIPList(transitGateway *TransitVpc) (*TransitGatewayBgpLanIpInfo, error) {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "list_aviatrix_transit_advanced_config",
		"gateway_name": transitGateway.GwName,
	}

	var data TransitGatewayBgpLanIpInfoResp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	var bgpLanIpList []string
	var haBgpLanIpList []string
	var azureBgpLanIpList []string
	var azureHaBgpLanIpList []string
	for _, bgpLanIp := range data.Results.BgpLanIpList {
		bgpLanIpList = append(bgpLanIpList, strings.Split(bgpLanIp, ":")[2])
	}
	for _, haBgpLanIp := range data.Results.HaBgpLanIpList {
		haBgpLanIpList = append(haBgpLanIpList, strings.Split(haBgpLanIp, ":")[2])
	}
	azureBgpLanIpList = append(azureBgpLanIpList, data.Results.AzureBgpLanIpList...)
	azureHaBgpLanIpList = append(azureHaBgpLanIpList, data.Results.AzureHaBgpLanIpList...)

	return &TransitGatewayBgpLanIpInfo{
		BgpLanIpList:        bgpLanIpList,
		HaBgpLanIpList:      haBgpLanIpList,
		AzureBgpLanIpList:   azureBgpLanIpList,
		AzureHaBgpLanIpList: azureHaBgpLanIpList,
	}, nil
}

func (c *Client) EnableS2CRxBalancing(gwName string) error {
	data := map[string]string{
		"action":           "enable_s2c_rx_balancing",
		"gateway_name":     gwName,
		"CID":              c.CID,
		"s2c_rx_balancing": "yes",
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) DisableS2CRxBalancing(gwName string) error {
	data := map[string]string{
		"action":           "disable_s2c_rx_balancing",
		"gateway_name":     gwName,
		"CID":              c.CID,
		"s2c_rx_balancing": "no",
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) EnableTransitPreserveAsPath(transitGateway *TransitVpc) error {
	action := "enable_transit_preserve_as_path"
	data := map[string]interface{}{
		"CID":          c.CID,
		"action":       action,
		"gateway_name": transitGateway.GwName,
	}
	return c.PostAPI(action, data, BasicCheck)
}

func (c *Client) DisableTransitPreserveAsPath(transitGateway *TransitVpc) error {
	action := "disable_transit_preserve_as_path"
	data := map[string]interface{}{
		"CID":          c.CID,
		"action":       action,
		"gateway_name": transitGateway.GwName,
	}
	return c.PostAPI(action, data, BasicCheck)
}

func processZtpFileContent(cloudInitTransit string) (string, error) {
	var jsonCloudInit map[string]interface{}
	err := json.Unmarshal([]byte(cloudInitTransit), &jsonCloudInit)
	if err != nil {
		return "", fmt.Errorf("failed to parse cloud_init_transit as JSON: %w", err)
	}

	// Extract the 'text' field from the cloudinit data
	text, ok := jsonCloudInit["text"].(string)
	if !ok {
		return "", fmt.Errorf("'text' field not found or is not a string in cloud_init_transit")
	}
	return text, nil
}

// createZtpFile creates a new ztp file and writes the given content.
func createZtpFile(filePath, content string) error {
	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create the file: %w", err)
	}
	defer outFile.Close()

	// Write the content to the file
	_, err = outFile.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to the file: %w", err)
	}
	return nil
}

func getFileName(ztpFileDownloadPath, gwName, vpcID string) string {
	return ztpFileDownloadPath + "/" + gwName + "-" + vpcID + "-cloud-init.txt"
}
