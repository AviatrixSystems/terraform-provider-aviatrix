package aviatrix

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

const defaultAwsGuardDutyScanningInterval = 60
const awsGuardDutyID = "aviatrix_aws_guard_duty"

var validAwsGuardDutyScanningIntervals = []int{5, 10, 15, 30, 60}

func resourceAviatrixAwsGuardDuty() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixAwsGuardDutyCreate,
		Read:   resourceAviatrixAwsGuardDutyRead,
		Update: resourceAviatrixAwsGuardDutyUpdate,
		Delete: resourceAviatrixAwsGuardDutyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"enabled_accounts": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Accounts to enable AWS GuardDuty.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Account name",
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Region.",
						},
						"excluded_ips": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.IsIPAddress,
							},
							Optional:    true,
							Description: "Excluded IPs.",
						},
					},
				},
			},
			"scanning_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "Scanning Interval",
				Default:      defaultAwsGuardDutyScanningInterval,
				ValidateFunc: validation.IntInSlice(validAwsGuardDutyScanningIntervals),
			},
		},
	}
}

func marshalAwsGuardDutyInput(d *schema.ResourceData) *goaviatrix.AwsGuardDuty {
	var accounts []*goaviatrix.AwsGuardDutyAccount
	for _, a := range d.Get("enabled_accounts").(*schema.Set).List() {
		aMap := a.(map[string]interface{})
		gdAccount := &goaviatrix.AwsGuardDutyAccount{
			AccountName: aMap["account_name"].(string),
			Region:      aMap["region"].(string),
		}
		var excludedIPs []string
		for _, ip := range aMap["excluded_ips"].(*schema.Set).List() {
			excludedIPs = append(excludedIPs, ip.(string))
		}
		gdAccount.ExcludedIPs = excludedIPs
		accounts = append(accounts, gdAccount)
	}
	return &goaviatrix.AwsGuardDuty{
		ScanningInterval: d.Get("scanning_interval").(int),
		EnabledAccounts:  accounts,
	}
}

func resourceAviatrixAwsGuardDutyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	guardDuty := marshalAwsGuardDutyInput(d)
	err := client.UpdateAwsGuardDutyPollInterval(guardDuty)
	if err != nil {
		return fmt.Errorf("could not set scanning interval: %v", err)
	}
	for _, account := range guardDuty.EnabledAccounts {
		err := client.EnableAwsGuardDuty(account)
		if err != nil {
			return fmt.Errorf("could not enable AWS GuardDuty for enabled_account %#v: %v", account, err)
		}
		err = client.UpdateAwsGuardDutyExcludedIPs(account)
		if err != nil {
			return fmt.Errorf("could not set excluded IPs for enabled_account %#v: %v", account, err)
		}
	}
	// Only 1 GuardDuty resource can exist per controller, so, the ID is static.
	d.SetId(awsGuardDutyID)
	return nil
}

func resourceAviatrixAwsGuardDutyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	guardDuty, err := client.GetAwsGuardDuty()
	if err != nil {
		return fmt.Errorf("could not get AWS GuardDuty: %v", err)
	}
	d.Set("scanning_interval", guardDuty.ScanningInterval)
	var enabledAccounts []map[string]interface{}
	for _, v := range guardDuty.EnabledAccounts {
		account := map[string]interface{}{
			"account_name": v.AccountName,
			"region":       strings.Fields(v.Region)[0],
			"excluded_ips": v.ExcludedIPs,
		}
		enabledAccounts = append(enabledAccounts, account)
	}
	if err := d.Set("enabled_accounts", enabledAccounts); err != nil {
		return fmt.Errorf("could not set enabled_accounts into state: %v", err)
	}
	d.SetId(awsGuardDutyID)
	return nil
}

func resourceAviatrixAwsGuardDutyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	guardDuty := marshalAwsGuardDutyInput(d)
	if d.HasChange("scanning_interval") {
		err := client.UpdateAwsGuardDutyPollInterval(guardDuty)
		if err != nil {
			return fmt.Errorf("could not update scanning interval: %v", err)
		}
	}
	if d.HasChange("enabled_accounts") {
		o, n := d.GetChange("enabled_accounts")
		removedAccounts := o.(*schema.Set).Difference(n.(*schema.Set))
		addedAccounts := n.(*schema.Set).Difference(o.(*schema.Set))
		for _, a := range removedAccounts.List() {
			aM := a.(map[string]interface{})
			account := &goaviatrix.AwsGuardDutyAccount{
				AccountName: aM["account_name"].(string),
				Region:      aM["region"].(string),
			}
			err := client.DisableAwsGuardDuty(account)
			if err != nil {
				return fmt.Errorf("could not disable GuardDuty for account %#v, %v", account, err)
			}
		}
		for _, a := range addedAccounts.List() {
			aM := a.(map[string]interface{})
			account := &goaviatrix.AwsGuardDutyAccount{
				AccountName: aM["account_name"].(string),
				Region:      aM["region"].(string),
			}
			var excludedIPs []string
			for _, v := range aM["excluded_ips"].(*schema.Set).List() {
				excludedIPs = append(excludedIPs, v.(string))
			}
			account.ExcludedIPs = excludedIPs
			err := client.EnableAwsGuardDuty(account)
			if err != nil {
				return fmt.Errorf("could not enable GuardDuty for account %#v: %v", account, err)
			}
			err = client.UpdateAwsGuardDutyExcludedIPs(account)
			if err != nil {
				return fmt.Errorf("could not edit GuardDuty excluded IPs for account %#v: %v", account, err)
			}
		}
	}
	return nil
}

func resourceAviatrixAwsGuardDutyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	guardDuty := marshalAwsGuardDutyInput(d)
	guardDuty.ScanningInterval = defaultAwsGuardDutyScanningInterval
	err := client.UpdateAwsGuardDutyPollInterval(guardDuty)
	if err != nil {
		return fmt.Errorf("could not set scanning interval back to default value: %v", err)
	}
	for _, account := range guardDuty.EnabledAccounts {
		err := client.DisableAwsGuardDuty(account)
		if err != nil {
			return fmt.Errorf("could not disable GuardDuty for account %#v: %v", account, err)
		}
	}
	return nil
}
