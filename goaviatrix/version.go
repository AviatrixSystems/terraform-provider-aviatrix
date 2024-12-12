package goaviatrix

import (
	"context"
	"strconv"
)

type Version struct {
	CID           string `form:"CID,omitempty"`
	Action        string `form:"action,omitempty"`
	TargetVersion string `form:"version,omitempty"`
	Version       string `json:"version,omitempty"`
}

type VersionInfoResults struct {
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
}

type VersionInfoResp struct {
	Return  bool               `json:"return"`
	Results VersionInfoResults `json:"results"`
	Reason  string             `json:"reason"`
}

type VersionInfo struct {
	Current  string
	Previous string
}

func (c *Client) GetCurrentVersion() (string, error) {
	form := map[string]string{
		"CID":    c.CID,
		"action": "list_version_info",
	}

	var data VersionInfoResp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return "", err
	}
	return data.Results.CurrentVersion, nil
}

func (c *Client) GetVersionInfo() (*VersionInfo, error) {
	form := map[string]string{
		"action": "list_version_info",
		"CID":    c.CID,
	}
	var data struct {
		Results struct {
			PreviousVersion string `json:"previous_version"`
			CurrentVersion  string `json:"current_version"`
		}
	}
	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	return &VersionInfo{
		Current:  data.Results.CurrentVersion,
		Previous: data.Results.PreviousVersion,
	}, nil
}

func (c *Client) GetCompatibleImageVersion(ctx context.Context, cloudType int, softwareVersion string) (string, error) {
	form := map[string]string{
		"action":           "get_compatible_image_version",
		"CID":              c.CID,
		"software_version": softwareVersion,
		"cloud_type":       strconv.Itoa(cloudType),
	}
	var data struct {
		Results struct {
			ImageVersion string `json:"image_version"`
		}
	}
	err := c.GetAPIContext(ctx, &data, form["action"], form, BasicCheck)
	if err != nil {
		return "", err
	}
	return data.Results.ImageVersion, nil
}
