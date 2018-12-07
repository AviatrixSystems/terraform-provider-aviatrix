package goaviatrix

import (
	"encoding/json"
	"errors"
)

// AwsTGW simple struct to hold aws_tgw details
type SecurityDomain struct {
	Action           string `form:"action,omitempty"`
	CID              string `form:"CID,omitempty"`
	Name 			 string `form:"route_domain_name,omitempty"`
	AccountName      string `form:"account_name,omitempty"`
	Region           string `form:"region,omitempty"`
	AwsTgwName 		 string `form:"tgw_name,omitempty"`
}

type SecurityDomainAPIResp struct {
	Return  bool   	 `json:"return"`
	Results []string `json:"results"`
	Reason  string   `json:"reason"`
}

type SecurityDomainList struct {
	Return  bool   			 `json:"return"`
	Results []SecurityDomain `json:"results"`
	Reason  string 			 `json:"reason"`
}

func (c *Client) CreateSecurityDomain(securityDomain *SecurityDomain) (error) {
	securityDomain.CID = c.CID
	securityDomain.Action = "add_route_domain"
	resp, err := c.Post(c.baseURL, securityDomain)
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

func (c *Client) GetSecurityDomain(securityDomain *SecurityDomain) (string, error) {
	securityDomain.CID = c.CID
	securityDomain.Action = "list_route_domain_names"

	resp, err := c.Post(c.baseURL, securityDomain)
	if err != nil {
		return "", err
	}

	data := SecurityDomainAPIResp{
		Return:  false,
		Results: make([]string, 0),
		Reason:  "",
	}
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if !data.Return {
		return "", errors.New(data.Reason)
	}

	securityDomainList := data.Results
	for i := range securityDomainList {
		if securityDomainList[i] == securityDomain.Name {
			return securityDomainList[i], nil
		}
	}

	return "", ErrNotFound
}

func (c *Client) UpdateSecurityDomain(securityDomain *SecurityDomain) (error) {
	return nil
}

func (c *Client) DeleteSecurityDomain(securityDomain *SecurityDomain) (error) {
	securityDomain.CID = c.CID
	securityDomain.Action = "delete_route_domain"
	resp, err := c.Post(c.baseURL, securityDomain)
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