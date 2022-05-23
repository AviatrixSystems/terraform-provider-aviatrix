package aviatrix

import (
	"context"
	"log"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixCopilotSecurityGroupManagement() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixCopilotSecurityGroupManagementCreate,
		ReadWithoutTimeout:   resourceAviatrixCopilotSecurityGroupManagementRead,
		DeleteWithoutTimeout: resourceAviatrixCopilotSecurityGroupManagementDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"cloud_type": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Cloud type.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Access account name.",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC ID.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Copilot instance ID.",
			},
		},
	}
}

func marshalCopilotSecurityGroupManagementInput(d *schema.ResourceData) *goaviatrix.CopilotSecurityGroupManagement {
	copilotSecurityGroupManagement := &goaviatrix.CopilotSecurityGroupManagement{
		//CloudType: strconv.Itoa(d.Get("cloud_type").(int)),
		CloudType:   d.Get("cloud_type").(int),
		AccountName: d.Get("account_name").(string),
		Region:      d.Get("region").(string),
		VpcId:       d.Get("vpc_id").(string),
		InstanceID:  d.Get("instance_id").(string),
	}

	return copilotSecurityGroupManagement
}

func resourceAviatrixCopilotSecurityGroupManagementCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	copilotSecurityGroupManagement := marshalCopilotSecurityGroupManagementInput(d)

	d.SetId(copilotSecurityGroupManagement.InstanceID)
	flag := false
	defer resourceAviatrixCopilotSecurityGroupManagementReadIfRequired(ctx, d, meta, &flag)

	if err := client.EnableCopilotSecurityGroupManagement(ctx, copilotSecurityGroupManagement); err != nil {
		return diag.Errorf("could not enable copilot security group management: %v", err)
	}

	return resourceAviatrixCopilotSecurityGroupManagementReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixCopilotSecurityGroupManagementReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixCopilotSecurityGroupManagementRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixCopilotSecurityGroupManagementRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Get("instance_id").(string) == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no name received. Import Id is %s", id)
		d.Set("instance_id", id)
		d.SetId(id)
	}

	copilotSecurityGroupManagement, err := client.GetCopilotSecurityGroupManagement(ctx)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read copilot security group management: %v", err)
	}

	d.Set("cloud_type", copilotSecurityGroupManagement.CloudType)
	d.Set("account_name", copilotSecurityGroupManagement.AccountName)
	d.Set("region", copilotSecurityGroupManagement.Region)
	d.Set("vpc_id", copilotSecurityGroupManagement.VpcId)
	d.Set("instance_id", copilotSecurityGroupManagement.InstanceID)

	d.SetId(copilotSecurityGroupManagement.InstanceID)
	return nil
}

func resourceAviatrixCopilotSecurityGroupManagementDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	err := client.DisableCopilotSecurityGroupManagement(ctx)
	if err != nil {
		return diag.Errorf("could not disable copilot security group management: %v", err)
	}

	return nil
}
