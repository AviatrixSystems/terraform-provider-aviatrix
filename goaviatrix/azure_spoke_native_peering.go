package goaviatrix

import (
	"strings"
)

type AzureSpokeNativePeering struct {
	CID                string `form:"CID,omitempty"`
	Action             string `form:"action,omitempty"`
	TransitGatewayName string `form:"transit_gateway_name,omitempty"`
	SpokeAccountName   string `form:"account_name,omitempty"`
	SpokeRegion        string `form:"region,omitempty"`
	SpokeVpcID         string `form:"vpc_id,omitempty"`
}

type AzureSpokeNativePeeringAPIResp struct {
	Return  bool                          `json:"return"`
	Results []AzureSpokeNativePeeringEdit `json:"results"`
	Reason  string                        `json:"reason"`
}

type AzureSpokeNativePeeringEdit struct {
	AccountName string `json:"account_original_name"`
	Name        string `json:"name"`
	Region      string `json:"region"`
	VpcID       string `json:"vpc_id"`
}

func (c *Client) CreateAzureSpokeNativePeering(azureSpokeNativePeering *AzureSpokeNativePeering) error {
	azureSpokeNativePeering.CID = c.CID
	azureSpokeNativePeering.Action = "attach_arm_native_spoke_to_transit"
	return c.PostAPI(azureSpokeNativePeering.Action, azureSpokeNativePeering, BasicCheck)
}

func (c *Client) GetAzureSpokeNativePeering(azureSpokeNativePeering *AzureSpokeNativePeering) (*AzureSpokeNativePeering, error) {
	var data AzureSpokeNativePeeringAPIResp
	form := map[string]string{
		"CID":                  c.CID,
		"action":               "list_arm_native_spokes",
		"transit_gateway_name": azureSpokeNativePeering.TransitGatewayName,
		"details":              "true",
	}
	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}
	if len(data.Results) == 0 {
		return nil, ErrNotFound
	}
	peeringList := data.Results
	for i := range peeringList {
		if peeringList[i].Name == "" || len(strings.Split(peeringList[i].Name, ":")) != 3 {
			continue
		}
		spokeAccountName := peeringList[i].AccountName
		spokeVpcID := peeringList[i].VpcID
		if azureSpokeNativePeering.SpokeAccountName != spokeAccountName || azureSpokeNativePeering.SpokeVpcID != spokeVpcID {
			continue
		}
		azureSpokeNativePeering.SpokeRegion = peeringList[i].Region
		return azureSpokeNativePeering, nil
	}
	return nil, ErrNotFound
}

func (c *Client) DeleteAzureSpokeNativePeering(azureSpokeNativePeering *AzureSpokeNativePeering) error {
	form := map[string]string{
		"CID":                  c.CID,
		"action":               "detach_arm_native_spoke_to_transit",
		"transit_gateway_name": azureSpokeNativePeering.TransitGatewayName,
		"spoke_name":           azureSpokeNativePeering.SpokeAccountName + ":" + strings.Replace(azureSpokeNativePeering.SpokeVpcID, ".", "-", -1),
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}
