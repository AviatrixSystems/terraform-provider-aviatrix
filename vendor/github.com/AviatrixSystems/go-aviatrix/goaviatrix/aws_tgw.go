package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
)

// AwsTGW simple struct to hold aws_tgw details
type AWSTgw struct {
	Action                    string               `form:"action,omitempty"`
	CID                       string               `form:"CID,omitempty"`
	Name                      string               `form:"tgw_name,omitempty"`
	AccountName               string               `form:"account_name,omitempty"`
	Region                    string               `form:"region,omitempty"`
	AwsSideAsNumber           string               `form:"aws_side_asn,omitempty"`
	AttachedAviatrixTransitGW []string             `form:"attached_aviatrix_transit_gateway,omitempty"`
	SecurityDomains           []SecurityDomainRule `form:"security_domains,omitempty"`
}

type AWSTgwAPIResp struct {
	Return  bool     `json:"return"`
	Results []string `json:"results"`
	Reason  string   `json:"reason"`
}

type AWSTgwList struct {
	Return  bool     `json:"return"`
	Results []AWSTgw `json:"results"`
	Reason  string   `json:"reason"`
}

type RouteDomainAPIResp struct {
	Return  bool                `json:"return"`
	Results []RouteDomainDetail `json:"results"`
	Reason  string              `json:"reason"`
}

type RouteDomainDetail struct {
	Associations         []string             `json:"associations"`
	Name                 string               `json:"name"`
	ConnectedRouteDomain []string             `json:"connected_route_domain"`
	AttachedVPC          []AttachedVPCDetail  `json:"attached_vpc"`
	RoutesInRouteTable   []RoutesInRouteTable `json:"routes_in_route_table"`
	RouteTableId         string               `json:"route_table_id"`
}

type AttachedVPCDetail struct {
	TgwName      string   `json:"tgw_name"`
	Region       string   `json:"region"`
	VPCName      string   `json:"vpc_name"`
	AttachmentId string   `json:"attachment_id"`
	RouteDomain  string   `json:"route_domain"`
	VPCCidr      []string `json:"vpc_cidr"`
	VPCId        string   `json:"vpc_id"`
	AccountName  string   `json:"account_name"`
}

type RoutesInRouteTable struct {
	VPCId           string `json:"vpc_id"`
	CidrBlock       string `json:"cidr_block"`
	Type            string `json:"type"`
	State           string `json:"state"`
	TgwAttachmentId string `json:"tgw_attachment_id"`
}

type VPCList struct {
	Return  bool      `json:"return"`
	Results []VPCInfo `json:"results"`
	Reason  string    `json:"reason"`
}

type VPCInfo struct {
	AccountName string `json:"account_name,omitempty"`
	CloudType   int    `json:"cloud_type,omitempty"`
	Region      string `json:"vpc_region,omitempty"`
	Name        string `json:"vpc_name,omitempty"`
	TransitVpc  string `json:"transit_vpc,omitempty"`
	VPCId       string `json:"vpc_id,omitempty"`
}

func (c *Client) CreateAWSTgw(awsTgw *AWSTgw) error {
	awsTgw.CID = c.CID
	awsTgw.Action = "add_aws_tgw"
	resp, err := c.Post(c.baseURL, awsTgw)
	if err != nil {
		return err
	}

	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		return errors.New(data.Reason)
	}

	return nil
}

func (c *Client) GetAWSTgw(awsTgw *AWSTgw) (*AWSTgw, error) {
	awsTgw.CID = c.CID
	path := c.baseURL + fmt.Sprintf("?action=list_route_domain_names&tgw_name=%s&CID=%s", awsTgw.Name,
		awsTgw.CID)

	resp, err := c.Get(path, nil)

	if err != nil {
		return nil, err
	}

	data := AWSTgwAPIResp{
		Return:  false,
		Results: make([]string, 0),
		Reason:  "",
	}
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if !data.Return {
		return nil, errors.New(data.Reason)
	}

	connectedDomainList := data.Results
	connectedDomainList = append([]string{"Aviatrix_Edge_Domain"}, connectedDomainList...)

	for i := range connectedDomainList {
		dm := connectedDomainList[i]

		path = c.baseURL + fmt.Sprintf("?action=view_route_domain_details&CID=%s&tgw_name=%s"+
			"&route_domain_name=%s", c.CID, awsTgw.Name, dm)
		resp, err = c.Get(path, nil)
		if err != nil {
			return nil, err
		}

		var data1 RouteDomainAPIResp
		if err = json.NewDecoder(resp.Body).Decode(&data1); err != nil {
			return nil, err
		}
		if !data1.Return {
			return nil, errors.New(data1.Reason)
		}
		routeDomainDetail := data1.Results

		sdr := SecurityDomainRule{
			Name: routeDomainDetail[0].Name,
		}
		for i := range routeDomainDetail[0].ConnectedRouteDomain {
			sdr.ConnectedDomain = append(sdr.ConnectedDomain, routeDomainDetail[0].ConnectedRouteDomain[i])
		}

		attachedVPCs := routeDomainDetail[0].AttachedVPC
		for i := range attachedVPCs {

			if dm != "Aviatrix_Edge_Domain" {
				vpcSolo := VPCSolo{
					Region:      attachedVPCs[i].Region,
					AccountName: attachedVPCs[i].AccountName,
					VpcID:       attachedVPCs[i].VPCId,
				}
				sdr.AttachedVPCs = append(sdr.AttachedVPCs, vpcSolo)
			} else {
				gateway := &Gateway{
					VpcID: attachedVPCs[i].VPCId,
				}
				gateway, err = c.GetTransitGwFromVpcID(gateway)
				if err != nil {
					return nil, err
				}
				awsTgw.AttachedAviatrixTransitGW = append(awsTgw.AttachedAviatrixTransitGW, gateway.GwName)
			}
		}

		awsTgw.SecurityDomains = append(awsTgw.SecurityDomains, sdr)
	}

	return awsTgw, nil
}

func (c *Client) UpdateAWSTgw(awsTgw *AWSTgw) error {
	return nil
}

func (c *Client) DeleteAWSTgw(awsTgw *AWSTgw) error {
	awsTgw.CID = c.CID
	awsTgw.Action = "delete_aws_tgw"
	resp, err := c.Post(c.baseURL, awsTgw)
	if err != nil {
		return err
	}

	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		return errors.New(data.Reason)
	}

	return nil
}

func (c *Client) ValidateAWSTgwDomains(domainsAll []string, domainConnAll [][]string, attachedVPCAll [][]string,
) ([]string, [][]string, [][]string, error) {

	sort.Strings(domainsAll)

	numOfDomains := len(domainsAll)
	matrix := make([][]int, numOfDomains)
	var domainsToCreate []string
	var domainConnPolicy [][]string
	var domainConnRemove [][]string

	for i := range matrix {
		matrix[i] = make([]int, numOfDomains)
	}

	m := make(map[string]int)
	for i := 1; i <= numOfDomains; i++ {
		if m[domainsAll[i-1]] != 0 {
			err := fmt.Errorf("duplicate domains (name: %v) to create", domainsAll[i-1])
			return domainsToCreate, domainConnPolicy, domainConnRemove, err
		}
		m[domainsAll[i-1]] = i
	}

	m1 := make(map[string]int)
	for i := 1; i <= len(attachedVPCAll); i++ {
		if m1[attachedVPCAll[i-1][1]] != 0 {
			err := fmt.Errorf("duplicate VPC IDs (ID: %v) to attach", attachedVPCAll[i-1][1])
			return domainsToCreate, domainConnPolicy, domainConnRemove, err
		}
		m1[attachedVPCAll[i-1][1]] = i
	}

	var dmConnections []string

	for i := range domainConnAll {
		x := m[domainConnAll[i][0]]
		y := m[domainConnAll[i][1]]

		temp := "" + domainConnAll[i][0] + " - " + domainConnAll[i][1]
		dmConnections = append(dmConnections, temp)

		if x == 0 {
			err := fmt.Errorf("unrecognized domain name (%v) in domain connection", domainConnAll[i][0])
			return domainsToCreate, domainConnPolicy, domainConnRemove, err
		}
		if y == 0 {
			err := fmt.Errorf("unrecognized domain name (%v) in domain connection", domainConnAll[i][1])
			return domainsToCreate, domainConnPolicy, domainConnRemove, err
		}
		if x == y {
			err := fmt.Errorf("connection between same domains (name: %v)", domainConnAll[i][0])
			return domainsToCreate, domainConnPolicy, domainConnRemove, err
		}

		matrix[x-1][y-1] = 1
	}

	m2 := make(map[string]int)
	for i := 1; i <= len(dmConnections); i++ {
		if m2[dmConnections[i-1]] != 0 {
			err := fmt.Errorf("duplicate domain connections (%v)", dmConnections[i-1])
			return domainsToCreate, domainConnPolicy, domainConnRemove, err
		}
		m2[dmConnections[i-1]] = i
	}

	for i := 0; i < numOfDomains; i++ {
		for j := i + 1; j < numOfDomains; j++ {
			if matrix[i][j] != matrix[j][i] {
				err := fmt.Errorf("unsymmetric domain connection (%v)", ""+domainsAll[i]+" - "+domainsAll[j])
				return domainsToCreate, domainConnPolicy, domainConnRemove, err
			}
		}
	}

	defaultX := [3]string{"Default_Domain", "Shared_Service_Domain", "Aviatrix_Edge_Domain"}

	for i := 0; i < 3; i++ {
		for j := i; j < 3; j++ {
			if i != j {
				if matrix[m[defaultX[i]]-1][m[defaultX[j]]-1] == 0 {
					temp := []string{defaultX[i], defaultX[j]}
					domainConnRemove = append(domainConnRemove, temp)
				}
				matrix[m[defaultX[i]]-1][m[defaultX[j]]-1] = 2
				matrix[m[defaultX[j]]-1][m[defaultX[i]]-1] = 2
			}
		}
	}

	for i := range domainConnAll {
		if matrix[m[domainConnAll[i][0]]-1][m[domainConnAll[i][1]]-1] == 1 {
			matrix[m[domainConnAll[i][0]]-1][m[domainConnAll[i][1]]-1] = 2
			matrix[m[domainConnAll[i][1]]-1][m[domainConnAll[i][0]]-1] = 2
			temp := []string{domainConnAll[i][0], domainConnAll[i][1]}
			domainConnPolicy = append(domainConnPolicy, temp)
		}
	}

	for i := range domainsAll {
		if domainsAll[i] != "Default_Domain" &&
			domainsAll[i] != "Shared_Service_Domain" &&
			domainsAll[i] != "Aviatrix_Edge_Domain" {
			domainsToCreate = append(domainsToCreate, domainsAll[i])
		}
	}

	return domainsToCreate, domainConnPolicy, domainConnRemove, nil
}

func (c *Client) AttachAviatrixTransitGWToAWSTgw(awsTgw *AWSTgw, gateway *Gateway, SecurityDomainName string) error {
	transitGw, err := c.GetGateway(gateway)
	if err != nil {
		return err
	}

	path := c.baseURL + fmt.Sprintf("?action=attach_vpc_to_tgw&CID=%s&region=%s&vpc_account_name=%s&vpc_name="+
		"%s&gateway_name=%s&tgw_account_name=%s&tgw_name=%s&route_domain_name=%s", c.CID, awsTgw.Region,
		transitGw.AccountName, transitGw.VpcID, transitGw.GwName, awsTgw.AccountName, awsTgw.Name, SecurityDomainName)
	resp, err := c.Get(path, nil)
	if err != nil {
		return err
	}

	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		return errors.New(data.Reason)
	}

	return nil
}

func (c *Client) DetachAviatrixTransitGWToAWSTgw(awsTgw *AWSTgw, gateway *Gateway, SecurityDomainName string) error {
	transitGw, err := c.GetGateway(gateway)

	if err != nil {
		return err
	}
	path := c.baseURL + fmt.Sprintf("?action=detach_vpc_from_tgw&CID=%s&tgw_name=%s&vpc_name=%s", c.CID,
		awsTgw.Name, transitGw.VpcID)

	resp, err := c.Get(path, nil)

	if err != nil {
		return err
	}

	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		return errors.New(data.Reason)
	}

	return nil
}

func (c *Client) AttachVpcToAWSTgw(awsTgw *AWSTgw, vpcSolo VPCSolo, SecurityDomainName string) error {
	path := c.baseURL + fmt.Sprintf("?action=attach_vpc_to_tgw&region=%s&vpc_account_name=%s&vpc_name=%s"+
		"&tgw_name=%s&route_domain_name=%s&CID=%s", awsTgw.Region, vpcSolo.AccountName, vpcSolo.VpcID, awsTgw.Name,
		SecurityDomainName, c.CID)

	resp, err := c.Get(path, nil)

	if err != nil {
		return err
	}

	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		return errors.New(data.Reason)
	}

	return nil
}

func (c *Client) DetachVpcFromAWSTgw(awsTgw *AWSTgw, vpcID string) error {
	path := c.baseURL + fmt.Sprintf("?action=detach_vpc_from_tgw&CID=%s&tgw_name=%s&vpc_name=%s", c.CID,
		awsTgw.Name, vpcID)
	resp, err := c.Get(path, nil)

	if err != nil {
		return err
	}

	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		return errors.New(data.Reason)
	}

	return nil
}

func (c *Client) GetTransitGwFromVpcID(gateway *Gateway) (*Gateway, error) {
	path := c.baseURL + fmt.Sprintf("?action=list_vpcs_summary&CID=%s", c.CID)
	resp, err := c.Get(path, nil)

	if err != nil {
		return nil, err
	}

	data := VPCList{
		Return:  false,
		Results: make([]VPCInfo, 0),
		Reason:  "",
	}

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if !data.Return {
		return nil, errors.New(data.Reason)
	}

	vpcLists := data.Results
	for i := range vpcLists {
		vpcId := vpcLists[i].VPCId
		if vpcLists[i].TransitVpc == "yes" && vpcId != "" {
			index := strings.Index(vpcId, "~~")
			if index > 0 {
				vpcId = vpcId[:index]
			}
			if vpcId == gateway.VpcID {
				gateway.GwName = vpcLists[i].Name
				return gateway, nil
			}
		}
	}
	log.Printf("Couldn't find transit gateway attached to vpc %s", gateway.VpcID)
	return nil, ErrNotFound
}
