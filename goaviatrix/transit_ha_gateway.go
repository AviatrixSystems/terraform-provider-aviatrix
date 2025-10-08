package goaviatrix

import (
	"encoding/base64"
	"fmt"
	"os"

	"golang.org/x/net/context"
)

type TransitHaGateway struct {
	Action                    string `json:"action"`
	CID                       string `json:"CID"`
	AccountName               string `json:"account_name"`
	CloudType                 int    `json:"cloud_type"`
	VpcID                     string `json:"vpc_id,omitempty"`
	VNetNameResourceGroup     string `json:"vnet_and_resource_group_names"`
	PrimaryGwName             string `json:"primary_gw_name"`
	GwName                    string `json:"ha_gw_name"`
	GwSize                    string `json:"gw_size"`
	Subnet                    string `json:"gw_subnet"`
	VpcRegion                 string `json:"region"`
	Zone                      string `json:"zone"`
	AvailabilityDomain        string `json:"availability_domain"`
	FaultDomain               string `json:"fault_domain"`
	BgpLanVpcID               string `json:"bgp_lan_vpc"`
	BgpLanSubnet              string `json:"bgp_lan_specify_subnet"`
	Eip                       string `json:"eip,omitempty"`
	InsaneMode                string `json:"insane_mode"`
	TagList                   string `json:"tag_string"`
	TagJSON                   string `json:"tag_json"`
	AutoGenHaGwName           string `json:"autogen_hagw_name"`
	BackupLinkList            []BackupLinkInterface
	BackupLinkConfig          string `json:"backup_link_config,omitempty"`
	InterfaceMapping          string `json:"interface_mapping,omitempty"`
	Interfaces                string `json:"interfaces,omitempty"`
	DeviceID                  string `json:"device_id,omitempty"`
	ZtpFileDownloadPath       string `json:"-"`
	ZtpFileType               string `json:"ztp_file_type,omitempty"`
	GatewayRegistrationMethod string `json:"gw_registration_method,omitempty"`
	ManagementEgressIPPrefix  string `json:"mgmt_egress_ip,omitempty"`
	Async                     bool   `json:"async,omitempty"`
}

type BackupLinkInterface struct {
	PeerGwName               string   `json:"peer_gw_name"`
	PeerBackupPort           string   `json:"peer_backup_port,omitempty"`
	SelfBackupPort           string   `json:"self_backup_port,omitempty"`
	ConnectionType           string   `json:"connection_type"`
	PeerBackupLogicalIfNames []string `json:"peer_backup_logical_ifnames,omitempty"`
	SelfBackupLogicalIfNames []string `json:"self_backup_logical_ifnames,omitempty"`
}

func (c *Client) CreateTransitHaGw(transitHaGateway *TransitHaGateway) (string, error) {
	transitHaGateway.CID = c.CID
	transitHaGateway.Action = "create_multicloud_ha_gateway"
	var data CreateEdgeEquinixResp
	var resp string
	if IsCloudType(transitHaGateway.CloudType, EdgeRelatedCloudTypes) {
		var err error
		resp, err = c.PostAPIContext2HaGw(context.Background(), &data, transitHaGateway.Action, transitHaGateway, BasicCheck)
		if err != nil {
			return "", err
		}
	} else {
		transitHaGateway.Async = true
		err := c.PostAsyncAPI(transitHaGateway.Action, transitHaGateway, BasicCheck)
		return "", err
	}

	// create the ZTP file for Equinix Edge transit gateway
	if transitHaGateway.CloudType == EDGEEQUINIX || transitHaGateway.CloudType == EDGEMEGAPORT {
		fileName := getFileName(transitHaGateway.ZtpFileDownloadPath, transitHaGateway.GwName, transitHaGateway.VpcID)
		err := createZtpFile(fileName, data.Result)
		if err != nil {
			return "", err
		}
	}

	if IsCloudType(transitHaGateway.CloudType, EDGESELFMANAGED) {
		// log the ztp file type
		var fileName string
		if transitHaGateway.ZtpFileType == "iso" {
			fileName = transitHaGateway.ZtpFileDownloadPath + "/" + transitHaGateway.GwName + "-" + transitHaGateway.VpcID + ".iso"

			// Decode base64 content (the data.Result should contain base64-encoded ISO data)
			decodedBytes, err := base64.StdEncoding.DecodeString(data.Result)
			if err != nil {
				return "", fmt.Errorf("failed to decode base64 content for ISO file: %w", err)
			}

			// Create and write the binary ISO file
			outFile, err := os.Create(fileName)
			if err != nil {
				return "", fmt.Errorf("failed to create ISO file %s: %w", fileName, err)
			}
			defer outFile.Close()

			// Write the decoded binary content to the file
			_, err = outFile.Write(decodedBytes)
			if err != nil {
				return "", fmt.Errorf("failed to write binary content to ISO file %s: %w", fileName, err)
			}

			fmt.Printf("[DEBUG] CreateTransitHaGw: Successfully wrote %d bytes (decoded from %d base64 chars) to %s\n",
				len(decodedBytes), len(data.Result), fileName)
		} else {
			fileName = getFileName(transitHaGateway.ZtpFileDownloadPath, transitHaGateway.GwName, transitHaGateway.VpcID)

			fileContent := data.Result
			err := createZtpFile(fileName, fileContent)
			if err != nil {
				return "", err
			}
		}
	}
	return resp, nil
}
