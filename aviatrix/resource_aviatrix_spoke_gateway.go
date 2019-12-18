package aviatrix

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixSpokeGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixSpokeGatewayCreate,
		Read:   resourceAviatrixSpokeGatewayRead,
		Update: resourceAviatrixSpokeGatewayUpdate,
		Delete: resourceAviatrixSpokeGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cloud_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Type of cloud service provider.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This parameter represents the name of a Cloud-Account in Aviatrix controller.",
			},
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the gateway which is going to be created.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC-ID/VNet-Name of cloud provider.",
			},
			"vpc_reg": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Region of cloud provider.",
			},
			"gw_size": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Size of the gateway instance.",
			},
			"subnet": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Public Subnet Info.",
			},
			"insane_mode_az": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "AZ of subnet being created for Insane Mode Spoke Gateway. Required if insane_mode is enabled for aws cloud.",
			},
			"enable_snat": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Specify whether enabling Source NAT feature on the gateway or not.",
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
						"src_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A source IP address range where the policy rule applies.",
						},
						"src_port": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A source port that the policy rule applies.",
						},
						"dst_ip": {
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
						"new_src_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The changed source IP address when all specified qualifier conditions meet. One of the rule fields must be specified for this rule to take effect.",
						},
						"new_src_port": {
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
			"dnat_policy": {
				Type:        schema.TypeList,
				Optional:    true,
				Default:     nil,
				Description: "Policy rule applied for enabling Destination NAT (DNAT), which allows you to change the destination to a virtual address range.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"src_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A source IP address range where the policy rule applies.",
						},
						"src_port": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A source port that the policy rule applies.",
						},
						"dst_ip": {
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
						"new_src_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The changed source IP address when all specified qualifier conditions meet. One of the rule fields must be specified for this rule to take effect.",
						},
						"new_src_port": {
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
			"allocate_new_eip": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				Description: "If false, reuse an idle address in Elastic IP pool for this gateway. " +
					"Otherwise, allocate a new Elastic IP and use it for this gateway.",
			},
			"eip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Required when allocate_new_eip is false. It uses specified EIP for this gateway.",
			},
			"ha_subnet": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "HA Subnet. Required if enabling HA for AWS/ARM.",
			},
			"ha_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "HA Zone. Required if enabling HA for GCP.",
			},
			"ha_insane_mode_az": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "AZ of subnet being created for Insane Mode Spoke HA Gateway. Required if insane_mode is true and ha_subnet is set.",
			},
			"ha_gw_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "HA Gateway Size.",
			},
			"ha_eip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Public IP address that you want assigned to the HA Spoke Gateway.",
			},
			"single_az_ha": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Set to 'enabled' if this feature is desired.",
			},
			"transit_gw": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Specify the transit Gateway.",
			},
			"tag_list": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Default:     nil,
				Description: "Instance tag of cloud provider.",
			},
			"insane_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable Insane Mode for Spoke Gateway. Valid values: true, false. If insane mode is enabled, gateway size has to at least be c5 size.",
			},
			"enable_active_mesh": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Switch to Enable/Disable Active Mesh Mode for Spoke Gateway. Valid values: true, false.",
			},
			"enable_vpc_dns_server": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable vpc_dns_server for Gateway. Only supports AWS. Valid values: true, false.",
			},
			"enable_encrypt_volume": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable encrypt gateway EBS volume. Only supported for AWS provider. Valid values: true, false. Default value: false.",
			},
			"customer_managed_keys": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Customer managed key ID.",
			},
			"cloud_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud instance ID.",
			},
		},
	}
}

func resourceAviatrixSpokeGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gateway := &goaviatrix.SpokeVpc{
		CloudType:      d.Get("cloud_type").(int),
		AccountName:    d.Get("account_name").(string),
		GwName:         d.Get("gw_name").(string),
		VpcSize:        d.Get("gw_size").(string),
		Subnet:         d.Get("subnet").(string),
		HASubnet:       d.Get("ha_subnet").(string),
		TransitGateway: d.Get("transit_gw").(string),
	}
	enableSNat := d.Get("enable_snat").(bool)
	snatMode := d.Get("snat_mode").(string)
	if enableSNat && snatMode == "primary" {
		if len(d.Get("snat_policy").([]interface{})) != 0 {
			return fmt.Errorf("'snat_policy' should be empty for 'snat_mode' of 'primary'")
		}
		gateway.EnableNat = "yes"
	}

	singleAZ := d.Get("single_az_ha").(bool)
	if singleAZ {
		gateway.SingleAzHa = "enabled"
	} else {
		gateway.SingleAzHa = "disabled"
	}

	allocateNewEip := d.Get("allocate_new_eip").(bool)
	if allocateNewEip {
		gateway.ReuseEip = "off"
	} else {
		gateway.ReuseEip = "on"
		gateway.Eip = d.Get("eip").(string)
	}

	if gateway.CloudType == 1 || gateway.CloudType == 4 || gateway.CloudType == 16 {
		gateway.VpcID = d.Get("vpc_id").(string)
		if gateway.VpcID == "" {
			return fmt.Errorf("'vpc_id' cannot be empty for creating a spoke gw")
		}
	} else if gateway.CloudType == 8 {
		gateway.VNetNameResourceGroup = d.Get("vpc_id").(string)
		if gateway.VNetNameResourceGroup == "" {
			return fmt.Errorf("'vpc_id' cannot be empty for creating a spoke gw")
		}
	} else {
		return fmt.Errorf("invalid cloud type, it can only be aws (1), gcp (4), arm (8)")
	}

	if gateway.CloudType == 1 || gateway.CloudType == 8 || gateway.CloudType == 16 {
		gateway.VpcRegion = d.Get("vpc_reg").(string)
	} else if gateway.CloudType == 4 {
		// for gcp, rest api asks for "zone" rather than vpc region
		gateway.Zone = d.Get("vpc_reg").(string)
	} else {
		return fmt.Errorf("invalid cloud type, it can only be AWS (1), GCP (4), or ARM (8)")
	}

	insaneMode := d.Get("insane_mode").(bool)
	insaneModeAz := d.Get("insane_mode_az").(string)
	haZone := d.Get("ha_zone").(string)
	haSubnet := d.Get("ha_subnet").(string)
	haGwSize := d.Get("ha_gw_size").(string)
	haInsaneModeAz := d.Get("ha_insane_mode_az").(string)
	if insaneMode {
		if gateway.CloudType != 1 && gateway.CloudType != 8 {
			return fmt.Errorf("insane_mode is only supported for aws and arm (cloud_type = 1 or 8)")
		}
		if gateway.CloudType == 1 {
			if insaneModeAz == "" {
				return fmt.Errorf("insane_mode_az needed if insane_mode is enabled for aws cloud")
			}
			if haSubnet != "" && haInsaneModeAz == "" {
				return fmt.Errorf("ha_insane_mode_az needed if insane_mode is enabled for aws cloud and ha_subnet is set")
			}
			// Append availability zone to subnet
			var strs []string
			strs = append(strs, gateway.Subnet, insaneModeAz)
			gateway.Subnet = strings.Join(strs, "~~")
		}
		gateway.InsaneMode = "on"
	} else {
		gateway.InsaneMode = "off"
	}
	if haZone != "" || haSubnet != "" {
		if haGwSize == "" {
			return fmt.Errorf("A valid non empty ha_gw_size parameter is mandatory for this resource if " +
				"ha_subnet or ha_zone is set. Example: t2.micro")
		}
	}

	enableEncryptVolume := d.Get("enable_encrypt_volume").(bool)
	customerManagedKeys := d.Get("customer_managed_keys").(string)
	if enableEncryptVolume && d.Get("cloud_type").(int) != 1 {
		return fmt.Errorf("'enable_encrypt_volume' is only supported for AWS provider (cloud_type: 1)")
	}
	if !enableEncryptVolume && customerManagedKeys != "" {
		return fmt.Errorf("'customer_managed_keys' should be empty since Encrypt Volume is not enabled")
	}
	if enableEncryptVolume && d.Get("single_az_ha").(bool) {
		return fmt.Errorf("'single_az_ha' needs to be disabled to encrypt gateway EBS volume")
	}

	log.Printf("[INFO] Creating Aviatrix Spoke Gateway: %#v", gateway)

	err := client.LaunchSpokeVpc(gateway)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix Spoke Gateway: %s", err)
	}

	d.SetId(gateway.GwName)

	flag := false
	defer resourceAviatrixSpokeGatewayReadIfRequired(d, meta, &flag)

	if enableActiveMesh := d.Get("enable_active_mesh").(bool); !enableActiveMesh {
		gw := &goaviatrix.Gateway{
			GwName: d.Get("gw_name").(string),
		}
		gw.EnableActiveMesh = "no"

		err := client.DisableActiveMesh(gw)
		if err != nil {
			return fmt.Errorf("couldn't disable Active Mode for Aviatrix Spoke Gateway: %s", err)
		}
	}

	if !singleAZ {
		singleAZGateway := &goaviatrix.Gateway{
			GwName:   d.Get("gw_name").(string),
			SingleAZ: "disabled",
		}

		log.Printf("[INFO] Disable Single AZ GW HA: %#v", singleAZGateway)

		err := client.DisableSingleAZGateway(singleAZGateway)
		if err != nil {
			return fmt.Errorf("failed to disable single AZ GW HA: %s", err)
		}
	}

	if haSubnet != "" || haZone != "" {
		//Enable HA
		haGateway := &goaviatrix.SpokeVpc{
			GwName:    d.Get("gw_name").(string),
			CloudType: d.Get("cloud_type").(int),
			HASubnet:  haSubnet,
			HAZone:    haZone,
		}

		haGateway.Eip = d.Get("ha_eip").(string)

		if insaneMode == true && haGateway.CloudType == 1 {
			var haStrs []string
			haStrs = append(haStrs, haSubnet, haInsaneModeAz)
			haSubnet = strings.Join(haStrs, "~~")
			haGateway.HASubnet = haSubnet
		}

		err = client.EnableHaSpokeVpc(haGateway)
		if err != nil {
			return fmt.Errorf("failed to enable HA Aviatrix Spoke Gateway: %s", err)
		}

		log.Printf("[INFO]Resizing Spoke HA Gateway: %#v", haGwSize)

		if haGwSize != gateway.VpcSize {
			if haGwSize == "" {
				return fmt.Errorf("A valid non empty ha_gw_size parameter is mandatory for this resource if " +
					"ha_subnet or ha_zone is set. Example: t2.micro us-west1-b")
			}

			haGateway := &goaviatrix.Gateway{
				CloudType: d.Get("cloud_type").(int),
				GwName:    d.Get("gw_name").(string) + "-hagw",
			}

			haGateway.GwSize = d.Get("ha_gw_size").(string)

			log.Printf("[INFO] Resizing Spoke HA Gateway size to: %s ", haGateway.GwSize)

			err := client.UpdateGateway(haGateway)
			log.Printf("[INFO] Resizing Spoke HA Gateway size to: %s ", haGateway.GwSize)
			if err != nil {
				return fmt.Errorf("failed to update Aviatrix Spoke HA Gateway size: %s", err)
			}

			d.Set("ha_gw_size", haGwSize)
		}
	}

	if _, ok := d.GetOk("tag_list"); ok && gateway.CloudType == 1 {
		tagList := d.Get("tag_list").([]interface{})
		tagListStr := goaviatrix.ExpandStringList(tagList)
		tagListStr = goaviatrix.TagListStrColon(tagListStr)
		gateway.TagList = strings.Join(tagListStr, ",")
		tags := &goaviatrix.Tags{
			CloudType:    1,
			ResourceType: "gw",
			ResourceName: d.Get("gw_name").(string),
			TagList:      gateway.TagList,
		}
		err = client.AddTags(tags)
		if err != nil {
			return fmt.Errorf("failed to add tags: %s", err)
		}
	} else if ok && gateway.CloudType != 1 {
		return fmt.Errorf("adding tags only supported for aws, cloud_type must be 1")
	}

	if transitGwName := d.Get("transit_gw").(string); transitGwName != "" {
		//No HA config, just return
		err := client.SpokeJoinTransit(gateway)
		if err != nil {
			return fmt.Errorf("failed to join Transit Gateway: %s", err)
		}
	}

	enableVpcDnsServer := d.Get("enable_vpc_dns_server").(bool)
	if d.Get("cloud_type").(int) == 1 && enableVpcDnsServer {
		gwVpcDnsServer := &goaviatrix.Gateway{
			GwName: d.Get("gw_name").(string),
		}

		log.Printf("[INFO] Enable VPC DNS Server: %#v", gwVpcDnsServer)

		err := client.EnableVpcDnsServer(gwVpcDnsServer)
		if err != nil {
			return fmt.Errorf("failed to enable VPC DNS Server: %s", err)
		}
	} else if enableVpcDnsServer {
		return fmt.Errorf("'enable_vpc_dns_server' only supports AWS(1)")
	}

	if enableEncryptVolume {
		gwEncVolume := &goaviatrix.Gateway{
			GwName:              d.Get("gw_name").(string),
			CustomerManagedKeys: customerManagedKeys,
		}
		err := client.EnableEncryptVolume(gwEncVolume)
		if err != nil {
			return fmt.Errorf("failed to enable encrypt gateway volume for %s due to %s", gwEncVolume.GwName, err)
		}
	}
	if enableSNat && snatMode == "secondary" {
		if len(d.Get("snat_policy").([]interface{})) != 0 {
			return fmt.Errorf("'snat_policy' should be empty for 'snat_mode' of 'secondary'")
		}
		gwToEnableSNat := &goaviatrix.Gateway{
			GatewayName: d.Get("gw_name").(string),
		}
		gwToEnableSNat.EnableNat = "yes"
		gwToEnableSNat.SnatMode = "secondary"
		err := client.EnableSNat(gwToEnableSNat)
		if err != nil {
			return fmt.Errorf("failed to enable SNAT of 'secondary': %s", err)
		}
	} else if enableSNat && snatMode == "custom" {
		if len(d.Get("snat_policy").([]interface{})) == 0 {
			return fmt.Errorf("please specify 'snat_policy' for 'snat_mode' of 'custom'")
		}
		gwToEnableSNat := &goaviatrix.Gateway{
			GatewayName: d.Get("gw_name").(string),
		}
		gwToEnableSNat.EnableNat = "yes"
		gwToEnableSNat.SnatMode = "custom"
		if _, ok := d.GetOk("snat_policy"); ok {
			policies := d.Get("snat_policy").([]interface{})
			for _, policy := range policies {
				pl := policy.(map[string]interface{})
				customPolicy := &goaviatrix.PolicyRule{
					SrcIP:      pl["src_ip"].(string),
					SrcPort:    pl["src_port"].(string),
					DstIP:      pl["dst_ip"].(string),
					DstPort:    pl["dst_port"].(string),
					Protocol:   pl["protocol"].(string),
					Interface:  pl["interface"].(string),
					Connection: pl["connection"].(string),
					Mark:       pl["mark"].(string),
					NewSrcIP:   pl["new_src_ip"].(string),
					NewSrcPort: pl["new_src_port"].(string),
					ExcludeRTB: pl["exclude_rtb"].(string),
				}
				gwToEnableSNat.SnatPolicy = append(gwToEnableSNat.SnatPolicy, *customPolicy)
			}
		}
		time.Sleep(60 * time.Second)

		err := client.EnableSNat(gwToEnableSNat)
		if err != nil {
			return fmt.Errorf("failed to enable SNAT of 'custom': %s", err)
		}
	}

	if _, ok := d.GetOk("dnat_policy"); ok {
		if len(d.Get("dnat_policy").([]interface{})) != 0 {
			gwToUpdateDNat := &goaviatrix.Gateway{
				GatewayName: d.Get("gw_name").(string),
			}
			policies := d.Get("dnat_policy").([]interface{})
			for _, policy := range policies {
				dP := policy.(map[string]interface{})
				dNatPolicy := &goaviatrix.PolicyRule{
					SrcIP:      dP["src_ip"].(string),
					SrcPort:    dP["src_port"].(string),
					DstIP:      dP["dst_ip"].(string),
					DstPort:    dP["dst_port"].(string),
					Protocol:   dP["protocol"].(string),
					Interface:  dP["interface"].(string),
					Connection: dP["connection"].(string),
					Mark:       dP["mark"].(string),
					NewSrcIP:   dP["new_src_ip"].(string),
					NewSrcPort: dP["new_src_port"].(string),
					ExcludeRTB: dP["exclude_rtb"].(string),
				}
				gwToUpdateDNat.DnatPolicy = append(gwToUpdateDNat.DnatPolicy, *dNatPolicy)
			}
			if !enableSNat {
				time.Sleep(60 * time.Second)
			}

			err := client.UpdateDNat(gwToUpdateDNat)
			if err != nil {
				return fmt.Errorf("failed to update DNAT: %s", err)
			}
		}
	}

	return resourceAviatrixSpokeGatewayReadIfRequired(d, meta, &flag)
}

func resourceAviatrixSpokeGatewayReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixSpokeGatewayRead(d, meta)
	}
	return nil
}

func resourceAviatrixSpokeGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gwName := d.Get("gw_name").(string)
	if gwName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no gateway name received. Import Id is %s", id)
		d.Set("gw_name", id)
		d.SetId(id)
	}

	gateway := &goaviatrix.Gateway{
		AccountName: d.Get("account_name").(string),
		GwName:      d.Get("gw_name").(string),
	}

	gw, err := client.GetGateway(gateway)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix Spoke Gateway: %s", err)
	}

	log.Printf("[TRACE] reading spoke gateway %s: %#v", d.Get("gw_name").(string), gw)

	if gw != nil {
		d.Set("cloud_type", gw.CloudType)
		d.Set("account_name", gw.AccountName)

		if gw.CloudType == 1 {
			d.Set("vpc_id", strings.Split(gw.VpcID, "~~")[0]) //aws vpc_id returns as <vpc_id>~~<other vpc info> in rest api
			d.Set("vpc_reg", gw.VpcRegion)                    //aws vpc_reg returns as vpc_region in rest api

			if gw.AllocateNewEipRead {
				d.Set("allocate_new_eip", true)
			} else {
				d.Set("allocate_new_eip", false)
			}
		} else if gw.CloudType == 4 {
			d.Set("vpc_id", strings.Split(gw.VpcID, "~-~")[0]) //gcp vpc_id returns as <vpc_id>~-~<other vpc info> in rest api
			d.Set("vpc_reg", gw.GatewayZone)                   //gcp vpc_reg returns as gateway_zone in json

			d.Set("allocate_new_eip", true)
		} else if gw.CloudType == 8 || gw.CloudType == 16 {
			d.Set("vpc_id", gw.VpcID)
			d.Set("vpc_reg", gw.VpcRegion)

			d.Set("allocate_new_eip", true)
		}
		d.Set("enable_encrypt_volume", gw.EnableEncryptVolume)
		d.Set("eip", gw.PublicIP)

		d.Set("subnet", gw.VpcNet)
		d.Set("gw_size", gw.GwSize)
		d.Set("cloud_instance_id", gw.CloudnGatewayInstID)

		if gw.SingleAZ == "yes" {
			d.Set("single_az_ha", true)
		} else {
			d.Set("single_az_ha", false)
		}

		if gw.InsaneMode == "yes" {
			d.Set("insane_mode", true)
			if gw.CloudType == 1 {
				d.Set("insane_mode_az", gw.GatewayZone)
			} else {
				d.Set("insane_mode_az", "")
			}
		} else {
			d.Set("insane_mode", false)
			d.Set("insane_mode_az", "")
		}

		if gw.EnableActiveMesh == "yes" {
			d.Set("enable_active_mesh", true)
		} else {
			d.Set("enable_active_mesh", false)
		}

		if gw.CloudType == 1 && gw.EnableVpcDnsServer == "Enabled" {
			d.Set("enable_vpc_dns_server", true)
		} else {
			d.Set("enable_vpc_dns_server", false)
		}

		gwDetail, err := client.GetGatewayDetail(gateway)
		if err != nil {
			return fmt.Errorf("couldn't get detail information of Aviatrix Spoke Gateway: %s due to: %s", gw.GwName, err)
		}
		if len(gwDetail.DnatPolicy) != 0 {
			var dnatPolicy []map[string]interface{}
			for _, policy := range gwDetail.DnatPolicy {
				dP := make(map[string]interface{})
				dP["src_ip"] = policy.SrcIP
				dP["src_port"] = policy.SrcPort
				dP["dst_ip"] = policy.DstIP
				dP["dst_port"] = policy.DstPort
				dP["protocol"] = policy.Protocol
				dP["interface"] = policy.Interface
				dP["connection"] = policy.Connection
				dP["mark"] = policy.Mark
				dP["new_src_ip"] = policy.NewSrcIP
				dP["new_src_port"] = policy.NewSrcPort
				dP["exclude_rtb"] = policy.ExcludeRTB
				dnatPolicy = append(dnatPolicy, dP)
			}

			if err := d.Set("dnat_policy", dnatPolicy); err != nil {
				log.Printf("[WARN] Error setting 'dnat_policy' for (%s): %s", d.Id(), err)
			}
		} else {
			d.Set("dnat_policy", nil)
		}

		if gw.EnableNat == "yes" {
			d.Set("enable_snat", true)
			if gw.SnatMode == "customized" {
				d.Set("snat_mode", "custom")
				var snatPolicy []map[string]interface{}
				for _, policy := range gwDetail.SnatPolicy {
					sP := make(map[string]interface{})
					sP["src_ip"] = policy.SrcIP
					sP["src_port"] = policy.SrcPort
					sP["dst_ip"] = policy.DstIP
					sP["dst_port"] = policy.DstPort
					sP["protocol"] = policy.Protocol
					sP["interface"] = policy.Interface
					sP["connection"] = policy.Connection
					sP["mark"] = policy.Mark
					sP["new_src_ip"] = policy.NewSrcIP
					sP["new_src_port"] = policy.NewSrcPort
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
			d.Set("enable_snat", false)
			d.Set("snat_mode", "primary")
			d.Set("snat_policy", nil)
		}
	}

	if gw.SpokeVpc == "yes" {
		d.Set("transit_gw", gw.TransitGwName)
	} else {
		d.Set("transit_gw", "")
	}

	if gw.CloudType == 1 {
		tags := &goaviatrix.Tags{
			CloudType:    1,
			ResourceType: "gw",
			ResourceName: d.Get("gw_name").(string),
		}
		tagList, err := client.GetTags(tags)
		if err != nil {
			return fmt.Errorf("unable to read tag_list for gateway: %v due to %v", gateway.GwName, err)
		}
		var tagListStr []string
		if _, ok := d.GetOk("tag_list"); ok {
			tagList1 := d.Get("tag_list").([]interface{})
			tagListStr = goaviatrix.ExpandStringList(tagList1)
		}
		if len(goaviatrix.Difference(tagListStr, tagList)) != 0 || len(goaviatrix.Difference(tagList, tagListStr)) != 0 {
			if err := d.Set("tag_list", tagList); err != nil {
				log.Printf("[WARN] Error setting tag_list for (%s): %s", d.Id(), err)
			}
		} else {
			if err := d.Set("tag_list", tagListStr); err != nil {
				log.Printf("[WARN] Error setting tag_list for (%s): %s", d.Id(), err)
			}
		}
	}

	haGateway := &goaviatrix.Gateway{
		AccountName: d.Get("account_name").(string),
		GwName:      d.Get("gw_name").(string) + "-hagw",
	}

	haGw, err := client.GetGateway(haGateway)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.Set("ha_gw_size", "")
			d.Set("ha_subnet", "")
			d.Set("ha_zone", "")
			d.Set("ha_eip", "")
			d.Set("ha_insane_mode_az", "")
		} else {
			return fmt.Errorf("couldn't find Aviatrix Spoke HA Gateway: %s", err)
		}
	} else {
		log.Printf("[INFO] Spoke HA Gateway size: %s", haGw.GwSize)
		if haGw.CloudType == 1 || haGw.CloudType == 8 || haGw.CloudType == 16 {
			d.Set("ha_subnet", haGw.VpcNet)
			d.Set("ha_zone", "")
		} else if haGw.CloudType == 4 {
			d.Set("ha_zone", haGw.GatewayZone)
			d.Set("ha_subnet", "")
		}

		d.Set("ha_eip", haGw.PublicIP)
		d.Set("ha_gw_size", haGw.GwSize)
		if haGw.InsaneMode == "yes" && haGw.CloudType == 1 {
			d.Set("ha_insane_mode_az", haGw.GatewayZone)
		} else {
			d.Set("ha_insane_mode_az", "")
		}
	}

	return nil
}

func resourceAviatrixSpokeGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gateway := &goaviatrix.Gateway{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string),
	}

	haGateway := &goaviatrix.Gateway{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string) + "-hagw",
	}

	log.Printf("[INFO] Updating Aviatrix gateway: %#v", gateway)

	d.Partial(true)

	if d.HasChange("cloud_type") {
		return fmt.Errorf("updating cloud_type is not allowed")
	}
	if d.HasChange("account_name") {
		return fmt.Errorf("updating account_name is not allowed")
	}
	if d.HasChange("gw_name") {
		return fmt.Errorf("updating gw_name is not allowed")
	}
	if d.HasChange("vpc_id") {
		return fmt.Errorf("updating vpc_id is not allowed")
	}
	if d.HasChange("vpc_reg") {
		return fmt.Errorf("updating vpc_reg is not allowed")
	}
	if d.HasChange("subnet") {
		return fmt.Errorf("updating subnet is not allowed")
	}
	if d.HasChange("insane_mode") {
		return fmt.Errorf("updating insane_mode is not allowed")
	}
	if d.HasChange("insane_mode_az") {
		return fmt.Errorf("updating insane_mode_az is not allowed")
	}
	if d.HasChange("single_az_ha") {
		singleAZGateway := &goaviatrix.Gateway{
			GwName: d.Get("gw_name").(string),
		}

		singleAZ := d.Get("single_az_ha").(bool)

		if singleAZ {
			singleAZGateway.SingleAZ = "enabled"
		} else {
			singleAZGateway.SingleAZ = "disabled"
		}

		if singleAZGateway.SingleAZ == "enabled" {
			log.Printf("[INFO] Enable Single AZ GW HA: %#v", singleAZGateway)

			err := client.EnableSingleAZGateway(singleAZGateway)
			if err != nil {
				return fmt.Errorf("failed to enable single AZ GW HA: %s", err)
			}
		} else if singleAZGateway.SingleAZ == "disabled" {
			log.Printf("[INFO] Disable Single AZ GW HA: %#v", singleAZGateway)
			err := client.DisableSingleAZGateway(singleAZGateway)
			if err != nil {
				return fmt.Errorf("failed to disable single AZ GW HA: %s", err)
			}
		}

		d.SetPartial("single_az_ha")
	}

	if d.HasChange("tag_list") && gateway.CloudType == 1 {
		tags := &goaviatrix.Tags{
			CloudType:    1,
			ResourceType: "gw",
			ResourceName: d.Get("gw_name").(string),
		}
		o, n := d.GetChange("tag_list")
		if o == nil {
			o = new([]interface{})
		}
		if n == nil {
			n = new([]interface{})
		}
		os := o.([]interface{})
		ns := n.([]interface{})
		oldList := goaviatrix.ExpandStringList(os)
		newList := goaviatrix.ExpandStringList(ns)
		oldTagList := goaviatrix.Difference(oldList, newList)
		newTagList := goaviatrix.Difference(newList, oldList)
		if len(oldTagList) != 0 || len(newTagList) != 0 {
			if len(oldTagList) != 0 {
				oldTagList = goaviatrix.TagListStrColon(oldTagList)
				tags.TagList = strings.Join(oldTagList, ",")
				err := client.DeleteTags(tags)
				if err != nil {
					return fmt.Errorf("failed to delete tags : %s", err)
				}
			}
			if len(newTagList) != 0 {
				newTagList = goaviatrix.TagListStrColon(newTagList)
				tags.TagList = strings.Join(newTagList, ",")
				err := client.AddTags(tags)
				if err != nil {
					return fmt.Errorf("failed to add tags : %s", err)
				}
			}
		}

		d.SetPartial("tag_list")
	} else if d.HasChange("tag_list") && gateway.CloudType != 1 {
		return fmt.Errorf("adding tags is only supported for aws, cloud_type must be set to 1")
	}

	//Get primary gw size if gw_size changed, to be used later on for ha gateway size update
	primaryGwSize := d.Get("gw_size").(string)
	if d.HasChange("gw_size") {
		old, _ := d.GetChange("gw_size")
		primaryGwSize = old.(string)
		gateway.GwSize = d.Get("gw_size").(string)
		err := client.UpdateGateway(gateway)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix Spoke Gateway: %s", err)
		}
		d.SetPartial("gw_size")
	}

	newHaGwEnabled := false
	if d.HasChange("ha_subnet") || d.HasChange("ha_zone") || d.HasChange("ha_insane_mode_az") {
		spokeGw := &goaviatrix.SpokeVpc{
			GwName:    d.Get("gw_name").(string),
			CloudType: d.Get("cloud_type").(int),
		}

		if spokeGw.CloudType == 1 {
			spokeGw.Eip = d.Get("ha_eip").(string)
		}

		if !d.HasChange("ha_subnet") && d.HasChange("ha_insane_mode_az") {
			return fmt.Errorf("ha_subnet must change if ha_insane_mode_az changes")
		}
		if d.Get("insane_mode").(bool) == true && spokeGw.CloudType == 1 {
			var haStrs []string
			insaneModeHaAz := d.Get("ha_insane_mode_az").(string)
			if insaneModeHaAz == "" {
				return fmt.Errorf("ha_insane_mode_az needed if insane_mode is enabled and ha_subnet is set")
			}
			haStrs = append(haStrs, spokeGw.HASubnet, insaneModeHaAz)
			spokeGw.HASubnet = strings.Join(haStrs, "~~")
		}

		oldSubnet, newSubnet := d.GetChange("ha_subnet")
		oldZone, newZone := d.GetChange("ha_zone")
		deleteHaGw := false
		changeHaGw := false
		if spokeGw.CloudType == 1 || spokeGw.CloudType == 8 || spokeGw.CloudType == 16 {
			spokeGw.HASubnet = d.Get("ha_subnet").(string)
			if oldSubnet == "" && newSubnet != "" {
				newHaGwEnabled = true
			} else if oldSubnet != "" && newSubnet == "" {
				deleteHaGw = true
			} else if oldSubnet != "" && newSubnet != "" {
				changeHaGw = true
			}
		} else if spokeGw.CloudType == 4 {
			spokeGw.HAZone = d.Get("ha_zone").(string)
			if oldZone == "" && newZone != "" {
				newHaGwEnabled = true
			} else if oldZone != "" && newZone == "" {
				deleteHaGw = true
			} else if oldZone != "" && newZone != "" {
				changeHaGw = true
			}
		}
		if newHaGwEnabled {
			//New configuration to enable HA
			err := client.EnableHaSpokeVpc(spokeGw)
			if err != nil {
				return fmt.Errorf("failed to enable HA Aviatrix Spoke Gateway: %s", err)
			}
			newHaGwEnabled = true
		} else if deleteHaGw {
			//Ha configuration has been deleted
			err := client.DeleteGateway(haGateway)
			if err != nil {
				return fmt.Errorf("failed to delete Aviatrix Spoke HA gateway: %s", err)
			}
		} else if changeHaGw {
			//HA subnet has been modified. Delete older HA GW,
			// and launch new HA GW in new subnet.
			err := client.DeleteGateway(haGateway)
			if err != nil {
				return fmt.Errorf("failed to delete Aviatrix Spoke HA gateway: %s", err)
			}

			gateway.GwName = d.Get("spokeGw_name").(string)
			//New configuration to enable HA
			haErr := client.EnableHaSpokeVpc(spokeGw)
			if haErr != nil {
				return fmt.Errorf("failed to enable HA Aviatrix Spoke Gateway: %s", err)
			}
		}
		d.SetPartial("ha_subnet")
		d.SetPartial("ha_zone")
		d.SetPartial("ha_insane_mode_az")
	}

	if d.HasChange("ha_gw_size") || newHaGwEnabled {
		newHaGwSize := d.Get("ha_gw_size").(string)
		if !newHaGwEnabled || (newHaGwSize != primaryGwSize) {
			// MODIFIES HA GW SIZE if
			// Ha gateway wasn't newly configured
			// OR
			// newly configured Ha gateway is set to be different size than primary gateway
			// (when ha gateway is enabled, it's size is by default the same as primary gateway)
			_, err := client.GetGateway(haGateway)
			if err != nil {
				if err == goaviatrix.ErrNotFound {
					d.Set("ha_gw_size", "")
					d.Set("ha_subnet", "")
					d.Set("ha_zone", "")
					d.Set("ha_insane_mode_az", "")
					return nil
				}
				return fmt.Errorf("couldn't find Aviatrix Spoke HA Gateway while trying to update HA Gw size: %s", err)
			}
			haGateway.GwSize = d.Get("ha_gw_size").(string)
			if haGateway.GwSize == "" {
				return fmt.Errorf("A valid non empty ha_gw_size parameter is mandatory for this resource if " +
					"ha_subnet or ha_zone is set. Example: t2.micro or us-west1-b")
			}
			err = client.UpdateGateway(haGateway)
			log.Printf("[INFO] Updating HA Gateway size to: %s ", haGateway.GwSize)
			if err != nil {
				return fmt.Errorf("failed to update Aviatrix Spoke HA Gateway size: %s", err)
			}
		}
		d.SetPartial("ha_gw_size")
	}

	if d.HasChange("enable_snat") {
		enableSNat := d.Get("enable_snat").(bool)
		if enableSNat {
			gw := &goaviatrix.Gateway{
				CloudType:   d.Get("cloud_type").(int),
				GatewayName: d.Get("gw_name").(string),
			}
			snatMode := d.Get("snat_mode").(string)
			if snatMode == "secondary" {
				if len(d.Get("snat_policy").([]interface{})) != 0 {
					return fmt.Errorf("'snat_policy' should be empty for 'snat_mode' of 'secondary'")
				}
				gw.EnableNat = "yes"
				gw.SnatMode = "secondary"
				err := client.EnableSNat(gw)
				if err != nil {
					return fmt.Errorf("failed to enable SNAT for 'snat_mode' of 'secondary': %s", err)
				}
			} else if snatMode == "custom" {
				if len(d.Get("snat_policy").([]interface{})) == 0 {
					return fmt.Errorf("please specify 'snat_policy' for 'snat_mode' of 'custom'")
				}
				gw.EnableNat = "yes"
				gw.SnatMode = "custom"
				if _, ok := d.GetOk("snat_policy"); ok {
					policies := d.Get("snat_policy").([]interface{})
					for _, policy := range policies {
						pl := policy.(map[string]interface{})
						customPolicy := &goaviatrix.PolicyRule{
							SrcIP:      pl["src_ip"].(string),
							SrcPort:    pl["src_port"].(string),
							DstIP:      pl["dst_ip"].(string),
							DstPort:    pl["dst_port"].(string),
							Protocol:   pl["protocol"].(string),
							Interface:  pl["interface"].(string),
							Connection: pl["connection"].(string),
							Mark:       pl["mark"].(string),
							NewSrcIP:   pl["new_src_ip"].(string),
							NewSrcPort: pl["new_src_port"].(string),
							ExcludeRTB: pl["exclude_rtb"].(string),
						}
						gw.SnatPolicy = append(gw.SnatPolicy, *customPolicy)
					}
				}
				err := client.EnableSNat(gw)
				if err != nil {
					return fmt.Errorf("failed to enable SNAT for 'snat_mode' of 'custom': %s", err)
				}
			} else if snatMode == "primary" {
				if len(d.Get("snat_policy").([]interface{})) != 0 {
					return fmt.Errorf("'snat_policy' should be empty for 'snat_mode' of 'primary'")
				}
				gw.EnableNat = "yes"
				gw.SnatMode = "primary"
				err := client.EnableSNat(gw)
				if err != nil {
					return fmt.Errorf("failed to enable SNAT of 'primary': %s", err)
				}
			}
		} else {
			if len(d.Get("snat_policy").([]interface{})) != 0 {
				return fmt.Errorf("'snat_policy' should be empty for disabling SNAT")
			}
			gw := &goaviatrix.Gateway{
				CloudType: d.Get("cloud_type").(int),
				GwName:    d.Get("gw_name").(string),
			}
			err := client.DisableSNat(gw)
			if err != nil {
				return fmt.Errorf("failed to disable SNAT: %s", err)
			}
		}
	} else if d.Get("enable_snat").(bool) {
		if !d.HasChange("snat_mode") && d.Get("snat_mode").(string) == "custom" && d.HasChange("snat_policy") {
			if len(d.Get("snat_policy").([]interface{})) == 0 {
				return fmt.Errorf("please specify 'snat_policy'for 'snat_mode' of 'custom'")
			}
			gw := &goaviatrix.Gateway{
				CloudType:   d.Get("cloud_type").(int),
				GatewayName: d.Get("gw_name").(string),
			}
			gw.EnableNat = "yes"
			gw.SnatMode = "custom"
			if _, ok := d.GetOk("snat_policy"); ok {
				policies := d.Get("snat_policy").([]interface{})
				for _, policy := range policies {
					pl := policy.(map[string]interface{})
					customPolicy := &goaviatrix.PolicyRule{
						SrcIP:      pl["src_ip"].(string),
						SrcPort:    pl["src_port"].(string),
						DstIP:      pl["dst_ip"].(string),
						DstPort:    pl["dst_port"].(string),
						Protocol:   pl["protocol"].(string),
						Interface:  pl["interface"].(string),
						Connection: pl["connection"].(string),
						Mark:       pl["mark"].(string),
						NewSrcIP:   pl["new_src_ip"].(string),
						NewSrcPort: pl["new_src_port"].(string),
						ExcludeRTB: pl["exclude_rtb"].(string),
					}
					gw.SnatPolicy = append(gw.SnatPolicy, *customPolicy)
				}
			}
			err := client.EnableSNat(gw)
			if err != nil {
				return fmt.Errorf("failed to update 'snat_policy' for 'snat_mode' of 'custom': %s", err)
			}
		} else if d.HasChange("snat_mode") || d.HasChange("snat_policy") {
			gw := &goaviatrix.Gateway{
				CloudType: d.Get("cloud_type").(int),
				GwName:    d.Get("gw_name").(string),
			}
			err := client.DisableSNat(gw)
			if err != nil {
				return fmt.Errorf("failed to disable SNAT: %s", err)
			}

			gw.GatewayName = d.Get("gw_name").(string)
			snatMode := d.Get("snat_mode").(string)
			if snatMode == "secondary" {
				if len(d.Get("snat_policy").([]interface{})) != 0 {
					return fmt.Errorf("'snat_policy' should be empty for 'snat_mode' of 'secondary'")
				}
				gw.EnableNat = "yes"
				gw.SnatMode = "secondary"
				err := client.EnableSNat(gw)
				if err != nil {
					return fmt.Errorf("failed to enable SNAT for 'snat_mode' of 'secondary': %s", err)
				}
			} else if snatMode == "custom" {
				if len(d.Get("snat_policy").([]interface{})) == 0 {
					return fmt.Errorf("please specify 'snat_policy' for 'snat_mode' of 'custom'")
				}
				gw.EnableNat = "yes"
				gw.SnatMode = "custom"

				if _, ok := d.GetOk("snat_policy"); ok {
					policies := d.Get("snat_policy").([]interface{})
					for _, policy := range policies {
						pl := policy.(map[string]interface{})
						customPolicy := &goaviatrix.PolicyRule{
							SrcIP:      pl["src_ip"].(string),
							SrcPort:    pl["src_port"].(string),
							DstIP:      pl["dst_ip"].(string),
							DstPort:    pl["dst_port"].(string),
							Protocol:   pl["protocol"].(string),
							Interface:  pl["interface"].(string),
							Connection: pl["connection"].(string),
							Mark:       pl["mark"].(string),
							NewSrcIP:   pl["new_src_ip"].(string),
							NewSrcPort: pl["new_src_port"].(string),
							ExcludeRTB: pl["exclude_rtb"].(string),
						}
						gw.SnatPolicy = append(gw.SnatPolicy, *customPolicy)
					}
				}
				err := client.EnableSNat(gw)
				if err != nil {
					return fmt.Errorf("failed to enable SNAT of 'secondary': %s", err)
				}
			} else if snatMode == "primary" {
				if len(d.Get("snat_policy").([]interface{})) != 0 {
					return fmt.Errorf("'snat_policy' should be empty for 'snat_mode' of 'primary'")
				}
				err := client.EnableSNat(gw)
				if err != nil {
					return fmt.Errorf("failed to enable SNAT of 'primary': %s", err)
				}
			}
		}
	} else if !d.Get("enable_snat").(bool) {
		if d.HasChange("snat_mode") || d.HasChange("snat_policy") {
			if d.Get("snat_policy").(interface{}) == nil {
				return fmt.Errorf("source NAT is disabled, can't update 'snat_mode' or 'snat_policy'")
			}
		}
	}

	if d.HasChange("transit_gw") {
		spokeVPC := &goaviatrix.SpokeVpc{
			CloudType:      d.Get("cloud_type").(int),
			GwName:         d.Get("gw_name").(string),
			HASubnet:       d.Get("ha_subnet").(string),
			TransitGateway: d.Get("transit_gw").(string),
		}

		o, n := d.GetChange("transit_gw")
		if o == "" {
			//New configuration to join to transit GW
			err := client.SpokeJoinTransit(spokeVPC)
			if err != nil {
				return fmt.Errorf("failed to join Transit Gateway: %s", err)
			}
		} else if n == "" {
			//Transit GW has been deleted, leave transit GW.
			err := client.SpokeLeaveTransit(spokeVPC)
			if err != nil {
				return fmt.Errorf("failed to leave Transit Gateway: %s", err)
			}
		} else {
			//Change transit GW
			err := client.SpokeLeaveTransit(spokeVPC)
			if err != nil {
				return fmt.Errorf("failed to leave Transit Gateway: %s", err)
			}

			err = client.SpokeJoinTransit(spokeVPC)
			if err != nil {
				return fmt.Errorf("failed to join Transit Gateway: %s", err)
			}
		}

		d.SetPartial("transit_gw")
	}

	if d.HasChange("enable_active_mesh") {
		gw := &goaviatrix.Gateway{
			GwName: d.Get("gw_name").(string),
		}

		enableActiveMesh := d.Get("enable_active_mesh").(bool)
		if enableActiveMesh {
			gw.EnableActiveMesh = "yes"
			err := client.EnableActiveMesh(gw)
			if err != nil {
				return fmt.Errorf("failed to enable Active Mesh Mode: %s", err)
			}
		} else {
			gw.EnableActiveMesh = "no"
			err := client.DisableActiveMesh(gw)
			if err != nil {
				return fmt.Errorf("failed to disable Active Mesh Mode: %s", err)
			}
		}
	}

	if d.HasChange("enable_vpc_dns_server") && d.Get("cloud_type").(int) == 1 {
		gw := &goaviatrix.Gateway{
			CloudType: d.Get("cloud_type").(int),
			GwName:    d.Get("gw_name").(string),
		}

		enableVpcDnsServer := d.Get("enable_vpc_dns_server").(bool)
		if enableVpcDnsServer {
			err := client.EnableVpcDnsServer(gw)
			if err != nil {
				return fmt.Errorf("failed to enable VPC DNS Server: %s", err)
			}
		} else {
			err := client.DisableVpcDnsServer(gw)
			if err != nil {
				return fmt.Errorf("failed to disable VPC DNS Server: %s", err)
			}
		}

		d.SetPartial("enable_vpc_dns_server")
	} else if d.HasChange("enable_vpc_dns_server") {
		return fmt.Errorf("'enable_vpc_dns_server' only supports AWS(1)")
	}

	if d.HasChange("enable_encrypt_volume") {
		if d.Get("enable_encrypt_volume").(bool) {
			if d.Get("cloud_type").(int) != 1 {
				return fmt.Errorf("'enable_encrypt_volume' is only supported for AWS provider (cloud_type: 1)")
			}
			if d.Get("single_az_ha").(bool) {
				return fmt.Errorf("'single_az_ha' needs to be disabled to encrypt gateway EBS volume")
			}
			gwEncVolume := &goaviatrix.Gateway{
				GwName:              d.Get("gw_name").(string),
				CustomerManagedKeys: d.Get("customer_managed_keys").(string),
			}
			err := client.EnableEncryptVolume(gwEncVolume)
			if err != nil {
				return fmt.Errorf("failed to enable encrypt gateway volume for %s due to %s", gwEncVolume.GwName, err)
			}
		} else {
			return fmt.Errorf("can't disable Encrypt Volume for gateway: %s", gateway.GwName)
		}
	} else if d.HasChange("customer_managed_keys") {
		return fmt.Errorf("updating customer_managed_keys only is not allowed")
	}

	if d.HasChange("dnat_policy") {
		gwToUpdateDNat := &goaviatrix.Gateway{
			GatewayName: d.Get("gw_name").(string),
		}
		if len(d.Get("dnat_policy").([]interface{})) != 0 {
			policies := d.Get("dnat_policy").([]interface{})
			for _, policy := range policies {
				dP := policy.(map[string]interface{})
				dNatPolicy := &goaviatrix.PolicyRule{
					SrcIP:      dP["src_ip"].(string),
					SrcPort:    dP["src_port"].(string),
					DstIP:      dP["dst_ip"].(string),
					DstPort:    dP["dst_port"].(string),
					Protocol:   dP["protocol"].(string),
					Interface:  dP["interface"].(string),
					Connection: dP["connection"].(string),
					Mark:       dP["mark"].(string),
					NewSrcIP:   dP["new_src_ip"].(string),
					NewSrcPort: dP["new_src_port"].(string),
					ExcludeRTB: dP["exclude_rtb"].(string),
				}
				gwToUpdateDNat.DnatPolicy = append(gwToUpdateDNat.DnatPolicy, *dNatPolicy)
			}
		} else {
			gwToUpdateDNat.DnatPolicy = make([]goaviatrix.PolicyRule, 0)
		}
		err := client.UpdateDNat(gwToUpdateDNat)
		if err != nil {
			return fmt.Errorf("failed to update DNAT: %s", err)
		}
	}

	d.Partial(false)
	d.SetId(gateway.GwName)
	return resourceAviatrixSpokeGatewayRead(d, meta)
}

func resourceAviatrixSpokeGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	gateway := &goaviatrix.Gateway{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix Spoke Gateway: %#v", gateway)

	if transitGw := d.Get("transit_gw").(string); transitGw != "" {
		spokeVPC := &goaviatrix.SpokeVpc{
			GwName: d.Get("gw_name").(string),
		}

		err := client.SpokeLeaveTransit(spokeVPC)
		if err != nil {
			return fmt.Errorf("failed to leave Transit Gateway: %s", err)
		}
	}

	//If HA is enabled, delete HA GW first.
	haSubnet := d.Get("ha_subnet").(string)
	haZone := d.Get("ha_zone").(string)
	if haSubnet != "" || haZone != "" {
		//Delete HA Gw too
		gateway.GwName += "-hagw"
		err := client.DeleteGateway(gateway)
		if err != nil {
			return fmt.Errorf("failed to delete Aviatrix Spoke HA gateway: %s", err)
		}
	}

	gateway.GwName = d.Get("gw_name").(string)

	err := client.DeleteGateway(gateway)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Spoke Gateway: %s", err)
	}

	return nil
}
