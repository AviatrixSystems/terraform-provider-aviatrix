package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixAzureVngConn() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixAzureVngConnCreate,
		Read:   resourceAviatrixAzureVngConnRead,
		Delete: resourceAviatrixAzureVngConnDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"primary_gateway_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Primary gateway name",
			},
			"connection_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Connection name",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VPC ID",
			},
			"vng_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VNG name",
			},
			"attached": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "VNG attached or not",
			},
		},
	}
}

func marshalAzureVngConnInput(d *schema.ResourceData) *goaviatrix.AzureVngConn {
	return &goaviatrix.AzureVngConn{
		PrimaryGatewayName: d.Get("primary_gateway_name").(string),
		ConnectionName:     d.Get("connection_name").(string),
	}
}

func resourceAviatrixAzureVngConnCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	azureVngConn := marshalAzureVngConnInput(d)

	if err := client.ConnectAzureVng(azureVngConn); err != nil {
		return fmt.Errorf("could not connect to azure vng: %v", err)
	}

	d.SetId(azureVngConn.ConnectionName)
	return nil
}

func resourceAviatrixAzureVngConnRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	connectionName := d.Get("connection_name").(string)

	if connectionName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)

		d.Set("connection_name", id)
		connectionName = id
	}

	azureVngConnStatus, err := client.GetAzureVngConnStatus(connectionName)
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not get azure vng conn status: %v", err)
	}

	d.Set("primary_gateway_name", azureVngConnStatus.PrimaryGatewayName)
	d.Set("vpc_id", azureVngConnStatus.VpcId)
	d.Set("vng_name", azureVngConnStatus.VngName)
	d.Set("attached", azureVngConnStatus.Attached)

	d.SetId(connectionName)
	return nil
}

func resourceAviatrixAzureVngConnDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	vpcId := d.Get("vpc_id").(string)
	connectionName := d.Get("connection_name").(string)

	if err := client.DisconnectAzureVng(vpcId, connectionName); err != nil {
		return fmt.Errorf("could not disconnect vng connection: %v", err)
	}

	return nil
}
