package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
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
	Tags         map[string]string
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
	data := map[string]string{
		"action":        "list_resource_tags",
		"CID":           c.CID,
		"cloud_type":    strconv.Itoa(tags.CloudType),
		"resource_type": tags.ResourceType,
		"resource_name": tags.ResourceName,
	}
	var resp TagAPIResp
	err := c.GetAPI(&resp, data["Action"], data, BasicCheck)
	if err != nil {
		return nil, err
	}

	var tagList []string
	if tagsMap, ok := resp.Results["usr_tags"]; ok {
		tags.Tags = tagsMap
		for key, val := range tagsMap {
			tagStr := key + ":" + val
			tagList = append(tagList, tagStr)
		}
	}

	return tagList, nil
}

func (c *Client) GetTagsMap(tags *Tags) (map[string]string, error) {
	data := map[string]string{
		"action":        "list_resource_tags",
		"CID":           c.CID,
		"cloud_type":    strconv.Itoa(tags.CloudType),
		"resource_type": tags.ResourceType,
		"resource_name": tags.ResourceName,
	}

	var resp TagAPIResp
	err := c.GetAPI(&resp, data["Action"], data, BasicCheck)
	if err != nil {
		return nil, err
	}
	if tagsMap, ok := resp.Results["usr_tags"]; ok {
		return tagsMap, nil
	}
	return nil, nil
}

func (c *Client) DeleteTags(tags *Tags) error {
	params := map[string]string{
		"action":        "delete_resource_tag",
		"CID":           c.CID,
		"cloud_type":    strconv.Itoa(tags.CloudType),
		"del_tag_list":  tags.TagList,
		"resource_name": tags.ResourceName,
		"resource_type": tags.ResourceType,
	}

	return c.PostAPI(params["action"], params, BasicCheck)
}

func (c *Client) UpdateTags(tags *Tags) error {
	tags.CID = c.CID
	tags.Action = "update_resource_tags"

	return c.PostAPI(tags.Action, tags, BasicCheck)
}
