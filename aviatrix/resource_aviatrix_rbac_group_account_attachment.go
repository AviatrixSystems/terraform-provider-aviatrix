package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixRbacGroupAccountAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixRbacGroupAccountAttachmentCreate,
		Read:   resourceAviatrixRbacGroupAccountAttachmentRead,
		Delete: resourceAviatrixRbacGroupAccountAttachmentDelete,
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
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Account name.",
			},
		},
	}
}

func resourceAviatrixRbacGroupAccountAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	attachment := &goaviatrix.RbacGroupAccountAttachment{
		GroupName:   d.Get("group_name").(string),
		AccountName: d.Get("account_name").(string),
	}

	log.Printf("[INFO] Creating Aviatrix RBAC permission group account attachment: %#v", attachment)

	err := client.CreateRbacGroupAccountAttachment(attachment)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix RBAC permission group account attachment: %s", err)
	}

	log.Printf("[DEBUG] Aviatrix RBAC permission group account attachment created")

	d.SetId(attachment.GroupName + "~" + attachment.AccountName)
	return resourceAviatrixRbacGroupAccountAttachmentRead(d, meta)
}

func resourceAviatrixRbacGroupAccountAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	groupName := d.Get("group_name").(string)
	accountName := d.Get("account_name").(string)
	if groupName == "" || accountName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no group name or account name received. Import Id is %s", id)
		d.Set("group_name", strings.Split(id, "~")[0])
		d.Set("account_name", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	attachment := &goaviatrix.RbacGroupAccountAttachment{
		GroupName:   d.Get("group_name").(string),
		AccountName: d.Get("account_name").(string),
	}

	log.Printf("[INFO] Looking for Aviatrix RBAC permission group account attachment: %#v", attachment)

	accountAttachment, err := client.GetRbacGroupAccountAttachment(attachment)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix RBAC permission group account attachment: %s", err)
	}
	if accountAttachment != nil {
		d.Set("group_name", accountAttachment.GroupName)
		d.Set("account_name", accountAttachment.AccountName)
		d.SetId(accountAttachment.GroupName + "~" + accountAttachment.AccountName)
	}

	return nil
}

func resourceAviatrixRbacGroupAccountAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	attachment := &goaviatrix.RbacGroupAccountAttachment{
		GroupName:   d.Get("group_name").(string),
		AccountName: d.Get("account_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix RBAC permission group account attachment: %#v", attachment)

	err := client.DeleteRbacGroupAccountAttachment(attachment)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix RBAC permission group account attachment: %s", err)
	}

	return nil
}
