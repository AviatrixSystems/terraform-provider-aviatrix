package aviatrix

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixGatewayDNat() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixGatewayDNatCreate,
		Read:   resourceAviatrixGatewayDNatRead,
		Update: resourceAviatrixGatewayDNatUpdate,
		Delete: resourceAviatrixGatewayDNatDelete,
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
				Type:             schema.TypeList,
				Required:         true,
				DiffSuppressFunc: goaviatrix.DiffSuppressFuncGatewayDNat,
				Description:      "Policy rule to be applied to gateway.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"src_cidr": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "This is a qualifier condition that specifies a source IP address range " +
								"where the rule applies. When left blank, this field is not used.",
						},
						"src_port": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "This is a qualifier condition that specifies a source port that the rule applies. " +
								"When left blank, this field is not used.",
						},
						"dst_cidr": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "This is a qualifier condition that specifies a destination IP address range " +
								"where the rule applies. When left blank, this field is not used.",
						},
						"dst_port": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "This is a qualifier condition that specifies a destination port " +
								"where the rule applies. When left blank, this field is not used.",
						},
						"interface": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: DiffSuppressFuncNatInterface,
							Description: "This is a qualifier condition that specifies output interface where the rule applies. " +
								"When left blank, this field is not used.",
						},
						"mark": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "This is a rule field that specifies a tag or mark of a TCP session " +
								"when all qualifier conditions meet. When left blank, this field is not used.",
						},
						"dnat_ips": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "This is a rule field that specifies the translated destination IP address " +
								"when all specified qualifier conditions meet. When left blank, this field is not used. " +
								"One of the rule field must be specified for this rule to take effect.",
						},
						"dnat_port": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "This is a rule field that specifies the translated destination port " +
								"when all specified qualifier conditions meet. When left blank, this field is not used. " +
								"One of the rule field must be specified for this rule to take effect.",
						},
						"exclude_rtb": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This field specifies which VPC private route table will not be programmed with the default route entry.",
						},
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "all",
							ValidateFunc: validation.StringInSlice([]string{"all", "tcp", "udp", "icmp"}, false),
							Description: "This is a qualifier condition that specifies a destination port protocol " +
								"where the rule applies. Default: all.",
						},
						"connection": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "None",
							Description: "This is a qualifier condition that specifies output connection where the rule applies. When left blank, this field is not used.",
						},
						"apply_route_entry": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "This is an option to program the route entry 'DST CIDR pointing to Aviatrix Gateway' into Cloud platform routing table. Type: Boolean. Default: True.",
						},
					},
				},
			},
			"sync_to_ha": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to sync the policies to the HA gateway.",
			},
			"connection_policy": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Computed attribute to store the previous connection policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"src_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"src_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dst_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dst_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"interface": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mark": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dnat_ips": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dnat_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"exclude_rtb": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"apply_route_entry": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"interface_policy": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Computed attribute to store the previous interface policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"src_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"src_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dst_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dst_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"interface": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mark": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dnat_ips": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dnat_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"exclude_rtb": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"apply_route_entry": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAviatrixGatewayDNatCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gateway := &goaviatrix.Gateway{
		GatewayName: d.Get("gw_name").(string),
	}
	gateway.SyncDNATToHA = strconv.FormatBool(d.Get("sync_to_ha").(bool))

	if _, ok := d.GetOk("dnat_policy"); ok {
		policies := d.Get("dnat_policy").([]interface{})
		for _, policy := range policies {
			pl := policy.(map[string]interface{})
			customPolicy := &goaviatrix.PolicyRule{
				SrcIP:           pl["src_cidr"].(string),
				SrcPort:         pl["src_port"].(string),
				DstIP:           pl["dst_cidr"].(string),
				DstPort:         pl["dst_port"].(string),
				Protocol:        pl["protocol"].(string),
				Interface:       pl["interface"].(string),
				Connection:      pl["connection"].(string),
				Mark:            pl["mark"].(string),
				NewDstIP:        pl["dnat_ips"].(string),
				NewDstPort:      pl["dnat_port"].(string),
				ExcludeRTB:      pl["exclude_rtb"].(string),
				ApplyRouteEntry: pl["apply_route_entry"].(bool),
			}
			gateway.DnatPolicy = append(gateway.DnatPolicy, *customPolicy)
		}
	}

	d.SetId(gateway.GatewayName)
	flag := false
	defer resourceAviatrixGatewayDNatReadIfRequired(d, meta, &flag)

	err := client.UpdateDNat(gateway)
	if err != nil {
		return fmt.Errorf("failed to update DNAT for gateway(name: %s) due to: %s", gateway.GatewayName, err)
	}

	return resourceAviatrixGatewayDNatReadIfRequired(d, meta, &flag)
}

func resourceAviatrixGatewayDNatReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixGatewayDNatRead(d, meta)
	}
	return nil
}

func resourceAviatrixGatewayDNatRead(d *schema.ResourceData, meta interface{}) error {
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
			var connectionPolicy []map[string]interface{}
			var interfacePolicy []map[string]interface{}

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
				dP["dnat_ips"] = policy.NewDstIP
				dP["dnat_port"] = policy.NewDstPort
				dP["exclude_rtb"] = policy.ExcludeRTB
				dP["apply_route_entry"] = policy.ApplyRouteEntry
				dnatPolicy = append(dnatPolicy, dP)

				if policy.Connection != "None" {
					connectionPolicy = append(connectionPolicy, dP)
				}
				if policy.Interface != "" {
					interfacePolicy = append(interfacePolicy, dP)
				}
			}

			if err := d.Set("dnat_policy", dnatPolicy); err != nil {
				log.Printf("[WARN] Error setting 'dnat_policy' for (%s): %s", d.Id(), err)
			}

			if err := d.Set("connection_policy", connectionPolicy); err != nil {
				log.Printf("[WARN] Error setting 'connection_policy' for (%s): %s", d.Id(), err)
			}

			if err := d.Set("interface_policy", interfacePolicy); err != nil {
				log.Printf("[WARN] Error setting 'interface_policy' for (%s): %s", d.Id(), err)
			}

			d.Set("sync_to_ha", gwDetail.SyncDNATToHA)
		} else {
			d.SetId("")
			return nil
		}
	}

	return nil
}

func resourceAviatrixGatewayDNatUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	log.Printf("[INFO] Updating Aviatrix gateway: %#v", d.Get("gw_name").(string))

	d.Partial(true)
	gateway := &goaviatrix.Gateway{
		GatewayName: d.Get("gw_name").(string),
	}

	gateway.SyncDNATToHA = strconv.FormatBool(d.Get("sync_to_ha").(bool))

	if d.HasChange("dnat_policy") || d.HasChange("sync_to_ha") {
		if _, ok := d.GetOk("dnat_policy"); ok {
			policies := d.Get("dnat_policy").([]interface{})
			for _, policy := range policies {
				pl := policy.(map[string]interface{})
				customPolicy := &goaviatrix.PolicyRule{
					SrcIP:           pl["src_cidr"].(string),
					SrcPort:         pl["src_port"].(string),
					DstIP:           pl["dst_cidr"].(string),
					DstPort:         pl["dst_port"].(string),
					Protocol:        pl["protocol"].(string),
					Interface:       pl["interface"].(string),
					Connection:      pl["connection"].(string),
					Mark:            pl["mark"].(string),
					NewDstIP:        pl["dnat_ips"].(string),
					NewDstPort:      pl["dnat_port"].(string),
					ExcludeRTB:      pl["exclude_rtb"].(string),
					ApplyRouteEntry: pl["apply_route_entry"].(bool),
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
	d.SetId(gateway.GatewayName)
	return resourceAviatrixGatewayDNatRead(d, meta)
}

func resourceAviatrixGatewayDNatDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	gateway := &goaviatrix.Gateway{
		GatewayName:  d.Get("gw_name").(string),
		DnatPolicy:   make([]goaviatrix.PolicyRule, 0),
		SyncDNATToHA: strconv.FormatBool(d.Get("sync_to_ha").(bool)),
	}

	err := client.UpdateDNat(gateway)
	if err != nil {
		return fmt.Errorf("failed to update DNAT to nil for Aviatrix gateway(name: %s) due to: %s", gateway.GatewayName, err)
	}

	return nil
}
