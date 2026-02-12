package goaviatrix

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
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
	BackupStorageName   string `json:"storage_name,omitempty"`
	BackupContainerName string `json:"container_name,omitempty"`
	BackupRegion        string `json:"region,omitempty"`
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
		"CID":        c.CID,
		"action":     "enable_cloudn_backup_config",
		"cloud_type": strconv.Itoa(cloudnBackupConfiguration.BackupCloudType),
		"acct_name":  cloudnBackupConfiguration.BackupAccountName,
		"region":     cloudnBackupConfiguration.BackupRegion,
	}
	// Azure has a set of different parameters that must be set.
	if IsCloudType(cloudnBackupConfiguration.BackupCloudType, AzureArmRelatedCloudTypes) {
		form["storage_name"] = cloudnBackupConfiguration.BackupStorageName
		form["container_name"] = cloudnBackupConfiguration.BackupContainerName
	} else {
		form["bucket_name"] = cloudnBackupConfiguration.BackupBucketName
	}

	if cloudnBackupConfiguration.MultipleBackups == "true" {
		form["multiple_bkup"] = "true"
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

	if result, ok := data.Results["enabled"]; ok {
		return result, nil
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
	return c.PostAsyncAPIContextSetCertDomain(ctx, params["action"], params, BasicCheck)
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
		return -1, fmt.Errorf("could not get gateway count: %w", err)
	}
	return time.Duration(20 * int(math.Ceil(float64(gatewayCount)/15.0))), nil
}

func (c *Client) PostAsyncAPIContextSetCertDomain(
	ctx context.Context,
	action string,
	i interface{},
	checkFunc CheckAPIResponseFunc,
) error {
	const maxBody = 256 << 10 // 256 KiB

	readAndClose := func(resp *http.Response) ([]byte, error) {
		if resp == nil || resp.Body == nil {
			return nil, fmt.Errorf("nil http response/body")
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(io.LimitReader(resp.Body, maxBody+1))
		if err != nil {
			return nil, err
		}
		if len(b) > maxBody {
			b = b[:maxBody]
		}
		return b, nil
	}

	// Kickoff request
	r, err := c.PostContext(ctx, c.baseURL, i)
	if err != nil {
		// Close response from this failed attempt, if any.
		if r != nil && r.Body != nil {
			_ = r.Body.Close()
		}
		return fmt.Errorf("HTTP POST %s failed: %w", action, err)
	}
	resp := r

	body, err := readAndClose(resp)
	if err != nil {
		return fmt.Errorf("read response body for %s failed: %w", action, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP POST %s returned %s\nBody: %s", action, resp.Status, string(body))
	}

	var kickoff struct {
		Return bool   `json:"return"`
		Result string `json:"results"`
		Reason string `json:"reason"`
	}
	if err := json.Unmarshal(body, &kickoff); err != nil {
		return fmt.Errorf("json Decode %s failed: %w\nBody: %s", action, err, string(body))
	}
	if !kickoff.Return {
		return fmt.Errorf("rest API %s POST failed to initiate async action: %s", action, kickoff.Reason)
	}

	requestID := kickoff.Result
	form := map[string]string{
		"action":     "check_task_status",
		"CID":        c.CID,
		"request_id": requestID,
	}

	const maxPoll = 360
	sleepDuration := 10 * time.Second

	for j := 0; j < maxPoll; j++ {
		r, err := c.PostContext(ctx, c.baseURL, form)
		if err != nil {
			// Could be transient HTTP error, e.g. EOF.
			// Also avoid leaking response from this failed attempt.
			if r != nil && r.Body != nil {
				_ = r.Body.Close()
			}

			select {
			case <-time.After(sleepDuration):
				continue
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		resp = r

		body, err = readAndClose(resp)
		if err != nil {
			return fmt.Errorf("read check_task_status body failed: %w", err)
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("check_task_status returned %s\nBody: %s", resp.Status, string(body))
		}

		var status struct {
			Return bool   `json:"return"`
			Reason string `json:"reason"`
		}
		if err := json.Unmarshal(body, &status); err != nil {
			return fmt.Errorf("decode check_task_status failed: %w\nBody: %s", err, string(body))
		}

		if !status.Return {
			if status.Reason != "REQUEST_IN_PROGRESS" {
				return fmt.Errorf("rest API %s POST failed: %s", action, status.Reason)
			}

			// Not done yet
			select {
			case <-time.After(sleepDuration):
				continue
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		return checkFunc(action, "Post", "", status.Return)
	}

	return fmt.Errorf(
		"waited %s but upgrade never finished. Please manually verify the upgrade status",
		time.Duration(maxPoll)*sleepDuration,
	)
}
