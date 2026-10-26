package aviatrix

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixControllerAccessAllowListConfig() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixControllerAccessAllowListConfigCreate,
		ReadWithoutTimeout:   resourceAviatrixControllerAccessAllowListConfigRead,
		UpdateWithoutTimeout: resourceAviatrixControllerAccessAllowListConfigUpdate,
		DeleteWithoutTimeout: resourceAviatrixControllerAccessAllowListConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"allow_list": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Set of IPs allowed access to the Controller.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IP address allowed access to the Controller.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description of the IP address.",
						},
					},
				},
			},
			"enable_enforce": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Set to true to enable enforcement of the `allow_list {}`'s IPs.",
			},
		},
	}
}

func marshalControllerAccessAllowListConfigInput(d *schema.ResourceData) *goaviatrix.AllowList {
	var allowList goaviatrix.AllowList

	al := d.Get("allow_list").(*schema.Set).List()
	for _, v0 := range al {
		v1 := v0.(map[string]interface{})

		ai := goaviatrix.AllowIp{
			IpAddress:   v1["ip_address"].(string),
			Description: v1["description"].(string),
		}

		allowList.AllowList = append(allowList.AllowList, ai)
	}

	allowList.Enforce = d.Get("enable_enforce").(bool)

	return &allowList
}

func resourceAviatrixControllerAccessAllowListConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	allowList := marshalControllerAccessAllowListConfigInput(d)

	flag := false
	defer resourceAviatrixControllerAccessAllowListConfigReadIfRequired(ctx, d, meta, &flag)

	err := client.CreateControllerAccessAllowList(ctx, allowList)
	if err != nil {
		return diag.Errorf("failed to create controller access allow list config: %s", err)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixControllerAccessAllowListConfigReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixControllerAccessAllowListConfigReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixControllerAccessAllowListConfigRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixControllerAccessAllowListConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return diag.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	allowList, err := client.GetControllerAccessAllowList(ctx)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read controller access allow list config: %s", err)
	}

	d.Set("enable_enforce", allowList.Enforce)

	var al []interface{}

	for _, v0 := range allowList.AllowList {
		v1 := make(map[string]interface{})

		v1["ip_address"] = v0.IpAddress
		v1["description"] = v0.Description
		al = append(al, v1)
	}

	if err = d.Set("allow_list", al); err != nil {
		return diag.Errorf("failed to set allow_list: %s", err)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixControllerAccessAllowListConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	d.Partial(true)
	if d.HasChanges("allow_list", "enable_enforce") {
		allowList := marshalControllerAccessAllowListConfigInput(d)

		err := client.CreateControllerAccessAllowList(ctx, allowList)
		if err != nil {
			return diag.Errorf("failed to update controller access allow list config: %s", err)
		}
	}

	d.Partial(false)
	return resourceAviatrixControllerAccessAllowListConfigRead(ctx, d, meta)
}

func resourceAviatrixControllerAccessAllowListConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	err := client.DeleteControllerAccessAllowList(ctx)
	if err != nil {
		return diag.Errorf("failed to delete controller access allow list config: %v", err)
	}

	return nil
}
