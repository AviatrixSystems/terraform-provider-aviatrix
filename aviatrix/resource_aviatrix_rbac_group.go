package aviatrix

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixRbacGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixRbacGroupCreate,
		Read:   resourceAviatrixRbacGroupRead,
		Delete: resourceAviatrixRbacGroupDelete,
		Update: resourceAviatrixRbacGroupUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
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
	client := mustClient(meta)

	groupName := getString(d, "group_name")
	group := &goaviatrix.RbacGroup{
		GroupName: groupName,
	}

	log.Printf("[INFO] Creating Aviatrix RBAC permission group: %#v", group)

	d.SetId(group.GroupName)
	flag := false
	defer func() { _ = resourceAviatrixRbacGroupReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	err := client.CreatePermissionGroup(group)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix RBAC permission group: %w", err)
	}

	log.Printf("[DEBUG] Aviatrix RBAC permission group %s created", group.GroupName)

	if getBool(d, "local_login") {
		err := client.EnableLocalLoginForRBACGroup(groupName)
		if err != nil {
			return fmt.Errorf("failed to enable local_login for Aviatrix RBAC permission group: %w", err)
		}
	}

	return resourceAviatrixRbacGroupReadIfRequired(d, meta, &flag)
}

func resourceAviatrixRbacGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)
	groupName := getString(d, "group_name")

	if getBool(d, "local_login") {
		err := client.EnableLocalLoginForRBACGroup(groupName)
		if err != nil {
			return fmt.Errorf("failed to enable local_login for Aviatrix RBAC permission group: %w", err)
		}
	} else {
		err := client.DisableLocalLoginForRBACGroup(groupName)
		if err != nil {
			return fmt.Errorf("failed to disable local_login for Aviatrix RBAC permission group: %w", err)
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
	client := mustClient(meta)

	groupName := getString(d, "group_name")
	if groupName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no group name received. Import Id is %s", id)
		mustSet(d, "group_name", id)
		d.SetId(id)
		groupName = id
	}

	group := &goaviatrix.RbacGroup{
		GroupName: groupName,
	}

	log.Printf("[INFO] Looking for Aviatrix RBAC permission group: %#v", group)

	rGroup, err := client.GetPermissionGroupDetails(groupName)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix RBAC permission group: %w", err)
	}
	if rGroup != nil {
		mustSet(d, "group_name", rGroup.GroupName)
		mustSet(d, "local_login", rGroup.LocalLogin)
		d.SetId(rGroup.GroupName)
	}

	return nil
}

func resourceAviatrixRbacGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	group := &goaviatrix.RbacGroup{
		GroupName: getString(d, "group_name"),
	}

	log.Printf("[INFO] Deleting Aviatrix RBAC permission group: %#v", group)

	err := client.DeletePermissionGroup(group)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix RBAC permission group: %w", err)
	}

	return nil
}
