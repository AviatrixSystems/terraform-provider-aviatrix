package aviatrix

import (
	"fmt"
	"log"

	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAviatrixVPNUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixVPNUserCreate,
		Read:   resourceAviatrixVPNUserRead,
		Update: resourceAviatrixVPNUserUpdate,
		Delete: resourceAviatrixVPNUserDelete,

		Schema: map[string]*schema.Schema{
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"gw_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"user_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"user_email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"saml_endpoint": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAviatrixVPNUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	vpn_user := &goaviatrix.VPNUser{
		VpcID:        d.Get("vpc_id").(string),
		GwName:       d.Get("gw_name").(string),
		UserName:     d.Get("user_name").(string),
		UserEmail:    d.Get("user_email").(string),
		SamlEndpoint: d.Get("saml_endpoint").(string),
	}

	log.Printf("[INFO] Creating Aviatrix VPN User: %#v", vpn_user)

	err := client.CreateVPNUser(vpn_user)
	if err != nil {
		return fmt.Errorf("Failed to create Aviatrix VPNUser: %s", err)
	}
	d.SetId(vpn_user.UserName)
	return nil
}

func resourceAviatrixVPNUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	vpn_user := &goaviatrix.VPNUser{
		UserName: d.Get("user_name").(string),
	}
	vu, err := client.GetVPNUser(vpn_user)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find Aviatrix VPNUser: %s", err)
	}
	log.Printf("[TRACE] Reading vpn_user %s: %#v",
		d.Get("user_name").(string), vu)
	if vu != nil {
		d.Set("vpc_id", vu.VpcID)
		d.Set("gw_name", vu.GwName)
	}
	return nil
}

func resourceAviatrixVPNUserUpdate(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("Aviatrix VPNUser doesn't support update.")
}

func resourceAviatrixVPNUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	vpn_user := &goaviatrix.VPNUser{
		UserName: d.Get("user_name").(string),
		VpcID:    d.Get("vpc_id").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix VPNUser: %#v", vpn_user)

	err := client.DeleteVPNUser(vpn_user)
	if err != nil {
		return fmt.Errorf("Failed to delete Aviatrix VPNUser: %s", err)
	}
	return nil
}
