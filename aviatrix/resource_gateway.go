package aviatrix

import (
	"fmt"
	"log"

	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAviatrixGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixGatewayCreate,
		Read:   resourceAviatrixGatewayRead,
		Update: resourceAviatrixGatewayUpdate,
		Delete: resourceAviatrixGatewayDelete,

		Schema: map[string]*schema.Schema{
			"cloud_type": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"account_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"gw_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_reg": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_size": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_net": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ha_subnet": &schema.Schema{
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
	}

	log.Printf("[INFO] Creating Aviatrix gateway: %#v", gateway)

	err := client.CreateGateway(gateway)
	if err != nil {
		log.Printf("[INFO] Failed to create Aviatrix gateway: %#v", gateway)
		return fmt.Errorf("Failed to create Aviatrix gateway: %s", err)
	}
	if enable_nat := d.Get("enable_nat").(string); enable_nat == "yes" {
		log.Printf("[INFO] Aviatrix NAT enabled gateway: %#v", gateway)
	}
	if dns_server := d.Get("dns_server").(string); dns_server != "" {
		log.Printf("[INFO] Aviatrix gateway DNS server: %#v", gateway)
	}
	// single_AZ enabled for Gateway. https://docs.aviatrix.com/HowTos/gateway.html#high-availability
	if single_az_ha := d.Get("single_az_ha").(string); single_az_ha == "enabled" {
		single_az_gateway := &goaviatrix.Gateway{
			GwName:   d.Get("gw_name").(string),
			SingleAZ: d.Get("single_az_ha").(string),
		}
		log.Printf("[INFO] Enable Single AZ GW HA: %#v", single_az_gateway)
		err := client.EnableSingleAZGateway(gateway)
		if err != nil {
			return fmt.Errorf("Failed to create single AZ GW HA: %s", err)
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
				return fmt.Errorf("Failed to auto-cleanup failed gateway: %s", del_err)
			}
			return fmt.Errorf("Failed to create GW HA: %s", err)
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
			return fmt.Errorf("Failed to create peering HA: %s", err)
		}
	}
	d.SetId(gateway.GwName)
	return resourceAviatrixGatewayRead(d, meta)
}

func resourceAviatrixGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	gateway := &goaviatrix.Gateway{
		AccountName: d.Get("account_name").(string),
		GwName:      d.Get("gw_name").(string),
		SingleAZ:    d.Get("single_az_ha").(string),
	}
	gw, err := client.GetGateway(gateway)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find Aviatrix Gateway: %s", err)
	}
	log.Printf("[TRACE] reading gateway %s: %#v", d.Get("gw_name").(string), gw)
	if gw != nil {
		d.Set("vpc_size", gw.VpcSize)
		d.Set("public_ip", gw.PublicIP)
		d.Set("cloud_instance_id", gw.CloudnGatewayInstID)
		d.Set("public_dns_server", gw.PublicDnsServer)
		if public_subnet := d.Get("public_subnet").(string); public_subnet != "" {
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
		return fmt.Errorf("Failed to update Aviatrix Gateway: %s", err)
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
	if ha_subnet := d.Get("ha_subnet").(string); ha_subnet != "" {
		log.Printf("[INFO] Deleting Aviatrix gateway HA: %#v", gateway)
		err := client.DisableHaGateway(gateway)
		if err != nil {
			return fmt.Errorf("Failed to disable Aviatrix gateway HA: %s", err)
		}
	}
	// public_subnet is for Peering HA
	if public_subnet := d.Get("public_subnet").(string); public_subnet != "" {
		//Delete backup gateway first
		gateway.GwName += "-hagw"
		log.Printf("[INFO] Deleting Aviatrix Backup Gateway [-hagw]: %#v", gateway)
		err := client.DeleteGateway(gateway)
		if err != nil {
			return fmt.Errorf("Failed to delete backup [-hgw] gateway: %s", err)
		}
	}
	gateway.GwName = d.Get("gw_name").(string)
	log.Printf("[INFO] Deleting Aviatrix gateway: %#v", gateway)
	err := client.DeleteGateway(gateway)
	if err != nil {
		return fmt.Errorf("Failed to delete Aviatrix Gateway: %s", err)
	}
	return nil
}
