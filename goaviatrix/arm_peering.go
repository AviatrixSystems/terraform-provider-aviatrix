package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"strings"
)

// ARMPeer simple struct to hold arm_peer details
type ARMPeer struct {
	Action       string `form:"action,omitempty"`
	CID          string `form:"CID,omitempty"`
	AccountName1 string `form:"req_account_name,omitempty"`
	AccountName2 string `form:"acc_account_name,omitempty"`
	VNet1        string `form:"req_vpc_id,omitempty"`
	VNet2        string `form:"acc_vpc_id,omitempty"`
	Region1      string `form:"req_region,omitempty"`
	Region2      string `form:"acc_region,omitempty"`
	VNetCidr1    []string
	VNetCidr2    []string
}

type ArmPeerAPIResp struct {
	Return  bool   `json:"return"`
	Reason  string `json:"reason"`
	Results string `json:"results"`
}

func (c *Client) CreateARMPeer(armPeer *ARMPeer) error {
	armPeer.CID = c.CID
	armPeer.Action = "arm_peer_vnet_pair"
	resp, err := c.Post(c.baseURL, armPeer)
	if err != nil {
		return errors.New("HTTP Post arm_peer_vnet_pair failed: " + err.Error())
	}
	var data ArmPeerAPIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode arm_peer_vnet_pair failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API arm_peer_vnet_pair Post failed: " + data.Reason)
	}

	return nil
}

func (c *Client) GetARMPeer(armPeer *ARMPeer) (*ARMPeer, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for list_arm_peer_vnet_pairs ") + err.Error())
	}
	listArmPeering := url.Values{}
	listArmPeering.Add("CID", c.CID)
	listArmPeering.Add("action", "list_arm_peer_vnet_pairs")
	Url.RawQuery = listArmPeering.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return nil, errors.New("HTTP Get list_arm_peer_vnet_pairs failed: " + err.Error())
	}
	var data map[string]interface{}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_arm_peer_vnet_pairs failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if _, ok := data["reason"]; ok {
		log.Printf("[INFO] Couldn't find ARM peering between VPCs %s and %s: %s", armPeer.VNet1, armPeer.VNet2, data["reason"])
		return nil, ErrNotFound
	}
	if val, ok := data["results"]; ok {
		pairList := val.(interface{}).([]interface{})
		for i := range pairList {
			if pairList[i].(map[string]interface{})["requester"].(map[string]interface{})["vpc_id"].(string) == armPeer.VNet1 &&
				pairList[i].(map[string]interface{})["accepter"].(map[string]interface{})["vpc_id"].(string) == armPeer.VNet2 {
				armPeer := &ARMPeer{
					VNet1:        pairList[i].(map[string]interface{})["requester"].(map[string]interface{})["vpc_id"].(string),
					VNet2:        pairList[i].(map[string]interface{})["accepter"].(map[string]interface{})["vpc_id"].(string),
					AccountName1: pairList[i].(map[string]interface{})["requester"].(map[string]interface{})["account_name"].(string),
					AccountName2: pairList[i].(map[string]interface{})["accepter"].(map[string]interface{})["account_name"].(string),
					Region1:      pairList[i].(map[string]interface{})["requester"].(map[string]interface{})["region"].(string),
					Region2:      pairList[i].(map[string]interface{})["accepter"].(map[string]interface{})["region"].(string),
				}

				vnetCidrList1 := pairList[i].(map[string]interface{})["requester"].(map[string]interface{})["vpc_cidr"].([]interface{})
				var vnetCidr1 []string
				for i := range vnetCidrList1 {
					vnetCidr1 = append(vnetCidr1, vnetCidrList1[i].(interface{}).(string))
				}
				armPeer.VNetCidr1 = vnetCidr1

				vnetCidrList2 := pairList[i].(map[string]interface{})["accepter"].(map[string]interface{})["vpc_cidr"].([]interface{})
				var vnetCidr2 []string
				for i := range vnetCidrList2 {
					vnetCidr2 = append(vnetCidr2, vnetCidrList2[i].(interface{}).(string))
				}
				armPeer.VNetCidr2 = vnetCidr2

				return armPeer, nil
			}
		}
	}
	return nil, ErrNotFound
}

func (c *Client) UpdateARMPeer(armPeer *ARMPeer) error {
	return nil
}

func (c *Client) DeleteARMPeer(armPeer *ARMPeer) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for arm_unpeer_vnet_pair") + err.Error())
	}
	armUnpeerVNetPair := url.Values{}
	armUnpeerVNetPair.Add("CID", c.CID)
	armUnpeerVNetPair.Add("action", "arm_unpeer_vnet_pair")
	armUnpeerVNetPair.Add("vpc_name1", armPeer.VNet1)
	armUnpeerVNetPair.Add("vpc_name2", armPeer.VNet2)
	Url.RawQuery = armUnpeerVNetPair.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Post arm_unpeer_vnet_pair failed: " + err.Error())
	}

	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode arm_unpeer_vnet_pair failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API arm_unpeer_vnet_pair Post failed: " + data.Reason)
	}
	return nil
}
