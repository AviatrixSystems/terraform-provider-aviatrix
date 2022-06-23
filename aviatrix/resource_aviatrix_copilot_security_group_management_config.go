package aviatrix

import (
	"context"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixCopilotSecurityGroupManagementConfig() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixCopilotSecurityGroupManagementConfigCreate,
		ReadWithoutTimeout:   resourceAviatrixCopilotSecurityGroupManagementConfigRead,
		UpdateWithoutTimeout: resourceAviatrixCopilotSecurityGroupManagementConfigUpdate,
		DeleteWithoutTimeout: resourceAviatrixCopilotSecurityGroupManagementConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"enable_copilot_security_group_management": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Switch to enable copilot security group management.",
			},
			"cloud_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Cloud type.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Access account name.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "VPC ID.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Copilot instance ID.",
			},
			"region": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Copilot region. Valid for AWS and Azure.",
				ConflictsWith: []string{"zone"},
			},
			"zone": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Copilot zone. Valid for GCP.",
				ConflictsWith: []string{"region"},
			},
		},
	}
}

func marshalCopilotSecurityGroupManagementConfigInput(d *schema.ResourceData) *goaviatrix.CopilotSecurityGroupManagementConfig {
	copilotSecurityGroupManagementConfig := &goaviatrix.CopilotSecurityGroupManagementConfig{
		CloudType:                            d.Get("cloud_type").(int),
		AccountName:                          d.Get("account_name").(string),
		Region:                               d.Get("region").(string),
		Zone:                                 d.Get("zone").(string),
		VpcId:                                d.Get("vpc_id").(string),
		InstanceId:                           d.Get("instance_id").(string),
		EnableCopilotSecurityGroupManagement: d.Get("enable_copilot_security_group_management").(bool),
	}

	return copilotSecurityGroupManagementConfig
}

func resourceAviatrixCopilotSecurityGroupManagementConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	copilotSecurityGroupManagementConfig := marshalCopilotSecurityGroupManagementConfigInput(d)

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	flag := false
	defer resourceAviatrixCopilotSecurityGroupManagementConfigReadIfRequired(ctx, d, meta, &flag)

	if copilotSecurityGroupManagementConfig.EnableCopilotSecurityGroupManagement {
		if copilotSecurityGroupManagementConfig.CloudType == 0 || copilotSecurityGroupManagementConfig.AccountName == "" ||
			copilotSecurityGroupManagementConfig.VpcId == "" || copilotSecurityGroupManagementConfig.InstanceId == "" ||
			(copilotSecurityGroupManagementConfig.Region == "" && copilotSecurityGroupManagementConfig.Zone == "") {
			return diag.Errorf("'cloud_type', 'account_name', 'region'/'zone', 'vpc_id' and 'instance_id' are required to enable copilot security group management")
		}

		if goaviatrix.IsCloudType(copilotSecurityGroupManagementConfig.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) &&
			copilotSecurityGroupManagementConfig.Zone != "" {
			return diag.Errorf("'zone' is only supported for 'GCP', please use 'region' for AWS and Azure")
		}

		if goaviatrix.IsCloudType(copilotSecurityGroupManagementConfig.CloudType, goaviatrix.GCPRelatedCloudTypes) &&
			copilotSecurityGroupManagementConfig.Region != "" {
			return diag.Errorf("'region' is only supported for AWS and Azure, please use 'zone' for GCP")
		}

		if err := client.EnableCopilotSecurityGroupManagement(ctx, copilotSecurityGroupManagementConfig); err != nil {
			return diag.Errorf("could not enable copilot security group management: %v", err)
		}
	} else {
		if copilotSecurityGroupManagementConfig.CloudType != 0 || copilotSecurityGroupManagementConfig.AccountName != "" ||
			copilotSecurityGroupManagementConfig.VpcId != "" || copilotSecurityGroupManagementConfig.InstanceId != "" ||
			copilotSecurityGroupManagementConfig.Region != "" || copilotSecurityGroupManagementConfig.Zone != "" {
			return diag.Errorf("'cloud_type', 'account_name', 'region'/'zone', 'vpc_id' and 'instance_id' are not needed to disable copilot security group management")
		}

		err := client.DisableCopilotSecurityGroupManagement(ctx)
		if err != nil {
			return diag.Errorf("could not disable copilot security group management: %v", err)
		}
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

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return diag.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	copilotSecurityGroupManagementConfig, err := client.GetCopilotSecurityGroupManagementConfig(ctx)
	if err != nil {
		return diag.Errorf("could not read copilot security group management config: %v", err)
	}

	if copilotSecurityGroupManagementConfig != nil {
		d.Set("enable_copilot_security_group_management", copilotSecurityGroupManagementConfig.State == "Enabled")
		d.Set("cloud_type", copilotSecurityGroupManagementConfig.CloudType)
		d.Set("account_name", copilotSecurityGroupManagementConfig.AccountName)
		d.Set("vpc_id", copilotSecurityGroupManagementConfig.VpcId)
		d.Set("instance_id", copilotSecurityGroupManagementConfig.InstanceIdReturn)

		if goaviatrix.IsCloudType(copilotSecurityGroupManagementConfig.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) {
			d.Set("region", copilotSecurityGroupManagementConfig.Region)
		} else if goaviatrix.IsCloudType(copilotSecurityGroupManagementConfig.CloudType, goaviatrix.GCPRelatedCloudTypes) {
			d.Set("zone", copilotSecurityGroupManagementConfig.Zone)
		}
	} else {
		return diag.Errorf("could not read copilot security group management config")
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixCopilotSecurityGroupManagementConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	copilotSecurityGroupManagementConfig := marshalCopilotSecurityGroupManagementConfigInput(d)

	if d.HasChange("enable_copilot_security_group_management") {
		if copilotSecurityGroupManagementConfig.EnableCopilotSecurityGroupManagement {
			if copilotSecurityGroupManagementConfig.CloudType == 0 || copilotSecurityGroupManagementConfig.AccountName == "" ||
				copilotSecurityGroupManagementConfig.VpcId == "" || copilotSecurityGroupManagementConfig.InstanceId == "" ||
				(copilotSecurityGroupManagementConfig.Region == "" && copilotSecurityGroupManagementConfig.Zone == "") {
				return diag.Errorf("'cloud_type', 'account_name', 'region'/'zone', 'vpc_id' and 'instance_id' are required to enable copilot security group management")
			}

			if goaviatrix.IsCloudType(copilotSecurityGroupManagementConfig.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) &&
				copilotSecurityGroupManagementConfig.Zone != "" {
				return diag.Errorf("'zone' is only supported for 'GCP', please use 'region' for AWS and Azure")
			}

			if goaviatrix.IsCloudType(copilotSecurityGroupManagementConfig.CloudType, goaviatrix.GCPRelatedCloudTypes) &&
				copilotSecurityGroupManagementConfig.Region != "" {
				return diag.Errorf("'region' is only supported for AWS and Azure, please use 'zone' for GCP")
			}

			if err := client.EnableCopilotSecurityGroupManagement(ctx, copilotSecurityGroupManagementConfig); err != nil {
				return diag.Errorf("could not enable copilot security group management: %v", err)
			}
		} else {
			if copilotSecurityGroupManagementConfig.CloudType != 0 || copilotSecurityGroupManagementConfig.AccountName != "" ||
				copilotSecurityGroupManagementConfig.VpcId != "" || copilotSecurityGroupManagementConfig.InstanceId != "" ||
				copilotSecurityGroupManagementConfig.Region != "" || copilotSecurityGroupManagementConfig.Zone != "" {
				return diag.Errorf("'cloud_type', 'account_name', 'region'/'zone', 'vpc_id' and 'instance_id' are not needed to disable copilot security group management")
			}

			err := client.DisableCopilotSecurityGroupManagement(ctx)
			if err != nil {
				return diag.Errorf("could not disable copilot security group management: %v", err)
			}
		}
	} else {
		if copilotSecurityGroupManagementConfig.EnableCopilotSecurityGroupManagement {
			if d.HasChanges("cloud_type", "account_name", "region", "zone", "vpc_id", "instance_id") {
				err := client.DisableCopilotSecurityGroupManagement(ctx)
				if err != nil {
					return diag.Errorf("could not disable copilot security group management: %v", err)
				}

				if copilotSecurityGroupManagementConfig.CloudType == 0 || copilotSecurityGroupManagementConfig.AccountName == "" ||
					copilotSecurityGroupManagementConfig.VpcId == "" || copilotSecurityGroupManagementConfig.InstanceId == "" ||
					(copilotSecurityGroupManagementConfig.Region == "" && copilotSecurityGroupManagementConfig.Zone == "") {
					return diag.Errorf("'cloud_type', 'account_name', 'region'/'zone', 'vpc_id' and 'instance_id' are required to enable copilot security group management")
				}

				if goaviatrix.IsCloudType(copilotSecurityGroupManagementConfig.CloudType, goaviatrix.AWSRelatedCloudTypes|goaviatrix.AzureArmRelatedCloudTypes) &&
					copilotSecurityGroupManagementConfig.Zone != "" {
					return diag.Errorf("'zone' is only supported for 'GCP', please use 'region' for AWS and Azure")
				}

				if goaviatrix.IsCloudType(copilotSecurityGroupManagementConfig.CloudType, goaviatrix.GCPRelatedCloudTypes) &&
					copilotSecurityGroupManagementConfig.Region != "" {
					return diag.Errorf("'region' is only supported for AWS and Azure, please use 'zone' for GCP")
				}

				if err := client.EnableCopilotSecurityGroupManagement(ctx, copilotSecurityGroupManagementConfig); err != nil {
					return diag.Errorf("could not enable copilot security group management: %v", err)
				}
			}
		} else {
			if d.HasChanges("cloud_type", "account_name", "region", "zone", "vpc_id", "instance_id") {
				return diag.Errorf("'cloud_type', 'account_name', 'region'/'zone', 'vpc_id' and 'instance_id' are not allowed to be changed when_copilot security group management is disabled")
			}
		}
	}

	return resourceAviatrixCopilotSecurityGroupManagementConfigRead(ctx, d, meta)
}

func resourceAviatrixCopilotSecurityGroupManagementConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	err := client.DisableCopilotSecurityGroupManagement(ctx)
	if err != nil {
		return diag.Errorf("could not disable copilot security group management: %v", err)
	}

	return nil
}
