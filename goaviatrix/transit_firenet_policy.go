package goaviatrix

import (
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
	form := map[string]string{
		"CID":                  c.CID,
		"action":               "add_spoke_to_transit_firenet_inspection",
		"firenet_gateway_name": transitFireNetPolicy.TransitFireNetGatewayName,
		"spoke_gateway_name":   transitFireNetPolicy.InspectedResourceName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) GetTransitFireNetPolicy(transitFireNetPolicy *TransitFireNetPolicy) error {
	form := map[string]string{
		"CID":    c.CID,
		"action": "list_transit_firenet_spoke_policies",
	}

	var data TransitFireNetPolicyAPIResp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return err
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
	form := map[string]string{
		"CID":                  c.CID,
		"action":               "delete_spoke_from_transit_firenet_inspection",
		"firenet_gateway_name": transitFireNetPolicy.TransitFireNetGatewayName,
		"spoke_gateway_name":   transitFireNetPolicy.InspectedResourceName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}
