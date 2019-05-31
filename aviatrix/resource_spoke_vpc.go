package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixSpokeVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSpokeVpcCreate,
		Read:   resourceAviatrixSpokeVpcRead,
		Update: resourceAviatrixSpokeVpcUpdate,
		Delete: resourceAviatrixSpokeVpcDelete,
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
				Optional:    true,
				Description: "VPC-ID/VNet-Name of cloud provider.",
			},
			"vnet_and_resource_group_names": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The string consisted of name of (Azure) VNet and name Resource-Group.",
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
				Description: "Public Subnet Info.",
			},
			"enable_nat": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "no",
				Description: "Specify whether enabling NAT feature on the gateway or not.",
			},
			"ha_subnet": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "HA Subnet. Required if enabling HA for AWS/ARM.",
			},
			"ha_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "HA Zone. Required if enabling HA for GCP.",
			},
			"ha_gw_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "HA Gateway Size.",
			},
			"single_az_ha": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "disabled",
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
			"cloud_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud instance ID.",
			},
		},
	}
}

func resourceAviatrixSpokeVpcCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	gateway := &goaviatrix.SpokeVpc{
		CloudType:      d.Get("cloud_type").(int),
		AccountName:    d.Get("account_name").(string),
		GwName:         d.Get("gw_name").(string),
		VpcID:          d.Get("vpc_id").(string),
		VnetRsrcGrp:    d.Get("vnet_and_resource_group_names").(string),
		VpcRegion:      d.Get("vpc_reg").(string),
		VpcSize:        d.Get("vpc_size").(string),
		Subnet:         d.Get("subnet").(string),
		HASubnet:       d.Get("ha_subnet").(string),
		EnableNAT:      d.Get("enable_nat").(string),
		SingleAzHa:     d.Get("single_az_ha").(string),
		TransitGateway: d.Get("transit_gw").(string),
	}

	if gateway.EnableNAT != "yes" {
		gateway.EnableNAT = "no"
	}
	if gateway.CloudType == 1 || gateway.CloudType == 4 {
		gateway.VnetRsrcGrp = ""
		d.Set("vnet_and_resource_group_names", gateway.VnetRsrcGrp)
	} else if gateway.CloudType == 8 {
		gateway.VpcID = ""
		d.Set("vpc_id", gateway.VpcID)
	} else {
		return fmt.Errorf("invalid cloud type, it can only be aws (1), gcp (4), arm (8)")
	}

	haZone := d.Get("ha_zone").(string)
	haSubnet := d.Get("ha_subnet").(string)
	haGwSize := d.Get("ha_gw_size").(string)
	if haZone != "" || haSubnet != "" {
		if haGwSize == "" {
			return fmt.Errorf("A valid non empty ha_gw_size parameter is mandatory for this resource if " +
				"ha_subnet or ha_zone is set. Example: t2.micro")
		}
	}

	log.Printf("[INFO] Creating Aviatrix Spoke VPC: %#v", gateway)

	err := client.LaunchSpokeVpc(gateway)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Spoke VPC: %s", err)
	}
	d.SetId(gateway.GwName)

	if enableNAT := d.Get("enable_nat").(string); enableNAT == "yes" {
		log.Printf("[INFO] Aviatrix NAT enabled gateway: %#v", gateway)
	}
	if singleAZHA := d.Get("single_az_ha").(string); singleAZHA == "enabled" {
		singleAZGateway := &goaviatrix.Gateway{
			GwName:   d.Get("gw_name").(string),
			SingleAZ: d.Get("single_az_ha").(string),
		}
		log.Printf("[INFO] Enable Single AZ GW HA: %#v", singleAZGateway)
		err := client.EnableSingleAZGateway(singleAZGateway)
		if err != nil {
			return fmt.Errorf("failed to create single AZ GW HA: %s", err)
		}
	}
	if haSubnet != "" || haZone != "" {
		//Enable HA
		haGateway := &goaviatrix.SpokeVpc{
			GwName:    d.Get("gw_name").(string),
			CloudType: d.Get("cloud_type").(int),
			HASubnet:  haSubnet,
			HAZone:    haZone,
		}
		err = client.EnableHaSpokeVpc(haGateway)
		if err != nil {
			return fmt.Errorf("failed to enable HA Aviatrix TransitVpc: %s", err)
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
			err := client.UpdateGateway(haGateway)
			log.Printf("[INFO] Resizing Spoke HA GAteway size to: %s ", haGateway.GwSize)
			if err != nil {
				return fmt.Errorf("failed to update Aviatrix Transit HA Gateway size: %s", err)
			}
			d.Set("ha_gw_size", haGwSize)
		}
	}

	if _, ok := d.GetOk("tag_list"); ok && gateway.CloudType == 1 {
		tagList := d.Get("tag_list").([]interface{})
		tagListStr := goaviatrix.ExpandStringList(tagList)
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
	} else if ok && gateway.CloudType != 1 {
		return fmt.Errorf("adding tags only supported for aws, cloud_type must be 1")
	}
	if transitGwName := d.Get("transit_gw").(string); transitGwName != "" {
		//No HA config, just return
		err := client.SpokeJoinTransit(gateway)
		if err != nil {
			return fmt.Errorf("failed to join TransitVpc: %s", err)
		}
	}

	return resourceAviatrixSpokeVpcRead(d, meta)
}

func resourceAviatrixSpokeVpcRead(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("couldn't find Aviatrix SpokeVpc: %s", err)
	}
	log.Printf("[TRACE] reading spoke gateway %s: %#v",
		d.Get("gw_name").(string), gw)
	if gw != nil {
		d.Set("cloud_type", gw.CloudType)
		d.Set("account_name", gw.AccountName)
		if gw.CloudType == 1 {
			d.Set("vpc_id", strings.Split(gw.VpcID, "~~")[0]) //aws vpc_id returns as <vpc_id>~~<other vpc info> in rest api
			d.Set("vpc_reg", gw.VpcRegion)                    //aws vpc_reg returns as vpc_region in rest api
		} else if gw.CloudType == 4 {
			d.Set("vpc_id", strings.Split(gw.VpcID, "~-~")[0]) //gcp vpc_id returns as <vpc_id>~-~<other vpc info> in rest api
			d.Set("vpc_reg", gw.GatewayZone)                   //gcp vpc_reg returns as gateway_zone in json
		} else if gw.CloudType == 8 {
			d.Set("vnet_and_resource_group_names", gw.VpcID)
			d.Set("vpc_reg", gw.VpcRegion)
		}
		d.Set("subnet", gw.VpcNet)
		d.Set("vpc_size", gw.GwSize)
		d.Set("public_ip", gw.PublicIP)
		d.Set("cloud_instance_id", gw.CloudnGatewayInstID)
		d.Set("enable_nat", gw.EnableNat)

		if gw.SingleAZ == "yes" {
			d.Set("single_az_ha", "enabled")
		} else {
			d.Set("single_az_ha", "disabled")
		}
	}

	if gw.SpokeVpc == "yes" {
		d.Set("transit_gw", gw.TransitGwName)
	} else {
		d.Set("transit_gw", "")
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
			d.Set("tag_list", tagList)
		} else {
			d.Set("tag_list", tagListStr)
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
		} else {
			return fmt.Errorf("couldn't find Aviatrix SpokeVpc HA Gateway: %s", err)
		}
	} else {
		log.Printf("[INFO] Spoke HA Gateway size: %s", haGw.GwSize)
		if haGw.CloudType == 1 || haGw.CloudType == 8 {
			d.Set("ha_subnet", haGw.VpcNet)
			d.Set("ha_zone", "")
		} else if haGw.CloudType == 4 {
			d.Set("ha_zone", haGw.GatewayZone)
			d.Set("ha_subnet", "")
		}
		d.Set("ha_gw_size", haGw.GwSize)
	}
	return nil
}

func resourceAviatrixSpokeVpcUpdate(d *schema.ResourceData, meta interface{}) error {
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
	if d.HasChange("vnet_and_resource_group_names") {
		return fmt.Errorf("updating vnet_and_resource_group_names is not allowed")
	}
	if d.HasChange("subnet") {
		return fmt.Errorf("updating subnet is not allowed")
	}
	if d.HasChange("single_az_ha") {
		_, singleAz := d.GetChange("single_az_ha")
		singleAZGateway := &goaviatrix.Gateway{
			GwName:   d.Get("gw_name").(string),
			SingleAZ: singleAz.(string),
		}
		if singleAz == "enabled" {
			log.Printf("[INFO] Enable Single AZ GW HA: %#v", singleAZGateway)
			err := client.EnableSingleAZGateway(singleAZGateway)
			if err != nil {
				return fmt.Errorf("failed to enable single AZ GW HA: %s", err)
			}
		} else if singleAz == "disabled" {
			log.Printf("[INFO] Enable Single AZ GW HA: %#v", singleAZGateway)
			err := client.DisableSingleAZGateway(singleAZGateway)
			if err != nil {
				return fmt.Errorf("failed to enable single AZ GW HA: %s", err)
			}
		} else {
			return fmt.Errorf("single_az_ha should be only 'enabled/disabled'")
		}
	}

	if d.HasChange("tag_list") && gateway.CloudType == 1 {
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
		d.SetPartial("tag_list")
	} else if d.HasChange("tag_list") && gateway.CloudType != 1 {
		return fmt.Errorf("adding tags is only supported for aws, cloud_type must be set to 1")
	}

	//Get primary gw size if vpc_size changed, to be used later on for ha gateway size update
	primaryGwSize := d.Get("vpc_size").(string)
	if d.HasChange("vpc_size") {
		old, _ := d.GetChange("vpc_size")
		primaryGwSize = old.(string)
		gateway.GwSize = d.Get("vpc_size").(string)
		err := client.UpdateGateway(gateway)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix SpokeVpc: %s", err)
		}
		d.SetPartial("vpc_size")
	}

	newHaGwEnabled := false
	if d.HasChange("ha_subnet") || d.HasChange("ha_zone") {
		spokeGw := &goaviatrix.SpokeVpc{
			GwName:    d.Get("gw_name").(string),
			CloudType: d.Get("cloud_type").(int),
		}

		oldSubnet, newSubnet := d.GetChange("ha_subnet")
		oldZone, newZone := d.GetChange("ha_zone")
		deleteHaGw := false
		changeHaGw := false
		if spokeGw.CloudType == 1 || spokeGw.CloudType == 8 {
			spokeGw.HASubnet = d.Get("ha_subnet").(string)
			if oldSubnet == "" && newSubnet != "" {
				newHaGwEnabled = true
			} else if oldSubnet != "" && newSubnet == "" {
				deleteHaGw = true
			} else if oldSubnet != "" && newSubnet != "" {
				changeHaGw = true
			}
		} else if spokeGw.CloudType == 4 {
			spokeGw.HAZone = d.Get("ha_zone").(string)
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
			err := client.EnableHaSpokeVpc(spokeGw)
			if err != nil {
				return fmt.Errorf("failed to enable HA Aviatrix SpokeVpc: %s", err)
			}
			newHaGwEnabled = true
		} else if deleteHaGw {
			//Ha configuration has been deleted
			err := client.DeleteGateway(haGateway)
			if err != nil {
				return fmt.Errorf("failed to delete Aviatrix SpokeVpc HA gateway: %s", err)
			}
		} else if changeHaGw {
			//HA subnet has been modified. Delete older HA GW,
			// and launch new HA GW in new subnet.
			err := client.DeleteGateway(haGateway)
			if err != nil {
				return fmt.Errorf("failed to delete Aviatrix SpokeVpc HA gateway: %s", err)
			}

			gateway.GwName = d.Get("spokeGw_name").(string)
			//New configuration to enable HA
			haErr := client.EnableHaSpokeVpc(spokeGw)
			if haErr != nil {
				return fmt.Errorf("failed to enable HA Aviatrix SpokeVpc: %s", err)
			}
		}
		d.SetPartial("ha_subnet")
		d.SetPartial("ha_zone")
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
					return nil
				}
				return fmt.Errorf("couldn't find Aviatrix Spoke HA Gateway while trying to update HA Gw "+
					"size: %s", err)
			}
			haGateway.GwSize = d.Get("ha_gw_size").(string)
			if haGateway.GwSize == "" {
				return fmt.Errorf("A valid non empty ha_gw_size parameter is mandatory for this resource if " +
					"ha_subnet or ha_zone is set. Example: t2.micro or us-west1-b")
			}
			err = client.UpdateGateway(haGateway)
			log.Printf("[INFO] Updating HA Gateway size to: %s ", haGateway.GwSize)
			if err != nil {
				return fmt.Errorf("failed to update Aviatrix Spoke HA Gw size: %s", err)
			}
		}
		d.SetPartial("ha_gw_size")
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
		d.SetPartial("vpc_size")
	}
	if d.HasChange("transit_gw") {
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
				return fmt.Errorf("failed to join transit VPC: %s", err)
			}
		} else if n == "" {
			//Transit GW has been deleted, leave transit GW.
			err := client.SpokeLeaveTransit(spokeVPC)
			if err != nil {
				return fmt.Errorf("failed to leave transit VPC: %s", err)
			}
		} else {
			//Change transit GW
			err := client.SpokeLeaveTransit(spokeVPC)
			if err != nil {
				return fmt.Errorf("failed to leave transit VPC: %s", err)
			}

			err = client.SpokeJoinTransit(spokeVPC)
			if err != nil {
				return fmt.Errorf("failed to join transit VPC: %s", err)
			}
		}
		d.SetPartial("transit_gw")
	}
	d.Partial(false)

	d.SetId(gateway.GwName)
	return resourceAviatrixSpokeVpcRead(d, meta)
}

func resourceAviatrixSpokeVpcDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	gateway := &goaviatrix.Gateway{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix Spoke VPC: %#v", gateway)

	if transitGw := d.Get("transit_gw").(string); transitGw != "" {
		spokeVPC := &goaviatrix.SpokeVpc{
			GwName: d.Get("gw_name").(string),
		}

		err := client.SpokeLeaveTransit(spokeVPC)
		if err != nil {
			return fmt.Errorf("failed to leave transit VPC: %s", err)
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
			return fmt.Errorf("failed to delete Aviatrix SpokeVpc HA gateway: %s", err)
		}
	}
	gateway.GwName = d.Get("gw_name").(string)
	err := client.DeleteGateway(gateway)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix SpokeVpc: %s", err)
	}
	return nil
}
