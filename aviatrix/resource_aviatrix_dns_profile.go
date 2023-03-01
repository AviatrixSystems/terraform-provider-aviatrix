package aviatrix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixDNSProfile() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAviatrixDNSProfileCreate,
		ReadWithoutTimeout:   resourceAviatrixDNSProfileRead,
		UpdateWithoutTimeout: resourceAviatrixDNSProfileUpdate,
		DeleteWithoutTimeout: resourceAviatrixDNSProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS profile name.",
			},
			"global_dns_servers": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of global DNS servers.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"local_domain_names": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of local domain names.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"lan_dns_servers": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of LAN DNS servers.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"wan_dns_servers": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of WAN DNS servers.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func marshalDNSProfileInput(d *schema.ResourceData) map[string]interface{} {
	data := make(map[string]interface{})
	var templateName []string

	dnsProfile := &goaviatrix.DNSProfile{
		Global:           getStringList(d, "global_dns_servers"),
		Lan:              getStringList(d, "lan_dns_servers"),
		LocalDomainNames: getStringList(d, "local_domain_names"),
		Wan:              getStringList(d, "wan_dns_servers"),
	}

	name := d.Get("name").(string)
	templateName = append(templateName, name)

	data["template_names"] = templateName
	data[name] = dnsProfile

	return data
}

func resourceAviatrixDNSProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	data := marshalDNSProfileInput(d)

	d.SetId(d.Get("name").(string))
	flag := false
	defer resourceAviatrixDNSProfileReadIfRequired(ctx, d, meta, &flag)

	if err := client.CreateDNSProfile(ctx, data); err != nil {
		return diag.Errorf("could not create DNS profile: %v", err)
	}

	return resourceAviatrixDNSProfileReadIfRequired(ctx, d, meta, &flag)
}

func resourceAviatrixDNSProfileReadIfRequired(ctx context.Context, d *schema.ResourceData, meta interface{}, flag *bool) diag.Diagnostics {
	if !(*flag) {
		*flag = true
		return resourceAviatrixDNSProfileRead(ctx, d, meta)
	}
	return nil
}

func resourceAviatrixDNSProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Get("name").(string) == "" {
		id := d.Id()
		d.Set("name", id)
		d.SetId(id)
	}

	profile, err := client.GetDNSProfile(ctx, d.Get("name").(string))
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read DNS profile: %s", err)
	}

	d.Set("global_dns_servers", profile["global"])
	d.Set("lan_dns_servers", profile["lan"])
	d.Set("local_domain_names", profile["local_domain_names"])
	d.Set("wan_dns_servers", profile["wan"])

	d.SetId(d.Get("name").(string))
	return nil
}

func resourceAviatrixDNSProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	d.Partial(true)

	data := marshalDNSProfileInput(d)

	err := client.UpdateDNSProfile(ctx, data)
	if err != nil {
		return diag.Errorf("could not update DNS profile: %v", err)
	}

	d.Partial(false)

	return resourceAviatrixDNSProfileRead(ctx, d, meta)
}

func resourceAviatrixDNSProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	data := marshalDNSProfileInput(d)

	err := client.DeleteDNSProfile(ctx, data)
	if err != nil {
		return diag.Errorf("could not delete DNS profile: %v", err)
	}

	return nil
}
