package aviatrix

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixBranchRouterVirtualWanAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixBranchRouterVirtualWanAttachmentCreate,
		Read:   resourceAviatrixBranchRouterVirtualWanAttachmentRead,
		Delete: resourceAviatrixBranchRouterVirtualWanAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"connection_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Connection name.",
			},
			"branch_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Branch router name.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Azure access account name.",
			},
			"resource_group": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ARM resource group name.",
			},
			"hub_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Virtual WAN vhub name.",
			},
			"branch_router_bgp_asn": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Branch router AS Number.",
			},
		},
	}
}

func marshalBranchRouterVirtualWanAttachmentInput(d *schema.ResourceData) *goaviatrix.BranchRouterVirtualWanAttachment {
	return &goaviatrix.BranchRouterVirtualWanAttachment{
		ConnectionName:  d.Get("connection_name").(string),
		BranchName:      d.Get("branch_name").(string),
		AccountName:     d.Get("account_name").(string),
		ResourceGroup:   d.Get("resource_group").(string),
		HubName:         d.Get("hub_name").(string),
		BranchRouterAsn: strconv.Itoa(d.Get("branch_router_bgp_asn").(int)),
	}
}

func resourceAviatrixBranchRouterVirtualWanAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	attachment := marshalBranchRouterVirtualWanAttachmentInput(d)

	if err := client.CreateBranchRouterVirtualWanAttachment(attachment); err != nil {
		return err
	}

	d.SetId(attachment.ConnectionName)
	return nil
}

func resourceAviatrixBranchRouterVirtualWanAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	connectionName := d.Get("connection_name").(string)
	if connectionName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no branch_router_virtual_wan_attachment connectionName received. Import Id is %s", id)
		d.SetId(id)
		connectionName = id
	}

	attachment := &goaviatrix.BranchRouterVirtualWanAttachment{
		ConnectionName: connectionName,
	}

	attachment, err := client.GetBranchRouterVirtualWanAttachment(attachment)
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not find branch_router_virtual_wan_attachment %s: %v", connectionName, err)
	}

	d.Set("connection_name", attachment.ConnectionName)
	d.Set("branch_name", attachment.BranchName)
	d.Set("account_name", attachment.AccountName)
	d.Set("resource_group", attachment.ResourceGroup)
	d.Set("hub_name", attachment.HubName)

	branchRouterAsn, err := strconv.Atoi(attachment.BranchRouterAsn)
	if err != nil {
		return fmt.Errorf("could not covert branch router asn to int: %v", err)
	}
	d.Set("branch_router_bgp_asn", branchRouterAsn)

	d.SetId(attachment.ConnectionName)
	return nil
}

func resourceAviatrixBranchRouterVirtualWanAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	attachment := marshalBranchRouterVirtualWanAttachmentInput(d)

	if err := client.DeleteBranchRouterAttachment(attachment.ConnectionName); err != nil {
		return err
	}

	return nil
}
