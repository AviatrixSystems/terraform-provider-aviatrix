package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type AwsTgwVpcAttachment struct {
	Action                       string `form:"action,omitempty"`
	CID                          string `form:"CID,omitempty"`
	TgwName                      string `form:"tgw_name"`
	Region                       string `form:"region"`
	SecurityDomainName           string `form:"security_domain_name"`
	VpcAccountName               string `form:"vpc_account_name"`
	VpcID                        string `form:"vpc_id"`
	CustomizedRoutes             string `form:"customized_routes, omitempty" json:"customized_routes, omitempty"`
	Subnets                      string
	RouteTables                  string
	CustomizedRouteAdvertisement string
	DisableLocalRoutePropagation bool `form:"disable_local_route_propagation, omitempty" json:"disable_local_route_propagation, omitempty"`
}

type DomainListResp struct {
	Return  bool     `json:"return,omitempty"`
	Results []string `json:"results,omitempty"`
	Reason  string   `json:"reason,omitempty"`
}

func (c *Client) CreateAwsTgwVpcAttachment(awsTgwVpcAttachment *AwsTgwVpcAttachment) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for attach_vpc_to_tgw") + err.Error())
	}
	attachVpcFromTgw := url.Values{}
	attachVpcFromTgw.Add("CID", c.CID)
	attachVpcFromTgw.Add("action", "attach_vpc_to_tgw")
	attachVpcFromTgw.Add("region", awsTgwVpcAttachment.Region)
	attachVpcFromTgw.Add("vpc_account_name", awsTgwVpcAttachment.VpcAccountName)
	attachVpcFromTgw.Add("vpc_name", awsTgwVpcAttachment.VpcID)
	attachVpcFromTgw.Add("tgw_name", awsTgwVpcAttachment.TgwName)
	attachVpcFromTgw.Add("route_domain_name", awsTgwVpcAttachment.SecurityDomainName)
	if awsTgwVpcAttachment.DisableLocalRoutePropagation {
		attachVpcFromTgw.Add("disable_local_route_propagation", "yes")
	}
	if awsTgwVpcAttachment.CustomizedRoutes != "" {
		attachVpcFromTgw.Add("customized_routes", awsTgwVpcAttachment.CustomizedRoutes)
	}
	if awsTgwVpcAttachment.CustomizedRouteAdvertisement != "" {
		attachVpcFromTgw.Add("customized_route_advertisement", awsTgwVpcAttachment.CustomizedRouteAdvertisement)
	}
	if awsTgwVpcAttachment.Subnets != "" {
		attachVpcFromTgw.Add("subnet_list", awsTgwVpcAttachment.Subnets)
	}
	if awsTgwVpcAttachment.RouteTables != "" {
		attachVpcFromTgw.Add("route_table_list", awsTgwVpcAttachment.RouteTables)
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

func (c *Client) CreateAwsTgwVpcAttachmentForFireNet(awsTgwVpcAttachment *AwsTgwVpcAttachment) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for attach_vpc_to_tgw: ") + err.Error())
	}
	connectFireNetWithTgw := url.Values{}
	connectFireNetWithTgw.Add("CID", c.CID)
	connectFireNetWithTgw.Add("action", "connect_firenet_with_tgw")
	connectFireNetWithTgw.Add("vpc_id", awsTgwVpcAttachment.VpcID)
	connectFireNetWithTgw.Add("tgw_name", awsTgwVpcAttachment.TgwName)
	connectFireNetWithTgw.Add("domain_name", awsTgwVpcAttachment.SecurityDomainName)
	Url.RawQuery = connectFireNetWithTgw.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get connect_firenet_with_tgw failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode connect_firenet_with_tgw failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API connect_firenet_with_tgw Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetAwsTgwVpcAttachment(awsTgwVpcAttachment *AwsTgwVpcAttachment) (*AwsTgwVpcAttachment, error) {
	awsTgw := &AWSTgw{
		Name: awsTgwVpcAttachment.TgwName,
	}
	awsTgw, err := c.ListTgwDetails(awsTgw)
	if err != nil {
		return nil, fmt.Errorf("couldn't find AWS TGW: %s", awsTgwVpcAttachment.TgwName)
	}
	awsTgwVpcAttachment.Region = awsTgw.Region

	err = c.GetAwsTgwDomain(awsTgw, awsTgwVpcAttachment.SecurityDomainName)
	if err != nil {
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

	return aTVA, nil
}

func (c *Client) UpdateAwsTgwVpcAttachment(awsTgwVpcAttachment *AwsTgwVpcAttachment) error {
	return nil
}

func (c *Client) DeleteAwsTgwVpcAttachment(awsTgwVpcAttachment *AwsTgwVpcAttachment) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for detach_vpc_from_tgw") + err.Error())
	}
	detachVpcFromTgw := url.Values{}
	detachVpcFromTgw.Add("CID", c.CID)
	detachVpcFromTgw.Add("action", "detach_vpc_from_tgw")
	detachVpcFromTgw.Add("tgw_name", awsTgwVpcAttachment.TgwName)
	detachVpcFromTgw.Add("vpc_name", awsTgwVpcAttachment.VpcID)
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

func (c *Client) DeleteAwsTgwVpcAttachmentForFireNet(awsTgwVpcAttachment *AwsTgwVpcAttachment) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for disconnect_firenet_with_tgw") + err.Error())
	}
	disconnectFireNetWithTgw := url.Values{}
	disconnectFireNetWithTgw.Add("CID", c.CID)
	disconnectFireNetWithTgw.Add("action", "disconnect_firenet_with_tgw")
	disconnectFireNetWithTgw.Add("vpc_id", awsTgwVpcAttachment.VpcID)
	Url.RawQuery = disconnectFireNetWithTgw.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get disconnect_firenet_with_tgw failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode disconnect_firenet_with_tgw failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API disconnect_firenet_with_tgw Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetAwsTgwDetail(awsTgw *AWSTgw) (*AWSTgw, error) {
	awsTgw, err := c.ListTgwDetails(awsTgw)
	if err != nil {
		return nil, fmt.Errorf("couldn't find AWS TGW: %s", awsTgw.Name)
	}

	return awsTgw, nil
}

func (c *Client) GetAwsTgwDomain(awsTgw *AWSTgw, sDM string) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for list_route_domain_names") + err.Error())
	}
	listRouteDomainNames := url.Values{}
	listRouteDomainNames.Add("CID", c.CID)
	listRouteDomainNames.Add("action", "list_route_domain_names")
	listRouteDomainNames.Add("tgw_name", awsTgw.Name)
	Url.RawQuery = listRouteDomainNames.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get list_route_domain_names failed: " + err.Error())
	}
	data := DomainListResp{
		Return:  false,
		Results: make([]string, 0),
		Reason:  "",
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode list_route_domain_names failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API list_route_domain_names Get failed: " + data.Reason)
	}
	mDomain := make(map[string]bool)
	for i := range data.Results {
		mDomain[data.Results[i]] = true
	}
	if !mDomain[sDM] {
		return errors.New(awsTgw.Name + " does not have security domain: " + sDM)
	}

	return nil
}

func (c *Client) GetVPCAttachmentRouteTableDetails(awsTgwVpcAttachment *AwsTgwVpcAttachment) (*AwsTgwVpcAttachment, error) {
	Url, err := url.Parse(c.baseURL)
	viewRouteDomainDetails := url.Values{}
	viewRouteDomainDetails.Add("CID", c.CID)
	viewRouteDomainDetails.Add("action", "view_route_domain_details")
	viewRouteDomainDetails.Add("tgw_name", awsTgwVpcAttachment.TgwName)
	viewRouteDomainDetails.Add("route_domain_name", awsTgwVpcAttachment.SecurityDomainName)
	Url.RawQuery = viewRouteDomainDetails.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return nil, errors.New("HTTP Get view_route_domain_details failed: " + err.Error())
	}

	var data RouteDomainAPIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return awsTgwVpcAttachment, errors.New("Json Decode view_route_domain_details failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return awsTgwVpcAttachment, errors.New("Rest API view_route_domain_details Get failed: " + data.Reason)
	}
	routeDomainDetail := data.Results
	attachedVPCs := routeDomainDetail[0].AttachedVPC
	for i := range attachedVPCs {
		if attachedVPCs[i].VPCId == awsTgwVpcAttachment.VpcID {
			awsTgwVpcAttachment.VpcAccountName = attachedVPCs[i].AccountName
			return awsTgwVpcAttachment, nil
		}
	}
	return nil, ErrNotFound
}
