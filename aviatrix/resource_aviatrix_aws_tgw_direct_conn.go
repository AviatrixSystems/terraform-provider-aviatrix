package aviatrix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixAWSTgwDirectConn() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixAWSTgwDirectConnCreate,
		Read:   resourceAviatrixAWSTgwDirectConnRead,
		Update: resourceAviatrixAWSTgwDirectConnUpdate,
		Delete: resourceAviatrixAWSTgwDirectConnDelete,
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
			"direct_conn_account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "",
			},
			"direct_conn_gw_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "",
			},
			"route_domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "AWS side as a number. Integer between 1-65535. Example: '12'.",
			},
			"allowed_prefix": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Public IP address. Example: '40.0.0.0'.",
			},
		},
	}
}

func resourceAviatrixAWSTgwDirectConnCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	awsTgwDirectConn := &goaviatrix.AwsTgwDirectConn{
		TgwName:               d.Get("tgw_name").(string),
		DirectConnAccountName: d.Get("direct_conn_account_name").(string),
		DirectConnGwID:        d.Get("direct_conn_gw_id").(string),
		RouteDomainName:       d.Get("route_domain_name").(string),
		AllowedPrefix:         d.Get("allowed_prefix").(string),
	}

	_, err := client.CreateAwsTgwDirectConn(awsTgwDirectConn)
	if err != nil {
		return fmt.Errorf("failed to create Aviatrix AWS TGW Direct Connection: %s", err)
	}

	d.SetId(awsTgwDirectConn.TgwName + "~" + awsTgwDirectConn.DirectConnGwID)
	return nil
	return resourceAviatrixAwsTgwVpnConnRead(d, meta)
}

func resourceAviatrixAWSTgwDirectConnRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	tgwName := d.Get("tgw_name").(string)
	directConnGwID := d.Get("direct_conn_gw_id").(string)

	if tgwName == "" || directConnGwID == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)
		if !strings.Contains(id, "~") {
			log.Printf("[DEBUG] Import Id: %s is invalid", id)
		}
		d.Set("tgw_name", strings.Split(id, "~")[0])
		d.Set("direct_conn_gw_id", strings.Split(id, "~")[1])
		d.SetId(id)
	}

	awsTgwDirectConn := &goaviatrix.AwsTgwDirectConn{
		TgwName:        d.Get("tgw_name").(string),
		DirectConnGwID: d.Get("direct_conn_gw_id").(string),
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
	d.Set("direct_conn_account_name", directConn.DirectConnAccountName)
	d.Set("direct_conn_gw_id", directConn.DirectConnGwID)
	d.Set("route_domain_name", directConn.RouteDomainName)
	d.Set("allowed_prefix", directConn.AllowedPrefix)
	d.SetId(directConn.TgwName + "~" + directConn.DirectConnGwID)

	return nil
}

func resourceAviatrixAWSTgwDirectConnUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	awsTgwDirectConn := &goaviatrix.AwsTgwDirectConn{
		TgwName:          d.Get("tgw_name").(string),
		DirectConnGwName: d.Get("direct_conn_gw_id").(string),
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

	return resourceAviatrixAWSTgwDirectConnRead(d, meta)
}

func resourceAviatrixAWSTgwDirectConnDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	awsTgwDirectConn := &goaviatrix.AwsTgwDirectConn{
		TgwName:          d.Get("tgw_name").(string),
		DirectConnGwName: d.Get("direct_conn_gw_id").(string),
	}

	log.Printf("[INFO] Deleting Aviatrix AWS TGW Direct Conn: %#v", awsTgwDirectConn)

	err := client.DeleteAwsTgwDirectConn(awsTgwDirectConn)
	if err != nil {
		return fmt.Errorf("failed to delete Aviatrix AWS TGW Direct Conn: %s", err)
	}

	return nil
}
