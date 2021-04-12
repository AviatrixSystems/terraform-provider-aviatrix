package aviatrix

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAviatrixAwsTgwSecurityDomainConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAviatrixAwsTgwSecurityDomainConnectionCreate,
		ReadContext:   resourceAviatrixAwsTgwSecurityDomainConnectionRead,
		DeleteContext: resourceAviatrixAwsTgwSecurityDomainConnectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"tgw_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "AWS TGW name.",
			},
			"domain_name1": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					domainName2Old, _ := d.GetChange("domain_name2")
					return old == d.Get("domain_name2").(string) && new == domainName2Old.(string)
				},
				Description: "Security domain name 1.",
			},
			"domain_name2": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					domainName1Old, _ := d.GetChange("domain_name1")
					return old == d.Get("domain_name1").(string) && new == domainName1Old.(string)
				},
				Description: "Security domain name 2.",
			},
		},
	}
}

func resourceAviatrixAwsTgwSecurityDomainConnectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	sourceDomainName := d.Get("domain_name1").(string)
	destinationDomainName := d.Get("domain_name2").(string)

	if sourceDomainName == destinationDomainName {
		return diag.Errorf("two domains cannot be the same")
	}

	awsTgw := &goaviatrix.AWSTgw{
		Name: d.Get("tgw_name").(string),
	}

	securityDomain := &goaviatrix.SecurityDomain{
		Name:       sourceDomainName,
		AwsTgwName: awsTgw.Name,
	}

	securityDomainDetails, err := client.GetSecurityDomainDetails(ctx, securityDomain)
	if err != nil && err != goaviatrix.ErrNotFound {
		return diag.Errorf("couldn't get the details of the security domain %s due to %v", securityDomain.Name, err)
	}

	for _, sd := range securityDomainDetails.ConnectedDomain {
		if sd == destinationDomainName {
			return diag.Errorf("the connection between %s and %s already exits, please import to manage with Terraform", sourceDomainName, destinationDomainName)
		}
	}

	if err := client.CreateDomainConnection(awsTgw, sourceDomainName, destinationDomainName); err != nil {
		return diag.Errorf("could not create the security domain connection: %v", err)
	}

	if sourceDomainName < destinationDomainName {
		d.SetId(awsTgw.Name + "~" + sourceDomainName + "~" + destinationDomainName)
	} else {
		d.SetId(awsTgw.Name + "~" + destinationDomainName + "~" + sourceDomainName)
	}

	return resourceAviatrixAwsTgwSecurityDomainConnectionRead(ctx, d, meta)
}

func resourceAviatrixAwsTgwSecurityDomainConnectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	if d.Get("tgw_name").(string) == "" {
		id := d.Id()
		log.Printf("[DEBUG] Looks like an import. Import Id is %s", id)

		parts := strings.Split(id, "~")
		if len(parts) != 3 {
			return diag.Errorf("invalid ID format, expected ID in format tgw_name~domain_name1~domain_name2, instead got %s", d.Id())
		}

		tgwName := parts[0]
		sourceDomainName := parts[1]
		destinationDomainName := parts[2]

		if tgwName == "" || sourceDomainName == "" || destinationDomainName == "" {
			return diag.Errorf("tgw_name, domain_name1, or domain_name2 cannot be empty")
		}

		if sourceDomainName == destinationDomainName {
			return diag.Errorf("two domains cannot be the same")
		}

		d.Set("tgw_name", tgwName)
		d.Set("domain_name1", sourceDomainName)
		d.Set("domain_name2", destinationDomainName)

		if sourceDomainName < destinationDomainName {
			d.SetId(tgwName + "~" + sourceDomainName + "~" + destinationDomainName)
		} else {
			d.SetId(tgwName + "~" + destinationDomainName + "~" + sourceDomainName)
		}
	}

	tgwName := d.Get("tgw_name").(string)
	sourceDomainName := d.Get("domain_name1").(string)
	destinationDomainName := d.Get("domain_name2").(string)

	securityDomain := &goaviatrix.SecurityDomain{
		Name:       sourceDomainName,
		AwsTgwName: tgwName,
	}

	securityDomainDetails, err := client.GetSecurityDomainDetails(ctx, securityDomain)
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("couldn't get the details of the security domain %s: %v", sourceDomainName, err)
	}

	for _, sd := range securityDomainDetails.ConnectedDomain {
		if sd == destinationDomainName {
			if sourceDomainName < destinationDomainName {
				d.SetId(tgwName + "~" + sourceDomainName + "~" + destinationDomainName)
			} else {
				d.SetId(tgwName + "~" + destinationDomainName + "~" + sourceDomainName)
			}
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceAviatrixAwsTgwSecurityDomainConnectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*goaviatrix.Client)

	awsTgw := &goaviatrix.AWSTgw{
		Name: d.Get("tgw_name").(string),
	}

	if err := client.DeleteDomainConnection(awsTgw, d.Get("domain_name1").(string), d.Get("domain_name2").(string)); err != nil {
		return diag.Errorf("could not delete security domain connection: %v", err)
	}

	return nil
}
