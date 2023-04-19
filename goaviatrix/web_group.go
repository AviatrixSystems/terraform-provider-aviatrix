package goaviatrix

import (
	"context"
	"fmt"
	"log"
)

type WebGroupMatchExpression struct {
	SniFilter string `json:"snifilter,omitempty"`
	UrlFilter string `json:"urlfilter,omitempty"`
}

type WebGroupSelector struct {
	Expressions []*WebGroupMatchExpression
}

type WebGroup struct {
	Name     string
	UUID     string
	Selector WebGroupSelector
}

func webGroupFilterToMap(filter *WebGroupMatchExpression) map[string]string {
	filterMap := make(map[string]string)

	if len(filter.SniFilter) > 0 {
		filterMap["snifilter"] = filter.SniFilter
	}
	if len(filter.UrlFilter) > 0 {
		filterMap["urlfilter"] = filter.UrlFilter
	}

	log.Printf("[DEBUG] POST filter map: %v\n", filterMap)
	return filterMap
}

func makeWebGroupForm(webGroup *WebGroup) map[string]interface{} {
	form := map[string]interface{}{
		"name": webGroup.Name,
	}

	var or []map[string]map[string]string
	for _, webGroupSelector := range webGroup.Selector.Expressions {
		and := map[string]map[string]string{
			"all": webGroupFilterToMap(webGroupSelector),
		}

		or = append(or, and)
	}

	form["selector"] = map[string]interface{}{
		"any": or,
	}

	return form
}

func (c *Client) CreateWebGroup(ctx context.Context, webGroup *WebGroup) (string, error) {
	endpoint := "app-domains"
	form := makeWebGroupForm(webGroup)

	type WebGroupResp struct {
		UUID string `json:"uuid"`
	}

	var data WebGroupResp
	err := c.PostAPIContext25(ctx, &data, endpoint, form)
	if err != nil {
		return "", err
	}
	return data.UUID, nil
}

func (c *Client) GetWebGroup(ctx context.Context, uuid string) (*WebGroup, error) {
	endpoint := "app-domains"

	type WebGroupMatchExpressionResult struct {
		All map[string]string `json:"all"`
	}

	type WebGroupAnyResult struct {
		Any []WebGroupMatchExpressionResult `json:"any"`
	}

	type WebGroupResult struct {
		UUID     string            `json:"uuid"`
		Name     string            `json:"name"`
		Selector WebGroupAnyResult `json:"selector"`
	}

	type WebGroupResp struct {
		WebGroups []WebGroupResult `json:"app_domains"`
	}

	var data WebGroupResp
	err := c.GetAPIContext25(ctx, &data, endpoint, nil)
	if err != nil {
		return nil, err
	}

	for _, webGroupResult := range data.WebGroups {
		if webGroupResult.UUID == uuid {
			webGroup := &WebGroup{
				Name: webGroupResult.Name,
				UUID: webGroupResult.UUID,
			}

			for _, filterResult := range webGroupResult.Selector.Any {
				filterMap := filterResult.All

				filter := &WebGroupMatchExpression{
					SniFilter: filterMap["snifilter"],
					UrlFilter: filterMap["urlfilter"],
				}

				webGroup.Selector.Expressions = append(webGroup.Selector.Expressions, filter)
			}
			return webGroup, nil
		}
	}
	return nil, ErrNotFound
}

func (c *Client) UpdateWebGroup(ctx context.Context, webGroup *WebGroup, uuid string) error {
	endpoint := fmt.Sprintf("app-domains/%s", uuid)
	form := makeWebGroupForm(webGroup)
	return c.PutAPIContext25(ctx, endpoint, form)
}

func (c *Client) DeleteWebGroup(ctx context.Context, uuid string) error {
	endpoint := fmt.Sprintf("app-domains/%s", uuid)
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}
