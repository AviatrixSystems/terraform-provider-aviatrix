package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

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
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Account name",
				ForceNew:    true,
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Region.",
				ForceNew:    true,
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
	}
}

func resourceAviatrixAwsGuardDutyCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*goaviatrix.Client)
	guardDuty := marshalAwsGuardDutyInput(d)

	err = client.EnableAwsGuardDuty(guardDuty)
	if err != nil {
		return fmt.Errorf("could not enable AWS GuardDuty: %v", err)
	}
	d.SetId(guardDuty.ID())
	defer captureErr(resourceAviatrixAwsGuardDutyRead, d, meta, &err)
	err = client.UpdateAwsGuardDutyExcludedIPs(guardDuty)
	if err != nil {
		return fmt.Errorf("could not set excluded IPs: %v", err)
	}
	return err
}

func resourceAviatrixAwsGuardDutyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	accName := d.Get("account_name").(string)
	region := d.Get("region").(string)
	if accName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no account_name received. Import Id is %s", id)
		parts := strings.Split(id, "~~")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return fmt.Errorf("invalid import ID: %q", id)
		}
		accName, region = parts[0], parts[1]
		d.SetId(id)
	}

	acc, err := client.GetAwsGuardDutyAccount(accName, region)
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not get guard duty account: %v", err)
	}

	d.Set("account_name", acc.AccountName)
	d.Set("region", acc.Region)
	if err := d.Set("excluded_ips", acc.ExcludedIPs); err != nil {
		return fmt.Errorf("setting excluded_ips: %v", err)
	}

	d.SetId(acc.ID())
	return nil
}

func resourceAviatrixAwsGuardDutyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	account := marshalAwsGuardDutyInput(d)

	if d.HasChange("excluded_ips") {
		err := client.UpdateAwsGuardDutyExcludedIPs(account)
		if err != nil {
			return fmt.Errorf("could not edit GuardDuty excluded IPs: %v", err)
		}
	}
	return nil
}

func resourceAviatrixAwsGuardDutyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	account := marshalAwsGuardDutyInput(d)

	err := client.DisableAwsGuardDuty(account)
	if err != nil {
		return fmt.Errorf("could not disable GuardDuty: %v", err)
	}
	return nil
}

func marshalAwsGuardDutyInput(d *schema.ResourceData) *goaviatrix.AwsGuardDutyAccount {
	var excludedIPs []string
	for _, ip := range d.Get("excluded_ips").(*schema.Set).List() {
		excludedIPs = append(excludedIPs, ip.(string))
	}
	return &goaviatrix.AwsGuardDutyAccount{
		AccountName: d.Get("account_name").(string),
		Region:      d.Get("region").(string),
		ExcludedIPs: excludedIPs,
	}
}
