package aviatrix

import (
	"context"
	"fmt"
	"strings"

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
			"policy_object": {
				Description: "Static list of DCF Policy Blocks or Lists.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "Name of the DCF Policy Block or List.",
							Required:    true,
							Type:        schema.TypeString,
						},
						"priority": {
							Description: "Priority of the DCF Policy Block or List.",
							Required:    true,
							Type:        schema.TypeInt,
						},
						"type": {
							Description: "Type of DCF Policy Object ('LIST' or 'BLOCK').",
							Required:    true,
							Type:        schema.TypeString,
						},
						"uuid": {
							Description: "UUID of the DCF Policy Block or List.",
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

	policyObjects, ok := d.Get("policy_object").(*schema.Set)
	if !ok {
		return nil, fmt.Errorf("Policy Block policy_object must be of type *schema.Set")
	}

	for _, policyObjectInterface := range policyObjects.List() {
		var ok bool

		policyObjectMap, ok := policyObjectInterface.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("policy_object must be of type map[string]interface{}")
		}

		subPolicy, err := marshalSubPolicyInput(policyObjectMap)
		if err != nil {
			return nil, err
		}

		policyBlock.SubPolicies = append(policyBlock.SubPolicies, *subPolicy)
	}

	policyBlock.UUID = d.Id()

	return policyBlock, nil
}

func marshalSubPolicyInput(subPolicyMap map[string]interface{}) (*goaviatrix.DCFSubPolicy, error) {
	var ok bool

	subPolicy := &goaviatrix.DCFSubPolicy{}

	if subPolicy.Name, ok = subPolicyMap["name"].(string); !ok {
		return nil, fmt.Errorf("policy_object name must be of type string")
	}

	if subPolicy.Priority, ok = subPolicyMap["priority"].(int); !ok {
		return nil, fmt.Errorf("policy_object priority must be of type string")
	}

	uuid, ok := subPolicyMap["uuid"].(string)
	if !ok {
		return nil, fmt.Errorf("policy_object uuid must be of type string")
	}

	subPolicyType, ok := subPolicyMap["type"].(string)
	if !ok {
		return nil, fmt.Errorf("policy_object type must be of type string")
	}

	switch strings.ToUpper(subPolicyType) {
	case "BLOCK":
		subPolicy.Block = uuid
	case "LIST":
		subPolicy.List = uuid
	default:
		return nil, fmt.Errorf("policy_object type must be either 'BLOCK' or 'List'")
	}

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
