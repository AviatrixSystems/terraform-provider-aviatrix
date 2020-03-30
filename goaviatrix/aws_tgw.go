package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"sort"
	"strconv"
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
	ManageVpcAttachment       string
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
	Associations           []string             `json:"associations"`
	Name                   string               `json:"name"`
	ConnectedRouteDomain   []string             `json:"connected_route_domain"`
	AttachedVPC            []AttachedVPCDetail  `json:"attached_vpc"`
	RoutesInRouteTable     []RoutesInRouteTable `json:"routes_in_route_table"`
	RouteTableId           string               `json:"route_table_id"`
	AviatrixFirewallDomain bool                 `json:"firewall_domain"`
	NativeEgressDomain     bool                 `json:"egress_domain"`
	NativeFirewallDomain   bool                 `json:"native_firewall_domain"`
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

type TGWInfoResp struct {
	Return  bool        `json:"return"`
	Results TGWInfoList `json:"results"`
	Reason  string      `json:"reason"`
}

type TGWInfoList struct {
	TgwInfo TgwInfoDetail `json:"tgw_info"`
	TgwID   string        `json:"_id"`
	Name    string        `json:"name"`
}

type TgwInfoDetail struct {
	AccountName     string `json:"acct_name"`
	Region          string `json:"region"`
	AwsSideAsNumber int    `json:"tgw_aws_asn"`
}

type listAttachedVpcNamesResp struct {
	Return  bool     `json:"return"`
	Results []string `json:"results"`
	Reason  string   `json:"reason"`
}

type AttachmentRouteTableDetailsAPIResp struct {
	Return  bool                        `json:"return"`
	Results AttachmentRouteTableDetails `json:"results"`
	Reason  string                      `json:"reason"`
}

type AttachmentRouteTableDetails struct {
	VpcId                        string   `json:"vpc_id"`
	VpcName                      string   `json:"vpc_name"`
	VpcRegion                    string   `json:"vpc_region"`
	VpcAccount                   string   `json:"vpc_account"`
	RouteDomainName              string   `json:"route_domain_name"`
	Subnets                      []string `json:"attach_subnet_list"`
	RouteTables                  string   `json:"route_table_list"`
	CustomizedRoutes             []string `json:"customized_routes"`
	CustomizedRouteAdvertisement []string `json:"customized_routes_advertise"`
	DisableLocalRoutePropagation bool     `json:"disable_local_propagation"`
}

type ListAwsTgwAttachmentAPIResp struct {
	Return  bool               `json:"return"`
	Results []AttachmentDetail `json:"results"`
	Reason  string             `json:"reason"`
}

type AttachmentDetail struct {
	VpcID   string `json:"vpc_id"`
	TgwName string `json:"tgw_name"`
	GwName  string `json:"avx_gw_name"`
}

func (c *Client) CreateAWSTgw(awsTgw *AWSTgw) error {
	awsTgw.CID = c.CID
	awsTgw.Action = "add_aws_tgw"
	resp, err := c.Post(c.baseURL, awsTgw)
	if err != nil {
		return errors.New("HTTP Post add_aws_tgw failed: " + err.Error())
	}

	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode add_aws_tgw failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API add_aws_tgw Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetAWSTgw(awsTgw *AWSTgw) (*AWSTgw, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for list_route_domain_names") + err.Error())
	}
	listRouteDomainNames := url.Values{}
	listRouteDomainNames.Add("CID", c.CID)
	listRouteDomainNames.Add("action", "list_route_domain_names")
	listRouteDomainNames.Add("tgw_name", awsTgw.Name)
	Url.RawQuery = listRouteDomainNames.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get list_route_domain_names failed: " + err.Error())
	}

	data := AWSTgwAPIResp{
		Return:  false,
		Results: make([]string, 0),
		Reason:  "",
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_route_domain_names failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return nil, errors.New("Rest API list_route_domain_names Get failed: " + data.Reason)
	}

	connectedDomainList := data.Results
	connectedDomainList = append([]string{"Aviatrix_Edge_Domain"}, connectedDomainList...)

	for i := range connectedDomainList {
		dm := connectedDomainList[i]

		viewRouteDomainDetails := url.Values{}
		viewRouteDomainDetails.Add("CID", c.CID)
		viewRouteDomainDetails.Add("action", "view_route_domain_details")
		viewRouteDomainDetails.Add("tgw_name", awsTgw.Name)
		viewRouteDomainDetails.Add("route_domain_name", dm)

		Url.RawQuery = viewRouteDomainDetails.Encode()
		resp, err := c.Get(Url.String(), nil)

		if err != nil {
			return nil, errors.New("HTTP Get view_route_domain_details failed: " + err.Error())
		}

		var data1 RouteDomainAPIResp
		buf = new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		bodyString = buf.String()
		bodyIoCopy = strings.NewReader(bodyString)
		if err = json.NewDecoder(bodyIoCopy).Decode(&data1); err != nil {
			return nil, errors.New("Json Decode view_route_domain_details failed: " + err.Error() + "\n Body: " + bodyString)
		}
		if !data1.Return {
			return nil, errors.New("Rest API view_route_domain_details Get failed: " + data1.Reason)
		}
		routeDomainDetail := data1.Results
		sdr := SecurityDomainRule{
			Name:                   routeDomainDetail[0].Name,
			AviatrixFirewallDomain: routeDomainDetail[0].AviatrixFirewallDomain,
			NativeEgressDomain:     routeDomainDetail[0].NativeEgressDomain,
			NativeFirewallDomain:   routeDomainDetail[0].NativeFirewallDomain,
		}
		for i := range routeDomainDetail[0].ConnectedRouteDomain {
			sdr.ConnectedDomain = append(sdr.ConnectedDomain, routeDomainDetail[0].ConnectedRouteDomain[i])
		}

		attachedVPCs := routeDomainDetail[0].AttachedVPC
		for i := range attachedVPCs {
			if strings.Contains(attachedVPCs[i].VPCId, "vpn-") {
				continue
			}

			if dm != "Aviatrix_Edge_Domain" {
				listAttachmentRouteTableDetails := url.Values{}
				listAttachmentRouteTableDetails.Add("CID", c.CID)
				listAttachmentRouteTableDetails.Add("action", "list_attachment_route_table_details")
				listAttachmentRouteTableDetails.Add("tgw_name", awsTgw.Name)
				listAttachmentRouteTableDetails.Add("attachment_name", attachedVPCs[i].VPCId)
				Url.RawQuery = listAttachmentRouteTableDetails.Encode()
				resp, err := c.Get(Url.String(), nil)
				if err != nil {
					return nil, errors.New("HTTP Get list_attachment_route_table_details failed: " + err.Error())
				}

				var data2 AttachmentRouteTableDetailsAPIResp
				buf = new(bytes.Buffer)
				buf.ReadFrom(resp.Body)
				bodyString = buf.String()
				bodyIoCopy = strings.NewReader(bodyString)
				if err = json.NewDecoder(bodyIoCopy).Decode(&data2); err != nil {
					return nil, errors.New("Json Decode list_attachment_route_table_details failed: " + err.Error() + "\n Body: " + bodyString)
				}
				if !data2.Return {
					return nil, errors.New("Rest API list_attachment_route_table_details Get failed: " + data2.Reason)
				}
				vpcSolo := VPCSolo{
					Region:                       attachedVPCs[i].Region,
					AccountName:                  attachedVPCs[i].AccountName,
					VpcID:                        attachedVPCs[i].VPCId,
					DisableLocalRoutePropagation: data2.Results.DisableLocalRoutePropagation,
					RouteTables:                  data2.Results.RouteTables,
				}
				if data2.Results.CustomizedRoutes != nil && len(data2.Results.CustomizedRoutes) != 0 {
					customizedRoutes := ""
					length := len(data2.Results.CustomizedRoutes)
					for i := 0; i < length-1; i++ {
						customizedRoutes += data2.Results.CustomizedRoutes[i] + ","
					}
					vpcSolo.CustomizedRoutes = customizedRoutes + data2.Results.CustomizedRoutes[length-1]
				}
				if data2.Results.CustomizedRouteAdvertisement != nil && len(data2.Results.CustomizedRouteAdvertisement) != 0 {
					customizedRouteAdvertisement := ""
					length := len(data2.Results.CustomizedRouteAdvertisement)
					for i := 0; i < length-1; i++ {
						customizedRouteAdvertisement += data2.Results.CustomizedRouteAdvertisement[i] + ","
					}
					vpcSolo.CustomizedRouteAdvertisement = customizedRouteAdvertisement + data2.Results.CustomizedRouteAdvertisement[length-1]
				}
				if data2.Results.Subnets != nil && len(data2.Results.Subnets) != 0 {
					subnets := ""
					length := len(data2.Results.Subnets)
					for i := 0; i < length-1; i++ {
						subnets += strings.Split(data2.Results.Subnets[i], "~~")[0] + ","
					}
					vpcSolo.Subnets = subnets + strings.Split(data2.Results.Subnets[length-1], "~~")[0]
				}
				sdr.AttachedVPCs = append(sdr.AttachedVPCs, vpcSolo)
			} else {
				gateway := &Gateway{
					VpcID: attachedVPCs[i].VPCId,
				}
				gateway, err = c.GetTransitGwFromVpcID(awsTgw, gateway)
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

func (c *Client) IsFirewallSecurityDomain(tgwName string, domainName string) (bool, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return false, errors.New(("url Parsing failed for list_route_domain_names") + err.Error())
	}
	viewRouteDomainDetails := url.Values{}
	viewRouteDomainDetails.Add("CID", c.CID)
	viewRouteDomainDetails.Add("action", "view_route_domain_details")
	viewRouteDomainDetails.Add("tgw_name", tgwName)
	viewRouteDomainDetails.Add("route_domain_name", domainName)
	Url.RawQuery = viewRouteDomainDetails.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return false, errors.New("HTTP Get view_route_domain_details failed: " + err.Error())
	}
	var data RouteDomainAPIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return false, errors.New("Json Decode view_route_domain_details failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return false, errors.New("Rest API view_route_domain_details Get failed: " + data.Reason)
	}
	routeDomainDetail := data.Results
	if routeDomainDetail != nil && len(routeDomainDetail) != 0 {
		if routeDomainDetail[0].AviatrixFirewallDomain {
			return true, nil
		}
		return false, nil
	}
	return false, ErrNotFound
}

func (c *Client) UpdateAWSTgw(awsTgw *AWSTgw) error {
	return nil
}

func (c *Client) DeleteAWSTgw(awsTgw *AWSTgw) error {
	awsTgw.CID = c.CID
	awsTgw.Action = "delete_aws_tgw"
	resp, err := c.Post(c.baseURL, awsTgw)
	if err != nil {
		return errors.New("HTTP Post delete_aws_tgw failed: " + err.Error())
	}

	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode delete_aws_tgw failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API delete_aws_tgw Get failed: " + data.Reason)
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

	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for attach_vpc_to_tgw") + err.Error())
	}
	attachVpcToTgw := url.Values{}
	attachVpcToTgw.Add("CID", c.CID)
	attachVpcToTgw.Add("action", "attach_vpc_to_tgw")
	attachVpcToTgw.Add("region", awsTgw.Region)
	attachVpcToTgw.Add("vpc_account_name", transitGw.AccountName)
	attachVpcToTgw.Add("vpc_name", transitGw.VpcID)
	attachVpcToTgw.Add("gateway_name", transitGw.GwName)
	attachVpcToTgw.Add("tgw_account_name", awsTgw.AccountName)
	attachVpcToTgw.Add("tgw_name", awsTgw.Name)
	attachVpcToTgw.Add("route_domain_name", SecurityDomainName)
	Url.RawQuery = attachVpcToTgw.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get attach_vpc_to_tgw failed: " + err.Error())
	}

	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode attach_vpc_to_tgw failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API attach_vpc_to_tgw Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DetachAviatrixTransitGWFromAWSTgw(awsTgw *AWSTgw, gateway *Gateway, SecurityDomainName string) error {
	transitGw, err := c.GetGateway(gateway)
	if err != nil {
		return err
	}

	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for detach_vpc_from_tgw") + err.Error())
	}
	detachVpcFromTgw := url.Values{}
	detachVpcFromTgw.Add("CID", c.CID)
	detachVpcFromTgw.Add("action", "detach_vpc_from_tgw")
	detachVpcFromTgw.Add("tgw_name", awsTgw.Name)
	detachVpcFromTgw.Add("vpc_name", transitGw.VpcID)
	Url.RawQuery = detachVpcFromTgw.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get detach_vpc_from_tgw failed: " + err.Error())
	}

	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode detach_vpc_from_tgw failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		if strings.Contains(data.Reason, "is not attached to") {
			return nil
		}
		return errors.New("Rest API detach_vpc_from_tgw Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) AttachVpcToAWSTgw(awsTgw *AWSTgw, vpcSolo VPCSolo, SecurityDomainName string) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for attach_vpc_to_tgw") + err.Error())
	}
	attachVpcFromTgw := url.Values{}
	attachVpcFromTgw.Add("CID", c.CID)
	attachVpcFromTgw.Add("action", "attach_vpc_to_tgw")
	attachVpcFromTgw.Add("region", awsTgw.Region)
	attachVpcFromTgw.Add("vpc_account_name", vpcSolo.AccountName)
	attachVpcFromTgw.Add("vpc_name", vpcSolo.VpcID)
	attachVpcFromTgw.Add("tgw_name", awsTgw.Name)
	attachVpcFromTgw.Add("route_domain_name", SecurityDomainName)
	if vpcSolo.DisableLocalRoutePropagation {
		attachVpcFromTgw.Add("disable_local_route_propagation", "yes")
	}
	if vpcSolo.CustomizedRoutes != "" {
		attachVpcFromTgw.Add("customized_routes", vpcSolo.CustomizedRoutes)
	}
	if vpcSolo.CustomizedRouteAdvertisement != "" {
		attachVpcFromTgw.Add("customized_route_advertisement", vpcSolo.CustomizedRouteAdvertisement)
	}
	if vpcSolo.Subnets != "" {
		attachVpcFromTgw.Add("subnet_list", vpcSolo.Subnets)
	}
	if vpcSolo.CustomizedRoutes != "" {
		attachVpcFromTgw.Add("route_table_list", vpcSolo.RouteTables)
	}
	Url.RawQuery = attachVpcFromTgw.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get attach_vpc_to_tgw failed: " + err.Error())
	}

	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode attach_vpc_to_tgw failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API attach_vpc_to_tgw Get failed: " + data.Reason)
	}

	return nil
}

func (c *Client) DetachVpcFromAWSTgw(awsTgw *AWSTgw, vpcID string) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for detach_vpc_from_tgw") + err.Error())
	}
	detachVpcFromTgw := url.Values{}
	detachVpcFromTgw.Add("CID", c.CID)
	detachVpcFromTgw.Add("action", "detach_vpc_from_tgw")
	detachVpcFromTgw.Add("tgw_name", awsTgw.Name)
	detachVpcFromTgw.Add("vpc_name", vpcID)
	Url.RawQuery = detachVpcFromTgw.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get detach_vpc_from_tgw failed: " + err.Error())
	}

	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode detach_vpc_from_tgw failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API detach_vpc_from_tgw Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetTransitGwFromVpcID(awsTgw *AWSTgw, gateway *Gateway) (*Gateway, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for 'list_all_tgw_attachments': ") + err.Error())
	}
	listTgwDetails := url.Values{}
	listTgwDetails.Add("CID", c.CID)
	listTgwDetails.Add("action", "list_all_tgw_attachments")
	Url.RawQuery = listTgwDetails.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get 'list_all_tgw_attachments' failed: " + err.Error())
	}

	var data ListAwsTgwAttachmentAPIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode 'list_all_tgw_attachments' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return nil, errors.New("Rest API 'list_all_tgw_attachments' Get failed: " + data.Reason)
	}
	for i := range data.Results {
		if data.Results[i].TgwName == awsTgw.Name && data.Results[i].VpcID == gateway.VpcID && data.Results[i].GwName != "" {
			gateway.GwName = data.Results[i].GwName
			return gateway, nil
		}
	}
	log.Printf("Couldn't find transit gateway attached to vpc %s", gateway.VpcID)
	return nil, ErrNotFound
}

func (c *Client) ListTgwDetails(awsTgw *AWSTgw) (*AWSTgw, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for list_tgw_details") + err.Error())
	}
	listTgwDetails := url.Values{}
	listTgwDetails.Add("CID", c.CID)
	listTgwDetails.Add("action", "list_tgw_details")
	listTgwDetails.Add("tgw_name", awsTgw.Name)
	Url.RawQuery = listTgwDetails.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get list_tgw_details failed: " + err.Error())
	}

	var data TGWInfoResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_tgw_details failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		if strings.Contains(data.Reason, "does not exist") {
			return nil, ErrNotFound
		}
		return nil, errors.New("Rest API list_tgw_details Get failed: " + data.Reason)
	}
	tgwInfoList := data.Results
	if tgwInfoList.Name == awsTgw.Name {
		tgwInfoDetail := tgwInfoList.TgwInfo
		awsTgw.AccountName = tgwInfoDetail.AccountName
		awsTgw.Region = tgwInfoDetail.Region
		awsTgw.AwsSideAsNumber = strconv.Itoa(tgwInfoDetail.AwsSideAsNumber)
		return awsTgw, nil
	}
	return nil, ErrNotFound
}

func (c *Client) IsVpcAttachedToTgw(awsTgw *AWSTgw, vpcSolo *VPCSolo) (bool, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return false, errors.New(("url Parsing failed for list_attached_vpc_names_to_route_domain") + err.Error())
	}
	listAttachedVpcNames := url.Values{}
	listAttachedVpcNames.Add("CID", c.CID)
	listAttachedVpcNames.Add("action", "list_attached_vpc_names_to_route_domain")
	listAttachedVpcNames.Add("tgw_name", awsTgw.Name)
	Url.RawQuery = listAttachedVpcNames.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return false, errors.New("HTTP Get list_attached_vpc_names_to_route_domain failed: " + err.Error())
	}

	data := listAttachedVpcNamesResp{
		Return:  false,
		Results: make([]string, 0),
		Reason:  "",
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return false, errors.New("Json Decode list_attached_vpc_names_to_route_domain failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return false, errors.New("Rest API list_attached_vpc_names_to_route_domain Get failed: " + data.Reason)
	}

	attachedVpcNames := data.Results
	for i := range attachedVpcNames {
		if strings.Split(attachedVpcNames[i], "~~")[0] == vpcSolo.VpcID {
			return true, nil
		}
	}
	return false, ErrNotFound
}

func (c *Client) GetAttachmentRouteTableDetails(tgwName string, attachmentName string) (*AttachmentRouteTableDetails, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for list_attachment_route_table_details") + err.Error())
	}
	listAttachmentRouteTableDetails := url.Values{}
	listAttachmentRouteTableDetails.Add("CID", c.CID)
	listAttachmentRouteTableDetails.Add("action", "list_attachment_route_table_details")
	listAttachmentRouteTableDetails.Add("tgw_name", tgwName)
	listAttachmentRouteTableDetails.Add("attachment_name", attachmentName)
	Url.RawQuery = listAttachmentRouteTableDetails.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get list_attachment_route_table_details failed: " + err.Error())
	}

	var data AttachmentRouteTableDetailsAPIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_attachment_route_table_details failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return nil, errors.New("Rest API list_attachment_route_table_details Get failed: " + data.Reason)
	}

	return &data.Results, nil
}
