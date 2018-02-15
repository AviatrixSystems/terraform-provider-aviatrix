package aviatrix

import (
	"fmt"
	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceAviatrixGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixGatewayCreate,
		Read:   resourceAviatrixGatewayRead,
		Update: resourceAviatrixGatewayUpdate,
		Delete: resourceAviatrixGatewayDelete,

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
			"vpc_net": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAviatrixGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	gateway := &goaviatrix.Gateway{
		CloudType:   d.Get("cloud_type").(int),
		AccountName: d.Get("account_name").(string),
		GwName:      d.Get("gw_name").(string),
		VpcID:       d.Get("vpc_id").(string),
		VpcRegion:   d.Get("vpc_reg").(string),
		VpcSize:     d.Get("vpc_size").(string),
		VpcNet:      d.Get("vpc_net").(string),
	}

	log.Printf("[INFO] Creating Aviatrix gateway: %#v", gateway)

	err := client.CreateGateway(gateway)
	if err != nil {
		return fmt.Errorf("Failed to create Aviatrix Gateway: %s", err)
	}
	d.SetId(gateway.GwName)
	return resourceAviatrixGatewayRead(d, meta)
}

func resourceAviatrixGatewayRead(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("Couldn't find Aviatrix Gateway: %s", err)
	}
	log.Printf("[TRACE] reading gateway %s: %#v",
		d.Get("gw_name").(string), gw)
	if gw != nil {
		d.Set("vpc_size", gw.VpcSize)
		d.Set("public_ip", gw.PublicIP)
	}
	return nil
}

func resourceAviatrixGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	gateway := &goaviatrix.Gateway{
		GwName: d.Get("gw_name").(string),
		GwSize: d.Get("vpc_size").(string),
	}

	log.Printf("[INFO] Updating Aviatrix gateway: %#v", gateway)

	err := client.UpdateGateway(gateway)
	if err != nil {
		return fmt.Errorf("Failed to update Aviatrix Gateway: %s", err)
	}
	d.SetId(gateway.GwName)
	return nil
}

func resourceAviatrixGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	gateway := &goaviatrix.Gateway{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix gateway: %#v", gateway)

	err := client.DeleteGateway(gateway)
	if err != nil {
		return fmt.Errorf("Failed to delete Aviatrix Gateway: %s", err)
	}
	return nil
}
