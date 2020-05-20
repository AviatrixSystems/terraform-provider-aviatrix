package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				ForceNew:    true,
				Description: "VPC Id of Aviatrix VPN gateway.",
			},
			"gw_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "If ELB is enabled, this will be the name of the ELB, " +
					"else it will be the name of the Aviatrix VPN gateway.",
			},
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPN user name.",
			},
			"user_email": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "VPN User's email.",
			},
			"saml_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				ForceNew:    true,
				Description: "This is the name of the SAML endpoint to which the user will be associated.",
			},
			"profiles": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of profiles for user to attach to.",
			},
			"manage_user_attachment": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAviatrixVPNUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	vpnUser := &goaviatrix.VPNUser{
		VpcID:        d.Get("vpc_id").(string),
		GwName:       d.Get("gw_name").(string),
		UserName:     d.Get("user_name").(string),
		UserEmail:    d.Get("user_email").(string),
		SamlEndpoint: d.Get("saml_endpoint").(string),
	}

	if vpnUser.VpcID == "" {
		return fmt.Errorf("invalid choice: vpc_id can't be empty")
	}
	if vpnUser.GwName == "" {
		return fmt.Errorf("invalid choice: gw_name can't be empty")
	}

	manageUserAttachment := d.Get("manage_user_attachment").(bool)
	if !manageUserAttachment && len(d.Get("profiles").([]interface{})) != 0 {
		return fmt.Errorf("'manage_user_attachment' is set false. Please empty 'profiles' and manage user attachment in other resource")
	}

	log.Printf("[INFO] Creating Aviatrix VPN User: %#v", vpnUser)

	err := client.CreateVPNUser(vpnUser)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix VPNUser: %s", err)
	}
	d.SetId(vpnUser.UserName)

	flag := false
	defer resourceAviatrixVPNUserReadIfRequired(d, meta, &flag)

	if manageUserAttachment && len(d.Get("profiles").([]interface{})) != 0 {
		for _, profileName := range d.Get("profiles").([]interface{}) {
			profile := &goaviatrix.Profile{
				Name: profileName.(string),
			}
			profile.UserList = append(profile.UserList, vpnUser.UserName)
			err := client.AttachUsers(profile)
			if err != nil {
				return fmt.Errorf("failed to attach User(%s) to Profile(%s) due to: %s", vpnUser.UserName, profile.Name, err)
			}
		}
	}

	return resourceAviatrixVPNUserReadIfRequired(d, meta, &flag)
}

func resourceAviatrixVPNUserReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixVPNUserRead(d, meta)
	}
	return nil
}

func resourceAviatrixVPNUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	userName := d.Get("user_name").(string)
	if userName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, user_name is empty. Id is %s", id)
		d.Set("user_name", id)
		d.Set("manage_user_attachment", true)
		d.SetId(id)
	}

	vpnUser := &goaviatrix.VPNUser{
		UserName: d.Get("user_name").(string),
	}
	vu, err := client.GetVPNUser(vpnUser)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix VPNUser: %s", err)
	}

	log.Printf("[TRACE] Reading vpn_user %s: %#v", userName, vu)

	if vu != nil {
		d.Set("vpc_id", vu.VpcID)
		d.Set("gw_name", vu.GwName)
		d.Set("user_name", vu.UserName)
		if vu.UserEmail != "" {
			d.Set("user_email", vu.UserEmail)
		}
		d.Set("saml_endpoint", vu.SamlEndpoint)

		manageUserAttachment := d.Get("manage_user_attachment").(bool)
		if manageUserAttachment {
			var profiles []string
			for _, profile := range d.Get("profiles").([]interface{}) {
				profiles = append(profiles, profile.(string))
			}
			profileFromRead := vu.Profiles
			if len(goaviatrix.Difference(profiles, profileFromRead)) == 0 && len(goaviatrix.Difference(profileFromRead, profiles)) == 0 {
				err := d.Set("profiles", profiles)
				if err != nil {
					return fmt.Errorf("couldn't set 'profiles' for vpn user: %s", vpnUser.UserName)
				}
			} else {
				err := d.Set("profiles", profileFromRead)
				if err != nil {
					return fmt.Errorf("couldn't set 'profiles' for vpn user: %s", vpnUser.UserName)
				}
			}
		}
	}

	return nil
}

func resourceAviatrixVPNUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	vpnUser := &goaviatrix.VPNUser{
		UserName: d.Get("user_name").(string),
	}
	d.Partial(true)

	manageUserAttachment := d.Get("manage_user_attachment").(bool)
	if d.HasChange("manage_user_attachment") {
		_, nMUA := d.GetChange("manage_user_attachment")
		newManageUserAttachment := nMUA.(bool)
		if newManageUserAttachment {
			d.Set("manage_user_attachment", true)
		} else {
			d.Set("manage_user_attachment", false)
		}
	}
	if manageUserAttachment {
		if d.HasChange("profiles") {
			oldU, newU := d.GetChange("profiles")
			if oldU == nil {
				oldU = new([]interface{})
			}
			if newU == nil {
				newU = new([]interface{})
			}
			oldString := oldU.([]interface{})
			newString := newU.([]interface{})
			oldProfileList := goaviatrix.ExpandStringList(oldString)
			newProfileList := goaviatrix.ExpandStringList(newString)
			toAddProfiles := goaviatrix.Difference(newProfileList, oldProfileList)
			for _, profileName := range toAddProfiles {
				profile := &goaviatrix.Profile{
					Name: profileName,
				}
				profile.UserList = []string{vpnUser.UserName}
				err := client.AttachUsers(profile)
				if err != nil {
					return fmt.Errorf("failed to attach User: %s", err)
				}
			}
			toDelProfiles := goaviatrix.Difference(oldProfileList, newProfileList)
			for _, profileName := range toDelProfiles {
				profile := &goaviatrix.Profile{
					Name: profileName,
				}
				profile.UserList = []string{vpnUser.UserName}
				err := client.DetachUsers(profile)
				if err != nil {
					return fmt.Errorf("failed to detach User: %s", err)
				}
			}
			d.SetPartial("profiles")
		}
	} else {
		if len(d.Get("profiles").([]interface{})) != 0 {
			return fmt.Errorf("'manage_user_attachment' is set false. Please empty 'profiles' and manage user attachment in other resource")
		}
	}

	d.Partial(false)
	return resourceAviatrixVPNUserRead(d, meta)
}

func resourceAviatrixVPNUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	vpnUser := &goaviatrix.VPNUser{
		UserName: d.Get("user_name").(string),
		VpcID:    d.Get("vpc_id").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix VPNUser: %#v", vpnUser)

	err := client.DeleteVPNUser(vpnUser)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix VPNUser: %s", err)
	}

	return nil
}
