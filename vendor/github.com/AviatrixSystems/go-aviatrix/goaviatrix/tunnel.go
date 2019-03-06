package goaviatrix

// Tunnel simple struct to hold tunnel details

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
)

type Tunnel struct {
	VpcName1        string `json:"vpc_name1"`
	VpcName2        string `json:"vpc_name2"`
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
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for peer_vpc_pair ") + err.Error())
	}
	peerVpcPair := url.Values{}
	peerVpcPair.Add("CID", c.CID)
	peerVpcPair.Add("action", "peer_vpc_pair")
	peerVpcPair.Add("vpc_name1", tunnel.VpcName1)
	peerVpcPair.Add("vpc_name2", tunnel.VpcName2)
	peerVpcPair.Add("ha_enabled", tunnel.EnableHA)
	Url.RawQuery = peerVpcPair.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get peer_vpc_pair failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode peer_vpc_pair failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API peer_vpc_pair Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetTunnel(tunnel *Tunnel) (*Tunnel, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for list_peer_vpc_pairs ") + err.Error())
	}
	listPeerVpcPairs := url.Values{}
	listPeerVpcPairs.Add("CID", c.CID)
	listPeerVpcPairs.Add("action", "list_peer_vpc_pairs")
	Url.RawQuery = listPeerVpcPairs.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return nil, errors.New("HTTP Get list_peer_vpc_pairs failed: " + err.Error())
	}
	var data TunnelListResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_peer_vpc_pairs failed: " + err.Error())
	}
	if !data.Return {
		return nil, errors.New("Rest API list_peer_vpc_pairs Get failed: " + data.Reason)
	}
	tunList := data.Results.PairList
	for i := range tunList {
		if tunList[i].VpcName1 == tunnel.VpcName1 && tunList[i].VpcName2 == tunnel.VpcName2 {
			log.Printf("[DEBUG] Found %s~%s tunnel: %#v", tunnel.VpcName1, tunnel.VpcName2, tunList[i])
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
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for unpeer_vpc_pair ") + err.Error())
	}
	unPeerVpcPair := url.Values{}
	unPeerVpcPair.Add("CID", c.CID)
	unPeerVpcPair.Add("action", "unpeer_vpc_pair")
	unPeerVpcPair.Add("vpc_name1", tunnel.VpcName1)
	unPeerVpcPair.Add("vpc_name2", tunnel.VpcName2)
	Url.RawQuery = unPeerVpcPair.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get unpeer_vpc_pair failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode unpeer_vpc_pair failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API unpeer_vpc_pair Get failed: " + data.Reason)
	}
	return nil
}
