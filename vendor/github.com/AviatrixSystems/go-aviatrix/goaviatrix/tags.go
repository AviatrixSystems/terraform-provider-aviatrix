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

func (c *Client) GetTags(tags *Tags) ([]string, error) {
	tags.CID = c.CID
	tags.Action = "list_resource_tags"
	resp, err := c.Post(c.baseURL, tags)
	if err != nil {
		return nil, err
	}

	var data TagAPIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if !data.Return {
		return nil, errors.New(data.Reason)
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
		return err
	}
	resp, err := c.HTTPClient.Do(req)
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
