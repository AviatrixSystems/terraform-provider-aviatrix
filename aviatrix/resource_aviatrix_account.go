package aviatrix

import (
	"fmt"
	"log"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixAccountCreate,
		Read:   resourceAviatrixAccountRead,
		Update: resourceAviatrixAccountUpdate,
		Delete: resourceAviatrixAccountDelete,
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
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Account number to associate with Aviatrix account. Should be 12 digits.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
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
				},
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
				Description: "AWS Access Key.",
			},
			"aws_secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "AWS Secret Key.",
			},
			"awsgov_account_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Gov Account number to associate with Aviatrix account.",
			},
			"awsgov_access_key": {
				Type:        schema.TypeString,
				Optional:    true,
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
			"azure_gov_subscription_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Azure Gov Subscription ID.",
			},
			"azure_gov_directory_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Azure Gov Directory ID.",
			},
			"azure_gov_application_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Azure Gov Application ID.",
			},
			"azure_gov_application_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Azure Gov Application Key.",
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
		},
	}
}

func resourceAviatrixAccountCreate(d *schema.ResourceData, meta interface{}) error {
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
		AzureGovSubscriptionId:                d.Get("azure_gov_subscription_id").(string),
		AzureGovApplicationEndpoint:           d.Get("azure_gov_directory_id").(string),
		AzureGovApplicationClientId:           d.Get("azure_gov_application_id").(string),
		AzureGovApplicationClientSecret:       d.Get("azure_gov_application_key").(string),
		OciTenancyID:                          d.Get("oci_tenancy_id").(string),
		OciUserID:                             d.Get("oci_user_id").(string),
		OciCompartmentID:                      d.Get("oci_compartment_id").(string),
		OciApiPrivateKeyFilePath:              d.Get("oci_api_private_key_filepath").(string),
	}

	awsIam := d.Get("aws_iam").(bool)
	if awsIam {
		account.AwsIam = "true"
	} else {
		account.AwsIam = "false"
	}

	if account.CloudType == goaviatrix.AWS {
		if account.AwsAccountNumber == "" {
			return fmt.Errorf("aws account number is needed for aws cloud")
		}
		if account.AwsIam != "true" && account.AwsIam != "false" {
			return fmt.Errorf("aws iam can only be 'true' or 'false'")
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

			_, gatewayRoleAppOk := d.GetOk("aws_gateway_role_app")
			_, gatewayRoleEc2Ok := d.GetOk("aws_gateway_role_ec2")
			if gatewayRoleAppOk != gatewayRoleEc2Ok {
				return fmt.Errorf("failed to create Aviatrix account: must provide both gateway app role ARN and gateway ec2 role ARN when using separate IAM role and policy for gateways")
			}
		}
	} else if account.CloudType == goaviatrix.GCP {
		if account.GcloudProjectCredentialsFilepathLocal == "" {
			return fmt.Errorf("gcloud project credentials local filepath needed to upload file to controller")
		}
		log.Printf("[INFO] Creating Aviatrix account: %#v", account)
	} else if account.CloudType == goaviatrix.AZURE {
		if account.ArmSubscriptionId == "" {
			return fmt.Errorf("arm subscription id needed for azure cloud")
		}
		if account.ArmApplicationEndpoint == "" {
			return fmt.Errorf("arm directory id needed for azure cloud")
		}
		if account.ArmApplicationClientId == "" {
			return fmt.Errorf("arm application id needed for azure cloud")
		}
		if account.ArmApplicationClientSecret == "" {
			return fmt.Errorf("arm application key needed for azure cloud")
		}
	} else if account.CloudType == goaviatrix.OCI {
		if account.OciTenancyID == "" {
			return fmt.Errorf("oci tenancy ocid needed for oracle cloud")
		}
		if account.OciUserID == "" {
			return fmt.Errorf("oci user id needed for oracle cloud")
		}
		if account.OciCompartmentID == "" {
			return fmt.Errorf("oci compartment ocid needed for oracle cloud")
		}
		if account.OciApiPrivateKeyFilePath == "" {
			return fmt.Errorf("oci api private key filepath needed to upload file to controller")
		}
	} else if account.CloudType == goaviatrix.AWSGOV {
		if account.AwsgovAccountNumber == "" {
			return fmt.Errorf("aws gov account number needed for aws gov cloud")
		}
		if account.AwsgovAccessKey == "" {
			return fmt.Errorf("aws gov access key needed for aws gov cloud")
		}
		if account.AwsgovSecretKey == "" {
			return fmt.Errorf("aws gov secret key needed for aws gov cloud")
		}
	} else if account.CloudType == goaviatrix.AZUREGOV {
		if account.AzureGovSubscriptionId == "" {
			return fmt.Errorf("azure gov subsription id needed when creating an account for arm gov cloud")
		}
		if account.AzureGovApplicationEndpoint == "" {
			return fmt.Errorf("azure gov directory id needed when creating an account for arm gov cloud")
		}
		if account.AzureGovApplicationClientId == "" {
			return fmt.Errorf("azure gov application id needed when creating an account for arm gov cloud")
		}
		if account.AzureGovApplicationClientSecret == "" {
			return fmt.Errorf("azure gov application key needed when creating an account for arm gov cloud")
		}
	} else {
		return fmt.Errorf("cloud type can only be either aws (1), gcp (4), azure (8), oci (16), azure gov (32) or aws gov (256)")
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
		return fmt.Errorf("failed to create Aviatrix Account: %s", err)
	}

	d.SetId(account.AccountName)
	return resourceAviatrixAccountRead(d, meta)
}

func resourceAviatrixAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	accountName := d.Get("account_name").(string)
	if accountName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no account name received. Import Id is %s", id)
		d.Set("account_name", id)
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
		return fmt.Errorf("aviatrix Account: %s", err)
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
				d.Set("aws_access_key", acc.AwsAccessKey)
				d.Set("aws_iam", false)
			}
		} else if acc.CloudType == goaviatrix.GCP {
			d.Set("gcloud_project_id", acc.GcloudProjectName)
		} else if acc.CloudType == goaviatrix.AZURE {
			d.Set("arm_subscription_id", acc.ArmSubscriptionId)
		} else if acc.CloudType == goaviatrix.AWSGOV {
			d.Set("awsgov_account_number", acc.AwsgovAccountNumber)
			d.Set("awsgov_access_key", acc.AwsgovAccessKey)
		} else if acc.CloudType == goaviatrix.AZUREGOV {
			d.Set("azure_gov_subscription_id", acc.AzureGovSubscriptionId)
		}
		d.SetId(acc.AccountName)
	}

	return nil
}

func resourceAviatrixAccountUpdate(d *schema.ResourceData, meta interface{}) error {
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
		AzureGovSubscriptionId:                d.Get("azure_gov_subscription_id").(string),
		AzureGovApplicationEndpoint:           d.Get("azure_gov_directory_id").(string),
		AzureGovApplicationClientId:           d.Get("azure_gov_application_id").(string),
		AzureGovApplicationClientSecret:       d.Get("azure_gov_application_key").(string),
		OciTenancyID:                          d.Get("oci_tenancy_id").(string),
		OciUserID:                             d.Get("oci_user_id").(string),
		OciCompartmentID:                      d.Get("oci_compartment_id").(string),
		OciApiPrivateKeyFilePath:              d.Get("oci_api_private_key_filepath").(string),
	}

	awsIam := d.Get("aws_iam").(bool)
	if awsIam {
		account.AwsIam = "true"
	} else {
		account.AwsIam = "false"
	}

	log.Printf("[INFO] Updating Aviatrix account: %#v", account)

	d.Partial(true)

	if d.HasChange("cloud_type") {
		return fmt.Errorf("update cloud_type is not allowed")
	}

	if d.HasChange("account_name") {
		return fmt.Errorf("update account name is not allowed")
	}

	if account.CloudType == goaviatrix.AWS {
		if d.HasChanges("aws_account_number", "aws_access_key", "aws_secret_key", "aws_iam", "aws_role_app", "aws_role_ec2", "aws_gateway_role_app", "aws_gateway_role_ec2") {
			_, gatewayRoleAppOk := d.GetOk("aws_gateway_role_app")
			_, gatewayRoleEc2Ok := d.GetOk("aws_gateway_role_ec2")
			if gatewayRoleAppOk != gatewayRoleEc2Ok {
				return fmt.Errorf("failed to update Aviatrix account: must provide both gateway app role ARN and gateway ec2 role ARN when using separate IAM role and policy for gateways")
			}

			err := client.UpdateAccount(account)
			if err != nil {
				return fmt.Errorf("failed to update Aviatrix Account: %s", err)
			}
		}
	} else if account.CloudType == goaviatrix.GCP {
		if d.HasChange("gcloud_project_id") || d.HasChange("gcloud_project_credentials_filepath") {
			err := client.UpdateGCPAccount(account)
			if err != nil {
				return fmt.Errorf("failed to update Aviatrix Account: %s", err)
			}
		}
	} else if account.CloudType == goaviatrix.AZURE {
		if d.HasChange("arm_subscription_id") || d.HasChange("arm_directory_id") || d.HasChange("arm_application_id") || d.HasChange("arm_application_key") {
			err := client.UpdateAccount(account)
			if err != nil {
				return fmt.Errorf("failed to update Aviatrix Account: %s", err)
			}
		}
	} else if account.CloudType == goaviatrix.OCI {
		if d.HasChange("oci_tenancy_id") || d.HasChange("oci_user_id") || d.HasChange("oci_compartment_id") || d.HasChange("oci_api_private_key_filepath") {
			return fmt.Errorf("updating OCI account is not supported")
		}
	} else if account.CloudType == goaviatrix.AWSGOV {
		if d.HasChange("awsgov_account_number") || d.HasChange("awsgov_access_key") || d.HasChange("awsgov_secret_key") {
			err := client.UpdateAccount(account)
			if err != nil {
				return fmt.Errorf("failed to update Aviatrix Account: %s", err)
			}
		}
	} else if account.CloudType == goaviatrix.AZUREGOV {
		if d.HasChanges("azure_gov_subscription_id", "azure_gov_directory_id", "azure_gov_application_id", "azure_gov_application_key") {
			err := client.UpdateAccount(account)
			if err != nil {
				return fmt.Errorf("failed to update Azure GOV Aviatrix Account: %v", err)
			}
		}
	}

	d.Partial(false)
	return resourceAviatrixAccountRead(d, meta)
}

//for now, deleting gcp account will not delete the credential file
func resourceAviatrixAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	account := &goaviatrix.Account{
		AccountName: d.Get("account_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix account: %#v", account)

	err := client.DeleteAccount(account)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Account: %s", err)
	}

	return nil
}
