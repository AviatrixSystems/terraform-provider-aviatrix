package goaviatrix

import (
	"fmt"
	"strings"
)

// Spoke gateway simple struct to hold spoke details
type SpokeVpc struct {
	AccountName                  string `form:"account_name,omitempty" json:"account_name,omitempty"`
	Action                       string `form:"action,omitempty"`
	CID                          string `form:"CID,omitempty"`
	CloudType                    int    `form:"cloud_type,omitempty" json:"cloud_type,omitempty"`
	DnsServer                    string `form:"dns_server,omitempty" json:"dns_server,omitempty"`
	GwName                       string `form:"gw_name,omitempty" json:"vpc_name,omitempty"`
	GwSize                       string `form:"gw_size,omitempty"`
	VpcID                        string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	VNetNameResourceGroup        string `form:"vnet_and_resource_group_names,omitempty"`
	Subnet                       string `form:"gw_subnet,omitempty" json:"gw_subnet,omitempty"`
	VpcRegion                    string `form:"vpc_region,omitempty" json:"vpc_region,omitempty"`
	VpcSize                      string `form:"gw_size,omitempty" json:"vpc_size,omitempty"`
	EnableNat                    string `form:"enable_nat,omitempty" json:"enable_nat,omitempty"`
	EnableVpcDnsServer           string `json:"use_vpc_dns,omitempty"`
	HASubnet                     string `form:"ha_subnet,omitempty"`
	HAZone                       string `form:"new_zone,omitempty"`
	HASubnetGCP                  string `form:"new_subnet,omitempty"`
	SingleAzHa                   string `form:"single_az_ha,omitempty"`
	TransitGateway               string `form:"transit_gw,omitempty"`
	TagJson                      string `form:"json_tags,omitempty"`
	ReuseEip                     string `form:"reuse_eip,omitempty"`
	AllocateNewEipRead           bool   `json:"newly_allocated_eip,omitempty"`
	Eip                          string `form:"eip,omitempty" json:"eip,omitempty"`
	InsaneMode                   string `form:"insane_mode,omitempty"`
	Zone                         string `form:"zone,omitempty" json:"zone,omitempty"`
	BgpManualSpokeAdvertiseCidrs string `form:"bgp_manual_spoke,omitempty"`
	EncVolume                    string `form:"enc_volume,omitempty"`
	CustomerManagedKeys          string `form:"cmk,omitempty"`
	EnablePrivateOob             string `form:"private_oob,omitempty"`
	OobManagementSubnet          string `form:"oob_mgmt_subnet,omitempty"`
	HAOobManagementSubnet        string
	AvailabilityDomain           string   `form:"availability_domain,omitempty"`
	FaultDomain                  string   `form:"fault_domain,omitempty"`
	EnableSpotInstance           bool     `form:"spot_instance,omitempty"`
	SpotPrice                    string   `form:"spot_price,omitempty"`
	DeleteSpot                   bool     `form:"delete_spot,omitempty"`
	EnableBgp                    string   `form:"enable_bgp"`
	LearnedCidrsApproval         string   `form:"learned_cidrs_approval,omitempty"`
	ApprovedLearnedCidrs         []string `form:"approved_learned_cidrs"`
	Async                        bool     `form:"async,omitempty"`
	BgpOverLan                   bool     `form:"bgp_lan,omitempty"`
	BgpLanInterfacesCount        int      `form:"bgp_lan_intf_count,omitempty"`
	LbVpcId                      string   `form:"lb_vpc_id,omitempty"`
	EnableGlobalVpc              bool     `form:"global_vpc"`
}

type SpokeGatewayAdvancedConfig struct {
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

type SpokeGatewayAdvancedConfigResp struct {
	Return  bool                                 `json:"return"`
	Results SpokeGatewayAdvancedConfigRespResult `json:"results"`
	Reason  string                               `json:"reason"`
}

type SpokeGatewayAdvancedConfigRespResult struct {
	BgpPollingTime                    int                       `json:"bgp_polling_time"`
	BgpBfdPollingTime                 int                       `json:"bgp_neighbor_status_polling_time,omitempty"`
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

func (c *Client) LaunchSpokeVpc(spoke *SpokeVpc) error {
	spoke.CID = c.CID
	spoke.Action = "create_multicloud_primary_gateway"
	spoke.Async = true

	return c.PostAsyncAPI(spoke.Action, spoke, BasicCheck)
}

func (c *Client) SpokeJoinTransit(spoke *SpokeVpc) error {
	form := map[string]string{
		"CID":        c.CID,
		"action":     "attach_spoke_to_transit_gw",
		"spoke_gw":   spoke.GwName,
		"transit_gw": spoke.TransitGateway,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) SpokeLeaveAllTransit(spoke *SpokeVpc) error {
	form := map[string]string{
		"CID":      c.CID,
		"action":   "detach_spoke_from_transit_gw",
		"spoke_gw": spoke.GwName,
	}

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "has not joined to any transit") {
				return nil
			}
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	return c.PostAPI(form["action"], form, checkFunc)
}

func (c *Client) SpokeLeaveTransit(spoke *SpokeVpc) error {
	action := "detach_spoke_from_transit_gw"
	data := map[string]interface{}{
		"CID":        c.CID,
		"action":     action,
		"spoke_gw":   spoke.GwName,
		"transit_gw": spoke.TransitGateway,
	}
	return c.PostAPI(action, data, BasicCheck)
}

func (c *Client) EnableHaSpokeVpc(spoke *SpokeVpc) error {
	form := map[string]string{
		"CID":     c.CID,
		"action":  "create_multicloud_ha_gateway",
		"gw_name": spoke.GwName,
		"eip":     spoke.Eip,
		"async":   "true",
	}

	if IsCloudType(spoke.CloudType, GCPRelatedCloudTypes) {
		form["new_zone"] = spoke.HAZone
	} else {
		form["public_subnet"] = spoke.HASubnet
		form["oob_mgmt_subnet"] = spoke.HAOobManagementSubnet
	}

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "HA GW already exists") {
				return nil
			}
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	return c.PostAsyncAPI(form["action"], form, checkFunc)
}

func (c *Client) EnableHaSpokeGateway(gateway *SpokeVpc) error {
	gateway.CID = c.CID
	gateway.Action = "create_peering_ha_gateway"

	return c.PostAPI(gateway.Action, gateway, BasicCheck)
}

func (c *Client) EnableAutoAdvertiseS2CCidrs(gateway *Gateway) error {
	form := map[string]string{
		"action":       "enable_auto_advertise_s2c_cidrs",
		"CID":          c.CID,
		"gateway_name": gateway.GwName,
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) DisableAutoAdvertiseS2CCidrs(gateway *Gateway) error {
	form := map[string]string{
		"action":       "disable_auto_advertise_s2c_cidrs",
		"CID":          c.CID,
		"gateway_name": gateway.GwName,
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) GetSpokeGatewayAdvancedConfig(spokeGateway *SpokeVpc) (*SpokeGatewayAdvancedConfig, error) {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "list_aviatrix_spoke_advanced_config",
		"gateway_name": spokeGateway.GwName,
	}

	var data SpokeGatewayAdvancedConfigResp

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

	return &SpokeGatewayAdvancedConfig{
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

func (c *Client) EnableSpokeConnectionLearnedCIDRApproval(gwName, connName string) error {
	data := map[string]string{
		"action":          "enable_bgp_connection_cidr_approval",
		"CID":             c.CID,
		"gateway_name":    gwName,
		"connection_name": connName,
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) DisableSpokeConnectionLearnedCIDRApproval(gwName, connName string) error {
	data := map[string]string{
		"action":          "disable_bgp_connection_cidr_approval",
		"CID":             c.CID,
		"gateway_name":    gwName,
		"connection_name": connName,
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) UpdateSpokeConnectionPendingApprovedCidrs(gwName, connName string, approvedCidrs []string) error {
	data := map[string]string{
		"action":          "set_bgp_connection_approved_cidr_rules",
		"CID":             c.CID,
		"gateway_name":    gwName,
		"connection_name": connName,
		"cidr_rules":      strings.Join(approvedCidrs, ","),
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) EditSpokeConnectionBGPManualAdvertiseCIDRs(gwName, connName string, cidrs []string) error {
	data := map[string]string{
		"action":                                "edit_spoke_connection_bgp_manual_advertise_cidrs",
		"CID":                                   c.CID,
		"gateway_name":                          gwName,
		"connection_name":                       connName,
		"connection_bgp_manual_advertise_cidrs": strings.Join(cidrs, ","),
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) SetBgpEcmpSpoke(spokeGateway *SpokeVpc, enabled bool) error {
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
		GatewayName: spokeGateway.GwName,
	}, BasicCheck)
}

func (c *Client) EnableActiveStandbySpoke(spokeGateway *SpokeVpc) error {
	action := "enable_active_standby"
	form := map[string]string{
		"CID":          c.CID,
		"action":       action,
		"gateway_name": spokeGateway.GwName,
	}
	return c.PostAPI(action, form, BasicCheck)
}

func (c *Client) DisableActiveStandbySpoke(spokeGateway *SpokeVpc) error {
	action := "disable_active_standby"
	form := map[string]string{
		"CID":          c.CID,
		"action":       action,
		"gateway_name": spokeGateway.GwName,
	}
	return c.PostAPI(action, form, BasicCheck)
}

func (c *Client) SetPrependASPathSpoke(spokeGateway *SpokeVpc, prependASPath []string) error {
	action, subaction := "edit_aviatrix_spoke_advanced_config", "prepend_as_path"
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
		GatewayName:   spokeGateway.GwName,
		PrependASPath: strings.Join(prependASPath, " "),
	}, BasicCheck)
}

func (c *Client) SetBgpPollingTimeSpoke(spokeGateway *SpokeVpc, newPollingTime int) error {
	action := "change_bgp_polling_time"
	return c.PostAPI(action, struct {
		CID         string `form:"CID"`
		Action      string `form:"action"`
		GatewayName string `form:"gateway_name"`
		PollingTime int    `form:"bgp_polling_time"`
	}{
		CID:         c.CID,
		Action:      action,
		GatewayName: spokeGateway.GwName,
		PollingTime: newPollingTime,
	}, BasicCheck)
}

func (c *Client) SetBgpBfdPollingTimeSpoke(spokeGateway *SpokeVpc, newPollingTime int) error {
	action := "change_bgp_neighbor_status_polling_time"
	return c.PostAPI(action, struct {
		CID         string `form:"CID"`
		Action      string `form:"action"`
		GatewayName string `form:"gateway_name"`
		PollingTime int    `form:"bgp_neighbor_status_polling_time"`
	}{
		CID:         c.CID,
		Action:      action,
		GatewayName: spokeGateway.GwName,
		PollingTime: newPollingTime,
	}, BasicCheck)
}

func (c *Client) SetSpokeBgpManualAdvertisedNetworks(spokeGateway *SpokeVpc) error {
	form := map[string]string{
		"CID":                              c.CID,
		"action":                           "edit_aviatrix_spoke_advanced_config",
		"subaction":                        "bgp_manual_spoke",
		"gateway_name":                     spokeGateway.GwName,
		"bgp_manual_spoke_advertise_cidrs": spokeGateway.BgpManualSpokeAdvertiseCidrs,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) EnableSpokeLearnedCidrsApproval(gateway *SpokeVpc) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "enable_bgp_gateway_cidr_approval",
		"gateway_name": gateway.GwName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) DisableSpokeLearnedCidrsApproval(gateway *SpokeVpc) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "disable_bgp_gateway_cidr_approval",
		"gateway_name": gateway.GwName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) UpdateSpokePendingApprovedCidrs(gateway *SpokeVpc) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "set_bgp_gateway_approved_cidr_rules",
		"gateway_name": gateway.GwName,
		"cidr_rules":   strings.Join(gateway.ApprovedLearnedCidrs, ","),
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) SetLocalASNumberSpoke(spokeGateway *SpokeVpc, localASNumber string) error {
	action := "edit_spoke_local_as_number"
	return c.PostAPI(action, struct {
		CID           string `form:"CID"`
		Action        string `form:"action"`
		GatewayName   string `form:"gateway_name"`
		LocalASNumber string `form:"local_as_num"`
	}{
		CID:           c.CID,
		Action:        action,
		GatewayName:   spokeGateway.GwName,
		LocalASNumber: localASNumber,
	}, func(action, method, reason string, ret bool) error {
		if !ret {
			// Tried to set ASN to the same value, don't fail
			if strings.Contains(reason, "No change on gateway") {
				return nil
			}
			return fmt.Errorf("rest API %s %s failed: %s", action, method, reason)
		}
		return nil
	})
}

func (c *Client) EnableActiveStandbyPreemptiveSpoke(spokeGateway *SpokeVpc) error {
	action := "enable_active_standby"
	form := map[string]string{
		"CID":          c.CID,
		"action":       action,
		"gateway_name": spokeGateway.GwName,
		"preemptive":   "true",
	}
	return c.PostAPI(action, form, BasicCheck)
}

func (c *Client) EnableSpokeOnpremRoutePropagation(spokeGateway *SpokeVpc) error {
	action := "enable_spoke_onprem_route_propagation"
	form := map[string]string{
		"CID":          c.CID,
		"action":       action,
		"gateway_name": spokeGateway.GwName,
	}
	return c.PostAPI(action, form, BasicCheck)
}

func (c *Client) DisableSpokeOnpremRoutePropagation(spokeGateway *SpokeVpc) error {
	action := "disable_spoke_onprem_route_propagation"
	form := map[string]string{
		"CID":          c.CID,
		"action":       action,
		"gateway_name": spokeGateway.GwName,
	}
	return c.PostAPI(action, form, BasicCheck)
}

func (c *Client) EnableSpokePreserveAsPath(spokeGateway *SpokeVpc) error {
	action := "enable_spoke_preserve_as_path"
	data := map[string]interface{}{
		"CID":          c.CID,
		"action":       action,
		"gateway_name": spokeGateway.GwName,
	}
	return c.PostAPI(action, data, BasicCheck)
}

func (c *Client) DisableSpokePreserveAsPath(spokeGateway *SpokeVpc) error {
	action := "disable_spoke_preserve_as_path"
	data := map[string]interface{}{
		"CID":          c.CID,
		"action":       action,
		"gateway_name": spokeGateway.GwName,
	}
	return c.PostAPI(action, data, BasicCheck)
}

func (c *Client) EnableGlobalVpc(gateway *Gateway) error {
	form := map[string]string{
		"action":       "enable_global_vpc",
		"CID":          c.CID,
		"gateway_name": gateway.GwName,
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) DisableGlobalVpc(gateway *Gateway) error {
	form := map[string]string{
		"action":       "disable_global_vpc",
		"CID":          c.CID,
		"gateway_name": gateway.GwName,
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}
