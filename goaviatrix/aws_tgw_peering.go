package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

type AwsTgwPeering struct {
	Action   string `form:"action,omitempty"`
	CID      string `form:"CID,omitempty"`
	TgwName1 string `form:"tgw_name1,omitempty" json:"tgw_name1,omitempty"`
	TgwName2 string `form:"tgw_name2,omitempty" json:"tgw_name2,omitempty"`
}

type AwsTgwPeeringAPIResp struct {
	Return  bool     `json:"return"`
	Results []string `json:"results"`
	Reason  string   `json:"reason"`
}

func (c *Client) CreateAwsTgwPeering(awsTgwPeering *AwsTgwPeering) error {
	awsTgwPeering.CID = c.CID
	awsTgwPeering.Action = "add_tgw_peering"
	resp, err := c.Post(c.baseURL, awsTgwPeering)
	if err != nil {
		return errors.New("HTTP Post 'add_tgw_peering' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'add_tgw_peering' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'add_tgw_peering' Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetAwsTgwPeering(awsTgwPeering *AwsTgwPeering) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'list_peered_tgw_names': ") + err.Error())
	}
	listPeeredTgwNames := url.Values{}
	listPeeredTgwNames.Add("CID", c.CID)
	listPeeredTgwNames.Add("action", "list_peered_tgw_names")
	listPeeredTgwNames.Add("tgw_name", awsTgwPeering.TgwName1)
	Url.RawQuery = listPeeredTgwNames.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get 'list_peered_tgw_names' failed: " + err.Error())
	}
	var data AwsTgwPeeringAPIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'list_peered_tgw_names' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		if strings.Contains(data.Reason, "does not exist") {
			return ErrNotFound
		}
		return errors.New("Rest API 'list_peered_tgw_names' Get failed: " + data.Reason)
	}
	if len(data.Results) == 0 {
		log.Errorf("Aws tgw peering with tgw: %s and tgw: %s not found", awsTgwPeering.TgwName1, awsTgwPeering.TgwName2)
		return ErrNotFound
	}
	peeringList := data.Results
	for i := range peeringList {
		if peeringList[i] == awsTgwPeering.TgwName2 {
			return nil
		}
	}
	return ErrNotFound
}

func (c *Client) DeleteAwsTgwPeering(awsTgwPeering *AwsTgwPeering) error {
	awsTgwPeering.CID = c.CID
	awsTgwPeering.Action = "delete_tgw_peering"
	resp, err := c.Post(c.baseURL, awsTgwPeering)
	if err != nil {
		return errors.New("HTTP Post 'delete_tgw_peering' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'delete_tgw_peering' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'delete_tgw_peering' Post failed: " + data.Reason)
	}
	return nil
}
