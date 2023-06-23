package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
)

func resourceAviatrixSpokeHaGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSpokeHaGatewayCreate,
		Read:   resourceAviatrixSpokeHaGatewayRead,
		Update: resourceAviatrixSpokeHaGatewayUpdate,
		Delete: resourceAviatrixSpokeHaGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
		},
	}
}

func resourceAviatrixSpokeHaGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gateway := &goaviatrix.SpokeHaGateway{
		PrimaryGwName:      d.Get("primary_gw_name").(string),
		GwName:             d.Get("gw_name").(string),
		GwSize:             d.Get("gw_size").(string),
		Subnet:             d.Get("subnet").(string),
		Zone:               d.Get("zone").(string),
		AvailabilityDomain: d.Get("availability_domain").(string),
		FaultDomain:        d.Get("fault_domain").(string),
		Eip:                d.Get("eip").(string),
	}

	primaryGw := &goaviatrix.Gateway{
		GwName: d.Get("primary_gw_name").(string),
	}
	gw, err := client.GetGateway(primaryGw)
	if err != nil {
		return fmt.Errorf("couldn't retrieve Aviatrix primary spoke gateway in spoke ha gateway creation: %s", err)
	}

	if gateway.GwName == "" {
		gateway.AutoGenHaGwName = "yes"
	}

	if d.Get("insane_mode").(bool) {
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
			gateway.Eip = fmt.Sprintf("%s:%s", azureEipName.(string), gateway.Eip)
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
			insaneModeAz := d.Get("insane_mode_az").(string)
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
		return fmt.Errorf("failed to create Aviatrix Spoke HA Gateway: %s", err)
	}

	d.SetId(spokeHaGwName)
	return resourceAviatrixSpokeHaGatewayRead(d, meta)
}

func resourceAviatrixSpokeHaGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	var isImport bool
	gwName := d.Get("gw_name").(string)
	if gwName == "" {
		isImport = true
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no gateway name received. Import Id is %s", id)
		d.Set("gw_name", id)
		d.SetId(id)
	}

	gateway := &goaviatrix.Gateway{
		AccountName: d.Get("account_name").(string),
		GwName:      d.Get("gw_name").(string),
	}

	gw, err := client.GetGateway(gateway)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix Spoke Gateway: %s", err)
	}

	log.Printf("[TRACE] reading spoke gateway %s: %#v", d.Get("gw_name").(string), gw)

	d.Set("primary_gw_name", gw.PrimaryGwName)
	d.Set("eip", gw.PublicIP)
	d.Set("subnet", gw.VpcNet)
	d.Set("gw_size", gw.GwSize)
	d.Set("cloud_type", gw.CloudType)
	d.Set("account_name", gw.AccountName)
	d.Set("cloud_instance_id", gw.CloudnGatewayInstID)
	d.Set("security_group_id", gw.GwSecurityGroupID)
	d.Set("private_ip", gw.PrivateIP)
	d.Set("public_ip", gw.PublicIP)
	d.Set("image_version", gw.ImageVersion)
	d.Set("software_version", gw.SoftwareVersion)

	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
		d.Set("vpc_reg", gw.VpcRegion)
	} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		d.Set("zone", gw.GatewayZone)
	} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		d.Set("vpc_reg", gw.VpcRegion)
		_, zoneIsSet := d.GetOk("zone")
		if (isImport || zoneIsSet) && gw.GatewayZone != "AvailabilitySet" {
			d.Set("zone", "az-"+gw.GatewayZone)
		}
		azureEip := strings.Split(gw.ReuseEip, ":")
		if len(azureEip) == 3 {
			d.Set("azure_eip_name_resource_group", fmt.Sprintf("%s:%s", azureEip[0], azureEip[1]))
		} else {
			log.Printf("[WARN] could not get Azure EIP name and resource group for the Spoke HA Gateway %s", gw.GwName)
		}
	} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.OCIRelatedCloudTypes) {
		d.Set("vpc_reg", gw.VpcRegion)
	} else if gw.CloudType == goaviatrix.AliCloud {
		d.Set("vpc_reg", gw.VpcRegion)
	}

	if gw.InsaneMode == "yes" {
		d.Set("insane_mode", true)
		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) {
			d.Set("insane_mode_az", gw.GatewayZone)
		}
	} else {
		d.Set("insane_mode", false)
	}

	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.OCIRelatedCloudTypes) {
		if gw.GatewayZone != "" {
			d.Set("availability_domain", gw.GatewayZone)
		} else {
			d.Set("availability_domain", d.Get("availability_domain").(string))
		}
		d.Set("fault_domain", gw.FaultDomain)
	}

	return nil
}

func resourceAviatrixSpokeHaGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gateway := &goaviatrix.Gateway{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string),
	}

	if d.HasChange("gw_size") {
		gateway.GwName = d.Get("gw_name").(string)
		gateway.VpcSize = d.Get("gw_size").(string)
		err := client.UpdateGateway(gateway)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix Spoke HA Gateway %s: %s", gateway.GwName, err)
		}
	}

	d.Partial(false)
	d.SetId(gateway.GwName)
	return resourceAviatrixSpokeHaGatewayRead(d, meta)
}

func resourceAviatrixSpokeHaGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gateway := &goaviatrix.Gateway{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix Spoke Ha Gateway: %#v", gateway)

	err := client.DeleteGateway(gateway)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Spoke HA Gateway %s : %s", gateway.GwName, err)
	}

	return nil
}
