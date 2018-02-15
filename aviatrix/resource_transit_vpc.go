package aviatrix

import (
	"fmt"
	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strings"
)

func resourceAviatrixTransitVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixTransitVpcCreate,
		Read:   resourceAviatrixTransitVpcRead,
		Update: resourceAviatrixTransitVpcUpdate,
		Delete: resourceAviatrixTransitVpcDelete,

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
				Required: true,
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
			"ha_subnet": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"dns_server": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag_list": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

func resourceAviatrixTransitVpcCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	gateway := &goaviatrix.TransitVpc{
		CloudType:   d.Get("cloud_type").(int),
		AccountName: d.Get("account_name").(string),
		GwName:      d.Get("gw_name").(string),
		VpcID:       d.Get("vpc_id").(string),
		VpcRegion:   d.Get("vpc_reg").(string),
		VpcSize:     d.Get("vpc_size").(string),
		Subnet:      d.Get("subnet").(string),
		DnsServer:   d.Get("dns_server").(string),
	}
	if _, ok := d.GetOk("tag_list"); ok {
		tag_list := d.Get("tag_list").([]interface{})
		tag_list_str := goaviatrix.ExpandStringList(tag_list)
		gateway.TagList = strings.Join(tag_list_str, ",")
	}

	log.Printf("[INFO] Creating Aviatrix TransitVpc: %#v", gateway)

	err := client.LaunchTransitVpc(gateway)
	if err != nil {
		return fmt.Errorf("Failed to create Aviatrix TransitVpc: %s", err)
	}
	if ha_subnet := d.Get("ha_subnet").(string); ha_subnet != "" {
		//Enable HA
		ha_gateway := &goaviatrix.TransitVpc{
			GwName:   d.Get("gw_name").(string),
			HASubnet: d.Get("ha_subnet").(string),
		}
		err = client.EnableHaTransitVpc(ha_gateway)
		if err != nil {
			return fmt.Errorf("Failed to enable2 HA Aviatrix TransitVpc: %s", err)
		}
	}
	d.SetId(gateway.GwName)
	return nil
	//return resourceAviatrixTransitVpcRead(d, meta)
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
		return fmt.Errorf("Couldn't find Aviatrix TransitVpc: %s", err)
	}
	log.Printf("[TRACE] reading gateway %s: %#v",
		d.Get("gw_name").(string), gw)
	if gw != nil {
		d.Set("account_name", gw.AccountName)
		d.Set("gw_name", gw.GwName)
		//d.Set("vpc_id", gw.VpcID)
		d.Set("vpc_reg", gw.VpcRegion)
		d.Set("vpc_size", gw.VpcSize)
	}
	return nil
}

func resourceAviatrixTransitVpcUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	gateway := &goaviatrix.Gateway{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string),
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
			return fmt.Errorf("Failed to update Aviatrix TransitVpc: %s", err)
		}
		d.SetPartial("vpc_size")
	}

	if d.HasChange("ha_subnet") {
		ha_gateway := &goaviatrix.TransitVpc{
			GwName:   d.Get("gw_name").(string),
			HASubnet: d.Get("ha_subnet").(string),
		}
		o, n := d.GetChange("ha_subnet")
		if o == "" {
			//New configuration to enable HA
			err := client.EnableHaTransitVpc(ha_gateway)
			if err != nil {
				return fmt.Errorf("Failed to enable HA Aviatrix TransitVpc: %s", err)
			}
		} else if n == "" {
			//Ha configuration has been deleted
			gateway.GwName += "-hagw"
			err := client.DeleteGateway(gateway)
			if err != nil {
				return fmt.Errorf("Failed to delete Aviatrix TransitVpc HA gateway: %s", err)
			}
		} else {
			//HA subnet has been modified. Delete older HA GW,
			// and launch new HA GW in new subnet.
			gateway.GwName += "-hagw"
			err := client.DeleteGateway(gateway)
			if err != nil {
				return fmt.Errorf("Failed to delete Aviatrix TransitVpc HA gateway: %s", err)
			}

			gateway.GwName = d.Get("gw_name").(string)
			//New configuration to enable HA
			ha_err := client.EnableHaTransitVpc(ha_gateway)
			if ha_err != nil {
				return fmt.Errorf("Failed to enable HA Aviatrix TransitVpc: %s", err)
			}
		}
		d.SetPartial("ha_subnet")
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
	if ha_subnet := d.Get("ha_subnet").(string); ha_subnet != "" {
		//Delete HA Gw too
		gateway.GwName += "-hagw"
		err := client.DeleteGateway(gateway)
		if err != nil {
			return fmt.Errorf("Failed to delete Aviatrix TransitVpc HA gateway: %s", err)
		}
	}
	gateway.GwName = d.Get("gw_name").(string)
	err := client.DeleteGateway(gateway)
	if err != nil {
		return fmt.Errorf("Failed to delete Aviatrix TransitVpc: %s", err)
	}
	return nil
}
