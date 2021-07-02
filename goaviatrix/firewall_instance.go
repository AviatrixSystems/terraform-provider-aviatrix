package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
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
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return "", errors.New(("url Parsing failed for add_firewall_instance: ") + err.Error())
	}
	addFirewallInstance := url.Values{}
	addFirewallInstance.Add("CID", c.CID)
	addFirewallInstance.Add("action", "add_firewall_instance")
	if firewallInstance.GwName != "" {
		addFirewallInstance.Add("gw_name", firewallInstance.GwName)
	} else {
		addFirewallInstance.Add("vpc_id", firewallInstance.VpcID)
	}
	addFirewallInstance.Add("firewall_name", firewallInstance.FirewallName)
	addFirewallInstance.Add("firewall_image", firewallInstance.FirewallImage)
	addFirewallInstance.Add("firewall_image_version", firewallInstance.FirewallImageVersion)
	addFirewallInstance.Add("firewall_image_id", firewallInstance.FirewallImageId)
	addFirewallInstance.Add("firewall_size", firewallInstance.FirewallSize)
	addFirewallInstance.Add("key_name", firewallInstance.KeyName)
	addFirewallInstance.Add("iam_role", firewallInstance.IamRole)
	addFirewallInstance.Add("bootstrap_bucket_name", firewallInstance.BootstrapBucketName)
	addFirewallInstance.Add("no_associate", strconv.FormatBool(true))
	addFirewallInstance.Add("username", firewallInstance.Username)
	addFirewallInstance.Add("password", firewallInstance.Password)
	if firewallInstance.EgressVpc != "" {
		addFirewallInstance.Add("cloud_type", strconv.Itoa(GCP))
		addFirewallInstance.Add("egress", firewallInstance.EgressSubnet)
		addFirewallInstance.Add("egress_vpc", firewallInstance.EgressVpc)
		addFirewallInstance.Add("management", firewallInstance.ManagementSubnet)
		addFirewallInstance.Add("management_vpc", firewallInstance.ManagementVpc)
		addFirewallInstance.Add("zone", firewallInstance.AvailabilityZone)
	} else {
		addFirewallInstance.Add("egress_subnet", firewallInstance.EgressSubnet)
		addFirewallInstance.Add("management_subnet", firewallInstance.ManagementSubnet)
	}
	if firewallInstance.SshPublicKey != "" {
		addFirewallInstance.Add("ssh_public_key", firewallInstance.SshPublicKey)
	}
	if firewallInstance.BootstrapStorageName != "" {
		addFirewallInstance.Add("bootstrap_storage_name", firewallInstance.BootstrapStorageName)
	}
	if firewallInstance.StorageAccessKey != "" {
		addFirewallInstance.Add("storage_access_key", firewallInstance.StorageAccessKey)
	}
	if firewallInstance.FileShareFolder != "" {
		addFirewallInstance.Add("file_share_folder", firewallInstance.FileShareFolder)
	}
	if firewallInstance.ShareDirectory != "" {
		addFirewallInstance.Add("share_directory", firewallInstance.ShareDirectory)
	}
	if firewallInstance.SicKey != "" {
		addFirewallInstance.Add("sic_key", firewallInstance.SicKey)
	}
	if firewallInstance.ContainerFolder != "" {
		addFirewallInstance.Add("container_folder", firewallInstance.ContainerFolder)
	}
	if firewallInstance.SasUrlConfig != "" {
		addFirewallInstance.Add("sas_url_config", firewallInstance.SasUrlConfig)
	}
	if firewallInstance.SasUriLicense != "" {
		addFirewallInstance.Add("sas_url_license", firewallInstance.SasUriLicense)
	}
	if firewallInstance.UserData != "" {
		addFirewallInstance.Add("user_data", firewallInstance.UserData)
	}
	if len(firewallInstance.Tags) > 0 {
		addFirewallInstance.Add("tag_json", firewallInstance.TagJson)
	}
	if firewallInstance.AvailabilityDomain != "" && firewallInstance.FaultDomain != "" {
		addFirewallInstance.Add("cloud_type", strconv.Itoa(OCI))
		addFirewallInstance.Add("availability_domain", firewallInstance.AvailabilityDomain)
		addFirewallInstance.Add("fault_domain", firewallInstance.FaultDomain)
	}

	Url.RawQuery = addFirewallInstance.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return "", errors.New("HTTP Get add_firewall_instance failed: " + err.Error())
	}

	var data FirewallInstanceCreateResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return "", errors.New("Json Decode add_firewall_instance failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return "", errors.New("Rest API add_firewall_instance Get failed: " + data.Reason)
	}
	if data.Results.FirewallID != "" {
		return data.Results.FirewallID, nil
	}
	return "", ErrNotFound
}

func (c *Client) GetFirewallInstance(firewallInstance *FirewallInstance) (*FirewallInstance, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for get_instance_by_id: ") + err.Error())
	}
	getInstanceById := url.Values{}
	getInstanceById.Add("CID", c.CID)
	getInstanceById.Add("action", "get_instance_by_id")
	getInstanceById.Add("instance_id", firewallInstance.InstanceID)
	Url.RawQuery = getInstanceById.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get get_instance_by_id failed: " + err.Error())
	}
	var data FirewallInstanceResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode get_instance_by_id failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		if strings.Contains(data.Reason, "Unrecognized firewall instance_id") {
			return nil, ErrNotFound
		}
		return nil, errors.New("Rest API get_instance_by_id Get failed: " + data.Reason)
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
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for delete_firenet_firewall_instance ") + err.Error())
	}
	deleteFirenetFirewallInstance := url.Values{}
	deleteFirenetFirewallInstance.Add("CID", c.CID)
	deleteFirenetFirewallInstance.Add("action", "delete_firenet_firewall_instance")
	deleteFirenetFirewallInstance.Add("vpc_id", firewallInstance.VpcID)
	deleteFirenetFirewallInstance.Add("firewall_id", firewallInstance.InstanceID)
	Url.RawQuery = deleteFirenetFirewallInstance.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get delete_firenet_firewall_instance failed: " + err.Error())
	}

	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode delete_firenet_firewall_instance failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API delete_firenet_firewall_instance Get failed: " + data.Reason)
	}
	return nil
}
