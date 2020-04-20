package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
)

type DomainConn struct {
	Action      string `form:"action,omitempty"`
	CID         string `form:"CID,omitempty"`
	TgwName1    string `form:"tgw_name1,omitempty" json:"tgw_name1,omitempty"`
	DomainName1 string `form:"tgw_name2,omitempty" json:"tgw_name2,omitempty"`
	TgwName2    string `form:"tgw_name1,omitempty" json:"tgw_name1,omitempty"`
	DomainName2 string `form:"tgw_name2,omitempty" json:"tgw_name2,omitempty"`
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
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'add_connection_between_route_domains': ") + err.Error())
	}
	addConnectionBetweenRouteDomains := url.Values{}
	addConnectionBetweenRouteDomains.Add("CID", c.CID)
	addConnectionBetweenRouteDomains.Add("action", "add_connection_between_route_domains")
	addConnectionBetweenRouteDomains.Add("tgw_name", domainConn.TgwName1)
	addConnectionBetweenRouteDomains.Add("source_route_domain_name", domainConn.DomainName1)
	addConnectionBetweenRouteDomains.Add("destination_route_domain_name", domainConn.TgwName2+":"+domainConn.DomainName2)
	Url.RawQuery = addConnectionBetweenRouteDomains.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get 'add_connection_between_route_domains' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'add_connection_between_route_domains' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'add_connection_between_route_domains' Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetDomainConn(domainConn *DomainConn) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'list_connected_route_domains': ") + err.Error())
	}
	listConnectedRouteDomains := url.Values{}
	listConnectedRouteDomains.Add("CID", c.CID)
	listConnectedRouteDomains.Add("action", "list_connected_route_domains")
	listConnectedRouteDomains.Add("tgw_name", domainConn.TgwName1)
	listConnectedRouteDomains.Add("route_domain_name", domainConn.DomainName1)
	Url.RawQuery = listConnectedRouteDomains.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get 'list_connected_route_domains' failed: " + err.Error())
	}
	var data ListConnectedRouteDomainsResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'list_connected_route_domains' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'list_connected_route_domains' Get failed: " + data.Reason)
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
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'delete_connection_between_route_domains': ") + err.Error())
	}
	deleteConnectionBetweenRouteDomains := url.Values{}
	deleteConnectionBetweenRouteDomains.Add("CID", c.CID)
	deleteConnectionBetweenRouteDomains.Add("action", "delete_connection_between_route_domains")
	deleteConnectionBetweenRouteDomains.Add("tgw_name", domainConn.TgwName1)
	deleteConnectionBetweenRouteDomains.Add("source_route_domain_name", domainConn.DomainName1)
	deleteConnectionBetweenRouteDomains.Add("destination_route_domain_name", domainConn.TgwName2+":"+domainConn.DomainName2)
	Url.RawQuery = deleteConnectionBetweenRouteDomains.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get 'delete_connection_between_route_domains' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'delete_connection_between_route_domains' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'delete_connection_between_route_domains' Get failed: " + data.Reason)
	}
	return nil
}
