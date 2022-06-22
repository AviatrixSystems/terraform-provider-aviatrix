package aviatrix

import (
	"context"
	"log"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixCopilotSecurityGroupManagementConfig() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixCopilotSecurityGroupManagementConfigCreate,
		ReadWithoutTimeout:   resourceAviatrixCopilotSecurityGroupManagementConfigRead,
		DeleteWithoutTimeout: resourceAviatrixCopilotSecurityGroupManagementConfigDelete,
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
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Copilot region. Valid for AWS and Azure.",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Copilot zone. Valid for GCP.",
			},
		},
	}
}

func marshalCopilotSecurityGroupManagementConfigInput(d *schema.ResourceData) *goaviatrix.CopilotSecurityGroupManagementConfig {
	copilotSecurityGroupManagementConfig := &goaviatrix.CopilotSecurityGroupManagementConfig{
		CloudType:   d.Get("cloud_type").(int),
		AccountName: d.Get("account_name").(string),
		Region:      d.Get("region").(string),
		Zone:        d.Get("zone").(string),
		VpcId:       d.Get("vpc_id").(string),
		InstanceID:  d.Get("instance_id").(string),
	}

	return copilotSecurityGroupManagementConfig
}

func resourceAviatrixCopilotSecurityGroupManagementConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	copilotSecurityGroupManagementConfig := marshalCopilotSecurityGroupManagementConfigInput(d)

	if goaviatrix.IsCloudType(copilotSecurityGroupManagementConfig.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) &&
		(copilotSecurityGroupManagementConfig.Region == "" || copilotSecurityGroupManagementConfig.Zone != "") {
		return diag.Errorf("'zone' is only supported for 'GCP', please use 'region' for AWS and Azure")
	}

	if goaviatrix.IsCloudType(copilotSecurityGroupManagementConfig.CloudType, goaviatrix.GCPRelatedCloudTypes) &&
		(copilotSecurityGroupManagementConfig.Region != "" || copilotSecurityGroupManagementConfig.Zone == "") {
		return diag.Errorf("'region' is only supported for AWS and Azure, please use 'zone' for GCP")
	}

	d.SetId(copilotSecurityGroupManagementConfig.InstanceID)
	flag := false
	defer resourceAviatrixCopilotSecurityGroupManagementConfigReadIfRequired(ctx, d, meta, &flag)

	if err := client.EnableCopilotSecurityGroupManagement(ctx, copilotSecurityGroupManagementConfig); err != nil {
		return diag.Errorf("could not enable copilot security group management: %v", err)
	}

	return resourceAviatrixCopilotSecurityGroupManagementConfigReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixCopilotSecurityGroupManagementConfigReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixCopilotSecurityGroupManagementConfigRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixCopilotSecurityGroupManagementConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Get("instance_id").(string) == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no instance id received. Import Id is %s", id)
		d.Set("instance_id", id)
		d.SetId(id)
	}

	copilotSecurityGroupManagementConfig, err := client.GetCopilotSecurityGroupManagementConfig(ctx)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read copilot security group management config: %v", err)
	}

	d.Set("cloud_type", copilotSecurityGroupManagementConfig.CloudType)
	d.Set("account_name", copilotSecurityGroupManagementConfig.AccountName)
	d.Set("vpc_id", copilotSecurityGroupManagementConfig.VpcId)
	d.Set("instance_id", copilotSecurityGroupManagementConfig.InstanceIDReturn)

	if goaviatrix.IsCloudType(copilotSecurityGroupManagementConfig.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
		d.Set("region", copilotSecurityGroupManagementConfig.Region)
	} else if goaviatrix.IsCloudType(copilotSecurityGroupManagementConfig.CloudType, goaviatrix.GCPRelatedCloudTypes) {
		d.Set("zone", copilotSecurityGroupManagementConfig.Zone)
	}

	d.SetId(copilotSecurityGroupManagementConfig.InstanceIDReturn)
	return nil
}

func resourceAviatrixCopilotSecurityGroupManagementConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	err := client.DisableCopilotSecurityGroupManagement(ctx)
	if err != nil {
		return diag.Errorf("could not disable copilot security group management: %v", err)
	}

	return nil
}
