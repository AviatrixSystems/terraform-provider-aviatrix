package aviatrix

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixProfileCreate,
		Read:   resourceAviatrixProfileRead,
		Update: resourceAviatrixProfileUpdate,
		Delete: resourceAviatrixProfileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceAviatrixVPNProfileResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceAviatrixVPNProfileStateUpgradeV0,
				Version: 0,
			},
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name for the VPN profile.",
			},
			"base_rule": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Base policy rule of the profile to be added. Enter 'allow_all' or 'deny_all'.",
			},
			"users": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of VPN users to attach to this profile.",
			},
			"policy": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "New security policy for the profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The opposite of the base rule for correct behaviour. 'allow' or 'deny'.",
						},
						"proto": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Protocol to allow or deny.",
						},
						"port": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Port to be allowed or denied.",
						},
						"target": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "CIDR to be allowed or denied.",
						},
					},
				},
			},
			"manage_user_attachment": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceAviatrixProfileCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	log.Printf("[INFO] Creating Aviatrix Profile: %v %T", d.Get("users"), d.Get("users"))

	profile := &goaviatrix.Profile{
		Name:     getString(d, "name"),
		BaseRule: getString(d, "base_rule"),
		Policy:   make([]goaviatrix.ProfileRule, 0),
	}
	if profile.Name == "" {
		return fmt.Errorf("profile name can't be empty string")
	}

	manageUserAttachment := getBool(d, "manage_user_attachment")
	if manageUserAttachment {
		for _, user := range getList(d, "users") {
			profile.UserList = append(profile.UserList, mustString(user))
		}
	} else {
		if len(getList(d, "users")) != 0 {
			return fmt.Errorf("'manage_user_attachment' is set false. Please empty 'users' and manage user attachment in other resource")
		}
	}

	log.Printf("[INFO] Creating Aviatrix Profile with users: %v", profile.UserList)

	names := getList(d, "policy")
	for _, domain := range names {
		if domain != nil {
			dn := mustMap(domain)
			profileRule := &goaviatrix.ProfileRule{
				Action:   mustString(dn["action"]),
				Protocol: mustString(dn["proto"]),
				Port:     mustString(dn["port"]),
				Target:   mustString(dn["target"]),
			}
			err := client.ValidateProfileRule(profileRule)
			if err != nil {
				return fmt.Errorf("policy validation failed: %w", err)
			}
			profile.Policy = append(profile.Policy, *profileRule)
		}
	}

	log.Printf("[INFO] Creating Aviatrix Profile with Policy: %v", profile.Policy)

	d.SetId(profile.Name)
	flag := false
	defer func() { _ = resourceAviatrixProfileReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	err := client.CreateProfile(profile)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Profile: %w", err)
	}

	return resourceAviatrixProfileReadIfRequired(d, meta, &flag)
}

func resourceAviatrixProfileReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixProfileRead(d, meta)
	}
	return nil
}

func resourceAviatrixProfileRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	profileName := getString(d, "name")
	if profileName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no profile name received. Import Id is %s", id)
		mustSet(d, "name", id)
		mustSet(d, "manage_user_attachment", true)
		d.SetId(id)
	}

	profile := &goaviatrix.Profile{
		Name:   getString(d, "name"),
		Policy: make([]goaviatrix.ProfileRule, 0),
	}

	profileBase, errBase := client.GetProfileBasePolicy(profile)
	if errBase != nil {
		return fmt.Errorf("can't get profile base policy for profile: %s", profile.Name)
	}
	mustSet(d, "base_rule", profileBase.BaseRule)

	log.Printf("[INFO] Reading Aviatrix Profile: %#v", profile)
	profile, err := client.GetProfile(profile)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find profile: %w", err)
	}
	mustSet(d, "name", profile.Name)
	log.Printf("[TRACE] Profile policy %v", profile.Policy)

	manageUserAttachment := getBool(d, "manage_user_attachment")
	if manageUserAttachment {
		var users []string
		for _, user := range getList(d, "users") {
			users = append(users, mustString(user))
		}
		if len(goaviatrix.Difference(users, profile.UserList)) == 0 &&
			len(goaviatrix.Difference(profile.UserList, users)) == 0 {
			mustSet(d, "users", users)
		} else {
			mustSet(d, "users", profile.UserList)
			log.Printf("[TRACE] Profile userlistnew %v", profile.UserList)
		}
	}
	log.Printf("[TRACE] Profile policy %v", profile.Policy)

	var Policies []map[string]interface{}
	for _, policy := range profile.Policy {
		policyDict := make(map[string]interface{})
		policyDict["action"] = policy.Action
		policyDict["target"] = policy.Target
		policyDict["proto"] = policy.Protocol
		policyDict["port"] = policy.Port
		Policies = append(Policies, policyDict)
	}

	if err := d.Set("policy", Policies); err != nil {
		log.Printf("[WARN] Error setting policy for (%s): %s", d.Id(), err)
	}
	log.Printf("[INFO] Generated policies: %v", Policies)

	d.SetId(profile.Name)
	return nil
}

func resourceAviatrixProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	profile := &goaviatrix.Profile{
		Name: getString(d, "name"),
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
		for _, user := range getList(d, "users") {
			profile.UserList = append(profile.UserList, mustString(user))
		}
		log.Printf("[INFO] Creating Aviatrix Profile with users: %v", profile.UserList)
	}
	names := getList(d, "policy")
	for _, domain := range names {
		dn := mustMap(domain)
		profileRule := &goaviatrix.ProfileRule{
			Action:   mustString(dn["action"]),
			Protocol: mustString(dn["proto"]),
			Port:     mustString(dn["port"]),
			Target:   mustString(dn["target"]),
		}
		err := client.ValidateProfileRule(profileRule)
		if err != nil {
			return fmt.Errorf("policy validation failed: %w", err)
		}
		profile.Policy = append(profile.Policy, *profileRule)
	}

	log.Printf("[INFO] Reading Aviatrix Profile: %#v", profile)

	if d.HasChange("name") {
		return fmt.Errorf("cannot change name of a profile")
	}
	if d.HasChange("base_rule") {
		return fmt.Errorf("cannot change base rule of a profile")
	}
	if manageUserAttachment {
		if d.HasChange("users") {
			oldU, newU := d.GetChange("users")
			log.Printf("[INFO] Users to be attached : %#v %#v ", oldU, newU)

			if oldU == nil {
				oldU = new([]interface{})
			}
			if newU == nil {
				newU = new([]interface{})
			}
			oldString := mustSlice(oldU)
			newString := mustSlice(newU)
			oldUserList := goaviatrix.ExpandStringList(oldString)
			newUserList := goaviatrix.ExpandStringList(newString)
			// Attach all the newly added Users
			toAddUsers := goaviatrix.Difference(newUserList, oldUserList)
			log.Printf("[INFO] Users to be attached : %#v", toAddUsers)
			profile.UserList = toAddUsers
			err := client.AttachUsers(profile)
			if err != nil {
				return fmt.Errorf("failed to attach User : %w", err)
			}
			// Detach all the removed Users
			toDelGws := goaviatrix.Difference(oldUserList, newUserList)
			log.Printf("[INFO] Users to be detached : %#v", toDelGws)
			profile.UserList = toDelGws
			err = client.DetachUsers(profile)
			if err != nil {
				return fmt.Errorf("failed to detach user : %w", err)
			}
		}
	} else {
		if len(getList(d, "users")) != 0 {
			return fmt.Errorf("'manage_user_attachment' is set false. Please empty 'users' and manage user attachment in other resource")
		}
	}

	log.Printf("[INFO] Checking for policy changes")
	if d.HasChange("policy") {
		err := client.UpdateProfilePolicy(profile)
		if err != nil {
			return fmt.Errorf("failed to create Aviatrix Profile: %w", err)
		}
	}

	d.Partial(false)
	return nil
}

func resourceAviatrixProfileDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	profile := &goaviatrix.Profile{
		Name: getString(d, "name"),
	}
	log.Printf("[INFO] Deleting Aviatrix Profile: %#v", profile)
	if _, ok := d.GetOk("users"); ok {
		log.Printf("[INFO] Found users: %#v", d.Get("users"))

		profile.UserList = goaviatrix.ExpandStringList(getList(d, "users"))
		err := client.DetachUsers(profile)
		if err != nil {
			return fmt.Errorf("failed to detach Users: %w", err)
		}
	}

	err := client.DeleteProfile(profile)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Profile: %w", err)
	}

	return nil
}
