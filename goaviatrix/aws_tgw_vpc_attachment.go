package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

type AwsTgwVpcAttachment struct {
	Action             string `form:"action,omitempty"`
	CID                string `form:"CID,omitempty"`
	TgwName            string `form:"tgw_name"`
	Region             string `form:"region"`
	SecurityDomainName string `form:"security_domain_name"`
	VpcAccountName     string `form:"vpc_account_name"`
	VpcID              string `form:"vpc_id"`
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
	Url.RawQuery = attachVpcFromTgw.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get attach_vpc_to_tgw failed: " + err.Error())
	}

	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode attach_vpc_to_tgw failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API attach_vpc_to_tgw Get failed: " + data.Reason)
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

	aTVA, err := c.GetAwsTgwDomainAttachedVpc(awsTgwVpcAttachment)
	if err != nil {
		if err == ErrNotFound {
			return nil, err
		}
		return nil, errors.New("could not get security domain details: " + err.Error())
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
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode detach_vpc_from_tgw failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API detach_vpc_from_tgw Get failed: " + data.Reason)
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
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode list_route_domain_names failed: " + err.Error())
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

func (c *Client) GetAwsTgwDomainAttachedVpc(awsTgwVpcAttachment *AwsTgwVpcAttachment) (*AwsTgwVpcAttachment, error) {
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
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return awsTgwVpcAttachment, errors.New("Json Decode view_route_domain_details failed: " + err.Error())
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
