package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Vpc struct {
	CloudType          int          `form:"cloud_type,omitempty" json:"cloud_type,omitempty"`
	AccountName        string       `form:"account_name,omitempty" json:"account_name,omitempty"`
	Region             string       `form:"region,omitempty" json:"vpc_region,omitempty"`
	Name               string       `form:"pool_name,omitempty" json:"pool_name,omitempty"`
	Cidr               string       `form:"vpc_cidr,omitempty" json:"vpc_cidr,omitempty"`
	AviatrixTransitVpc string       `form:"aviatrix_transit_vpc,omitempty"`
	AviatrixFireNetVpc string       `form:"aviatrix_firenet_vpc,omitempty"`
	VpcID              string       `json:"vpc_list,omitempty"`
	Subnets            []SubnetInfo `form:"subnet_list,omitempty" json:"subnets,omitempty"`
	PublicSubnets      []SubnetInfo
	PrivateSubnets     []SubnetInfo
	PublicRoutesOnly   bool
}

type VpcEdit struct {
	CloudType          int          `form:"cloud_type,omitempty" json:"cloud_type,omitempty"`
	AccountName        string       `form:"account_name,omitempty" json:"account_name,omitempty"`
	Region             string       `form:"region,omitempty" json:"vpc_region,omitempty"`
	Name               string       `form:"pool_name,omitempty" json:"pool_name,omitempty"`
	Cidr               string       `form:"vpc_cidr,omitempty" json:"vpc_cidr,omitempty"`
	AviatrixTransitVpc bool         `json:"avx_transit_vpc,omitempty"`
	AviatrixFireNetVpc bool         `json:"avx_firenet_vpc,omitempty"`
	VpcID              []string     `json:"vpc_list,omitempty"`
	Subnets            []SubnetInfo `json:"subnets,omitempty"`
}

type VpcResp struct {
	Return  bool                  `json:"return"`
	Results AllVpcPoolVpcListResp `json:"results"`
	Reason  string                `json:"reason"`
}

type AllVpcPoolVpcListResp struct {
	AllVpcPoolVpcList []VpcEdit `json:"all_vpc_pool_vpc_list, omitempty"`
}

type SubnetInfo struct {
	Region   string `json:"region, omitempty"`
	Cidr     string `json:"cidr, omitempty"`
	Name     string `json:"name, omitempty"`
	SubnetID string `json:"id, omitempty"`
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

func (c *Client) GetVpc(vpc *Vpc) (*Vpc, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for list_custom_vpcs ") + err.Error())
	}
	listPeerVpcPairs := url.Values{}
	listPeerVpcPairs.Add("CID", c.CID)
	listPeerVpcPairs.Add("action", "list_custom_vpcs")
	Url.RawQuery = listPeerVpcPairs.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return nil, errors.New("HTTP Get list_custom_vpcs failed: " + err.Error())
	}
	var data VpcResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_custom_vpcs failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return nil, errors.New("Rest API list_custom_vpcs Get failed: " + data.Reason)
	}
	allVpcPoolVpcListResp := data.Results.AllVpcPoolVpcList
	for i := range allVpcPoolVpcListResp {
		if allVpcPoolVpcListResp[i].Name == vpc.Name {
			log.Debugf("Found VPC: %#v", allVpcPoolVpcListResp[i])

			vpc.CloudType = allVpcPoolVpcListResp[i].CloudType
			vpc.AccountName = allVpcPoolVpcListResp[i].AccountName
			vpc.Region = allVpcPoolVpcListResp[i].Region
			vpc.Cidr = allVpcPoolVpcListResp[i].Cidr
			if allVpcPoolVpcListResp[i].AviatrixTransitVpc {
				vpc.AviatrixTransitVpc = "yes"
			} else {
				vpc.AviatrixTransitVpc = "no"
			}
			if allVpcPoolVpcListResp[i].AviatrixFireNetVpc {
				vpc.AviatrixFireNetVpc = "yes"
			} else {
				vpc.AviatrixFireNetVpc = "no"
			}
			vpc.VpcID = allVpcPoolVpcListResp[i].VpcID[0]
			vpc.Subnets = allVpcPoolVpcListResp[i].Subnets

			return vpc, nil
		}
	}
	log.Error("VPC not found")
	return nil, ErrNotFound
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
