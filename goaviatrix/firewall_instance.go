package goaviatrix

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type FirewallInstance struct {
	CID    string `form:"CID,omitempty"`
	Action string `form:"action,omitempty"`

	VpcID                string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	FirenetVpc           string `json:"firenet_vpc"`
	GwName               string `form:"gw_name,omitempty" json:"gw_name,omitempty"`
	FirewallName         string `form:"firewall_name,omitempty" json:"instance_name,omitempty"`
	FirewallImage        string `form:"firewall_image,omitempty" json:"firewall_image,omitempty"`
	FirewallImageVersion string `form:"firewall_image_version,omitempty" json:"firewall_image_version,omitempty"`
	FirewallImageId      string `form:"firewall_image_id,omitempty" json:"firewall_image_id"`
	FirewallSize         string `form:"firewall_size,omitempty" json:"instance_size,omitempty"`
	EgressSubnet         string `form:"egress_subnet,omitempty" json:"egress_subnet,omitempty"`
	ManagementSubnet     string `form:"management_subnet,omitempty" json:"management_subnet,omitempty"`
	KeyName              string `form:"key_name,omitempty" json:"key_name,omitempty"`
	KeyFile              string `json:"key_file"`
	IamRole              string `form:"iam_role,omitempty" json:"iam_role,omitempty"`
	BootstrapBucketName  string `form:"bootstrap_bucket_name,omitempty" json:"bootstrap_bucket_name,omitempty"`
	InstanceID           string `form:"firewall_id,omitempty" json:"instance_id,omitempty"`
	Attached             bool
	LanInterface         string `form:"lan_interface,omitempty" json:"lan_interface_id,omitempty"`
	ManagementInterface  string `form:"management_interface,omitempty" json:"management_interface_id,omitempty"`
	ManagementVpc        string `form:"management_vpc,omitempty" json:"management_vpc"`
	ManagementSubnetID   string `json:"management_subnet_id"`
	EgressInterface      string `form:"egress_interface,omitempty" json:"egress_interface_id,omitempty"`
	EgressVpc            string `form:"egress_vpc,omitempty" json:"egress_vpc"`
	EgressSubnetID       string `json:"egress_subnet_id"`
	ManagementPublicIP   string `json:"management_public_ip,omitempty"`
	VendorType           string
	Username             string          `form:"username,omitempty"`
	Password             string          `form:"password,omitempty"`
	AvailabilityZone     string          `form:"zone,omitempty" json:"availability_zone,omitempty"`
	CloudVendor          string          `json:"cloud_vendor,omitempty"`
	SshPublicKey         string          `form:"ssh_public_key,omitempty" json:"ssh_public_key,omitempty"`
	BootstrapStorageName string          `form:"bootstrap_storage_name,omitempty" json:"bootstrap_storage_name,omitempty"`
	StorageAccessKey     string          `form:"storage_access_key,omitempty" json:"storage_access_key,omitempty"`
	FileShareFolder      string          `form:"file_share_folder,omitempty" json:"file_share_folder,omitempty"`
	ShareDirectory       string          `form:"share_directory,omitempty" json:"share_directory,omitempty"`
	SicKey               string          `form:"sic_key,omitempty" json:"sic_key,omitempty"`
	ContainerFolder      string          `form:"container_folder,omitempty" json:"container_folder,omitempty"`
	SasUrlConfig         string          `form:"sas_url_config,omitempty" json:"sas_url_config,omitempty"`
	SasUriLicense        string          `form:"sas_url_license,omitempty" json:"sas_url_license,omitempty"`
	UserData             string          `form:"user_data,omitempty" json:"user_data,omitempty"`
	TagsMessage          json.RawMessage `json:"usr_tags"`
	Tags                 map[string]string
	TagJson              string
	AvailabilityDomain   string `form:"availability_domain,omitempty"`
	FaultDomain          string `form:"fault_domain,omitempty" json:"fault_domain"`
}

type FirewallInstanceResp struct {
	Return  bool             `json:"return"`
	Results FirewallInstance `json:"results"`
	Reason  string           `json:"reason"`
}

type FirewallInstanceCreateResp struct {
	Return  bool                         `json:"return"`
	Results FirewallInstanceCreateResult `json:"results"`
	Reason  string                       `json:"reason"`
}

type FirewallInstanceCreateResult struct {
	Text       string `json:"text,omitempty"`
	FirewallID string `json:"firewall_id,omitempty"`
}

func (c *Client) CreateFirewallInstance(firewallInstance *FirewallInstance) (string, error) {
	form := map[string]string{
		"CID":                    c.CID,
		"action":                 "add_firewall_instance",
		"firewall_name":          firewallInstance.FirewallName,
		"firewall_image":         firewallInstance.FirewallImage,
		"firewall_image_version": firewallInstance.FirewallImageVersion,
		"firewall_image_id":      firewallInstance.FirewallImageId,
		"firewall_size":          firewallInstance.FirewallSize,
		"key_name":               firewallInstance.KeyName,
		"iam_role":               firewallInstance.IamRole,
		"bootstrap_bucket_name":  firewallInstance.BootstrapBucketName,
		"no_associate":           strconv.FormatBool(true),
		"username":               firewallInstance.Username,
		"password":               firewallInstance.Password,
	}

	if firewallInstance.GwName != "" {
		form["gw_name"] = firewallInstance.GwName
	} else {
		form["vpc_id"] = firewallInstance.VpcID
	}
	if firewallInstance.EgressVpc != "" {
		form["cloud_type"] = strconv.Itoa(GCP)
		form["egress"] = firewallInstance.EgressSubnet
		form["egress_vpc"] = firewallInstance.EgressVpc
		form["management"] = firewallInstance.ManagementSubnet
		form["management_vpc"] = firewallInstance.ManagementVpc
		form["zone"] = firewallInstance.AvailabilityZone
	} else {
		form["egress_subnet"] = firewallInstance.EgressSubnet
		form["management_subnet"] = firewallInstance.ManagementSubnet
	}
	if firewallInstance.SshPublicKey != "" {
		form["ssh_public_key"] = firewallInstance.SshPublicKey
	}
	if firewallInstance.BootstrapStorageName != "" {
		form["bootstrap_storage_name"] = firewallInstance.BootstrapStorageName
	}
	if firewallInstance.StorageAccessKey != "" {
		form["storage_access_key"] = firewallInstance.StorageAccessKey
	}
	if firewallInstance.FileShareFolder != "" {
		form["file_share_folder"] = firewallInstance.FileShareFolder
	}
	if firewallInstance.ShareDirectory != "" {
		form["share_directory"] = firewallInstance.ShareDirectory
	}
	if firewallInstance.SicKey != "" {
		form["sic_key"] = firewallInstance.SicKey
	}
	if firewallInstance.ContainerFolder != "" {
		form["container_folder"] = firewallInstance.ContainerFolder
	}
	if firewallInstance.SasUrlConfig != "" {
		form["sas_url_config"] = firewallInstance.SasUrlConfig
	}
	if firewallInstance.SasUriLicense != "" {
		form["sas_url_license"] = firewallInstance.SasUriLicense
	}
	if firewallInstance.UserData != "" {
		form["user_data"] = firewallInstance.UserData
	}
	if len(firewallInstance.Tags) > 0 {
		form["tag_json"] = firewallInstance.TagJson
	}
	if firewallInstance.AvailabilityDomain != "" && firewallInstance.FaultDomain != "" {
		form["cloud_type"] = strconv.Itoa(OCI)
		form["availability_domain"] = firewallInstance.AvailabilityDomain
		form["fault_domain"] = firewallInstance.FaultDomain
	}

	var data FirewallInstanceCreateResp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return "", err
	}

	if data.Results.FirewallID != "" {
		return data.Results.FirewallID, nil
	}

	return "", ErrNotFound
}

func (c *Client) GetFirewallInstance(firewallInstance *FirewallInstance) (*FirewallInstance, error) {
	form := map[string]string{
		"CID":         c.CID,
		"action":      "get_instance_by_id",
		"instance_id": firewallInstance.InstanceID,
	}

	var data FirewallInstanceResp

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "Unrecognized firewall instance_id") {
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	err := c.GetAPI(&data, form["action"], form, checkFunc)
	if err != nil {
		return nil, err
	}

	if data.Results.InstanceID == firewallInstance.InstanceID {
		// Only try to decode if tags are not empty string
		if string(data.Results.TagsMessage) != `""` {
			err = json.Unmarshal(data.Results.TagsMessage, &data.Results.Tags)
			if err != nil {
				return nil, fmt.Errorf("json Decode get_instance_by_id failed: %v", err)
			}
		}
		return &data.Results, nil
	}

	return nil, ErrNotFound
}

func (c *Client) DeleteFirewallInstance(firewallInstance *FirewallInstance) error {
	form := map[string]string{
		"CID":    c.CID,
		"action": "delete_firenet_firewall_instance",
		"vpc_id": firewallInstance.VpcID,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}
