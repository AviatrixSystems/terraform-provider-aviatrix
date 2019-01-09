package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
)

// AWSPeer simple struct to hold aws_peer details
type AWSPeer struct {
	Action       string `form:"action,omitempty"`
	CID          string `form:"CID,omitempty"`
	AccountName1 string `form:"peer1_account_name,omitempty"`
	AccountName2 string `form:"peer2_account_name,omitempty"`
	VpcID1       string `form:"peer1_vpc_id,omitempty"`
	VpcID2       string `form:"peer2_vpc_id,omitempty"`
	Region1      string `form:"peer1_region,omitempty"`
	Region2      string `form:"peer2_region,omitempty"`
	RtbList1     string `form:"peer1_rtb_id,omitempty"`
	RtbList2     string `form:"peer2_rtb_id,omitempty"`
}

type AwsPeerAPIResp struct {
	Return  bool              `json:"return"`
	Reason  string            `json:"reason"`
	Results map[string]string `json:"results"`
}

func (c *Client) CreateAWSPeer(aws_peer *AWSPeer) (string, error) {
	aws_peer.CID = c.CID
	aws_peer.Action = "create_aws_peering"
	resp, err := c.Post(c.baseURL, aws_peer)
	if err != nil {
		return "", err
	}
	var data AwsPeerAPIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	if !data.Return {
		return "", errors.New(data.Reason)
	}
	r, _ := regexp.Compile(`pcx-\w+`)
	id := r.FindString(data.Results["text"])
	return id, nil
}

func (c *Client) GetAWSPeer(aws_peer *AWSPeer) (*AWSPeer, error) {
	aws_peer.CID = c.CID
	aws_peer.Action = "list_aws_peerings"
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=%s", aws_peer.CID, aws_peer.Action)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	//Output result for this query cannot be unmarshalled
	//easily into our defined struct AWSPeer.
	//So using a map of string->interface{}
	var data map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if _, ok := data["reason"]; ok {
		log.Printf("[INFO] Couldn't find AWS peering between VPCs %s and %s: %s", aws_peer.VpcID1, aws_peer.VpcID2, data["reason"])
		return nil, ErrNotFound
	}
	if val, ok := data["results"]; ok {
		if pair_list, ok1 := val.(map[string]interface{})["pair_list"].([]interface{}); ok1 {
			for i := range pair_list {
				if pair_list[i].(map[string]interface{})["requester"].(map[string]interface{})["vpc_id"].(string) == aws_peer.VpcID1 && pair_list[i].(map[string]interface{})["accepter"].(map[string]interface{})["vpc_id"].(string) == aws_peer.VpcID2 {
					aws_peer := &AWSPeer{
						VpcID1:       pair_list[i].(map[string]interface{})["requester"].(map[string]interface{})["vpc_id"].(string),
						VpcID2:       pair_list[i].(map[string]interface{})["accepter"].(map[string]interface{})["vpc_id"].(string),
						AccountName1: pair_list[i].(map[string]interface{})["requester"].(map[string]interface{})["account_name"].(string),
						AccountName2: pair_list[i].(map[string]interface{})["accepter"].(map[string]interface{})["account_name"].(string),
						Region1:      pair_list[i].(map[string]interface{})["requester"].(map[string]interface{})["region"].(string),
						Region2:      pair_list[i].(map[string]interface{})["accepter"].(map[string]interface{})["region"].(string),
					}
					return aws_peer, nil
				}
			}
		}
	}
	log.Printf("[INFO] No AWS peering between VPC %s and %s is present.", aws_peer.VpcID1, aws_peer.VpcID2)
	return nil, ErrNotFound
}

func (c *Client) UpdateAWSPeer(aws_peer *AWSPeer) error {
	return nil
}

func (c *Client) DeleteAWSPeer(aws_peer *AWSPeer) error {
	aws_peer.CID = c.CID
	aws_peer.Action = "delete_aws_peering"
	resp, err := c.Post(c.baseURL, aws_peer)
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
