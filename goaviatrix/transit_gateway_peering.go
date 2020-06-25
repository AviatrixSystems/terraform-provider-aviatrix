package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

type TransitGatewayPeering struct {
	TransitGatewayName1                 string `form:"gateway1,omitempty" json:"gateway_1,omitempty"`
	TransitGatewayName2                 string `form:"gateway2,omitempty" json:"gateway_2,omitempty"`
	Gateway1ExcludedCIDRs               string `form:"source_filter_cidrs,omitempty"`
	Gateway2ExcludedCIDRs               string `form:"destination_filter_cidrs,omitempty"`
	Gateway1ExcludedTGWConnections      string `form:"source_exclude_connections,omitempty"`
	Gateway2ExcludedTGWConnections      string `form:"destination_exclude_connections,omitempty"`
	Gateway1ExcludedCIDRsSlice          []string
	Gateway2ExcludedCIDRsSlice          []string
	Gateway1ExcludedTGWConnectionsSlice []string
	Gateway2ExcludedTGWConnectionsSlice []string
	CID                                 string `form:"CID,omitempty"`
	Action                              string `form:"action,omitempty"`
}

type TransitGatewayPeeringAPIResp struct {
	Return  bool                      `json:"return"`
	Results [][]TransitGatewayPeering `json:"results"`
	Reason  string                    `json:"reason"`
}

type TransitGatewayPeeringDetailsAPIResp struct {
	Return  bool                                `json:"return"`
	Results TransitGatewayPeeringDetailsResults `json:"results"`
	Reason  string                              `json:"reason"`
}

type TransitGatewayPeeringDetailsResults struct {
	Site1 TransitGatewayPeeringDetail `json:"site_1"`
	Site2 TransitGatewayPeeringDetail `json:"site_2"`
}

type TransitGatewayPeeringDetail struct {
	ExcludedCIDRs          []string `json:"exclude_filter_list"`
	ExcludedTGWConnections []string `json:"exclude_connections"`
}

func (c *Client) CreateTransitGatewayPeering(transitGatewayPeering *TransitGatewayPeering) error {
	transitGatewayPeering.CID = c.CID
	transitGatewayPeering.Action = "create_inter_transit_gateway_peering"
	resp, err := c.Post(c.baseURL, transitGatewayPeering)
	if err != nil {
		return errors.New("HTTP POST create_inter_transit_gateway_peering failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode create_inter_transit_gateway_peering failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API create_inter_transit_gateway_peering Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetTransitGatewayPeering(transitGatewayPeering *TransitGatewayPeering) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for list_inter_transit_gateway_peering ") + err.Error())
	}
	listInterTransitGwPeering := url.Values{}
	listInterTransitGwPeering.Add("CID", c.CID)
	listInterTransitGwPeering.Add("action", "list_inter_transit_gateway_peering")
	Url.RawQuery = listInterTransitGwPeering.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get list_inter_transit_gateway_peering failed: " + err.Error())
	}
	var data TransitGatewayPeeringAPIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode list_inter_transit_gateway_peering failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API list_inter_transit_gateway_peering Get failed: " + data.Reason)
	}
	if len(data.Results) == 0 {
		log.Errorf("Transit gateway peering with gateways %s and %s not found",
			transitGatewayPeering.TransitGatewayName1, transitGatewayPeering.TransitGatewayName2)
		return ErrNotFound
	}
	peeringList := data.Results
	for i := range peeringList {
		for j := range peeringList[i] {
			if peeringList[i][j].TransitGatewayName1 == transitGatewayPeering.TransitGatewayName1 &&
				peeringList[i][j].TransitGatewayName2 == transitGatewayPeering.TransitGatewayName2 ||
				peeringList[i][j].TransitGatewayName1 == transitGatewayPeering.TransitGatewayName2 &&
					peeringList[i][j].TransitGatewayName2 == transitGatewayPeering.TransitGatewayName1 {
				log.Debugf("Found %s<->%s transit gateway peering: %#v",
					transitGatewayPeering.TransitGatewayName1,
					transitGatewayPeering.TransitGatewayName2, peeringList[i][j])
				return nil
			}
		}
	}
	return ErrNotFound
}

func (c *Client) GetTransitGatewayPeeringDetails(transitGatewayPeering *TransitGatewayPeering) (*TransitGatewayPeering, error) {
	transitGatewayPeering.CID = c.CID
	transitGatewayPeering.Action = "get_inter_transit_gateway_peering_details"
	resp, err := c.Post(c.baseURL, transitGatewayPeering)
	if err != nil {
		return nil, errors.New("HTTP POST get_inter_transit_gateway_peering_details failed: " + err.Error())
	}
	var data TransitGatewayPeeringDetailsAPIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode get_inter_transit_gateway_peering_details failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return nil, errors.New("Rest API get_inter_transit_gateway_peering_details Get failed: " + data.Reason)
	}

	transitGatewayPeering.Gateway1ExcludedCIDRsSlice = data.Results.Site1.ExcludedCIDRs
	transitGatewayPeering.Gateway1ExcludedTGWConnectionsSlice = data.Results.Site1.ExcludedTGWConnections
	transitGatewayPeering.Gateway2ExcludedCIDRsSlice = data.Results.Site2.ExcludedCIDRs
	transitGatewayPeering.Gateway2ExcludedTGWConnectionsSlice = data.Results.Site2.ExcludedTGWConnections

	return transitGatewayPeering, nil
}

func (c *Client) UpdateTransitGatewayPeering(transitGatewayPeering *TransitGatewayPeering) error {
	transitGatewayPeering.CID = c.CID
	transitGatewayPeering.Action = "edit_inter_transit_gateway_peering"
	resp, err := c.Post(c.baseURL, transitGatewayPeering)
	if err != nil {
		return errors.New("HTTP POST edit_inter_transit_gateway_peering failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode edit_inter_transit_gateway_peering failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API edit_inter_transit_gateway_peering Get failed: " + data.Reason)
	}

	return nil
}

func (c *Client) DeleteTransitGatewayPeering(transitGatewayPeering *TransitGatewayPeering) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for delete_inter_transit_gateway_peering ") + err.Error())
	}
	deleteInterTransitGwPeering := url.Values{}
	deleteInterTransitGwPeering.Add("CID", c.CID)
	deleteInterTransitGwPeering.Add("action", "delete_inter_transit_gateway_peering")
	deleteInterTransitGwPeering.Add("gateway1", transitGatewayPeering.TransitGatewayName1)
	deleteInterTransitGwPeering.Add("gateway2", transitGatewayPeering.TransitGatewayName2)
	Url.RawQuery = deleteInterTransitGwPeering.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get delete_inter_transit_gateway_peering failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode delete_inter_transit_gateway_peering failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API delete_inter_transit_gateway_peering Get failed: " + data.Reason)
	}
	return nil
}
