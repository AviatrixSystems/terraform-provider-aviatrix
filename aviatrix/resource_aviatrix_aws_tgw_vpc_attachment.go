package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixAwsTgwVpcAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixAwsTgwVpcAttachmentCreate,
		Read:   resourceAviatrixAwsTgwVpcAttachmentRead,
		Update: resourceAviatrixAwsTgwVpcAttachmentUpdate,
		Delete: resourceAviatrixAwsTgwVpcAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"tgw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the AWS TGW.",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Region of cloud provider.",
			},
			"security_domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the security domain.",
			},
			"vpc_account_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This parameter represents the name of a Cloud-Account in Aviatrix controller.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "This parameter represents the ID of the VPC.",
			},
			"customized_routes": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Advanced option. Customized Spoke VPC Routes. It allows the admin to enter non-RFC1918 routes in the VPC route table targeting the TGW.",
			},
			"subnets": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				DiffSuppressFunc: DiffSuppressFuncIgnoreSpaceInString,
				Description:      "Advanced option. VPC subnets separated by ',' to attach to the VPC. If left blank, Aviatrix Controller automatically selects a subnet representing each AZ for the VPC attachment.",
			},
			"route_tables": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
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
				ForceNew:    true,
				Default:     false,
				Description: "Advanced option. If set to true, it disables automatic route propagation of this VPC to other VPCs within the same security domain.",
			},
			"edge_attachment": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Edge attachment ID. To allow access to the private IP of the MGMT interface of the " +
					"Firewalls, set this attribute to enable Management Access From Onprem. This feature advertises " +
					"the Firewalls private MGMT subnet to your Edge domain.",
			},
		},
	}
}

func resourceAviatrixAwsTgwVpcAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	awsTgwVpcAttachment := &goaviatrix.AwsTgwVpcAttachment{
		TgwName:                      d.Get("tgw_name").(string),
		Region:                       d.Get("region").(string),
		SecurityDomainName:           d.Get("security_domain_name").(string),
		VpcAccountName:               d.Get("vpc_account_name").(string),
		VpcID:                        d.Get("vpc_id").(string),
		CustomizedRoutes:             d.Get("customized_routes").(string),
		Subnets:                      d.Get("subnets").(string),
		RouteTables:                  d.Get("route_tables").(string),
		CustomizedRouteAdvertisement: d.Get("customized_route_advertisement").(string),
		DisableLocalRoutePropagation: d.Get("disable_local_route_propagation").(bool),
		EdgeAttachment:               d.Get("edge_attachment").(string),
	}

	isFirewallSecurityDomain, err := client.IsFirewallSecurityDomain(awsTgwVpcAttachment.TgwName, awsTgwVpcAttachment.SecurityDomainName)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			return fmt.Errorf("could not find Security Domain: " + awsTgwVpcAttachment.SecurityDomainName)
		}
		return fmt.Errorf("could not find Security Domain due to: %v", err)
	}

	log.Printf("[INFO] Attaching vpc: %s to tgw %s", awsTgwVpcAttachment.VpcID, awsTgwVpcAttachment.TgwName)

	flag := false

	if isFirewallSecurityDomain {
		err = client.CreateAwsTgwVpcAttachmentForFireNet(awsTgwVpcAttachment)
		if err != nil {
			return fmt.Errorf("failed to create Aviatrix Aws Tgw Vpc Attach for FireNet: %s", err)
		}

		d.SetId(awsTgwVpcAttachment.TgwName + "~" + awsTgwVpcAttachment.SecurityDomainName + "~" + awsTgwVpcAttachment.VpcID)

		defer resourceAviatrixAwsTgwVpcAttachmentReadIfRequired(d, meta, &flag)

		if awsTgwVpcAttachment.EdgeAttachment != "" {
			err = client.UpdateFirewallAttachmentAccessFromOnprem(awsTgwVpcAttachment)
			if err != nil {
				return fmt.Errorf("failed to enable firewall attachment access from onprem: %s", err)
			}
		}
	} else {
		if awsTgwVpcAttachment.EdgeAttachment != "" {
			return fmt.Errorf("management access from onprem only works for FireNet")
		}

		err = client.CreateAwsTgwVpcAttachment(awsTgwVpcAttachment)
		if err != nil {
			return fmt.Errorf("failed to create Aviatrix Aws Tgw Vpc Attach: %s", err)
		}

		d.SetId(awsTgwVpcAttachment.TgwName + "~" + awsTgwVpcAttachment.SecurityDomainName + "~" + awsTgwVpcAttachment.VpcID)
	}

	return resourceAviatrixAwsTgwVpcAttachmentReadIfRequired(d, meta, &flag)
}

func resourceAviatrixAwsTgwVpcAttachmentReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixAwsTgwVpcAttachmentRead(d, meta)
	}
	return nil
}

func resourceAviatrixAwsTgwVpcAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	tgwName := d.Get("tgw_name").(string)
	securityDomainName := d.Get("security_domain_name").(string)
	vpcID := d.Get("vpc_id").(string)

	if tgwName == "" || securityDomainName == "" || vpcID == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no vpc names received. Import Id is %s", id)
		d.Set("tgw_name", strings.Split(id, "~")[0])
		d.Set("security_domain_name", strings.Split(id, "~")[1])
		d.Set("vpc_id", strings.Split(id, "~")[2])
	}
	awsTgwVpcAttachment := &goaviatrix.AwsTgwVpcAttachment{
		TgwName:            d.Get("tgw_name").(string),
		SecurityDomainName: d.Get("security_domain_name").(string),
		VpcID:              d.Get("vpc_id").(string),
	}

	aTVA, err := client.GetAwsTgwVpcAttachment(awsTgwVpcAttachment)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to get Aviatrix Aws Tgw Vpc Attach: %s", err)
	}
	if aTVA != nil {
		d.Set("tgw_name", aTVA.TgwName)
		d.Set("region", aTVA.Region)
		d.Set("security_domain_name", aTVA.SecurityDomainName)
		d.Set("vpc_account_name", aTVA.VpcAccountName)
		d.Set("vpc_id", aTVA.VpcID)
		d.Set("disable_local_route_propagation", aTVA.DisableLocalRoutePropagation)

		if d.Get("subnets").(string) != "" {
			subnetsFromConfigList := strings.Split(d.Get("subnets").(string), ",")
			subnetsFromReadList := strings.Split(aTVA.Subnets, ",")
			if len(goaviatrix.Difference(subnetsFromConfigList, subnetsFromReadList)) == 0 ||
				len(goaviatrix.Difference(subnetsFromReadList, subnetsFromConfigList)) == 0 {
				d.Set("subnets", d.Get("subnets").(string))
			} else {
				d.Set("subnets", aTVA.Subnets)
			}
		} else {
			d.Set("subnets", aTVA.Subnets)
		}
		if d.Get("route_tables").(string) != "" {
			routeTablesFromConfigList := strings.Split(d.Get("route_tables").(string), ",")
			for i := 0; i < len(routeTablesFromConfigList); i++ {
				routeTablesFromConfigList[i] = strings.TrimSpace(routeTablesFromConfigList[i])
			}
			routeTablesFromReadList := strings.Split(aTVA.RouteTables, ",")
			for i := 0; i < len(routeTablesFromReadList); i++ {
				routeTablesFromReadList[i] = strings.TrimSpace(routeTablesFromReadList[i])
			}
			if (len(goaviatrix.Difference(routeTablesFromConfigList, routeTablesFromReadList)) != 0 ||
				len(goaviatrix.Difference(routeTablesFromReadList, routeTablesFromConfigList)) != 0) &&
				aTVA.RouteTables != "ALL" &&
				aTVA.RouteTables != "All" {
				d.Set("route_tables", aTVA.RouteTables)
			} else {
				d.Set("route_tables", d.Get("route_tables").(string))
			}
		} else {
			d.Set("route_tables", aTVA.RouteTables)
		}

		d.Set("customized_routes", aTVA.CustomizedRoutes)
		d.Set("customized_route_advertisement", aTVA.CustomizedRouteAdvertisement)
		d.Set("edge_attachment", aTVA.EdgeAttachment)
		d.SetId(awsTgwVpcAttachment.TgwName + "~" + awsTgwVpcAttachment.SecurityDomainName + "~" + awsTgwVpcAttachment.VpcID)
		return nil
	}

	return fmt.Errorf("no Aviatrix Aws Tgw Vpc Attach found")
}

func resourceAviatrixAwsTgwVpcAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	flag := false
	defer resourceAviatrixAwsTgwVpcAttachmentReadIfRequired(d, meta, &flag)

	client := meta.(*goaviatrix.Client)

	d.Partial(true)
	if d.HasChange("region") {
		return fmt.Errorf("updating region is not allowed")
	}
	if d.HasChange("vpc_account_name") {
		return fmt.Errorf("updating vpc_account_name is not allowed")
	}
	if d.HasChange("customized_routes") {
		awsTgwVpcAttachment := &goaviatrix.AwsTgwVpcAttachment{
			TgwName:          d.Get("tgw_name").(string),
			VpcID:            d.Get("vpc_id").(string),
			CustomizedRoutes: d.Get("customized_routes").(string),
		}
		err := client.EditTgwSpokeVpcCustomizedRoutes(awsTgwVpcAttachment)
		if err != nil {
			return fmt.Errorf("failed to update spoke vpc customized routes: %s", err)
		}
	}
	if d.HasChange("customized_route_advertisement") {
		awsTgwVpcAttachment := &goaviatrix.AwsTgwVpcAttachment{
			TgwName:                      d.Get("tgw_name").(string),
			VpcID:                        d.Get("vpc_id").(string),
			CustomizedRouteAdvertisement: d.Get("customized_route_advertisement").(string),
		}
		err := client.EditTgwSpokeVpcCustomizedRouteAdvertisement(awsTgwVpcAttachment)
		if err != nil {
			return fmt.Errorf("failed to update spoke vpc customized routes advertisement: %s", err)
		}
	}

	if d.HasChange("edge_attachment") {
		awsTgwVpcAttachment := &goaviatrix.AwsTgwVpcAttachment{
			TgwName:            d.Get("tgw_name").(string),
			VpcID:              d.Get("vpc_id").(string),
			SecurityDomainName: d.Get("security_domain_name").(string),
			EdgeAttachment:     d.Get("edge_attachment").(string),
		}

		isFirewallSecurityDomain, err := client.IsFirewallSecurityDomain(awsTgwVpcAttachment.TgwName, awsTgwVpcAttachment.SecurityDomainName)
		if err != nil {
			if err == goaviatrix.ErrNotFound {
				return fmt.Errorf("could not find Security Domain: " + awsTgwVpcAttachment.SecurityDomainName)
			}
			return fmt.Errorf("could not find Security Domain due to: %v", err)
		}

		oldEA, newEA := d.GetChange("edge_attachment")
		oldEAString := oldEA.(string)
		newEAString := newEA.(string)

		if isFirewallSecurityDomain {
			if oldEAString != "" && newEAString != "" {
				awsTgwVpcAttachment.EdgeAttachment = ""

				err := client.UpdateFirewallAttachmentAccessFromOnprem(awsTgwVpcAttachment)
				if err != nil {
					return fmt.Errorf("failed to disable firewall attachment access from onprem while updating: %s", err)
				}

				awsTgwVpcAttachment.EdgeAttachment = newEAString

				err = client.UpdateFirewallAttachmentAccessFromOnprem(awsTgwVpcAttachment)
				if err != nil {
					return fmt.Errorf("failed to enable firewall attachment access from onprem while updating: %s", err)
				}
			} else {
				err := client.UpdateFirewallAttachmentAccessFromOnprem(awsTgwVpcAttachment)
				if err != nil {
					return fmt.Errorf("failed to update firewall attachment access from onprem: %s", err)
				}
			}
		} else {
			if newEAString != "" {
				return fmt.Errorf("management access from onprem only works for FireNet")
			}
		}

	}

	d.Partial(false)
	return resourceAviatrixAwsTgwVpcAttachmentReadIfRequired(d, meta, &flag)
}

func resourceAviatrixAwsTgwVpcAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	awsTgwVpcAttachment := &goaviatrix.AwsTgwVpcAttachment{
		TgwName:            d.Get("tgw_name").(string),
		Region:             d.Get("region").(string),
		SecurityDomainName: d.Get("security_domain_name").(string),
		VpcAccountName:     d.Get("vpc_account_name").(string),
		VpcID:              d.Get("vpc_id").(string),
	}

	isFirewallSecurityDomain, err := client.IsFirewallSecurityDomain(awsTgwVpcAttachment.TgwName, awsTgwVpcAttachment.SecurityDomainName)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			return fmt.Errorf("could not find Security Domain: " + awsTgwVpcAttachment.VpcID)
		}
		return fmt.Errorf(("could not find Security Domain due to: ") + err.Error())
	}

	if isFirewallSecurityDomain {
		err := client.DeleteAwsTgwVpcAttachmentForFireNet(awsTgwVpcAttachment)
		if err != nil {
			return fmt.Errorf("failed to detach FireNet VPC from TGW: %s", err)
		}
	} else {
		err := client.DeleteAwsTgwVpcAttachment(awsTgwVpcAttachment)
		if err != nil {
			return fmt.Errorf("failed to detach VPC from TGW: %s", err)
		}
	}

	return nil
}
