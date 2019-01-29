package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

// Gateway simple struct to hold gateway details
type Gateway struct {
	AccountName             string `form:"account_name,omitempty" json:"account_name,omitempty"`
	Action                  string `form:"action,omitempty"`
	AdditionalCidrs         string `form:"additional_cidrs,omitempty"`
	AuthMethod              string `form:"auth_method,omitempty" json:"auth_method,omitempty"`
	AllocateNewEip          string `form:"allocate_new_eip,omitempty" json:"allocate_new_eip,omitempty"`
	BkupGatewayZone         string `form:"bkup_gateway_zone,omitempty" json:"bkup_gateway_zone,omitempty"`
	BkupPrivateIP           string `form:"bkup_private_ip,omitempty" json:"bkup_private_ip,omitempty"`
	CID                     string `form:"CID,omitempty"`
	CIDR                    string `form:"cidr,omitempty"`
	ClientCertAuth          string `form:"client_cert_auth,omitempty" json:"client_cert_auth,omitempty"`
	ClientCertSharing       string `form:"client_cert_sharing,omitempty" json:"client_cert_sharing,omitempty"`
	CloudType               int    `form:"cloud_type,omitempty" json:"cloud_type,omitempty"`
	CloudnBkupGatewayInstID string `form:"cloudn_bkup_gateway_inst_id,omitempty" json:"cloudn_bkup_gateway_inst_id,omitempty"`
	CloudnGatewayInstID     string `form:"cloudn_gateway_inst_id,omitempty" json:"cloudn_gateway_inst_id,omitempty"`
	DirectInternet          string `form:"direct_internet,omitempty" json:"direct_internet,omitempty"`
	DockerConsulIP          string `form:"docker_consul_ip,omitempty" json:"docker_consul_ip,omitempty"`
	DockerNtwkCidr          string `form:"docker_ntwk_cidr,omitempty" json:"docker_ntwk_cidr,omitempty"`
	DockerNtwkName          string `form:"docker_ntwk_name,omitempty" json:"docker_ntwk_name,omitempty"`
	DuoAPIHostname          string `form:"duo_api_hostname,omitempty"`
	DuoIntegrationKey       string `form:"duo_integration_key,omitempty"`
	DuoPushMode             string `form:"duo_push_mode,omitempty"`
	DuoSecretKey            string `form:"duo_secret_key,omitempty"`
	Eip                     string `form:"eip,omitempty" json:"eip,omitempty"`
	ElbDNSName              string `form:"elb_dns_name,omitempty" json:"elb_dns_name,omitempty"`
	ElbName                 string `form:"elb_name,omitempty" json:"lb_name,omitempty"`
	ElbState                string `form:"elb_state,omitempty" json:"elb_state,omitempty"`
	EnableClientCertSharing string `form:"enable_client_cert_sharing,omitempty"`
	EnableElb               string `form:"enable_elb,omitempty"`
	EnableLdap              string `form:"enable_ldap,omitempty"`
	DnsServer               string `form:"dns_server,omitempty"`
	PublicDnsServer         string `form:"public_dns_server,omitempty" json:"public_dns_server,omitempty"`
	EnableNat               string `form:"enable_nat,omitempty" json:"enable_nat,omitempty"`
	SingleAZ                string `form:"single_az_ha,omitempty"`
	EnableHybridConnection  bool   `json:"tgw_enabled,omitempty"`
	EnablePbr               string `form:"enable_pbr,omitempty"`
	Expiration              string `form:"expiration,omitempty" json:"expiration,omitempty"`
	GatewayZone             string `form:"gateway_zone,omitempty" json:"gateway_zone,omitempty"`
	GwName                  string `form:"gw_name,omitempty" json:"vpc_name,omitempty"`
	GwSecurityGroupID       string `form:"gw_security_group_id,omitempty" json:"gw_security_group_id,omitempty"`
	GwSize                  string `form:"gw_size,omitempty" json:"vpc_size,omitempty"`
	GwSubnetID              string `form:"gw_subnet_id,omitempty" json:"gw_subnet_id,omitempty"`
	HASubnet                string `form:"ha_subnet,omitempty"`
	PeeringHASubnet         string `form:"public_subnet,omitempty"`
	NewZone                 string `form:"new_zone,omitempty"`
	InstState               string `form:"inst_state,omitempty" json:"inst_state,omitempty"`
	IntraVMRoute            string `form:"intra_vm_route,omitempty" json:"intra_vm_route,omitempty"`
	IsHagw                  string `form:"is_hagw,omitempty" json:"is_hagw,omitempty"`
	LdapAdditionalReq       string `form:"ldap_additional_req,omitempty"`
	LdapBaseDn              string `form:"ldap_base_dn,omitempty" json:"ldap_base_dn,omitempty"`
	LdapBindDn              string `form:"ldap_bind_dn,omitempty" json:"ldap_bind_dn,omitempty"`
	LdapCaCert              string `form:"ldap_ca_cert,omitempty" json:"ldap_ca_cert,omitempty"`
	LdapClientCert          string `form:"ldap_client_cert,omitempty" json:"ldap_client_cert,omitempty"`
	LdapPassword            string `form:"ldap_password,omitempty" json:"ldap_password,omitempty"`
	LdapServer              string `form:"ldap_server,omitempty" json:"ldap_server,omitempty"`
	LdapUseSsl              string `form:"ldap_use_ssl,omitempty" json:"ldap_use_ssl,omitempty"`
	LdapUserAttr            string `form:"ldap_username_attribute,omitempty" json:"ldap_user_attr,omitempty`
	LicenseID               string `form:"license_id,omitempty" json:"license_id,omitempty"`
	MaxConn                 string `form:"max_conn,omitempty"`
	//MaxConnections          string `form:"max_connections,omitempty" json:"max_connections,omitempty"`
	Nameservers        string `form:"nameservers,omitempty"`
	OktaToken          string `form:"okta_token,omitempty" json:"okta_token,omitempty"`
	OktaURL            string `form:"okta_url,omitempty" json:"okta_url,omitempty"`
	OktaUsernameSuffix string `form:"okta_username_suffix,omitempty" json:"okta_username_suffix,omitempty"`
	OtpMode            string `form:"otp_mode,omitempty" json:"otp_mode,omitempty"`
	PbrDefaultGateway  string `form:"pbr_default_gateway,omitempty"`
	PbrEnabled         string `form:"pbr_enabled,omitempty" json:"pbr_enabled,omitempty"`
	PbrLogging         string `form:"pbr_logging,omitempty"`
	PbrSubnet          string `form:"pbr_subnet,omitempty"`
	PrivateIP          string `form:"private_ip,omitempty" json:"private_ip,omitempty"`
	PublicIP           string `form:"public_ip,omitempty" json:"public_ip,omitempty"`
	SamlEnabled        string `form:"saml_enabled,omitempty" json:"saml_enabled,omitempty"`
	SandboxIP          string `form:"sandbox_ip,omitempty" json:"sandbox_ip,omitempty"`
	SaveTemplate       string `form:"save_template,omitempty"`
	SearchDomains      string `form:"search_domains,omitempty"`
	SplitTunnel        string `form:"split_tunnel,omitempty" json:"split_tunnel,omitempty"`
	TagList            string `form:"tags,omitempty"`
	TunnelName         string `form:"tunnel_name,omitempty" json:"tunnel_name,omitempty"`
	TunnelType         string `form:"tunnel_type,omitempty" json:"tunnel_type,omitempty"`
	VendorName         string `form:"vendor_name,omitempty" json:"vendor_name,omitempty"`
	VpcID              string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	VpcNet             string `form:"vpc_net,omitempty" json:"vpc_net,omitempty"`
	VpcRegion          string `form:"vpc_reg,omitempty" json:"vpc_region,omitempty"`
	VpcSplunkIPPort    string `form:"vpc_splunk_ip_port,omitempty" json:"vpc_splunk_ip_port,omitempty"`
	VpcState           string `form:"vpc_state,omitempty" json:"vpc_state,omitempty"`
	VpcType            string `form:"vpc_type,omitempty" json:"vpc_type,omitempty"`
	VpnCidr            string `form:"cidr,omitempty" json:"cidr,omitempty"`
	VpnStatus          string `form:"vpn_access,omitempty" json:"vpn_status,omitempty"`
	Zone               string `form:"zone,omitempty" json:"zone,omitempty"`

	VpcSize string `form:"vpc_size,omitempty" ` //Only use for gateway create
}

type GatewayListResp struct {
	Return  bool      `json:"return"`
	Results []Gateway `json:"results"`
	Reason  string    `json:"reason"`
}

func (c *Client) CreateGateway(gateway *Gateway) error {
	gateway.CID = c.CID
	gateway.Action = "connect_container"
	resp, err := c.Post(c.baseURL, gateway)
	if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		return errors.New(data.Reason)
	}
	return nil
}

func (c *Client) EnableNatGateway(gateway *Gateway) error {
	gateway.CID = c.CID
	gateway.Action = "enable_nat"
	resp, err := c.Post(c.baseURL, gateway)
	if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		return errors.New(data.Reason)
	}
	return nil
}
func (c *Client) EnableSingleAZGateway(gateway *Gateway) error {
	gateway.CID = c.CID
	gateway.Action = "enable_single_az_ha"
	resp, err := c.Post(c.baseURL, gateway)
	if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		return errors.New(data.Reason)
	}
	return nil
}
func (c *Client) EnablePeeringHaGateway(gateway *Gateway) error {
	gateway.CID = c.CID
	gateway.Action = "create_peering_ha_gateway"
	resp, err := c.Post(c.baseURL, gateway)
	if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		return errors.New(data.Reason)
	}
	return nil
}
func (c *Client) EnableHaGateway(gateway *Gateway) error {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=enable_vpc_ha&vpc_name=%s&specific_subnet=%s", c.CID, gateway.GwName, gateway.HASubnet)
	resp, err := c.Get(path, nil)
	if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		return errors.New(data.Reason)
	}
	return nil
}

func (c *Client) DisableHaGateway(gateway *Gateway) error {
	path := c.baseURL + fmt.Sprintf("?CID=%s&action=disable_vpc_ha&vpc_name=%s", c.CID, gateway.GwName)
	resp, err := c.Get(path, nil)
	if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		return errors.New(data.Reason)
	}
	return nil
}

func (c *Client) GetGateway(gateway *Gateway) (*Gateway, error) {
	url := "?CID=%s&action=list_vpcs_summary"
	path := c.baseURL + fmt.Sprintf(url, c.CID)

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	var data GatewayListResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if !data.Return {
		return nil, errors.New(data.Reason)
	}

	gwlist := data.Results
	for i := range gwlist {
		if gwlist[i].GwName == gateway.GwName {
			return &gwlist[i], nil
		}
	}
	log.Printf("Couldn't find Aviatrix gateway %s", gateway.GwName)
	return nil, ErrNotFound
}

func (c *Client) UpdateGateway(gateway *Gateway) error {
	gateway.CID = c.CID
	gateway.Action = "edit_gw_config"
	resp, err := c.Post(c.baseURL, gateway)
	if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		return errors.New(data.Reason)
	}
	return nil
}

func (c *Client) DeleteGateway(gateway *Gateway) error {
	path := c.baseURL + fmt.Sprintf("?action=delete_container&CID=%s&cloud_type=%d&gw_name=%s",
		c.CID, gateway.CloudType, gateway.GwName)
	resp, err := c.Delete(path, nil)

	if err != nil {
		return err
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if !data.Return {
		return errors.New(data.Reason)
	}
	return nil
}
