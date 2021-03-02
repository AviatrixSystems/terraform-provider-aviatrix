package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixAWSTgw() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixAWSTgwCreate,
		Read:   resourceAviatrixAWSTgwRead,
		Update: resourceAviatrixAWSTgwUpdate,
		Delete: resourceAviatrixAWSTgwDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 2,
		MigrateState:  resourceAviatrixAWSTgwMigrateState,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceAviatrixAWSTgwResourceV1().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceAviatrixAWSTgwStateUpgradeV1,
				Version: 1,
			},
		},

		Schema: map[string]*schema.Schema{
			"tgw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the AWS TGW which is going to be created.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This parameter represents the name of a Cloud-Account in Aviatrix controller.",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Region of cloud provider.",
			},
			"aws_side_as_number": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "BGP Local ASN (Autonomous System Number), Integer between 1-4294967294.",
				ValidateFunc: goaviatrix.ValidateASN,
			},
			"security_domains": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Security Domains to create together with AWS TGW's creation.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_domain_name": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Name of the security domain created.",
							ValidateFunc: validation.StringDoesNotContainAny(":"),
						},
						"connected_domains": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "A list of domains connected to the domain.",
						},
						"attached_vpc": {
							Type:        schema.TypeList,
							Optional:    true,
							Deprecated:  "Please set `manage_vpc_attachment` to false, and use the standalone aviatrix_aws_tgw_vpc_attachment resource instead.",
							Description: "A list of VPCs attached to the domain.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_region": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Region of the vpc.",
									},
									"vpc_account_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The name of a Cloud-Account in Aviatrix controller associated with this VPC.",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "This parameter represents the ID of the VPC.",
									},
									"customized_routes": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "",
										Description: "Advanced option. Customized Spoke VPC Routes. It allows the admin to enter non-RFC1918 routes in the VPC route table targeting the TGW.",
									},
									"subnets": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "",
										Description: "Advanced option. VPC subnets separated by ',' to attach to the VPC. If left blank, Aviatrix Controller automatically selects a subnet representing each AZ for the VPC attachment.",
									},
									"route_tables": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "",
										Description: "Advanced option. Route tables separated by ',' to participate in TGW Orchestrator, i.e., learned routes will be propagated to these route tables.",
									},
									"customized_route_advertisement": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "",
										Description: "Advanced option. Customized route(s) to be advertised to other VPCs that are connected to the same TGW.",
									},
									"disable_local_route_propagation": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Advanced option. If set to true, it disables automatic route propagation of this VPC to other VPCs within the same security domain.",
									},
								},
							},
						},
						"aviatrix_firewall": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Set to true if the security domain is an aviatrix firewall domain.",
						},
						"native_egress": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Set to true if the security domain is a native egress domain.",
						},
						"native_firewall": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Set to true if the security domain is a native firewall domain.",
						},
					},
				},
			},
			"cloud_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				Description:  "Type of cloud service provider, requires an integer value. Supported for AWS (1) and AWS GOV (256). Default value: 1.",
				ValidateFunc: validation.IntInSlice([]int{goaviatrix.AWS, goaviatrix.AWSGOV}),
			},
			"attached_aviatrix_transit_gateway": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Deprecated:  "Please set `manage_transit_gateway_attachment` to false, and use the standalone aviatrix_aws_tgw_transit_gateway_attachment resource instead.",
				Description: "A list of Names of Aviatrix Transit Gateway to attach to one of the three default domains.",
			},
			"manage_vpc_attachment": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				Description: "This parameter is a switch used to determine whether or not to manage VPC attachments to " +
					"the TGW using the aviatrix_aws_tgw resource. If this is set to false, attachment of VPCs must be done " +
					"using the aviatrix_aws_tgw_vpc_attachment resource. Valid values: true, false. Default value: true.",
			},
			"manage_transit_gateway_attachment": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				Description: "This parameter is a switch used to determine whether or not to manage transit gateway " +
					"attachments to the TGW using the aviatrix_aws_tgw resource. If this is set to false, attachment of " +
					"transit gateways must be done using the aviatrix_aws_tgw_transit_gateway_attachment resource. " +
					"Valid values: true, false. Default value: true.",
			},
			"enable_multicast": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable Multicast.",
			},
			"cidrs": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
				Description: "TGW CIDRs.",
			},
		},
	}
}

func resourceAviatrixAWSTgwCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	awsTgw := &goaviatrix.AWSTgw{
		Name:                      d.Get("tgw_name").(string),
		AccountName:               d.Get("account_name").(string),
		Region:                    d.Get("region").(string),
		AwsSideAsNumber:           d.Get("aws_side_as_number").(string),
		CloudType:                 d.Get("cloud_type").(int),
		AttachedAviatrixTransitGW: make([]string, 0),
		SecurityDomains:           make([]goaviatrix.SecurityDomainRule, 0),
		EnableMulticast:           d.Get("enable_multicast").(bool),
	}

	manageVpcAttachment := d.Get("manage_vpc_attachment").(bool)
	manageTransitGwAttachment := d.Get("manage_transit_gateway_attachment").(bool)

	if awsTgw.Name == "" {
		return fmt.Errorf("tgw name can't be empty string")
	}
	if awsTgw.AccountName == "" {
		return fmt.Errorf("account name can't be empty string")
	}
	if awsTgw.Region == "" {
		return fmt.Errorf("tgw region can't be empty string")
	}
	if awsTgw.AwsSideAsNumber == "" {
		return fmt.Errorf("aws side number can't be empty string")
	}

	log.Printf("[INFO] Creating AWS TGW")

	var domainsAll []string
	var domainConnAll [][]string
	var attachedGWAll []string
	var attachedVPCAll [][]string

	mapSecurityDomainRule := make(map[string][3]bool)
	mapFireNetVpc := make(map[string]bool)

	domains := d.Get("security_domains").([]interface{})
	for _, domain := range domains {

		dn := domain.(map[string]interface{})
		domainsAll = append(domainsAll, dn["security_domain_name"].(string))

		securityDomainRule := goaviatrix.SecurityDomainRule{
			Name: dn["security_domain_name"].(string),
		}
		if dn["aviatrix_firewall"].(bool) {
			securityDomainRule.AviatrixFirewallDomain = true
			mapFireNetVpc[securityDomainRule.Name] = true
		}
		if dn["native_egress"].(bool) {
			securityDomainRule.NativeEgressDomain = true
		}
		if dn["native_firewall"].(bool) {
			securityDomainRule.NativeFirewallDomain = true
		}

		if !client.SecurityDomainRuleValidation(&securityDomainRule) {
			return fmt.Errorf("only one or none of 'firewall_domain', 'native_egress' and 'native_firewall' could be set true")
		}

		mapSecurityDomainRule[securityDomainRule.Name] = [3]bool{securityDomainRule.AviatrixFirewallDomain, securityDomainRule.NativeEgressDomain, securityDomainRule.NativeFirewallDomain}

		for _, connectedDomain := range dn["connected_domains"].([]interface{}) {
			securityDomainRule.ConnectedDomain = append(securityDomainRule.ConnectedDomain, connectedDomain.(string))
			tempDomainConn := []string{dn["security_domain_name"].(string), connectedDomain.(string)}
			domainConnAll = append(domainConnAll, tempDomainConn)
		}

		for _, attachedVPCs := range dn["attached_vpc"].([]interface{}) {

			attachedVPC := attachedVPCs.(map[string]interface{})

			if !manageVpcAttachment && attachedVPC != nil {
				return fmt.Errorf("manage_vpc_attachment is set to false. 'attached_vpc' should be empty")
			}

			if dn["security_domain_name"].(string) == "Aviatrix_Edge_Domain" && attachedVPC != nil {
				return fmt.Errorf("validation of source file failed: no VPCs should be attached to 'Aviatrix_Edge_Domain'")
			}

			vpcSolo := goaviatrix.VPCSolo{
				Region:                       attachedVPC["vpc_region"].(string),
				AccountName:                  attachedVPC["vpc_account_name"].(string),
				VpcID:                        attachedVPC["vpc_id"].(string),
				Subnets:                      attachedVPC["subnets"].(string),
				RouteTables:                  attachedVPC["route_tables"].(string),
				CustomizedRoutes:             attachedVPC["customized_routes"].(string),
				CustomizedRouteAdvertisement: attachedVPC["customized_route_advertisement"].(string),
				DisableLocalRoutePropagation: attachedVPC["disable_local_route_propagation"].(bool),
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

			tempAttachedVPC := []string{dn["security_domain_name"].(string), attachedVPC["vpc_id"].(string),
				attachedVPC["vpc_account_name"].(string), attachedVPC["vpc_region"].(string)}

			if attachedVPC["disable_local_route_propagation"].(bool) {
				tempAttachedVPC = append(tempAttachedVPC, "yes")
			} else {
				tempAttachedVPC = append(tempAttachedVPC, "no")
			}

			if attachedVPC["customized_routes"].(string) != "" {
				tempAttachedVPC = append(tempAttachedVPC, attachedVPC["customized_routes"].(string))
			} else {
				tempAttachedVPC = append(tempAttachedVPC, "")
			}

			if attachedVPC["customized_route_advertisement"].(string) != "" {
				tempAttachedVPC = append(tempAttachedVPC, attachedVPC["customized_route_advertisement"].(string))
			} else {
				tempAttachedVPC = append(tempAttachedVPC, "")
			}

			if attachedVPC["subnets"].(string) != "" {
				tempAttachedVPC = append(tempAttachedVPC, attachedVPC["subnets"].(string))
			} else {
				tempAttachedVPC = append(tempAttachedVPC, "")
			}

			if attachedVPC["route_tables"].(string) != "" {
				tempAttachedVPC = append(tempAttachedVPC, attachedVPC["route_tables"].(string))
			} else {
				tempAttachedVPC = append(tempAttachedVPC, "")
			}

			attachedVPCAll = append(attachedVPCAll, tempAttachedVPC)
		}

		awsTgw.SecurityDomains = append(awsTgw.SecurityDomains, securityDomainRule)
	}

	defaultDomainsWithCreation := []string{"Aviatrix_Edge_Domain", "Default_Domain", "Shared_Service_Domain"}
	if len(goaviatrix.Difference(defaultDomainsWithCreation, domainsAll)) != 0 {
		return fmt.Errorf("one or more of the three default domains are missing")
	}

	attachedGWs := d.Get("attached_aviatrix_transit_gateway").([]interface{})
	if manageTransitGwAttachment {
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
	} else if len(attachedGWs) != 0 {
		return fmt.Errorf("'manage_transit_gateway_attachment' is set to false. Please set it to true, or use " +
			"'aviatrix_aws_tgw_transit_gateway_attachment' to manage transit gateway attachments")
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

	flag := false
	defer resourceAviatrixAWSTgwReadIfRequired(d, meta, &flag)

	for i := range domainsToCreate {
		securityDomain := &goaviatrix.SecurityDomain{
			Name:                   domainsToCreate[i],
			AccountName:            d.Get("account_name").(string),
			Region:                 d.Get("region").(string),
			AwsTgwName:             d.Get("tgw_name").(string),
			AviatrixFirewallDomain: mapSecurityDomainRule[domainsToCreate[i]][0],
			NativeEgressDomain:     mapSecurityDomainRule[domainsToCreate[i]][1],
			NativeFirewallDomain:   mapSecurityDomainRule[domainsToCreate[i]][2],
		}
		err := client.CreateSecurityDomain(securityDomain)
		if err != nil {
			return fmt.Errorf("failed to create Security Domain: %s", err)
		}
	}

	for i := range domainConnPolicy {
		if len(domainConnPolicy[i]) == 2 {
			err := client.CreateDomainConnection(awsTgw, domainConnPolicy[i][0], domainConnPolicy[i][1])
			if err != nil {
				return fmt.Errorf("failed to create security domain connection: %s", err)
			}
		}
	}

	for i := range domainConnRemove {
		if len(domainConnRemove[i]) == 2 {
			err := client.DeleteDomainConnection(awsTgw, domainConnRemove[i][0], domainConnRemove[i][1])
			if err != nil {
				return fmt.Errorf("failed to delete domain connection: %s", err)
			}
		}
	}

	if manageTransitGwAttachment {
		for i := range attachedGWAll {
			gateway := &goaviatrix.Gateway{
				GwName: attachedGWAll[i],
			}
			err := client.AttachAviatrixTransitGWToAWSTgw(awsTgw, gateway, "Aviatrix_Edge_Domain")
			if err != nil {
				return fmt.Errorf("failed to attach transit GW: %s", err)
			}
		}
	}

	for i := range attachedVPCAll {
		if len(attachedVPCAll[i]) == 9 {
			vpcSolo := goaviatrix.VPCSolo{
				Region:                       attachedVPCAll[i][3],
				AccountName:                  attachedVPCAll[i][2],
				VpcID:                        attachedVPCAll[i][1],
				CustomizedRoutes:             attachedVPCAll[i][5],
				CustomizedRouteAdvertisement: attachedVPCAll[i][6],
				Subnets:                      attachedVPCAll[i][7],
				RouteTables:                  attachedVPCAll[i][8],
			}
			if attachedVPCAll[i][4] == "yes" {
				vpcSolo.DisableLocalRoutePropagation = true
			} else {
				vpcSolo.DisableLocalRoutePropagation = false
			}

			if mapFireNetVpc[attachedVPCAll[i][0]] {
				err := client.ConnectFireNetWithTgw(awsTgw, vpcSolo, attachedVPCAll[i][0])
				if err != nil {
					return fmt.Errorf("failed to attach FireNet VPC: %s", err)
				}
			} else {
				err := client.AttachVpcToAWSTgw(awsTgw, vpcSolo, attachedVPCAll[i][0])
				if err != nil {
					return fmt.Errorf("failed to attach VPC: %s", err)
				}
			}
		}
	}

	if cidrs := getStringSet(d, "cidrs"); len(cidrs) != 0 {
		err := client.UpdateTGWCidrs(awsTgw.Name, cidrs)
		if err != nil {
			return fmt.Errorf("could not update TGW CIDRs after creation: %v", err)
		}
	}

	return resourceAviatrixAWSTgwReadIfRequired(d, meta, &flag)
}

func resourceAviatrixAWSTgwReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixAWSTgwRead(d, meta)
	}
	return nil
}

func resourceAviatrixAWSTgwRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	tgwName := d.Get("tgw_name").(string)
	if tgwName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no aws tgw name received. Import Id is %s", id)
		d.Set("tgw_name", id)
		d.Set("manage_vpc_attachment", true)
		d.Set("manage_transit_gateway_attachment", true)
		d.SetId(id)
	}

	awsTgw := &goaviatrix.AWSTgw{
		Name: d.Get("tgw_name").(string),
	}
	awsTgw, err := client.ListTgwDetails(awsTgw)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find AWS TGW %s: %v", awsTgw.Name, err)
	}
	d.Set("account_name", awsTgw.AccountName)
	d.Set("tgw_name", awsTgw.Name)
	d.Set("region", awsTgw.Region)
	d.Set("cloud_type", awsTgw.CloudType)
	d.Set("enable_multicast", awsTgw.EnableMulticast)
	if err := d.Set("cidrs", awsTgw.CidrList); err != nil {
		return fmt.Errorf("could not set aws_tgw.cidrs into state: %v", err)
	}

	log.Printf("[INFO] Reading AWS TGW")

	awsTgw, err2 := client.GetAWSTgw(awsTgw)
	if err2 != nil {
		return fmt.Errorf("couldn't find AWS TGW %s: %v", tgwName, err2)
	}

	d.Set("aws_side_as_number", awsTgw.AwsSideAsNumber)

	manageTransitGwAttachment := d.Get("manage_transit_gateway_attachment").(bool)
	if manageTransitGwAttachment {
		if err := d.Set("attached_aviatrix_transit_gateway", awsTgw.AttachedAviatrixTransitGW); err != nil {
			log.Printf("[WARN] Error setting 'attached_aviatrix_transit_gateway' for (%s): %s", d.Id(), err)
		}
	}

	manageVpcAttachment := d.Get("manage_vpc_attachment").(bool)

	mSecurityDomain := make(map[string]map[string]interface{})
	for _, sd := range awsTgw.SecurityDomains {
		sdr := make(map[string]interface{})
		sdr["security_domain_name"] = sd.Name
		sdr["connected_domains"] = sd.ConnectedDomain
		sdr["aviatrix_firewall"] = sd.AviatrixFirewallDomain
		sdr["native_egress"] = sd.NativeEgressDomain
		sdr["native_firewall"] = sd.NativeFirewallDomain

		if manageVpcAttachment {
			var aVPCs []interface{}
			for _, attachedVPC := range sd.AttachedVPCs {
				vpcSolo := make(map[string]interface{})
				vpcSolo["vpc_region"] = attachedVPC.Region
				vpcSolo["vpc_account_name"] = attachedVPC.AccountName
				vpcSolo["vpc_id"] = attachedVPC.VpcID
				vpcSolo["customized_routes"] = attachedVPC.CustomizedRoutes
				vpcSolo["customized_route_advertisement"] = attachedVPC.CustomizedRouteAdvertisement
				vpcSolo["subnets"] = attachedVPC.Subnets
				vpcSolo["route_tables"] = attachedVPC.RouteTables
				vpcSolo["disable_local_route_propagation"] = attachedVPC.DisableLocalRoutePropagation
				aVPCs = append(aVPCs, vpcSolo)
			}
			sdr["attached_vpc"] = aVPCs
		}

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

			if manageVpcAttachment {
				mVPC := make(map[string]bool)
				var aVPCNew []map[string]interface{}

				for _, attachedVPCs := range mSecurityDomain[dn["security_domain_name"].(string)]["attached_vpc"].([]interface{}) {
					attachedVPC := attachedVPCs.(map[string]interface{})
					mVPC[attachedVPC["vpc_id"].(string)] = true
				}

				for _, attachedVPCs := range dn["attached_vpc"].([]interface{}) {
					attachedVPC := attachedVPCs.(map[string]interface{})
					if mVPC[attachedVPC["vpc_id"].(string)] {
						for _, attachedVPCsFromRefresh := range mSecurityDomain[dn["security_domain_name"].(string)]["attached_vpc"].([]interface{}) {
							attachedVPCFromRefresh := attachedVPCsFromRefresh.(map[string]interface{})
							if attachedVPCFromRefresh["vpc_id"] == attachedVPC["vpc_id"] {
								attachedVPC["vpc_account_name"] = attachedVPCFromRefresh["vpc_account_name"]
								if attachedVPC["subnets"].(string) != "" {
									subnetsFromConfigList := strings.Split(attachedVPC["subnets"].(string), ",")
									var subnetsFromReadList []string
									subnetsFromReadList = strings.Split(attachedVPCFromRefresh["subnets"].(string), ",")
									if len(goaviatrix.Difference(subnetsFromConfigList, subnetsFromReadList)) != 0 ||
										len(goaviatrix.Difference(subnetsFromReadList, subnetsFromConfigList)) != 0 {
										attachedVPC["subnets"] = attachedVPCFromRefresh["subnets"]
									}
								} else {
									attachedVPC["subnets"] = ""
								}
								if attachedVPC["route_tables"].(string) != "" {
									routeTablesFromConfigList := strings.Split(attachedVPC["route_tables"].(string), ",")
									for i := 0; i < len(routeTablesFromConfigList); i++ {
										routeTablesFromConfigList[i] = strings.TrimSpace(routeTablesFromConfigList[i])
									}
									var routeTablesFromReadList []string
									routeTablesFromReadList = strings.Split(attachedVPCFromRefresh["route_tables"].(string), ",")
									for i := 0; i < len(routeTablesFromReadList); i++ {
										routeTablesFromReadList[i] = strings.TrimSpace(routeTablesFromReadList[i])
									}
									if (len(goaviatrix.Difference(routeTablesFromConfigList, routeTablesFromReadList)) != 0 ||
										len(goaviatrix.Difference(routeTablesFromReadList, routeTablesFromConfigList)) != 0) &&
										attachedVPCFromRefresh["route_tables"] != "ALL" &&
										attachedVPCFromRefresh["route_tables"] != "All" {
										attachedVPC["route_tables"] = attachedVPCFromRefresh["route_tables"]
									}
								} else {
									attachedVPC["route_tables"] = ""
								}
								attachedVPC["vpc_region"] = attachedVPCFromRefresh["vpc_region"]
								attachedVPC["customized_routes"] = attachedVPCFromRefresh["customized_routes"]
								attachedVPC["customized_route_advertisement"] = attachedVPCFromRefresh["customized_route_advertisement"]
								attachedVPC["disable_local_route_propagation"] = attachedVPCFromRefresh["disable_local_route_propagation"]
							}
						}
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
			}

			securityDomains = append(securityDomains, mSecurityDomain[dn["security_domain_name"].(string)])
		}
	}

	for _, dn := range awsTgw.SecurityDomains {
		if !mOld[dn.Name] {
			securityDomains = append(securityDomains, mSecurityDomain[dn.Name])
		}
	}

	if err := d.Set("security_domains", securityDomains); err != nil {
		log.Printf("[WARN] Error setting security_domains for (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceAviatrixAWSTgwUpdate(d *schema.ResourceData, meta interface{}) error {
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
	var toUpdateCustomizedRoutesOnly [][]string
	var toUpdateCustomizedRoutesAdOnly [][]string
	mapOldFireNetVpc := make(map[string]bool)
	mapNewFireNetVpc := make(map[string]bool)

	d.Partial(true)

	if d.HasChange("account_name") {
		return fmt.Errorf("updating account_name is not allowed")
	}
	if d.HasChange("region") {
		return fmt.Errorf("updating region is not allowed")
	}
	if d.HasChange("aws_side_as_number") {
		return fmt.Errorf("updating aws_side_as_number is not allowed")
	}
	if d.HasChange("cloud_type") {
		return fmt.Errorf("updating cloud_type is not allowed")
	}
	if d.HasChange("enable_multicast") {
		return fmt.Errorf("updating enable_multicast is not allowed")
	}

	manageVpcAttachment := d.Get("manage_vpc_attachment").(bool)
	if d.HasChange("manage_vpc_attachment") {
		_, nMVA := d.GetChange("manage_vpc_attachment")
		newManageVpcAttachment := nMVA.(bool)
		if newManageVpcAttachment {
			d.Set("manage_vpc_attachment", true)
		} else {
			d.Set("manage_vpc_attachment", false)
		}
	}

	manageTransitGwAttachment := d.Get("manage_transit_gateway_attachment").(bool)
	if d.HasChange("manage_transit_gateway_attachment") {
		_, nMTGA := d.GetChange("manage_transit_gateway_attachment")
		newManageTransitGwAttachment := nMTGA.(bool)
		if newManageTransitGwAttachment {
			d.Set("manage_transit_gateway_attachment", true)
		} else {
			d.Set("manage_transit_gateway_attachment", false)
		}
	}
	if !manageTransitGwAttachment && len(d.Get("attached_aviatrix_transit_gateway").([]interface{})) != 0 {
		return fmt.Errorf("'manage_transit_gateway_attachment' is set to false. Please set it to true, or use " +
			"'aviatrix_aws_tgw_transit_gateway_attachment' to manage transit gateway attachments")
	}

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

	mapSecurityDomainsOld := make(map[string][3]bool)
	mapSecurityDomainsNew := make(map[string][3]bool)
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

			if dn["aviatrix_firewall"].(bool) {
				securityDomainRule.AviatrixFirewallDomain = true
				mapOldFireNetVpc[securityDomainRule.Name] = true
			}
			if dn["native_egress"].(bool) {
				securityDomainRule.NativeEgressDomain = true
			}
			if dn["native_firewall"].(bool) {
				securityDomainRule.NativeFirewallDomain = true
			}

			mapSecurityDomainsOld[securityDomainRule.Name] = [3]bool{securityDomainRule.AviatrixFirewallDomain, securityDomainRule.NativeEgressDomain, securityDomainRule.NativeFirewallDomain}

			for _, connectedDomain := range dn["connected_domains"].([]interface{}) {
				securityDomainRule.ConnectedDomain = append(securityDomainRule.ConnectedDomain, connectedDomain.(string))
				tempDomainConn := []string{dn["security_domain_name"].(string), connectedDomain.(string)}
				domainConnOld = append(domainConnOld, tempDomainConn)
			}

			for _, attachedVPCs := range dn["attached_vpc"].([]interface{}) {

				attachedVPC := attachedVPCs.(map[string]interface{})

				if dn["security_domain_name"].(string) == "Aviatrix_Edge_Domain" && attachedVPC != nil {
					return fmt.Errorf("validation of source file failed: no VPCs should be attached to 'Aviatrix_Edge_Domain'")
				}

				vpcSolo := goaviatrix.VPCSolo{
					Region:                       attachedVPC["vpc_region"].(string),
					AccountName:                  attachedVPC["vpc_account_name"].(string),
					VpcID:                        attachedVPC["vpc_id"].(string),
					Subnets:                      attachedVPC["subnets"].(string),
					RouteTables:                  attachedVPC["route_tables"].(string),
					CustomizedRoutes:             attachedVPC["customized_routes"].(string),
					CustomizedRouteAdvertisement: attachedVPC["customized_route_advertisement"].(string),
					DisableLocalRoutePropagation: attachedVPC["disable_local_route_propagation"].(bool),
				}
				securityDomainRule.AttachedVPCs = append(securityDomainRule.AttachedVPCs, vpcSolo)

				tempAttachedVPC := []string{dn["security_domain_name"].(string), attachedVPC["vpc_id"].(string),
					attachedVPC["vpc_account_name"].(string), attachedVPC["vpc_region"].(string)}
				if attachedVPC["disable_local_route_propagation"].(bool) {
					tempAttachedVPC = append(tempAttachedVPC, "yes")
				} else {
					tempAttachedVPC = append(tempAttachedVPC, "no")
				}

				if attachedVPC["customized_routes"].(string) != "" {
					tempAttachedVPC = append(tempAttachedVPC, attachedVPC["customized_routes"].(string))
				} else {
					tempAttachedVPC = append(tempAttachedVPC, "")
				}

				if attachedVPC["customized_route_advertisement"].(string) != "" {
					tempAttachedVPC = append(tempAttachedVPC, attachedVPC["customized_route_advertisement"].(string))
				} else {
					tempAttachedVPC = append(tempAttachedVPC, "")
				}

				if attachedVPC["subnets"].(string) != "" {
					tempAttachedVPC = append(tempAttachedVPC, attachedVPC["subnets"].(string))
				} else {
					tempAttachedVPC = append(tempAttachedVPC, "")
				}

				if attachedVPC["route_tables"].(string) != "" {
					tempAttachedVPC = append(tempAttachedVPC, attachedVPC["route_tables"].(string))
				} else {
					tempAttachedVPC = append(tempAttachedVPC, "")
				}

				attachedVPCOld = append(attachedVPCOld, tempAttachedVPC)
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
			if dn["aviatrix_firewall"].(bool) {
				securityDomainRule.AviatrixFirewallDomain = true
				mapNewFireNetVpc[securityDomainRule.Name] = true
			}
			if dn["native_egress"].(bool) {
				securityDomainRule.NativeEgressDomain = true
			}
			if dn["native_firewall"].(bool) {
				securityDomainRule.NativeFirewallDomain = true
			}

			if !client.SecurityDomainRuleValidation(&securityDomainRule) {
				return fmt.Errorf("only one or none of 'firewall_domain', 'native_egress' and 'native_firewall' could be set true")
			}

			mapSecurityDomainsNew[securityDomainRule.Name] = [3]bool{securityDomainRule.AviatrixFirewallDomain, securityDomainRule.NativeEgressDomain, securityDomainRule.NativeFirewallDomain}

			if val, ok := mapSecurityDomainsOld[securityDomainRule.Name]; ok {
				if val[0] != securityDomainRule.AviatrixFirewallDomain {
					return fmt.Errorf("cannot update 'aviatrix_firewall'")
				}
				if val[1] != securityDomainRule.NativeEgressDomain {
					return fmt.Errorf("cannot update 'native_egress'")
				}
				if val[2] != securityDomainRule.NativeFirewallDomain {
					return fmt.Errorf("cannot update 'native_firewall'")
				}
			}

			for _, connectedDomain := range dn["connected_domains"].([]interface{}) {
				securityDomainRule.ConnectedDomain = append(securityDomainRule.ConnectedDomain, connectedDomain.(string))
				tempDomainConn := []string{dn["security_domain_name"].(string), connectedDomain.(string)}
				domainConnNew = append(domainConnNew, tempDomainConn)
			}

			for _, attachedVPCs := range dn["attached_vpc"].([]interface{}) {

				attachedVPC := attachedVPCs.(map[string]interface{})

				if !manageVpcAttachment && attachedVPC != nil {
					return fmt.Errorf("manage_vpc_attachment is set to false. 'attached_vpc' should be empty")
				}

				if dn["security_domain_name"].(string) == "Aviatrix_Edge_Domain" && attachedVPC != nil {
					return fmt.Errorf("validation of source file failed: no VPCs should be attached to 'Aviatrix_Edge_Domain'")
				}

				vpcSolo := goaviatrix.VPCSolo{
					Region:                       attachedVPC["vpc_region"].(string),
					AccountName:                  attachedVPC["vpc_account_name"].(string),
					VpcID:                        attachedVPC["vpc_id"].(string),
					Subnets:                      attachedVPC["subnets"].(string),
					RouteTables:                  attachedVPC["route_tables"].(string),
					CustomizedRoutes:             attachedVPC["customized_routes"].(string),
					CustomizedRouteAdvertisement: attachedVPC["customized_route_advertisement"].(string),
					DisableLocalRoutePropagation: attachedVPC["disable_local_route_propagation"].(bool),
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

				tempAttachedVPC := []string{dn["security_domain_name"].(string), attachedVPC["vpc_id"].(string),
					attachedVPC["vpc_account_name"].(string), attachedVPC["vpc_region"].(string)}
				if attachedVPC["disable_local_route_propagation"].(bool) {
					tempAttachedVPC = append(tempAttachedVPC, "yes")
				} else {
					tempAttachedVPC = append(tempAttachedVPC, "no")
				}

				if attachedVPC["customized_routes"].(string) != "" {
					tempAttachedVPC = append(tempAttachedVPC, attachedVPC["customized_routes"].(string))
				} else {
					tempAttachedVPC = append(tempAttachedVPC, "")
				}

				if attachedVPC["customized_route_advertisement"].(string) != "" {
					tempAttachedVPC = append(tempAttachedVPC, attachedVPC["customized_route_advertisement"].(string))
				} else {
					tempAttachedVPC = append(tempAttachedVPC, "")
				}

				if attachedVPC["subnets"].(string) != "" {
					tempAttachedVPC = append(tempAttachedVPC, attachedVPC["subnets"].(string))
				} else {
					tempAttachedVPC = append(tempAttachedVPC, "")
				}

				if attachedVPC["route_tables"].(string) != "" {
					tempAttachedVPC = append(tempAttachedVPC, attachedVPC["route_tables"].(string))
				} else {
					tempAttachedVPC = append(tempAttachedVPC, "")
				}

				attachedVPCNew = append(attachedVPCNew, tempAttachedVPC)
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

		toAttachVPCs = goaviatrix.DifferenceSliceAttachedVPC(attachedVPCNew, attachedVPCOld)
		toDetachVPCs = goaviatrix.DifferenceSliceAttachedVPC(attachedVPCOld, attachedVPCNew)

		toUpdateCustomizedRoutesOnly, toUpdateCustomizedRoutesAdOnly = goaviatrix.ValidateAttachedVPCsForCustomizedRoutes(attachedVPCOld, attachedVPCNew)

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

	if manageTransitGwAttachment {
		for i := range toDetachGWs {
			gateway := &goaviatrix.Gateway{
				GwName: toDetachGWs[i],
			}

			err := client.DetachAviatrixTransitGWFromAWSTgw(awsTgw, gateway, "Aviatrix_Edge_Domain")
			if err != nil {
				resourceAviatrixAWSTgwRead(d, meta)
				return fmt.Errorf("failed to detach transit GW: %s", err)
			}
		}
	}

	for i := range domainsToCreate {
		securityDomain := &goaviatrix.SecurityDomain{
			Name:                   domainsToCreate[i],
			AccountName:            d.Get("account_name").(string),
			Region:                 d.Get("region").(string),
			AwsTgwName:             d.Get("tgw_name").(string),
			AviatrixFirewallDomain: mapSecurityDomainsNew[domainsToCreate[i]][0],
			NativeEgressDomain:     mapSecurityDomainsNew[domainsToCreate[i]][1],
			NativeFirewallDomain:   mapSecurityDomainsNew[domainsToCreate[i]][2],
		}

		err := client.CreateSecurityDomain(securityDomain)
		if err != nil {
			resourceAviatrixAWSTgwRead(d, meta)
			return fmt.Errorf("failed to create Security Domain: %s", err)
		}
	}

	for i := range domainConnRemove {
		if len(domainConnRemove[i]) == 2 {
			err := client.DeleteDomainConnection(awsTgw, domainConnRemove[i][0], domainConnRemove[i][1])
			if err != nil {
				resourceAviatrixAWSTgwRead(d, meta)
				return fmt.Errorf("failed to delete domain connection: %s", err)
			}
		}
	}

	for i := range domainConnPolicy {
		if len(domainConnPolicy[i]) == 2 {
			err := client.CreateDomainConnection(awsTgw, domainConnPolicy[i][0], domainConnPolicy[i][1])
			if err != nil {
				resourceAviatrixAWSTgwRead(d, meta)
				return fmt.Errorf("failed to create security domain connection: %s", err)
			}
		}
	}

	if manageVpcAttachment {
		for i := range toDetachVPCs {
			if len(toDetachVPCs[i]) == 9 {
				if mapOldFireNetVpc[toDetachVPCs[i][0]] {
					err := client.DisconnectFireNetFromTgw(awsTgw, toDetachVPCs[i][1])
					if err != nil {
						return fmt.Errorf("failed to detach FireNet VPC: %s", err)
					}
				} else {
					err := client.DetachVpcFromAWSTgw(awsTgw, toDetachVPCs[i][1])
					if err != nil {
						resourceAviatrixAWSTgwRead(d, meta)
						return fmt.Errorf("failed to detach VPC: %s", err)
					}
				}
			}
		}
	}

	if manageTransitGwAttachment {
		for i := range toAttachGWs {
			gateway := &goaviatrix.Gateway{
				GwName: toAttachGWs[i],
			}

			err := client.AttachAviatrixTransitGWToAWSTgw(awsTgw, gateway, "Aviatrix_Edge_Domain")
			if err != nil {
				resourceAviatrixAWSTgwRead(d, meta)
				return fmt.Errorf("failed to attach transit GW: %s", err)
			}
		}
	}

	if manageVpcAttachment {
		for i := range toAttachVPCs {
			if len(toAttachVPCs[i]) == 9 {
				vpcSolo := goaviatrix.VPCSolo{
					Region:                       toAttachVPCs[i][3],
					AccountName:                  toAttachVPCs[i][2],
					VpcID:                        toAttachVPCs[i][1],
					CustomizedRoutes:             toAttachVPCs[i][5],
					CustomizedRouteAdvertisement: toAttachVPCs[i][6],
					Subnets:                      toAttachVPCs[i][7],
					RouteTables:                  toAttachVPCs[i][8],
				}
				if toAttachVPCs[i][4] == "yes" {
					vpcSolo.DisableLocalRoutePropagation = true
				} else {
					vpcSolo.DisableLocalRoutePropagation = false
				}
				res, _ := client.IsVpcAttachedToTgw(awsTgw, &vpcSolo)
				if !res {
					if mapNewFireNetVpc[toAttachVPCs[i][0]] {
						err := client.ConnectFireNetWithTgw(awsTgw, vpcSolo, toAttachVPCs[i][0])
						if err != nil {
							return fmt.Errorf("failed to attach FireNet VPC: %s", err)
						}
					} else {
						err := client.AttachVpcToAWSTgw(awsTgw, vpcSolo, toAttachVPCs[i][0])
						if err != nil {
							return fmt.Errorf("failed to attach VPC: %s", err)
						}
					}
				}
			}
		}
	}

	if manageVpcAttachment {
		for i := range toUpdateCustomizedRoutesOnly {
			if len(toUpdateCustomizedRoutesOnly[i]) == 9 {
				awsTgwVpcAttachment := &goaviatrix.AwsTgwVpcAttachment{
					TgwName:          d.Get("tgw_name").(string),
					VpcID:            toUpdateCustomizedRoutesOnly[i][1],
					CustomizedRoutes: toUpdateCustomizedRoutesOnly[i][5],
				}
				err := client.EditTgwSpokeVpcCustomizedRoutes(awsTgwVpcAttachment)
				if err != nil {
					return fmt.Errorf("failed to update spoke vpc customized routes: %s", err)
				}
			}
		}
	}

	if manageVpcAttachment {
		for i := range toUpdateCustomizedRoutesAdOnly {
			if len(toUpdateCustomizedRoutesAdOnly[i]) == 9 {
				awsTgwVpcAttachment := &goaviatrix.AwsTgwVpcAttachment{
					TgwName:                      d.Get("tgw_name").(string),
					VpcID:                        toUpdateCustomizedRoutesAdOnly[i][1],
					CustomizedRouteAdvertisement: toUpdateCustomizedRoutesAdOnly[i][6],
				}
				err := client.EditTgwSpokeVpcCustomizedRouteAdvertisement(awsTgwVpcAttachment)
				if err != nil {
					return fmt.Errorf("failed to update spoke vpc customized routes advertisement: %s", err)
				}
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
			resourceAviatrixAWSTgwRead(d, meta)
			return fmt.Errorf("failed to delete Security Domain: %s", err)
		}
	}

	if d.HasChange("cidrs") {
		cidrs := getStringSet(d, "cidrs")
		err := client.UpdateTGWCidrs(awsTgw.Name, cidrs)
		if err != nil {
			return fmt.Errorf("could not update TGW CIDRs during update: %v", err)
		}
	}

	d.Partial(false)
	d.SetId(awsTgw.Name)
	return resourceAviatrixAWSTgwRead(d, meta)
}

func resourceAviatrixAWSTgwDelete(d *schema.ResourceData, meta interface{}) error {
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

	mapFireNetVpc := make(map[string]bool)

	manageVpcAttachment := d.Get("manage_vpc_attachment").(bool)
	if manageVpcAttachment {
		var attachedVPCs [][]string
		domains := d.Get("security_domains").([]interface{})
		for _, domain := range domains {
			dn := domain.(map[string]interface{})

			for _, aVPCs := range dn["attached_vpc"].([]interface{}) {
				aVPC := aVPCs.(map[string]interface{})

				tempAttachedVPC := []string{dn["security_domain_name"].(string), aVPC["vpc_id"].(string),
					aVPC["vpc_account_name"].(string), aVPC["vpc_region"].(string)}

				if dn["aviatrix_firewall"].(bool) {
					mapFireNetVpc[dn["security_domain_name"].(string)] = true
				}

				attachedVPCs = append(attachedVPCs, tempAttachedVPC)
			}
		}
		for i := range attachedVPCs {
			if len(attachedVPCs[i]) == 4 {
				if mapFireNetVpc[attachedVPCs[i][0]] {
					err := client.DisconnectFireNetFromTgw(awsTgw, attachedVPCs[i][1])
					if err != nil {
						return fmt.Errorf("failed to detach FireNet VPC: %s", err)
					}
				} else {
					err := client.DetachVpcFromAWSTgw(awsTgw, attachedVPCs[i][1])
					if err != nil {
						resourceAviatrixAWSTgwRead(d, meta)
						return fmt.Errorf("failed to detach VPC: %s", err)
					}
				}
			}
		}
	}

	manageTransitGwAttachment := d.Get("manage_transit_gateway_attachment").(bool)
	if manageTransitGwAttachment {
		var attachedGWs []string
		transitGWs := d.Get("attached_aviatrix_transit_gateway").([]interface{})
		for _, transitGW := range transitGWs {
			attachedGWs = append(attachedGWs, transitGW.(string))
		}
		for i := range attachedGWs {
			gateway := &goaviatrix.Gateway{
				GwName: attachedGWs[i],
			}

			err := client.DetachAviatrixTransitGWFromAWSTgw(awsTgw, gateway, "Aviatrix_Edge_Domain")
			if err != nil {
				resourceAviatrixAWSTgwRead(d, meta)
				return fmt.Errorf("failed to detach transit GW: %s", err)
			}
		}
	}

	err := client.DeleteAWSTgw(awsTgw)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't destroy AWS TGW %s: %v", awsTgw.Name, err)
	}

	return nil
}
