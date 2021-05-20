package goaviatrix

import (
	"fmt"
	"strings"
)

type DomainConn struct {
	Action      string `form:"action,omitempty"`
	CID         string `form:"CID,omitempty"`
	TgwName1    string `form:"tgw_name1,omitempty" json:"tgw_name1,omitempty"`
	DomainName1 string
	TgwName2    string `form:"tgw_name2,omitempty" json:"tgw_name2,omitempty"`
	DomainName2 string
}

type ListConnectedRouteDomainsResp struct {
	Return  bool                       `json:"return"`
	Results ConnectedRouteDomainDetail `json:"results"`
	Reason  string                     `json:"reason"`
}

type ConnectedRouteDomainDetail struct {
	ConnectedDomainNames    []string `json:"connected_domain_names"`
	NotConnectedDomainNames []string `json:"not_connected_domain_names"`
	Egress                  string   `json:"egress"`
}

func (c *Client) CreateDomainConn(domainConn *DomainConn) error {
	form := map[string]string{
		"CID":                           c.CID,
		"action":                        "add_connection_between_route_domains",
		"tgw_name":                      domainConn.TgwName1,
		"source_route_domain_name":      domainConn.DomainName1,
		"destination_route_domain_name": domainConn.TgwName2 + ":" + domainConn.DomainName2,
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) GetDomainConn(domainConn *DomainConn) error {
	var data ListConnectedRouteDomainsResp
	form := map[string]string{
		"CID":               c.CID,
		"action":            "list_connected_route_domains",
		"tgw_name":          domainConn.TgwName1,
		"route_domain_name": domainConn.DomainName1,
	}
	check := func(action, reason string, ret bool) error {
		if !ret {
			if strings.Contains(data.Reason, "does not exist") {
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s Post failed: %s", action, reason)
		}
		return nil
	}
	err := c.GetAPI(&data, form["action"], form, check)
	if err != nil {
		return err
	}
	connectedDomains := data.Results.ConnectedDomainNames
	for i := range connectedDomains {
		if connectedDomains[i] == domainConn.TgwName2+":"+domainConn.DomainName2 {
			return nil
		}
	}
	return ErrNotFound
}

func (c *Client) DeleteDomainConn(domainConn *DomainConn) error {
	form := map[string]string{
		"CID":                           c.CID,
		"action":                        "delete_connection_between_route_domains",
		"tgw_name":                      domainConn.TgwName1,
		"source_route_domain_name":      domainConn.DomainName1,
		"destination_route_domain_name": domainConn.TgwName2 + ":" + domainConn.DomainName2,
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}
