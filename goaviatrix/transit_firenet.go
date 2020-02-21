package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"strings"
)

type TransitFireNetInspection struct {
	TransitFireNetGatewayName string `form:"gateway_1,omitempty" json:"gateway_1,omitempty"`
	InspectedResourceName     string `form:"gateway_2,omitempty" json:"gateway_2,omitempty"`
}

type TransitFireNetInspectionAPIResp struct {
	Return  bool                           `json:"return"`
	Results []TransitFireNetInspectionEdit `json:"results"`
	Reason  string                         `json:"reason"`
}

type TransitFireNetInspectionEdit struct {
	TransitFireNetGwName      string   `json:"gw_name,omitempty"`
	InspectedResourceNameList []string `json:"inspected,omitempty"`
}

func (c *Client) CreateTransitFireNetInspection(transitFireNetInspection *TransitFireNetInspection) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'CreateTransitFireNetInspection': ") + err.Error())
	}
	CreateTransitFireNetInspection := url.Values{}
	CreateTransitFireNetInspection.Add("CID", c.CID)
	CreateTransitFireNetInspection.Add("action", "add_spoke_to_transit_firenet_inspection")
	CreateTransitFireNetInspection.Add("firenet_gateway_name", transitFireNetInspection.TransitFireNetGatewayName)
	CreateTransitFireNetInspection.Add("spoke_gateway_name", transitFireNetInspection.InspectedResourceName)
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

func (c *Client) GetTransitFireNetInspection(transitFireNetInspection *TransitFireNetInspection) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for GetTransitFireNetInspection: ") + err.Error())
	}
	listTransitFireNetSpokePolicies := url.Values{}
	listTransitFireNetSpokePolicies.Add("CID", c.CID)
	listTransitFireNetSpokePolicies.Add("action", "list_transit_firenet_spoke_policies")
	Url.RawQuery = listTransitFireNetSpokePolicies.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get list_transit_firenet_spoke_policies failed: " + err.Error())
	}
	var data TransitFireNetInspectionAPIResp
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
		log.Printf("Transit gateway peering with transit firenet gateway: %s and inspected resource name: %s not found",
			transitFireNetInspection.TransitFireNetGatewayName, transitFireNetInspection.InspectedResourceName)
		return ErrNotFound
	}
	inspectionList := data.Results
	for i := range inspectionList {
		if inspectionList[i].TransitFireNetGwName != transitFireNetInspection.TransitFireNetGatewayName {
			continue
		}
		for j := range inspectionList[i].InspectedResourceNameList {
			if inspectionList[i].InspectedResourceNameList[j] == transitFireNetInspection.InspectedResourceName {
				return nil
			}
		}
	}
	return ErrNotFound
}

func (c *Client) DeleteTransitFireNetInspection(transitFireNetInspection *TransitFireNetInspection) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'DeleteTransitFireNetInspection': ") + err.Error())
	}
	CreateTransitFireNetInspection := url.Values{}
	CreateTransitFireNetInspection.Add("CID", c.CID)
	CreateTransitFireNetInspection.Add("action", "delete_spoke_from_transit_firenet_inspection")
	CreateTransitFireNetInspection.Add("firenet_gateway_name", transitFireNetInspection.TransitFireNetGatewayName)
	CreateTransitFireNetInspection.Add("spoke_gateway_name", transitFireNetInspection.InspectedResourceName)
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
