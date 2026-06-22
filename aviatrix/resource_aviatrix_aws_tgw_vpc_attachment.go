package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixAwsTgwVpcAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixAwsTgwVpcAttachmentCreate,
		Read:   resourceAviatrixAwsTgwVpcAttachmentRead,
		Update: resourceAviatrixAwsTgwVpcAttachmentUpdate,
		Delete: resourceAviatrixAwsTgwVpcAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
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
			"network_domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the network domain.",
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
				Description: "Advanced option. If set to true, it disables automatic route propagation of this VPC to other VPCs within the same network domain.",
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
	client := mustClient(meta)

	awsTgwVpcAttachment := &goaviatrix.AwsTgwVpcAttachment{
		TgwName:                      getString(d, "tgw_name"),
		Region:                       getString(d, "region"),
		VpcAccountName:               getString(d, "vpc_account_name"),
		VpcID:                        getString(d, "vpc_id"),
		CustomizedRoutes:             getString(d, "customized_routes"),
		Subnets:                      getString(d, "subnets"),
		RouteTables:                  getString(d, "route_tables"),
		CustomizedRouteAdvertisement: getString(d, "customized_route_advertisement"),
		DisableLocalRoutePropagation: getBool(d, "disable_local_route_propagation"),
		EdgeAttachment:               getString(d, "edge_attachment"),
		SecurityDomainName:           getString(d, "network_domain_name"),
	}

	isFirewallSecurityDomain, err := client.IsFirewallSecurityDomain(awsTgwVpcAttachment.TgwName, awsTgwVpcAttachment.SecurityDomainName)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("could not find Security Domain: %s", awsTgwVpcAttachment.SecurityDomainName)
		}
		return fmt.Errorf("could not find Security Domain due to: %w", err)
	}

	log.Printf("[INFO] Attaching vpc: %s to tgw %s", awsTgwVpcAttachment.VpcID, awsTgwVpcAttachment.TgwName)

	d.SetId(awsTgwVpcAttachment.TgwName + "~" + awsTgwVpcAttachment.SecurityDomainName + "~" + awsTgwVpcAttachment.VpcID)
	flag := false
	defer func() { _ = resourceAviatrixAwsTgwVpcAttachmentReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	if isFirewallSecurityDomain {
		err = client.CreateAwsTgwVpcAttachmentForFireNet(awsTgwVpcAttachment)
		if err != nil {
			return fmt.Errorf("failed to create Aviatrix Aws Tgw Vpc Attach for FireNet: %w", err)
		}

		if awsTgwVpcAttachment.EdgeAttachment != "" {
			err = client.UpdateFirewallAttachmentAccessFromOnprem(awsTgwVpcAttachment)
			if err != nil {
				return fmt.Errorf("failed to enable firewall attachment access from onprem: %w", err)
			}
		}
	} else {
		if awsTgwVpcAttachment.EdgeAttachment != "" {
			return fmt.Errorf("management access from onprem only works for FireNet")
		}

		err = client.CreateAwsTgwVpcAttachment(awsTgwVpcAttachment)
		if err != nil {
			return fmt.Errorf("failed to create Aviatrix Aws Tgw Vpc Attach: %w", err)
		}
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
	client := mustClient(meta)

	tgwName := getString(d, "tgw_name")
	vpcID := getString(d, "vpc_id")

	if tgwName == "" || getString(d, "network_domain_name") == "" || vpcID == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no vpc names received. Import Id is %s", id)
		mustSet(d, "tgw_name", strings.Split(id, "~")[0])
		mustSet(d, "network_domain_name", strings.Split(id, "~")[1])
		mustSet(d, "vpc_id", strings.Split(id, "~")[2])
	}
	awsTgwVpcAttachment := &goaviatrix.AwsTgwVpcAttachment{
		TgwName:            getString(d, "tgw_name"),
		VpcID:              getString(d, "vpc_id"),
		SecurityDomainName: getString(d, "network_domain_name"),
	}

	aTVA, err := client.GetAwsTgwVpcAttachment(awsTgwVpcAttachment)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to get Aviatrix Aws Tgw Vpc Attach: %w", err)
	}
	if aTVA != nil {
		mustSet(d, "tgw_name", aTVA.TgwName)
		mustSet(d, "region", aTVA.Region)
		mustSet(d, "network_domain_name", aTVA.SecurityDomainName)
		mustSet(d, "vpc_account_name", aTVA.VpcAccountName)
		mustSet(d, "vpc_id", aTVA.VpcID)
		mustSet(d, "disable_local_route_propagation", aTVA.DisableLocalRoutePropagation)

		if getString(d, "subnets") != "" {
			subnetsFromConfigList := strings.Split(getString(d, "subnets"), ",")
			subnetsFromReadList := strings.Split(aTVA.Subnets, ",")
			if len(goaviatrix.Difference(subnetsFromConfigList, subnetsFromReadList)) == 0 ||
				len(goaviatrix.Difference(subnetsFromReadList, subnetsFromConfigList)) == 0 {
				mustSet(d, "subnets", getString(d, "subnets"))
			} else {
				mustSet(d, "subnets", aTVA.Subnets)
			}
		} else {
			mustSet(d, "subnets", aTVA.Subnets)
		}
		if getString(d, "route_tables") != "" {
			routeTablesFromConfigList := strings.Split(getString(d, "route_tables"), ",")
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
				mustSet(d, "route_tables", aTVA.RouteTables)
			} else {
				mustSet(d, "route_tables", getString(d, "route_tables"))
			}
		} else {
			mustSet(d, "route_tables", aTVA.RouteTables)
		}
		mustSet(d, "customized_routes", aTVA.CustomizedRoutes)
		mustSet(d, "customized_route_advertisement", aTVA.CustomizedRouteAdvertisement)
		mustSet(d, "edge_attachment", aTVA.EdgeAttachment)
		d.SetId(awsTgwVpcAttachment.TgwName + "~" + awsTgwVpcAttachment.SecurityDomainName + "~" + awsTgwVpcAttachment.VpcID)
		return nil
	}

	return fmt.Errorf("no Aviatrix Aws Tgw Vpc Attach found")
}

func resourceAviatrixAwsTgwVpcAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	flag := false
	defer func() { _ = resourceAviatrixAwsTgwVpcAttachmentReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	client := mustClient(meta)

	d.Partial(true)
	if d.HasChange("region") {
		return fmt.Errorf("updating region is not allowed")
	}
	if d.HasChange("vpc_account_name") {
		return fmt.Errorf("updating vpc_account_name is not allowed")
	}
	if d.HasChange("customized_routes") {
		awsTgwVpcAttachment := &goaviatrix.AwsTgwVpcAttachment{
			TgwName:          getString(d, "tgw_name"),
			VpcID:            getString(d, "vpc_id"),
			CustomizedRoutes: getString(d, "customized_routes"),
		}
		err := client.EditTgwSpokeVpcCustomizedRoutes(awsTgwVpcAttachment)
		if err != nil {
			return fmt.Errorf("failed to update spoke vpc customized routes: %w", err)
		}
	}
	if d.HasChange("customized_route_advertisement") {
		awsTgwVpcAttachment := &goaviatrix.AwsTgwVpcAttachment{
			TgwName:                      getString(d, "tgw_name"),
			VpcID:                        getString(d, "vpc_id"),
			CustomizedRouteAdvertisement: getString(d, "customized_route_advertisement"),
		}
		err := client.EditTgwSpokeVpcCustomizedRouteAdvertisement(awsTgwVpcAttachment)
		if err != nil {
			return fmt.Errorf("failed to update spoke vpc customized routes advertisement: %w", err)
		}
	}

	if d.HasChange("edge_attachment") {
		awsTgwVpcAttachment := &goaviatrix.AwsTgwVpcAttachment{
			TgwName:            getString(d, "tgw_name"),
			VpcID:              getString(d, "vpc_id"),
			EdgeAttachment:     getString(d, "edge_attachment"),
			SecurityDomainName: getString(d, "network_domain_name"),
		}

		isFirewallSecurityDomain, err := client.IsFirewallSecurityDomain(awsTgwVpcAttachment.TgwName, awsTgwVpcAttachment.SecurityDomainName)
		if err != nil {
			if errors.Is(err, goaviatrix.ErrNotFound) {
				return fmt.Errorf("could not find Network Domain: %s", awsTgwVpcAttachment.SecurityDomainName)
			}
			return fmt.Errorf("could not find Network Domain due to: %w", err)
		}

		oldEA, newEA := d.GetChange("edge_attachment")
		oldEAString := mustString(oldEA)
		newEAString := mustString(newEA)

		if isFirewallSecurityDomain {
			if oldEAString != "" && newEAString != "" {
				awsTgwVpcAttachment.EdgeAttachment = ""

				err := client.UpdateFirewallAttachmentAccessFromOnprem(awsTgwVpcAttachment)
				if err != nil {
					return fmt.Errorf("failed to disable firewall attachment access from onprem while updating: %w", err)
				}

				awsTgwVpcAttachment.EdgeAttachment = newEAString

				err = client.UpdateFirewallAttachmentAccessFromOnprem(awsTgwVpcAttachment)
				if err != nil {
					return fmt.Errorf("failed to enable firewall attachment access from onprem while updating: %w", err)
				}
			} else {
				err := client.UpdateFirewallAttachmentAccessFromOnprem(awsTgwVpcAttachment)
				if err != nil {
					return fmt.Errorf("failed to update firewall attachment access from onprem: %w", err)
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
	client := mustClient(meta)

	awsTgwVpcAttachment := &goaviatrix.AwsTgwVpcAttachment{
		TgwName:            getString(d, "tgw_name"),
		Region:             getString(d, "region"),
		VpcAccountName:     getString(d, "vpc_account_name"),
		VpcID:              getString(d, "vpc_id"),
		SecurityDomainName: getString(d, "network_domain_name"),
	}

	isFirewallSecurityDomain, err := client.IsFirewallSecurityDomain(awsTgwVpcAttachment.TgwName, awsTgwVpcAttachment.SecurityDomainName)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			return fmt.Errorf("could not find Network Domain: %s", awsTgwVpcAttachment.VpcID)
		}
		return fmt.Errorf("could not find Network Domain due to: %w", err)
	}

	if isFirewallSecurityDomain {
		err := client.DeleteAwsTgwVpcAttachmentForFireNet(awsTgwVpcAttachment)
		if err != nil {
			return fmt.Errorf("failed to detach FireNet VPC from TGW: %w", err)
		}
	} else {
		err := client.DeleteAwsTgwVpcAttachment(awsTgwVpcAttachment)
		if err != nil {
			return fmt.Errorf("failed to detach VPC from TGW: %w", err)
		}
	}

	return nil
}
