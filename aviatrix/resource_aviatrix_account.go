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
			"aws_orange_account_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Top Secret Region Account Number.",
			},
			"aws_orange_cap_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Top Secret Region CAP Endpoint URL.",
			},
			"aws_orange_cap_agency": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Top Secret Region CAP Agency.",
			},
			"aws_orange_cap_mission": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Top Secret Region CAP Mission.",
			},
			"aws_orange_cap_role_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Top Secret Region CAP Role Name.",
			},
			"aws_orange_cap_cert": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "AWS Top Secret Region CAP Certificate file path.",
			},
			"aws_orange_cap_cert_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "AWS Top Secret Region CAP Certificate Key file path.",
			},
			"aws_orange_ca_chain_cert": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "AWS Top Secret Region Custom Certificate Authority file path.",
			},
			"aws_red_account_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Secret Region Account Number.",
			},
			"aws_red_cap_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Secret Region CAP Endpoint URL.",
			},
			"aws_red_cap_agency": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Secret Region CAP Agency.",
			},
			"aws_red_cap_account_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Secret Region CAP Account Name.",
			},
			"aws_red_cap_role_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Secret Region CAP Role Name.",
			},
			"aws_red_cap_cert": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "AWS Secret Region CAP Certificate file path.",
			},
			"aws_red_cap_cert_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "AWS Secret Region CAP Certificate Key file path.",
			},
			"aws_red_ca_chain_cert": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "AWS Secret Region Custom Certificate Authority file path.",
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
			"aws_orange_cap_cert_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS Top Secret Region CAP Certificate file path on the controller.",
			},
			"aws_orange_cap_cert_key_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS Top Secret Region CAP Certificate Key file path on the controller.",
			},
			"aws_ca_cert_path": { //TODO
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS Top Secret Region or Secret Region Custom Certificate Authority file path on the controller.",
			},
			"aws_red_cap_cert_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS Secret Region CAP Certificate file path on the controller.",
			},
			"aws_red_cap_cert_key_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "AWS Secret Region CAP Certificate Key file path on the controller.",
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
		AwsOrangeAccountNumber:                d.Get("aws_orange_account_number").(string),
		AwsOrangeCapUrl:                       d.Get("aws_orange_cap_url").(string),
		AwsOrangeCapAgency:                    d.Get("aws_orange_cap_agency").(string),
		AwsOrangeCapMission:                   d.Get("aws_orange_cap_mission").(string),
		AwsOrangeCapRoleName:                  d.Get("aws_orange_cap_role_name").(string),
		AwsOrangeCapCert:                      d.Get("aws_orange_cap_cert").(string),
		AwsOrangeCapCertKey:                   d.Get("aws_orange_cap_cert_key").(string),
		AwsOrangeCaChainCert:                  d.Get("aws_orange_ca_chain_cert").(string),
		AwsOrangeCapCertPath:                  d.Get("aws_orange_cap_cert_path").(string),
		AwsOrangeCapCertKeyPath:               d.Get("aws_orange_cap_cert_key_path").(string),
		AwsCaCertPath:                         d.Get("aws_ca_cert_path").(string),
		AwsRedAccountNumber:                   d.Get("aws_red_account_number").(string),
		AwsRedCapUrl:                          d.Get("aws_red_cap_url").(string),
		AwsRedCapAgency:                       d.Get("aws_red_cap_agency").(string),
		AwsRedCapAccountName:                  d.Get("aws_red_cap_account_name").(string),
		AwsRedCapRoleName:                     d.Get("aws_red_cap_role_name").(string),
		AwsRedCapCert:                         d.Get("aws_red_cap_cert").(string),
		AwsRedCapCertKey:                      d.Get("aws_red_cap_cert_key").(string),
		AwsRedCaChainCert:                     d.Get("aws_red_ca_chain_cert").(string),
		AwsRedCapCertPath:                     d.Get("aws_red_cap_cert_path").(string),
		AwsRedCapCertKeyPath:                  d.Get("aws_red_cap_cert_key_path").(string),
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
	} else if account.CloudType == goaviatrix.AWSC2S {
		if account.AwsOrangeAccountNumber == "" {
			return fmt.Errorf("AWS orange account number is needed when creating an account for AWS orange")
		}
		if account.AwsOrangeCapUrl == "" {
			return fmt.Errorf("AWS orange CAP endpoint url is needed when creating an account for AWS orange")
		}
		if account.AwsOrangeCapAgency == "" {
			return fmt.Errorf("AWS orange CAP agency is needed when creating an account for AWS orange")
		}
		if account.AwsOrangeCapMission == "" {
			return fmt.Errorf("AWS orange CAP mission is needed when creating an account for AWS orange")
		}
		if account.AwsOrangeCapRoleName == "" {
			return fmt.Errorf("AWS orange CAP role name is needed when creating an account for AWS orange")
		}
		if account.AwsOrangeCapCert == "" {
			return fmt.Errorf("AWS orange CAP cert file is needed when creating an account for AWS orange")
		}
		if account.AwsOrangeCapCertKey == "" {
			return fmt.Errorf("AWS orange CAP cert key file is needed when creating an account for AWS orange")
		}
		if account.AwsOrangeCaChainCert == "" {
			return fmt.Errorf("AWS orange custom CA Chain file  is needed when creating an account for AWS orange")
		}
	} else if account.CloudType == goaviatrix.AWSSC2S {
		if account.AwsRedAccountNumber == "" {
			return fmt.Errorf("AWS red account number is needed when creating an account for AWS red")
		}
		if account.AwsRedCapUrl == "" {
			return fmt.Errorf("AWS red CAP endpoint url is needed when creating an account for AWS red")
		}
		if account.AwsRedCapAgency == "" {
			return fmt.Errorf("AWS red CAP agency is needed when creating an account for AWS red")
		}
		if account.AwsRedCapAccountName == "" {
			return fmt.Errorf("AWS red CAP Account Name is needed when creating an account for AWS red")
		}
		if account.AwsRedCapRoleName == "" {
			return fmt.Errorf("AWS red CAP role name is needed when creating an account for AWS red")
		}
		if account.AwsRedCapCert == "" {
			return fmt.Errorf("AWS red CAP cert file is needed when creating an account for AWS red")
		}
		if account.AwsRedCapCertKey == "" {
			return fmt.Errorf("AWS red CAP cert key file is needed when creating an account for AWS red")
		}
		if account.AwsRedCaChainCert == "" {
			return fmt.Errorf("AWS red custom CA Chain file  is needed when creating an account for AWS red")
		}
	} else {
		return fmt.Errorf("cloud type can only be either AWS (1), GCP (4), AZURE (8), OCI (16), AZURE Gov (32), AWS GOV (256) or AWS Orange (16384)")
	}

	var err error
	if account.CloudType == goaviatrix.GCP {
		err = client.CreateGCPAccount(account)
	} else if account.CloudType == goaviatrix.OCI {
		err = client.CreateOCIAccount(account)
	} else if account.CloudType == goaviatrix.AWSC2S {
		err = client.CreateAWSC2SAccount(account)
	} else if account.CloudType == goaviatrix.AWSSC2S {
		err = client.CreateAWSSC2SAccount(account)
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
		} else if acc.CloudType == goaviatrix.AWSC2S {
			d.Set("aws_orange_account_number", acc.AwsOrangeAccountNumber)
			d.Set("aws_orange_cap_url", acc.AwsOrangeCapUrl)
			d.Set("aws_orange_cap_agency", acc.AwsOrangeCapAgency)
			d.Set("aws_orange_cap_mission", acc.AwsOrangeCapMission)
			d.Set("aws_orange_cap_role_name", acc.AwsOrangeCapRoleName)

			d.Set("aws_orange_cap_cert_path", acc.AwsOrangeCapCertPath)
			d.Set("aws_orange_cap_cert_key_path", acc.AwsOrangeCapCertKeyPath)
			d.Set("aws_ca_cert_path", acc.AwsCaCertPath)
		} else if acc.CloudType == goaviatrix.AWSSC2S {
			d.Set("aws_red_account_number", acc.AwsRedAccountNumber)
			d.Set("aws_red_cap_url", acc.AwsRedCapUrl)
			d.Set("aws_red_cap_agency", acc.AwsRedCapAgency)
			d.Set("aws_red_cap_account_name", acc.AwsRedCapAccountName)
			d.Set("aws_red_cap_role_name", acc.AwsRedCapRoleName)

			d.Set("aws_red_cap_cert_path", acc.AwsRedCapCertPath)
			d.Set("aws_red_cap_cert_key_path", acc.AwsRedCapCertKeyPath)
			d.Set("aws_ca_cert_path", acc.AwsCaCertPath)
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
		AwsOrangeAccountNumber:                d.Get("aws_orange_account_number").(string),
		AwsOrangeCapUrl:                       d.Get("aws_orange_cap_url").(string),
		AwsOrangeCapAgency:                    d.Get("aws_orange_cap_agency").(string),
		AwsOrangeCapMission:                   d.Get("aws_orange_cap_mission").(string),
		AwsOrangeCapRoleName:                  d.Get("aws_orange_cap_role_name").(string),
		AwsOrangeCapCert:                      d.Get("aws_orange_cap_cert").(string),
		AwsOrangeCapCertKey:                   d.Get("aws_orange_cap_cert_key").(string),
		AwsOrangeCaChainCert:                  d.Get("aws_orange_ca_chain_cert").(string),
		AwsOrangeCapCertPath:                  d.Get("aws_orange_cap_cert_path").(string),
		AwsOrangeCapCertKeyPath:               d.Get("aws_orange_cap_cert_key_path").(string),
		AwsCaCertPath:                         d.Get("aws_ca_cert_path").(string),
		AwsRedAccountNumber:                   d.Get("aws_red_account_number").(string),
		AwsRedCapUrl:                          d.Get("aws_red_cap_url").(string),
		AwsRedCapAgency:                       d.Get("aws_red_cap_agency").(string),
		AwsRedCapAccountName:                  d.Get("aws_red_cap_account_name").(string),
		AwsRedCapRoleName:                     d.Get("aws_red_cap_role_name").(string),
		AwsRedCapCert:                         d.Get("aws_red_cap_cert").(string),
		AwsRedCapCertKey:                      d.Get("aws_red_cap_cert_key").(string),
		AwsRedCaChainCert:                     d.Get("aws_red_ca_chain_cert").(string),
		AwsRedCapCertPath:                     d.Get("aws_red_cap_cert_path").(string),
		AwsRedCapCertKeyPath:                  d.Get("aws_red_cap_cert_key_path").(string),
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
	} else if account.CloudType == goaviatrix.AWSC2S {
		fileChanges := map[string]bool{
			"aws_orange_cap_cert":      d.HasChange("aws_orange_cap_cert") && account.AwsOrangeCapCert != "",
			"aws_orange_cap_cert_key":  d.HasChange("aws_orange_cap_cert_key") && account.AwsOrangeCapCertKey != "",
			"aws_orange_ca_chain_cert": d.HasChange("aws_orange_ca_chain_cert") && account.AwsOrangeCaChainCert != "",
		}
		hasFileChanges := fileChanges["aws_orange_cap_cert"] || fileChanges["aws_orange_cap_cert_key"] || fileChanges["aws_orange_ca_chain_cert"]

		if d.HasChanges("aws_orange_account_number", "aws_orange_cap_url", "aws_orange_cap_agency", "aws_orange_cap_mission", "aws_orange_cap_role_name") || hasFileChanges {
			err := client.UpdateAWSC2SAccount(account, fileChanges)
			if err != nil {
				return fmt.Errorf("failed to update AWS Secret Aviatrix Account: %v", err)
			}
		}
	} else if account.CloudType == goaviatrix.AWSSC2S {
		fileChanges := map[string]bool{
			"aws_red_cap_cert":      d.HasChange("aws_red_cap_cert") && account.AwsRedCapCert != "",
			"aws_red_cap_cert_key":  d.HasChange("aws_red_cap_cert_key") && account.AwsRedCapCertKey != "",
			"aws_red_ca_chain_cert": d.HasChange("aws_red_ca_chain_cert") && account.AwsRedCaChainCert != "",
		}
		hasFileChanges := fileChanges["aws_red_cap_cert"] || fileChanges["aws_red_cap_cert_key"] || fileChanges["aws_red_ca_chain_cert"]

		if d.HasChanges("aws_red_account_number", "aws_red_cap_url", "aws_red_cap_agency", "aws_red_cap_account_name", "aws_red_cap_role_name") || hasFileChanges {
			err := client.UpdateAWSSC2SAccount(account, fileChanges)
			if err != nil {
				return fmt.Errorf("failed to update AWS Top Secret Aviatrix Account: %v", err)
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
