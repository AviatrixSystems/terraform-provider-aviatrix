package aviatrix

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixAWSTgwDirectConnect() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixAWSTgwDirectConnectCreate,
		Read:   resourceAviatrixAWSTgwDirectConnectRead,
		Update: resourceAviatrixAWSTgwDirectConnectUpdate,
		Delete: resourceAviatrixAWSTgwDirectConnectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
		},

		Schema: map[string]*schema.Schema{
			"tgw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "This parameter represents the name of an AWS TGW.",
			},
			"directconnect_account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "This parameter represents the name of an Account in Aviatrix controller.",
			},
			"dx_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "This parameter represents the name of a Direct Connect Gateway ID.",
			},
			"network_domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of an Aviatrix network domain, to which the direct connect gateway will be attached.",
			},
			"allowed_prefix": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Public IP address. Example: '40.0.0.0'.",
			},
			"enable_learned_cidrs_approval": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Switch to enable/disable encrypted transit approval for direct connection. Valid values: true, false.",
			},
		},
	}
}

func resourceAviatrixAWSTgwDirectConnectCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	awsTgwDirectConnect := &goaviatrix.AwsTgwDirectConnect{
		TgwName:                  getString(d, "tgw_name"),
		DirectConnectAccountName: getString(d, "directconnect_account_name"),
		DxGatewayID:              getString(d, "dx_gateway_id"),
		AllowedPrefix:            getString(d, "allowed_prefix"),
		SecurityDomainName:       getString(d, "network_domain_name"),
	}

	learnedCidrsApproval := getBool(d, "enable_learned_cidrs_approval")
	if learnedCidrsApproval {
		awsTgwDirectConnect.LearnedCidrsApproval = "yes"
	} else {
		awsTgwDirectConnect.LearnedCidrsApproval = "no"
	}

	d.SetId(awsTgwDirectConnect.TgwName + "~" + awsTgwDirectConnect.DxGatewayID)
	flag := false
	defer func() { _ = resourceAviatrixAWSTgwDirectConnectReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	err := client.CreateAwsTgwDirectConnect(awsTgwDirectConnect)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix AWS TGW Direct Connect: %w", err)
	}

	return resourceAviatrixAWSTgwDirectConnectReadIfRequired(d, meta, &flag)
}

func resourceAviatrixAWSTgwDirectConnectReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixAWSTgwDirectConnectRead(d, meta)
	}
	return nil
}

func resourceAviatrixAWSTgwDirectConnectRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	tgwName := getString(d, "tgw_name")
	directConnectGatewayID := getString(d, "dx_gateway_id")

	if tgwName == "" || directConnectGatewayID == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)
		if !strings.Contains(id, "~") {
			log.Printf("[DEBUG] Import Id: %s is invalid", id)
		}
		mustSet(d, "tgw_name", strings.Split(id, "~")[0])
		mustSet(d, "dx_gateway_id", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	awsTgwDirectConnect := &goaviatrix.AwsTgwDirectConnect{
		TgwName:     getString(d, "tgw_name"),
		DxGatewayID: getString(d, "dx_gateway_id"),
	}

	directConnect, err := client.GetAwsTgwDirectConnect(awsTgwDirectConnect)
	if err != nil {
		if errors.Is(err, goaviatrix.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix Aws Tgw Direct Connect: %w", err)
	}
	log.Printf("[INFO] Found Aviatrix Aws Tgw Direct Connect: %#v", directConnect)
	mustSet(d, "tgw_name", directConnect.TgwName)
	mustSet(d, "directconnect_account_name", directConnect.DirectConnectAccountName)
	mustSet(d, "dx_gateway_id", directConnect.DxGatewayID)
	mustSet(d, "allowed_prefix", directConnect.AllowedPrefix)
	mustSet(d, "network_domain_name", directConnect.SecurityDomainName)

	if directConnect.LearnedCidrsApproval == "yes" {
		mustSet(d, "enable_learned_cidrs_approval", true)
	} else {
		mustSet(d, "enable_learned_cidrs_approval", false)
	}

	d.SetId(directConnect.TgwName + "~" + directConnect.DxGatewayID)
	return nil
}

func resourceAviatrixAWSTgwDirectConnectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	awsTgwDirectConnect := &goaviatrix.AwsTgwDirectConnect{
		TgwName:       getString(d, "tgw_name"),
		DxGatewayName: getString(d, "dx_gateway_id"),
	}

	d.Partial(true)

	log.Printf("[INFO] Updating Aviatrix Site2Cloud: %#v", awsTgwDirectConnect)
	if ok := d.HasChange("allowed_prefix"); ok {
		awsTgwDirectConnect.AllowedPrefix = getString(d, "allowed_prefix")
		err := client.UpdateDirectConnAllowedPrefix(awsTgwDirectConnect)
		if err != nil {
			return fmt.Errorf("failed to update Aws Tgw Direct Connect Allowed Prefix: %w", err)
		}
	}

	if d.HasChange("enable_learned_cidrs_approval") {
		learnedCidrsApproval := getBool(d, "enable_learned_cidrs_approval")
		if learnedCidrsApproval {
			awsTgwDirectConnect.LearnedCidrsApproval = "yes"
			err := client.EnableDirectConnectLearnedCidrsApproval(awsTgwDirectConnect)
			if err != nil {
				return fmt.Errorf("failed to enable learned cidrs approval: %w", err)
			}
		} else {
			awsTgwDirectConnect.LearnedCidrsApproval = "no"
			err := client.DisableDirectConnectLearnedCidrsApproval(awsTgwDirectConnect)
			if err != nil {
				return fmt.Errorf("failed to disable learned cidrs approval: %w", err)
			}
		}
	}

	d.Partial(false)
	return resourceAviatrixAWSTgwDirectConnectRead(d, meta)
}

func resourceAviatrixAWSTgwDirectConnectDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)
	awsTgwDirectConnect := &goaviatrix.AwsTgwDirectConnect{
		TgwName:         getString(d, "tgw_name"),
		DirectConnectID: getString(d, "dx_gateway_id"),
	}

	log.Printf("[INFO] Deleting Aviatrix AWS TGW Direct Connect: %#v", awsTgwDirectConnect)

	err := client.DeleteAwsTgwDirectConnect(awsTgwDirectConnect)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix AWS TGW Direct Connect: %w", err)
	}

	return nil
}
