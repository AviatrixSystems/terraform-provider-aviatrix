package aviatrix

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixAccount() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixAccountCreate,
		ReadWithoutTimeout:   resourceAviatrixAccountRead,
		UpdateWithoutTimeout: resourceAviatrixAccountUpdate,
		DeleteWithoutTimeout: resourceAviatrixAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},

		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
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
			"awsgov_iam": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "AWSGov IAM-role based flag",
			},
			"awsgov_role_app": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "AWSGov App role ARN",
			},
			"awsgov_role_ec2": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "AWSGov EC2 role ARN",
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
				Default:     false,
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
			"awsts_account_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Top Secret Region Account Number.",
			},
			"awsts_cap_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Top Secret Region CAP Endpoint URL.",
			},
			"awsts_cap_agency": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Top Secret Region CAP Agency.",
			},
			"awsts_cap_mission": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Top Secret Region CAP Mission.",
			},
			"awsts_cap_role_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Top Secret Region CAP Role Name.",
			},
			"awsts_cap_cert": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "AWS Top Secret Region CAP Certificate file path.",
			},
			"awsts_cap_cert_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "AWS Top Secret Region CAP Certificate Key file path.",
			},
			"awsts_ca_chain_cert": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "AWS Top Secret Region Custom Certificate Authority file path.",
			},
			"awss_account_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Secret Region Account Number.",
			},
			"awss_cap_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Secret Region CAP Endpoint URL.",
			},
			"awss_cap_agency": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Secret Region CAP Agency.",
			},
			"awss_cap_account_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Secret Region CAP Account Name.",
			},
			"awss_cap_role_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Secret Region CAP Role Name.",
			},
			"awss_cap_cert": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "AWS Secret Region CAP Certificate file path.",
			},
			"awss_cap_cert_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "AWS Secret Region CAP Certificate Key file path.",
			},
			"awss_ca_chain_cert": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "AWS Secret Region Custom Certificate Authority file path.",
			},
			"edge_csp_username": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"edge_zededa_username"},
				Description:   "Edge CSP username.",
				Deprecated: "Since V3.1.1+, please use edge_zededa_username instead, edge_csp_username will be " +
					"deprecated in the V3.2.0 release.",
			},
			"edge_csp_password": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"edge_zededa_password"},
				Description:   "Edge CSP password.",
				Deprecated: "Since V3.1.1+, please use edge_zededa_password instead, edge_csp_password will be " +
					"deprecated in the V3.2.0 release.",
			},
			"edge_csp_api_endpoint": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "Edge CSP API endpoint URL",
			},
			"edge_zededa_username": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"edge_csp_username"},
				Description:   "Edge Zededa username.",
			},
			"edge_zededa_password": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"edge_csp_password"},
				Description:   "Edge Zededa password.",
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
			"rbac_groups": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					groupNameOld, _ := d.GetChange("rbac_groups")
					return len(mustSlice(groupNameOld)) != 0
				},
				Description: "List of RBAC permission group names.",
			},
		},
	}
}

func resourceAviatrixAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)
	account := &goaviatrix.Account{
		AccountName:                           getString(d, "account_name"),
		CloudType:                             getInt(d, "cloud_type"),
		AwsAccountNumber:                      getString(d, "aws_account_number"),
		AwsRoleApp:                            getString(d, "aws_role_app"),
		AwsRoleEc2:                            getString(d, "aws_role_ec2"),
		AwsGatewayRoleApp:                     getString(d, "aws_gateway_role_app"),
		AwsGatewayRoleEc2:                     getString(d, "aws_gateway_role_ec2"),
		AwsAccessKey:                          getString(d, "aws_access_key"),
		AwsSecretKey:                          getString(d, "aws_secret_key"),
		AwsgovAccountNumber:                   getString(d, "awsgov_account_number"),
		AwsgovRoleApp:                         getString(d, "awsgov_role_app"),
		AwsgovRoleEc2:                         getString(d, "awsgov_role_ec2"),
		AwsgovAccessKey:                       getString(d, "awsgov_access_key"),
		AwsgovSecretKey:                       getString(d, "awsgov_secret_key"),
		GcloudProjectName:                     getString(d, "gcloud_project_id"),
		GcloudProjectCredentialsFilepathLocal: getString(d, "gcloud_project_credentials_filepath"),
		ArmSubscriptionId:                     getString(d, "arm_subscription_id"),
		ArmApplicationEndpoint:                getString(d, "arm_directory_id"),
		ArmApplicationClientId:                getString(d, "arm_application_id"),
		ArmApplicationClientSecret:            getString(d, "arm_application_key"),
		AzuregovSubscriptionId:                getString(d, "azuregov_subscription_id"),
		AzuregovApplicationEndpoint:           getString(d, "azuregov_directory_id"),
		AzuregovApplicationClientId:           getString(d, "azuregov_application_id"),
		AzuregovApplicationClientSecret:       getString(d, "azuregov_application_key"),
		OciTenancyID:                          getString(d, "oci_tenancy_id"),
		OciUserID:                             getString(d, "oci_user_id"),
		OciCompartmentID:                      getString(d, "oci_compartment_id"),
		OciApiPrivateKeyFilePath:              getString(d, "oci_api_private_key_filepath"),
		AlicloudAccountId:                     getString(d, "alicloud_account_id"),
		AlicloudAccessKey:                     getString(d, "alicloud_access_key"),
		AlicloudSecretKey:                     getString(d, "alicloud_secret_key"),
		AwsChinaAccountNumber:                 getString(d, "awschina_account_number"),
		AwsChinaRoleApp:                       getString(d, "awschina_role_app"),
		AwsChinaRoleEc2:                       getString(d, "awschina_role_ec2"),
		AwsChinaAccessKey:                     getString(d, "awschina_access_key"),
		AwsChinaSecretKey:                     getString(d, "awschina_secret_key"),
		AzureChinaSubscriptionId:              getString(d, "azurechina_subscription_id"),
		AzureChinaApplicationEndpoint:         getString(d, "azurechina_directory_id"),
		AzureChinaApplicationClientId:         getString(d, "azurechina_application_id"),
		AzureChinaApplicationClientSecret:     getString(d, "azurechina_application_key"),
	}

	edgeAccount := &goaviatrix.EdgeAccount{
		AccountName:        getString(d, "account_name"),
		CloudType:          getInt(d, "cloud_type"),
		EdgeCSPUsername:    getString(d, "edge_csp_username"),
		EdgeCSPPassword:    getString(d, "edge_csp_password"),
		EdgeCSPApiEndpoint: getString(d, "edge_csp_api_endpoint"),
	}
	if getString(d, "edge_zededa_username") != "" && getString(d, "edge_zededa_password") != "" {
		edgeAccount.EdgeCSPUsername = getString(d, "edge_zededa_username")
		edgeAccount.EdgeCSPPassword = getString(d, "edge_zededa_password")
	}

	if _, ok := d.GetOk("rbac_groups"); ok {
		account.GroupNames = strings.Join(goaviatrix.ExpandStringList(getList(d, "rbac_groups")), ",")
	}

	awsIam := getBool(d, "aws_iam")
	account.AwsIam = strconv.FormatBool(awsIam)

	awsGovIam := getBool(d, "awsgov_iam")
	account.AwsgovIam = strconv.FormatBool(awsGovIam)

	awsChinaIam := getBool(d, "awschina_iam")
	account.AwsChinaIam = strconv.FormatBool(awsChinaIam)

	_, gatewayRoleAppOk := d.GetOk("aws_gateway_role_app")
	_, gatewayRoleEc2Ok := d.GetOk("aws_gateway_role_ec2")

	if gatewayRoleAppOk || gatewayRoleEc2Ok {
		if !goaviatrix.IsCloudType(account.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			return diag.Errorf("could not create Aviatrix Account: aws_gateway_role_app and aws_gateway_role_ec2 can only be used with AWS (1), AWSGov (256) and AWSChina (1024)")
		}
		if goaviatrix.IsCloudType(account.CloudType, goaviatrix.AWS) && !awsIam {
			return diag.Errorf("could not create Aviatrix Account: aws_gateway_role_app and aws_gateway_role_ec2 can only be used with AWS (1) when awsIam is enabled")
		}
		if goaviatrix.IsCloudType(account.CloudType, goaviatrix.AWSGov) && !awsGovIam {
			return diag.Errorf("could not create Aviatrix Account: aws_gateway_role_app and aws_gateway_role_ec2 can only be used with AWSGov (256) when awsGovIam is enabled")
		}
		if goaviatrix.IsCloudType(account.CloudType, goaviatrix.AWSChina) && !awsChinaIam {
			return diag.Errorf("could not create Aviatrix Account: aws_gateway_role_app and aws_gateway_role_ec2 can only be used with AWSChina (1024) when awsChinaIam is enabled")
		}
		if !(gatewayRoleAppOk && gatewayRoleEc2Ok) {
			return diag.Errorf("could not create Aviatrix Account: must provide both aws_gateway_role_app and aws_gateway_role_ec2 when using separate IAM role and policy for gateways")
		}
	}

	if !goaviatrix.IsCloudType(account.CloudType, goaviatrix.AWS) && (awsIam || account.AwsAccountNumber != "" || account.AwsRoleEc2 != "" || account.AwsRoleApp != "" || account.AwsAccessKey != "" || account.AwsSecretKey != "") {
		return diag.Errorf("could not create Aviatrix Account: 'aws_iam', 'aws_account_number', 'aws_role_app', aws_role_ec2', 'aws_access_key' and 'aws_secret_key' can only be set when 'cloud_type' is AWS (1)")
	}
	if !goaviatrix.IsCloudType(account.CloudType, goaviatrix.GCP) && (account.GcloudProjectName != "" || account.GcloudProjectCredentialsFilepathLocal != "") {
		return diag.Errorf("could not create Aviatrix Account: 'gcloud_project_id' and 'gcloud_project_credentials_filepath' can only be set when 'cloud_type' is GCP (4)")
	}
	if !goaviatrix.IsCloudType(account.CloudType, goaviatrix.Azure) && (account.ArmSubscriptionId != "" || account.ArmApplicationClientId != "" || account.ArmApplicationClientSecret != "" || account.ArmApplicationEndpoint != "") {
		return diag.Errorf("could not create Aviatrix Account: 'arm_subscription_id', 'arm_directory_id', 'arm_application_id' and 'arm_application_key' can only be set when 'cloud_type' is Azure (8)")
	}
	if !goaviatrix.IsCloudType(account.CloudType, goaviatrix.OCI) && (account.OciCompartmentID != "" || account.OciApiPrivateKeyFilePath != "" || account.OciTenancyID != "" || account.OciUserID != "") {
		return diag.Errorf("could not create Aviatrix Account: 'oci_compartment_id', oci_api_private_key_file_path', 'oci_tenancy_id' and 'oci_uesr_id' can only be set when cloud_type is OCI (16)")
	}
	if !goaviatrix.IsCloudType(account.CloudType, goaviatrix.AzureGov) && (account.AzuregovSubscriptionId != "" || account.AzuregovApplicationEndpoint != "" || account.AzuregovApplicationClientId != "" || account.AzuregovApplicationClientSecret != "") {
		return diag.Errorf("could not create Aviatrix Account: 'azuregov_subscription_id', 'azuregov_directory_id', 'azuregov_application_id' and 'azuregov_application_key' can only be set when 'cloud_type' is AzureGov (32)")
	}
	if !goaviatrix.IsCloudType(account.CloudType, goaviatrix.AWSGov) && (awsGovIam || account.AwsgovAccountNumber != "" || account.AwsgovRoleApp != "" || account.AwsgovRoleEc2 != "" || account.AwsgovAccessKey != "" || account.AwsgovSecretKey != "") {
		return diag.Errorf("could not create Aviatrix Account: 'awsgov_iam', 'awsgov_account_number', 'awsgov_role_app', 'awsgov_role_ec2', 'awsgov_access_key' and 'awsgov_secret_key' can only be set when 'cloud_type' is AWSGov (256)")
	}
	if !goaviatrix.IsCloudType(account.CloudType, goaviatrix.AWSChina) && (awsChinaIam || account.AwsChinaAccountNumber != "" || account.AwsChinaRoleApp != "" || account.AwsChinaRoleEc2 != "" || account.AwsChinaAccessKey != "" || account.AwsChinaSecretKey != "") {
		return diag.Errorf("could not create Aviatrix Account: 'awschina_iam', 'awschina_account_number', 'awschina_role_app', 'awschina_role_ec2', 'awschina_access_key' and 'awschina_secret_key' can only be set when 'cloud_type' is AWSChina (1024)")
	}
	if !goaviatrix.IsCloudType(account.CloudType, goaviatrix.AzureChina) && (account.AzureChinaSubscriptionId != "" || account.AzureChinaApplicationEndpoint != "" || account.AzureChinaApplicationClientId != "" || account.AzureChinaApplicationClientSecret != "") {
		return diag.Errorf("could not create Aviatrix Account: `azurechina_subscription_id', 'azurechina_directory_id', 'azurechina_application_id' and 'azurechina_application_key' can only be set when 'cloud_type' is AzureChina (2048)")
	}
	if !goaviatrix.IsCloudType(account.CloudType, goaviatrix.AliCloud) && (account.AlicloudAccountId != "" || account.AlicloudAccessKey != "" || account.AlicloudSecretKey != "") {
		return diag.Errorf("could not create Aviatrix Account: 'aliyun_account_id', 'aliyun_access_key' and 'aliyun_secret_key' can only be set when 'cloud_type' is Alibaba Cloud (8192)")
	}
	if !goaviatrix.IsCloudType(account.CloudType, goaviatrix.EDGECSP) && (getString(d, "edge_csp_username") != "" || getString(d, "edge_csp_password") != "") {
		return diag.Errorf("could not create Aviatrix Account: 'edge_csp_username' and 'edge_csp_password' can only be set when 'cloud_type' is Edge CSP (65536)")
	}
	if !goaviatrix.IsCloudType(account.CloudType, goaviatrix.EDGECSP) && (getString(d, "edge_zededa_username") != "" || getString(d, "edge_zededa_password") != "") {
		return diag.Errorf("could not create Aviatrix Account: 'edge_zededa_username' and 'edge_zededa_password' can only be set when 'cloud_type' is Edge Zededa (65536)")
	}

	if !goaviatrix.IsCloudType(account.CloudType, goaviatrix.EDGENEO) {
		cspEndpoint := getString(d, "edge_csp_api_endpoint")
		if cspEndpoint != "" {
			return diag.Errorf("couldn't create Aviatrix Account: 'edge_csp_api_endpoint' can only be set when 'cloud_type' is Edge NEO/Platform (262144)")
		}
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
			if account.AwsAccessKey != "" || account.AwsSecretKey != "" {
				return diag.Errorf("could not create Aviatrix Account: 'aws_access_key' and 'aws_secret_key' can only be set when 'aws_iam' is false and 'cloud_type' is AWS (1)")
			}

			if _, ok := d.GetOk("aws_role_app"); !ok {
				account.AwsRoleApp = fmt.Sprintf("arn:aws:iam::%s:role/aviatrix-role-app", account.AwsAccountNumber)
			}
			if _, ok := d.GetOk("aws_role_ec2"); !ok {
				account.AwsRoleEc2 = fmt.Sprintf("arn:aws:iam::%s:role/aviatrix-role-ec2", account.AwsAccountNumber)
			}
			log.Printf("[TRACE] Reading Aviatrix account aws_role_app: [%s]", getString(d, "aws_role_app"))
			log.Printf("[TRACE] Reading Aviatrix account aws_role_ec2: [%s]", getString(d, "aws_role_ec2"))
		} else {
			if account.AwsRoleApp != "" || account.AwsRoleEc2 != "" {
				return diag.Errorf("could not create Aviatrix Account: 'aws_role_app' and 'aws_role_ec2' can only be set when 'aws_iam' is true and 'cloud_type' is AWS (1)")
			}
			if account.AwsAccessKey == "" || account.AwsSecretKey == "" {
				return diag.Errorf("could not create Aviatrix Account: 'aws_access_key' and 'aws_secret_key' must be set when 'aws_iam' is false and 'cloud_type' is AWS (1)")
			}
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
		if awsGovIam {
			if account.AwsgovAccessKey != "" || account.AwsgovSecretKey != "" {
				return diag.Errorf("could not create Aviatrix Account: 'awsgov_access_key' and 'awsgov_secret_key' can only be set when 'awsgov_iam' is false and 'cloud_type' is AWSGov (256)")
			}
			if account.AwsgovRoleApp == "" {
				account.AwsgovRoleApp = fmt.Sprintf("arn:aws-us-gov:iam::%s:role/aviatrix-role-app", account.AwsgovAccountNumber)
			}
			if account.AwsgovRoleEc2 == "" {
				account.AwsgovRoleEc2 = fmt.Sprintf("arn:aws-us-gov:iam::%s:role/aviatrix-role-ec2", account.AwsgovAccountNumber)
			}
		} else {
			if account.AwsgovRoleApp != "" || account.AwsgovRoleEc2 != "" {
				return diag.Errorf("could not create Aviatrix Account: 'awsgov_role_app' and 'awsgov_role_ec2' can only be set when 'awsgov_iam' is true and 'cloud_type' is AWSGov (256)")
			}
			if account.AwsgovAccessKey == "" || account.AwsgovSecretKey == "" {
				return diag.Errorf("could not create Aviatrix Account: 'awsgov_access_key' and 'awsgov_secret_key' must be set when 'awsgov_iam' is false and 'cloud_type' is AWSGov (256)")
			}
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
			if account.AwsChinaAccessKey != "" || account.AwsChinaSecretKey != "" {
				return diag.Errorf("could not create Aviatrix Account: 'awschina_access_key' and 'awschina_secret_key' can only be set when 'awschina_iam' is false and 'cloud_type' is AWSChina (1024)")
			}
			if account.AwsChinaRoleApp == "" {
				account.AwsChinaRoleApp = fmt.Sprintf("arn:aws-cn:iam::%s:role/aviatrix-role-app", account.AwsChinaAccountNumber)
			}
			if account.AwsChinaRoleEc2 == "" {
				account.AwsChinaRoleEc2 = fmt.Sprintf("arn:aws-cn:iam::%s:role/aviatrix-role-ec2", account.AwsChinaAccountNumber)
			}
		} else {
			if account.AwsChinaRoleApp != "" || account.AwsChinaRoleEc2 != "" {
				return diag.Errorf("could not create Aviatrix Account: 'awschina_role_app' and 'awschina_role_ec2' can only be set when 'awschina_iam' is true and 'cloud_type' is AWSChina (1024)")
			}
			if account.AwsChinaAccessKey == "" || account.AwsChinaSecretKey == "" {
				return diag.Errorf("could not create Aviatrix Account: 'awschina_access_key' and 'awschina_secret_key' must be set when 'awschina_iam' is false and 'cloud_type' is AWSChina (1024)")
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
	} else if account.CloudType == goaviatrix.EDGECSP {
		//if edgeAccount.EdgeCSPUsername == "" {
		//	return diag.Errorf("edge_csp_username is required to create an Aviatrix account for Edge CSP")
		//}
		//if edgeAccount.EdgeCSPPassword == "" {
		//	return diag.Errorf("edge_csp_password is required to create an Aviatrix account for Edge CSP")
		//}
		if !((getString(d, "edge_csp_username") != "" && getString(d, "edge_csp_password") != "") ||
			(getString(d, "edge_zededa_username") != "" && getString(d, "edge_zededa_password") != "")) {
			return diag.Errorf("edge_csp_username and edge_csp_password are required to create an Aviatrix account for Edge CSP, " +
				"edge_zededa_username and edge_zededa_password are required to create an Aviatrix account for Edge Zededa")
		}
	} else if goaviatrix.IsCloudType(account.CloudType, goaviatrix.EdgeRelatedCloudTypes) {
		log.Print("no check is needed to create an Aviatrix account for Edge Equinix, Edge NEO and Edge Megaport")
	} else {
		return diag.Errorf("cloud type can only be either AWS (1), GCP (4), Azure (8), OCI (16), AzureGov (32), AWSGov (256), AWSChina (1024), AzureChina (2048), Alibaba Cloud (8192), AWS Top Secret (16384), AWS Secret (32768), Edge CSP/Zededa (65536), Edge Equinix (524288), Edge Megaport(1048576) or Edge NEO/Platform (262144)")
	}

	var err error
	if account.CloudType == goaviatrix.GCP {
		err = client.CreateGCPAccount(account)
	} else if account.CloudType == goaviatrix.OCI {
		err = client.CreateOCIAccount(account)
	} else if account.CloudType == goaviatrix.EDGECSP || account.CloudType == goaviatrix.EDGEEQUINIX || account.CloudType == goaviatrix.EDGENEO {
		err = client.CreateEdgeAccount(edgeAccount)
	} else {
		err = client.CreateAccount(account)
	}

	if err != nil {
		return diag.Errorf("failed to create Aviatrix Account: %s", err)
	}

	d.SetId(account.AccountName)
	client.InvalidateCache()
	return resourceAviatrixAccountRead(ctx, d, meta)
}

func resourceAviatrixAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(goaviatrix.ClientInterface)
	if !ok {
		return diag.Errorf("internal error: meta is not a valid goaviatrix.ClientInterface")
	}
	var diags diag.Diagnostics

	accountName := getString(d, "account_name")
	isImport := accountName == ""
	if isImport {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no account name received. Import Id is %s", id)
		mustSet(d, "account_name", id)
		d.SetId(id)
	}

	account := &goaviatrix.Account{
		AccountName: getString(d, "account_name"),
	}

	log.Printf("[INFO] Looking for Aviatrix account: %#v", account)

	acc, err := client.GetAccount(account)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("aviatrix Account: %s", err)
	}
	mustSet(d, "account_name", acc.AccountName)
	mustSet(d, "cloud_type", acc.CloudType)
	if acc.CloudType == goaviatrix.AWS {
		mustSet(d, "aws_account_number", acc.AwsAccountNumber)
		if acc.AwsRoleEc2 != "" {
			mustSet(
				// force default setting and save to .tfstate file
				d, "aws_access_key", "")
			mustSet(d, "aws_secret_key", "")
			mustSet(d, "aws_iam", true)
			mustSet(d, "aws_role_app", acc.AwsRoleApp)
			mustSet(d, "aws_role_ec2", acc.AwsRoleEc2)
			mustSet(d, "aws_gateway_role_app", acc.AwsGatewayRoleApp)
			mustSet(d, "aws_gateway_role_ec2", acc.AwsGatewayRoleEc2)
		} else {
			mustSet(d, "aws_iam", false)
		}
	} else if acc.CloudType == goaviatrix.GCP {
		mustSet(d, "gcloud_project_id", acc.GcloudProjectName)
	} else if acc.CloudType == goaviatrix.Azure {
		mustSet(d, "arm_subscription_id", acc.ArmSubscriptionId)
	} else if acc.CloudType == goaviatrix.AWSGov {
		mustSet(d, "awsgov_account_number", acc.AwsgovAccountNumber)
		if acc.AwsgovRoleEc2 != "" {
			mustSet(d, "awsgov_access_key", "")
			mustSet(d, "awsgov_secret_key", "")
			mustSet(d, "awsgov_iam", true)
			mustSet(d, "awsgov_role_app", acc.AwsgovRoleApp)
			mustSet(d, "awsgov_role_ec2", acc.AwsgovRoleEc2)
			mustSet(d, "aws_gateway_role_app", acc.AwsGatewayRoleApp)
			mustSet(d, "aws_gateway_role_ec2", acc.AwsGatewayRoleEc2)
		} else {
			mustSet(d, "awsgov_iam", false)
		}
	} else if acc.CloudType == goaviatrix.AzureGov {
		mustSet(d, "azuregov_subscription_id", acc.AzuregovSubscriptionId)
	} else if goaviatrix.IsCloudType(acc.CloudType, goaviatrix.AWSChina) {
		mustSet(d, "awschina_account_number", acc.AwsChinaAccountNumber)
		if acc.AwsChinaRoleEc2 != "" {
			mustSet(
				// Force access key and secret key to be empty
				d, "awschina_access_key", "")
			mustSet(d, "awschina_secret_key", "")
			mustSet(d, "awschina_iam", true)
			mustSet(d, "awschina_role_app", acc.AwsChinaRoleApp)
			mustSet(d, "awschina_role_ec2", acc.AwsChinaRoleEc2)
			mustSet(d, "aws_gateway_role_app", acc.AwsGatewayRoleApp)
			mustSet(d, "aws_gateway_role_ec2", acc.AwsGatewayRoleEc2)
		} else {
			mustSet(d, "awschina_iam", false)
		}
	} else if goaviatrix.IsCloudType(acc.CloudType, goaviatrix.AzureChina) {
		mustSet(d, "azurechina_subscription_id", acc.AzureChinaSubscriptionId)
	} else if acc.CloudType == goaviatrix.AliCloud {
		mustSet(d, "alicloud_account_id", acc.AwsAccountNumber)
	} else if acc.CloudType == goaviatrix.EDGECSP {
		if getString(d, "edge_csp_password") != "" {
			mustSet(d, "edge_csp_username", acc.EdgeCSPUsername)
		} else if getString(d, "edge_zededa_password") != "" {
			mustSet(d, "edge_zededa_username", acc.EdgeCSPUsername)
		} else {
			mustSet(
				// let user choose when importing CSP/Zededa account
				d, "edge_csp_username", acc.EdgeCSPUsername)
			mustSet(d, "edge_zededa_username", acc.EdgeCSPUsername)
		}
		err := d.Set("edge_csp_api_endpoint", acc.EdgeCSPApiEndpoint)
		if err != nil {
			return diag.Errorf("Setting CSP API endpoint: %s", err)
		}
	}
	mustSet(d, "rbac_groups", acc.GroupNamesRead)
	d.SetId(acc.AccountName)

	// Don't check account audit during import. In terraform version 0.14.11 or earlier, returning diag.Warning during import
	// will cause it to fail silently. It will not return an error, but the state file will not be updated after.
	if getBool(d, "audit_account") && !isImport {
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
	client := mustClient(meta)
	defer client.InvalidateCache()
	account := &goaviatrix.Account{
		AccountName:                           getString(d, "account_name"),
		CloudType:                             getInt(d, "cloud_type"),
		AwsAccountNumber:                      getString(d, "aws_account_number"),
		AwsRoleApp:                            getString(d, "aws_role_app"),
		AwsRoleEc2:                            getString(d, "aws_role_ec2"),
		AwsGatewayRoleApp:                     getString(d, "aws_gateway_role_app"),
		AwsGatewayRoleEc2:                     getString(d, "aws_gateway_role_ec2"),
		AwsAccessKey:                          getString(d, "aws_access_key"),
		AwsSecretKey:                          getString(d, "aws_secret_key"),
		AwsgovAccountNumber:                   getString(d, "awsgov_account_number"),
		AwsgovRoleApp:                         getString(d, "awsgov_role_app"),
		AwsgovRoleEc2:                         getString(d, "awsgov_role_ec2"),
		AwsgovAccessKey:                       getString(d, "awsgov_access_key"),
		AwsgovSecretKey:                       getString(d, "awsgov_secret_key"),
		GcloudProjectName:                     getString(d, "gcloud_project_id"),
		GcloudProjectCredentialsFilepathLocal: getString(d, "gcloud_project_credentials_filepath"),
		ArmSubscriptionId:                     getString(d, "arm_subscription_id"),
		ArmApplicationEndpoint:                getString(d, "arm_directory_id"),
		ArmApplicationClientId:                getString(d, "arm_application_id"),
		ArmApplicationClientSecret:            getString(d, "arm_application_key"),
		AzuregovSubscriptionId:                getString(d, "azuregov_subscription_id"),
		AzuregovApplicationEndpoint:           getString(d, "azuregov_directory_id"),
		AzuregovApplicationClientId:           getString(d, "azuregov_application_id"),
		AzuregovApplicationClientSecret:       getString(d, "azuregov_application_key"),
		OciTenancyID:                          getString(d, "oci_tenancy_id"),
		OciUserID:                             getString(d, "oci_user_id"),
		OciCompartmentID:                      getString(d, "oci_compartment_id"),
		OciApiPrivateKeyFilePath:              getString(d, "oci_api_private_key_filepath"),
		AlicloudAccountId:                     getString(d, "alicloud_account_id"),
		AlicloudAccessKey:                     getString(d, "alicloud_access_key"),
		AlicloudSecretKey:                     getString(d, "alicloud_secret_key"),
		AwsChinaAccountNumber:                 getString(d, "awschina_account_number"),
		AwsChinaRoleApp:                       getString(d, "awschina_role_app"),
		AwsChinaRoleEc2:                       getString(d, "awschina_role_ec2"),
		AwsChinaAccessKey:                     getString(d, "awschina_access_key"),
		AwsChinaSecretKey:                     getString(d, "awschina_secret_key"),
		AzureChinaSubscriptionId:              getString(d, "azurechina_subscription_id"),
		AzureChinaApplicationEndpoint:         getString(d, "azurechina_directory_id"),
		AzureChinaApplicationClientId:         getString(d, "azurechina_application_id"),
		AzureChinaApplicationClientSecret:     getString(d, "azurechina_application_key"),
	}

	edgeAccount := &goaviatrix.EdgeAccount{
		AccountName:        getString(d, "account_name"),
		CloudType:          getInt(d, "cloud_type"),
		EdgeCSPUsername:    getString(d, "edge_csp_username"),
		EdgeCSPPassword:    getString(d, "edge_csp_password"),
		EdgeCSPApiEndpoint: getString(d, "edge_csp_api_endpoint"),
	}
	if getString(d, "edge_zededa_username") != "" && getString(d, "edge_zededa_password") != "" {
		edgeAccount.EdgeCSPUsername = getString(d, "edge_zededa_username")
		edgeAccount.EdgeCSPPassword = getString(d, "edge_zededa_password")
	}

	awsIam := getBool(d, "aws_iam")
	account.AwsIam = strconv.FormatBool(awsIam)

	awsGovIam := getBool(d, "awsgov_iam")
	account.AwsgovIam = strconv.FormatBool(awsGovIam)

	awsChinaIam := getBool(d, "awschina_iam")
	account.AwsChinaIam = strconv.FormatBool(awsChinaIam)

	log.Printf("[INFO] Updating Aviatrix account: %#v", account)

	d.Partial(true)

	if d.HasChange("cloud_type") {
		return diag.Errorf("update cloud_type is not allowed")
	}

	if d.HasChanges("aws_gateway_role_app", "aws_gateway_role_ec2") {
		_, gatewayRoleAppOk := d.GetOk("aws_gateway_role_app")
		_, gatewayRoleEc2Ok := d.GetOk("aws_gateway_role_ec2")

		if !goaviatrix.IsCloudType(account.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			return diag.Errorf("could not update Aviatrix Account: aws_gateway_role_app and aws_gateway_role_ec2 can only be used with AWS (1), AWSGov (256) and AWSChina (1024)")
		}
		if goaviatrix.IsCloudType(account.CloudType, goaviatrix.AWS) && !awsIam {
			return diag.Errorf("could not update Aviatrix Account: aws_gateway_role_app and aws_gateway_role_ec2 can only be used with AWS (1) when awsIam is enabled")
		}
		if goaviatrix.IsCloudType(account.CloudType, goaviatrix.AWSGov) && !awsGovIam {
			return diag.Errorf("could not update Aviatrix Account: aws_gateway_role_app and aws_gateway_role_ec2 can only be used with AWSGov (256) when awsGovIam is enabled")
		}
		if goaviatrix.IsCloudType(account.CloudType, goaviatrix.AWSChina) && !awsChinaIam {
			return diag.Errorf("could not update Aviatrix Account: aws_gateway_role_app and aws_gateway_role_ec2 can only be used with AWSChina (1024) when awsChinaIam is enabled")
		}
		if gatewayRoleAppOk != gatewayRoleEc2Ok {
			return diag.Errorf("could not update Aviatrix account: must provide both aws_gateway_role_app and aws_gateway_role_ec2 when using separate IAM role and policy for gateways")
		}
	}

	if !goaviatrix.IsCloudType(account.CloudType, goaviatrix.AWS) && (awsIam || account.AwsAccountNumber != "" || account.AwsRoleEc2 != "" || account.AwsRoleApp != "" || account.AwsAccessKey != "" || account.AwsSecretKey != "") {
		return diag.Errorf("could not update Aviatrix Account: 'aws_iam', 'aws_account_number', 'aws_role_app', aws_role_ec2', 'aws_access_key' and 'aws_secret_key' can only be set when 'cloud_type' is AWS (1)")
	}
	if !goaviatrix.IsCloudType(account.CloudType, goaviatrix.GCP) && (account.GcloudProjectName != "" || account.GcloudProjectCredentialsFilepathLocal != "") {
		return diag.Errorf("could not update Aviatrix Account: 'gcloud_project_id' and 'gcloud_project_credentials_filepath' can only be set when 'cloud_type' is GCP (4)")
	}
	if !goaviatrix.IsCloudType(account.CloudType, goaviatrix.Azure) && (account.ArmSubscriptionId != "" || account.ArmApplicationClientId != "" || account.ArmApplicationClientSecret != "" || account.ArmApplicationEndpoint != "") {
		return diag.Errorf("could not update Aviatrix Account: 'arm_subscription_id', 'arm_directory_id', 'arm_application_id' and 'arm_application_key' can only be set when 'cloud_type' is Azure (8)")
	}
	if !goaviatrix.IsCloudType(account.CloudType, goaviatrix.OCI) && (account.OciCompartmentID != "" || account.OciApiPrivateKeyFilePath != "" || account.OciTenancyID != "" || account.OciUserID != "") {
		return diag.Errorf("could not update Aviatrix Account: 'oci_compartment_id', oci_api_private_key_file_path', 'oci_tenancy_id' and 'oci_uesr_id' can only be set when cloud_type is OCI (16)")
	}
	if !goaviatrix.IsCloudType(account.CloudType, goaviatrix.AzureGov) && (account.AzuregovSubscriptionId != "" || account.AzuregovApplicationEndpoint != "" || account.AzuregovApplicationClientId != "" || account.AzuregovApplicationClientSecret != "") {
		return diag.Errorf("could not update Aviatrix Account: 'azuregov_subscription_id', 'azuregov_directory_id', 'azuregov_application_id' and 'azuregov_application_key' can only be set when 'cloud_type' is AzureGov (32)")
	}
	if !goaviatrix.IsCloudType(account.CloudType, goaviatrix.AWSGov) && (awsGovIam || account.AwsgovAccountNumber != "" || account.AwsgovRoleApp != "" || account.AwsgovRoleEc2 != "" || account.AwsgovAccessKey != "" || account.AwsgovSecretKey != "") {
		return diag.Errorf("could not update Aviatrix Account: 'awsgov_iam', 'awsgov_account_number', 'awsgov_role_app', 'awsgov_role_ec2', 'awsgov_access_key' and 'awsgov_secret_key' can only be set when 'cloud_type' is AWSGov (256)")
	}
	if !goaviatrix.IsCloudType(account.CloudType, goaviatrix.AWSChina) && (awsChinaIam || account.AwsChinaAccountNumber != "" || account.AwsChinaRoleApp != "" || account.AwsChinaRoleEc2 != "" || account.AwsChinaAccessKey != "" || account.AwsChinaSecretKey != "") {
		return diag.Errorf("could not update Aviatrix Account: 'awschina_iam', 'awschina_account_number', 'awschina_role_app', 'awschina_role_ec2', 'awschina_access_key' and 'awschina_secret_key' can only be set when 'cloud_type' is AWSChina (1024)")
	}
	if !goaviatrix.IsCloudType(account.CloudType, goaviatrix.AzureChina) && (account.AzureChinaSubscriptionId != "" || account.AzureChinaApplicationEndpoint != "" || account.AzureChinaApplicationClientId != "" || account.AzureChinaApplicationClientSecret != "") {
		return diag.Errorf("could not update Aviatrix Account: `azurechina_subscription_id', 'azurechina_directory_id', 'azurechina_application_id' and 'azurechina_application_key' can only be set when 'cloud_type' is AzureChina (2048)")
	}
	if !goaviatrix.IsCloudType(account.CloudType, goaviatrix.AliCloud) && (account.AlicloudAccountId != "" || account.AlicloudAccessKey != "" || account.AlicloudSecretKey != "") {
		return diag.Errorf("could not update Aviatrix Account: 'aliyun_account_id', 'aliyun_access_key' and 'aliyun_secret_key' can only be set when 'cloud_type' is Alibaba Cloud (8192)")
	}
	if !goaviatrix.IsCloudType(account.CloudType, goaviatrix.EDGECSP) && (getString(d, "edge_csp_username") != "" || getString(d, "edge_csp_password") != "") {
		return diag.Errorf("could not update Aviatrix Account: 'edge_csp_username' and 'edge_csp_password' can only be set when 'cloud_type' is Edge CSP (65536)")
	}
	if !goaviatrix.IsCloudType(account.CloudType, goaviatrix.EDGECSP) && (getString(d, "edge_zededa_username") != "" || getString(d, "edge_zededa_password") != "") {
		return diag.Errorf("could not create Aviatrix Account: 'edge_zededa_username' and 'edge_zededa_password' can only be set when 'cloud_type' is Edge Zededa (65536)")
	}

	if account.CloudType == goaviatrix.AWS {
		if awsIam && (account.AwsAccessKey != "" || account.AwsSecretKey != "") {
			return diag.Errorf("could not update Aviatrix Account: 'aws_access_key' and 'aws_secret_key' can only be set when 'aws_iam' is false and 'cloud_type' is AWS (1)")
		} else if !awsIam && (account.AwsRoleApp != "" || account.AwsRoleEc2 != "") {
			return diag.Errorf("could not update Aviatrix Account: 'aws_role_app' and 'aws_role_ec2' can only be set when 'aws_iam' is true and 'cloud_type' is AWS (1)")
		}

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
		if awsGovIam && (account.AwsgovAccessKey != "" || account.AwsgovSecretKey != "") {
			return diag.Errorf("could not update Aviatrix Account: 'awsgov_access_key' and 'awsgov_secret_key' can only be set when 'awsgov_iam' is false and 'cloud_type' is AWSGov (256)")
		} else if !awsGovIam && (account.AwsgovRoleApp != "" || account.AwsgovRoleEc2 != "") {
			return diag.Errorf("could not update Aviatrix Account: 'awsgov_role_app' and 'awsgov_role_ec2' can only be set when 'awsgov_iam' is true and 'cloud_type' is AWSGov (256)")
		}

		if d.HasChanges("awsgov_account_number", "awsgov_access_key", "awsgov_secret_key", "awsgov_iam", "awsgov_role_app", "awsgov_role_ec2", "aws_gateway_role_app", "aws_gateway_role_ec2") {
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
		if d.HasChanges("awschina_iam", "awschina_role_app", "awschina_role_ec2", "awschina_access_key", "awschina_secret_key", "aws_gateway_role_app", "aws_gateway_role_ec2") {
			err := client.UpdateAccount(account)
			if err != nil {
				return diag.Errorf("failed to update AWSChina Aviatrix Account: %v", err)
			}
		}
	} else if goaviatrix.IsCloudType(account.CloudType, goaviatrix.AzureChina) {
		if awsChinaIam && (account.AwsChinaAccessKey != "" || account.AwsChinaSecretKey != "") {
			return diag.Errorf("could not update Aviatrix Account: 'awschina_access_key' and 'awschina_secret_key' can only be set when 'awschina_iam' is false and 'cloud_type' is AWSChina (1024)")
		} else if !awsChinaIam && (account.AwsChinaRoleApp != "" || account.AwsChinaRoleEc2 != "") {
			return diag.Errorf("could not update Aviatrix Account: 'awschina_role_app' and 'awschina_role_ec2' can only be set when 'awschina_iam' is true and 'cloud_type' is AWSChina (1024)")
		}

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
	} else if account.CloudType == goaviatrix.EDGECSP {
		if d.HasChange("edge_csp_username") || d.HasChange("edge_csp_password") {
			err := client.UpdateEdgeAccount(edgeAccount)
			if err != nil {
				return diag.Errorf("failed to update Edge CSP Account: %s", err)
			}
		}
		if d.HasChange("edge_zededa_username") || d.HasChange("edge_zededa_password") {
			err := client.UpdateEdgeAccount(edgeAccount)
			if err != nil {
				return diag.Errorf("failed to update Edge Zededa Account: %s", err)
			}
		}
	}
	d.Partial(false)
	return resourceAviatrixAccountRead(ctx, d, meta)
}

// for now, deleting gcp account will not delete the credential file
func resourceAviatrixAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(goaviatrix.ClientInterface)
	if !ok {
		return diag.Errorf("internal error: meta is not a valid goaviatrix.ClientInterface")
	}
	account := &goaviatrix.Account{
		AccountName: getString(d, "account_name"),
	}

	log.Printf("[INFO] Deleting Aviatrix account: %#v", account)
	defer client.InvalidateCache()
	err := client.DeleteAccount(account)
	if err != nil {
		return diag.Errorf("failed to delete Aviatrix Account: %s", err)
	}
	return nil
}

// Validate account number string is 12 digits
func validateAwsAccountNumber(val interface{}, key string) (warns []string, errs []error) {
	v := mustString(val)
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
