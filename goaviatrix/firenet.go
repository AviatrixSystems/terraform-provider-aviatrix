package goaviatrix

import (
	"fmt"
	"strconv"
	"strings"
)

type FireNet struct {
	CID               string `form:"CID,omitempty"`
	Action            string `form:"action,omitempty"`
	VpcID             string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	GwName            string `form:"gw_name,omitempty" json:"gw_name,omitempty"`
	FirewallInstance  []FirewallInstance
	FirewallEgress    bool   `form:"firewall_egress,omitempty" json:"firewall_egress,omitempty"`
	Inspection        bool   `form:"inspection,omitempty" json:"inspection,omitempty"`
	HashingAlgorithm  string `json:"firewall_hashing,omitempty"`
	EgressStaticCidrs string
	ExcludedCidrs     string
}

type FireNetDetail struct {
	CloudType                string                 `json:"cloud_type"`
	Region                   string                 `json:"region,omitempty"`
	VpcID                    string                 `json:"vpc_id,omitempty"`
	FirewallInstance         []FirewallInstanceInfo `json:"firewall,omitempty"`
	Gateway                  []GatewayInfo          `json:"gateway,omitempty"`
	FirewallEgress           string                 `json:"firewall_egress,omitempty"`
	NativeGwlb               bool                   `json:"native_gwlb"`
	Inspection               string                 `json:"inspection,omitempty"`
	HashingAlgorithm         string                 `json:"firewall_hashing,omitempty"`
	LanPing                  string                 `json:"lan_ping"`
	TgwSegmentationForEgress string                 `json:"tgw_segmentation"`
	EgressStaticCidrs        []string               `json:"egress_static_cidr"`
	ExcludedCidrs            []string               `json:"exclude_cidr"`
	FailClose                string                 `json:"fail_close"`
}

type GetFireNetResp struct {
	Return  bool          `json:"return"`
	Results FireNetDetail `json:"results"`
	Reason  string        `json:"reason"`
}

type ListFireNetResp struct {
	Return  bool              `json:"return"`
	Results FirewallInterface `json:"results"`
	Reason  string            `json:"reason"`
}

type FirewallInterface struct {
	Instances  []string            `json:"instances"`
	Interfaces map[string][]string `json:"interfaces"`
}

type GatewayInfo struct {
	DomainName string `json:"domain_name"`
	HaStatus   string `json:"ha_status"`
	GwName     string `json:"name"`
	TgwID      string `json:"tgw_id"`
}

type FirewallInstanceInfo struct {
	Enabled             bool   `json:"enabled"`
	GwName              string `json:"gateway"`
	InstanceID          string `json:"id"`
	FirewallName        string `json:"name"`
	LanInterface        string `json:"lan_interface_id,omitempty"`
	ManagementInterface string `json:"management_interface_id,omitempty"`
	EgressInterface     string `json:"egress_interface_id,omitempty"`
	VendorType          string `json:"vendor,omitempty"`
}

func (c *Client) CreateFireNet(fireNet *FireNet) error {
	return nil
}

func (c *Client) GetFireNet(fireNet *FireNet) (*FireNetDetail, error) {
	form := map[string]string{
		"CID":    c.CID,
		"action": "show_firenet_detail",
		"vpc_id": fireNet.VpcID,
	}

	var data GetFireNetResp

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "not found in DB") {
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	err := c.GetAPI(&data, form["action"], form, checkFunc)
	if err != nil {
		return nil, err
	}

	if strings.Split(data.Results.VpcID, "~~")[0] == fireNet.VpcID {
		data.Results.VpcID = fireNet.VpcID
		return &data.Results, nil
	}
	return nil, ErrNotFound
}

func (c *Client) AssociateFirewallWithFireNet(firewallInstance *FirewallInstance) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "associate_firewall_with_firenet",
		"vpc_id":       firewallInstance.VpcID,
		"gateway_name": firewallInstance.GwName,
		"firewall_id":  firewallInstance.InstanceID,
	}

	if firewallInstance.LanInterface != "" {
		form["lan_interface"] = firewallInstance.LanInterface
	}
	if firewallInstance.VendorType == "Generic" {
		form["firewall_name"] = firewallInstance.FirewallName
		form["management_interface"] = firewallInstance.ManagementInterface
		form["egress_interface"] = firewallInstance.EgressInterface
	}

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "already associated") {
				return nil
			}
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	return c.PostAPI(form["action"], form, checkFunc)
}

func (c *Client) DisassociateFirewallFromFireNet(firewallInstance *FirewallInstance) error {
	form := map[string]string{
		"CID":         c.CID,
		"action":      "disassociate_firewall_with_firenet",
		"vpc_id":      firewallInstance.VpcID,
		"firewall_id": firewallInstance.InstanceID,
	}

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "not found") {
				return nil
			}
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	return c.PostAPI(form["action"], form, checkFunc)
}

func (c *Client) AttachFirewallToFireNet(firewallInstance *FirewallInstance) error {
	form := map[string]string{
		"CID":         c.CID,
		"action":      "attach_firewall_to_firenet",
		"vpc_id":      firewallInstance.VpcID,
		"firewall_id": firewallInstance.InstanceID,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) DetachFirewallFromFireNet(firewallInstance *FirewallInstance) error {
	form := map[string]string{
		"CID":         c.CID,
		"action":      "detach_firewall_from_firenet",
		"vpc_id":      firewallInstance.VpcID,
		"firewall_id": firewallInstance.InstanceID,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) ConnectFireNetWithTgw(awsTgw *AWSTgw, vpcSolo VPCSolo, SecurityDomainName string) error {
	form := map[string]string{
		"CID":         c.CID,
		"action":      "connect_firenet_with_tgw",
		"vpc_id":      vpcSolo.VpcID,
		"tgw_name":    awsTgw.Name,
		"domain_name": SecurityDomainName,
		"async":       "true",
	}

	return c.PostAsyncAPI(form["action"], form, BasicCheck)
}

func (c *Client) DisconnectFireNetFromTgw(awsTgw *AWSTgw, vpcID string) error {
	form := map[string]string{
		"CID":    c.CID,
		"action": "disconnect_firenet_with_tgw",
		"vpc_id": vpcID,
		"async":  "true",
	}

	return c.PostAsyncAPI(form["action"], form, BasicCheck)
}

func (c *Client) EditFireNetInspection(fireNet *FireNet) error {
	form := map[string]string{
		"CID":        c.CID,
		"action":     "edit_firenet",
		"vpc_id":     fireNet.VpcID,
		"inspection": strconv.FormatBool(fireNet.Inspection),
	}

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "configuration not changed") {
				return nil
			}
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	return c.PostAPI(form["action"], form, checkFunc)
}

func (c *Client) EditFireNetEgress(fireNet *FireNet) error {
	form := map[string]string{
		"CID":             c.CID,
		"action":          "edit_firenet",
		"vpc_id":          fireNet.VpcID,
		"firewall_egress": strconv.FormatBool(fireNet.FirewallEgress),
	}

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "configuration not changed") {
				return nil
			}
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	return c.PostAPI(form["action"], form, checkFunc)
}

func (c *Client) EditFireNetHashingAlgorithm(fireNet *FireNet) error {
	data := map[string]string{
		"action":           "edit_firenet",
		"CID":              c.CID,
		"vpc_id":           fireNet.VpcID,
		"firewall_hashing": fireNet.HashingAlgorithm,
	}

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "configuration not changed") {
				return nil
			}
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	return c.PostAPI("edit_firenet", data, checkFunc)
}

func (c *Client) EnableFireNetLanKeepAlive(net *FireNet) error {
	data := map[string]string{
		"action":   "edit_firenet",
		"CID":      c.CID,
		"vpc_id":   net.VpcID,
		"lan_ping": "true",
	}

	customCheck := func(action, method, reason string, ret bool) error {
		if !ret {
			// [AVXERR-FIRENET-0029] No change in firenet attribute
			if strings.Contains(reason, "AVXERR-FIRENET-0029") {
				return nil
			}
			return fmt.Errorf("rest API %s Post failed: %s", action, reason)
		}
		return nil
	}

	return c.PostAPI("edit_firenet(lan_ping=true)", data, customCheck)
}

func (c *Client) DisableFireNetLanKeepAlive(net *FireNet) error {
	data := map[string]string{
		"action":   "edit_firenet",
		"CID":      c.CID,
		"vpc_id":   net.VpcID,
		"lan_ping": "false",
	}

	customCheck := func(action, method, reason string, ret bool) error {
		if !ret {
			// [AVXERR-FIRENET-0029] No change in firenet attribute
			if strings.Contains(reason, "AVXERR-FIRENET-0029") {
				return nil
			}
			return fmt.Errorf("rest API %s Post failed: %s", action, reason)
		}
		return nil
	}

	return c.PostAPI("edit_firenet(lan_ping=false)", data, customCheck)
}

func (c *Client) EnableTgwSegmentationForEgress(net *FireNet) error {
	data := map[string]string{
		"action": "enable_firenet_tgw_segmentation_for_egress",
		"CID":    c.CID,
		"vpc_id": net.VpcID,
	}

	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) DisableTgwSegmentationForEgress(net *FireNet) error {
	data := map[string]string{
		"action": "disable_firenet_tgw_segmentation_for_egress",
		"CID":    c.CID,
		"vpc_id": net.VpcID,
	}

	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) EditFirenetEgressStaticCidr(net *FireNet) error {
	data := map[string]string{
		"action":             "edit_firenet_egress_static_cidr",
		"CID":                c.CID,
		"vpc_id":             net.VpcID,
		"egress_static_cidr": net.EgressStaticCidrs,
	}

	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) EditFirenetExcludedCidr(net *FireNet) error {
	form := map[string]string{
		"action":       "edit_firenet_excluded_cidr",
		"CID":          c.CID,
		"vpc_id":       net.VpcID,
		"exclude_cidr": net.ExcludedCidrs,
	}
	check := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "list is not changed") {
				return nil
			}
			return fmt.Errorf("rest API edit_firenet_excluded_cidr Post failed: %s", reason)
		}
		return nil
	}
	return c.PostAPI(form["action"], form, check)
}

func (c *Client) EnableFirenetFailClose(net *FireNet) error {
	form := map[string]string{
		"action": "enable_firenet_fail_close",
		"CID":    c.CID,
		"vpc_id": net.VpcID,
	}
	check := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "configuration not changed") {
				return nil
			}
			return fmt.Errorf("rest API enable_firenet_fail_close Post failed: %s", reason)
		}
		return nil
	}
	return c.PostAPI(form["action"], form, check)
}

func (c *Client) DisableFirenetFailClose(net *FireNet) error {
	form := map[string]string{
		"action": "disable_firenet_fail_close",
		"CID":    c.CID,
		"vpc_id": net.VpcID,
	}
	check := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "configuration not changed") {
				return nil
			}
			return fmt.Errorf("rest API disable_firenet_fail_close Post failed: %s", reason)
		}
		return nil
	}
	return c.PostAPI(form["action"], form, check)
}
