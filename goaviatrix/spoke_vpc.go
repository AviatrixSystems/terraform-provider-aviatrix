package goaviatrix

import (
	"fmt"
	"strconv"
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
	Subnet                       string `form:"public_subnet,omitempty" json:"public_subnet,omitempty"`
	VpcRegion                    string `form:"region,omitempty" json:"vpc_region,omitempty"`
	VpcSize                      string `form:"gw_size,omitempty" json:"vpc_size,omitempty"`
	EnableNat                    string `form:"nat_enabled,omitempty" json:"enable_nat,omitempty"`
	EnableVpcDnsServer           string `json:"use_vpc_dns,omitempty"`
	HASubnet                     string `form:"ha_subnet,omitempty"`
	HAZone                       string `form:"new_zone,omitempty"`
	HASubnetGCP                  string `form:"new_subnet,omitempty"`
	SingleAzHa                   string `form:"single_az_ha,omitempty"`
	TransitGateway               string `form:"transit_gw,omitempty"`
	TagList                      string `form:"tag_string,omitempty"`
	TagJson                      string `form:"tag_json,omitempty"`
	ReuseEip                     string `form:"reuse_eip,omitempty"`
	AllocateNewEipRead           bool   `json:"newly_allocated_eip,omitempty"`
	Eip                          string `form:"eip,omitempty" json:"eip,omitempty"`
	InsaneMode                   string `form:"insane_mode,omitempty"`
	Zone                         string `form:"zone,omitempty" json:"zone,omitempty"`
	BgpManualSpokeAdvertiseCidrs string `form:"bgp_manual_spoke,omitempty"`
	EncVolume                    string `form:"enc_volume,omitempty"`
	EnablePrivateOob             string `form:"private_oob,omitempty"`
	OobManagementSubnet          string `form:"oob_mgmt_subnet,omitempty"`
	HAOobManagementSubnet        string
	AvailabilityDomain           string   `form:"availability_domain,omitempty"`
	FaultDomain                  string   `form:"fault_domain,omitempty"`
	EnableSpotInstance           bool     `form:"spot_instance,omitempty"`
	SpotPrice                    string   `form:"spot_price,omitempty"`
	EnableBgp                    string   `form:"enable_bgp"`
	LearnedCidrsApproval         string   `form:"learned_cidrs_approval,omitempty"`
	ApprovedLearnedCidrs         []string `form:"approved_learned_cidrs"`
}

type SpokeGatewayAdvancedConfig struct {
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

type SpokeGatewayAdvancedConfigResp struct {
	Return  bool                                 `json:"return"`
	Results SpokeGatewayAdvancedConfigRespResult `json:"results"`
	Reason  string                               `json:"reason"`
}

type SpokeGatewayAdvancedConfigRespResult struct {
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

func (c *Client) LaunchSpokeVpc(spoke *SpokeVpc) error {
	spoke.CID = c.CID
	spoke.Action = "create_spoke_gw"

	return c.PostAPI(spoke.Action, spoke, BasicCheck)
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
		"action":  "enable_spoke_ha",
		"gw_name": spoke.GwName,
		"eip":     spoke.Eip,
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

	return c.PostAPI(form["action"], form, checkFunc)
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

func (c *Client) EnableSpokeConnectionLearnedCIDRApproval(gwName, connName string) error {
	data := map[string]string{
		"action":          "enable_transit_connection_learned_cidrs_approval",
		"CID":             c.CID,
		"gateway_name":    gwName,
		"connection_name": connName,
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) DisableSpokeConnectionLearnedCIDRApproval(gwName, connName string) error {
	data := map[string]string{
		"action":          "disable_transit_connection_learned_cidrs_approval",
		"CID":             c.CID,
		"gateway_name":    gwName,
		"connection_name": connName,
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) UpdateSpokeConnectionPendingApprovedCidrs(gwName, connName string, approvedCidrs []string) error {
	data := map[string]string{
		"action":                            "update_transit_connection_pending_approved_cidrs",
		"CID":                               c.CID,
		"gateway_name":                      gwName,
		"connection_name":                   connName,
		"connection_approved_learned_cidrs": strings.Join(approvedCidrs, ","),
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

func (c *Client) SetBgpPollingTimeSpoke(spokeGateway *SpokeVpc, newPollingTime string) error {
	action := "change_bgp_polling_time"
	return c.PostAPI(action, struct {
		CID         string `form:"CID"`
		Action      string `form:"action"`
		GatewayName string `form:"gateway_name"`
		PollingTime string `form:"bgp_polling_time"`
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
		"action":       "enable_transit_learned_cidrs_approval",
		"gateway_name": gateway.GwName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) DisableSpokeLearnedCidrsApproval(gateway *SpokeVpc) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "disable_transit_learned_cidrs_approval",
		"gateway_name": gateway.GwName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) UpdateSpokePendingApprovedCidrs(gateway *SpokeVpc) error {
	form := map[string]string{
		"CID":                    c.CID,
		"action":                 "update_transit_pending_approved_cidrs",
		"gateway_name":           gateway.GwName,
		"approved_learned_cidrs": strings.Join(gateway.ApprovedLearnedCidrs, ","),
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
