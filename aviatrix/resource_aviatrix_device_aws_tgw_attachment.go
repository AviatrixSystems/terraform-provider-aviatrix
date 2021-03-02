package aviatrix

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixDeviceAwsTgwAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixDeviceAwsTgwAttachmentCreate,
		Read:   resourceAviatrixDeviceAwsTgwAttachmentRead,
		Delete: resourceAviatrixDeviceAwsTgwAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"connection_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Connection name.",
			},
			"device_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Device name.",
			},
			"aws_tgw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "AWS TGW name.",
			},
			"device_bgp_asn": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Device BGP AS Number.",
			},
			"security_domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Security domain name.",
			},
		},
	}
}

func marshalDeviceAwsTgwAttachmentInput(d *schema.ResourceData) *goaviatrix.DeviceAwsTgwAttachment {
	return &goaviatrix.DeviceAwsTgwAttachment{
		ConnectionName:     d.Get("connection_name").(string),
		DeviceName:         d.Get("device_name").(string),
		AwsTgwName:         d.Get("aws_tgw_name").(string),
		DeviceAsn:          strconv.Itoa(d.Get("device_bgp_asn").(int)),
		SecurityDomainName: d.Get("security_domain_name").(string),
	}
}

func resourceAviatrixDeviceAwsTgwAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	attachment := marshalDeviceAwsTgwAttachmentInput(d)

	if err := client.CreateDeviceAwsTgwAttachment(attachment); err != nil {
		return fmt.Errorf("could not create aws tgw attachment: %v", err)
	}

	d.SetId(attachment.ID())
	return nil
}

func resourceAviatrixDeviceAwsTgwAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	connectionName := d.Get("connection_name").(string)
	deviceName := d.Get("device_name").(string)
	tgwName := d.Get("aws_tgw_name").(string)
	if connectionName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no device_aws_tgw_attachment connection_name received. Import Id is %s", id)
		d.SetId(id)
		parts := strings.Split(id, "~")
		if len(parts) != 3 {
			return fmt.Errorf("import id is invalid, expecting connection_name~device_name~aws_tgw_name: %s", id)
		}
		connectionName = parts[0]
		deviceName = parts[1]
		tgwName = parts[2]
	}

	attachment := &goaviatrix.DeviceAwsTgwAttachment{
		ConnectionName: connectionName,
		DeviceName:     deviceName,
		AwsTgwName:     tgwName,
	}

	attachment, err := client.GetDeviceAwsTgwAttachment(attachment)
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not find device_aws_tgw_attachment %s: %v", connectionName, err)
	}

	d.Set("connection_name", attachment.ConnectionName)
	d.Set("device_name", attachment.DeviceName)
	d.Set("aws_tgw_name", attachment.AwsTgwName)
	d.Set("security_domain_name", attachment.SecurityDomainName)

	asn, err := strconv.Atoi(attachment.DeviceAsn)
	if err != nil {
		return fmt.Errorf("could not convert DeviceAsn to int: %v", err)
	}
	d.Set("device_bgp_asn", asn)

	d.SetId(attachment.ID())
	return nil
}

func resourceAviatrixDeviceAwsTgwAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	cn := d.Get("connection_name").(string)

	if err := client.DeleteDeviceAttachment(cn); err != nil {
		return err
	}

	return nil
}
