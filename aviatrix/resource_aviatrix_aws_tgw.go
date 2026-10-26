package aviatrix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func resourceAviatrixAWSTgw() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviatrixAWSTgwCreate,
		Read:   resourceAviatrixAWSTgwRead,
		Update: resourceAviatrixAWSTgwUpdate,
		Delete: resourceAviatrixAWSTgwDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 3,
		MigrateState:  resourceAviatrixAWSTgwMigrateState,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceAviatrixAWSTgwResourceV1().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceAviatrixAWSTgwStateUpgradeV1,
				Version: 1,
			},
			{
				Type:    resourceAviatrixAWSTgwResourceV2().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceAviatrixAWSTgwStateUpgradeV2,
				Version: 2,
			},
		},

		Schema: map[string]*schema.Schema{
			"tgw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the AWS TGW which is going to be created.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This parameter represents the name of a Cloud-Account in Aviatrix controller.",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Region of cloud provider.",
			},
			"aws_side_as_number": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "BGP Local ASN (Autonomous System Number), Integer between 1-4294967294.",
				ValidateFunc: goaviatrix.ValidateASN,
			},
			"cloud_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				Description:  "Type of cloud service provider, requires an integer value. Supported for AWS (1) and AWS GOV (256). Default value: 1.",
				ValidateFunc: validation.IntInSlice([]int{goaviatrix.AWS, goaviatrix.AWSGov}),
			},
			"enable_multicast": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable Multicast.",
			},
			"cidrs": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
				Description: "TGW CIDRs.",
			},
			"inspection_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Domain-based",
				Description:  "Inspection mode. Valid values: 'Domain-based' and 'Connection-based'.",
				ValidateFunc: validation.StringInSlice([]string{"Domain-based", "Connection-based"}, false),
			},
			"tgw_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "TGW ID.",
			},
		},
	}
}

func resourceAviatrixAWSTgwCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	awsTgw := &goaviatrix.AWSTgw{
		Name:                    d.Get("tgw_name").(string),
		AccountName:             d.Get("account_name").(string),
		Region:                  d.Get("region").(string),
		AwsSideAsNumber:         d.Get("aws_side_as_number").(string),
		CloudType:               d.Get("cloud_type").(int),
		EnableMulticast:         d.Get("enable_multicast").(bool),
		InspectionMode:          d.Get("inspection_mode").(string),
		NotCreateDefaultDomains: true,
	}

	if awsTgw.Name == "" {
		return fmt.Errorf("tgw name can't be empty string")
	}
	if awsTgw.AccountName == "" {
		return fmt.Errorf("account name can't be empty string")
	}
	if awsTgw.Region == "" {
		return fmt.Errorf("tgw region can't be empty string")
	}
	if awsTgw.AwsSideAsNumber == "" {
		return fmt.Errorf("aws side number can't be empty string")
	}

	log.Printf("[INFO] Creating AWS TGW")

	d.SetId(awsTgw.Name)
	flag := false
	defer func() { _ = resourceAviatrixAWSTgwReadIfRequired(d, meta, &flag) }()

	err1 := client.CreateAWSTgw(awsTgw)
	if err1 != nil {
		return fmt.Errorf("failed to create AWS TGW: %s", err1)
	}

	if cidrs := getStringSet(d, "cidrs"); len(cidrs) != 0 {
		err := client.UpdateTGWCidrs(awsTgw.Name, cidrs)
		if err != nil {
			return fmt.Errorf("could not update TGW CIDRs after creation: %v", err)
		}
	}

	if awsTgw.InspectionMode == "Connection-based" {
		err := client.UpdateTGWInspectionMode(awsTgw.Name, awsTgw.InspectionMode)
		if err != nil {
			return fmt.Errorf("could not update TGW inspection mode after creation: %v", err)
		}
	}

	return resourceAviatrixAWSTgwReadIfRequired(d, meta, &flag)
}

func resourceAviatrixAWSTgwReadIfRequired(d *schema.ResourceData, meta interface{}, flag *bool) error {
	if !(*flag) {
		*flag = true
		return resourceAviatrixAWSTgwRead(d, meta)
	}
	return nil
}

func resourceAviatrixAWSTgwRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)

	tgwName := d.Get("tgw_name").(string)
	if tgwName == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import, no aws tgw name received. Import Id is %s", id)
		d.Set("tgw_name", id)
		d.SetId(id)
	}

	awsTgw := &goaviatrix.AWSTgw{
		Name: d.Get("tgw_name").(string),
	}
	awsTgw, err := client.ListTgwDetails(awsTgw)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't find AWS TGW %s: %v", awsTgw.Name, err)
	}
	d.Set("account_name", awsTgw.AccountName)
	d.Set("tgw_name", awsTgw.Name)
	d.Set("region", awsTgw.Region)
	d.Set("cloud_type", awsTgw.CloudType)
	d.Set("aws_side_as_number", awsTgw.AwsSideAsNumber)
	d.Set("enable_multicast", awsTgw.EnableMulticast)
	d.Set("tgw_id", awsTgw.TgwId)
	d.Set("inspection_mode", awsTgw.InspectionMode)
	if err := d.Set("cidrs", awsTgw.CidrList); err != nil {
		return fmt.Errorf("could not set aws_tgw.cidrs into state: %v", err)
	}

	return nil
}

func resourceAviatrixAWSTgwUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating AWS TGW")

	client := meta.(*goaviatrix.Client)
	awsTgw := &goaviatrix.AWSTgw{
		Name:        d.Get("tgw_name").(string),
		AccountName: d.Get("account_name").(string),
		Region:      d.Get("region").(string),
	}

	d.Partial(true)

	if d.HasChange("account_name") {
		return fmt.Errorf("updating account_name is not allowed")
	}
	if d.HasChange("region") {
		return fmt.Errorf("updating region is not allowed")
	}
	if d.HasChange("cloud_type") {
		return fmt.Errorf("updating cloud_type is not allowed")
	}
	if d.HasChange("enable_multicast") {
		return fmt.Errorf("updating enable_multicast is not allowed")
	}

	if d.HasChange("cidrs") {
		cidrs := getStringSet(d, "cidrs")
		err := client.UpdateTGWCidrs(awsTgw.Name, cidrs)
		if err != nil {
			return fmt.Errorf("could not update TGW CIDRs during update: %v", err)
		}
	}

	if d.HasChange("inspection_mode") {
		err := client.UpdateTGWInspectionMode(awsTgw.Name, d.Get("inspection_mode").(string))
		if err != nil {
			return fmt.Errorf("could not update TGW inspection mode during update: %v", err)
		}
	}

	d.Partial(false)
	d.SetId(awsTgw.Name)
	return resourceAviatrixAWSTgwRead(d, meta)
}

func resourceAviatrixAWSTgwDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goaviatrix.Client)
	awsTgw := &goaviatrix.AWSTgw{
		Name:                      d.Get("tgw_name").(string),
		AccountName:               d.Get("account_name").(string),
		Region:                    d.Get("region").(string),
		AwsSideAsNumber:           d.Get("aws_side_as_number").(string),
		AttachedAviatrixTransitGW: make([]string, 0),
		SecurityDomains:           make([]goaviatrix.SecurityDomainRule, 0),
	}

	log.Printf("[INFO] Deleting AWS TGW")

	err := client.DeleteAWSTgw(awsTgw)
	if err != nil {
		if err == goaviatrix.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("couldn't destroy AWS TGW %s: %v", awsTgw.Name, err)
	}

	return nil
}
