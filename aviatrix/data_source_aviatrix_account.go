package aviatrix

import (
	"fmt"
	"log"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
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
			"awsts_account_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS Top Secret Region Account Number.",
			},
			"awsts_cap_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS Top Secret Region CAP Endpoint URL.",
			},
			"awsts_cap_agency": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS Top Secret Region CAP Agency.",
			},
			"awsts_cap_mission": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS Top Secret Region CAP Mission.",
			},
			"awsts_cap_role_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS Top Secret Region CAP Role Name.",
			},
			"awss_account_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS Secret Region Account Number.",
			},
			"awss_cap_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS Secret Region CAP Endpoint URL.",
			},
			"awss_cap_agency": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS Secret Region CAP Agency.",
			},
			"awss_cap_account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS Secret Region CAP Account Name.",
			},
			"awss_cap_role_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS Secret Region CAP Role Name.",
			},
			"awsts_cap_cert_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS Top Secret Region CAP Certificate file path on the controller.",
			},
			"awsts_cap_cert_key_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS Top Secret Region CAP Certificate Key file path on the controller.",
			},
			"aws_ca_cert_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS Top Secret Region or Secret Region Custom Certificate Authority file path on the controller.",
			},
			"awss_cap_cert_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS Secret Region CAP Certificate file path on the controller.",
			},
			"awss_cap_cert_key_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS Secret Region CAP Certificate Key file path on the controller.",
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
	if acc.CloudType == goaviatrix.AWS {
		d.Set("aws_account_number", acc.AwsAccountNumber)
	}
	d.Set("aws_role_arn", acc.AwsRoleApp)
	d.Set("aws_role_ec2", acc.AwsRoleEc2)
	d.Set("aws_gateway_role_app", acc.AwsGatewayRoleApp)
	d.Set("aws_gateway_role_ec2", acc.AwsGatewayRoleEc2)
	d.Set("gcloud_project_id", acc.GcloudProjectName)
	d.Set("arm_subscription_id", acc.ArmSubscriptionId)
	d.Set("azuregov_subscription_id", acc.AzuregovSubscriptionId)
	if goaviatrix.IsCloudType(acc.CloudType, goaviatrix.AWSGov) {
		d.Set("awsgov_account_number", acc.AwsgovAccountNumber)
		if acc.AwsgovRoleEc2 != "" {
			d.Set("awsgov_iam", true)
			d.Set("awsgov_role_app", acc.AwsgovRoleApp)
			d.Set("awsgov_role_ec2", acc.AwsgovRoleEc2)
		} else {
			d.Set("awsgov_iam", false)
		}
	} else if goaviatrix.IsCloudType(acc.CloudType, goaviatrix.AWSChina) {
		d.Set("awschina_account_number", acc.AwsChinaAccountNumber)
		d.Set("awschina_role_app", acc.AwsChinaRoleApp)
		d.Set("awschina_role_ec2", acc.AwsChinaRoleEc2)
		if acc.AwsChinaRoleEc2 == "" {
			d.Set("awschina_iam", true)
		} else {
			d.Set("awschina_iam", false)
		}
	} else if goaviatrix.IsCloudType(acc.CloudType, goaviatrix.AzureChina) {
		d.Set("azurechina_subscription_id", acc.AzureChinaSubscriptionId)
	}

	if acc.CloudType == goaviatrix.AliCloud {
		d.Set("alicloud_account_id", acc.AwsAccountNumber)
	}

	if goaviatrix.IsCloudType(acc.CloudType, goaviatrix.AWSTS) {
		d.Set("awsts_account_number", acc.AwsTsAccountNumber)
		d.Set("awsts_cap_url", acc.AwsTsCapUrl)
		d.Set("awsts_cap_agency", acc.AwsTsCapAgency)
		d.Set("awsts_cap_mission", acc.AwsTsCapMission)
		d.Set("awsts_cap_role_name", acc.AwsTsCapRoleName)
		d.Set("awsts_cap_cert_path", acc.AwsTsCapCertPath)
		d.Set("awsts_cap_cert_key_path", acc.AwsTsCapCertKeyPath)
		d.Set("aws_ca_cert_path", acc.AwsCaCertPath)
	} else if goaviatrix.IsCloudType(acc.CloudType, goaviatrix.AWSS) {
		d.Set("awss_account_number", acc.AwsSAccountNumber)
		d.Set("awss_cap_url", acc.AwsSCapUrl)
		d.Set("awss_cap_agency", acc.AwsSCapAgency)
		d.Set("awss_cap_account_name", acc.AwsSCapAccountName)
		d.Set("awss_cap_role_name", acc.AwsSCapRoleName)
		d.Set("awss_cap_cert_path", acc.AwsSCapCertPath)
		d.Set("awss_cap_cert_key_path", acc.AwsSCapCertKeyPath)
		d.Set("aws_ca_cert_path", acc.AwsCaCertPath)
	}

	d.SetId(acc.AccountName)

	return nil
}
