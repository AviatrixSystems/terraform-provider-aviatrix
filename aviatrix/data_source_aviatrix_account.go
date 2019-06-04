package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func dataSourceAviatrixAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAviatrixAccountRead,

		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Account name. This can be used for logging in to CloudN console or UserConnect controller.",
			},
			"cloud_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Type of cloud service provider.",
			},
			"aws_account_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS Account number to associate with Aviatrix account.",
			},
			"aws_role_arn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS App role ARN.",
			},
			"aws_role_ec2": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS EC2 role ARN",
			},
			"aws_access_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS Access Key.",
			},
		},
	}
}

func dataSourceAviatrixAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	account := &goaviatrix.Account{
		AccountName: d.Get("account_name").(string),
	}
	log.Printf("[INFO] Looking for Aviatrix account: %#v", account)
	acc, err := client.GetAccount(account)
	if err != nil {
		return fmt.Errorf("aviatrix Account: %s", err)
	}

	if acc != nil {
		d.Set("account_name", acc.AccountName)
		d.Set("cloud_type", acc.CloudType)
		d.Set("aws_account_number", acc.AwsAccountNumber)
		d.Set("aws_access_key", acc.AwsAccessKey)
		d.Set("aws_secret_key", acc.AwsSecretKey)
		d.Set("aws_role_arn", acc.AwsRoleApp)
		d.Set("aws_role_ec2", acc.AwsRoleEc2)
		d.SetId(acc.AccountName)
	} else {
		d.SetId("")
	}
	return nil
}
