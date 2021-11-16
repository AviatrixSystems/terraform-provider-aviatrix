package goaviatrix

import (
	"fmt"
	"strings"
)

// Spoke gateway simple struct to hold spoke details
type SpokeVpc struct {
	AccountName           string `form:"account_name,omitempty" json:"account_name,omitempty"`
	Action                string `form:"action,omitempty"`
	CID                   string `form:"CID,omitempty"`
	CloudType             int    `form:"cloud_type,omitempty" json:"cloud_type,omitempty"`
	DnsServer             string `form:"dns_server,omitempty" json:"dns_server,omitempty"`
	GwName                string `form:"gw_name,omitempty" json:"vpc_name,omitempty"`
	GwSize                string `form:"gw_size,omitempty"`
	VpcID                 string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	VNetNameResourceGroup string `form:"vnet_and_resource_group_names,omitempty"`
	Subnet                string `form:"public_subnet,omitempty" json:"public_subnet,omitempty"`
	VpcRegion             string `form:"region,omitempty" json:"vpc_region,omitempty"`
	VpcSize               string `form:"gw_size,omitempty" json:"vpc_size,omitempty"`
	EnableNat             string `form:"nat_enabled,omitempty" json:"enable_nat,omitempty"`
	EnableVpcDnsServer    string `json:"use_vpc_dns,omitempty"`
	HASubnet              string `form:"ha_subnet,omitempty"`
	HAZone                string `form:"new_zone,omitempty"`
	HASubnetGCP           string `form:"new_subnet,omitempty"`
	SingleAzHa            string `form:"single_az_ha,omitempty"`
	TransitGateway        string `form:"transit_gw,omitempty"`
	TagList               string `form:"tag_string,omitempty"`
	TagJson               string `form:"tag_json,omitempty"`
	ReuseEip              string `form:"reuse_eip,omitempty"`
	AllocateNewEipRead    bool   `json:"newly_allocated_eip,omitempty"`
	Eip                   string `form:"eip,omitempty" json:"eip,omitempty"`
	InsaneMode            string `form:"insane_mode,omitempty"`
	Zone                  string `form:"zone,omitempty" json:"zone,omitempty"`
	EncVolume             string `form:"enc_volume,omitempty"`
	EnablePrivateOob      string `form:"private_oob,omitempty"`
	OobManagementSubnet   string `form:"oob_mgmt_subnet,omitempty"`
	HAOobManagementSubnet string
	StorageName           string `form:"storage_name"`
	AvailabilityDomain    string `form:"availability_domain,omitempty"`
	FaultDomain           string `form:"fault_domain,omitempty"`
	EnableSpotInstance    bool   `form:"spot_instance,omitempty"`
	SpotPrice             string `form:"spot_price,omitempty"`
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
