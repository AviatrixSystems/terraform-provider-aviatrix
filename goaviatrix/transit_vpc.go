package goaviatrix

import (
	"fmt"
	"strconv"
	"strings"
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
	VNetNameResourceGroup        string `form:"vnet_and_resource_group_names,omitempty"`
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
	TagList                      string `form:"tag_string,omitempty"`
	TagJson                      string `form:"tag_json,omitempty"`
	EnableHybridConnection       bool   `form:"enable_hybrid_connection" json:"tgw_enabled,omitempty"`
	ConnectedTransit             string `form:"connected_transit" json:"connected_transit,omitempty"`
	InsaneMode                   string `form:"insane_mode,omitempty"`
	ReuseEip                     string `form:"reuse_eip,omitempty"`
	AllocateNewEipRead           bool   `json:"newly_allocated_eip,omitempty"`
	Eip                          string `form:"eip,omitempty"`
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
	EnableSummarizeCidrToTgw     bool
	AvailabilityDomain           string   `form:"availability_domain,omitempty"`
	FaultDomain                  string   `form:"fault_domain,omitempty"`
	EnableSpotInstance           bool     `form:"spot_instance,omitempty"`
	SpotPrice                    string   `form:"spot_price,omitempty"`
	ApprovedLearnedCidrs         []string `form:"approved_learned_cidrs"`
	BgpLanVpcID                  string   `form:"bgp_lan_vpc"`
	BgpLanSpecifySubnet          string   `form:"bgp_lan_specify_subnet"`
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
	EnableSummarizeCidrToTgw          bool
	ApprovedLearnedCidrs              []string
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

func (c *Client) LaunchTransitVpc(gateway *TransitVpc) error {
	gateway.CID = c.CID
	gateway.Action = "create_transit_gw"

	return c.PostAPI(gateway.Action, gateway, BasicCheck)
}

func (c *Client) EnableHaTransitGateway(gateway *TransitVpc) error {
	gateway.CID = c.CID
	gateway.Action = "create_peering_ha_gateway"

	return c.PostAPI(gateway.Action, gateway, BasicCheck)
}

func (c *Client) EnableHaTransitVpc(gateway *TransitVpc) error {
	form := map[string]string{
		"CID":     c.CID,
		"action":  "enable_transit_ha",
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
		"action":       "enable_transit_learned_cidrs_approval",
		"gateway_name": gateway.GwName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) DisableTransitLearnedCidrsApproval(gateway *TransitVpc) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "disable_transit_learned_cidrs_approval",
		"gateway_name": gateway.GwName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) UpdateTransitPendingApprovedCidrs(gateway *TransitVpc) error {
	form := map[string]string{
		"CID":                    c.CID,
		"action":                 "update_transit_pending_approved_cidrs",
		"gateway_name":           gateway.GwName,
		"approved_learned_cidrs": strings.Join(gateway.ApprovedLearnedCidrs, ","),
	}

	return c.PostAPI(form["action"], form, BasicCheck)
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
		"CID":                  c.CID,
		"action":               "list_aviatrix_transit_advanced_config",
		"transit_gateway_name": transitGateway.GwName,
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
		EnableSummarizeCidrToTgw:          data.Results.EnableSummarizeCidrToTgw == "yes",
		ApprovedLearnedCidrs:              data.Results.ApprovedLearnedCidrs,
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

func (c *Client) UpdateTransitConnectionPendingApprovedCidrs(gwName, connName string, approvedCidrs []string) error {
	data := map[string]string{
		"action":                            "update_transit_connection_pending_approved_cidrs",
		"CID":                               c.CID,
		"gateway_name":                      gwName,
		"connection_name":                   connName,
		"connection_approved_learned_cidrs": strings.Join(approvedCidrs, ","),
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
		"action":            "edit_site2cloud_conn",
		"CID":               c.CID,
		"vpc_id":            vpcId,
		"conn_name":         connName,
		"network_type":      "2",
		"cloud_subnet_cidr": remoteSubnet,
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}
