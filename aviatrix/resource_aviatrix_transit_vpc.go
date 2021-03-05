package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixTransitVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixTransitVpcCreate,
		Read:   resourceAviatrixTransitVpcRead,
		Update: resourceAviatrixTransitVpcUpdate,
		Delete: resourceAviatrixTransitVpcDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 2,
		MigrateState:  resourceTransitVpcMigrateState,

		Schema: map[string]*schema.Schema{
			"cloud_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateCloudType,
				Description:  "Type of cloud service provider, requires an integer value. Use 1 for AWS.",
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
			"vpc_size": {
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
			"ha_subnet": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "HA Subnet.",
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
			"enable_nat": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "no",
				Description: "Enable NAT for this container.",
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
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "no",
				Description: "Specify Connected Transit status.",
			},
			"insane_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable Insane Mode for Transit. Valid values: true, false. If insane mode is enabled, gateway size has to at least be c5 size.",
			},
			"enable_firenet_interfaces": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Specify whether to enable firenet interfaces or not.",
			},
		},
	}
}

func resourceAviatrixTransitVpcCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	gateway := &goaviatrix.TransitVpc{
		CloudType:              d.Get("cloud_type").(int),
		AccountName:            d.Get("account_name").(string),
		GwName:                 d.Get("gw_name").(string),
		VpcID:                  d.Get("vpc_id").(string),
		VpcRegion:              d.Get("vpc_reg").(string),
		VpcSize:                d.Get("vpc_size").(string),
		Subnet:                 d.Get("subnet").(string),
		EnableNAT:              d.Get("enable_nat").(string),
		EnableHybridConnection: d.Get("enable_hybrid_connection").(bool),
		ConnectedTransit:       d.Get("connected_transit").(string),
	}

	cloudType := d.Get("cloud_type").(int)
	if cloudType == goaviatrix.AWS {
		gateway.VpcID = d.Get("vpc_id").(string)
		if gateway.VpcID == "" {
			return fmt.Errorf("'vpc_id' cannot be empty for creating a transit gw for aws vpc")
		}
	} else if cloudType == goaviatrix.AZURE {
		gateway.VNetNameResourceGroup = d.Get("vpc_id").(string)
		if gateway.VNetNameResourceGroup == "" {
			return fmt.Errorf("'vpc_id' cannot be empty for creating a transit gw for azure vnet")
		}
	}

	if gateway.EnableNAT != "" && gateway.EnableNAT != "yes" && gateway.EnableNAT != "no" {
		return fmt.Errorf("enable_nat can only be empty string, 'yes', or 'no'")
	}
	enableNat := gateway.EnableNAT

	insaneMode := d.Get("insane_mode").(bool)
	if insaneMode {
		if cloudType != goaviatrix.AWS {
			return fmt.Errorf("insane_mode is only support for aws (cloud_type = 1)")
		}
		if d.Get("insane_mode_az").(string) == "" {
			return fmt.Errorf("insane_mode_az needed if insane_mode is enabled")
		}
		if d.Get("ha_subnet").(string) != "" && d.Get("ha_insane_mode_az").(string) == "" {
			return fmt.Errorf("ha_insane_mode_az needed if insane_mode is enabled and ha_subnet is set")
		}

		gateway.InsaneMode = "on"

		// Append availability zone to subnet
		var strs []string
		insaneModeAz := d.Get("insane_mode_az").(string)
		strs = append(strs, gateway.Subnet, insaneModeAz)
		gateway.Subnet = strings.Join(strs, "~~")
	} else {
		gateway.InsaneMode = "off"
	}

	haSubnet := d.Get("ha_subnet").(string)
	haGwSize := d.Get("ha_gw_size").(string)
	if haGwSize == "" && haSubnet != "" {
		return fmt.Errorf("A valid non empty ha_gw_size parameter is mandatory for this resource if " +
			"ha_subnet is set. Example: t2.micro")
	}

	log.Printf("[INFO] Creating Aviatrix TransitVpc: %#v", gateway)

	err := client.LaunchTransitVpc(gateway)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix TransitVpc: %s", err)
	}

	d.SetId(gateway.GwName)

	flag := false
	defer resourceAviatrixTransitVpcReadIfRequired(d, meta, &flag)

	if haSubnet != "" {
		//Enable HA
		transitGateway := &goaviatrix.TransitVpc{
			GwName:   d.Get("gw_name").(string),
			HASubnet: haSubnet,
		}
		if insaneMode {
			var haStrs []string
			insaneModeHaAz := d.Get("ha_insane_mode_az").(string)
			haStrs = append(haStrs, haSubnet, insaneModeHaAz)
			haSubnet = strings.Join(haStrs, "~~")
			transitGateway.HASubnet = haSubnet
		}

		log.Printf("[INFO] Enabling HA on Transit Gateway: %#v", haSubnet)

		err = client.EnableHaTransitVpc(transitGateway)
		if err != nil {
			return fmt.Errorf("failed to enable2 HA Aviatrix TransitVpc: %s", err)
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
			err := client.UpdateGateway(haGateway)
			log.Printf("[INFO] Resizing Transit HA GAteway size to: %s ", haGateway.GwSize)
			if err != nil {
				return fmt.Errorf("failed to update Aviatrix Transit HA Gateway size: %s", err)
			}
		}
	}

	if _, ok := d.GetOk("tag_list"); ok {
		if cloudType != goaviatrix.AWS {
			return fmt.Errorf("'tag_list' is only supported for AWS cloud type 1")
		}
		tagList := d.Get("tag_list").([]interface{})
		tagListStr := goaviatrix.ExpandStringList(tagList)
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
	}

	enableHybridConnection := d.Get("enable_hybrid_connection").(bool)
	if enableHybridConnection && cloudType != goaviatrix.AWS {
		return fmt.Errorf("'enable_hybrid_connection' is only supported for AWS cloud type 1")
	}
	if enableHybridConnection {
		if cloudType != goaviatrix.AWS {
			return fmt.Errorf("'enable_hybrid_connection' is only supported for AWS cloud type 1")
		}
		err := client.AttachTransitGWForHybrid(gateway)
		if err != nil {
			return fmt.Errorf("failed to enable transit GW for Hybrid: %s", err)
		}
	}

	connectedTransit := d.Get("connected_transit").(string)
	if connectedTransit == "yes" {
		err := client.EnableConnectedTransit(gateway)
		if err != nil {
			return fmt.Errorf("failed to enable connected transit: %s", err)
		}
	}

	if enableNat == "yes" {
		gw := &goaviatrix.Gateway{
			GwName: gateway.GwName,
		}
		err := client.EnableSNat(gw)
		if err != nil {
			return fmt.Errorf("failed to enable SNAT: %s", err)
		}
	}

	enableFireNetInterfaces := d.Get("enable_firenet_interfaces").(bool)
	if enableFireNetInterfaces {
		err := client.EnableGatewayFireNetInterfaces(gateway)
		if err != nil {
			return fmt.Errorf("failed to enable transit GW for FireNet Interfaces: %s", err)
		}
	}

	return resourceAviatrixTransitVpcReadIfRequired(d, meta, &flag)
}

func resourceAviatrixTransitVpcReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixTransitVpcRead(d, meta)
	}
	return nil
}

func resourceAviatrixTransitVpcRead(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("couldn't find Aviatrix TransitVpc: %s", err)
	}

	log.Printf("[TRACE] reading gateway %s: %#v", d.Get("gw_name").(string), gw)

	if gw != nil {
		d.Set("cloud_type", gw.CloudType)
		d.Set("account_name", gw.AccountName)
		d.Set("gw_name", gw.GwName)
		d.Set("subnet", gw.VpcNet)
		if gw.CloudType == goaviatrix.AWS {
			d.Set("vpc_id", strings.Split(gw.VpcID, "~~")[0])
		} else if gw.CloudType == goaviatrix.AZURE {
			d.Set("vpc_id", gw.VpcID)
		}
		d.Set("vpc_reg", gw.VpcRegion)
		d.Set("vpc_size", gw.GwSize)
		d.Set("enable_nat", gw.EnableNat)
		if gw.CloudType == goaviatrix.AWS {
			d.Set("enable_hybrid_connection", gw.EnableHybridConnection)
		} else {
			d.Set("enable_hybrid_connection", false)
		}
		d.Set("connected_transit", gw.ConnectedTransit)
		if gw.InsaneMode == "yes" {
			d.Set("insane_mode", true)
			d.Set("insane_mode_az", gw.GatewayZone)
		} else {
			d.Set("insane_mode", false)
			d.Set("insane_mode_az", "")
		}

		gwDetail, err := client.GetGatewayDetail(gw)
		if err != nil {
			return fmt.Errorf("couldn't get Aviatrix TransitVpc: %s", err)
		}
		d.Set("enable_firenet_interfaces", gwDetail.DMZEnabled)
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
			d.Set("ha_insane_mode_az", "")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix Transit HA Gateway: %s", err)
	}
	d.Set("ha_subnet", haGw.VpcNet)
	d.Set("ha_gw_size", haGw.GwSize)
	if haGw.InsaneMode == "yes" {
		d.Set("ha_insane_mode_az", haGw.GatewayZone)
	} else {
		d.Set("ha_insane_mode_az", "")
	}

	return nil
}

func resourceAviatrixTransitVpcUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	gateway := &goaviatrix.Gateway{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string),
	}
	haGateway := &goaviatrix.Gateway{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string) + "-hagw",
	}
	log.Printf("[INFO] Updating Aviatrix TransitVpc: %#v", gateway)

	d.Partial(true)
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

	if d.HasChange("vpc_size") {
		gateway.GwSize = d.Get("vpc_size").(string)
		err := client.UpdateGateway(gateway)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix TransitVpc: %s", err)
		}
	}
	if d.HasChange("ha_subnet") || d.HasChange("ha_insane_mode_az") {
		transitGateway := &goaviatrix.TransitVpc{
			GwName:   d.Get("gw_name").(string),
			HASubnet: d.Get("ha_subnet").(string),
		}
		if d.Get("insane_mode").(bool) {
			var haStrs []string
			insaneModeHaAz := d.Get("ha_insane_mode_az").(string)

			if insaneModeHaAz == "" {
				return fmt.Errorf("ha_insane_mode_az needed if insane_mode is enabled and ha_subnet is set")
			}
			haStrs = append(haStrs, transitGateway.HASubnet, insaneModeHaAz)
			transitGateway.HASubnet = strings.Join(haStrs, "~~")
		}

		o, n := d.GetChange("ha_subnet")
		if o == "" {
			//New configuration to enable HA
			err := client.EnableHaTransitVpc(transitGateway)
			if err != nil {
				return fmt.Errorf("failed to enable HA Aviatrix TransitVpc: %s", err)
			}
		} else if n == "" {
			//Ha configuration has been deleted
			err := client.DeleteGateway(haGateway)
			if err != nil {
				return fmt.Errorf("failed to delete Aviatrix TransitVpc HA gateway: %s", err)
			}
		} else {
			//HA subnet has been modified. Delete older HA GW, and launch new HA GW in new subnet.
			err := client.DeleteGateway(haGateway)
			if err != nil {
				return fmt.Errorf("failed to delete Aviatrix TransitVpc HA gateway: %s", err)
			}

			gateway.GwName = d.Get("gw_name").(string)
			//New configuration to enable HA
			haErr := client.EnableHaTransitVpc(transitGateway)
			if haErr != nil {
				return fmt.Errorf("failed to enable HA Aviatrix TransitVpc: %s", err)
			}
		}
	}

	if gateway.CloudType == goaviatrix.AWS {
		if d.HasChange("tag_list") {
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
					tags.TagList = strings.Join(oldTagList, ",")
					err := client.DeleteTags(tags)
					if err != nil {
						return fmt.Errorf("failed to delete tags : %s", err)
					}
				}
				if len(newTagList) != 0 {
					tags.TagList = strings.Join(newTagList, ",")
					err := client.AddTags(tags)
					if err != nil {
						return fmt.Errorf("failed to add tags : %s", err)
					}
				}
			}
		}
	} else {
		if d.HasChange("tag_list") {
			return fmt.Errorf("'tag_list' is only supported for AWS cloud type 1")
		}
	}

	if gateway.CloudType == goaviatrix.AWS {
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

	if d.HasChange("connected_transit") {
		transitGateway := &goaviatrix.TransitVpc{
			CloudType:   d.Get("cloud_type").(int),
			AccountName: d.Get("account_name").(string),
			GwName:      d.Get("gw_name").(string),
			VpcID:       d.Get("vpc_id").(string),
			VpcRegion:   d.Get("vpc_reg").(string),
		}
		connectedTransit := d.Get("connected_transit").(string)
		if connectedTransit != "yes" && connectedTransit != "no" {
			return fmt.Errorf("connected_transit is not set correctly")
		}
		if connectedTransit == "yes" {
			err := client.EnableConnectedTransit(transitGateway)
			if err != nil {
				return fmt.Errorf("failed to enable connected transit: %s", err)
			}
		}
		if connectedTransit == "no" {
			err := client.DisableConnectedTransit(transitGateway)
			if err != nil {
				return fmt.Errorf("failed to disable connected transit: %s", err)
			}
		}
	}

	if d.HasChange("ha_gw_size") {
		_, err := client.GetGateway(haGateway)
		if err != nil {
			if err == goaviatrix.ErrNotFound {
				d.Set("ha_gw_size", "")
				d.Set("ha_subnet", "")
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
	}

	if d.HasChange("enable_nat") {
		gw := &goaviatrix.Gateway{
			CloudType: d.Get("cloud_type").(int),
			GwName:    d.Get("gw_name").(string),
		}
		o, n := d.GetChange("enable_nat")
		if o == "yes" && n == "no" {
			err := client.DisableSNat(gw)
			if err != nil {
				return fmt.Errorf("failed to disable SNAT: %s", err)
			}
		}
		if o == "no" && n == "yes" {
			err := client.EnableSNat(gw)
			if err != nil {
				return fmt.Errorf("failed to enable SNAT: %s", err)
			}
		}
	}

	if d.HasChange("enable_firenet_interfaces") {
		transitGW := &goaviatrix.TransitVpc{
			GwName: gateway.GwName,
		}
		enableFireNetInterfaces := d.Get("enable_firenet_interfaces").(bool)
		if enableFireNetInterfaces {
			err := client.EnableGatewayFireNetInterfaces(transitGW)
			if err != nil {
				return fmt.Errorf("failed to enable transit GW for FireNet Interfaces: %s", err)
			}
		} else {
			err := client.DisableGatewayFireNetInterfaces(transitGW)
			if err != nil {
				return fmt.Errorf("failed to remove transit GW for FireNet Interfaces: %s", err)
			}
		}
	}

	d.Partial(false)
	return nil
}

func resourceAviatrixTransitVpcDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	gateway := &goaviatrix.Gateway{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix TransitVpc: %#v", gateway)

	enableFireNetInterfaces := d.Get("enable_firenet_interfaces").(bool)
	if enableFireNetInterfaces {
		gw := &goaviatrix.TransitVpc{
			CloudType: d.Get("cloud_type").(int),
			GwName:    d.Get("gw_name").(string),
		}
		err := client.DisableGatewayFireNetInterfaces(gw)
		if err != nil {
			return fmt.Errorf("failed to disable transit GW for FireNet Interfaces: %s", err)
		}
	}

	//If HA is enabled, delete HA GW first.
	if haSubnet := d.Get("ha_subnet").(string); haSubnet != "" {
		gateway.GwName += "-hagw"
		err := client.DeleteGateway(gateway)
		if err != nil {
			return fmt.Errorf("failed to delete Aviatrix TransitVpc HA gateway: %s", err)
		}
	}

	gateway.GwName = d.Get("gw_name").(string)
	err := client.DeleteGateway(gateway)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix TransitVpc: %s", err)
	}

	return nil
}
