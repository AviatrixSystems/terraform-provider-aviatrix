package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixRbacGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixRbacGroupCreate,
		Read:   resourceAviatrixRbacGroupRead,
		Delete: resourceAviatrixRbacGroupDelete,
		Update: resourceAviatrixRbacGroupUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "RBAC permission group name.",
			},
			"local_login": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to allow members of an RBAC group to bypass LDAP/MFA for Duo login",
			},
		},
	}
}

func resourceAviatrixRbacGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	groupName := d.Get("group_name").(string)
	group := &goaviatrix.RbacGroup{
		GroupName: groupName,
	}

	log.Printf("[INFO] Creating Aviatrix RBAC permission group: %#v", group)

	err := client.CreatePermissionGroup(group)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix RBAC permission group: %s", err)
	}

	log.Printf("[DEBUG] Aviatrix RBAC permission group %s created", group.GroupName)
	d.SetId(group.GroupName)
	flag := false
	defer resourceAviatrixRbacGroupReadIfRequired(d, meta, &flag)
	if d.Get("local_login").(bool) {
		err := client.EnableLocalLoginForRBACGroup(groupName)
		if err != nil {
			return fmt.Errorf("failed to enable local_login for Aviatrix RBAC permission group: %s", err)
		}
	}

	return resourceAviatrixRbacGroupReadIfRequired(d, meta, &flag)
}

func resourceAviatrixRbacGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	groupName := d.Get("group_name").(string)

	if d.Get("local_login").(bool) {
		err := client.EnableLocalLoginForRBACGroup(groupName)
		if err != nil {
			return fmt.Errorf("failed to enable local_login for Aviatrix RBAC permission group: %s", err)
		}
	} else {
		err := client.DisableLocalLoginForRBACGroup(groupName)
		if err != nil {
			return fmt.Errorf("failed to disable local_login for Aviatrix RBAC permission group: %s", err)
		}
	}
	return resourceAviatrixRbacGroupRead(d, meta)
}

func resourceAviatrixRbacGroupReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixRbacGroupRead(d, meta)
	}
	return nil
}

func resourceAviatrixRbacGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	groupName := d.Get("group_name").(string)
	if groupName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no group name received. Import Id is %s", id)
		d.Set("group_name", id)
		d.SetId(id)
		groupName = id
	}

	group := &goaviatrix.RbacGroup{
		GroupName: groupName,
	}

	log.Printf("[INFO] Looking for Aviatrix RBAC permission group: %#v", group)

	rGroup, err := client.GetPermissionGroupDetails(groupName)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix RBAC permission group: %s", err)
	}
	if rGroup != nil {
		d.Set("group_name", rGroup.GroupName)
		d.Set("local_login", rGroup.LocalLogin)
		d.SetId(rGroup.GroupName)
	}

	return nil
}

func resourceAviatrixRbacGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	group := &goaviatrix.RbacGroup{
		GroupName: d.Get("group_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix RBAC permission group: %#v", group)

	err := client.DeletePermissionGroup(group)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix RBAC permission group: %s", err)
	}

	return nil
}
