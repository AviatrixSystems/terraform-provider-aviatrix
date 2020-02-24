package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"strings"
)

type FirewallManagementAccess struct {
	CID                          string `form:"CID,omitempty"`
	Action                       string `form:"action,omitempty"`
	TransitFireNetGatewayName    string `form:"transit_firenet_gateway_name,omitempty" json:"gw_name,omitempty"`
	ManagementAccessResourceName string `form:"management_access,omitempty" json:"management_access,omitempty"`
}

type FirewallManagementAccessAPIResp struct {
	Return  bool                       `json:"return"`
	Results []FirewallManagementAccess `json:"results"`
	Reason  string                     `json:"reason"`
}

func (c *Client) CreateFirewallManagementAccess(firewallManagementAccess *FirewallManagementAccess) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'CreateFirewallManagementAccess': ") + err.Error())
	}
	editTransitFireNetManagementAccess := url.Values{}
	editTransitFireNetManagementAccess.Add("CID", c.CID)
	editTransitFireNetManagementAccess.Add("action", "edit_transit_firenet_management_access")
	editTransitFireNetManagementAccess.Add("gateway_name", firewallManagementAccess.TransitFireNetGatewayName)
	editTransitFireNetManagementAccess.Add("management_access", firewallManagementAccess.ManagementAccessResourceName)
	Url.RawQuery = editTransitFireNetManagementAccess.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get 'edit_transit_firenet_management_access' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'edit_transit_firenet_management_access' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'edit_transit_firenet_management_access' Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetFirewallManagementAccess(firewallManagementAccess *FirewallManagementAccess) (*FirewallManagementAccess, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for GetFirewallManagementAccess: ") + err.Error())
	}
	listTransitFireNetSpokePolicies := url.Values{}
	listTransitFireNetSpokePolicies.Add("CID", c.CID)
	listTransitFireNetSpokePolicies.Add("action", "list_transit_firenet_spoke_policies")
	Url.RawQuery = listTransitFireNetSpokePolicies.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return nil, errors.New("HTTP Get list_transit_firenet_spoke_policies failed: " + err.Error())
	}
	var data FirewallManagementAccessAPIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_transit_firenet_spoke_policies failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return nil, errors.New("Rest API list_transit_firenet_spoke_policies Get failed: " + data.Reason)
	}
	if len(data.Results) == 0 {
		log.Printf("Transit gateway peering with transit firenet gateway: %s and inspected resource name: %s not found",
			firewallManagementAccess.TransitFireNetGatewayName, firewallManagementAccess.ManagementAccessResourceName)
		return nil, ErrNotFound
	}
	firewallManagementAccessList := data.Results
	for i := range firewallManagementAccessList {
		if firewallManagementAccessList[i].TransitFireNetGatewayName != firewallManagementAccess.TransitFireNetGatewayName {
			continue
		}
		if firewallManagementAccessList[i].ManagementAccessResourceName == "no" {
			return nil, ErrNotFound
		}
		firewallManagementAccess.ManagementAccessResourceName = firewallManagementAccessList[i].ManagementAccessResourceName
		return firewallManagementAccess, nil
	}
	return nil, ErrNotFound
}

func (c *Client) DestroyFirewallManagementAccess(firewallManagementAccess *FirewallManagementAccess) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'DestroyFirewallManagementAccess': ") + err.Error())
	}
	editTransitFireNetManagementAccess := url.Values{}
	editTransitFireNetManagementAccess.Add("CID", c.CID)
	editTransitFireNetManagementAccess.Add("action", "edit_transit_firenet_management_access")
	editTransitFireNetManagementAccess.Add("gateway_name", firewallManagementAccess.TransitFireNetGatewayName)
	editTransitFireNetManagementAccess.Add("management_access", firewallManagementAccess.ManagementAccessResourceName)
	Url.RawQuery = editTransitFireNetManagementAccess.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get 'edit_transit_firenet_management_access' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'edit_transit_firenet_management_access' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'edit_transit_firenet_management_access' Get failed: " + data.Reason)
	}
	return nil
}
