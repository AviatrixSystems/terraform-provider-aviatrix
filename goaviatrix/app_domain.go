package goaviatrix

import (
	"context"
)

type AppDomain struct {
	Filters []*AppDomainFilter
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
		"CID":     c.CID,
		"filters": appDomain.Filters,
	}

	type AppDomainResults struct {
		UUID string `json:"uuid"`
	}

	type AppDomainResp struct {
		Return  bool
		Reason  string
		Results AppDomainResults
	}

	var data AppDomainResp
	err := c.PostAPI25Context(ctx, &data, endpoint, form, BasicCheck)
	if err != nil {
		return "", err
	}
	return data.Results.UUID, nil
}

func (c *Client) GetAppDomain(ctx context.Context, uuid string) (*AppDomain, error) {
	//action := ""
	//data := map[string]string {
	//	"CID": c.CID,
	//}
	return nil, ErrNotFound
}

func (c *Client) UpdateAppDomain(ctx context.Context, appDomain *AppDomain, uuid string) error {
	return nil
}

func (c *Client) DeleteAppDomain(ctx context.Context, uuid string) error {
	return nil
}
