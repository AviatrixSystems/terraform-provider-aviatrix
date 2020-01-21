package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixSourceNat() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSourceNatCreate,
		Read:   resourceAviatrixSourceNatRead,
		Update: resourceAviatrixSourceNatUpdate,
		Delete: resourceAviatrixSourceNatDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the gateway which is going to be created.",
			},
			"snat_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "primary",
				Description: "Valid values: 'primary', 'secondary' and 'custom'.",
			},
			"snat_policy": {
				Type:        schema.TypeList,
				Optional:    true,
				Default:     nil,
				Description: "Policy rule applied for 'snat_mode'' of 'custom'.'",
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
							Description: "A source port that the policy rule applies.",
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
						"snat_ips": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The changed source IP address when all specified qualifier conditions meet. One of the rule fields must be specified for this rule to take effect.",
						},
						"snat_port": {
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

func resourceAviatrixSourceNatCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gateway := &goaviatrix.Gateway{
		GatewayName: d.Get("gw_name").(string),
	}

	snatMode := d.Get("snat_mode").(string)
	if snatMode == "primary" {
		gateway.EnableNat = "yes"
		if len(d.Get("snat_policy").([]interface{})) != 0 {
			return fmt.Errorf("'snat_policy' should be empty for 'snat_mode' of 'primary'")
		}
	} else if snatMode == "secondary" {
		if len(d.Get("snat_policy").([]interface{})) != 0 {
			return fmt.Errorf("'snat_policy' should be empty for 'snat_mode' of 'secondary'")
		}
		gateway.EnableNat = "yes"
		gateway.SnatMode = "secondary"
	} else if snatMode == "custom" {
		if len(d.Get("snat_policy").([]interface{})) == 0 {
			return fmt.Errorf("please specify 'snat_policy' for 'snat_mode' of 'custom'")
		}
		gateway.EnableNat = "yes"
		gateway.SnatMode = "custom"
		if _, ok := d.GetOk("snat_policy"); ok {
			policies := d.Get("snat_policy").([]interface{})
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
					NewSrcIP:   pl["snat_ips"].(string),
					NewSrcPort: pl["snat_port"].(string),
					ExcludeRTB: pl["exclude_rtb"].(string),
				}
				gateway.SnatPolicy = append(gateway.SnatPolicy, *customPolicy)
			}
		}
	} else {
		return fmt.Errorf("please specify valid value for 'snat_mode'('primary', 'secondary' or 'custom')")
	}
	err := client.EnableSNat(gateway)
	if err != nil {
		return fmt.Errorf("failed to enable SNAT of mode: %s for gateway(name: %s) due to: %s", snatMode, gateway.GatewayName, err)
	}

	d.SetId(gateway.GatewayName)
	return resourceAviatrixSourceNatRead(d, meta)
}

func resourceAviatrixSourceNatRead(d *schema.ResourceData, meta interface{}) error {
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
		if gw.EnableNat == "yes" {
			if gw.SnatMode == "customized" {
				d.Set("snat_mode", "custom")
				var snatPolicy []map[string]interface{}
				for _, policy := range gwDetail.SnatPolicy {
					sP := make(map[string]interface{})
					sP["src_cidr"] = policy.SrcIP
					sP["src_port"] = policy.SrcPort
					sP["dst_cidr"] = policy.DstIP
					sP["dst_port"] = policy.DstPort
					sP["protocol"] = policy.Protocol
					sP["interface"] = policy.Interface
					sP["connection"] = policy.Connection
					sP["mark"] = policy.Mark
					sP["snat_ips"] = policy.NewSrcIP
					sP["snat_port"] = policy.NewSrcPort
					sP["exclude_rtb"] = policy.ExcludeRTB
					snatPolicy = append(snatPolicy, sP)
				}

				if err := d.Set("snat_policy", snatPolicy); err != nil {
					log.Printf("[WARN] Error setting 'snat_policy' for (%s): %s", d.Id(), err)
				}
			} else if gw.SnatMode == "secondary" {
				d.Set("snat_mode", "secondary")
				d.Set("snat_policy", nil)
			} else {
				d.Set("snat_mode", "primary")
				d.Set("snat_policy", nil)
			}
		} else {
			return fmt.Errorf("snat is not enabled for Aviatrix gateway: %s", gw.GwName)
		}
	}

	return nil
}

func resourceAviatrixSourceNatUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	log.Printf("[INFO] Updating Aviatrix gateway: %#v", d.Get("gw_name").(string))

	d.Partial(true)
	gateway := &goaviatrix.Gateway{
		GatewayName: d.Get("gw_name").(string),
	}

	if d.HasChange("snat_mode") {
		snatMode := d.Get("snat_mode").(string)
		if snatMode == "secondary" {
			if len(d.Get("snat_policy").([]interface{})) != 0 {
				return fmt.Errorf("'snat_policy' should be empty for 'snat_mode' of 'secondary'")
			}
			gateway.SnatMode = "secondary"
		} else if snatMode == "custom" {
			if len(d.Get("snat_policy").([]interface{})) == 0 {
				return fmt.Errorf("please specify 'snat_policy' for 'snat_mode' of 'custom'")
			}
			gateway.SnatMode = "custom"
			if _, ok := d.GetOk("snat_policy"); ok {
				policies := d.Get("snat_policy").([]interface{})
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
						NewSrcIP:   pl["snat_ips"].(string),
						NewSrcPort: pl["snat_port"].(string),
						ExcludeRTB: pl["exclude_rtb"].(string),
					}
					gateway.SnatPolicy = append(gateway.SnatPolicy, *customPolicy)
				}
			}
		} else if snatMode == "primary" {
			if len(d.Get("snat_policy").([]interface{})) != 0 {
				return fmt.Errorf("'snat_policy' should be empty for 'snat_mode' of 'primary'")
			}
		} else {
			return fmt.Errorf("please specify valid value for 'snat_mode'('primary', 'secondary' or 'custom')")
		}
		err := client.DisableSNat(gateway)
		if err != nil {
			return fmt.Errorf("failed to disable SNAT for gateway(name: %s) due to: %s", gateway.GatewayName, err)
		}
		err = client.EnableSNat(gateway)
		if err != nil {
			return fmt.Errorf("failed to enable SNAT of 'primary' for gateway(name: %s) due to: %s", gateway.GatewayName, err)
		}
	}

	if d.HasChange("snat_policy") {
		if !d.HasChange("snat_mode") {
			snatMode := d.Get("snat_mode").(string)
			if snatMode != "custom" {
				return fmt.Errorf("cann't update 'snat_policy' for 'snat_mode': %s", snatMode)
			}
			if len(d.Get("snat_policy").([]interface{})) == 0 {
				return fmt.Errorf("please specify 'snat_policy' for 'snat_mode' of 'custom'")
			}

			gateway.SnatMode = "custom"
			if _, ok := d.GetOk("snat_policy"); ok {
				policies := d.Get("snat_policy").([]interface{})
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
						NewSrcIP:   pl["snat_ips"].(string),
						NewSrcPort: pl["snat_port"].(string),
						ExcludeRTB: pl["exclude_rtb"].(string),
					}
					gateway.SnatPolicy = append(gateway.SnatPolicy, *customPolicy)
				}
			}

			err := client.EnableSNat(gateway)
			if err != nil {
				return fmt.Errorf("failed to enable SNAT of 'custom': %s", err)
			}
		}
	}

	d.Partial(false)
	d.SetId(d.Get("gw_name").(string))
	return resourceAviatrixSourceNatRead(d, meta)
}

func resourceAviatrixSourceNatDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	gateway := &goaviatrix.Gateway{
		GatewayName: d.Get("gw_name").(string),
	}

	err := client.DisableSNat(gateway)
	if err != nil {
		return fmt.Errorf("failed to disable SNAT for Aviatrix gateway(name: %s) due to: %s", gateway.GatewayName, err)
	}

	return nil
}
