package aviatrix

import (
	"fmt"
	"log"

	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAWSTgw() *schema.Resource {
	return &schema.Resource{
		Create: resourceAWSTgwCreate,
		Read:   resourceAWSTgwRead,
		Update: resourceAWSTgwUpdate,
		Delete: resourceAWSTgwDelete,

		Schema: map[string]*schema.Schema{
			"aws_tgw_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aws_side_as_number": {
				Type:     schema.TypeString,
				Required: true,
			},
			"security_domains": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_domain_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"connected_domains": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
						"attached_vpc": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAWSTgwCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	awsTgw := &goaviatrix.AWSTgw{
		Name:            d.Get("aws_tgw_name").(string),
		AccountName:     d.Get("account_name").(string),
		Region:          d.Get("region").(string),
		AwsSideAsNumber: d.Get("aws_side_as_number").(string),
		SecurityDomains: make([]goaviatrix.SecurityDomainRule, 0),
	}
	domains := d.Get("security_domains").([]interface{})

	log.Printf("[INFO] zjin: 001")

	var domainsAll []string
	var domainConnAll [][]string

	for _, domain := range domains {
		dn := domain.(map[string]interface{})

		domainsAll = append(domainsAll, dn["security_domain_name"].(string))

		securityDomainRule := goaviatrix.SecurityDomainRule{
			Name: dn["security_domain_name"].(string),
		}
		for _, connectedDomain := range dn["connected_domains"].([]interface{}) {
			securityDomainRule.ConnectedDomain = append(securityDomainRule.ConnectedDomain, connectedDomain.(string))
			temp := []string{dn["security_domain_name"].(string), connectedDomain.(string)}
			domainConnAll = append(domainConnAll, temp)
		}
		for _, attachedVPC := range dn["attached_vpc"].([]interface{}) {
			securityDomainRule.AttachedVPCs = append(securityDomainRule.AttachedVPCs, attachedVPC.(string))
		}
		awsTgw.SecurityDomains = append(awsTgw.SecurityDomains, securityDomainRule)
	}

	domainsToCreate, domainConnPolicy, domainConnRemove, err := client.ValidateAWSTgwDomains(domainsAll, domainConnAll)
	if err != nil {
		return fmt.Errorf("validation of source file failed")
	}

	err1 := client.CreateAWSTgw(awsTgw)
	if err1 != nil {
		return fmt.Errorf("failed to create AWS TGW: %s", err1)
	}

	d.SetId(awsTgw.Name)

	for i := range domainsToCreate {
		securityDomain := &goaviatrix.SecurityDomain{
			Name:        domainsToCreate[i],
			AccountName: d.Get("account_name").(string),
			Region:      d.Get("region").(string),
			AwsTgwName:  d.Get("aws_tgw_name").(string),
		}

		err := client.CreateSecurityDomain(securityDomain)
		if err != nil {
			return fmt.Errorf("failed to create Security Domain: %s", err)
		}
	}

	for i := range domainConnPolicy {
		if len(domainConnPolicy[i]) == 2 {
			err := client.CreateDomainConnectionPolicy(awsTgw, domainConnPolicy[i][0], domainConnPolicy[i][1])
			if err != nil {
				return fmt.Errorf("failed to create security domain connection: %s", err)
			}
		}
	}

	for i := range domainConnRemove {
		if len(domainConnRemove[i]) == 2 {
			err := client.DeleteDomainConnectionPolicy(awsTgw, domainConnRemove[i][0], domainConnRemove[i][1])
			if err != nil {
				return fmt.Errorf("failed to delete domain connection: %s", err)
			}
		}
	}
	
	return resourceAWSTgwRead(d, meta)
}

func resourceAWSTgwRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	awsTgw := &goaviatrix.AWSTgw{
		Name:            d.Get("aws_tgw_name").(string),
		AccountName:     d.Get("account_name").(string),
		Region:          d.Get("region").(string),
		AwsSideAsNumber: d.Get("aws_side_as_number").(string),
	}

	awsTgwResp, err := client.GetAWSTgw(awsTgw)
	if err != nil {
		return fmt.Errorf("couldn't find AWS TGW: %s", err)
	}
	log.Printf("[TRACE] reading AWS TGW %s: %#v", d.Get("aws_tgw_name").(string), awsTgwResp)

	d.Set("account_name", awsTgw.AccountName)
	d.Set("aws_tgw_name", awsTgw.Name)
	d.Set("region", awsTgw.Region)
	d.Set("aws_side_as_number", awsTgw.AwsSideAsNumber)

	return nil
}

func resourceAWSTgwUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAWSTgwDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	awsTgw := &goaviatrix.AWSTgw{
		AccountName: d.Get("account_name").(string),
		Name:        d.Get("aws_tgw_name").(string),
		Region:      d.Get("region").(string),
	}

	err := client.DeleteAWSTgw(awsTgw)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find AWS TGW: %s", err)
	}

	return nil
}
