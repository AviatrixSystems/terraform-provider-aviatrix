package aviatrix

import (
	"fmt"
	"github.com/AviatrixSystems/go-aviatrix/goaviatrix"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
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
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_region": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"vpc_account_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"vpc_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
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

	log.Printf("[INFO] Creating AWS TGW")

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

		for _, attachedVPCs := range dn["attached_vpc"].([]interface{}) {

			attachedVPC := attachedVPCs.(map[string]interface{})

			if dn["security_domain_name"].(string) == "Aviatrix_Edge_Domain" && attachedVPC != nil {
				return fmt.Errorf("validation of source file failed: no VPCs should be attached to 'Aviatrix_Edge_Domain'")
			}

			vpcSolo := goaviatrix.VPCSolo{
				Region:      attachedVPC["vpc_region"].(string),
				AccountName: attachedVPC["vpc_account_name"].(string),
				VpcID:       attachedVPC["vpc_id"].(string),
			}

			if vpcSolo.Region == "" {
				return fmt.Errorf("validation of source file failed: region of VPC (ID: %v) is not given",
					vpcSolo.VpcID)
			} else if vpcSolo.Region != awsTgw.Region {
				return fmt.Errorf("validation of source file failed: region of VPC (ID: %v) is different than "+
					"AWS_TGW", vpcSolo.VpcID)
			}

			if vpcSolo.AccountName == "" {
				return fmt.Errorf("validation of source file failed: account of VPC (ID: %v) is not given",
					vpcSolo.VpcID)
			}

			securityDomainRule.AttachedVPCs = append(securityDomainRule.AttachedVPCs, vpcSolo)

			temp := []string{dn["security_domain_name"].(string), attachedVPC["vpc_id"].(string),
				attachedVPC["vpc_account_name"].(string), attachedVPC["vpc_region"].(string)}

			attachedVPCAll = append(attachedVPCAll, temp)
		}

		awsTgw.SecurityDomains = append(awsTgw.SecurityDomains, securityDomainRule)
	}

	attachedGWs := d.Get("attached_aviatrix_transit_gateway").([]interface{})
	for _, attachedGW := range attachedGWs {
		attachedGWAll = append(attachedGWAll, attachedGW.(string))
		awsTgw.AttachedAviatrixTransitGW = append(awsTgw.AttachedAviatrixTransitGW, attachedGW.(string))
	}

	mAttachedGW := make(map[string]int)
	for i := 1; i <= len(attachedGWAll); i++ {
		if mAttachedGW[attachedGWAll[i-1]] != 0 {
			return fmt.Errorf("validation of source file failed: duplicate transit gateways (ID: %v) to attach",
				attachedGWAll[i-1])
		}
		mAttachedGW[attachedGWAll[i-1]] = i
	}

	domainsToCreate, domainConnPolicy, domainConnRemove, err := client.ValidateAWSTgwDomains(domainsAll, domainConnAll,
		attachedVPCAll)
	if err != nil {
		return fmt.Errorf("validation of source file failed: %v", err)
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

	for i := range attachedGWAll {
		gateway := &goaviatrix.Gateway{
			GwName: attachedGWAll[i],
		}
		err := client.AttachAviatrixTransitGWToAWSTgw(awsTgw, gateway, "Aviatrix_Edge_Domain")
		if err != nil {
			resourceAWSTgwRead(d, meta)
			return fmt.Errorf("failed to attach transit GW: %s", err)
		}
	}

	for i := range attachedVPCAll {
		if len(attachedVPCAll[i]) == 4 {
			vpcSolo := goaviatrix.VPCSolo{
				Region:      attachedVPCAll[i][3],
				AccountName: attachedVPCAll[i][2],
				VpcID:       attachedVPCAll[i][1],
			}

			err := client.AttachVpcToAWSTgw(awsTgw, vpcSolo, attachedVPCAll[i][0])
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

	log.Printf("[INFO] Reading AWS TGW")

	awsTgw, err := client.GetAWSTgw(awsTgw)
	if err != nil {
		return fmt.Errorf("couldn't find AWS TGW: %s", awsTgw.Name)
	}

	d.Set("account_name", awsTgw.AccountName)
	d.Set("tgw_name", awsTgw.Name)
	d.Set("region", awsTgw.Region)
	d.Set("aws_side_as_number", awsTgw.AwsSideAsNumber)
	d.Set("attached_aviatrix_transit_gateway", awsTgw.AttachedAviatrixTransitGW)

	mSecurityDomain := make(map[string]map[string]interface{})

	for _, sd := range awsTgw.SecurityDomains {
		sdr := make(map[string]interface{})
		sdr["security_domain_name"] = sd.Name
		sdr["connected_domains"] = sd.ConnectedDomain

		var aVPCs []interface{}
		for _, attachedVPC := range sd.AttachedVPCs {
			vpcSolo := make(map[string]interface{})
			vpcSolo["vpc_region"] = attachedVPC.Region
			vpcSolo["vpc_account_name"] = attachedVPC.AccountName
			vpcSolo["vpc_id"] = attachedVPC.VpcID
			aVPCs = append(aVPCs, vpcSolo)
		}
		sdr["attached_vpc"] = aVPCs

		mSecurityDomain[sd.Name] = sdr
	}

	var securityDomains []map[string]interface{}
	domains := d.Get("security_domains").([]interface{})

	mOld := make(map[string]bool)

	for _, domain := range domains {
		dn := domain.(map[string]interface{})

		mOld[dn["security_domain_name"].(string)] = true

		if mSecurityDomain[dn["security_domain_name"].(string)] != nil {
			mADm := make(map[string]bool)
			aDmNew := make([]string, 0)
			attachedDomains := mSecurityDomain[dn["security_domain_name"].(string)]["connected_domains"].([]string)

			for i := 0; i < len(attachedDomains); i++ {
				mADm[attachedDomains[i]] = true
			}
			attachedDomains1 := dn["connected_domains"].([]interface{})

			for i := 0; i < len(attachedDomains1); i++ {
				if mADm[attachedDomains1[i].(string)] {
					aDmNew = append(aDmNew, attachedDomains1[i].(string))
					mADm[attachedDomains1[i].(string)] = false
				}
			}

			for i := 0; i < len(attachedDomains); i++ {
				if mADm[attachedDomains[i]] {
					aDmNew = append(aDmNew, attachedDomains[i])
				}
			}

			mSecurityDomain[dn["security_domain_name"].(string)]["connected_domains"] = aDmNew

			mVPC := make(map[string]bool)
			var aVPCNew []map[string]interface{}

			for _, attachedVPCs := range mSecurityDomain[dn["security_domain_name"].(string)]["attached_vpc"].([]interface{}) {
				attachedVPC := attachedVPCs.(map[string]interface{})
				mVPC[attachedVPC["vpc_id"].(string)] = true
			}

			for _, attachedVPCs := range dn["attached_vpc"].([]interface{}) {
				attachedVPC := attachedVPCs.(map[string]interface{})
				if mVPC[attachedVPC["vpc_id"].(string)] {
					aVPCNew = append(aVPCNew, attachedVPC)
					mVPC[attachedVPC["vpc_id"].(string)] = false
				}
			}

			for _, attachedVPCs := range mSecurityDomain[dn["security_domain_name"].(string)]["attached_vpc"].([]interface{}) {
				attachedVPC := attachedVPCs.(map[string]interface{})
				if mVPC[attachedVPC["vpc_id"].(string)] {
					aVPCNew = append(aVPCNew, attachedVPC)
				}
			}

			mSecurityDomain[dn["security_domain_name"].(string)]["attached_vpc"] = aVPCNew

			securityDomains = append(securityDomains, mSecurityDomain[dn["security_domain_name"].(string)])
		}
	}

	for _, dn := range awsTgw.SecurityDomains {
		if !mOld[dn.Name] {
			securityDomains = append(securityDomains, mSecurityDomain[dn.Name])
		}
	}

	d.Set("security_domains", securityDomains)

	return nil
}

func resourceAWSTgwUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating AWS TGW")

	client := meta.(*goaviatrix.Client)
	awsTgw := &goaviatrix.AWSTgw{
		Name:        d.Get("tgw_name").(string),
		AccountName: d.Get("account_name").(string),
		Region:      d.Get("region").(string),
	}

	var toAttachGWs []string
	var toDetachGWs []string
	var domainsToCreate []string
	var domainsToRemove []string
	var domainConnPolicy [][]string
	var domainConnRemove [][]string
	var toAttachVPCs [][]string
	var toDetachVPCs [][]string

	d.Partial(true)

	mAttachedGWNew := make(map[string]int)

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

		for i := 1; i <= len(newAGWList); i++ {
			if mAttachedGWNew[newAGWList[i-1]] != 0 {
				return fmt.Errorf("validation of source file failed: duplicate transit gateways (ID: %v) to attach", newAGWList[i-1])
			}
			mAttachedGWNew[newAGWList[i-1]] = i
		}

		toAttachGWs = goaviatrix.Difference(newAGWList, oldAGWList)
		toDetachGWs = goaviatrix.Difference(oldAGWList, newAGWList)
	}

	if d.HasChange("security_domains") {
		oldSD, newSD := d.GetChange("security_domains")
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

			for _, attachedVPCs := range dn["attached_vpc"].([]interface{}) {

				attachedVPC := attachedVPCs.(map[string]interface{})

				if dn["security_domain_name"].(string) == "Aviatrix_Edge_Domain" && attachedVPC != nil {
					return fmt.Errorf("validation of source file failed: no VPCs should be attached to 'Aviatrix_Edge_Domain'")
				}

				vpcSolo := goaviatrix.VPCSolo{
					Region:      attachedVPC["vpc_region"].(string),
					AccountName: attachedVPC["vpc_account_name"].(string),
					VpcID:       attachedVPC["vpc_id"].(string),
				}
				securityDomainRule.AttachedVPCs = append(securityDomainRule.AttachedVPCs, vpcSolo)

				temp := []string{dn["security_domain_name"].(string), attachedVPC["vpc_id"].(string),
					attachedVPC["vpc_account_name"].(string), attachedVPC["vpc_region"].(string)}

				attachedVPCOld = append(attachedVPCOld, temp)
			}

		}

		domainsToCreateOld, domainConnPolicyOld, domainConnRemoveOld, _ := client.ValidateAWSTgwDomains(domainsOld,
			domainConnOld, attachedVPCOld)

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

			for _, attachedVPCs := range dn["attached_vpc"].([]interface{}) {

				attachedVPC := attachedVPCs.(map[string]interface{})

				if dn["security_domain_name"].(string) == "Aviatrix_Edge_Domain" && attachedVPC != nil {
					return fmt.Errorf("validation of source file failed: no VPCs should be attached to 'Aviatrix_Edge_Domain'")
				}

				vpcSolo := goaviatrix.VPCSolo{
					Region:      attachedVPC["vpc_region"].(string),
					AccountName: attachedVPC["vpc_account_name"].(string),
					VpcID:       attachedVPC["vpc_id"].(string),
				}

				if vpcSolo.Region == "" {
					return fmt.Errorf("validation of source file failed: region of VPC (ID: %v) is not given",
						vpcSolo.VpcID)
				} else if vpcSolo.Region != awsTgw.Region {
					return fmt.Errorf("validation of source file failed: region of VPC (ID: %v) is different than "+
						"AWS_TGW", vpcSolo.VpcID)
				}

				if vpcSolo.AccountName == "" {
					return fmt.Errorf("validation of source file failed: account of VPC (ID: %v) is not given",
						vpcSolo.VpcID)
				}

				securityDomainRule.AttachedVPCs = append(securityDomainRule.AttachedVPCs, vpcSolo)

				temp := []string{dn["security_domain_name"].(string), attachedVPC["vpc_id"].(string),
					attachedVPC["vpc_account_name"].(string), attachedVPC["vpc_region"].(string)}

				attachedVPCNew = append(attachedVPCNew, temp)
			}

		}

		domainsToCreateNew, domainConnPolicyNew, domainConnRemoveNew, err := client.ValidateAWSTgwDomains(domainsNew,
			domainConnNew, attachedVPCNew)
		if err != nil {
			return fmt.Errorf("validation of source file failed: %v", err)
		}

		domainsToCreate = goaviatrix.Difference(domainsToCreateNew, domainsToCreateOld)
		domainsToRemove = goaviatrix.Difference(domainsToCreateOld, domainsToCreateNew)

		domainConnPolicy = goaviatrix.DifferenceSlice(domainConnPolicyNew, domainConnPolicyOld)
		domainConnRemove = goaviatrix.DifferenceSlice(domainConnPolicyOld, domainConnPolicyNew)

		domainConnPolicy1 := goaviatrix.DifferenceSlice(domainConnRemoveOld, domainConnRemoveNew)
		domainConnRemove1 := goaviatrix.DifferenceSlice(domainConnRemoveNew, domainConnRemoveOld)

		toAttachVPCs = goaviatrix.DifferenceSlice(attachedVPCNew, attachedVPCOld)
		toDetachVPCs = goaviatrix.DifferenceSlice(attachedVPCOld, attachedVPCNew)

		if domainConnPolicy1 != nil || len(domainConnPolicy1) != 0 {
			for i := range domainConnPolicy1 {
				domainConnPolicy = append(domainConnPolicy, domainConnPolicy1[i])
			}
		}

		if domainConnRemove1 != nil || len(domainConnRemove1) != 0 {
			for i := range domainConnRemove1 {
				domainConnRemove = append(domainConnRemove, domainConnRemove1[i])
			}
		}
	}

	for i := range toDetachGWs {
		gateway := &goaviatrix.Gateway{
			GwName: toDetachGWs[i],
		}

		err := client.DetachAviatrixTransitGWToAWSTgw(awsTgw, gateway, "Aviatrix_Edge_Domain")
		if err != nil {
			resourceAWSTgwRead(d, meta)
			return fmt.Errorf("failed to detach transit GW: %s", err)
		}
	}

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

	for i := range toDetachVPCs {
		if len(toDetachVPCs[i]) == 4 {
			err := client.DetachVpcFromAWSTgw(awsTgw, toDetachVPCs[i][1])
			if err != nil {
				resourceAWSTgwRead(d, meta)
				return fmt.Errorf("failed to detach VPC: %s", err)
			}
		}
	}

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

	for i := range toAttachVPCs {
		if len(toAttachVPCs[i]) == 4 {
			vpcSolo := goaviatrix.VPCSolo{
				Region:      toAttachVPCs[i][3],
				AccountName: toAttachVPCs[i][2],
				VpcID:       toAttachVPCs[i][1],
			}

			err := client.AttachVpcToAWSTgw(awsTgw, vpcSolo, toAttachVPCs[i][0])
			if err != nil {
				resourceAWSTgwRead(d, meta)
				return fmt.Errorf("failed to attach VPC: %s", err)
			}
		}
	}

	for i := range domainsToRemove {
		securityDomain := &goaviatrix.SecurityDomain{
			Name:        domainsToRemove[i],
			AccountName: d.Get("account_name").(string),
			Region:      d.Get("region").(string),
			AwsTgwName:  d.Get("tgw_name").(string),
		}

		err := client.DeleteSecurityDomain(securityDomain)
		if err != nil {
			resourceAWSTgwRead(d, meta)
			return fmt.Errorf("failed to delete Security Domain: %s", err)
		}
	}

	d.Partial(false)

	return nil
}

func resourceAWSTgwDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	awsTgw := &goaviatrix.AWSTgw{
		Name:                      d.Get("tgw_name").(string),
		AccountName:               d.Get("account_name").(string),
		Region:                    d.Get("region").(string),
		AwsSideAsNumber:           d.Get("aws_side_as_number").(string),
		AttachedAviatrixTransitGW: make([]string, 0),
		SecurityDomains:           make([]goaviatrix.SecurityDomainRule, 0),
	}

	log.Printf("[INFO] Deleting AWS TGW")

	var attachedGWs []string
	var attachedVPCs [][]string

	domains := d.Get("security_domains").([]interface{})
	for _, domain := range domains {
		dn := domain.(map[string]interface{})

		for _, aVPCs := range dn["attached_vpc"].([]interface{}) {
			aVPC := aVPCs.(map[string]interface{})

			temp := []string{dn["security_domain_name"].(string), aVPC["vpc_id"].(string),
				aVPC["vpc_account_name"].(string), aVPC["vpc_region"].(string)}
			attachedVPCs = append(attachedVPCs, temp)
		}
	}

	for i := range attachedVPCs {
		if len(attachedVPCs[i]) == 4 {
			err := client.DetachVpcFromAWSTgw(awsTgw, attachedVPCs[i][1])
			if err != nil {
				resourceAWSTgwRead(d, meta)
				return fmt.Errorf("failed to detach VPC: %s", err)
			}
		}
	}

	transitGWs := d.Get("attached_aviatrix_transit_gateway").([]interface{})

	for _, transitGW := range transitGWs {
		attachedGWs = append(attachedGWs, transitGW.(string))
	}

	for i := range attachedGWs {
		gateway := &goaviatrix.Gateway{
			GwName: attachedGWs[i],
		}

		err := client.DetachAviatrixTransitGWToAWSTgw(awsTgw, gateway, "Aviatrix_Edge_Domain")
		if err != nil {
			resourceAWSTgwRead(d, meta)
			return fmt.Errorf("failed to detach transit GW: %s", err)
		}
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
