package goaviatrix

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

	// Use async API internally but maintain the same external interface
	err := c.PostAsyncAPI(spokeHaGateway.Action, spokeHaGateway, BasicCheck)
	if err != nil {
		return "", err
	}

	// Determine the gateway name for the return value
	// This follows the same logic as the original synchronous implementation
	gwName := spokeHaGateway.GwName
	if gwName == "" {
		// When AutoGenHaGwName is "yes", the controller generates the name
		// following the pattern: primary_gateway_name + "-hagw"
		gwName = spokeHaGateway.PrimaryGwName + "-hagw"
	}

	return gwName, nil
}
