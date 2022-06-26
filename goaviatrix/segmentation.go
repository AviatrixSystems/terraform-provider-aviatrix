package goaviatrix

import "strings"

type SegmentationSecurityDomain struct {
	DomainName string
}

type SegmentationSecurityDomainConnectionPolicy struct {
	Domain1 *SegmentationSecurityDomain
	Domain2 *SegmentationSecurityDomain
}

type SegmentationSecurityDomainAssociation struct {
	TransitGatewayName string
	SecurityDomainName string
	AttachmentName     string
}

func (c *Client) CreateSegmentationSecurityDomain(domain *SegmentationSecurityDomain) error {
	action := "add_multi_cloud_security_domain"
	data := map[string]interface{}{
		"action":      action,
		"CID":         c.CID,
		"domain_name": domain.DomainName,
	}
	return c.PostAPI(action, data, BasicCheck)
}

func (c *Client) DeleteSegmentationSecurityDomain(domain *SegmentationSecurityDomain) error {
	action := "delete_multi_cloud_security_domain"
	data := map[string]interface{}{
		"action":      action,
		"CID":         c.CID,
		"domain_name": domain.DomainName,
	}
	return c.PostAPI(action, data, BasicCheck)
}

func (c *Client) GetSegmentationSecurityDomain(domain *SegmentationSecurityDomain) (*SegmentationSecurityDomain, error) {
	form := map[string]string{
		"CID":    c.CID,
		"action": "list_multi_cloud_security_domain_names",
	}

	type Resp struct {
		Return  bool     `json:"return"`
		Results []string `json:"results"`
		Reason  string   `json:"reason"`
	}

	var data Resp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	if !Contains(data.Results, domain.DomainName) {
		return nil, ErrNotFound
	}

	return domain, nil
}

func (c *Client) CreateSegmentationSecurityDomainConnectionPolicy(policy *SegmentationSecurityDomainConnectionPolicy) error {
	action := "connect_multi_cloud_security_domains"
	data := map[string]interface{}{
		"action":            action,
		"CID":               c.CID,
		"domain_name":       policy.Domain1.DomainName,
		"other_domain_name": policy.Domain2.DomainName,
	}
	return c.PostAPI(action, data, BasicCheck)
}

func (c *Client) DeleteSegmentationSecurityDomainConnectionPolicy(policy *SegmentationSecurityDomainConnectionPolicy) error {
	action := "disconnect_multi_cloud_security_domains"
	data := map[string]interface{}{
		"action":            action,
		"CID":               c.CID,
		"domain_name":       policy.Domain1.DomainName,
		"other_domain_name": policy.Domain2.DomainName,
	}
	return c.PostAPI(action, data, BasicCheck)
}

func (c *Client) GetSegmentationSecurityDomainConnectionPolicy(policy *SegmentationSecurityDomainConnectionPolicy) (*SegmentationSecurityDomainConnectionPolicy, error) {
	form := map[string]string{
		"CID":         c.CID,
		"action":      "list_multi_cloud_security_domain_connection_policy",
		"domain_name": policy.Domain1.DomainName,
	}

	type Result struct {
		ConnectedDomains    []string `json:"connected_domains"`
		NotConnectedDomains []string `json:"not_connected_domains"`
	}

	type Resp struct {
		Return  bool   `json:"return"`
		Results Result `json:"results"`
		Reason  string `json:"reason"`
	}

	var data Resp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	// Check if the other domain is included in the list of connected domains
	if !Contains(data.Results.ConnectedDomains, policy.Domain2.DomainName) {
		return nil, ErrNotFound
	}

	return policy, nil
}

func (c *Client) CreateSegmentationSecurityDomainAssociation(association *SegmentationSecurityDomainAssociation) error {
	action := "associate_attachment_to_multi_cloud_security_domain"
	data := map[string]interface{}{
		"action":          action,
		"CID":             c.CID,
		"attachment_name": association.AttachmentName,
		"domain_name":     association.SecurityDomainName,
	}
	return c.PostAPI(action, data, BasicCheck)
}

func (c *Client) DeleteSegmentationSecurityDomainAssociation(association *SegmentationSecurityDomainAssociation) error {
	action := "disassociate_attachment_from_multi_cloud_security_domain"
	data := map[string]interface{}{
		"action":          action,
		"CID":             c.CID,
		"attachment_name": association.AttachmentName,
		"domain_name":     association.SecurityDomainName,
	}
	return c.PostAPI(action, data, BasicCheck)
}

func (c *Client) GetSegmentationSecurityDomainAssociation(association *SegmentationSecurityDomainAssociation) (*SegmentationSecurityDomainAssociation, error) {
	form := map[string]string{
		"CID":    c.CID,
		"action": "list_multi_cloud_domain_attachments",
	}

	type Attachment struct {
		Name        string `json:"name"`
		Domain      string `json:"domain"`
		TransitName string `json:"transit_name"`
		Type        string `json:"type"`
	}

	type Result struct {
		Attachments []Attachment `json:"attachments"`
	}

	type Resp struct {
		Return  bool   `json:"return"`
		Results Result `json:"results"`
		Reason  string `json:"reason"`
	}

	var data Resp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	found := false
	for _, attachment := range data.Results.Attachments {
		if attachment.Type == "EDGESPOKE" {
			attachment.Name = strings.Split(attachment.Name, ":")[0]
		}

		if attachment.Domain == association.SecurityDomainName &&
			attachment.Name == association.AttachmentName &&
			attachment.TransitName == association.TransitGatewayName {
			found = true
		}
	}

	if !found {
		return nil, ErrNotFound
	}

	return association, nil
}
