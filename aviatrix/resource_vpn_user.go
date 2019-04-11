package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixVPNUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixVPNUserCreate,
		Read:   resourceAviatrixVPNUserRead,
		Update: resourceAviatrixVPNUserUpdate,
		Delete: resourceAviatrixVPNUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC Id of Aviatrix VPN gateway.",
			},
			"gw_name": {
				Type:     schema.TypeString,
				Required: true,
				Description: "If ELB is enabled, this will be the name of the ELB, " +
					"else it will be the name of the Aviatrix VPN gateway.",
			},
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPN user name.",
			},
			"user_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "VPN User's email.",
			},
			"saml_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "This is the name of the SAML endpoint to which the user is to be associated.",
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
	if vpn_user.VpcID == "" {
		return fmt.Errorf("invalid choice: vpc_id can't be empty")
	}
	if vpn_user.GwName == "" {
		return fmt.Errorf("invalid choice: gw_name can't be empty")
	}
	log.Printf("[INFO] Creating Aviatrix VPN User: %#v", vpn_user)

	err := client.CreateVPNUser(vpn_user)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix VPNUser: %s", err)
	}
	d.SetId(vpn_user.UserName)
	return resourceAviatrixVPNUserRead(d, meta)
}

func resourceAviatrixVPNUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	user_name := d.Get("user_name").(string)
	if user_name == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, user_name is empty. Id is %s", id)
		user_name = id
	}
	vpn_user := &goaviatrix.VPNUser{
		UserName: user_name,
	}
	vu, err := client.GetVPNUser(vpn_user)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix VPNUser: %s", err)
	}
	log.Printf("[TRACE] Reading vpn_user %s: %#v",
		user_name, vu)
	if vu != nil {
		d.Set("vpc_id", vu.VpcID)
		d.Set("gw_name", vu.GwName)
		d.Set("user_name", user_name)
		if vu.UserEmail != "" {
			d.Set("user_email", vu.UserEmail)
		}
		d.Set("saml_endpoint", vu.SamlEndpoint)
	}
	return nil
}

func resourceAviatrixVPNUserUpdate(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("the AviatrixVPNUser resource doesn't support update")
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
		return fmt.Errorf("failed to delete Aviatrix VPNUser: %s", err)
	}
	return nil
}
