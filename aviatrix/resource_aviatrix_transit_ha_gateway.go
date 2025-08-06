package aviatrix

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAviatrixTransitHaGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixTransitHaGatewayCreate,
		Read:   resourceAviatrixTransitHaGatewayRead,
		Update: resourceAviatrixTransitHaGatewayUpdate,
		Delete: resourceAviatrixTransitHaGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"primary_gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the primary transit gateway for which HA is being enabled.",
			},
			"ha_gw_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Name of the HA transit gateway. If not specified, defaults to <primary_gw_name>-hagw.",
			},
			"cloud_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				Description:  "Type of cloud service provider.",
				ValidateFunc: validateCloudType,
			},
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "This parameter represents the name of a Cloud-Account in Aviatrix controller.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC-ID/VNet-Name/Site-ID of cloud provider.",
			},
			"ha_subnet": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDR,
				Description: "HA Subnet. Required for enabling HA for AWS/AWSGov/AWSChina/Azure/OCI/Alibaba Cloud. " +
					"Optional for enabling HA for GCP gateway.",
			},
			"ha_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "",
				Description: "HA Zone. Required if enabling HA for GCP. Optional for Azure.",
			},
			"ha_insane_mode_az": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "",
				Description: "AZ of subnet being created for Insane Mode Transit HA Gateway. Required for AWS if insane_mode is enabled.",
			},
			"ha_gw_size": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "HA Gateway Size.",
			},
			"ha_device_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Device ID for HA AEP EAT gateway.",
			},
			"insane_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Enable Insane Mode for Transit HA Gateway. Valid values: true, false.",
			},
			"ha_oob_management_subnet": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDR,
				Description:  "OOB HA management subnet.",
			},
			"ha_oob_availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "OOB HA availability zone.",
			},
			"ha_availability_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "HA availability domain for OCI.",
			},
			"ha_fault_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "HA fault domain for OCI.",
			},
			"ha_eip": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPAddress,
				Description:  "Public IP address that you want assigned to the HA Transit Gateway.",
			},
			"ha_azure_eip_name_resource_group": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "The name of the public IP address and its resource group in Azure to assign to the HA Transit Gateway.",
				ValidateFunc: validateAzureEipNameResourceGroup,
			},
			"ha_software_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "ha_software_version can be used to set the desired software version of the HA gateway. " +
					"If set, we will attempt to update the gateway to the specified version.",
			},
			"ha_image_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "ha_image_version can be used to set the desired image version of the HA gateway. " +
					"If set, we will attempt to update the gateway to the specified version.",
			},
			"ha_bgp_lan_interfaces": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Interfaces to run BGP protocol on top of the ethernet interface, to connect to the onprem/remote peer. Only available for GCP HA Transit.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: DiffSuppressFuncGCPVpcId,
							Description:      "VPC-ID of GCP cloud provider.",
						},
						"subnet": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsCIDR,
							Description:  "Subnet Info.",
						},
					},
				},
			},
			"private_mode_subnet_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Private Mode HA subnet availability zone.",
			},
			"enable_encrypt_volume": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Enable encrypt gateway EBS volume. Only supported for AWS and AWSGov providers. Valid values: true, false. Default value: false.",
			},
			"customer_managed_keys": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "Customer managed key ID.",
			},
			// Edge Transit specific attributes
			"ztp_file_download_path": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "ZTP file download path for Edge transit gateways.",
			},
			"interfaces": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "WAN/LAN/MANAGEMENT interfaces for Edge transit gateways.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ifname": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Interface name.",
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"WAN", "LAN", "MANAGEMENT"}, false),
							Description:  "Interface type. Valid values: WAN, LAN, MANAGEMENT.",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Interface bandwidth in Mbps.",
						},
						"public_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Interface public IP.",
						},
						"tag": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Interface tag.",
						},
						"dhcp": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Interface DHCP.",
						},
						"cidr": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Interface CIDR.",
						},
						"gateway_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Interface gateway IP.",
						},
					},
				},
			},
			"interface_mapping": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Interface mapping for ESXI devices.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Interface name.",
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"MANAGEMENT", "WAN"}, false),
							Description:  "Interface type.",
						},
						"index": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Interface index.",
						},
					},
				},
			},
			"management_egress_ip_prefix_list": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Set of management egress IP prefixes for the HA gateway.",
			},
			// Computed attributes
			"ha_security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "HA security group used for the transit gateway.",
			},
			"ha_cloud_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud instance ID of HA transit gateway.",
			},
			"ha_private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private IP address of the HA transit gateway created.",
			},
		},
	}
}

func resourceAviatrixTransitHaGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	primaryGwName := d.Get("primary_gw_name").(string)
	cloudType := d.Get("cloud_type").(int)

	// Check if primary gateway exists
	primaryGateway := &goaviatrix.Gateway{
		GwName: primaryGwName,
	}

	_, err := client.GetGateway(primaryGateway)
	if err != nil {
		return fmt.Errorf("primary transit gateway %s does not exist: %s", primaryGwName, err)
	}

	// Determine HA gateway name
	haGwName := d.Get("ha_gw_name").(string)
	if haGwName == "" {
		haGwName = primaryGwName + "-hagw"
		d.Set("ha_gw_name", haGwName)
	}

	// Create HA gateway structure
	haSubnet := d.Get("ha_subnet").(string)
	haZone := d.Get("ha_zone").(string)
	insaneMode := d.Get("insane_mode").(bool)

	// Validate requirements based on cloud type
	if err := validateHaGatewayRequirements(d, cloudType); err != nil {
		return err
	}

	if goaviatrix.IsCloudType(cloudType, goaviatrix.EdgeRelatedCloudTypes) {
		return createEdgeTransitHaGateway(d, client, cloudType)
	}

	// Create standard cloud provider HA gateway
	transitHaGw := &goaviatrix.TransitHaGateway{
		PrimaryGwName: primaryGwName,
		GwName:        haGwName,
		AccountName:   d.Get("account_name").(string),
		CloudType:     cloudType,
		VpcID:         d.Get("vpc_id").(string),
		Subnet:        haSubnet,
		Zone:          haZone,
		GwSize:        d.Get("ha_gw_size").(string),
		Eip:           d.Get("ha_eip").(string),
		InsaneMode:    "no",
	}

	// Set insane mode
	if insaneMode {
		transitHaGw.InsaneMode = "yes"
		if goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes) {
			insaneModeAz := d.Get("ha_insane_mode_az").(string)
			if insaneModeAz == "" {
				return fmt.Errorf("ha_insane_mode_az is required when insane_mode is enabled for AWS")
			}
			transitHaGw.Subnet = fmt.Sprintf("%s~~%s", haSubnet, insaneModeAz)
		}
	}

	// Handle zone for Azure
	if goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes) && haZone != "" {
		transitHaGw.Subnet = fmt.Sprintf("%s~~%s~~", haSubnet, haZone)
	}

	// Handle OCI specific attributes
	if goaviatrix.IsCloudType(cloudType, goaviatrix.OCIRelatedCloudTypes) {
		transitHaGw.AvailabilityDomain = d.Get("ha_availability_domain").(string)
		transitHaGw.FaultDomain = d.Get("ha_fault_domain").(string)
	}

	// Handle Azure EIP
	if goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		if haAzureEipName, ok := d.GetOk("ha_azure_eip_name_resource_group"); ok && transitHaGw.Eip != "" {
			transitHaGw.Eip = fmt.Sprintf("%s:%s", haAzureEipName.(string), transitHaGw.Eip)
		}
	}

	// Handle BGP LAN interfaces for GCP
	if goaviatrix.IsCloudType(cloudType, goaviatrix.GCP) {
		var haBgpLanVpcID []string
		var haBgpLanSpecifySubnet []string
		for _, haBgpInterface := range d.Get("ha_bgp_lan_interfaces").([]interface{}) {
			item := haBgpInterface.(map[string]interface{})
			haBgpLanVpcID = append(haBgpLanVpcID, item["vpc_id"].(string))
			haBgpLanSpecifySubnet = append(haBgpLanSpecifySubnet, item["subnet"].(string))
		}
		transitHaGw.BgpLanVpcID = strings.Join(haBgpLanVpcID, ",")
		transitHaGw.BgpLanSubnet = strings.Join(haBgpLanSpecifySubnet, ",")
	}

	// Handle private mode
	privateModeInfo, _ := client.GetPrivateModeInfo(context.Background())
	if privateModeInfo.EnablePrivateMode {
		if privateModeSubnetZone := d.Get("private_mode_subnet_zone").(string); privateModeSubnetZone != "" {
			transitHaGw.Subnet = fmt.Sprintf("%s~~%s", transitHaGw.Subnet, privateModeSubnetZone)
		}
	}

	// Handle OOB management
	if enablePrivateOob := d.Get("enable_private_oob").(bool); enablePrivateOob {
		haOobManagementSubnet := d.Get("ha_oob_management_subnet").(string)
		haOobAvailabilityZone := d.Get("ha_oob_availability_zone").(string)
		if haOobManagementSubnet != "" && haOobAvailabilityZone != "" {
			transitHaGw.Subnet = transitHaGw.Subnet + "~~" + haOobAvailabilityZone
		}
	}

	log.Printf("[INFO] Creating Aviatrix Transit HA Gateway: %#v", transitHaGw)

	d.SetId(haGwName)

	// Create the HA gateway
	_, err = client.CreateTransitHaGw(transitHaGw)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Transit HA Gateway: %s", err)
	}

	// Update gateway size if different from primary
	if transitHaGw.GwSize != "" {
		haGateway := &goaviatrix.Gateway{
			CloudType: cloudType,
			GwName:    haGwName,
			VpcSize:   transitHaGw.GwSize,
		}

		log.Printf("[INFO] Updating Transit HA Gateway size to: %s", haGateway.VpcSize)

		err = client.UpdateGateway(haGateway)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix Transit HA Gateway size: %s", err)
		}
	}

	return resourceAviatrixTransitHaGatewayRead(d, meta)
}

func resourceAviatrixTransitHaGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	haGwName := d.Id()
	if haGwName == "" {
		return fmt.Errorf("invalid empty gateway name")
	}

	gateway := &goaviatrix.Gateway{
		GwName: haGwName,
	}

	gw, err := client.GetGateway(gateway)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix Transit HA Gateway: %s", err)
	}

	log.Printf("[TRACE] Reading Aviatrix Transit HA Gateway: %s", gw.GwName)

	d.Set("ha_gw_name", gw.GwName)
	d.Set("cloud_type", gw.CloudType)
	d.Set("account_name", gw.AccountName)
	d.Set("vpc_id", gw.VpcID)
	d.Set("ha_gw_size", gw.VpcSize)
	d.Set("ha_private_ip", gw.PrivateIP)
	d.Set("ha_cloud_instance_id", gw.CloudnGatewayInstID)
	d.Set("ha_security_group_id", gw.GwSecurityGroupID)

	if gw.InsaneMode == "yes" {
		d.Set("insane_mode", true)
	} else {
		d.Set("insane_mode", false)
	}

	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
		d.Set("ha_zone", gw.GatewayZone)
		if gw.Eip != "" {
			parts := strings.Split(gw.Eip, ":")
			if len(parts) == 3 {
				d.Set("ha_azure_eip_name_resource_group", parts[0]+":"+parts[1])
				d.Set("ha_eip", parts[2])
			}
		}
	} else {
		d.Set("ha_eip", gw.Eip)
	}

	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		d.Set("ha_zone", gw.GatewayZone)
	}

	if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.OCIRelatedCloudTypes) {
		d.Set("ha_availability_domain", gw.AvailabilityDomain)
		d.Set("ha_fault_domain", gw.FaultDomain)
	}

	// Parse subnet information
	if gw.VpcNet != "" {
		if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AWSRelatedCloudTypes) && gw.InsaneMode == "yes" {
			// For AWS insane mode, subnet format is "subnet~~az"
			parts := strings.Split(gw.VpcNet, "~~")
			if len(parts) >= 1 {
				d.Set("ha_subnet", parts[0])
			}
			if len(parts) >= 2 {
				d.Set("ha_insane_mode_az", parts[1])
			}
		} else if goaviatrix.IsCloudType(gw.CloudType, goaviatrix.AzureArmRelatedCloudTypes) {
			// For Azure, subnet format might be "subnet~~zone~~"
			parts := strings.Split(gw.VpcNet, "~~")
			if len(parts) >= 1 {
				d.Set("ha_subnet", parts[0])
			}
		} else {
			d.Set("ha_subnet", gw.VpcNet)
		}
	}

	return nil
}

func resourceAviatrixTransitHaGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	haGwName := d.Id()
	log.Printf("[INFO] Updating Aviatrix Transit HA Gateway: %s", haGwName)

	// Handle size changes
	if d.HasChange("ha_gw_size") {
		gateway := &goaviatrix.Gateway{
			CloudType: d.Get("cloud_type").(int),
			GwName:    haGwName,
			VpcSize:   d.Get("ha_gw_size").(string),
		}

		err := client.UpdateGateway(gateway)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix Transit HA Gateway size: %s", err)
		}
	}

	// Handle software version updates
	if d.HasChange("ha_software_version") {
		gateway := &goaviatrix.Gateway{
			GwName:          haGwName,
			SoftwareVersion: d.Get("ha_software_version").(string),
		}

		err := client.UpdateGateway(gateway)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix Transit HA Gateway software version: %s", err)
		}
	}

	// Handle image version updates
	if d.HasChange("ha_image_version") {
		gateway := &goaviatrix.Gateway{
			GwName:       haGwName,
			ImageVersion: d.Get("ha_image_version").(string),
		}

		err := client.UpdateGateway(gateway)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix Transit HA Gateway image version: %s", err)
		}
	}

	return resourceAviatrixTransitHaGatewayRead(d, meta)
}

func resourceAviatrixTransitHaGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	haGwName := d.Id()

	log.Printf("[INFO] Deleting Aviatrix Transit HA Gateway: %s", haGwName)

	gateway := &goaviatrix.Gateway{
		GwName: haGwName,
	}

	err := client.DeleteGateway(gateway)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Transit HA Gateway: %s", err)
	}

	return nil
}

// Helper functions
func validateHaGatewayRequirements(d *schema.ResourceData, cloudType int) error {
	haSubnet := d.Get("ha_subnet").(string)
	haZone := d.Get("ha_zone").(string)
	insaneMode := d.Get("insane_mode").(bool)
	haAvailabilityDomain := d.Get("ha_availability_domain").(string)
	haFaultDomain := d.Get("ha_fault_domain").(string)

	// GCP requirements
	if goaviatrix.IsCloudType(cloudType, goaviatrix.GCPRelatedCloudTypes) && haZone == "" {
		return fmt.Errorf("'ha_zone' must be set to enable HA on GCP")
	}

	// Azure requirements
	if goaviatrix.IsCloudType(cloudType, goaviatrix.AzureArmRelatedCloudTypes) && haSubnet == "" {
		return fmt.Errorf("'ha_subnet' must be provided to enable HA on Azure")
	}

	// AWS insane mode requirements
	if goaviatrix.IsCloudType(cloudType, goaviatrix.AWSRelatedCloudTypes) && insaneMode {
		if d.Get("ha_insane_mode_az").(string) == "" {
			return fmt.Errorf("ha_insane_mode_az needed if insane_mode is enabled for AWS")
		}
	}

	// OCI requirements
	if goaviatrix.IsCloudType(cloudType, goaviatrix.OCIRelatedCloudTypes) {
		if haAvailabilityDomain == "" || haFaultDomain == "" {
			return fmt.Errorf("'ha_availability_domain' and 'ha_fault_domain' are required to enable HA on OCI")
		}
	}

	return nil
}

func createEdgeTransitHaGateway(d *schema.ResourceData, client *goaviatrix.Client, cloudType int) error {
	primaryGwName := d.Get("primary_gw_name").(string)
	haGwName := d.Get("ha_gw_name").(string)
	if haGwName == "" {
		haGwName = primaryGwName + "-hagw"
		d.Set("ha_gw_name", haGwName)
	}

	transitHaGw := &goaviatrix.TransitHaGateway{
		PrimaryGwName:       primaryGwName,
		GwName:              haGwName,
		VpcID:               d.Get("vpc_id").(string),
		ZtpFileDownloadPath: d.Get("ztp_file_download_path").(string),
		CloudType:           cloudType,
		InsaneMode:          "yes",
	}

	// Handle interfaces for Edge transit
	if haInterfaces, ok := d.GetOk("interfaces"); ok {
		interfacesList, err := getEdgeInterfaceDetails(haInterfaces.([]interface{}), cloudType)
		if err != nil {
			return fmt.Errorf("failed to get the interface details for HA Edge Transit Gateway: %w", err)
		}
		transitHaGw.Interfaces = interfacesList
	}

	// Handle interface mapping for AEP
	if goaviatrix.IsCloudType(cloudType, goaviatrix.EDGENEO) {
		interfaceMappingInput := d.Get("interface_mapping").([]interface{})
		interfaceMapping, err := getEdgeInterfaceMappingDetails(interfaceMappingInput)
		if err != nil {
			return fmt.Errorf("failed to get the interface mapping details: %w", err)
		}
		transitHaGw.InterfaceMapping = interfaceMapping

		// Set device_id for AEP ha gateway
		if deviceID, ok := d.GetOk("ha_device_id"); ok {
			transitHaGw.DeviceID = deviceID.(string)
		} else {
			return fmt.Errorf("ha_device_id is required for AEP HA Edge Transit Gateway")
		}
	}

	// Handle management egress IP prefix list
	haManagementEgressIPPrefixList := getStringSet(d, "management_egress_ip_prefix_list")
	if len(haManagementEgressIPPrefixList) > 0 {
		transitHaGw.ManagementEgressIPPrefix = strings.Join(haManagementEgressIPPrefixList, ",")
	}

	log.Printf("[INFO] Creating Edge Transit HA Gateway: %#v", transitHaGw)

	// Create the HA gateway
	_, err := client.CreateTransitHaGw(transitHaGw)
	if err != nil {
		return fmt.Errorf("failed to create Edge Transit HA Gateway: %s", err)
	}

	return nil
}

func getEdgeInterfaceDetails(interfaces []interface{}, cloudType int) (string, error) {
	var interfaceDetailsList []map[string]interface{}

	for _, interfaceItem := range interfaces {
		interfaceMap, ok := interfaceItem.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("invalid interface configuration")
		}

		interfaceDetails := map[string]interface{}{
			"ifname": interfaceMap["ifname"].(string),
			"type":   interfaceMap["type"].(string),
		}

		if bandwidth, ok := interfaceMap["bandwidth"]; ok {
			interfaceDetails["bandwidth"] = bandwidth.(int)
		}

		if publicIP, ok := interfaceMap["public_ip"]; ok && publicIP.(string) != "" {
			interfaceDetails["public_ip"] = publicIP.(string)
		}

		if tag, ok := interfaceMap["tag"]; ok && tag.(string) != "" {
			interfaceDetails["tag"] = tag.(string)
		}

		if dhcp, ok := interfaceMap["dhcp"]; ok {
			interfaceDetails["dhcp"] = dhcp.(bool)
		}

		if cidr, ok := interfaceMap["cidr"]; ok && cidr.(string) != "" {
			interfaceDetails["cidr"] = cidr.(string)
		}

		if gatewayIP, ok := interfaceMap["gateway_ip"]; ok && gatewayIP.(string) != "" {
			interfaceDetails["gateway_ip"] = gatewayIP.(string)
		}

		interfaceDetailsList = append(interfaceDetailsList, interfaceDetails)
	}

	interfacesJSON, err := json.Marshal(interfaceDetailsList)
	if err != nil {
		return "", fmt.Errorf("failed to marshal interfaces to JSON: %w", err)
	}

	return b64.StdEncoding.EncodeToString(interfacesJSON), nil
}

func getEdgeInterfaceMappingDetails(interfaceMappingInput []interface{}) (string, error) {
	interfaceMapping := map[string][]string{}

	if len(interfaceMappingInput) > 0 {
		// Set the interface mapping for ESXI devices
		for _, value := range interfaceMappingInput {
			mappingMap, ok := value.(map[string]interface{})
			if !ok {
				return "", fmt.Errorf("invalid type %T for interface mapping, expected a map", value)
			}
			interfaceName, ok1 := mappingMap["name"].(string)
			interfaceType, ok2 := mappingMap["type"].(string)
			interfaceIndex, ok3 := mappingMap["index"].(int)
			if !ok1 || !ok2 || !ok3 {
				return "", fmt.Errorf("invalid interface mapping, 'name', 'type', and 'index' must be provided")
			}

			var updatedInterfaceType string
			if interfaceType == "MANAGEMENT" {
				updatedInterfaceType = "mgmt"
			} else if interfaceType == "WAN" {
				updatedInterfaceType = "wan"
			} else {
				return "", fmt.Errorf("invalid interface type %s", interfaceType)
			}

			interfaceMapping[interfaceName] = []string{updatedInterfaceType, strconv.Itoa(interfaceIndex)}
		}
	} else {
		// Set the interface mapping for Dell devices
		interfaceMapping["eth0"] = []string{"mgmt", "0"}
		interfaceMapping["eth5"] = []string{"wan", "0"}
		interfaceMapping["eth2"] = []string{"wan", "1"}
		interfaceMapping["eth3"] = []string{"wan", "2"}
		interfaceMapping["eth4"] = []string{"wan", "3"}
	}

	// Convert interfaceMapping to JSON byte slice
	interfaceMappingJSON, err := json.Marshal(interfaceMapping)
	if err != nil {
		return "", fmt.Errorf("failed to marshal interface mapping to json: %w", err)
	}

	return string(interfaceMappingJSON), nil
}
