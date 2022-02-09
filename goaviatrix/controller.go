package goaviatrix

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
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
	form := map[string]string{
		"CID":       c.CID,
		"action":    "config_http_access",
		"operation": "enable",
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) DisableHttpAccess() error {
	form := map[string]string{
		"CID":       c.CID,
		"action":    "config_http_access",
		"operation": "disable",
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) GetHttpAccessEnabled() (string, error) {
	var data ControllerHttpAccessResp
	form := map[string]string{
		"CID":       c.CID,
		"action":    "config_http_access",
		"operation": "get",
	}
	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return "", err
	}
	return data.Result, nil
}

func (c *Client) EnableExceptionRule() error {
	form := map[string]string{
		"CID":    c.CID,
		"action": "enable_fqdn_exception_rule",
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) DisableExceptionRule() error {
	form := map[string]string{
		"CID":    c.CID,
		"action": "disable_fqdn_exception_rule",
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) GetExceptionRuleStatus() (bool, error) {
	var data GetFqdnExceptionRuleResp
	form := map[string]string{
		"CID":    c.CID,
		"action": "get_fqdn_exception_rule_status",
	}
	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return false, err
	}
	if data.Results == "disabled" {
		return false, nil
	}
	return true, nil
}

func (c *Client) EnableSecurityGroupManagement(account string) error {
	form := map[string]string{
		"CID":                 c.CID,
		"action":              "enable_controller_security_group_management",
		"access_account_name": account,
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) DisableSecurityGroupManagement() error {
	form := map[string]string{
		"CID":    c.CID,
		"action": "disable_controller_security_group_management",
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) GetSecurityGroupManagementStatus() (*SecurityGroupInfo, error) {
	var data GetSecurityGroupManagementResp
	form := map[string]string{
		"CID":    c.CID,
		"action": "get_controller_security_group_management_status",
	}
	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}
	return &data.Results, nil
}

func (c *Client) EnableCloudnBackupConfig(cloudnBackupConfiguration *CloudnBackupConfiguration) error {
	form := map[string]string{
		"CID":            c.CID,
		"action":         "enable_cloudn_backup_config",
		"cloud_type":     strconv.Itoa(cloudnBackupConfiguration.BackupCloudType),
		"account_name":   cloudnBackupConfiguration.BackupAccountName,
		"bucket_name":    cloudnBackupConfiguration.BackupBucketName,
		"storage_name":   cloudnBackupConfiguration.BackupStorageName,
		"container_name": cloudnBackupConfiguration.BackupContainerName,
		"region":         cloudnBackupConfiguration.BackupRegion,
	}
	if cloudnBackupConfiguration.MultipleBackups == "true" {
		form["multiple"] = "true"
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) DisableCloudnBackupConfig() error {
	form := map[string]string{
		"CID":    c.CID,
		"action": "disable_cloudn_backup_config",
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) GetCloudnBackupConfig() (*CloudnBackupConfiguration, error) {
	var data GetCloudnBackupConfigResp
	form := map[string]string{
		"CID":    c.CID,
		"action": "get_cloudn_backup_config",
	}
	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}
	return &data.Results, nil
}

func (c *Client) GetControllerVpcDnsServerStatus() (bool, error) {
	type Resp struct {
		Return  bool   `json:"return"`
		Results string `json:"results"`
		Reason  string `json:"reason"`
	}
	var data Resp
	form := map[string]string{
		"CID":    c.CID,
		"action": "get_controller_vpc_dns_server_status",
	}
	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return false, err
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
		"async":       "true",
	}
	return c.PostAsyncAPIContext(ctx, params["action"], params, BasicCheck)
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
