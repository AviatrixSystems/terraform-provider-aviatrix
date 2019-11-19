package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
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
		},
	}
}

func resourceAviatrixAWSTgwDirectConnectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	awsTgwDirectConn := &goaviatrix.AwsTgwDirectConn{
		TgwName:                  d.Get("tgw_name").(string),
		DirectConnectAccountName: d.Get("directconnect_account_name").(string),
		DxGatewayID:              d.Get("dx_gateway_id").(string),
		SecurityDomainName:       d.Get("security_domain_name").(string),
		AllowedPrefix:            d.Get("allowed_prefix").(string),
	}

	err := client.CreateAwsTgwDirectConn(awsTgwDirectConn)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix AWS TGW Direct Connection: %s", err)
	}

	d.SetId(awsTgwDirectConn.TgwName + "~" + awsTgwDirectConn.DxGatewayID)
	return resourceAviatrixAwsTgwVpnConnRead(d, meta)
}

func resourceAviatrixAWSTgwDirectConnectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	tgwName := d.Get("tgw_name").(string)
	directConnGwID := d.Get("dx_gateway_id").(string)

	if tgwName == "" || directConnGwID == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)
		if !strings.Contains(id, "~") {
			log.Printf("[DEBUG] Import Id: %s is invalid", id)
		}
		d.Set("tgw_name", strings.Split(id, "~")[0])
		d.Set("dx_gateway_id", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	awsTgwDirectConn := &goaviatrix.AwsTgwDirectConn{
		TgwName:     d.Get("tgw_name").(string),
		DxGatewayID: d.Get("dx_gateway_id").(string),
	}

	directConn, err := client.GetAwsTgwDirectConn(awsTgwDirectConn)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find Aviatrix Aws Tgw Direct Connection: %s", err)
	}
	log.Printf("[INFO] Found Aviatrix Aws Tgw Direct Connection: %#v", directConn)

	d.Set("tgw_name", directConn.TgwName)
	d.Set("directconnect_account_name", directConn.DirectConnectAccountName)
	d.Set("dx_gateway_id", directConn.DxGatewayID)
	d.Set("security_domain_name", directConn.SecurityDomainName)
	d.Set("allowed_prefix", directConn.AllowedPrefix)
	d.SetId(directConn.TgwName + "~" + directConn.DxGatewayID)

	return nil
}

func resourceAviatrixAWSTgwDirectConnectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	awsTgwDirectConn := &goaviatrix.AwsTgwDirectConn{
		TgwName:       d.Get("tgw_name").(string),
		DxGatewayName: d.Get("dx_gateway_id").(string),
	}

	d.Partial(true)

	log.Printf("[INFO] Updating Aviatrix Site2Cloud: %#v", awsTgwDirectConn)
	if ok := d.HasChange("allowed_prefix"); ok {
		awsTgwDirectConn.AllowedPrefix = d.Get("allowed_prefix").(string)
		err := client.UpdateDirectConnAllowedPrefix(awsTgwDirectConn)
		if err != nil {
			return fmt.Errorf("failed to update Aws Tgw Direct Conn Allowed Prefix: %s", err)
		}
		d.SetPartial("allowed_prefix")
	}

	return resourceAviatrixAWSTgwDirectConnectRead(d, meta)
}

func resourceAviatrixAWSTgwDirectConnectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	awsTgwDirectConn := &goaviatrix.AwsTgwDirectConn{
		TgwName:         d.Get("tgw_name").(string),
		DirectConnectID: d.Get("dx_gateway_id").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix AWS TGW Direct Conn: %#v", awsTgwDirectConn)

	err := client.DeleteAwsTgwDirectConn(awsTgwDirectConn)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix AWS TGW Direct Conn: %s", err)
	}

	return nil
}
