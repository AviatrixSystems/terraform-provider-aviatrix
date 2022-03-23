package goaviatrix

import (
	"context"
	"fmt"
	"log"
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

type AppDomain struct {
	Name      string
	IpFilter  *AppDomainIPFilter
	TagFilter *AppDomainTagFilter
}

type AppDomainFilter struct {
	Type      string            `json:"type"`
	Ips       []string          `json:"ips,omitempty"`
	Tags      map[string]string `json:"tags,omitempty"`
	Resources []string          `json:"resources,omitempty"`
}

func (c *Client) CreateAppDomain(ctx context.Context, appDomain *AppDomain) (string, error) {
	endpoint := "app-domains"

	form := map[string]interface{}{
		"name": appDomain.Name,
	}

	if appDomain.IpFilter != nil {
		form["ip_filter"] = map[string]interface{}{
			"ip_or_cidrs": appDomain.IpFilter.Ips,
		}
	}

	if appDomain.TagFilter != nil {
		var tagList []map[string]string
		for key, value := range appDomain.TagFilter.Tags {
			tagList = append(tagList, map[string]string{
				"key": key,
				"val": value,
			})
		}

		form["tag_filter"] = map[string]interface{}{
			"tags": tagList,
		}
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

	type AppDomainResult struct {
		UUID      string              `json:"uuid"`
		Name      string              `json:"name"`
		IpFilter  AppDomainResultIps  `json:"ip_filter,omitempty"`
		TagFilter AppDomainResultTags `json:"tag_filter,omitempty"`
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
			}

			if len(appDomainResult.IpFilter.Ips) > 0 {
				appDomain.IpFilter = &AppDomainIPFilter{
					Ips: appDomainResult.IpFilter.Ips,
				}
			}

			if len(appDomainResult.TagFilter.Tags) > 0 {
				appDomain.TagFilter = &AppDomainTagFilter{
					Tags: make(map[string]string),
				}

				for _, keyValPair := range appDomainResult.TagFilter.Tags {
					key, keyOk := keyValPair["key"]
					val, valOk := keyValPair["val"]
					if keyOk && valOk {
						appDomain.TagFilter.Tags[key] = val
					} else {
						log.Printf("[TRACE] Invalid App Domain tag filter: %v\n", keyValPair)
					}
				}

			}

			return appDomain, nil
		}
	}
	//action := ""
	//data := map[string]string {
	//	"CID": c.CID,
	//}
	return nil, ErrNotFound
}

func (c *Client) UpdateAppDomain(ctx context.Context, appDomain *AppDomain, uuid string) error {
	endpoint := fmt.Sprintf("app-domains/%s", uuid)

	form := map[string]interface{}{
		"name": appDomain.Name,
	}

	if appDomain.IpFilter != nil {
		form["ip_filter"] = map[string]interface{}{
			"ip_or_cidrs": appDomain.IpFilter.Ips,
		}
	}

	if appDomain.TagFilter != nil {
		var tagList []map[string]string
		for key, value := range appDomain.TagFilter.Tags {
			tagList = append(tagList, map[string]string{
				"key": key,
				"val": value,
			})
		}

		form["tag_filter"] = map[string]interface{}{
			"tags": tagList,
		}
	}

	return c.PutAPIContext25(ctx, endpoint, form)
}

func (c *Client) DeleteAppDomain(ctx context.Context, uuid string) error {
	endpoint := fmt.Sprintf("app-domains/%s", uuid)
	return c.DeleteAPIContext25(ctx, endpoint, nil)
}
