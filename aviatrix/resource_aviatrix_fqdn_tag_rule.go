package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixFQDNTagRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixFQDNTagRuleCreate,
		Read:   resourceAviatrixFQDNTagRuleRead,
		Delete: resourceAviatrixFQDNTagRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"fqdn_tag_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "FQDN Filter Tag Name to attach this domain.",
			},
			"fqdn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "FQDN.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Protocol.",
			},
			"port": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Port.",
			},
			"action": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Base Policy",
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Base Policy", "Allow", "Deny"}, false),
				Description: "What action should happen to matching requests. " +
					"Possible values are: 'Base Policy', 'Allow' or 'Deny'. Defaults to 'Base Policy' if no value is provided.",
			},
		},
	}
}

func marshalFQDNTagRuleInput(d *schema.ResourceData) *goaviatrix.FQDN {
	return &goaviatrix.FQDN{
		FQDNTag: d.Get("fqdn_tag_name").(string),
		DomainList: []*goaviatrix.Filters{
			{
				FQDN:     d.Get("fqdn").(string),
				Protocol: d.Get("protocol").(string),
				Port:     d.Get("port").(string),
				Verdict:  d.Get("action").(string),
			},
		},
	}
}

func getFQDNTagRuleID(fqdn *goaviatrix.FQDN) string {
	return fmt.Sprintf("%s~%s~%s~%s~%s",
		fqdn.FQDNTag, fqdn.DomainList[0].FQDN, fqdn.DomainList[0].Protocol, fqdn.DomainList[0].Port, fqdn.DomainList[0].Verdict)
}

func resourceAviatrixFQDNTagRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	fqdn := marshalFQDNTagRuleInput(d)

	if err := client.AddFQDNTagRule(fqdn); err != nil {
		return err
	}

	id := getFQDNTagRuleID(fqdn)
	d.SetId(id)
	return nil
}

func resourceAviatrixFQDNTagRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	fqdnTag := d.Get("fqdn_tag_name").(string)
	fqdnDomain := d.Get("fqdn").(string)
	protocol := d.Get("protocol").(string)
	port := d.Get("port").(string)
	action := d.Get("action").(string)

	if fqdnTag == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no id received. Import Id is %s", id)

		parts := strings.Split(id, "~")
		if len(parts) != 5 {
			return fmt.Errorf("invalid fqdn_tag_rule import id: %q, "+
				"import id must be in the form fqdn_tag_name~fqdn~protocol~port~action", id)
		}
		d.SetId(id)

		fqdnTag, fqdnDomain, protocol, port, action = parts[0], parts[1], parts[2], parts[3], parts[4]
	}

	fqdn := &goaviatrix.FQDN{
		FQDNTag: fqdnTag,
		DomainList: []*goaviatrix.Filters{
			{
				FQDN:     fqdnDomain,
				Protocol: protocol,
				Port:     port,
				Verdict:  action,
			},
		},
	}

	fqdn, err := client.GetFQDNTagRule(fqdn)
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not find fqdn_tag_rule %s: %v", fqdnDomain, err)
	}

	d.Set("fqdn_tag_name", fqdnTag)
	d.Set("fqdn", fqdnDomain)
	d.Set("protocol", protocol)
	d.Set("port", port)
	d.Set("action", action)

	id := getFQDNTagRuleID(fqdn)
	d.SetId(id)
	return nil
}

func resourceAviatrixFQDNTagRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	fqdn := marshalFQDNTagRuleInput(d)

	err := client.DeleteFQDNTagRule(fqdn)
	if err == goaviatrix.ErrNotFound {
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}
