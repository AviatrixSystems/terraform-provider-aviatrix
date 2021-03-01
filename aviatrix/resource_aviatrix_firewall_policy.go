package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixFirewallPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixFirewallPolicyCreate,
		Read:   resourceAviatrixFirewallPolicyRead,
		Delete: resourceAviatrixFirewallPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of gateway.",
			},
			"src_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "CIDRs separated by comma or tag names such 'HR' or 'marketing' etc.",
			},
			"dst_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "CIDRs separated by comma or tag names such 'HR' or 'marketing' etc.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "all",
				ForceNew:    true,
				Description: "'all', 'tcp', 'udp', 'icmp', 'sctp', 'rdp', 'dccp'.",
			},
			"port": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A single port or a range of port numbers.",
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Valid values: 'allow', 'deny' or 'force-drop'(in stateful firewall rule to allow immediate packet dropping on established sessions).",
			},
			"log_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Valid values: true or false.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				ForceNew:    true,
				Description: "Description of this firewall policy.",
			},
		},
	}
}

func getFirewallPolicyID(fw *goaviatrix.Firewall) string {
	return fmt.Sprintf("%s~%s~%s~%s~%s~%s",
		fw.GwName, fw.PolicyList[0].SrcIP, fw.PolicyList[0].DstIP, fw.PolicyList[0].Protocol, fw.PolicyList[0].Port, fw.PolicyList[0].Action)

}

func marshalFirewallPolicyInput(d *schema.ResourceData) *goaviatrix.Firewall {
	logEnabled := "on"
	if !d.Get("log_enabled").(bool) {
		logEnabled = "off"
	}

	return &goaviatrix.Firewall{
		GwName: d.Get("gw_name").(string),
		PolicyList: []*goaviatrix.Policy{
			{
				SrcIP:       d.Get("src_ip").(string),
				DstIP:       d.Get("dst_ip").(string),
				Protocol:    d.Get("protocol").(string),
				Port:        d.Get("port").(string),
				Action:      d.Get("action").(string),
				LogEnabled:  logEnabled,
				Description: d.Get("description").(string),
			},
		},
	}
}

func resourceAviatrixFirewallPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	fw := marshalFirewallPolicyInput(d)

	if err := client.AddFirewallPolicy(fw); err != nil {
		return err
	}

	d.SetId(getFirewallPolicyID(fw))
	return nil
}

func resourceAviatrixFirewallPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gwName := d.Get("gw_name").(string)
	srcIP := d.Get("src_ip").(string)
	dstIP := d.Get("dst_ip").(string)
	protocol := d.Get("protocol").(string)
	port := d.Get("port").(string)
	action := d.Get("action").(string)
	logEnabled := "on"
	if !d.Get("log_enabled").(bool) {
		logEnabled = "off"
	}
	description := d.Get("description").(string)
	if gwName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no firewall_policy received. Import Id is %s", id)

		parts := strings.Split(id, "~")
		if len(parts) != 6 {
			return fmt.Errorf("invalid firewall_policy import id: %q, "+
				"import id must be in the form gw_name~src_ip~dst_ip~protocol~port~action", id)
		}
		d.SetId(id)

		gwName, srcIP, dstIP, protocol, port, action = parts[0], parts[1], parts[2], parts[3], parts[4], parts[5]
	}

	fw := &goaviatrix.Firewall{
		GwName: gwName,
		PolicyList: []*goaviatrix.Policy{
			{
				SrcIP:       srcIP,
				DstIP:       dstIP,
				Protocol:    protocol,
				Port:        port,
				Action:      action,
				LogEnabled:  logEnabled,
				Description: description,
			},
		},
	}

	fw, err := client.GetFirewallPolicy(fw)
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	id := getFirewallPolicyID(fw)
	if err != nil {
		return fmt.Errorf("could not find firewall_policy %s: %v", id, err)
	}

	d.Set("gw_name", fw.GwName)
	d.Set("src_ip", fw.PolicyList[0].SrcIP)
	d.Set("dst_ip", fw.PolicyList[0].DstIP)
	d.Set("protocol", fw.PolicyList[0].Protocol)
	d.Set("port", fw.PolicyList[0].Port)
	d.Set("action", fw.PolicyList[0].Action)
	d.Set("log_enabled", fw.PolicyList[0].LogEnabled == "on")
	d.Set("description", fw.PolicyList[0].Description)

	d.SetId(id)
	return nil
}

func resourceAviatrixFirewallPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	fw := marshalFirewallPolicyInput(d)

	if err := client.DeleteFirewallPolicy(fw); err != nil {
		return err
	}

	d.SetId(getFirewallPolicyID(fw))
	return nil
}
