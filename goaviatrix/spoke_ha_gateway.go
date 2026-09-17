package goaviatrix

import (
	"fmt"
	"log"
)

type SpokeHaGateway struct {
	Action                string `form:"action" json:"action"`
	CID                   string `form:"CID" json:"CID"`
	AccountName           string `form:"account_name,omitempty" json:"account_name"`
	CloudType             int    `form:"cloud_type,omitempty" json:"cloud_type"`
	GroupUUID             string `form:"group_uuid,omitempty" json:"group_uuid,omitempty"`
	VpcID                 string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	VNetNameResourceGroup string `form:"vnet_and_resource_group_names,omitempty" json:"vnet_and_resource_group_names"`
	PrimaryGwName         string `form:"primary_gw_name,omitempty" json:"primary_gw_name"`
	GwName                string `form:"ha_gw_name,omitempty" json:"ha_gw_name"`
	GwSize                string `form:"gw_size,omitempty" json:"gw_size"`
	Subnet                string `form:"gw_subnet,omitempty" json:"gw_subnet"`
	VpcRegion             string `form:"region,omitempty" json:"region"`
	Zone                  string `form:"zone,omitempty" json:"zone"`
	AvailabilityDomain    string `form:"availability_domain,omitempty" json:"availability_domain"`
	FaultDomain           string `form:"fault_domain,omitempty" json:"fault_domain"`
	BgpLanVpcID           string `form:"bgp_lan_vpc,omitempty" json:"bgp_lan_vpc"`
	BgpLanSubnet          string `form:"bgp_lan_subnet,omitempty" json:"bgp_lan_subnet"`
	Eip                   string `form:"eip,omitempty" json:"eip,omitempty"`
	InsaneMode            string `form:"insane_mode,omitempty" json:"insane_mode"`
	TagList               string `form:"tag_string,omitempty" json:"tag_string"`
	TagJSON               string `form:"tag_json,omitempty" json:"tag_json"`
	AutoGenHaGwName       string `form:"autogen_hagw_name,omitempty" json:"autogen_hagw_name"`
	Async                 bool   `form:"async,omitempty" json:"async"`
	InsertionGateway      bool   `form:"insertion_gateway,omitempty" json:"insertion_gateway,omitempty"`
}

type APIRespHaGw struct {
	Return   bool   `json:"return"`
	Results  string `json:"results"`
	Reason   string `json:"reason"`
	HaGwName string `json:"ha_gw_name"`
}

func (c *Client) CreateSpokeHaGw(spokeHaGateway *SpokeHaGateway) (string, error) {
	spokeHaGateway.CID = c.CID
	spokeHaGateway.Action = "create_multicloud_ha_gateway"
	spokeHaGateway.Async = true // Enable async mode

	// Capture ha_gw_name from the async response using a hook
	var haGwName string
	hook := WithResponseHook(func(raw map[string]interface{}) {
		if name, ok := raw["ha_gw_name"].(string); ok {
			haGwName = name
		}
	})

	err := c.PostAsyncAPI(spokeHaGateway.Action, spokeHaGateway, BasicCheck, hook)
	if err != nil {
		return "", err
	}

	// If async API returned the HA gateway name, use it
	if haGwName != "" {
		log.Printf("[INFO] HA gateway name from async response: %s", haGwName)
		return haGwName, nil
	}

	// If user provided a specific HA gateway name, use it
	if spokeHaGateway.GwName != "" {
		log.Printf("[INFO] Using user-provided HA gateway name: %s", spokeHaGateway.GwName)
		return spokeHaGateway.GwName, nil
	}

	return "", fmt.Errorf("HA gateway name not found")
}
