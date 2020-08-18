package aviatrix

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixDeviceVirtualWanAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixDeviceVirtualWanAttachmentCreate,
		Read:   resourceAviatrixDeviceVirtualWanAttachmentRead,
		Delete: resourceAviatrixDeviceVirtualWanAttachmentDelete,
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
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Azure access account name.",
			},
			"resource_group": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ARM resource group name.",
			},
			"hub_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Virtual WAN vhub name.",
			},
			"device_bgp_asn": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Device AS Number.",
			},
		},
	}
}

func marshalDeviceVirtualWanAttachmentInput(d *schema.ResourceData) *goaviatrix.DeviceVirtualWanAttachment {
	return &goaviatrix.DeviceVirtualWanAttachment{
		ConnectionName: d.Get("connection_name").(string),
		DeviceName:     d.Get("device_name").(string),
		AccountName:    d.Get("account_name").(string),
		ResourceGroup:  d.Get("resource_group").(string),
		HubName:        d.Get("hub_name").(string),
		DeviceAsn:      strconv.Itoa(d.Get("device_bgp_asn").(int)),
	}
}

func resourceAviatrixDeviceVirtualWanAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	attachment := marshalDeviceVirtualWanAttachmentInput(d)

	if err := client.CreateDeviceVirtualWanAttachment(attachment); err != nil {
		return fmt.Errorf("could not create virtual wan and device attachment: %v", err)
	}

	d.SetId(attachment.ConnectionName)
	return nil
}

func resourceAviatrixDeviceVirtualWanAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	connectionName := d.Get("connection_name").(string)
	if connectionName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no device_virtual_wan_attachment connectionName received. Import Id is %s", id)
		d.SetId(id)
		connectionName = id
	}

	attachment := &goaviatrix.DeviceVirtualWanAttachment{
		ConnectionName: connectionName,
	}

	attachment, err := client.GetDeviceVirtualWanAttachment(attachment)
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not find device_virtual_wan_attachment %s: %v", connectionName, err)
	}

	d.Set("connection_name", attachment.ConnectionName)
	d.Set("device_name", attachment.DeviceName)
	d.Set("account_name", attachment.AccountName)
	d.Set("resource_group", attachment.ResourceGroup)
	d.Set("hub_name", attachment.HubName)

	deviceAsn, err := strconv.Atoi(attachment.DeviceAsn)
	if err != nil {
		return fmt.Errorf("could not covert device asn to int: %v", err)
	}
	d.Set("device_bgp_asn", deviceAsn)

	d.SetId(attachment.ConnectionName)
	return nil
}

func resourceAviatrixDeviceVirtualWanAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	attachment := marshalDeviceVirtualWanAttachmentInput(d)

	if err := client.DeleteDeviceAttachment(attachment.ConnectionName); err != nil {
		return fmt.Errorf("could not delete virtual wan and device attachment: %v", err)
	}

	return nil
}
