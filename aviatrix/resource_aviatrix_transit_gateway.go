package aviatrix

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixTransitGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixTransitGatewayCreate,
		Read:   resourceAviatrixTransitGatewayRead,
		Update: resourceAviatrixTransitGatewayUpdate,
		Delete: resourceAviatrixTransitGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,
		MigrateState:  resourceAviatrixTransitGatewayMigrateState,

		Schema: map[string]*schema.Schema{
			"cloud_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Type of cloud service provider, requires an integer value. Use 1 for AWS.",
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
				Description: "Public Subnet Name.",
			},
			"insane_mode_az": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "AZ of subnet being created for Insane Mode Transit Gateway. Required if insane_mode is enabled.",
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
				Default:     "",
				Description: "HA Subnet. Required for enabling HA for AWS/ARM gateway.",
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
				Description: "AZ of subnet being created for Insane Mode Transit HA Gateway. Required if insane_mode is enabled and ha_subnet is set.",
			},
			"ha_gw_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "HA Gateway Size. Mandatory if HA is enabled (ha_subnet is set).",
			},
			"ha_eip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Public IP address that you want assigned to the HA Transit Gateway.",
			},
			"single_az_ha": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Set to 'enabled' if this feature is desired.",
			},
			"single_ip_snat": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable or disable Source NAT feature in 'single_ip' mode for this container.",
			},
			"tag_list": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Default:     nil,
				Description: "Instance tag of cloud provider.",
			},
			"enable_hybrid_connection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Sign of readiness for TGW connection.",
			},
			"connected_transit": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Specify Connected Transit status.",
			},
			"insane_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable Insane Mode for Transit. Valid values: true, false. If insane mode is enabled, gateway size has to at least be c5 size.",
			},
			"enable_firenet": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Specify whether to enable firenet interfaces or not.",
			},
			"enable_active_mesh": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Switch to Enable/Disable Active Mesh Mode for Transit Gateway. Valid values: true, false.",
			},
			"enable_vpc_dns_server": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable vpc_dns_server for Gateway. Only supports AWS. Valid values: true, false.",
			},
			"enable_advertise_transit_cidr": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Switch to Enable/Disable advertise transit VPC network CIDR.",
			},
			"bgp_manual_spoke_advertise_cidrs": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Intended CIDR list to advertise to VGW.",
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
					"It applies to all spoke gateways attached to this transit gateway.",
			},
			"filtered_spoke_vpc_routes": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				Description: "A list of comma separated CIDRs to be filtered from the spoke VPC route table. When configured, " +
					"filtering CIDR(s) or it’s subnet will be deleted from VPC routing tables as well as from spoke gateway’s " +
					"routing table. It applies to all spoke gateways attached to this transit gateway.",
			},
			"excluded_advertised_spoke_routes": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				Description: "A list of comma separated CIDRs to be advertised to on-prem as 'Excluded CIDR List'. " +
					"When configured, it inspects all the advertised CIDRs from its spoke gateways and " +
					"remove those included in the 'Excluded CIDR List'.",
			},
			"customer_managed_keys": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Customer managed key ID.",
			},
			"enable_transit_firenet": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Specify whether to enable transit firenet interfaces or not.",
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Security group used for the transit gateway.",
			},
			"cloud_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance ID of the transit gateway.",
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private IP address of the transit gateway created.",
			},
			"ha_cloud_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud instance ID of ha transit gateway.",
			},
		},
	}
}

func resourceAviatrixTransitGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gateway := &goaviatrix.TransitVpc{
		CloudType:              d.Get("cloud_type").(int),
		AccountName:            d.Get("account_name").(string),
		GwName:                 d.Get("gw_name").(string),
		VpcID:                  d.Get("vpc_id").(string),
		VpcSize:                d.Get("gw_size").(string),
		Subnet:                 d.Get("subnet").(string),
		EnableHybridConnection: d.Get("enable_hybrid_connection").(bool),
	}

	enableNAT := d.Get("single_ip_snat").(bool)
	if enableNAT {
		gateway.EnableNAT = "yes"
	} else {
		gateway.EnableNAT = "no"
	}

	singleAZ := d.Get("single_az_ha").(bool)
	if singleAZ {
		gateway.SingleAzHa = "enabled"
	} else {
		gateway.SingleAzHa = "disabled"
	}

	connectedTransit := d.Get("connected_transit").(bool)
	if connectedTransit {
		gateway.ConnectedTransit = "yes"
	} else {
		gateway.ConnectedTransit = "no"
	}

	allocateNewEip := d.Get("allocate_new_eip").(bool)
	if allocateNewEip {
		gateway.ReuseEip = "off"
	} else {
		gateway.ReuseEip = "on"
		gateway.Eip = d.Get("eip").(string)
	}

	cloudType := d.Get("cloud_type").(int)
	if cloudType == 1 || cloudType == 4 || cloudType == 16 {
		gateway.VpcID = d.Get("vpc_id").(string)
		if gateway.VpcID == "" {
			return fmt.Errorf("'vpc_id' cannot be empty for creating a transit gw for aws vpc")
		}
	} else if cloudType == 8 {
		gateway.VNetNameResourceGroup = d.Get("vpc_id").(string)
		if gateway.VNetNameResourceGroup == "" {
			return fmt.Errorf("'vpc_id' cannot be empty for creating a transit gw for azure vnet")
		}
	}

	if gateway.CloudType == 1 || gateway.CloudType == 8 || gateway.CloudType == 16 {
		gateway.VpcRegion = d.Get("vpc_reg").(string)
	} else if gateway.CloudType == 4 {
		// for gcp, rest api asks for "zone" rather than vpc region
		gateway.Zone = d.Get("vpc_reg").(string)
	} else {
		return fmt.Errorf("invalid cloud type, it can only be AWS (1), GCP (4), or ARM (8)")
	}

	insaneMode := d.Get("insane_mode").(bool)
	if insaneMode {
		if cloudType != 1 && cloudType != 8 {
			return fmt.Errorf("insane_mode is only supported for aws and arm (cloud_type = 1 or 8)")
		}
		if cloudType == 1 {
			if d.Get("insane_mode_az").(string) == "" {
				return fmt.Errorf("insane_mode_az needed if insane_mode is enabled for aws cloud")
			}
			if d.Get("ha_subnet").(string) != "" && d.Get("ha_insane_mode_az").(string) == "" {
				return fmt.Errorf("ha_insane_mode_az needed if insane_mode is enabled for aws cloud and ha_subnet is set")
			}
			// Append availability zone to subnet
			var strs []string
			insaneModeAz := d.Get("insane_mode_az").(string)
			strs = append(strs, gateway.Subnet, insaneModeAz)
			gateway.Subnet = strings.Join(strs, "~~")
		}
		gateway.InsaneMode = "on"
	} else {
		gateway.InsaneMode = "off"
	}

	haSubnet := d.Get("ha_subnet").(string)
	if haSubnet != "" && gateway.CloudType != 1 && gateway.CloudType != 8 && gateway.CloudType != 16 && gateway.CloudType != 256 {
		return fmt.Errorf("'ha_subnet' is only required for AWS/ARM/OCI providers if enabling HA")
	}
	haZone := d.Get("ha_zone").(string)
	if haZone != "" && gateway.CloudType != 4 {
		return fmt.Errorf("'ha_zone' is only required for GCP provider if enabling HA")
	}
	haGwSize := d.Get("ha_gw_size").(string)
	if haSubnet == "" && haZone == "" && haGwSize != "" {
		return fmt.Errorf("'ha_gw_size' is only required if enabling HA")
	}
	if haGwSize == "" && haSubnet != "" {
		return fmt.Errorf("A valid non empty ha_gw_size parameter is mandatory for this resource if " +
			"ha_subnet is set. Example: t2.micro")
	}

	enableEncryptVolume := d.Get("enable_encrypt_volume").(bool)
	customerManagedKeys := d.Get("customer_managed_keys").(string)
	if enableEncryptVolume && d.Get("cloud_type").(int) != 1 {
		return fmt.Errorf("'enable_encrypt_volume' is only supported for AWS provider (cloud_type: 1)")
	}
	if !enableEncryptVolume && customerManagedKeys != "" {
		return fmt.Errorf("'customer_managed_keys' should be empty since Encrypt Volume is not enabled")
	}
	if enableEncryptVolume && d.Get("single_az_ha").(bool) {
		return fmt.Errorf("'single_az_ha' needs to be disabled to encrypt gateway EBS volume")
	}

	enableFireNet := d.Get("enable_firenet").(bool)
	enableTransitFireNet := d.Get("enable_transit_firenet").(bool)
	if enableFireNet && enableTransitFireNet {
		return fmt.Errorf("can't enable firenet function and transit firenet function at the same time")
	}
	if enableTransitFireNet && gateway.CloudType != 1 && gateway.CloudType != 8 {
		return fmt.Errorf("'enable_transit_firenet' is only supported in AWS and AZURE providers")
	} else if enableTransitFireNet && gateway.CloudType == 8 {
		gateway.EnableTransitFireNet = "on"
	}

	log.Printf("[INFO] Creating Aviatrix Transit Gateway: %#v", gateway)

	err := client.LaunchTransitVpc(gateway)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Transit Gateway: %s", err)
	}

	d.SetId(gateway.GwName)

	flag := false
	defer resourceAviatrixTransitGatewayReadIfRequired(d, meta, &flag)

	if enableActiveMesh := d.Get("enable_active_mesh").(bool); !enableActiveMesh {
		gw := &goaviatrix.Gateway{
			GwName: d.Get("gw_name").(string),
		}
		gw.EnableActiveMesh = "no"

		err := client.DisableActiveMesh(gw)
		if err != nil {
			return fmt.Errorf("couldn't disable Active Mode for Aviatrix Transit Gateway: %s", err)
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
		transitGateway := &goaviatrix.TransitVpc{
			CloudType: d.Get("cloud_type").(int),
			GwName:    d.Get("gw_name").(string),
			HASubnet:  haSubnet,
			Eip:       d.Get("ha_eip").(string),
		}

		if insaneMode {
			var haStrs []string
			insaneModeHaAz := d.Get("ha_insane_mode_az").(string)
			haStrs = append(haStrs, haSubnet, insaneModeHaAz)
			haSubnet = strings.Join(haStrs, "~~")
			transitGateway.HASubnet = haSubnet
		}

		if transitGateway.CloudType == 4 && haZone == "" {
			return fmt.Errorf("no ha_zone is provided for enabling Transit HA gateway: %s", transitGateway.GwName)
		} else if transitGateway.CloudType == 4 {
			transitGateway.HAZone = haZone
		}

		log.Printf("[INFO] Enabling HA on Transit Gateway: %#v", haSubnet)

		err = client.EnableHaTransitVpc(transitGateway)
		if err != nil {
			return fmt.Errorf("failed to enable HA Aviatrix Transit Gateway: %s", err)
		}

		//Resize HA Gateway
		log.Printf("[INFO]Resizing Transit HA Gateway: %#v", haGwSize)

		if haGwSize != gateway.VpcSize {
			if haGwSize == "" {
				return fmt.Errorf("A valid non empty ha_gw_size parameter is mandatory for this resource if " +
					"ha_subnet is set. Example: t2.micro")
			}

			haGateway := &goaviatrix.Gateway{
				CloudType: d.Get("cloud_type").(int),
				GwName:    d.Get("gw_name").(string) + "-hagw",
			}
			haGateway.GwSize = d.Get("ha_gw_size").(string)

			log.Printf("[INFO] Resizing Transit HA GAteway size to: %s ", haGateway.GwSize)

			err := client.UpdateGateway(haGateway)
			if err != nil {
				return fmt.Errorf("failed to update Aviatrix Transit HA Gateway size: %s", err)
			}
		}
	}

	if _, ok := d.GetOk("tag_list"); ok {
		if cloudType != 1 {
			return fmt.Errorf("'tag_list' is only supported for AWS cloud type 1")
		}
		tagList := d.Get("tag_list").([]interface{})
		tagListStr := goaviatrix.ExpandStringList(tagList)
		tagListStr = goaviatrix.TagListStrColon(tagListStr)
		gateway.TagList = strings.Join(tagListStr, ",")
		tags := &goaviatrix.Tags{
			CloudType:    1,
			ResourceType: "gw",
			ResourceName: d.Get("gw_name").(string),
			TagList:      gateway.TagList,
		}

		err = client.AddTags(tags)
		if err != nil {
			return fmt.Errorf("failed to add tags: %s", err)
		}
	}

	enableHybridConnection := d.Get("enable_hybrid_connection").(bool)
	if enableHybridConnection {
		if cloudType != 1 {
			return fmt.Errorf("'enable_hybrid_connection' is only supported for AWS cloud type 1")
		}

		err := client.AttachTransitGWForHybrid(gateway)
		if err != nil {
			return fmt.Errorf("failed to enable transit GW for Hybrid: %s", err)
		}
	}

	if connectedTransit {
		err := client.EnableConnectedTransit(gateway)
		if err != nil {
			return fmt.Errorf("failed to enable connected transit: %s", err)
		}
	}

	if enableNAT {
		gw := &goaviatrix.Gateway{
			GatewayName: gateway.GwName,
		}

		err := client.EnableSNat(gw)
		if err != nil {
			return fmt.Errorf("failed to enable SNAT: %s", err)
		}
	}

	if enableFireNet {
		err := client.EnableGatewayFireNetInterfaces(gateway)
		if err != nil {
			return fmt.Errorf("failed to enable transit GW for FireNet Interfaces: %s", err)
		}
	}

	enableVpcDnsServer := d.Get("enable_vpc_dns_server").(bool)
	if d.Get("cloud_type").(int) == 1 && enableVpcDnsServer {
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

	enableAdvertiseTransitCidr := d.Get("enable_advertise_transit_cidr").(bool)
	if enableAdvertiseTransitCidr {
		err := client.EnableAdvertiseTransitCidr(gateway)
		if err != nil {
			return fmt.Errorf("failed to enable advertise transit CIDR: %s", err)
		}
	}

	bgpManualSpokeAdvertiseCidrs := d.Get("bgp_manual_spoke_advertise_cidrs").(string)
	if bgpManualSpokeAdvertiseCidrs != "" {
		gateway.BgpManualSpokeAdvertiseCidrs = bgpManualSpokeAdvertiseCidrs
		err := client.SetBgpManualSpokeAdvertisedNetworks(gateway)
		if err != nil {
			return fmt.Errorf("failed to set BGP Manual Spoke Advertise Cidrs: %s", err)
		}
	}

	if enableEncryptVolume {
		gwEncVolume := &goaviatrix.Gateway{
			GwName:              d.Get("gw_name").(string),
			CustomerManagedKeys: customerManagedKeys,
		}
		err := client.EnableEncryptVolume(gwEncVolume)
		if err != nil {
			return fmt.Errorf("failed to enable encrypt gateway volume for %s due to %s", gwEncVolume.GwName, err)
		}
	}

	if customizedSpokeVpcRoutes := d.Get("customized_spoke_vpc_routes").(string); customizedSpokeVpcRoutes != "" {
		transitGateway := &goaviatrix.Gateway{
			GwName:                   d.Get("gw_name").(string),
			CustomizedSpokeVpcRoutes: strings.Split(customizedSpokeVpcRoutes, ","),
		}
		for i := 0; ; i++ {
			log.Printf("[INFO] Editing customized routes of transit gateway: %s ", transitGateway.GwName)
			err := client.EditGatewayCustomRoutes(transitGateway)
			if err == nil {
				break
			}
			if i <= 10 && strings.Contains(err.Error(), "when it is down") {
				time.Sleep(10 * time.Second)
			} else {
				return fmt.Errorf("failed to customize spoke vpc routes of transit gateway: %s due to: %s", transitGateway.GwName, err)
			}
		}
	}

	if filteredSpokeVpcRoutes := d.Get("filtered_spoke_vpc_routes").(string); filteredSpokeVpcRoutes != "" {
		transitGateway := &goaviatrix.Gateway{
			GwName:                 d.Get("gw_name").(string),
			FilteredSpokeVpcRoutes: strings.Split(filteredSpokeVpcRoutes, ","),
		}
		for i := 0; ; i++ {
			log.Printf("[INFO] Editing filtered routes of transit gateway: %s ", transitGateway.GwName)
			err := client.EditGatewayFilterRoutes(transitGateway)
			if err == nil {
				break
			}
			if i <= 10 && strings.Contains(err.Error(), "when it is down") {
				time.Sleep(10 * time.Second)
			} else {
				return fmt.Errorf("failed to edit filtered spoke vpc routes of transit gateway: %s due to: %s", transitGateway.GwName, err)
			}
		}
	}

	if advertisedSpokeRoutesExclude := d.Get("excluded_advertised_spoke_routes").(string); advertisedSpokeRoutesExclude != "" {
		transitGateway := &goaviatrix.Gateway{
			GwName:                d.Get("gw_name").(string),
			AdvertisedSpokeRoutes: strings.Split(advertisedSpokeRoutesExclude, ","),
		}
		for i := 0; ; i++ {
			log.Printf("[INFO] Editing customized routes advertisement of transit gateway: %s ", transitGateway.GwName)
			err := client.EditGatewayAdvertisedCidr(transitGateway)
			if err == nil {
				break
			}
			if i <= 10 && strings.Contains(err.Error(), "when it is down") {
				time.Sleep(10 * time.Second)
			} else {
				return fmt.Errorf("failed to edit advertised spoke vpc routes of transit gateway: %s due to: %s", transitGateway.GwName, err)
			}
		}
	}

	if enableTransitFireNet && gateway.CloudType == 1 {
		enableActiveMesh := d.Get("enable_active_mesh").(bool)
		if !connectedTransit || !enableActiveMesh {
			return fmt.Errorf("both active_mesh and connected_transit need to be enabled to enable transit firenet")
		}
		gwTransitFireNet := &goaviatrix.Gateway{
			GwName: d.Get("gw_name").(string),
		}
		err := client.EnableTransitFireNet(gwTransitFireNet)
		if err != nil {
			return fmt.Errorf("failed to enable transit firenet for %s due to %s", gwTransitFireNet.GwName, err)
		}
	}

	return resourceAviatrixTransitGatewayReadIfRequired(d, meta, &flag)
}

func resourceAviatrixTransitGatewayReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixTransitGatewayRead(d, meta)
	}
	return nil
}

func resourceAviatrixTransitGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gwName := d.Get("gw_name").(string)
	if gwName == "" {
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
		return fmt.Errorf("couldn't find Aviatrix Transit Gateway: %s", err)
	}

	log.Printf("[TRACE] reading gateway %s: %#v", d.Get("gw_name").(string), gw)

	if gw != nil {
		d.Set("cloud_type", gw.CloudType)
		d.Set("account_name", gw.AccountName)
		d.Set("gw_name", gw.GwName)
		d.Set("subnet", gw.VpcNet)

		if gw.CloudType == 1 {
			d.Set("vpc_id", strings.Split(gw.VpcID, "~~")[0])
			d.Set("vpc_reg", gw.VpcRegion)
			if gw.AllocateNewEipRead {
				d.Set("allocate_new_eip", true)
			} else {
				d.Set("allocate_new_eip", false)
			}
		} else if gw.CloudType == 4 {
			d.Set("vpc_id", strings.Split(gw.VpcID, "~-~")[0])
			d.Set("vpc_reg", gw.GatewayZone)
			d.Set("allocate_new_eip", true)
		} else if gw.CloudType == 8 || gw.CloudType == 16 {
			d.Set("vpc_id", gw.VpcID)
			d.Set("vpc_reg", gw.VpcRegion)
			d.Set("allocate_new_eip", true)
		}

		d.Set("enable_encrypt_volume", gw.EnableEncryptVolume)
		d.Set("eip", gw.PublicIP)
		d.Set("gw_size", gw.GwSize)
		d.Set("cloud_instance_id", gw.CloudnGatewayInstID)
		d.Set("security_group_id", gw.GwSecurityGroupID)
		d.Set("private_ip", gw.PrivateIP)

		if gw.EnableNat == "yes" && gw.SnatMode == "primary" {
			d.Set("single_ip_snat", true)
		} else {
			d.Set("single_ip_snat", false)
		}

		if gw.SingleAZ == "yes" {
			d.Set("single_az_ha", true)
		} else {
			d.Set("single_az_ha", false)
		}

		if gw.CloudType == 1 {
			d.Set("enable_hybrid_connection", gw.EnableHybridConnection)
		} else {
			d.Set("enable_hybrid_connection", false)
		}

		if gw.ConnectedTransit == "yes" {
			d.Set("connected_transit", true)
		} else {
			d.Set("connected_transit", false)
		}

		if gw.InsaneMode == "yes" {
			d.Set("insane_mode", true)
			if gw.CloudType == 1 {
				d.Set("insane_mode_az", gw.GatewayZone)
			} else {
				d.Set("insane_mode_az", "")
			}
		} else {
			d.Set("insane_mode", false)
			d.Set("insane_mode_az", "")
		}

		if len(gw.CustomizedSpokeVpcRoutes) != 0 {
			if customizedRoutes := d.Get("customized_spoke_vpc_routes").(string); customizedRoutes != "" {
				customizedRoutesArray := strings.Split(customizedRoutes, ",")
				if len(goaviatrix.Difference(customizedRoutesArray, gw.CustomizedSpokeVpcRoutes)) == 0 &&
					len(goaviatrix.Difference(gw.CustomizedSpokeVpcRoutes, customizedRoutesArray)) == 0 {
					d.Set("customized_spoke_vpc_routes", customizedRoutes)
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

		if len(gw.ExcludeCidrList) != 0 {
			if advertisedSpokeRoutes := d.Get("excluded_advertised_spoke_routes").(string); advertisedSpokeRoutes != "" {
				advertisedSpokeRoutesArray := strings.Split(advertisedSpokeRoutes, ",")
				if len(goaviatrix.Difference(advertisedSpokeRoutesArray, gw.ExcludeCidrList)) == 0 &&
					len(goaviatrix.Difference(gw.ExcludeCidrList, advertisedSpokeRoutesArray)) == 0 {
					d.Set("excluded_advertised_spoke_routes", advertisedSpokeRoutes)
				} else {
					d.Set("excluded_advertised_spoke_routes", strings.Join(gw.ExcludeCidrList, ","))
				}
			} else {
				d.Set("excluded_advertised_spoke_routes", strings.Join(gw.ExcludeCidrList, ","))
			}
		} else {
			d.Set("excluded_advertised_spoke_routes", "")
		}

		gwDetail, err := client.GetGatewayDetail(gw)
		if err != nil {
			return fmt.Errorf("couldn't get Aviatrix Transit Gateway: %s", err)
		}

		d.Set("enable_firenet", gwDetail.EnableFireNet)
		d.Set("enable_transit_firenet", gwDetail.EnableTransitFireNet)

		if gw.EnableActiveMesh == "yes" {
			d.Set("enable_active_mesh", true)
		} else {
			d.Set("enable_active_mesh", false)
		}

		if gw.CloudType == 1 && gw.EnableVpcDnsServer == "Enabled" {
			d.Set("enable_vpc_dns_server", true)
		} else {
			d.Set("enable_vpc_dns_server", false)
		}

		if gwDetail.EnableAdvertiseTransitCidr == "yes" {
			d.Set("enable_advertise_transit_cidr", true)
		} else {
			d.Set("enable_advertise_transit_cidr", false)
		}

		var bgpManualSpokeAdvertiseCidrs []string
		if _, ok := d.GetOk("bgp_manual_spoke_advertise_cidrs"); ok {
			bgpManualSpokeAdvertiseCidrs = strings.Split(d.Get("bgp_manual_spoke_advertise_cidrs").(string), ",")
		}
		if len(goaviatrix.Difference(bgpManualSpokeAdvertiseCidrs, gwDetail.BgpManualSpokeAdvertiseCidrs)) != 0 ||
			len(goaviatrix.Difference(gwDetail.BgpManualSpokeAdvertiseCidrs, bgpManualSpokeAdvertiseCidrs)) != 0 {
			bgpMSAN := ""
			for i := range gwDetail.BgpManualSpokeAdvertiseCidrs {
				if i == 0 {
					bgpMSAN = bgpMSAN + gwDetail.BgpManualSpokeAdvertiseCidrs[i]
				} else {
					bgpMSAN = bgpMSAN + "," + gwDetail.BgpManualSpokeAdvertiseCidrs[i]
				}
			}
			d.Set("bgp_manual_spoke_advertise_cidrs", bgpMSAN)
		} else {
			d.Set("bgp_manual_spoke_advertise_cidrs", d.Get("bgp_manual_spoke_advertise_cidrs").(string))
		}
	}

	if gw.CloudType == 1 {
		tags := &goaviatrix.Tags{
			CloudType:    1,
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
			d.Set("ha_insane_mode_az", "")
			d.Set("ha_eip", "")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix Transit HA Gateway: %s", err)
	} else {
		if haGw.CloudType == 1 || haGw.CloudType == 8 || haGw.CloudType == 16 {
			d.Set("ha_subnet", haGw.VpcNet)
			d.Set("ha_zone", "")
		} else if haGw.CloudType == 4 {
			d.Set("ha_zone", haGw.GatewayZone)
			d.Set("ha_subnet", "")
		}
		d.Set("ha_eip", haGw.PublicIP)
		d.Set("ha_gw_size", haGw.GwSize)
		d.Set("ha_cloud_instance_id", haGw.CloudnGatewayInstID)
	}

	if haGw.InsaneMode == "yes" && haGw.CloudType == 1 {
		d.Set("ha_insane_mode_az", haGw.GatewayZone)
	} else {
		d.Set("ha_insane_mode_az", "")
	}

	return nil
}

func resourceAviatrixTransitGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	gateway := &goaviatrix.Gateway{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string),
	}
	haGateway := &goaviatrix.Gateway{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string) + "-hagw",
	}
	log.Printf("[INFO] Updating Aviatrix Transit Gateway: %#v", gateway)

	d.Partial(true)
	if d.HasChange("ha_subnet") {
		haSubnet := d.Get("ha_subnet").(string)
		if haSubnet != "" && gateway.CloudType != 1 && gateway.CloudType != 8 && gateway.CloudType != 16 && gateway.CloudType != 256 {
			return fmt.Errorf("'ha_subnet' is only required for AWS/ARM/OCI providers if enabling HA")
		}
	}
	if d.HasChange("ha_zone") {
		haZone := d.Get("ha_zone").(string)
		if haZone != "" && gateway.CloudType != 4 {
			return fmt.Errorf("'ha_zone' is only required for GCP provider if enabling HA")
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
	if d.HasChange("enable_transit_firenet") && d.Get("cloud_type").(int) == 8 {
		return fmt.Errorf("editing 'enable_transit_firenet' in AZURE is not supported")
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

	if d.HasChange("gw_size") {
		gateway.GwSize = d.Get("gw_size").(string)

		err := client.UpdateGateway(gateway)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix Transit Gateway: %s", err)
		}

		d.SetPartial("gw_size")
	}

	newHaGwEnabled := false
	if d.HasChange("ha_subnet") || d.HasChange("ha_zone") || d.HasChange("ha_insane_mode_az") {
		transitGw := &goaviatrix.TransitVpc{
			GwName:    d.Get("gw_name").(string),
			CloudType: d.Get("cloud_type").(int),
		}

		if transitGw.CloudType == 1 {
			transitGw.Eip = d.Get("ha_eip").(string)
		}

		if !d.HasChange("ha_subnet") && d.HasChange("ha_insane_mode_az") {
			return fmt.Errorf("ha_subnet must change if ha_insane_mode_az changes")
		}
		if d.Get("insane_mode").(bool) && transitGw.CloudType == 1 {
			var haStrs []string
			insaneModeHaAz := d.Get("ha_insane_mode_az").(string)
			if insaneModeHaAz == "" {
				return fmt.Errorf("ha_insane_mode_az needed if insane_mode is enabled and ha_subnet is set")
			}
			haStrs = append(haStrs, transitGw.HASubnet, insaneModeHaAz)
			transitGw.HASubnet = strings.Join(haStrs, "~~")
		}

		oldSubnet, newSubnet := d.GetChange("ha_subnet")
		oldZone, newZone := d.GetChange("ha_zone")
		deleteHaGw := false
		changeHaGw := false
		if transitGw.CloudType == 1 || transitGw.CloudType == 8 {
			transitGw.HASubnet = d.Get("ha_subnet").(string)
			if oldSubnet == "" && newSubnet != "" {
				newHaGwEnabled = true
			} else if oldSubnet != "" && newSubnet == "" {
				deleteHaGw = true
			} else if oldSubnet != "" && newSubnet != "" {
				changeHaGw = true
			}
		} else if transitGw.CloudType == 4 {
			transitGw.HAZone = d.Get("ha_zone").(string)
			if oldZone == "" && newZone != "" {
				newHaGwEnabled = true
			} else if oldZone != "" && newZone == "" {
				deleteHaGw = true
			} else if oldZone != "" && newZone != "" {
				changeHaGw = true
			}
		}
		if newHaGwEnabled {
			err := client.EnableHaTransitVpc(transitGw)
			if err != nil {
				return fmt.Errorf("failed to enable HA Aviatrix Transit Gateway: %s", err)
			}
			newHaGwEnabled = true
		} else if deleteHaGw {
			err := client.DeleteGateway(haGateway)
			if err != nil {
				return fmt.Errorf("failed to delete Aviatrix Transit HA gateway: %s", err)
			}
		} else if changeHaGw {
			err := client.DeleteGateway(haGateway)
			if err != nil {
				return fmt.Errorf("failed to delete Aviatrix Transit HA gateway: %s", err)
			}

			haErr := client.EnableHaTransitVpc(transitGw)
			if haErr != nil {
				return fmt.Errorf("failed to enable HA Aviatrix Transit Gateway: %s", err)
			}
		}
		d.SetPartial("ha_subnet")
		d.SetPartial("ha_zone")
		d.SetPartial("ha_insane_mode_az")
	}

	if gateway.CloudType == 1 {
		if d.HasChange("tag_list") {
			tags := &goaviatrix.Tags{
				CloudType:    1,
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
		}
	} else {
		if d.HasChange("tag_list") {
			return fmt.Errorf("'tag_list' is only supported for AWS cloud type 1")
		}
	}

	if d.HasChange("connected_transit") {
		transitGateway := &goaviatrix.TransitVpc{
			CloudType:   d.Get("cloud_type").(int),
			AccountName: d.Get("account_name").(string),
			GwName:      d.Get("gw_name").(string),
			VpcID:       d.Get("vpc_id").(string),
			VpcRegion:   d.Get("vpc_reg").(string),
		}
		connectedTransit := d.Get("connected_transit").(bool)
		if connectedTransit {
			err := client.EnableConnectedTransit(transitGateway)
			if err != nil {
				return fmt.Errorf("failed to enable connected transit: %s", err)
			}
		} else {
			err := client.DisableConnectedTransit(transitGateway)
			if err != nil {
				return fmt.Errorf("failed to disable connected transit: %s", err)
			}
		}

		d.SetPartial("connected_transit")
	}

	if d.HasChange("ha_gw_size") {
		_, err := client.GetGateway(haGateway)
		if err != nil {
			if err == goaviatrix.ErrNotFound {
				d.Set("ha_gw_size", "")
				d.Set("ha_subnet", "")
				d.Set("ha_insane_mode_az", "")
				return nil
			}
			return fmt.Errorf("couldn't find Aviatrix Transit HA Gateway while trying to update HA Gw "+
				"size: %s", err)
		}

		haGateway.GwSize = d.Get("ha_gw_size").(string)
		if haGateway.GwSize == "" {
			return fmt.Errorf("A valid non empty ha_gw_size parameter is mandatory for this resource if " +
				"ha_subnet is set. Example: t2.micro")
		}

		err = client.UpdateGateway(haGateway)
		log.Printf("[INFO] Updating Transit HA GAteway size to: %s ", haGateway.GwSize)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix Transit HA Gw size: %s", err)
		}

		d.SetPartial("ha_gw_size")
	}

	if d.HasChange("single_ip_snat") {
		gw := &goaviatrix.Gateway{
			CloudType:   d.Get("cloud_type").(int),
			GatewayName: d.Get("gw_name").(string),
		}
		enableNat := d.Get("single_ip_snat").(bool)

		if enableNat {
			err := client.EnableSNat(gw)
			if err != nil {
				return fmt.Errorf("failed to enable 'single_ip' mode SNAT feature: %s", err)
			}
		} else {
			err := client.DisableSNat(gw)
			if err != nil {
				return fmt.Errorf("failed to disable 'single_ip' mode SNAT: %s", err)
			}
		}

		d.SetPartial("single_ip_snat")
	}

	if gateway.CloudType == 1 {
		if d.HasChange("enable_hybrid_connection") {
			transitGateway := &goaviatrix.TransitVpc{
				CloudType:   d.Get("cloud_type").(int),
				AccountName: d.Get("account_name").(string),
				GwName:      d.Get("gw_name").(string),
				VpcID:       d.Get("vpc_id").(string),
				VpcRegion:   d.Get("vpc_reg").(string),
			}
			enableHybridConnection := d.Get("enable_hybrid_connection").(bool)
			if enableHybridConnection {
				err := client.AttachTransitGWForHybrid(transitGateway)
				if err != nil {
					return fmt.Errorf("failed to enable transit GW for Hybrid: %s", err)
				}
			} else {
				err := client.DetachTransitGWForHybrid(transitGateway)
				if err != nil {
					return fmt.Errorf("failed to disable transit GW for Hybrid: %s", err)
				}
			}
		}
	} else {
		if d.HasChange("enable_hybrid_connection") {
			return fmt.Errorf("'enable_hybrid_connection' is only supported for AWS cloud type 1")
		}
	}

	if d.HasChange("enable_active_mesh") {
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
	}

	if d.HasChange("enable_firenet") && d.HasChange("enable_transit_firenet") {
		transitGW := &goaviatrix.TransitVpc{
			GwName: gateway.GwName,
			VpcID:  d.Get("vpc_id").(string),
		}
		enableFireNet := d.Get("enable_firenet").(bool)
		enableTransitFireNet := d.Get("enable_transit_firenet").(bool)
		if enableFireNet && enableTransitFireNet {
			return fmt.Errorf("can't enable firenet function and transit firenet function at the same time")
		}
		if !enableFireNet {
			err := client.DisableGatewayFireNetInterfaces(transitGW)
			if err != nil {
				return fmt.Errorf("failed to disable transit GW for FireNet Interfaces: %s", err)
			}
		}
		if !enableTransitFireNet {
			gwTransitFireNet := &goaviatrix.Gateway{
				GwName: d.Get("gw_name").(string),
			}
			err := client.DisableTransitFireNet(gwTransitFireNet)
			if err != nil {
				return fmt.Errorf("failed to disable transit firenet for %s due to %s", gwTransitFireNet.GwName, err)
			}
		}
		if enableFireNet {
			err := client.EnableGatewayFireNetInterfaces(transitGW)
			if err != nil {
				return fmt.Errorf("failed to enable transit GW for FireNet Interfaces: %s", err)
			}
		}
		if enableTransitFireNet {
			enableActiveMesh := d.Get("enable_active_mesh").(bool)
			connectedTransit := d.Get("connected_transit").(bool)
			if !connectedTransit || !enableActiveMesh {
				return fmt.Errorf("both active_mesh and connected_transit need to be enabled to enable transit firenet")
			}
			gwTransitFireNet := &goaviatrix.Gateway{
				GwName: d.Get("gw_name").(string),
			}
			err := client.EnableTransitFireNet(gwTransitFireNet)
			if err != nil {
				return fmt.Errorf("failed to enable transit firenet for %s due to %s", gwTransitFireNet.GwName, err)
			}
		}

		d.SetPartial("enable_firenet")
		d.SetPartial("enable_transit_firenet")
	} else if d.HasChange("enable_firenet") {
		transitGW := &goaviatrix.TransitVpc{
			GwName: gateway.GwName,
			VpcID:  d.Get("vpc_id").(string),
		}
		enableFireNet := d.Get("enable_firenet").(bool)
		if enableFireNet {
			err := client.EnableGatewayFireNetInterfaces(transitGW)
			if err != nil {
				return fmt.Errorf("failed to enable transit GW for FireNet Interfaces: %s", err)
			}
		} else {
			err := client.DisableGatewayFireNetInterfaces(transitGW)
			if err != nil {
				return fmt.Errorf("failed to disable transit GW for FireNet Interfaces: %s", err)
			}
		}

		d.SetPartial("enable_firenet")
	} else if d.HasChange("enable_transit_firenet") {
		enableTransitFireNet := d.Get("enable_transit_firenet").(bool)
		if enableTransitFireNet {
			enableActiveMesh := d.Get("enable_active_mesh").(bool)
			connectedTransit := d.Get("connected_transit").(bool)
			if !connectedTransit || !enableActiveMesh {
				return fmt.Errorf("both active_mesh and connected_transit need to be enabled to enable transit firenet")
			}
			gwTransitFireNet := &goaviatrix.Gateway{
				GwName: d.Get("gw_name").(string),
			}
			err := client.EnableTransitFireNet(gwTransitFireNet)
			if err != nil {
				return fmt.Errorf("failed to enable transit firenet for %s due to %s", gwTransitFireNet.GwName, err)
			}
		} else {
			gwTransitFireNet := &goaviatrix.Gateway{
				GwName: d.Get("gw_name").(string),
			}
			err := client.DisableTransitFireNet(gwTransitFireNet)
			if err != nil {
				return fmt.Errorf("failed to disable transit firenet for %s due to %s", gwTransitFireNet.GwName, err)
			}
		}
		d.SetPartial("enable_transit_firenet")
	}

	if d.HasChange("enable_vpc_dns_server") && d.Get("cloud_type").(int) == 1 {
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

	if d.HasChange("enable_advertise_transit_cidr") {
		transitGw := &goaviatrix.TransitVpc{
			GwName: d.Get("gw_name").(string),
		}
		enableAdvertiseTransitCidr := d.Get("enable_advertise_transit_cidr").(bool)
		if enableAdvertiseTransitCidr {
			transitGw.EnableAdvertiseTransitCidr = true
			err := client.EnableAdvertiseTransitCidr(transitGw)
			if err != nil {
				return fmt.Errorf("failed to enable advertise transit CIDR: %s", err)
			}
		} else {
			transitGw.EnableAdvertiseTransitCidr = false
			err := client.DisableAdvertiseTransitCidr(transitGw)
			if err != nil {
				return fmt.Errorf("failed to disable advertise transit CIDR: %s", err)
			}
		}
		d.SetPartial("enable_advertise_transit_cidr")
	}

	if d.HasChange("bgp_manual_spoke_advertise_cidrs") {
		transitGw := &goaviatrix.TransitVpc{
			GwName: d.Get("gw_name").(string),
		}
		bgpManualSpokeAdvertiseCidrs := d.Get("bgp_manual_spoke_advertise_cidrs").(string)
		transitGw.BgpManualSpokeAdvertiseCidrs = bgpManualSpokeAdvertiseCidrs
		err := client.SetBgpManualSpokeAdvertisedNetworks(transitGw)
		if err != nil {
			return fmt.Errorf("failed to set bgp manual spoke advertise CIDRs: %s", err)
		}

		d.SetPartial("bgp_manual_spoke_advertise_cidrs")
	}

	if d.HasChange("enable_encrypt_volume") {
		if d.Get("enable_encrypt_volume").(bool) {
			if d.Get("cloud_type").(int) != 1 {
				return fmt.Errorf("'enable_encrypt_volume' is only supported for AWS provider (cloud_type: 1)")
			}
			if d.Get("single_az_ha").(bool) {
				return fmt.Errorf("'single_az_ha' needs to be disabled to encrypt gateway EBS volume")
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
		d.SetPartial("enable_encrypt_volume")
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
			log.Printf("[INFO] Customizeing routes of transit gateway: %s ", transitGateway.GwName)
			if err != nil {
				return fmt.Errorf("failed to customize spoke vpc routes of transit gateway: %s due to: %s", transitGateway.GwName, err)
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
			log.Printf("[INFO] Editing filtered spoke vpc routes of transit gateway: %s ", transitGateway.GwName)
			if err != nil {
				return fmt.Errorf("failed to edit filtered spoke vpc routes of transit gateway: %s due to: %s", transitGateway.GwName, err)
			}
		}
		d.SetPartial("filtered_spoke_vpc_routes")
	}

	if d.HasChange("excluded_advertised_spoke_routes") {
		o, n := d.GetChange("excluded_advertised_spoke_routes")
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
			log.Printf("[INFO] Editing excluded advertised spoke vpc routes of transit gateway: %s ", transitGateway.GwName)
			if err != nil {
				return fmt.Errorf("failed to edit excluded advertised spoke vpc routes of transit gateway: %s due to: %s", transitGateway.GwName, err)
			}
		}
		d.SetPartial("excluded_advertised_spoke_routes")
	}

	d.Partial(false)
	return resourceAviatrixTransitGatewayRead(d, meta)
}

func resourceAviatrixTransitGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gateway := &goaviatrix.Gateway{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix Transit Gateway: %#v", gateway)

	enableFireNet := d.Get("enable_firenet").(bool)
	if enableFireNet {
		gw := &goaviatrix.TransitVpc{
			CloudType: d.Get("cloud_type").(int),
			GwName:    d.Get("gw_name").(string),
		}

		err := client.DisableGatewayFireNetInterfaces(gw)
		if err != nil {
			return fmt.Errorf("failed to disable Aviatrix Transit Gateway for FireNet Interfaces: %s", err)
		}
	}

	enableTransitFireNet := d.Get("enable_transit_firenet").(bool)
	if enableTransitFireNet && gateway.CloudType == 1 {
		err := client.DisableTransitFireNet(gateway)
		if err != nil {
			return fmt.Errorf("failed to disable transit firenet for %s due to %s", gateway.GwName, err)
		}
	} else if enableTransitFireNet && gateway.CloudType == 8 {
		err := client.IsTransitFireNetReadyToBeDisabled(gateway)
		if err != nil {
			return fmt.Errorf("failed to disable transit firenet for %s due to %s", gateway.GwName, err)
		}
	}

	//If HA is enabled, delete HA GW first.
	haSubnet := d.Get("ha_subnet").(string)
	haZone := d.Get("ha_zone").(string)
	if haSubnet != "" || haZone != "" {
		gateway.GwName += "-hagw"

		err := client.DeleteGateway(gateway)
		if err != nil {
			return fmt.Errorf("failed to delete Aviatrix Transit Gateway HA gateway: %s", err)
		}
	}

	gateway.GwName = d.Get("gw_name").(string)

	err := client.DeleteGateway(gateway)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Transit Gateway: %s", err)
	}

	return nil
}
