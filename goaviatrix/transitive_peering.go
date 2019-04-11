package goaviatrix

import (
	"encoding/json"
	"errors"
	"log"
	//"github.com/davecgh/go-spew/spew"
)

// TransPeer simple struct to hold transitive peering details

type TransPeer struct {
	CID           string `form:"CID,omitempty"`
	Action        string `form:"action,omitempty"`
	Source        string `form:"source" json:"source"`
	Nexthop       string `form:"nexthop" json:"nexthop"`
	ReachableCidr string `form:"reachable_cidr" json:"reachable_cidr"`
}

type TransPeerListResp struct {
	Return  bool        `json:"return"`
	Results []TransPeer `json:"results"`
	Reason  string      `json:"reason"`
}

func (c *Client) CreateTransPeer(transPeer *TransPeer) error {
	transPeer.CID = c.CID
	transPeer.Action = "add_extended_vpc_peer"
	resp, err := c.Post(c.baseURL, transPeer)
	if err != nil {
		return errors.New("HTTP Post add_extended_vpc_peer failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode add_extended_vpc_peer failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API add_extended_vpc_peer Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetTransPeer(transPeer *TransPeer) (*TransPeer, error) {
	transPeer.CID = c.CID
	transPeer.Action = "list_extended_vpc_peer"
	resp, err := c.Post(c.baseURL, transPeer)
	if err != nil {
		return nil, errors.New("HTTP Post list_extended_vpc_peer failed: " + err.Error())
	}
	var data TransPeerListResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_extended_vpc_peer failed: " + err.Error())
	}
	if !data.Return {
		return nil, errors.New("Rest API list_extended_vpc_peer Post failed: " + data.Reason)
	}
	transPeerList := data.Results
	for i := range transPeerList {
		if transPeerList[i].Source == transPeer.Source && transPeerList[i].Nexthop == transPeer.Nexthop {
			return &transPeerList[i], nil
		}
	}
	log.Printf("Transitive peering with gateways %s and %s with subnet %s not found",
		transPeer.Source, transPeer.Nexthop, transPeer.ReachableCidr)
	return nil, ErrNotFound
}

func (c *Client) UpdateTransPeer(transPeer *TransPeer) error {
	return nil
}

func (c *Client) DeleteTransPeer(transPeer *TransPeer) error {
	transPeer.CID = c.CID
	transPeer.Action = "delete_extended_vpc_peer"
	resp, err := c.Post(c.baseURL, transPeer)
	if err != nil {
		return errors.New("HTTP Post delete_extended_vpc_peer failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode delete_extended_vpc_peer failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API delete_extended_vpc_peer Post failed: " + data.Reason)
	}
	return nil
}
