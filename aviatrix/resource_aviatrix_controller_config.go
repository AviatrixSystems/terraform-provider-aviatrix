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
			"backup_config_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Switch to enable/disable controller cloudn backup config.",
			},
			"cloud_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Type of cloud service provider, requires an integer value. Use 1 for AWS.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This parameter represents the name of a Cloud-Account in Aviatrix controller.",
			},
			"bucket_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "S3 Bucket Name for AWS.",
			},
			"multiple": {
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

	backupConfigEnabled := d.Get("backup_config_enabled").(bool)
	cloudType := d.Get("cloud_type").(int)
	accountName := d.Get("account_name").(string)
	bucketName := d.Get("bucket_name").(string)
	multiple := d.Get("multiple").(bool)
	if backupConfigEnabled {
		if cloudType == 0 || accountName == "" || bucketName == "" {
			return fmt.Errorf("please specifdy 'cloud_type', 'account_name' and 'bucket_name' to enable backup configuration")
		}
		cloudnBackupConfig := &goaviatrix.CloudnBackupConfig{
			CloudType:   cloudType,
			AccountName: accountName,
			BucketName:  bucketName,
		}
		if multiple {
			cloudnBackupConfig.Multiple = "true"
		}
		err := client.EnableCloudnBackupConfig(cloudnBackupConfig)
		if err != nil {
			return fmt.Errorf("failed to enable backup configuration: %s", err)
		}
	} else {
		if cloudType != 0 || accountName != "" || bucketName != "" {
			return fmt.Errorf("'cloud_type', 'account_name' and 'bucket_name' should all be empty for not enabling backup configuration")
		}
		if multiple {
			return fmt.Errorf("'multiple' should be empty or set false for not enabling backup configuration")
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
		if cloudnBackupConfig.BackupConfigEnabled == "yes" {
			d.Set("backup_config_enabled", true)
			d.Set("cloud_type", cloudnBackupConfig.CloudType)
			d.Set("account_name", cloudnBackupConfig.AccountName)
			d.Set("bucket_name", cloudnBackupConfig.BucketName)
			if cloudnBackupConfig.Multiple == "yes" {
				d.Set("multiple", true)
			} else {
				d.Set("multiple", false)
			}
		} else {
			d.Set("backup_config_enabled", false)
			d.Set("multiple", false)
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

	if d.HasChange("backup_config_enabled") {
		backupConfigEnabled := d.Get("backup_config_enabled").(bool)
		cloudType := d.Get("cloud_type").(int)
		accountName := d.Get("account_name").(string)
		bucketName := d.Get("bucket_name").(string)
		multiple := d.Get("multiple").(bool)
		if backupConfigEnabled {
			if cloudType == 0 || accountName == "" || bucketName == "" {
				return fmt.Errorf("please specifdy 'cloud_type', 'account_name' and 'bucket_name' to enable backup configuration")
			}
			cloudnBackupConfig := &goaviatrix.CloudnBackupConfig{
				CloudType:   cloudType,
				AccountName: accountName,
				BucketName:  bucketName,
			}
			err := client.EnableCloudnBackupConfig(cloudnBackupConfig)
			if err != nil {
				return fmt.Errorf("failed to enable backup configuration: %s", err)
			}
		} else {
			if cloudType != 0 || accountName != "" || bucketName != "" {
				return fmt.Errorf("'cloud_type', 'account_name' and 'bucket_name' should all be empty for not enabling backup configuration")
			}
			if multiple {
				return fmt.Errorf("'multiple' should be empty or set false for not enabling backup configuration")
			}
			err := client.DisableCloudnBackupConfig()
			if err != nil {
				return fmt.Errorf("failed to disable backup configuration: %s", err)
			}
		}
		d.SetPartial("backup_config_enabled")
	}

	if d.HasChange("cloud_type") && !d.HasChange("backup_config_enabled") {
		return fmt.Errorf("updating 'cloud_type' without updating 'backup_config_enabled' is not allowed")
	}

	if d.HasChange("account_name") && !d.HasChange("backup_config_enabled") {
		return fmt.Errorf("updating 'account_name' without updating 'backup_config_enabled' is not allowed")
	}

	if d.HasChange("bucket_name") && !d.HasChange("backup_config_enabled") {
		return fmt.Errorf("updating 'bucket_name' without updating 'backup_config_enabled' is not allowed")
	}

	if d.HasChange("multiple") && !d.HasChange("backup_config_enabled") {
		return fmt.Errorf("updating 'multiple' without updating 'backup_config_enabled' is not allowed")
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

	d.Set("backup_config_enabled", false)
	cloudnBackupConfig, _ := client.GetCloudnBackupConfig()
	if cloudnBackupConfig.BackupConfigEnabled == "yes" {
		err := client.DisableCloudnBackupConfig()
		if err != nil {
			log.Printf("[ERROR] Failed to disable cloudn backup config on controller %s", d.Id())
			return err
		}
	}

	return nil
}
