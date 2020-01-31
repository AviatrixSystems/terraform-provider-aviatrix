package aviatrix

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixControllerConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixControllerConfigCreate,
		Read:   resourceAviatrixControllerConfigRead,
		Update: resourceAviatrixControllerConfigUpdate,
		Delete: resourceAviatrixControllerConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"sg_management_account_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cloud account name of user.",
			},
			"security_group_management": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Used to manage the Controller instance’s inbound rules from gateways.",
			},
			"http_access": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Switch for http access. Default: false.",
			},
			"fqdn_exception_rule": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "A system-wide mode. Default: true.",
			},
			"target_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The release version number to which the controller will be upgraded to.",
			},
			"backup_configuration": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Switch to enable/disable controller cloudn backup config.",
			},
			"backup_cloud_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Type of cloud service provider, requires an integer value. Use 1 for AWS.",
			},
			"backup_account_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "This parameter represents the name of a Cloud-Account in Aviatrix controller.",
			},
			"backup_bucket_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "S3 Bucket Name for AWS.",
			},
			"multiple_backups": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Switch to enable the controller to backup up to a maximum of 3 rotating backups.",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current version of the controller.",
			},
		},
	}
}

func resourceAviatrixControllerConfigCreate(d *schema.ResourceData, meta interface{}) error {
	var err error

	account := d.Get("sg_management_account_name").(string)

	client := meta.(*goaviatrix.Client)

	log.Printf("[INFO] Configuring Aviatrix controller : %#v", d)

	httpAccess := d.Get("http_access").(bool)
	if httpAccess {
		curStatus, _ := client.GetHttpAccessEnabled()
		if curStatus == "True" {
			log.Printf("[INFO] Http Access is already enabled")
		} else {
			err = client.EnableHttpAccess()
			time.Sleep(10 * time.Second)
		}
	} else {
		curStatus, _ := client.GetHttpAccessEnabled()
		if curStatus == "False" {
			log.Printf("[INFO] Http Access is already disabled")
		} else {
			err = client.DisableHttpAccess()
			time.Sleep(10 * time.Second)
		}
	}
	if err != nil {
		return fmt.Errorf("failed to configure controller http access: %s", err)
	}

	fqdnExceptionRule := d.Get("fqdn_exception_rule").(bool)
	if fqdnExceptionRule {
		curStatus, _ := client.GetExceptionRuleStatus()
		if curStatus {
			log.Printf("[INFO] FQDN Exception Rule is already enabled")
		} else {
			err = client.EnableExceptionRule()
		}
	} else {
		curStatus, _ := client.GetExceptionRuleStatus()
		if !curStatus {
			log.Printf("[INFO] FQDN Exception Rule is already disabled")
		} else {
			err = client.DisableExceptionRule()
		}
	}
	if err != nil {
		return fmt.Errorf("failed to configure controller exception rule: %s", err)
	}

	securityGroupManagement := d.Get("security_group_management").(bool)
	if securityGroupManagement {
		curStatus, _ := client.GetSecurityGroupManagementStatus()
		if curStatus.State == "Enabled" {
			log.Printf("[INFO] Security Group Management is already enabled")
		} else {
			err = client.EnableSecurityGroupManagement(account)
		}
	} else {
		curStatus, _ := client.GetSecurityGroupManagementStatus()
		if curStatus.State == "Disabled" {
			log.Printf("[INFO] Security Group Management is already disabled")
		} else {
			err = client.DisableSecurityGroupManagement()
		}
	}
	if err != nil {
		return fmt.Errorf("failed to configure controller Security Group Management: %s", err)
	}

	version := &goaviatrix.Version{
		Version: d.Get("target_version").(string),
	}
	if version.Version != "" {
		err := client.Upgrade(version)
		if err != nil {
			return fmt.Errorf("failed to upgrade Aviatrix Controller: %s", err)
		}

		newCurrent, _, _ := client.GetCurrentVersion()
		log.Printf("Upgrade complete (now %s)", newCurrent)
	}

	backupConfiguration := d.Get("backup_configuration").(bool)
	backupCloudType := d.Get("backup_cloud_type").(int)
	backupAccountName := d.Get("backup_account_name").(string)
	backupBucketName := d.Get("backup_bucket_name").(string)
	multipleBackups := d.Get("multiple_backups").(bool)
	if backupConfiguration {
		if backupCloudType == 0 || backupAccountName == "" || backupBucketName == "" {
			return fmt.Errorf("please specifdy 'backup_cloud_type', 'backup_account_name' and 'backup_bucket_name'" +
				" to enable backup configuration")
		}
		cloudnBackupConfiguration := &goaviatrix.CloudnBackupConfiguration{
			BackupCloudType:   backupCloudType,
			BackupAccountName: backupAccountName,
			BackupBucketName:  backupBucketName,
		}
		if multipleBackups {
			cloudnBackupConfiguration.MultipleBackups = "true"
		}
		err := client.EnableCloudnBackupConfig(cloudnBackupConfiguration)
		if err != nil {
			return fmt.Errorf("failed to enable backup configuration: %s", err)
		}
	} else {
		if backupCloudType != 0 || backupAccountName != "" || backupBucketName != "" {
			return fmt.Errorf("'backup_cloud_type', 'backup_account_name' and 'backup_bucket_name' should all be empty for not enabling backup configuration")
		}
		if multipleBackups {
			return fmt.Errorf("'multiple_backups' should be empty or set false for not enabling backup configuration")
		}
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixControllerConfigRead(d, meta)
}

func resourceAviatrixControllerConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	log.Printf("[INFO] Getting controller %s configuration", d.Id())
	result, err := client.GetHttpAccessEnabled()
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read Aviatrix Controller Config: %s", err)
	}

	if result[1:5] == "True" {
		d.Set("http_access", true)
	} else {
		d.Set("http_access", false)
	}

	res, err := client.GetExceptionRuleStatus()
	if err != nil {
		return fmt.Errorf("could not read Aviatrix Controller Exception Rule Status: %s", err)
	}
	if res {
		d.Set("fqdn_exception_rule", true)
	} else {
		d.Set("fqdn_exception_rule", false)
	}

	sgm, err := client.GetSecurityGroupManagementStatus()
	if err != nil {
		return fmt.Errorf("could not read Aviatrix Controller Security Group Management Status: %s", err)
	}
	if sgm != nil {
		if sgm.State == "Enabled" {
			d.Set("security_group_management", true)
		} else {
			d.Set("security_group_management", false)
		}
		d.Set("sg_management_account_name", sgm.AccountName)
	} else {
		return fmt.Errorf("could not read Aviatrix Controller Security Group Management Status")
	}

	current, _, err := client.GetCurrentVersion()
	if err != nil {
		return fmt.Errorf("unable to read current Controller version: %s (%s)", err, current)
	}

	targetVersion := d.Get("target_version")
	if targetVersion == "latest" {
		d.Set("target_version", current)
	} else {
		d.Set("target_version", targetVersion)
	}

	d.Set("version", current)

	cloudnBackupConfig, err := client.GetCloudnBackupConfig()
	if err != nil {
		return fmt.Errorf("unable to read current controller cloudn backup config: %s", err)
	}
	if cloudnBackupConfig != nil {
		if cloudnBackupConfig.BackupConfiguration == "yes" {
			d.Set("backup_configuration", true)
			d.Set("backup_cloud_type", cloudnBackupConfig.BackupCloudType)
			d.Set("backup_account_name", cloudnBackupConfig.BackupAccountName)
			d.Set("backup_bucket_name", cloudnBackupConfig.BackupBucketName)
			if cloudnBackupConfig.MultipleBackups == "yes" {
				d.Set("multiple_backups", true)
			} else {
				d.Set("multiple_backups", false)
			}
		} else {
			d.Set("backup_configuration", false)
			d.Set("multiple_backups", false)
		}
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixControllerConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	account := d.Get("sg_management_account_name").(string)

	log.Printf("[INFO] Updating Controller configuration: %#v", d)
	d.Partial(true)

	if d.HasChange("http_access") {
		httpAccess := d.Get("http_access").(bool)
		if httpAccess {
			err := client.EnableHttpAccess()
			time.Sleep(10 * time.Second)
			if err != nil {
				log.Printf("[ERROR] Failed to enable http access on controller %s", d.Id())
				return err
			}
		} else {
			err := client.DisableHttpAccess()
			time.Sleep(10 * time.Second)
			if err != nil {
				log.Printf("[ERROR] Failed to disable http access on controller %s", d.Id())
				return err
			}
		}
		d.SetPartial("http_access")
	}

	if d.HasChange("fqdn_exception_rule") {
		fqdnExceptionRule := d.Get("fqdn_exception_rule").(bool)
		if fqdnExceptionRule {
			err := client.EnableExceptionRule()
			if err != nil {
				log.Printf("[ERROR] Failed to enable exception rule on controller %s", d.Id())
				return err
			}
		} else {
			err := client.DisableExceptionRule()
			if err != nil {
				log.Printf("[ERROR] Failed to disable exception rule on controller %s", d.Id())
				return err
			}
		}
		d.SetPartial("fqdn_exception_rule")
	}

	if d.HasChange("security_group_management") {
		securityGroupManagement := d.Get("security_group_management").(bool)
		if securityGroupManagement {
			err := client.EnableSecurityGroupManagement(account)
			if err != nil {
				log.Printf("[ERROR] Failed to enable Security Group Management on controller %s", d.Id())
				return err
			}
		} else {
			err := client.DisableSecurityGroupManagement()
			if err != nil {
				log.Printf("[ERROR] Failed to disable Security Group Management on controller %s", d.Id())
				return err
			}
		}
		d.SetPartial("security_group_management")
	}

	if d.HasChange("target_version") {
		curVersion := d.Get("version").(string)
		cur := strings.Split(curVersion, ".")
		latestVersion, _ := client.GetLatestVersion()
		latest := strings.Split(latestVersion, ".")
		version := &goaviatrix.Version{
			Version: d.Get("target_version").(string),
		}

		targetVersion := d.Get("target_version").(string)
		if targetVersion == "latest" {
			if latestVersion != "" {
				for i := range cur {
					if cur[i] != latest[i] {
						err := client.Upgrade(version)
						if err != nil {
							return fmt.Errorf("failed to upgrade Aviatrix Controller: %s", err)
						}
						break
					}
				}
			}
		} else {
			err := client.Upgrade(version)
			if err != nil {
				return fmt.Errorf("failed to upgrade Aviatrix Controller: %s", err)
			}
		}
		d.SetPartial("target_version")
	}

	if d.HasChange("backup_configuration") {
		backupConfiguration := d.Get("backup_configuration").(bool)
		backupCloudType := d.Get("backup_cloud_type").(int)
		backupAccountName := d.Get("backup_account_name").(string)
		backupBucketName := d.Get("backup_bucket_name").(string)
		multipleBackups := d.Get("multiple_backups").(bool)
		if backupConfiguration {
			if backupCloudType == 0 || backupAccountName == "" || backupBucketName == "" {
				return fmt.Errorf("please specifdy 'backup_cloud_type', 'backup_account_name' and 'backup_bucket_name' to enable backup configuration")
			}
			cloudnBackupConfig := &goaviatrix.CloudnBackupConfiguration{
				BackupCloudType:   backupCloudType,
				BackupAccountName: backupAccountName,
				BackupBucketName:  backupBucketName,
			}
			err := client.EnableCloudnBackupConfig(cloudnBackupConfig)
			if err != nil {
				return fmt.Errorf("failed to enable backup configuration: %s", err)
			}
		} else {
			if backupCloudType != 0 || backupAccountName != "" || backupBucketName != "" {
				return fmt.Errorf("'backup_cloud_type', 'backup_account_name' and 'backup_bucket_name' should all be empty for not enabling backup configuration")
			}
			if multipleBackups {
				return fmt.Errorf("'multiple_backups' should be empty or set false for not enabling backup configuration")
			}
			err := client.DisableCloudnBackupConfig()
			if err != nil {
				return fmt.Errorf("failed to disable backup configuration: %s", err)
			}
		}
		d.SetPartial("backup_configuration")
	}

	if d.HasChange("backup_cloud_type") && !d.HasChange("backup_configuration") {
		return fmt.Errorf("updating 'backup_cloud_type' without updating 'backup_configuration' is not allowed")
	}

	if d.HasChange("backup_account_name") && !d.HasChange("backup_configuration") {
		return fmt.Errorf("updating 'backup_account_name' without updating 'backup_configuration' is not allowed")
	}

	if d.HasChange("backup_bucket_name") && !d.HasChange("backup_configuration") {
		return fmt.Errorf("updating 'backup_bucket_name' without updating 'backup_configuration' is not allowed")
	}

	if d.HasChange("multiple_backups") && !d.HasChange("backup_configuration") {
		return fmt.Errorf("updating 'multiple_backups' without updating 'backup_configuration' is not allowed")
	}

	d.Partial(false)
	return resourceAviatrixControllerConfigRead(d, meta)
}

func resourceAviatrixControllerConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	d.Set("http_access", false)
	curStatusHttp, _ := client.GetHttpAccessEnabled()
	if curStatusHttp != "Disabled" {
		err := client.DisableHttpAccess()
		time.Sleep(10 * time.Second)
		if err != nil {
			log.Printf("[ERROR] Failed to disable http access on controller %s", d.Id())
			return err
		}
	}

	d.Set("fqdn_exception_rule", true)
	curStatusException, _ := client.GetExceptionRuleStatus()
	if !curStatusException {
		err := client.EnableExceptionRule()
		if err != nil {
			log.Printf("[ERROR] Failed to enable exception rule on controller %s", d.Id())
			return err
		}
	}

	d.Set("security_group_management", false)
	curStatusSG, _ := client.GetSecurityGroupManagementStatus()
	if curStatusSG.State != "Disabled" {
		err := client.DisableSecurityGroupManagement()
		if err != nil {
			log.Printf("[ERROR] Failed to disable security group management on controller %s", d.Id())
			return err
		}
	}

	d.Set("backup_configuration", false)
	cloudnBackupConfig, _ := client.GetCloudnBackupConfig()
	if cloudnBackupConfig.BackupConfiguration == "yes" {
		err := client.DisableCloudnBackupConfig()
		if err != nil {
			log.Printf("[ERROR] Failed to disable cloudn backup config on controller %s", d.Id())
			return err
		}
	}

	return nil
}
