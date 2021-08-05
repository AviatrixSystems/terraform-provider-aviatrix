package goaviatrix

import (
	"math"

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
	InstanceCount interface{}        `json:"inst_num,omitempty"`
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
	form := map[string]string{
		"CID":      c.CID,
		"action":   "cloud_network_info",
		"cache":    "no",
		"show_all": "yes",
	}

	var data VpcTrackerResp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	var vpcList []*VpcTracker
	for _, vpc := range data.Results {
		cidr := ""
		// GCP vpc's will not send any CIDR's
		if len(vpc.CIDRs) > 0 {
			cidr = vpc.CIDRs[0]
		}

		actualInstCount := 0
		instCount, ok := vpc.InstanceCount.(float64)
		if ok {
			instCount = math.Round(instCount)
			actualInstCount = int(instCount)
		}

		vpcList = append(vpcList, &VpcTracker{
			CloudType:     vendorNameToCloudType(vpc.VendorName),
			VpcID:         vpc.VpcID,
			AccountName:   vpc.AccountName,
			Region:        vpc.Region,
			Name:          vpc.VpcName,
			InstanceCount: actualInstCount,
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
		"CLOUD_AZURE":   Azure,
		"CLOUD_ORACLE":  OCI,
		"CLOUD_AWS_GOV": AWSGov,
	}
	ct, ok := vendorToCloud[v]
	if !ok {
		log.Errorf("could not map vendor name to cloud type with vendor=%s", v)
		return 0
	}
	return ct
}
