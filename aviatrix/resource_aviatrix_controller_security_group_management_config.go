package aviatrix

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
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
		},
	}
}

func resourceAviatrixControllerSecurityGroupManagementConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	account := d.Get("account_name").(string)
	enableSecurityGroupManagement := d.Get("enable_security_group_management").(bool)

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
				return fmt.Errorf("failed to enable controller Security Group Management: %s", err)
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
				return fmt.Errorf("failed to disable controller Security Group Management: %s", err)
			}
		}
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceAviatrixControllerSecurityGroupManagementConfigRead(d, meta)
}

func resourceAviatrixControllerSecurityGroupManagementConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	if d.Id() != strings.Replace(client.ControllerIP, ".", "-", -1) {
		return fmt.Errorf("ID: %s does not match controller IP. Please provide correct ID for importing", d.Id())
	}

	sgm, err := client.GetSecurityGroupManagementStatus()
	if err != nil {
		return fmt.Errorf("could not read Aviatrix Controller Security Group Management Status: %s", err)
	}
	if sgm != nil {
		d.Set("enable_security_group_management", sgm.State == "Enabled")
		d.Set("account_name", sgm.AccountName)
	} else {
		return fmt.Errorf("could not read Aviatrix Controller Security Group Management Status")
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceAviatrixControllerSecurityGroupManagementConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	if d.HasChange("account_name") || d.HasChange("enable_security_group_management") {
		oldAccount, newAccount := d.GetChange("account_name")
		securityGroupManagement := d.Get("enable_security_group_management").(bool)

		if oldAccount.(string) != "" && newAccount.(string) != "" && securityGroupManagement {
			err := client.DisableSecurityGroupManagement()
			if err != nil {
				if err != nil {
					return fmt.Errorf("failed to disable Security Group Management on controller %s: %s", d.Id(), err)
				}
			}
			err = client.EnableSecurityGroupManagement(newAccount.(string))
			if err != nil {
				return fmt.Errorf("failed to enable Security Group Management on controller %s: %s", d.Id(), err)
			}
		} else {
			return resourceAviatrixControllerSecurityGroupManagementConfigCreate(d, meta)
		}
	}

	return resourceAviatrixControllerSecurityGroupManagementConfigRead(d, meta)
}

func resourceAviatrixControllerSecurityGroupManagementConfigDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
