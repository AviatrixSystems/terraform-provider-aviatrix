package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixAWSTgwDirectConnect() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixAWSTgwDirectConnectCreate,
		Read:   resourceAviatrixAWSTgwDirectConnectRead,
		Update: resourceAviatrixAWSTgwDirectConnectUpdate,
		Delete: resourceAviatrixAWSTgwDirectConnectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
			"security_domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of an Aviatrix security domain, to which the direct connect gateway will be attached.",
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
	client := meta.(*goaviatrix.Client)

	awsTgwDirectConnect := &goaviatrix.AwsTgwDirectConnect{
		TgwName:                  d.Get("tgw_name").(string),
		DirectConnectAccountName: d.Get("directconnect_account_name").(string),
		DxGatewayID:              d.Get("dx_gateway_id").(string),
		SecurityDomainName:       d.Get("security_domain_name").(string),
		AllowedPrefix:            d.Get("allowed_prefix").(string),
	}

	learnedCidrsApproval := d.Get("enable_learned_cidrs_approval").(bool)
	if learnedCidrsApproval {
		awsTgwDirectConnect.LearnedCidrsApproval = "yes"
	} else {
		awsTgwDirectConnect.LearnedCidrsApproval = "no"
	}

	d.SetId(awsTgwDirectConnect.TgwName + "~" + awsTgwDirectConnect.DxGatewayID)
	flag := false
	defer resourceAviatrixAWSTgwDirectConnectReadIfRequired(d, meta, &flag)

	err := client.CreateAwsTgwDirectConnect(awsTgwDirectConnect)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix AWS TGW Direct Connect: %s", err)
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
	client := meta.(*goaviatrix.Client)

	tgwName := d.Get("tgw_name").(string)
	directConnectGatewayID := d.Get("dx_gateway_id").(string)

	if tgwName == "" || directConnectGatewayID == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)
		if !strings.Contains(id, "~") {
			log.Printf("[DEBUG] Import Id: %s is invalid", id)
		}
		d.Set("tgw_name", strings.Split(id, "~")[0])
		d.Set("dx_gateway_id", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	awsTgwDirectConnect := &goaviatrix.AwsTgwDirectConnect{
		TgwName:     d.Get("tgw_name").(string),
		DxGatewayID: d.Get("dx_gateway_id").(string),
	}

	directConnect, err := client.GetAwsTgwDirectConnect(awsTgwDirectConnect)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix Aws Tgw Direct Connect: %s", err)
	}
	log.Printf("[INFO] Found Aviatrix Aws Tgw Direct Connect: %#v", directConnect)

	d.Set("tgw_name", directConnect.TgwName)
	d.Set("directconnect_account_name", directConnect.DirectConnectAccountName)
	d.Set("dx_gateway_id", directConnect.DxGatewayID)
	d.Set("security_domain_name", directConnect.SecurityDomainName)
	d.Set("allowed_prefix", directConnect.AllowedPrefix)
	if directConnect.LearnedCidrsApproval == "yes" {
		d.Set("enable_learned_cidrs_approval", true)
	} else {
		d.Set("enable_learned_cidrs_approval", false)
	}

	d.SetId(directConnect.TgwName + "~" + directConnect.DxGatewayID)
	return nil
}

func resourceAviatrixAWSTgwDirectConnectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	awsTgwDirectConnect := &goaviatrix.AwsTgwDirectConnect{
		TgwName:       d.Get("tgw_name").(string),
		DxGatewayName: d.Get("dx_gateway_id").(string),
	}

	d.Partial(true)

	log.Printf("[INFO] Updating Aviatrix Site2Cloud: %#v", awsTgwDirectConnect)
	if ok := d.HasChange("allowed_prefix"); ok {
		awsTgwDirectConnect.AllowedPrefix = d.Get("allowed_prefix").(string)
		err := client.UpdateDirectConnAllowedPrefix(awsTgwDirectConnect)
		if err != nil {
			return fmt.Errorf("failed to update Aws Tgw Direct Connect Allowed Prefix: %s", err)
		}
	}

	if d.HasChange("enable_learned_cidrs_approval") {
		learnedCidrsApproval := d.Get("enable_learned_cidrs_approval").(bool)
		if learnedCidrsApproval {
			awsTgwDirectConnect.LearnedCidrsApproval = "yes"
			err := client.EnableDirectConnectLearnedCidrsApproval(awsTgwDirectConnect)
			if err != nil {
				return fmt.Errorf("failed to enable learned cidrs approval: %s", err)
			}
		} else {
			awsTgwDirectConnect.LearnedCidrsApproval = "no"
			err := client.DisableDirectConnectLearnedCidrsApproval(awsTgwDirectConnect)
			if err != nil {
				return fmt.Errorf("failed to disable learned cidrs approval: %s", err)
			}
		}
	}

	d.Partial(false)
	return resourceAviatrixAWSTgwDirectConnectRead(d, meta)
}

func resourceAviatrixAWSTgwDirectConnectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	awsTgwDirectConnect := &goaviatrix.AwsTgwDirectConnect{
		TgwName:         d.Get("tgw_name").(string),
		DirectConnectID: d.Get("dx_gateway_id").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix AWS TGW Direct Connect: %#v", awsTgwDirectConnect)

	err := client.DeleteAwsTgwDirectConnect(awsTgwDirectConnect)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix AWS TGW Direct Connect: %s", err)
	}

	return nil
}
