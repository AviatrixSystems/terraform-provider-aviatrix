package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixGatewayCreate,
		Read:   resourceAviatrixGatewayRead,
		Update: resourceAviatrixGatewayUpdate,
		Delete: resourceAviatrixGatewayDelete,
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
				Description: "Account name. This account will be used to launch Aviatrix gateway.",
			},
			"gw_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Aviatrix gateway unique name.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of legacy VPC/Vnet to be connected.",
			},
			"vpc_reg": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Region where this gateway will be launched.",
			},
			"vpc_size": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Size of Gateway Instance.",
			},
			"vpc_net": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A VPC Network address range selected from one of the available network ranges.",
			},
			"ha_subnet": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This is for Gateway HA. Deprecated.",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP address of the Gateway created.",
			},
			"backup_public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private IP address of the Gateway created.",
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Security group used for the gateway.",
			},
			"enable_nat": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "no",
				Description: "Enable NAT for this container.",
			},
			"public_dns_server": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "NS server used by the gateway.",
			},
			"vpn_access": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "no",
				Description: "Enable user access through VPN to this container.",
			},
			"vpn_cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "VPN CIDR block for the container.",
			},
			"enable_elb": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "no",
				Description: "Specify whether to enable ELB or not.",
			},
			"elb_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "A name for the ELB that is created.",
			},
			"split_tunnel": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "yes",
				Description: "Specify split tunnel mode.",
			},
			"name_servers": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				Description: "A list of DNS servers used to resolve domain names by " +
					"a connected VPN user when Split Tunnel Mode is enabled.",
			},
			"search_domains": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				Description: "A list of domain names that will use the NameServer " +
					"when a specific name is not in the destination when Split Tunnel Mode is enabled.",
			},
			"additional_cidrs": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				Description: "A list of destination CIDR ranges that will also go through the VPN tunnel " +
					"when Split Tunnel Mode is enabled.",
			},
			"otp_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Two step authentication mode.",
			},
			"saml_enabled": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "no",
				Description: "This field indicates whether enabling SAML or not.",
			},
			"okta_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Token for Okta auth mode.",
			},
			"okta_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "URL for Okta auth mode.",
			},
			"okta_username_suffix": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Username suffix for Okta auth mode.",
			},
			"duo_integration_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Integration key for DUO auth mode.",
			},
			"duo_secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Default:     "",
				Description: "Secret key for DUO auth mode.",
			},
			"duo_api_hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "API hostname for DUO auth mode.",
			},
			"duo_push_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Push mode for DUO auth.",
			},
			"enable_ldap": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "no",
				Description: "Specify whether to enable LDAP or not. Supported values: 'yes' and 'no'.",
			},
			"ldap_server": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "LDAP server address. Required: Yes if enable_ldap is 'yes'.",
			},
			"ldap_bind_dn": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "LDAP bind DN. Required: Yes if enable_ldap is 'yes'.",
			},
			"ldap_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Default:     "",
				Description: "LDAP password. Required: Yes if enable_ldap is 'yes'.",
			},
			"ldap_base_dn": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "LDAP base DN. Required: Yes if enable_ldap is 'yes'.",
			},
			"ldap_username_attribute": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "LDAP user attribute. Required: Yes if enable_ldap is 'yes'.",
			},
			"peering_ha_subnet": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Public Subnet Information while creating Peering HA Gateway, only subnet is accepted.",
			},
			"peering_ha_eip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Public IP address that you want assigned to the HA peering instance.",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A GCE zone where this gateway will be launched.",
			},
			"cloud_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance ID of the gateway.",
			},
			"cloudn_bkup_gateway_inst_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance ID of the backup gateway.",
			},
			"single_az_ha": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "disabled",
				Description: "Set to 'enabled' if this feature is desired.",
			},
			"allocate_new_eip": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "on",
				Description: "When value is off, reuse an idle address in Elastic IP pool for this gateway. " +
					"Otherwise, allocate a new Elastic IP and use it for this gateway.",
			},
			"eip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Required when allocate_new_eip is 'off'. It uses specified EIP for this gateway.",
			},
			"tag_list": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Default:     nil,
				Description: "Instance tag of cloud provider.",
			},
		},
	}
}

func resourceAviatrixGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	gateway := &goaviatrix.Gateway{
		CloudType:          d.Get("cloud_type").(int),
		AccountName:        d.Get("account_name").(string),
		GwName:             d.Get("gw_name").(string),
		VpcID:              d.Get("vpc_id").(string),
		VpcRegion:          d.Get("vpc_reg").(string),
		VpcSize:            d.Get("vpc_size").(string),
		VpcNet:             d.Get("vpc_net").(string),
		EnableNat:          d.Get("enable_nat").(string),
		VpnStatus:          d.Get("vpn_access").(string),
		VpnCidr:            d.Get("vpn_cidr").(string),
		EnableElb:          d.Get("enable_elb").(string),
		ElbName:            d.Get("elb_name").(string),
		SplitTunnel:        d.Get("split_tunnel").(string),
		OtpMode:            d.Get("otp_mode").(string),
		SamlEnabled:        d.Get("saml_enabled").(string),
		OktaToken:          d.Get("okta_token").(string),
		OktaURL:            d.Get("okta_url").(string),
		OktaUsernameSuffix: d.Get("okta_username_suffix").(string),
		DuoIntegrationKey:  d.Get("duo_integration_key").(string),
		DuoSecretKey:       d.Get("duo_secret_key").(string),
		DuoAPIHostname:     d.Get("duo_api_hostname").(string),
		DuoPushMode:        d.Get("duo_push_mode").(string),
		EnableLdap:         d.Get("enable_ldap").(string),
		LdapServer:         d.Get("ldap_server").(string),
		LdapBindDn:         d.Get("ldap_bind_dn").(string),
		LdapPassword:       d.Get("ldap_password").(string),
		LdapBaseDn:         d.Get("ldap_base_dn").(string),
		LdapUserAttr:       d.Get("ldap_username_attribute").(string),
		HASubnet:           d.Get("ha_subnet").(string),
		PeeringHASubnet:    d.Get("peering_ha_subnet").(string),
		NewZone:            d.Get("zone").(string),
		SingleAZ:           d.Get("single_az_ha").(string),
		AllocateNewEip:     d.Get("allocate_new_eip").(string),
		Eip:                d.Get("eip").(string),
	}
	if gateway.OtpMode != "" && gateway.OtpMode != "2" && gateway.OtpMode != "3" {
		return fmt.Errorf("otp_mode can only be '2' or '3' or empty string")
	}
	if gateway.SamlEnabled == "yes" {
		if gateway.EnableLdap == "yes" || gateway.OtpMode != "" {
			return fmt.Errorf("ldap and mfa can't be configured if saml is enabled")
		}
	}
	if gateway.EnableLdap == "yes" && gateway.OtpMode == "3" {
		return fmt.Errorf("ldap can't be configured along with okta authentication")
	}
	if gateway.EnableLdap == "yes" {
		if gateway.LdapServer == "" {
			return fmt.Errorf("ldap server must be set if ldap is enabled")
		}
		if gateway.LdapBindDn == "" {
			return fmt.Errorf("ldap bind dn must be set if ldap is enabled")
		}
		if gateway.LdapPassword == "" {
			return fmt.Errorf("ldap password must be set if ldap is enabled")
		}
		if gateway.LdapBaseDn == "" {
			return fmt.Errorf("ldap base dn must be set if ldap is enabled")
		}
		if gateway.LdapUserAttr == "" {
			return fmt.Errorf("ldap user attribute must be set if ldap is enabled")
		}
	}
	if gateway.OtpMode == "2" {
		if gateway.DuoIntegrationKey == "" {
			return fmt.Errorf("duo integration key required if otp_mode set to 2")
		}
		if gateway.DuoSecretKey == "" {
			return fmt.Errorf("duo secret key required if otp_mode set to 2")
		}
		if gateway.DuoAPIHostname == "" {
			return fmt.Errorf("duo api hostname required if otp_mode set to 2")
		}
		if gateway.DuoPushMode != "auto" && gateway.DuoPushMode != "token" && gateway.DuoPushMode != "selective" {
			return fmt.Errorf("duo push mode must be set to a valid value (auto, selective, or token)")
		}
	} else if gateway.OtpMode == "3" {
		if gateway.OktaToken == "" {
			return fmt.Errorf("okta token must be set if otp_mode is set to 3")
		}
		if gateway.OktaURL == "" {
			return fmt.Errorf("okta url must be set if otp_mode is set to 3")
		}
	}

	if gateway.EnableElb != "yes" {
		gateway.EnableElb = "no"
	}
	if gateway.EnableNat != "yes" {
		gateway.EnableNat = "no"
	}
	if gateway.VpnStatus != "yes" {
		gateway.VpnStatus = "no"
	}
	if gateway.SplitTunnel != "no" {
		gateway.SplitTunnel = "yes"
	}
	if gateway.EnableElb == "yes" && gateway.VpnStatus != "yes" {
		return fmt.Errorf("can not enable elb without vpn access set to yes")
	}
	log.Printf("[INFO] Creating Aviatrix gateway: %#v", gateway)

	err := client.CreateGateway(gateway)
	if err != nil {
		log.Printf("[INFO] failed to create Aviatrix gateway: %#v", gateway)
		return fmt.Errorf("failed to create Aviatrix gateway: %s", err)
	}
	if enableNAT := d.Get("enable_nat").(string); enableNAT == "yes" {
		log.Printf("[INFO] Aviatrix NAT enabled gateway: %#v", gateway)
	}

	// single_AZ enabled for Gateway. https://docs.aviatrix.com/HowTos/gateway.html#high-availability
	if singleAZHA := d.Get("single_az_ha").(string); singleAZHA == "enabled" {
		singleAZGateway := &goaviatrix.Gateway{
			GwName:   d.Get("gw_name").(string),
			SingleAZ: d.Get("single_az_ha").(string),
		}
		log.Printf("[INFO] Enable Single AZ GW HA: %#v", singleAZGateway)
		err := client.EnableSingleAZGateway(gateway)
		if err != nil {
			return fmt.Errorf("failed to create single AZ GW HA: %s", err)
		}
	}

	// ha_subnet is for Gateway HA. Deprecated. https://docs.aviatrix.com/HowTos/gateway.html#high-availability
	if ha_subnet := d.Get("ha_subnet").(string); ha_subnet != "" {
		ha_gateway := &goaviatrix.Gateway{
			GwName:   d.Get("gw_name").(string),
			HASubnet: d.Get("ha_subnet").(string),
		}
		log.Printf("[INFO] Enable gateway HA: %#v", ha_gateway)
		err := client.EnableHaGateway(ha_gateway)
		if err != nil {
			del_err := client.DeleteGateway(gateway)
			if del_err != nil {
				return fmt.Errorf("failed to auto-cleanup failed gateway: %s", del_err)
			}
			return fmt.Errorf("failed to create GW HA: %s", err)
		}
	}
	// peering_ha_subnet is for Peering HA Gateway. https://docs.aviatrix.com/HowTos/gateway.html#high-availability
	if peeringHaSubnet := d.Get("peering_ha_subnet").(string); peeringHaSubnet != "" {
		peeringHaGateway := &goaviatrix.Gateway{
			Eip:             d.Get("peering_ha_eip").(string),
			GwName:          d.Get("gw_name").(string),
			PeeringHASubnet: d.Get("peering_ha_subnet").(string),
			NewZone:         d.Get("zone").(string),
		}
		log.Printf("[INFO] Enable peering HA: %#v", peeringHaGateway)
		err := client.EnablePeeringHaGateway(peeringHaGateway)
		if err != nil {
			return fmt.Errorf("failed to create peering HA: %s", err)
		}
	}
	d.SetId(gateway.GwName)

	if _, ok := d.GetOk("tag_list"); ok {
		tagList := d.Get("tag_list").([]interface{})
		tagListStr := goaviatrix.ExpandStringList(tagList)
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
	}
	if vpnAccess, ok := d.GetOk("vpn_access"); ok && vpnAccess == "yes" {
		gw := &goaviatrix.Gateway{
			GwName: gateway.GwName,
		}
		gw1, err := client.GetGateway(gw)
		if err != nil {
			return fmt.Errorf("couldn't find Aviatrix Gateway: %s due to %v", gw.GwName, err)
		}
		sTunnel := &goaviatrix.SplitTunnel{
			SplitTunnel:     "no",
			VpcID:           gateway.VpcID,
			AdditionalCidrs: d.Get("additional_cidrs").(string),
			NameServers:     d.Get("name_servers").(string),
			SearchDomains:   d.Get("search_domains").(string),
			SaveTemplate:    "no",
		}
		if gw1.EnableElb != "yes" {
			sTunnel.ElbName = gw1.GwName
		} else {
			sTunnel.ElbName = gw1.ElbName
		}
		if gateway.SplitTunnel != "" {
			sTunnel.SplitTunnel = gateway.SplitTunnel
		}
		if sTunnel.SplitTunnel != "" && sTunnel.SplitTunnel != "no" && sTunnel.SplitTunnel != "yes" {
			return fmt.Errorf("split_tunnel is not set correctly")
		}
		if sTunnel.SplitTunnel == "yes" {
			if sTunnel.AdditionalCidrs != "" || sTunnel.NameServers != "" || sTunnel.SearchDomains != "" {
				err = client.ModifySplitTunnel(sTunnel)
				if err != nil {
					return fmt.Errorf("failed to modify split tunnel: %s", err)
				}
			}
		}
	}
	return resourceAviatrixGatewayRead(d, meta)
}

func resourceAviatrixGatewayRead(d *schema.ResourceData, meta interface{}) error {
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
	if d.Get("single_az_ha") != nil {
		gateway.SingleAZ = d.Get("single_az_ha").(string)
	}
	gw, err := client.GetGateway(gateway)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix Gateway: %s", err)
	}
	log.Printf("[TRACE] reading gateway %s: %#v", d.Get("gw_name").(string), gw)
	if gw != nil {
		d.Set("cloud_type", gw.CloudType)
		d.Set("account_name", gw.AccountName)
		d.Set("gw_name", gw.GwName)
		d.Set("vpc_id", strings.Split(gw.VpcID, "~~")[0])
		d.Set("vpc_reg", gw.VpcRegion)
		d.Set("vpc_net", gw.VpcNet)
		if gw.EnableNat != "" {
			d.Set("enable_nat", gw.EnableNat)
		}
		if gw.AllocateNewEipRead {
			d.Set("allocate_new_eip", "on")
		} else {
			d.Set("allocate_new_eip", "off")
		}
		if gw.EnableLdapRead {
			d.Set("enable_ldap", "yes")
		} else {
			d.Set("enable_ldap", "no")
		}
		if gw.VpnStatus != "" {
			if gw.VpnStatus == "disabled" {
				d.Set("vpn_access", "no")
			} else if gw.VpnStatus == "enabled" {
				d.Set("vpn_access", "yes")
			} else {
				d.Set("vpn_access", gw.VpnStatus)
			}
		}
		vpnAccess := d.Get("vpn_access")
		if vpnAccess == "no" {
			d.Set("split_tunnel", "yes")
		} else {
			d.Set("split_tunnel", gw.SplitTunnel)
		}
		d.Set("vpn_cidr", gw.VpnCidr)
		if gw.ElbState == "enabled" {
			d.Set("enable_elb", "yes")
			elb_name := d.Get("elb_name")
			// Versions prior to 4.0 won't return elb_name, so deduce it from elb_dns_name
			if elb_name == "" {
				elb_dns_name := gw.ElbDNSName
				log.Printf("[INFO] Controllers prior to 4.0 do not return elb_name. Deducing from elb_dns_name")
				if elb_dns_name != "" {
					hostname := strings.Split(elb_dns_name, ".")[0]
					// AWS adds - followed by a random string after the name given to the ELB
					parts := strings.Split(hostname, "-")
					// Remove random string added by AWS
					parts[len(parts)-1] = ""
					parts = parts[:len(parts)-1]
					// Join again using - in case it was used in the ELB Name
					elb_name = strings.Join(parts, "-")
				} else {
					return fmt.Errorf("neither elb_name or elb_dns_name returned by the API in an ELB enabled gateway")
				}
			}
			d.Set("elb_name", elb_name)
		} else {
			d.Set("enable_elb", "no")
			d.Set("elb_name", "")
		}
		if gw.SamlEnabled != "" {
			d.Set("saml_enabled", gw.SamlEnabled)
		}
		if gw.AuthMethod == "duo_auth" || gw.AuthMethod == "duo_auth+LDAP" {
			d.Set("otp_mode", "2")
		} else if gw.AuthMethod == "okta_auth" {
			d.Set("otp_mode", "3")
		} else {
			d.Set("otp_mode", "")
		}
		d.Set("okta_token", gw.OktaToken)
		d.Set("okta_url", gw.OktaURL)
		d.Set("okta_username_suffix", gw.OktaUsernameSuffix)
		d.Set("duo_integration_key", gw.DuoIntegrationKey)
		//d.Set("duo_secret_key", gw.DuoSecretKey)		//prevent from reading sensitive info
		d.Set("duo_api_hostname", gw.DuoAPIHostname)
		d.Set("duo_push_mode", gw.DuoPushMode)
		d.Set("ldap_server", gw.LdapServer)
		d.Set("ldap_bind_dn", gw.LdapBindDn)
		//d.Set("ldap_password", gw.LdapPassword)		//prevent from reading sensitive info
		d.Set("ldap_base_dn", gw.LdapBaseDn)
		d.Set("ldap_username_attribute", gw.LdapUserAttr)
		if gw.HASubnet != "" {
			d.Set("ha_subnet", gw.HASubnet)
		}
		if gw.NewZone != "" {
			d.Set("zone", gw.NewZone)
		}
		if gw.SingleAZ != "" {
			if gw.SingleAZ == "yes" {
				d.Set("single_az_ha", "enabled")
			} else if gw.SingleAZ == "no" {
				d.Set("single_az_ha", "disabled")
			} else {
				d.Set("single_az_ha", gw.SingleAZ)
			}
		}
		d.Set("eip", gw.PublicIP)

		// Though go_aviatrix Gateway struct declares VpcSize as only used on gateway creation
		// it is the attribute receiving the instance size of an existing gateway instead of
		// GwSize. (at least in v3.5)
		if gw.GwSize != "" {
			d.Set("vpc_size", gw.GwSize)
		} else {
			if gw.VpcSize != "" {
				d.Set("vpc_size", gw.VpcSize)
			}
		}
		d.Set("public_ip", gw.PublicIP)
		d.Set("cloud_instance_id", gw.CloudnGatewayInstID)
		d.Set("public_dns_server", gw.PublicDnsServer)
		d.Set("security_group_id", gw.GwSecurityGroupID)

		peeringHaSubnet := d.Get("peering_ha_subnet").(string)
		if peeringHaSubnet != "" || gwName == "" {
			peeringHaGateway := &goaviatrix.Gateway{
				AccountName: d.Get("account_name").(string),
				GwName:      d.Get("gw_name").(string) + "-hagw",
			}
			gwHaGw, err := client.GetGateway(peeringHaGateway)
			if err == nil {
				d.Set("cloudn_bkup_gateway_inst_id", gwHaGw.CloudnGatewayInstID)
				d.Set("backup_public_ip", gwHaGw.PublicIP)
				d.Set("peering_ha_subnet", gwHaGw.VpcNet)
				d.Set("peering_ha_eip", gwHaGw.PublicIP)
			} else {
				d.Set("cloudn_bkup_gateway_inst_id", "")
				d.Set("backup_public_ip", "")
				d.Set("peering_ha_subnet", "")
				d.Set("peering_ha_eip", "")
			}
			log.Printf("[TRACE] reading peering HA gateway %s: %#v", d.Get("gw_name").(string), gwHaGw)
		} else {
			d.Set("cloudn_bkup_gateway_inst_id", "")
			d.Set("backup_public_ip", "")
			d.Set("peering_ha_subnet", "")
			d.Set("peering_ha_eip", "")
		}

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
			d.Set("tag_list", tagList)
		} else {
			d.Set("tag_list", tagListStr)
		}

		if gw.VpnStatus == "enabled" && gw.SplitTunnel == "yes" {
			splitTunnel := &goaviatrix.SplitTunnel{
				VpcID: gw.VpcID,
			}
			if gw.EnableElb != "yes" {
				splitTunnel.ElbName = gw.GwName
			} else {
				splitTunnel.ElbName = gw.ElbName
			}
			splitTunnel1, err := client.GetSplitTunnel(splitTunnel)
			if err != nil {
				return fmt.Errorf("unable to read split information for gateway: %v due to %v", gw.GwName, err)
			}
			d.Set("name_servers", splitTunnel1.NameServers)
			d.Set("search_domains", splitTunnel1.SearchDomains)
			d.Set("additional_cidrs", splitTunnel1.AdditionalCidrs)
		} else {
			d.Set("name_servers", "")
			d.Set("search_domains", "")
			d.Set("additional_cidrs", "")
		}
	}
	return nil
}

func resourceAviatrixGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	log.Printf("[INFO] Updating Aviatrix gateway: %#v", d.Get("gw_name").(string))

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
	if d.HasChange("vpc_net") {
		return fmt.Errorf("updating vpc_net is not allowed")
	}
	if d.HasChange("vpn_access") {
		return fmt.Errorf("updating vpn_access is not allowed")
	}
	if d.HasChange("peering_ha_eip") {
		return fmt.Errorf("updating peering_ha_eip is not allowed")
	}
	if d.HasChange("enable_elb") {
		return fmt.Errorf("updating enable_elb is not allowed")
	}
	if d.HasChange("elb_name") {
		return fmt.Errorf("updating elb_name is not allowed")
	}
	//	if d.HasChange("otp_mode") {
	//		return fmt.Errorf("updating otp_mode is not allowed")
	//	}

	gateway := &goaviatrix.Gateway{
		GwName:   d.Get("gw_name").(string),
		GwSize:   d.Get("vpc_size").(string),
		SingleAZ: d.Get("single_az_ha").(string),
	}
	err := client.UpdateGateway(gateway)
	if err != nil {
		return fmt.Errorf("failed to update Aviatrix Gateway: %s", err)
	}
	if d.HasChange("otp_mode") || d.HasChange("enable_ldap") || d.HasChange("saml_enabled") ||
		d.HasChange("okta_token") || d.HasChange("okta_url") || d.HasChange("okta_username_suffix") ||
		d.HasChange("duo_integration_key") || d.HasChange("duo_secret_key") || d.HasChange("duo_api_hostname") ||
		d.HasChange("duo_push_mode") || d.HasChange("ldap_server") || d.HasChange("ldap_bind_dn") ||
		d.HasChange("ldap_password") || d.HasChange("ldap_base_dn") || d.HasChange("ldap_username_attribute") {

		if vpnAccess := d.Get("vpn_access").(string); vpnAccess != "yes" {
			return fmt.Errorf("vpn_access must be set to yes to modify vpn authentication")
		}

		vpn_gw := &goaviatrix.VpnGatewayAuth{
			GwName:             d.Get("gw_name").(string),
			ElbName:            d.Get("elb_name").(string),
			VpcID:              d.Get("vpc_id").(string),
			OtpMode:            d.Get("otp_mode").(string),
			SamlEnabled:        d.Get("saml_enabled").(string),
			OktaToken:          d.Get("okta_token").(string),
			OktaURL:            d.Get("okta_url").(string),
			OktaUsernameSuffix: d.Get("okta_username_suffix").(string),
			DuoIntegrationKey:  d.Get("duo_integration_key").(string),
			DuoSecretKey:       d.Get("duo_secret_key").(string),
			DuoAPIHostname:     d.Get("duo_api_hostname").(string),
			DuoPushMode:        d.Get("duo_push_mode").(string),
			EnableLdap:         d.Get("enable_ldap").(string),
			LdapServer:         d.Get("ldap_server").(string),
			LdapBindDn:         d.Get("ldap_bind_dn").(string),
			LdapPassword:       d.Get("ldap_password").(string),
			LdapBaseDn:         d.Get("ldap_base_dn").(string),
			LdapUserAttr:       d.Get("ldap_username_attribute").(string),
		}

		if vpn_gw.OtpMode != "" && vpn_gw.OtpMode != "2" && vpn_gw.OtpMode != "3" {
			return fmt.Errorf("otp_mode can only be '2' or '3' or empty string")
		}
		if vpn_gw.SamlEnabled == "yes" {
			if vpn_gw.EnableLdap == "yes" || vpn_gw.OtpMode != "" {
				return fmt.Errorf("ldap and mfa can't be configured if saml is enabled")
			}
		}
		if vpn_gw.EnableLdap == "yes" && vpn_gw.OtpMode == "3" {
			return fmt.Errorf("ldap can't be configured along with okta authentication")
		}
		if vpn_gw.EnableLdap == "yes" {
			if vpn_gw.LdapServer == "" {
				return fmt.Errorf("ldap server must be set if ldap is enabled")
			}
			if vpn_gw.LdapBindDn == "" {
				return fmt.Errorf("ldap bind dn must be set if ldap is enabled")
			}
			if vpn_gw.LdapPassword == "" {
				return fmt.Errorf("ldap password must be set if ldap is enabled")
			}
			if vpn_gw.LdapBaseDn == "" {
				return fmt.Errorf("ldap base dn must be set if ldap is enabled")
			}
			if vpn_gw.LdapUserAttr == "" {
				return fmt.Errorf("ldap user attribute must be set if ldap is enabled")
			}
		}
		if vpn_gw.OtpMode == "2" {
			if vpn_gw.DuoIntegrationKey == "" {
				return fmt.Errorf("duo integration key required if otp_mode set to 2")
			}
			if vpn_gw.DuoSecretKey == "" {
				return fmt.Errorf("duo secret key required if otp_mode set to 2")
			}
			if vpn_gw.DuoAPIHostname == "" {
				return fmt.Errorf("duo api hostname required if otp_mode set to 2")
			}
			if vpn_gw.DuoPushMode != "auto" && vpn_gw.DuoPushMode != "token" && vpn_gw.DuoPushMode != "selective" {
				return fmt.Errorf("duo push mode must be set to a valid value (auto, selective, or token)")
			}
			if vpn_gw.EnableLdap == "yes" {
				vpn_gw.AuthType = "duo_ldap_auth"
			} else {
				vpn_gw.AuthType = "duo_auth"
			}
		} else if vpn_gw.OtpMode == "3" {
			if vpn_gw.OktaToken == "" {
				return fmt.Errorf("okta token must be set if otp_mode is set to 3")
			}
			if vpn_gw.OktaURL == "" {
				return fmt.Errorf("okta url must be set if otp_mode is set to 3")
			}
			vpn_gw.AuthType = "okta_auth"
		} else {
			if vpn_gw.EnableLdap == "yes" {
				vpn_gw.AuthType = "ldap_auth"
			} else if vpn_gw.SamlEnabled == "yes" {
				vpn_gw.AuthType = "saml_auth"
			} else {
				vpn_gw.AuthType = "none"
			}
		}
		if vpn_gw.ElbName != "" {
			vpn_gw.LbOrGatewayName = vpn_gw.ElbName
		} else {
			vpn_gw.LbOrGatewayName = vpn_gw.GwName
		}
		err := client.SetVpnGatewayAuthentication(vpn_gw)
		if err != nil {
			return fmt.Errorf("failed to update Aviatrix VPN Gateway Authentication: %s", err)
		}
	}
	if d.HasChange("tag_list") {
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
				tags.TagList = strings.Join(oldTagList, ",")
				err := client.DeleteTags(tags)
				if err != nil {
					return fmt.Errorf("failed to delete tags : %s", err)
				}
			}
			if len(newTagList) != 0 {
				tags.TagList = strings.Join(newTagList, ",")
				err := client.AddTags(tags)
				if err != nil {
					return fmt.Errorf("failed to add tags : %s", err)
				}
			}
		}
		d.SetPartial("tag_list")
	}
	if d.HasChange("split_tunnel") || d.HasChange("additional_cidrs") ||
		d.HasChange("name_servers") || d.HasChange("search_domains") {
		o, n := d.GetChange("split_tunnel")
		if o == nil {
			o = new([]interface{})
		}
		if n == nil {
			n = new([]interface{})
		}
		oST := o.(string)
		nST := n.(string)
		if oST == "" {
			oST = "no"
		}
		if nST == "" {
			nST = "no"
		}
		if nST != "no" && nST != "yes" {
			return fmt.Errorf("split_tunnel is not set correctly")
		}
		if oST != nST || (nST == "yes" && (d.HasChange("additional_cidrs") || d.HasChange("name_servers") || d.HasChange("search_domains"))) {
			sTunnel := &goaviatrix.SplitTunnel{
				SplitTunnel:     nST,
				VpcID:           d.Get("vpc_id").(string),
				ElbName:         d.Get("elb_name").(string),
				AdditionalCidrs: d.Get("additional_cidrs").(string),
				NameServers:     d.Get("name_servers").(string),
				SearchDomains:   d.Get("search_domains").(string),
				SaveTemplate:    "no",
			}
			err = client.ModifySplitTunnel(sTunnel)
			if err != nil {
				return fmt.Errorf("failed to modify split tunnel: %s", err)
			}
		}
	}
	if d.HasChange("single_az_ha") {
		singleAZGateway := &goaviatrix.Gateway{
			GwName:   d.Get("gw_name").(string),
			SingleAZ: d.Get("single_az_ha").(string),
		}
		if singleAZGateway.SingleAZ != "enabled" && singleAZGateway.SingleAZ != "disabled" {
			return fmt.Errorf("[INFO] single_az_ha of gateway: %v is not set correctly", singleAZGateway.GwName)
		}
		if singleAZGateway.SingleAZ == "enabled" {
			log.Printf("[INFO] Enable Single AZ GW HA: %#v", singleAZGateway)
			err := client.EnableSingleAZGateway(gateway)
			if err != nil {
				return fmt.Errorf("failed to create single AZ GW HA: %s", err)
			}
		}
		if singleAZGateway.SingleAZ == "disabled" {
			log.Printf("[INFO] Disable Single AZ GW HA: %#v", singleAZGateway)
			err := client.DisableSingleAZGateway(gateway)
			if err != nil {
				return fmt.Errorf("failed to disable single AZ GW HA: %s", err)
			}
		}
	}
	if d.HasChange("enable_nat") {
		gw := &goaviatrix.Gateway{
			CloudType: d.Get("cloud_type").(int),
			GwName:    d.Get("gw_name").(string),
		}
		o, n := d.GetChange("enable_nat")
		if o == "yes" && n == "no" {
			err := client.DisableSNat(gw)
			if err != nil {
				return fmt.Errorf("failed to disable SNAT: %s", err)
			}
		}
		if o == "no" && n == "yes" {
			err := client.EnableSNat(gw)
			if err != nil {
				return fmt.Errorf("failed to enable SNAT: %s", err)
			}
		}
		d.SetPartial("enable_nat")
	}
	if d.HasChange("vpn_cidr") {
		if d.Get("vpn_access").(string) == "yes" && d.Get("enable_elb").(string) == "yes" {
			gw := &goaviatrix.Gateway{
				CloudType: d.Get("cloud_type").(int),
				GwName:    d.Get("gw_name").(string),
				VpcID:     d.Get("vpc_id").(string),
				ElbName:   d.Get("elb_name").(string),
			}
			_, n := d.GetChange("vpn_cidr")
			gw.VpnCidr = n.(string)
			err := client.UpdateVpnCidr(gw)
			if err != nil {
				return fmt.Errorf("failed to update vpn cidr: %s", err)
			}
		} else {
			log.Printf("[INFO] can't update vpn cidr because elb is disabled for gateway: %#v", gateway.GwName)
		}
		d.SetPartial("enable_nat")
	}
	d.Partial(false)

	d.SetId(gateway.GwName)
	return resourceAviatrixGatewayRead(d, meta)

}

func resourceAviatrixGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	gateway := &goaviatrix.Gateway{
		CloudType: d.Get("cloud_type").(int),
		GwName:    d.Get("gw_name").(string),
	}
	// ha_subnet is for Gateway HA
	if HASubnet := d.Get("ha_subnet").(string); HASubnet != "" {
		log.Printf("[INFO] Deleting Aviatrix gateway HA: %#v", gateway)
		err := client.DisableHaGateway(gateway)
		if err != nil {
			return fmt.Errorf("failed to disable Aviatrix gateway HA: %s", err)
		}
	}
	// peering_ha_subnet is for Peering HA
	if peeringHaSubnet := d.Get("peering_ha_subnet").(string); peeringHaSubnet != "" {
		//Delete backup gateway first
		gateway.GwName += "-hagw"
		log.Printf("[INFO] Deleting Aviatrix Backup Gateway [-hagw]: %#v", gateway)
		err := client.DeleteGateway(gateway)
		if err != nil {
			return fmt.Errorf("failed to delete backup [-hgw] gateway: %s", err)
		}
	}
	gateway.GwName = d.Get("gw_name").(string)
	log.Printf("[INFO] Deleting Aviatrix gateway: %#v", gateway)
	err := client.DeleteGateway(gateway)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix Gateway: %s", err)
	}
	return nil
}
