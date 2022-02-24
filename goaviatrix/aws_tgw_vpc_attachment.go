package goaviatrix

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type AwsTgwVpcAttachment struct {
	Action                       string `form:"action,omitempty"`
	CID                          string `form:"CID,omitempty"`
	TgwName                      string `form:"tgw_name"`
	Region                       string `form:"region"`
	SecurityDomainName           string `form:"security_domain_name"`
	VpcAccountName               string `form:"vpc_account_name"`
	VpcID                        string `form:"vpc_id"`
	CustomizedRoutes             string `form:"customized_routes,omitempty" json:"customized_routes,omitempty"`
	Subnets                      string
	RouteTables                  string
	CustomizedRouteAdvertisement string
	DisableLocalRoutePropagation bool `form:"disable_local_route_propagation,omitempty" json:"disable_local_route_propagation,omitempty"`
	EdgeAttachment               string
}

type DomainListResp struct {
	Return  bool     `json:"return,omitempty"`
	Results []string `json:"results,omitempty"`
	Reason  string   `json:"reason,omitempty"`
}

func (c *Client) CreateAwsTgwVpcAttachment(awsTgwVpcAttachment *AwsTgwVpcAttachment) error {
	form := map[string]string{
		"CID":               c.CID,
		"action":            "attach_vpc_to_tgw",
		"region":            awsTgwVpcAttachment.Region,
		"vpc_account_name":  awsTgwVpcAttachment.VpcAccountName,
		"vpc_name":          awsTgwVpcAttachment.VpcID,
		"tgw_name":          awsTgwVpcAttachment.TgwName,
		"route_domain_name": awsTgwVpcAttachment.SecurityDomainName,
		"async":             "true",
	}
	if awsTgwVpcAttachment.DisableLocalRoutePropagation {
		form["disable_local_route_propagation"] = "yes"
	}
	if awsTgwVpcAttachment.CustomizedRoutes != "" {
		form["customized_routes"] = awsTgwVpcAttachment.CustomizedRoutes
	}
	if awsTgwVpcAttachment.CustomizedRouteAdvertisement != "" {
		form["customized_route_advertisement"] = awsTgwVpcAttachment.CustomizedRouteAdvertisement
	}
	if awsTgwVpcAttachment.Subnets != "" {
		form["subnet_list"] = awsTgwVpcAttachment.Subnets
	}
	if awsTgwVpcAttachment.RouteTables != "" {
		form["route_table_list"] = awsTgwVpcAttachment.RouteTables
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) CreateAwsTgwVpcAttachmentForFireNet(awsTgwVpcAttachment *AwsTgwVpcAttachment) error {
	form := map[string]string{
		"CID":         c.CID,
		"action":      "connect_firenet_with_tgw",
		"vpc_id":      awsTgwVpcAttachment.VpcID,
		"tgw_name":    awsTgwVpcAttachment.TgwName,
		"domain_name": awsTgwVpcAttachment.SecurityDomainName,
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) GetAwsTgwVpcAttachment(awsTgwVpcAttachment *AwsTgwVpcAttachment) (*AwsTgwVpcAttachment, error) {
	awsTgw := &AWSTgw{
		Name: awsTgwVpcAttachment.TgwName,
	}
	awsTgw, err := c.ListTgwDetails(awsTgw)
	if err != nil {
		if err == ErrNotFound {
			return nil, err
		}
		return nil, fmt.Errorf("couldn't find AWS TGW %s: %v", awsTgwVpcAttachment.TgwName, err)
	}
	awsTgwVpcAttachment.Region = awsTgw.Region

	err = c.GetAwsTgwDomain(awsTgw, awsTgwVpcAttachment.SecurityDomainName)
	if err != nil {
		if err == ErrNotFound {
			return nil, err
		}
		return nil, errors.New("aws tgw does not have security domain: " + err.Error())
	}

	aTVA, err := c.GetVPCAttachmentRouteTableDetails(awsTgwVpcAttachment)
	if err != nil {
		if err == ErrNotFound {
			return nil, err
		}
		return nil, errors.New("could not get security domain details: " + err.Error())
	}

	ARTDetail, err := c.GetAttachmentRouteTableDetails(awsTgwVpcAttachment.TgwName, awsTgwVpcAttachment.VpcID)
	if err != nil {
		return nil, errors.New("could not get Attachment Route Table Details: " + err.Error())
	}
	if ARTDetail != nil {
		aTVA.DisableLocalRoutePropagation = ARTDetail.DisableLocalRoutePropagation
		if ARTDetail.CustomizedRoutes != nil && len(ARTDetail.CustomizedRoutes) != 0 {
			customizedRoutes := ""
			length := len(ARTDetail.CustomizedRoutes)
			for i := 0; i < length-1; i++ {
				customizedRoutes += ARTDetail.CustomizedRoutes[i] + ","
			}
			aTVA.CustomizedRoutes = customizedRoutes + ARTDetail.CustomizedRoutes[length-1]
		}
		if ARTDetail.Subnets != nil && len(ARTDetail.Subnets) != 0 {
			subnets := ""
			length := len(ARTDetail.Subnets)
			for i := 0; i < length-1; i++ {
				subnets += strings.Split(ARTDetail.Subnets[i], "~~")[0] + ","
			}
			aTVA.Subnets = subnets + strings.Split(ARTDetail.Subnets[length-1], "~~")[0]
		}
		aTVA.RouteTables = ARTDetail.RouteTables
		if ARTDetail.CustomizedRouteAdvertisement != nil && len(ARTDetail.CustomizedRouteAdvertisement) != 0 {
			customizedRouteAdvertisement := ""
			length := len(ARTDetail.CustomizedRouteAdvertisement)
			for i := 0; i < length-1; i++ {
				customizedRouteAdvertisement += ARTDetail.CustomizedRouteAdvertisement[i] + ","
			}
			aTVA.CustomizedRouteAdvertisement = customizedRouteAdvertisement + ARTDetail.CustomizedRouteAdvertisement[length-1]
		}
	} else {
		aTVA.DisableLocalRoutePropagation = false
		aTVA.CustomizedRoutes = ""
		aTVA.Subnets = ""
		aTVA.RouteTables = ""
		aTVA.CustomizedRouteAdvertisement = ""
	}

	firenetManagementDetails, err := c.GetFirenetManagementDetails(awsTgwVpcAttachment)
	if err != nil {
		return nil, fmt.Errorf("could not get firenet management details: %s", err)
	}
	if len(firenetManagementDetails) != 0 {
		aTVA.EdgeAttachment = firenetManagementDetails[1]
	}

	return aTVA, nil
}

func (c *Client) DeleteAwsTgwVpcAttachment(awsTgwVpcAttachment *AwsTgwVpcAttachment) error {
	form := map[string]string{
		"CID":      c.CID,
		"action":   "detach_vpc_from_tgw",
		"tgw_name": awsTgwVpcAttachment.TgwName,
		"vpc_name": awsTgwVpcAttachment.VpcID,
		"async":    "true",
	}
	return c.PostAsyncAPI(form["action"], form, BasicCheck)
}

func (c *Client) DeleteAwsTgwVpcAttachmentForFireNet(awsTgwVpcAttachment *AwsTgwVpcAttachment) error {
	form := map[string]string{
		"CID":    c.CID,
		"action": "disconnect_firenet_with_tgw",
		"vpc_id": awsTgwVpcAttachment.VpcID,
		"async":  "true",
	}
	return c.PostAsyncAPI(form["action"], form, BasicCheck)
}

func (c *Client) GetAwsTgwDetail(awsTgw *AWSTgw) (*AWSTgw, error) {
	awsTgw, err := c.ListTgwDetails(awsTgw)
	if err != nil {
		return nil, fmt.Errorf("couldn't find AWS TGW %s: %v", awsTgw.Name, err)
	}

	return awsTgw, nil
}

func (c *Client) GetAwsTgwDomain(awsTgw *AWSTgw, sDM string) error {
	var data DomainListResp
	form := map[string]string{
		"CID":      c.CID,
		"action":   "list_route_domain_names",
		"tgw_name": awsTgw.Name,
	}
	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return err
	}
	mDomain := make(map[string]bool)
	for i := range data.Results {
		mDomain[data.Results[i]] = true
	}
	if !mDomain[sDM] {
		return ErrNotFound
	}

	return nil
}

func (c *Client) GetVPCAttachmentRouteTableDetails(awsTgwVpcAttachment *AwsTgwVpcAttachment) (*AwsTgwVpcAttachment, error) {
	var data RouteDomainAPIResp
	form := map[string]string{
		"CID":               c.CID,
		"action":            "view_route_domain_details",
		"tgw_name":          awsTgwVpcAttachment.TgwName,
		"route_domain_name": awsTgwVpcAttachment.SecurityDomainName,
	}

	numberOfRetries := 3
	retryInterval := 5
	for i := 0; i < numberOfRetries; i++ {
		err := c.GetAPI(&data, form["action"], form, BasicCheck)
		if err != nil {
			return nil, err
		}
		routeDomainDetail := data.Results
		attachedVPCs := routeDomainDetail[0].AttachedVPC
		for j := range attachedVPCs {
			if attachedVPCs[j].VPCId == awsTgwVpcAttachment.VpcID {
				awsTgwVpcAttachment.VpcAccountName = attachedVPCs[j].AccountName
				return awsTgwVpcAttachment, nil
			}
		}
		time.Sleep(time.Duration(retryInterval) * time.Second)
	}

	return nil, ErrNotFound
}

func (c *Client) EditTgwSpokeVpcCustomizedRoutes(awsTgwVpcAttachment *AwsTgwVpcAttachment) error {
	form := map[string]string{
		"CID":        c.CID,
		"action":     "edit_tgw_spoke_vpc_customized_routes",
		"tgw_name":   awsTgwVpcAttachment.TgwName,
		"vpc_id":     awsTgwVpcAttachment.VpcID,
		"route_list": awsTgwVpcAttachment.CustomizedRoutes,
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) EditTgwSpokeVpcCustomizedRouteAdvertisement(awsTgwVpcAttachment *AwsTgwVpcAttachment) error {
	form := map[string]string{
		"CID":             c.CID,
		"action":          "update_customized_route_advertisement",
		"tgw_name":        awsTgwVpcAttachment.TgwName,
		"attachment_name": awsTgwVpcAttachment.VpcID,
		"cidr_list":       awsTgwVpcAttachment.CustomizedRouteAdvertisement,
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) UpdateFirewallAttachmentAccessFromOnprem(awsTgwVpcAttachment *AwsTgwVpcAttachment) error {
	params := map[string]string{
		"action":          "update_firewall_attachment_access_from_onprem",
		"CID":             c.CID,
		"tgw_name":        awsTgwVpcAttachment.TgwName,
		"attachment_name": awsTgwVpcAttachment.VpcID,
		"edge_attachment": awsTgwVpcAttachment.EdgeAttachment,
	}

	if awsTgwVpcAttachment.EdgeAttachment == "" {
		params["edge_attachment"] = "no"
	}

	return c.PostAPI(params["action"], params, BasicCheck)
}

func (c *Client) GetFirenetManagementDetails(awsTgwVpcAttachment *AwsTgwVpcAttachment) ([]string, error) {
	params := map[string]string{
		"action":          "get_tgw_attachment_details",
		"CID":             c.CID,
		"tgw_name":        awsTgwVpcAttachment.TgwName,
		"attachment_name": awsTgwVpcAttachment.VpcID,
	}

	type AccessFromEdge struct {
		AccessFromEdge []string `json:"access_from_edge"`
	}

	type Resp struct {
		Results []AccessFromEdge
	}

	var data Resp

	err := c.GetAPI(&data, params["action"], params, BasicCheck)
	if err != nil {
		return nil, err
	}

	return data.Results[0].AccessFromEdge, nil
}
