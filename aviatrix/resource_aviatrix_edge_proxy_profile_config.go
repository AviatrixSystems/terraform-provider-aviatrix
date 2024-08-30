package aviatrix

import (
	"context"
	"fmt"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixEdgeProxyProfileConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixEdgeProxyProfileConfigCreate,
		Read:   resourceAviatrixEdgeProxyProfileConfigRead,
		Delete: resourceAviatrixEdgeProxyProfileConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				ForceNew:    true,
				Description: "HTTPS proxy IP.",
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "HTTPS proxy Port.",
			},
			"ca_certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Server CA Certificate file.",
			},
			"proxy_profile_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Edge Proxy Profile ID.",
			},
		},
	}
}

func resourceAviatrixEdgeProxyProfileConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	proxy := marshalEdgeProxyProfileConfigInput(d)
	createdProxy, err := client.CreateEdgeProxyProfile(context.Background(), proxy)
	if err != nil {
		return fmt.Errorf("could not config proxy: %v", err)
	}

	d.SetId(createdProxy.ProxyID)
	return nil
}

func resourceAviatrixEdgeProxyProfileConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	proxy, err := client.GetEdgePlatformProxyProfile(context.Background(), d.Get("account_name").(string), d.Get("proxy_profile_name").(string))
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't get proxy configuration: %s", err)
	}

	d.Set("ip_address", proxy.IPAddress)
	d.Set("port", proxy.Port)
	d.Set("proxy_profile_id", proxy.ProxyID)
	d.Set("name", proxy.Name)
	if proxy.CaCert != nil && *proxy.CaCert != "" {
		d.Set("ca_certificate", proxy.CaCert)
	}

	fmt.Printf("read proxy profile: %v\n", proxy)

	d.SetId(proxy.ProxyID)
	return nil
}

func resourceAviatrixEdgeProxyProfileConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	err := client.DeleteEdgePlatformProxyProfile(context.Background(), d.Get("account_name").(string), d.Get("proxy_profile_name").(string))
	if err != nil {
		return fmt.Errorf("failed to delete proxy profile: %s", err)
	}

	return nil
}

func marshalEdgeProxyProfileConfigInput(d *schema.ResourceData) *goaviatrix.EdgePlatformProxyProfile {
	return &goaviatrix.EdgePlatformProxyProfile{
		Address:     d.Get("ip_address").(string),
		Port:        d.Get("port").(int),
		CACert:      d.Get("ca_certificate").(string),
		AccountName: d.Get("account_name").(string),
		Name:        d.Get("proxy_profile_name").(string),
	}
}
