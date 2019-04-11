package goaviatrix

import (
	"encoding/json"
	"errors"
	"net/url"
)

// AwsTGW simple struct to hold aws_tgw details
type SecurityDomain struct {
	Action      string `form:"action, omitempty"`
	CID         string `form:"CID, omitempty"`
	Name        string `form:"route_domain_name, omitempty"`
	AccountName string `form:"account_name, omitempty"`
	Region      string `form:"region, omitempty"`
	AwsTgwName  string `form:"tgw_name, omitempty"`
}

type SecurityDomainAPIResp struct {
	Return  bool     `json:"return"`
	Results []string `json:"results"`
	Reason  string   `json:"reason"`
}

type SecurityDomainList struct {
	Return  bool     `json:"return"`
	Results []string `json:"results"`
	Reason  string   `json:"reason"`
}

type SecurityDomainRule struct {
	Name            string    `json:"security_domain_name, omitempty"`
	ConnectedDomain []string  `json:"connected_domains, omitempty"`
	AttachedVPCs    []VPCSolo `json:"attached_vpc, omitempty"`
}

type VPCSolo struct {
	Region      string `json:"vpc_region, omitempty"`
	AccountName string `json:"vpc_account_name, omitempty"`
	VpcID       string `json:"vpc_id, omitempty"`
}

func (c *Client) CreateSecurityDomain(securityDomain *SecurityDomain) error {
	securityDomain.CID = c.CID
	securityDomain.Action = "add_route_domain"
	resp, err := c.Post(c.baseURL, securityDomain)
	if err != nil {
		return errors.New("HTTP Post add_route_domain failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode add_route_domain failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API add_route_domain Post failed: " + data.Reason)
	}

	return nil
}

func (c *Client) GetSecurityDomain(securityDomain *SecurityDomain) (string, error) {
	securityDomain.CID = c.CID
	securityDomain.Action = "list_route_domain_names"
	resp, err := c.Post(c.baseURL, securityDomain)
	if err != nil {
		return "", errors.New("HTTP Post list_route_domain_names failed: " + err.Error())
	}

	data := SecurityDomainAPIResp{
		Return:  false,
		Results: make([]string, 0),
		Reason:  "",
	}
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", errors.New("Json Decode list_route_domain_names failed: " + err.Error())
	}
	if !data.Return {
		return "", errors.New("Rest API list_route_domain_names Post failed: " + data.Reason)
	}

	securityDomainList := data.Results
	for i := range securityDomainList {
		if securityDomainList[i] == securityDomain.Name {
			return securityDomainList[i], nil
		}
	}

	return "", ErrNotFound
}

func (c *Client) UpdateSecurityDomain(securityDomain *SecurityDomain) error {
	return nil
}

func (c *Client) DeleteSecurityDomain(securityDomain *SecurityDomain) error {
	securityDomain.CID = c.CID
	securityDomain.Action = "delete_route_domain"
	resp, err := c.Post(c.baseURL, securityDomain)
	if err != nil {
		return errors.New("HTTP Post delete_route_domain failed: " + err.Error())
	}

	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode delete_route_domain failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API delete_route_domain Post failed: " + data.Reason)
	}

	return nil
}

func (c *Client) CreateDomainConnection(awsTgw *AWSTgw, sourceDomain string, destinationDomain string) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for add_connection_between_route_domains") + err.Error())
	}
	addConnectionBetweenRouteDomains := url.Values{}
	addConnectionBetweenRouteDomains.Add("CID", c.CID)
	addConnectionBetweenRouteDomains.Add("action", "add_connection_between_route_domains")
	addConnectionBetweenRouteDomains.Add("account_name", awsTgw.AccountName)
	addConnectionBetweenRouteDomains.Add("region", awsTgw.Region)
	addConnectionBetweenRouteDomains.Add("tgw_name", awsTgw.Name)
	addConnectionBetweenRouteDomains.Add("source_route_domain_name", sourceDomain)
	addConnectionBetweenRouteDomains.Add("destination_route_domain_name", destinationDomain)
	Url.RawQuery = addConnectionBetweenRouteDomains.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get add_connection_between_route_domains failed: " + err.Error())
	}

	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode add_connection_between_route_domains failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API add_connection_between_route_domains Get failed: " + data.Reason)
	}

	return nil
}

func (c *Client) DeleteDomainConnection(awsTgw *AWSTgw, sourceDomain string, destinationDomain string) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for delete_connection_between_route_domains") + err.Error())
	}
	deleteConnectionBetweenRouteDomains := url.Values{}
	deleteConnectionBetweenRouteDomains.Add("CID", c.CID)
	deleteConnectionBetweenRouteDomains.Add("action", "delete_connection_between_route_domains")
	deleteConnectionBetweenRouteDomains.Add("account_name", awsTgw.AccountName)
	deleteConnectionBetweenRouteDomains.Add("region", awsTgw.Region)
	deleteConnectionBetweenRouteDomains.Add("tgw_name", awsTgw.Name)
	deleteConnectionBetweenRouteDomains.Add("source_route_domain_name", sourceDomain)
	deleteConnectionBetweenRouteDomains.Add("destination_route_domain_name", destinationDomain)
	Url.RawQuery = deleteConnectionBetweenRouteDomains.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get delete_connection_between_route_domains failed: " + err.Error())
	}

	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode delete_connection_between_route_domains failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API delete_connection_between_route_domains Get failed: " + data.Reason)
	}

	return nil
}
