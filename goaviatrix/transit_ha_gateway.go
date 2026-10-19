package goaviatrix

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
)

type TransitHaGateway struct {
	Action                    string `form:"action" json:"action"`
	CID                       string `form:"CID" json:"CID"`
	AccountName               string `form:"account_name,omitempty" json:"account_name"`
	CloudType                 int    `form:"cloud_type,omitempty" json:"cloud_type"`
	GroupUUID                 string `form:"group_uuid,omitempty" json:"group_uuid,omitempty"`
	VpcID                     string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	VNetNameResourceGroup     string `form:"vnet_and_resource_group_names,omitempty" json:"vnet_and_resource_group_names"`
	PrimaryGwName             string `form:"primary_gw_name,omitempty" json:"primary_gw_name,omitempty"`
	GwName                    string `form:"ha_gw_name,omitempty" json:"ha_gw_name"`
	GwSize                    string `form:"gw_size,omitempty" json:"gw_size"`
	Subnet                    string `form:"gw_subnet,omitempty" json:"gw_subnet"`
	VpcRegion                 string `form:"region,omitempty" json:"region"`
	Zone                      string `form:"zone,omitempty" json:"zone"`
	AvailabilityDomain        string `form:"availability_domain,omitempty" json:"availability_domain"`
	FaultDomain               string `form:"fault_domain,omitempty" json:"fault_domain"`
	BgpLanVpcID               string `form:"bgp_lan_vpc,omitempty" json:"bgp_lan_vpc"`
	BgpLanSubnet              string `form:"bgp_lan_specify_subnet,omitempty" json:"bgp_lan_specify_subnet"`
	Eip                       string `form:"eip,omitempty" json:"eip,omitempty"`
	InsaneMode                string `form:"insane_mode,omitempty" json:"insane_mode"`
	TagList                   string `form:"tag_string,omitempty" json:"tag_string"`
	TagJSON                   string `form:"tag_json,omitempty" json:"tag_json"`
	AutoGenHaGwName           string `form:"autogen_hagw_name,omitempty" json:"autogen_hagw_name"`
	BackupLinkList            []BackupLinkInterface
	BackupLinkConfig          string `form:"backup_link_config,omitempty" json:"backup_link_config,omitempty"`
	InterfaceMapping          string `form:"interface_mapping,omitempty" json:"interface_mapping,omitempty"`
	Interfaces                string `form:"interfaces,omitempty" json:"interfaces,omitempty"`
	DeviceID                  string `form:"device_id,omitempty" json:"device_id,omitempty"`
	ZtpFileDownloadPath       string `form:"-" json:"-"`
	ZtpFileType               string `form:"ztp_file_type,omitempty" json:"ztp_file_type,omitempty"`
	GatewayRegistrationMethod string `form:"gw_registration_method,omitempty" json:"gw_registration_method,omitempty"`
	ManagementEgressIPPrefix  string `form:"mgmt_egress_ip,omitempty" json:"mgmt_egress_ip,omitempty"`
	Async                     bool   `form:"async,omitempty" json:"async,omitempty"`
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

		// Capture ha_gw_name from the async response using a hook
		var haGwName string
		hook := WithResponseHook(func(raw map[string]interface{}) {
			if name, ok := raw["ha_gw_name"].(string); ok {
				haGwName = name
			}
		})

		err := c.PostAsyncAPI(transitHaGateway.Action, transitHaGateway, BasicCheck, hook)
		if err != nil {
			return "", err
		}
		// If async API returned the HA gateway name, use it
		if haGwName != "" {
			return haGwName, nil
		}
		// If user provided a specific HA gateway name, use it
		if transitHaGateway.GwName != "" {
			return transitHaGateway.GwName, nil
		}
		return "", nil
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
