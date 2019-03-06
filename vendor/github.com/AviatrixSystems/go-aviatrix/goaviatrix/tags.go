package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
)

// Tags simple struct to hold tag details
type Tags struct {
	Action       string `form:"action,omitempty"`
	CID          string `form:"CID,omitempty"`
	CloudType    int    `form:"cloud_type,omitempty"`
	ResourceType string `form:"resource_type,omitempty"`
	ResourceName string `form:"resource_name,omitempty"`
	TagList      string `form:"new_tag_list,omitempty"`
}

type TagAPIResp struct {
	Return  bool                         `json:"return"`
	Results map[string]map[string]string `json:"results"`
	Reason  string                       `json:"reason"`
}

func (c *Client) AddTags(tags *Tags) error {
	tags.CID = c.CID
	tags.Action = "add_resource_tags"
	resp, err := c.Post(c.baseURL, tags)
	if err != nil {
		return errors.New("HTTP Post add_resource_tags failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode add_resource_tags failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API add_resource_tags Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetTags(tags *Tags) ([]string, error) {
	tags.CID = c.CID
	tags.Action = "list_resource_tags"
	resp, err := c.Post(c.baseURL, tags)
	if err != nil {
		return nil, errors.New("HTTP Post list_resource_tags failed: " + err.Error())
	}

	var data TagAPIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_resource_tags failed: " + err.Error())
	}
	if !data.Return {
		return nil, errors.New("Rest API list_resource_tags Post failed: " + data.Reason)
	}

	var tagList []string
	keys := reflect.ValueOf(data.Results).MapKeys()
	strKeys := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		strKeys[i] = keys[i].String()
	}
	for i := 0; i < len(keys); i++ {
		if strKeys[i] == "tags" {
			allKeys := reflect.ValueOf(data.Results["tags"]).MapKeys()
			for i := 0; i < len(allKeys); i++ {
				if allKeys[i].String() == "Aviatrix-Created-Resource" &&
					data.Results["tags"][allKeys[i].String()] == "Do-Not-Delete-Aviatrix-Created-Resource" {
					continue
				}
				str := allKeys[i].String() + ":" + data.Results["tags"][allKeys[i].String()]
				tagList = append(tagList, str)
			}
			return tagList, nil
		}
	}

	return nil, nil
}

func (c *Client) DeleteTags(tags *Tags) error {
	tags.CID = c.CID
	tags.Action = "delete_resource_tags"
	verb := "POST"
	body := fmt.Sprintf("CID=%s&action=%s&cloud_type=%d&resource_type=%s&resource_name=%s&del_tag_list=%s",
		c.CID, tags.Action, tags.CloudType, tags.ResourceType, tags.ResourceName, tags.TagList)
	log.Printf("[TRACE] %s %s Body: %s", verb, c.baseURL, body)
	req, err := http.NewRequest(verb, c.baseURL, strings.NewReader(body))
	if err == nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		return errors.New("HTTP New Request Post delete_resource_tags failed: " + err.Error())
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return errors.New("HTTP Post delete_resource_tags failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode delete_resource_tags failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API delete_resource_tags Post failed: " + data.Reason)
	}
	return nil
}
