package goaviatrix

import (
	"context"
)

type LogProfile struct {
	ProfileName  string `json:"profile_name"`
	ProfileID    string `json:"profile_id"`
	SessionEnd   bool   `json:"session_end"`
	SessionStart bool   `json:"session_start"`
}

func (c *Client) GetLogProfileByName(ctx context.Context, profileName string) (*LogProfile, error) {
	endpoint := "dcf/log-profile"

	var data LogProfile
	err := c.GetAPIContext25(ctx, &data, endpoint, map[string]string{"profile_name": profileName})
	if err != nil {
		return nil, err
	}

	if data.ProfileName == profileName {
		return &data, nil
	}
	return nil, ErrNotFound
}
