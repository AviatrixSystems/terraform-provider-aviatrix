package goaviatrix

import (
	"context"
	"fmt"
	"log"
	"strings"
)

const (
	SniFilterKey   = "snifilter"
	UrlFilterKey   = "urlfilter"
	ExternalURLKey = "external_url"
	ExternalSNIKey = "external_sni"
)

type WebGroupMatchExpression struct {
	SniFilter   string `json:"snifilter,omitempty"`
	UrlFilter   string `json:"urlfilter,omitempty"`
	ExternalURL string `json:"external_url,omitempty"`
	ExternalSNI string `json:"external_sni,omitempty"`
	ExtArgs     map[string]string
}

type WebGroupSelector struct {
	Expressions []*WebGroupMatchExpression
}

type WebGroup struct {
	Name     string
	UUID     string
	Selector WebGroupSelector
}

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

func webGroupFilterToMap(filter *WebGroupMatchExpression) map[string]string {
	filterMap := make(map[string]string)

	// External feed arguments sit alongside the feed key in the same "all" map.
	for key, value := range filter.ExtArgs {
		filterMap[key] = value
	}
	if len(filter.SniFilter) > 0 {
		filterMap[SniFilterKey] = filter.SniFilter
	}
	if len(filter.UrlFilter) > 0 {
		filterMap[UrlFilterKey] = filter.UrlFilter
	}
	if len(filter.ExternalURL) > 0 {
		filterMap[ExternalURLKey] = filter.ExternalURL
	}
	if len(filter.ExternalSNI) > 0 {
		filterMap[ExternalSNIKey] = filter.ExternalSNI
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

// newWebGroupMatchExpression converts one API "all" map into a match expression.
func newWebGroupMatchExpression(filterMap map[string]string) *WebGroupMatchExpression {
	filter := &WebGroupMatchExpression{
		ExternalURL: filterMap[ExternalURLKey],
		ExternalSNI: filterMap[ExternalSNIKey],
	}
	if filter.ExternalURL == "" && filter.ExternalSNI == "" {
		filter.SniFilter = filterMap[SniFilterKey]
		filter.UrlFilter = filterMap[UrlFilterKey]
		return filter
	}

	// Every key other than the feed key is an argument of the external feed.
	extArgs := make(map[string]string)
	for key, value := range filterMap {
		if key == ExternalURLKey || key == ExternalSNIKey {
			continue
		}
		extArgs[key] = value
	}
	if len(extArgs) > 0 {
		filter.ExtArgs = extArgs
	}
	return filter
}

// newWebGroup converts an API app domain response into a WebGroup.
func newWebGroup(data WebGroupResult) *WebGroup {
	webGroup := &WebGroup{
		Name: data.Name,
		UUID: data.UUID,
	}
	for _, filterResult := range data.Selector.Any {
		webGroup.Selector.Expressions = append(webGroup.Selector.Expressions, newWebGroupMatchExpression(filterResult.All))
	}
	return webGroup
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

func (c *Client) GetWebGroupByName(ctx context.Context, name string) (*WebGroup, error) {
	endpoint := fmt.Sprintf("app-domains/name/%s", name)

	var data WebGroupResult
	err := c.GetAPIContext25(ctx, &data, endpoint, nil)
	if err != nil {
		if strings.Contains(err.Error(), "App domain not found") {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if data.Name == name {
		return newWebGroup(data), nil
	}
	return nil, ErrNotFound
}

func (c *Client) GetWebGroup(ctx context.Context, uuid string) (*WebGroup, error) {
	endpoint := fmt.Sprintf("app-domains/%s", uuid)

	var data WebGroupResult
	err := c.GetAPIContext25(ctx, &data, endpoint, nil)
	if err != nil {
		if strings.Contains(err.Error(), "App domain not found") {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if data.UUID == uuid {
		return newWebGroup(data), nil
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
