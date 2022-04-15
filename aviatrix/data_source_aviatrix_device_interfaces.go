package aviatrix

import (
	"fmt"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAviatrixDeviceInterfaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAviatrixDeviceInterfaceConfigRead,

		Schema: map[string]*schema.Schema{
			"device_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of device.",
			},
			"wan_interfaces": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of WAN primary interface and WAN primary interface public IP.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"wan_primary_interface": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "WAN primary interface of the device.",
						},
						"wan_primary_interface_public_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "WAN primary interface public IP address.",
						},
					},
				},
			},
		},
	}
}

func dataSourceAviatrixDeviceInterfaceConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	deviceName := d.Get("device_name").(string)

	deviceWanInterfaces, err := client.GetDeviceInterfaces(deviceName)
	if err != nil {
		return fmt.Errorf("couldn't get device wan interfaces: %s", err)
	}

	var wanInterfaces []map[string]interface{}
	for _, wanInterface := range *deviceWanInterfaces {
		wI := make(map[string]interface{})
		wI["wan_primary_interface"] = wanInterface.Interface
		wI["wan_primary_interface_public_ip"] = wanInterface.IP
		wanInterfaces = append(wanInterfaces, wI)
	}

	if err = d.Set("wan_interfaces", wanInterfaces); err != nil {
		return fmt.Errorf("couldn't set wan_interfaces: %s", err)
	}

	d.SetId(deviceName)
	return nil
}
