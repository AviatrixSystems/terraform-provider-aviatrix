package aviatrix

import (
	"fmt"
	"log"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"aws_gateway_role_app": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS App role ARN for gateways.",
			},
			"aws_gateway_role_ec2": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS EC2 role ARN for gateways.",
			},
			"aws_access_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS Access Key.",
			},
			"gcloud_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GCloud Project ID.",
			},
			"arm_subscription_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Azure Subscription ID.",
			},
			"awsgov_account_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS Gov Account number to associate with Aviatrix account.",
			},
			"awsgov_access_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS Gov Access Key.",
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
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("aviatrix Account: %s", err)
	}

	d.Set("account_name", acc.AccountName)
	d.Set("cloud_type", acc.CloudType)
	d.Set("aws_account_number", acc.AwsAccountNumber)
	d.Set("aws_access_key", acc.AwsAccessKey)
	d.Set("aws_secret_key", acc.AwsSecretKey)
	d.Set("aws_role_arn", acc.AwsRoleApp)
	d.Set("aws_role_ec2", acc.AwsRoleEc2)
	d.Set("aws_gateway_role_app", acc.AwsGatewayRoleApp)
	d.Set("aws_gateway_role_ec2", acc.AwsGatewayRoleEc2)
	d.Set("gcloud_project_id", acc.GcloudProjectName)
	d.Set("arm_subscription_id", acc.ArmSubscriptionId)
	d.Set("awsgov_account_number", acc.AwsgovAccountNumber)
	d.Set("awsgov_access_key", acc.AwsgovAccessKey)
	d.SetId(acc.AccountName)

	return nil
}
