package aviatrix

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"aviatrix.com/terraform-provider-aviatrix/goaviatrix"
)

func TestAccDataSourceAviatrixFireNetFirewallManager_basic(t *testing.T) {
	rName := acctest.RandString(5)
	resourceName := "data.aviatrix_firenet_firewall_manager.test"

	skipAcc := os.Getenv("SKIP_DATA_FIRENET_FIREWALL_MANAGER")
	if skipAcc == "yes" {
		t.Skip("Skipping Data Source FireNet Firewall Manager test as SKIP_DATA_FIRENET_FIREWALL_MANAGER is set")
	}
	msg := ". Set SKIP_DATA_FIRENET_FIREWALL_MANAGER to yes to skip Data Source FireNet FIREWALL MANAGER tests"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, msg)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAviatrixFireNetFirewallManagerConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAviatrixFireNetFirewallManager(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gateway_name", fmt.Sprintf("tftg-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "vendor_type", "Generic"),
				),
			},
		},
	})
}

func TestAccDataSourceAviatrixFireNetFirewallManager_panorama(t *testing.T) {
	if os.Getenv("SKIP_DATA_FIRENET_FIREWALL_MANAGER") == "yes" {
		t.Skip("Skipping Data Source FireNet Firewall Manager test as SKIP_DATA_FIRENET_FIREWALL_MANAGER is set")
	}

	rName := acctest.RandString(5)
	resourceName := "data.aviatrix_firenet_firewall_manager.test_panorama"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, ". Set SKIP_DATA_FIRENET_FIREWALL_MANAGER to yes to skip Data Source FireNet FIREWALL MANAGER tests")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAviatrixFireNetFirewallManagerConfigPanorama(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAviatrixFireNetFirewallManager(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gateway_name", fmt.Sprintf("tftg-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "vendor_type", "Palo Alto Networks Panorama"),
					resource.TestCheckResourceAttr(resourceName, "public_ip", "192.168.1.10"),
					resource.TestCheckResourceAttr(resourceName, "username", "admin"),
					resource.TestCheckResourceAttr(resourceName, "template", "test-template"),
					resource.TestCheckResourceAttr(resourceName, "template_stack", "test-stack"),
				),
			},
		},
	})
}

func TestAccDataSourceAviatrixFireNetFirewallManager_advancedConfig(t *testing.T) {
	if os.Getenv("SKIP_DATA_FIRENET_FIREWALL_MANAGER") == "yes" {
		t.Skip("Skipping Data Source FireNet Firewall Manager test as SKIP_DATA_FIRENET_FIREWALL_MANAGER is set")
	}

	rName := acctest.RandString(5)
	resourceName := "data.aviatrix_firenet_firewall_manager.test_advanced"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, ". Set SKIP_DATA_FIRENET_FIREWALL_MANAGER to yes to skip Data Source FireNet FIREWALL MANAGER tests")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAviatrixFireNetFirewallManagerConfigAdvanced(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAviatrixFireNetFirewallManager(resourceName),
					resource.TestCheckResourceAttr(resourceName, "gateway_name", fmt.Sprintf("tftg-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "config_mode", "ADVANCE"),
					resource.TestCheckResourceAttr(resourceName, "vendor_type", "Palo Alto Networks Panorama"),
				),
			},
		},
	})
}

// Unit test for parseFirewallTemplateConfig function
func TestParseFirewallTemplateConfig(t *testing.T) {
	tests := []struct {
		name     string
		input    *schema.Set
		expected map[string]goaviatrix.FirewallTemplateConfig
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty set",
			input:    schema.NewSet(schema.HashResource(&schema.Resource{}), []interface{}{}),
			expected: nil,
		},
		{
			name: "valid configuration",
			input: schema.NewSet(schema.HashResource(&schema.Resource{
				Schema: map[string]*schema.Schema{
					"firewall_id":    {Type: schema.TypeString},
					"template":       {Type: schema.TypeString},
					"template_stack": {Type: schema.TypeString},
					"route_table":    {Type: schema.TypeString},
				},
			}), []interface{}{
				map[string]interface{}{
					"firewall_id":    "firewall1",
					"template":       "template1",
					"template_stack": "stack1",
					"route_table":    "vrouter1",
				},
				map[string]interface{}{
					"firewall_id":    "firewall2",
					"template":       "template2",
					"template_stack": "stack2",
					"route_table":    "vrouter2",
				},
			}),
			expected: map[string]goaviatrix.FirewallTemplateConfig{
				"firewall1": {
					Template:      "template1",
					TemplateStack: "stack1",
					RouteTable:    "vrouter1",
				},
				"firewall2": {
					Template:      "template2",
					TemplateStack: "stack2",
					RouteTable:    "vrouter2",
				},
			},
		},
		{
			name: "partial configuration",
			input: schema.NewSet(schema.HashResource(&schema.Resource{
				Schema: map[string]*schema.Schema{
					"firewall_id":    {Type: schema.TypeString},
					"template":       {Type: schema.TypeString},
					"template_stack": {Type: schema.TypeString},
					"route_table":    {Type: schema.TypeString},
				},
			}), []interface{}{
				map[string]interface{}{
					"firewall_id": "firewall1",
					"template":    "template1",
					// missing template_stack and route_table
				},
			}),
			expected: map[string]goaviatrix.FirewallTemplateConfig{
				"firewall1": {
					Template:      "template1",
					TemplateStack: "",
					RouteTable:    "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseFirewallTemplateConfig(tt.input)

			if tt.expected == nil && result != nil {
				t.Errorf("Expected nil, got %v", result)
				return
			}

			if tt.expected != nil && result == nil {
				t.Errorf("Expected %v, got nil", tt.expected)
				return
			}

			if tt.expected == nil && result == nil {
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d items, got %d", len(tt.expected), len(result))
				return
			}

			for key, expectedConfig := range tt.expected {
				actualConfig, exists := result[key]
				if !exists {
					t.Errorf("Expected key %s not found in result", key)
					continue
				}

				if actualConfig.Template != expectedConfig.Template {
					t.Errorf("For key %s, expected template %s, got %s", key, expectedConfig.Template, actualConfig.Template)
				}
				if actualConfig.TemplateStack != expectedConfig.TemplateStack {
					t.Errorf("For key %s, expected template_stack %s, got %s", key, expectedConfig.TemplateStack, actualConfig.TemplateStack)
				}
				if actualConfig.RouteTable != expectedConfig.RouteTable {
					t.Errorf("For key %s, expected route_table %s, got %s", key, expectedConfig.RouteTable, actualConfig.RouteTable)
				}
			}
		})
	}
}

func testAccDataSourceAviatrixFireNetFirewallManagerConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_vpc" "test_vpc" {
	cloud_type           = 1
	account_name         = aviatrix_account.test_account.account_name
	region               = "%s"
	name                 = "vpc-for-firenet"
	cidr                 = "10.10.0.0/24"
	aviatrix_firenet_vpc = true
}
resource "aviatrix_transit_gateway" "test_transit_gateway" {
	cloud_type               = aviatrix_vpc.test_vpc.cloud_type
	account_name             = aviatrix_account.test_account.account_name
	gw_name                  = "tftg-%s"
	vpc_id                   = aviatrix_vpc.test_vpc.vpc_id
	vpc_reg                  = aviatrix_vpc.test_vpc.region
	gw_size                  = "c5.xlarge"
	subnet                   = aviatrix_vpc.test_vpc.subnets[0].cidr
	enable_hybrid_connection = true
	enable_firenet           = true
}
data "aviatrix_firenet_firewall_manager" "test" {
	vpc_id       = aviatrix_vpc.test_vpc.vpc_id
	gateway_name = aviatrix_transit_gateway.test_transit_gateway.gw_name
	vendor_type  = "Generic"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_REGION"), rName)
}

func testAccDataSourceAviatrixFireNetFirewallManagerConfigPanorama(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_vpc" "test_vpc" {
	cloud_type           = 1
	account_name         = aviatrix_account.test_account.account_name
	region               = "%s"
	name                 = "vpc-for-firenet-panorama"
	cidr                 = "10.11.0.0/24"
	aviatrix_firenet_vpc = true
}
resource "aviatrix_transit_gateway" "test_transit_gateway" {
	cloud_type               = aviatrix_vpc.test_vpc.cloud_type
	account_name             = aviatrix_account.test_account.account_name
	gw_name                  = "tftg-%s"
	vpc_id                   = aviatrix_vpc.test_vpc.vpc_id
	vpc_reg                  = aviatrix_vpc.test_vpc.region
	gw_size                  = "c5.xlarge"
	subnet                   = aviatrix_vpc.test_vpc.subnets[0].cidr
	enable_hybrid_connection = true
	enable_firenet           = true
}
data "aviatrix_firenet_firewall_manager" "test_panorama" {
	vpc_id         = aviatrix_vpc.test_vpc.vpc_id
	gateway_name   = aviatrix_transit_gateway.test_transit_gateway.gw_name
	vendor_type    = "Palo Alto Networks Panorama"
	public_ip      = "192.168.1.10"
	username       = "admin"
	password       = "password123"
	template       = "test-template"
	template_stack = "test-stack"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_REGION"), rName)
}

func testAccDataSourceAviatrixFireNetFirewallManagerConfigAdvanced(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_account" {
	account_name       = "tfa-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_vpc" "test_vpc" {
	cloud_type           = 1
	account_name         = aviatrix_account.test_account.account_name
	region               = "%s"
	name                 = "vpc-for-firenet-advanced"
	cidr                 = "10.12.0.0/24"
	aviatrix_firenet_vpc = true
}
resource "aviatrix_transit_gateway" "test_transit_gateway" {
	cloud_type               = aviatrix_vpc.test_vpc.cloud_type
	account_name             = aviatrix_account.test_account.account_name
	gw_name                  = "tftg-%s"
	vpc_id                   = aviatrix_vpc.test_vpc.vpc_id
	vpc_reg                  = aviatrix_vpc.test_vpc.region
	gw_size                  = "c5.xlarge"
	subnet                   = aviatrix_vpc.test_vpc.subnets[0].cidr
	enable_hybrid_connection = true
	enable_firenet           = true
}
data "aviatrix_firenet_firewall_manager" "test_advanced" {
	vpc_id         = aviatrix_vpc.test_vpc.vpc_id
	gateway_name   = aviatrix_transit_gateway.test_transit_gateway.gw_name
	vendor_type    = "Palo Alto Networks Panorama"
	public_ip      = "192.168.1.20"
	username       = "admin"
	password       = "password456"
	config_mode    = "ADVANCE"
	firewall_template_config {
		firewall_id    = "firewall-1"
		template       = "advanced-template-1"
		template_stack = "advanced-stack-1"
		route_table    = "vrouter-1"
	}
	firewall_template_config {
		firewall_id    = "firewall-2"
		template       = "advanced-template-2"  
		template_stack = "advanced-stack-2"
		route_table    = "vrouter-2"
	}
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_REGION"), rName)
}

func testAccDataSourceAviatrixFireNetFirewallManager(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
