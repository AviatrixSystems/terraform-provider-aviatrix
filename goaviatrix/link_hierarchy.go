package goaviatrix

import (
	"context"
	"fmt"
	"reflect"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type WanLink struct {
	WanTag string `json:"wan_tag"`
}

type Link struct {
	Name        string    `json:"name"`
	WanLinkList []WanLink `json:"wan_link"`
}

type LinkHierarchy struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Links []Link `json:"links"`
}

type LinkHierarchyResp struct {
	LinkHierarchy []LinkHierarchy `json:"link_hierarchy"`
}

func (c *Client) CreateLinkHierarchy(ctx context.Context, linkHierarchy map[string]interface{}) (string, error) {
	endpoint := "ipsla/link-hierarchy"

	type resp struct {
		UUID string `json:"uuid"`
	}

	var data resp
	err := c.PostAPIContext25(ctx, &data, endpoint, linkHierarchy)
	if err != nil {
		return "", err
	}

	return data.UUID, nil
}

func (c *Client) GetLinkHierarchy(ctx context.Context, uuid string) (*LinkHierarchy, error) {
	endpoint := fmt.Sprintf("ipsla/link-hierarchy/%s", uuid)

	var data LinkHierarchyResp
	err := c.GetAPIContext25(ctx, &data, endpoint, nil)
	if err != nil {
		return nil, err
	}

	for _, linkHierarchyResult := range data.LinkHierarchy {
		if linkHierarchyResult.UUID == uuid {
			return &linkHierarchyResult, nil
		}
	}

	return nil, ErrNotFound
}

func (c *Client) UpdateLinkHierarchy(ctx context.Context, linkHierarchy map[string]interface{}, uuid string) error {
	endpoint := fmt.Sprintf("ipsla/link-hierarchy/%s", uuid)
	return c.PutAPIContext25(ctx, endpoint, linkHierarchy)
}

func (c *Client) DeleteLinkHierarchy(ctx context.Context, uuid string) error {
	endpoint := fmt.Sprintf("ipsla/link-hierarchy/%s", uuid)
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}

func DiffSuppressFuncLinkHierarchy(k, old, new string, d *schema.ResourceData) bool {
	lOld, lNew := d.GetChange("links")

	oldList, ok := lOld.([]interface{})
	if !ok {
		return false
	}
	newList, ok := lNew.([]interface{})
	if !ok {
		return false
	}

	linksOld := make([]map[string]interface{}, 0, len(oldList))
	for _, l0 := range oldList {
		l1, ok := l0.(map[string]interface{})
		if !ok {
			return false
		}
		linksOld = append(linksOld, l1)
	}

	linksNew := make([]map[string]interface{}, 0, len(newList))
	for _, l0 := range newList {
		l1, ok := l0.(map[string]interface{})
		if !ok {
			return false
		}
		linksNew = append(linksNew, l1)
	}

	sort.Slice(linksOld, func(i, j int) bool {
		ni, ok := linksOld[i]["name"].(string)
		if !ok {
			return false
		}
		nj, ok := linksOld[j]["name"].(string)
		if !ok {
			return false
		}
		return ni < nj
	})

	sort.Slice(linksNew, func(i, j int) bool {
		ni, ok := linksNew[i]["name"].(string)
		if !ok {
			return false
		}
		nj, ok := linksNew[j]["name"].(string)
		if !ok {
			return false
		}
		return ni < nj
	})

	return reflect.DeepEqual(linksOld, linksNew)
}
