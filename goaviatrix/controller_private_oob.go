package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type PrivateOobResp struct {
	Return  bool   `json:"return"`
	Results string `json:"results"`
	Reason  string `json:"reason"`
}

func (c *Client) EnablePrivateOob() error {
	data := map[string]string{
		"action": "enable_private_oob",
		"CID":    c.CID,
	}
	checkFunc := func(action, reason string, ret bool) error {
		if !ret && !strings.HasPrefix(reason, "enable already") {
			return fmt.Errorf("rest API %s Post failed: %s", action, reason)
		}
		return nil
	}
	return c.PostAPI(data["action"], data, checkFunc)
}

func (c *Client) GetPrivateOobState() (bool, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return false, errors.New(("url Parsing failed for get_private_oob_state") + err.Error())
	}
	getPrivateOobState := url.Values{}
	getPrivateOobState.Add("CID", c.CID)
	getPrivateOobState.Add("action", "get_private_oob_state")
	Url.RawQuery = getPrivateOobState.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return false, errors.New("HTTP Get get_private_oob_state failed: " + err.Error())
	}
	var data PrivateOobResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return false, errors.New("Json Decode get_private_oob_state failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return false, errors.New("Rest API get_private_oob_state Get failed: " + data.Reason)
	}
	if data.Results == "Enabled" || data.Results == "enabled" {
		return true, nil
	} else if data.Results == "Disabled" || data.Results == "disabled" {
		return false, nil
	}
	return false, ErrNotFound
}

func (c *Client) DisablePrivateOob() error {
	data := map[string]string{
		"action": "disable_private_oob",
		"CID":    c.CID,
	}
	checkFunc := func(action, reason string, ret bool) error {
		if !ret && !strings.HasPrefix(reason, "disable already") {
			return fmt.Errorf("rest API %s Post failed: %s", action, reason)
		}
		return nil
	}
	return c.PostAPI(data["action"], data, checkFunc)
}
