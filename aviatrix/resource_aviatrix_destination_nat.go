package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixDestinationNat() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixDestinationNatCreate,
		Read:   resourceAviatrixDestinationNatRead,
		Update: resourceAviatrixDestinationNatUpdate,
		Delete: resourceAviatrixDestinationNatDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the gateway.",
			},
			"dnat_policy": {
				Type:        schema.TypeList,
				Required:    true,
				Default:     nil,
				Description: "Policy rule to be applied to gateway.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"src_cidr": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A source IP address range where the policy rule applies.",
						},
						"src_port": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A source port where the policy rule applies.",
						},
						"dst_cidr": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A destination IP address range where the policy rule applies.",
						},
						"dst_port": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A destination port where the policy rule applies.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A destination port protocol where the policy rule applies.",
						},
						"interface": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "An output interface where the policy rule applies.",
						},
						"connection": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "None",
							Description: "None.",
						},
						"mark": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A tag or mark of a TCP session where the policy rule applies.",
						},
						"dnat_ips": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The changed source IP address when all specified qualifier conditions meet. One of the rule fields must be specified for this rule to take effect.",
						},
						"dnat_port": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The translated destination port when all specified qualifier conditions meet. One of the rule field must be specified for this rule to take effect.",
						},
						"exclude_rtb": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This field specifies which VPC private route table will not be programmed with the default route entry.",
						},
					},
				},
			},
		},
	}
}

func resourceAviatrixDestinationNatCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gateway := &goaviatrix.Gateway{
		GatewayName: d.Get("gw_name").(string),
	}

	if _, ok := d.GetOk("dnat_policy"); ok {
		policies := d.Get("dnat_policy").([]interface{})
		for _, policy := range policies {
			pl := policy.(map[string]interface{})
			customPolicy := &goaviatrix.PolicyRule{
				SrcIP:      pl["src_cidr"].(string),
				SrcPort:    pl["src_port"].(string),
				DstIP:      pl["dst_cidr"].(string),
				DstPort:    pl["dst_port"].(string),
				Protocol:   pl["protocol"].(string),
				Interface:  pl["interface"].(string),
				Connection: pl["connection"].(string),
				Mark:       pl["mark"].(string),
				NewSrcIP:   pl["dnat_ips"].(string),
				NewSrcPort: pl["dnat_port"].(string),
				ExcludeRTB: pl["exclude_rtb"].(string),
			}
			gateway.DnatPolicy = append(gateway.DnatPolicy, *customPolicy)
		}
	}
	err := client.UpdateDNat(gateway)
	if err != nil {
		return fmt.Errorf("failed to update DNAT for gateway(name: %s) due to: %s", gateway.GatewayName, err)
	}

	d.SetId(gateway.GatewayName)
	return resourceAviatrixDestinationNatRead(d, meta)
}

func resourceAviatrixDestinationNatRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gwName := d.Get("gw_name").(string)
	if gwName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no gateway name received. Import Id is %s", id)
		d.Set("gw_name", id)
		d.SetId(id)
	}

	gateway := &goaviatrix.Gateway{
		GwName: d.Get("gw_name").(string),
	}

	gw, err := client.GetGateway(gateway)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix gateway: %s", err)
	}

	log.Printf("[TRACE] reading gateway %s: %#v", d.Get("gw_name").(string), gw)
	if gw != nil {
		d.Set("gw_name", gw.GwName)

		gwDetail, err := client.GetGatewayDetail(gateway)
		if err != nil {
			return fmt.Errorf("couldn't get detail information of Aviatrix gateway(name: %s) due to: %s", gw.GwName, err)
		}
		if len(gwDetail.DnatPolicy) != 0 {
			var dnatPolicy []map[string]interface{}
			for _, policy := range gwDetail.DnatPolicy {
				dP := make(map[string]interface{})
				dP["src_cidr"] = policy.SrcIP
				dP["src_port"] = policy.SrcPort
				dP["dst_cidr"] = policy.DstIP
				dP["dst_port"] = policy.DstPort
				dP["protocol"] = policy.Protocol
				dP["interface"] = policy.Interface
				dP["connection"] = policy.Connection
				dP["mark"] = policy.Mark
				dP["dnat_ips"] = policy.NewSrcIP
				dP["dnat_port"] = policy.NewSrcPort
				dP["exclude_rtb"] = policy.ExcludeRTB
				dnatPolicy = append(dnatPolicy, dP)
			}

			if err := d.Set("dnat_policy", dnatPolicy); err != nil {
				log.Printf("[WARN] Error setting 'dnat_policy' for (%s): %s", d.Id(), err)
			}
		} else {
			d.SetId("")
			return nil
		}
	}

	return nil
}

func resourceAviatrixDestinationNatUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	log.Printf("[INFO] Updating Aviatrix gateway: %#v", d.Get("gw_name").(string))

	d.Partial(true)
	gateway := &goaviatrix.Gateway{
		GatewayName: d.Get("gw_name").(string),
	}

	if d.HasChange("dnat_policy") {
		if _, ok := d.GetOk("dnat_policy"); ok {
			policies := d.Get("dnat_policy").([]interface{})
			for _, policy := range policies {
				pl := policy.(map[string]interface{})
				customPolicy := &goaviatrix.PolicyRule{
					SrcIP:      pl["src_cidr"].(string),
					SrcPort:    pl["src_port"].(string),
					DstIP:      pl["dst_cidr"].(string),
					DstPort:    pl["dst_port"].(string),
					Protocol:   pl["protocol"].(string),
					Interface:  pl["interface"].(string),
					Connection: pl["connection"].(string),
					Mark:       pl["mark"].(string),
					NewSrcIP:   pl["dnat_ips"].(string),
					NewSrcPort: pl["dnat_port"].(string),
					ExcludeRTB: pl["exclude_rtb"].(string),
				}
				gateway.DnatPolicy = append(gateway.DnatPolicy, *customPolicy)
			}
		}
		err := client.UpdateDNat(gateway)
		if err != nil {
			return fmt.Errorf("failed to update DNAT for gateway(name: %s) due to: %s", gateway.GatewayName, err)
		}
	}

	d.Partial(false)
	d.SetId(d.Get("gw_name").(string))
	return resourceAviatrixDestinationNatRead(d, meta)
}

func resourceAviatrixDestinationNatDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	gateway := &goaviatrix.Gateway{
		GatewayName: d.Get("gw_name").(string),
		DnatPolicy:  make([]goaviatrix.PolicyRule, 0),
	}

	err := client.UpdateDNat(gateway)
	if err != nil {
		return fmt.Errorf("failed to update DNAT to nil for Aviatrix gateway(name: %s) due to: %s", gateway.GatewayName, err)
	}

	return nil
}
