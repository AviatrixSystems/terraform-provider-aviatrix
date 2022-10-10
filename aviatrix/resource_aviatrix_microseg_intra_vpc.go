package aviatrix

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixMicrosegIntraVpc() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixMicrosegIntraVpcCreate,
		ReadWithoutTimeout:   resourceAviatrixMicrosegIntraVpcRead,
		UpdateWithoutTimeout: resourceAviatrixMicrosegIntraVpcUpdate,
		DeleteWithoutTimeout: resourceAviatrixMicrosegIntraVpcDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpcs": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Set of micro-segmentation VPCs.",
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

func resourceAviatrixMicrosegIntraVpcCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	vpcList, err := marshalMicrosegIntraVpcListInput(d)
	if err != nil {
		return diag.Errorf("invalid inputs for Micro-segmentation Intra VPC during create: %s\n", err)
	}

	flag := false
	defer resourceAviatrixMicrosegIntraVpcReadIfRequired(ctx, d, meta, &flag)

	err = client.CreateMicrosegIntraVpc(ctx, vpcList)
	if err != nil {
		return diag.Errorf("failed to create Micro-segmentation Intra VPC: %s", err)
	}
	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixMicrosegIntraVpcReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixMicrosegIntraVpcReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixMicrosegIntraVpcRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixMicrosegIntraVpcRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	vpcList, err := client.GetMicrosegIntraVpc(ctx)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read Micro-segmentation Intra VPC list: %s", err)
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
		return diag.Errorf("failed to set vpcs during Micro-segmentation Intra VPC read: %s\n", err)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixMicrosegIntraVpcUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	d.Partial(true)
	if d.HasChange("vpcs") {
		vpcList, err := marshalMicrosegIntraVpcListInput(d)
		if err != nil {
			return diag.Errorf("invalid inputs for Micro-segmentation Intra VPC during update: %s\n", err)
		}
		err = client.CreateMicrosegIntraVpc(ctx, vpcList)
		if err != nil {
			return diag.Errorf("failed to update Micro-segmentation Intra VPC: %s", err)
		}
	}

	d.Partial(false)
	return resourceAviatrixMicrosegIntraVpcRead(ctx, d, meta)
}

func resourceAviatrixMicrosegIntraVpcDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	err := client.DeleteMicrosegIntraVpc(ctx)
	if err != nil {
		return diag.Errorf("failed to delete Micro-segmentation Intra VPC: %v", err)
	}

	return nil
}

func marshalMicrosegIntraVpcListInput(d *schema.ResourceData) (*goaviatrix.MicrosegIntraVpcList, error) {
	vpcList := &goaviatrix.MicrosegIntraVpcList{}

	vpcs := d.Get("vpcs").([]interface{})
	for _, vpcInterface := range vpcs {
		vpc := vpcInterface.(map[string]interface{})

		microsegVpc := &goaviatrix.MicrosegIntraVpc{
			AccountName: vpc["account_name"].(string),
			VpcId:       vpc["vpc_id"].(string),
			Region:      vpc["region"].(string),
		}

		vpcList.VPCs = append(vpcList.VPCs, *microsegVpc)
	}

	return vpcList, nil
}
