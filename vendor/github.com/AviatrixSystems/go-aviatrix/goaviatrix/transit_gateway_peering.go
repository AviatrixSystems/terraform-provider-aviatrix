package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=create_inter_transit_gateway_peering&gateway1=%s&gateway2=%s",
		c.CID, transitGatewayPeering.TransitGatewayName1, transitGatewayPeering.TransitGatewayName2)
	resp, err := c.Get(path, nil)
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

func (c *Client) GetTransitGatewayPeering(transitGatewayPeering *TransitGatewayPeering) error {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=list_inter_transit_gateway_peering", c.CID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return err
	}
	var data TransitGatewayPeeringAPIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		return errors.New(data.Reason)
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
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=delete_inter_transit_gateway_peering&gateway1=%s&gateway2=%s",
		c.CID, transitGatewayPeering.TransitGatewayName1, transitGatewayPeering.TransitGatewayName2)
	resp, err := c.Delete(path, nil)
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
