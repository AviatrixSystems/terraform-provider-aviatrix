package goaviatrix

import (
	"fmt"
	"log"
)

type SpokeHaGateway struct {
	Action                string `json:"action"`
	CID                   string `json:"CID"`
	AccountName           string `json:"account_name"`
	CloudType             int    `json:"cloud_type"`
	VpcID                 string `json:"vpc_id,omitempty"`
	VNetNameResourceGroup string `json:"vnet_and_resource_group_names"`
	PrimaryGwName         string `json:"primary_gw_name"`
	GwName                string `json:"ha_gw_name"`
	GwSize                string `json:"gw_size"`
	Subnet                string `json:"gw_subnet"`
	VpcRegion             string `json:"region"`
	Zone                  string `json:"zone"`
	AvailabilityDomain    string `json:"availability_domain"`
	FaultDomain           string `json:"fault_domain"`
	BgpLanVpcId           string `json:"bpg_lan_vpc_id"`
	BgpLanSubnet          string `json:"bgp_lan_subnet"`
	Eip                   string `json:"eip,omitempty"`
	InsaneMode            string `json:"insane_mode"`
	TagList               string `json:"tag_string"`
	TagJson               string `json:"tag_json"`
	AutoGenHaGwName       string `json:"autogen_hagw_name"`
	Async                 bool   `json:"async"`
	InsertionGateway      bool   `json:"insertion_gateway,omitempty"`
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
