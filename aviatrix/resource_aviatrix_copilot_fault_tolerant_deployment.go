package aviatrix

import (
	"context"
	"strings"
	"time"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixCopilotFaultTolerantDeployment() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixCopilotFaultTolerantDeploymentCreate,
		ReadWithoutTimeout:   resourceAviatrixCopilotFaultTolerantDeploymentRead,
		DeleteWithoutTimeout: resourceAviatrixCopilotFaultTolerantDeploymentDelete,
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
			"main_copilot_vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC ID.",
			},
			"main_copilot_subnet": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Subnet CIDR.",
			},
			"main_copilot_instance_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "t3.2xlarge",
				ForceNew:    true,
				Description: "Instance size.",
			},
			"cluster_data_nodes": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster data nodes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "VPC ID.",
						},
						"subnet": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Subnet CIDR.",
						},
						"instance_size": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "t3.2xlarge",
							Description: "Instance size.",
						},
						"data_volume_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     100,
							Description: "Data volume size.",
						},
					},
				},
			},
			"main_copilot_public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Copilot public IP.",
			},
			"main_copilot_private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Copilot private IP.",
			},
		},
	}
}

func marshalCopilotFaultTolerantDeploymentInput(d *schema.ResourceData) *goaviatrix.CopilotFaultTolerantDeployment {
	copilotFaultTolerantDeployment := &goaviatrix.CopilotFaultTolerantDeployment{
		CloudType:                        d.Get("cloud_type").(int),
		AccountName:                      d.Get("account_name").(string),
		Region:                           d.Get("region").(string),
		ControllerServiceAccountUsername: d.Get("controller_service_account_username").(string),
		ControllerServiceAccountPassword: d.Get("controller_service_account_password").(string),
	}

	mainCopilot := &goaviatrix.MainCopilot{
		VpcId:        d.Get("main_copilot_vpc_id").(string),
		Subnet:       d.Get("main_copilot_subnet").(string),
		InstanceSize: d.Get("main_copilot_instance_size").(string),
	}
	copilotFaultTolerantDeployment.MainCopilot = mainCopilot

	clusterDataNode := d.Get("cluster_data_nodes").([]interface{})
	for _, clusterDataNode0 := range clusterDataNode {
		clusterDataNode1 := clusterDataNode0.(map[string]interface{})

		clusterDataNode2 := &goaviatrix.ClusterDataNode{
			VpcId:          clusterDataNode1["vpc_id"].(string),
			Subnet:         clusterDataNode1["subnet"].(string),
			InstanceSize:   clusterDataNode1["instance_size"].(string),
			DataVolumeSize: clusterDataNode1["data_volume_size"].(int),
		}

		copilotFaultTolerantDeployment.ClusterDataNodes = append(copilotFaultTolerantDeployment.ClusterDataNodes, clusterDataNode2)
	}

	return copilotFaultTolerantDeployment
}

func resourceAviatrixCopilotFaultTolerantDeploymentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	copilotFaultTolerantDeployment := marshalCopilotFaultTolerantDeploymentInput(d)

	if len(copilotFaultTolerantDeployment.ClusterDataNodes) < 3 {
		return diag.Errorf("at least three cluster data nodes are required")
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	flag := false
	defer resourceAviatrixCopilotFaultTolerantDeploymentReadIfRequired(ctx, d, meta, &flag)

	if err := client.CreateCopilotFaultTolerant(ctx, copilotFaultTolerantDeployment); err != nil {
		return diag.Errorf("could not start to deploy copilot: %v", err)
	}

	for i := 0; ; i++ {
		copilotAssociationStatus, err := client.GetCopilotAssociationStatus(ctx)

		if err != nil && err != goaviatrix.ErrNotFound {
			return diag.Errorf("could not get copilot association status: %v", err)
		}

		if err != goaviatrix.ErrNotFound && copilotAssociationStatus.Status {
			break
		}

		if i < 90 {
			time.Sleep(time.Duration(20) * time.Second)
		} else {
			return diag.Errorf("could not deploy copilot: %s", err)
		}
	}

	return resourceAviatrixCopilotFaultTolerantDeploymentReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixCopilotFaultTolerantDeploymentReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixCopilotFaultTolerantDeploymentRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixCopilotFaultTolerantDeploymentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return diag.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	copilotAssociationStatus, err := client.GetCopilotAssociationStatus(ctx)
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("could not get copilot association status: %v", err)
	}

	d.Set("main_copilot_private_ip", copilotAssociationStatus.IP)
	d.Set("main_copilot_public_ip", copilotAssociationStatus.PublicIp)

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixCopilotFaultTolerantDeploymentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	err := client.DeleteCopilotFaultTolerant(ctx)
	if err != nil {
		return diag.Errorf("could not delete copilot: %v", err)
	}

	return nil
}
