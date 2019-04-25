package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceControllerConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceControllerConfigCreate,
		Read:   resourceControllerConfigRead,
		Update: resourceControllerConfigUpdate,
		Delete: resourceControllerConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"http_access": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Switch for http access. Default: false.",
			},
			"fqdn_exception_rule": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "A system-wide mode. Default: 'true'.",
			},
		},
	}
}

func resourceControllerConfigCreate(d *schema.ResourceData, meta interface{}) error {
	var err error

	client := meta.(*goaviatrix.Client)

	log.Printf("[INFO] Configuring Aviatrix controller : %#v", d)

	if httpAccess := d.Get("http_access").(bool); httpAccess {
		curStatus, _ := client.GetHttpAccessEnabled()
		if curStatus == "True" {
			log.Printf("[INFO] Http Access is already enabled")
		} else {
			err = client.EnableHttpAccess()
		}
	} else {
		curStatus, _ := client.GetHttpAccessEnabled()
		if curStatus == "False" {
			log.Printf("[INFO] Http Access is already disabled")
		} else {
			err = client.DisableHttpAccess()
		}
	}
	if err != nil {
		return fmt.Errorf("failed to configure controller http access: %s", err)
	}

	if fqdnExceptionRule := d.Get("fqdn_exception_rule").(bool); fqdnExceptionRule {
		curStatus, _ := client.GetExceptionRuleStatus()
		if curStatus {
			log.Printf("[INFO] FQDN Exception Rule is already enabled")
		} else {
			err = client.EnableExceptionRule()
		}
	} else {
		curStatus, _ := client.GetExceptionRuleStatus()
		if !curStatus {
			log.Printf("[INFO] FQDN Exception Rule is already disabled")
		} else {
			err = client.DisableExceptionRule()
		}
	}
	if err != nil {
		return fmt.Errorf("failed to configure controller exception rule: %s", err)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return resourceControllerConfigRead(d, meta)
}

func resourceControllerConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	log.Printf("[INFO] Getting controller %s configuration", d.Id())
	result, err := client.GetHttpAccessEnabled()
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read Aviatrix Controller Config: %s", err)
	}

	if result[1:5] == "True" {
		d.Set("http_access", true)
	} else {
		d.Set("http_access", false)
	}

	res, err := client.GetExceptionRuleStatus()
	if err != nil {
		return fmt.Errorf("could not read Aviatrix Controller Exception Rule Status: %s", err)
	}
	if res {
		d.Set("fqdn_exception_rule", true)
	} else {
		d.Set("fqdn_exception_rule", false)
	}

	d.SetId(strings.Replace(client.ControllerIP, ".", "-", -1))
	return nil
}

func resourceControllerConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	log.Printf("[INFO] Updating Controller configuration: %#v", d)
	d.Partial(true)

	if d.HasChange("http_access") {
		if httpAccess := d.Get("http_access").(bool); httpAccess {
			err := client.EnableHttpAccess()
			if err != nil {
				log.Printf("[ERROR] Failed to enable http access on controller %s", d.Id())
				return err
			}
		} else {
			err := client.DisableHttpAccess()
			if err != nil {
				log.Printf("[ERROR] Failed to disable http access on controller %s", d.Id())
				return err
			}
		}
		d.SetPartial("http_access")
	}

	if d.HasChange("fqdn_exception_rule") {
		if httpAccess := d.Get("fqdn_exception_rule").(bool); httpAccess {
			err := client.EnableExceptionRule()
			if err != nil {
				log.Printf("[ERROR] Failed to enable exception rule on controller %s", d.Id())
				return err
			}
		} else {
			err := client.DisableExceptionRule()
			if err != nil {
				log.Printf("[ERROR] Failed to disable exception rule on controller %s", d.Id())
				return err
			}
		}
		d.SetPartial("fqdn_exception_rule")
	}

	d.Partial(false)
	return nil
}

func resourceControllerConfigDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
