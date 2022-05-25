package goaviatrix

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// AwsTGW simple struct to hold aws_tgw details
type AWSTgw struct {
	Action                    string               `form:"action,omitempty"`
	CID                       string               `form:"CID,omitempty"`
	Name                      string               `form:"tgw_name,omitempty"`
	CloudType                 int                  `form:"cloud_type,omitempty" json:"cloud_type,omitempty"`
	AccountName               string               `form:"account_name,omitempty"`
	Region                    string               `form:"region,omitempty"`
	AwsSideAsNumber           string               `form:"aws_side_asn,omitempty"`
	AttachedAviatrixTransitGW []string             `form:"attached_aviatrix_transit_gateway,omitempty"`
	SecurityDomains           []SecurityDomainRule `form:"security_domains,omitempty"`
	ManageVpcAttachment       string
	EnableMulticast           bool `form:"multicast_enable"`
	CidrList                  []string
	NotCreateDefaultDomains   bool `form:"not_create_default_domains,omitempty"`
	TgwId                     string
	Async                     bool `form:"async,omitempty"`
	InspectionMode            string
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
	VPCId           []string `json:"vpc_id"`
	CidrBlock       string   `json:"cidr_block"`
	Type            string   `json:"type"`
	State           string   `json:"state"`
	TgwAttachmentId []string `json:"tgw_attachment_id"`
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
	AccountName               string   `json:"acct_name"`
	Region                    string   `json:"region"`
	AwsSideAsNumber           int      `json:"tgw_aws_asn"`
	CloudType                 int      `json:"cloud_type"`
	EnableMulticast           bool     `json:"multicast_enable"`
	CidrList                  []string `json:"tgw_cidr_list"`
	TgwId                     string   `json:"tgw_id"`
	ConnectionBasedInspection bool     `json:"connection_based_inspection"`
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
	awsTgw.Async = true
	return c.PostAsyncAPI(awsTgw.Action, awsTgw, BasicCheck)
}

func (c *Client) GetAWSTgw(awsTgw *AWSTgw) (*AWSTgw, error) {
	form := map[string]string{
		"CID":      c.CID,
		"action":   "list_route_domain_names",
		"tgw_name": awsTgw.Name,
	}
	var data AWSTgwAPIResp
	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	connectedDomainList := data.Results
	connectedDomainList = append([]string{"Aviatrix_Edge_Domain"}, connectedDomainList...)

	for i := range connectedDomainList {
		dm := connectedDomainList[i]
		if strings.HasPrefix(dm, "peering_") || strings.Contains(dm, ":") {
			continue
		}

		form = map[string]string{
			"CID":               c.CID,
			"action":            "view_route_domain_details",
			"tgw_name":          awsTgw.Name,
			"route_domain_name": dm,
		}
		var data1 RouteDomainAPIResp
		err = c.GetAPI(&data1, form["action"], form, BasicCheck)
		if err != nil {
			return nil, err
		}

		routeDomainDetail := data1.Results
		if len(routeDomainDetail) == 0 {
			continue
		}
		sdr := SecurityDomainRule{
			Name:                   routeDomainDetail[0].Name,
			AviatrixFirewallDomain: routeDomainDetail[0].AviatrixFirewallDomain,
			NativeEgressDomain:     routeDomainDetail[0].NativeEgressDomain,
			NativeFirewallDomain:   routeDomainDetail[0].NativeFirewallDomain,
		}
		for i := range routeDomainDetail[0].ConnectedRouteDomain {
			if strings.HasPrefix(routeDomainDetail[0].ConnectedRouteDomain[i], "peering_") || strings.Contains(routeDomainDetail[0].ConnectedRouteDomain[i], ":") {
				continue
			}
			sdr.ConnectedDomain = append(sdr.ConnectedDomain, routeDomainDetail[0].ConnectedRouteDomain[i])
		}

		attachedVPCs := routeDomainDetail[0].AttachedVPC
		for i := range attachedVPCs {
			if strings.Contains(attachedVPCs[i].VPCId, "vpn-") {
				continue
			}

			if dm != "Aviatrix_Edge_Domain" {
				form = map[string]string{
					"CID":             c.CID,
					"action":          "list_attachment_route_table_details",
					"tgw_name":        awsTgw.Name,
					"attachment_name": attachedVPCs[i].VPCId,
				}
				var data2 AttachmentRouteTableDetailsAPIResp
				err = c.GetAPI(&data2, form["action"], form, BasicCheck)
				if err != nil {
					return nil, err
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
	form := map[string]string{
		"CID":               c.CID,
		"action":            "view_route_domain_details",
		"tgw_name":          tgwName,
		"route_domain_name": domainName,
	}
	var data RouteDomainAPIResp
	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return false, err
	}
	routeDomainDetail := data.Results
	if len(routeDomainDetail) != 0 {
		if routeDomainDetail[0].AviatrixFirewallDomain {
			return true, nil
		}
		return false, nil
	}
	return false, ErrNotFound
}

func (c *Client) DeleteAWSTgw(awsTgw *AWSTgw) error {
	awsTgw.CID = c.CID
	awsTgw.Action = "delete_aws_tgw"
	return c.PostAPI(awsTgw.Action, awsTgw, BasicCheck)
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
		return fmt.Errorf("could not get transit gateway to attach to AWS TGW: %v", err)
	}
	form := map[string]string{
		"CID":               c.CID,
		"action":            "attach_vpc_to_tgw",
		"region":            awsTgw.Region,
		"vpc_account_name":  transitGw.AccountName,
		"vpc_name":          transitGw.VpcID,
		"gateway_name":      transitGw.GwName,
		"tgw_account_name":  awsTgw.AccountName,
		"tgw_name":          awsTgw.Name,
		"route_domain_name": SecurityDomainName,
		"async":             "true",
	}
	return c.PostAsyncAPI(form["action"], form, BasicCheck)
}

func (c *Client) DetachAviatrixTransitGWFromAWSTgw(awsTgw *AWSTgw, gateway *Gateway, SecurityDomainName string) error {
	transitGw, err := c.GetGateway(gateway)
	if err != nil {
		return fmt.Errorf("could not get transit gateway to detach from AWS TGW: %v", err)
	}
	form := map[string]string{
		"CID":      c.CID,
		"action":   "detach_vpc_from_tgw",
		"tgw_name": awsTgw.Name,
		"vpc_name": transitGw.VpcID,
		"async":    "true",
	}
	check := func(action, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "is not attached to") {
				return nil
			}
			return fmt.Errorf("rest API %s %s failed: %s", action, method, reason)
		}
		return nil
	}
	return c.PostAsyncAPI(form["action"], form, check)
}

func (c *Client) AttachVpcToAWSTgw(awsTgw *AWSTgw, vpcSolo VPCSolo, SecurityDomainName string) error {
	form := map[string]string{
		"CID":               c.CID,
		"action":            "attach_vpc_to_tgw",
		"region":            awsTgw.Region,
		"vpc_account_name":  vpcSolo.AccountName,
		"vpc_name":          vpcSolo.VpcID,
		"tgw_name":          awsTgw.Name,
		"route_domain_name": SecurityDomainName,
		"async":             "true",
	}
	if vpcSolo.DisableLocalRoutePropagation {
		form["disable_local_route_propagation"] = "yes"
	}
	if vpcSolo.CustomizedRoutes != "" {
		form["customized_routes"] = vpcSolo.CustomizedRoutes
	}
	if vpcSolo.CustomizedRouteAdvertisement != "" {
		form["customized_route_advertisement"] = vpcSolo.CustomizedRouteAdvertisement
	}
	if vpcSolo.Subnets != "" {
		form["subnet_list"] = vpcSolo.Subnets
	}
	if vpcSolo.CustomizedRoutes != "" {
		form["route_table_list"] = vpcSolo.RouteTables
	}
	return c.PostAsyncAPI(form["action"], form, BasicCheck)
}

func (c *Client) DetachVpcFromAWSTgw(awsTgw *AWSTgw, vpcID string) error {
	form := map[string]string{
		"CID":      c.CID,
		"action":   "detach_vpc_from_tgw",
		"tgw_name": awsTgw.Name,
		"vpc_name": vpcID,
		"async":    "true",
	}
	return c.PostAsyncAPI(form["action"], form, BasicCheck)
}

func (c *Client) GetTransitGwFromVpcID(awsTgw *AWSTgw, gateway *Gateway) (*Gateway, error) {
	var data ListAwsTgwAttachmentAPIResp
	form := map[string]string{
		"CID":    c.CID,
		"action": "list_all_tgw_attachments",
	}
	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	for i := range data.Results {
		if data.Results[i].TgwName == awsTgw.Name && data.Results[i].VpcID == gateway.VpcID && data.Results[i].GwName != "" {
			gateway.GwName = data.Results[i].GwName
			return gateway, nil
		}
	}
	log.Errorf("Couldn't find transit gateway attached to vpc %s", gateway.VpcID)
	return nil, ErrNotFound
}

func (c *Client) ListTgwDetails(awsTgw *AWSTgw) (*AWSTgw, error) {
	var data TGWInfoResp
	form := map[string]string{
		"CID":      c.CID,
		"action":   "list_tgw_details",
		"tgw_name": awsTgw.Name,
	}
	check := func(action, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "does not exist") {
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s %s failed: %s", action, method, reason)
		}
		return nil
	}
	err := c.GetAPI(&data, form["action"], form, check)
	if err != nil {
		return nil, err
	}

	tgwInfoList := data.Results
	if tgwInfoList.Name == awsTgw.Name {
		tgwInfoDetail := tgwInfoList.TgwInfo
		awsTgw.AccountName = tgwInfoDetail.AccountName
		awsTgw.Region = tgwInfoDetail.Region
		awsTgw.AwsSideAsNumber = strconv.Itoa(tgwInfoDetail.AwsSideAsNumber)
		awsTgw.CloudType = tgwInfoDetail.CloudType
		awsTgw.EnableMulticast = tgwInfoDetail.EnableMulticast
		awsTgw.CidrList = tgwInfoDetail.CidrList
		awsTgw.TgwId = tgwInfoDetail.TgwId
		if tgwInfoDetail.ConnectionBasedInspection {
			awsTgw.InspectionMode = "Connection-based"
		} else {
			awsTgw.InspectionMode = "Domain-based"
		}
		return awsTgw, nil
	}
	return nil, ErrNotFound
}

func (c *Client) IsVpcAttachedToTgw(awsTgw *AWSTgw, vpcSolo *VPCSolo) (bool, error) {
	var data listAttachedVpcNamesResp
	form := map[string]string{
		"CID":      c.CID,
		"action":   "list_attached_vpc_names_to_route_domain",
		"tgw_name": awsTgw.Name,
	}
	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return false, err
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
	var data AttachmentRouteTableDetailsAPIResp
	form := map[string]string{
		"CID":             c.CID,
		"action":          "list_attachment_route_table_details",
		"tgw_name":        tgwName,
		"attachment_name": attachmentName,
	}
	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}
	return &data.Results, nil
}

func (c *Client) UpdateTGWCidrs(tgwName string, cidrs []string) error {
	data := map[string]string{
		"action":    "update_tgw_cidrs",
		"CID":       c.CID,
		"tgw_name":  tgwName,
		"cidr_list": strings.Join(cidrs, ","),
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) UpdateTGWInspectionMode(tgwName, inspectionMode string) error {
	data := map[string]string{
		"action":   "edit_aws_tgw_inspection_mode",
		"CID":      c.CID,
		"tgw_name": tgwName,
		"mode":     inspectionMode,
	}

	return c.PostAPI(data["action"], data, BasicCheck)
}
