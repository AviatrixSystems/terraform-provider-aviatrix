package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixRbacGroupUserAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixRbacGroupUserAttachmentCreate,
		Read:   resourceAviatrixRbacGroupUserAttachmentRead,
		Delete: resourceAviatrixRbacGroupUserAttachmentDelete,
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
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Account user name.",
			},
		},
	}
}

func resourceAviatrixRbacGroupUserAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	attachment := &goaviatrix.RbacGroupUserAttachment{
		GroupName: getString(d, "group_name"),
		UserName:  getString(d, "user_name"),
	}

	log.Printf("[INFO] Creating Aviatrix RBAC permission group user attachment: %#v", attachment)

	d.SetId(attachment.GroupName + "~" + attachment.UserName)
	flag := false
	defer func() { _ = resourceAviatrixRbacGroupUserAttachmentReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	err := client.CreateRbacGroupUserAttachment(attachment)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix RBAC permission group user attachment: %w", err)
	}

	log.Printf("[DEBUG] Aviatrix RBAC permission group user attachment created")

	return resourceAviatrixRbacGroupUserAttachmentReadIfRequired(d, meta, &flag)
}

func resourceAviatrixRbacGroupUserAttachmentReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixRbacGroupUserAttachmentRead(d, meta)
	}
	return nil
}

func resourceAviatrixRbacGroupUserAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	groupName := getString(d, "group_name")
	userName := getString(d, "user_name")
	if groupName == "" || userName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no group name or account user name received. Import Id is %s", id)
		mustSet(d, "group_name", strings.Split(id, "~")[0])
		mustSet(d, "user_name", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	attachment := &goaviatrix.RbacGroupUserAttachment{
		GroupName: getString(d, "group_name"),
		UserName:  getString(d, "user_name"),
	}

	log.Printf("[INFO] Looking for Aviatrix RBAC permission group user attachment: %#v", attachment)

	userAttachment, err := client.GetRbacGroupUserAttachment(attachment)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix RBAC permission group user attachment: %w", err)
	}
	if userAttachment != nil {
		mustSet(d, "group_name", userAttachment.GroupName)
		mustSet(d, "user_name", userAttachment.UserName)
		d.SetId(userAttachment.GroupName + "~" + userAttachment.UserName)
	}

	return nil
}

func resourceAviatrixRbacGroupUserAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	attachment := &goaviatrix.RbacGroupUserAttachment{
		GroupName: getString(d, "group_name"),
		UserName:  getString(d, "user_name"),
	}

	log.Printf("[INFO] Deleting Aviatrix RBAC permission group user attachment: %#v", attachment)

	err := client.DeleteRbacGroupUserAttachment(attachment)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix RBAC permission group user attachment: %w", err)
	}

	return nil
}
