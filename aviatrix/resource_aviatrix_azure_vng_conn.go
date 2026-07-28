package aviatrix

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixAzureVngConn() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixAzureVngConnCreate,
		Read:   resourceAviatrixAzureVngConnRead,
		Delete: resourceAviatrixAzureVngConnDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough, //nolint:staticcheck // SA1019: deprecated but requires structural changes to migrate,
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
		PrimaryGatewayName: getString(d, "primary_gateway_name"),
		ConnectionName:     getString(d, "connection_name"),
	}
}

func resourceAviatrixAzureVngConnCreate(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	azureVngConn := marshalAzureVngConnInput(d)

	d.SetId(azureVngConn.ConnectionName)
	flag := false
	defer func() { _ = resourceAviatrixAzureVngConnReadIfRequired(d, meta, &flag) }() //nolint:errcheck // read on deferred path

	if err := client.ConnectAzureVng(azureVngConn); err != nil {
		return fmt.Errorf("could not connect to azure vng: %w", err)
	}

	return resourceAviatrixAzureVngConnReadIfRequired(d, meta, &flag)
}

func resourceAviatrixAzureVngConnReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixAzureVngConnRead(d, meta)
	}
	return nil
}

func resourceAviatrixAzureVngConnRead(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	connectionName := getString(d, "connection_name")

	if connectionName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)
		mustSet(d, "connection_name", id)
		connectionName = id
	}

	azureVngConnStatus, err := client.GetAzureVngConnStatus(connectionName)
	if errors.Is(err, goaviatrix.ErrNotFound) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not get azure vng conn status: %w", err)
	}
	mustSet(d, "primary_gateway_name", azureVngConnStatus.PrimaryGatewayName)
	mustSet(d, "vpc_id", azureVngConnStatus.VpcId)
	mustSet(d, "vng_name", azureVngConnStatus.VngName)
	mustSet(d, "attached", azureVngConnStatus.Attached)

	d.SetId(connectionName)
	return nil
}

func resourceAviatrixAzureVngConnDelete(d *schema.ResourceData, meta interface{}) error {
	client := mustClient(meta)

	vpcId := getString(d, "vpc_id")
	connectionName := getString(d, "connection_name")

	if err := client.DisconnectAzureVng(vpcId, connectionName); err != nil {
		return fmt.Errorf("could not disconnect vng connection: %w", err)
	}

	return nil
}
