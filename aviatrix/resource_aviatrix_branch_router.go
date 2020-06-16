package aviatrix

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixBranchRouter() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixBranchRouterCreate,
		Read:   resourceAviatrixBranchRouterRead,
		Update: resourceAviatrixBranchRouterUpdate,
		Delete: resourceAviatrixBranchRouterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the router.",
			},
			"public_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPAddress,
				Description:  "Public IP address of the router.",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username to use to connect to the router.",
			},
			"key_file": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"password", "key_file"},
				Description:  "Path to private key file.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: envDefaultFunc("AVIATRIX_BRANCH_ROUTER_PASSWORD"),
				Description: "Password to connect to the router. " +
					"This attribute can also be set via environment variable 'AVIATRIX_BRANCH_ROUTER_PASSWORD'. " +
					"If both are set the value in the config file will be used.",
			},
			"wan_primary_interface": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Primary WAN interface of the branch router. For example, 'GigabitEthernet1'.",
			},
			"wan_primary_interface_public_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsIPAddress,
				Description:  "Primary WAN interface public IP address.",
			},
			"wan_backup_interface": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Backup WAN interface of the branch router. For example, 'GigabitEthernet2'.",
			},
			"wan_backup_interface_public_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPAddress,
				Description:  "Backup WAN interface public IP address.",
			},
			"host_os": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ios",
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ios"}, false),
				Description:  "Router host OS. Default value is 'ios'. Only valid value is 'ios'.",
			},
			"ssh_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     22,
				Description: "SSH port to use to connect to the router. Defaults to 22 if not set.",
			},
			"address_1": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Address line 1.",
			},
			"address_2": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Address line 2.",
			},
			"city": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "City",
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "State",
			},
			"country": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ISO two-letter country code.",
			},
			"zip_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Zip code.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description.",
			},
		},
	}
}

// marshalBranchRouterInput marshals the ResourceData into a BranchRouter struct.
func marshalBranchRouterInput(d *schema.ResourceData) *goaviatrix.BranchRouter {
	return &goaviatrix.BranchRouter{
		Name:               d.Get("name").(string),
		PublicIP:           d.Get("public_ip").(string),
		Username:           d.Get("username").(string),
		KeyFile:            d.Get("key_file").(string),
		Password:           d.Get("password").(string),
		PrimaryInterface:   d.Get("wan_primary_interface").(string),
		PrimaryInterfaceIP: d.Get("wan_primary_interface_public_ip").(string),
		BackupInterface:    d.Get("wan_backup_interface").(string),
		BackupInterfaceIP:  d.Get("wan_backup_interface_public_ip").(string),
		HostOS:             d.Get("host_os").(string),
		SshPort:            d.Get("ssh_port").(int),
		SshPortStr:         strconv.Itoa(d.Get("ssh_port").(int)),
		Address1:           d.Get("address_1").(string),
		Address2:           d.Get("address_2").(string),
		City:               d.Get("city").(string),
		State:              d.Get("state").(string),
		Country:            d.Get("country").(string),
		ZipCode:            d.Get("zip_code").(string),
		Description:        d.Get("description").(string),
	}
}

func resourceAviatrixBranchRouterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	br := marshalBranchRouterInput(d)

	if err := client.CreateBranchRouter(br); err != nil {
		return err
	}

	if err := client.ConfigureBranchRouterInterfaces(br); err != nil {
		return err
	}

	d.SetId(br.Name)
	return nil
}

func resourceAviatrixBranchRouterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	name := d.Get("name").(string)
	if name == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no branch router name received. Import Id is %s", id)
		d.SetId(id)
		name = id
	}

	br := &goaviatrix.BranchRouter{
		Name: name,
	}

	br, err := client.GetBranchRouter(br)
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not find branch router %s: %v", name, err)
	}

	d.Set("name", br.Name)
	d.Set("public_ip", br.PublicIP)
	d.Set("username", br.Username)
	d.Set("wan_primary_interface", br.PrimaryInterface)
	d.Set("wan_primary_interface_public_ip", br.PrimaryInterfaceIP)
	d.Set("wan_backup_interface", br.BackupInterface)
	d.Set("wan_backup_interface_public_ip", br.BackupInterfaceIP)
	d.Set("host_os", br.HostOS)
	d.Set("ssh_port", br.SshPort)
	d.Set("address_1", br.Address1)
	d.Set("address_2", br.Address2)
	d.Set("city", br.City)
	d.Set("state", br.State)
	d.Set("country", br.Country)
	d.Set("zip_code", br.ZipCode)
	d.Set("description", br.Description)

	d.SetId(br.Name)
	return nil
}

func resourceAviatrixBranchRouterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	br := marshalBranchRouterInput(d)

	if err := client.UpdateBranchRouter(br); err != nil {
		return err
	}

	if err := client.ConfigureBranchRouterInterfaces(br); err != nil {
		return err
	}

	d.SetId(br.Name)
	return nil
}

func resourceAviatrixBranchRouterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	br := marshalBranchRouterInput(d)

	if err := client.DeleteBranchRouter(br); err != nil {
		return err
	}

	d.SetId(br.Name)
	return nil
}
