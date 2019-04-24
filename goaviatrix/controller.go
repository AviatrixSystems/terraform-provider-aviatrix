package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
)

// Controller Http Access enabled get result struct
type ControllerHttpAccessResp struct {
	Return bool   `json:"return"`
	Result string `json:"results"`
	Reason string `json:"reason"`
}

func (c *Client) EnableHttpAccess() error {
	url := "?CID=%s&action=config_http_access&operation=enable"
	path := c.baseURL + fmt.Sprintf(url, c.CID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		log.Printf("[ERROR] Error invoking controller %s", data.Reason)
		return errors.New(data.Reason)
	}
	return nil
}

func (c *Client) DisableHttpAccess() error {
	url := "?CID=%s&action=config_http_access&operation=disable"
	path := c.baseURL + fmt.Sprintf(url, c.CID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		log.Printf("[ERROR] Error invoking controller %s", data.Reason)
		return errors.New(data.Reason)
	}
	return nil
}

func (c *Client) GetHttpAccessEnabled() (string, error) {
	url := "?CID=%s&action=config_http_access&operation=get"
	path := c.baseURL + fmt.Sprintf(url, c.CID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return "", err
	}
	var data ControllerHttpAccessResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	if !data.Return {
		log.Printf("[ERROR] Error invoking controller %s", data.Reason)
		return "", errors.New(data.Reason)
	}
	result := data.Result
	return result, nil
}

func (c *Client) EnableExceptionRule() error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New("url Parsing failed for enable_fqdn_exception_rule " + err.Error())
	}
	enableFqdnExceptionRule := url.Values{}
	enableFqdnExceptionRule.Add("CID", c.CID)
	enableFqdnExceptionRule.Add("action", "enable_fqdn_exception_rule")
	Url.RawQuery = enableFqdnExceptionRule.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get enable_fqdn_exception_rule failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode detach_fqdn_filter_tag_from_gw failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API detach_fqdn_filter_tag_from_gw Get failed: " + data.Reason)
	}

	return nil
}

func (c *Client) DisableExceptionRule() error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New("url Parsing failed for disable_fqdn_exception_rule " + err.Error())
	}
	disableFqdnExceptionRule := url.Values{}
	disableFqdnExceptionRule.Add("CID", c.CID)
	disableFqdnExceptionRule.Add("action", "disable_fqdn_exception_rule")
	Url.RawQuery = disableFqdnExceptionRule.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get disable_fqdn_exception_rule failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode disable_fqdn_exception_rule failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API disable_fqdn_exception_rule Get failed: " + data.Reason)
	}

	return nil
}

func (c *Client) GetExceptionRuleStatus() (bool, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return false, errors.New("url Parsing failed for get_fqdn_exception_rule_status " + err.Error())
	}
	getFqdnExceptionRuleStatus := url.Values{}
	getFqdnExceptionRuleStatus.Add("CID", c.CID)
	getFqdnExceptionRuleStatus.Add("action", "get_fqdn_exception_rule_status")
	Url.RawQuery = getFqdnExceptionRuleStatus.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return false, errors.New("HTTP Get get_fqdn_exception_rule_status failed: " + err.Error())
	}
	data := GetFqdnExceptionRuleResp{
		Return:  false,
		Results: "",
		Reason:  "",
	}
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return false, errors.New("Json Decode get_fqdn_exception_rule_status failed: " + err.Error())
	}
	if !data.Return {
		return false, errors.New("Rest API get_fqdn_exception_rule_status Get failed: " + data.Reason)
	}

	if data.Results == "disabled" {
		return false, nil
	}
	return true, nil
}
