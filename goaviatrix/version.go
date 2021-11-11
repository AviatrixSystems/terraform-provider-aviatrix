package goaviatrix

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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
	}
	resp, err := c.Post(c.baseURL, form)
	if err != nil {
		return fmt.Errorf("HTTP POST %s failed: %v", form["action"], err)
	}
	var data struct {
		Return bool `json:"return"`
		Result int  `json:"results"`
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode " + form["action"] + " failed: " + err.Error() + "\n Body: " + bodyString)
	}

	requestID := data.Result
	if requestID == 0 {
		// Could not decode as integer, so something went wrong
		return fmt.Errorf("rest API %s POST failed to initiate async action", form["action"])
	}

	form = map[string]string{
		"action": "check_upgrade_status",
		"CID":    c.CID,
		"id":     strconv.Itoa(requestID),
		"pos":    "0",
	}
	backendURL := fmt.Sprintf("https://%s/v1/backend1", c.ControllerIP)
	const maxPoll = 180
	sleepDuration := time.Second * 10
	var i int
	for ; i < maxPoll; i++ {
		resp, err = c.Post(backendURL, form)
		if err != nil {
			// Could be transient HTTP error, e.g. EOF error
			time.Sleep(sleepDuration)
			continue
		}
		buf = new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		var data struct {
			Pos    int    `json:"pos"`
			Done   bool   `json:"done"`
			Status bool   `json:"status"`
			Result string `json:"result"`
		}
		err = json.Unmarshal(buf.Bytes(), &data)
		if err != nil {
			return fmt.Errorf("decode check_upgrade_status failed: %v\n Body: %s", err, buf.String())
		}
		if !data.Done {
			// Not done yet
			time.Sleep(sleepDuration)
			continue
		}

		// Upgrade is done, check for error
		if strings.HasPrefix(data.Result, "Error") {
			return fmt.Errorf("post check_upgrade_status failed: %s", data.Result)
		}
		break
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
