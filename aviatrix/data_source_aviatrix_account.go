package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
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
			"azuregov_subscription_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Azure Gov Subscription ID.",
			},
			"awsgov_account_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS Gov Account number to associate with Aviatrix account.",
			},
			"awsgov_iam": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "AWSGov IAM role based flag.",
			},
			"awsgov_role_app": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWSGov App role ARN.",
			},
			"awsgov_role_ec2": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWSGov EC2 role ARN.",
			},
			"awschina_account_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS China Account number.",
			},
			"awschina_iam": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "AWS China IAM-role based flag.",
			},
			"awschina_role_app": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS China App role ARN.",
			},
			"awschina_role_ec2": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS China EC2 role ARN.",
			},
			"azurechina_subscription_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Azure China subscription ID.",
			},
			"alicloud_account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Alibaba Cloud Account ID to associate with Aviatrix account.",
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

	_ = d.Set("account_name", acc.AccountName)
	_ = d.Set("cloud_type", acc.CloudType)
	if acc.CloudType == goaviatrix.AWS {
		_ = d.Set("aws_account_number", acc.AwsAccountNumber)
	}
	_ = d.Set("aws_role_arn", acc.AwsRoleApp)
	_ = d.Set("aws_role_ec2", acc.AwsRoleEc2)
	_ = d.Set("aws_gateway_role_app", acc.AwsGatewayRoleApp)
	_ = d.Set("aws_gateway_role_ec2", acc.AwsGatewayRoleEc2)
	_ = d.Set("gcloud_project_id", acc.GcloudProjectName)
	_ = d.Set("arm_subscription_id", acc.ArmSubscriptionId)
	_ = d.Set("azuregov_subscription_id", acc.AzuregovSubscriptionId)
	if goaviatrix.IsCloudType(acc.CloudType, goaviatrix.AWSGov) {
		_ = d.Set("awsgov_account_number", acc.AwsgovAccountNumber)
		if acc.AwsgovRoleEc2 != "" {
			_ = d.Set("awsgov_iam", true)
			_ = d.Set("awsgov_role_app", acc.AwsgovRoleApp)
			_ = d.Set("awsgov_role_ec2", acc.AwsgovRoleEc2)
		} else {
			_ = d.Set("awsgov_iam", false)
		}
	} else if goaviatrix.IsCloudType(acc.CloudType, goaviatrix.AWSChina) {
		_ = d.Set("awschina_account_number", acc.AwsChinaAccountNumber)
		_ = d.Set("awschina_role_app", acc.AwsChinaRoleApp)
		_ = d.Set("awschina_role_ec2", acc.AwsChinaRoleEc2)
		if acc.AwsChinaRoleEc2 == "" {
			_ = d.Set("awschina_iam", true)
		} else {
			_ = d.Set("awschina_iam", false)
		}
	} else if goaviatrix.IsCloudType(acc.CloudType, goaviatrix.AzureChina) {
		_ = d.Set("azurechina_subscription_id", acc.AzureChinaSubscriptionId)
	}

	if acc.CloudType == goaviatrix.AliCloud {
		_ = d.Set("alicloud_account_id", acc.AwsAccountNumber)
	}

	d.SetId(acc.AccountName)

	return nil
}
