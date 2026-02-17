package aviatrix

import (
	"context"
	"errors"
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
	client := mustClient(meta)

	proxy := marshalEdgeProxyProfileConfigInput(d)
	createdProxy, err := client.CreateEdgeProxyProfile(ctx, edgePlatformProxyProfileFromProxyProfile(proxy))
	if err != nil {
		return diag.Errorf("could not config proxy: %v", err)
	}
	mustSet(d, "proxy_profile_id", createdProxy.ProxyID)
	d.SetId(createdProxy.ProxyID)
	return nil
}

func resourceAviatrixEdgeProxyProfileConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	proxy := marshalEdgeProxyProfileConfigInput(d)

	client := mustClient(meta)
	existingProxy, err := client.GetEdgePlatformProxyProfile(ctx, proxy.AccountName, proxy.Name)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
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
	client := mustClient(meta)

	accountName := getString(d, "account_name")
	proxyProfileName := getString(d, "proxy_profile_name")
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
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("couldn't get proxy configuration: %s", err)
	}
	mustSet(d, "ip_address", proxy.IPAddress)
	mustSet(d, "port", proxy.Port)
	mustSet(d, "proxy_profile_id", proxy.ProxyID)
	mustSet(d, "name", proxy.Name)
	mustSet(d, "ca_certificate", proxy.CaCert)

	d.SetId(proxy.ProxyID)
	return nil
}

func resourceAviatrixEdgeProxyProfileConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := mustClient(meta)

	err := client.DeleteEdgePlatformProxyProfile(ctx, getString(d, "account_name"), getString(d, "proxy_profile_name"))
	if err != nil {
		return diag.Errorf("failed to delete proxy profile: %s", err)
	}

	return nil
}

func marshalEdgeProxyProfileConfigInput(d *schema.ResourceData) *proxyProfile {
	profile := &proxyProfile{
		AccountName: getString(d, "account_name"),
		Name:        getString(d, "proxy_profile_name"),
	}
	if v, ok := d.GetOk("ip_address"); ok {
		addr := mustString(v)
		profile.Address = &addr
	}
	if v, ok := d.GetOk("port"); ok {
		port := mustInt(v)
		profile.Port = &port
	}
	if v, ok := d.GetOk("ca_certificate"); ok {
		cert := mustString(v)
		profile.CACert = &cert
	}
	return profile
}
