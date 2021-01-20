package aviatrix

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixProxyConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixProxyConfigCreate,
		Read:   resourceAviatrixProxyConfigRead,
		Delete: resourceAviatrixProxyConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				Sensitive:   true,
				ForceNew:    true,
				Description: "Server CA Certificate file.",
			},
		},
	}
}

func resourceAviatrixProxyConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	proxy := marshalProxyConfigInput(d)
	if err := client.CreateProxyConfig(proxy); err != nil {
		return fmt.Errorf("could not config proxy: %v", err)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixProxyConfigRead(d, meta)
}

func resourceAviatrixProxyConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return fmt.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	proxy, err := client.GetProxyConfig()
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't get proxy configuration: %s", err)
	}

	d.Set("http_proxy", proxy.HttpProxy)
	d.Set("https_proxy", proxy.HttpsProxy)

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixProxyConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	err := client.DeleteProxyConfig()
	if err != nil {
		return fmt.Errorf("failed to delete proxy configration: %s", err)
	}

	return nil
}

func marshalProxyConfigInput(d *schema.ResourceData) *goaviatrix.ProxyConfig {
	return &goaviatrix.ProxyConfig{
		HttpProxy:          d.Get("http_proxy").(string),
		HttpsProxy:         d.Get("https_proxy").(string),
		ProxyCaCertificate: d.Get("proxy_ca_certificate").(string),
	}
}
