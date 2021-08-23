package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixRbacGroupUserAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixRbacGroupUserAttachmentCreate,
		Read:   resourceAviatrixRbacGroupUserAttachmentRead,
		Delete: resourceAviatrixRbacGroupUserAttachmentDelete,
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
	client := meta.(*goaviatrix.Client)

	attachment := &goaviatrix.RbacGroupUserAttachment{
		GroupName: d.Get("group_name").(string),
		UserName:  d.Get("user_name").(string),
	}

	log.Printf("[INFO] Creating Aviatrix RBAC permission group user attachment: %#v", attachment)

	d.SetId(attachment.GroupName + "~" + attachment.UserName)
	flag := false
	defer resourceAviatrixRbacGroupUserAttachmentReadIfRequired(d, meta, &flag)

	err := client.CreateRbacGroupUserAttachment(attachment)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix RBAC permission group user attachment: %s", err)
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
	client := meta.(*goaviatrix.Client)

	groupName := d.Get("group_name").(string)
	userName := d.Get("user_name").(string)
	if groupName == "" || userName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no group name or account user name received. Import Id is %s", id)
		d.Set("group_name", strings.Split(id, "~")[0])
		d.Set("user_name", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	attachment := &goaviatrix.RbacGroupUserAttachment{
		GroupName: d.Get("group_name").(string),
		UserName:  d.Get("user_name").(string),
	}

	log.Printf("[INFO] Looking for Aviatrix RBAC permission group user attachment: %#v", attachment)

	userAttachment, err := client.GetRbacGroupUserAttachment(attachment)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix RBAC permission group user attachment: %s", err)
	}
	if userAttachment != nil {
		d.Set("group_name", userAttachment.GroupName)
		d.Set("user_name", userAttachment.UserName)
		d.SetId(userAttachment.GroupName + "~" + userAttachment.UserName)
	}

	return nil
}

func resourceAviatrixRbacGroupUserAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	attachment := &goaviatrix.RbacGroupUserAttachment{
		GroupName: d.Get("group_name").(string),
		UserName:  d.Get("user_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix RBAC permission group user attachment: %#v", attachment)

	err := client.DeleteRbacGroupUserAttachment(attachment)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix RBAC permission group user attachment: %s", err)
	}

	return nil
}
