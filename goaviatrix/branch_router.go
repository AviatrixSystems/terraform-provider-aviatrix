package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
)

// BranchRouter represents a branch router used in CloudWAN
type BranchRouter struct {
	Action             string                     `form:"action,omitempty" map:"action" json:"-"`
	CID                string                     `form:"CID,omitempty" map:"CID" json:"-"`
	Name               string                     `form:"branch_name,omitempty" map:"branch_name" json:"rgw_name"`
	PublicIP           string                     `form:"public_ip,omitempty" map:"public_ip" json:"hostname"`
	Username           string                     `form:"username,omitempty" map:"username" json:"username"`
	KeyFile            string                     `form:"-" map:"-" json:"-"`
	Password           string                     `form:"password,omitempty" map:"password" json:"-"`
	HostOS             string                     `form:"host_os,omitempty" map:"host_os" json:"host_os"`
	SshPort            int                        `form:"-" map:"-" json:"ssh_port"`
	SshPortStr         string                     `form:"port,omitempty" map:"port" json:"-"`
	Address1           string                     `form:"addr_1,omitempty" map:"addr_1" json:"-"`
	Address2           string                     `form:"addr_2,omitempty" map:"addr_2" json:"-"`
	City               string                     `form:"city,omitempty" map:"city" json:"-"`
	State              string                     `form:"state,omitempty" map:"state" json:"-"`
	Country            string                     `form:"country,omitempty" map:"country" json:"-"`
	ZipCode            string                     `form:"zipcode,omitempty" map:"zipcode" json:"-"`
	Description        string                     `form:"description,omitempty" map:"description" json:"description"`
	Address            GetBranchRouterRespAddress `form:"-" map:"-" json:"address"`
	CheckReason        string                     `form:"-" map:"-" json:"check_reason"`
	BranchState        string                     `form:"-" map:"-" json:"registered"`
	PrimaryInterface   string                     `form:"-" map:"-" json:"wan_if_primary"`
	PrimaryInterfaceIP string                     `form:"-" map:"-" json:"wan_if_primary_public_ip"`
	BackupInterface    string                     `form:"-" map:"-" json:"wan_if_backup"`
	BackupInterfaceIP  string                     `form:"-" map:"-" json:"wan_if_backup_public_ip"`
	ConnectionName     string                     `form:"-" map:"-" json:"conn_name"`
}

type GetBranchRouterRespAddress struct {
	Address1 string `json:"addr_1"`
	Address2 string `json:"addr_2"`
	City     string `json:"city"`
	State    string `json:"state"`
	Country  string `json:"country"`
	ZipCode  string `json:"zipcode"`
}

func (c *Client) CreateBranchRouter(br *BranchRouter) error {
	br.Action = "register_cloudwan_branch"
	br.CID = c.CID
	files := []File{
		{
			Path:      br.KeyFile,
			ParamName: "private_key_file",
		},
	}
	resp, err := c.PostFile(c.baseURL, br.toMap(), files)
	if err != nil {
		return errors.New("HTTP Post register_cloudwan_branch failed: " + err.Error())
	}

	type Resp struct {
		Return  bool   `json:"return,omitempty"`
		Results string `json:"results,omitempty"`
		Reason  string `json:"reason,omitempty"`
	}
	var data Resp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return errors.New("Reading response body register_cloudwan_branch failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return errors.New("Json Decode register_cloudwan_branch failed: " + err.Error() + "\n Body: " + b.String())
	}
	if !data.Return {
		return errors.New("Rest API register_cloudwan_branch Post failed: " + data.Reason)
	}
	return nil
}

// toMap converts the struct to a map[string]string
// The 'map' tags on the struct tell us what the key name should be.
func (br *BranchRouter) toMap() map[string]string {
	out := make(map[string]string)
	v := reflect.ValueOf(br).Elem()
	tag := "map"
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := typ.Field(i)
		if tagv := fi.Tag.Get(tag); tagv != "" && tagv != "-" {
			out[tagv] = v.Field(i).String()
		}
	}
	return out
}

func (c *Client) GetBranchRouter(br *BranchRouter) (*BranchRouter, error) {
	resp, err := c.Post(c.baseURL, struct {
		CID    string `form:"CID"`
		Action string `form:"action"`
	}{
		CID:    c.CID,
		Action: "list_cloudwan_branches_summary",
	})
	if err != nil {
		return nil, errors.New("HTTP POST list_cloudwan_branches_summary failed: " + err.Error())
	}

	type Resp struct {
		Return  bool           `json:"return"`
		Results []BranchRouter `json:"results"`
		Reason  string         `json:"reason"`
	}
	var data Resp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return nil, errors.New("Reading response body list_cloudwan_branches_summary failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return nil, errors.New("Json Decode list_cloudwan_branches_summary failed: " + err.Error() +
			"\n Body: " + b.String())
	}
	if !data.Return {
		return nil, errors.New("Rest API list_cloudwan_branches_summary Post failed: " + data.Reason)
	}
	var foundBr *BranchRouter
	for _, branch := range data.Results {
		if branch.Name == br.Name {
			foundBr = &branch
			break
		}
	}
	if foundBr == nil {
		log.Errorf("Could not find Aviatrix branch router %s", br.Name)
		return nil, ErrNotFound
	}

	foundBr.Address1 = foundBr.Address.Address1
	foundBr.Address2 = foundBr.Address.Address2
	foundBr.City = foundBr.Address.City
	foundBr.State = foundBr.Address.State
	foundBr.Country = foundBr.Address.Country
	foundBr.ZipCode = foundBr.Address.ZipCode

	return foundBr, nil
}

func (c *Client) GetBranchRouterName(connName string) (string, error) {
	resp, err := c.Post(c.baseURL, struct {
		CID    string `form:"CID"`
		Action string `form:"action"`
	}{
		CID:    c.CID,
		Action: "list_cloudwan_branches_summary",
	})
	if err != nil {
		return "", errors.New("HTTP POST list_cloudwan_branches_summary failed: " + err.Error())
	}

	type Resp struct {
		Return  bool           `json:"return"`
		Results []BranchRouter `json:"results"`
		Reason  string         `json:"reason"`
	}
	var data Resp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return "", errors.New("Reading response body list_cloudwan_branches_summary failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return "", errors.New("Json Decode list_cloudwan_branches_summary failed: " + err.Error() +
			"\n Body: " + b.String())
	}
	if !data.Return {
		return "", errors.New("Rest API list_cloudwan_branches_summary Post failed: " + data.Reason)
	}

	for _, branch := range data.Results {
		if branch.ConnectionName == connName {
			return branch.Name, nil
		}
	}

	return "", ErrNotFound
}

func (c *Client) UpdateBranchRouter(br *BranchRouter) error {
	br.Action = "update_cloudwan_branch_info"
	br.CID = c.CID

	files := []File{
		{
			Path:      br.KeyFile,
			ParamName: "private_key_file",
		},
	}
	resp, err := c.PostFile(c.baseURL, br.toMap(), files)
	if err != nil {
		return errors.New("HTTP Post update_cloudwan_branch_info failed: " + err.Error())
	}

	type Resp struct {
		Return  bool   `json:"return,omitempty"`
		Results string `json:"results,omitempty"`
		Reason  string `json:"reason,omitempty"`
	}
	var data Resp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return errors.New("Reading response body update_cloudwan_branch_info failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return errors.New("Json Decode update_cloudwan_branch_info failed: " + err.Error() + "\n Body: " + b.String())
	}
	if !data.Return {
		return errors.New("Rest API update_cloudwan_branch_info Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DeleteBranchRouter(br *BranchRouter) error {
	resp, err := c.Post(c.baseURL, struct {
		CID    string `form:"CID"`
		Action string `form:"action"`
		Name   string `form:"branch_name"`
	}{
		CID:    c.CID,
		Action: "deregister_cloudwan_branch",
		Name:   br.Name,
	})
	if err != nil {
		return errors.New("HTTP POST deregister_cloudwan_branch failed: " + err.Error())
	}

	type Resp struct {
		Return  bool   `json:"return"`
		Results string `json:"results"`
		Reason  string `json:"reason"`
	}
	var data Resp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return errors.New("Reading response body deregister_cloudwan_branch failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return errors.New("Json Decode deregister_cloudwan_branch failed: " + err.Error() +
			"\n Body: " + b.String())
	}
	if !data.Return {
		return errors.New("Rest API deregister_cloudwan_branch Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) ConfigureBranchRouterInterfaces(br *BranchRouter) error {
	availableInterfaces, err := c.GetBranchRouterInterfaces(br)
	if err != nil {
		return err
	}

	if !Contains(availableInterfaces, br.PrimaryInterface) {
		return fmt.Errorf("branch router does not have the given primary interface '%s'. "+
			"Possible interfaces are [%s]", br.PrimaryInterface, strings.Join(availableInterfaces, ", "))
	}

	if br.BackupInterface != "" && !Contains(availableInterfaces, br.BackupInterface) {
		return fmt.Errorf("branch router does not have the given backup interface '%s'. "+
			"Possible interfaces are [%s]", br.BackupInterface, strings.Join(availableInterfaces, ", "))
	}

	resp, err := c.Post(c.baseURL, struct {
		CID                string `form:"CID"`
		Action             string `form:"action"`
		Name               string `form:"branch_name"`
		PrimaryInterface   string `form:"wan_primary_if"`
		PrimaryInterfaceIP string `form:"wan_primary_ip"`
		BackupInterface    string `form:"wan_backup_if"`
		BackupInterfaceIP  string `form:"wan_backup_ip"`
	}{
		CID:                c.CID,
		Action:             "config_cloudwan_branch_wan_interfaces",
		Name:               br.Name,
		PrimaryInterface:   br.PrimaryInterface,
		PrimaryInterfaceIP: br.PrimaryInterfaceIP,
		BackupInterface:    br.BackupInterface,
		BackupInterfaceIP:  br.BackupInterfaceIP,
	})
	if err != nil {
		return errors.New("HTTP POST config_cloudwan_branch_wan_interfaces failed: " + err.Error())
	}

	type Resp struct {
		Return  bool   `json:"return"`
		Results string `json:"results"`
		Reason  string `json:"reason"`
	}
	var data Resp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return errors.New("Reading response body config_cloudwan_branch_wan_interfaces failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return errors.New("Json Decode config_cloudwan_branch_wan_interfaces failed: " + err.Error() +
			"\n Body: " + b.String())
	}
	if !data.Return {
		return errors.New("Rest API config_cloudwan_branch_wan_interfaces Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetBranchRouterInterfaces(br *BranchRouter) ([]string, error) {
	resp, err := c.Post(c.baseURL, struct {
		CID    string `form:"CID"`
		Action string `form:"action"`
		Name   string `form:"branch_name"`
	}{
		CID:    c.CID,
		Action: "get_cloudwan_branch_wan_interfaces",
		Name:   br.Name,
	})
	if err != nil {
		return nil, errors.New("HTTP POST get_cloudwan_branch_wan_interfaces failed: " + err.Error())
	}

	type Resp struct {
		Return  bool                   `json:"return"`
		Results map[string]interface{} `json:"results"`
		Reason  string                 `json:"reason"`
	}
	var data Resp
	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return nil, errors.New("Reading response body get_cloudwan_branch_wan_interfaces failed: " + err.Error())
	}

	if err = json.NewDecoder(&b).Decode(&data); err != nil {
		return nil, errors.New("Json Decode get_cloudwan_branch_wan_interfaces failed: " + err.Error() +
			"\n Body: " + b.String())
	}
	if !data.Return {
		return nil, errors.New("Rest API get_cloudwan_branch_wan_interfaces Post failed: " + data.Reason)
	}
	var interfaces []string

	for k, _ := range data.Results {
		interfaces = append(interfaces, k)
	}

	return interfaces, nil
}
