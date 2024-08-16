package goaviatrix

import (
	"context"
	"time"
)

type EdgeNEOProxyProfile struct {
	Action      string `json:"action"`
	CID         string `json:"CID"`
	AccountName string `json:"account_name"`
	Name        string `json:"proxy_name"`
	Address     string `json:"address"`
	Port        int    `json:"port"`
	CACert      string `json:"ca_cert"`
}

type EdgeNEOProxyProfileResp struct {
	ProxyID           string     `json:"proxyID"`
	Name              string     `json:"name"`
	IPAddress         string     `json:"address"`
	Port              int64      `json:"port"`
	ProxyProfileCount int64      `json:"deviceCount"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
	CaCert            *string    `json:"caCert,omitempty"`
	ExpiredAt         *time.Time `json:"expiredAt,omitempty"`
}

type EdgeNEOProxyProfileListResponse struct {
	Return    bool `json:"return"`
	Results   []EdgeNEOProxyProfileResp
	Reason    string `json:"reason"`
	Errortype string `json:"errortype"`
}

func (c *Client) CreateEdgeProxyProfile(ctx context.Context, edgeNEOProxyProfile *EdgeNEOProxyProfile) error {
	edgeNEOProxyProfile.Action = "create_edge_csp_proxy_profile"
	edgeNEOProxyProfile.CID = c.CID

	err := c.PostAPIContext2(ctx, nil, edgeNEOProxyProfile.Action, edgeNEOProxyProfile, BasicCheck)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetEdgeNEOProxyProfile(ctx context.Context, accountName, profileName string) (*EdgeNEOProxyProfileResp, error) {
	form := map[string]string{
		"action":       "list_edge_csp_proxy_profiles",
		"CID":          c.CID,
		"account_name": accountName,
	}

	var data EdgeNEOProxyProfileListResponse

	err := c.PostAPIContext2(ctx, &data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	edgeNEOProxyProfileList := data.Results
	for _, edgeNEOProxyProfile := range edgeNEOProxyProfileList {
		if edgeNEOProxyProfile.Name == profileName {
			return &edgeNEOProxyProfile, nil
		}
	}

	return nil, ErrNotFound
}

func (c *Client) DeleteEdgeNEOProxyProfile(ctx context.Context, accountName, profileName string) error {
	proxyProfile, err := c.GetEdgeNEOProxyProfile(ctx, accountName, profileName)
	if err == ErrNotFound {
		return nil
	} else if err != nil {
		return err
	}

	form := map[string]string{
		"action":       "delete_edge_csp_proxy_profile",
		"CID":          c.CID,
		"account_name": accountName,
		"proxy_id":     proxyProfile.ProxyID,
	}
	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}
