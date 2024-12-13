package aviatrix

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var interfaces = []interface{}{
	map[string]interface{}{
		"gateway_ip":                  "192.168.20.1",
		"ip_address":                  "192.168.20.11/24",
		"type":                        "WAN",
		"index":                       0,
		"secondary_private_cidr_list": []interface{}{"192.168.19.16/29"},
	},
	map[string]interface{}{
		"gateway_ip":                  "192.168.21.1",
		"ip_address":                  "192.168.21.11/24",
		"type":                        "WAN",
		"index":                       1,
		"secondary_private_cidr_list": []interface{}{"192.168.21.16/29"},
	},
	map[string]interface{}{
		"dhcp":  true,
		"type":  "MANAGEMENT",
		"index": 0,
	},
	map[string]interface{}{
		"gateway_ip": "192.168.22.1",
		"ip_address": "192.168.22.11/24",
		"type":       "WAN",
		"index":      2,
	},
	map[string]interface{}{
		"gateway_ip": "192.168.23.1",
		"ip_address": "192.168.23.11/24",
		"type":       "WAN",
		"index":      3,
	},
}

var expectedInterfaceDetails = []goaviatrix.EdgeTransitInterface{
	{
		GatewayIp:      "192.168.20.1",
		IpAddress:      "192.168.20.11/24",
		Name:           "eth0",
		Type:           "WAN",
		SecondaryCIDRs: []string{"192.168.19.16/29"},
	},
	{
		GatewayIp:      "192.168.21.1",
		IpAddress:      "192.168.21.11/24",
		Name:           "eth1",
		Type:           "WAN",
		SecondaryCIDRs: []string{"192.168.21.16/29"},
	},
	{
		Dhcp: true,
		Name: "eth2",
		Type: "MANAGEMENT",
	},
	{
		GatewayIp: "192.168.22.1",
		IpAddress: "192.168.22.11/24",
		Name:      "eth3",
		Type:      "WAN",
	},
	{
		GatewayIp: "192.168.23.1",
		IpAddress: "192.168.23.11/24",
		Name:      "eth4",
		Type:      "WAN",
	},
}

func TestAccAviatrixTransitGateway_basic(t *testing.T) {
	var gateway goaviatrix.Gateway

	rName := acctest.RandString(5)

	skipGw := os.Getenv("SKIP_TRANSIT_GATEWAY")
	if skipGw == "yes" {
		t.Skip("Skipping Transit gateway test as SKIP_TRANSIT_GATEWAY is set")
	}

	skipGwAWS := os.Getenv("SKIP_TRANSIT_GATEWAY_AWS")
	skipGwAZURE := os.Getenv("SKIP_TRANSIT_GATEWAY_AZURE")
	skipGwGCP := os.Getenv("SKIP_TRANSIT_GATEWAY_GCP")
	skipGwOCI := os.Getenv("SKIP_TRANSIT_GATEWAY_OCI")
	skipGwAEP := os.Getenv("SKIP_TRANSIT_GATEWAY_AEP")

	if skipGwAWS == "yes" && skipGwAZURE == "yes" && skipGwGCP == "yes" && skipGwOCI == "yes" && skipGwAEP == "yes" {
		t.Skip("Skipping Transit gateway test as SKIP_TRANSIT_GATEWAY_AWS, SKIP_TRANSIT_GATEWAY_AZURE, " +
			"SKIP_TRANSIT_GATEWAY_GCP and SKIP_TRANSIT_GATEWAY_OCI are all set")
	}

	if skipGwAWS != "yes" {
		resourceNameAws := "aviatrix_transit_gateway.test_transit_gateway_aws"
		msgCommonAws := ". Set SKIP_TRANSIT_GATEWAY_AWS to yes to skip Transit Gateway tests in aws"

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheck(t, msgCommonAws)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckTransitGatewayDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccTransitGatewayConfigBasicAWS(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckTransitGatewayExists(resourceNameAws, &gateway),
						resource.TestCheckResourceAttr(resourceNameAws, "gw_name", fmt.Sprintf("tfg-aws-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameAws, "gw_size", "t2.micro"),
						resource.TestCheckResourceAttr(resourceNameAws, "account_name", fmt.Sprintf("tfa-aws-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameAws, "vpc_id", os.Getenv("AWS_VPC_ID")),
						resource.TestCheckResourceAttr(resourceNameAws, "subnet", os.Getenv("AWS_SUBNET")),
						resource.TestCheckResourceAttr(resourceNameAws, "vpc_reg", os.Getenv("AWS_REGION")),
						resource.TestCheckResourceAttr(resourceNameAws, "bgp_polling_time", "50"),
						resource.TestCheckResourceAttr(resourceNameAws, "bgp_neighbor_status_polling_time", "0"),
					),
				},
				{
					ResourceName:      resourceNameAws,
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccTransitGatewayConfigAWSBasicBgpBfd(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckTransitGatewayExists(resourceNameAws, &gateway),
						resource.TestCheckResourceAttr(resourceNameAws, "gw_name", fmt.Sprintf("tfg-aws-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameAws, "gw_size", "t2.micro"),
						resource.TestCheckResourceAttr(resourceNameAws, "account_name", fmt.Sprintf("tfa-aws-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameAws, "vpc_id", os.Getenv("AWS_VPC_ID")),
						resource.TestCheckResourceAttr(resourceNameAws, "subnet", os.Getenv("AWS_SUBNET")),
						resource.TestCheckResourceAttr(resourceNameAws, "vpc_reg", os.Getenv("AWS_REGION")),
						resource.TestCheckResourceAttr(resourceNameAws, "bgp_polling_time", "50"),
						resource.TestCheckResourceAttr(resourceNameAws, "bgp_neighbor_status_polling_time", "7"),
					),
				},
			},
		})
	} else {
		t.Log("Skipping Transit gateway test in aws as SKIP_TRANSIT_GATEWAY_AWS is set")
	}

	if skipGwAZURE != "yes" {
		resourceNameAzure := "aviatrix_transit_gateway.test_transit_gateway_azure"

		msgCommonAzure := ". Set SKIP_TRANSIT_GATEWAY_AZURE to yes to skip Transit Gateway tests in Azure"

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheckAZURE(t, msgCommonAzure)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckTransitGatewayDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccTransitGatewayConfigBasicAZURE(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckTransitGatewayExists(resourceNameAzure, &gateway),
						resource.TestCheckResourceAttr(resourceNameAzure, "gw_name", fmt.Sprintf("tfg-azure-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameAzure, "gw_size", os.Getenv("AZURE_GW_SIZE")),
						resource.TestCheckResourceAttr(resourceNameAzure, "account_name", fmt.Sprintf("tfa-azure-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameAzure, "vpc_id", os.Getenv("AZURE_VNET_ID")),
						resource.TestCheckResourceAttr(resourceNameAzure, "subnet", os.Getenv("AZURE_SUBNET")),
						resource.TestCheckResourceAttr(resourceNameAzure, "vpc_reg", os.Getenv("AZURE_REGION")),
					),
				},
				{
					ResourceName:      resourceNameAzure,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	} else {
		t.Log("Skipping Transit gateway test in aws as SKIP_TRANSIT_GATEWAY_Azure is set")
	}

	if skipGwGCP != "yes" {
		resourceNameGCP := "aviatrix_transit_gateway.test_transit_gateway_gcp"
		gcpGwSize := os.Getenv("GCP_GW_SIZE")
		if gcpGwSize == "" {
			gcpGwSize = "n1-standard-1"
		}

		msgCommonGCP := ". Set SKIP_TRANSIT_GATEWAY_GCP to yes to skip Transit Gateway tests in GCP"

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheckGCP(t, msgCommonGCP)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckTransitGatewayDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccTransitGatewayConfigBasicGCP(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckTransitGatewayExists(resourceNameGCP, &gateway),
						resource.TestCheckResourceAttr(resourceNameGCP, "gw_name", fmt.Sprintf("tfg-gcp-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameGCP, "gw_size", gcpGwSize),
						resource.TestCheckResourceAttr(resourceNameGCP, "account_name", fmt.Sprintf("tfa-gcp-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameGCP, "vpc_id", os.Getenv("GCP_VPC_ID")),
						resource.TestCheckResourceAttr(resourceNameGCP, "subnet", os.Getenv("GCP_SUBNET")),
						resource.TestCheckResourceAttr(resourceNameGCP, "vpc_reg", os.Getenv("GCP_ZONE")),
					),
				},
				{
					ResourceName:      resourceNameGCP,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	} else {
		t.Log("Skipping Transit gateway test in aws as SKIP_TRANSIT_GATEWAY_GCP is set")
	}

	if skipGwOCI != "yes" {
		resourceNameOCI := "aviatrix_transit_gateway.test_transit_gateway_oci"
		ociGwSize := os.Getenv("OCI_GW_SIZE")
		if ociGwSize == "" {
			ociGwSize = "VM.Standard2.2"
		}

		msgCommonOCI := ". Set SKIP_TRANSIT_GATEWAY_OCI to yes to skip Transit Gateway tests in OCI"

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheckGCP(t, msgCommonOCI)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckTransitGatewayDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccTransitGatewayConfigBasicOCI(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckTransitGatewayExists(resourceNameOCI, &gateway),
						resource.TestCheckResourceAttr(resourceNameOCI, "gw_name", fmt.Sprintf("tfg-oci-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameOCI, "gw_size", ociGwSize),
						resource.TestCheckResourceAttr(resourceNameOCI, "account_name", fmt.Sprintf("tfa-oci-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameOCI, "vpc_id", os.Getenv("OCI_VPC_ID")),
						resource.TestCheckResourceAttr(resourceNameOCI, "subnet", os.Getenv("OCI_SUBNET")),
						resource.TestCheckResourceAttr(resourceNameOCI, "vpc_reg", os.Getenv("OCI_REGION")),
					),
				},
				{
					ResourceName:      resourceNameOCI,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	} else {
		t.Log("Skipping Transit gateway test in aws as SKIP_TRANSIT_GATEWAY_OCI is set")
	}

	if skipGwAEP != "yes" {
		resourceNameAEP := "aviatrix_transit_gateway.test_transit_gateway_aep"
		msgCommonAEP := ". Set SKIP_TRANSIT_GATEWAY_AEP to yes to skip Transit Gateway tests in edge AEP"

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheckEdge(t, msgCommonAEP)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckTransitGatewayDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccTransitGatewayConfigBasicAEP(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckTransitGatewayExists(resourceNameAEP, &gateway),
						resource.TestCheckResourceAttr(resourceNameAEP, "gw_name", fmt.Sprintf("tfg-aep-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameAEP, "gw_size", "SMALL"),
						resource.TestCheckResourceAttr(resourceNameAEP, "vpc_id", os.Getenv("AEP_VPC_ID")),
						resource.TestCheckResourceAttr(resourceNameAEP, "device_id", os.Getenv("AEP_DEVICE_ID")),
						resource.TestCheckResourceAttr(resourceNameAEP, "interfaces.0.gateway_ip", "192.168.20.1"),
						resource.TestCheckResourceAttr(resourceNameAEP, "interfaces.0.ip_address", "192.168.20.11/24"),
						resource.TestCheckResourceAttr(resourceNameAEP, "interfaces.0.type", "WAN"),
						resource.TestCheckResourceAttr(resourceNameAEP, "interfaces.0.index", "0"),
					),
				},
				{
					ResourceName:      resourceNameAEP,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	} else {
		t.Log("Skipping Transit gateway test in edge AEP as SKIP_TRANSIT_GATEWAY_AEP is set")
	}
}

func testAccTransitGatewayConfigBasicAWS(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_transit_gateway" "test_transit_gateway_aws" {
	cloud_type                       = 1
	account_name                     = aviatrix_account.test_acc_aws.account_name
	gw_name                          = "tfg-aws-%[1]s"
	vpc_id                           = "%[5]s"
	vpc_reg                          = "%[6]s"
	gw_size                          = "t2.micro"
	subnet                           = "%[7]s"
	bgp_polling_time                 = 50
	bgp_neighbor_status_polling_time = 5
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"))
}

func testAccTransitGatewayConfigAWSBasicBgpBfd(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_transit_gateway" "test_transit_gateway_aws" {
	cloud_type                       = 1
	account_name                     = aviatrix_account.test_acc_aws.account_name
	gw_name                          = "tfg-aws-%s"
	vpc_id                           = "%s"
	vpc_reg                          = "%s"
	gw_size                          = "t2.micro"
	subnet                           = "%s"
	bgp_polling_time                 = 50
	bgp_neighbor_status_polling_time = 7
}
	`, rName, os.Getenv("AWS_VPC_ID"), os.Getenv("AWS_REGION"), os.Getenv("AWS_SUBNET"))
}

func testAccTransitGatewayConfigBasicAZURE(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_azure" {
	account_name        = "tfa-azure-%s"
	cloud_type          = 8
	arm_subscription_id = "%s"
	arm_directory_id    = "%s"
	arm_application_id  = "%s"
	arm_application_key = "%s"
}
resource "aviatrix_transit_gateway" "test_transit_gateway_azure" {
	cloud_type   = 8
	account_name = aviatrix_account.test_acc_azure.account_name
	gw_name      = "tfg-azure-%[1]s"
	vpc_id       = "%[6]s"
	vpc_reg      = "%[7]s"
	gw_size      = "%[8]s"
	subnet       = "%[9]s"
}
	`, rName, os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_DIRECTORY_ID"), os.Getenv("ARM_APPLICATION_ID"),
		os.Getenv("ARM_APPLICATION_KEY"), os.Getenv("AZURE_VNET_ID"), os.Getenv("AZURE_REGION"),
		os.Getenv("AZURE_GW_SIZE"), os.Getenv("AZURE_SUBNET"))
}

func testAccTransitGatewayConfigBasicGCP(rName string) string {
	gcpGwSize := os.Getenv("GCP_GW_SIZE")
	if gcpGwSize == "" {
		gcpGwSize = "n1-standard-1"
	}
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_gcp" {
	account_name                        = "tfa-gcp-%s"
	cloud_type                          = 4
	gcloud_project_id                   = "%s"
	gcloud_project_credentials_filepath = "%s"
}
resource "aviatrix_transit_gateway" "test_transit_gateway_gcp" {
	cloud_type   = 4
	account_name = aviatrix_account.test_acc_gcp.account_name
	gw_name      = "tfg-gcp-%[1]s"
	vpc_id       = "%[4]s"
	vpc_reg      = "%[5]s"
	gw_size      = "%[6]s"
	subnet       = "%[7]s"
}
	`, rName, os.Getenv("GCP_ID"), os.Getenv("GCP_CREDENTIALS_FILEPATH"),
		os.Getenv("GCP_VPC_ID"), os.Getenv("GCP_ZONE"), gcpGwSize, os.Getenv("GCP_SUBNET"))
}

func testAccTransitGatewayConfigBasicOCI(rName string) string {
	ociGwSize := os.Getenv("OCI_GW_SIZE")
	if ociGwSize == "" {
		ociGwSize = "VM.Standard2.2"
	}
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_oci" {
	account_name                 = "tfa-oci-%s"
	cloud_type                   = 16
	oci_tenancy_id               = "%s"
	oci_user_id                  = "%s"
	oci_compartment_id           = "%s"
	oci_api_private_key_filepath = "%s"
}
resource "aviatrix_transit_gateway" "test_transit_gateway_oci" {
	cloud_type   = 16
	account_name = aviatrix_account.test_acc_oci.account_name
	gw_name      = "tfg-oci-%[1]s"
	vpc_id       = "%[6]s"
	vpc_reg      = "%[7]s"
	gw_size      = "%[8]s"
	subnet       = "%[9]s"
}
	`, rName, os.Getenv("OCI_TENANCY_ID"), os.Getenv("OCI_USER_ID"), os.Getenv("OCI_COMPARTMENT_ID"),
		os.Getenv("OCI_API_KEY_FILEPATH"), os.Getenv("OCI_VPC_ID"), os.Getenv("OCI_REGION"),
		ociGwSize, os.Getenv("OCI_SUBNET"))
}

func testAccTransitGatewayConfigBasicAEP(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_edge_aep" {
	account_name       = "edge-%s"
	cloud_type         = 262144
}
resource "aviatrix_transit_gateway" "test_transit_gateway_aep" {
	cloud_type   = 262144
	account_name = aviatrix_account.test_acc_edge_aep.account_name
	gw_name      = "tfg-edge-aep-%[1]s"
	vpc_id       = "%[2]s"
	device_id = "%[3]s"
	gw_size      = "SMALL"
	interfaces {
        gateway_ip    = "192.168.20.1"
        ip_address    = "192.168.20.11/24"
        type          = "WAN"
        index         = 0
        secondary_private_cidr_list = ["192.168.19.16/29"]
    }

    interfaces {
        gateway_ip    = "192.168.21.1"
        ip_address    = "192.168.21.11/24"
        type          = "WAN"
        index         = 1
        secondary_private_cidr_list = ["192.168.21.16/29"]
    }

    interfaces {
        dhcp   = true
        type   = "MANAGEMENT"
        index  = 0
    }

    interfaces {
        gateway_ip  = "192.168.22.1"
        ip_address  = "192.168.22.11/24"
        type        = "WAN"
        index       = 2
    }

    interfaces {
        gateway_ip = "192.168.23.1"
        ip_address = "192.168.23.11/24"
        type       = "WAN"
        index      = 3
    }
}
	`, rName, os.Getenv("AEP_VPC_ID"), os.Getenv("AEP_DEVICE_ID"))
}

func testAccCheckTransitGatewayExists(n string, gateway *goaviatrix.Gateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("transit gateway Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no transit gateway ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		foundGateway := &goaviatrix.Gateway{
			GwName:      rs.Primary.Attributes["gw_name"],
			AccountName: rs.Primary.Attributes["account_name"],
		}
		_, err := client.GetGateway(foundGateway)
		if err != nil {
			return err
		}
		if foundGateway.GwName != rs.Primary.ID {
			return fmt.Errorf("transit gateway not found")
		}

		*gateway = *foundGateway
		return nil
	}
}

func testAccCheckTransitGatewayDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_transit_vpc" {
			continue
		}

		foundGateway := &goaviatrix.Gateway{
			GwName:      rs.Primary.Attributes["gw_name"],
			AccountName: rs.Primary.Attributes["account_name"],
		}

		_, err := client.GetGateway(foundGateway)
		if err == nil {
			return fmt.Errorf("transit gateway still exists")
		}
	}

	return nil
}

func TestGetInterfaceMappingDetails(t *testing.T) {
	tests := []struct {
		name                  string
		interfaceMappingInput []interface{}
		expectedOutput        string
		expectedError         error
	}{
		{
			name: "Valid input for ESXI devices",
			interfaceMappingInput: []interface{}{
				map[string]interface{}{
					"name":  "eth0",
					"type":  "mgmt",
					"index": 0,
				},
				map[string]interface{}{
					"name":  "eth1",
					"type":  "wan",
					"index": 1,
				},
			},
			expectedOutput: `{"eth0":["mgmt","0"],"eth1":["wan","1"]}`,
			expectedError:  nil,
		},
		{
			name:                  "Empty input (default Dell device mapping)",
			interfaceMappingInput: []interface{}{},
			expectedOutput:        `{"eth0":["mgmt","0"],"eth2":["wan","1"],"eth3":["wan","2"],"eth4":["wan","3"],"eth5":["wan","0"]}`,
			expectedError:         nil,
		},
		{
			name: "Invalid input type (non-map element)",
			interfaceMappingInput: []interface{}{
				"invalid_type", // This is not a map
			},
			expectedOutput: "",
			expectedError:  fmt.Errorf("invalid type string for interface mapping, expected a map"),
		},
		{
			name: "Invalid map fields (missing required keys)",
			interfaceMappingInput: []interface{}{
				map[string]interface{}{
					"name": "eth0", // Missing 'type' and 'index'
				},
			},
			expectedOutput: "",
			expectedError:  fmt.Errorf("invalid interface mapping, 'name', 'type', and 'index' must be strings"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := getInterfaceMappingDetails(tt.interfaceMappingInput)
			if err != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			}
			if result != tt.expectedOutput {
				t.Errorf("expected output: %s, got: %s", tt.expectedOutput, result)
			}
		})
	}
}

func TestGetInterfaceName(t *testing.T) {
	tests := []struct {
		name        string
		intfType    string
		intfIndex   int
		wanCount    int
		expected    string
		expectedErr error
	}{
		{
			name:        "Valid WAN interface with index 0",
			intfType:    "WAN",
			intfIndex:   0,
			wanCount:    3,
			expected:    "eth0",
			expectedErr: nil,
		},
		{
			name:        "Valid WAN interface with index 1",
			intfType:    "WAN",
			intfIndex:   1,
			wanCount:    3,
			expected:    "eth1",
			expectedErr: nil,
		},
		{
			name:        "Valid WAN interface with index 2",
			intfType:    "WAN",
			intfIndex:   2,
			wanCount:    3,
			expected:    "eth3",
			expectedErr: nil,
		},
		{
			name:        "Valid MANAGEMENT interface with index 0",
			intfType:    "MANAGEMENT",
			intfIndex:   0,
			wanCount:    3,
			expected:    "eth2",
			expectedErr: nil,
		},
		{
			name:        "Valid MANAGEMENT interface with index 1",
			intfType:    "MANAGEMENT",
			intfIndex:   1,
			wanCount:    3,
			expected:    "eth4",
			expectedErr: nil,
		},
		{
			name:        "Invalid interface type",
			intfType:    "INVALID",
			intfIndex:   0,
			wanCount:    3,
			expected:    "",
			expectedErr: errors.New("invalid interface type INVALID"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := getInterfaceName(tt.intfType, tt.intfIndex, tt.wanCount)

			// Check error
			if err != nil && tt.expectedErr != nil {
				if err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
				}
			} else if (err != nil && tt.expectedErr == nil) || (err == nil && tt.expectedErr != nil) {
				t.Errorf("unexpected error: %v", err)
			}

			// Check result
			if result != tt.expected {
				t.Errorf("expected result: %s, got: %s", tt.expected, result)
			}
		})
	}
}

func TestGetEipMapDetails(t *testing.T) {
	tests := []struct {
		name        string
		eipMap      []interface{}
		wanCount    int
		expected    string
		expectedErr error
	}{
		{
			name: "Valid EIP map with WAN and MANAGEMENT interfaces",
			eipMap: []interface{}{
				map[string]interface{}{
					"interface_type":  "WAN",
					"interface_index": 0,
					"private_ip":      "192.168.0.10",
					"public_ip":       "203.0.113.10",
				},
				map[string]interface{}{
					"interface_type":  "MANAGEMENT",
					"interface_index": 0,
					"private_ip":      "192.168.1.10",
					"public_ip":       "203.0.113.11",
				},
			},
			wanCount:    3,
			expected:    `{"eth0":[{"private_ip":"192.168.0.10","public_ip":"203.0.113.10"}],"eth2":[{"private_ip":"192.168.1.10","public_ip":"203.0.113.11"}]}`,
			expectedErr: nil,
		},
		{
			name: "Invalid EIP map: missing interface type",
			eipMap: []interface{}{
				map[string]interface{}{
					"interface_index": 0,
					"private_ip":      "192.168.0.10",
					"public_ip":       "203.0.113.10",
				},
			},
			wanCount:    3,
			expected:    "",
			expectedErr: errors.New("interface_type must be a string"),
		},
		{
			name: "Invalid EIP map: invalid interface type",
			eipMap: []interface{}{
				map[string]interface{}{
					"interface_type":  "INVALID",
					"interface_index": 0,
					"private_ip":      "192.168.0.10",
					"public_ip":       "203.0.113.10",
				},
			},
			wanCount:    3,
			expected:    "",
			expectedErr: errors.New("failed to get the interface name using type and index for eip_map: invalid interface type INVALID"),
		},
		{
			name:        "Empty EIP map",
			eipMap:      []interface{}{},
			wanCount:    3,
			expected:    `{}`,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := getEipMapDetails(tt.eipMap, tt.wanCount)

			// Check for errors
			if err != nil && tt.expectedErr != nil {
				if err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
				}
			} else if (err != nil && tt.expectedErr == nil) || (err == nil && tt.expectedErr != nil) {
				t.Errorf("unexpected error: %v", err)
			}

			// Check result
			if result != tt.expected {
				t.Errorf("expected result: %s, got: %s", tt.expected, result)
			}
		})
	}
}

// test to count the interface types in the gateway
func TestCountInterfaceTypes(t *testing.T) {
	// count the WAN interfaces
	wanCount, err := countInterfaceTypes(interfaces)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	// Check that the WAN count matches the expected value
	expectedWANCount := 4
	if wanCount != expectedWANCount {
		t.Errorf("Expected %d WAN interfaces, got %d", expectedWANCount, wanCount)
	}
}

// test to get the interface details from the resource
func TestGetInterfaceDetails(t *testing.T) {
	// get the interface details
	interfaceDetails, err := getInterfaceDetails(interfaces)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	// base64 encode the expected string
	expectedInterfaceDetailsJson, err := json.Marshal(expectedInterfaceDetails)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	// base64 encode the expected string
	expectedInterfaceDetailsEncoded := base64.StdEncoding.EncodeToString(expectedInterfaceDetailsJson)
	// Check that the interface details are as expected
	if interfaceDetails != expectedInterfaceDetailsEncoded {
		t.Errorf("Expected %s, got %s", expectedInterfaceDetailsEncoded, interfaceDetails)
	}
}
