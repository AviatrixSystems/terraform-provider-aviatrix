package aviatrix

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixClientProxyConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixClientProxyConfigCreate,
		Read:   resourceAviatrixClientProxyConfigRead,
		Delete: resourceAviatrixClientProxyConfigDelete,
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
				ForceNew:    true,
				Description: "Server CA Certificate local file path.",
			},
		},
	}
}

func resourceAviatrixClientProxyConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	clientProxy := marshalClientProxyConfigInput(d)
	if err := client.CreateClientProxyConfig(clientProxy); err != nil {
		return fmt.Errorf("could not config client proxy: %v", err)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixClientProxyConfigRead(d, meta)
}

func resourceAviatrixClientProxyConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return fmt.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	clientProxy, err := client.GetClientProxyConfig()
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't get client proxy configuration: %s", err)
	}

	d.Set("http_proxy", clientProxy.HttpProxy)
	d.Set("https_proxy", clientProxy.HttpsProxy)

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixClientProxyConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	err := client.DeleteClientProxyConfig()
	if err != nil {
		return fmt.Errorf("failed to delete client proxy configration: %s", err)
	}

	return nil
}

func marshalClientProxyConfigInput(d *schema.ResourceData) *goaviatrix.ClientProxyConfig {
	return &goaviatrix.ClientProxyConfig{
		HttpProxy:          d.Get("http_proxy").(string),
		HttpsProxy:         d.Get("https_proxy").(string),
		ProxyCaCertificate: d.Get("proxy_ca_certificate").(string),
	}
}
