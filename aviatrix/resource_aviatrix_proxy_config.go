package aviatrix

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixProxyConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixProxyConfigCreate,
		Read:   resourceAviatrixProxyConfigRead,
		Delete: resourceAviatrixProxyConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},

		Schema: map[string]*schema.Schema{
			"http_proxy": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "http proxy URL.",
			},
			"https_proxy": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "https proxy URL.",
			},
			"proxy_ca_certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Server CA Certificate file.",
			},
		},
	}
}

func resourceAviatrixProxyConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	proxy := marshalProxyConfigInput(d)
	if err := client.CreateProxyConfig(proxy); err != nil {
		return fmt.Errorf("could not config proxy: %w", err)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixProxyConfigRead(d, meta)
}

func resourceAviatrixProxyConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return fmt.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	proxy, err := client.GetProxyConfig()
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't get proxy configuration: %w", err)
	}
	mustSet(d, "http_proxy", proxy.HttpProxy)
	mustSet(d, "https_proxy", proxy.HttpsProxy)
	if proxy.ProxyCaCertificate != "" {
		mustSet(d, "proxy_ca_certificate", proxy.ProxyCaCertificate)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixProxyConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	err := client.DeleteProxyConfig()
	if err != nil {
		return fmt.Errorf("failed to delete proxy configuration: %w", err)
	}

	return nil
}

func marshalProxyConfigInput(d *schema.ResourceData) *goaviatrix.ProxyConfig {
	return &goaviatrix.ProxyConfig{
		HttpProxy:          getString(d, "http_proxy"),
		HttpsProxy:         getString(d, "https_proxy"),
		ProxyCaCertificate: getString(d, "proxy_ca_certificate"),
	}
}
