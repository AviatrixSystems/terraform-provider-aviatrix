package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAviatrixTransitVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixTransitVpcCreate,
		Read:   resourceAviatrixTransitVpcRead,
		Update: resourceAviatrixTransitVpcUpdate,
		Delete: resourceAviatrixTransitVpcDelete,

		Schema: map[string]*schema.Schema{
			"cloud_type": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"account_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"gw_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_reg": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_size": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ha_subnet": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ha_gw_size": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dns_server": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag_list": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"enable_hybrid_connection": {
				Type:     schema.TypeBool,
				Optional: true,
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
		DnsServer:              d.Get("dns_server").(string),
		EnableHybridConnection: d.Get("enable_hybrid_connection").(bool),
	}
	if _, ok := d.GetOk("tag_list"); ok {
		tagList := d.Get("tag_list").([]interface{})
		tagListStr := goaviatrix.ExpandStringList(tagList)
		gateway.TagList = strings.Join(tagListStr, ",")
	}

	log.Printf("[INFO] Creating Aviatrix TransitVpc: %#v", gateway)

	err := client.LaunchTransitVpc(gateway)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix TransitVpc: %s", err)
	}
	haSubnet := d.Get("ha_subnet").(string)
	haGwSize := d.Get("ha_gw_size").(string)
	enableHybridConnection := d.Get("enable_hybrid_connection").(bool)
	d.SetId(gateway.GwName)
	d.Set("ha_subnet", "")
	d.Set("ha_gw_size", "")
	d.Set("enableHybridConnection", false)
	if haSubnet != "" {
		//Enable HA
		transitGateway := &goaviatrix.TransitVpc{
			GwName:   d.Get("gw_name").(string),
			HASubnet: haSubnet,
		}
		log.Printf("[INFO] Enabling HA on Transit Gateway: %#v", haSubnet)
		err = client.EnableHaTransitVpc(transitGateway)
		if err != nil {
			return fmt.Errorf("failed to enable2 HA Aviatrix TransitVpc: %s", err)
		}
		d.Set("ha_subnet", haSubnet)
		d.Set("ha_gw_size", gateway.VpcSize)

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
			d.Set("ha_gw_size", haGwSize)
		}
	}

	if enableHybridConnection == true {
		err := client.AttachTransitGWForHybrid(gateway)
		if err != nil {
			return fmt.Errorf("failed to enable transit GW for Hybird: %s", err)
		}
		d.Set("enableHybridConnection", true)
	}

	return resourceAviatrixTransitVpcRead(d, meta)
}

func resourceAviatrixTransitVpcRead(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("couldn't find Aviatrix TransitVpc: %s", err)
	}
	log.Printf("[TRACE] reading gateway %s: %#v",
		d.Get("gw_name").(string), gw)
	if gw != nil {
		d.Set("account_name", gw.AccountName)
		d.Set("gw_name", gw.GwName)
		//d.Set("vpc_id", gw.VpcID)
		d.Set("vpc_reg", gw.VpcRegion)
		d.Set("vpc_size", gw.GwSize)
		d.Set("enable_hybrid_connection", gw.EnableHybridConnection)
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
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix Transit HA Gateway: %s", err)
	}
	log.Printf("[INFO] Transit HA Gateway size: %s", haGw.GwSize)
	d.Set("ha_gw_size", haGw.GwSize)

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
		oldTagList := goaviatrix.ExpandStringList(os)
		tags.TagList = strings.Join(oldTagList, ",")
		err := client.DeleteTags(tags)
		if err != nil {
			return fmt.Errorf("failed to delete tags : %s", err)
		}
		newTagList := goaviatrix.ExpandStringList(ns)
		tags.TagList = strings.Join(newTagList, ",")
		err = client.AddTags(tags)
		if err != nil {
			return fmt.Errorf("failed to add tags : %s", err)
		}
		d.SetPartial("tag_list")
	}
	if d.HasChange("vpc_size") {
		gateway.GwSize = d.Get("vpc_size").(string)
		err := client.UpdateGateway(gateway)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix TransitVpc: %s", err)
		}
		d.SetPartial("vpc_size")
	}

	if d.HasChange("ha_subnet") {
		transitGateway := &goaviatrix.TransitVpc{
			GwName:   d.Get("gw_name").(string),
			HASubnet: d.Get("ha_subnet").(string),
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
			//HA subnet has been modified. Delete older HA GW,
			// and launch new HA GW in new subnet.
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
		d.SetPartial("ha_subnet")
	}

	if d.HasChange("enable_hybrid_connection") {
		transitGateway := &goaviatrix.TransitVpc{
			CloudType:   d.Get("cloud_type").(int),
			AccountName: d.Get("account_name").(string),
			GwName:      d.Get("gw_name").(string),
			VpcID:       d.Get("vpc_id").(string),
			VpcRegion:   d.Get("vpc_reg").(string),
		}
		enableHybridConnection := d.Get("enable_hybrid_connection").(bool)
		if enableHybridConnection == true {
			err := client.AttachTransitGWForHybrid(transitGateway)
			if err != nil {
				return fmt.Errorf("failed to enable transit GW for Hybird: %s", err)
			}
		}

		if enableHybridConnection == false {
			err := client.DetachTransitGWForHybrid(transitGateway)
			if err != nil {
				return fmt.Errorf("failed to disable transit GW for Hybird: %s", err)
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
		d.SetPartial("ha_gw_size")
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

	//If HA is enabled, delete HA GW first.
	if haSubnet := d.Get("ha_subnet").(string); haSubnet != "" {
		//Delete HA Gw too
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
