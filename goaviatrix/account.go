package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
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
	AwsgovAccessKey                       string `form:"awsgov_access_key,omitempty" json:"awsgovcloud_access_key,omitempty"`
	AwsgovSecretKey                       string `form:"awsgov_secret_key,omitempty" json:"awsgovcloud_secret_key,omitempty"`
	AwsgovCloudtrailBucket                string `form:"awsgov_cloudtrail_bucket,omitempty" json:"awsgov_cloudtrail_bucket,omitempty"`
	AzurechinaSubscriptionId              string `form:"azurechina_subscription_id,omitempty" json:"azurechina_subscription_id,omitempty"`
	AwschinaAccountNumber                 string `form:"awschina_account_number,omitempty" json:"awschina_account_number,omitempty"`
	AwschinaAccessKey                     string `form:"awschina_access_key,omitempty" json:"awschinacloud_access_key,omitempty"`
	AwschinaSecretKey                     string `form:"awschina_secret_key,omitempty" json:"awschinacloud_secret_key,omitempty"`
	ArmChinaSubscriptionId                string `form:"arm_china_subscription_id,omitempty" json:"arm_china_subscription_id,omitempty"`
	ArmChinaApplicationEndpoint           string `form:"arm_china_application_endpoint,omitempty" json:"arm_china_application_endpoint,omitempty"`
	ArmChinaApplicationClientId           string `form:"arm_china_application_client_id,omitempty" json:"arm_china_application_client_id,omitempty"`
	ArmChinaApplicationClientSecret       string `form:"arm_china_application_client_secret,omitempty" json:"arm_china_application_client_secret,omitempty"`
	ProjectCredentialsFilename            string `form:"filename,omitempty"` //Applies for both GCP and OCI
	ProjectCredentialsContents            string `form:"contents,omitempty"` //Applies for both GCP and OCI
	GcloudProjectCredentialsFilepathLocal string `form:"gcloud_project_credentials_local,omitempty"`
	GcloudProjectName                     string `form:"gcloud_project_name,omitempty" json:"project,omitempty"`
	OciTenancyID                          string `form:"oci_tenancy_id" json:"oci_tenancy_id,omitempty"`
	OciUserID                             string `form:"oci_user_id" json:"oci_user_id,omitempty"`
	OciCompartmentID                      string `form:"oci_compartment_id" json:"oci_compartment_id,omitempty"`
	OciApiPrivateKeyFilePath              string `form:"oci_api_key_path" json:"oci_api_private_key_filepath,omitempty"`
	AzureGovSubscriptionId                string `form:"azure_gov_subscription_id,omitempty" json:"arm_gov_subscription_id,omitempty"`
	AzureGovApplicationEndpoint           string `form:"azure_gov_application_endpoint,omitempty" json:"arm_gov_ad_tenant_id,omitempty"`
	AzureGovApplicationClientId           string `form:"azure_gov_application_client_id,omitempty" json:"arm_gov_ad_client_id,omitempty"`
	AzureGovApplicationClientSecret       string `form:"azure_gov_application_client_secret,omitempty" json:"azure_gov_application_client_secret,omitempty"`
	AwsOrangeAccountNumber                string `json:"awsorangecloud_account_number,omitempty"`
	AwsOrangeCapUrl                       string `json:"aws_orange_cap_url,omitempty"`
	AwsOrangeCapAgency                    string `json:"aws_orange_cap_agency,omitempty"`
	AwsOrangeCapMission                   string `json:"aws_orange_cap_mission,omitempty"`
	AwsOrangeCapRoleName                  string `json:"aws_orange_cap_role_name,omitempty"`
	AwsOrangeCapCert                      string
	AwsOrangeCapCertKey                   string
	AwsOrangeCaChainCert                  string
	AwsOrangeCapCertPath                  string `json:"aws_orange_cap_cert_path,omitempty"`
	AwsOrangeCapCertKeyPath               string `json:"aws_orange_cap_key_path,omitempty"`
	AwsCaCertPath                         string `json:"aws_ca_cert_path,omitempty"`
	AwsRedAccountNumber                   string `form:"aws_red_account_number,omitempty" json:"awsredcloud_account_number,omitempty"`
	AwsRedCapUrl                          string `form:"aws_red_cap_url,omitempty" json:"aws_red_cap_url,omitempty"`
	AwsRedCapAgency                       string `form:"aws_red_cap_agency,omitempty" json:"aws_red_cap_agency,omitempty"`
	AwsRedCapAccountName                  string `form:"aws_red_cap_account_name,omitempty" json:"aws_red_cap_account_name,omitempty"`
	AwsRedCapRoleName                     string `json:"aws_red_cap_role_name,omitempty"`
	AwsRedCapCert                         string
	AwsRedCapCertKey                      string
	AwsRedCaChainCert                     string
	AwsRedCapCertPath                     string `form:"aws_red_cap_cert_path,omitempty" json:"aws_red_cap_cert_path,omitempty"`
	AwsRedCapCertKeyPath                  string `json:"aws_red_cap_key_path,omitempty"`
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
	resp, err := c.Post(c.baseURL, account)
	if err != nil {
		return errors.New("HTTP Post setup_account_profile failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode setup_account_profile failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API setup_account_profile Post failed: " + data.Reason)
	}
	return nil
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

	resp, err := c.PostFile(c.baseURL, params, files)
	if err != nil {
		return errors.New("HTTP Post setup_account_profile failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode setup_account_profile failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API setup_account_profile Post failed: " + data.Reason)
	}
	return nil
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

	resp, err := c.PostFile(c.baseURL, params, files)
	if err != nil {
		return errors.New("HTTP Post setup_account_profile failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode setup_account_profile failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API setup_account_profile Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) CreateAWSC2SAccount(account *Account) error {
	params := map[string]string{
		"CID":                       c.CID,
		"action":                    "setup_account_profile",
		"account_name":              account.AccountName,
		"cloud_type":                strconv.Itoa(account.CloudType),
		"aws_orange_account_number": account.AwsOrangeAccountNumber,
		"aws_orange_cap_url":        account.AwsOrangeCapUrl,
		"aws_orange_cap_agency":     account.AwsOrangeCapAgency,
		"aws_orange_cap_mission":    account.AwsOrangeCapMission,
		"aws_orange_cap_role_name":  account.AwsOrangeCapRoleName,
	}

	files := []File{
		{
			Path:      account.AwsOrangeCapCert,
			ParamName: "aws_orange_cap_cert",
		},
		{
			Path:      account.AwsOrangeCapCertKey,
			ParamName: "aws_orange_cap_cert_key",
		},
		{
			Path:      account.AwsOrangeCaChainCert,
			ParamName: "aws_orange_ca_chain_cert",
		},
	}

	resp, err := c.PostFile(c.baseURL, params, files)
	if err != nil {
		return errors.New("HTTP Post setup_account_profile failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode setup_account_profile failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API setup_account_profile Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) CreateAWSSC2SAccount(account *Account) error {
	params := map[string]string{
		"CID":                      c.CID,
		"action":                   "setup_account_profile",
		"account_name":             account.AccountName,
		"cloud_type":               strconv.Itoa(account.CloudType),
		"aws_red_account_number":   account.AwsRedAccountNumber,
		"aws_red_cap_url":          account.AwsRedCapUrl,
		"aws_red_cap_agency":       account.AwsRedCapAgency,
		"aws_red_cap_account_name": account.AwsRedCapAccountName,
		"aws_red_cap_role_name":    account.AwsRedCapRoleName,
	}

	files := []File{
		{
			Path:      account.AwsRedCapCert,
			ParamName: "aws_red_cap_cert",
		},
		{
			Path:      account.AwsRedCapCertKey,
			ParamName: "aws_red_cap_cert_key",
		},
		{
			Path:      account.AwsRedCaChainCert,
			ParamName: "aws_red_ca_chain_cert",
		},
	}

	resp, err := c.PostFile(c.baseURL, params, files)
	if err != nil {
		return errors.New("HTTP Post setup_account_profile failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode setup_account_profile failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API setup_account_profile Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetAccount(account *Account) (*Account, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New("url Parsing failed for list_accounts " + err.Error())
	}
	listAccounts := url.Values{}
	listAccounts.Add("CID", c.CID)
	listAccounts.Add("action", "list_accounts")
	Url.RawQuery = listAccounts.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get list_accounts failed: " + err.Error())
	}

	var data AccountListResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_accounts failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return nil, errors.New("Rest API list_accounts Get failed: " + data.Reason)
	}
	accList := data.Results.AccountList
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
	resp, err := c.Post(c.baseURL, account)
	if err != nil {
		return errors.New("HTTP Post edit_account_profile failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode edit_account_profile failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API edit_account_profile Post failed: " + data.Reason)
	}
	return nil
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

	resp, err := c.PostFile(c.baseURL, params, files)
	if err != nil {
		return errors.New("HTTP Post edit_account_profile failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode edit_account_profile failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API edit_account_profile Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) UpdateAWSC2SAccount(account *Account, fileChanges map[string]bool) error {
	params := map[string]string{
		"CID":                       c.CID,
		"action":                    "edit_account_profile",
		"account_name":              account.AccountName,
		"cloud_type":                strconv.Itoa(account.CloudType),
		"aws_orange_account_number": account.AwsOrangeAccountNumber,
		"aws_orange_cap_url":        account.AwsOrangeCapUrl,
		"aws_orange_cap_agency":     account.AwsOrangeCapAgency,
		"aws_orange_cap_mission":    account.AwsOrangeCapMission,
		"aws_orange_cap_role_name":  account.AwsOrangeCapRoleName,
	}

	files := make([]File, 0, 3)
	if fileChanges["aws_orange_cap_cert"] {
		files = append(files, File{
			Path:      account.AwsOrangeCapCert,
			ParamName: "aws_orange_cap_cert",
		})
	} else {
		params["aws_orange_cap_cert_path"] = account.AwsOrangeCapCertPath
	}

	if fileChanges["aws_orange_cap_cert_key"] {
		files = append(files, File{
			Path:      account.AwsOrangeCapCertKey,
			ParamName: "aws_orange_cap_cert_key",
		})
	} else {
		params["aws_orange_cap_key_path"] = account.AwsOrangeCapCertKeyPath
	}

	if fileChanges["aws_orange_ca_chain_cert"] {
		files = append(files, File{
			Path:      account.AwsOrangeCaChainCert,
			ParamName: "aws_orange_ca_chain_cert",
		})
	} else {
		params["aws_ca_cert_path"] = account.AwsCaCertPath
	}

	resp, err := c.PostFile(c.baseURL, params, files)
	if err != nil {
		return errors.New("HTTP Post edit_account_profile failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode edit_account_profile failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API edit_account_profile Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) UpdateAWSSC2SAccount(account *Account, fileChanges map[string]bool) error {
	params := map[string]string{
		"CID":                      c.CID,
		"action":                   "edit_account_profile",
		"account_name":             account.AccountName,
		"cloud_type":               strconv.Itoa(account.CloudType),
		"aws_red_account_number":   account.AwsRedAccountNumber,
		"aws_red_cap_url":          account.AwsRedCapUrl,
		"aws_red_cap_agency":       account.AwsRedCapAgency,
		"aws_red_cap_account_name": account.AwsRedCapAccountName,
		"aws_red_cap_role_name":    account.AwsRedCapRoleName,
	}

	files := make([]File, 0, 3)
	if fileChanges["aws_red_cap_cert"] {
		files = append(files, File{
			Path:      account.AwsRedCapCert,
			ParamName: "aws_red_cap_cert",
		})
	} else {
		params["aws_red_cap_cert_path"] = account.AwsRedCapCertPath
	}

	if fileChanges["aws_red_cap_cert_key"] {
		files = append(files, File{
			Path:      account.AwsRedCapCertKey,
			ParamName: "aws_red_cap_cert_key",
		})
	} else {
		params["aws_red_cap_key_path"] = account.AwsRedCapCertKeyPath
	}

	if fileChanges["aws_red_ca_chain_cert"] {
		files = append(files, File{
			Path:      account.AwsRedCaChainCert,
			ParamName: "aws_red_ca_chain_cert",
		})
	} else {
		params["aws_ca_cert_path"] = account.AwsCaCertPath
	}

	resp, err := c.PostFile(c.baseURL, params, files)
	if err != nil {
		return errors.New("HTTP Post edit_account_profile failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode edit_account_profile failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API edit_account_profile Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DeleteAccount(account *Account) error {
	path := c.baseURL + fmt.Sprintf("?action=delete_account_profile&CID=%s&account_name=%s",
		c.CID, account.AccountName)
	resp, err := c.Delete(path, nil)
	if err != nil {
		return errors.New("HTTP delete_account_profile failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode delete_account_profile failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API delete_account_profile Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) UploadOciApiPrivateKeyFile(account *Account) error {
	account.CID = c.CID
	account.Action = "upload_file"
	resp, err := c.Post(c.baseURL, account)
	if err != nil {
		return errors.New("HTTP Post upload_file failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode upload_file failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API upload_file Post failed: " + data.Reason)
	}
	return nil
}
