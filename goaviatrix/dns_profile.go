package goaviatrix

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
)

type DNSProfile struct {
	Global           []string `json:"global"`
	Lan              []string `json:"lan"`
	LocalDomainNames []string `json:"local_domain_names"`
	Wan              []string `json:"wan"`
}

type DNSProfileListResp struct {
	Return  bool                   `json:"return"`
	Results map[string]interface{} `json:"results"`
	Reason  string                 `json:"reason"`
}

func (c *Client) CreateDNSProfile(ctx context.Context, data map[string]interface{}) error {
	form := map[string]string{
		"action": "create_dns_profile",
		"CID":    c.CID,
	}

	profiles, err := json.Marshal(data)
	if err != nil {
		return err
	}

	form["data"] = b64.StdEncoding.EncodeToString(profiles)

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

func (c *Client) GetDNSProfile(ctx context.Context, name string) (map[string]interface{}, error) {
	form := map[string]string{
		"action": "list_dns_profile",
		"CID":    c.CID,
	}

	var data DNSProfileListResp

	err := c.PostAPIContext2(ctx, &data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	templateNames := data.Results["template_names"].([]interface{})

	for _, n := range templateNames {
		if n.(string) == name {
			return data.Results[name].(map[string]interface{}), nil
		}
	}

	return nil, ErrNotFound
}

func (c *Client) UpdateDNSProfile(ctx context.Context, data map[string]interface{}) error {
	form := map[string]string{
		"action": "update_dns_profile",
		"CID":    c.CID,
	}

	profiles, err := json.Marshal(data)
	if err != nil {
		return err
	}

	form["data"] = b64.StdEncoding.EncodeToString(profiles)

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}

func (c *Client) DeleteDNSProfile(ctx context.Context, data map[string]interface{}) error {
	form := map[string]string{
		"action": "delete_dns_profile",
		"CID":    c.CID,
	}

	profiles, err := json.Marshal(data)
	if err != nil {
		return err
	}

	form["data"] = b64.StdEncoding.EncodeToString(profiles)

	return c.PostAPIContext2(ctx, nil, form["action"], form, BasicCheck)
}
