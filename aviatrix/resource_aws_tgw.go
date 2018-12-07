package aviatrix

import (
	"fmt"
	"log"

	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAWSTgw() *schema.Resource {
	return &schema.Resource{
		Create: resourceAWSTgwCreate,
		Read:   resourceAWSTgwRead,
		Update: resourceAWSTgwUpdate,
		Delete: resourceAWSTgwDelete,

		Schema: map[string]*schema.Schema{
			"aws_tgw_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aws_side_as_number": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAWSTgwCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	awsTgw := &goaviatrix.AWSTgw{
		Name:            d.Get("aws_tgw_name").(string),
		AccountName:     d.Get("account_name").(string),
		Region:          d.Get("region").(string),
		AwsSideAsNumber: d.Get("aws_side_as_number").(string),
	}

	log.Printf("[INFO] Creating AWS TGW: %#v", awsTgw)

	err := client.CreateAWSTgw(awsTgw)
	if err != nil {
		return fmt.Errorf("failed to create AWS TGW: %s", err)
	}

	d.SetId(awsTgw.Name)

	return resourceAWSTgwRead(d, meta)
}

func resourceAWSTgwRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	awsTgw := &goaviatrix.AWSTgw{
		Name:            d.Get("aws_tgw_name").(string),
		AccountName:     d.Get("account_name").(string),
		Region:          d.Get("region").(string),
		AwsSideAsNumber: d.Get("aws_side_as_number").(string),
	}

	awsTgwResp, err := client.GetAWSTgw(awsTgw)
	if err != nil {
		return fmt.Errorf("couldn't find AWS TGW: %s", err)
	}
	log.Printf("[TRACE] reading AWS TGW %s: %#v", d.Get("aws_tgw_name").(string), awsTgwResp)

	d.Set("account_name", awsTgw.AccountName)
	d.Set("aws_tgw_name", awsTgw.Name)
	d.Set("region", awsTgw.Region)
	d.Set("aws_side_as_number", awsTgw.AwsSideAsNumber)

	return nil
}

func resourceAWSTgwUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAWSTgwDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	awsTgw := &goaviatrix.AWSTgw{
		AccountName: d.Get("account_name").(string),
		Name:        d.Get("aws_tgw_name").(string),
		Region:      d.Get("region").(string),
	}

	err := client.DeleteAWSTgw(awsTgw)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find AWS TGW: %s", err)
	}

	return nil
}
