package aviatrix

import (
	"log"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
)

func dataSourceAviatrixAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAviatrixAccountRead,

		Schema: map[string]*schema.Schema {
			"account_name": {
				Type: schema.TypeString,
				Required: true,
			},
			"account_email": &schema.Schema{
				Type:     schema.TypeString,
				Optional:    true,
			},
			"cloud_type": &schema.Schema{
				Type:     schema.TypeInt,
				Optional:    true,
			},
			"aws_account_number": &schema.Schema{
				Type:     schema.TypeString,
				Optional:    true,
			},
			"aws_iam": &schema.Schema{
				Type:     schema.TypeString,
				Optional:    true,
			},
			"aws_role_arn": &schema.Schema{
				Type:     schema.TypeString,
				Optional:    true,
			},
			"aws_role_ec2": &schema.Schema{
				Type:     schema.TypeString,
				Optional:    true,
			},
			"aws_access_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional:    true,
			},
			"aws_secret_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func dataSourceAviatrixAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	account := &goaviatrix.Account{
		AccountName:  d.Get("account_name").(string),
	}
	log.Printf("[INFO] Looking for Aviatrix account: %#v", account)
	acc, err := client.GetAccount(account)
	if err != nil {
		return fmt.Errorf("Aviatrix Account: %s", err)
	}

	if acc != nil {
		d.Set("account_name", acc.AccountName)
		d.Set("account_email", acc.AccountEmail)
		d.Set("cloud_type", acc.CloudType)
		d.Set("aws_account_number", acc.AwsAccountNumber)
		d.Set("aws_access_key", acc.AwsAccessKey)
		d.Set("aws_secret_key", acc.AwsSecretKey)
		d.SetId(acc.AccountName)
	} else {
		d.SetId("")
	}
	return nil
}
