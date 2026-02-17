package aviatrix

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

const defaultAwsGuardDutyScanningInterval = 60

var validAwsGuardDutyScanningIntervals = []int{5, 10, 15, 30, 60}

func resourceAviatrixControllerConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixControllerConfigCreate,
		Read:   resourceAviatrixControllerConfigRead,
		Update: resourceAviatrixControllerConfigUpdate,
		Delete: resourceAviatrixControllerConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},

		Schema: map[string]*schema.Schema{
			// The http_access API hasn't worked correctly
			// since at least version 7.1. We leave the schema
			// element here so that we don't break any existing
			// Terraform code, but it doesn't do anything.
			"http_access": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Switch for HTTP access. Default: false.",
				Deprecated:  "HTTP access is no longer supported and will not be configured on the controller.",
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
			"manage_gateway_upgrades": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				Description: "If true, aviatrix_controller_config will upgrade all gateways when target_version " +
					"is set. If false, only the controller will be upgraded when target_version is set. In that " +
					"case gateway upgrades should be handled in each gateway resource individually using the " +
					"software_version and image_version attributes.",
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
				Description: "Bucket name. Required for AWS, AWSGov, GCP and OCI.",
			},
			"backup_storage_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Storage name. Required for Azure.",
			},
			"backup_container_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Container name. Required for Azure.",
			},
			"backup_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Name of region. Required for Azure and OCI.",
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
				Description: "Current version of the controller without the build number.",
			},
			"current_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current version of the controller.",
			},
			"previous_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Previous version of the controller.",
			},
			"enable_vpc_dns_server": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable VPC/VNET DNS Server.",
			},
			"ca_certificate_file_path": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "File path to the CA certificate.",
				RequiredWith:  []string{"server_public_certificate_file_path", "server_private_key_file_path"},
				ConflictsWith: []string{"ca_certificate_file"},
			},
			"server_public_certificate_file_path": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "File path to the server public certificate.",
				RequiredWith:  []string{"ca_certificate_file_path", "server_private_key_file_path"},
				ConflictsWith: []string{"server_public_certificate_file"},
			},
			"server_private_key_file_path": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "File path to the server private key.",
				RequiredWith:  []string{"server_public_certificate_file_path", "ca_certificate_file_path"},
				ConflictsWith: []string{"server_private_key_file"},
			},
			"ca_certificate_file": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				Description:   "CA certificate file.",
				RequiredWith:  []string{"server_public_certificate_file", "server_private_key_file"},
				ConflictsWith: []string{"ca_certificate_file_path"},
			},
			"server_public_certificate_file": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				Description:   "Server public certificate file.",
				RequiredWith:  []string{"ca_certificate_file", "server_private_key_file"},
				ConflictsWith: []string{"server_public_certificate_file_path"},
			},
			"server_private_key_file": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				Description:   "Server private key file.",
				RequiredWith:  []string{"server_public_certificate_file", "ca_certificate_file"},
				ConflictsWith: []string{"server_private_key_file_path"},
			},
			"aws_guard_duty_scanning_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "Scanning Interval for AWS Guard Duty.",
				Default:      defaultAwsGuardDutyScanningInterval,
				ValidateFunc: validation.IntInSlice(validAwsGuardDutyScanningIntervals),
			},
		},
	}
}

func resourceAviatrixControllerConfigCreate(d *schema.ResourceData, meta interface{}) error {
	var err error

	client := mustClient(meta)

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	flag := false
	defer func() { _ = resourceAviatrixControllerConfigReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	log.Printf("[INFO] Configuring Aviatrix controller : %#v", d)

	fqdnExceptionRule := getBool(d, "fqdn_exception_rule")
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
		return fmt.Errorf("failed to configure controller exception rule: %w", err)
	}

	backupConfiguration := getBool(d, "backup_configuration")
	backupCloudType := getInt(d, "backup_cloud_type")
	backupAccountName := getString(d, "backup_account_name")
	backupBucketName := getString(d, "backup_bucket_name")
	backupStorageName := getString(d, "backup_storage_name")
	backupContainerName := getString(d, "backup_container_name")
	backupRegion := getString(d, "backup_region")
	multipleBackups := getBool(d, "multiple_backups")

	if backupConfiguration {
		err = validateBackupConfig(d)
		if err != nil {
			return err
		}

		cloudnBackupConfiguration := &goaviatrix.CloudnBackupConfiguration{
			BackupCloudType:     backupCloudType,
			BackupAccountName:   backupAccountName,
			BackupBucketName:    backupBucketName,
			BackupStorageName:   backupStorageName,
			BackupContainerName: backupContainerName,
			BackupRegion:        backupRegion,
		}
		if multipleBackups {
			cloudnBackupConfiguration.MultipleBackups = "true"
		}

		err = client.EnableCloudnBackupConfig(cloudnBackupConfiguration)
		if err != nil {
			return fmt.Errorf("failed to enable backup configuration: %w", err)
		}
	} else {
		if backupCloudType != 0 || backupAccountName != "" || backupBucketName != "" || backupStorageName != "" ||
			backupContainerName != "" || backupRegion != "" || multipleBackups {
			return fmt.Errorf("'backup_cloud_type', 'backup_account_name', 'backup_bucket_name'," +
				" 'backup_storage_name', 'backup_container_name' and 'backup_region' should all be empty," +
				" 'multiple_backups' should be empty or false for not enabling backup configuration")
		}
	}

	enableVpcDnsServer := getBool(d, "enable_vpc_dns_server")
	err = client.SetControllerVpcDnsServer(enableVpcDnsServer)
	if err != nil {
		return fmt.Errorf("could not toggle controller vpc dns server: %w", err)
	}

	if _, useFilePath := d.GetOk("ca_certificate_file_path"); useFilePath {
		certConfig := &goaviatrix.HTTPSCertConfig{
			CACertificateFilePath:     getString(d, "ca_certificate_file_path"),
			ServerCertificateFilePath: getString(d, "server_public_certificate_file_path"),
			ServerPrivateKeyFilePath:  getString(d, "server_private_key_file_path"),
		}
		err = client.ImportNewHTTPSCerts(certConfig)
		if err != nil {
			return fmt.Errorf("could not import HTTPS certs: %w", err)
		}
	} else if _, useFileContent := d.GetOk("ca_certificate_file"); useFileContent {
		certConfig := &goaviatrix.HTTPSCertConfig{
			CACertificateFile:     getString(d, "ca_certificate_file"),
			ServerCertificateFile: getString(d, "server_public_certificate_file"),
			ServerPrivateKeyFile:  getString(d, "server_private_key_file"),
		}
		err = client.ImportNewHTTPSCerts(certConfig)
		if err != nil {
			return fmt.Errorf("could not import HTTPS certs: %w", err)
		}
	}

	scanningInterval := getInt(d, "aws_guard_duty_scanning_interval")
	err = client.UpdateAwsGuardDutyPollInterval(scanningInterval)
	if err != nil {
		return fmt.Errorf("could not update scanning interval: %w", err)
	}

	return resourceAviatrixControllerConfigReadIfRequired(d, meta, &flag)
}

func resourceAviatrixControllerConfigReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixControllerConfigRead(d, meta)
	}
	return nil
}

func resourceAviatrixControllerConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	log.Printf("[INFO] Getting controller %s configuration", d.Id())

	res, err := client.GetExceptionRuleStatus()
	if err != nil {
		return fmt.Errorf("could not read Aviatrix Controller Exception Rule Status: %w", err)
	}
	if res {
		mustSet(d, "fqdn_exception_rule", true)
	} else {
		mustSet(d, "fqdn_exception_rule", false)
	}

	var versionInfo *goaviatrix.VersionInfo
	try, maxTries, backoff := 0, 3, 1000*time.Millisecond
	for {
		try++
		versionInfo, err = client.GetVersionInfo()
		if err != nil {
			if try == maxTries {
				return fmt.Errorf("unable to read Controller version information: %w", err)
			}
			time.Sleep(backoff)
			// Double the backoff time after each failed try
			backoff *= 2
			continue
		}
		break
	}
	mustSet(d, "version", versionInfo.Current)
	mustSet(d, "previous_version", versionInfo.Previous)

	cloudnBackupConfig, err := client.GetCloudnBackupConfig()
	if err != nil {
		return fmt.Errorf("unable to read current controller cloudn backup config: %w", err)
	}
	if cloudnBackupConfig != nil && cloudnBackupConfig.BackupConfiguration == "yes" {
		mustSet(d, "backup_configuration", true)
		mustSet(d, "backup_cloud_type", cloudnBackupConfig.BackupCloudType)
		mustSet(d, "backup_account_name", cloudnBackupConfig.BackupAccountName)
		mustSet(d, "backup_bucket_name", cloudnBackupConfig.BackupBucketName)
		mustSet(d, "backup_storage_name", cloudnBackupConfig.BackupStorageName)
		mustSet(d, "backup_container_name", cloudnBackupConfig.BackupContainerName)
		mustSet(d, "backup_region", cloudnBackupConfig.BackupRegion)
		if cloudnBackupConfig.MultipleBackups == "yes" {
			mustSet(d, "multiple_backups", true)
		} else {
			mustSet(d, "multiple_backups", false)
		}
	} else {
		mustSet(d, "backup_cloud_type", 0)
		mustSet(d, "backup_configuration", false)
		mustSet(d, "multiple_backups", false)
	}

	vpcDnsServerEnabled, err := client.GetControllerVpcDnsServerStatus()
	if err != nil {
		return fmt.Errorf("could not get controller vpc dns server status: %w", err)
	}
	mustSet(d, "enable_vpc_dns_server", vpcDnsServerEnabled)

	httpsCertsImported, err := client.GetHTTPSCertsStatus()
	if err != nil {
		return fmt.Errorf("could not get HTTPS Certificate status: %w", err)
	}
	if !httpsCertsImported {
		mustSet(d, "ca_certificate_file_path", "")
		mustSet(d, "server_public_certificate_file_path", "")
		mustSet(d, "server_private_key_file_path", "")
		mustSet(d, "ca_certificate_file", "")
		mustSet(d, "server_public_certificate_file", "")
		mustSet(d, "server_private_key_file", "")
	}

	guardDuty, err := client.GetAwsGuardDuty()
	if err != nil {
		return fmt.Errorf("could not get aws guard duty scanning interval: %w", err)
	}
	mustSet(d, "aws_guard_duty_scanning_interval", guardDuty.ScanningInterval)

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixControllerConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	log.Printf("[INFO] Updating Controller configuration: %#v", d)
	d.Partial(true)

	if d.HasChange("fqdn_exception_rule") {
		fqdnExceptionRule := getBool(d, "fqdn_exception_rule")
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
	}

	backupConfiguration := getBool(d, "backup_configuration")
	backupCloudType := getInt(d, "backup_cloud_type")
	backupAccountName := getString(d, "backup_account_name")
	backupBucketName := getString(d, "backup_bucket_name")
	backupStorageName := getString(d, "backup_storage_name")
	backupContainerName := getString(d, "backup_container_name")
	backupRegion := getString(d, "backup_region")
	multipleBackups := getBool(d, "multiple_backups")

	if d.HasChange("backup_configuration") {
		if backupConfiguration {
			err := validateBackupConfig(d)
			if err != nil {
				return err
			}

			cloudnBackupConfiguration := &goaviatrix.CloudnBackupConfiguration{
				BackupCloudType:     backupCloudType,
				BackupAccountName:   backupAccountName,
				BackupBucketName:    backupBucketName,
				BackupStorageName:   backupStorageName,
				BackupContainerName: backupContainerName,
				BackupRegion:        backupRegion,
			}
			if multipleBackups {
				cloudnBackupConfiguration.MultipleBackups = "true"
			}

			err = client.EnableCloudnBackupConfig(cloudnBackupConfiguration)
			if err != nil {
				return fmt.Errorf("failed to enable backup configuration: %w", err)
			}
		} else {
			if backupCloudType != 0 || backupAccountName != "" || backupBucketName != "" || backupStorageName != "" ||
				backupContainerName != "" || backupRegion != "" || multipleBackups {
				return fmt.Errorf("'backup_cloud_type', 'backup_account_name', 'backup_bucket_name'," +
					" 'backup_storage_name', 'backup_container_name' and 'backup_region' should all be empty," +
					" 'multiple_backups' should be empty or false for not enabling backup configuration")
			}

			err := client.DisableCloudnBackupConfig()
			if err != nil {
				return fmt.Errorf("failed to disable backup configuration: %w", err)
			}
		}
	} else {
		if d.HasChange("backup_cloud_type") || d.HasChange("backup_account_name") ||
			d.HasChange("backup_bucket_name") || d.HasChange("backup_storage_name") ||
			d.HasChange("backup_container_name") || d.HasChange("backup_region") ||
			d.HasChange("multiple_backups") {

			if backupConfiguration {
				err := validateBackupConfig(d)
				if err != nil {
					return err
				}

				err = client.DisableCloudnBackupConfig()
				if err != nil {
					return fmt.Errorf("failed to disable backup configuration: %w", err)
				}

				cloudnBackupConfiguration := &goaviatrix.CloudnBackupConfiguration{
					BackupCloudType:     backupCloudType,
					BackupAccountName:   backupAccountName,
					BackupBucketName:    backupBucketName,
					BackupStorageName:   backupStorageName,
					BackupContainerName: backupContainerName,
					BackupRegion:        backupRegion,
				}
				if multipleBackups {
					cloudnBackupConfiguration.MultipleBackups = "true"
				}

				err = client.EnableCloudnBackupConfig(cloudnBackupConfiguration)
				if err != nil {
					return fmt.Errorf("failed to enable backup configuration: %w", err)
				}
			} else {
				if backupCloudType != 0 || backupAccountName != "" || backupBucketName != "" || backupStorageName != "" ||
					backupContainerName != "" || backupRegion != "" || multipleBackups {
					return fmt.Errorf("'backup_cloud_type', 'backup_account_name', 'backup_bucket_name'," +
						" 'backup_storage_name', 'backup_container_name' and 'backup_region' should all be empty," +
						" 'multiple_backups' should be empty or false for not enabling backup configuration")
				}
			}
		}
	}

	if d.HasChange("enable_vpc_dns_server") {
		enableVpcDnsServer := getBool(d, "enable_vpc_dns_server")
		err := client.SetControllerVpcDnsServer(enableVpcDnsServer)
		if err != nil {
			return fmt.Errorf("could not toggle controller vpc dns server: %w", err)
		}
	}

	if d.HasChange("ca_certificate_file_path") || d.HasChange("server_public_certificate_file_path") || d.HasChange("server_private_key_file_path") ||
		d.HasChange("ca_certificate_file") || d.HasChange("server_public_certificate_file") || d.HasChange("server_private_key_file") {
		if _, useFilePath := d.GetOk("ca_certificate_file_path"); useFilePath {
			certConfig := &goaviatrix.HTTPSCertConfig{
				CACertificateFilePath:     getString(d, "ca_certificate_file_path"),
				ServerCertificateFilePath: getString(d, "server_public_certificate_file_path"),
				ServerPrivateKeyFilePath:  getString(d, "server_private_key_file_path"),
			}

			err := client.ImportNewHTTPSCerts(certConfig)
			if err != nil {
				return fmt.Errorf("could not import new HTTPS certs: %w", err)
			}
		} else if _, useFileContent := d.GetOk("ca_certificate_file"); useFileContent {
			certConfig := &goaviatrix.HTTPSCertConfig{
				CACertificateFile:     getString(d, "ca_certificate_file"),
				ServerCertificateFile: getString(d, "server_public_certificate_file"),
				ServerPrivateKeyFile:  getString(d, "server_private_key_file"),
			}

			err := client.ImportNewHTTPSCerts(certConfig)
			if err != nil {
				return fmt.Errorf("could not import new HTTPS certs: %w", err)
			}
		}
	}

	if d.HasChange("aws_guard_duty_scanning_interval") {
		scanningInterval := getInt(d, "aws_guard_duty_scanning_interval")
		err := client.UpdateAwsGuardDutyPollInterval(scanningInterval)
		if err != nil {
			return fmt.Errorf("could not update scanning interval: %w", err)
		}
	}

	d.Partial(false)
	return resourceAviatrixControllerConfigRead(d, meta)
}

func resourceAviatrixControllerConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)
	mustSet(d, "fqdn_exception_rule", true)
	curStatusException, _ := client.GetExceptionRuleStatus()
	if !curStatusException {
		err := client.EnableExceptionRule()
		if err != nil {
			log.Printf("[ERROR] Failed to enable exception rule on controller %s", d.Id())
			return err
		}
	}
	mustSet(d, "backup_configuration", false)
	cloudnBackupConfig, _ := client.GetCloudnBackupConfig()
	if cloudnBackupConfig.BackupConfiguration == "yes" {
		err := client.DisableCloudnBackupConfig()
		if err != nil {
			log.Printf("[ERROR] Failed to disable cloudn backup config on controller %s", d.Id())
			return err
		}
	}

	err := client.SetControllerVpcDnsServer(false)
	if err != nil {
		return fmt.Errorf("could not disable controller vpc dns server: %w", err)
	}

	err = client.DisableImportedHTTPSCerts()
	if err != nil {
		return fmt.Errorf("could not disable imported certs: %w", err)
	}

	err = client.UpdateAwsGuardDutyPollInterval(defaultAwsGuardDutyScanningInterval)
	if err != nil {
		return fmt.Errorf("could not update scanning interval: %w", err)
	}

	return nil
}

func validateBackupConfig(d *schema.ResourceData) error {
	backupCloudType := getInt(d, "backup_cloud_type")
	backupAccountName := getString(d, "backup_account_name")
	backupBucketName := getString(d, "backup_bucket_name")
	backupStorageName := getString(d, "backup_storage_name")
	backupContainerName := getString(d, "backup_container_name")
	backupRegion := getString(d, "backup_region")

	if backupCloudType == 0 || backupAccountName == "" {
		return fmt.Errorf("please specify 'backup_cloud_type' and 'backup_account_name'" +
			" to enable backup configuration")
	}

	switch backupCloudType {
	case goaviatrix.AWS, goaviatrix.AWSGov, goaviatrix.GCP:
		if backupBucketName == "" {
			return fmt.Errorf("please specify 'backup_bucket_name' to enable backup configuration for AWS, AWSGov and GCP")
		}
		if backupStorageName != "" || backupContainerName != "" || backupRegion != "" {
			return fmt.Errorf("'backup_storage_name', 'backup_container_name' and 'backup_region'" +
				" should be empty for AWS, AWSGov and GCP")
		}
	case goaviatrix.Azure:
		if backupStorageName == "" || backupContainerName == "" || backupRegion == "" {
			return fmt.Errorf("please specify 'backup_storage_name', 'backup_container_name' and" +
				" 'backup_region' to enable backup configuration for Azure")
		}
		if backupBucketName != "" {
			return fmt.Errorf("'backup_bucket_name' should be empty for Azure")
		}
	case goaviatrix.OCI:
		if backupBucketName == "" || backupRegion == "" {
			return fmt.Errorf("please specify 'backup_bucket_name' and 'backup_region'" +
				" to enable backup configuration for OCI")
		}
		if backupStorageName != "" || backupContainerName != "" {
			return fmt.Errorf("'backup_storage_name' and 'backup_container_name'" +
				" should be empty for OCI")
		}
	}

	return nil
}
