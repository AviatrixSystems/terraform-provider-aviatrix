package aviatrix

import (
	"fmt"
	"log"
	//"strings"

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
			"tgw_name": {
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
			"attached_aviatrix_transit_gateway": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
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
		Name:                      d.Get("tgw_name").(string),
		AccountName:               d.Get("account_name").(string),
		Region:                    d.Get("region").(string),
		AwsSideAsNumber:           d.Get("aws_side_as_number").(string),
		AttachedAviatrixTransitGW: make([]string, 0),
		SecurityDomains:           make([]goaviatrix.SecurityDomainRule, 0),
	}

	var domainsAll []string
	var domainConnAll [][]string
	var attachedGWAll []string
	var attachedVPCAll [][]string

	domains := d.Get("security_domains").([]interface{})
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
			temp := []string{dn["security_domain_name"].(string), attachedVPC.(string)}
			attachedVPCAll = append(attachedVPCAll, temp)
		}
		awsTgw.SecurityDomains = append(awsTgw.SecurityDomains, securityDomainRule)
	}

	attachedGWs := d.Get("attached_aviatrix_transit_gateway").([]interface{})
	for _, attachedGW := range attachedGWs {
		attachedGWAll = append(attachedGWAll, attachedGW.(string))
		awsTgw.AttachedAviatrixTransitGW = append(awsTgw.AttachedAviatrixTransitGW, attachedGW.(string))
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
			AwsTgwName:  d.Get("tgw_name").(string),
		}

		err := client.CreateSecurityDomain(securityDomain)
		if err != nil {
			resourceAWSTgwRead(d, meta)
			return fmt.Errorf("failed to create Security Domain: %s", err)
		}
	}

	for i := range domainConnPolicy {
		if len(domainConnPolicy[i]) == 2 {
			err := client.CreateDomainConnection(awsTgw, domainConnPolicy[i][0], domainConnPolicy[i][1])
			if err != nil {
				resourceAWSTgwRead(d, meta)
				return fmt.Errorf("failed to create security domain connection: %s", err)
			}
		}
	}

	for i := range domainConnRemove {
		if len(domainConnRemove[i]) == 2 {
			err := client.DeleteDomainConnection(awsTgw, domainConnRemove[i][0], domainConnRemove[i][1])
			if err != nil {
				resourceAWSTgwRead(d, meta)
				return fmt.Errorf("failed to delete domain connection: %s", err)
			}
		}
	}

	if attachedGWAll != nil {
		for i := range awsTgw.AttachedAviatrixTransitGW {
			gateway := &goaviatrix.Gateway{
				GwName: awsTgw.AttachedAviatrixTransitGW[i],
			}
			err := client.AttachAviatrixTransitGWToAWSTgw(awsTgw, gateway, "Aviatrix_Edge_Domain")
			if err != nil {
				resourceAWSTgwRead(d, meta)
				return fmt.Errorf("failed to attach transit GW: %s", err)
			}
		}
	}

	if len(attachedVPCAll) > 0 {
		for i := range attachedVPCAll {
			gateway := &goaviatrix.Gateway{
				VpcID: attachedVPCAll[i][1],
			}

			err := client.AttachAviatrixTransitGWToAWSTgw(awsTgw, gateway, attachedVPCAll[i][0])
			if err != nil {
				resourceAWSTgwRead(d, meta)
				return fmt.Errorf("failed to attach VPC: %s", err)
			}
		}
	}

	return resourceAWSTgwRead(d, meta)
}

func resourceAWSTgwRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	awsTgw := &goaviatrix.AWSTgw{
		Name:            d.Get("tgw_name").(string),
		AccountName:     d.Get("account_name").(string),
		Region:          d.Get("region").(string),
		AwsSideAsNumber: d.Get("aws_side_as_number").(string),
	}

	awsTgw, err := client.GetAWSTgw(awsTgw)
	if err != nil {
		return fmt.Errorf("couldn't find AWS TGW: %s", err)
	}

	d.Set("account_name", awsTgw.AccountName)
	d.Set("tgw_name", awsTgw.Name)
	d.Set("region", awsTgw.Region)
	d.Set("aws_side_as_number", awsTgw.AwsSideAsNumber)
	d.Set("attached_aviatrix_transit_gateway", awsTgw.AttachedAviatrixTransitGW)

	var securityDomains []map[string]interface{}
	for _, sd := range awsTgw.SecurityDomains {
		sdr := make(map[string]interface{})
		sdr["security_domain_name"] = sd.Name
		sdr["connected_domains"] = sd.ConnectedDomain
		sdr["attached_vpc"] = sd.AttachedVPCs
		securityDomains = append(securityDomains, sdr)
	}
	d.Set("security_domains", securityDomains)

	return nil
}

func resourceAWSTgwUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] zjin 001: Updating AWS TGW")

	client := meta.(*goaviatrix.Client)
	awsTgw := &goaviatrix.AWSTgw{
		Name:        d.Get("tgw_name").(string),
		AccountName: d.Get("account_name").(string),
		Region:      d.Get("region").(string),
	}

	d.Partial(true)

	if d.HasChange("attached_aviatrix_transit_gateway") {
		oldAGW, newAGW := d.GetChange("attached_aviatrix_transit_gateway")

		if oldAGW == nil {
			oldAGW = new([]interface{})
		}
		if newAGW == nil {
			newAGW = new([]interface{})
		}

		oldString := oldAGW.([]interface{})
		newString := newAGW.([]interface{})
		oldAGWList := goaviatrix.ExpandStringList(oldString)
		newAGWList := goaviatrix.ExpandStringList(newString)
		toAttachGWs := goaviatrix.Difference(newAGWList, oldAGWList)
		toDetachGWs := goaviatrix.Difference(oldAGWList, newAGWList)

		for i := range toAttachGWs {
			gateway := &goaviatrix.Gateway{
				GwName: toAttachGWs[i],
			}

			err := client.AttachAviatrixTransitGWToAWSTgw(awsTgw, gateway, "Aviatrix_Edge_Domain")
			if err != nil {
				resourceAWSTgwRead(d, meta)
				return fmt.Errorf("failed to attach transit GW: %s", err)
			}
		}

		for i := range toDetachGWs {
			gateway := &goaviatrix.Gateway{
				VpcID: toDetachGWs[i],
			}

			err := client.DetachAviatrixTransitGWToAWSTgw(awsTgw, gateway, "Aviatrix_Edge_Domain")
			if err != nil {
				resourceAWSTgwRead(d, meta)
				return fmt.Errorf("failed to detach transit GW: %s", err)
			}
		}
	}

	if d.HasChange("security_domains") {
		log.Printf("[INFO] zjin 002: Security Domains")

		oldSD, newSD := d.GetChange("security_domains")

		log.Printf("[INFO] zjin 003: oldSD = %v", oldSD)
		log.Printf("[INFO] zjin 003: newSD = %v", newSD)

		if oldSD == nil {
			oldSD = new([]interface{})
		}
		if newSD == nil {
			newSD = new([]interface{})
		}

		var domainsOld []string
		var domainConnOld [][]string
		var attachedVPCOld [][]string

		for _, domain := range oldSD.([]interface{}) {
			dn := domain.(map[string]interface{})

			domainsOld = append(domainsOld, dn["security_domain_name"].(string))

			securityDomainRule := goaviatrix.SecurityDomainRule{
				Name: dn["security_domain_name"].(string),
			}
			for _, connectedDomain := range dn["connected_domains"].([]interface{}) {
				securityDomainRule.ConnectedDomain = append(securityDomainRule.ConnectedDomain, connectedDomain.(string))
				temp := []string{dn["security_domain_name"].(string), connectedDomain.(string)}
				domainConnOld = append(domainConnOld, temp)
			}
			for _, attachedVPC := range dn["attached_vpc"].([]interface{}) {
				securityDomainRule.AttachedVPCs = append(securityDomainRule.AttachedVPCs, attachedVPC.(string))
				temp := []string{dn["security_domain_name"].(string), attachedVPC.(string)}
				attachedVPCOld = append(attachedVPCOld, temp)
			}
		}

		domainsToCreateOld, domainConnPolicyOld, domainConnRemoveOld, err := client.ValidateAWSTgwDomains(domainsOld,
			domainConnOld)

		log.Printf("[INFO] zjin 004: domainsToCreateOld = %v", domainsToCreateOld)
		log.Printf("[INFO] zjin 004: domainConnPolicyOld = %v", domainConnPolicyOld)
		log.Printf("[INFO] zjin 004: domainConnRemoveOld = %v", domainConnRemoveOld)
		log.Printf("[INFO] zjin 004: err = %v", err)

		var domainsNew []string
		var domainConnNew [][]string
		var attachedVPCNew [][]string

		for _, domain := range newSD.([]interface{}) {
			dn := domain.(map[string]interface{})

			domainsNew = append(domainsNew, dn["security_domain_name"].(string))

			securityDomainRule := goaviatrix.SecurityDomainRule{
				Name: dn["security_domain_name"].(string),
			}
			for _, connectedDomain := range dn["connected_domains"].([]interface{}) {
				securityDomainRule.ConnectedDomain = append(securityDomainRule.ConnectedDomain, connectedDomain.(string))
				temp := []string{dn["security_domain_name"].(string), connectedDomain.(string)}
				domainConnNew = append(domainConnNew, temp)
			}
			for _, attachedVPC := range dn["attached_vpc"].([]interface{}) {
				securityDomainRule.AttachedVPCs = append(securityDomainRule.AttachedVPCs, attachedVPC.(string))
				temp := []string{dn["security_domain_name"].(string), attachedVPC.(string)}
				attachedVPCNew = append(attachedVPCNew, temp)
			}
		}

		domainsToCreateNew, domainConnPolicyNew, domainConnRemoveNew, err1 := client.ValidateAWSTgwDomains(domainsNew,
			domainConnNew)

		log.Printf("[INFO] zjin 005: domainsToCreateNew = %v", domainsToCreateNew)
		log.Printf("[INFO] zjin 005: domainConnPolicyNew = %v", domainConnPolicyNew)
		log.Printf("[INFO] zjin 005: domainConnRemoveNew = %v", domainConnRemoveNew)
		log.Printf("[INFO] zjin 005: err1 = %v", err1)

		domainsToCreate := goaviatrix.Difference(domainsToCreateNew, domainsToCreateOld)
		//domainsToRemove := goaviatrix.Difference(domainsToCreateOld, domainsToCreateNew)

		domainConnPolicy := goaviatrix.DifferenceSlice(domainConnPolicyNew, domainConnPolicyOld)
		domainConnRemove := goaviatrix.DifferenceSlice(domainConnRemoveNew, domainConnRemoveOld)

		log.Printf("[INFO] zjin 006: domainsToCreate = %v", domainsToCreate)
		log.Printf("[INFO] zjin 006: domainConnPolicy = %v", domainConnPolicy)
		log.Printf("[INFO] zjin 006: domainConnRemove = %v", domainConnRemove)

		log.Printf("[INFO] zjin 007: attachedVPCNew size = %v", len(attachedVPCNew))
		log.Printf("[INFO] zjin 007: attachedVPCOld size = %v", len(attachedVPCOld))

		VPCToCreate := goaviatrix.DifferenceSlice(attachedVPCNew, attachedVPCOld)
		log.Printf("[INFO] zjin 007: VPCToCreate = %v", VPCToCreate)

	}

	return nil
}

func resourceAWSTgwDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	awsTgw := &goaviatrix.AWSTgw{
		AccountName: d.Get("account_name").(string),
		Name:        d.Get("tgw_name").(string),
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
