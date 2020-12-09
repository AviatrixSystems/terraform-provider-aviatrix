package aviatrix

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

var remoteSyslogMatcher = regexp.MustCompile(`\bremote_syslog_[0-9]\b`)

func resourceAviatrixRemoteSyslog() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixRemoteSyslogCreate,
		Read:   resourceAviatrixRemoteSyslogRead,
		Delete: resourceAviatrixRemoteSyslogDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"index": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 9),
				Description:  "Profile index: a total of 10 profiles from index 0 to 9 are supported for remote syslog, while index 9 is reserved for CoPilot.",
			},
			"server": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Server: FQDN or IP address of the remote syslog server",
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Port: Listening port of the remote syslog server (6514 by default)",
			},
			"ca_certificate_file_path": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "CA Certificate: Certificate Authority (CA) certificate",
			},
			"public_certificate_file_path": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Server Public Certificate: Public certificate of the controller signed by the same CA",
			},
			"private_key_file_path": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Server Private Key: Private key of the controller that pairs with the public certificate",
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "TCP",
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP"}, false),
				Description:  "Protocol: TCP or UDP (TCP by default)",
			},
			"template": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional Custom Template: Useful when forwarding to 3rd party servers like Datadog or Sumo",
			},
			"excluded_gateways": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
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

func marshalRemoteSyslogInput(d *schema.ResourceData) *goaviatrix.RemoteSyslog {
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

	var excludeGateways []string
	for _, v := range d.Get("excluded_gateways").(*schema.Set).List() {
		excludeGateways = append(excludeGateways, v.(string))
	}
	if len(excludeGateways) != 0 {
		remoteSyslog.ExcludeGatewayInput = strings.Join(excludeGateways, ",")
	}

	return remoteSyslog
}

func resourceAviatrixRemoteSyslogCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	remoteSyslog := marshalRemoteSyslogInput(d)

	if !((remoteSyslog.CaCertificate != "" && remoteSyslog.PublicCertificate != "" && remoteSyslog.PrivateKey != "") ||
		(remoteSyslog.CaCertificate == "" && remoteSyslog.PublicCertificate == "" && remoteSyslog.PrivateKey == "")) {
		return fmt.Errorf("one or more certificates missing")
	}

	if err := client.EnableRemoteSyslog(remoteSyslog); err != nil {
		return fmt.Errorf("could not enable remote syslog: %v", err)
	}

	d.SetId("remote_syslog_" + strconv.Itoa(remoteSyslog.Index))
	return resourceAviatrixRemoteSyslogRead(d, meta)
}

func resourceAviatrixRemoteSyslogRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	server := d.Get("server").(string)

	if server == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)

		match := remoteSyslogMatcher.Match([]byte(id))
		if !match {
			return fmt.Errorf("invalid ID format, expected ID in format \"remote_syslog_{index}\", instead got %s", id)
		}

		index, _ := strconv.Atoi(id[len(id)-1:])
		d.Set("index", index)
		d.SetId(id)
	}

	remoteSyslogStatus, err := client.GetRemoteSyslogStatus(d.Get("index").(int))
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not get remote syslog status: %v", err)
	}

	idx, _ := strconv.Atoi(remoteSyslogStatus.Index)
	d.Set("index", idx)
	d.Set("server", remoteSyslogStatus.Server)
	port, _ := strconv.Atoi(remoteSyslogStatus.Port)
	d.Set("port", port)
	d.Set("protocol", remoteSyslogStatus.Protocol)
	d.Set("template", remoteSyslogStatus.Template)
	d.Set("notls", remoteSyslogStatus.Notls)
	d.Set("status", remoteSyslogStatus.Status)
	d.Set("excluded_gateways", remoteSyslogStatus.ExcludedGateways)

	d.SetId("remote_syslog_" + remoteSyslogStatus.Index)
	return nil
}

func resourceAviatrixRemoteSyslogDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	if err := client.DisableRemoteSyslog(d.Get("index").(int)); err != nil {
		return fmt.Errorf("could not disable remote syslog: %v", err)
	}

	return nil
}
