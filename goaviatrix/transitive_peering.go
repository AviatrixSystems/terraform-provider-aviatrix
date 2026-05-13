package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"

	//"github.com/davecgh/go-spew/spew"

	log "github.com/sirupsen/logrus"
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

	return c.PostAPI(transPeer.Action, transPeer, BasicCheck)
}

func (c *Client) GetTransPeer(transPeer *TransPeer) (*TransPeer, error) {
	// TODO: use GetAPI - need API details
	transPeer.CID = c.CID
	transPeer.Action = "list_extended_vpc_peer"
	resp, err := c.Post(c.baseURL, transPeer)
	if err != nil {
		return nil, errors.New("HTTP Post list_extended_vpc_peer failed: " + err.Error())
	}
	var data TransPeerListResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_extended_vpc_peer failed: " + err.Error() + "\n Body: " + bodyString)
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
	log.Errorf("Transitive peering with gateways %s and %s with subnet %s not found",
		transPeer.Source, transPeer.Nexthop, transPeer.ReachableCidr)
	return nil, ErrNotFound
}

func (c *Client) UpdateTransPeer(transPeer *TransPeer) error {
	return nil
}

func (c *Client) DeleteTransPeer(transPeer *TransPeer) error {
	transPeer.CID = c.CID
	transPeer.Action = "delete_extended_vpc_peer"

	return c.PostAPI(transPeer.Action, transPeer, BasicCheck)
}
