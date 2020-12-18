package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
)

type PrivateOobResp struct {
	Return  bool   `json:"return"`
	Results string `json:"results"`
	Reason  string `json:"reason"`
}

func (c *Client) EnablePrivateOob() error {
	resp, err := c.Post(c.baseURL, struct {
		CID    string `form:"CID"`
		Action string `form:"action"`
	}{
		CID:    c.CID,
		Action: "enable_private_oob",
	})
	if err != nil {
		return errors.New("HTTP POST enable_private_oob failed: " + err.Error())
	}

	var data PrivateOobResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode enable_private_oob failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		if strings.Contains(data.Results, "enable already") {
			return nil
		}
		return errors.New("Rest API enable_private_oob Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetPrivateOobState() (bool, error) {
	resp, err := c.Get(c.baseURL, struct {
		CID    string `form:"CID"`
		Action string `form:"action"`
	}{
		CID:    c.CID,
		Action: "get_private_oob_state",
	})
	if err != nil {
		return false, errors.New("HTTP POST get_private_oob_state failed: " + err.Error())
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
		return false, errors.New("Rest API get_private_oob_state Post failed: " + data.Reason)
	}
	if data.Results == "Enabled" || data.Results == "enabled" {
		return true, nil
	} else if data.Results == "Disabled" || data.Results == "disabled" {
		return false, nil
	}
	return false, ErrNotFound
}

func (c *Client) DisablePrivateOob() error {
	resp, err := c.Post(c.baseURL, struct {
		CID    string `form:"CID"`
		Action string `form:"action"`
	}{
		CID:    c.CID,
		Action: "disable_private_oob",
	})
	if err != nil {
		return errors.New("HTTP POST disable_private_oob failed: " + err.Error())
	}

	var data PrivateOobResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode disable_private_oob failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		if strings.Contains(data.Results, "disable already") {
			return nil
		}
		return errors.New("Rest API disable_private_oob Post failed: " + data.Reason)
	}
	return nil
}
