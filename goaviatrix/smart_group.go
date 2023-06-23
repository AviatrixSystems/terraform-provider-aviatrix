package goaviatrix

import (
	"context"
	"fmt"
	"log"
	"strings"
)

type SmartGroupMatchExpression struct {
	CIDR        string `json:"cidr,omitempty"`
	FQDN        string `json:"fqdn,omitempty"`
	Type        string `json:"type,omitempty"`
	Site        string `json:"site,omitempty"`
	ResId       string `json:"res_id,omitempty"`
	AccountId   string `json:"account_id,omitempty"`
	AccountName string `json:"account_name,omitempty"`
	Name        string `json:"name,omitempty"`
	Region      string `json:"region,omitempty"`
	Zone        string `json:"zone,omitempty"`
	Tags        map[string]string
}

type SmartGroupSelector struct {
	Expressions []*SmartGroupMatchExpression
}

type SmartGroup struct {
	Name     string
	UUID     string
	Selector SmartGroupSelector
}

func smartGroupFilterToMap(filter *SmartGroupMatchExpression) map[string]string {
	filterMap := make(map[string]string)

	if len(filter.Type) > 0 {
		filterMap["type"] = filter.Type
	}
	if len(filter.CIDR) > 0 {
		filterMap["cidr"] = filter.CIDR
	}
	if len(filter.FQDN) > 0 {
		filterMap["fqdn"] = filter.FQDN
	}
	if len(filter.Site) > 0 {
		filterMap["site"] = filter.Site
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
	if len(filter.Name) > 0 {
		filterMap["name"] = filter.Name
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

func makeSmartGroupForm(smartGroup *SmartGroup) map[string]interface{} {
	form := map[string]interface{}{
		"name": smartGroup.Name,
	}

	var or []map[string]map[string]string
	for _, smartGroupSelector := range smartGroup.Selector.Expressions {
		and := map[string]map[string]string{
			"all": smartGroupFilterToMap(smartGroupSelector),
		}

		or = append(or, and)
	}

	form["selector"] = map[string]interface{}{
		"any": or,
	}

	return form
}

func (c *Client) CreateSmartGroup(ctx context.Context, smartGroup *SmartGroup) (string, error) {
	endpoint := "app-domains"
	form := makeSmartGroupForm(smartGroup)

	type SmartGroupResp struct {
		UUID string `json:"uuid"`
	}

	var data SmartGroupResp
	err := c.PostAPIContext25(ctx, &data, endpoint, form)
	if err != nil {
		return "", err
	}
	return data.UUID, nil
}

func (c *Client) GetSmartGroup(ctx context.Context, uuid string) (*SmartGroup, error) {
	endpoint := "app-domains"

	type SmartGroupMatchExpressionResult struct {
		All map[string]string `json:"all"`
	}

	type SmartGroupAnyResult struct {
		Any []SmartGroupMatchExpressionResult `json:"any"`
	}

	type SmartGroupResult struct {
		UUID     string              `json:"uuid"`
		Name     string              `json:"name"`
		Selector SmartGroupAnyResult `json:"selector"`
	}

	type SmartGroupResp struct {
		SmartGroups []SmartGroupResult `json:"app_domains"`
	}

	var data SmartGroupResp
	err := c.GetAPIContext25(ctx, &data, endpoint, nil)
	if err != nil {
		return nil, err
	}

	for _, smartGroupResult := range data.SmartGroups {
		if smartGroupResult.UUID == uuid {
			smartGroup := &SmartGroup{
				Name: smartGroupResult.Name,
				UUID: smartGroupResult.UUID,
			}

			for _, filterResult := range smartGroupResult.Selector.Any {
				filterMap := filterResult.All

				filter := &SmartGroupMatchExpression{
					CIDR:        filterMap["cidr"],
					FQDN:        filterMap["fqdn"],
					Type:        filterMap["type"],
					Site:        filterMap["site"],
					ResId:       filterMap["res_id"],
					AccountId:   filterMap["account_id"],
					AccountName: filterMap["account_name"],
					Name:        filterMap["name"],
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

				smartGroup.Selector.Expressions = append(smartGroup.Selector.Expressions, filter)
			}
			return smartGroup, nil
		}
	}
	return nil, ErrNotFound
}

func (c *Client) UpdateSmartGroup(ctx context.Context, smartGroup *SmartGroup, uuid string) error {
	endpoint := fmt.Sprintf("app-domains/%s", uuid)
	form := makeSmartGroupForm(smartGroup)
	return c.PutAPIContext25(ctx, endpoint, form)
}

func (c *Client) DeleteSmartGroup(ctx context.Context, uuid string) error {
	endpoint := fmt.Sprintf("app-domains/%s", uuid)
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}

func (c *Client) GetSmartGroups(ctx context.Context) ([]*SmartGroup, error) {
	endpoint := "app-domains"

	type SmartGroupMatchExpressionResult struct {
		All map[string]string `json:"all"`
	}

	type SmartGroupAnyResult struct {
		Any []SmartGroupMatchExpressionResult `json:"any"`
	}

	type SmartGroupResult struct {
		UUID     string              `json:"uuid"`
		Name     string              `json:"name"`
		Selector SmartGroupAnyResult `json:"selector"`
	}

	type SmartGroupResp struct {
		SmartGroups []SmartGroupResult `json:"app_domains"`
	}

	var data SmartGroupResp
	err := c.GetAPIContext25(ctx, &data, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var smartGroups []*SmartGroup
	for _, smartGroupResult := range data.SmartGroups {
		if smartGroupResult.UUID != "" {
			smartGroup := &SmartGroup{
				Name: smartGroupResult.Name,
				UUID: smartGroupResult.UUID,
			}

			for _, filterResult := range smartGroupResult.Selector.Any {
				filterMap := filterResult.All

				filter := &SmartGroupMatchExpression{
					CIDR:        filterMap["cidr"],
					FQDN:        filterMap["fqdn"],
					Type:        filterMap["type"],
					Site:        filterMap["site"],
					ResId:       filterMap["res_id"],
					AccountId:   filterMap["account_id"],
					AccountName: filterMap["account_name"],
					Name:        filterMap["name"],
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

				smartGroup.Selector.Expressions = append(smartGroup.Selector.Expressions, filter)
			}
			smartGroups = append(smartGroups, smartGroup)
		}
	}
	return smartGroups, nil
}
