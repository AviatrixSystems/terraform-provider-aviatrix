package aviatrix

import (
	"context"
	"fmt"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAviatrixDistributedFirewallingZeroTrustRule() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixDistributedFirewallingZeroTrustRuleCreate,
		UpdateWithoutTimeout: resourceAviatrixDistributedFirewallingZeroTrustRuleUpdate,
		ReadWithoutTimeout:   resourceAviatrixDistributedFirewallingZeroTrustRuleRead,
		DeleteWithoutTimeout: resourceAviatrixDistributedFirewallingZeroTrustRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The action for the zero trust rule, DENY or PERMIT.",
			},
			"logging": {
				Type:        schema.TypeBool,
				Required:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "Boolean value to enable or disable logging for the zero trust rule.",
			},
		},
	}
}

func resourceAviatrixDistributedFirewallingZeroTrustRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	proxyCaConfig := &goaviatrix.ProxyCaConfig{
		CaCert: d.Get("ca_cert").(string),
		CaKey:  d.Get("ca_key").(string),
	}

	if err := client.SetNewCertificate(ctx, proxyCaConfig); err != nil {
		return diag.Errorf("failed to set new Distributed-firewalling proxy ca certificate: %v", err)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixDistributedFirewallingProxyCaConfigRead(ctx, d, meta)
}

func resourceAviatrixDistributedFirewallingZeroTrustRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return diag.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	proxyCaConfig, err := client.GetCaCertificate(ctx)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read Distributed-firewalling proxy ca cert config: %s", err)
	}
	d.Set("ca_cert", proxyCaConfig.CaCert)

	proxyCaCertInstance, err := client.GetMetaCaCertificate(ctx)
	if err == nil {
		d.Set("common_name", proxyCaCertInstance.CommonName)
		d.Set("expiration_time", proxyCaCertInstance.ExpirationDate)
		d.Set("issuer_name", proxyCaCertInstance.Issuer)
		d.Set("unique_serial", proxyCaCertInstance.SerialNumber)
		d.Set("upload_info", proxyCaCertInstance.UploadInfo)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixDistributedFirewallingZeroTrustRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	err := client.DeleteCaCertificate(ctx)
	if err != nil {
		return diag.Errorf("failed to delete Distributed-firewalling proxy ca cert : %s", err)
	}

	return nil
}

func resourceAviatrixDistributedFirewallingZeroTrustRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	d.Partial(true)
	if d.HasChange("policies") {
		policyList, err := marshalDistributedFirewallingPolicyListInput(d)
		if err != nil {
			return diag.Errorf("invalid inputs for Distributed-firewalling Policy during update: %s\n", err)
		}
		err = client.UpdateDistributedFirewallingPolicyList(ctx, policyList)
		if err != nil {
			return diag.Errorf("failed to update Distributed-firewalling policies: %s", err)
		}
	}

	d.Partial(false)
	return resourceAviatrixDistributedFirewallingPolicyListRead(ctx, d, meta)
}
