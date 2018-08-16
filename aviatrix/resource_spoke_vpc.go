package aviatrix

import (
	"fmt"
	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strings"
)

func resourceAviatrixSpokeVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSpokeVpcCreate,
		Read:   resourceAviatrixSpokeVpcRead,
		Update: resourceAviatrixSpokeVpcUpdate,
		Delete: resourceAviatrixSpokeVpcDelete,

		Schema: map[string]*schema.Schema{
			"cloud_type": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"account_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"gw_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"vnet_and_resource_group_names": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_reg": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_size": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_nat": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"ha_subnet": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"dns_server": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"transit_gw": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag_list": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"cloud_instance_id": {
				Type:     schema.TypeString,
				Computed: true,
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
		DnsServer:      d.Get("dns_server").(string),
		TransitGateway: d.Get("transit_gw").(string),
	}
	if cloud_type := d.Get("cloud_type").(int); cloud_type == 1 {
		gateway.VnetRsrcGrp = ""
		d.Set("vnet_and_resource_group_names", gateway.VnetRsrcGrp)
	}
	if cloud_type := d.Get("cloud_type").(int); cloud_type == 8 {
		gateway.VpcID = ""
		d.Set("vpc_id", gateway.VpcID)
	}
	if _, ok := d.GetOk("tag_list"); ok {
		tag_list := d.Get("tag_list").([]interface{})
		tag_list_str := goaviatrix.ExpandStringList(tag_list)
		gateway.TagList = strings.Join(tag_list_str, ",")
	}
	log.Printf("[INFO] Creating Aviatrix Spoke VPC: %#v", gateway)

	err := client.LaunchSpokeVpc(gateway)
	if err != nil {
		return fmt.Errorf("Failed to create Aviatrix Spoke VPC: %s", err)
	}
	if ha_subnet := d.Get("ha_subnet").(string); ha_subnet != "" {
		//Enable HA
		ha_gateway := &goaviatrix.SpokeVpc{
			GwName:   d.Get("gw_name").(string),
			HASubnet: d.Get("ha_subnet").(string),
		}
		err = client.EnableHaSpokeVpc(ha_gateway)
		if err != nil {
			return fmt.Errorf("Failed to enable HA Aviatrix TransitVpc: %s", err)
		}
	}

	if transit_gw := d.Get("transit_gw").(string); transit_gw != "" {
		//No HA config, just return
		err := client.SpokeJoinTransit(gateway)
		if err != nil {
			return fmt.Errorf("Failed to join TransitVpc: %s", err)
		}
	}

	d.SetId(gateway.GwName)
	return nil
	//return resourceAviatrixSpokeVpcRead(d, meta)
}

func resourceAviatrixSpokeVpcRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
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
		return fmt.Errorf("Couldn't find Aviatrix SpokeVpc: %s", err)
	}
	log.Printf("[TRACE] reading spoke gateway %s: %#v",
		d.Get("gw_name").(string), gw)
	if gw != nil {
		d.Set("vpc_size", gw.VpcSize)
		d.Set("public_ip", gw.PublicIP)
		d.Set("cloud_instance_id", gw.CloudnGatewayInstID)
	}
	return nil
}

func resourceAviatrixSpokeVpcUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	gateway := &goaviatrix.Gateway{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string),
	}
	log.Printf("[INFO] Updating Aviatrix gateway: %#v", gateway)

	d.Partial(true)
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
		old_tag_list := goaviatrix.ExpandStringList(os)
		tags.TagList = strings.Join(old_tag_list, ",")
		err := client.DeleteTags(tags)
		if err != nil {
			return fmt.Errorf("Failed to delete tags : %s", err)
		}
		new_tag_list := goaviatrix.ExpandStringList(ns)
		tags.TagList = strings.Join(new_tag_list, ",")
		err = client.AddTags(tags)
		if err != nil {
			return fmt.Errorf("Failed to add tags : %s", err)
		}
		d.SetPartial("tag_list")
	}
	if d.HasChange("vpc_size") {
		gateway.VpcSize = d.Get("vpc_size").(string)
		err := client.UpdateGateway(gateway)
		if err != nil {
			return fmt.Errorf("Failed to update Aviatrix SpokeVpc: %s", err)
		}
		d.SetPartial("vpc_size")
	}

	if d.HasChange("ha_subnet") {
		ha_gateway := &goaviatrix.SpokeVpc{
			GwName:   d.Get("gw_name").(string),
			HASubnet: d.Get("ha_subnet").(string),
		}
		o, n := d.GetChange("ha_subnet")
		if o == "" {
			//New configuration to enable HA
			err := client.EnableHaSpokeVpc(ha_gateway)
			if err != nil {
				return fmt.Errorf("Failed to enable HA Aviatrix SpokeVpc: %s", err)
			}
		} else if n == "" {
			//Ha configuration has been deleted
			gateway.GwName += "-hagw"
			err := client.DeleteGateway(gateway)
			if err != nil {
				return fmt.Errorf("Failed to delete Aviatrix SpokeVpc HA gateway: %s", err)
			}
		} else {
			//HA subnet has been modified. Delete older HA GW,
			// and launch new HA GW in new subnet.
			gateway.GwName += "-hagw"
			err := client.DeleteGateway(gateway)
			if err != nil {
				return fmt.Errorf("Failed to delete Aviatrix SpokeVpc HA gateway: %s", err)
			}

			gateway.GwName = d.Get("gw_name").(string)
			//New configuration to enable HA
			ha_err := client.EnableHaSpokeVpc(ha_gateway)
			if ha_err != nil {
				return fmt.Errorf("Failed to enable HA Aviatrix SpokeVpc: %s", err)
			}
		}
		d.SetPartial("ha_subnet")
	}
	if d.HasChange("transit_gw") {
		spoke_vpc := &goaviatrix.SpokeVpc{
			CloudType:      d.Get("cloud_type").(int),
			GwName:         d.Get("gw_name").(string),
			HASubnet:       d.Get("ha_subnet").(string),
			TransitGateway: d.Get("transit_gw").(string),
		}

		o, n := d.GetChange("transit_gw")
		if o == "" {
			//New configuration to join to transit GW
			err := client.SpokeJoinTransit(spoke_vpc)
			if err != nil {
				return fmt.Errorf("Failed to join transit VPC: %s", err)
			}
		} else if n == "" {
			//Transit GW has been deleted, leave transit GW.
			err := client.SpokeLeaveTransit(spoke_vpc)
			if err != nil {
				return fmt.Errorf("Failed to leave transit VPC: %s", err)
			}
		} else {
			//Change transit GW
			err := client.SpokeLeaveTransit(spoke_vpc)
			if err != nil {
				return fmt.Errorf("Failed to leave transit VPC: %s", err)
			}

			err = client.SpokeJoinTransit(spoke_vpc)
			if err != nil {
				return fmt.Errorf("Failed to join transit VPC: %s", err)
			}
		}
		d.SetPartial("transit_gw")

	}
	d.Partial(false)
	//d.SetId(gateway.GwName)
	return nil
}

func resourceAviatrixSpokeVpcDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	gateway := &goaviatrix.Gateway{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix Spoke VPC: %#v", gateway)

	if transit_gw := d.Get("transit_gw").(string); transit_gw != "" {
		spoke_vpc := &goaviatrix.SpokeVpc{
			GwName: d.Get("gw_name").(string),
		}

		err := client.SpokeLeaveTransit(spoke_vpc)
		if err != nil {
			return fmt.Errorf("Failed to leave transit VPC: %s", err)
		}
	}

	//If HA is enabled, delete HA GW first.
	if ha_subnet := d.Get("ha_subnet").(string); ha_subnet != "" {
		//Delete HA Gw too
		gateway.GwName += "-hagw"
		err := client.DeleteGateway(gateway)
		if err != nil {
			return fmt.Errorf("Failed to delete Aviatrix SpokeVpc HA gateway: %s", err)
		}
	}
	gateway.GwName = d.Get("gw_name").(string)
	err := client.DeleteGateway(gateway)
	if err != nil {
		return fmt.Errorf("Failed to delete Aviatrix SpokeVpc: %s", err)
	}
	return nil
}
