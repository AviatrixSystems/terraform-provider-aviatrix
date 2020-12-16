package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixFirewallInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixFirewallInstanceCreate,
		Read:   resourceAviatrixFirewallInstanceRead,
		Delete: resourceAviatrixFirewallInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the Security VPC.",
			},
			"firenet_gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the primary FireNet gateway.",
			},
			"firewall_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the firewall instance to be created.",
			},
			"firewall_image": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "One of the AWS AMIs from Palo Alto Networks.",
			},
			"firewall_size": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance size of the firewall.",
			},
			"egress_subnet": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Egress Interface Subnet.",
			},
			"management_subnet": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "Management Interface Subnet. Required for Palo Alto Networks VM-Series, " +
					"and required to be empty for Check Point or Fortinet series.",
			},
			"firewall_image_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Version of firewall image.",
			},
			"key_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The .pem file name for SSH access to the firewall instance.",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Applicable to Azure deployment only. 'admin' as a username is not accepted.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "Authentication method. Applicable to Azure deployment only.",
			},
			"ssh_public_key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "Authentication method. Applicable to Azure deployment only.",
			},
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureAZ,
				Description:  "Availability Zone. Only available for AZURE. Must be in the form 'az-n', for example, 'az-2'.",
			},
			"iam_role": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Advanced option. IAM role. Only available for AWS.",
			},
			"bootstrap_bucket_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Advanced option. Bootstrap bucket name. Only available for AWS.",
			},
			"bootstrap_storage_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "Advanced option. Bootstrap storage name. Applicable to Azure and Palo Alto Networks " +
					"VM-Series/Fortinet Series deployment only.",
			},
			"storage_access_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				ForceNew:  true,
				Description: "Advanced option. Storage access key. Applicable to Azure and Palo Alto Networks " +
					"VM-Series deployment only.",
			},
			"file_share_folder": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "Advanced option. File share folder. Applicable to Azure and Palo Alto Networks " +
					"VM-Series deployment only.",
			},
			"share_directory": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "Advanced option. Share directory. Applicable to Azure and Palo Alto Networks " +
					"VM-Series deployment only.",
			},
			"sic_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "Advanced option. Bic key. Applicable to Azure and Check Point Series deployment only.",
			},
			"container_folder": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Advanced option. Bootstrap storage name. Applicable to Azure and Fortinet Series deployment only.",
			},
			"sas_url_config": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Advanced option. Bootstrap storage name. Applicable to Azure and Fortinet Series deployment only.",
			},
			"sas_url_license": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Advanced option. Bootstrap storage name. Applicable to Azure and Fortinet Series deployment only.",
			},
			"user_data": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Advanced option. Bootstrap storage name. Applicable to Check Point Series and Fortinet Series deployment only.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the firewall instance created.",
			},
			"lan_interface": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of Lan Interface created.",
			},
			"management_interface": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of Management Interface created.",
			},
			"egress_interface": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of Egress Interface created.",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Management Public IP.",
			},
		},
	}
}

func resourceAviatrixFirewallInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	firewallInstance := &goaviatrix.FirewallInstance{
		VpcID:                d.Get("vpc_id").(string),
		GwName:               d.Get("firenet_gw_name").(string),
		FirewallName:         d.Get("firewall_name").(string),
		FirewallImage:        d.Get("firewall_image").(string),
		FirewallImageVersion: d.Get("firewall_image_version").(string),
		FirewallSize:         d.Get("firewall_size").(string),
		EgressSubnet:         d.Get("egress_subnet").(string),
		ManagementSubnet:     d.Get("management_subnet").(string),
		KeyName:              d.Get("key_name").(string),
		IamRole:              d.Get("iam_role").(string),
		BootstrapBucketName:  d.Get("bootstrap_bucket_name").(string),
		Username:             d.Get("username").(string),
		Password:             d.Get("password").(string),
		SshPublicKey:         d.Get("ssh_public_key").(string),
		BootstrapStorageName: d.Get("bootstrap_storage_name").(string),
		StorageAccessKey:     d.Get("storage_access_key").(string),
		FileShareFolder:      d.Get("file_share_folder").(string),
		ShareDirectory:       d.Get("share_directory").(string),
		SicKey:               d.Get("sic_key").(string),
		ContainerFolder:      d.Get("container_folder").(string),
		SasUrlConfig:         d.Get("sas_url_config").(string),
		SasUriLicense:        d.Get("sas_url_license").(string),
		UserData:             d.Get("user_data").(string),
	}

	if strings.HasPrefix(firewallInstance.FirewallImage, "Palo Alto Networks") {
		if firewallInstance.ManagementSubnet == "" {
			return fmt.Errorf("'management_subnet' is required unempty for Palo Alto Networks VM-Series")
		}
	} else if strings.HasPrefix(firewallInstance.FirewallImage, "Check Point CloudGuard") ||
		strings.HasPrefix(firewallInstance.FirewallImage, "Fortinet FortiGate") {
		if firewallInstance.ManagementSubnet != "" {
			return fmt.Errorf("'management_subnet' is required to be empty for Check Point or Fortinet FortiGate series")
		}
	} else {
		return fmt.Errorf("firewall image: %s is not supported", firewallInstance.FirewallImage)
	}

	// For additional config validation we try to get the cloud_type from the given
	// gateway name. If there is an issue, we will just continue on without the additional
	// validation.
	var cloudType int
	gw, err := client.GetGateway(&goaviatrix.Gateway{GwName: firewallInstance.GwName})
	if err != nil {
		log.Printf("[WARN] Could not get cloud_type from firenet_gw_name: %v", err)
	} else {
		cloudType = gw.CloudType
	}

	if err != nil {
		if err == goaviatrix.ErrNotFound {
			return fmt.Errorf("could not find the vpc with vpc_id=%s: %v", firewallInstance.VpcID, err)
		}
		return fmt.Errorf("could get the cloud type from the vpc_id=%s: %v", firewallInstance.VpcID, err)
	}
	zone := d.Get("zone").(string)
	if zone != "" && cloudType != goaviatrix.AZURE {
		return fmt.Errorf("'zone' attribute is only valid for AZURE")
	}
	if zone != "" {
		firewallInstance.EgressSubnet = fmt.Sprintf("%s~~%s~~", firewallInstance.EgressSubnet, zone)
		firewallInstance.ManagementSubnet = fmt.Sprintf("%s~~%s~~", firewallInstance.ManagementSubnet, zone)
	}

	if firewallInstance.Username != "" || firewallInstance.Password != "" || firewallInstance.SshPublicKey != "" {
		if cloudType != 0 && cloudType != goaviatrix.AZURE {
			return fmt.Errorf("'username' and 'password' or 'ssh_public_key' are only supported for Azure")
		}
	}
	if firewallInstance.Password != "" && firewallInstance.SshPublicKey != "" {
		return fmt.Errorf("anthentication method can be either a password or an SSH public key. Please specify one of them and set the other one to empty")
	}
	if firewallInstance.IamRole != "" || firewallInstance.BootstrapBucketName != "" {
		if cloudType != 0 && cloudType != goaviatrix.AWS {
			return fmt.Errorf("advanced options of 'iam_role' and 'bootstrap_bucket_name' are only supported for AWS provider, please set them to empty")
		}
	}
	if firewallInstance.UserData != "" && strings.HasPrefix(firewallInstance.FirewallImage, "Palo Alto Networks") {
		return fmt.Errorf("advanced option of 'user_data' is only supported for Check Point Series and Fortinet FortiGate Series, not for %s", firewallInstance.FirewallImage)
	}
	if firewallInstance.StorageAccessKey != "" || firewallInstance.FileShareFolder != "" || firewallInstance.ShareDirectory != "" {
		if !strings.HasPrefix(firewallInstance.FirewallImage, "Palo Alto Networks") || (cloudType != 0 && cloudType != goaviatrix.AZURE) {
			return fmt.Errorf("advanced options of 'storage_access_key', 'file_share_folder' and 'share_directory' are only supported for Azure and Palo Alto Networks VM-Series")
		}
	}
	if firewallInstance.ContainerFolder != "" || firewallInstance.SasUrlConfig != "" || firewallInstance.SasUriLicense != "" {
		if !strings.HasPrefix(firewallInstance.FirewallImage, "Fortinet FortiGate") || (cloudType != 0 && cloudType != goaviatrix.AZURE) {
			return fmt.Errorf("advanced options of 'container_folder', 'sas_url_config' and 'sas_url_license' are only supported for Azure and Fortinet FortiGate series")
		}
	}
	if firewallInstance.BootstrapStorageName != "" {
		if strings.HasPrefix(firewallInstance.FirewallImage, "Check Point CloudGuard") || (cloudType != 0 && cloudType != goaviatrix.AZURE) {
			return fmt.Errorf("advanced option of 'bootstrap_storage_name' is only supported for Azure and Palo Alto Networks VM-Series/Fortinet FortiGate series")
		}
	}
	if firewallInstance.SicKey != "" {
		if !strings.HasPrefix(firewallInstance.FirewallImage, "Check Point CloudGuard") || (cloudType != 0 && cloudType != goaviatrix.AZURE) {
			return fmt.Errorf("advanced option of 'bootstrap_storage_name' is only supported for Azure and Check Point Series")
		}
	}

	instanceID, err := client.CreateFirewallInstance(firewallInstance)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			return fmt.Errorf("failed to get firewall instance information")
		}
		return fmt.Errorf("failed to create a new firewall instance: %s", err)
	}

	d.SetId(instanceID)
	return resourceAviatrixFirewallInstanceRead(d, meta)
}

func resourceAviatrixFirewallInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	instanceID := d.Get("instance_id").(string)
	var isImport bool
	if instanceID == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no firewall names received. Import Id is %s", id)
		d.Set("instance_id", id)
		d.SetId(id)
		isImport = true
	}

	firewallInstance := &goaviatrix.FirewallInstance{
		InstanceID: d.Get("instance_id").(string),
	}

	fI, err := client.GetFirewallInstance(firewallInstance)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Firewall Instance: %s", err)
	}

	log.Printf("[INFO] Found Firewall Instance: %#v", firewallInstance)

	d.Set("vpc_id", fI.VpcID)
	d.Set("firenet_gw_name", fI.GwName)
	d.Set("firewall_name", fI.FirewallName)
	d.Set("firewall_image", fI.FirewallImage)
	d.Set("firewall_size", fI.FirewallSize)
	d.Set("instance_id", fI.InstanceID)
	d.Set("egress_subnet", fI.EgressSubnet)
	if strings.HasPrefix(fI.FirewallImage, "Palo Alto Networks") {
		d.Set("management_subnet", fI.ManagementSubnet)
	}

	if (d.Get("zone").(string) != "" || isImport) && fI.AvailabilityZone != "AvailabilitySet" &&
		fI.AvailabilityZone != "" && fI.CloudVendor == "Azure ARM" {
		d.Set("zone", "az-"+fI.AvailabilityZone)
	}

	d.Set("lan_interface", fI.LanInterface)
	d.Set("management_interface", fI.ManagementInterface)
	d.Set("egress_interface", fI.EgressInterface)
	d.Set("public_ip", fI.ManagementPublicIP)

	if fI.FirewallImageVersion != "" {
		d.Set("firewall_image_version", fI.FirewallImageVersion)
	}
	if d.Get("key_name") != "" {
		d.Set("key_name", fI.KeyName)
	}
	if fI.IamRole != "" {
		d.Set("iam_role", fI.IamRole)
	}
	if fI.BootstrapBucketName != "" {
		d.Set("bootstrap_bucket_name", fI.BootstrapBucketName)
	}
	if fI.Username != "" {
		d.Set("username", fI.Username)
	}
	if fI.BootstrapStorageName != "" {
		d.Set("bootstrap_storage_name", fI.BootstrapStorageName)
	}
	if fI.FileShareFolder != "" {
		d.Set("file_share_folder", fI.FileShareFolder)
	}
	if fI.ShareDirectory != "" {
		d.Set("share_directory", fI.ShareDirectory)
	}
	if fI.ContainerFolder != "" {
		d.Set("container_folder", fI.ContainerFolder)
	}
	if fI.SasUrlConfig != "" {
		d.Set("sas_url_config", fI.SasUrlConfig)
	}
	if fI.SasUriLicense != "" {
		d.Set("sas_url_license", fI.SasUriLicense)
	}
	if fI.UserData != "" {
		d.Set("user_data", fI.UserData)
	}

	return nil
}

func resourceAviatrixFirewallInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	firewallInstance := &goaviatrix.FirewallInstance{
		VpcID:      d.Get("vpc_id").(string),
		InstanceID: d.Get("instance_id").(string),
	}

	log.Printf("[INFO] Deleting firewall instance: %#v", firewallInstance)

	err := client.DeleteFirewallInstance(firewallInstance)
	if err != nil {
		return fmt.Errorf("failed to delete firewall instance: %s", err)
	}

	return nil
}
