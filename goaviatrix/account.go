package goaviatrix

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Account struct {
	CID                                   string   `form:"CID,omitempty"`
	Action                                string   `form:"action,omitempty"`
	AccountName                           string   `form:"account_name,omitempty" json:"account_name,omitempty"`
	CloudType                             int      `form:"cloud_type,omitempty" json:"cloud_type,omitempty"`
	AwsAccountNumber                      string   `form:"aws_account_number,omitempty" json:"account_number,omitempty"`
	AwsIam                                string   `form:"aws_iam,omitempty" json:"aws_iam,omitempty"`
	AwsAccessKey                          string   `form:"aws_access_key,omitempty" json:"account_access_key,omitempty"`
	AwsSecretKey                          string   `form:"aws_secret_key,omitempty" json:"account_secret_access_key,omitempty"`
	AwsRoleApp                            string   `form:"aws_role_arn,omitempty" json:"aws_role_arn,omitempty"`
	AwsRoleEc2                            string   `form:"aws_role_ec2,omitempty" json:"aws_role_ec2,omitempty"`
	AwsGatewayRoleApp                     string   `form:"aws_gateway_role_app,omitempty" json:"aws_gateway_role_app,omitempty"`
	AwsGatewayRoleEc2                     string   `form:"aws_gateway_role_ec2,omitempty" json:"aws_gateway_role_ec2,omitempty"`
	AzureSubscriptionId                   string   `form:"azure_subscription_id,omitempty" json:"azure_subscription_id,omitempty"`
	ArmSubscriptionId                     string   `form:"arm_subscription_id,omitempty" json:"arm_subscription_id,omitempty"`
	ArmApplicationEndpoint                string   `form:"arm_application_endpoint,omitempty" json:"arm_ad_tenant_id,omitempty"`
	ArmApplicationClientId                string   `form:"arm_application_client_id,omitempty" json:"arm_ad_client_id,omitempty"`
	ArmApplicationClientSecret            string   `form:"arm_application_client_secret,omitempty" json:"arm_ad_client_secret,omitempty"`
	AwsgovAccountNumber                   string   `form:"awsgov_account_number,omitempty" json:"awsgovcloud_account_number,omitempty"`
	AwsgovIam                             string   `form:"awsgov_iam,omitempty"`
	AwsgovRoleApp                         string   `form:"awsgov_role_arn,omitempty" json:"aws_gov_aws_role_arn,omitempty"`
	AwsgovRoleEc2                         string   `form:"awsgov_role_ec2,omitempty" json:"aws_gov_aws_role_ec2,omitempty"`
	AwsgovAccessKey                       string   `form:"awsgov_access_key,omitempty" json:"awsgovcloud_access_key,omitempty"`
	AwsgovSecretKey                       string   `form:"awsgov_secret_key,omitempty" json:"awsgovcloud_secret_key,omitempty"`
	AwsgovCloudtrailBucket                string   `form:"awsgov_cloudtrail_bucket,omitempty" json:"awsgov_cloudtrail_bucket,omitempty"`
	ProjectCredentialsFilename            string   `form:"filename,omitempty"` // Applies for both GCP and OCI
	ProjectCredentialsContents            string   `form:"contents,omitempty"` // Applies for both GCP and OCI
	GcloudProjectCredentialsFilepathLocal string   `form:"gcloud_project_credentials_local,omitempty"`
	GcloudProjectName                     string   `form:"gcloud_project_name,omitempty" json:"project,omitempty"`
	OciTenancyID                          string   `form:"oci_tenancy_id" json:"oci_tenancy_id,omitempty"`
	OciUserID                             string   `form:"oci_user_id" json:"oci_user_id,omitempty"`
	OciCompartmentID                      string   `form:"oci_compartment_id" json:"oci_compartment_id,omitempty"`
	OciApiPrivateKeyFilePath              string   `form:"oci_api_key_path" json:"oci_api_private_key_filepath,omitempty"`
	AzuregovSubscriptionId                string   `form:"azure_gov_subscription_id,omitempty" json:"arm_gov_subscription_id,omitempty"`
	AzuregovApplicationEndpoint           string   `form:"azure_gov_application_endpoint,omitempty" json:"arm_gov_ad_tenant_id,omitempty"`
	AzuregovApplicationClientId           string   `form:"azure_gov_application_client_id,omitempty" json:"arm_gov_ad_client_id,omitempty"`
	AzuregovApplicationClientSecret       string   `form:"azure_gov_application_client_secret,omitempty" json:"azure_gov_application_client_secret,omitempty"`
	AlicloudAccountId                     string   `form:"aliyun_account_id,omitempty"`
	AlicloudAccessKey                     string   `form:"aliyun_access_key,omitempty"`
	AlicloudSecretKey                     string   `form:"aliyun_secret_key,omitempty"`
	AwsChinaAccountNumber                 string   `form:"aws_china_account_number,omitempty" json:"aws_china_account_number,omitempty"`
	AwsChinaIam                           string   `form:"aws_china_iam,omitempty"`
	AwsChinaRoleApp                       string   `form:"aws_china_role_arn,omitempty" json:"aws_china_aws_role_arn,omitempty"`
	AwsChinaRoleEc2                       string   `form:"aws_china_role_ec2,omitempty" json:"aws_china_aws_role_ec2,omitempty"`
	AwsChinaAccessKey                     string   `form:"aws_china_access_key,omitempty" json:"aws_china_access_key,omitempty"`
	AwsChinaSecretKey                     string   `form:"aws_china_secret_key,omitempty"`
	AzureChinaSubscriptionId              string   `form:"arm_china_subscription_id,omitempty" json:"arm_china_subscription_id,omitempty"`
	AzureChinaApplicationEndpoint         string   `form:"arm_china_application_endpoint,omitempty"`
	AzureChinaApplicationClientId         string   `form:"arm_china_application_client_id,omitempty"`
	AzureChinaApplicationClientSecret     string   `form:"arm_china_application_client_secret,omitempty"`
	GroupNames                            string   `form:"groups,omitempty"`
	GroupNamesRead                        []string `json:"rbac_groups,omitempty"`
	EdgeCSPUsername                       string   `json:"edge_csp_username"`
	EdgeCSPApiEndpoint                    string   `json:"edge_csp_api_endpoint,omitempty"`
	EdgeEquinixUsername                   string   `json:"equinix_username"`
}

type EdgeAccount struct {
	CID                 string `json:"CID,omitempty"`
	Action              string `json:"action,omitempty"`
	AccountName         string `json:"account_name,omitempty"`
	CloudType           int    `json:"cloud_type,omitempty"`
	EdgeCSPUsername     string `json:"edge_csp_username,omitempty"`
	EdgeCSPPassword     string `json:"edge_csp_password,omitempty"`
	EdgeCSPApiEndpoint  string `json:"edge_csp_api_endpoint,omitempty"`
	EdgeEquinixUsername string `json:"equinix_username,omitempty"`
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
	return c.PostAPI(account.Action, account, DuplicateBasicCheck)
}

func (c *Client) CreateGCPAccount(account *Account) error {
	params := map[string]string{
		"CID":                 c.CID,
		"action":              "setup_account_profile",
		"account_name":        account.AccountName,
		"cloud_type":          strconv.Itoa(account.CloudType),
		"gcloud_project_name": account.GcloudProjectName,
		"groups":              account.GroupNames,
	}

	files := []File{
		{
			Path:      account.GcloudProjectCredentialsFilepathLocal,
			ParamName: "gcloud_project_credentials",
		},
	}

	return c.PostFileAPI(params, files, DuplicateBasicCheck)
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
		"groups":             account.GroupNames,
	}

	files := []File{
		{
			Path:      account.OciApiPrivateKeyFilePath,
			ParamName: "oci_api_key",
		},
	}

	return c.PostFileAPI(params, files, DuplicateBasicCheck)
}

func (c *Client) InvalidateCache() {
	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()
	c.cachedAccounts = nil
}

func (c *Client) ListAccounts() ([]Account, error) {
	// If we have cached accounts, return them.
	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()
	if c.cachedAccounts != nil {
		return c.cachedAccounts, nil
	}

	// Otherwise, fetch from the backend
	form := map[string]string{
		"CID":    c.CID,
		"action": "list_accounts",
	}
	var resp AccountListResp
	err := c.GetAPI(&resp, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}
	c.cachedAccounts = resp.Results.AccountList

	return c.cachedAccounts, nil
}

// GetAccount returns the account from the cache. We return an account object
// instead of a pointer to ensure that the caller receives a copy of the
// account data rather than a reference to the cached object. This helps avoid
// potential issues with unintended modifications to the cached account data by
// the caller.
func (c *Client) GetAccount(account *Account) (Account, error) {
	accList, err := c.ListAccounts()
	if err != nil {
		return Account{}, err
	}
	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()
	for i := range accList {
		if accList[i].AccountName == account.AccountName {
			log.Infof("Found Aviatrix Account %s", account.AccountName)
			return accList[i], nil
		}
	}
	log.Errorf("Couldn't find Aviatrix account %s", account.AccountName)
	return Account{}, ErrNotFound
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

func (c *Client) CreateEdgeAccount(edgeAccount *EdgeAccount) error {
	edgeAccount.CID = c.CID
	edgeAccount.Action = "setup_account_profile"
	if edgeAccount.CloudType == EDGEEQUINIX {
		edgeAccount.EdgeEquinixUsername = "no-reply@aviatrix.com"
	}

	return c.PostAPIContext2(context.Background(), nil, edgeAccount.Action, edgeAccount, DuplicateBasicCheck)
}

func (c *Client) UpdateEdgeAccount(edgeAccount *EdgeAccount) error {
	edgeAccount.CID = c.CID
	edgeAccount.Action = "edit_account_profile"
	return c.PostAPIContext2(context.Background(), nil, edgeAccount.Action, edgeAccount, BasicCheck)
}
