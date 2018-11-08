package goaviatrix

// Tunnel simple struct to hold tunnel details

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

type Tunnel struct {
	VpcName1        string `json:"vpc_name1"`
	VpcName2        string `json:"vpc_name2"`
	OverAwsPeering  string `json:"over_aws_peering"`
	PeeringState    string `json:"peering_state"`
	PeeringHaStatus string `json:"peering_ha_status"`
	Cluster         string `json:"cluster"`
	PeeringLink     string `json:"peering_link"`
	EnableHA        string `json:"enable_ha"`
}

type TunnelResult struct {
	PairList []Tunnel `json:"pair_list"`
}

type TunnelListResp struct {
	Return  bool         `json:"return"`
	Results TunnelResult `json:"results"`
	Reason  string       `json:"reason"`
}

func (c *Client) CreateTunnel(tunnel *Tunnel) error {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=peer_vpc_pair&vpc_name1=%s&vpc_name2=%s&ha_enabled=%s", c.CID, tunnel.VpcName1, tunnel.VpcName2, tunnel.EnableHA)
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

func (c *Client) GetTunnel(tunnel *Tunnel) (*Tunnel, error) {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=list_peer_vpc_pairs", c.CID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	var data TunnelListResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if !data.Return {
		return nil, errors.New(data.Reason)
	}
	tunList := data.Results.PairList
	for i := range tunList {
		if tunList[i].VpcName1 == tunnel.VpcName1 && tunList[i].VpcName2 == tunnel.VpcName2 {
			log.Printf("[DEBUG] Found %s<->%s tunnel: %#v", tunnel.VpcName1, tunnel.VpcName2, tunList[i])
			return &tunList[i], nil
		}
	}
	log.Printf("Tunnel with gateways %s and %s not found", tunnel.VpcName1, tunnel.VpcName2)
	return nil, ErrNotFound
}

func (c *Client) UpdateTunnel(tunnel *Tunnel) error {
	return nil
}

func (c *Client) DeleteTunnel(tunnel *Tunnel) error {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=unpeer_vpc_pair&vpc_name1=%s&vpc_name2=%s", c.CID, tunnel.VpcName1, tunnel.VpcName2)
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
