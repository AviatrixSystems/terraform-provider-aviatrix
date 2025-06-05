package aviatrix

import (
	"context"
	"errors"
	"fmt"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixDCFPolicyBlock() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixDCFPolicyBlockCreate,
		ReadWithoutTimeout:   resourceAviatrixDCFPolicyBlockRead,
		UpdateWithoutTimeout: resourceAviatrixDCFPolicyBlockUpdate,
		DeleteWithoutTimeout: resourceAviatrixDCFPolicyBlockDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of the DCF Policy Block.",
				Required:    true,
				Type:        schema.TypeString,
			},
			"policy_block": {
				Description: "Static set of DCF Policy Blocks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "Name of the DCF Policy Block.",
							Required:    true,
							Type:        schema.TypeString,
						},
						"priority": {
							Description: "Priority of the DCF Policy Block.",
							Required:    true,
							Type:        schema.TypeInt,
						},
						"target_uuid": {
							Description: "Target UUID of the DCF Policy Block.",
							Required:    true,
							Type:        schema.TypeString,
						},
					},
				},
				Optional: true,
				Type:     schema.TypeSet,
			},
			"policy_list": {
				Description: "Static set of DCF Policy Lists.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "Name of the DCF Policy List.",
							Required:    true,
							Type:        schema.TypeString,
						},
						"priority": {
							Description: "Priority of the DCF Policy List.",
							Required:    true,
							Type:        schema.TypeInt,
						},
						"target_uuid": {
							Description: "Target UUID of the DCF Policy List.",
							Required:    true,
							Type:        schema.TypeString,
						},
					},
				},
				Optional: true,
				Type:     schema.TypeSet,
			},
		},
	}
}

func marshalDCFPolicyBlockInput(d *schema.ResourceData) (*goaviatrix.DCFPolicyBlock, error) {
	policyBlock := &goaviatrix.DCFPolicyBlock{}

	name, ok := d.Get("name").(string)
	if !ok {
		return nil, fmt.Errorf("PolicyBlock name must be of type string")
	}
	policyBlock.Name = name

	policyBlocks, ok := d.Get("policy_block").(*schema.Set)
	if !ok {
		return nil, fmt.Errorf("policy_block must be of type *schema.Set")
	}

	for _, policyBlockInterface := range policyBlocks.List() {
		var ok bool

		policyBlockMap, ok := policyBlockInterface.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("policy_block interface must be of type map[string]interface{}")
		}

		subPolicy, err := marshalSubPolicyBlockInput(policyBlockMap)
		if err != nil {
			return nil, err
		}

		policyBlock.SubPolicies = append(policyBlock.SubPolicies, *subPolicy)
	}

	policyLists, ok := d.Get("policy_list").(*schema.Set)
	if !ok {
		return nil, fmt.Errorf("policy_list must be of type *schema.Set")
	}

	for _, policyListInterface := range policyLists.List() {
		var ok bool

		policyListMap, ok := policyListInterface.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("policy_list interface must be of type map[string]interface{}")
		}

		subPolicy, err := marshalSubPolicyListInput(policyListMap)
		if err != nil {
			return nil, err
		}

		policyBlock.SubPolicies = append(policyBlock.SubPolicies, *subPolicy)
	}

	if len(policyBlock.SubPolicies) == 0 {
		return nil, fmt.Errorf("policy_block type must contain a sub-policy (block or list)")
	}

	policyBlock.UUID = d.Id()

	return policyBlock, nil
}

func marshalSubPolicyBlockInput(subPolicyMap map[string]interface{}) (*goaviatrix.DCFSubPolicy, error) {
	var ok bool

	subPolicy := &goaviatrix.DCFSubPolicy{}

	if subPolicy.Name, ok = subPolicyMap["name"].(string); !ok {
		return nil, fmt.Errorf("policy_block name must be of type string")
	}

	if subPolicy.Priority, ok = subPolicyMap["priority"].(int); !ok {
		return nil, fmt.Errorf("policy_block priority must be of type string")
	}

	target_uuid, ok := subPolicyMap["target_uuid"].(string)
	if !ok {
		return nil, fmt.Errorf("policy_block target_uuid must be of type string")
	}

	subPolicy.Block = target_uuid

	return subPolicy, nil
}

func marshalSubPolicyListInput(subPolicyMap map[string]interface{}) (*goaviatrix.DCFSubPolicy, error) {
	var ok bool

	subPolicy := &goaviatrix.DCFSubPolicy{}

	if subPolicy.Name, ok = subPolicyMap["name"].(string); !ok {
		return nil, fmt.Errorf("policy_list name must be of type string")
	}

	if subPolicy.Priority, ok = subPolicyMap["priority"].(int); !ok {
		return nil, fmt.Errorf("policy_list priority must be of type string")
	}

	target_uuid, ok := subPolicyMap["target_uuid"].(string)
	if !ok {
		return nil, fmt.Errorf("policy_list target_uuid must be of type string")
	}

	subPolicy.List = target_uuid

	return subPolicy, nil
}

func resourceAviatrixDCFPolicyBlockCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*goaviatrix.Client)
	if !ok {
		return diag.Errorf("client must be of type *goaviatrix.Client")
	}

	policyBlock, err := marshalDCFPolicyBlockInput(d)
	if err != nil {
		return diag.Errorf("invalid inputs for DCF Policy during create: %s\n", err)
	}

	uuid, err := client.CreateDCFPolicyBlock(ctx, policyBlock)
	if err != nil {
		return diag.Errorf("failed to create DCF Policy List: %s", err)
	}

	d.SetId(uuid)

	return nil
}

func resourceAviatrixDCFPolicyBlockRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*goaviatrix.Client)
	if !ok {
		return diag.Errorf("client must be of type *goaviatrix.Client")
	}

	uuid := d.Id()

	policyBlock, err := client.GetDCFPolicyBlock(ctx, uuid)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read DCF Policy List: %s", err)
	}

	if err := d.Set("name", policyBlock.Name); err != nil {
		return diag.Errorf("failed to set name during DCF Policy List read: %s\n", err)
	}

	d.SetId(policyBlock.UUID)

	return nil
}

func resourceAviatrixDCFPolicyBlockUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*goaviatrix.Client)
	if !ok {
		return diag.Errorf("client must be of type *goaviatrix.Client")
	}

	policyBlock, err := marshalDCFPolicyBlockInput(d)
	if err != nil {
		return diag.Errorf("invalid inputs for DCF Policy during update: %s\n", err)
	}

	err = client.UpdateDCFPolicyBlock(ctx, policyBlock)
	if err != nil {
		return diag.Errorf("failed to update DCF Policy List: %s", err)
	}

	return nil
}

func resourceAviatrixDCFPolicyBlockDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*goaviatrix.Client)
	if !ok {
		return diag.Errorf("client must be of type *goaviatrix.Client")
	}

	uuid := d.Id()

	err := client.DeleteDCFPolicyBlock(ctx, uuid)
	if err != nil {
		return diag.Errorf("failed to delete DCF Policy List: %v", err)
	}

	return nil
}
