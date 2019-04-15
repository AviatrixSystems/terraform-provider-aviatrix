package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Version struct {
	CID           string `form:"CID,omitempty"`
	Action        string `form:"action,omitempty"`
	TargetVersion string `form:"version,omitempty"`
	Version       string `json:"version,omitempty"`
}

type UpgradeResp struct {
	Return  bool   `json:"return"`
	Results string `json:"results"`
	Reason  string `json:"reason"`
}

type VersionInfo struct {
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
}

type VersionInfoResp struct {
	Return  bool        `json:"return"`
	Results VersionInfo `json:"results"`
	Reason  string      `json:"reason"`
}

type AviatrixVersion struct {
	Major int64
	Minor int64
	Build int64
}

func (c *Client) Upgrade(version *Version) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for upgrade") + err.Error())
	}
	attachSpokeToTransitGw := url.Values{}
	attachSpokeToTransitGw.Add("CID", c.CID)
	attachSpokeToTransitGw.Add("action", "upgrade")

	if version.Version == "" {
		return errors.New("no target version is set")
	} else if version.Version != "latest" {
		attachSpokeToTransitGw.Add("version", version.Version)
	}
	for i := 0; ; i++ {
		Url.RawQuery = attachSpokeToTransitGw.Encode()
		resp, err := c.Get(Url.String(), nil)
		if err != nil {
			return errors.New("HTTP Get upgrade failed: " + err.Error())
		}
		var data UpgradeResp
		if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return errors.New("Json Decode upgrade failed: " + err.Error())
		}
		if !data.Return {
			if strings.Contains(data.Reason, "Active upgrade in progress.") && i < 3 {
				log.Printf("[INFO] Active upgrade is in progress. Retry after 60 secs...")
				time.Sleep(60 * time.Second)
				continue
			}
			return errors.New("Rest API upgrade Get failed: " + data.Reason)
		}
		break
	}
	return nil
}

func (c *Client) GetCurrentVersion() (string, *AviatrixVersion, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return "", nil, errors.New(("url Parsing failed for list_version_info") + err.Error())
	}
	listVersionInfo := url.Values{}
	listVersionInfo.Add("CID", c.CID)
	listVersionInfo.Add("action", "list_version_info")
	Url.RawQuery = listVersionInfo.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return "", nil, errors.New("HTTP Get list_version_info failed: " + err.Error())
	}
	var data VersionInfoResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", nil, errors.New("Json Decode list_version_info failed: " + err.Error())
	}

	if !data.Return {
		return "", nil, errors.New("Rest API list_version_info Get failed: " + data.Reason)
	}

	// strip off "UserConnect-"
	parts := strings.Split(data.Results.CurrentVersion[12:], ".")
	aver := &AviatrixVersion{}
	var err1, err2, err3 error
	aver.Major, err1 = strconv.ParseInt(parts[0], 10, 0)
	aver.Minor, err2 = strconv.ParseInt(parts[1], 10, 0)
	aver.Build, err3 = strconv.ParseInt(parts[2], 10, 0)
	if err1 != nil || err2 != nil || err3 != nil {
		log.Printf("[WARN] Unable to get current version: %s|%s|%s (when parsing '%s')", err1, err2, err3, data.Results.CurrentVersion[11:])
		return data.Results.CurrentVersion, nil, err
	}
	return data.Results.CurrentVersion, aver, nil
}

func (c *Client) Pre32Upgrade() error {
	privateBaseURL := strings.Replace(c.baseURL, "/v1/api", "/v1/backend1", 1)
	params := &Version{
		Action: "userconnect_release",
		CID:    c.CID,
	}
	path := privateBaseURL
	for i := 0; ; i++ {
		resp, err := c.Post(path, params)
		if err != nil {
			return errors.New("HTTP Post userconnect_release failed: " + err.Error())
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			body, _ := ioutil.ReadAll(resp.Body)
			log.Printf("[TRACE] response %s", body)
			if strings.Contains(string(body), "in progress") && i < 3 {
				log.Printf("[INFO] Active upgrade is in progress. Retry after 60 secs...")
				time.Sleep(60 * time.Second)
			} else {
				break
			}
		} else {
			return fmt.Errorf("status code %d", resp.StatusCode)
		}
	}
	return nil
}

func (c *Client) GetLatestVersion() (string, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return "", errors.New(("url Parsing failed for list_version_info") + err.Error())
	}
	listVersionInfo := url.Values{}
	listVersionInfo.Add("CID", c.CID)
	listVersionInfo.Add("action", "list_version_info")
	Url.RawQuery = listVersionInfo.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return "", errors.New("HTTP Get list_version_info failed: " + err.Error())
	}
	var data VersionInfoResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", errors.New("Json Decode list_version_info failed: " + err.Error())
	}

	if !data.Return {
		return "", errors.New("Rest API list_version_info Get failed: " + data.Reason)
	}

	if data.Results.CurrentVersion != "" {
		return data.Results.LatestVersion, nil
	}
	return "", nil
}
