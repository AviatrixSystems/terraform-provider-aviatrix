package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Tags simple struct to hold tag details
type Tags struct {
	Action       string `form:"action,omitempty"`
	CID          string `form:"CID,omitempty"`
	CloudType    int    `form:"cloud_type,omitempty"`
	ResourceType string `form:"resource_type,omitempty"`
	ResourceName string `form:"resource_name,omitempty"`
	TagList      string `form:"new_tag_list,omitempty"`
	Tags         map[string]string
}

type TagAPIResp struct {
	Return  bool                         `json:"return"`
	Results map[string]map[string]string `json:"results"`
	Reason  string                       `json:"reason"`
}

func (tags *Tags) ComputeTagList() {
	tagList := make([]string, 0, len(tags.Tags))
	for key, val := range tags.Tags {
		tagList = append(tagList, key+":"+val)
	}
	tagListStr := strings.Join(tagList, ",")
	tags.TagList = tagListStr
}

func (c *Client) AddTags(tags *Tags) error {
	tags.CID = c.CID
	tags.Action = "add_resource_tags"
	resp, err := c.Post(c.baseURL, tags)
	if err != nil {
		return errors.New("HTTP Post add_resource_tags failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode add_resource_tags failed: " + err.Error() + "\n Body: " + bodyString)
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
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_resource_tags failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return nil, errors.New("Rest API list_resource_tags Post failed: " + data.Reason)
	}
	if tagsMap, ok := data.Results["usr_tags"]; ok {
		tags.Tags = tagsMap
	}
	var tagList []string
	keys := reflect.ValueOf(data.Results).MapKeys()
	strKeys := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		strKeys[i] = keys[i].String()
	}
	for i := 0; i < len(keys); i++ {
		if strKeys[i] == "usr_tags" {
			allKeys := reflect.ValueOf(data.Results["usr_tags"]).MapKeys()
			for i := 0; i < len(allKeys); i++ {
				if (allKeys[i].String() == "Key" && data.Results["usr_tags"][allKeys[i].String()] == "Aviatrix-Created-Resource") ||
					(allKeys[i].String() == "Value" && data.Results["usr_tags"][allKeys[i].String()] == "Do-Not-Delete-Aviatrix-Created-Resource") {
					continue
				}
				str := allKeys[i].String() + ":" + data.Results["usr_tags"][allKeys[i].String()]
				tagList = append(tagList, str)
			}
			return tagList, nil
		}
	}

	return nil, nil
}

func (c *Client) GetTagsMap(tags *Tags) error {
	tags.CID = c.CID
	tags.Action = "list_resource_tags"

	resp, err := c.Post(c.baseURL, tags)
	if err != nil {
		return fmt.Errorf("HTTP Post list_resource_tags failed: %v", err)
	}
	var data TagAPIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return fmt.Errorf("Json Decode list_resource_tags failed: %v\nBody: %s", err, bodyString)
	}
	if !data.Return {
		return fmt.Errorf("rest API list_resource_tags Post failed: %s", data.Reason)
	}

	if tagsMap, ok := data.Results["usr_tags"]; ok {
		tags.Tags = tagsMap
	}
	return nil
}

func (c *Client) DeleteTags(tags *Tags) error {
	tags.CID = c.CID
	tags.Action = "delete_resource_tag"
	verb := "POST"
	body := fmt.Sprintf("CID=%s&action=%s&cloud_type=%d&resource_type=%s&resource_name=%s&del_tag_list=%s",
		c.CID, tags.Action, tags.CloudType, tags.ResourceType, tags.ResourceName, tags.TagList)
	log.Tracef("%s %s Body: %s", verb, c.baseURL, body)
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
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode delete_resource_tags failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API delete_resource_tags Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) UpdateTags(tags *Tags) error {
	tags.CID = c.CID
	tags.Action = "update_resource_tags"

	return c.PostAPI(tags.Action, tags, BasicCheck)
}
