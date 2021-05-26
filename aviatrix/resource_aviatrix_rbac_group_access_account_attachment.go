package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixRbacGroupAccessAccountAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixRbacGroupAccessAccountAttachmentCreate,
		Read:   resourceAviatrixRbacGroupAccessAccountAttachmentRead,
		Delete: resourceAviatrixRbacGroupAccessAccountAttachmentDelete,
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
	client := meta.(*goaviatrix.Client)

	attachment := &goaviatrix.RbacGroupAccessAccountAttachment{
		GroupName:         d.Get("group_name").(string),
		AccessAccountName: d.Get("access_account_name").(string),
	}

	log.Printf("[INFO] Creating Aviatrix RBAC permission group access account attachment: %#v", attachment)

	err := client.CreateRbacGroupAccessAccountAttachment(attachment)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix RBAC permission group access account attachment: %s", err)
	}

	log.Printf("[DEBUG] Aviatrix RBAC permission group access account attachment created")

	d.SetId(attachment.GroupName + "~" + attachment.AccessAccountName)
	return resourceAviatrixRbacGroupAccessAccountAttachmentRead(d, meta)
}

func resourceAviatrixRbacGroupAccessAccountAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	groupName := d.Get("group_name").(string)
	accessAccountName := d.Get("access_account_name").(string)
	if groupName == "" || accessAccountName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no group name or access account name received. Import Id is %s", id)
		d.Set("group_name", strings.Split(id, "~")[0])
		d.Set("access_account_name", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	attachment := &goaviatrix.RbacGroupAccessAccountAttachment{
		GroupName:         d.Get("group_name").(string),
		AccessAccountName: d.Get("access_account_name").(string),
	}

	log.Printf("[INFO] Looking for Aviatrix RBAC permission group access account attachment: %#v", attachment)

	accessAccountAttachment, err := client.GetRbacGroupAccessAccountAttachment(attachment)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix RBAC permission group access account attachment: %s", err)
	}
	if accessAccountAttachment != nil {
		d.Set("group_name", accessAccountAttachment.GroupName)
		d.Set("access_account_name", accessAccountAttachment.AccessAccountName)
		d.SetId(accessAccountAttachment.GroupName + "~" + accessAccountAttachment.AccessAccountName)
	}

	return nil
}

func resourceAviatrixRbacGroupAccessAccountAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	attachment := &goaviatrix.RbacGroupAccessAccountAttachment{
		GroupName:         d.Get("group_name").(string),
		AccessAccountName: d.Get("access_account_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix RBAC permission group access account attachment: %#v", attachment)

	err := client.DeleteRbacGroupAccessAccountAttachment(attachment)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix RBAC permission group access account attachment: %s", err)
	}

	return nil
}
