package aviatrix

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixGatewaySNat() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixGatewaySNatCreate,
		Read:   resourceAviatrixGatewaySNatRead,
		Update: resourceAviatrixGatewaySNatUpdate,
		Delete: resourceAviatrixGatewaySNatDelete,
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
			"snat_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "customized_snat",
				ValidateFunc: validation.StringInSlice([]string{"customized_snat"}, false),
				Description:  "Nat mode. Currently only supports 'customized_snat'.",
			},
			"snat_policy": {
				Type:        schema.TypeList,
				Optional:    true,
				Default:     nil,
				Description: "Policy rules applied for 'snat_mode'' of 'customized_snat'.'",
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
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "This is a qualifier condition that specifies a destination port protocol " +
								"where the rule applies. When left blank, this field is not used.",
						},
						"interface": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "eth0",
							Description: "This is a qualifier condition that specifies output interface " +
								"where the rule applies. When left blank, this field is not used. Default value: 'eth0'. " +
								"Empty string is not a valid value.",
						},
						"connection": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "None",
							Description: "None.",
						},
						"mark": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "This is a qualifier condition that specifies a tag or mark of a TCP session " +
								"where the rule applies. When left blank, this field is not used.",
						},
						"snat_ips": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "This is a rule field that specifies the changed source IP address " +
								"when all specified qualifier conditions meet. When left blank, this field is not used. " +
								"One of the rule fields must be specified for this rule to take effect.",
						},
						"snat_port": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "This is a rule field that specifies the changed source port " +
								"when all specified qualifier conditions meet. When left blank, this field is not used. " +
								"One of the rule fields must be specified for this rule to take effect.",
						},
						"exclude_rtb": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This field specifies which VPC private route table will not be programmed with the default route entry.",
						},
					},
				},
			},
			"sync_to_ha": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to sync the policies to the HA gateway.",
			},
		},
	}
}

func resourceAviatrixGatewaySNatCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gateway := &goaviatrix.Gateway{
		GatewayName: d.Get("gw_name").(string),
	}

	if len(d.Get("snat_policy").([]interface{})) == 0 {
		return fmt.Errorf("please specify 'snat_policy' for 'snat_mode' of 'customized_snat'")
	}
	gateway.EnableNat = "yes"
	gateway.SnatMode = "custom"
	gateway.SyncSNATToHA = strconv.FormatBool(d.Get("sync_to_ha").(bool))
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
		return fmt.Errorf("failed to configure policies for 'customized_snat' mode due to: %s", err)
	}

	d.SetId(gateway.GatewayName)
	return resourceAviatrixGatewaySNatRead(d, meta)
}

func resourceAviatrixGatewaySNatRead(d *schema.ResourceData, meta interface{}) error {
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
		if gw.NatEnabled && gw.SnatMode == "customized" {
			d.Set("snat_mode", "customized_snat")
			var snatPolicy []map[string]interface{}

			// Duplicate SNAT policies can be returned from the API.
			// Before we save the policies to state we need to deduplicate.
			dedupMap := make(map[string]struct{})
			for _, policy := range gwDetail.SnatPolicy {
				sP := make(map[string]interface{})
				sP["src_cidr"] = policy.SrcIP
				sP["src_port"] = policy.SrcPort
				sP["dst_cidr"] = policy.DstIP
				sP["dst_port"] = policy.DstPort
				sP["protocol"] = policy.Protocol
				sP["interface"] = "eth0"
				sP["connection"] = policy.Connection
				sP["mark"] = policy.Mark
				sP["snat_ips"] = policy.NewSrcIP
				sP["snat_port"] = policy.NewSrcPort
				sP["exclude_rtb"] = policy.ExcludeRTB

				// To deduplicate we will generate a unique key for each policy.
				key := fmt.Sprintf("%s~%s~%s~%s~%s~%s~%s~%s~%s~%s~%s", policy.SrcIP, policy.SrcPort, policy.DstIP,
					policy.DstPort, policy.Protocol, "eth0", policy.Connection, policy.Mark, policy.NewSrcIP, policy.NewSrcPort, policy.ExcludeRTB)
				// If the map already contains the unique key then we know this policy is a duplicate.
				if _, ok := dedupMap[key]; ok {
					continue
				}

				// Otherwise, its a unique policy so we write it to state and the dedupMap.
				dedupMap[key] = struct{}{}
				snatPolicy = append(snatPolicy, sP)
			}

			if err := d.Set("snat_policy", snatPolicy); err != nil {
				log.Printf("[WARN] Error setting 'snat_policy' for (%s): %s", d.Id(), err)
			}

			d.Set("sync_to_ha", gwDetail.SyncSNATToHA)
		} else {
			d.SetId("")
			return nil
		}
	}

	return nil
}

func resourceAviatrixGatewaySNatUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	log.Printf("[INFO] Updating Aviatrix gateway: %#v", d.Get("gw_name").(string))

	d.Partial(true)
	gateway := &goaviatrix.Gateway{
		GatewayName: d.Get("gw_name").(string),
	}

	gateway.SyncSNATToHA = strconv.FormatBool(d.Get("sync_to_ha").(bool))

	if d.HasChange("snat_policy") || d.HasChange("sync_to_ha") {
		if len(d.Get("snat_policy").([]interface{})) == 0 {
			return fmt.Errorf("please specify 'snat_policy' for 'snat_mode' of 'customized_snat'")
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
			return fmt.Errorf("failed to enable SNAT of 'customized_snat': %s", err)
		}
	}

	d.Partial(false)
	d.SetId(d.Get("gw_name").(string))
	return resourceAviatrixGatewaySNatRead(d, meta)
}

func resourceAviatrixGatewaySNatDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	gateway := &goaviatrix.Gateway{
		GatewayName:  d.Get("gw_name").(string),
		SnatMode:     "custom",
		SyncSNATToHA: strconv.FormatBool(d.Get("sync_to_ha").(bool)),
	}

	err := client.DisableCustomSNat(gateway)
	if err != nil {
		return fmt.Errorf("failed to disable SNAT for Aviatrix gateway(name: %s) due to: %s", gateway.GatewayName, err)
	}

	return nil
}
