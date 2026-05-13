package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixRbacGroupAccessAccountAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixRbacGroupAccessAccountAttachmentCreate,
		Read:   resourceAviatrixRbacGroupAccessAccountAttachmentRead,
		Delete: resourceAviatrixRbacGroupAccessAccountAttachmentDelete,
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
			"access_account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Access account name.",
			},
		},
	}
}

func resourceAviatrixRbacGroupAccessAccountAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	attachment := &goaviatrix.RbacGroupAccessAccountAttachment{
		GroupName:         getString(d, "group_name"),
		AccessAccountName: getString(d, "access_account_name"),
	}

	log.Printf("[INFO] Creating Aviatrix RBAC permission group access account attachment: %#v", attachment)

	d.SetId(attachment.GroupName + "~" + attachment.AccessAccountName)
	flag := false
	defer func() { _ = resourceAviatrixRbacGroupAccessAccountAttachmentReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	err := client.CreateRbacGroupAccessAccountAttachment(attachment)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix RBAC permission group access account attachment: %w", err)
	}

	log.Printf("[DEBUG] Aviatrix RBAC permission group access account attachment created")

	return resourceAviatrixRbacGroupAccessAccountAttachmentReadIfRequired(d, meta, &flag)
}

func resourceAviatrixRbacGroupAccessAccountAttachmentReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixRbacGroupAccessAccountAttachmentRead(d, meta)
	}
	return nil
}

func resourceAviatrixRbacGroupAccessAccountAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	groupName := getString(d, "group_name")
	accessAccountName := getString(d, "access_account_name")
	if groupName == "" || accessAccountName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no group name or access account name received. Import Id is %s", id)
		mustSet(d, "group_name", strings.Split(id, "~")[0])
		mustSet(d, "access_account_name", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	attachment := &goaviatrix.RbacGroupAccessAccountAttachment{
		GroupName:         getString(d, "group_name"),
		AccessAccountName: getString(d, "access_account_name"),
	}

	log.Printf("[INFO] Looking for Aviatrix RBAC permission group access account attachment: %#v", attachment)

	accessAccountAttachment, err := client.GetRbacGroupAccessAccountAttachment(attachment)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix RBAC permission group access account attachment: %w", err)
	}
	if accessAccountAttachment != nil {
		mustSet(d, "group_name", accessAccountAttachment.GroupName)
		mustSet(d, "access_account_name", accessAccountAttachment.AccessAccountName)
		d.SetId(accessAccountAttachment.GroupName + "~" + accessAccountAttachment.AccessAccountName)
	}

	return nil
}

func resourceAviatrixRbacGroupAccessAccountAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	attachment := &goaviatrix.RbacGroupAccessAccountAttachment{
		GroupName:         getString(d, "group_name"),
		AccessAccountName: getString(d, "access_account_name"),
	}

	log.Printf("[INFO] Deleting Aviatrix RBAC permission group access account attachment: %#v", attachment)

	err := client.DeleteRbacGroupAccessAccountAttachment(attachment)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix RBAC permission group access account attachment: %w", err)
	}

	return nil
}
