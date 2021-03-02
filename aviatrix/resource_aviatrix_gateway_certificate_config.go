package aviatrix

import (
	"context"
	"fmt"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixGatewayCertificateConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAviatrixGatewayCertificateConfigCreate,
		ReadContext:   resourceAviatrixGatewayCertificateConfigRead,
		DeleteContext: resourceAviatrixGatewayCertificateConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"ca_certificate": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "CA Certificate.",
			},
			"ca_private_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "CA Private Key.",
			},
		},
	}
}

func marshalGatewayCertificateConfigInput(d *schema.ResourceData) *goaviatrix.GatewayCertificate {
	return &goaviatrix.GatewayCertificate{
		CaCertificate: d.Get("ca_certificate").(string),
		CaPrivateKey:  d.Get("ca_private_key").(string),
	}
}

func resourceAviatrixGatewayCertificateConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	gwCert := marshalGatewayCertificateConfigInput(d)

	if err := client.ConfigureGatewayCertificate(ctx, gwCert); err != nil {
		return diag.FromErr(fmt.Errorf("could not configure gateway certificates: %v", err))
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixGatewayCertificateConfigRead(ctx, d, meta)
}

func resourceAviatrixGatewayCertificateConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	gwCertStatus, err := client.GetGatewayCertificateStatus(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not find gateway_certificate: %v", err))
	}

	if gwCertStatus == "enabled" {
		d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	} else {
		d.SetId("")
	}

	return nil
}

func resourceAviatrixGatewayCertificateConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if err := client.DisableGatewayCertificate(ctx); err != nil {
		return diag.FromErr(fmt.Errorf("could not disable gateway certificate checking: %v", err))
	}

	return nil
}
