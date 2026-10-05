package aviatrix

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

type proxyProfile struct {
	Action      string  `json:"action"`
	CID         string  `json:"CID"`
	AccountName string  `json:"account_name"`
	Name        string  `json:"proxy_name"`
	Address     *string `json:"address"`
	Port        *int    `json:"port"`
	CACert      *string `json:"ca_cert"`
}

func resourceAviatrixEdgeProxyProfileConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAviatrixEdgeProxyProfileConfigCreate,
		UpdateContext: resourceAviatrixEdgeProxyProfileConfigUpdate,
		ReadContext:   resourceAviatrixEdgeProxyProfileConfigRead,
		DeleteContext: resourceAviatrixEdgeProxyProfileConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Edge NEO account name.",
			},
			"proxy_profile_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Edge Proxy Profile name.",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "HTTPS proxy IP.",
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "HTTPS proxy Port.",
			},
			"ca_certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Server CA Certificate in base64 encoded PEM format",
			},
			"proxy_profile_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Edge Proxy Profile ID.",
			},
		},
	}
}

func edgePlatformProxyProfileFromProxyProfile(proxy *proxyProfile) *goaviatrix.EdgePlatformProxyProfile {
	return &goaviatrix.EdgePlatformProxyProfile{
		AccountName: proxy.AccountName,
		Name:        proxy.Name,
		Address: func() string {
			if proxy.Address == nil {
				return ""
			}
			return *proxy.Address
		}(),
		Port: func() int {
			if proxy.Port == nil {
				return 0
			}
			return *proxy.Port
		}(),
		CACert: func() string {
			if proxy.CACert == nil {
				return ""
			}
			return *proxy.CACert
		}(),
	}
}

func resourceAviatrixEdgeProxyProfileConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	proxy := marshalEdgeProxyProfileConfigInput(d)
	createdProxy, err := client.CreateEdgeProxyProfile(ctx, edgePlatformProxyProfileFromProxyProfile(proxy))
	if err != nil {
		return diag.Errorf("could not config proxy: %v", err)
	}

	d.Set("proxy_profile_id", createdProxy.ProxyID)
	d.SetId(createdProxy.ProxyID)
	return nil
}

func resourceAviatrixEdgeProxyProfileConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	proxy := marshalEdgeProxyProfileConfigInput(d)

	client := meta.(*goaviatrix.Client)
	existingProxy, err := client.GetEdgePlatformProxyProfile(ctx, proxy.AccountName, proxy.Name)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("couldn't get proxy configuration: %s", err)
	}

	modifiedProxy := &goaviatrix.EdgePlatformProxyProfileUpdate{
		EdgePlatformProxyProfile: goaviatrix.EdgePlatformProxyProfile{
			AccountName: proxy.AccountName,
			Name:        proxy.Name,
			Address:     existingProxy.IPAddress,
			Port:        existingProxy.Port,
			CACert:      existingProxy.CaCert,
		},
		ProxyID: existingProxy.ProxyID,
	}

	if proxy.Address != nil {
		modifiedProxy.Address = *proxy.Address
	}
	if proxy.Port != nil {
		modifiedProxy.Port = *proxy.Port
	}
	if proxy.CACert != nil {
		modifiedProxy.CACert = *proxy.CACert
	}

	if err := client.UpdateEdgeProxyProfile(ctx, modifiedProxy); err != nil {
		return diag.Errorf("could not config proxy: %v", err)
	}

	d.SetId(modifiedProxy.ProxyID)
	return nil
}

func resourceAviatrixEdgeProxyProfileConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	accountName := d.Get("account_name").(string)
	proxyProfileName := d.Get("proxy_profile_name").(string)
	if accountName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no account name received. Import Id is %s", id)
		parts := strings.Split(id, "~")
		if len(parts) != 2 {
			return diag.Errorf("Invalid Import ID received, ID must be in the format account_name~proxy_profile_name")
		}
		accountName = parts[0]
		proxyProfileName = parts[1]
		d.SetId(id)
	}

	proxy, err := client.GetEdgePlatformProxyProfile(ctx, accountName, proxyProfileName)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("couldn't get proxy configuration: %s", err)
	}

	d.Set("ip_address", proxy.IPAddress)
	d.Set("port", proxy.Port)
	d.Set("proxy_profile_id", proxy.ProxyID)
	d.Set("name", proxy.Name)
	d.Set("ca_certificate", proxy.CaCert)

	d.SetId(proxy.ProxyID)
	return nil
}

func resourceAviatrixEdgeProxyProfileConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	err := client.DeleteEdgePlatformProxyProfile(ctx, d.Get("account_name").(string), d.Get("proxy_profile_name").(string))
	if err != nil {
		return diag.Errorf("failed to delete proxy profile: %s", err)
	}

	return nil
}

func marshalEdgeProxyProfileConfigInput(d *schema.ResourceData) *proxyProfile {
	profile := &proxyProfile{
		AccountName: d.Get("account_name").(string),
		Name:        d.Get("proxy_profile_name").(string),
	}
	if v, ok := d.GetOk("ip_address"); ok {
		addr := v.(string)
		profile.Address = &addr
	}
	if v, ok := d.GetOk("port"); ok {
		port := v.(int)
		profile.Port = &port
	}
	if v, ok := d.GetOk("ca_certificate"); ok {
		cert := v.(string)
		profile.CACert = &cert
	}
	return profile
}
