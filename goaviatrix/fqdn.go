package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Filters struct {
	FQDN     string `form:"fqdn,omitempty" json:"fqdn,omitempty"`
	Protocol string `form:"proto,omitempty" json:"proto,omitempty"`
	Port     string `form:"port,omitempty" json:"port,omitempty"`
	Verdict  string `form:"verdict,omitempty" json:"verdict,omitempty"`
}

// Gateway simple struct to hold fqdn details
type FQDN struct {
	FQDNTag         string        `form:"tag_name,omitempty" json:"tag_name,omitempty"`
	Action          string        `form:"action,omitempty"`
	CID             string        `form:"CID,omitempty"`
	FQDNStatus      string        `form:"status,omitempty" json:"status,omitempty"`
	FQDNMode        string        `form:"color,omitempty" json:"color,omitempty"`
	GwFilterTagList []GwFilterTag `form:"gw_name,omitempty" json:"gw_name,omitempty"`
	DomainList      []*Filters    `form:"domain_names[],omitempty" json:"domain_names,omitempty"`
}

type GwFilterTag struct {
	Name         string   `json:"gw_name,omitempty"`
	SourceIPList []string `json:"source_ip_list,omitempty"`
}

type ResultListResp struct {
	Return  bool     `json:"return"`
	Results []string `json:"results"`
	Reason  string   `json:"reason"`
}

type GetFqdnExceptionRuleResp struct {
	Return  bool   `json:"return"`
	Results string `json:"results"`
	Reason  string `json:"reason"`
}

type ResultListSourceIPResp struct {
	Return  bool       `json:"return"`
	Results GwSourceIP `json:"results"`
	Reason  string     `json:"reason"`
}

type GwSourceIP struct {
	ConfiguredIPs []string `json:"configured_ips"`
	VpcSubnets    []string `json:"vpc_subnets"`
}

type FQDNPassThroughResp struct {
	Return  bool            `json:"return"`
	Results FQDNPassThrough `json:"results"`
	Reason  string          `json:"reason"`
}

type FQDNPassThrough struct {
	ConfiguredIPs []string `json:"configured_ips"`
}

func (c *Client) CreateFQDN(fqdn *FQDN) error {
	form := map[string]string{
		"CID":      c.CID,
		"action":   "add_fqdn_filter_tag",
		"tag_name": fqdn.FQDNTag,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) DeleteFQDN(fqdn *FQDN) error {
	form := map[string]string{
		"CID":      c.CID,
		"action":   "del_fqdn_filter_tag",
		"tag_name": fqdn.FQDNTag,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

//change state to 'enabled' or 'disabled'
func (c *Client) UpdateFQDNStatus(fqdn *FQDN) error {
	form := map[string]string{
		"CID":      c.CID,
		"action":   "set_fqdn_filter_tag_state",
		"tag_name": fqdn.FQDNTag,
		"status":   fqdn.FQDNStatus,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

//Change default mode to 'white' or 'black'
func (c *Client) UpdateFQDNMode(fqdn *FQDN) error {
	form := map[string]string{
		"CID":      c.CID,
		"action":   "set_fqdn_filter_tag_color",
		"tag_name": fqdn.FQDNTag,
		"color":    fqdn.FQDNMode,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) UpdateDomains(fqdn *FQDN) error {
	// TODO: use PostAPI - domain names need special processing
	fqdn.CID = c.CID
	fqdn.Action = "set_fqdn_filter_tag_domain_names"
	log.Infof("Update domains: %#v", fqdn)

	verb := "POST"
	body := fmt.Sprintf("CID=%s&action=%s&tag_name=%s", c.CID, fqdn.Action, fqdn.FQDNTag)
	for i, dn := range fqdn.DomainList {
		body = body + fmt.Sprintf("&domain_names[%d][fqdn]=%s&domain_names[%d]"+
			"[proto]=%s&domain_names[%d][port]=%s&domain_names[%d][verdict]=%s", i, dn.FQDN, i, dn.Protocol, i, dn.Port, i, dn.Verdict)
	}
	log.Tracef("%s %s Body: %s", verb, c.baseURL, body)
	req, err := http.NewRequest(verb, c.baseURL, strings.NewReader(body))
	if err == nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		return errors.New("HTTP NewRequest set_fqdn_filter_tag_domain_names failed: " + err.Error())
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return errors.New("HTTP Post set_fqdn_filter_tag_domain_names failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode set_fqdn_filter_tag_domain_names failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API set_fqdn_filter_tag_domain_names Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) AttachGws(fqdn *FQDN) error {
	return nil
}

func (c *Client) DetachGws(fqdn *FQDN, gwList []string) error {
	form := map[string]string{
		"CID":      c.CID,
		"action":   "detach_fqdn_filter_tag_from_gw",
		"tag_name": fqdn.FQDNTag,
	}

	for i := range gwList {
		form["gw_name"] = gwList[i]
		err := c.PostAPI(form["action"], form, BasicCheck)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) ListFQDNTags() ([]*FQDN, error) {
	form := map[string]string{
		"CID":    c.CID,
		"action": "list_fqdn_filter_tags",
	}

	var data map[string]interface{}

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	tags := make([]*FQDN, 0)
	if val, ok := data["results"]; ok {
		for tag, data := range val.(map[string]interface{}) {
			tagData := data.(map[string]interface{})
			fqdn := &FQDN{
				FQDNTag:    tag,
				FQDNMode:   tagData["wbmode"].(string),
				FQDNStatus: tagData["state"].(string),
			}
			tags = append(tags, fqdn)
		}
	}
	return tags, nil
}

func (c *Client) GetFQDNTag(fqdn *FQDN) (*FQDN, error) {
	tags, err := c.ListFQDNTags()
	if err != nil {
		return nil, err
	}

	for _, tag := range tags {
		if tag.FQDNTag == fqdn.FQDNTag {
			fqdn.FQDNMode = tag.FQDNMode
			fqdn.FQDNStatus = tag.FQDNStatus
			return fqdn, nil
		}
	}
	log.Errorf("Couldn't find Aviatrix FQDN tag %s", fqdn.FQDNTag)
	return nil, ErrNotFound
}

func (c *Client) ListDomains(fqdn *FQDN) (*FQDN, error) {
	form := map[string]string{
		"CID":      c.CID,
		"action":   "list_fqdn_filter_tag_domain_names",
		"tag_name": fqdn.FQDNTag,
	}

	var data map[string]interface{}

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	names := data["results"].([]interface{})
	for _, domain := range names {
		dn := domain.(map[string]interface{})
		fqdnFilter := Filters{
			FQDN:     dn["fqdn"].(string),
			Protocol: dn["proto"].(string),
			Port:     dn["port"].(string),
			Verdict:  dn["verdict"].(string),
		}
		fqdn.DomainList = append(fqdn.DomainList, &fqdnFilter)
	}

	return fqdn, nil
}

func (c *Client) ListGws(fqdn *FQDN) ([]string, error) {
	form := map[string]string{
		"CID":      c.CID,
		"action":   "list_fqdn_filter_tag_attached_gws",
		"tag_name": fqdn.FQDNTag,
	}

	var data ResultListResp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	return data.Results, nil
}

func (c *Client) AttachTagToGw(fqdn *FQDN, gateway *Gateway) error {
	form := map[string]string{
		"CID":      c.CID,
		"action":   "attach_fqdn_filter_tag_to_gw",
		"tag_name": fqdn.FQDNTag,
		"gw_name":  gateway.GwName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) UpdateSourceIPFilters(fqdn *FQDN, gateway *Gateway, sourceIPs []string) error {
	form := map[string]string{
		"CID":      c.CID,
		"action":   "update_fqdn_filter_tag_source_ip_filters",
		"tag_name": fqdn.FQDNTag,
		"gw_name":  gateway.GwName,
	}

	if len(sourceIPs) != 0 {
		for i := range sourceIPs {
			form["source_ips["+strconv.Itoa(i)+"]"] = sourceIPs[i]
		}
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) GetGwFilterTagList(fqdn *FQDN) (*FQDN, error) {
	listGws, err := c.ListGws(fqdn)
	if err != nil {
		return nil, errors.New("failed for list_fqdn_filter_tag_source_ip_filters: " + err.Error())
	}

	var gwFilterTagList []GwFilterTag
	for i := range listGws {
		form := map[string]string{
			"CID":          c.CID,
			"action":       "list_fqdn_filter_tag_source_ip_filters",
			"tag_name":     fqdn.FQDNTag,
			"gateway_name": listGws[i],
		}

		var data ResultListSourceIPResp

		err = c.GetAPI(&data, form["action"], form, BasicCheck)
		if err != nil {
			return nil, err
		}

		var gwFilterTag GwFilterTag
		gwFilterTag.Name = listGws[i]
		sourceIPs := make([]string, 0)
		for j := range data.Results.ConfiguredIPs {
			sourceIPs = append(sourceIPs, strings.Split(data.Results.ConfiguredIPs[j], "~~")[0])
		}
		gwFilterTag.SourceIPList = sourceIPs
		gwFilterTagList = append(gwFilterTagList, gwFilterTag)
	}

	fqdn.GwFilterTagList = gwFilterTagList
	return fqdn, nil
}

func (c *Client) GetFQDNPassThroughCIDRs(gw *Gateway) ([]string, error) {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "list_fqdn_pass_through_cidrs",
		"gateway_name": gw.GwName,
	}

	var data FQDNPassThroughResp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	if len(data.Results.ConfiguredIPs) < 1 {
		return nil, ErrNotFound
	}

	return data.Results.ConfiguredIPs, nil
}

func (c *Client) ConfigureFQDNPassThroughCIDRs(gw *Gateway, IPs []string) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "update_fqdn_pass_through_cidrs",
		"gateway_name": gw.GwName,
	}
	for i, ip := range IPs {
		key := fmt.Sprintf("source_cidrs[%d]", i)
		form[key] = ip
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) DisableFQDNPassThrough(gw *Gateway) error {
	form := map[string]string{
		"CID":          c.CID,
		"action":       "update_fqdn_pass_through_cidrs",
		"gateway_name": gw.GwName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) AddFQDNTagRule(fqdn *FQDN) error {
	policies, err := json.Marshal(fqdn.DomainList)
	if err != nil {
		return fmt.Errorf("could not marshal fqdn domain: %v", err)
	}

	form := map[string]string{
		"CID":      c.CID,
		"action":   "add_fqdn_policies_to_tag",
		"tag_name": fqdn.FQDNTag,
		"policies": string(policies),
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) GetFQDNTagRule(fqdn *FQDN) (*FQDN, error) {
	foundFQDN, err := c.ListDomains(&FQDN{FQDNTag: fqdn.FQDNTag})
	if err != nil {
		return nil, fmt.Errorf("could not list fqdn domains: %v", err)
	}
	domain := fqdn.DomainList[0]
	found := false
	for _, d := range foundFQDN.DomainList {
		if d.FQDN == domain.FQDN && d.Protocol == domain.Protocol && d.Port == domain.Port && d.Verdict == domain.Verdict {
			found = true
		}
	}
	if !found {
		return nil, ErrNotFound
	}

	return fqdn, nil
}

func (c *Client) DeleteFQDNTagRule(fqdn *FQDN) error {
	policies, err := json.Marshal(fqdn.DomainList)
	if err != nil {
		return fmt.Errorf("could not marshal fqdn domain: %v", err)
	}

	form := map[string]string{
		"CID":      c.CID,
		"action":   "delete_fqdn_policies_to_tag",
		"tag_name": fqdn.FQDNTag,
		"policies": string(policies),
	}

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "the following rules were not found") {
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	return c.PostAPI(form["action"], form, checkFunc)
}
