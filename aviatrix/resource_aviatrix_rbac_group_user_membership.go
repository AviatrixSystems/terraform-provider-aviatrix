package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixRbacGroupUserMembership() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixRbacGroupUserMembershipCreate,
		Read:   resourceAviatrixRbacGroupUserMembershipRead,
		Update: resourceAviatrixRbacGroupUserMembershipUpdate,
		Delete: resourceAviatrixRbacGroupUserMembershipDelete,

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
				if d.Id() == "" {
					return nil, fmt.Errorf("import requires group_name as ID")
				}
				_ = d.Set("group_name", d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "RBAC permission group name. This resource is authoritative for the group's user membership.",
			},
			"user_names": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "Complete set of user names that must be members of the group (authoritative).",
			},
			"remove_users_on_destroy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, deleting this resource will remove all users from the group. Default is false (the users are left in place).",
			},
		},
	}
}

func resourceAviatrixRbacGroupUserMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	group := getString(d, "group_name")
	users := expandStringSet(getSet(d, "user_names"))

	log.Printf("[INFO] Creating (authoritative) user membership for group %q: %v", group, users)

	if err := client.SetRbacGroupUsers(group, users); err != nil {
		return fmt.Errorf("failed to set user for RBAC group %q: %w", group, err)
	}

	d.SetId(group)
	return resourceAviatrixRbacGroupUserMembershipRead(d, meta)
}

func resourceAviatrixRbacGroupUserMembershipRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	group := getString(d, "group_name")
	if group == "" {
		group = d.Id()
		_ = d.Set("group_name", group)
	}

	log.Printf("[INFO] Reading (authoritative) user membership for group %q", group)

	current, err := client.ListRbacGroupUsers(group)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			log.Printf("[WARN] RBAC group %q not found; removing from state", group)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to list users for RBAC group %q: %w", group, err)
	}

	sort.Strings(current)
	if err := d.Set("user_names", stringSliceToIfaceSlice(current)); err != nil {
		return fmt.Errorf("failed to set user_names for %q: %w", group, err)
	}

	d.SetId(group)
	return nil
}

func resourceAviatrixRbacGroupUserMembershipUpdate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)
	group := getString(d, "group_name")

	if d.HasChange("user_names") {
		desired := expandStringSet(getSet(d, "user_names"))
		if err := client.SetRbacGroupUsers(group, desired); err != nil {
			return fmt.Errorf("failed to update users for RBAC group %q: %w", group, err)
		}
	}

	return resourceAviatrixRbacGroupUserMembershipRead(d, meta)
}

func resourceAviatrixRbacGroupUserMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	if v, ok := d.GetOk("remove_users_on_destroy"); !ok || !mustBool(v) {
		return nil
	}

	group := getString(d, "group_name")

	users, err := client.ListRbacGroupUsers(group)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("failed to list users before delete for group %q: %w", group, err)
	}

	if len(users) == 0 {
		return nil
	}

	if err := client.DeleteRbacGroupUsers(group, users); err != nil {
		return fmt.Errorf("failed to remove users for group %q: %w", group, err)
	}

	return nil
}
