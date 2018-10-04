package aviatrix

import (
	"fmt"
	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceAviatrixProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixProfileCreate,
		Read:   resourceAviatrixProfileRead,
		Update: resourceAviatrixProfileUpdate,
		Delete: resourceAviatrixProfileDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"base_rule": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"users": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"policy": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"proto": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"target": {
							Type:     schema.TypeString,
							Optional: true,
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
	}
	for _, user := range d.Get("users").([]interface{}) {
		profile.UserList = append(profile.UserList, user.(string))
	}
	log.Printf("[INFO] Creating Aviatrix Profile with users: %v", profile.UserList)
	names := d.Get("policy").([]interface{})
	for _, domain := range names {
		dn := domain.(map[string]interface{})
		profileRule := goaviatrix.ProfileRule{
			Action:   dn["action"].(string),
			Protocol: dn["proto"].(string),
			Port:     dn["port"].(string),
			Target:   dn["target"].(string),
		}
		profile.Policy = append(profile.Policy, profileRule)
	}
	log.Printf("[INFO] Creating Aviatrix Profile with Policy: %v", profile.Policy)
	err := client.CreateProfile(profile)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Profile: %s", err)
	}
	d.SetId(profile.Name)

	return nil
}

func resourceAviatrixProfileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	profile := &goaviatrix.Profile{
		Name: d.Get("name").(string),
	}

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

	d.Set("users", profile.UserList)
	log.Printf("[TRACE] Profile userlistnew %v", profile.UserList)

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
		profileRule := goaviatrix.ProfileRule{
			Action:   dn["action"].(string),
			Protocol: dn["proto"].(string),
			Port:     dn["port"].(string),
			Target:   dn["target"].(string),
		}
		profile.Policy = append(profile.Policy, profileRule)
	}

	log.Printf("[INFO] Reading Aviatrix Profile: %#v", profile)

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
