package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AwsTGW simple struct to hold aws_tgw details
type AWSTgw struct {
	Action          		  string  				  `form:"action,omitempty"`
	CID             		  string  				  `form:"CID,omitempty"`
	Name 		    		  string  			      `form:"tgw_name,omitempty"`
	AccountName     		  string  				  `form:"account_name,omitempty"`
	Region          		  string  				  `form:"region,omitempty"`
	AwsSideAsNumber 		  string                  `form:"aws_side_asn,omitempty"`
	AttachedAviatrixTransitGW []string			      `form:"attached_aviatrix_transit_gateway,omitempty"`
	SecurityDomains 		  []SecurityDomainRule    `form:"security_domains,omitempty"`
}

type AWSTgwAPIResp struct {
	Return  bool        `json:"return"`
	Results []string	`json:"results"`
	Reason  string		`json:"reason"`
}

type AWSTgwList struct {
	Return  bool   	    `json:"return"`
	Results []AWSTgw    `json:"results"`
	Reason  string      `json:"reason"`
}

type RouteDomainAPIResp struct {
	Return  bool				   `json:"return"`
	Results []RouteDomainDetail	   `json:"results"`
	Reason  string				   `json:"reason"`
}

type RouteDomainDetail struct {
	Associations 		 []string				`json:"associations"`
	Name 		 		 string 				`json:"name"`
	ConnectedRouteDomain []string   			`json:"connected_route_domain"`
	AttachedVPC 		 []AttachedVPCDetail	`json:"attached_vpc"`
	RoutesInRouteTable 	 []RoutesInRouteTable	`json:"routes_in_route_table"`
	RouteTableId 		 string					`json:"route_table_id"`
}

type AttachedVPCDetail struct {
	TgwName 	 string		`json:"tgw_name"`
	Region		 string		`json:"regioin"`
	VPCName      string		`json:"vpc_name"`
	AttachmentId string		`json:"attachment_id"`
	RouteDomain  string		`json:"route_domain"`
	VPCCidr 	 []string	`json:"vpc_cidr"`
	VPCId        string     `json:"vpc_id"`
}

type RoutesInRouteTable struct {
	VPCId 	        string    `json:"vpc_id"`
	CidrBlock		string	  `json:"cidr_block"`
	Type      		string	  `json:"type"`
	State 			string	  `json:"state"`
	TgwAttachmentId string	  `json:"tgw_attachment_id"`
}

func (c *Client) CreateAWSTgw(awsTgw *AWSTgw) (error) {
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
		path = c.baseURL + fmt.Sprintf("?action=view_route_domain_details&tgw_name=%s" +
			"&route_domain_name=%s&CID=%s", awsTgw.Name, dm, c.CID)
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
				sdr.AttachedVPCs = append(sdr.AttachedVPCs, attachedVPCs[i].VPCId + "~~" + attachedVPCs[i].VPCName)
			} else {
				awsTgw.AttachedAviatrixTransitGW = append(awsTgw.AttachedAviatrixTransitGW, attachedVPCs[i].VPCId +
					"~~" + attachedVPCs[i].VPCName)
			}
		}

		awsTgw.SecurityDomains = append(awsTgw.SecurityDomains, sdr)
	}

	return awsTgw, nil
}

func (c *Client) UpdateAWSTgw(awsTgw *AWSTgw) (error) {
	return nil
}

func (c *Client) DeleteAWSTgw(awsTgw *AWSTgw) (error) {
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

func (c *Client) ValidateAWSTgwDomains(domainsAll []string, domainConnAll [][]string) ([]string, [][]string,
	[][]string, error) {
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
		m[domainsAll[i - 1]] = i
	}

	for i := range domainConnAll {
		x := m[domainConnAll[i][0]]
		y := m[domainConnAll[i][1]]
		if x == 0 || y == 0 || x == y {
			return domainsToCreate, domainConnPolicy, domainConnRemove, ErrNotFound
		}
		matrix[x - 1][y - 1] = 1
	}

	for i := 0; i < numOfDomains; i++ {
		for j := i + 1; j < numOfDomains; j++ {
			if matrix[i][j] != matrix[j][i] {
				return domainsToCreate, domainConnPolicy, domainConnRemove, ErrNotFound
			}
		}
	}

	defaultX := [3]string{"Default_Domain", "Shared_Service_Domain", "Aviatrix_Edge_Domain"}

	for i := 0; i < 3; i++ {
		for j := i; j < 3; j++ {
			if i != j {

				if matrix[m[defaultX[i]] - 1][m[defaultX[j]] - 1] == 0 {
					temp := []string{defaultX[i], defaultX[j]}
					domainConnRemove = append(domainConnRemove, temp)
				}
				matrix[m[defaultX[i]] - 1][m[defaultX[j]] - 1] = 2
				matrix[m[defaultX[j]] - 1][m[defaultX[i]] - 1] = 2
			}
		}
	}

	for i := range domainConnAll {
		if matrix[m[domainConnAll[i][0]] - 1][m[domainConnAll[i][1]] - 1] == 1 {
			matrix[m[domainConnAll[i][0]] - 1][m[domainConnAll[i][1]] - 1] = 2
			matrix[m[domainConnAll[i][1]] - 1][m[domainConnAll[i][0]] - 1] = 2
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

func (c *Client) AttachAviatrixTransitGWToAWSTgw(awsTgw *AWSTgw, gateway *Gateway, SecurityDomainName string) (error) {
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

func (c *Client) DetachAviatrixTransitGWToAWSTgw(awsTgw *AWSTgw, gateway *Gateway, SecurityDomainName string) (error) {
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

func (c *Client) AttachVpcToAWSTgw(awsTgw *AWSTgw) (error) {
	return nil
}

func (c *Client) DetachVpcFromAWSTgw(awsTgw *AWSTgw) (error) {
	return nil
}