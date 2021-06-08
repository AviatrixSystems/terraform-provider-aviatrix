package aviatrix

import (
	"context"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixControllerCertDomainConfig() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixControllerCertDomainConfigCreate,
		ReadWithoutTimeout:   resourceAviatrixControllerCertDomainConfigRead,
		UpdateWithoutTimeout: resourceAviatrixControllerCertDomainConfigUpdate,
		DeleteWithoutTimeout: resourceAviatrixControllerCertDomainConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"cert_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "aviatrixnetwork.com",
				Description: "Domain name that is used in FQDN for generating cert.",
			},
		},
	}
}

func resourceAviatrixControllerCertDomainConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	certDomain := d.Get("cert_domain").(string)

	err := client.SetCertDomain(ctx, certDomain)
	if err != nil {
		if strings.Contains(err.Error(), "EOF") {
			certDomainConfig, err := client.GetCertDomain(ctx)
			if err != nil {
				return diag.Errorf("could not confirm if cert domain is updated: %v", err)
			}
			if certDomainConfig.CertDomain != certDomain {
				return diag.Errorf("could not set cert domain: %v", err)
			}
		} else {
			return diag.Errorf("could not set cert domain: %v", err)
		}
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixControllerCertDomainConfigRead(ctx, d, meta)
}

func resourceAviatrixControllerCertDomainConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return diag.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	certDomainConfig, err := client.GetCertDomain(ctx)
	if err != nil {
		return diag.Errorf("could not get cert domain: %v", err)
	}

	d.Set("cert_domain", certDomainConfig.CertDomain)

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixControllerCertDomainConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.HasChange("cert_domain") {
		err := client.SetCertDomain(ctx, d.Get("cert_domain").(string))
		if err != nil {
			if strings.Contains(err.Error(), "EOF") {
				certDomainConfig, err := client.GetCertDomain(ctx)
				if err != nil {
					return diag.Errorf("could not confirm if cert domain is updated: %v", err)
				}
				if certDomainConfig.CertDomain != d.Get("cert_domain") {
					return diag.Errorf("could not update cert domain: %v", err)
				}
			} else {
				return diag.Errorf("could not update cert domain: %v", err)
			}
		}
	}

	return resourceAviatrixControllerCertDomainConfigRead(ctx, d, meta)
}

func resourceAviatrixControllerCertDomainConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	err := client.SetCertDomain(ctx, "aviatrixnetwork.com")
	if err != nil {
		if strings.Contains(err.Error(), "EOF") {
			certDomainConfig, err := client.GetCertDomain(ctx)
			if err != nil {
				return diag.Errorf("could not confirm if cert domain is updated: %v", err)
			}
			if certDomainConfig.CertDomain != "aviatrixnetwork.com" {
				return diag.Errorf("could not reset cert domain: %v", err)
			}
		} else {
			return diag.Errorf("could not reset cert domain: %v", err)
		}
	}

	return nil
}
