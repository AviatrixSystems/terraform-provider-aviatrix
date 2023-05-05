package goaviatrix

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

type VersionInfoResults struct {
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
}

type VersionInfoResp struct {
	Return  bool               `json:"return"`
	Results VersionInfoResults `json:"results"`
	Reason  string             `json:"reason"`
}

type AviatrixVersion struct {
	Major        int64
	Minor        int64
	Build        int64
	MinorBuildID string
	HasBuild     bool // HasBuild indicates if the version string originally included a build number.
}

type VersionInfo struct {
	Current  *AviatrixVersion
	Previous *AviatrixVersion
}

type PlatformUpgradeResp struct {
	Return  bool                  `json:"return"`
	Results PlatformUpgradeStatus `json:"results"`
	Reason  string                `json:"reason"`
}

type PlatformUpgradeStatus struct {
	PlatformUpgrade     PlatformUpgradeInfo `json:"platform_upgrade,omitempty"`
	IsUpgradeInProgress bool                `json:"is_upgrade_in_progress,omitempty"`
}

type PlatformUpgradeInfo struct {
	Progmsg string                `json:"progmsg,omitempty"`
	Results PlatformUpgradeResult `json:"results,omitempty"`
	Reason  string                `json:"reason,omitempty"`
}

type PlatformUpgradeResult struct {
	Msg    string `json:"msg,omitempty"`
	Reason string `json:"reason,omitempty"`
	Status string `json:"status,omitempty"`
}

func (av *AviatrixVersion) String(includeBuild bool) string {
	version := fmt.Sprintf("%d.%d", av.Major, av.Minor)
	if av.MinorBuildID != "" {
		version += "-" + av.MinorBuildID
	}
	if includeBuild {
		version += "." + strconv.Itoa(int(av.Build))
	}
	return version
}

// AsyncUpgrade will upgrade controller asynchronously
func (c *Client) AsyncUpgrade(version *Version, upgradeGateways bool) error {
	form := map[string]string{
		"CID":   c.CID,
		"async": "true", // indicates an async command
	}
	if upgradeGateways {
		form["action"] = "upgrade"
		if version.Version != "latest" {
			form["version"] = version.Version
		}
	} else {
		form["action"] = "upgrade_platform"
		form["gateway_list"] = ""
		form["software_version"] = version.Version
		form["version"] = version.Version
	}
	resp, err := c.Post(c.baseURL, form)
	if err != nil {
		return fmt.Errorf("HTTP POST %s failed: %v", form["action"], err)
	}
	var data struct {
		Return bool   `json:"return"`
		Result string `json:"results"`
		Reason string `json:"reason"`
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode " + form["action"] + " failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return fmt.Errorf("rest API %s POST failed to initiate async action: %s", form["action"], data.Reason)
	}

	time.Sleep(time.Second * 90)

	form1 := map[string]interface{}{
		"action":        "platform_upgrade_status",
		"CID":           c.CID,
		"platform_only": true,
	}

	const maxPoll = 180
	sleepDuration := time.Second * 10
	var i int
	for ; i < maxPoll; i++ {
		resp, err = c.Post(c.baseURL, form1)
		if err != nil {
			// Could be transient HTTP error, e.g. EOF error
			time.Sleep(sleepDuration)
			continue
		}
		var data1 PlatformUpgradeResp
		buf = new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		err = json.Unmarshal(buf.Bytes(), &data1)
		if err != nil {
			if strings.Contains(buf.String(), "503 Service Unavailable") || strings.Contains(buf.String(), "502 Proxy Error") {
				time.Sleep(sleepDuration)
				continue
			}
			return fmt.Errorf("decode platform_upgrade_status failed: %v\n Body: %s", err, buf.String())
		}
		if !data1.Return {
			return fmt.Errorf("rest API %s POST failed to initiate async action: %s", form1["action"], data1.Reason)
		}
		if data1.Results.IsUpgradeInProgress {
			time.Sleep(sleepDuration)
			continue
		}
		if data1.Results.PlatformUpgrade.Results.Status == "complete" {
			c.Login()
			return nil
		} else if data1.Results.PlatformUpgrade.Results.Status == "in progress" {
			time.Sleep(sleepDuration)
			continue
		}
		return fmt.Errorf(data1.Results.PlatformUpgrade.Results.Reason)
	}
	// Waited for too long and upgrade never finished
	if i == maxPoll {
		return fmt.Errorf("waited %s but upgrade never finished. Please manually verify the upgrade status", maxPoll*sleepDuration)
	}
	c.Login()
	return nil
}

func (c *Client) UpgradeGateway(gateway *Gateway) error {
	form := map[string]string{
		"action":           "upgrade_selected_gateway",
		"CID":              c.CID,
		"gateway_list":     gateway.GwName,
		"software_version": gateway.SoftwareVersion,
		"image_version":    gateway.ImageVersion,
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) GetCurrentVersion() (string, *AviatrixVersion, error) {
	form := map[string]string{
		"CID":    c.CID,
		"action": "list_version_info",
	}

	var data VersionInfoResp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return "", nil, err
	}

	curVersion, aVer, err := ParseVersion(data.Results.CurrentVersion)
	if err != nil {
		return "", aVer, err
	}

	return curVersion, aVer, nil
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
	_, current, err := ParseVersion(data.Results.CurrentVersion)
	if err != nil {
		return nil, fmt.Errorf("could not parse current version: %v", err)
	}
	_, previous, err := ParseVersion(data.Results.PreviousVersion)
	if err != nil {
		return nil, fmt.Errorf("could not parse previous version: %v", err)
	}
	return &VersionInfo{
		Current:  current,
		Previous: previous,
	}, nil
}

func (c *Client) GetLatestVersion() (string, error) {
	form := map[string]string{
		"CID":            c.CID,
		"action":         "list_version_info",
		"latest_version": strconv.FormatBool(true),
	}

	var data VersionInfoResp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return "", err
	}

	latestVersion, _, err := ParseVersion(data.Results.LatestVersion)
	if err != nil {
		return "", err
	}
	return latestVersion, nil
}

func ParseVersion(version string) (string, *AviatrixVersion, error) {
	version = strings.TrimPrefix(version, "UserConnect-")
	if version == "" {
		log.Infof("version is an empty string")
		return "", &AviatrixVersion{}, nil
	}

	parts := strings.Split(version, ".")
	aver := &AviatrixVersion{}
	var err1, err2, err3 error
	aver.Major, err1 = strconv.ParseInt(parts[0], 10, 0)
	if len(parts) >= 2 {
		if strings.Contains(parts[1], "-") {
			minorParts := strings.Split(parts[1], "-")
			aver.MinorBuildID = minorParts[1]
			aver.Minor, err2 = strconv.ParseInt(minorParts[0], 10, 0)
		} else {
			aver.Minor, err2 = strconv.ParseInt(parts[1], 10, 0)
		}
	} else {
		return "", aver, errors.New("unable to get latest version when parsing version information")
	}
	if len(parts) >= 3 {
		aver.Build, err3 = strconv.ParseInt(parts[2], 10, 0)
		aver.HasBuild = true
	}
	if err1 != nil || err2 != nil || err3 != nil {
		return "", aver, errors.New("unable to get latest version when parsing version information")
	}
	return strconv.FormatInt(aver.Major, 10) + "." + strconv.FormatInt(aver.Minor, 10), aver, nil
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

// CompareSoftwareVersions first return value will be
// less than 0 if a < b
// equal to 0  if a == b
// more than 0 if a > b
func CompareSoftwareVersions(a, b string) (int, error) {
	_, versionA, errA := ParseVersion(a)
	_, versionB, errB := ParseVersion(b)
	if errA != nil || errB != nil {
		return 0, fmt.Errorf("invalid software version")
	}
	// Versions are exactly equal
	if a == b {
		return 0, nil
	}
	// Major versions are different
	if versionA.Major-versionB.Major != 0 {
		return int(versionA.Major - versionB.Major), nil
	}
	// Minor versions are different
	if versionA.Minor-versionB.Minor != 0 {
		return int(versionA.Minor - versionB.Minor), nil
	}
	// MinorBuildIDs are different
	if versionA.MinorBuildID != versionB.MinorBuildID {
		// A does not have a minor build ID but B does
		if versionA.MinorBuildID == "" {
			return -1, nil
		}
		// B does not have a minor build ID but A does
		if versionB.MinorBuildID == "" {
			return 1, nil
		}
		// Otherwise, both A and B have minor build IDs but they are different.
		// e.g. 6.5-patch vs 6.5-cyruspatch
		// We will just consider them as equal.
	}
	if versionA.HasBuild && versionB.HasBuild {
		// Build versions are different
		if versionA.Build-versionB.Build != 0 {
			return int(versionA.Build - versionB.Build), nil
		}
	} else if versionA.HasBuild { // Having a build is always less than not having a build. e.g a=6.5.100 b=6.5, b is considered higher
		return -1, nil
	} else if versionB.HasBuild {
		return 1, nil
	}
	// Versions are the same
	return 0, nil
}
