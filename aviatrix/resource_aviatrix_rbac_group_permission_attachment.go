package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixRbacGroupPermissionAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixRbacGroupPermissionAttachmentCreate,
		Read:   resourceAviatrixRbacGroupPermissionAttachmentRead,
		Delete: resourceAviatrixRbacGroupPermissionAttachmentDelete,
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
			"permission_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Permission name.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					mapPermissionName := map[string]bool{
						"all_dashboard_write":        true,
						"all_accounts_write":         true,
						"all_gateway_write":          true,
						"all_tgw_orchestrator_write": true,
						"all_transit_network_write":  true,
						"all_firewall_network_write": true,
						"all_cloud_wan_write":        true,
						"all_peering_write":          true,
						"all_site2cloud_write":       true,
						"all_openvpn_write":          true,
						"all_security_write":         true,
						"all_useful_tools_write":     true,
						"all_troubleshoot_write":     true,
						"all_write":                  true,
					}
					v := val.(string)
					if _, ok := mapPermissionName[v]; !ok {
						errs = append(errs, fmt.Errorf("permission_name: '%s' is not supported", val))
					}
					return
				},
			},
		},
	}
}

func resourceAviatrixRbacGroupPermissionAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	attachment := &goaviatrix.RbacGroupPermissionAttachment{
		GroupName:      d.Get("group_name").(string),
		PermissionName: d.Get("permission_name").(string),
	}

	log.Printf("[INFO] Creating Aviatrix RBAC group permission attachment: %#v", attachment)

	err := client.CreateRbacGroupPermissionAttachment(attachment)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix RBAC group permission attachment: %s", err)
	}

	log.Printf("[DEBUG] Aviatrix RBAC group permission attachment created")

	d.SetId(attachment.GroupName + "~" + attachment.PermissionName)
	return resourceAviatrixRbacGroupPermissionAttachmentRead(d, meta)
}

func resourceAviatrixRbacGroupPermissionAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	groupName := d.Get("group_name").(string)
	permissionName := d.Get("permission_name").(string)
	if groupName == "" || permissionName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no group name or permission name received. Import Id is %s", id)
		d.Set("group_name", strings.Split(id, "~")[0])
		d.Set("permission_name", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	attachment := &goaviatrix.RbacGroupPermissionAttachment{
		GroupName:      d.Get("group_name").(string),
		PermissionName: d.Get("permission_name").(string),
	}

	log.Printf("[INFO] Looking for Aviatrix RBAC group permission attachment: %#v", attachment)

	permissionAttachment, err := client.GetRbacGroupPermissionAttachment(attachment)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix RBAC group permission attachment: %s", err)
	}
	if permissionAttachment != nil {
		d.Set("group_name", permissionAttachment.GroupName)
		d.Set("permission_name", permissionAttachment.PermissionName)
		d.SetId(permissionAttachment.GroupName + "~" + permissionAttachment.PermissionName)
	}

	return nil
}

func resourceAviatrixRbacGroupPermissionAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	attachment := &goaviatrix.RbacGroupPermissionAttachment{
		GroupName:      d.Get("group_name").(string),
		PermissionName: d.Get("permission_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix RBAC group permission attachment: %#v", attachment)

	err := client.DeleteRbacGroupPermissionAttachment(attachment)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix RBAC group permission attachment: %s", err)
	}

	return nil
}
