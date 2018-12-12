package goaviatrix

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AwsTGW simple struct to hold aws_tgw details
type AWSTgw struct {
	Action          string  				`form:"action,omitempty"`
	CID             string  				`form:"CID,omitempty"`
	Name 		    string  				`form:"tgw_name,omitempty"`
	AccountName     string  				`form:"account_name,omitempty"`
	Region          string  				`form:"region,omitempty"`
	AwsSideAsNumber string                  `form:"aws_side_asn,omitempty"`
	SecurityDomains []SecurityDomainRule    `form:"security_domains,omitempty"`
}

type AWSTgwAPIResp struct {
	Return  bool        `json:"return"`
	Results []string    `json:"results"`
	Reason  string      `json:"reason"`
}

type AWSTgwList struct {
	Return  bool   	    `json:"return"`
	Results []AWSTgw    `json:"results"`
	Reason  string      `json:"reason"`
}

func (c *Client) CreateAWSTgw(awsTgw *AWSTgw) (error) {
	awsTgw.CID = c.CID
	awsTgw.Action = "add_aws_tgw"
	resp, err := c.Post(c.baseURL, awsTgw)
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

func (c *Client) GetAWSTgw(awsTgw *AWSTgw) (string, error) {
	awsTgw.CID = c.CID
	awsTgw.Action = "list_tgw_names"
	path := c.baseURL + fmt.Sprintf("?action=%s&CID=%s&account_name=%s&region=%s", awsTgw.Action, awsTgw.CID,
		awsTgw.AccountName, awsTgw.Region)
	resp, err := c.Get(path, nil)
	if err != nil {
		return "", err
	}

	data := AWSTgwAPIResp{
		Return:  false,
		Results: make([]string, 0),
		Reason:  "",
	}
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	if !data.Return {
		return "", errors.New(data.Reason)
	}

	awsTgwList := data.Results
	for i := range awsTgwList {
		if awsTgwList[i] == awsTgw.Name {
			return awsTgwList[i], nil
		}
	}

	return "", ErrNotFound
}

func (c *Client) UpdateAWSTgw(awsTgw *AWSTgw) (error) {
	return nil
}

func (c *Client) DeleteAWSTgw(awsTgw *AWSTgw) (error) {
	awsTgw.CID = c.CID
	awsTgw.Action = "delete_aws_tgw"
	resp, err := c.Post(c.baseURL, awsTgw)
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

func (c *Client) ValidateAWSTgwDomains(domainsAll []string, domainConnAll [][]string) ([]string, [][]string,
	[][]string, error) {
	numOfDomains := len(domainsAll)
	matrix := make([][]int, numOfDomains)
	var domainsToCreate []string
	var domainConnPolicy [][]string
	var domainConnRemove [][]string

	for i := range matrix {
		matrix[i] = make([]int, numOfDomains)
	}

	m := make(map[string]int)
	for i := 1; i <= numOfDomains; i++ {
		m[domainsAll[i - 1]] = i
	}

	for i := range domainConnAll {
		x := m[domainConnAll[i][0]]
		y := m[domainConnAll[i][1]]
		if x == 0 || y == 0 || x == y {
			return domainsToCreate, domainConnPolicy, domainConnRemove, ErrNotFound
		}
		matrix[x - 1][y - 1] = 1
	}

	for i := 0; i < numOfDomains; i++ {
		for j := i + 1; j < numOfDomains; j++ {
			if matrix[i][j] != matrix[j][i] {
				return domainsToCreate, domainConnPolicy, domainConnRemove, ErrNotFound
			}
		}
	}

	defaultX := [3]string{"Default_Domain", "Shared_Service_Domain", "Aviatrix_Edge_Domain"}

	for i := 0; i < 3; i++ {
		for j := i; j < 3; j++ {
			if i != j {

				if matrix[m[defaultX[i]] - 1][m[defaultX[j]] - 1] == 0 {
					temp := []string{defaultX[i], defaultX[j]}
					domainConnRemove = append(domainConnRemove, temp)
				}
				matrix[m[defaultX[i]] - 1][m[defaultX[j]] - 1] = 2
				matrix[m[defaultX[j]] - 1][m[defaultX[i]] - 1] = 2
			}
		}
	}

	for i := range domainConnAll {
		if matrix[m[domainConnAll[i][0]] - 1][m[domainConnAll[i][1]] - 1] == 1 {
			matrix[m[domainConnAll[i][0]] - 1][m[domainConnAll[i][1]] - 1] = 2
			matrix[m[domainConnAll[i][1]] - 1][m[domainConnAll[i][0]] - 1] = 2
			temp := []string{domainConnAll[i][0], domainConnAll[i][1]}
			domainConnPolicy = append(domainConnPolicy, temp)
		}
	}

	for i := range domainsAll {
		if domainsAll[i] != "Default_Domain" &&
			domainsAll[i] != "Shared_Service_Domain" &&
			domainsAll[i] != "Aviatrix_Edge_Domain" {
			domainsToCreate = append(domainsToCreate, domainsAll[i])
		}
	}

	return domainsToCreate, domainConnPolicy, domainConnRemove, nil
}

func (c *Client) AttachVpcToAWSTgw(awsTgw *AWSTgw) (error) {
	return nil
}

func (c *Client) DetachVpcFromAWSTgw(awsTgw *AWSTgw) (error) {
	return nil
}

func (c *Client) AttachAviatrixTgwToAWSTgw(awsTgw *AWSTgw) (error) {
	return nil
}

func (c *Client) DetachAviatrixTgwFromAWSTgw(awsTgw *AWSTgw) (error) {
	return nil
}