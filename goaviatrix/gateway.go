package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Gateway simple struct to hold gateway details
type Gateway struct {
	AccountName                 string `form:"account_name,omitempty" json:"account_name,omitempty"`
	Action                      string `form:"action,omitempty"`
	AdditionalCidrs             string `form:"additional_cidrs,omitempty"`
	AuthMethod                  string `form:"auth_method,omitempty" json:"auth_method,omitempty"`
	AllocateNewEip              string `form:"allocate_new_eip,omitempty"`
	AllocateNewEipRead          bool   `json:"newly_allocated_eip,omitempty"`
	BkupGatewayZone             string `form:"bkup_gateway_zone,omitempty" json:"bkup_gateway_zone,omitempty"`
	BkupPrivateIP               string `form:"bkup_private_ip,omitempty" json:"bkup_private_ip,omitempty"`
	CID                         string `form:"CID,omitempty"`
	CIDR                        string `form:"cidr,omitempty"`
	ClientCertAuth              string `form:"client_cert_auth,omitempty" json:"client_cert_auth,omitempty"`
	ClientCertSharing           string `form:"client_cert_sharing,omitempty" json:"client_cert_sharing,omitempty"`
	CloudType                   int    `form:"cloud_type,omitempty" json:"cloud_type,omitempty"`
	CloudnBkupGatewayInstID     string `form:"cloudn_bkup_gateway_inst_id,omitempty" json:"cloudn_bkup_gateway_inst_id,omitempty"`
	CloudnGatewayInstID         string `form:"cloudn_gateway_inst_id,omitempty" json:"cloudn_gateway_inst_id,omitempty"`
	ConnectedTransit            string `json:"connected_transit,omitempty"`
	DirectInternet              string `form:"direct_internet,omitempty" json:"direct_internet,omitempty"`
	DockerConsulIP              string `form:"docker_consul_ip,omitempty" json:"docker_consul_ip,omitempty"`
	DockerNtwkCidr              string `form:"docker_ntwk_cidr,omitempty" json:"docker_ntwk_cidr,omitempty"`
	DockerNtwkName              string `form:"docker_ntwk_name,omitempty" json:"docker_ntwk_name,omitempty"`
	DuoAPIHostname              string `form:"duo_api_hostname,omitempty" json:"duo_api_hostname,omitempty"`
	DuoIntegrationKey           string `form:"duo_integration_key,omitempty" json:"duo_integration_key,omitempty"`
	DuoPushMode                 string `form:"duo_push_mode,omitempty" json:"duo_push_mode,omitempty"`
	DuoSecretKey                string `form:"duo_secret_key,omitempty" json:"duo_secret_key,omitempty"`
	Eip                         string `form:"eip,omitempty" json:"eip,omitempty"`
	ElbDNSName                  string `form:"elb_dns_name,omitempty" json:"elb_dns_name,omitempty"`
	ElbName                     string `form:"elb_name,omitempty" json:"lb_name,omitempty"`
	ElbState                    string `form:"elb_state,omitempty" json:"elb_state,omitempty"`
	VpnProtocol                 string `form:"elb_protocol,omitempty"`
	EnableClientCertSharing     string `form:"enable_client_cert_sharing,omitempty"`
	EnableElb                   string `form:"enable_elb,omitempty"`
	EnableLdap                  string `form:"enable_ldap,omitempty"`
	EnableLdapRead              bool   `json:"enable_ldap,omitempty"`
	EnableVpcDnsServer          string `json:"use_vpc_dns,omitempty"`
	DnsServer                   string `form:"dns_server,omitempty"`
	PublicDnsServer             string `form:"public_dns_server,omitempty" json:"public_dns_server,omitempty"`
	EnableNat                   string `form:"enable_nat,omitempty" json:"enable_nat,omitempty"`
	SingleAZ                    string `form:"single_az_ha,omitempty" json:"single_az_ha,omitempty"`
	EnableHybridConnection      bool   `json:"tgw_enabled,omitempty"`
	EnablePbr                   string `form:"enable_pbr,omitempty"`
	Expiration                  string `form:"expiration,omitempty" json:"expiration,omitempty"`
	GatewayZone                 string `form:"gateway_zone,omitempty" json:"gateway_zone,omitempty"`
	GwName                      string `form:"gw_name,omitempty" json:"vpc_name,omitempty"`
	GwSecurityGroupID           string `form:"gw_security_group_id,omitempty" json:"gw_security_group_id,omitempty"`
	GwSize                      string `form:"gw_size,omitempty" json:"vpc_size,omitempty"`
	GwSubnetID                  string `form:"gw_subnet_id,omitempty" json:"gw_subnet_id,omitempty"`
	PeeringHASubnet             string `form:"public_subnet,omitempty"`
	NewZone                     string `form:"new_zone,omitempty"`
	NewSubnet                   string `form:"new_subnet,omitempty"`
	InsaneMode                  string `form:"insane_mode,omitempty" json:"high_perf,omitempty"`
	InstState                   string `form:"inst_state,omitempty" json:"inst_state,omitempty"`
	IntraVMRoute                string `form:"intra_vm_route,omitempty" json:"intra_vm_route,omitempty"`
	IsHagw                      string `form:"is_hagw,omitempty" json:"is_hagw,omitempty"`
	LdapAdditionalReq           string `form:"ldap_additional_req,omitempty"`
	LdapBaseDn                  string `form:"ldap_base_dn,omitempty" json:"ldap_base_dn,omitempty"`
	LdapBindDn                  string `form:"ldap_bind_dn,omitempty" json:"ldap_bind_dn,omitempty"`
	LdapCaCert                  string `form:"ldap_ca_cert,omitempty" json:"ldap_ca_cert,omitempty"`
	LdapClientCert              string `form:"ldap_client_cert,omitempty" json:"ldap_client_cert,omitempty"`
	LdapPassword                string `form:"ldap_password,omitempty" json:"ldap_password,omitempty"`
	LdapServer                  string `form:"ldap_server,omitempty" json:"ldap_server,omitempty"`
	LdapUseSsl                  string `form:"ldap_use_ssl,omitempty" json:"ldap_use_ssl,omitempty"`
	LdapUserAttr                string `form:"ldap_username_attribute,omitempty" json:"ldap_username_attribute,omitempty"`
	LicenseID                   string `form:"license_id,omitempty" json:"license_id,omitempty"`
	MaxConn                     string `form:"max_conn,omitempty" json:"max_connections,omitempty"`
	NameServers                 string `form:"nameservers,omitempty"`
	OktaToken                   string `form:"okta_token,omitempty" json:"okta_token,omitempty"`
	OktaURL                     string `form:"okta_url,omitempty" json:"okta_url,omitempty"`
	OktaUsernameSuffix          string `form:"okta_username_suffix,omitempty" json:"okta_username_suffix,omitempty"`
	OtpMode                     string `form:"otp_mode,omitempty" json:"otp_mode,omitempty"`
	PbrDefaultGateway           string `form:"pbr_default_gateway,omitempty"`
	PbrEnabled                  string `form:"pbr_enabled,omitempty" json:"pbr_enabled,omitempty"`
	PbrLogging                  string `form:"pbr_logging,omitempty"`
	PbrSubnet                   string `form:"pbr_subnet,omitempty"`
	PrivateIP                   string `form:"private_ip,omitempty" json:"private_ip,omitempty"`
	PublicIP                    string `form:"public_ip,omitempty" json:"public_ip,omitempty"`
	SamlEnabled                 string `form:"saml_enabled,omitempty" json:"saml_enabled,omitempty"`
	SandboxIP                   string `form:"sandbox_ip,omitempty" json:"sandbox_ip,omitempty"`
	SaveTemplate                string `form:"save_template,omitempty"`
	SearchDomains               string `form:"search_domains,omitempty"`
	SplitTunnel                 string `form:"split_tunnel,omitempty" json:"split_tunnel,omitempty"`
	SpokeVpc                    string `json:"spoke_vpc,omitempty"`
	TagList                     string `form:"tags,omitempty"`
	TransitGwName               string `form:"transit_gw_name,omitempty" json:"transit_gw_name,omitempty"`
	EgressTransitGwName         string `form:"egress_transit_gw_name,omitempty" json:"egress_transit_gw_name,omitempty"`
	TunnelName                  string `form:"tunnel_name,omitempty" json:"tunnel_name,omitempty"`
	TunnelType                  string `form:"tunnel_type,omitempty" json:"tunnel_type,omitempty"`
	VendorName                  string `form:"vendor_name,omitempty" json:"vendor_name,omitempty"`
	VpcID                       string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	VpcNet                      string `form:"vpc_net,omitempty" json:"public_subnet,omitempty"`
	VpcRegion                   string `form:"vpc_reg,omitempty" json:"vpc_region,omitempty"`
	VpcSplunkIPPort             string `form:"vpc_splunk_ip_port,omitempty" json:"vpc_splunk_ip_port,omitempty"`
	VpcState                    string `form:"vpc_state,omitempty" json:"vpc_state,omitempty"`
	VpcType                     string `form:"vpc_type,omitempty" json:"vpc_type,omitempty"`
	VpnCidr                     string `form:"cidr,omitempty" json:"vpn_cidr,omitempty"`
	VpnStatus                   string `form:"vpn_access,omitempty" json:"vpn_status,omitempty"`
	Zone                        string `form:"zone,omitempty" json:"zone,omitempty"`
	VpcSize                     string `form:"vpc_size,omitempty" ` //Only use for gateway create
	DMZEnabled                  string `json:"dmz_enabled,omitempty"`
	EnableActiveMesh            string `form:"enable_activemesh,omitempty" json:"enable_activemesh,omitempty"`
	EnableVpnNat                bool   `form:"vpn_nat,omitempty" `
	EnableDesignatedGateway     string `form:"designated_gateway,omitempty" json:"designated_gateway,omitempty"`
	AdditionalCidrsDesignatedGw string `form:"additional_cidr_list,omitempty" json:"summarized_cidrs,omitempty"`
	EnableEncryptVolume         bool   `json:"gw_enc,omitempty"`
	CustomerManagedKeys         string `form:"customer_managed_keys,omitempty" json:"customer_managed_keys,omitempty"`
	SnatMode                    string `form:"mode,omitempty" json:"snat_target,omitempty"`
	SnatPolicy                  []PolicyRule
	PolicyList                  string `form:"policy_list,omitempty"`
	GatewayName                 string `form:"gateway_name,omitempty"`
	DnatPolicy                  []PolicyRule
	CustomizedSpokeVpcRoutes    []string `json:"customized_cidr_list,omitempty"`
	FilteredSpokeVpcRoutes      []string `json:"filtering_cidr_list,omitempty"`
	AdvertisedSpokeRoutes       []string `json:"advertise_cidr_list,omitempty"`
	IncludeCidrList             []string `json:"include_cidr_list,omitempty"`
	ExcludeCidrList             []string `json:"exclude_cidr_list,omitempty"`
	EnableTransitFireNet        string   `json:"firenet_enabled,omitempty"`
	LearnedCidrsApproval        string   `json:"learned_cidrs_approval,omitempty"`
	Dns                         string   `json:"dns,omitempty"`
	EncVolume                   string   `form:"enc_volume,omitempty"`
	SyncSNATToHA                string   `form:"sync_snat_to_ha,omitempty"`
	SyncDNATToHA                string   `form:"sync_dnat_to_ha,omitempty"`
	EnableMonitorGWSubnets      bool     `form:"enable_monitor_gw_subnets,omitempty" json:"enable_monitor_gw_subnets,omitempty"`
	MonitorExcludeList          string   `form:"monitor_exclude_list,omitempty"`
}

type PolicyRule struct {
	SrcIP      string `form:"src_ip,omitempty" json:"src_ip,omitempty"`
	SrcPort    string `form:"src_port,omitempty" json:"src_port,omitempty"`
	DstIP      string `form:"dst_ip,omitempty" json:"dst_ip,omitempty"`
	DstPort    string `form:"dst_port,omitempty" json:"dst_port,omitempty"`
	Protocol   string `form:"protocol,omitempty" json:"protocol,omitempty"`
	Interface  string `form:"interface,omitempty" json:"interface,omitempty"`
	Connection string `form:"connection,omitempty" json:"connection,omitempty"`
	Mark       string `form:"mark,omitempty" json:"mark,omitempty"`
	NewSrcIP   string `form:"new_src_ip,omitempty" json:"new_src_ip,omitempty"`
	NewSrcPort string `form:"new_src_port,omitempty" json:"new_src_port,omitempty"`
	ExcludeRTB string `form:"exclude_rtb,omitempty" json:"exclude_rtb,omitempty"`
	NewDstIP   string `form:"new_dst_ip,omitempty" json:"new_dst_ip,omitempty"`
	NewDstPort string `form:"new_dst_port,omitempty" json:"new_dst_port,omitempty"`
}

type GatewayDetail struct {
	AccountName                  string       `form:"account_name,omitempty" json:"account_name,omitempty"`
	Action                       string       `form:"action,omitempty"`
	GwName                       string       `form:"gw_name,omitempty" json:"vpc_name,omitempty"`
	DMZEnabled                   bool         `json:"dmz_enabled,omitempty"`
	EnableAdvertiseTransitCidr   string       `json:"advertise_transit_cidr,omitempty"`
	BgpManualSpokeAdvertiseCidrs []string     `json:"bgp_manual_spoke_advertise_cidrs,omitempty"`
	VpnNat                       bool         `json:"vpn_nat,omitempty"`
	SnatPolicy                   []PolicyRule `json:"snat_ip_port_list,omitempty"`
	DnatPolicy                   []PolicyRule `json:"dnat_ip_port_list,omitempty"`
	Elb                          ElbDetail    `json:"elb,omitempty"`
	EnableEgressTransitFireNet   bool         `json:"egress_transit,omitempty"`
	EnableFireNet                bool         `json:"firenet_enabled,omitempty"`
	EnableTransitFireNet         bool         `json:"transit_firenet_enabled,omitempty"`
	LearnedCidrsApproval         string       `json:"learned_cidrs_approval,omitempty"`
	SyncSNATToHA                 bool         `json:"sync_snat_to_ha,omitempty"`
	SyncDNATToHA                 bool         `json:"sync_dnat_to_ha,omitempty"`
	GwZone                       string       `json:"gw_zone,omitempty"`
	TransitGwName                string       `json:"transit_gw_name,omitempty"`
	EgressTransitGwName          string       `json:"egress_transit_gw_name,omitempty"`
	RouteTables                  []string     `json:"spoke_rtb_list,omitempty"`
}

type ElbDetail struct {
	VpnProtocol string `json:"elb_protocol,omitempty"`
}

type ListTransitFireNetPolicyResp struct {
	Return  bool                       `json:"return"`
	Results []TransitFireNetPolicyEdit `json:"results"`
	Reason  string                     `json:"reason"`
}

type VpnGatewayAuth struct { // Used for set_vpn_gateway_authentication rest api call
	Action             string `form:"action,omitempty"`
	AuthType           string `form:"auth_type,omitempty" json:"auth_type,omitempty"`
	CID                string `form:"CID,omitempty"`
	DuoAPIHostname     string `form:"duo_api_hostname,omitempty" json:"duo_api_hostname,omitempty"`
	DuoIntegrationKey  string `form:"duo_integration_key,omitempty" json:"duo_integration_key,omitempty"`
	DuoPushMode        string `form:"duo_push_mode,omitempty" json:"duo_push_mode,omitempty"`
	DuoSecretKey       string `form:"duo_secret_key,omitempty" json:"duo_secret_key,omitempty"`
	EnableLdap         string `form:"enable_ldap,omitempty"`
	LbOrGatewayName    string `form:"lb_or_gateway_name,omitempty" json:"lb_or_gateway_name,omitempty"`
	LdapAdditionalReq  string `form:"ldap_additional_req,omitempty"`
	LdapBaseDn         string `form:"ldap_base_dn,omitempty" json:"ldap_base_dn,omitempty"`
	LdapBindDn         string `form:"ldap_bind_dn,omitempty" json:"ldap_bind_dn,omitempty"`
	LdapCaCert         string `form:"ldap_ca_cert,omitempty" json:"ldap_ca_cert,omitempty"`
	LdapClientCert     string `form:"ldap_client_cert,omitempty" json:"ldap_client_cert,omitempty"`
	LdapPassword       string `form:"ldap_password,omitempty" json:"ldap_password,omitempty"`
	LdapServer         string `form:"ldap_server,omitempty" json:"ldap_server,omitempty"`
	LdapUseSsl         string `form:"ldap_use_ssl,omitempty" json:"ldap_use_ssl,omitempty"`
	LdapUserAttr       string `form:"ldap_username_attribute,omitempty" json:"ldap_username_attribute,omitempty"`
	OktaToken          string `form:"okta_token,omitempty" json:"okta_token,omitempty"`
	OktaURL            string `form:"okta_url,omitempty" json:"okta_url,omitempty"`
	OktaUsernameSuffix string `form:"okta_username_suffix,omitempty" json:"okta_username_suffix,omitempty"`
	OtpMode            string `form:"otp_mode,omitempty" json:"otp_mode,omitempty"`
	SamlEnabled        string `form:"saml_enabled,omitempty" json:"saml_enabled,omitempty"`
	VpcID              string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
}

type GatewayListResp struct {
	Return  bool      `json:"return"`
	Results []Gateway `json:"results"`
	Reason  string    `json:"reason"`
}

type GatewayDetailApiResp struct {
	Return  bool          `json:"return"`
	Results GatewayDetail `json:"results"`
	Reason  string        `json:"reason"`
}

func (c *Client) CreateGateway(gateway *Gateway) error {
	gateway.CID = c.CID
	gateway.Action = "connect_container"
	resp, err := c.Post(c.baseURL, gateway)
	if err != nil {
		return errors.New("HTTP Post connect_container failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode connect_container failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API connect_container Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) EnableNatGateway(gateway *Gateway) error {
	gateway.CID = c.CID
	gateway.Action = "enable_nat"
	resp, err := c.Post(c.baseURL, gateway)
	if err != nil {
		return errors.New("HTTP Post enable_nat failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode enable_nat failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API enable_nat Post failed: " + data.Reason)
	}
	return nil
}
func (c *Client) EnableSingleAZGateway(gateway *Gateway) error {
	gateway.CID = c.CID
	gateway.Action = "enable_single_az_ha"
	resp, err := c.Post(c.baseURL, gateway)
	if err != nil {
		return errors.New("HTTP Post enable_single_az_ha failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode enable_single_az_ha failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API enable_single_az_ha Post failed: " + data.Reason)
	}
	return nil
}
func (c *Client) EnablePeeringHaGateway(gateway *Gateway) error {
	gateway.CID = c.CID
	gateway.Action = "create_peering_ha_gateway"
	resp, err := c.Post(c.baseURL, gateway)
	if err != nil {
		return errors.New("HTTP Post create_peering_ha_gateway failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode create_peering_ha_gateway failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API create_peering_ha_gateway Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DisableSingleAZGateway(gateway *Gateway) error {
	gateway.CID = c.CID
	gateway.Action = "disable_single_az_ha"
	resp, err := c.Post(c.baseURL, gateway)
	if err != nil {
		return errors.New("HTTP Post disable_single_az_ha failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode disable_single_az_ha failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API disable_single_az_ha Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetGateway(gateway *Gateway) (*Gateway, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for list_vpcs_summary") + err.Error())
	}
	listVpcSummary := url.Values{}
	listVpcSummary.Add("CID", c.CID)
	listVpcSummary.Add("action", "list_vpcs_summary")
	Url.RawQuery = listVpcSummary.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get list_vpcs_summary failed: " + err.Error())
	}
	var data GatewayListResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_vpcs_summary failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return nil, errors.New("Rest API list_vpcs_summary Get failed: " + data.Reason)
	}

	gwList := data.Results
	for i := range gwList {
		if gwList[i].GwName == gateway.GwName {
			return &gwList[i], nil
		}
	}
	log.Errorf("Couldn't find Aviatrix gateway %s", gateway.GwName)
	return nil, ErrNotFound
}

func (c *Client) GetGatewayDetail(gateway *Gateway) (*GatewayDetail, error) {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.New(("url Parsing failed for list_vpc_by_name") + err.Error())
	}
	listVpcByName := url.Values{}
	listVpcByName.Add("CID", c.CID)
	listVpcByName.Add("action", "list_vpc_by_name")
	listVpcByName.Add("vpc_name", gateway.GwName)
	Url.RawQuery = listVpcByName.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return nil, errors.New("HTTP Get list_vpc_by_name failed: " + err.Error())
	}
	var data GatewayDetailApiResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_vpc_by_name failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return nil, errors.New("Rest API list_vpc_by_name Get failed: " + data.Reason)
	}
	if data.Results.GwName == gateway.GwName {
		return &data.Results, nil
	}

	log.Errorf("Couldn't find Aviatrix gateway %s", gateway.GwName)
	return nil, ErrNotFound
}

func (c *Client) UpdateGateway(gateway *Gateway) error {
	gateway.CID = c.CID
	gateway.Action = "edit_gw_config"
	resp, err := c.Post(c.baseURL, gateway)
	if err != nil {
		return errors.New("HTTP Post edit_gw_config failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode edit_gw_config failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API edit_gw_config Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DeleteGateway(gateway *Gateway) error {
	path := c.baseURL + fmt.Sprintf("?action=delete_container&CID=%s&cloud_type=%d&gw_name=%s",
		c.CID, gateway.CloudType, gateway.GwName)
	resp, err := c.Delete(path, nil)
	if err != nil {
		return errors.New("HTTP Get delete_container failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode delete_container failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API delete_container Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) EnableSNat(gateway *Gateway) error {
	gateway.CID = c.CID
	gateway.Action = "enable_snat"
	args, err := json.Marshal(gateway.SnatPolicy)
	if err != nil {
		return err
	}
	gateway.PolicyList = string(args)
	resp, err := c.Post(c.baseURL, gateway)
	if err != nil {
		return errors.New("HTTP Get enable_snat failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode enable_snat failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API enable_snat Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DisableSNat(gateway *Gateway) error {
	gateway.CID = c.CID
	gateway.Action = "disable_snat"
	resp, err := c.Post(c.baseURL, gateway)
	if err != nil {
		return errors.New("HTTP Get 'disable_snat' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'disable_snat' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'disable_snat' Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DisableCustomSNat(gateway *Gateway) error {
	gateway.CID = c.CID
	gateway.Action = "enable_snat"
	resp, err := c.Post(c.baseURL, gateway)
	if err != nil {
		return errors.New("HTTP Get 'enable_snat' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'enable_snat' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'enable_snat' Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) UpdateDNat(gateway *Gateway) error {
	gateway.CID = c.CID
	gateway.Action = "update_dnat_config"
	args, err := json.Marshal(gateway.DnatPolicy)
	if err != nil {
		return err
	}
	gateway.PolicyList = string(args)
	resp, err := c.Post(c.baseURL, gateway)
	if err != nil {
		return errors.New("HTTP Get update_dnat_config failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode update_dnat_config failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API update_dnat_config Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) UpdateVpnCidr(gateway *Gateway) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for edit_vpn_gateway_virtual_address_range") + err.Error())
	}
	setVpnClientCIDR := url.Values{}
	setVpnClientCIDR.Add("CID", c.CID)
	setVpnClientCIDR.Add("action", "edit_vpn_gateway_virtual_address_range")
	setVpnClientCIDR.Add("vpn_cidr", gateway.VpnCidr)
	setVpnClientCIDR.Add("gateway_name", gateway.GwName)
	Url.RawQuery = setVpnClientCIDR.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get edit_vpn_gateway_virtual_address_range failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode edit_vpn_gateway_virtual_address_range failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API edit_vpn_gateway_virtual_address_range Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) UpdateMaxVpnConn(gateway *Gateway) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for set_vpn_client_cidr") + err.Error())
	}
	setMaxVpnConn := url.Values{}
	setMaxVpnConn.Add("CID", c.CID)
	setMaxVpnConn.Add("action", "set_vpn_max_connection")
	setMaxVpnConn.Add("max_connections", gateway.MaxConn)
	setMaxVpnConn.Add("vpc_id", gateway.VpcID)
	setMaxVpnConn.Add("lb_or_gateway_name", gateway.ElbName)
	if gateway.Dns == "true" {
		setMaxVpnConn.Add("dns", "true")
	}
	Url.RawQuery = setMaxVpnConn.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get set_vpn_max_connection failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode set_vpn_max_connection failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API set_vpn_max_connection Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) SetVpnGatewayAuthentication(gateway *VpnGatewayAuth) error {
	gateway.CID = c.CID
	gateway.Action = "set_vpn_gateway_authentication"
	resp, err := c.Post(c.baseURL, gateway)

	if err != nil {
		return errors.New("HTTP Post set_vpn_gateway_authentication failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode set_vpn_gateway_authentication failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API set_vpn_gateway_authentication Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) EnableActiveMesh(gateway *Gateway) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for enable_gateway_activemesh") + err.Error())
	}
	enableSNat := url.Values{}
	enableSNat.Add("CID", c.CID)
	enableSNat.Add("action", "enable_gateway_activemesh")
	enableSNat.Add("gateway_name", gateway.GwName)
	Url.RawQuery = enableSNat.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get enable_gateway_activemesh failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode enable_gateway_activemesh failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API enable_gateway_activemesh Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DisableActiveMesh(gateway *Gateway) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for disable_gateway_activemesh") + err.Error())
	}
	enableSNat := url.Values{}
	enableSNat.Add("CID", c.CID)
	enableSNat.Add("action", "disable_gateway_activemesh")
	enableSNat.Add("gateway_name", gateway.GwName)
	Url.RawQuery = enableSNat.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get disable_gateway_activemesh failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode disable_gateway_activemesh failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API disable_gateway_activemesh Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) EnableVpcDnsServer(gateway *Gateway) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for enable_vpc_dns_server") + err.Error())
	}
	enableSNat := url.Values{}
	enableSNat.Add("CID", c.CID)
	enableSNat.Add("action", "enable_vpc_dns_server")
	enableSNat.Add("gateway_name", gateway.GwName)
	Url.RawQuery = enableSNat.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get enable_vpc_dns_server failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode enable_vpc_dns_server failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API enable_vpc_dns_server Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DisableVpcDnsServer(gateway *Gateway) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for disable_vpc_dns_server") + err.Error())
	}
	disableSNat := url.Values{}
	disableSNat.Add("CID", c.CID)
	disableSNat.Add("action", "disable_vpc_dns_server")
	disableSNat.Add("gateway_name", gateway.GwName)
	Url.RawQuery = disableSNat.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get disable_vpc_dns_server failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode disable_vpc_dns_server failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API disable_vpc_dns_server Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) EnableVpnNat(gateway *Gateway) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for enable_nat_on_vpn_gateway") + err.Error())
	}
	enableVpnNat := url.Values{}
	enableVpnNat.Add("CID", c.CID)
	enableVpnNat.Add("action", "enable_nat_on_vpn_gateway")
	enableVpnNat.Add("vpc_id", gateway.VpcID)
	if gateway.ElbName != "" {
		enableVpnNat.Add("lb_or_gateway_name", gateway.ElbName)
	} else {
		enableVpnNat.Add("lb_or_gateway_name", gateway.GwName)
	}
	if gateway.Dns == "true" {
		enableVpnNat.Add("dns", "true")
	}
	Url.RawQuery = enableVpnNat.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get enable_nat_on_vpn_gateway failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode enable_nat_on_vpn_gateway failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API enable_nat_on_vpn_gateway Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DisableVpnNat(gateway *Gateway) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for disable_nat_on_vpn_gateway") + err.Error())
	}
	disableVpnNat := url.Values{}
	disableVpnNat.Add("CID", c.CID)
	disableVpnNat.Add("action", "disable_nat_on_vpn_gateway")
	disableVpnNat.Add("vpc_id", gateway.VpcID)
	if gateway.ElbName != "" {
		disableVpnNat.Add("lb_or_gateway_name", gateway.ElbName)
	} else {
		disableVpnNat.Add("lb_or_gateway_name", gateway.GwName)
	}
	if gateway.Dns == "true" {
		disableVpnNat.Add("dns", "true")
	}
	Url.RawQuery = disableVpnNat.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get disable_nat_on_vpn_gateway failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode disable_nat_on_vpn_gateway failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API disable_nat_on_vpn_gateway Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) EditDesignatedGateway(gateway *Gateway) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'EditDesignatedGateway': ") + err.Error())
	}
	editDesignatedGateway := url.Values{}
	editDesignatedGateway.Add("CID", c.CID)
	editDesignatedGateway.Add("action", "set_designated_gateway_additional_cidr_list")
	editDesignatedGateway.Add("gateway_name", gateway.GwName)
	editDesignatedGateway.Add("additional_cidr_list", gateway.AdditionalCidrsDesignatedGw)
	Url.RawQuery = editDesignatedGateway.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get 'set_designated_gateway_additional_cidr_list' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'set_designated_gateway_additional_cidr_list' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'set_designated_gateway_additional_cidr_list' Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) EnableEncryptVolume(gateway *Gateway) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'EnableEncryptVolume': ") + err.Error())
	}
	encryptGatewayVolume := url.Values{}
	encryptGatewayVolume.Add("CID", c.CID)
	encryptGatewayVolume.Add("action", "encrypt_gateway_volume")
	encryptGatewayVolume.Add("gateway_name", gateway.GwName)
	if gateway.CustomerManagedKeys != "" {
		encryptGatewayVolume.Add("customer_managed_keys", gateway.CustomerManagedKeys)
	}
	Url.RawQuery = encryptGatewayVolume.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get 'encrypt_gateway_volume' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'encrypt_gateway_volume' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		if strings.Contains(data.Reason, "already encrypted") {
			return nil
		}
		return errors.New("Rest API 'encrypt_gateway_volume' Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) EditGatewayCustomRoutes(gateway *Gateway) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'EditCustomRoutes': ") + err.Error())
	}
	editGatewayCustomRoutes := url.Values{}
	editGatewayCustomRoutes.Add("CID", c.CID)
	editGatewayCustomRoutes.Add("action", "edit_gateway_custom_routes")
	editGatewayCustomRoutes.Add("gateway_name", gateway.GwName)
	editGatewayCustomRoutes.Add("cidr", strings.Join(gateway.CustomizedSpokeVpcRoutes, ","))
	Url.RawQuery = editGatewayCustomRoutes.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get 'edit_gateway_custom_routes' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'edit_gateway_custom_routes' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'edit_gateway_custom_routes' Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) EditGatewayFilterRoutes(gateway *Gateway) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'EditFilterRoutes': ") + err.Error())
	}
	editGatewayFilterRoutes := url.Values{}
	editGatewayFilterRoutes.Add("CID", c.CID)
	editGatewayFilterRoutes.Add("action", "edit_gateway_filter_routes")
	editGatewayFilterRoutes.Add("gateway_name", gateway.GwName)
	editGatewayFilterRoutes.Add("cidr", strings.Join(gateway.FilteredSpokeVpcRoutes, ","))
	Url.RawQuery = editGatewayFilterRoutes.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get 'edit_gateway_filter_routes' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'edit_gateway_filter_routes' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'edit_gateway_filter_routes' Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) EditGatewayAdvertisedCidr(gateway *Gateway) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'EditGatewayAdveretisedCidr': ") + err.Error())
	}
	editGatewayFilterRoutes := url.Values{}
	editGatewayFilterRoutes.Add("CID", c.CID)
	editGatewayFilterRoutes.Add("action", "edit_gateway_advertised_cidr")
	editGatewayFilterRoutes.Add("gateway_name", gateway.GwName)
	editGatewayFilterRoutes.Add("cidr", strings.Join(gateway.AdvertisedSpokeRoutes, ","))
	Url.RawQuery = editGatewayFilterRoutes.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get 'edit_gateway_advertised_cidr' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'edit_gateway_advertised_cidr' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'edit_gateway_advertised_cidr' Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) EnableTransitFireNet(gateway *Gateway) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'EnableTransitFireNet': ") + err.Error())
	}
	enableTransitFireNet := url.Values{}
	enableTransitFireNet.Add("CID", c.CID)
	enableTransitFireNet.Add("action", "enable_gateway_for_transit_firenet")
	enableTransitFireNet.Add("gateway_name", gateway.GwName)
	Url.RawQuery = enableTransitFireNet.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get 'enable_gateway_for_transit_firenet' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'enable_gateway_for_transit_firenet' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'enable_gateway_for_transit_firenet' Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DisableTransitFireNet(gateway *Gateway) error {
	err := c.IsTransitFireNetReadyToBeDisabled(gateway)
	if err != nil {
		return err
	}

	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for 'DisableTransitFireNet': ") + err.Error())
	}
	enableTransitFireNet := url.Values{}
	enableTransitFireNet.Add("CID", c.CID)
	enableTransitFireNet.Add("action", "disable_gateway_for_transit_firenet")
	enableTransitFireNet.Add("gateway_name", gateway.GwName)
	Url.RawQuery = enableTransitFireNet.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get 'disable_gateway_for_transit_firenet' failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode 'disable_gateway_for_transit_firenet' failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API 'disable_gateway_for_transit_firenet' Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) IsTransitFireNetReadyToBeDisabled(gateway *Gateway) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for IsTransitFireNetReadyToBeDisabled: ") + err.Error())
	}
	listTransitFireNetSpokePolicies := url.Values{}
	listTransitFireNetSpokePolicies.Add("CID", c.CID)
	listTransitFireNetSpokePolicies.Add("action", "list_transit_firenet_spoke_policies")
	Url.RawQuery = listTransitFireNetSpokePolicies.Encode()
	resp, err := c.Get(Url.String(), nil)

	if err != nil {
		return errors.New("HTTP Get list_transit_firenet_spoke_policies failed: " + err.Error())
	}
	var data ListTransitFireNetPolicyResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode list_transit_firenet_spoke_policies failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API list_transit_firenet_spoke_policies Get failed: " + data.Reason)
	}
	if len(data.Results) == 0 {
		return nil
	}
	policyList := data.Results
	for i := range policyList {
		if policyList[i].TransitFireNetGwName != gateway.GwName {
			continue
		}
		if policyList[i].ManagementAccessResourceName != "no" && len(policyList[i].InspectedResourceNameList) != 0 {
			return fmt.Errorf("%s is still firewall management access enabled and has transit firenet policy/policies", gateway.GwName)
		} else if policyList[i].ManagementAccessResourceName != "no" {
			return fmt.Errorf("%s is still firewall management access enabled", gateway.GwName)
		} else if len(policyList[i].InspectedResourceNameList) != 0 {
			return fmt.Errorf("%s still has transit firenet policy/policies", gateway.GwName)
		}
	}
	return nil
}

func (c *Client) EnableSegmentation(transitGateway *TransitVpc) error {
	action := "enable_transit_gateway_for_multi_cloud_security_domain"
	form := map[string]interface{}{
		"CID":                  c.CID,
		"action":               action,
		"transit_gateway_name": transitGateway.GwName,
	}
	return c.PostAPI(action, form, BasicCheck)
}

func (c *Client) DisableSegmentation(transitGateway *TransitVpc) error {
	action := "disable_transit_gateway_for_multi_cloud_security_domain"
	form := map[string]interface{}{
		"CID":                  c.CID,
		"action":               action,
		"transit_gateway_name": transitGateway.GwName,
	}
	return c.PostAPI(action, form, BasicCheck)
}

func (c *Client) IsSegmentationEnabled(transitGateway *TransitVpc) (bool, error) {
	action := "list_transit_gateways_for_multi_cloud_domains"
	resp, err := c.Post(c.baseURL, &APIRequest{
		CID:    c.CID,
		Action: action,
	})
	if err != nil {
		return false, fmt.Errorf("HTTP POST %q failed: %v", action, err)
	}

	type Result struct {
		EnabledDomains  []string `json:"domain_enabled_list"`
		DisabledDomains []string `json:"domain_disabled_list"`
	}

	type Resp struct {
		Return  bool   `json:"return"`
		Results Result `json:"results"`
		Reason  string `json:"reason"`
	}

	var data Resp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return false, fmt.Errorf("reading response body %q failed: %v", action, err)
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return false, fmt.Errorf("json decode %q failed: %v\nBody: %s", action, err, b.String())
	}
	if !data.Return {
		return false, fmt.Errorf("rest API %q Post failed: %s", action, data.Reason)
	}

	return Contains(data.Results.EnabledDomains, transitGateway.GwName), nil
}

func (c *Client) EnableEgressTransitFirenet(transitGateway *TransitVpc) error {
	action := "enable_transit_firenet_on_egress_transit_gateway"
	data := map[string]interface{}{
		"CID":          c.CID,
		"action":       action,
		"gateway_name": transitGateway.GwName,
	}
	return c.PostAPI(action, data, BasicCheck)
}

func (c *Client) DisableEgressTransitFirenet(transitGateway *TransitVpc) error {
	action := "disable_transit_firenet_on_egress_transit_gateway"
	data := map[string]interface{}{
		"CID":          c.CID,
		"action":       action,
		"gateway_name": transitGateway.GwName,
	}
	return c.PostAPI(action, data, BasicCheck)
}

func (c *Client) EnableMonitorGatewaySubnets(gateway *Gateway) error {
	action := "enable_monitor_gateway_subnets"
	MonitorInstanceList := gateway.MonitorExcludeList
	form := map[string]interface{}{
		"CID":          c.CID,
		"action":       action,
		"gateway_name": gateway.GwName,
	}
	if len(MonitorInstanceList) != 0 {
		form["monitor_exclude_gateway_list"] = MonitorInstanceList
	}
	return c.PostAPI(action, form, BasicCheck)
}

func (c *Client) DisableMonitorGatewaySubnets(gateway *Gateway) error {
	action := "disable_monitor_gateway_subnets"
	form := map[string]interface{}{
		"CID":          c.CID,
		"action":       action,
		"gateway_name": gateway.GwName,
	}
	return c.PostAPI(action, form, BasicCheck)
}
