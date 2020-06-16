package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"
)

type VpcTracker struct {
	CloudType     int
	AccountName   string
	Region        string
	Name          string
	Cidr          string
	InstanceCount int
	VpcID         string
	Subnets       []VPCTrackerSubnet
}

type VPCTrackerSubnet struct {
	Region    string `json:"region,omitempty"`
	Cidr      string `json:"cidr,omitempty"`
	Name      string `json:"name,omitempty"`
	GatewayIP string `json:"gw_ip,omitempty"`
}

type VPCTrackerItemResp struct {
	VendorName    string             `json:"vendor_name,omitempty"`
	VpcID         string             `json:"vpc_id,omitempty"`
	AccountName   string             `json:"vpc_account_name,omitempty"`
	VpcName       string             `json:"vpc_name,omitempty"`
	Region        string             `json:"region,omitempty"`
	InstanceCount int                `json:"inst_num,omitempty"`
	CIDRs         []string           `json:"vpc_cidrs,omitempty"`
	Subnets       []VPCTrackerSubnet `json:"subnets,omitempty"`
}

type VpcTrackerResp struct {
	Return  bool                 `json:"return"`
	Results []VPCTrackerItemResp `json:"results"`
	Reason  string               `json:"reason"`
}

// GetVpcTracker retrieves the list of VPC's from the 'VPC Tracker' feature.
func (c *Client) GetVpcTracker() ([]*VpcTracker, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New("could not parse url for cloud_network_info: " + err.Error())
	}

	v := url.Values{}
	v.Add("action", "cloud_network_info")
	v.Add("CID", c.CID)
	v.Add("cache", "no")
	v.Add("show_all", "yes")
	Url.RawQuery = v.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get cloud_network_info failed: " + err.Error())
	}

	var data VpcTrackerResp
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Json Decode cloud_network_info failed: %s \n Body: %s", err.Error(), resp.Body))
	}
	if !data.Return {
		return nil, errors.New("Rest API cloud_network_info Get failed: " + data.Reason)
	}

	var vpcList []*VpcTracker
	for _, vpc := range data.Results {
		cidr := ""
		// GCP vpc's will not send any CIDR's
		if len(vpc.CIDRs) > 0 {
			cidr = vpc.CIDRs[0]
		}
		vpcList = append(vpcList, &VpcTracker{
			CloudType:     vendorNameToCloudType(vpc.VendorName),
			VpcID:         vpc.VpcID,
			AccountName:   vpc.AccountName,
			Region:        vpc.Region,
			Name:          vpc.VpcName,
			InstanceCount: vpc.InstanceCount,
			Subnets:       vpc.Subnets,
			Cidr:          cidr,
		})
	}

	return vpcList, nil
}

func vendorNameToCloudType(v string) int {
	vendorToCloud := map[string]int{
		"CLOUD_AWS":     AWS,
		"CLOUD_GOOGLE":  GCP,
		"CLOUD_AZURE":   AZURE,
		"CLOUD_ORACLE":  OCI,
		"CLOUD_AWS_GOV": AWSGOV,
	}
	ct, ok := vendorToCloud[v]
	if !ok {
		log.Errorf("could not map vendor name to cloud type with vendor=%s", v)
		return 0
	}
	return ct
}
