package aviatrix

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixRemoteSyslog() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixRemoteSyslogCreate,
		Read:   resourceAviatrixRemoteSyslogRead,
		Update: resourceAviatrixRemoteSyslogUpdate,
		Delete: resourceAviatrixRemoteSyslogDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"index": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 9),
				Description:  "Profile index: a total of 10 profiles from index 0 to 9 are supported for remote syslog, while index 9 is reserved for CoPilot.",
			},
			"server": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Server: FQDN or IP address of the remote syslog server",
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Port: Listening port of the remote syslog server (6514 by default)",
			},
			"ca_certificate_file_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CA Certificate: Certificate Authority (CA) certificate",
			},
			"public_certificate_file_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Server Public Certificate: Public certificate of the controller signed by the same CA",
			},
			"private_key_file_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Server Private Key: Private key of the controller that pairs with the public certificate",
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "TCP",
				ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP"}, false),
				Description:  "Protocol: TCP or UDP (TCP by default)",
			},
			"template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional Custom Template: Useful when forwarding to 3rd party servers like Datadog or Sumo",
			},
			"exclude_gateway": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of excluded gateways.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"notls": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "True if not protected by TLS.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Enabled or not.",
			},
		},
	}
}

func resourceAviatrixRemoteSyslogCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	remoteSyslog := &goaviatrix.RemoteSyslog{
		Server:            d.Get("server").(string),
		Port:              d.Get("port").(int),
		Protocol:          d.Get("protocol").(string),
		Index:             d.Get("index").(int),
		Template:          d.Get("template").(string),
		CaCertificate:     d.Get("ca_certificate_file_path").(string),
		PublicCertificate: d.Get("public_certificate_file_path").(string),
		PrivateKey:        d.Get("private_key_file_path").(string),
	}

	if !((remoteSyslog.CaCertificate != "" && remoteSyslog.PublicCertificate != "" && remoteSyslog.PrivateKey != "") ||
		(remoteSyslog.CaCertificate == "" && remoteSyslog.PublicCertificate == "" && remoteSyslog.PrivateKey == "")) {
		return fmt.Errorf("one or more certificates missing")
	}

	var excludeGateways []string
	for _, v := range d.Get("exclude_gateway").(*schema.Set).List() {
		excludeGateways = append(excludeGateways, v.(string))
	}
	if len(excludeGateways) != 0 {
		remoteSyslog.ExcludeGatewayInput = strings.Join(excludeGateways, ",")
	}

	if err := client.EnableRemoteSyslog(remoteSyslog); err != nil {
		return fmt.Errorf("could not enable remote syslog: %v", err)
	}

	d.SetId("remote_syslog_" + strconv.Itoa(remoteSyslog.Index))
	return nil
}

func resourceAviatrixRemoteSyslogRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	remoteSyslogStatus, err := client.GetRemoteSyslogStatus(d.Get("index").(int))
	if err != nil {
		return fmt.Errorf("could not get remote syslog status: %v", err)
	}

	d.Set("notls", remoteSyslogStatus.Notls)
	d.Set("status", remoteSyslogStatus.Status)

	var excludedGateways []interface{}
	for _, v := range remoteSyslogStatus.ExcludedGateway {
		excludedGateways = append(excludedGateways, v)
	}
	if err := d.Set("exclude_gateway", excludedGateways); err != nil {
		return fmt.Errorf("could not set excluded_gateway: %v", err)
	}

	d.SetId("remote_syslog_" + remoteSyslogStatus.Index)
	return nil
}

func resourceAviatrixRemoteSyslogUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAviatrixRemoteSyslogCreate(d, meta)
}

func resourceAviatrixRemoteSyslogDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	if err := client.DisableRemoteSyslog(d.Get("index").(int)); err != nil {
		return fmt.Errorf("could not disable remote syslog: %v", err)
	}

	d.SetId("")
	return nil
}
