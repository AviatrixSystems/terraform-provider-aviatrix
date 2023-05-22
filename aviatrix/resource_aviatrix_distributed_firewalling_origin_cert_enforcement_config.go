package aviatrix

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixDistributedFirewallingOriginCertEnforcementConfig() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixDistributedFirewallingOriginCertEnforcementConfigCreate,
		ReadWithoutTimeout:   resourceAviatrixDistributedFirewallingOriginCertEnforcementConfigRead,
		UpdateWithoutTimeout: resourceAviatrixDistributedFirewallingOriginCertEnforcementConfigUpdate,
		DeleteWithoutTimeout: resourceAviatrixDistributedFirewallingOriginCertEnforcementConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"enforcement_level": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Permissive",
				ValidateFunc: validation.StringInSlice([]string{"Strict", "Permissive", "Ignore"}, false),
				Description:  "Which origin cert enforcement level to set to for distributed firewalling.",
			},
		},
	}
}

func resourceAviatrixDistributedFirewallingOriginCertEnforcementConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	enforcementLevel := &goaviatrix.EnforcementLevel{
		Level: d.Get("enforcement_level").(string),
	}

	flag := false
	defer resourceAviatrixDistributedFirewallingOriginCertEnforcementConfigReadIfRequired(ctx, d, meta, &flag)

	err := client.SetEnforcementLevel(ctx, enforcementLevel)
	if err != nil {
		return diag.Errorf("failed to config Distributed-firewalling origin cert enforcement level: %s", err)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixDistributedFirewallingOriginCertEnforcementConfigReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixDistributedFirewallingOriginCertEnforcementConfigReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixDistributedFirewallingOriginCertEnforcementConfigRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixDistributedFirewallingOriginCertEnforcementConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return diag.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	enforcementLevel, err := client.GetEnforcementLevel(ctx)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read Distributed-firewalling origin cert enforcement level config: %s", err)
	}

	if enforcementLevel.Level == "ENFORCED" {
		d.Set("enforcement_level", "Strict")
	} else if enforcementLevel.Level == "DISABLED" {
		d.Set("enforcement_level", "Ignore")
	} else {
		d.Set("enforcement_level", "Permissive")
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixDistributedFirewallingOriginCertEnforcementConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	d.Partial(true)
	if d.HasChange("enforcement_level") {
		enforcementLevel := &goaviatrix.EnforcementLevel{
			Level: d.Get("enforcement_level").(string),
		}

		err := client.UpdateEnforcementLevel(ctx, enforcementLevel)
		if err != nil {
			return diag.Errorf("failed to set Distributed-firewalling origin cert enforcement config to %s in update: %s", enforcementLevel.Level, err)
		}
	}

	d.Partial(false)
	return resourceAviatrixDistributedFirewallingOriginCertEnforcementConfigRead(ctx, d, meta)
}

func resourceAviatrixDistributedFirewallingOriginCertEnforcementConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	err := client.DeleteEnforcementLevel(ctx)
	if err != nil {
		return diag.Errorf("failed to delete Distributed-firewalling origin cert enforcement level config: %v", err)
	}

	return nil
}
