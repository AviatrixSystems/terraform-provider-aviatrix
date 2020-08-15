package aviatrix

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixBranchRouterAwsTgwAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixBranchRouterAwsTgwAttachmentCreate,
		Read:   resourceAviatrixBranchRouterAwsTgwAttachmentRead,
		Delete: resourceAviatrixBranchRouterAwsTgwAttachmentDelete,
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
				Description: "Branch name.",
			},
			"aws_tgw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "AWS TGW name.",
			},
			"branch_router_bgp_asn": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Branch router BGP AS Number.",
			},
			"security_domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Security domain name.",
			},
		},
	}
}

func marshalBranchRouterAwsTgwAttachmentInput(d *schema.ResourceData) *goaviatrix.BranchRouterAwsTgwAttachment {
	return &goaviatrix.BranchRouterAwsTgwAttachment{
		ConnectionName:     d.Get("connection_name").(string),
		BranchName:         d.Get("branch_name").(string),
		AwsTgwName:         d.Get("aws_tgw_name").(string),
		BranchRouterAsn:    strconv.Itoa(d.Get("branch_router_bgp_asn").(int)),
		SecurityDomainName: d.Get("security_domain_name").(string),
	}
}

func resourceAviatrixBranchRouterAwsTgwAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	brata := marshalBranchRouterAwsTgwAttachmentInput(d)

	if err := client.CreateBranchRouterAwsTgwAttachment(brata); err != nil {
		return err
	}

	d.SetId(brata.ID())
	return nil
}

func resourceAviatrixBranchRouterAwsTgwAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	connectionName := d.Get("connection_name").(string)
	branchName := d.Get("branch_name").(string)
	tgwName := d.Get("aws_tgw_name").(string)
	if connectionName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no branch_router_aws_tgw_attachment connection_name received. Import Id is %s", id)
		d.SetId(id)
		parts := strings.Split(id, "~")
		if len(parts) != 3 {
			return fmt.Errorf("import id is invalid, expecting connection_name~branch_name~aws_tgw_name: %s", id)
		}
		connectionName = parts[0]
		branchName = parts[1]
		tgwName = parts[2]
	}

	brata := &goaviatrix.BranchRouterAwsTgwAttachment{
		ConnectionName: connectionName,
		BranchName:     branchName,
		AwsTgwName:     tgwName,
	}

	brata, err := client.GetBranchRouterAwsTgwAttachment(brata)
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not find branch_router_aws_tgw_attachment %s: %v", connectionName, err)
	}

	d.Set("connection_name", brata.ConnectionName)
	d.Set("branch_name", brata.BranchName)
	d.Set("aws_tgw_name", brata.AwsTgwName)
	d.Set("security_domain_name", brata.SecurityDomainName)

	asn, err := strconv.Atoi(brata.BranchRouterAsn)
	if err != nil {
		return fmt.Errorf("could not convert BranchRouterAsn to int: %v", err)
	}
	d.Set("branch_router_bgp_asn", asn)

	d.SetId(brata.ID())
	return nil
}

func resourceAviatrixBranchRouterAwsTgwAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	cn := d.Get("connection_name").(string)

	if err := client.DeleteDeviceAttachment(cn); err != nil {
		return err
	}

	return nil
}
