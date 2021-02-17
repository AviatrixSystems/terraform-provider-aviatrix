package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
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
	upgradeController := url.Values{}
	upgradeController.Add("CID", c.CID)
	upgradeController.Add("action", "upgrade")

	if version.Version == "" {
		return errors.New("no target version is set")
	} else if version.Version != "latest" {
		upgradeController.Add("version", version.Version)
	}
	for i := 0; ; i++ {
		Url.RawQuery = upgradeController.Encode()
		resp, err := c.Get(Url.String(), nil)
		if err != nil {
			return errors.New("HTTP Get upgrade failed: " + err.Error())
		}
		var data UpgradeResp
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		bodyString := buf.String()
		bodyIoCopy := strings.NewReader(bodyString)
		if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
			return errors.New("Json Decode upgrade failed: " + err.Error() + "\n Body: " + bodyString)
		}
		if !data.Return {
			if strings.Contains(data.Reason, "Active upgrade in progress.") && i < 3 {
				log.Infof("Active upgrade is in progress. Retry after 60 secs...")
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
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return "", nil, errors.New("Json Decode list_version_info failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return "", nil, errors.New("Rest API list_version_info Get failed: " + data.Reason)
	}

	curVersion, aVer, err := ParseVersion(data.Results.CurrentVersion)
	if err != nil {
		return "", aVer, err
	}

	return curVersion, aVer, nil
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
			log.Tracef("response %s", body)
			if strings.Contains(string(body), "in progress") && i < 3 {
				log.Infof("Active upgrade is in progress. Retry after 60 secs...")
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
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return "", errors.New("Json Decode list_version_info failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return "", errors.New("Rest API list_version_info Get failed: " + data.Reason)
	}

	latestVersion, _, err := ParseVersion(data.Results.LatestVersion)
	if err != nil {
		return "", err
	}
	return latestVersion, nil
}

func ParseVersion(version string) (string, *AviatrixVersion, error) {
	if strings.HasPrefix(version, "UserConnect-") {
		version = version[12:]
	}
	if version == "" {
		return "", nil, errors.New("unable to parse version information since it is empty")
	}

	parts := strings.Split(version, ".")
	aver := &AviatrixVersion{}
	var err1, err2, err3 error
	aver.Major, err1 = strconv.ParseInt(parts[0], 10, 0)
	if len(parts) >= 2 {
		if strings.Contains(parts[1], "-") {
			aver.Minor, err2 = strconv.ParseInt(strings.Split(parts[1], "-")[0], 10, 0)
		} else {
			aver.Minor, err2 = strconv.ParseInt(parts[1], 10, 0)
		}
	} else {
		return "", aver, errors.New("unable to get latest version when parsing version information")
	}
	if len(parts) >= 3 {
		aver.Build, err3 = strconv.ParseInt(parts[2], 10, 0)
	}
	if err1 != nil || err2 != nil || err3 != nil {
		return "", aver, errors.New("unable to get latest version when parsing version information")
	}
	return strconv.FormatInt(aver.Major, 10) + "." + strconv.FormatInt(aver.Minor, 10), aver, nil
}
