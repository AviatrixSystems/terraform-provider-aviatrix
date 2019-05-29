package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAwsTgwVpcAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsTgwVpcAttachmentCreate,
		Read:   resourceAwsTgwVpcAttachmentRead,
		Update: resourceAwsTgwVpcAttachmentUpdate,
		Delete: resourceAwsTgwVpcAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"tgw_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the AWS TGW.",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Region of cloud provider.",
			},
			"security_domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the security domain.",
			},
			"vpc_account_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This parameter represents the name of a Cloud-Account in Aviatrix controller.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This parameter represents the ID of the VPC.",
			},
		},
	}
}

func resourceAwsTgwVpcAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	awsTgwVpcAttachment := &goaviatrix.AwsTgwVpcAttachment{
		TgwName:            d.Get("tgw_name").(string),
		Region:             d.Get("region").(string),
		SecurityDomainName: d.Get("security_domain_name").(string),
		VpcAccountName:     d.Get("vpc_account_name").(string),
		VpcID:              d.Get("vpc_id").(string),
	}

	log.Printf("[INFO] Attaching vpc: %s to tgw %s", awsTgwVpcAttachment.VpcID, awsTgwVpcAttachment.TgwName)

	err := client.CreateAwsTgwVpcAttachment(awsTgwVpcAttachment)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Aws Tgw Vpc Attach: %s", err)
	}

	d.SetId(awsTgwVpcAttachment.TgwName + "~" + awsTgwVpcAttachment.SecurityDomainName + "~" + awsTgwVpcAttachment.VpcID)

	return resourceAwsTgwVpcAttachmentRead(d, meta)
}

func resourceAwsTgwVpcAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	tgwName := d.Get("tgw_name").(string)
	securityDomainName := d.Get("security_domain_name").(string)
	vpcID := d.Get("vpc_id").(string)

	if tgwName == "" || securityDomainName == "" || vpcID == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no vpc names received. Import Id is %s", id)
		d.Set("tgw_name", strings.Split(id, "~")[0])
		d.Set("security_domain_name", strings.Split(id, "~")[1])
		d.Set("vpc_id", strings.Split(id, "~")[2])
	}
	awsTgwVpcAttachment := &goaviatrix.AwsTgwVpcAttachment{
		TgwName:            d.Get("tgw_name").(string),
		SecurityDomainName: d.Get("security_domain_name").(string),
		VpcID:              d.Get("vpc_id").(string),
	}

	aTVA, err := client.GetAwsTgwVpcAttachment(awsTgwVpcAttachment)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to get Aviatrix Aws Tgw Vpc Attach: %s", err)
	}
	if aTVA != nil {
		d.Set("tgw_name", awsTgwVpcAttachment.TgwName)
		d.Set("region", awsTgwVpcAttachment.Region)
		d.Set("security_domain_name", awsTgwVpcAttachment.SecurityDomainName)
		d.Set("vpc_account_name", awsTgwVpcAttachment.VpcAccountName)
		d.Set("vpc_id", awsTgwVpcAttachment.VpcID)
		d.SetId(awsTgwVpcAttachment.TgwName + "~" + awsTgwVpcAttachment.SecurityDomainName + "~" + awsTgwVpcAttachment.VpcID)
		return nil
	}

	return fmt.Errorf("no Aviatrix Aws Tgw Vpc Attach found")
}

func resourceAwsTgwVpcAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	if d.HasChange("tgw_name") {
		return fmt.Errorf("updating tgw_name is not allowed")
	}
	if d.HasChange("region") {
		return fmt.Errorf("updating region is not allowed")
	}
	if d.HasChange("vpc_account_name") {
		return fmt.Errorf("updating vpc_account_name is not allowed")
	}
	if d.HasChange("vpc_id") {
		return fmt.Errorf("updating vpc_id is not allowed")
	}
	d.Partial(false)

	return resourceAwsTgwVpcAttachmentRead(d, meta)
}

func resourceAwsTgwVpcAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	awsTgwVpcAttachment := &goaviatrix.AwsTgwVpcAttachment{
		TgwName:            d.Get("tgw_name").(string),
		Region:             d.Get("region").(string),
		SecurityDomainName: d.Get("security_domain_name").(string),
		VpcAccountName:     d.Get("vpc_account_name").(string),
		VpcID:              d.Get("vpc_id").(string),
	}

	log.Printf("[INFO] Detaching vpc: %s from tgw %s", awsTgwVpcAttachment.VpcID, awsTgwVpcAttachment.TgwName)

	err := client.DeleteAwsTgwVpcAttachment(awsTgwVpcAttachment)
	if err != nil {
		return fmt.Errorf("failed to detach vpc from tgw: %s", err)
	}

	return nil
}
