package goaviatrix

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
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

type CloudnBackupConfiguration struct {
	CID                 string `form:"CID,omitempty"`
	Action              string `form:"action,omitempty"`
	BackupConfiguration string `json:"enabled,omitempty"`
	BackupAccountName   string `json:"acct_name,omitempty"`
	BackupCloudType     int    `json:"cloud_type,omitempty"`
	BackupBucketName    string `json:"bucket_name,omitempty"`
	BackupStorageName   string `json:"storage_name"`
	BackupContainerName string `json:"container_name"`
	BackupRegion        string `json:"region"`
	MultipleBackups     string `json:"multiple_bkup,omitempty"`
}

type GetCloudnBackupConfigResp struct {
	Return  bool                      `json:"return"`
	Results CloudnBackupConfiguration `json:"results"`
	Reason  string                    `json:"reason"`
}

type CertDomainConfig struct {
	CertDomain string `json:"cert_domain"`
	IsDefault  bool   `json:"is_default"`
}

type ResourceCounts struct {
	Name  string `json:"Name"`
	Count int    `json:"Count"`
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
		log.Errorf("Error invoking controller %s", data.Reason)
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
		log.Errorf("Error invoking controller %s", data.Reason)
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
		log.Errorf("Error invoking controller %s", data.Reason)
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

func (c *Client) EnableCloudnBackupConfig(cloudnBackupConfiguration *CloudnBackupConfiguration) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'enable_cloudn_backup_config': %s") + err.Error())
	}
	enableCloudnBackupConfig := url.Values{}
	enableCloudnBackupConfig.Add("CID", c.CID)
	enableCloudnBackupConfig.Add("action", "enable_cloudn_backup_config")
	enableCloudnBackupConfig.Add("cloud_type", strconv.Itoa(cloudnBackupConfiguration.BackupCloudType))
	enableCloudnBackupConfig.Add("account_name", cloudnBackupConfiguration.BackupAccountName)
	enableCloudnBackupConfig.Add("bucket_name", cloudnBackupConfiguration.BackupBucketName)
	enableCloudnBackupConfig.Add("storage_name", cloudnBackupConfiguration.BackupStorageName)
	enableCloudnBackupConfig.Add("container_name", cloudnBackupConfiguration.BackupContainerName)
	enableCloudnBackupConfig.Add("region", cloudnBackupConfiguration.BackupRegion)
	if cloudnBackupConfiguration.MultipleBackups == "true" {
		enableCloudnBackupConfig.Add("multiple", "true")
	}
	Url.RawQuery = enableCloudnBackupConfig.Encode()
	resp, err := c.Get(Url.String(), nil)
	log.Infof("Enabling cloudn backup config: %#v", cloudnBackupConfiguration)
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
	log.Infof("Disabling cloudn backup config")
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

func (c *Client) GetCloudnBackupConfig() (*CloudnBackupConfiguration, error) {
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

func (c *Client) GetControllerVpcDnsServerStatus() (bool, error) {
	action := "get_controller_vpc_dns_server_status"
	resp, err := c.Post(c.baseURL, &APIRequest{
		CID:    c.CID,
		Action: action,
	})
	if err != nil {
		return false, fmt.Errorf("HTTP POST %q failed: %v", action, err)
	}

	type Resp struct {
		Return  bool   `json:"return"`
		Results string `json:"results"`
		Reason  string `json:"reason"`
	}
	var data Resp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return false, fmt.Errorf("reading response body %q failed: %v", action, err)
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return false, fmt.Errorf("json decode %q failed: %v\n Body: %s", action, err, b.String())
	}
	if !data.Return {
		return false, fmt.Errorf("rest api %q post failed: %s", action, data.Reason)
	}

	return data.Results == "Enabled", nil
}

func (c *Client) SetControllerVpcDnsServer(enabled bool) error {
	action := "enable_controller_vpc_dns_server"
	if !enabled {
		action = "disable_controller_vpc_dns_server"
	}
	return c.PostAPI(action, &APIRequest{
		CID:    c.CID,
		Action: action,
	}, BasicCheck)
}

func (c *Client) SetEmailExceptionNotification(ctx context.Context, enabled bool) error {
	action := "enable_exception_email_notification"
	if !enabled {
		action = "disable_exception_email_notification"
	}
	return c.PostAPIContext(ctx, action, &APIRequest{
		CID:    c.CID,
		Action: action,
	}, BasicCheck)
}

func (c *Client) GetEmailExceptionNotificationStatus(ctx context.Context) (bool, error) {
	params := map[string]string{
		"action": "get_exception_email_notification_status",
		"CID":    c.CID,
	}

	type Resp struct {
		Return  bool            `json:"return"`
		Results map[string]bool `json:"results"`
		Reason  string          `json:"reason"`
	}

	var data Resp

	err := c.GetAPIContext(ctx, &data, params["action"], params, BasicCheck)
	if err != nil {
		return true, err
	}

	if ans, ok := data.Results["enabled"]; ok {
		return ans, nil
	} else {
		return true, fmt.Errorf("response doesn't contain the key \"enabled\"")
	}
}

func (c *Client) SetCertDomain(ctx context.Context, certDomain string) error {
	params := map[string]string{
		"action":      "set_cert_domain",
		"CID":         c.CID,
		"cert_domain": certDomain,
	}
	return c.PostAPIContext(ctx, params["action"], params, BasicCheck)
}

func (c *Client) GetCertDomain(ctx context.Context) (*CertDomainConfig, error) {
	params := map[string]string{
		"action": "list_cert_domain",
		"CID":    c.CID,
	}

	type Resp struct {
		Return  bool             `json:"return"`
		Results CertDomainConfig `json:"results"`
		Reason  string           `json:"reason"`
	}

	var data Resp

	err := c.GetAPIContext(ctx, &data, params["action"], params, BasicCheck)
	if err != nil {
		return nil, err
	}

	return &data.Results, nil
}

func (c *Client) GetGatewayCount(ctx context.Context) (int, error) {
	params := map[string]string{
		"action": "list_resource_counts",
		"CID":    c.CID,
	}

	type Resp struct {
		Return  bool             `json:"return"`
		Results []ResourceCounts `json:"results"`
	}

	var data Resp
	err := c.GetAPIContext(ctx, &data, params["action"], params, BasicCheck)
	if err != nil {
		return -1, err
	}

	var gatewayCount int
	for _, resourceCount := range data.Results {
		if strings.Contains(resourceCount.Name, "Gateways") {
			gatewayCount += resourceCount.Count
		}
	}

	return gatewayCount, nil
}

func (c *Client) GetSleepTime(ctx context.Context) (time.Duration, error) {
	gatewayCount, err := c.GetGatewayCount(ctx)
	if err != nil {
		return -1, fmt.Errorf("could not get gateway count: %v", err)
	}
	return time.Duration(20 * int(math.Ceil(float64(gatewayCount)/15.0))), nil
}
