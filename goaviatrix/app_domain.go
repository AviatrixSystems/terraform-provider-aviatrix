package goaviatrix

import (
	"context"
	"fmt"
	"log"
	"strings"
)

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

func makeAppDomainForm(appDomain *AppDomain) map[string]interface{} {
	form := map[string]interface{}{
		"name": appDomain.Name,
	}

	var or []map[string]map[string]string
	for _, appDomainSelector := range appDomain.Selector.Expressions {
		and := map[string]map[string]string{
			"all": appDomainFilterToMap(appDomainSelector),
		}

		or = append(or, and)
	}

	form["selector"] = map[string]interface{}{
		"any": or,
	}

	return form
}

func (c *Client) CreateAppDomain(ctx context.Context, appDomain *AppDomain) (string, error) {
	endpoint := "app-domains"
	form := makeAppDomainForm(appDomain)

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

	type AppDomainMatchExpressionResult struct {
		All map[string]string `json:"all"`
	}

	type AppDomainAnyResult struct {
		Any []AppDomainMatchExpressionResult `json:"any"`
	}

	type AppDomainResult struct {
		UUID     string             `json:"uuid"`
		Name     string             `json:"name"`
		Selector AppDomainAnyResult `json:"selector"`
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

			for _, filterResult := range appDomainResult.Selector.Any {
				filterMap := filterResult.All

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
					}
				}

				if len(tags) > 0 {
					filter.Tags = tags
				}

				appDomain.Selector.Expressions = append(appDomain.Selector.Expressions, filter)
			}
			return appDomain, nil
		}
	}
	return nil, ErrNotFound
}

func (c *Client) UpdateAppDomain(ctx context.Context, appDomain *AppDomain, uuid string) error {
	endpoint := fmt.Sprintf("app-domains/%s", uuid)
	form := makeAppDomainForm(appDomain)
	return c.PutAPIContext25(ctx, endpoint, form)
}

func (c *Client) DeleteAppDomain(ctx context.Context, uuid string) error {
	endpoint := fmt.Sprintf("app-domains/%s", uuid)
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}
