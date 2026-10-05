package goaviatrix

// Tunnel simple struct to hold tunnel details

import (
	log "github.com/sirupsen/logrus"
)

type Tunnel struct {
	VpcName1        string `json:"vpc_name1"`
	VpcName2        string `json:"vpc_name2"`
	PeeringState    string `json:"peering_state"`
	PeeringHaStatus string `json:"peering_ha_status"`
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
	form := map[string]string{
		"CID":        c.CID,
		"action":     "peer_vpc_pair",
		"vpc_name1":  tunnel.VpcName1,
		"vpc_name2":  tunnel.VpcName2,
		"ha_enabled": tunnel.EnableHA,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) GetTunnel(tunnel *Tunnel) (*Tunnel, error) {
	form := map[string]string{
		"CID":    c.CID,
		"action": "list_peer_vpc_pairs",
	}

	var data TunnelListResp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	tunList := data.Results.PairList
	for i := range tunList {
		if tunList[i].VpcName1 == tunnel.VpcName1 && tunList[i].VpcName2 == tunnel.VpcName2 {
			log.Debugf("Found %s~%s tunnel: %#v", tunnel.VpcName1, tunnel.VpcName2, tunList[i])
			return &tunList[i], nil
		}
	}
	log.Errorf("Tunnel with gateways %s and %s not found", tunnel.VpcName1, tunnel.VpcName2)
	return nil, ErrNotFound
}

func (c *Client) UpdateTunnel(tunnel *Tunnel) error {
	return nil
}

func (c *Client) DeleteTunnel(tunnel *Tunnel) error {
	form := map[string]string{
		"CID":       c.CID,
		"action":    "unpeer_vpc_pair",
		"vpc_name1": tunnel.VpcName1,
		"vpc_name2": tunnel.VpcName2,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}
