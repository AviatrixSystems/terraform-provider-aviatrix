package goaviatrix

import (
	"encoding/json"
	"errors"
	"log"
	//"github.com/davecgh/go-spew/spew"
)

// Transpeer simple struct to hold transitive peering details

type Transpeer struct {
	CID           	string `form:"CID,omitempty"`
	Action          string `form:"action,omitempty"`
	Source          string `form:"source" json:"source"`
	Nexthop         string `form:"nexthop" json:"nexthop"`
	ReachableCidr   string `form:"reachable_cidr" json:"reachable_cidr"`
}

type TranspeerListResp struct {
	Return  bool   `json:"return"`
	Results []Transpeer `json:"results"`
	Reason  string `json:"reason"`
}

func (c *Client) CreateTranspeer(transpeer *Transpeer) (error) {
	transpeer.CID=c.CID
	transpeer.Action="add_extended_vpc_peer"
	resp,err := c.Post(c.baseURL, transpeer)
		if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if(!data.Return){
		return errors.New(data.Reason)
	}
	return nil
}

func (c *Client) GetTranspeer(transpeer *Transpeer) (*Transpeer, error) {
	transpeer.CID=c.CID
	transpeer.Action="list_extended_vpc_peer"
	resp,err := c.Post(c.baseURL, transpeer)
		if err != nil {
		return nil, err
	}
	var data TranspeerListResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if(!data.Return){
		return nil, errors.New(data.Reason)
	}
	transpeerList:= data.Results
	for i := range transpeerList {
		if transpeerList[i].Source == transpeer.Source && transpeerList[i].Nexthop == transpeer.Nexthop {
			return &transpeerList[i], nil
		}
	}
	log.Printf("Transitive peering with gateways %s and %s with subnet %s not found", transpeer.Source, transpeer.Nexthop, transpeer.ReachableCidr)
	return nil, ErrNotFound
}

func (c *Client) UpdateTranspeer(transpeer *Transpeer) (error) {
	return nil
}

func (c *Client) DeleteTranspeer(transpeer *Transpeer) (error) {
	transpeer.CID=c.CID
	transpeer.Action="delete_extended_vpc_peer"
	resp,err := c.Post(c.baseURL, transpeer)
		if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if(!data.Return){
		return errors.New(data.Reason)
	}
	return nil
}
