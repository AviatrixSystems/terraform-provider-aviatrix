package goaviatrix

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"reflect"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DNSProfileListResp struct {
	Return  bool                   `json:"return"`
	Results map[string]interface{} `json:"results"`
	Reason  string                 `json:"reason"`
}

func (c *Client) CreateDNSProfileList(ctx context.Context, data map[string]interface{}) error {
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

func (c *Client) GetDNSProfileList(ctx context.Context) (map[string]interface{}, error) {
	form := map[string]string{
		"action": "list_dns_profile",
		"CID":    c.CID,
	}

	var data DNSProfileListResp

	err := c.PostAPIContext2(ctx, &data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	} else if len(data.Results) == 0 {
		return nil, ErrNotFound
	}

	return data.Results, nil
}

func (c *Client) UpdateDNSProfileList(ctx context.Context, data map[string]interface{}) error {
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

func (c *Client) DeleteDNSProfileList(ctx context.Context, data map[string]interface{}) error {
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

func DiffSuppressFuncDNSProfileList(k, old, new string, d *schema.ResourceData) bool {
	pOld, pNew := d.GetChange("profiles")
	var profilesOld []map[string]interface{}

	for _, p0 := range pOld.([]interface{}) {
		p1 := p0.(map[string]interface{})
		profilesOld = append(profilesOld, p1)
	}

	var profilesNew []map[string]interface{}

	for _, p0 := range pNew.([]interface{}) {
		p1 := p0.(map[string]interface{})
		profilesNew = append(profilesNew, p1)
	}

	sort.Slice(profilesOld, func(i, j int) bool {
		return profilesOld[i]["name"].(string) < profilesOld[j]["name"].(string)
	})

	sort.Slice(profilesNew, func(i, j int) bool {
		return profilesNew[i]["name"].(string) < profilesNew[j]["name"].(string)
	})

	return reflect.DeepEqual(profilesOld, profilesNew)
}
