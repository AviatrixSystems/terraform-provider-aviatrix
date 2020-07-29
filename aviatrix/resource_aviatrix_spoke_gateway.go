package aviatrix

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixSpokeGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSpokeGatewayCreate,
		Read:   resourceAviatrixSpokeGatewayRead,
		Update: resourceAviatrixSpokeGatewayUpdate,
		Delete: resourceAviatrixSpokeGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cloud_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Type of cloud service provider.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This parameter represents the name of a Cloud-Account in Aviatrix controller.",
			},
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the gateway which is going to be created.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC-ID/VNet-Name of cloud provider.",
			},
			"vpc_reg": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Region of cloud provider.",
			},
			"gw_size": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Size of the gateway instance.",
			},
			"subnet": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Public Subnet Info.",
			},
			"insane_mode_az": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "AZ of subnet being created for Insane Mode Spoke Gateway. Required if insane_mode is enabled for aws cloud.",
			},
			"single_ip_snat": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Specify whether to enable Source NAT feature in 'single_ip' mode on the gateway or not.",
			},
			"allocate_new_eip": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				Description: "If false, reuse an idle address in Elastic IP pool for this gateway. " +
					"Otherwise, allocate a new Elastic IP and use it for this gateway.",
			},
			"eip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Required when allocate_new_eip is false. It uses specified EIP for this gateway.",
			},
			"ha_subnet": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "HA Subnet. Required if enabling HA for AWS/AZURE. Optional if enabling HA for GCP.",
			},
			"ha_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "HA Zone. Required if enabling HA for GCP.",
			},
			"ha_insane_mode_az": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "AZ of subnet being created for Insane Mode Spoke HA Gateway. Required if insane_mode is true and ha_subnet is set.",
			},
			"ha_gw_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "HA Gateway Size.",
			},
			"ha_eip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Public IP address that you want assigned to the HA Spoke Gateway.",
			},
			"single_az_ha": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Set to 'enabled' if this feature is desired.",
			},
			"transit_gw": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Specify the transit Gateway.",
			},
			"tag_list": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Default:     nil,
				Description: "Instance tag of cloud provider.",
			},
			"insane_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable Insane Mode for Spoke Gateway. Valid values: true, false. If insane mode is enabled, gateway size has to at least be c5 size.",
			},
			"enable_active_mesh": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Switch to Enable/Disable Active Mesh Mode for Spoke Gateway. Valid values: true, false.",
			},
			"enable_vpc_dns_server": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable vpc_dns_server for Gateway. Only supports AWS. Valid values: true, false.",
			},
			"enable_encrypt_volume": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable encrypt gateway EBS volume. Only supported for AWS provider. Valid values: true, false. Default value: false.",
			},
			"customized_spoke_vpc_routes": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				Description: "A list of comma separated CIDRs to be customized for the spoke VPC routes. When configured, " +
					"it will replace all learned routes in VPC routing tables, including RFC1918 and non-RFC1918 CIDRs. " +
					"It applies to this spoke gateway only​.",
			},
			"filtered_spoke_vpc_routes": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				Description: "A list of comma separated CIDRs to be filtered from the spoke VPC route table. When configured, " +
					"filtering CIDR(s) or it’s subnet will be deleted from VPC routing tables as well as from spoke gateway’s " +
					"routing table. It applies to this spoke gateway only.",
			},
			"included_advertised_spoke_routes": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A list of comma separated CIDRs to be advertised to on-prem as 'Included CIDR List'. When configured, it will replace all advertised routes from this VPC.​",
			},
			"customer_managed_keys": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Customer managed key ID.",
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Security group used for the spoke gateway.",
			},
			"cloud_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud instance ID.",
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private IP address of the spoke gateway created.",
			},
			"ha_cloud_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud instance ID of HA spoke gateway.",
			},
			"ha_gw_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Aviatrix spoke gateway unique name of HA spoke gateway.",
			},
			"ha_private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private IP address of the spoke gateway created.",
			},
		},
	}
}

func resourceAviatrixSpokeGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gateway := &goaviatrix.SpokeVpc{
		CloudType:      d.Get("cloud_type").(int),
		AccountName:    d.Get("account_name").(string),
		GwName:         d.Get("gw_name").(string),
		VpcSize:        d.Get("gw_size").(string),
		Subnet:         d.Get("subnet").(string),
		HASubnet:       d.Get("ha_subnet").(string),
		TransitGateway: d.Get("transit_gw").(string),
	}
	enableSNat := d.Get("single_ip_snat").(bool)
	if enableSNat {
		gateway.EnableNat = "yes"
	}

	singleAZ := d.Get("single_az_ha").(bool)
	if singleAZ {
		gateway.SingleAzHa = "enabled"
	} else {
		gateway.SingleAzHa = "disabled"
	}

	allocateNewEip := d.Get("allocate_new_eip").(bool)
	if allocateNewEip {
		gateway.ReuseEip = "off"
	} else {
		gateway.ReuseEip = "on"
		gateway.Eip = d.Get("eip").(string)
	}

	if gateway.CloudType == goaviatrix.AWS || gateway.CloudType == goaviatrix.GCP || gateway.CloudType == goaviatrix.OCI {
		gateway.VpcID = d.Get("vpc_id").(string)
		if gateway.VpcID == "" {
			return fmt.Errorf("'vpc_id' cannot be empty for creating a spoke gw")
		}
	} else if gateway.CloudType == goaviatrix.AZURE {
		gateway.VNetNameResourceGroup = d.Get("vpc_id").(string)
		if gateway.VNetNameResourceGroup == "" {
			return fmt.Errorf("'vpc_id' cannot be empty for creating a spoke gw")
		}
	} else {
		return fmt.Errorf("invalid cloud type, it can only be aws (1), gcp (4), azure (8)")
	}

	if gateway.CloudType == goaviatrix.AWS || gateway.CloudType == goaviatrix.AZURE || gateway.CloudType == goaviatrix.OCI {
		gateway.VpcRegion = d.Get("vpc_reg").(string)
	} else if gateway.CloudType == goaviatrix.GCP {
		// for gcp, rest api asks for "zone" rather than vpc region
		gateway.Zone = d.Get("vpc_reg").(string)
	} else {
		return fmt.Errorf("invalid cloud type, it can only be AWS (1), GCP (4), or AZURE (8)")
	}

	insaneMode := d.Get("insane_mode").(bool)
	insaneModeAz := d.Get("insane_mode_az").(string)
	haSubnet := d.Get("ha_subnet").(string)
	haZone := d.Get("ha_zone").(string)
	if haZone != "" && gateway.CloudType != goaviatrix.GCP {
		return fmt.Errorf("'ha_zone' is only required for GCP provider if enabling HA")
	}
	if gateway.CloudType == goaviatrix.GCP && haSubnet != "" && haZone == "" {
		return fmt.Errorf("'ha_zone' must be set to enable HA on GCP, cannot enable HA with only 'ha_subnet' enabled")
	}
	haGwSize := d.Get("ha_gw_size").(string)
	if haSubnet == "" && haZone == "" && haGwSize != "" {
		return fmt.Errorf("'ha_gw_size' is only required if enabling HA")
	}
	haInsaneModeAz := d.Get("ha_insane_mode_az").(string)
	if insaneMode {
		if gateway.CloudType != goaviatrix.AWS && gateway.CloudType != goaviatrix.AZURE {
			return fmt.Errorf("insane_mode is only supported for aws and azure (cloud_type = 1 or 8)")
		}
		if gateway.CloudType == goaviatrix.AWS {
			if insaneModeAz == "" {
				return fmt.Errorf("insane_mode_az needed if insane_mode is enabled for aws cloud")
			}
			if haSubnet != "" && haInsaneModeAz == "" {
				return fmt.Errorf("ha_insane_mode_az needed if insane_mode is enabled for aws cloud and ha_subnet is set")
			}
			// Append availability zone to subnet
			var strs []string
			strs = append(strs, gateway.Subnet, insaneModeAz)
			gateway.Subnet = strings.Join(strs, "~~")
		}
		gateway.InsaneMode = "on"
	} else {
		gateway.InsaneMode = "off"
	}
	if haZone != "" || haSubnet != "" {
		if haGwSize == "" {
			return fmt.Errorf("A valid non empty ha_gw_size parameter is mandatory for this resource if " +
				"ha_subnet or ha_zone is set. Example: t2.micro")
		}
	}

	enableEncryptVolume := d.Get("enable_encrypt_volume").(bool)
	customerManagedKeys := d.Get("customer_managed_keys").(string)
	if enableEncryptVolume && d.Get("cloud_type").(int) != goaviatrix.AWS {
		return fmt.Errorf("'enable_encrypt_volume' is only supported for AWS provider (cloud_type: 1)")
	}
	if !enableEncryptVolume && customerManagedKeys != "" {
		return fmt.Errorf("'customer_managed_keys' should be empty since Encrypt Volume is not enabled")
	}
	if !enableEncryptVolume && gateway.CloudType == goaviatrix.AWS {
		gateway.EncVolume = "no"
	}

	log.Printf("[INFO] Creating Aviatrix Spoke Gateway: %#v", gateway)

	err := client.LaunchSpokeVpc(gateway)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Spoke Gateway: %s", err)
	}

	d.SetId(gateway.GwName)

	flag := false
	defer resourceAviatrixSpokeGatewayReadIfRequired(d, meta, &flag)

	if enableActiveMesh := d.Get("enable_active_mesh").(bool); !enableActiveMesh {
		gw := &goaviatrix.Gateway{
			GwName: d.Get("gw_name").(string),
		}
		gw.EnableActiveMesh = "no"

		err := client.DisableActiveMesh(gw)
		if err != nil {
			return fmt.Errorf("couldn't disable Active Mode for Aviatrix Spoke Gateway: %s", err)
		}
	}

	if !singleAZ {
		singleAZGateway := &goaviatrix.Gateway{
			GwName:   d.Get("gw_name").(string),
			SingleAZ: "disabled",
		}

		log.Printf("[INFO] Disable Single AZ GW HA: %#v", singleAZGateway)

		err := client.DisableSingleAZGateway(singleAZGateway)
		if err != nil {
			return fmt.Errorf("failed to disable single AZ GW HA: %s", err)
		}
	}

	if haSubnet != "" || haZone != "" {
		//Enable HA
		haGateway := &goaviatrix.SpokeVpc{
			GwName:    d.Get("gw_name").(string),
			CloudType: d.Get("cloud_type").(int),
			HASubnet:  haSubnet,
			HAZone:    haZone,
			Eip:       d.Get("ha_eip").(string),
		}

		if d.Get("cloud_type").(int) == goaviatrix.GCP {
			haGateway.HASubnetGCP = haSubnet
		}

		if insaneMode == true && haGateway.CloudType == goaviatrix.AWS {
			var haStrs []string
			haStrs = append(haStrs, haSubnet, haInsaneModeAz)
			haSubnet = strings.Join(haStrs, "~~")
			haGateway.HASubnet = haSubnet
		}

		if haGateway.CloudType == goaviatrix.GCP {
			err = client.EnableHaSpokeGateway(haGateway)
		} else {
			err = client.EnableHaSpokeVpc(haGateway)
		}
		if err != nil {
			return fmt.Errorf("failed to enable HA Aviatrix Spoke Gateway: %s", err)
		}

		log.Printf("[INFO]Resizing Spoke HA Gateway: %#v", haGwSize)

		if haGwSize != gateway.VpcSize {
			if haGwSize == "" {
				return fmt.Errorf("A valid non empty ha_gw_size parameter is mandatory for this resource if " +
					"ha_subnet or ha_zone is set. Example: t2.micro us-west1-b")
			}

			haGateway := &goaviatrix.Gateway{
				CloudType: d.Get("cloud_type").(int),
				GwName:    d.Get("gw_name").(string) + "-hagw",
			}

			haGateway.GwSize = d.Get("ha_gw_size").(string)

			log.Printf("[INFO] Resizing Spoke HA Gateway size to: %s ", haGateway.GwSize)

			err := client.UpdateGateway(haGateway)
			log.Printf("[INFO] Resizing Spoke HA Gateway size to: %s ", haGateway.GwSize)
			if err != nil {
				return fmt.Errorf("failed to update Aviatrix Spoke HA Gateway size: %s", err)
			}

			d.Set("ha_gw_size", haGwSize)
		}
	}

	if _, ok := d.GetOk("tag_list"); ok && gateway.CloudType == goaviatrix.AWS {
		tagList := d.Get("tag_list").([]interface{})
		tagListStr := goaviatrix.ExpandStringList(tagList)
		tagListStr = goaviatrix.TagListStrColon(tagListStr)
		gateway.TagList = strings.Join(tagListStr, ",")
		tags := &goaviatrix.Tags{
			CloudType:    goaviatrix.AWS,
			ResourceType: "gw",
			ResourceName: d.Get("gw_name").(string),
			TagList:      gateway.TagList,
		}
		err = client.AddTags(tags)
		if err != nil {
			return fmt.Errorf("failed to add tags: %s", err)
		}
	} else if ok && gateway.CloudType != goaviatrix.AWS {
		return fmt.Errorf("adding tags only supported for aws, cloud_type must be 1")
	}

	if transitGwName := d.Get("transit_gw").(string); transitGwName != "" {
		//No HA config, just return
		err := client.SpokeJoinTransit(gateway)
		if err != nil {
			return fmt.Errorf("failed to join Transit Gateway: %s", err)
		}
	}

	enableVpcDnsServer := d.Get("enable_vpc_dns_server").(bool)
	if d.Get("cloud_type").(int) == goaviatrix.AWS && enableVpcDnsServer {
		gwVpcDnsServer := &goaviatrix.Gateway{
			GwName: d.Get("gw_name").(string),
		}

		log.Printf("[INFO] Enable VPC DNS Server: %#v", gwVpcDnsServer)

		err := client.EnableVpcDnsServer(gwVpcDnsServer)
		if err != nil {
			return fmt.Errorf("failed to enable VPC DNS Server: %s", err)
		}
	} else if enableVpcDnsServer {
		return fmt.Errorf("'enable_vpc_dns_server' only supports AWS(1)")
	}

	if customizedSpokeVpcRoutes := d.Get("customized_spoke_vpc_routes").(string); customizedSpokeVpcRoutes != "" {
		transitGateway := &goaviatrix.Gateway{
			GwName:                   d.Get("gw_name").(string),
			CustomizedSpokeVpcRoutes: strings.Split(customizedSpokeVpcRoutes, ","),
		}
		for i := 0; ; i++ {
			log.Printf("[INFO] Editing customized routes of spoke gateway: %s ", transitGateway.GwName)
			err := client.EditGatewayCustomRoutes(transitGateway)
			if err == nil {
				break
			}
			if i <= 18 && (strings.Contains(err.Error(), "when it is down") || strings.Contains(err.Error(), "hagw is down")) {
				time.Sleep(10 * time.Second)
			} else {
				return fmt.Errorf("failed to customize spoke vpc routes of spoke gateway: %s due to: %s", transitGateway.GwName, err)
			}
		}
	}

	if filteredSpokeVpcRoutes := d.Get("filtered_spoke_vpc_routes").(string); filteredSpokeVpcRoutes != "" {
		transitGateway := &goaviatrix.Gateway{
			GwName:                 d.Get("gw_name").(string),
			FilteredSpokeVpcRoutes: strings.Split(filteredSpokeVpcRoutes, ","),
		}
		for i := 0; ; i++ {
			log.Printf("[INFO] Editing filtered routes of spoke gateway: %s ", transitGateway.GwName)
			err := client.EditGatewayFilterRoutes(transitGateway)
			if err == nil {
				break
			}
			if i <= 18 && (strings.Contains(err.Error(), "when it is down") || strings.Contains(err.Error(), "hagw is down")) {
				time.Sleep(10 * time.Second)
			} else {
				return fmt.Errorf("failed to edit filtered spoke vpc routes of spoke gateway: %s due to: %s", transitGateway.GwName, err)
			}
		}
	}

	if includedAdvertisedSpokeRoutes := d.Get("included_advertised_spoke_routes").(string); includedAdvertisedSpokeRoutes != "" {
		transitGateway := &goaviatrix.Gateway{
			GwName:                d.Get("gw_name").(string),
			AdvertisedSpokeRoutes: strings.Split(includedAdvertisedSpokeRoutes, ","),
		}
		for i := 0; ; i++ {
			log.Printf("[INFO] Editing customized routes advertisement of spoke gateway: %s ", transitGateway.GwName)
			err := client.EditGatewayAdvertisedCidr(transitGateway)
			if err == nil {
				break
			}
			if i <= 30 && (strings.Contains(err.Error(), "when it is down") || strings.Contains(err.Error(), "hagw is down")) {
				time.Sleep(10 * time.Second)
			} else {
				return fmt.Errorf("failed to edit advertised spoke vpc routes of spoke gateway: %s due to: %s", transitGateway.GwName, err)
			}
		}
	}

	return resourceAviatrixSpokeGatewayReadIfRequired(d, meta, &flag)
}

func resourceAviatrixSpokeGatewayReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixSpokeGatewayRead(d, meta)
	}
	return nil
}

func resourceAviatrixSpokeGatewayRead(d *schema.ResourceData, meta interface{}) error {
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

	if gw != nil {
		d.Set("cloud_type", gw.CloudType)
		d.Set("account_name", gw.AccountName)

		if gw.CloudType == goaviatrix.AWS {
			d.Set("vpc_id", strings.Split(gw.VpcID, "~~")[0]) //aws vpc_id returns as <vpc_id>~~<other vpc info> in rest api
			d.Set("vpc_reg", gw.VpcRegion)                    //aws vpc_reg returns as vpc_region in rest api

			if gw.AllocateNewEipRead {
				d.Set("allocate_new_eip", true)
			} else {
				d.Set("allocate_new_eip", false)
			}
		} else if gw.CloudType == goaviatrix.GCP {
			d.Set("vpc_id", strings.Split(gw.VpcID, "~-~")[0]) //gcp vpc_id returns as <vpc_id>~-~<other vpc info> in rest api
			d.Set("vpc_reg", gw.GatewayZone)                   //gcp vpc_reg returns as gateway_zone in json

			if gw.AllocateNewEipRead {
				d.Set("allocate_new_eip", true)
			} else {
				d.Set("allocate_new_eip", false)
			}
		} else if gw.CloudType == goaviatrix.AZURE || gw.CloudType == goaviatrix.OCI {
			d.Set("vpc_id", gw.VpcID)
			d.Set("vpc_reg", gw.VpcRegion)

			d.Set("allocate_new_eip", true)
		}
		d.Set("enable_encrypt_volume", gw.EnableEncryptVolume)
		d.Set("eip", gw.PublicIP)

		d.Set("subnet", gw.VpcNet)
		d.Set("gw_size", gw.GwSize)
		d.Set("cloud_instance_id", gw.CloudnGatewayInstID)
		d.Set("security_group_id", gw.GwSecurityGroupID)
		d.Set("private_ip", gw.PrivateIP)

		if gw.SingleAZ == "yes" {
			d.Set("single_az_ha", true)
		} else {
			d.Set("single_az_ha", false)
		}

		if gw.InsaneMode == "yes" {
			d.Set("insane_mode", true)
			if gw.CloudType == goaviatrix.AWS {
				d.Set("insane_mode_az", gw.GatewayZone)
			} else {
				d.Set("insane_mode_az", "")
			}
		} else {
			d.Set("insane_mode", false)
			d.Set("insane_mode_az", "")
		}

		if gw.EnableActiveMesh == "yes" {
			d.Set("enable_active_mesh", true)
		} else {
			d.Set("enable_active_mesh", false)
		}

		if gw.CloudType == goaviatrix.AWS && gw.EnableVpcDnsServer == "Enabled" {
			d.Set("enable_vpc_dns_server", true)
		} else {
			d.Set("enable_vpc_dns_server", false)
		}

		if gw.EnableNat == "yes" && gw.SnatMode == "primary" {
			d.Set("single_ip_snat", true)
		} else {
			d.Set("single_ip_snat", false)
		}

		if len(gw.CustomizedSpokeVpcRoutes) != 0 {
			if customizedSpokeVpcRoutes := d.Get("customized_spoke_vpc_routes").(string); customizedSpokeVpcRoutes != "" {
				customizedRoutesArray := strings.Split(customizedSpokeVpcRoutes, ",")
				if len(goaviatrix.Difference(customizedRoutesArray, gw.CustomizedSpokeVpcRoutes)) == 0 &&
					len(goaviatrix.Difference(gw.CustomizedSpokeVpcRoutes, customizedRoutesArray)) == 0 {
					d.Set("customized_spoke_vpc_routes", customizedSpokeVpcRoutes)
				} else {
					d.Set("customized_spoke_vpc_routes", strings.Join(gw.CustomizedSpokeVpcRoutes, ","))
				}
			} else {
				d.Set("customized_spoke_vpc_routes", strings.Join(gw.CustomizedSpokeVpcRoutes, ","))
			}
		} else {
			d.Set("customized_spoke_vpc_routes", "")
		}

		if len(gw.FilteredSpokeVpcRoutes) != 0 {
			if filteredSpokeVpcRoutes := d.Get("filtered_spoke_vpc_routes").(string); filteredSpokeVpcRoutes != "" {
				filteredSpokeVpcRoutesArray := strings.Split(filteredSpokeVpcRoutes, ",")
				if len(goaviatrix.Difference(filteredSpokeVpcRoutesArray, gw.FilteredSpokeVpcRoutes)) == 0 &&
					len(goaviatrix.Difference(gw.FilteredSpokeVpcRoutes, filteredSpokeVpcRoutesArray)) == 0 {
					d.Set("filtered_spoke_vpc_routes", filteredSpokeVpcRoutes)
				} else {
					d.Set("filtered_spoke_vpc_routes", strings.Join(gw.FilteredSpokeVpcRoutes, ","))
				}
			} else {
				d.Set("filtered_spoke_vpc_routes", strings.Join(gw.FilteredSpokeVpcRoutes, ","))
			}
		} else {
			d.Set("filtered_spoke_vpc_routes", "")
		}

		if len(gw.IncludeCidrList) != 0 {
			if includedAdvertisedSpokeRoutes := d.Get("included_advertised_spoke_routes").(string); includedAdvertisedSpokeRoutes != "" {
				advertisedSpokeRoutesArray := strings.Split(includedAdvertisedSpokeRoutes, ",")
				if len(goaviatrix.Difference(advertisedSpokeRoutesArray, gw.IncludeCidrList)) == 0 &&
					len(goaviatrix.Difference(gw.IncludeCidrList, advertisedSpokeRoutesArray)) == 0 {
					d.Set("included_advertised_spoke_routes", includedAdvertisedSpokeRoutes)
				} else {
					d.Set("included_advertised_spoke_routes", strings.Join(gw.IncludeCidrList, ","))
				}
			} else {
				d.Set("included_advertised_spoke_routes", strings.Join(gw.AdvertisedSpokeRoutes, ","))
			}
		} else {
			d.Set("included_advertised_spoke_routes", "")
		}
	}

	if gw.SpokeVpc == "yes" {
		d.Set("transit_gw", gw.TransitGwName)
	} else {
		d.Set("transit_gw", "")
	}

	if gw.CloudType == goaviatrix.AWS {
		tags := &goaviatrix.Tags{
			CloudType:    goaviatrix.AWS,
			ResourceType: "gw",
			ResourceName: d.Get("gw_name").(string),
		}
		tagList, err := client.GetTags(tags)
		if err != nil {
			return fmt.Errorf("unable to read tag_list for gateway: %v due to %v", gateway.GwName, err)
		}
		var tagListStr []string
		if _, ok := d.GetOk("tag_list"); ok {
			tagList1 := d.Get("tag_list").([]interface{})
			tagListStr = goaviatrix.ExpandStringList(tagList1)
		}
		if len(goaviatrix.Difference(tagListStr, tagList)) != 0 || len(goaviatrix.Difference(tagList, tagListStr)) != 0 {
			if err := d.Set("tag_list", tagList); err != nil {
				log.Printf("[WARN] Error setting tag_list for (%s): %s", d.Id(), err)
			}
		} else {
			if err := d.Set("tag_list", tagListStr); err != nil {
				log.Printf("[WARN] Error setting tag_list for (%s): %s", d.Id(), err)
			}
		}
	}

	haGateway := &goaviatrix.Gateway{
		AccountName: d.Get("account_name").(string),
		GwName:      d.Get("gw_name").(string) + "-hagw",
	}

	haGw, err := client.GetGateway(haGateway)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.Set("ha_gw_size", "")
			d.Set("ha_subnet", "")
			d.Set("ha_zone", "")
			d.Set("ha_eip", "")
			d.Set("ha_insane_mode_az", "")
		} else {
			return fmt.Errorf("couldn't find Aviatrix Spoke HA Gateway: %s", err)
		}
	} else {
		log.Printf("[INFO] Spoke HA Gateway size: %s", haGw.GwSize)
		if haGw.CloudType == goaviatrix.AWS || haGw.CloudType == goaviatrix.AZURE || haGw.CloudType == goaviatrix.OCI {
			d.Set("ha_subnet", haGw.VpcNet)
			d.Set("ha_zone", "")
		} else if haGw.CloudType == goaviatrix.GCP {
			d.Set("ha_zone", haGw.GatewayZone)
			if d.Get("ha_subnet") != "" || isImport {
				d.Set("ha_subnet", haGw.VpcNet)
			} else {
				d.Set("ha_subnet", "")
			}
		}

		d.Set("ha_eip", haGw.PublicIP)
		d.Set("ha_gw_size", haGw.GwSize)
		d.Set("ha_cloud_instance_id", haGw.CloudnGatewayInstID)
		d.Set("ha_gw_name", haGw.GwName)
		d.Set("ha_private_ip", haGw.PrivateIP)
		if haGw.InsaneMode == "yes" && haGw.CloudType == goaviatrix.AWS {
			d.Set("ha_insane_mode_az", haGw.GatewayZone)
		} else {
			d.Set("ha_insane_mode_az", "")
		}
	}

	return nil
}

func resourceAviatrixSpokeGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gateway := &goaviatrix.Gateway{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string),
	}

	haGateway := &goaviatrix.Gateway{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string) + "-hagw",
	}

	log.Printf("[INFO] Updating Aviatrix gateway: %#v", gateway)

	d.Partial(true)
	if d.HasChange("ha_zone") {
		haZone := d.Get("ha_zone").(string)
		if haZone != "" && gateway.CloudType != goaviatrix.GCP {
			return fmt.Errorf("'ha_zone' is only required for GCP provider if enabling HA")
		}
	}
	if d.HasChange("ha_zone") || d.HasChange("ha_subnet") {
		haZone := d.Get("ha_zone").(string)
		haSubnet := d.Get("ha_subnet").(string)
		if gateway.CloudType == goaviatrix.GCP && haSubnet != "" && haZone == "" {
			return fmt.Errorf("'ha_zone' must be set to enable HA on GCP, cannot enable HA with only 'ha_subnet' enabled")
		}
	}
	if d.HasChange("cloud_type") {
		return fmt.Errorf("updating cloud_type is not allowed")
	}
	if d.HasChange("account_name") {
		return fmt.Errorf("updating account_name is not allowed")
	}
	if d.HasChange("gw_name") {
		return fmt.Errorf("updating gw_name is not allowed")
	}
	if d.HasChange("vpc_id") {
		return fmt.Errorf("updating vpc_id is not allowed")
	}
	if d.HasChange("vpc_reg") {
		return fmt.Errorf("updating vpc_reg is not allowed")
	}
	if d.HasChange("subnet") {
		return fmt.Errorf("updating subnet is not allowed")
	}
	if d.HasChange("insane_mode") {
		return fmt.Errorf("updating insane_mode is not allowed")
	}
	if d.HasChange("insane_mode_az") {
		return fmt.Errorf("updating insane_mode_az is not allowed")
	}
	if d.HasChange("single_az_ha") {
		singleAZGateway := &goaviatrix.Gateway{
			GwName: d.Get("gw_name").(string),
		}

		singleAZ := d.Get("single_az_ha").(bool)

		if singleAZ {
			singleAZGateway.SingleAZ = "enabled"
		} else {
			singleAZGateway.SingleAZ = "disabled"
		}

		if singleAZGateway.SingleAZ == "enabled" {
			log.Printf("[INFO] Enable Single AZ GW HA: %#v", singleAZGateway)

			err := client.EnableSingleAZGateway(singleAZGateway)
			if err != nil {
				return fmt.Errorf("failed to enable single AZ GW HA: %s", err)
			}
		} else if singleAZGateway.SingleAZ == "disabled" {
			log.Printf("[INFO] Disable Single AZ GW HA: %#v", singleAZGateway)
			err := client.DisableSingleAZGateway(singleAZGateway)
			if err != nil {
				return fmt.Errorf("failed to disable single AZ GW HA: %s", err)
			}
		}

		d.SetPartial("single_az_ha")
	}

	if d.HasChange("tag_list") && gateway.CloudType == goaviatrix.AWS {
		tags := &goaviatrix.Tags{
			CloudType:    goaviatrix.AWS,
			ResourceType: "gw",
			ResourceName: d.Get("gw_name").(string),
		}
		o, n := d.GetChange("tag_list")
		if o == nil {
			o = new([]interface{})
		}
		if n == nil {
			n = new([]interface{})
		}
		os := o.([]interface{})
		ns := n.([]interface{})
		oldList := goaviatrix.ExpandStringList(os)
		newList := goaviatrix.ExpandStringList(ns)
		oldTagList := goaviatrix.Difference(oldList, newList)
		newTagList := goaviatrix.Difference(newList, oldList)
		if len(oldTagList) != 0 || len(newTagList) != 0 {
			if len(oldTagList) != 0 {
				oldTagList = goaviatrix.TagListStrColon(oldTagList)
				tags.TagList = strings.Join(oldTagList, ",")
				err := client.DeleteTags(tags)
				if err != nil {
					return fmt.Errorf("failed to delete tags : %s", err)
				}
			}
			if len(newTagList) != 0 {
				newTagList = goaviatrix.TagListStrColon(newTagList)
				tags.TagList = strings.Join(newTagList, ",")
				err := client.AddTags(tags)
				if err != nil {
					return fmt.Errorf("failed to add tags : %s", err)
				}
			}
		}

		d.SetPartial("tag_list")
	} else if d.HasChange("tag_list") && gateway.CloudType != goaviatrix.AWS {
		return fmt.Errorf("adding tags is only supported for aws, cloud_type must be set to 1")
	}

	//Get primary gw size if gw_size changed, to be used later on for ha gateway size update
	primaryGwSize := d.Get("gw_size").(string)
	if d.HasChange("gw_size") {
		old, _ := d.GetChange("gw_size")
		primaryGwSize = old.(string)
		gateway.GwSize = d.Get("gw_size").(string)
		err := client.UpdateGateway(gateway)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix Spoke Gateway: %s", err)
		}
		d.SetPartial("gw_size")
	}

	newHaGwEnabled := false
	if d.HasChange("ha_subnet") || d.HasChange("ha_zone") || d.HasChange("ha_insane_mode_az") {
		spokeGw := &goaviatrix.SpokeVpc{
			GwName:    d.Get("gw_name").(string),
			CloudType: d.Get("cloud_type").(int),
		}

		if spokeGw.CloudType == goaviatrix.AWS || spokeGw.CloudType == goaviatrix.GCP {
			spokeGw.Eip = d.Get("ha_eip").(string)
		}

		if !d.HasChange("ha_subnet") && d.HasChange("ha_insane_mode_az") {
			return fmt.Errorf("ha_subnet must change if ha_insane_mode_az changes")
		}
		if d.Get("insane_mode").(bool) == true && spokeGw.CloudType == goaviatrix.AWS {
			var haStrs []string
			insaneModeHaAz := d.Get("ha_insane_mode_az").(string)
			if insaneModeHaAz == "" {
				return fmt.Errorf("ha_insane_mode_az needed if insane_mode is enabled and ha_subnet is set")
			}
			haStrs = append(haStrs, spokeGw.HASubnet, insaneModeHaAz)
			spokeGw.HASubnet = strings.Join(haStrs, "~~")
		}

		oldSubnet, newSubnet := d.GetChange("ha_subnet")
		oldZone, newZone := d.GetChange("ha_zone")
		deleteHaGw := false
		changeHaGw := false
		if spokeGw.CloudType == goaviatrix.AWS || spokeGw.CloudType == goaviatrix.AZURE || spokeGw.CloudType == goaviatrix.OCI {
			spokeGw.HASubnet = d.Get("ha_subnet").(string)
			if oldSubnet == "" && newSubnet != "" {
				newHaGwEnabled = true
			} else if oldSubnet != "" && newSubnet == "" {
				deleteHaGw = true
			} else if oldSubnet != "" && newSubnet != "" {
				changeHaGw = true
			}
		} else if spokeGw.CloudType == goaviatrix.GCP {
			spokeGw.HAZone = d.Get("ha_zone").(string)
			spokeGw.HASubnetGCP = d.Get("ha_subnet").(string)
			if oldZone == "" && newZone != "" {
				newHaGwEnabled = true
			} else if oldZone != "" && newZone == "" {
				deleteHaGw = true
			} else if oldZone != "" && newZone != "" {
				changeHaGw = true
			}
		}
		if newHaGwEnabled {
			//New configuration to enable HA
			if haGateway.CloudType == goaviatrix.GCP {
				err := client.EnableHaSpokeGateway(spokeGw)
				if err != nil {
					return fmt.Errorf("failed to enable HA Aviatrix Spoke Gateway: %s", err)
				}
			} else {
				err := client.EnableHaSpokeVpc(spokeGw)
				if err != nil {
					return fmt.Errorf("failed to enable HA Aviatrix Spoke Gateway: %s", err)
				}
			}
			newHaGwEnabled = true
		} else if deleteHaGw {
			//Ha configuration has been deleted
			err := client.DeleteGateway(haGateway)
			if err != nil {
				return fmt.Errorf("failed to delete Aviatrix Spoke HA gateway: %s", err)
			}
		} else if changeHaGw {
			//HA subnet has been modified. Delete older HA GW,
			// and launch new HA GW in new subnet.
			err := client.DeleteGateway(haGateway)
			if err != nil {
				return fmt.Errorf("failed to delete Aviatrix Spoke HA gateway: %s", err)
			}

			//New configuration to enable HA
			if haGateway.CloudType == goaviatrix.GCP {
				err := client.EnableHaSpokeGateway(spokeGw)
				if err != nil {
					return fmt.Errorf("failed to enable HA Aviatrix Spoke Gateway: %s", err)
				}
			} else {
				err := client.EnableHaSpokeVpc(spokeGw)
				if err != nil {
					return fmt.Errorf("failed to enable HA Aviatrix Spoke Gateway: %s", err)
				}
			}
		}
		d.SetPartial("ha_subnet")
		d.SetPartial("ha_zone")
		d.SetPartial("ha_insane_mode_az")
	}

	if d.HasChange("ha_gw_size") || newHaGwEnabled {
		newHaGwSize := d.Get("ha_gw_size").(string)
		if !newHaGwEnabled || (newHaGwSize != primaryGwSize) {
			// MODIFIES HA GW SIZE if
			// Ha gateway wasn't newly configured
			// OR
			// newly configured Ha gateway is set to be different size than primary gateway
			// (when ha gateway is enabled, it's size is by default the same as primary gateway)
			_, err := client.GetGateway(haGateway)
			if err != nil {
				if err == goaviatrix.ErrNotFound {
					d.Set("ha_gw_size", "")
					d.Set("ha_subnet", "")
					d.Set("ha_zone", "")
					d.Set("ha_insane_mode_az", "")
					return nil
				}
				return fmt.Errorf("couldn't find Aviatrix Spoke HA Gateway while trying to update HA Gw size: %s", err)
			}
			haGateway.GwSize = d.Get("ha_gw_size").(string)
			if haGateway.GwSize == "" {
				return fmt.Errorf("A valid non empty ha_gw_size parameter is mandatory for this resource if " +
					"ha_subnet or ha_zone is set. Example: t2.micro or us-west1-b")
			}
			err = client.UpdateGateway(haGateway)
			log.Printf("[INFO] Updating HA Gateway size to: %s ", haGateway.GwSize)
			if err != nil {
				return fmt.Errorf("failed to update Aviatrix Spoke HA Gateway size: %s", err)
			}
		}
		d.SetPartial("ha_gw_size")
	}

	if d.HasChange("single_ip_snat") {
		enableSNat := d.Get("single_ip_snat").(bool)
		gw := &goaviatrix.Gateway{
			CloudType:   d.Get("cloud_type").(int),
			GatewayName: d.Get("gw_name").(string),
		}
		if enableSNat {
			err := client.EnableSNat(gw)
			if err != nil {
				return fmt.Errorf("failed to enable single_ip' mode SNAT: %s", err)
			}
		} else {
			err := client.DisableSNat(gw)
			if err != nil {
				return fmt.Errorf("failed to enable 'single_ip' mode SNAT: %s", err)
			}
		}
	}

	if d.HasChange("enable_active_mesh") && d.HasChange("transit_gw") {
		spokeVPC := &goaviatrix.SpokeVpc{
			CloudType:      d.Get("cloud_type").(int),
			GwName:         d.Get("gw_name").(string),
			HASubnet:       d.Get("ha_subnet").(string),
			TransitGateway: d.Get("transit_gw").(string),
		}

		o, n := d.GetChange("transit_gw")
		if o != "" {
			//New configuration to join to transit GW
			err := client.SpokeLeaveTransit(spokeVPC)
			if err != nil {
				return fmt.Errorf("failed to leave Transit Gateway: %s", err)
			}
		}

		gw := &goaviatrix.Gateway{
			GwName: d.Get("gw_name").(string),
		}
		enableActiveMesh := d.Get("enable_active_mesh").(bool)
		if enableActiveMesh {
			gw.EnableActiveMesh = "yes"
			err := client.EnableActiveMesh(gw)
			if err != nil {
				return fmt.Errorf("failed to enable Active Mesh Mode: %s", err)
			}
		} else {
			gw.EnableActiveMesh = "no"
			err := client.DisableActiveMesh(gw)
			if err != nil {
				return fmt.Errorf("failed to disable Active Mesh Mode: %s", err)
			}
		}

		if n != "" {
			err := client.SpokeJoinTransit(spokeVPC)
			if err != nil {
				return fmt.Errorf("failed to join Transit Gateway: %s", err)
			}
		}

		d.SetPartial("transit_gw")
		d.SetPartial("enable_active_mesh")
	} else if d.HasChange("enable_active_mesh") {
		gw := &goaviatrix.Gateway{
			GwName: d.Get("gw_name").(string),
		}

		enableActiveMesh := d.Get("enable_active_mesh").(bool)
		if enableActiveMesh {
			gw.EnableActiveMesh = "yes"
			err := client.EnableActiveMesh(gw)
			if err != nil {
				return fmt.Errorf("failed to enable Active Mesh Mode: %s", err)
			}
		} else {
			gw.EnableActiveMesh = "no"
			err := client.DisableActiveMesh(gw)
			if err != nil {
				return fmt.Errorf("failed to disable Active Mesh Mode: %s", err)
			}
		}

		d.SetPartial("enable_active_mesh")
	} else if d.HasChange("transit_gw") {
		spokeVPC := &goaviatrix.SpokeVpc{
			CloudType:      d.Get("cloud_type").(int),
			GwName:         d.Get("gw_name").(string),
			HASubnet:       d.Get("ha_subnet").(string),
			TransitGateway: d.Get("transit_gw").(string),
		}

		o, n := d.GetChange("transit_gw")
		if o == "" {
			//New configuration to join to transit GW
			err := client.SpokeJoinTransit(spokeVPC)
			if err != nil {
				return fmt.Errorf("failed to join Transit Gateway: %s", err)
			}
		} else if n == "" {
			//Transit GW has been deleted, leave transit GW.
			err := client.SpokeLeaveTransit(spokeVPC)
			if err != nil {
				return fmt.Errorf("failed to leave Transit Gateway: %s", err)
			}
		} else {
			//Change transit GW
			err := client.SpokeLeaveTransit(spokeVPC)
			if err != nil {
				return fmt.Errorf("failed to leave Transit Gateway: %s", err)
			}

			err = client.SpokeJoinTransit(spokeVPC)
			if err != nil {
				return fmt.Errorf("failed to join Transit Gateway: %s", err)
			}
		}

		d.SetPartial("transit_gw")
	}

	if d.HasChange("enable_vpc_dns_server") && d.Get("cloud_type").(int) == goaviatrix.AWS {
		gw := &goaviatrix.Gateway{
			CloudType: d.Get("cloud_type").(int),
			GwName:    d.Get("gw_name").(string),
		}

		enableVpcDnsServer := d.Get("enable_vpc_dns_server").(bool)
		if enableVpcDnsServer {
			err := client.EnableVpcDnsServer(gw)
			if err != nil {
				return fmt.Errorf("failed to enable VPC DNS Server: %s", err)
			}
		} else {
			err := client.DisableVpcDnsServer(gw)
			if err != nil {
				return fmt.Errorf("failed to disable VPC DNS Server: %s", err)
			}
		}

		d.SetPartial("enable_vpc_dns_server")
	} else if d.HasChange("enable_vpc_dns_server") {
		return fmt.Errorf("'enable_vpc_dns_server' only supports AWS(1)")
	}

	if d.HasChange("enable_encrypt_volume") {
		if d.Get("enable_encrypt_volume").(bool) {
			if d.Get("cloud_type").(int) != goaviatrix.AWS {
				return fmt.Errorf("'enable_encrypt_volume' is only supported for AWS provider (cloud_type: 1)")
			}
			gwEncVolume := &goaviatrix.Gateway{
				GwName:              d.Get("gw_name").(string),
				CustomerManagedKeys: d.Get("customer_managed_keys").(string),
			}
			err := client.EnableEncryptVolume(gwEncVolume)
			if err != nil {
				return fmt.Errorf("failed to enable encrypt gateway volume for %s due to %s", gwEncVolume.GwName, err)
			}
		} else {
			return fmt.Errorf("can't disable Encrypt Volume for gateway: %s", gateway.GwName)
		}
	} else if d.HasChange("customer_managed_keys") {
		return fmt.Errorf("updating customer_managed_keys only is not allowed")
	}

	if d.HasChange("customized_spoke_vpc_routes") {
		o, n := d.GetChange("customized_spoke_vpc_routes")
		if o == nil {
			o = new(interface{})
		}
		if n == nil {
			n = new(interface{})
		}
		os := o.(interface{})
		ns := n.(interface{})
		oldRouteList := strings.Split(os.(string), ",")
		newRouteList := strings.Split(ns.(string), ",")
		if len(goaviatrix.Difference(oldRouteList, newRouteList)) != 0 || len(goaviatrix.Difference(newRouteList, oldRouteList)) != 0 {
			transitGateway := &goaviatrix.Gateway{
				GwName:                   d.Get("gw_name").(string),
				CustomizedSpokeVpcRoutes: newRouteList,
			}
			err := client.EditGatewayCustomRoutes(transitGateway)
			log.Printf("[INFO] Customizeing routes of spoke gateway: %s ", transitGateway.GwName)
			if err != nil {
				return fmt.Errorf("failed to customize spoke vpc routes of spoke gateway: %s due to: %s", transitGateway.GwName, err)
			}
		}
		d.SetPartial("customized_spoke_vpc_routes")
	}

	if d.HasChange("filtered_spoke_vpc_routes") {
		o, n := d.GetChange("filtered_spoke_vpc_routes")
		if o == nil {
			o = new(interface{})
		}
		if n == nil {
			n = new(interface{})
		}
		os := o.(interface{})
		ns := n.(interface{})
		oldRouteList := strings.Split(os.(string), ",")
		newRouteList := strings.Split(ns.(string), ",")
		if len(goaviatrix.Difference(oldRouteList, newRouteList)) != 0 || len(goaviatrix.Difference(newRouteList, oldRouteList)) != 0 {
			transitGateway := &goaviatrix.Gateway{
				GwName:                 d.Get("gw_name").(string),
				FilteredSpokeVpcRoutes: newRouteList,
			}
			err := client.EditGatewayFilterRoutes(transitGateway)
			log.Printf("[INFO] Editing filtered spoke vpc routes of spoke gateway: %s ", transitGateway.GwName)
			if err != nil {
				return fmt.Errorf("failed to edit filtered spoke vpc routes of spoke gateway: %s due to: %s", transitGateway.GwName, err)
			}
		}
		d.SetPartial("filtered_spoke_vpc_routes")
	}

	if d.HasChange("included_advertised_spoke_routes") {
		o, n := d.GetChange("included_advertised_spoke_routes")
		if o == nil {
			o = new(interface{})
		}
		if n == nil {
			n = new(interface{})
		}
		os := o.(interface{})
		ns := n.(interface{})
		oldRouteList := strings.Split(os.(string), ",")
		newRouteList := strings.Split(ns.(string), ",")
		if len(goaviatrix.Difference(oldRouteList, newRouteList)) != 0 || len(goaviatrix.Difference(newRouteList, oldRouteList)) != 0 {
			transitGateway := &goaviatrix.Gateway{
				GwName:                d.Get("gw_name").(string),
				AdvertisedSpokeRoutes: newRouteList,
			}
			err := client.EditGatewayAdvertisedCidr(transitGateway)
			log.Printf("[INFO] Editing included advertised spoke vpc routes of spoke gateway: %s ", transitGateway.GwName)
			if err != nil {
				return fmt.Errorf("failed to edit included advertised spoke vpc routes of spoke gateway: %s due to: %s", transitGateway.GwName, err)
			}
		}
		d.SetPartial("included_advertised_spoke_routes")
	}

	d.Partial(false)
	d.SetId(gateway.GwName)
	return resourceAviatrixSpokeGatewayRead(d, meta)
}

func resourceAviatrixSpokeGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gateway := &goaviatrix.Gateway{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix Spoke Gateway: %#v", gateway)

	if transitGw := d.Get("transit_gw").(string); transitGw != "" {
		spokeVPC := &goaviatrix.SpokeVpc{
			GwName: d.Get("gw_name").(string),
		}

		err := client.SpokeLeaveTransit(spokeVPC)
		if err != nil {
			return fmt.Errorf("failed to leave Transit Gateway: %s", err)
		}
	}

	//If HA is enabled, delete HA GW first.
	haSubnet := d.Get("ha_subnet").(string)
	haZone := d.Get("ha_zone").(string)
	if haSubnet != "" || haZone != "" {
		//Delete HA Gw too
		gateway.GwName += "-hagw"
		err := client.DeleteGateway(gateway)
		if err != nil {
			return fmt.Errorf("failed to delete Aviatrix Spoke HA gateway: %s", err)
		}
	}

	gateway.GwName = d.Get("gw_name").(string)

	err := client.DeleteGateway(gateway)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Spoke Gateway: %s", err)
	}

	return nil
}
