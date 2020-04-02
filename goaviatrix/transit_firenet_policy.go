package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

type TransitFireNetPolicy struct {
	TransitFireNetGatewayName string `form:"gateway_1,omitempty" json:"gateway_1,omitempty"`
	InspectedResourceName     string `form:"gateway_2,omitempty" json:"gateway_2,omitempty"`
}

type TransitFireNetPolicyAPIResp struct {
	Return  bool                       `json:"return"`
	Results []TransitFireNetPolicyEdit `json:"results"`
	Reason  string                     `json:"reason"`
}

type TransitFireNetPolicyEdit struct {
	TransitFireNetGwName         string   `json:"gw_name,omitempty"`
	InspectedResourceNameList    []string `json:"inspected,omitempty"`
	ManagementAccessResourceName string   `json:"management_access,omitempty"`
}

func (c *Client) CreateTransitFireNetPolicy(transitFireNetPolicy *TransitFireNetPolicy) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'CreateTransitFireNetInspection': ") + err.Error())
	}
	CreateTransitFireNetInspection := url.Values{}
	CreateTransitFireNetInspection.Add("CID", c.CID)
	CreateTransitFireNetInspection.Add("action", "add_spoke_to_transit_firenet_inspection")
	CreateTransitFireNetInspection.Add("firenet_gateway_name", transitFireNetPolicy.TransitFireNetGatewayName)
	CreateTransitFireNetInspection.Add("spoke_gateway_name", transitFireNetPolicy.InspectedResourceName)
	Url.RawQuery = CreateTransitFireNetInspection.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get add_spoke_to_transit_firenet_inspection failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode add_spoke_to_transit_firenet_inspection failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API add_spoke_to_transit_firenet_inspection Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetTransitFireNetPolicy(transitFireNetPolicy *TransitFireNetPolicy) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for GetTransitFireNetPolicy: ") + err.Error())
	}
	listTransitFireNetSpokePolicies := url.Values{}
	listTransitFireNetSpokePolicies.Add("CID", c.CID)
	listTransitFireNetSpokePolicies.Add("action", "list_transit_firenet_spoke_policies")
	Url.RawQuery = listTransitFireNetSpokePolicies.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get list_transit_firenet_spoke_policies failed: " + err.Error())
	}
	var data TransitFireNetPolicyAPIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode list_transit_firenet_spoke_policies failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API list_transit_firenet_spoke_policies Get failed: " + data.Reason)
	}
	if len(data.Results) == 0 {
		log.Errorf("transit firenet policy between transit firenet gateway: %s and inspected resource name: %s not found",
			transitFireNetPolicy.TransitFireNetGatewayName, transitFireNetPolicy.InspectedResourceName)
		return ErrNotFound
	}
	policyList := data.Results
	for i := range policyList {
		if policyList[i].TransitFireNetGwName != transitFireNetPolicy.TransitFireNetGatewayName {
			continue
		}
		for j := range policyList[i].InspectedResourceNameList {
			if policyList[i].InspectedResourceNameList[j] == transitFireNetPolicy.InspectedResourceName {
				return nil
			}
		}
	}
	return ErrNotFound
}

func (c *Client) DeleteTransitFireNetPolicy(transitFireNetPolicy *TransitFireNetPolicy) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'DeleteTransitFireNetPolicy': ") + err.Error())
	}
	CreateTransitFireNetInspection := url.Values{}
	CreateTransitFireNetInspection.Add("CID", c.CID)
	CreateTransitFireNetInspection.Add("action", "delete_spoke_from_transit_firenet_inspection")
	CreateTransitFireNetInspection.Add("firenet_gateway_name", transitFireNetPolicy.TransitFireNetGatewayName)
	CreateTransitFireNetInspection.Add("spoke_gateway_name", transitFireNetPolicy.InspectedResourceName)
	Url.RawQuery = CreateTransitFireNetInspection.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get delete_spoke_from_transit_firenet_inspection failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode delete_spoke_from_transit_firenet_inspection failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API delete_spoke_from_transit_firenet_inspection Get failed: " + data.Reason)
	}
	return nil
}
