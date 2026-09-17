package goaviatrix

import (
	"context"
	"errors"
	"time"
)

type EdgePlatformProxyProfile struct {
	Action      string `json:"action"`
	CID         string `json:"CID"`
	AccountName string `json:"account_name"`
	Name        string `json:"proxy_name"`
	Address     string `json:"address"`
	Port        int    `json:"port"`
	CACert      string `json:"ca_cert"`
}

type EdgePlatformProxyProfileUpdate struct {
	EdgePlatformProxyProfile
	ProxyID string `json:"proxy_id"`
}

type EdgePlatformProxyProfileResp struct {
	ProxyID           string    `json:"proxyID"`
	Name              string    `json:"name"`
	IPAddress         string    `json:"address"`
	Port              int       `json:"port"`
	ProxyProfileCount int       `json:"deviceCount"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	CaCert            string    `json:"caCert,omitempty"`
	ExpiredAt         time.Time `json:"expiredAt,omitempty"`
}

type EdgePlatformProxyProfileCreateResponse struct {
	Return    bool `json:"return"`
	Results   EdgePlatformProxyProfileResp
	Reason    string `json:"reason"`
	Errortype string `json:"errortype"`
}

type EdgePlatformProxyProfileUpdateResponse struct {
	Return    bool   `json:"return"`
	Results   string `json:"results"`
	Reason    string `json:"reason"`
	Errortype string `json:"errortype"`
}

type EdgePlatformProxyProfileListResponse struct {
	Return    bool `json:"return"`
	Results   []EdgePlatformProxyProfileResp
	Reason    string `json:"reason"`
	Errortype string `json:"errortype"`
}

func (c *Client) CreateEdgeProxyProfile(ctx context.Context, edgeNEOProxyProfile *EdgePlatformProxyProfile) (*EdgePlatformProxyProfileResp, error) {
	edgeNEOProxyProfile.Action = "create_edge_csp_proxy_profile"
	edgeNEOProxyProfile.CID = c.CID

	var data EdgePlatformProxyProfileCreateResponse
	err := c.PostAPIContext2(ctx, &data, edgeNEOProxyProfile.Action, edgeNEOProxyProfile, BasicCheck)
	if err != nil {
		return nil, err
	}

	return &data.Results, nil
}

func (c *Client) UpdateEdgeProxyProfile(ctx context.Context, edgeNEOProxyProfile *EdgePlatformProxyProfileUpdate) error {
	edgeNEOProxyProfile.Action = "update_edge_csp_proxy_profile"
	edgeNEOProxyProfile.CID = c.CID

	var data EdgePlatformProxyProfileUpdateResponse
	err := c.PostAPIContext2(ctx, &data, edgeNEOProxyProfile.Action, edgeNEOProxyProfile, BasicCheck)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetEdgePlatformProxyProfile(ctx context.Context, accountName, profileName string) (*EdgePlatformProxyProfileResp, error) {
	form := map[string]string{
		"action":       "list_edge_csp_proxy_profiles",
		"CID":          c.CID,
		"account_name": accountName,
	}

	var data EdgePlatformProxyProfileListResponse

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

func (c *Client) DeleteEdgePlatformProxyProfile(ctx context.Context, accountName, profileName string) error {
	proxyProfile, err := c.GetEdgePlatformProxyProfile(ctx, accountName, profileName)
	if errors.Is(err, ErrNotFound) {
		return nil
	}
	if err != nil {
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
