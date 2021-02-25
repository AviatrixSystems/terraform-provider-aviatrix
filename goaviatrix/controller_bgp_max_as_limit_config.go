package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/types"
	"strings"
)

type ControllerBgpMaxAsLimitConfig struct {
	Action string `form:"action,omitempty"`
	CID string `form:"CID,omitempty"`
	MaxAsLimit string `json:"max_as_limit"`
}

type SetBgpMaxAsLimitResult struct {
	Result string `json:"result"`
}

type SetBgpMaxAsLimitResponse struct {
	Return bool   `json:"return"`
	Reason string `json:"reason"`
	Results SetBgpMaxAsLimitResult `json:"results"`
	ErrorCode int `json:"errorcode"`
}

func (c *Client) CreateControllerBgpMaxAsLimitConfig(controllerBgpMaxAsLimitConfig *ControllerBgpMaxAsLimitConfig) error {
	controllerBgpMaxAsLimitConfig.CID = c.CID
	controllerBgpMaxAsLimitConfig.Action = "set_bgp_max_as_limit"

	resp, err := c.Post(c.baseURL, controllerBgpMaxAsLimitConfig)
	if err != nil {
		return errors.New("HTTP Post set_bgp_max_as_limit failed: " + err.Error())
	}

	var data SetBgpMaxAsLimitResponse
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode set_bgp_max_as_limit failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if data.ErrorCode != 0 {
		return fmt.Errorf("rest API set_bgp_max_as_limit failed with error code %d: %s", data.ErrorCode, data.Reason)
	} else if !data.Return {
		return errors.New("Rest API connect_container Post failed: " + data.Reason)
	}

	return nil
}

func (c *Client) GetControllerBgpMaxAsLimitConfig() (*ControllerBgpMaxAsLimitConfig, error) {
	data := map[string]string{
		"action":       "show_bgp_max_as_limit",
		"CID":          c.CID,
	}

	var resp SetBgpMaxAsLimitResponse
	err := c.GetAPI(&resp, data["action"], data, BasicCheck)
	if err != nil {
		return nil, err
	}
	if resp.ErrorCode != 0 {
		return nil, fmt.Errorf("rest API set_bgp_max_as_limit failed with error code %d: %s", resp.ErrorCode, resp.Reason)
	} else if !resp.Return {
		return nil, errors.New("Rest API connect_container Post failed: " + resp.Reason)
	}

	return nil, nil //TODO
}

func (c *Client) UpdateControllerBgpMaxAsLimitConfig(controllerBgpMaxAsLimitConfig *ControllerBgpMaxAsLimitConfig) error {
	controllerBgpMaxAsLimitConfig.CID = c.CID
	controllerBgpMaxAsLimitConfig.Action = "set_bgp_max_as_limit"


}
