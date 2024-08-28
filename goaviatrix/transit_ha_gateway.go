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
	Subnet                string `json:"gw_subnet,omitempty"`
	VpcRegion             string `json:"region,omitempty"`
	Zone                  string `json:"zone,omitempty"`
	AvailabilityDomain    string `json:"availability_domain,omitempty"`
	FaultDomain           string `json:"fault_domain,omitempty"`
	BgpLanVpcId           string `json:"bgp_lan_vpc,omitempty"`
	BgpLanSubnet          string `json:"bgp_lan_specify_subnet,omitempty"`
	Eip                   string `json:"eip,omitempty"`
	InsaneMode            string `json:"insane_mode"`
	TagList               string `json:"tag_string,omitempty"`
	TagJson               string `json:"tag_json,omitempty"`
	AutoGenHaGwName       string `json:"autogen_hagw_name,omitempty"`
	BackupLinkList        []BackupLinkInterface
	BackupLinkConfig      string `json:"backup_link_config,omitempty"`
	InterfaceMapping      string `json:"interface_mapping,omitempty"`
	Interfaces            string `json:"interfaces,omitempty"`
	DeviceID              string `json:"device_id,omitempty"`
}

type BackupLinkInterface struct {
	PeerGwName     string `json:"peer_gw_name"`
	PeerBackupPort string `json:"peer_backup_port"`
	SelfBackupPort string `json:"self_backup_port"`
	ConnectionType string `json:"connection_type"`
}

func (c *Client) CreateTransitHaGw(transitHaGateway *TransitHaGateway) (string, error) {
	transitHaGateway.CID = c.CID
	transitHaGateway.Action = "create_multicloud_ha_gateway"

	return c.PostAPIContext2HaGw(context.Background(), nil, transitHaGateway.Action, transitHaGateway, BasicCheck)
}
