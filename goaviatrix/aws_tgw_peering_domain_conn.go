package goaviatrix

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DomainConn struct {
	Action      string `form:"action,omitempty"`
	CID         string `form:"CID,omitempty"`
	TgwName1    string `form:"tgw_name1,omitempty" json:"tgw_name1,omitempty"`
	DomainName1 string
	TgwName2    string `form:"tgw_name2,omitempty" json:"tgw_name2,omitempty"`
	DomainName2 string
}

type ListConnectedRouteDomainsResp struct {
	Return  bool                       `json:"return"`
	Results ConnectedRouteDomainDetail `json:"results"`
	Reason  string                     `json:"reason"`
}

type ConnectedRouteDomainDetail struct {
	ConnectedDomainNames    []string `json:"connected_domain_names"`
	NotConnectedDomainNames []string `json:"not_connected_domain_names"`
	Egress                  string   `json:"egress"`
}

func (c *Client) CreateDomainConn(domainConn *DomainConn) error {
	form := map[string]string{
		"CID":                           c.CID,
		"action":                        "add_connection_between_route_domains",
		"tgw_name":                      domainConn.TgwName1,
		"source_route_domain_name":      domainConn.DomainName1,
		"destination_route_domain_name": domainConn.TgwName2 + ":" + domainConn.DomainName2,
		"async":                         "true",
	}
	return c.PostAsyncAPI(form["action"], form, BasicCheck)
}

func (c *Client) GetDomainConn(domainConn *DomainConn) error {
	var data ListConnectedRouteDomainsResp
	form := map[string]string{
		"CID":               c.CID,
		"action":            "list_connected_route_domains",
		"tgw_name":          domainConn.TgwName1,
		"route_domain_name": domainConn.DomainName1,
	}
	check := func(action, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "does not exist") {
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s %s failed: %s", action, method, reason)
		}
		return nil
	}
	err := c.GetAPI(&data, form["action"], form, check)
	if err != nil {
		return err
	}
	connectedDomains := data.Results.ConnectedDomainNames
	for i := range connectedDomains {
		if connectedDomains[i] == domainConn.TgwName2+":"+domainConn.DomainName2 {
			return nil
		}
	}
	return ErrNotFound
}

func (c *Client) DeleteDomainConn(domainConn *DomainConn) error {
	form := map[string]string{
		"CID":                           c.CID,
		"action":                        "delete_connection_between_route_domains",
		"tgw_name":                      domainConn.TgwName1,
		"source_route_domain_name":      domainConn.DomainName1,
		"destination_route_domain_name": domainConn.TgwName2 + ":" + domainConn.DomainName2,
		"async":                         "true",
	}
	return c.PostAsyncAPI(form["action"], form, BasicCheck)
}

func DiffSuppressFuncAwsTgwPeeringDomainConnTgwName1(k, old, new string, d *schema.ResourceData) bool {
	tgwName2Old, _ := d.GetChange("tgw_name2")
	domainName1Old, _ := d.GetChange("domain_name1")
	domainName2Old, _ := d.GetChange("domain_name2")

	tgw2Cur, ok := d.Get("tgw_name2").(string)
	if !ok {
		return false
	}
	tgw2Old, ok := tgwName2Old.(string)
	if !ok {
		return false
	}

	dn1Cur, ok := d.Get("domain_name1").(string)
	if !ok {
		return false
	}
	dn2Cur, ok := d.Get("domain_name2").(string)
	if !ok {
		return false
	}
	dn1Old, ok := domainName1Old.(string)
	if !ok {
		return false
	}
	dn2Old, ok := domainName2Old.(string)
	if !ok {
		return false
	}

	return old == tgw2Cur &&
		new == tgw2Old &&
		dn1Cur == dn2Old &&
		dn2Cur == dn1Old
}

func DiffSuppressFuncAwsTgwPeeringDomainConnTgwName2(k, old, new string, d *schema.ResourceData) bool {
	tgwName1Old, _ := d.GetChange("tgw_name1")
	domainName1Old, _ := d.GetChange("domain_name1")
	domainName2Old, _ := d.GetChange("domain_name2")

	tgw1Cur, ok := d.Get("tgw_name1").(string)
	if !ok {
		return false
	}
	tgw1Old, ok := tgwName1Old.(string)
	if !ok {
		return false
	}

	dn1Cur, ok := d.Get("domain_name1").(string)
	if !ok {
		return false
	}
	dn2Cur, ok := d.Get("domain_name2").(string)
	if !ok {
		return false
	}
	dn1Old, ok := domainName1Old.(string)
	if !ok {
		return false
	}
	dn2Old, ok := domainName2Old.(string)
	if !ok {
		return false
	}

	return old == tgw1Cur &&
		new == tgw1Old &&
		dn1Cur == dn2Old &&
		dn2Cur == dn1Old
}

func DiffSuppressFuncAwsTgwPeeringDomainConnDomainName1(k, old, new string, d *schema.ResourceData) bool {
	tgwName1Old, _ := d.GetChange("tgw_name1")
	tgwName2Old, _ := d.GetChange("tgw_name2")
	domainName2Old, _ := d.GetChange("domain_name2")

	dn2Cur, ok := d.Get("domain_name2").(string)
	if !ok {
		return false
	}
	dn2Old, ok := domainName2Old.(string)
	if !ok {
		return false
	}

	tgw1Cur, ok := d.Get("tgw_name1").(string)
	if !ok {
		return false
	}
	tgw2Cur, ok := d.Get("tgw_name2").(string)
	if !ok {
		return false
	}
	tgw1Old, ok := tgwName1Old.(string)
	if !ok {
		return false
	}
	tgw2Old, ok := tgwName2Old.(string)
	if !ok {
		return false
	}

	return old == dn2Cur &&
		new == dn2Old &&
		tgw1Cur == tgw2Old &&
		tgw2Cur == tgw1Old
}

func DiffSuppressFuncAwsTgwPeeringDomainConnDomainName2(k, old, new string, d *schema.ResourceData) bool {
	tgwName1Old, _ := d.GetChange("tgw_name1")
	tgwName2Old, _ := d.GetChange("tgw_name2")
	domainName1Old, _ := d.GetChange("domain_name1")

	dn1Cur, ok := d.Get("domain_name1").(string)
	if !ok {
		return false
	}
	dn1Old, ok := domainName1Old.(string)
	if !ok {
		return false
	}

	tgw1Cur, ok := d.Get("tgw_name1").(string)
	if !ok {
		return false
	}
	tgw2Cur, ok := d.Get("tgw_name2").(string)
	if !ok {
		return false
	}
	tgw1Old, ok := tgwName1Old.(string)
	if !ok {
		return false
	}
	tgw2Old, ok := tgwName2Old.(string)
	if !ok {
		return false
	}

	return old == dn1Cur &&
		new == dn1Old &&
		tgw1Cur == tgw2Old &&
		tgw2Cur == tgw1Old
}
