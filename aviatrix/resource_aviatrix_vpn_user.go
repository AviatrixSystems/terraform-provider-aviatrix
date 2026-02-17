package aviatrix

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixVPNUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixVPNUserCreate,
		Read:   resourceAviatrixVPNUserRead,
		Update: resourceAviatrixVPNUserUpdate,
		Delete: resourceAviatrixVPNUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceAviatrixVPNUserResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceAviatrixVPNUserStateUpgradeV0,
				Version: 0,
			},
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "VPC Id of Aviatrix VPN gateway.",
			},
			"gw_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "If ELB is enabled, this will be the name of the ELB, " +
					"else it will be the name of the Aviatrix VPN gateway.",
			},
			"dns_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "FQDN of a DNS based VPN service such as GeoVPN or UDP load balancer.",
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
	client := mustClient(meta)

	vpnUser := &goaviatrix.VPNUser{
		VpcID:        getString(d, "vpc_id"),
		GwName:       getString(d, "gw_name"),
		DnsName:      getString(d, "dns_name"),
		UserName:     getString(d, "user_name"),
		UserEmail:    getString(d, "user_email"),
		SamlEndpoint: getString(d, "saml_endpoint"),
	}

	if vpnUser.DnsName == "" {
		if vpnUser.VpcID == "" && vpnUser.GwName == "" {
			return fmt.Errorf("please set 'vpc_id' and 'gw_name', or 'dns_name' alone")
		} else if vpnUser.VpcID == "" || vpnUser.GwName == "" {
			return fmt.Errorf("please set both 'vpc_id' and 'gw_name'")
		}
	} else {
		if vpnUser.VpcID != "" || vpnUser.GwName != "" {
			return fmt.Errorf("DNS is enabled. Please set 'vpc_id' and 'gw_name' to be empty")
		}
		vpnUser.DnsEnabled = true
	}

	manageUserAttachment := getBool(d, "manage_user_attachment")
	if !manageUserAttachment && len(getList(d, "profiles")) != 0 {
		return fmt.Errorf("'manage_user_attachment' is set false. Please empty 'profiles' and manage user attachment in other resource")
	}

	log.Printf("[INFO] Creating Aviatrix VPN User: %#v", vpnUser)

	d.SetId(vpnUser.UserName)

	flag := false
	defer func() { _ = resourceAviatrixVPNUserReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	err := client.CreateVPNUser(vpnUser)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix VPNUser: %w", err)
	}

	if manageUserAttachment && len(getList(d, "profiles")) != 0 {
		for _, profileName := range getList(d, "profiles") {
			profile := &goaviatrix.Profile{
				Name: mustString(profileName),
			}
			profile.UserList = append(profile.UserList, vpnUser.UserName)
			err := client.AttachUsers(profile)
			if err != nil {
				return fmt.Errorf("failed to attach User(%s) to Profile(%s) due to: %w", vpnUser.UserName, profile.Name, err)
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
	client := mustClient(meta)

	userName := getString(d, "user_name")
	if userName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, user_name is empty. Id is %s", id)
		mustSet(d, "user_name", id)
		mustSet(d, "manage_user_attachment", true)
		d.SetId(id)
	}

	vpnUser := &goaviatrix.VPNUser{
		UserName: getString(d, "user_name"),
	}
	vu, err := client.GetVPNUser(vpnUser)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix VPNUser: %w", err)
	}

	log.Printf("[TRACE] Reading vpn_user %s: %#v", userName, vu)

	if vu != nil {
		if vu.DnsEnabled {
			mustSet(d, "dns_name", vu.DnsName)
		} else {
			mustSet(d, "vpc_id", vu.VpcID)
			mustSet(d, "gw_name", vu.GwName)
		}
		mustSet(d, "user_name", vu.UserName)
		if vu.UserEmail != "" {
			mustSet(d, "user_email", vu.UserEmail)
		}
		mustSet(d, "saml_endpoint", vu.SamlEndpoint)

		manageUserAttachment := getBool(d, "manage_user_attachment")
		if manageUserAttachment {
			var profiles []string
			for _, profile := range getList(d, "profiles") {
				profiles = append(profiles, mustString(profile))
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
	client := mustClient(meta)

	vpnUser := &goaviatrix.VPNUser{
		UserName: getString(d, "user_name"),
	}
	d.Partial(true)

	manageUserAttachment := getBool(d, "manage_user_attachment")
	if d.HasChange("manage_user_attachment") {
		_, nMUA := d.GetChange("manage_user_attachment")
		newManageUserAttachment := mustBool(nMUA)
		if newManageUserAttachment {
			mustSet(d, "manage_user_attachment", true)
		} else {
			mustSet(d, "manage_user_attachment", false)
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
			oldString := mustSlice(oldU)
			newString := mustSlice(newU)
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
					return fmt.Errorf("failed to attach User: %w", err)
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
					return fmt.Errorf("failed to detach User: %w", err)
				}
			}
		}
	} else {
		if len(getList(d, "profiles")) != 0 {
			return fmt.Errorf("'manage_user_attachment' is set false. Please empty 'profiles' and manage user attachment in other resource")
		}
	}

	d.Partial(false)
	return resourceAviatrixVPNUserRead(d, meta)
}

func resourceAviatrixVPNUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	vpnUser := &goaviatrix.VPNUser{
		UserName: getString(d, "user_name"),
		VpcID:    getString(d, "vpc_id"),
		DnsName:  getString(d, "dns_name"),
	}
	if vpnUser.DnsName != "" {
		vpnUser.DnsEnabled = true
	}

	log.Printf("[INFO] Deleting Aviatrix VPNUser: %#v", vpnUser)

	err := client.DeleteVPNUser(vpnUser)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix VPNUser: %w", err)
	}

	return nil
}
