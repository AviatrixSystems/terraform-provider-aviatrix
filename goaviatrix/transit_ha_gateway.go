package goaviatrix

import "golang.org/x/net/context"

type TransitHaGateway struct {
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
	BgpLanVpcId           string `json:"bgp_lan_vpc"`
	BgpLanSubnet          string `json:"bgp_lan_specify_subnet"`
	Eip                   string `json:"eip,omitempty"`
	InsaneMode            string `json:"insane_mode"`
	TagList               string `json:"tag_string"`
	TagJson               string `json:"tag_json"`
	AutoGenHaGwName       string `json:"autogen_hagw_name"`
}

func (c *Client) CreateTransitHaGw(transitHaGateway *TransitHaGateway) (string, error) {
	transitHaGateway.CID = c.CID
	transitHaGateway.Action = "create_multicloud_ha_gateway"

	return c.PostAPIContext2HaGw(context.Background(), nil, transitHaGateway.Action, transitHaGateway, BasicCheck)
}
