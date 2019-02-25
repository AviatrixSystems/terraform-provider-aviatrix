package goaviatrix

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
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

func (c *Client) CreateAWSPeer(awsPeer *AWSPeer) (string, error) {
	awsPeer.CID = c.CID
	awsPeer.Action = "create_aws_peering"
	resp, err := c.Post(c.baseURL, awsPeer)
	if err != nil {
		return "", errors.New("HTTP Post create_aws_peering failed: " + err.Error())
	}
	var data AwsPeerAPIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", errors.New("Json Decode create_aws_peering failed: " + err.Error())
	}
	if !data.Return {
		return "", errors.New("Rest API create_aws_peering Post failed: " + data.Reason)
	}
	r, _ := regexp.Compile(`pcx-\w+`)
	id := r.FindString(data.Results["text"])
	return id, nil
}

func (c *Client) GetAWSPeer(awsPeer *AWSPeer) (*AWSPeer, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for list_aws_peerings ") + err.Error())
	}
	listAwsPeering := url.Values{}
	listAwsPeering.Add("CID", c.CID)
	listAwsPeering.Add("action", "list_aws_peerings")
	Url.RawQuery = listAwsPeering.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return nil, errors.New("HTTP Get list_aws_peerings failed: " + err.Error())
	}

	//Output result for this query cannot be unmarshalled
	//easily into our defined struct AWSPeer.
	//So using a map of string->interface{}
	var data map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_aws_peerings failed: " + err.Error())
	}
	if _, ok := data["reason"]; ok {
		log.Printf("[INFO] Couldn't find AWS peering between VPCs %s and %s: %s", awsPeer.VpcID1, awsPeer.VpcID2, data["reason"])
		return nil, ErrNotFound
	}
	if val, ok := data["results"]; ok {
		if pairList, ok1 := val.(map[string]interface{})["pair_list"].([]interface{}); ok1 {
			for i := range pairList {
				if pairList[i].(map[string]interface{})["requester"].(map[string]interface{})["vpc_id"].(string) == awsPeer.VpcID1 &&
					pairList[i].(map[string]interface{})["accepter"].(map[string]interface{})["vpc_id"].(string) == awsPeer.VpcID2 {
					awsPeer := &AWSPeer{
						VpcID1:       pairList[i].(map[string]interface{})["requester"].(map[string]interface{})["vpc_id"].(string),
						VpcID2:       pairList[i].(map[string]interface{})["accepter"].(map[string]interface{})["vpc_id"].(string),
						AccountName1: pairList[i].(map[string]interface{})["requester"].(map[string]interface{})["account_name"].(string),
						AccountName2: pairList[i].(map[string]interface{})["accepter"].(map[string]interface{})["account_name"].(string),
						Region1:      pairList[i].(map[string]interface{})["requester"].(map[string]interface{})["region"].(string),
						Region2:      pairList[i].(map[string]interface{})["accepter"].(map[string]interface{})["region"].(string),
					}
					return awsPeer, nil
				}
			}
		}
	}
	log.Printf("[INFO] No AWS peering between VPC %s and %s is present.", awsPeer.VpcID1, awsPeer.VpcID2)
	return nil, ErrNotFound
}

func (c *Client) UpdateAWSPeer(awsPeer *AWSPeer) error {
	return nil
}

func (c *Client) DeleteAWSPeer(awsPeer *AWSPeer) error {
	awsPeer.CID = c.CID
	awsPeer.Action = "delete_aws_peering"
	resp, err := c.Post(c.baseURL, awsPeer)
	if err != nil {
		return errors.New("HTTP Post delete_aws_peering failed: " + err.Error())
	}

	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode delete_aws_peering failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API delete_aws_peering Post failed: " + data.Reason)
	}
	return nil
}
