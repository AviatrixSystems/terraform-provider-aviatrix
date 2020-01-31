package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
)

// Controller Http Access enabled get result struct
type ControllerHttpAccessResp struct {
	Return bool   `json:"return"`
	Result string `json:"results"`
	Reason string `json:"reason"`
}

type GetSecurityGroupManagementResp struct {
	Return  bool              `json:"return"`
	Results SecurityGroupInfo `json:"results"`
	Reason  string            `json:"reason"`
}

type SecurityGroupInfo struct {
	State       string `json:"state"`
	AccountName string `json:"account_name"`
	Response    string `json:"response"`
}

type CloudnBackupConfig struct {
	CID                 string `form:"CID,omitempty"`
	Action              string `form:"action,omitempty"`
	BackupConfigEnabled string `json:"enabled,omitempty"`
	AccountName         string `json:"acct_name,omitempty"`
	CloudType           int    `json:"cloud_type,omitempty"`
	BucketName          string `json:"bucket_name,omitempty"`
	Multiple            string `json:"multiple_bkup,omitempty"`
}

type GetCloudnBackupConfigResp struct {
	Return  bool               `json:"return"`
	Results CloudnBackupConfig `json:"results"`
	Reason  string             `json:"reason"`
}

func (c *Client) EnableHttpAccess() error {
	url := "?CID=%s&action=config_http_access&operation=enable"
	path := c.baseURL + fmt.Sprintf(url, c.CID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return err
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
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
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
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
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
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
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode detach_fqdn_filter_tag_from_gw failed: " + err.Error() + "\n Body: " + bodyString)
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
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode disable_fqdn_exception_rule failed: " + err.Error() + "\n Body: " + bodyString)
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
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return false, errors.New("Json Decode get_fqdn_exception_rule_status failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return false, errors.New("Rest API get_fqdn_exception_rule_status Get failed: " + data.Reason)
	}
	if data.Results == "disabled" {
		return false, nil
	}
	return true, nil
}

func (c *Client) EnableSecurityGroupManagement(account string) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New("url Parsing failed for enable_controller_security_group_management " + err.Error())
	}
	enableSecurityGroupManagement := url.Values{}
	enableSecurityGroupManagement.Add("CID", c.CID)
	enableSecurityGroupManagement.Add("action", "enable_controller_security_group_management")
	enableSecurityGroupManagement.Add("access_account_name", account)
	Url.RawQuery = enableSecurityGroupManagement.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get enable_controller_security_group_management failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode enable_controller_security_group_management failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API enable_controller_security_group_management Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DisableSecurityGroupManagement() error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New("url Parsing failed for disable_controller_security_group_management " + err.Error())
	}
	disableSecurityGroupManagement := url.Values{}
	disableSecurityGroupManagement.Add("CID", c.CID)
	disableSecurityGroupManagement.Add("action", "disable_controller_security_group_management")
	Url.RawQuery = disableSecurityGroupManagement.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get disable_controller_security_group_management failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode disable_controller_security_group_management failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API disable_controller_security_group_management Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetSecurityGroupManagementStatus() (*SecurityGroupInfo, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New("url Parsing failed for get_controller_security_group_management_status " + err.Error())
	}
	getSecurityGroupManagementStatus := url.Values{}
	getSecurityGroupManagementStatus.Add("CID", c.CID)
	getSecurityGroupManagementStatus.Add("action", "get_controller_security_group_management_status")
	Url.RawQuery = getSecurityGroupManagementStatus.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get get_controller_security_group_management_status failed: " + err.Error())
	}
	var data GetSecurityGroupManagementResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode get_controller_security_group_management_status failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return nil, errors.New("Rest API get_controller_security_group_management_status Get failed: " + data.Reason)
	}
	return &data.Results, nil
}

func (c *Client) EnableCloudnBackupConfig(cloudnBackupConfig *CloudnBackupConfig) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'enable_cloudn_backup_config': %s") + err.Error())
	}
	enableCloudnBackupConfig := url.Values{}
	enableCloudnBackupConfig.Add("CID", c.CID)
	enableCloudnBackupConfig.Add("action", "enable_cloudn_backup_config")
	enableCloudnBackupConfig.Add("cloud_type", strconv.Itoa(cloudnBackupConfig.CloudType))
	enableCloudnBackupConfig.Add("account_name", cloudnBackupConfig.AccountName)
	enableCloudnBackupConfig.Add("bucket_name", cloudnBackupConfig.BucketName)
	if cloudnBackupConfig.Multiple == "true" {
		enableCloudnBackupConfig.Add("multiple", "true")
	}
	Url.RawQuery = enableCloudnBackupConfig.Encode()
	resp, err := c.Get(Url.String(), nil)
	log.Printf("[INFO] Enabling cloudn backup config: %#v", cloudnBackupConfig)
	if err != nil {
		return errors.New("HTTP Get 'enable_cloudn_backup_config' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'enable_cloudn_backup_config' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'enable_cloudn_backup_config' Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DisableCloudnBackupConfig() error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'disable_cloudn_backup_config': %s") + err.Error())
	}
	enableCloudnBackupConfig := url.Values{}
	enableCloudnBackupConfig.Add("CID", c.CID)
	enableCloudnBackupConfig.Add("action", "disable_cloudn_backup_config")
	Url.RawQuery = enableCloudnBackupConfig.Encode()
	resp, err := c.Get(Url.String(), nil)
	log.Printf("[INFO] Disabling cloudn backup config")
	if err != nil {
		return errors.New("HTTP Get 'disable_cloudn_backup_config' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'disable_cloudn_backup_config' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'disable_cloudn_backup_config' Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetCloudnBackupConfig() (*CloudnBackupConfig, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for 'get_cloudn_backup_config': %s") + err.Error())
	}
	enableCloudnBackupConfig := url.Values{}
	enableCloudnBackupConfig.Add("CID", c.CID)
	enableCloudnBackupConfig.Add("action", "get_cloudn_backup_config")
	Url.RawQuery = enableCloudnBackupConfig.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get 'get_cloudn_backup_config' failed: " + err.Error())
	}
	var data GetCloudnBackupConfigResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode 'get_cloudn_backup_config' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return nil, errors.New("Rest API 'get_cloudn_backup_config' Get failed: " + data.Reason)
	}
	return &data.Results, nil
}
