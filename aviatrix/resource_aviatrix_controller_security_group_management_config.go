package aviatrix

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixControllerSecurityGroupManagementConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixControllerSecurityGroupManagementConfigCreate,
		Read:   resourceAviatrixControllerSecurityGroupManagementConfigRead,
		Update: resourceAviatrixControllerSecurityGroupManagementConfigUpdate,
		Delete: resourceAviatrixControllerSecurityGroupManagementConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cloud account name of user.",
			},
			"enable_security_group_management": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Used to manage the Controller instance’s inbound rules from gateways.",
			},
			"gateway_egress_cidrs": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of egress CIDRs for private_network gateways that reach the controller through NAT gateways or firewalls.",
			},
		},
	}
}

func resourceAviatrixControllerSecurityGroupManagementConfigCreate(d *schema.ResourceData, meta any) error {
	client := mustClient(meta)

	account := getString(d, "account_name")
	enableSecurityGroupManagement := getBool(d, "enable_security_group_management")

	if enableSecurityGroupManagement {
		if account == "" {
			return fmt.Errorf("account_name is needed to enable controller Security Group Management")
		}
		curStatus, _ := client.GetSecurityGroupManagementStatus()
		if curStatus.State == "Enabled" {
			log.Printf("[INFO] Security Group Management is already enabled")
		} else {
			err := client.EnableSecurityGroupManagement(account)
			if err != nil {
				return fmt.Errorf("failed to enable controller Security Group Management: %w", err)
			}
		}
	} else {
		if account != "" {
			return fmt.Errorf("account_name isn't needed to disable controller Security Group Management")
		}
		curStatus, _ := client.GetSecurityGroupManagementStatus()
		if curStatus.State == "Disabled" {
			log.Printf("[INFO] Security Group Management is already disabled")
		} else {
			err := client.DisableSecurityGroupManagement()
			if err != nil {
				return fmt.Errorf("failed to disable controller Security Group Management: %w", err)
			}
		}
	}

	if enableSecurityGroupManagement {
		egressCidrs := getStringList(d, "gateway_egress_cidrs")
		if len(egressCidrs) > 0 {
			err := client.UpdateSecurityGroupGatewayEgressCidrs(strings.Join(egressCidrs, ","))
			if err != nil {
				return fmt.Errorf("failed to update gateway egress CIDRs: %w", err)
			}
		}
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixControllerSecurityGroupManagementConfigRead(d, meta)
}

func resourceAviatrixControllerSecurityGroupManagementConfigRead(d *schema.ResourceData, meta any) error {
	client := mustClient(meta)

	sgm, err := client.GetSecurityGroupManagementStatus()
	if err != nil {
		return fmt.Errorf("could not read Aviatrix Controller Security Group Management Status: %w", err)
	}
	if sgm != nil {
		mustSet(d, "enable_security_group_management", sgm.State == "Enabled")
		mustSet(d, "account_name", sgm.AccountName)
		if sgm.GatewayEgressCidrs != nil {
			mustSet(d, "gateway_egress_cidrs", sgm.GatewayEgressCidrs)
		}
	} else {
		return fmt.Errorf("could not read Aviatrix Controller Security Group Management Status")
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixControllerSecurityGroupManagementConfigUpdate(d *schema.ResourceData, meta any) error {
	client := mustClient(meta)

	if d.HasChange("account_name") || d.HasChange("enable_security_group_management") {
		oldAccount, newAccount := d.GetChange("account_name")
		securityGroupManagement := getBool(d, "enable_security_group_management")

		if mustString(oldAccount) != "" && mustString(newAccount) != "" && securityGroupManagement {
			err := client.DisableSecurityGroupManagement()
			if err != nil {
				return fmt.Errorf("failed to disable Security Group Management on controller %s: %w", d.Id(), err)
			}
			err = client.EnableSecurityGroupManagement(mustString(newAccount))
			if err != nil {
				return fmt.Errorf("failed to enable Security Group Management on controller %s: %w", d.Id(), err)
			}
		} else {
			return resourceAviatrixControllerSecurityGroupManagementConfigCreate(d, meta)
		}
	}

	if d.HasChange("gateway_egress_cidrs") {
		egressCidrs := getStringList(d, "gateway_egress_cidrs")
		cidrsStr := strings.Join(egressCidrs, ",")
		err := client.UpdateSecurityGroupGatewayEgressCidrs(cidrsStr)
		if err != nil {
			return fmt.Errorf("failed to update gateway egress CIDRs on controller %s: %w", d.Id(), err)
		}
	}

	return resourceAviatrixControllerSecurityGroupManagementConfigRead(d, meta)
}

func resourceAviatrixControllerSecurityGroupManagementConfigDelete(d *schema.ResourceData, meta any) error {
	return nil
}
