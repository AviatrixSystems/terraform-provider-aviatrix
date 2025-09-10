package aviatrix

import (
	"context"
	"errors"
	"fmt"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//nolint:funlen
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
			"attach_to": {
				Description: "Attach the DCF Policy Block to an attachment point.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"system_resource": {
				Description: "Indicates if the DCF Policy Block is a system resource.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"name": {
				Description: "Name of the DCF Policy Block.",
				Required:    true,
				Type:        schema.TypeString,
			},
			"policy_block_reference": {
				Description: "Static set of DCF Policy Blocks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
			"policy_list_reference": {
				Description: "Static set of DCF Policy Lists.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
			"attachment_point": {
				Description: "Static set of DCF Attachment Points.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"priority": {
							Description: "Priority of the DCF Policy List.",
							Required:    true,
							Type:        schema.TypeInt,
						},
						"name": {
							Description: "Name of the DCF Attachment Point.",
							Required:    true,
							Type:        schema.TypeString,
						},
						"target_uuid": {
							Description: "Target UUID of the DCF Attachment Point.",
							Computed:    true,
							Type:        schema.TypeString,
						},
						"uuid": {
							Description: "UUID of the DCF Attachment Point.",
							Computed:    true,
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

//nolint:cyclop
func marshalDCFPolicyBlockInput(d *schema.ResourceData) (*goaviatrix.DCFPolicyBlock, error) {
	policyBlock := &goaviatrix.DCFPolicyBlock{}

	name, ok := d.Get("name").(string)
	if !ok {
		return nil, fmt.Errorf("PolicyBlock name must be of type string")
	}
	policyBlock.Name = name

	attachTo, ok := d.Get("attach_to").(string)
	if !ok {
		return nil, fmt.Errorf("PolicyBlock attach_to must be of type string")
	}
	policyBlock.AttachTo = attachTo

	policyBlocks, ok := d.Get("policy_block_reference").(*schema.Set)
	if !ok {
		return nil, fmt.Errorf("policy_block_reference must be of type *schema.Set")
	}

	for _, policyBlockInterface := range policyBlocks.List() {
		var ok bool

		policyBlockMap, ok := policyBlockInterface.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("policy_block_reference interface must be of type map[string]interface{}")
		}

		subPolicy, err := marshalSubPolicyBlockInput(policyBlockMap)
		if err != nil {
			return nil, err
		}

		policyBlock.SubPolicies = append(policyBlock.SubPolicies, *subPolicy)
	}

	policyLists, ok := d.Get("policy_list_reference").(*schema.Set)
	if !ok {
		return nil, fmt.Errorf("policy_list_reference must be of type *schema.Set")
	}

	for _, policyListInterface := range policyLists.List() {
		var ok bool

		policyListMap, ok := policyListInterface.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("policy_list_reference interface must be of type map[string]interface{}")
		}

		subPolicy, err := marshalSubPolicyListInput(policyListMap)
		if err != nil {
			return nil, err
		}

		policyBlock.SubPolicies = append(policyBlock.SubPolicies, *subPolicy)
	}

	attachmentPoints, ok := d.Get("attachment_point").(*schema.Set)
	if !ok {
		return nil, fmt.Errorf("attachment_point must be of type *schema.Set")
	}
	for _, attachmentPointInterface := range attachmentPoints.List() {
		var ok bool
		attachmentPointMap, ok := attachmentPointInterface.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("attachment_point interface must be of type map[string]interface{}")
		}
		subPolicy := &goaviatrix.DCFSubPolicy{}
		if subPolicy.Priority, ok = attachmentPointMap["priority"].(int); !ok {
			return nil, fmt.Errorf("attachment_point priority must be of type int")
		}
		attachmentPoint := &goaviatrix.AttachmentPoint{}
		if attachmentPoint.Name, ok = attachmentPointMap["name"].(string); !ok {
			return nil, fmt.Errorf("attachment_point name must be of type string")
		}
		if attachmentPoint.UUID, ok = attachmentPointMap["uuid"].(string); !ok {
			return nil, fmt.Errorf("attachment_point uuid must be of type string")
		}
		if attachmentPoint.TargetUUID, ok = attachmentPointMap["target_uuid"].(string); !ok {
			return nil, fmt.Errorf("attachment_point target_uuid must be of type string")
		}
		subPolicy.AttachmentPoint = attachmentPoint
		policyBlock.SubPolicies = append(policyBlock.SubPolicies, *subPolicy)
	}
	systemResource, ok := d.Get("system_resource").(bool)
	if !ok {
		return nil, fmt.Errorf("PolicyBlock system_resource must be of type bool")
	}
	policyBlock.SystemResource = systemResource

	if len(policyBlock.SubPolicies) == 0 {
		return nil, fmt.Errorf("policy_block_reference type must contain a sub-policy (block or list)")
	}

	policyBlock.UUID = d.Id()

	return policyBlock, nil
}

func marshalSubPolicyBlockInput(subPolicyMap map[string]interface{}) (*goaviatrix.DCFSubPolicy, error) {
	var ok bool

	subPolicy := &goaviatrix.DCFSubPolicy{}

	if subPolicy.Priority, ok = subPolicyMap["priority"].(int); !ok {
		return nil, fmt.Errorf("policy_block_reference priority must be of type string")
	}

	targetUUID, ok := subPolicyMap["target_uuid"].(string)
	if !ok {
		return nil, fmt.Errorf("policy_block_reference target_uuid must be of type string")
	}

	subPolicy.Block = targetUUID

	return subPolicy, nil
}

func marshalSubPolicyListInput(subPolicyMap map[string]interface{}) (*goaviatrix.DCFSubPolicy, error) {
	var ok bool

	subPolicy := &goaviatrix.DCFSubPolicy{}

	if subPolicy.Priority, ok = subPolicyMap["priority"].(int); !ok {
		return nil, fmt.Errorf("policy_list_reference priority must be of type string")
	}

	targetUUID, ok := subPolicyMap["target_uuid"].(string)
	if !ok {
		return nil, fmt.Errorf("policy_list_reference target_uuid must be of type string")
	}

	subPolicy.List = targetUUID

	return subPolicy, nil
}

func resourceAviatrixDCFPolicyBlockCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*goaviatrix.Client)
	if !ok {
		return diag.Errorf("client must be of type *goaviatrix.Client")
	}

	policyBlock, err := marshalDCFPolicyBlockInput(d)
	if err != nil {
		return diag.Errorf("invalid inputs for DCF Policy during create: %s", err)
	}
	uuid, err := client.CreateDCFPolicyBlock(ctx, policyBlock)
	if err != nil {
		return diag.Errorf("failed to create DCF Policy List: %s", err)
	}

	d.SetId(uuid)

	return nil
}

//nolint:cyclop
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
		return diag.Errorf("failed to set name during DCF Policy Block read: %s", err)
	}
	if err := d.Set("system_resource", policyBlock.SystemResource); err != nil {
		return diag.Errorf("failed to set system_resource during DCF Policy Block read: %s", err)
	}
	if err := d.Set("attach_to", policyBlock.AttachTo); err != nil {
		return diag.Errorf("failed to set attach_to during DCF Policy Block read: %s", err)
	}
	policyLists, ok := d.Get("policy_list_reference").(*schema.Set)
	if !ok {
		return diag.Errorf("policy_list_reference must be of type *schema.Set")
	}
	policyBlocks, ok := d.Get("policy_block_reference").(*schema.Set)
	if !ok {
		return diag.Errorf("policy_block_reference must be of type *schema.Set")
	}
	policyAttachmentPoints, ok := d.Get("attachment_point").(*schema.Set)
	if !ok {
		return diag.Errorf("attachment_point must be of type *schema.Set")
	}
	for _, subPolicy := range policyBlock.SubPolicies {
		if subPolicy.List != "" {
			policyLists.Add(map[string]interface{}{
				"priority":    subPolicy.Priority,
				"target_uuid": subPolicy.List,
			})
			if err := d.Set("policy_list_reference", policyLists); err != nil {
				return diag.Errorf("failed to set policy_list_reference during DCF Policy List read: %s", err)
			}
		} else if subPolicy.Block != "" {
			policyBlocks.Add(map[string]interface{}{
				"priority":    subPolicy.Priority,
				"target_uuid": subPolicy.Block,
			})
			if err := d.Set("policy_block_reference", policyBlocks); err != nil {
				return diag.Errorf("failed to set policy_block_reference during DCF Policy List read: %s", err)
			}
		} else if subPolicy.AttachmentPoint != (&goaviatrix.AttachmentPoint{}) {
			policyAttachmentPoints.Add(map[string]interface{}{
				"name":        subPolicy.AttachmentPoint.Name,
				"uuid":        subPolicy.AttachmentPoint.UUID,
				"target_uuid": subPolicy.AttachmentPoint.TargetUUID,
				"priority":    subPolicy.Priority,
			})
			if err := d.Set("attachment_point", policyAttachmentPoints); err != nil {
				return diag.Errorf("failed to set attachment_point during DCF Policy List read: %s", err)
			}
		}
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
		return diag.Errorf("invalid inputs for DCF Policy during update: %s", err)
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
