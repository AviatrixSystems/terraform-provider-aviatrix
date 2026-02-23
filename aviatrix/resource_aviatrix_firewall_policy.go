package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixFirewallPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixFirewallPolicyCreate,
		Read:   resourceAviatrixFirewallPolicyRead,
		Delete: resourceAviatrixFirewallPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
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
			"position": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old != ""
				},
				Description: "Position in the policy list, where the firewall policy will be inserted to.",
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
	if !getBool(d, "log_enabled") {
		logEnabled = "off"
	}

	return &goaviatrix.Firewall{
		GwName: getString(d, "gw_name"),
		PolicyList: []*goaviatrix.Policy{
			{
				SrcIP:       getString(d, "src_ip"),
				DstIP:       getString(d, "dst_ip"),
				Protocol:    getString(d, "protocol"),
				Port:        getString(d, "port"),
				Action:      getString(d, "action"),
				LogEnabled:  logEnabled,
				Description: getString(d, "description"),
				Position:    getInt(d, "position"),
			},
		},
	}
}

func resourceAviatrixFirewallPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	fw := marshalFirewallPolicyInput(d)

	d.SetId(getFirewallPolicyID(fw))
	flag := false
	defer func() { _ = resourceAviatrixFirewallPolicyReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path
	if fw.PolicyList[0].Position == 0 {
		if err := client.AddFirewallPolicy(fw); err != nil {
			return err
		}
	} else {
		if err := client.InsertFirewallPolicy(fw); err != nil {
			return err
		}
	}

	return resourceAviatrixFirewallPolicyReadIfRequired(d, meta, &flag)
}

func resourceAviatrixFirewallPolicyReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixFirewallPolicyRead(d, meta)
	}
	return nil
}

func resourceAviatrixFirewallPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	gwName := getString(d, "gw_name")
	srcIP := getString(d, "src_ip")
	dstIP := getString(d, "dst_ip")
	protocol := getString(d, "protocol")
	port := getString(d, "port")
	action := getString(d, "action")
	logEnabled := "on"
	if !getBool(d, "log_enabled") {
		logEnabled = "off"
	}
	description := getString(d, "description")
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
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return err
	}
	id := getFirewallPolicyID(fw)
	mustSet(d, "gw_name", fw.GwName)
	mustSet(d, "src_ip", fw.PolicyList[0].SrcIP)
	mustSet(d, "dst_ip", fw.PolicyList[0].DstIP)
	mustSet(d, "protocol", fw.PolicyList[0].Protocol)
	mustSet(d, "port", fw.PolicyList[0].Port)
	mustSet(d, "action", fw.PolicyList[0].Action)
	mustSet(d, "log_enabled", fw.PolicyList[0].LogEnabled == "on")
	mustSet(d, "description", fw.PolicyList[0].Description)
	mustSet(d, "position", fw.PolicyList[0].Position)

	d.SetId(id)
	return nil
}

func resourceAviatrixFirewallPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	fw := marshalFirewallPolicyInput(d)

	if err := client.DeleteFirewallPolicy(fw); err != nil {
		return err
	}

	d.SetId(getFirewallPolicyID(fw))
	return nil
}
