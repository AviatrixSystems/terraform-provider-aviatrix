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
			},
			"cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_elb": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"elb_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"split_tunnel": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"otp_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"saml_enabled": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"okta_token": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"okta_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"okta_username_suffix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"duo_integration_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"duo_secret_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"duo_api_hostname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"duo_push_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_ldap": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ldap_server": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ldap_bind_dn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ldap_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ldap_base_dn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ldap_username_attribute": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_subnet": {
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
			},
			"allocate_new_eip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"eip": {
				Type:     schema.TypeString,
				Optional: true,
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
		VpnCidr:            d.Get("cidr").(string),
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
		PeeringHASubnet:    d.Get("public_subnet").(string),
		NewZone:            d.Get("zone").(string),
		SingleAZ:           d.Get("single_az_ha").(string),
		AllocateNewEip:     d.Get("allocate_new_eip").(string),
		Eip:                d.Get("eip").(string),
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
	// public_subnet is for Peering HA Gateway. https://docs.aviatrix.com/HowTos/gateway.html#high-availability
	if public_subnet := d.Get("public_subnet").(string); public_subnet != "" {
		ha_gateway := &goaviatrix.Gateway{
			GwName:          d.Get("gw_name").(string),
			PeeringHASubnet: d.Get("public_subnet").(string),
			NewZone:         d.Get("zone").(string),
		}
		log.Printf("[INFO] Enable peering HA: %#v", ha_gateway)
		err := client.EnablePeeringHaGateway(ha_gateway)
		if err != nil {
			return fmt.Errorf("failed to create peering HA: %s", err)
		}
	}
	d.SetId(gateway.GwName)

	return resourceAviatrixGatewayRead(d, meta)
}

func resourceAviatrixGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	
	gwname := d.Get("gw_name").(string)
	// If it is an import only Id is set
	if gwname == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no gateway name received. Import Id is %s",id)
		if strings.Contains(id,"@") {
		    substr := strings.Split(id,"@")
		    account_name := substr[1]
		    gateway_name := substr[0]
		    log.Printf("[INFO] Importing %s gateway in %s account",gateway_name,account_name)
		    d.Set("account_name",account_name)
		    d.Set("gw_name",gateway_name)
		    // Terraform must locate a resource declared with the same Id returned
		    d.SetId(gateway_name)
		} else {
			return fmt.Errorf("Id must be in the following format: <gateway name>@<account name>")
		}
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

		d.Set("cloud_type",gw.CloudType)
		d.Set("account_name",gw.AccountName)
		d.Set("gw_name",gw.GwName)
		d.Set("vpc_id",strings.Split(gw.VpcID,"~~")[0])
		d.Set("vpc_reg",gw.VpcRegion)
		if gw.VpcNet != "" { d.Set("vpc_net",gw.VpcNet) }
		if gw.EnableNat != ""  { d.Set("enable_nat",gw.EnableNat) }
		if gw.DnsServer != "" { d.Set("dns_server",gw.DnsServer) }
		if gw.VpnStatus != "" { d.Set("vpn_access",gw.VpnStatus) }
		if gw.VpnCidr != "" { d.Set("cidr",gw.VpnCidr) }
		if gw.ElbState == "enabled" {
		    d.Set("enable_elb","yes")
		    elb_name := d.Get("elb_name")
		    // Versions prior to 4.0 won't return elb_name, so deduce it from elb_dns_name
		    if elb_name == "" {
		    	elb_dns_name := gw.ElbDNSName
		    	log.Printf("[INFO] Controllers prior to 4.0 do not return elb_name. Deducing from elb_dns_name")
		    	if elb_dns_name != "" {
		    		hostname := strings.Split(elb_dns_name,".")[0]
		    		// AWS adds - followed by a random string after the name given to the ELB
		    		parts := strings.Split(hostname,"-")
		    		// Remove random string added by AWS 
		    		parts[len(parts)-1] = ""
		    		parts = parts[:len(parts)-1]
		    		// Join again using - in case it was used in the ELB Name
		    		elb_name = strings.Join(parts,"-")
		    	} else {
		    		return fmt.Errorf("Neither elb_name or elb_dns_name returned by the API in an ELB enabled gateway")
		    	}
		    }
		    d.Set("elb_name", elb_name)
		} else {
		    d.Set("enable_elb","no")
		}
		if gw.SplitTunnel != "" { d.Set("split_tunnel",gw.SplitTunnel) }
		if gw.OtpMode != "" { d.Set("otp_mode",gw.OtpMode) }
		if gw.SamlEnabled != "" { d.Set("saml_enabled",gw.SamlEnabled) }
		if gw.OktaToken != "" { d.Set("okta_token",gw.OktaToken) }
		if gw.OktaURL != "" { d.Set("okta_url",gw.OktaURL) }
		if gw.OktaUsernameSuffix != "" { d.Set("okta_username_suffix",gw.OktaUsernameSuffix) }
		if gw.DuoIntegrationKey != "" { d.Set("duo_integration_key",gw.DuoIntegrationKey) }
		if gw.DuoSecretKey != "" { d.Set("duo_secret_key",gw.DuoSecretKey) }
		if gw.DuoAPIHostname != "" { d.Set("duo_api_hostname",gw.DuoAPIHostname) }
		if gw.DuoPushMode != "" { d.Set("duo_push_mode",gw.DuoPushMode) }
		if gw.EnableLdap != "" { d.Set("enable_ldap",gw.EnableLdap) }
		if gw.LdapServer != "" { d.Set("ldap_server",gw.LdapServer) }
		if gw.LdapBindDn!= "" { d.Set("ldap_bind_dn",gw.LdapBindDn) }
		if gw.LdapPassword != "" { d.Set("ldap_password",gw.LdapPassword) }
		if gw.LdapBaseDn != "" { d.Set("ldap_base_dn",gw.LdapBaseDn) }
		if gw.LdapUserAttr != "" { d.Set("ldap_username_attribute",gw.LdapUserAttr) }
		if gw.HASubnet != "" { d.Set("ha_subnet",gw.HASubnet) }
		if gw.PeeringHASubnet != "" { d.Set("public_subnet",gw.PeeringHASubnet) }
		if gw.NewZone != "" { d.Set("zone",gw.NewZone) }
		if gw.SingleAZ != "" { d.Set("single_az_ha",gw.SingleAZ) }
	//	d.Set("allocate_new_eip",gw.AllocateNewEip)
		if gw.Eip != "" { d.Set("eip",gw.Eip) }

        // Though go_aviatrix Gateway struct declares VpcSize as only used on gateway creation
        // it is the attribute receiving the instance size of an existing gateway instead of
        // GwSize. (at least in v3.5)
        if gw.GwSize != "" {
        	d.Set("vpc_size",gw.GwSize)
        } else {
		   if gw.VpcSize != "" { d.Set("vpc_size", gw.VpcSize) }
        }   
		d.Set("public_ip", gw.PublicIP)
		d.Set("cloud_instance_id", gw.CloudnGatewayInstID)
		d.Set("public_dns_server", gw.PublicDnsServer)
		d.Set("security_group_id", gw.GwSecurityGroupID)

		if publicSubnet := d.Get("public_subnet").(string); publicSubnet != "" {
			gateway.GwName += "-hagw"
			gw, err := client.GetGateway(gateway)
			if err == nil {
				d.Set("cloudn_bkup_gateway_inst_id", gw.CloudnGatewayInstID)
				d.Set("backup_public_ip", gw.PublicIP)
			}
			log.Printf("[TRACE] reading peering HA gateway %s: %#v", d.Get("gw_name").(string), gw)
		}
	}
	return nil
}

func resourceAviatrixGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	gateway := &goaviatrix.Gateway{
		GwName:   d.Get("gw_name").(string),
		GwSize:   d.Get("vpc_size").(string),
		SingleAZ: d.Get("single_az_ha").(string),
	}

	log.Printf("[INFO] Updating Aviatrix gateway: %#v", gateway)

	err := client.UpdateGateway(gateway)
	if err != nil {
		return fmt.Errorf("failed to update Aviatrix Gateway: %s", err)
	}
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
	// public_subnet is for Peering HA
	if publicSubnet := d.Get("public_subnet").(string); publicSubnet != "" {
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
