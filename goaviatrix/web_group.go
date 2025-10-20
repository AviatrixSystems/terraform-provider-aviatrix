package goaviatrix

import (
	"context"
	"encoding/json"
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
	form := makeWebGroupForm(webGroup)
	return c.appdomainCache.Create(ctx, c, form)
}

func (c *Client) GetWebGroup(ctx context.Context, uuid string) (*WebGroup, error) {
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

	g, err := c.appdomainCache.Get(ctx, c, uuid)
	if err != nil {
		return nil, err
	}

	selector := WebGroupAnyResult{}
	if err := json.Unmarshal(g.Selector, &selector); err != nil {
		return nil, err
	}

	webGroup := &WebGroup{
		Name: g.Name,
		UUID: g.UUID,
	}

	for _, filterResult := range selector.Any {
		filterMap := filterResult.All

		filter := &WebGroupMatchExpression{
			SniFilter: filterMap["snifilter"],
			UrlFilter: filterMap["urlfilter"],
		}

		webGroup.Selector.Expressions = append(webGroup.Selector.Expressions, filter)
	}

	return webGroup, nil
}

func (c *Client) UpdateWebGroup(ctx context.Context, webGroup *WebGroup, uuid string) error {
	form := makeWebGroupForm(webGroup)
	return c.appdomainCache.Update(ctx, c, uuid, form)
}

func (c *Client) DeleteWebGroup(ctx context.Context, uuid string) error {
	return c.appdomainCache.Delete(ctx, c, uuid)
}
