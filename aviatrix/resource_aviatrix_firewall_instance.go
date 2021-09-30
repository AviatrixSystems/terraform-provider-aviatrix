package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixFirewallInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixFirewallInstanceCreate,
		Read:   resourceAviatrixFirewallInstanceRead,
		Update: resourceAviatrixFirewallInstanceUpdate,
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
				Optional:    true,
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
			"egress_vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Egress VPC ID. Required for GCP.",
			},
			"management_subnet": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "Management Interface Subnet. Required for Palo Alto Networks VM-Series, " +
					"and required to be empty for Check Point or Fortinet series.",
			},
			"management_vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "Management VPC ID. Required for GCP Palo Alto Networks VM-Series. " +
					"Required to be empty for GCP Check Point or Fortinet series.",
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
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Availability Zone. Only available for AWS, GCP and Azure.",
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
				Description: "Advanced option. Bootstrap bucket name. Only available for AWS and GCP.",
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
				DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
					return strings.TrimSpace(o) == strings.TrimSpace(n)
				},
			},
			"availability_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Availability domain for OCI.",
			},
			"fault_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Fault domain for OCI.",
			},
			"firewall_image_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Firewall image ID.",
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
			"gcp_vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GCP VPC ID",
			},
			"cloud_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Cloud Type",
			},
			"tags": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				ForceNew:    true,
				Description: "A map of tags to assign to the firewall instance.",
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
		FirewallImageId:      d.Get("firewall_image_id").(string),
		FirewallSize:         d.Get("firewall_size").(string),
		EgressSubnet:         d.Get("egress_subnet").(string),
		EgressVpc:            d.Get("egress_vpc_id").(string),
		ManagementSubnet:     d.Get("management_subnet").(string),
		ManagementVpc:        d.Get("management_vpc_id").(string),
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
		AvailabilityDomain:   d.Get("availability_domain").(string),
		FaultDomain:          d.Get("fault_domain").(string),
	}

	// For additional config validation we try to get the cloud_type from the given
	// gateway name or vpc_id. If there is an issue, we will just continue on without
	// the additional validation.
	var cloudType int
	if firewallInstance.GwName == "" {
		var err error
		cloudType, err = client.GetCloudTypeFromVpcID(firewallInstance.VpcID)
		if err != nil {
			log.Printf("[WARN] Could not get cloud_type from vpc_id: %v", err)
		}
	} else {
		gw, err := client.GetGateway(&goaviatrix.Gateway{GwName: firewallInstance.GwName})
		if err != nil {
			log.Printf("[WARN] Could not get cloud_type from firenet_gw_name: %v", err)
		} else {
			cloudType = gw.CloudType
		}
	}

	if strings.HasPrefix(firewallInstance.FirewallImage, "Palo Alto Networks") {
		if firewallInstance.ManagementSubnet == "" {
			return fmt.Errorf("'management_subnet' is required for Palo Alto Networks VM-Series")
		}
	} else if strings.Contains(firewallInstance.FirewallImage, "CloudGuard") {
		if goaviatrix.IsCloudType(cloudType, goaviatrix.OCIRelatedCloudTypes) && firewallInstance.ManagementSubnet == "" {
			return fmt.Errorf("'management_subnet' is required for Check Point CloudGuard for OCI")
		}
		if !goaviatrix.IsCloudType(cloudType, goaviatrix.OCIRelatedCloudTypes) && firewallInstance.ManagementSubnet != "" {
			return fmt.Errorf("'management_subnet' is required to be empty for Check Point CloudGuard except for OCI")
		}
	} else if strings.HasPrefix(firewallInstance.FirewallImage, "Fortinet FortiGate") {
		if firewallInstance.ManagementSubnet != "" {
			return fmt.Errorf("'management_subnet' is required to be empty for Fortinet FortiGate series")
		}
	}

	firenetDetail, err := client.GetFireNet(&goaviatrix.FireNet{VpcID: firewallInstance.VpcID})
	var isNativeGWLBVpc bool
	if err != nil {
		log.Printf("[INFO] Could not get FireNet detail for vpc_id(%s) because of (%v),"+
			" assuming this is a non-GWLB vpc", firewallInstance.VpcID, err)
	} else {
		isNativeGWLBVpc = firenetDetail.NativeGwlb
	}
	if isNativeGWLBVpc {
		if firewallInstance.GwName != "" {
			return fmt.Errorf("VPC %s has Native GWLB enabled but a 'firenet_gw_name' was provided. "+
				"Please remove 'firenet_gw_name' when using a Native GWLB enabled VPC", firewallInstance.VpcID)
		}
		if d.Get("zone") == "" {
			return fmt.Errorf("VPC %s has Native GWLB enabled but a 'zone' was not provided. "+
				"Please provide a 'zone' in your terraform config", firewallInstance.VpcID)
		}
	} else {
		if firewallInstance.GwName == "" {
			return fmt.Errorf("'firenet_gw_name' is required when using a non Native GWLB VPC. " +
				"Please provide a 'firenet_gw_name' in your terraform config")
		}
	}

	zone := d.Get("zone").(string)
	if zone != "" && !goaviatrix.IsCloudType(cloudType, goaviatrix.Azure|goaviatrix.AWS|goaviatrix.GCPRelatedCloudTypes) {
		return fmt.Errorf("'zone' attribute is only valid for AWS, GCP or Azure")
	}
	if zone != "" {
		firewallInstance.AvailabilityZone = zone
		if !goaviatrix.IsCloudType(cloudType, goaviatrix.GCPRelatedCloudTypes) {
			firewallInstance.EgressSubnet = fmt.Sprintf("%s~~%s~~", firewallInstance.EgressSubnet, zone)
			firewallInstance.ManagementSubnet = fmt.Sprintf("%s~~%s~~", firewallInstance.ManagementSubnet, zone)
		}
	}

	if !goaviatrix.IsCloudType(cloudType, goaviatrix.GCPRelatedCloudTypes) {
		if firewallInstance.ManagementVpc != "" {
			return fmt.Errorf("'management_vpc_id' is only valid for GCP")
		}
		if firewallInstance.EgressVpc != "" {
			return fmt.Errorf("'egress_vpc_id' is only valid for GCP")
		}
	} else if goaviatrix.IsCloudType(cloudType, goaviatrix.GCPRelatedCloudTypes) {
		if firewallInstance.ManagementVpc == "" && strings.HasPrefix(firewallInstance.FirewallImage, "Palo Alto Networks") {
			return fmt.Errorf("'management_vpc_id' is required for GCP with Palo Alto Networks Firewall")
		}
		if firewallInstance.ManagementVpc != "" && !strings.HasPrefix(firewallInstance.FirewallImage, "Palo Alto Networks") {
			return fmt.Errorf("'management_vpc_id' is required to be empty for GCP Check Point or FortiGate firewall")
		}
		if firewallInstance.EgressVpc == "" {
			return fmt.Errorf("'egress_vpc_id' is required for GCP")
		}
	}

	if firewallInstance.Username != "" || firewallInstance.Password != "" || firewallInstance.SshPublicKey != "" {
		if !goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			return fmt.Errorf("'username' and 'password' or 'ssh_public_key' are only supported for Azure")
		}
	}
	if firewallInstance.Password != "" && firewallInstance.SshPublicKey != "" {
		return fmt.Errorf("anthentication method can be either a password or an SSH public key. Please specify one of them and set the other one to empty")
	}
	if firewallInstance.IamRole != "" {
		if !goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes) {
			return fmt.Errorf("advanced option 'iam_role' is only supported for AWS provider, please set to empty")
		}
	}
	if firewallInstance.UserData != "" && strings.HasPrefix(firewallInstance.FirewallImage, "Palo Alto Networks") {
		return fmt.Errorf("advanced option of 'user_data' is only supported for Check Point Series and Fortinet FortiGate Series, not for %s", firewallInstance.FirewallImage)
	}
	if firewallInstance.StorageAccessKey != "" || firewallInstance.FileShareFolder != "" || firewallInstance.ShareDirectory != "" {
		if !strings.HasPrefix(firewallInstance.FirewallImage, "Palo Alto Networks") || !goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			return fmt.Errorf("advanced options of 'storage_access_key', 'file_share_folder' and 'share_directory' are only supported for Azure and Palo Alto Networks VM-Series")
		}
	}
	if firewallInstance.ContainerFolder != "" || firewallInstance.SasUrlConfig != "" || firewallInstance.SasUriLicense != "" {
		if !strings.HasPrefix(firewallInstance.FirewallImage, "Fortinet FortiGate") || !goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			return fmt.Errorf("advanced options of 'container_folder', 'sas_url_config' and 'sas_url_license' are only supported for Azure and Fortinet FortiGate series")
		}
	}
	if firewallInstance.BootstrapStorageName != "" {
		if strings.HasPrefix(firewallInstance.FirewallImage, "Check Point CloudGuard") || !goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			return fmt.Errorf("advanced option of 'bootstrap_storage_name' is only supported for Azure and Palo Alto Networks VM-Series/Fortinet FortiGate series")
		}
	}
	if firewallInstance.SicKey != "" {
		if !strings.HasPrefix(firewallInstance.FirewallImage, "Check Point CloudGuard") {
			return fmt.Errorf("advanced option of 'sic_key' is only supported for Check Point Series")
		}
	}

	if firewallInstance.FirewallImageId != "" && !goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
		return fmt.Errorf("'firewall_image_id' is only supported for AWS and Azure")
	}

	tags, err := extractTags(d, cloudType)
	if err != nil {
		return fmt.Errorf("error creating tags for firewall instance: %v", err)
	}
	tagJson, err := TagsMapToJson(tags)
	if err != nil {
		return fmt.Errorf("failed to add tags when creating firewall instance: %v", err)
	}
	firewallInstance.Tags = tags
	firewallInstance.TagJson = tagJson

	if goaviatrix.IsCloudType(cloudType, goaviatrix.OCIRelatedCloudTypes) && (firewallInstance.AvailabilityDomain == "" || firewallInstance.FaultDomain == "") {
		return fmt.Errorf("'availability_domain' and 'fault_domain' are required for OCI")
	}
	if !goaviatrix.IsCloudType(cloudType, goaviatrix.OCIRelatedCloudTypes) && (firewallInstance.AvailabilityDomain != "" || firewallInstance.FaultDomain != "") {
		return fmt.Errorf("'availability_domain' and 'fault_domain' are only valid for OCI")
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
	if instanceID == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no firewall names received. Import Id is %s", id)
		d.Set("instance_id", id)
		d.SetId(id)
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

	cloudType := goaviatrix.VendorToCloudType(fI.CloudVendor)

	d.Set("cloud_type", cloudType)
	if goaviatrix.IsCloudType(cloudType, goaviatrix.GCPRelatedCloudTypes) {
		d.Set("vpc_id", fI.FirenetVpc)
		d.Set("gcp_vpc_id", fI.VpcID)
	} else {
		d.Set("vpc_id", fI.VpcID)
	}
	d.Set("firenet_gw_name", fI.GwName)
	d.Set("firewall_name", fI.FirewallName)
	d.Set("firewall_image", fI.FirewallImage)
	d.Set("firewall_size", fI.FirewallSize)
	d.Set("instance_id", fI.InstanceID)
	if goaviatrix.IsCloudType(cloudType, goaviatrix.GCPRelatedCloudTypes) {
		d.Set("egress_subnet", fI.EgressSubnetID)
		d.Set("egress_vpc_id", fI.EgressVpc)
	} else {
		d.Set("egress_subnet", fI.EgressSubnet)
	}
	if strings.HasPrefix(fI.FirewallImage, "Palo Alto Networks") {
		if goaviatrix.IsCloudType(cloudType, goaviatrix.GCPRelatedCloudTypes) {
			d.Set("management_subnet", fI.ManagementSubnetID)
			d.Set("management_vpc_id", fI.ManagementVpc)
		} else {
			d.Set("management_subnet", fI.ManagementSubnet)
		}
	}

	if fI.AvailabilityZone != "" {
		if goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes) && fI.AvailabilityZone != "AvailabilitySet" {
			d.Set("zone", "az-"+fI.AvailabilityZone)
		} else if (goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes) && fI.GwName == "") || goaviatrix.IsCloudType(cloudType, goaviatrix.GCPRelatedCloudTypes) {
			d.Set("zone", fI.AvailabilityZone)
		}
	}

	d.Set("lan_interface", fI.LanInterface)
	d.Set("management_interface", fI.ManagementInterface)
	d.Set("egress_interface", fI.EgressInterface)
	d.Set("public_ip", fI.ManagementPublicIP)

	if fI.FirewallImageVersion != "" {
		d.Set("firewall_image_version", fI.FirewallImageVersion)
	}
	if fI.KeyFile == "" && goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes) {
		d.Set("key_name", fI.KeyName)
	}
	if fI.IamRole != "" {
		d.Set("iam_role", fI.IamRole)
	}
	if fI.BootstrapBucketName != "" {
		d.Set("bootstrap_bucket_name", fI.BootstrapBucketName)
	}
	if fI.Username != "" && !goaviatrix.IsCloudType(cloudType, goaviatrix.GCPRelatedCloudTypes) {
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
	if len(fI.Tags) > 0 {
		err := d.Set("tags", fI.Tags)
		if err != nil {
			return fmt.Errorf("failed to set tags for firewall_instance on read: %v", err)
		}
	}
	if fI.FirewallImageId != "" && goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
		d.Set("firewall_image_id", fI.FirewallImageId)
	}
	if goaviatrix.IsCloudType(cloudType, goaviatrix.OCIRelatedCloudTypes) {
		d.Set("availability_domain", fI.AvailabilityZone)
		d.Set("fault_domain", fI.FaultDomain)
	}

	return nil
}

func resourceAviatrixFirewallInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("firewall_image_id") {
		return fmt.Errorf("can not change firewall_image_id")
	}

	return resourceAviatrixFirewallInstanceRead(d, meta)
}

func resourceAviatrixFirewallInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	firewallInstance := &goaviatrix.FirewallInstance{
		VpcID:      d.Get("vpc_id").(string),
		InstanceID: d.Get("instance_id").(string),
	}
	if goaviatrix.IsCloudType(d.Get("cloud_type").(int), goaviatrix.GCPRelatedCloudTypes) {
		firewallInstance.VpcID = d.Get("gcp_vpc_id").(string)
	}

	log.Printf("[INFO] Deleting firewall instance: %#v", firewallInstance)

	err := client.DeleteFirewallInstance(firewallInstance)
	if err != nil {
		return fmt.Errorf("failed to delete firewall instance: %s", err)
	}

	return nil
}
