package goaviatrix

import (
	"context"
)

// AwsTGW simple struct to hold aws_tgw details
type SecurityDomain struct {
	Action                 string `form:"action, omitempty"`
	CID                    string `form:"CID, omitempty"`
	Name                   string `form:"route_domain_name, omitempty"`
	AccountName            string `form:"account_name, omitempty"`
	Region                 string `form:"region, omitempty"`
	AwsTgwName             string `form:"tgw_name, omitempty"`
	AviatrixFirewallDomain bool   `form:"firewall_domain, omitempty"`
	NativeEgressDomain     bool   `form:"native_egress_domain, omitempty"`
	NativeFirewallDomain   bool   `form:"native_firewall_domain, omitempty"`
	ForceDelete            bool   `form:"force,omitempty"`
	Async                  bool   `form:"async,omitempty"`
}

type SecurityDomainAPIResp struct {
	Return  bool     `json:"return"`
	Results []string `json:"results"`
	Reason  string   `json:"reason"`
}

type SecurityDomainList struct {
	Return  bool     `json:"return"`
	Results []string `json:"results"`
	Reason  string   `json:"reason"`
}

type SecurityDomainRule struct {
	Name                   string    `json:"security_domain_name,omitempty"`
	ConnectedDomain        []string  `json:"connected_domains,omitempty"`
	AttachedVPCs           []VPCSolo `json:"attached_vpc,omitempty"`
	AviatrixFirewallDomain bool      `json:"firewall_domain,omitempty"`
	NativeEgressDomain     bool      `json:"egress_domain,omitempty"`
	NativeFirewallDomain   bool      `json:"native_firewall_domain,omitempty"`
}

type SecurityDomainDetails struct {
	Name                   string    `json:"name"`
	ConnectedDomain        []string  `json:"connected_route_domain,omitempty"`
	AttachedVPCs           []VPCSolo `json:"attached_vpc,omitempty"`
	AviatrixFirewallDomain bool      `json:"firewall_domain,omitempty"`
	NativeEgressDomain     bool      `json:"egress_domain,omitempty"`
	NativeFirewallDomain   bool      `json:"native_firewall_domain,omitempty"`
}

type VPCSolo struct {
	Region                       string `json:"vpc_region,omitempty"`
	AccountName                  string `json:"vpc_account_name,omitempty"`
	VpcID                        string `json:"vpc_id,omitempty"`
	Subnets                      string
	RouteTables                  string
	CustomizedRoutes             string `json:",omitempty"`
	CustomizedRouteAdvertisement string
	DisableLocalRoutePropagation bool `json:",omitempty"`
}

type IntraDomainInspection struct {
	TgwName                      string `form:"tgw_name,omitempty" json:"tgw_name"`
	RouteDomainName              string `form:"route_domain_name,omitempty" json:"name"`
	FirewallDomainName           string `form:"firewall_domain_name,omitempty" json:"intra_domain_inspection_name"`
	IntraDomainInspectionEnabled bool   `json:"intra_domain_inspection"`
}

type DomainDetails struct {
	Name                         string `json:"name"`
	TgwName                      string `json:"tgw_name"`
	RouteTableId                 string `json:"route_table_id"`
	Account                      string `json:"account"`
	CouldType                    string `json:"cloud_type"`
	Region                       string `json:"region"`
	IntraDomainInspectionEnabled bool   `json:"intra_domain_inspection"`
	EgressInspection             bool   `json:"egress_inspection"`
	InspectionPolicy             string `json:"inspection_policy"`
	IntraDomainInspectionName    string `json:"intra_domain_inspection_name"`
	EgressInspectionName         string `json:"egress_inspection_name"`
	Type                         string `json:"type"`
}

func (c *Client) CreateSecurityDomain(securityDomain *SecurityDomain) error {
	securityDomain.CID = c.CID
	securityDomain.Action = "add_route_domain"
	securityDomain.Async = true

	return c.PostAsyncAPI(securityDomain.Action, securityDomain, BasicCheck)
}

func (c *Client) GetSecurityDomain(securityDomain *SecurityDomain) (string, error) {
	form := map[string]string{
		"CID":      c.CID,
		"action":   "list_route_domain_names",
		"tgw_name": securityDomain.AwsTgwName,
	}

	data := SecurityDomainAPIResp{
		Return:  false,
		Results: make([]string, 0),
		Reason:  "",
	}

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return "", err
	}

	securityDomainList := data.Results
	for i := range securityDomainList {
		if securityDomainList[i] == securityDomain.Name {
			return securityDomainList[i], nil
		}
	}
	return "", ErrNotFound
}

func (c *Client) UpdateSecurityDomain(securityDomain *SecurityDomain) error {
	return nil
}

func (c *Client) DeleteSecurityDomain(securityDomain *SecurityDomain) error {
	securityDomain.CID = c.CID
	securityDomain.Action = "delete_route_domain"

	return c.PostAPI(securityDomain.Action, securityDomain, BasicCheck)
}

func (c *Client) CreateDomainConnection(awsTgw *AWSTgw, sourceDomain string, destinationDomain string) error {
	form := map[string]string{
		"CID":                           c.CID,
		"action":                        "add_connection_between_route_domains",
		"account_name":                  awsTgw.AccountName,
		"region":                        awsTgw.Region,
		"tgw_name":                      awsTgw.Name,
		"source_route_domain_name":      sourceDomain,
		"destination_route_domain_name": destinationDomain,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) DeleteDomainConnection(awsTgw *AWSTgw, sourceDomain string, destinationDomain string) error {
	form := map[string]string{
		"CID":                           c.CID,
		"action":                        "delete_connection_between_route_domains",
		"tgw_name":                      awsTgw.Name,
		"source_route_domain_name":      sourceDomain,
		"destination_route_domain_name": destinationDomain,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) SecurityDomainRuleValidation(securityDomainRule *SecurityDomainRule) bool {
	num := 0
	if securityDomainRule.AviatrixFirewallDomain {
		num += 1
	}
	if securityDomainRule.NativeEgressDomain {
		num += 1
	}
	if securityDomainRule.NativeFirewallDomain {
		num += 1
	}
	if num > 1 {
		return false
	}
	return true
}

func (c *Client) GetSecurityDomainDetails(ctx context.Context, domain *SecurityDomain) (*SecurityDomainDetails, error) {
	params := map[string]string{
		"action":            "list_tgw_security_domain_details",
		"CID":               c.CID,
		"tgw_name":          domain.AwsTgwName,
		"route_domain_name": domain.Name,
	}

	type Resp struct {
		Return  bool                    `json:"return"`
		Results []SecurityDomainDetails `json:"results"`
		Reason  string                  `json:"reason"`
	}

	var data Resp

	err := c.GetAPIContext(ctx, &data, params["action"], params, BasicCheck)
	if err != nil {
		return nil, err
	}

	if len(data.Results) == 0 {
		return nil, ErrNotFound
	}

	return &data.Results[0], nil
}

func (c *Client) EnableIntraDomainInspection(ctx context.Context, intraDomainInspection *IntraDomainInspection) error {
	params := map[string]string{
		"action":               "enable_tgw_intra_domain_inspection",
		"CID":                  c.CID,
		"tgw_name":             intraDomainInspection.TgwName,
		"route_domain_name":    intraDomainInspection.RouteDomainName,
		"firewall_domain_name": intraDomainInspection.FirewallDomainName,
	}

	return c.PostAPIContext(ctx, params["action"], params, BasicCheck)
}

func (c *Client) DisableIntraDomainInspection(ctx context.Context, intraDomainInspection *IntraDomainInspection) error {
	params := map[string]string{
		"action":            "disable_tgw_intra_domain_inspection",
		"CID":               c.CID,
		"tgw_name":          intraDomainInspection.TgwName,
		"route_domain_name": intraDomainInspection.RouteDomainName,
	}

	return c.PostAPIContext(ctx, params["action"], params, BasicCheck)
}

func (c *Client) GetAllSecurityDomains() ([]DomainDetails, error) {
	params := map[string]string{
		"action": "list_all_tgw_security_domains",
		"CID":    c.CID,
	}

	type DomainDetail struct {
		Domains []DomainDetails `json:"domains"`
	}

	type Resp struct {
		Return  bool         `json:"return"`
		Results DomainDetail `json:"results"`
		Reason  string       `json:"reason"`
	}

	var data Resp
	err := c.GetAPI(&data, params["action"], params, BasicCheck)
	if err != nil {
		return nil, err
	}
	domainList := data.Results.Domains
	return domainList, nil
}

func (c *Client) GetIntraDomainInspectionStatus(ctx context.Context, intraDomainInspection *IntraDomainInspection) error {
	params := map[string]string{
		"action": "list_all_tgw_security_domains",
		"CID":    c.CID,
	}

	type DomainDetails struct {
		Domains []IntraDomainInspection `json:"domains"`
	}

	type Resp struct {
		Return  bool          `json:"return"`
		Results DomainDetails `json:"results"`
		Reason  string        `json:"reason"`
	}

	var data Resp

	err := c.GetAPIContext(ctx, &data, params["action"], params, BasicCheck)
	if err != nil {
		return err
	}

	for _, domain := range data.Results.Domains {
		if domain.TgwName == intraDomainInspection.TgwName && domain.RouteDomainName == intraDomainInspection.RouteDomainName {
			if !domain.IntraDomainInspectionEnabled {
				return ErrNotFound
			}

			intraDomainInspection.IntraDomainInspectionEnabled = domain.IntraDomainInspectionEnabled
			intraDomainInspection.FirewallDomainName = domain.FirewallDomainName

			return nil
		}
	}

	return ErrNotFound
}
