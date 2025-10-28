package goaviatrix

import (
	"context"
	"fmt"
	"net/url"
)

type LogProfile struct {
	ProfileName  string `json:"profile_name"`
	ProfileID    string `json:"profile_id"`
	SessionEnd   bool   `json:"session_end"`
	SessionStart bool   `json:"session_start"`
}

func (c *Client) GetLogProfileByName(ctx context.Context, profileName string) (*LogProfile, error) {
	profileName = url.QueryEscape(profileName)
	endpoint := fmt.Sprintf("dcf/log-profile/name/%s", profileName)

	var data LogProfile
	err := c.GetAPIContext25(ctx, &data, endpoint, nil)
	if err != nil {
		return nil, err
	}

	if data.ProfileName == profileName {
		return &data, nil
	}
	return nil, ErrNotFound
}
