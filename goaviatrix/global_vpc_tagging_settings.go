package goaviatrix

import (
	"context"
)

type GlobalVpcTaggingSettings struct {
	ServiceState string `json:"service_state"`
	EnableAlert  bool   `json:"alert"`
}

type GlobalVpcTaggingSettingsResp struct {
	Return  bool                     `json:"return"`
	Results GlobalVpcTaggingSettings `json:"results"`
	Reason  string                   `json:"reason"`
}

func (c *Client) UpdateGlobalVpcTaggingSettings(ctx context.Context, globalVpcTaggingSettings *GlobalVpcTaggingSettings) error {
	endpoint := "globalvpc/tagging_settings"

	err := c.PutAPIContext25(ctx, endpoint, globalVpcTaggingSettings)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetGlobalVpcTaggingSettings(ctx context.Context) (*GlobalVpcTaggingSettings, error) {
	endpoint := "globalvpc/tagging_settings"

	var data GlobalVpcTaggingSettingsResp
	err := c.GetAPIContext25(ctx, &data, endpoint, nil)
	if err != nil {
		return nil, err
	}

	return &data.Results, nil
}
