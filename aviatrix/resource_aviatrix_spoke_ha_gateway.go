package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixSpokeHaGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSpokeHaGatewayCreate,
		Read:   resourceAviatrixSpokeHaGatewayRead,
		Update: resourceAviatrixSpokeHaGatewayUpdate,
		Delete: resourceAviatrixSpokeHaGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},

		Schema: map[string]*schema.Schema{
			"primary_gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the primary gateway.",
			},
			"gw_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Name of the HA gateway which is going to be created.",
			},
			"gw_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Size of the gateway instance.",
			},
			"subnet": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDR,
				Description:  "Public Subnet Info.",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Availability Zone. Required for GCP gateway, example: 'us-west1-c'. Optional for Azure / Azure GOV / Azure CHINA gateway in the form 'az-n', example: 'az-2'.",
			},
			"insane_mode": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
				Description: "Enable Insane Mode for Spoke Gateway. Valid values: true, false. Supported for AWS/AWSGov, GCP, Azure and OCI. " +
					"If insane mode is enabled, gateway size has to at least be c5 size for AWS and Standard_D3_v2 size for Azure.",
			},
			"insane_mode_az": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				ForceNew:    true,
				Description: "AZ of subnet being created for Insane Mode Spoke Gateway. Required if insane_mode is enabled for AWS cloud.",
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
			"eip": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPAddress,
				Description:  "If set, the specified EIP is used for this gateway.",
			},
			"azure_eip_name_resource_group": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "The name of the public IP address and its resource group in Azure to assign to this Gateway.",
				ValidateFunc: validateAzureEipNameResourceGroup,
				RequiredWith: []string{"eip"},
			},
			"cloud_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Type of cloud service provider.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This parameter represents the name of a Cloud-Account in Aviatrix controller.",
			},
			"software_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Software version of the gateway.",
			},
			"image_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Image version of the gateway.",
			},
			"vpc_reg": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region of cloud provider.",
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Security group used for the spoke ha gateway.",
			},
			"cloud_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud instance ID.",
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private IP address of the spoke ha gateway created.",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP address of the spoke ha gateway created.",
			},
			"single_az_ha": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set to true if this feature is desired.",
			},
		},
	}
}

func resourceAviatrixSpokeHaGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	gateway := &goaviatrix.SpokeHaGateway{
		PrimaryGwName:      getString(d, "primary_gw_name"),
		GwName:             getString(d, "gw_name"),
		GwSize:             getString(d, "gw_size"),
		Subnet:             getString(d, "subnet"),
		Zone:               getString(d, "zone"),
		AvailabilityDomain: getString(d, "availability_domain"),
		FaultDomain:        getString(d, "fault_domain"),
		Eip:                getString(d, "eip"),
	}

	primaryGw := &goaviatrix.Gateway{
		GwName: getString(d, "primary_gw_name"),
	}
	gw, err := client.GetGateway(primaryGw)
	if err != nil {
		return fmt.Errorf("couldn't retrieve Aviatrix primary spoke gateway in spoke ha gateway creation: %w", err)
	}

	if gateway.GwName == "" {
		gateway.AutoGenHaGwName = "yes"
	}

	if getBool(d, "insane_mode") {
		gateway.InsaneMode = "yes"
	} else {
		gateway.InsaneMode = "no"
	}

	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		if gateway.Zone != "" || gateway.AvailabilityDomain != "" || gateway.FaultDomain != "" {
			return fmt.Errorf("'zone', 'availability_domain' and 'fault_domain' are required to be empty for creating an AWS related cloud type spoke ha gateway")
		}
	}

	azureEipName, azureEipNameOk := d.GetOk("azure_eip_name_resource_group")
	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		if gateway.AvailabilityDomain != "" || gateway.FaultDomain != "" {
			return fmt.Errorf("'availability_domain' and 'fault_domain' are required to be empty for creating an Azure related cloud type spoke ha gateway")
		}
		if gateway.Zone != "" {
			gateway.Subnet = fmt.Sprintf("%s~~%s~~", gateway.Subnet, gateway.Zone)
		}
		if gateway.Eip != "" {
			if !azureEipNameOk {
				return fmt.Errorf("'azure_eip_name_resource_group' must be set when 'eip' is set for Azure (8), AzureGov (32) or AzureChina (2048)")
			}
			gateway.Eip = fmt.Sprintf("%s:%s", mustString(azureEipName), gateway.Eip)
		}
	} else {
		if azureEipNameOk {
			return fmt.Errorf("'azure_eip_name_resource_group' only supports Azure clouds including Azure (8), AzureGov (32) or AzureChina (2048)")
		}
	}

	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		if gateway.Zone == "" {
			return fmt.Errorf("'zone' is required for creating an GCP spoke ha gateway")
		}
		if gateway.AvailabilityDomain != "" || gateway.FaultDomain != "" {
			return fmt.Errorf("'availability_domain' and 'fault_domain' are required to be empty for creating an GCP related cloud type spoke ha gateway")
		}
	}

	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.OCIRelatedCloudTypes) {
		if gateway.AvailabilityDomain == "" || gateway.FaultDomain == "" {
			return fmt.Errorf("'availability_domain' and 'fault_domain' are required for creating an OCI related cloud type spoke ha gateway")
		}
		if gateway.Zone != "" {
			return fmt.Errorf("'zone' is required to be empty for creating an OCI related cloud type spoke ha gateway")
		}
	}

	if gateway.InsaneMode == "yes" {
		if !goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.GCPRelatedCloudTypes|
			goaviatrix.AzureArmRelatedCloudTypes|goaviatrix.OCIRelatedCloudTypes) {
			return fmt.Errorf("insane_mode is only supported for AWS (1), GCP (4), Azure (8), OCI (16), AzureGov (32), AWSGov (256), AWS China (1024), AzureChina (2048), AWS Top Secret (16384) and AWS Secret (32768)")
		}

		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			insaneModeAz := getString(d, "insane_mode_az")
			if insaneModeAz == "" {
				return fmt.Errorf("'insane_mode_az' is required if insane_mode is enabled for AWS (1), AWSGov (256), AWS China (1024), AWS Top Secret (16384) or AWS Secret (32768)")
			}
			var insaneModeSubnet []string
			insaneModeSubnet = append(insaneModeSubnet, gateway.Subnet, insaneModeAz)
			gateway.Subnet = strings.Join(insaneModeSubnet, "~~")
		}
	}

	spokeHaGwName, err := client.CreateSpokeHaGw(gateway)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Spoke HA Gateway: %w", err)
	}

	// Handle single_az_ha setting after gateway creation
	// Need to differentiate between: not set (inherit from primary) vs explicitly set false
	// d.GetOk() doesn't seem to differentiate between explicit set to false vs not set at all
	rawConfig := d.GetRawConfig()
	if !rawConfig.GetAttr("single_az_ha").IsNull() {
		// User explicitly set single_az_ha in config
		singleAZ := getBool(d, "single_az_ha")
		singleAZGateway := &goaviatrix.Gateway{
			GwName: spokeHaGwName,
		}

		if singleAZ {
			log.Printf("[INFO] Enable Single AZ GW HA for spoke HA gateway: %#v", singleAZGateway)
			err := client.EnableSingleAZGateway(singleAZGateway)
			if err != nil {
				return fmt.Errorf("failed to enable single AZ GW HA for spoke HA gateway %s: %w", spokeHaGwName, err)
			}
		} else {
			log.Printf("[INFO] Disable Single AZ GW HA for spoke HA gateway: %#v", singleAZGateway)
			err := client.DisableSingleAZGateway(singleAZGateway)
			if err != nil {
				return fmt.Errorf("failed to disable single AZ GW HA for spoke HA gateway %s: %w", spokeHaGwName, err)
			}
		}
	} else {
		log.Printf("[INFO] single_az_ha not set - inheriting from primary gateway")
	}

	d.SetId(spokeHaGwName)
	return resourceAviatrixSpokeHaGatewayRead(d, meta)
}

func resourceAviatrixSpokeHaGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	var isImport bool
	gwName := getString(d, "gw_name")
	if gwName == "" {
		isImport = true
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no gateway name received. Import Id is %s", id)
		mustSet(d, "gw_name", id)
		d.SetId(id)
	}

	gateway := &goaviatrix.Gateway{
		AccountName: getString(d, "account_name"),
		GwName:      getString(d, "gw_name"),
	}

	gw, err := client.GetGateway(gateway)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix Spoke Gateway: %w", err)
	}

	log.Printf("[TRACE] reading spoke gateway %s: %#v", getString(d, "gw_name"), gw)
	mustSet(d, "primary_gw_name", gw.PrimaryGwName)
	mustSet(d, "eip", gw.PublicIP)
	mustSet(d, "subnet", gw.VpcNet)
	mustSet(d, "gw_size", gw.GwSize)
	mustSet(d, "cloud_type", gw.CloudType)
	mustSet(d, "account_name", gw.AccountName)
	mustSet(d, "cloud_instance_id", gw.CloudnGatewayInstID)
	mustSet(d, "security_group_id", gw.GwSecurityGroupID)
	mustSet(d, "private_ip", gw.PrivateIP)
	mustSet(d, "public_ip", gw.PublicIP)
	mustSet(d, "image_version", gw.ImageVersion)
	mustSet(d, "software_version", gw.SoftwareVersion)
	mustSet(d, "single_az_ha", gw.SingleAZ == "yes")

	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		mustSet(d, "vpc_reg", gw.VpcRegion)
	} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		mustSet(d, "zone", gw.GatewayZone)
	} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		mustSet(d, "vpc_reg", gw.VpcRegion)
		_, zoneIsSet := d.GetOk("zone")
		if (isImport || zoneIsSet) && gw.GatewayZone != "AvailabilitySet" {
			mustSet(d, "zone", "az-"+gw.GatewayZone)
		}
		azureEip := strings.Split(gw.ReuseEip, ":")
		if len(azureEip) == 3 {
			mustSet(d, "azure_eip_name_resource_group", fmt.Sprintf("%s:%s", azureEip[0], azureEip[1]))
		} else {
			log.Printf("[WARN] could not get Azure EIP name and resource group for the Spoke HA Gateway %s", gw.GwName)
		}
	} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.OCIRelatedCloudTypes) {
		mustSet(d, "vpc_reg", gw.VpcRegion)
	} else if gw.CloudType == goaviatrix.AliCloud {
		mustSet(d, "vpc_reg", gw.VpcRegion)
	}

	if gw.InsaneMode == "yes" {
		mustSet(d, "insane_mode", true)
		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			mustSet(d, "insane_mode_az", gw.GatewayZone)
		}
	} else {
		mustSet(d, "insane_mode", false)
	}

	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.OCIRelatedCloudTypes) {
		if gw.GatewayZone != "" {
			mustSet(d, "availability_domain", gw.GatewayZone)
		} else {
			mustSet(d, "availability_domain", getString(d, "availability_domain"))
		}
		mustSet(d, "fault_domain", gw.FaultDomain)
	}

	return nil
}

func resourceAviatrixSpokeHaGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	gateway := &goaviatrix.Gateway{
		CloudType: getInt(d, "cloud_type"),
		GwName:    getString(d, "gw_name"),
	}

	if d.HasChange("gw_size") {
		gateway.GwName = getString(d, "gw_name")
		gateway.VpcSize = getString(d, "gw_size")
		err := client.UpdateGateway(gateway)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix Spoke HA Gateway %s: %w", gateway.GwName, err)
		}
	}

	if d.HasChange("single_az_ha") {
		singleAZGateway := &goaviatrix.Gateway{
			GwName: getString(d, "gw_name"),
		}

		singleAZ := getBool(d, "single_az_ha")
		if singleAZ {
			singleAZGateway.SingleAZ = "yes"
		} else {
			singleAZGateway.SingleAZ = "no"
		}

		if singleAZ {
			log.Printf("[INFO] Enable Single AZ GW HA: %#v", singleAZGateway)
			err := client.EnableSingleAZGateway(singleAZGateway)
			if err != nil {
				return fmt.Errorf("failed to enable single AZ GW HA for %s: %w", singleAZGateway.GwName, err)
			}
		} else {
			log.Printf("[INFO] Disable Single AZ GW HA: %#v", singleAZGateway)
			err := client.DisableSingleAZGateway(singleAZGateway)
			if err != nil {
				return fmt.Errorf("failed to disable single AZ GW HA for %s: %w", singleAZGateway.GwName, err)
			}
		}
	}

	d.Partial(false)
	d.SetId(gateway.GwName)
	return resourceAviatrixSpokeHaGatewayRead(d, meta)
}

func resourceAviatrixSpokeHaGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	gateway := &goaviatrix.Gateway{
		CloudType: getInt(d, "cloud_type"),
		GwName:    getString(d, "gw_name"),
	}

	log.Printf("[INFO] Deleting Aviatrix Spoke Ha Gateway: %#v", gateway)

	err := client.DeleteGateway(gateway)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Spoke HA Gateway %s : %w", gateway.GwName, err)
	}

	return nil
}
