package aviatrix

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixDistributedFirewallingIntraVpc() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixDistributedFirewallingIntraVpcCreate,
		ReadWithoutTimeout:   resourceAviatrixDistributedFirewallingIntraVpcRead,
		UpdateWithoutTimeout: resourceAviatrixDistributedFirewallingIntraVpcUpdate,
		DeleteWithoutTimeout: resourceAviatrixDistributedFirewallingIntraVpcDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},

		Schema: map[string]*schema.Schema{
			"vpcs": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Set of distributed-firewalling VPCs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Account Name of the VPC.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "vpc_id of the VPC.",
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Region of the VPC.",
						},
					},
				},
			},
		},
	}
}

func resourceAviatrixDistributedFirewallingIntraVpcCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	vpcList, err := marshalDistributedFirewallingIntraVpcListInput(d)
	if err != nil {
		return diag.Errorf("invalid inputs for Distributed-firewalling Intra VPC during create: %s\n", err)
	}

	flag := false
	defer resourceAviatrixDistributedFirewallingIntraVpcReadIfRequired(ctx, d, meta, &flag)

	err = client.CreateDistributedFirewallingIntraVpc(ctx, vpcList)
	if err != nil {
		return diag.Errorf("failed to create Distributed-firewalling Intra VPC: %s", err)
	}
	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixDistributedFirewallingIntraVpcReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixDistributedFirewallingIntraVpcReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixDistributedFirewallingIntraVpcRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixDistributedFirewallingIntraVpcRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	vpcList, err := client.GetDistributedFirewallingIntraVpc(ctx)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read Distributed-firewalling Intra VPC list: %s", err)
	}

	var vpcs []map[string]interface{}
	for _, vpc := range vpcList.VPCs {
		v := make(map[string]interface{})
		v["account_name"] = vpc.AccountName
		v["vpc_id"] = vpc.VpcId
		v["region"] = vpc.Region

		vpcs = append(vpcs, v)
	}

	if err := d.Set("vpcs", vpcs); err != nil {
		return diag.Errorf("failed to set vpcs during Distributed-firewalling Intra VPC read: %s\n", err)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixDistributedFirewallingIntraVpcUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	d.Partial(true)
	if d.HasChange("vpcs") {
		vpcList, err := marshalDistributedFirewallingIntraVpcListInput(d)
		if err != nil {
			return diag.Errorf("invalid inputs for Distributed-firewalling Intra VPC during update: %s\n", err)
		}
		err = client.CreateDistributedFirewallingIntraVpc(ctx, vpcList)
		if err != nil {
			return diag.Errorf("failed to update Distributed-firewalling Intra VPC: %s", err)
		}
	}

	d.Partial(false)
	return resourceAviatrixDistributedFirewallingIntraVpcRead(ctx, d, meta)
}

func resourceAviatrixDistributedFirewallingIntraVpcDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	err := client.DeleteDistributedFirewallingIntraVpc(ctx)
	if err != nil {
		return diag.Errorf("failed to delete Distributed-firewalling Intra VPC: %v", err)
	}

	return nil
}

func marshalDistributedFirewallingIntraVpcListInput(d *schema.ResourceData) (*goaviatrix.DistributedFirewallingIntraVpcList, error) {
	vpcList := &goaviatrix.DistributedFirewallingIntraVpcList{}

	vpcs := getList(d, "vpcs")
	for _, vpcInterface := range vpcs {
		vpc := mustMap(vpcInterface)

		distributedFirewallingVpc := &goaviatrix.DistributedFirewallingIntraVpc{
			AccountName: mustString(vpc["account_name"]),
			VpcId:       mustString(vpc["vpc_id"]),
			Region:      mustString(vpc["region"]),
		}

		vpcList.VPCs = append(vpcList.VPCs, *distributedFirewallingVpc)
	}

	return vpcList, nil
}
