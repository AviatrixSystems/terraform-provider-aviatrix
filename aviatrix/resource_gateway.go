package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
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
				Type:     schema.TypeInt,
				Required: true,
			},
			"account_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"gw_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_reg": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_size": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_net": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ha_subnet": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_nat": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "no",
			},
			"dns_server": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_dns_server": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpn_access": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "no",
			},
			"vpn_cidr": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"enable_elb": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "no",
			},
			"elb_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"split_tunnel": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "yes",
			},
			"name_servers": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"search_domains": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"additional_cidrs": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"otp_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"saml_enabled": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "no",
			},
			"okta_token": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"okta_url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"okta_username_suffix": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"duo_integration_key": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"duo_secret_key": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"duo_api_hostname": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"duo_push_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"enable_ldap": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "no",
			},
			"ldap_server": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"ldap_bind_dn": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"ldap_password": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"ldap_base_dn": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"ldap_username_attribute": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"peering_ha_subnet": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cloud_instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloudn_bkup_gateway_inst_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"single_az_ha": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "disabled",
			},
			"allocate_new_eip": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "on",
			},
			"eip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tag_list": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Default:  nil,
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
		DnsServer:          d.Get("dns_server").(string),
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
	log.Printf("[INFO] Creating Aviatrix gateway: %#v", gateway)

	err := client.CreateGateway(gateway)
	if err != nil {
		log.Printf("[INFO] failed to create Aviatrix gateway: %#v", gateway)
		return fmt.Errorf("failed to create Aviatrix gateway: %s", err)
	}
	if enableNAT := d.Get("enable_nat").(string); enableNAT == "yes" {
		log.Printf("[INFO] Aviatrix NAT enabled gateway: %#v", gateway)
	}
	if DNSServer := d.Get("dns_server").(string); DNSServer != "" {
		log.Printf("[INFO] Aviatrix gateway DNS server: %#v", gateway)
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
		d.Set("enable_nat", gw.EnableNat)
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
		if gw.DnsServer != "" {
			d.Set("dns_server", gw.DnsServer)
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
		if gw.SplitTunnel != "" {
			d.Set("split_tunnel", gw.SplitTunnel)
		}
		if gw.OtpMode != "" {
			d.Set("otp_mode", gw.OtpMode)
		}
		if gw.SamlEnabled != "" {
			d.Set("saml_enabled", gw.SamlEnabled)
		}
		d.Set("okta_token", gw.OktaToken)
		d.Set("okta_url", gw.OktaURL)
		d.Set("okta_username_suffix", gw.OktaUsernameSuffix)
		d.Set("duo_integration_key", gw.DuoIntegrationKey)
		d.Set("duo_secret_key", gw.DuoSecretKey)
		d.Set("duo_api_hostname", gw.DuoAPIHostname)
		d.Set("duo_push_mode", gw.DuoPushMode)
		d.Set("ldap_server", gw.LdapServer)
		d.Set("ldap_bind_dn", gw.LdapBindDn)
		d.Set("ldap_password", gw.LdapPassword)
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

		if peeringHaSubnet := d.Get("peering_ha_subnet").(string); peeringHaSubnet != "" {
			peeringHaGateway := &goaviatrix.Gateway{
				AccountName: d.Get("account_name").(string),
				GwName:      d.Get("gw_name").(string) + "-hagw",
			}
			gwHaGw, err := client.GetGateway(peeringHaGateway)
			if err == nil {
				d.Set("cloudn_bkup_gateway_inst_id", gwHaGw.CloudnGatewayInstID)
				d.Set("backup_public_ip", gwHaGw.PublicIP)
				d.Set("peering_ha_subnet", gwHaGw.VpcNet)
			} else {
				d.Set("cloudn_bkup_gateway_inst_id", "")
				d.Set("backup_public_ip", "")
				d.Set("peering_ha_subnet", "")
			}
			log.Printf("[TRACE] reading peering HA gateway %s: %#v", d.Get("gw_name").(string), gwHaGw)
		} else {
			d.Set("cloudn_bkup_gateway_inst_id", "")
			d.Set("backup_public_ip", "")
			d.Set("peering_ha_subnet", "")
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
		d.Set("tag_list", tagList)

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
	if d.HasChange("enable_elb") {
		return fmt.Errorf("updating enable_elb is not allowed")
	}
	if d.HasChange("elb_name") {
		return fmt.Errorf("updating elb_name is not allowed")
	}
	if d.HasChange("otp_mode") {
		return fmt.Errorf("updating otp_mode is not allowed")
	}

	gateway := &goaviatrix.Gateway{
		GwName:   d.Get("gw_name").(string),
		GwSize:   d.Get("vpc_size").(string),
		SingleAZ: d.Get("single_az_ha").(string),
	}
	err := client.UpdateGateway(gateway)
	if err != nil {
		return fmt.Errorf("failed to update Aviatrix Gateway: %s", err)
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
		oldTagList := goaviatrix.ExpandStringList(os)
		if len(oldTagList) != 0 {
			tags.TagList = strings.Join(oldTagList, ",")
			err := client.DeleteTags(tags)
			if err != nil {
				return fmt.Errorf("failed to delete tags : %s", err)
			}
		}
		newTagList := goaviatrix.ExpandStringList(ns)
		if len(newTagList) != 0 {
			tags.TagList = strings.Join(newTagList, ",")
			err = client.AddTags(tags)
			if err != nil {
				return fmt.Errorf("failed to add tags : %s", err)
			}
		}
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
			log.Printf("[INFO] elb is not enabled for gateway: %#v", gateway.GwName)
		}
		d.SetPartial("enable_nat")
	}
	d.Partial(false)

	d.SetId(gateway.GwName)
	return nil
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
