package goaviatrix

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Gateway simple struct to hold gateway details
type Gateway struct {
	AccountName             string `form:"account_name,omitempty" json:"account_name,omitempty"`
	Action                  string `form:"action,omitempty"`
	AdditionalCidrs         string `form:"additional_cidrs,omitempty" json:"additional_cidrs"`
	AuthMethod              string `form:"auth_method,omitempty" json:"auth_method,omitempty"`
	AllocateNewEip          string `form:"allocate_new_eip,omitempty"`
	AllocateNewEipReadPtr   *bool  `json:"newly_allocated_eip,omitempty"`
	AllocateNewEipRead      bool
	BkupGatewayZone         string `form:"bkup_gateway_zone,omitempty" json:"bkup_gateway_zone,omitempty"`
	BkupPrivateIP           string `form:"bkup_private_ip,omitempty" json:"bkup_private_ip,omitempty"`
	CID                     string `form:"CID,omitempty"`
	CIDR                    string `form:"cidr,omitempty"`
	ClientCertAuth          string `form:"client_cert_auth,omitempty" json:"client_cert_auth,omitempty"`
	ClientCertSharing       string `form:"client_cert_sharing,omitempty" json:"client_cert_sharing,omitempty"`
	CloudType               int    `form:"cloud_type,omitempty" json:"cloud_type,omitempty"`
	CloudnBkupGatewayInstID string `form:"cloudn_bkup_gateway_inst_id,omitempty" json:"cloudn_bkup_gateway_inst_id,omitempty"`
	CloudnGatewayInstID     string `form:"cloudn_gateway_inst_id,omitempty" json:"cloudn_gateway_inst_id,omitempty"`
	ConnectedTransit        string `json:"connected_transit,omitempty"`
	DirectInternet          string `form:"direct_internet,omitempty" json:"direct_internet,omitempty"`
	DockerConsulIP          string `form:"docker_consul_ip,omitempty" json:"docker_consul_ip,omitempty"`
	DockerNtwkCidr          string `form:"docker_ntwk_cidr,omitempty" json:"docker_ntwk_cidr,omitempty"`
	DockerNtwkName          string `form:"docker_ntwk_name,omitempty" json:"docker_ntwk_name,omitempty"`
	DuoAPIHostname          string `form:"duo_api_hostname,omitempty" json:"duo_api_hostname,omitempty"`
	DuoIntegrationKey       string `form:"duo_integration_key,omitempty" json:"duo_integration_key,omitempty"`
	DuoPushMode             string `form:"duo_push_mode,omitempty" json:"duo_push_mode,omitempty"`
	DuoSecretKey            string `form:"duo_secret_key,omitempty" json:"duo_secret_key,omitempty"`
	Eip                     string `form:"eip,omitempty" json:"eip,omitempty"`
	ReuseEip                string `json:"reuse_eip,omitempty"`
	ElbDNSName              string `form:"elb_dns_name,omitempty" json:"elb_dns_name,omitempty"`
	ElbName                 string `form:"elb_name,omitempty" json:"lb_name,omitempty"`
	ElbState                string `form:"elb_state,omitempty" json:"elb_state,omitempty"`
	VpnProtocol             string `json:"vpn_protocol" form:"elb_protocol,omitempty"`
	EnableClientCertSharing string `form:"enable_client_cert_sharing,omitempty"`
	EnableElb               string `form:"enable_elb,omitempty"`
	EnableLdap              string `form:"enable_ldap,omitempty"`
	EnableLdapRead          bool   `json:"enable_ldap,omitempty"`
	EnableVpcDnsServer      string `json:"use_vpc_dns,omitempty"`
	DnsServer               string `form:"dns_server,omitempty"`
	PublicDnsServer         string `form:"public_dns_server,omitempty" json:"public_dns_server,omitempty"`

	// These two are very similar but have a slight difference
	// EnableNat - will be "yes" if single/multiple SNAT is enabled
	// NatEnabled - will be true if single/multiple/customized SNAT is enabled
	EnableNat  string `form:"enable_nat,omitempty" json:"enable_nat,omitempty"`
	NatEnabled bool   `json:"nat_enabled,omitempty"`

	SingleAZ                        string            `form:"single_az_ha,omitempty" json:"single_az_ha,omitempty"`
	EnableHybridConnection          bool              `json:"tgw_enabled,omitempty"`
	EnablePbr                       string            `form:"enable_pbr,omitempty"`
	Expiration                      string            `form:"expiration,omitempty" json:"expiration,omitempty"`
	GatewayZone                     string            `form:"gateway_zone,omitempty" json:"gateway_zone,omitempty"`
	GwName                          string            `form:"gw_name,omitempty" json:"vpc_name,omitempty"`
	GwSecurityGroupID               string            `form:"gw_security_group_id,omitempty" json:"gw_security_group_id,omitempty"`
	GwSize                          string            `form:"gw_size,omitempty" json:"vpc_size,omitempty"`
	GwSubnetID                      string            `form:"gw_subnet_id,omitempty" json:"gw_subnet_id,omitempty"`
	PeeringHASubnet                 string            `form:"public_subnet,omitempty"`
	NewZone                         string            `form:"new_zone,omitempty"`
	NewSubnet                       string            `form:"new_subnet,omitempty"`
	InsaneMode                      string            `form:"insane_mode,omitempty" json:"high_perf,omitempty"`
	InstState                       string            `form:"inst_state,omitempty" json:"inst_state,omitempty"`
	IntraVMRoute                    string            `form:"intra_vm_route,omitempty" json:"intra_vm_route,omitempty"`
	IsHagw                          string            `form:"is_hagw,omitempty" json:"is_hagw,omitempty"`
	JumboFrame                      bool              `json:"jumbo_frame"`
	LdapAdditionalReq               string            `form:"ldap_additional_req,omitempty"`
	LdapBaseDn                      string            `form:"ldap_base_dn,omitempty" json:"ldap_base_dn,omitempty"`
	LdapBindDn                      string            `form:"ldap_bind_dn,omitempty" json:"ldap_bind_dn,omitempty"`
	LdapCaCert                      string            `form:"ldap_ca_cert,omitempty" json:"ldap_ca_cert,omitempty"`
	LdapClientCert                  string            `form:"ldap_client_cert,omitempty" json:"ldap_client_cert,omitempty"`
	LdapPassword                    string            `form:"ldap_password,omitempty" json:"ldap_password,omitempty"`
	LdapServer                      string            `form:"ldap_server,omitempty" json:"ldap_server,omitempty"`
	LdapUseSsl                      string            `form:"ldap_use_ssl,omitempty" json:"ldap_use_ssl,omitempty"`
	LdapUserAttr                    string            `form:"ldap_username_attribute,omitempty" json:"ldap_username_attribute,omitempty"`
	LicenseID                       string            `form:"license_id,omitempty" json:"license_id,omitempty"`
	MaxConn                         string            `form:"max_conn,omitempty" json:"max_connections,omitempty"`
	NameServers                     string            `form:"nameservers,omitempty" json:"name_servers"`
	OktaToken                       string            `form:"okta_token,omitempty" json:"okta_token,omitempty"`
	OktaURL                         string            `form:"okta_url,omitempty" json:"okta_url,omitempty"`
	OktaUsernameSuffix              string            `form:"okta_username_suffix,omitempty" json:"okta_username_suffix,omitempty"`
	OtpMode                         string            `form:"otp_mode,omitempty" json:"otp_mode,omitempty"`
	PbrDefaultGateway               string            `form:"pbr_default_gateway,omitempty"`
	PbrEnabled                      string            `form:"pbr_enabled,omitempty" json:"pbr_enabled,omitempty"`
	PbrLogging                      string            `form:"pbr_logging,omitempty"`
	PbrSubnet                       string            `form:"pbr_subnet,omitempty"`
	PrivateIP                       string            `form:"private_ip,omitempty" json:"private_ip,omitempty"`
	PublicIP                        string            `form:"public_ip,omitempty" json:"public_ip,omitempty"`
	SamlEnabled                     string            `form:"saml_enabled,omitempty" json:"saml_enabled,omitempty"`
	SandboxIP                       string            `form:"sandbox_ip,omitempty" json:"sandbox_ip,omitempty"`
	SaveTemplate                    string            `form:"save_template,omitempty"`
	SearchDomains                   string            `form:"search_domains,omitempty" json:"search_domains"`
	SplitTunnel                     string            `form:"split_tunnel,omitempty" json:"split_tunnel,omitempty"`
	SpokeVpc                        string            `json:"spoke_vpc,omitempty"`
	TagList                         string            `form:"tag_string,omitempty"`
	Tags                            map[string]string `json:"tags,omitempty"`
	TagJson                         string            `form:"tag_json,omitempty"`
	TransitGwName                   string            `form:"transit_gw_name,omitempty" json:"transit_gw_name,omitempty"`
	EgressTransitGwName             string            `form:"egress_transit_gw_name,omitempty" json:"egress_transit_gw_name,omitempty"`
	TunnelName                      string            `form:"tunnel_name,omitempty" json:"tunnel_name,omitempty"`
	TunnelType                      string            `form:"tunnel_type,omitempty" json:"tunnel_type,omitempty"`
	VendorName                      string            `form:"vendor_name,omitempty" json:"vendor_name,omitempty"`
	VpcID                           string            `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	VpcNet                          string            `form:"vpc_net,omitempty" json:"public_subnet,omitempty"`
	VpcRegion                       string            `form:"vpc_reg,omitempty" json:"vpc_region,omitempty"`
	VpcSplunkIPPort                 string            `form:"vpc_splunk_ip_port,omitempty" json:"vpc_splunk_ip_port,omitempty"`
	VpcState                        string            `form:"vpc_state,omitempty" json:"vpc_state,omitempty"`
	VpcType                         string            `form:"vpc_type,omitempty" json:"vpc_type,omitempty"`
	VpnCidr                         string            `form:"cidr,omitempty" json:"vpn_cidr,omitempty"`
	VpnStatus                       string            `form:"vpn_access,omitempty" json:"vpn_status,omitempty"`
	Zone                            string            `form:"zone,omitempty" json:"zone,omitempty"`
	VpcSize                         string            `form:"vpc_size,omitempty" ` //Only use for gateway create
	DMZEnabled                      string            `json:"dmz_enabled,omitempty"`
	EnableActiveMesh                string            `form:"enable_activemesh,omitempty" json:"enable_activemesh,omitempty"`
	EnableVpnNat                    bool              `form:"vpn_nat,omitempty" json:"vpn_nat"`
	EnableDesignatedGateway         string            `form:"designated_gateway,omitempty" json:"designated_gateway,omitempty"`
	AdditionalCidrsDesignatedGw     string            `form:"additional_cidr_list,omitempty" json:"summarized_cidrs,omitempty"`
	EnableEncryptVolume             bool              `json:"gw_enc,omitempty"`
	CustomerManagedKeys             string            `form:"customer_managed_keys,omitempty" json:"customer_managed_keys,omitempty"`
	SnatMode                        string            `form:"mode,omitempty" json:"snat_target,omitempty"`
	SnatPolicy                      []PolicyRule
	PolicyList                      string `form:"policy_list,omitempty"`
	GatewayName                     string `form:"gateway_name,omitempty"`
	DnatPolicy                      []PolicyRule
	CustomizedSpokeVpcRoutes        []string `json:"customized_cidr_list,omitempty"`
	FilteredSpokeVpcRoutes          []string `json:"filtering_cidr_list,omitempty"`
	AdvertisedSpokeRoutes           []string `json:"advertise_cidr_list,omitempty"`
	IncludeCidrList                 []string `json:"include_cidr_list,omitempty"`
	ExcludeCidrList                 []string `json:"exclude_cidr_list,omitempty"`
	LearnedCidrsApproval            string   `json:"learned_cidrs_approval,omitempty"`
	Dns                             string   `json:"dns,omitempty"`
	EncVolume                       string   `form:"enc_volume,omitempty"`
	SyncSNATToHA                    string   `form:"sync_snat_to_ha,omitempty"`
	SyncDNATToHA                    string   `form:"sync_dnat_to_ha,omitempty"`
	MonitorSubnetsAction            string   `form:"monitor_subnets_action,omitempty" json:"monitor_subnets_action,omitempty"`
	MonitorExcludeGWList            []string `form:"monitor_exclude_gw_list,omitempty" json:"monitor_exclude_gw_list,omitempty"`
	FqdnLanCidr                     string   `form:"fqdn_lan_cidr,omitempty"`
	RouteTable                      string
	EnablePrivateOob                bool                                `json:"private_oob"`
	OobManagementSubnet             string                              `json:"oob_mgmt_subnet"`
	LanVpcID                        string                              `form:"lan_vpc_id"`
	LanPrivateSubnet                string                              `form:"lan_private_subnet"`
	CreateFQDNGateway               bool                                `form:"create_firewall_gw"`
	PrivateVpcDefaultEnabled        bool                                `json:"private_vpc_default_enabled"`
	SkipPublicVpcUpdateEnabled      bool                                `json:"skip_public_vpc_update_enabled"`
	EnableMultitierTransit          bool                                `json:"multitier_transit"`
	AutoAdvertiseCidrsEnabled       bool                                `json:"auto_advertise_s2c_cidrs"`
	TunnelDetectionTime             int                                 `json:"detection_time"`
	StorageName                     string                              `form:"storage_name,omitempty" json:"storage_name,omitempty"`
	BgpHoldTime                     int                                 `json:"bgp_hold_time"`
	BgpPollingTime                  int                                 `json:"bgp_polling_time"`
	PrependASPath                   string                              `json:"prepend_as_path"`
	LocalASNumber                   string                              `json:"local_as_number"`
	BgpEcmp                         bool                                `json:"bgp_ecmp"`
	EnableActiveStandby             bool                                `json:"enable_active_standby"`
	EnableBgpOverLan                bool                                `json:"enable_bgp_over_lan"`
	EnableTransitSummarizeCidrToTgw bool                                `json:"enable_transit_summarize_cidr_to_tgw"`
	EnableSegmentation              bool                                `json:"enable_segmentation"`
	LearnedCidrsApprovalMode        string                              `json:"learned_cidrs_approval_mode"`
	EnableFirenet                   bool                                `json:"enable_firenet"`
	EnableTransitFirenet            bool                                `json:"enable_transit_firenet"`
	EnableGatewayLoadBalancer       bool                                `json:"enable_gateway_load_balancer"`
	EnableEgressTransitFirenet      bool                                `json:"enable_egress_transit_firenet"`
	CustomizedTransitVpcRoutes      []string                            `json:"customized_transit_vpc_routes"`
	EnableAdvertiseTransitCidr      bool                                `json:"enable_advertise_transit_cidr"`
	EnableLearnedCidrsApproval      bool                                `json:"enable_learned_cidrs_approval"`
	BgpManualSpokeAdvertiseCidrs    []string                            `json:"bgp_manual_spoke_advertise_cidrs"`
	IdleTimeout                     string                              `json:"idle_timeout"`
	RenegotiationInterval           string                              `json:"renegotiation_interval"`
	FqdnInterfaces                  map[string][]string                 `json:"fqdn_interfaces"`
	ArmFqdnLanCidr                  map[string]string                   `json:"fqdn_fqdn_lan_cidr"`
	IsPsfGateway                    bool                                `json:"is_psf_gw"`
	PsfDetails                      PublicSubnetFilteringGatewayDetails `json:"psf_details"`
	BundleVpcInfo                   BundleVpcInfo                       `json:"bundle_vpc_info"`
	HaGw                            HaGateway                           `json:"hagw_details"`
	AvailabilityDomain              string                              `form:"availability_domain,omitempty"`
	FaultDomain                     string                              `form:"fault_domain,omitempty" json:"fault_domain"`
	EnableSpotInstance              bool                                `form:"spot_instance,omitempty" json:"spot_instance"`
	SpotPrice                       string                              `form:"spot_price,omitempty" json:"spot_price"`
	ImageVersion                    string                              `json:"gw_image_name"`
	SoftwareVersion                 string                              `json:"gw_software_version"`
}

type HaGateway struct {
	GwName              string `json:"vpc_name"`
	CloudType           int    `json:"cloud_type"`
	GwSize              string `json:"vpc_size"`
	VpcNet              string `json:"public_subnet"`
	PublicIP            string `json:"public_ip"`
	PrivateIP           string `json:"private_ip"`
	ReuseEip            string `json:"reuse_eip,omitempty"`
	CloudnGatewayInstID string `json:"cloudn_gateway_inst_id"`
	GatewayZone         string `json:"gateway_zone"`
	InsaneMode          string `json:"high_perf"`
	EnablePrivateOob    bool   `json:"private_oob"`
	OobManagementSubnet string `json:"oob_mgmt_subnet"`
	GwSecurityGroupID   string `json:"gw_security_group_id"`
	FaultDomain         string `json:"fault_domain"`
	ImageVersion        string `json:"gw_image_name"`
	SoftwareVersion     string `json:"gw_software_version"`
}

type PolicyRule struct {
	SrcIP           string `form:"src_ip,omitempty" json:"src_ip,omitempty"`
	SrcPort         string `form:"src_port,omitempty" json:"src_port,omitempty"`
	DstIP           string `form:"dst_ip,omitempty" json:"dst_ip,omitempty"`
	DstPort         string `form:"dst_port,omitempty" json:"dst_port,omitempty"`
	Protocol        string `form:"protocol,omitempty" json:"protocol,omitempty"`
	Interface       string `form:"interface,omitempty" json:"interface,omitempty"`
	Connection      string `form:"connection,omitempty" json:"connection,omitempty"`
	Mark            string `form:"mark,omitempty" json:"mark,omitempty"`
	NewSrcIP        string `form:"new_src_ip,omitempty" json:"new_src_ip,omitempty"`
	NewSrcPort      string `form:"new_src_port,omitempty" json:"new_src_port,omitempty"`
	ExcludeRTB      string `form:"exclude_rtb,omitempty" json:"exclude_rtb,omitempty"`
	ApplyRouteEntry bool   `form:"apply_route_entry,omitempty" json:"apply_route_entry"`
	NewDstIP        string `form:"new_dst_ip,omitempty" json:"new_dst_ip,omitempty"`
	NewDstPort      string `form:"new_dst_port,omitempty" json:"new_dst_port,omitempty"`
}

type GatewayDetail struct {
	AccountName                  string        `form:"account_name,omitempty" json:"account_name,omitempty"`
	Action                       string        `form:"action,omitempty"`
	GwName                       string        `form:"gw_name,omitempty" json:"vpc_name,omitempty"`
	DMZEnabled                   bool          `json:"dmz_enabled,omitempty"`
	EnableAdvertiseTransitCidr   string        `json:"advertise_transit_cidr,omitempty"`
	BgpManualSpokeAdvertiseCidrs []string      `json:"bgp_manual_spoke_advertise_cidrs,omitempty"`
	VpnNat                       bool          `json:"vpn_nat,omitempty"`
	SnatPolicy                   []PolicyRule  `json:"snat_ip_port_list,omitempty"`
	DnatPolicy                   []PolicyRule  `json:"dnat_ip_port_list,omitempty"`
	Elb                          ElbDetail     `json:"elb,omitempty"`
	EnableEgressTransitFireNet   bool          `json:"egress_transit,omitempty"`
	EnableFireNet                bool          `json:"firenet_enabled,omitempty"`
	EnabledGatewayLoadBalancer   bool          `json:"gwlb_enabled,omitempty"`
	EnableTransitFireNet         bool          `json:"transit_firenet_enabled,omitempty"`
	LearnedCidrsApproval         string        `json:"learned_cidrs_approval,omitempty"`
	SyncSNATToHA                 bool          `json:"sync_snat_to_ha,omitempty"`
	SyncDNATToHA                 bool          `json:"sync_dnat_to_ha,omitempty"`
	GwZone                       string        `json:"gw_zone,omitempty"`
	TransitGwName                string        `json:"transit_gw_name,omitempty"`
	EgressTransitGwName          string        `json:"egress_transit_gw_name,omitempty"`
	RouteTables                  []string      `json:"spoke_rtb_list,omitempty"`
	CustomizedTransitVpcRoutes   []string      `json:"customized_transit_vpc_cidrs"`
	BundleVpcInfo                BundleVpcInfo `json:"bundle_vpc_info"`
}

type BundleVpcInfo struct {
	LAN BundleVpcLanInfo
}

type BundleVpcLanInfo struct {
	VpcID  string `json:"vpc_id"`
	Subnet string
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

type VPNConfigListResp struct {
	Return  bool        `json:"return"`
	Results []VPNConfig `json:"results"`
	Reason  string      `json:"reason"`
}

type VPNConfig struct {
	Name   string `form:"name,omitempty" json:"name,omitempty"`
	Value  string `form:"value,omitempty" json:"value,omitempty"`
	Status string `form:"status,omitempty" json:"status,omitempty"`
}

type FQDNGatewayInfoResp struct {
	Return  bool           `json:"return"`
	Results FQDNGatwayInfo `json:"results"`
	Reason  string         `json:"reason"`
}

type FQDNGatwayInfo struct {
	Instances      []string            `json:"instances"`
	Interface      map[string][]string `json:"interfaces"`
	ArmFqdnLanCidr map[string]string   `json:"arm_fqdn_lan_cidr"`
}

func (c *Client) CreateGateway(gateway *Gateway) error {
	gateway.CID = c.CID
	gateway.Action = "connect_container"

	return c.PostAPI(gateway.Action, gateway, BasicCheck)
}

func (c *Client) CreatePublicSubnetFilteringGateway(gateway *Gateway) error {
	data := map[string]string{
		"action":         "add_public_subnet_filtering_gateway",
		"CID":            c.CID,
		"cloud_type":     strconv.Itoa(gateway.CloudType),
		"account_name":   gateway.AccountName,
		"region":         gateway.VpcRegion,
		"vpc_id":         gateway.VpcID,
		"gateway_name":   gateway.GwName,
		"gateway_size":   gateway.VpcSize,
		"gateway_subnet": gateway.VpcNet,
		"route_table":    gateway.RouteTable,
		"tag":            "",
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) DeletePublicSubnetFilteringGateway(gateway *Gateway) error {
	data := map[string]string{
		"action":       "delete_public_subnet_filtering_gateway",
		"CID":          c.CID,
		"gateway_name": gateway.GwName,
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) EnablePublicSubnetFilteringHAGateway(gateway *Gateway) error {
	data := map[string]string{
		"action":         "enable_ha_for_public_subnet_filtering_gateway",
		"CID":            c.CID,
		"gateway_name":   gateway.GwName,
		"gateway_subnet": gateway.PeeringHASubnet,
		"route_tables":   gateway.RouteTable,
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

type PublicSubnetFilteringGatewayDetails struct {
	RouteTableList    []string `json:"rtb_list"`
	HaRouteTableList  []string `json:"ha_rtb_list"`
	GuardDutyEnforced string   `json:"guard_duty_enforced"`
	GwSubnetCidr      string   `json:"gw_subnet_cidr"`
	GwSubnetAz        string   `json:"gw_subnet_az"`
	HaGwSubnetCidr    string   `json:"ha_gw_subnet_cidr"`
	HaGwSubnetAz      string   `json:"ha_gw_subnet_az"`
}

type PublicSubnetFilteringGatewayDetailsResp struct {
	Return  bool                                `json:"return"`
	Results PublicSubnetFilteringGatewayDetails `json:"results"`
	Reason  string                              `json:"reason"`
}

func (c *Client) GetPublicSubnetFilteringGatewayDetails(gateway *Gateway) (*PublicSubnetFilteringGatewayDetails, error) {
	data := map[string]string{
		"action":       "get_public_subnet_filtering_gateway_details",
		"CID":          c.CID,
		"gateway_name": gateway.GwName,
	}
	var resp PublicSubnetFilteringGatewayDetailsResp
	err := c.GetAPI(&resp, data["action"], data, BasicCheck)
	if err != nil {
		return nil, err
	}
	return &resp.Results, nil
}

func (c *Client) EditPublicSubnetFilteringRouteTableList(gateway *Gateway, routeTables []string) error {
	data := map[string]string{
		"action":       "edit_public_subnet_filtering_enforced_route_table_list",
		"CID":          c.CID,
		"gateway_name": gateway.GwName,
		"route_table":  strings.Join(routeTables, ", "),
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) EnableGuardDutyEnforcement(gateway *Gateway) error {
	data := map[string]string{
		"action":       "enable_public_subnet_filtering_guard_duty_enforced_mode",
		"CID":          c.CID,
		"gateway_name": gateway.GwName,
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) DisableGuardDutyEnforcement(gateway *Gateway) error {
	data := map[string]string{
		"action":       "disable_public_subnet_filtering_guard_duty_enforced_mode",
		"CID":          c.CID,
		"gateway_name": gateway.GwName,
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) EnableNatGateway(gateway *Gateway) error {
	gateway.CID = c.CID
	gateway.Action = "enable_nat"

	return c.PostAPI(gateway.Action, gateway, BasicCheck)
}

func (c *Client) EnableSingleAZGateway(gateway *Gateway) error {
	gateway.CID = c.CID
	gateway.Action = "enable_single_az_ha"

	return c.PostAPI(gateway.Action, gateway, BasicCheck)
}

func (c *Client) EnablePeeringHaGateway(gateway *Gateway) error {
	gateway.CID = c.CID
	gateway.Action = "create_peering_ha_gateway"

	return c.PostAPI(gateway.Action, gateway, BasicCheck)
}

func (c *Client) DisableSingleAZGateway(gateway *Gateway) error {
	gateway.CID = c.CID
	gateway.Action = "disable_single_az_ha"

	return c.PostAPI(gateway.Action, gateway, BasicCheck)
}

func (c *Client) GetGateway(gateway *Gateway) (*Gateway, error) {
	action := "list_vpcs_summary"
	params := map[string]string{
		"CID":          c.CID,
		"action":       action,
		"gateway_name": gateway.GwName,
	}

	var data GatewayListResp
	err := c.GetAPI(&data, action, params, BasicCheck)
	if err != nil {
		return nil, err
	}

	gwList := data.Results
	for i := range gwList {
		if gwList[i].GwName == gateway.GwName {
			gw := &gwList[i]
			// AllocateNewEipRead should default to true when not set by backend
			gw.AllocateNewEipRead = gw.AllocateNewEipReadPtr == nil || *gw.AllocateNewEipReadPtr
			return &gwList[i], nil
		}
	}
	log.Errorf("Couldn't find Aviatrix gateway %s", gateway.GwName)
	return nil, ErrNotFound
}

func (c *Client) GetGatewayDetail(gateway *Gateway) (*GatewayDetail, error) {
	form := map[string]string{
		"CID":      c.CID,
		"action":   "list_vpc_by_name",
		"vpc_name": gateway.GwName,
	}

	var data GatewayDetailApiResp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
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

	return c.PostAPI(gateway.Action, gateway, BasicCheck)
}

func (c *Client) DeleteGateway(gateway *Gateway) error {
	form := map[string]string{
		"CID":        c.CID,
		"action":     "delete_container",
		"cloud_type": strconv.Itoa(gateway.CloudType),
		"gw_name":    gateway.GwName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) EnableSNat(gateway *Gateway) error {
	gateway.CID = c.CID
	gateway.Action = "enable_snat"
	args, err := json.Marshal(gateway.SnatPolicy)
	if err != nil {
		return err
	}
	gateway.PolicyList = string(args)

	return c.PostAPI(gateway.Action, gateway, BasicCheck)
}

func (c *Client) DisableSNat(gateway *Gateway) error {
	gateway.CID = c.CID
	gateway.Action = "disable_snat"

	return c.PostAPI(gateway.Action, gateway, BasicCheck)
}

func (c *Client) DisableCustomSNat(gateway *Gateway) error {
	gateway.CID = c.CID
	gateway.Action = "enable_snat"

	return c.PostAPI(gateway.Action, gateway, BasicCheck)
}

func (c *Client) UpdateDNat(gateway *Gateway) error {
	gateway.CID = c.CID
	gateway.Action = "update_dnat_config"
	args, err := json.Marshal(gateway.DnatPolicy)
	if err != nil {
		return err
	}
	gateway.PolicyList = string(args)

	return c.PostAPI(gateway.Action, gateway, BasicCheck)
}

func (c *Client) UpdateVpnCidr(gateway *Gateway) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "edit_vpn_gateway_virtual_address_range",
		"vpn_cidr":     gateway.VpnCidr,
		"gateway_name": gateway.GwName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) UpdateMaxVpnConn(gateway *Gateway) error {
	form := map[string]string{
		"CID":                c.CID,
		"action":             "set_vpn_max_connection",
		"max_connections":    gateway.MaxConn,
		"vpc_id":             gateway.VpcID,
		"lb_or_gateway_name": gateway.ElbName,
	}

	if gateway.Dns == "true" {
		form["dns"] = "true"
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) SetVpnGatewayAuthentication(gateway *VpnGatewayAuth) error {
	gateway.CID = c.CID
	gateway.Action = "set_vpn_gateway_authentication"

	return c.PostAPI(gateway.Action, gateway, BasicCheck)
}

func (c *Client) EnableActiveMesh(gateway *Gateway) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "enable_gateway_activemesh",
		"gateway_name": gateway.GwName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) DisableActiveMesh(gateway *Gateway) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "disable_gateway_activemesh",
		"gateway_name": gateway.GwName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) EnableVpcDnsServer(gateway *Gateway) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "enable_vpc_dns_server",
		"gateway_name": gateway.GwName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) DisableVpcDnsServer(gateway *Gateway) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "disable_vpc_dns_server",
		"gateway_name": gateway.GwName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) EnableVpnNat(gateway *Gateway) error {
	form := map[string]string{
		"CID":    c.CID,
		"action": "enable_nat_on_vpn_gateway",
		"vpc_id": gateway.VpcID,
	}

	if gateway.ElbName != "" {
		form["lb_or_gateway_name"] = gateway.ElbName
	} else {
		form["lb_or_gateway_name"] = gateway.GwName
	}
	if gateway.Dns == "true" {
		form["dns"] = "true"
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) DisableVpnNat(gateway *Gateway) error {
	form := map[string]string{
		"CID":    c.CID,
		"action": "disable_nat_on_vpn_gateway",
		"vpc_id": gateway.VpcID,
	}

	if gateway.ElbName != "" {
		form["lb_or_gateway_name"] = gateway.ElbName
	} else {
		form["lb_or_gateway_name"] = gateway.GwName
	}
	if gateway.Dns == "true" {
		form["dns"] = "true"
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) EditDesignatedGateway(gateway *Gateway) error {
	form := map[string]string{
		"CID":                  c.CID,
		"action":               "set_designated_gateway_additional_cidr_list",
		"gateway_name":         gateway.GwName,
		"additional_cidr_list": gateway.AdditionalCidrsDesignatedGw,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) EnableEncryptVolume(gateway *Gateway) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "encrypt_gateway_volume",
		"gateway_name": gateway.GwName,
	}

	if gateway.CustomerManagedKeys != "" {
		form["customer_managed_keys"] = gateway.CustomerManagedKeys
	}

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "already encrypted") {
				return nil
			}
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	return c.PostAPI(form["action"], form, checkFunc)
}

func (c *Client) EditGatewayCustomRoutes(gateway *Gateway) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "edit_gateway_custom_routes",
		"gateway_name": gateway.GwName,
		"cidr":         strings.Join(gateway.CustomizedSpokeVpcRoutes, ","),
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) EditGatewayFilterRoutes(gateway *Gateway) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "edit_gateway_filter_routes",
		"gateway_name": gateway.GwName,
		"cidr":         strings.Join(gateway.FilteredSpokeVpcRoutes, ","),
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) EditGatewayAdvertisedCidr(gateway *Gateway) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "edit_gateway_advertised_cidr",
		"gateway_name": gateway.GwName,
		"cidr":         strings.Join(gateway.AdvertisedSpokeRoutes, ","),
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) EnableTransitFireNet(gateway *Gateway) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "enable_gateway_for_transit_firenet",
		"gateway_name": gateway.GwName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) EnableTransitFireNetWithGWLB(gateway *Gateway) error {
	data := map[string]string{
		"CID":          c.CID,
		"action":       "enable_gateway_for_transit_firenet",
		"gateway_name": gateway.GwName,
		"mode":         "gwlb",
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) DisableTransitFireNet(gateway *Gateway) error {
	err := c.IsTransitFireNetReadyToBeDisabled(gateway)
	if err != nil {
		return err
	}

	form := map[string]string{
		"CID":          c.CID,
		"action":       "disable_gateway_for_transit_firenet",
		"gateway_name": gateway.GwName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) IsTransitFireNetReadyToBeDisabled(gateway *Gateway) error {
	form := map[string]string{
		"CID":    c.CID,
		"action": "list_transit_firenet_spoke_policies",
	}

	var data ListTransitFireNetPolicyResp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return err
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
	form := map[string]string{
		"CID":    c.CID,
		"action": "list_transit_gateways_for_multi_cloud_domains",
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

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return false, err
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

func (c *Client) EnableMonitorGatewaySubnets(gwName string, excludedInstances []string) error {
	action := "enable_monitor_gateway_subnets"
	form := map[string]string{
		"CID":          c.CID,
		"action":       action,
		"gateway_name": gwName,
	}
	if len(excludedInstances) != 0 {
		form["monitor_exclude_gateway_list"] = strings.Join(excludedInstances, ",")
	}
	return c.PostAPI(action, form, BasicCheck)
}

func (c *Client) DisableMonitorGatewaySubnets(gwName string) error {
	action := "disable_monitor_gateway_subnets"
	form := map[string]string{
		"CID":          c.CID,
		"action":       action,
		"gateway_name": gwName,
	}
	check := func(action, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "no change needed") {
				return nil
			}
			return fmt.Errorf("rest API %s %s failed: %s", action, method, reason)
		}
		return nil
	}
	return c.PostAPI(action, form, check)
}

func (c *Client) EnableVPNConfig(gateway *Gateway, vpnConfig *VPNConfig) error {
	action := "edit_vpn_config"
	form := map[string]interface{}{
		"CID":     c.CID,
		"action":  action,
		"command": "enable",
		"vpc_id":  gateway.VpcID,
		"lb_name": gateway.GwName,
		"key":     vpnConfig.Name,
		"value":   vpnConfig.Value,
	}
	return c.PostAPI(action, form, BasicCheck)
}

func (c *Client) DisableVPNConfig(gateway *Gateway, vpnConfig *VPNConfig) error {
	action := "edit_vpn_config"
	form := map[string]interface{}{
		"CID":     c.CID,
		"action":  action,
		"command": "disable",
		"vpc_id":  gateway.VpcID,
		"lb_name": gateway.GwName,
		"key":     vpnConfig.Name,
		"value":   "-1",
	}
	return c.PostAPI(action, form, BasicCheck)
}

func (c *Client) GetVPNConfigList(gateway *Gateway) ([]VPNConfig, error) {
	form := map[string]string{
		"CID":     c.CID,
		"action":  "edit_vpn_config",
		"command": "show",
		"vpc_id":  gateway.VpcID,
		"lb_name": gateway.GwName,
	}

	var data VPNConfigListResp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	return data.Results, ErrNotFound
}

func (c *Client) EnableActiveStandby(transitGateway *TransitVpc) error {
	action := "enable_active_standby"
	form := map[string]string{
		"CID":          c.CID,
		"action":       action,
		"gateway_name": transitGateway.GwName,
	}
	return c.PostAPI(action, form, BasicCheck)
}

func (c *Client) DisableActiveStandby(transitGateway *TransitVpc) error {
	action := "disable_active_standby"
	form := map[string]string{
		"CID":          c.CID,
		"action":       action,
		"gateway_name": transitGateway.GwName,
	}
	return c.PostAPI(action, form, BasicCheck)
}

func (c *Client) SwitchActiveTransitGateway(gwName, connName string) error {
	action := "active_standby_connection_switchover"
	form := map[string]string{
		"CID":             c.CID,
		"action":          action,
		"gateway_name":    gwName,
		"connection_name": connName,
	}
	return c.PostAPI(action, form, BasicCheck)
}

func (c *Client) GetTransitGatewayLanCidr(gatewayName string) (string, error) {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "get_firewall_lan_cidr",
		"gateway_name": gatewayName,
	}

	type LANCidr struct {
		FirewallLanCidr string `form:"firewall_lan_cidr,omitempty" json:"firewall_lan_cidr,omitempty"`
	}

	type LANCidrResp struct {
		Return  bool    `json:"return"`
		Results LANCidr `json:"results"`
		Reason  string  `json:"reason"`
	}

	var data LANCidrResp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return "", err
	}

	return data.Results.FirewallLanCidr, ErrNotFound
}

func (c *Client) GetFqdnGatewayInfo(gateway *Gateway) (*FQDNGatwayInfo, error) {
	params := map[string]string{
		"action":    "list_firenet",
		"subaction": "instance",
		"vpc_id":    gateway.VpcID,
		"CID":       c.CID,
	}
	var data FQDNGatewayInfoResp
	err := c.GetAPI(&data, params["action"], params, BasicCheck)
	if err != nil {
		return nil, err
	}
	return &data.Results, nil
}

func (c *Client) UpdateTransitGatewayCustomizedVpcRoute(gateway string, customizedTransitVpcRoutes []string) error {
	params := map[string]string{
		"action":            "edit_transit_gateway_customized_vpc_route",
		"CID":               c.CID,
		"gateway_name":      gateway,
		"customized_routes": strings.Join(customizedTransitVpcRoutes, ","),
	}

	return c.PostAPI(params["action"], params, BasicCheck)
}

func (c *Client) EnableJumboFrame(gateway *Gateway) error {
	action := "enable_jumbo_frame"
	form := map[string]string{
		"CID":          c.CID,
		"action":       action,
		"gateway_name": gateway.GwName,
	}
	return c.PostAPI(action, form, BasicCheck)
}

func (c *Client) DisableJumboFrame(gateway *Gateway) error {
	action := "disable_jumbo_frame"
	form := map[string]string{
		"CID":          c.CID,
		"action":       action,
		"gateway_name": gateway.GwName,
	}
	return c.PostAPI(action, form, BasicCheck)
}

func (c *Client) GetJumboFrameStatus(gateway *Gateway) (bool, error) {
	action := "get_jumbo_frame_status"
	form := map[string]string{
		"CID":          c.CID,
		"action":       action,
		"gateway_name": gateway.GwName,
	}

	type JumboFrameResult struct {
		Return  bool   `json:"return"`
		Results string `json:"results"`
		Reason  string `json:"reason"`
	}

	var resp JumboFrameResult
	err := c.GetAPI(&resp, form["action"], form, BasicCheck)
	if err != nil {
		return false, err
	}
	return strings.Contains(resp.Results, "Jumbo frame is enabled"), nil
}

func (c *Client) EnablePrivateVpcDefaultRoute(gw *Gateway) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "enable_private_vpc_default_route",
		"gateway_name": gw.GwName,
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) DisablePrivateVpcDefaultRoute(gw *Gateway) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "disable_private_vpc_default_route",
		"gateway_name": gw.GwName,
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) EnableSkipPublicRouteUpdate(gw *Gateway) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "enable_skip_public_route_table_update",
		"gateway_name": gw.GwName,
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) DisableSkipPublicRouteUpdate(gw *Gateway) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "disable_skip_public_route_table_update",
		"gateway_name": gw.GwName,
	}
	return c.PostAPI(form["action"], form, BasicCheck)
}

// Entity should be gateway name or "Controller"
func (c *Client) GetTunnelDetectionTime(entity string) (int, error) {
	form := map[string]string{
		"CID":    c.CID,
		"action": "show_tunnel_status_change_detection_time",
		"entity": entity,
	}

	type DetectionTimeResults struct {
		DetectionTime int `json:"detection_time"`
	}

	type DetectionTimeResp struct {
		Return  bool                 `json:"return"`
		Results DetectionTimeResults `json:"results"`
		Reason  string               `json:"reason"`
	}

	var resp DetectionTimeResp
	err := c.GetAPI(&resp, form["action"], form, BasicCheck)
	if err != nil {
		return 0, err
	}
	return resp.Results.DetectionTime, err
}

func (c *Client) ModifyTunnelDetectionTime(entity string, detectionTime int) error {
	form := map[string]string{
		"CID":            c.CID,
		"action":         "modify_detection_time",
		"detection_time": strconv.Itoa(detectionTime),
		"entity":         entity,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) GetBgpManualSpokeAdvertiseCidrs(gateway *Gateway) ([]string, error) {
	form := map[string]string{
		"CID":                  c.CID,
		"action":               "list_aviatrix_transit_advanced_config",
		"transit_gateway_name": gateway.GwName,
	}

	type AdvancedConfigResults struct {
		BgpManualSpokeAdvertiseCidrs []string `json:"bgp_Manual_spoke_advertise_cidrs"`
	}

	type AdvancedConfigResp struct {
		Return  bool
		Results AdvancedConfigResults
		Reason  string
	}

	var resp AdvancedConfigResp
	err := c.GetAPI(&resp, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}
	return resp.Results.BgpManualSpokeAdvertiseCidrs, nil
}
