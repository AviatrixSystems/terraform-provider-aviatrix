package goaviatrix

import (
	"encoding/json"
	"fmt"
	"io"

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
		// If Post can return resp alongside err, avoid leaks.
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
		return nil, fmt.Errorf("HTTP Post list_extended_vpc_peer failed: %w", err)
	}
	if resp == nil || resp.Body == nil {
		return nil, fmt.Errorf("HTTP Post list_extended_vpc_peer returned nil response/body")
	}
	defer resp.Body.Close()

	const maxBody = 256 << 10 // 256 KiB
	body, readErr := io.ReadAll(io.LimitReader(resp.Body, maxBody+1))
	if readErr != nil {
		return nil, fmt.Errorf("read list_extended_vpc_peer response body failed: %w", readErr)
	}

	if len(body) > maxBody {
		return nil, fmt.Errorf("list_extended_vpc_peer: response body too large (>%d bytes)", maxBody)
	}

	var data TransPeerListResp
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("json Decode list_extended_vpc_peer failed: %w\nBody: %s", err, string(body))
	}
	if !data.Return {
		return nil, fmt.Errorf("list_extended_vpc_peer: post failed: %s", data.Reason)
	}

	for i := range data.Results {
		if data.Results[i].Source == transPeer.Source && data.Results[i].Nexthop == transPeer.Nexthop {
			found := data.Results[i]
			return &found, nil
		}
	}

	log.Errorf(
		"Transitive peering with gateways %s and %s with subnet %s not found",
		transPeer.Source, transPeer.Nexthop, transPeer.ReachableCidr,
	)
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
