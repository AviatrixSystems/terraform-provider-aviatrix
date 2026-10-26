package aviatrix

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixCopilotSimpleDeployment() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixCopilotSimpleDeploymentCreate,
		ReadWithoutTimeout:   resourceAviatrixCopilotSimpleDeploymentRead,
		DeleteWithoutTimeout: resourceAviatrixCopilotSimpleDeploymentDelete,
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
				Description: "Aviatrix access account name.",
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
			"subnet": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Subnet CIDR.",
			},
			"controller_service_account_username": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Controller service account username.",
			},
			"controller_service_account_password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "Controller service account password.",
			},
			"instance_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "t3.2xlarge",
				ForceNew:    true,
				Description: "Instance size.",
			},
			"data_volume_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				ForceNew:    true,
				Description: "Data volume size.",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Copilot public IP.",
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Copilot private IP.",
			},
		},
	}
}

func marshalCopilotSimpleDeploymentInput(d *schema.ResourceData) *goaviatrix.CopilotSimpleDeployment {
	copilotSimpleDeployment := &goaviatrix.CopilotSimpleDeployment{
		CloudType:                        getInt(d, "cloud_type"),
		AccountName:                      getString(d, "account_name"),
		Region:                           getString(d, "region"),
		VpcId:                            getString(d, "vpc_id"),
		Subnet:                           getString(d, "subnet"),
		ControllerServiceAccountUsername: getString(d, "controller_service_account_username"),
		ControllerServiceAccountPassword: getString(d, "controller_service_account_password"),
		InstanceSize:                     getString(d, "instance_size"),
		DataVolumeSize:                   getInt(d, "data_volume_size"),
	}

	return copilotSimpleDeployment
}

func resourceAviatrixCopilotSimpleDeploymentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	copilotSimpleDeployment := marshalCopilotSimpleDeploymentInput(d)

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	flag := false
	defer resourceAviatrixCopilotSimpleDeploymentReadIfRequired(ctx, d, meta, &flag)

	if err := client.CreateCopilotSimple(ctx, copilotSimpleDeployment); err != nil {
		return diag.Errorf("could not start to deploy copilot: %v", err)
	}

	for i := 0; ; i++ {
		copilotAssociationStatus, err := client.GetCopilotAssociationStatus(ctx)

		if err != nil && !errors.Is(err, goaviatrix.ErrNotFound) {
			return diag.Errorf("could not get copilot association status: %v", err)
		}

		if !errors.Is(err, goaviatrix.ErrNotFound) && copilotAssociationStatus.Status {
			break
		}

		if i < 90 {
			time.Sleep(time.Duration(20) * time.Second)
		} else {
			return diag.Errorf("could not deploy copilot: %s", err)
		}
	}

	return resourceAviatrixCopilotSimpleDeploymentReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixCopilotSimpleDeploymentReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixCopilotSimpleDeploymentRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixCopilotSimpleDeploymentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	copilotAssociationStatus, err := client.GetCopilotAssociationStatus(ctx)
	if errors.Is(err, goaviatrix.ErrNotFound) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("could not get copilot association status: %v", err)
	}
	mustSet(d, "private_ip", copilotAssociationStatus.IP)
	mustSet(d, "public_ip", copilotAssociationStatus.PublicIp)

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixCopilotSimpleDeploymentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	err := client.DeleteCopilotSimple(ctx)
	if err != nil {
		return diag.Errorf("could not delete copilot: %v", err)
	}

	return nil
}
