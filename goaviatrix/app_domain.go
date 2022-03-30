package goaviatrix

import (
	"context"
	"fmt"
	"log"
	"strings"
)

type AppDomainIPFilter struct {
	Type string
	Ips  []string
}

type AppDomainTagFilter struct {
	Type string
	Tags map[string]string
	//Resources []string
}

type AppDomainMatchExpression struct {
	CIDR        string `json:"cidr,omitempty"`
	Type        string `json:"type,omitempty"`
	ResId       string `json:"res_id,omitempty"`
	AccountId   string `json:"account_id,omitempty"`
	AccountName string `json:"account_name,omitempty"`
	Region      string `json:"region,omitempty"`
	Zone        string `json:"zone,omitempty"`
	Tags        map[string]string
}

type AppDomainSelector struct {
	Expressions []*AppDomainMatchExpression
}

type AppDomain struct {
	Name     string
	UUID     string
	Selector AppDomainSelector
	//IpFilter  *AppDomainIPFilter
	//TagFilter *AppDomainTagFilter
}

func appDomainFilterToMap(filter *AppDomainMatchExpression) map[string]string {
	filterMap := make(map[string]string)

	if len(filter.Type) > 0 {
		filterMap["type"] = filter.Type
	}
	if len(filter.CIDR) > 0 {
		filterMap["cidr"] = filter.CIDR
	}
	if len(filter.ResId) > 0 {
		filterMap["res_id"] = filter.ResId
	}
	if len(filter.AccountId) > 0 {
		filterMap["account_id"] = filter.AccountId
	}
	if len(filter.AccountName) > 0 {
		filterMap["account_name"] = filter.AccountName
	}
	if len(filter.Region) > 0 {
		filterMap["region"] = filter.Region
	}
	if len(filter.Zone) > 0 {
		filterMap["zone"] = filter.Zone
	}

	if len(filter.Tags) > 0 {
		for key, value := range filter.Tags {
			filterMap[fmt.Sprintf("tags.%s", key)] = value
		}
	}

	log.Printf("[DEBUG] POST filter map: %v\n", filterMap)
	return filterMap
}

func (c *Client) CreateAppDomain(ctx context.Context, appDomain *AppDomain) (string, error) {
	endpoint := "app-domains"

	form := map[string]interface{}{
		"name": appDomain.Name,
	}

	var or []map[string]map[string]string
	for _, appDomainSelector := range appDomain.Selector.Expressions {
		and := map[string]map[string]string{
			"_and": appDomainFilterToMap(appDomainSelector),
		}

		or = append(or, and)
	}

	form["selector"] = map[string]interface{}{
		"_or": or,
	}

	type AppDomainResp struct {
		UUID string `json:"uuid"`
	}

	var data AppDomainResp
	err := c.PostAPIContext25(ctx, &data, endpoint, form)
	if err != nil {
		return "", err
	}
	return data.UUID, nil
}

func (c *Client) GetAppDomain(ctx context.Context, uuid string) (*AppDomain, error) {
	endpoint := "app-domains"

	type AppDomainResultIps struct {
		Ips []string `json:"ip_or_cidrs"`
	}

	type AppDomainResultTags struct {
		Tags []map[string]string `json:"tags"`
	}

	type AppDomainFilter struct {
		And map[string]string `json:"_and"`
	}

	type AppDomainOrResult struct {
		Or []AppDomainFilter `json:"_or"`
	}

	type AppDomainResult struct {
		UUID     string            `json:"uuid"`
		Name     string            `json:"name"`
		Selector AppDomainOrResult `json:"selector"`
		//IpFilter  AppDomainResultIps  `json:"ip_filter,omitempty"`
		//TagFilter AppDomainResultTags `json:"tag_filter,omitempty"`
	}

	type AppDomainResp struct {
		AppDomains []AppDomainResult `json:"app_domains"`
	}

	var data AppDomainResp
	err := c.GetAPIContext25(ctx, &data, endpoint, nil)
	if err != nil {
		return nil, err
	}

	for _, appDomainResult := range data.AppDomains {
		if appDomainResult.UUID == uuid {
			appDomain := &AppDomain{
				Name: appDomainResult.Name,
				UUID: appDomainResult.UUID,
			}

			for _, filterResult := range appDomainResult.Selector.Or {
				filterMap := filterResult.And

				filter := &AppDomainMatchExpression{
					CIDR:        filterMap["cidr"],
					Type:        filterMap["type"],
					ResId:       filterMap["res_id"],
					AccountId:   filterMap["account_id"],
					AccountName: filterMap["account_name"],
					Region:      filterMap["region"],
					Zone:        filterMap["zone"],
				}

				tags := make(map[string]string)
				for key, value := range filterMap {
					if strings.HasPrefix(key, "tags.") {
						tags[strings.TrimPrefix(key, "tags.")] = value
					} else {
					}
				}

				if len(tags) > 0 {
					filter.Tags = tags
				}

				appDomain.Selector.Expressions = append(appDomain.Selector.Expressions, filter)

				//var filter AppDomainMatchExpression
				//json.Unmarshal()
				//filter := make(map[string])

			}

			//if len(appDomainResult.IpFilter.Ips) > 0 {
			//	appDomain.IpFilter = &AppDomainIPFilter{
			//		Ips: appDomainResult.IpFilter.Ips,
			//	}
			//}
			//
			//if len(appDomainResult.TagFilter.Tags) > 0 {
			//	appDomain.TagFilter = &AppDomainTagFilter{
			//		Tags: make(map[string]string),
			//	}
			//
			//	for _, keyValPair := range appDomainResult.TagFilter.Tags {
			//		key, keyOk := keyValPair["key"]
			//		val, valOk := keyValPair["val"]
			//		if keyOk && valOk {
			//			appDomain.TagFilter.Tags[key] = val
			//		} else {
			//			log.Printf("[TRACE] Invalid App Domain tag filter: %v\n", keyValPair)
			//		}
			//	}
			//
			//}

			return appDomain, nil
		}
	}
	return nil, ErrNotFound
}

func (c *Client) UpdateAppDomain(ctx context.Context, appDomain *AppDomain, uuid string) error {
	endpoint := fmt.Sprintf("app-domains/%s", uuid)

	form := map[string]interface{}{
		"name": appDomain.Name,
	}

	var or []map[string]map[string]string
	for _, appDomainSelector := range appDomain.Selector.Expressions {
		and := map[string]map[string]string{
			"_and": appDomainFilterToMap(appDomainSelector),
		}

		or = append(or, and)
	}

	form["selector"] = map[string]interface{}{
		"_or": or,
	}

	return c.PutAPIContext25(ctx, endpoint, form)
}

func (c *Client) DeleteAppDomain(ctx context.Context, uuid string) error {
	endpoint := fmt.Sprintf("app-domains/%s", uuid)
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}
