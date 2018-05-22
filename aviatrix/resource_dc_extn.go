package aviatrix

import (
	"fmt"
	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceDCExtn() *schema.Resource {
	return &schema.Resource{
		Create: resourceDCExtnCreate,
		Read:   resourceDCExtnRead,
		Update: resourceDCExtnUpdate,
		Delete: resourceDCExtnDelete,

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
			"vpc_reg": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"gw_size": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet_cidr": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"internet_access": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_subnet": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tunnel_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceDCExtnCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	dc_extn := &goaviatrix.DCExtn{
		CloudType:      d.Get("cloud_type").(int),
		AccountName:    d.Get("account_name").(string),
		GwName:         d.Get("gw_name").(string),
		VpcRegion:      d.Get("vpc_reg").(string),
		GwSize:         d.Get("gw_size").(string),
		SubnetCIDR:     d.Get("subnet_cidr").(string),
		InternetAccess: d.Get("internet_access").(string),
		PublicSubnet:   d.Get("public_subnet").(string),
		TunnelType:     d.Get("tunnel_type").(string),
	}

	log.Printf("[INFO] Creating Aviatrix DC Extension: %#v", dc_extn)

	err := client.CreateDCExtn(dc_extn)
	if err != nil {
		return fmt.Errorf("Failed to create Aviatrix DC Extension: %s", err)
	}
	d.SetId(dc_extn.GwName)
	return nil
	//return resourceDCExtnRead(d, meta)
}

func resourceDCExtnRead(d *schema.ResourceData, meta interface{}) error {
	return resourceAviatrixGatewayRead(d, meta)
}

func resourceDCExtnUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	dc_extn := &goaviatrix.DCExtn{}
	log.Printf("[INFO] Update available public subnet CIDR: %#v", dc_extn)
	err := client.UpdateDCExtn(dc_extn)
	if err != nil {
		return fmt.Errorf("No available public CIDR or fully exhausted: %s", err)
	}

	return nil
	//return resourceAviatrixGatewayUpdate(d, meta)
}

func resourceDCExtnDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	dc_extn := &goaviatrix.DCExtn{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string),
	}
	//If HA is enabled, delete HA GW first.
	//if ha_subnet := d.Get("ha_subnet").(string); ha_subnet != "" {
	//Delete HA Gw first
	//        log.Printf("[INFO] Deleting Aviatrix HA gateway: %#v", gateway)
	//        err := client.DisableHaGateway(gateway)
	//        if err != nil {
	//                return fmt.Errorf("Failed to delete Aviatrix HA gateway: %s", err)
	//        }
	//}
	log.Printf("[INFO] Deleting Aviatrix datacenter extension gateway: %#v", dc_extn)
	err := client.DeleteDCExtn(dc_extn)
	if err != nil {
		return fmt.Errorf("Failed to delete Aviatrix Gateway: %s", err)
	}
	d.SetId(dc_extn.GwName)
	return nil
}
