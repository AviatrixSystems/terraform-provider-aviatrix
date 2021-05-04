package aviatrix

import (
	"context"
	"fmt"
	"log"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAviatrixAccountCreate,
		ReadContext:   resourceAviatrixAccountRead,
		UpdateContext: resourceAviatrixAccountUpdate,
		DeleteContext: resourceAviatrixAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Account name. This can be used for logging in to CloudN console or UserConnect controller.",
			},
			"cloud_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateCloudType,
				Description:  "Type of cloud service provider.",
			},
			"aws_account_number": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAwsAccountNumber,
				Description:  "AWS Account number to associate with Aviatrix account. Should be 12 digits.",
			},
			"aws_iam": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "AWS IAM-role based flag.",
			},
			"aws_gateway_role_app": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS App role ARN for gateways.",
			},
			"aws_gateway_role_ec2": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS EC2 role ARN for gateways.",
			},
			"aws_access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "AWS Access Key.",
			},
			"aws_secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "AWS Secret Key.",
			},
			"awsgov_account_number": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAwsAccountNumber,
				Description:  "AWS Gov Account number to associate with Aviatrix account.",
			},
			"awsgov_access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "AWS Gov Access Key.",
			},
			"awsgov_secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "AWS Gov Secret Key.",
			},
			"gcloud_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "GCloud Project ID.",
			},
			"gcloud_project_credentials_filepath": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "GCloud Project credentials local file path.",
			},
			"arm_subscription_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Azure Subscription ID.",
			},
			"arm_directory_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Azure Directory ID.",
			},
			"arm_application_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Azure Application ID.",
			},
			"arm_application_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Azure Application Key.",
			},
			"oci_tenancy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "OCI Tenancy OCID.",
			},
			"oci_user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "OCI User OCID.",
			},
			"oci_compartment_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "OCI Compartment OCID.",
			},
			"oci_api_private_key_filepath": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "OCI API Private Key local file path.",
			},
			"azuregov_subscription_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Azure Gov Subscription ID.",
			},
			"azuregov_directory_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Azure Gov Directory ID.",
			},
			"azuregov_application_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Azure Gov Application ID.",
			},
			"azuregov_application_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Azure Gov Application Key.",
			},
			"alicloud_account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Alibaba Cloud Account ID to associate with Aviatrix account.",
			},
			"alicloud_access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Alibaba Cloud Access Key.",
			},
			"alicloud_secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Alibaba Cloud Secret Key.",
			},
			"audit_account": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable account audit.",
			},
			"awschina_account_number": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAwsAccountNumber,
				Description:  "AWS China Account Number.",
			},
			"awschina_iam": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "AWS China IAM-role based flag.",
			},
			"awschina_access_key": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"awschina_role_app", "awschina_role_ec2"},
				Description:   "AWS China Access Key.",
			},
			"awschina_secret_key": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"awschina_role_app", "awschina_role_ec2"},
				Description:   "AWS China Secret Key.",
			},
			"azurechina_subscription_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Azure China Subscription ID.",
			},
			"azurechina_directory_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Azure China Directory ID.",
			},
			"azurechina_application_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Azure China Application ID.",
			},
			"azurechina_application_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Azure China Application Key.",
			},
			"aws_role_app": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "AWS App role ARN.",
			},
			"aws_role_ec2": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "AWS EC2 role ARN.",
			},
			"awschina_role_app": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "AWS China App Role ARN.",
			},
			"awschina_role_ec2": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "AWS China EC2 Role ARN.",
			},
		},
	}
}

func resourceAviatrixAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	account := &goaviatrix.Account{
		AccountName:                           d.Get("account_name").(string),
		CloudType:                             d.Get("cloud_type").(int),
		AwsAccountNumber:                      d.Get("aws_account_number").(string),
		AwsRoleApp:                            d.Get("aws_role_app").(string),
		AwsRoleEc2:                            d.Get("aws_role_ec2").(string),
		AwsGatewayRoleApp:                     d.Get("aws_gateway_role_app").(string),
		AwsGatewayRoleEc2:                     d.Get("aws_gateway_role_ec2").(string),
		AwsAccessKey:                          d.Get("aws_access_key").(string),
		AwsSecretKey:                          d.Get("aws_secret_key").(string),
		AwsgovAccountNumber:                   d.Get("awsgov_account_number").(string),
		AwsgovAccessKey:                       d.Get("awsgov_access_key").(string),
		AwsgovSecretKey:                       d.Get("awsgov_secret_key").(string),
		GcloudProjectName:                     d.Get("gcloud_project_id").(string),
		GcloudProjectCredentialsFilepathLocal: d.Get("gcloud_project_credentials_filepath").(string),
		ArmSubscriptionId:                     d.Get("arm_subscription_id").(string),
		ArmApplicationEndpoint:                d.Get("arm_directory_id").(string),
		ArmApplicationClientId:                d.Get("arm_application_id").(string),
		ArmApplicationClientSecret:            d.Get("arm_application_key").(string),
		AzuregovSubscriptionId:                d.Get("azuregov_subscription_id").(string),
		AzuregovApplicationEndpoint:           d.Get("azuregov_directory_id").(string),
		AzuregovApplicationClientId:           d.Get("azuregov_application_id").(string),
		AzuregovApplicationClientSecret:       d.Get("azuregov_application_key").(string),
		OciTenancyID:                          d.Get("oci_tenancy_id").(string),
		OciUserID:                             d.Get("oci_user_id").(string),
		OciCompartmentID:                      d.Get("oci_compartment_id").(string),
		OciApiPrivateKeyFilePath:              d.Get("oci_api_private_key_filepath").(string),
		AlicloudAccountId:                     d.Get("alicloud_account_id").(string),
		AlicloudAccessKey:                     d.Get("alicloud_access_key").(string),
		AlicloudSecretKey:                     d.Get("alicloud_secret_key").(string),
		AwsChinaAccountNumber:                 d.Get("awschina_account_number").(string),
		AwsChinaRoleApp:                       d.Get("awschina_role_app").(string),
		AwsChinaRoleEc2:                       d.Get("awschina_role_ec2").(string),
		AwsChinaAccessKey:                     d.Get("awschina_access_key").(string),
		AwsChinaSecretKey:                     d.Get("awschina_secret_key").(string),
		AzureChinaSubscriptionId:              d.Get("azurechina_subscription_id").(string),
		AzureChinaApplicationEndpoint:         d.Get("azurechina_directory_id").(string),
		AzureChinaApplicationClientId:         d.Get("azurechina_application_id").(string),
		AzureChinaApplicationClientSecret:     d.Get("azurechina_application_key").(string),
	}

	awsIam := d.Get("aws_iam").(bool)
	if awsIam {
		account.AwsIam = "true"
	} else {
		account.AwsIam = "false"
	}

	_, gatewayRoleAppOk := d.GetOk("aws_gateway_role_app")
	_, gatewayRoleEc2Ok := d.GetOk("aws_gateway_role_ec2")

	if gatewayRoleAppOk || gatewayRoleEc2Ok {
		if !goaviatrix.IsCloudType(account.CloudType, goaviatrix.AWS) {
			return diag.Errorf("could not create Aviatrix Account: aws_gateway_role_app and aws_gateway_role_ec2 can only be used with AWS (1)")
		}
		if !awsIam {
			return diag.Errorf("could not create Aviatrix Account: aws_gateway_role_app and aws_gateway_role_ec2 can only be used when awsIam is enabled")
		}
		if !(gatewayRoleAppOk && gatewayRoleEc2Ok) {
			return diag.Errorf("could not create Aviatrix Account: must provide both aws_gateway_role_app and aws_gateway_role_ec2 when using separate IAM role and policy for gateways")
		}
	}

	awsChinaIam := d.Get("awschina_iam").(bool)
	if awsChinaIam {
		account.AwsChinaIam = "true"
	} else {
		account.AwsChinaIam = "false"
	}

	if !goaviatrix.IsCloudType(account.CloudType, goaviatrix.AWSChina) && (awsChinaIam || account.AwsChinaRoleApp != "" || account.AwsChinaRoleEc2 != "" || account.AwsChinaAccessKey != "" || account.AwsChinaSecretKey != "") {
		return diag.Errorf("could not create Aviatrix Account: 'awschina_iam', 'awschina_role_app', 'awschina_role_ec2', 'awschina_access_key' and 'awschina_secret_key' can only be set when cloud_type is AWSChina (1024)")
	}
	if !goaviatrix.IsCloudType(account.CloudType, goaviatrix.AzureChina) && (account.AzureChinaSubscriptionId != "" || account.AzureChinaApplicationEndpoint != "" || account.AzureChinaApplicationClientId != "" || account.AzureChinaApplicationClientSecret != "") {
		return diag.Errorf("could not create Aviatrix Account: `azurechina_subscription_id', 'azurechina_directory_id', 'azurechina_application_id' and 'azurechina_application_key' can only be set when cloud_type is AzureChina (2048)")
	}

	if account.CloudType == goaviatrix.AWS {
		if account.AwsAccountNumber == "" {
			return diag.Errorf("aws account number is needed for aws cloud")
		}
		if account.AwsIam != "true" && account.AwsIam != "false" {
			return diag.Errorf("aws iam can only be 'true' or 'false'")
		}

		log.Printf("[INFO] Creating Aviatrix account: %#v", account)
		if awsIam {
			if _, ok := d.GetOk("aws_role_app"); !ok {
				account.AwsRoleApp = fmt.Sprintf("arn:aws:iam::%s:role/aviatrix-role-app", account.AwsAccountNumber)
			}
			if _, ok := d.GetOk("aws_role_ec2"); !ok {
				account.AwsRoleEc2 = fmt.Sprintf("arn:aws:iam::%s:role/aviatrix-role-ec2", account.AwsAccountNumber)
			}
			log.Printf("[TRACE] Reading Aviatrix account aws_role_app: [%s]", d.Get("aws_role_app").(string))
			log.Printf("[TRACE] Reading Aviatrix account aws_role_ec2: [%s]", d.Get("aws_role_ec2").(string))
		}
	} else if account.CloudType == goaviatrix.GCP {
		if account.GcloudProjectCredentialsFilepathLocal == "" {
			return diag.Errorf("gcloud project credentials local filepath needed to upload file to controller")
		}
		log.Printf("[INFO] Creating Aviatrix account: %#v", account)
	} else if account.CloudType == goaviatrix.Azure {
		if account.ArmSubscriptionId == "" {
			return diag.Errorf("arm subscription id needed for azure cloud")
		}
		if account.ArmApplicationEndpoint == "" {
			return diag.Errorf("arm directory id needed for azure cloud")
		}
		if account.ArmApplicationClientId == "" {
			return diag.Errorf("arm application id needed for azure cloud")
		}
		if account.ArmApplicationClientSecret == "" {
			return diag.Errorf("arm application key needed for azure cloud")
		}
	} else if account.CloudType == goaviatrix.OCI {
		if account.OciTenancyID == "" {
			return diag.Errorf("oci tenancy ocid needed for oracle cloud")
		}
		if account.OciUserID == "" {
			return diag.Errorf("oci user id needed for oracle cloud")
		}
		if account.OciCompartmentID == "" {
			return diag.Errorf("oci compartment ocid needed for oracle cloud")
		}
		if account.OciApiPrivateKeyFilePath == "" {
			return diag.Errorf("oci api private key filepath needed to upload file to controller")
		}
	} else if account.CloudType == goaviatrix.AWSGov {
		if account.AwsgovAccountNumber == "" {
			return diag.Errorf("aws gov account number needed for aws gov cloud")
		}
		if account.AwsgovAccessKey == "" {
			return diag.Errorf("aws gov access key needed for aws gov cloud")
		}
		if account.AwsgovSecretKey == "" {
			return diag.Errorf("aws gov secret key needed for aws gov cloud")
		}
	} else if account.CloudType == goaviatrix.AzureGov {
		if account.AzuregovSubscriptionId == "" {
			return diag.Errorf("azure gov subsription id needed when creating an account for arm gov cloud")
		}
		if account.AzuregovApplicationEndpoint == "" {
			return diag.Errorf("azure gov directory id needed when creating an account for arm gov cloud")
		}
		if account.AzuregovApplicationClientId == "" {
			return diag.Errorf("azure gov application id needed when creating an account for arm gov cloud")
		}
		if account.AzuregovApplicationClientSecret == "" {
			return diag.Errorf("azure gov application key needed when creating an account for arm gov cloud")
		}
	} else if goaviatrix.IsCloudType(account.CloudType, goaviatrix.AWSChina) {
		if account.AwsChinaAccountNumber == "" {
			return diag.Errorf("could not create Aviatrix Account in AWS China (1024): 'awschina_account_number' is required")
		}
		if awsChinaIam {
			if account.AwsChinaRoleApp == "" {
				account.AwsChinaRoleApp = fmt.Sprintf("arn:aws-cn:iam::%s:role/aviatrix-role-app", account.AwsChinaAccountNumber)
			}
			if account.AwsChinaRoleEc2 == "" {
				account.AwsChinaRoleEc2 = fmt.Sprintf("arn:aws-cn:iam::%s:role/aviatrix-role-ec2", account.AwsChinaAccountNumber)
			}
		} else {
			if account.AwsChinaAccessKey == "" {
				return diag.Errorf("could not create Aviatrix Account in AWSChina (1024): 'awschina_access_key' is required when 'awschina_iam' is false")
			}
			if account.AwsChinaSecretKey == "" {
				return diag.Errorf("could not create Aviatrix Account in AWSChina (1024): 'awschina_secret_key' is required when 'awschina_iam' is false")
			}
		}
	} else if goaviatrix.IsCloudType(account.CloudType, goaviatrix.AzureChina) {
		if account.AzureChinaSubscriptionId == "" {
			return diag.Errorf("could not create Aviatrix Account in AzureChina (2048): 'azurechina_subscription_id' is required")
		}
		if account.AzureChinaApplicationEndpoint == "" {
			return diag.Errorf("could not create Aviatrix Account in AzureChina (2048): 'azurechina_directory_id' is required")
		}
		if account.AzureChinaApplicationClientId == "" {
			return diag.Errorf("could not create Aviatrix Account in AzureChina (2048): 'azurechina_application_id' is required")
		}
		if account.AzureChinaApplicationClientSecret == "" {
			return diag.Errorf("could not create Aviatrix Account in AzureChina (2048): 'azurechina_application_key' is required")
		}
	} else if account.CloudType == goaviatrix.AliCloud {
		if account.AlicloudAccountId == "" {
			return diag.Errorf("alicloud_account_id is required for alibaba cloud")
		}
		if account.AlicloudAccessKey == "" {
			return diag.Errorf("alicloud_access_key is required for alibaba cloud")
		}
		if account.AlicloudSecretKey == "" {
			return diag.Errorf("alicloud_secret_key is required for alibaba cloud")
		}
	} else {
		return diag.Errorf("cloud type can only be either AWS (1), GCP (4), Azure (8), OCI (16), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048) or Alibaba Cloud (8192)")
	}

	var err error
	if account.CloudType == goaviatrix.GCP {
		err = client.CreateGCPAccount(account)
	} else if account.CloudType == goaviatrix.OCI {
		err = client.CreateOCIAccount(account)
	} else {
		err = client.CreateAccount(account)
	}
	if err != nil {
		return diag.Errorf("failed to create Aviatrix Account: %s", err)
	}

	d.SetId(account.AccountName)
	return resourceAviatrixAccountRead(ctx, d, meta)
}

func resourceAviatrixAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	var diags diag.Diagnostics

	accountName := d.Get("account_name").(string)
	if accountName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no account name received. Import Id is %s", id)
		d.Set("account_name", id)
		d.Set("audit_account", true)
		d.SetId(id)
	}

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
		return diag.Errorf("aviatrix Account: %s", err)
	}

	if acc != nil {
		d.Set("account_name", acc.AccountName)
		d.Set("cloud_type", acc.CloudType)
		if acc.CloudType == goaviatrix.AWS {
			d.Set("aws_account_number", acc.AwsAccountNumber)
			if acc.AwsRoleEc2 != "" {
				//force default setting and save to .tfstate file
				d.Set("aws_access_key", "")
				d.Set("aws_secret_key", "")
				d.Set("aws_iam", true)
				d.Set("aws_role_app", acc.AwsRoleApp)
				d.Set("aws_role_ec2", acc.AwsRoleEc2)
				d.Set("aws_gateway_role_app", acc.AwsGatewayRoleApp)
				d.Set("aws_gateway_role_ec2", acc.AwsGatewayRoleEc2)
			} else {
				d.Set("aws_iam", false)
			}
		} else if acc.CloudType == goaviatrix.GCP {
			d.Set("gcloud_project_id", acc.GcloudProjectName)
		} else if acc.CloudType == goaviatrix.Azure {
			d.Set("arm_subscription_id", acc.ArmSubscriptionId)
		} else if acc.CloudType == goaviatrix.AWSGov {
			d.Set("awsgov_account_number", acc.AwsgovAccountNumber)
		} else if acc.CloudType == goaviatrix.AzureGov {
			d.Set("azuregov_subscription_id", acc.AzuregovSubscriptionId)
		} else if goaviatrix.IsCloudType(acc.CloudType, goaviatrix.AWSChina) {
			d.Set("awschina_account_number", acc.AwsChinaAccountNumber)
			d.Set("awschina_role_app", acc.AwsChinaRoleApp)
			d.Set("awschina_role_ec2", acc.AwsChinaRoleEc2)
			if acc.AwsChinaRoleEc2 != "" {
				// Force access key and secret key to be empty
				d.Set("awschina_access_key", "")
				d.Set("awschina_secret_key", "")
				d.Set("awschina_iam", true)
			} else {
				d.Set("awschina_iam", false)
			}
		} else if goaviatrix.IsCloudType(acc.CloudType, goaviatrix.AzureChina) {
			d.Set("azurechina_subscription_id", acc.AzureChinaSubscriptionId)
		} else if acc.CloudType == goaviatrix.AliCloud {
			d.Set("alicloud_account_id", acc.AwsAccountNumber)
		}
		d.SetId(acc.AccountName)
	}

	if d.Get("audit_account").(bool) {
		err = client.AuditAccount(ctx, account)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Aviatrix Account failed audit",
				Detail:   fmt.Sprintf("%v", err),
			})
		}
	}

	return diags
}

func resourceAviatrixAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	account := &goaviatrix.Account{
		AccountName:                           d.Get("account_name").(string),
		CloudType:                             d.Get("cloud_type").(int),
		AwsAccountNumber:                      d.Get("aws_account_number").(string),
		AwsRoleApp:                            d.Get("aws_role_app").(string),
		AwsRoleEc2:                            d.Get("aws_role_ec2").(string),
		AwsGatewayRoleApp:                     d.Get("aws_gateway_role_app").(string),
		AwsGatewayRoleEc2:                     d.Get("aws_gateway_role_ec2").(string),
		AwsAccessKey:                          d.Get("aws_access_key").(string),
		AwsSecretKey:                          d.Get("aws_secret_key").(string),
		AwsgovAccountNumber:                   d.Get("awsgov_account_number").(string),
		AwsgovAccessKey:                       d.Get("awsgov_access_key").(string),
		AwsgovSecretKey:                       d.Get("awsgov_secret_key").(string),
		GcloudProjectName:                     d.Get("gcloud_project_id").(string),
		GcloudProjectCredentialsFilepathLocal: d.Get("gcloud_project_credentials_filepath").(string),
		ArmSubscriptionId:                     d.Get("arm_subscription_id").(string),
		ArmApplicationEndpoint:                d.Get("arm_directory_id").(string),
		ArmApplicationClientId:                d.Get("arm_application_id").(string),
		ArmApplicationClientSecret:            d.Get("arm_application_key").(string),
		AzuregovSubscriptionId:                d.Get("azuregov_subscription_id").(string),
		AzuregovApplicationEndpoint:           d.Get("azuregov_directory_id").(string),
		AzuregovApplicationClientId:           d.Get("azuregov_application_id").(string),
		AzuregovApplicationClientSecret:       d.Get("azuregov_application_key").(string),
		OciTenancyID:                          d.Get("oci_tenancy_id").(string),
		OciUserID:                             d.Get("oci_user_id").(string),
		OciCompartmentID:                      d.Get("oci_compartment_id").(string),
		OciApiPrivateKeyFilePath:              d.Get("oci_api_private_key_filepath").(string),
		AlicloudAccountId:                     d.Get("alicloud_account_id").(string),
		AlicloudAccessKey:                     d.Get("alicloud_access_key").(string),
		AlicloudSecretKey:                     d.Get("alicloud_secret_key").(string),
		AwsChinaAccountNumber:                 d.Get("awschina_account_number").(string),
		AwsChinaRoleApp:                       d.Get("awschina_role_app").(string),
		AwsChinaRoleEc2:                       d.Get("awschina_role_ec2").(string),
		AwsChinaAccessKey:                     d.Get("awschina_access_key").(string),
		AwsChinaSecretKey:                     d.Get("awschina_secret_key").(string),
		AzureChinaSubscriptionId:              d.Get("azurechina_subscription_id").(string),
		AzureChinaApplicationEndpoint:         d.Get("azurechina_directory_id").(string),
		AzureChinaApplicationClientId:         d.Get("azurechina_application_id").(string),
		AzureChinaApplicationClientSecret:     d.Get("azurechina_application_key").(string),
	}

	awsIam := d.Get("aws_iam").(bool)
	if awsIam {
		account.AwsIam = "true"
	} else {
		account.AwsIam = "false"
	}

	awsChinaIam := d.Get("awschina_iam").(bool)
	if awsChinaIam {
		account.AwsChinaIam = "true"
	} else {
		account.AwsChinaIam = "false"
	}

	log.Printf("[INFO] Updating Aviatrix account: %#v", account)

	d.Partial(true)

	if d.HasChange("cloud_type") {
		return diag.Errorf("update cloud_type is not allowed")
	}

	if d.HasChange("account_name") {
		return diag.Errorf("update account name is not allowed")
	}

	if d.HasChanges("aws_gateway_role_app", "aws_gateway_role_ec2") {
		if !goaviatrix.IsCloudType(account.CloudType, goaviatrix.AWS) {
			return diag.Errorf("could not update Aviatrix Account: aws_gateway_role_app and aws_gateway_role_ec2 can only be used with AWS (1)")
		}
		_, gatewayRoleAppOk := d.GetOk("aws_gateway_role_app")
		_, gatewayRoleEc2Ok := d.GetOk("aws_gateway_role_ec2")
		if gatewayRoleAppOk != gatewayRoleEc2Ok {
			return diag.Errorf("failed to update Aviatrix account: must provide both aws_gateway_role_app and aws_gateway_role_ec2 when using separate IAM role and policy for gateways")
		}
	}

	if account.CloudType == goaviatrix.AWS {
		if d.HasChanges("aws_account_number", "aws_access_key", "aws_secret_key", "aws_iam", "aws_role_app", "aws_role_ec2", "aws_gateway_role_app", "aws_gateway_role_ec2") {
			err := client.UpdateAccount(account)
			if err != nil {
				return diag.Errorf("failed to update Aviatrix Account: %s", err)
			}
		}
	} else if account.CloudType == goaviatrix.GCP {
		if d.HasChange("gcloud_project_id") || d.HasChange("gcloud_project_credentials_filepath") {
			err := client.UpdateGCPAccount(account)
			if err != nil {
				return diag.Errorf("failed to update Aviatrix Account: %s", err)
			}
		}
	} else if account.CloudType == goaviatrix.Azure {
		if d.HasChange("arm_subscription_id") || d.HasChange("arm_directory_id") || d.HasChange("arm_application_id") || d.HasChange("arm_application_key") {
			err := client.UpdateAccount(account)
			if err != nil {
				return diag.Errorf("failed to update Aviatrix Account: %s", err)
			}
		}
	} else if account.CloudType == goaviatrix.OCI {
		if d.HasChange("oci_tenancy_id") || d.HasChange("oci_user_id") || d.HasChange("oci_compartment_id") || d.HasChange("oci_api_private_key_filepath") {
			return diag.Errorf("updating OCI account is not supported")
		}
	} else if account.CloudType == goaviatrix.AWSGov {
		if d.HasChange("awsgov_account_number") || d.HasChange("awsgov_access_key") || d.HasChange("awsgov_secret_key") {
			err := client.UpdateAccount(account)
			if err != nil {
				return diag.Errorf("failed to update Aviatrix Account: %s", err)
			}
		}
	} else if account.CloudType == goaviatrix.AzureGov {
		if d.HasChanges("azuregov_subscription_id", "azuregov_directory_id", "azuregov_application_id", "azuregov_application_key") {
			err := client.UpdateAccount(account)
			if err != nil {
				return diag.Errorf("failed to update Azure GOV Aviatrix Account: %v", err)
			}
		}
	} else if goaviatrix.IsCloudType(account.CloudType, goaviatrix.AWSChina) {
		if d.HasChanges("awschina_iam", "awschina_role_app", "awschina_role_ec2", "awschina_access_key", "awschina_secret_key") {
			err := client.UpdateAccount(account)
			if err != nil {
				return diag.Errorf("failed to update AWSChina Aviatrix Account: %v", err)
			}
		}
	} else if goaviatrix.IsCloudType(account.CloudType, goaviatrix.AzureChina) {
		if d.HasChanges("azurechina_subscription_id", "azurechina_directory_id", "azurechina_application_id", "azurechina_application_key") {
			err := client.UpdateAccount(account)
			if err != nil {
				return diag.Errorf("failed to update AzureChina Aviatrix Account: %v", err)
			}
		}
	} else if account.CloudType == goaviatrix.AliCloud {
		if d.HasChange("alicloud_account_id") || d.HasChange("alicloud_access_key") || d.HasChange("alicloud_secret_key") {
			err := client.UpdateAccount(account)
			if err != nil {
				return diag.Errorf("failed to update Aviatrix Account: %s", err)
			}
		}
	}

	d.Partial(false)
	return resourceAviatrixAccountRead(ctx, d, meta)
}

//for now, deleting gcp account will not delete the credential file
func resourceAviatrixAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)
	account := &goaviatrix.Account{
		AccountName: d.Get("account_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix account: %#v", account)

	err := client.DeleteAccount(account)
	if err != nil {
		return diag.Errorf("failed to delete Aviatrix Account: %s", err)
	}

	return nil
}

// Validate account number string is 12 digits
func validateAwsAccountNumber(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	if len(v) != 12 {
		errs = append(errs, fmt.Errorf("%q must be 12 digits, got: %s", key, val))
	} else {
		for _, r := range v {
			if r-'0' < 0 || r-'0' > 9 {
				errs = append(errs, fmt.Errorf("%q must be 12 digits, got: %s", key, val))
				break
			}
		}
	}
	return
}
