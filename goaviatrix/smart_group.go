package goaviatrix

import (
	"context"
	"fmt"
	"log"
	"strings"
)

type SmartGroupMatchExpression struct {
	CIDR         string `json:"cidr,omitempty"`
	FQDN         string `json:"fqdn,omitempty"`
	Type         string `json:"type,omitempty"`
	Site         string `json:"site,omitempty"`
	ResId        string `json:"res_id,omitempty"`
	AccountId    string `json:"account_id,omitempty"`
	AccountName  string `json:"account_name,omitempty"`
	Name         string `json:"name,omitempty"`
	Region       string `json:"region,omitempty"`
	Zone         string `json:"zone,omitempty"`
	K8sService   string `json:"k8s_service,omitempty"`
	K8sNamespace string `json:"k8s_namespace,omitempty"`
	K8sClusterId string `json:"k8s_cluster_id,omitempty"`
	K8sPodName   string `json:"k8s_pod,omitempty"`
	S2C          string `json:"s2c,omitempty"`
	External     string `json:"external,omitempty"`
	Tags         map[string]string
	ExtArgs      map[string]string
}

type SmartGroupSelector struct {
	Expressions []*SmartGroupMatchExpression
}

type SmartGroup struct {
	Name     string
	UUID     string
	Selector SmartGroupSelector
}

type SmartGroupResult struct {
	UUID     string              `json:"uuid"`
	Name     string              `json:"name"`
	Selector SmartGroupAnyResult `json:"selector"`
}

type SmartGroupAnyResult struct {
	Any []SmartGroupMatchExpressionResult `json:"any"`
}

type SmartGroupMatchExpressionResult struct {
	All map[string]interface{} `json:"all"`
}

type SmartGroupResp struct {
	UUID        string             `json:"uuid"`
	SmartGroups []SmartGroupResult `json:"app_domains"`
}

const (
	TagsPrefix      = "tags"
	ExtArgsPrefix   = "ext_args"
	CidrKey         = "cidr"
	FqdnKey         = "fqdn"
	TypeKey         = "type"
	SiteKey         = "site"
	ResIdKey        = "res_id"
	AccountIdKey    = "account_id"
	AccountNameKey  = "account_name"
	NameKey         = "name"
	RegionKey       = "region"
	ZoneKey         = "zone"
	K8sClusterIdKey = "k8s_cluster_id"
	K8sNamespaceKey = "k8s_namespace"
	K8sServiceKey   = "k8s_service"
	K8sPodNameKey   = "k8s_pod"
	S2CKey          = "s2c"
	ExternalKey     = "external"

	AnyKey      = "any"
	AllKey      = "all"
	SelectorKey = "selector"

	ApiEndpoint = "app-domains"
)

func NewSmartGroupMatchExpression(filterMap map[string]interface{}) *SmartGroupMatchExpression {
	smartGroup := &SmartGroupMatchExpression{}
	setFilterInterface(&smartGroup.CIDR, filterMap, CidrKey)
	setFilterInterface(&smartGroup.FQDN, filterMap, FqdnKey)
	setFilterInterface(&smartGroup.Type, filterMap, TypeKey)
	setFilterInterface(&smartGroup.Site, filterMap, SiteKey)
	setFilterInterface(&smartGroup.ResId, filterMap, ResIdKey)
	setFilterInterface(&smartGroup.AccountId, filterMap, AccountIdKey)
	setFilterInterface(&smartGroup.AccountName, filterMap, AccountNameKey)
	setFilterInterface(&smartGroup.Name, filterMap, NameKey)
	setFilterInterface(&smartGroup.Region, filterMap, RegionKey)
	setFilterInterface(&smartGroup.Zone, filterMap, ZoneKey)
	setFilterInterface(&smartGroup.K8sClusterId, filterMap, K8sClusterIdKey)
	setFilterInterface(&smartGroup.K8sNamespace, filterMap, K8sNamespaceKey)
	setFilterInterface(&smartGroup.K8sService, filterMap, K8sServiceKey)
	setFilterInterface(&smartGroup.K8sPodName, filterMap, K8sPodNameKey)
	setFilterInterface(&smartGroup.S2C, filterMap, S2CKey)
	setFilterInterface(&smartGroup.External, filterMap, ExternalKey)
	return smartGroup
}

func setFilterInterface(filterField *string, filterMap map[string]interface{}, fieldKey string) {
	val, ok := filterMap[fieldKey]
	if !ok || val == nil {
		return
	}

	s, ok := val.(string)
	if !ok {
		return
	}

	if fieldKey == RegionKey {
		// Ensure that the region is always in lowercase, no-space
		s = strings.ToLower(strings.ReplaceAll(s, " ", ""))
	}

	*filterField = s
}

func setFilter(filterField string, filterMap map[string]interface{}, fieldKey string) {
	if len(filterField) > 0 {
		filterMap[fieldKey] = filterField
	}
}

// SmartGroupFilterToResource returns the contents of the filter structure in a map that will
// pass the smart group TF resource schema
func SmartGroupFilterToResource(filter *SmartGroupMatchExpression) map[string]interface{} {
	return smartGroupFilterToMapBasic(filter, true)
}

// SmartGroupFilterToAPIMap returns the contents of the filter structure in a map that can
// be passed to the Smart Group API
func SmartGroupFilterToAPIMap(filter *SmartGroupMatchExpression) map[string]interface{} {
	return smartGroupFilterToMapBasic(filter, false)
}

// smartGroupFilterToMapBasic is underlying function to SmartGroupFilterToResource and SmartGroupFiltertoAPIMap
// The keepMaps argument dictates how the map keys ("tags" and "ext_args") are translated for each case.
func smartGroupFilterToMapBasic(filter *SmartGroupMatchExpression, keepMaps bool) map[string]interface{} {
	filterMap := make(map[string]interface{})

	setFilter(filter.Type, filterMap, TypeKey)
	setFilter(filter.CIDR, filterMap, CidrKey)
	setFilter(filter.FQDN, filterMap, FqdnKey)
	setFilter(filter.Site, filterMap, SiteKey)
	setFilter(filter.ResId, filterMap, ResIdKey)
	setFilter(filter.AccountId, filterMap, AccountIdKey)
	setFilter(filter.AccountName, filterMap, AccountNameKey)
	setFilter(filter.Name, filterMap, NameKey)
	setFilter(filter.Region, filterMap, RegionKey)
	setFilter(filter.Zone, filterMap, ZoneKey)
	setFilter(filter.K8sClusterId, filterMap, K8sClusterIdKey)
	setFilter(filter.K8sNamespace, filterMap, K8sNamespaceKey)
	setFilter(filter.K8sService, filterMap, K8sServiceKey)
	setFilter(filter.K8sPodName, filterMap, K8sPodNameKey)
	setFilter(filter.S2C, filterMap, S2CKey)
	setFilter(filter.External, filterMap, ExternalKey)
	if keepMaps {
		if len(filter.Tags) > 0 {
			filterMap[TagsPrefix] = filter.Tags
		}
		if len(filter.ExtArgs) > 0 {
			filterMap[ExtArgsPrefix] = filter.ExtArgs
		}
	} else {
		if len(filter.Tags) > 0 {
			for key, value := range filter.Tags {
				filterMap[fmt.Sprintf("%s.%s", TagsPrefix, key)] = value
			}
		}
		if len(filter.ExtArgs) > 0 {
			for key, value := range filter.ExtArgs {
				filterMap[key] = value
			}
		}
	}

	log.Printf("[DEBUG] POST filter map: %v\n", filterMap)
	return filterMap
}

func makeSmartGroupForm(smartGroup *SmartGroup) map[string]interface{} {
	form := map[string]interface{}{
		NameKey: smartGroup.Name,
	}

	var orGroup []map[string]map[string]interface{}
	for _, smartGroupSelector := range smartGroup.Selector.Expressions {
		andGroup := map[string]map[string]interface{}{
			AllKey: SmartGroupFilterToAPIMap(smartGroupSelector),
		}

		orGroup = append(orGroup, andGroup)
	}

	form[SelectorKey] = map[string]interface{}{
		AnyKey: orGroup,
	}

	return form
}

func (c *Client) CreateSmartGroup(ctx context.Context, smartGroup *SmartGroup) (string, error) {
	form := makeSmartGroupForm(smartGroup)
	var data SmartGroupResp
	if err := c.PostAPIContext25(ctx, &data, ApiEndpoint, form); err != nil {
		return "", err
	}
	return data.UUID, nil
}

func (c *Client) GetSmartGroup(ctx context.Context, uuid string) (*SmartGroup, error) {
	endpoint := fmt.Sprintf("app-domains/%s", uuid)

	var response SmartGroupResult

	err := c.GetAPIContext25(ctx, &response, endpoint, nil)
	if err != nil {
		return nil, err
	}

	return createSmartGroup(response), nil
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

	var data SmartGroupResp
	err := c.GetAPIContext25(ctx, &data, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var smartGroups []*SmartGroup
	for _, smartGroupResult := range data.SmartGroups {
		if smartGroupResult.UUID != "" {
			smartGroups = append(smartGroups, createSmartGroup(smartGroupResult))
		}
	}
	return smartGroups, nil
}

// createSmartGroup transforms the result returned by the API into the structure expected by the
// rest of the terraform smart group resource logic
func createSmartGroup(smartGroupResult SmartGroupResult) *SmartGroup {
	smartGroup := &SmartGroup{
		Name: smartGroupResult.Name,
		UUID: smartGroupResult.UUID,
	}

	for _, filterResult := range smartGroupResult.Selector.Any {
		filterMap := filterResult.All
		var filter *SmartGroupMatchExpression

		if MapContains(filterMap, ExternalKey) {
			filter = &SmartGroupMatchExpression{}

			raw, ok := filterMap[ExternalKey]
			if ok && raw != nil {
				if s, ok := raw.(string); ok {
					filter.External = s
				}
			}
		} else {
			filter = NewSmartGroupMatchExpression(filterMap)
		}

		if MapContains(filterMap, ExternalKey) {
			extArgs := make(map[string]string)
			for key, value := range filterMap {
				if key == ExternalKey {
					continue
				}
				if value == nil {
					continue
				}
				s, ok := value.(string)
				if !ok {
					continue
				}
				extArgs[key] = s
			}
			if len(extArgs) > 0 {
				filter.ExtArgs = extArgs
			}
		} else if MapContains(filterMap, TypeKey) {
			tags := make(map[string]string)
			for key, value := range filterMap {
				if !strings.HasPrefix(key, TagsPrefix+".") {
					continue
				}
				if value == nil {
					continue
				}
				s, ok := value.(string)
				if !ok {
					continue
				}
				tags[strings.TrimPrefix(key, TagsPrefix+".")] = s
			}
			if len(tags) > 0 {
				filter.Tags = tags
			}
		}

		smartGroup.Selector.Expressions = append(smartGroup.Selector.Expressions, filter)
	}

	return smartGroup
}
