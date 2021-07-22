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

type Vpc struct {
	CloudType              int          `form:"cloud_type,omitempty" json:"cloud_type,omitempty"`
	AccountName            string       `form:"account_name,omitempty" json:"account_name,omitempty"`
	Region                 string       `form:"region,omitempty" json:"vpc_region,omitempty"`
	Name                   string       `form:"pool_name,omitempty" json:"pool_name,omitempty"`
	Cidr                   string       `form:"vpc_cidr,omitempty" json:"vpc_cidr,omitempty"`
	SubnetSize             int          `form:"num_of_subnets,omitempty"`
	NumOfSubnetPairs       int          `form:"num_of_zone,omitempty"`
	EnablePrivateOobSubnet bool         `form:"private_oob_subnet,omitempty"`
	AviatrixTransitVpc     string       `form:"aviatrix_transit_vpc,omitempty"`
	AviatrixFireNetVpc     string       `form:"aviatrix_firenet_vpc,omitempty"`
	VpcID                  string       `json:"vpc_list,omitempty"`
	Subnets                []SubnetInfo `form:"subnet_list,omitempty" json:"subnets,omitempty"`
	PublicSubnets          []SubnetInfo
	PrivateSubnets         []SubnetInfo
	PublicRoutesOnly       bool
	ResourceGroup          string `json:"resource_group,omitempty"`
}

type VpcEdit struct {
	CloudType              int          `form:"cloud_type,omitempty" json:"cloud_type,omitempty"`
	AccountName            string       `form:"account_name,omitempty" json:"account_name,omitempty"`
	Region                 string       `form:"region,omitempty" json:"vpc_region,omitempty"`
	Name                   string       `form:"pool_name,omitempty" json:"pool_name,omitempty"`
	Cidr                   string       `form:"vpc_cidr,omitempty" json:"vpc_cidr,omitempty"`
	SubnetSize             int          `json:"subnet_size,omitempty"`
	NumOfSubnetPairs       int          `json:"num_of_subnet_pairs,omitempty"`
	EnablePrivateOobSubnet bool         `json:"private_oob_subnets,omitempty"`
	AviatrixTransitVpc     bool         `json:"avx_transit_vpc,omitempty"`
	AviatrixFireNetVpc     bool         `json:"avx_firenet_vpc,omitempty"`
	VpcID                  []string     `json:"vpc_list,omitempty"`
	Subnets                []SubnetInfo `json:"subnets,omitempty"`
	PublicSubnets          []SubnetInfo `json:"public_subnets,omitempty"`
	PrivateSubnets         []SubnetInfo `json:"private_subnets,omitempty"`
}

type VpcResp struct {
	Return  bool                  `json:"return"`
	Results AllVpcPoolVpcListResp `json:"results"`
	Reason  string                `json:"reason"`
}

type GetVpcByNameResp struct {
	Return  bool    `json:"return"`
	Results VpcEdit `json:"results"`
	Reason  string  `json:"reason"`
}

type AllVpcPoolVpcListResp struct {
	AllVpcPoolVpcList []VpcEdit `json:"all_vpc_pool_vpc_list,omitempty"`
}

type SubnetInfo struct {
	Region   string `json:"region,omitempty"`
	Cidr     string `json:"cidr,omitempty"`
	Name     string `json:"name,omitempty"`
	SubnetID string `json:"id,omitempty"`
}

func (c *Client) CreateVpc(vpc *Vpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for create_custom_vpc ") + err.Error())
	}
	createCustomVpc := url.Values{}
	createCustomVpc.Add("CID", c.CID)
	createCustomVpc.Add("action", "create_custom_vpc")
	createCustomVpc.Add("cloud_type", strconv.Itoa(vpc.CloudType))
	createCustomVpc.Add("account_name", vpc.AccountName)
	createCustomVpc.Add("pool_name", vpc.Name)
	if vpc.CloudType != GCP {
		createCustomVpc.Add("region", vpc.Region)
		createCustomVpc.Add("vpc_cidr", vpc.Cidr)
		createCustomVpc.Add("aviatrix_transit_vpc", vpc.AviatrixTransitVpc)
		createCustomVpc.Add("aviatrix_firenet_vpc", vpc.AviatrixFireNetVpc)
	} else {
		if vpc.Subnets != nil && len(vpc.Subnets) != 0 {
			i := 0
			for _, subnetInfo := range vpc.Subnets {
				createCustomVpc.Add("subnet_list["+strconv.Itoa(i)+"][name]", subnetInfo.Name)
				createCustomVpc.Add("subnet_list["+strconv.Itoa(i)+"][region]", subnetInfo.Region)
				createCustomVpc.Add("subnet_list["+strconv.Itoa(i)+"][cidr]", subnetInfo.Cidr)
				i++
			}
		}
	}
	if vpc.SubnetSize != 0 {
		createCustomVpc.Add("subnet_size", strconv.Itoa(vpc.SubnetSize))
	}
	if vpc.NumOfSubnetPairs != 0 {
		if IsCloudType(vpc.CloudType, AWSRelatedCloudTypes) {
			createCustomVpc.Add("num_of_zones", strconv.Itoa(vpc.NumOfSubnetPairs))
		} else if IsCloudType(vpc.CloudType, AzureArmRelatedCloudTypes) {
			createCustomVpc.Add("num_of_subnets", strconv.Itoa(vpc.NumOfSubnetPairs))
		}
	}
	if vpc.EnablePrivateOobSubnet {
		createCustomVpc.Add("private_oob_subnet", "true")
	}
	if vpc.ResourceGroup != "" {
		createCustomVpc.Add("resource_group", vpc.ResourceGroup)
	}
	Url.RawQuery = createCustomVpc.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get create_custom_vpc failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode create_custom_vpc failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API create_custom_vpc Get failed: " + data.Reason)
	}
	return nil
}

// GetVpcCloudTypeById returns the cloud_type of the vpc with the given ID.
// If the vpc does not exist, ErrNotFound is returned.
func (c *Client) GetVpcCloudTypeById(ID string) (int, error) {
	action := "list_custom_vpcs"
	d := map[string]string{
		"CID":    c.CID,
		"action": action,
	}
	resp, err := c.Post(c.baseURL, d)

	if err != nil {
		return 0, fmt.Errorf("HTTP Get %s failed: %v", action, err)
	}
	var data VpcResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return 0, fmt.Errorf("Json Decode %s failed: %v\n Body: %s", action, err, bodyString)
	}
	if !data.Return {
		return 0, fmt.Errorf("rest API %s Get failed: %s", action, data.Reason)
	}
	allVpcPoolVpcListResp := data.Results.AllVpcPoolVpcList
	for _, vpcPool := range allVpcPoolVpcListResp {
		for _, vpcID := range vpcPool.VpcID {
			if vpcID == ID {
				return vpcPool.CloudType, nil
			}
		}
	}
	return 0, ErrNotFound
}

func (c *Client) GetCloudTypeFromVpcID(vpcID string) (int, error) {
	data := map[string]string{
		"action": "list_custom_vpcs",
		"CID":    c.CID,
	}
	var respData VpcResp
	err := c.GetAPI(&respData, data["action"], data, BasicCheck)
	if err != nil {
		return 0, err
	}
	for _, vpcPool := range respData.Results.AllVpcPoolVpcList {
		for _, id := range vpcPool.VpcID {
			if id == vpcID {
				return vpcPool.CloudType, nil
			}
		}
	}
	return 0, ErrNotFound
}

func (c *Client) GetVpc(vpc *Vpc) (*Vpc, error) {
	form := map[string]string{
		"action":   "get_custom_vpc_by_name",
		"CID":      c.CID,
		"vpc_name": vpc.Name,
	}
	var data GetVpcByNameResp
	check := func(action, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "does not exist") {
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s Post failed: %s", action, reason)
		}
		return nil
	}
	err := c.GetAPI(&data, form["action"], form, check)
	if err != nil {
		return nil, err
	}
	vpc.CloudType = data.Results.CloudType
	vpc.AccountName = data.Results.AccountName
	vpc.Region = data.Results.Region
	vpc.Cidr = data.Results.Cidr
	if data.Results.AviatrixTransitVpc {
		vpc.AviatrixTransitVpc = "yes"
	} else {
		vpc.AviatrixTransitVpc = "no"
	}
	if data.Results.AviatrixFireNetVpc {
		vpc.AviatrixFireNetVpc = "yes"
	} else {
		vpc.AviatrixFireNetVpc = "no"
	}
	vpc.VpcID = data.Results.VpcID[0]
	vpc.Subnets = data.Results.Subnets
	vpc.PrivateSubnets = data.Results.PrivateSubnets
	vpc.PublicSubnets = data.Results.PublicSubnets
	vpc.SubnetSize = data.Results.SubnetSize
	vpc.NumOfSubnetPairs = data.Results.NumOfSubnetPairs
	vpc.EnablePrivateOobSubnet = data.Results.EnablePrivateOobSubnet
	return vpc, nil
}

func (c *Client) GetVpcRouteTableIDs(vpc *Vpc) ([]string, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for list_vpc_route_tables ") + err.Error())
	}
	listRouteTables := url.Values{}
	listRouteTables.Add("CID", c.CID)
	listRouteTables.Add("action", "list_vpc_route_tables")
	listRouteTables.Add("vpc_id", vpc.VpcID)
	listRouteTables.Add("account_name", vpc.AccountName)
	listRouteTables.Add("vpc_region", vpc.Region)
	if vpc.PublicRoutesOnly {
		listRouteTables.Add("public_only", "yes")
	}

	Url.RawQuery = listRouteTables.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return nil, errors.New("HTTP Get list_vpc_route_tables failed: " + err.Error())
	}

	type RespResults struct {
		RouteTables []string `json:"vpc_rtbs_list"`
	}
	type Resp struct {
		Return  bool        `json:"return"`
		Results RespResults `json:"results"`
		Reason  string      `json:"reason"`
	}
	var data Resp

	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_vpc_route_tables failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return nil, errors.New("Rest API list_vpc_route_tables Get failed: " + data.Reason)
	}

	var rtbs []string
	for _, id := range data.Results.RouteTables {
		rtbs = append(rtbs, strings.Split(id, "~~")[0])
	}

	return rtbs, nil
}

func (c *Client) UpdateVpc(vpc *Vpc) error {
	return nil
}

func (c *Client) DeleteVpc(vpc *Vpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for delete_custom_vpc ") + err.Error())
	}
	createCustomVpc := url.Values{}
	createCustomVpc.Add("CID", c.CID)
	createCustomVpc.Add("action", "delete_custom_vpc")
	createCustomVpc.Add("account_name", vpc.AccountName)
	createCustomVpc.Add("pool_name", vpc.Name)
	Url.RawQuery = createCustomVpc.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get delete_custom_vpc failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode delete_custom_vpc failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API delete_custom_vpc Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) EnableNativeAwsGwlbFirenet(vpc *Vpc) error {
	data := map[string]string{
		"action":       "enable_native_aws_gwlb_firenet",
		"CID":          c.CID,
		"account_name": vpc.AccountName,
		"region":       vpc.Region,
		"vpc_id":       vpc.VpcID,
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) DisableNativeAwsGwlbFirenet(vpc *Vpc) error {
	data := map[string]string{
		"action": "disable_native_aws_gwlb_firenet",
		"CID":    c.CID,
		"vpc_id": vpc.VpcID,
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) ListOciVpcAvailabilityDomains(vpc *Vpc) ([]string, error) {
	params := map[string]string{
		"action":       "list_oci_vpc_availability_domains",
		"CID":          c.CID,
		"account_name": vpc.AccountName,
		"region":       vpc.Region,
		"vpc_id":       vpc.VpcID,
	}

	type Resp struct {
		Return  bool     `json:"return"`
		Results []string `json:"results"`
	}

	var data Resp

	err := c.GetAPI(&data, params["action"], params, BasicCheck)

	if err != nil {
		return nil, err
	}

	return data.Results, nil
}

func (c *Client) ListOciVpcFaultDomains(vpc *Vpc) ([]string, error) {
	params := map[string]string{
		"action":       "list_oci_vpc_fault_domains",
		"CID":          c.CID,
		"account_name": vpc.AccountName,
		"region":       vpc.Region,
		"vpc_id":       vpc.VpcID,
	}

	type Resp struct {
		Return  bool     `json:"return"`
		Results []string `json:"results"`
	}

	var data Resp

	err := c.GetAPI(&data, params["action"], params, BasicCheck)

	if err != nil {
		return nil, err
	}

	return data.Results, nil
}
