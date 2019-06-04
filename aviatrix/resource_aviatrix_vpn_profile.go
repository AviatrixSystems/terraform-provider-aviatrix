package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixProfileCreate,
		Read:   resourceAviatrixProfileRead,
		Update: resourceAviatrixProfileUpdate,
		Delete: resourceAviatrixProfileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				Description: "Base policy rule of  the profile to be added. Enter 'allow_all' or 'deny_all'.",
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
							Optional:    true,
							Description: "The opposite of the base rule for correct behaviour. 'allow' or 'deny'.",
						},
						"proto": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Protocol to allow or deny.",
						},
						"port": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Port to be allowed or denied.",
						},
						"target": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "CIDR to be allowed or denied.",
						},
					},
				},
			},
		},
	}
}

func resourceAviatrixProfileCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*goaviatrix.Client)
	log.Printf("[INFO] Creating Aviatrix Profile: %v %T", d.Get("users"), d.Get("users"))

	profile := &goaviatrix.Profile{
		Name:     d.Get("name").(string),
		BaseRule: d.Get("base_rule").(string),
		Policy:   make([]goaviatrix.ProfileRule, 0),
	}
	if profile.Name == "" {
		return fmt.Errorf("profile name can't be empty string")
	}
	for _, user := range d.Get("users").([]interface{}) {
		profile.UserList = append(profile.UserList, user.(string))
	}
	log.Printf("[INFO] Creating Aviatrix Profile with users: %v", profile.UserList)
	names := d.Get("policy").([]interface{})
	for _, domain := range names {
		if domain != nil {
			dn := domain.(map[string]interface{})
			profileRule := &goaviatrix.ProfileRule{
				Action:   dn["action"].(string),
				Protocol: dn["proto"].(string),
				Port:     dn["port"].(string),
				Target:   dn["target"].(string),
			}
			err := client.ValidateProfileRule(profileRule)
			if err != nil {
				return fmt.Errorf("policy validation failed: %v", err)
			}
			profile.Policy = append(profile.Policy, *profileRule)
		}
	}
	log.Printf("[INFO] Creating Aviatrix Profile with Policy: %v", profile.Policy)
	err := client.CreateProfile(profile)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Profile: %s", err)
	}
	d.SetId(profile.Name)

	return resourceAviatrixProfileRead(d, meta)
}

func resourceAviatrixProfileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	profileName := d.Get("name").(string)
	if profileName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no profile name received. Import Id is %s", id)
		d.Set("name", id)
		d.SetId(id)
	}

	profile := &goaviatrix.Profile{
		Name:   d.Get("name").(string),
		Policy: make([]goaviatrix.ProfileRule, 0),
	}

	profileBase, errBase := client.GetProfileBasePolicy(profile)
	if errBase != nil {
		return fmt.Errorf("can't get profile base policy for profile: %s", profile.Name)
	}
	d.Set("base_rule", profileBase.BaseRule)

	log.Printf("[INFO] Reading Aviatrix Profile: %#v", profile)
	profile, err := client.GetProfile(profile)

	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find profile: %s", err)
	}
	d.Set("name", profile.Name)
	log.Printf("[TRACE] Profile policy %v", profile.Policy)
	log.Printf("[TRACE] Profile users %v", d.Get("users"))

	var users []string
	for _, user := range d.Get("users").([]interface{}) {
		users = append(users, user.(string))
	}
	if len(goaviatrix.Difference(users, profile.UserList)) == 0 &&
		len(goaviatrix.Difference(profile.UserList, users)) == 0 {
		d.Set("users", users)
	} else {
		d.Set("users", profile.UserList)
		log.Printf("[TRACE] Profile userlistnew %v", profile.UserList)
	}
	log.Printf("[TRACE] Profile policy %v", profile.Policy)

	var Policies []map[string]interface{}
	if profile != nil {
		for _, policy := range profile.Policy {
			policyDict := make(map[string]interface{})
			policyDict["action"] = policy.Action
			policyDict["target"] = policy.Target
			policyDict["proto"] = policy.Protocol
			policyDict["port"] = policy.Port
			Policies = append(Policies, policyDict)
		}
		d.Set("policy", Policies)
	}
	log.Printf("[INFO] Generated policies: %v", Policies)

	d.SetId(profile.Name)

	return nil
}

func resourceAviatrixProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	profile := &goaviatrix.Profile{
		Name: d.Get("name").(string),
	}
	d.Partial(true)

	for _, user := range d.Get("users").([]interface{}) {
		profile.UserList = append(profile.UserList, user.(string))
	}
	log.Printf("[INFO] Creating Aviatrix Profile with users: %v", profile.UserList)
	names := d.Get("policy").([]interface{})
	for _, domain := range names {
		dn := domain.(map[string]interface{})
		profileRule := &goaviatrix.ProfileRule{
			Action:   dn["action"].(string),
			Protocol: dn["proto"].(string),
			Port:     dn["port"].(string),
			Target:   dn["target"].(string),
		}
		err := client.ValidateProfileRule(profileRule)
		if err != nil {
			return fmt.Errorf("policy validation failed: %v", err)
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
	if d.HasChange("users") {

		oldU, newU := d.GetChange("users")
		log.Printf("[INFO] Users to be attached : %#v %#v ", oldU, newU)

		if oldU == nil {
			oldU = new([]interface{})
		}
		if newU == nil {
			newU = new([]interface{})
		}
		oldString := oldU.([]interface{})
		newString := newU.([]interface{})
		oldUserList := goaviatrix.ExpandStringList(oldString)
		newUserList := goaviatrix.ExpandStringList(newString)
		//Attach all the newly added Users
		toAddUsers := goaviatrix.Difference(newUserList, oldUserList)
		log.Printf("[INFO] Users to be attached : %#v", toAddUsers)
		profile.UserList = toAddUsers
		err := client.AttachUsers(profile)
		if err != nil {
			return fmt.Errorf("failed to attach User : %s", err)
		}
		//Detach all the removed Users
		toDelGws := goaviatrix.Difference(oldUserList, newUserList)
		log.Printf("[INFO] Users to be detached : %#v", toDelGws)
		profile.UserList = toDelGws
		err = client.DetachUsers(profile)
		if err != nil {
			return fmt.Errorf("failed to detach user : %s", err)
		}
		d.SetPartial("users")

	}
	log.Printf("[INFO] Checking for policy changes")

	if d.HasChange("policy") {

		err := client.UpdateProfilePolicy(profile)
		if err != nil {
			return fmt.Errorf("failed to create Aviatrix Profile: %s", err)
		}
		d.SetPartial("policy")

	}
	d.Partial(false)

	return nil
}

func resourceAviatrixProfileDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	profile := &goaviatrix.Profile{
		Name: d.Get("name").(string),
	}
	log.Printf("[INFO] Deleting Aviatrix Profile: %#v", profile)
	if _, ok := d.GetOk("users"); ok {
		log.Printf("[INFO] Found users: %#v", d.Get("users"))
		profile.UserList = goaviatrix.ExpandStringList(d.Get("users").([]interface{}))
		err := client.DetachUsers(profile)
		if err != nil {
			return fmt.Errorf("failed to detach Users: %s", err)
		}
	}
	err := client.DeleteProfile(profile)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Profile: %s", err)
	}
	return nil
}
