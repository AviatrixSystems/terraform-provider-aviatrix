package goaviatrix

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
)

type TransitGatewayPeering struct {
	//Name                string `form:"name,omitempty" json:"name,omitempty"`
	TransitGatewayName1 string `form:"gateway_1,omitempty" json:"gateway_1,omitempty"`
	TransitGatewayName2 string `form:"gateway_2,omitempty" json:"gateway_2,omitempty"`
	//OverAwsPeering      string `form:"over_aws_peering,omitempty" json:"over_aws_peering,omitempty"`
	//PeeringHaStatus     string `form:"peering_ha_status,omitempty" json:"peering_ha_status,omitempty"`
	//HighPerf            string `form:"high_perf,omitempty" json:"high_perf,omitempty"`
	//TimeStamp           string `form:"timestamp,omitempty" json:"timestamp,omitempty"`
	//State               string `form:"state,omitempty" json:"state,omitempty"`
	//Modified            string `form:"modified,omitempty" json:"modified,omitempty"`
	//Expiration          string `form:"expiration,omitempty" json:"expiration,omitempty"`
	//TxBytes             string `form:"tx_bytes,omitempty" json:"tx_bytes,omitempty"`
	//LicenseId           string `form:"license_id,omitempty" json:"license_id,omitempty"`
	//RxBytes             string `form:"rx_bytes,omitempty" json:"rx_bytes,omitempty"`
}

type TransitGatewayPeeringAPIResp struct {
	Return  bool                      `json:"return"`
	Results [][]TransitGatewayPeering `json:"results"`
	Reason  string                    `json:"reason"`
}

func (c *Client) CreateTransitGatewayPeering(transitGatewayPeering *TransitGatewayPeering) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for create_inter_transit_gateway_peering ") + err.Error())
	}
	createInterTransitGwPeering := url.Values{}
	createInterTransitGwPeering.Add("CID", c.CID)
	createInterTransitGwPeering.Add("action", "create_inter_transit_gateway_peering")
	createInterTransitGwPeering.Add("gateway1", transitGatewayPeering.TransitGatewayName1)
	createInterTransitGwPeering.Add("gateway2", transitGatewayPeering.TransitGatewayName2)
	Url.RawQuery = createInterTransitGwPeering.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get create_inter_transit_gateway_peering failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode create_inter_transit_gateway_peering failed: " + err.Error())
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
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode list_inter_transit_gateway_peering failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API list_inter_transit_gateway_peering Get failed: " + data.Reason)
	}
	if len(data.Results) == 0 {
		log.Printf("Transit gateway peering with gateways %s and %s not found",
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
				log.Printf("[DEBUG] Found %s<->%s transit gateway peering: %#v",
					transitGatewayPeering.TransitGatewayName1,
					transitGatewayPeering.TransitGatewayName2, peeringList[i][j])
				return nil
			}
		}
	}
	return ErrNotFound
}

func (c *Client) UpdateTransitGatewayPeering(transitGatewayPeering *TransitGatewayPeering) error {
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
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode delete_inter_transit_gateway_peering failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API delete_inter_transit_gateway_peering Get failed: " + data.Reason)
	}
	return nil
}
