package goaviatrix

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Account struct {
	CID                                   string `form:"CID,omitempty"`
	Action                                string `form:"action,omitempty"`
	AccountName                           string `form:"account_name,omitempty" json:"account_name,omitempty"`
	CloudType                             int    `form:"cloud_type,omitempty" json:"cloud_type,omitempty"`
	AwsAccountNumber                      string `form:"aws_account_number,omitempty" json:"account_number,omitempty"`
	AwsIam                                string `form:"aws_iam,omitempty" json:"aws_iam,omitempty"`
	AwsAccessKey                          string `form:"aws_access_key,omitempty" json:"account_access_key,omitempty"`
	AwsSecretKey                          string `form:"aws_secret_key,omitempty" json:"account_secret_access_key,omitempty"`
	AwsRoleApp                            string `form:"aws_role_arn,omitempty" json:"aws_role_arn,omitempty"`
	AwsRoleEc2                            string `form:"aws_role_ec2,omitempty" json:"aws_role_ec2,omitempty"`
	AwsGatewayRoleApp                     string `form:"aws_gateway_role_app,omitempty" json:"aws_gateway_role_app,omitempty"`
	AwsGatewayRoleEc2                     string `form:"aws_gateway_role_ec2,omitempty" json:"aws_gateway_role_ec2,omitempty"`
	AzureSubscriptionId                   string `form:"azure_subscription_id,omitempty" json:"azure_subscription_id,omitempty"`
	ArmSubscriptionId                     string `form:"arm_subscription_id,omitempty" json:"arm_subscription_id,omitempty"`
	ArmApplicationEndpoint                string `form:"arm_application_endpoint,omitempty" json:"arm_ad_tenant_id,omitempty"`
	ArmApplicationClientId                string `form:"arm_application_client_id,omitempty" json:"arm_ad_client_id,omitempty"`
	ArmApplicationClientSecret            string `form:"arm_application_client_secret,omitempty" json:"arm_ad_client_secret,omitempty"`
	AwsgovAccountNumber                   string `form:"awsgov_account_number,omitempty" json:"awsgovcloud_account_number,omitempty"`
	AwsgovIam                             string `form:"awsgov_iam,omitempty"`
	AwsgovRoleApp                         string `form:"awsgov_role_arn,omitempty" json:"aws_gov_aws_role_arn,omitempty"`
	AwsgovRoleEc2                         string `form:"awsgov_role_ec2,omitempty" json:"aws_gov_aws_role_ec2,omitempty"`
	AwsgovAccessKey                       string `form:"awsgov_access_key,omitempty" json:"awsgovcloud_access_key,omitempty"`
	AwsgovSecretKey                       string `form:"awsgov_secret_key,omitempty" json:"awsgovcloud_secret_key,omitempty"`
	AwsgovCloudtrailBucket                string `form:"awsgov_cloudtrail_bucket,omitempty" json:"awsgov_cloudtrail_bucket,omitempty"`
	ProjectCredentialsFilename            string `form:"filename,omitempty"` //Applies for both GCP and OCI
	ProjectCredentialsContents            string `form:"contents,omitempty"` //Applies for both GCP and OCI
	GcloudProjectCredentialsFilepathLocal string `form:"gcloud_project_credentials_local,omitempty"`
	GcloudProjectName                     string `form:"gcloud_project_name,omitempty" json:"project,omitempty"`
	OciTenancyID                          string `form:"oci_tenancy_id" json:"oci_tenancy_id,omitempty"`
	OciUserID                             string `form:"oci_user_id" json:"oci_user_id,omitempty"`
	OciCompartmentID                      string `form:"oci_compartment_id" json:"oci_compartment_id,omitempty"`
	OciApiPrivateKeyFilePath              string `form:"oci_api_key_path" json:"oci_api_private_key_filepath,omitempty"`
	AzuregovSubscriptionId                string `form:"azure_gov_subscription_id,omitempty" json:"arm_gov_subscription_id,omitempty"`
	AzuregovApplicationEndpoint           string `form:"azure_gov_application_endpoint,omitempty" json:"arm_gov_ad_tenant_id,omitempty"`
	AzuregovApplicationClientId           string `form:"azure_gov_application_client_id,omitempty" json:"arm_gov_ad_client_id,omitempty"`
	AzuregovApplicationClientSecret       string `form:"azure_gov_application_client_secret,omitempty" json:"azure_gov_application_client_secret,omitempty"`
	AlicloudAccountId                     string `form:"aliyun_account_id,omitempty"`
	AlicloudAccessKey                     string `form:"aliyun_access_key,omitempty"`
	AlicloudSecretKey                     string `form:"aliyun_secret_key,omitempty"`
	AwsChinaAccountNumber                 string `form:"aws_china_account_number,omitempty" json:"aws_china_account_number,omitempty"`
	AwsChinaIam                           string `form:"aws_china_iam,omitempty"`
	AwsChinaRoleApp                       string `form:"aws_china_role_arn,omitempty" json:"aws_china_aws_role_arn,omitempty"`
	AwsChinaRoleEc2                       string `form:"aws_china_role_ec2,omitempty" json:"aws_china_aws_role_ec2,omitempty"`
	AwsChinaAccessKey                     string `form:"aws_china_access_key,omitempty" json:"aws_china_access_key,omitempty"`
	AwsChinaSecretKey                     string `form:"aws_china_secret_key,omitempty"`
	AzureChinaSubscriptionId              string `form:"arm_china_subscription_id,omitempty" json:"arm_china_subscription_id,omitempty"`
	AzureChinaApplicationEndpoint         string `form:"arm_china_application_endpoint,omitempty"`
	AzureChinaApplicationClientId         string `form:"arm_china_application_client_id,omitempty"`
	AzureChinaApplicationClientSecret     string `form:"arm_china_application_client_secret,omitempty"`
	AwsTsAccountNumber                    string `json:"awsorangecloud_account_number,omitempty"`
	AwsTsCapUrl                           string `json:"aws_orange_cap_url,omitempty"`
	AwsTsCapAgency                        string `json:"aws_orange_cap_agency,omitempty"`
	AwsTsCapMission                       string `json:"aws_orange_cap_mission,omitempty"`
	AwsTsCapRoleName                      string `json:"aws_orange_cap_role_name,omitempty"`
	AwsTsCapCert                          string
	AwsTsCapCertKey                       string
	AwsTsCaChainCert                      string
	AwsTsCapCertPath                      string `json:"aws_orange_cap_cert_path,omitempty"`
	AwsTsCapCertKeyPath                   string `json:"aws_orange_cap_key_path,omitempty"`
	AwsCaCertPath                         string `json:"aws_ca_cert_path,omitempty"`
	AwsSAccountNumber                     string `form:"aws_red_account_number,omitempty" json:"awsredcloud_account_number,omitempty"`
	AwsSCapUrl                            string `form:"aws_red_cap_url,omitempty" json:"aws_red_cap_url,omitempty"`
	AwsSCapAgency                         string `form:"aws_red_cap_agency,omitempty" json:"aws_red_cap_agency,omitempty"`
	AwsSCapAccountName                    string `form:"aws_red_cap_account_name,omitempty" json:"aws_red_cap_account_name,omitempty"`
	AwsSCapRoleName                       string `json:"aws_red_cap_role_name,omitempty"`
	AwsSCapCert                           string
	AwsSCapCertKey                        string
	AwsSCaChainCert                       string
	AwsSCapCertPath                       string `form:"aws_red_cap_cert_path,omitempty" json:"aws_red_cap_cert_path,omitempty"`
	AwsSCapCertKeyPath                    string `json:"aws_red_cap_key_path,omitempty"`
}

type AccountResult struct {
	AccountList []Account `json:"account_list"`
}

type AccountListResp struct {
	Return  bool          `json:"return"`
	Results AccountResult `json:"results"`
	Reason  string        `json:"reason"`
}

func (c *Client) CreateAccount(account *Account) error {
	account.CID = c.CID
	account.Action = "setup_account_profile"
	return c.PostAPI(account.Action, account, BasicCheck)
}

func (c *Client) CreateGCPAccount(account *Account) error {
	params := map[string]string{
		"CID":                 c.CID,
		"action":              "setup_account_profile",
		"account_name":        account.AccountName,
		"cloud_type":          strconv.Itoa(account.CloudType),
		"gcloud_project_name": account.GcloudProjectName,
	}

	files := []File{
		{
			Path:      account.GcloudProjectCredentialsFilepathLocal,
			ParamName: "gcloud_project_credentials",
		},
	}

	return c.PostFileAPI(params, files, BasicCheck)
}

func (c *Client) CreateOCIAccount(account *Account) error {
	params := map[string]string{
		"CID":                c.CID,
		"action":             "setup_account_profile",
		"account_name":       account.AccountName,
		"cloud_type":         strconv.Itoa(account.CloudType),
		"oci_tenancy_id":     account.OciTenancyID,
		"oci_user_id":        account.OciUserID,
		"oci_compartment_id": account.OciCompartmentID,
	}

	files := []File{
		{
			Path:      account.OciApiPrivateKeyFilePath,
			ParamName: "oci_api_key",
		},
	}

	return c.PostFileAPI(params, files, BasicCheck)
}

func (c *Client) CreateAWSTSAccount(account *Account) error {
	params := map[string]string{
		"CID":                       c.CID,
		"action":                    "setup_account_profile",
		"account_name":              account.AccountName,
		"cloud_type":                strconv.Itoa(account.CloudType),
		"aws_orange_account_number": account.AwsTsAccountNumber,
		"aws_orange_cap_url":        account.AwsTsCapUrl,
		"aws_orange_cap_agency":     account.AwsTsCapAgency,
		"aws_orange_cap_mission":    account.AwsTsCapMission,
		"aws_orange_cap_role_name":  account.AwsTsCapRoleName,
	}

	files := []File{
		{
			Path:      account.AwsTsCapCert,
			ParamName: "aws_orange_cap_cert",
		},
		{
			Path:      account.AwsTsCapCertKey,
			ParamName: "aws_orange_cap_cert_key",
		},
		{
			Path:      account.AwsTsCaChainCert,
			ParamName: "aws_orange_ca_chain_cert",
		},
	}

	return c.PostFileAPI(params, files, BasicCheck)
}

func (c *Client) CreateAWSSAccount(account *Account) error {
	params := map[string]string{
		"CID":                      c.CID,
		"action":                   "setup_account_profile",
		"account_name":             account.AccountName,
		"cloud_type":               strconv.Itoa(account.CloudType),
		"aws_red_account_number":   account.AwsSAccountNumber,
		"aws_red_cap_url":          account.AwsSCapUrl,
		"aws_red_cap_agency":       account.AwsSCapAgency,
		"aws_red_cap_account_name": account.AwsSCapAccountName,
		"aws_red_cap_role_name":    account.AwsSCapRoleName,
	}

	files := []File{
		{
			Path:      account.AwsSCapCert,
			ParamName: "aws_red_cap_cert",
		},
		{
			Path:      account.AwsSCapCertKey,
			ParamName: "aws_red_cap_cert_key",
		},
		{
			Path:      account.AwsSCaChainCert,
			ParamName: "aws_red_ca_chain_cert",
		},
	}

	return c.PostFileAPI(params, files, BasicCheck)
}

func (c *Client) GetAccount(account *Account) (*Account, error) {
	form := map[string]string{
		"CID":    c.CID,
		"action": "list_accounts",
	}

	var resp AccountListResp
	err := c.GetAPI(&resp, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}
	accList := resp.Results.AccountList
	for i := range accList {
		if accList[i].AccountName == account.AccountName {
			log.Infof("Found Aviatrix Account %s", account.AccountName)
			return &accList[i], nil
		}
	}
	log.Errorf("Couldn't find Aviatrix account %s", account.AccountName)
	return nil, ErrNotFound
}

func (c *Client) UpdateAccount(account *Account) error {
	account.CID = c.CID
	account.Action = "edit_account_profile"
	return c.PostAPI(account.Action, account, BasicCheck)
}

func (c *Client) UpdateGCPAccount(account *Account) error {
	params := map[string]string{
		"CID":                 c.CID,
		"action":              "edit_account_profile",
		"account_name":        account.AccountName,
		"cloud_type":          strconv.Itoa(account.CloudType),
		"gcloud_project_name": account.GcloudProjectName,
	}

	files := []File{
		{
			Path:      account.GcloudProjectCredentialsFilepathLocal,
			ParamName: "gcloud_project_credentials",
		},
	}

	return c.PostFileAPI(params, files, BasicCheck)
}

func (c *Client) UpdateAWSTSAccount(account *Account, fileChanges map[string]bool) error {
	params := map[string]string{
		"CID":                       c.CID,
		"action":                    "edit_account_profile",
		"account_name":              account.AccountName,
		"cloud_type":                strconv.Itoa(account.CloudType),
		"aws_orange_account_number": account.AwsTsAccountNumber,
		"aws_orange_cap_url":        account.AwsTsCapUrl,
		"aws_orange_cap_agency":     account.AwsTsCapAgency,
		"aws_orange_cap_mission":    account.AwsTsCapMission,
		"aws_orange_cap_role_name":  account.AwsTsCapRoleName,
	}

	files := make([]File, 0, 3)
	if fileChanges["aws_orange_cap_cert"] {
		files = append(files, File{
			Path:      account.AwsTsCapCert,
			ParamName: "aws_orange_cap_cert",
		})
	} else {
		params["aws_orange_cap_cert_path"] = account.AwsTsCapCertPath
	}

	if fileChanges["aws_orange_cap_cert_key"] {
		files = append(files, File{
			Path:      account.AwsTsCapCertKey,
			ParamName: "aws_orange_cap_cert_key",
		})
	} else {
		params["aws_orange_cap_key_path"] = account.AwsTsCapCertKeyPath
	}

	if fileChanges["aws_orange_ca_chain_cert"] {
		files = append(files, File{
			Path:      account.AwsTsCaChainCert,
			ParamName: "aws_orange_ca_chain_cert",
		})
	} else {
		params["aws_orange_ca_cert_path"] = account.AwsCaCertPath
	}

	return c.PostFileAPI(params, files, BasicCheck)
}

func (c *Client) UpdateAWSSAccount(account *Account, fileChanges map[string]bool) error {
	params := map[string]string{
		"CID":                      c.CID,
		"action":                   "edit_account_profile",
		"account_name":             account.AccountName,
		"cloud_type":               strconv.Itoa(account.CloudType),
		"aws_red_account_number":   account.AwsSAccountNumber,
		"aws_red_cap_url":          account.AwsSCapUrl,
		"aws_red_cap_agency":       account.AwsSCapAgency,
		"aws_red_cap_account_name": account.AwsSCapAccountName,
		"aws_red_cap_role_name":    account.AwsSCapRoleName,
	}

	files := make([]File, 0, 3)
	if fileChanges["aws_red_cap_cert"] {
		files = append(files, File{
			Path:      account.AwsSCapCert,
			ParamName: "aws_red_cap_cert",
		})
	} else {
		params["aws_red_cap_cert_path"] = account.AwsSCapCertPath
	}

	if fileChanges["aws_red_cap_cert_key"] {
		files = append(files, File{
			Path:      account.AwsSCapCertKey,
			ParamName: "aws_red_cap_cert_key",
		})
	} else {
		params["aws_red_cap_key_path"] = account.AwsSCapCertKeyPath
	}

	if fileChanges["aws_red_ca_chain_cert"] {
		files = append(files, File{
			Path:      account.AwsSCaChainCert,
			ParamName: "aws_red_ca_chain_cert",
		})
	} else {
		params["aws_red_ca_cert_path"] = account.AwsCaCertPath
	}

	return c.PostFileAPI(params, files, BasicCheck)
}

func (c *Client) DeleteAccount(account *Account) error {
	account.CID = c.CID
	account.Action = "delete_account_profile"
	return c.PostAPI(account.Action, account, BasicCheck)
}

func (c *Client) UploadOciApiPrivateKeyFile(account *Account) error {
	account.CID = c.CID
	account.Action = "upload_file"
	return c.PostAPI(account.Action, account, BasicCheck)
}

func (c *Client) AuditAccount(ctx context.Context, account *Account) error {
	form := map[string]string{
		"CID":    c.CID,
		"action": "get_account_audit_records",
	}

	type AccountAuditResult struct {
		AccountName string `json:"account_name"`
		Status      string `json:"status"`
		Comment     string `json:"comment"`
	}

	type AccountAuditResponse struct {
		Return  bool                 `json:"return"`
		Results []AccountAuditResult `json:"results"`
	}

	var resp AccountAuditResponse
	err := c.GetAPIContext(ctx, &resp, form["action"], form, BasicCheck)
	if err != nil {
		return err
	}

	for _, accountAuditResult := range resp.Results {
		if accountAuditResult.AccountName == account.AccountName && !strings.Contains(accountAuditResult.Status, "Pass") {
			return fmt.Errorf("%s", accountAuditResult.Comment)
		}
	}
	return nil
}
