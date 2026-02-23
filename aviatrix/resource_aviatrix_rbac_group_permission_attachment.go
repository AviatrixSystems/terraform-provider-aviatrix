package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixRbacGroupPermissionAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixRbacGroupPermissionAttachmentCreate,
		Read:   resourceAviatrixRbacGroupPermissionAttachmentRead,
		Delete: resourceAviatrixRbacGroupPermissionAttachmentDelete,
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
			"permission_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"all_dashboard_write",
					"all_accounts_write",
					"all_gateway_write",
					"all_tgw_orchestrator_write",
					"all_transit_network_write",
					"all_firewall_network_write",
					"all_cloudn_write",
					"all_peering_write",
					"all_site2cloud_write",
					"all_openvpn_write",
					"all_security_write",
					"all_useful_tools_write",
					"all_troubleshoot_write",
					"all_write",
				}, false),
				Description: "Permission name.",
			},
		},
	}
}

func resourceAviatrixRbacGroupPermissionAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	attachment := &goaviatrix.RbacGroupPermissionAttachment{
		GroupName:      getString(d, "group_name"),
		PermissionName: getString(d, "permission_name"),
	}

	log.Printf("[INFO] Creating Aviatrix RBAC group permission attachment: %#v", attachment)

	d.SetId(attachment.GroupName + "~" + attachment.PermissionName)
	flag := false
	defer func() { _ = resourceAviatrixRbacGroupPermissionAttachmentReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	err := client.CreateRbacGroupPermissionAttachment(attachment)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix RBAC group permission attachment: %w", err)
	}

	log.Printf("[DEBUG] Aviatrix RBAC group permission attachment created")

	return resourceAviatrixRbacGroupPermissionAttachmentReadIfRequired(d, meta, &flag)
}

func resourceAviatrixRbacGroupPermissionAttachmentReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixRbacGroupPermissionAttachmentRead(d, meta)
	}
	return nil
}

func resourceAviatrixRbacGroupPermissionAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	groupName := getString(d, "group_name")
	permissionName := getString(d, "permission_name")
	if groupName == "" || permissionName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no group name or permission name received. Import Id is %s", id)
		mustSet(d, "group_name", strings.Split(id, "~")[0])
		mustSet(d, "permission_name", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	attachment := &goaviatrix.RbacGroupPermissionAttachment{
		GroupName:      getString(d, "group_name"),
		PermissionName: getString(d, "permission_name"),
	}

	log.Printf("[INFO] Looking for Aviatrix RBAC group permission attachment: %#v", attachment)

	permissionAttachment, err := client.GetRbacGroupPermissionAttachment(attachment)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix RBAC group permission attachment: %w", err)
	}
	if permissionAttachment != nil {
		mustSet(d, "group_name", permissionAttachment.GroupName)
		mustSet(d, "permission_name", permissionAttachment.PermissionName)
		d.SetId(permissionAttachment.GroupName + "~" + permissionAttachment.PermissionName)
	}

	return nil
}

func resourceAviatrixRbacGroupPermissionAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	attachment := &goaviatrix.RbacGroupPermissionAttachment{
		GroupName:      getString(d, "group_name"),
		PermissionName: getString(d, "permission_name"),
	}

	log.Printf("[INFO] Deleting Aviatrix RBAC group permission attachment: %#v", attachment)

	err := client.DeleteRbacGroupPermissionAttachment(attachment)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix RBAC group permission attachment: %w", err)
	}

	return nil
}
