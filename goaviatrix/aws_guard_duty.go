package goaviatrix

import (
	"fmt"
	"strconv"
	"strings"
)

type AwsGuardDuty struct {
	ScanningInterval int
	EnabledAccounts  []*AwsGuardDutyAccount
}

type AwsGuardDutyAccount struct {
	AccountName string `json:"account"`
	Region      string
	ExcludedIPs []string `json:"exempt_ips"`
}

func (acc *AwsGuardDutyAccount) ID() string {
	return acc.AccountName + "~~" + acc.Region
}

func (c *Client) UpdateAwsGuardDutyPollInterval(scanningInterval int) error {
	data := map[string]string{
		"action":   "update_aws_guard_duty_poll_interval",
		"CID":      c.CID,
		"interval": strconv.Itoa(scanningInterval),
	}
	checkFunc := func(action, method, reason string, ret bool) error {
		// AVXERR-SECURITY-0093 is returned if you try to update the interval to its currently configured value.
		if !ret && !strings.HasPrefix(reason, "[AVXERR-SECURITY-0093]") {
			return fmt.Errorf("rest API %s %s failed: %s", action, method, reason)
		}
		return nil
	}
	return c.PostAPI(data["action"], data, checkFunc)
}

func (c *Client) EnableAwsGuardDuty(account *AwsGuardDutyAccount) error {
	data := map[string]string{
		"action":       "enable_aws_guard_duty",
		"CID":          c.CID,
		"account_name": account.AccountName,
		"region":       account.Region,
	}
	checkFunc := func(action, method, reason string, ret bool) error {
		if !ret && !strings.HasPrefix(reason, "[AVXERR-SECURITY-0089] GuardDuty is already enabled for the account") {
			return fmt.Errorf("rest API %s %s failed: %s", action, method, reason)
		}
		return nil
	}
	return c.PostAPI(data["action"], data, checkFunc)
}

func (c *Client) DisableAwsGuardDuty(account *AwsGuardDutyAccount) error {
	data := map[string]string{
		"action":       "disable_aws_guard_duty",
		"CID":          c.CID,
		"account_name": account.AccountName,
		"region":       account.Region,
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) UpdateAwsGuardDutyExcludedIPs(account *AwsGuardDutyAccount) error {
	data := map[string]string{
		"action":       "update_aws_guard_duty_excluded_ips",
		"CID":          c.CID,
		"account_name": account.AccountName,
		"region":       account.Region,
		"excluded_ips": strings.Join(account.ExcludedIPs, ","),
	}
	checkFunc := func(action, method, reason string, ret bool) error {
		if !ret && !strings.HasPrefix(reason, "No change in the exclude ip list for account") {
			return fmt.Errorf("rest API %s %s failed: %s", action, method, reason)
		}
		return nil
	}
	return c.PostAPI(data["action"], data, checkFunc)
}

type ListAwsGuardDutyResp struct {
	Return  bool
	Reason  string
	Results ListAwsGuardDutyResults
}
type ListAwsGuardDutyResults struct {
	IntervalInfo  ListAwsGuardDutyIntervalInfo `json:"interval_info"`
	GuardDutyList []*AwsGuardDutyAccount       `json:"guard_duty"`
}
type ListAwsGuardDutyIntervalInfo struct {
	Interval int
	Options  []int
}

func (c *Client) GetAwsGuardDuty() (*AwsGuardDuty, error) {
	formData := map[string]string{
		"action": "list_aws_guard_duty",
		"CID":    c.CID,
	}
	var data ListAwsGuardDutyResp
	err := c.GetAPI(&data, formData["action"], formData, BasicCheck)
	if err != nil {
		return nil, err
	}
	return &AwsGuardDuty{
		ScanningInterval: data.Results.IntervalInfo.Interval,
		EnabledAccounts:  data.Results.GuardDutyList,
	}, nil
}

func (c *Client) GetAwsGuardDutyAccount(accountName, region string) (*AwsGuardDutyAccount, error) {
	guardDuty, err := c.GetAwsGuardDuty()
	if err != nil {
		return nil, err
	}
	var acc *AwsGuardDutyAccount
	for _, v := range guardDuty.EnabledAccounts {
		if v.AccountName != accountName || strings.Fields(v.Region)[0] != region {
			continue
		}
		v.Region = strings.Fields(v.Region)[0]
		acc = v
	}
	if acc == nil {
		return nil, ErrNotFound
	}
	return acc, nil
}
