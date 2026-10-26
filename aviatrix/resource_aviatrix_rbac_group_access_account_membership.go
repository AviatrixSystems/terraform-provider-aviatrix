package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixRbacGroupAccessAccountMembership() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixRbacGroupAccessAccountMembershipCreate,
		Read:   resourceAviatrixRbacGroupAccessAccountMembershipRead,
		Update: resourceAviatrixRbacGroupAccessAccountMembershipUpdate,
		Delete: resourceAviatrixRbacGroupAccessAccountMembershipDelete,

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
				Description: "RBAC permission group name. This resource is authoritative for the group's access account membership.",
			},
			"access_account_names": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "Complete set of access account names that must be members of the group (authoritative).",
			},
			"remove_access_accounts_on_destroy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, deleting this resource will remove all access accounts from the group. Default is false (the access accounts are left in place).",
			},
		},
	}
}

func resourceAviatrixRbacGroupAccessAccountMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	group := getString(d, "group_name")
	accessAccounts := expandStringSet(getSet(d, "access_account_names"))

	log.Printf("[INFO] Creating (authoritative) access account membership for group %q: %v", group, accessAccounts)

	if err := client.SetRbacGroupAccessAccounts(group, accessAccounts); err != nil {
		return fmt.Errorf("failed to set access account for RBAC group %q: %w", group, err)
	}

	d.SetId(group)
	return resourceAviatrixRbacGroupAccessAccountMembershipRead(d, meta)
}

func resourceAviatrixRbacGroupAccessAccountMembershipRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	group := getString(d, "group_name")
	if group == "" {
		group = d.Id()
		_ = d.Set("group_name", group)
	}

	log.Printf("[INFO] Reading (authoritative) access account membership for group %q", group)

	current, err := client.ListRbacGroupAccessAccounts(group)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			log.Printf("[WARN] RBAC group %q not found; removing from state", group)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to list access accounts for RBAC group %q: %w", group, err)
	}

	sort.Strings(current)
	if err := d.Set("access_account_names", stringSliceToIfaceSlice(current)); err != nil {
		return fmt.Errorf("failed to set access_account_names for %q: %w", group, err)
	}

	d.SetId(group)
	return nil
}

func resourceAviatrixRbacGroupAccessAccountMembershipUpdate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)
	group := getString(d, "group_name")

	if d.HasChange("access_account_names") {
		desired := expandStringSet(getSet(d, "access_account_names"))
		if err := client.SetRbacGroupAccessAccounts(group, desired); err != nil {
			return fmt.Errorf("failed to update access accounts for RBAC group %q: %w", group, err)
		}
	}

	return resourceAviatrixRbacGroupAccessAccountMembershipRead(d, meta)
}

func resourceAviatrixRbacGroupAccessAccountMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	if v, ok := d.GetOk("remove_access_accounts_on_destroy"); !ok || !mustBool(v) {
		return nil
	}

	group := getString(d, "group_name")

	accessAccounts, err := client.ListRbacGroupAccessAccounts(group)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("failed to list access accounts before delete for group %q: %w", group, err)
	}

	if len(accessAccounts) == 0 {
		return nil
	}

	if err := client.DeleteRbacGroupAccessAccounts(group, accessAccounts); err != nil {
		return fmt.Errorf("failed to remove access accounts for group %q: %w", group, err)
	}

	return nil
}
