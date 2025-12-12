package aviatrix

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

var interfaces = []interface{}{
	map[string]interface{}{
		"gateway_ip":                  "192.168.20.1",
		"ip_address":                  "192.168.20.11/24",
		"logical_ifname":              "wan0",
		"secondary_private_cidr_list": []interface{}{"192.168.19.16/29"},
	},
	map[string]interface{}{
		"gateway_ip":                  "192.168.21.1",
		"ip_address":                  "192.168.21.11/24",
		"logical_ifname":              "wan1",
		"secondary_private_cidr_list": []interface{}{"192.168.21.16/29"},
	},
	map[string]interface{}{
		"dhcp":           true,
		"logical_ifname": "mgmt0",
	},
	map[string]interface{}{
		"gateway_ip":     "192.168.22.1",
		"ip_address":     "192.168.22.11/24",
		"logical_ifname": "wan2",
	},
	map[string]interface{}{
		"gateway_ip":     "192.168.23.1",
		"ip_address":     "192.168.23.11/24",
		"logical_ifname": "wan3",
	},
	map[string]interface{}{
		"gateway_ip":     "169.254.100.1",
		"ip_address":     "10.0.1.10/24",
		"logical_ifname": "wan4",
		"underlay_cidr":  "169.254.100.2/30",
	},
}

var expectedInterfaceDetails = []goaviatrix.EdgeTransitInterface{
	{
		GatewayIp:      "192.168.20.1",
		IpAddress:      "192.168.20.11/24",
		LogicalIfName:  "wan0",
		SecondaryCIDRs: []string{"192.168.19.16/29"},
	},
	{
		GatewayIp:      "192.168.21.1",
		IpAddress:      "192.168.21.11/24",
		LogicalIfName:  "wan1",
		SecondaryCIDRs: []string{"192.168.21.16/29"},
	},
	{
		Dhcp:          true,
		LogicalIfName: "mgmt0",
	},
	{
		GatewayIp:     "192.168.22.1",
		IpAddress:     "192.168.22.11/24",
		LogicalIfName: "wan2",
	},
	{
		GatewayIp:     "192.168.23.1",
		IpAddress:     "192.168.23.11/24",
		LogicalIfName: "wan3",
	},
	{
		GatewayIp:     "169.254.100.1",
		IpAddress:     "10.0.1.10/24",
		LogicalIfName: "wan4",
		UnderlayCidr:  "169.254.100.2/30",
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
	skipGwEQUINIX := os.Getenv("SKIP_TRANSIT_GATEWAY_EQUINIX")
	skipGwSelfManaged := os.Getenv("SKIP_TRANSIT_GATEWAY_SELF_MANAGED")

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
						resource.TestCheckResourceAttr(resourceNameAws, "enable_jumbo_frame", "true"),
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
						resource.TestCheckResourceAttr(resourceNameAEP, "interfaces.0.logical_ifname", "wan0"),
						resource.TestCheckResourceAttr(resourceNameAEP, "ha_interfaces.0.gateway_ip", "192.168.20.1"),
						resource.TestCheckResourceAttr(resourceNameAEP, "ha_interfaces.0.ip_address", "192.168.20.12/24"),
						resource.TestCheckResourceAttr(resourceNameAEP, "ha_interfaces.0.logical_ifname", "wan0"),
						resource.TestCheckResourceAttr(resourceNameAEP, "peer_backup_logical_ifname", "wan1"),
						resource.TestCheckResourceAttr(resourceNameAEP, "peer_connection_type", "private"),
						resource.TestCheckResourceAttr(resourceNameAEP, "eip_map.0.logical_ifname", "wan0"),
						resource.TestCheckResourceAttr(resourceNameAEP, "eip_map.0.private_ip", "192.168.19.18"),
						resource.TestCheckResourceAttr(resourceNameAEP, "eip_map.0.public_ip", "35.0.16.1"),
						resource.TestCheckResourceAttr(resourceNameAEP, "bgp_polling_time", "55"),
						resource.TestCheckResourceAttr(resourceNameAEP, "local_as_number", "65000"),
						// enable_jumbo_frame is false by default for edge EAT gateways
						resource.TestCheckResourceAttr(resourceNameAEP, "enable_jumbo_frame", "false"),
						resource.TestCheckResourceAttr(resourceNameAEP, "tunnel_detection_time", "10"),
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

	if skipGwEQUINIX != "yes" {
		resourceNameEquinix := "aviatrix_transit_gateway.test_transit_gateway_equinix"
		msgCommonEquinix := ". Set SKIP_TRANSIT_GATEWAY_AEP to yes to skip Transit Gateway tests in edge AEP"

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheckEdge(t, msgCommonEquinix)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckTransitGatewayDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccTransitGatewayConfigBasicEquinix(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckTransitGatewayExists(resourceNameEquinix, &gateway),
						resource.TestCheckResourceAttr(resourceNameEquinix, "gw_name", fmt.Sprintf("tfg-aep-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameEquinix, "gw_size", "SMALL"),
						resource.TestCheckResourceAttr(resourceNameEquinix, "vpc_id", os.Getenv("AEP_VPC_ID")),
						resource.TestCheckResourceAttr(resourceNameEquinix, "ztp_file_download_path", "/tmp"),
						resource.TestCheckResourceAttr(resourceNameEquinix, "interfaces.0.gateway_ip", "192.168.20.1"),
						resource.TestCheckResourceAttr(resourceNameEquinix, "interfaces.0.ip_address", "192.168.20.11/24"),
						resource.TestCheckResourceAttr(resourceNameEquinix, "interfaces.0.logical_ifname", "wan0"),
						resource.TestCheckResourceAttr(resourceNameEquinix, "ha_interfaces.0.gateway_ip", "192.168.20.1"),
						resource.TestCheckResourceAttr(resourceNameEquinix, "ha_interfaces.0.ip_address", "192.168.20.12/24"),
						resource.TestCheckResourceAttr(resourceNameEquinix, "ha_interfaces.0.logical_ifname", "wan0"),
						resource.TestCheckResourceAttr(resourceNameEquinix, "peer_backup_logical_ifname", "wan1"),
						resource.TestCheckResourceAttr(resourceNameEquinix, "peer_connection_type", "private"),
						resource.TestCheckResourceAttr(resourceNameEquinix, "enable_jumbo_frame", "false"),
					),
				},
				{
					ResourceName:      resourceNameEquinix,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	} else {
		t.Log("Skipping Transit gateway test in edge Equinix as SKIP_TRANSIT_GATEWAY_EQUINIX is set")
	}

	if skipGwSelfManaged != "yes" {
		resourceNameSelfManaged := "aviatrix_transit_gateway.test_transit_gateway_selfmanaged"
		msgCommonSelfManaged := ". Set SKIP_TRANSIT_GATEWAY_SELF_MANAGED to yes to skip Transit Gateway tests in edge SelfManaged"

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
				preGatewayCheckEdge(t, msgCommonSelfManaged)
			},
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckTransitGatewayDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccTransitGatewayConfigBasicSelfManaged(rName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckTransitGatewayExists(resourceNameSelfManaged, &gateway),
						resource.TestCheckResourceAttr(resourceNameSelfManaged, "gw_name", fmt.Sprintf("tfg-edge-self-%s", rName)),
						resource.TestCheckResourceAttr(resourceNameSelfManaged, "vpc_id", os.Getenv("SEFLMANAGED_VPC_ID")),
						resource.TestCheckResourceAttr(resourceNameSelfManaged, "ztp_file_download_path", "/tmp"),
						resource.TestCheckResourceAttr(resourceNameSelfManaged, "interfaces.0.gateway_ip", "192.168.20.1"),
						resource.TestCheckResourceAttr(resourceNameSelfManaged, "interfaces.0.ip_address", "192.168.20.11/24"),
						resource.TestCheckResourceAttr(resourceNameSelfManaged, "interfaces.0.logical_ifname", "wan0"),
						resource.TestCheckResourceAttr(resourceNameSelfManaged, "ztp_file_type", "iso"),
					),
				},
				{
					ResourceName:      resourceNameSelfManaged,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	} else {
		t.Log("Skipping Transit gateway test in edge Self Managed as SKIP_TRANSIT_GATEWAY_SELF_MANAGED is set")
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
        gateway_ip     = "192.168.20.1"
        ip_address     = "192.168.20.11/24"
        logical_ifname = "wan0"
        secondary_private_cidr_list = ["192.168.19.16/29"]
    }

    interfaces {
        gateway_ip     = "192.168.21.1"
        ip_address     = "192.168.21.11/24"
        logical_ifname = "wan1"
        secondary_private_cidr_list = ["192.168.21.16/29"]
    }

    interfaces {
        dhcp   = true
        logical_ifname = "mgmt0"
    }

    interfaces {
        gateway_ip     = "192.168.22.1"
        ip_address     = "192.168.22.11/24"
        logical_ifname = "wan2"
    }

    interfaces {
        gateway_ip     = "192.168.23.1"
        ip_address     = "192.168.23.11/24"
        logical_ifname = "wan3"
    }

	ha_device_id = "a20c75c0-06c2-4102-9df1-b00b85e89eac"
    peer_backup_logical_ifname = "wan1"
    peer_connection_type = "private"
    ha_interfaces {
        gateway_ip    = "192.168.20.1"
        ip_address    = "192.168.20.12/24"
		logical_ifname = "wan0"
    }

    ha_interfaces {
        gateway_ip     = "192.168.21.1"
        ip_address     = "192.168.21.12/24"
        logical_ifname = "wan1"
        secondary_private_cidr_list = ["192.168.21.32/29"]
    }

    ha_interfaces {
        dhcp           = true
        logical_ifname = "mgmt0"
    }

    ha_interfaces {
        gateway_ip     = "192.168.22.1"
        ip_address     = "192.168.22.12/24"
        logical_ifname = "wan2"
    }

    ha_interfaces {
        gateway_ip     = "192.168.23.1"
        ip_address     = "192.168.23.12/24"
        logical_ifname = "wan3"
    }

	eip_map {
        logical_ifname = "wan0"
        private_ip     = "192.168.19.18"
        public_ip      = "35.0.16.1"
    }
	bgp_polling_time   = 55
	local_as_number    = 65000
	tunnel_detection_time = 10
}
	`, rName, os.Getenv("AEP_VPC_ID"), os.Getenv("AEP_DEVICE_ID"))
}

func testAccTransitGatewayConfigBasicEquinix(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_edge_equinix" {
	account_name       = "edge-%s"
	cloud_type         = 524288
}
resource "aviatrix_transit_gateway" "test_transit_gateway_equinix" {
	cloud_type   = 524288
	account_name = aviatrix_account.test_acc_edge_equinix.account_name
	gw_name      = "tfg-edge-equinix-%[1]s"
	vpc_id       = "%[2]s"
	gw_size      = "SMALL"
	ztp_file_download_path = "/tmp"
	interfaces {
        gateway_ip     = "192.168.20.1"
        ip_address     = "192.168.20.11/24"
        logical_ifname = "wan0"
        secondary_private_cidr_list = ["192.168.19.16/29"]
    }

    interfaces {
        gateway_ip     = "192.168.21.1"
        ip_address     = "192.168.21.11/24"
        logical_ifname = "wan1"
        secondary_private_cidr_list = ["192.168.21.16/29"]
    }

    interfaces {
        dhcp           = true
        logical_ifname = "mgmt0"
    }

    interfaces {
        gateway_ip     = "192.168.22.1"
        ip_address     = "192.168.22.11/24"
        logical_ifname = "wan2"
    }

    interfaces {
        gateway_ip     = "192.168.23.1"
        ip_address     = "192.168.23.11/24"
        logical_ifname = "wan3"
    }

    peer_backup_logical_ifname = "wan1"
    peer_connection_type = "private"
    ha_interfaces {
        gateway_ip    = "192.168.20.1"
        ip_address    = "192.168.20.12/24"
		logical_ifname = "wan0"
    }

    ha_interfaces {
        gateway_ip   = "192.168.21.1"
        ip_address   = "192.168.21.12/24"
		logical_ifname = "wan1"
        secondary_private_cidr_list = ["192.168.21.32/29"]
    }

    ha_interfaces {
        dhcp   = true
		logical_ifname = "mgmt0"
    }

    ha_interfaces {
        gateway_ip   = "192.168.22.1"
        ip_address   = "192.168.22.12/24"
		logical_ifname = "wan2"
    }

    ha_interfaces {
        gateway_ip     = "192.168.23.1"
        ip_address     = "192.168.23.12/24"
		logical_ifname = "wan3"
    }
}
	`, rName, os.Getenv("EQUINIX_VPC_ID"))
}

func testAccTransitGatewayConfigBasicSelfManaged(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_edge_selfmanaged" {
	account_name       = "edge-%s"
	cloud_type         = 4096
}
resource "aviatrix_transit_gateway" "test_transit_gateway_selfmanaged" {
	cloud_type   = 4096
	account_name = aviatrix_account.test_acc_edge_selmanaged.account_name
	gw_name      = "tfg-edge-self-%[1]s"
	vpc_id       = "%[2]s"
	gw_size      = ""
	ztp_file_download_path = "/tmp"
	ztp_file_type = "iso"
	interfaces {
        gateway_ip     = "192.168.20.1"
        ip_address     = "192.168.20.11/24"
        logical_ifname = "wan0"
        secondary_private_cidr_list = ["192.168.19.16/29"]
    }

    interfaces {
        gateway_ip     = "192.168.21.1"
        ip_address     = "192.168.21.11/24"
        logical_ifname = "wan1"
        secondary_private_cidr_list = ["192.168.21.16/29"]
    }

    interfaces {
        dhcp           = true
        logical_ifname = "mgmt0"
    }
}
	`, rName, os.Getenv("SELFMANAGED_VPC_ID"))
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

// TestAccAviatrixTransitGateway_ipv6AWS tests IPv6 CIDR fields for AWS transit gateway
func TestAccAviatrixTransitGateway_ipv6AWS(t *testing.T) {
	var gateway goaviatrix.Gateway

	rName := acctest.RandString(5)
	resourceName := "aviatrix_transit_gateway.test_transit_gateway_ipv6"

	msgCommon := ". Set SKIP_TRANSIT_GATEWAY_IPV6 to yes to skip Transit Gateway IPv6 tests"

	skipGwIPv6 := os.Getenv("SKIP_TRANSIT_GATEWAY_IPV6")
	if skipGwIPv6 == "yes" {
		t.Skip("Skipping Transit Gateway IPv6 test as SKIP_TRANSIT_GATEWAY_IPV6 is set")
	}

	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t2.micro"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAwsTransitGatewayIPv6Check(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayConfigAWSIPv6(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-aws-transit-ipv6-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "gw_size", awsGwSize),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-aws-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AWS_VPC_ID4")),
					resource.TestCheckResourceAttr(resourceName, "subnet", os.Getenv("AWS_SUBNET4")),
					resource.TestCheckResourceAttr(resourceName, "vpc_reg", os.Getenv("AWS_REGION")),
					resource.TestCheckResourceAttr(resourceName, "enable_ipv6", "true"),
					resource.TestCheckResourceAttr(resourceName, "subnet_ipv6_cidr", os.Getenv("AWS_SUBNET_IPV6_CIDR")),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"gcloud_project_credentials_filepath",
					"vnet_and_resource_group_names",
				},
			},
		},
	})
}

// TestAccAviatrixTransitGateway_ipv6WithHA tests IPv6 CIDR fields with HA enabled
func TestAccAviatrixTransitGateway_ipv6WithHA(t *testing.T) {
	var gateway goaviatrix.Gateway

	rName := acctest.RandString(5)
	resourceName := "aviatrix_transit_gateway.test_transit_gateway_ipv6_ha"

	msgCommon := ". Set SKIP_TRANSIT_GATEWAY_IPV6 to yes to skip Transit Gateway IPv6 tests"

	skipGwIPv6 := os.Getenv("SKIP_TRANSIT_GATEWAY_IPV6")
	if skipGwIPv6 == "yes" {
		t.Skip("Skipping Transit Gateway IPv6 HA test as SKIP_TRANSIT_GATEWAY_IPV6 is set")
	}

	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t2.micro"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAwsTransitGatewayIPv6HACheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayConfigAWSIPv6WithHA(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-aws-transit-ipv6-ha-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "gw_size", awsGwSize),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-aws-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AWS_VPC_ID4")),
					resource.TestCheckResourceAttr(resourceName, "subnet", os.Getenv("AWS_SUBNET4")),
					resource.TestCheckResourceAttr(resourceName, "vpc_reg", os.Getenv("AWS_REGION")),
					resource.TestCheckResourceAttr(resourceName, "enable_ipv6", "true"),
					resource.TestCheckResourceAttr(resourceName, "subnet_ipv6_cidr", os.Getenv("AWS_SUBNET_IPV6_CIDR")),
					resource.TestCheckResourceAttr(resourceName, "ha_subnet", os.Getenv("AWS_HA_SUBNET")),
					resource.TestCheckResourceAttr(resourceName, "ha_subnet_ipv6_cidr", os.Getenv("AWS_HA_SUBNET_IPV6_CIDR")),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"gcloud_project_credentials_filepath",
					"vnet_and_resource_group_names",
				},
			},
		},
	})
}

// TestAccAviatrixTransitGateway_ipv6Azure tests IPv6 CIDR fields for Azure transit gateway
func TestAccAviatrixTransitGateway_ipv6Azure(t *testing.T) {
	var gateway goaviatrix.Gateway

	rName := acctest.RandString(5)
	resourceName := "aviatrix_transit_gateway.test_transit_gateway_ipv6_azure"

	msgCommon := ". Set SKIP_TRANSIT_GATEWAY_IPV6_AZURE to yes to skip Azure Transit Gateway IPv6 tests"

	skipGwIPv6Azure := os.Getenv("SKIP_TRANSIT_GATEWAY_IPV6_AZURE")
	if skipGwIPv6Azure == "yes" {
		t.Skip("Skipping Azure Transit Gateway IPv6 test as SKIP_TRANSIT_GATEWAY_IPV6_AZURE is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAzureTransitGatewayIPv6Check(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayConfigAzureIPv6(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-azure-transit-ipv6-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "gw_size", os.Getenv("AZURE_GW_SIZE")),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-azure-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AZURE_VNET_ID")),
					resource.TestCheckResourceAttr(resourceName, "subnet", os.Getenv("AZURE_SUBNET")),
					resource.TestCheckResourceAttr(resourceName, "vpc_reg", os.Getenv("AZURE_REGION")),
					resource.TestCheckResourceAttr(resourceName, "enable_ipv6", "true"),
					resource.TestCheckResourceAttr(resourceName, "subnet_ipv6_cidr", os.Getenv("AZURE_SUBNET_IPV6_CIDR")),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"gcloud_project_credentials_filepath",
					"vnet_and_resource_group_names",
					"vpc_id",
				},
			},
		},
	})
}

// TestAccAviatrixTransitGateway_ipv6GCP tests IPv6 for GCP transit gateway (no subnet_ipv6_cidr required)
func TestAccAviatrixTransitGateway_ipv6GCP(t *testing.T) {
	var gateway goaviatrix.Gateway

	rName := acctest.RandString(5)
	resourceName := "aviatrix_transit_gateway.test_transit_gateway_ipv6_gcp"

	msgCommon := ". Set SKIP_TRANSIT_GATEWAY_IPV6_GCP to yes to skip GCP Transit Gateway IPv6 tests"

	skipGwIPv6GCP := os.Getenv("SKIP_TRANSIT_GATEWAY_IPV6_GCP")
	if skipGwIPv6GCP == "yes" {
		t.Skip("Skipping GCP Transit Gateway IPv6 test as SKIP_TRANSIT_GATEWAY_IPV6_GCP is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGCPTransitGatewayIPv6Check(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayConfigGCPIPv6(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-gcp-ipv6-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "gw_size", "n1-standard-1"),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-gcp-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("GCP_VPC_ID")),
					resource.TestCheckResourceAttr(resourceName, "vpc_reg", os.Getenv("GCP_ZONE")),
					resource.TestCheckResourceAttr(resourceName, "enable_ipv6", "true"),
					// For GCP, subnet_ipv6_cidr should be computed/optional, not required
					resource.TestCheckResourceAttrSet(resourceName, "subnet_ipv6_cidr"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"gcloud_project_credentials_filepath",
					"vnet_and_resource_group_names",
				},
			},
		},
	})
}

// TestAccAviatrixTransitGateway_ipv6GCPWithInsaneMode tests IPv6 with GCP and insane mode (no subnet_ipv6_cidr required)
func TestAccAviatrixTransitGateway_ipv6GCPWithInsaneMode(t *testing.T) {
	var gateway goaviatrix.Gateway

	rName := acctest.RandString(5)
	resourceName := "aviatrix_transit_gateway.test_transit_gateway_ipv6_gcp_insane"

	msgCommon := ". Set SKIP_TRANSIT_GATEWAY_IPV6_GCP to yes to skip GCP Transit Gateway IPv6 tests"

	skipGwIPv6GCP := os.Getenv("SKIP_TRANSIT_GATEWAY_IPV6_GCP")
	if skipGwIPv6GCP == "yes" {
		t.Skip("Skipping GCP Transit Gateway IPv6 with insane mode test as SKIP_TRANSIT_GATEWAY_IPV6_GCP is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preGCPTransitGatewayIPv6Check(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayConfigGCPIPv6WithInsaneMode(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-gcp-ipv6-insane-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "gw_size", "n1-highmem-4"),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-gcp-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("GCP_VPC_ID")),
					resource.TestCheckResourceAttr(resourceName, "vpc_reg", os.Getenv("GCP_ZONE")),
					resource.TestCheckResourceAttr(resourceName, "enable_ipv6", "true"),
					resource.TestCheckResourceAttr(resourceName, "insane_mode", "true"),
					// For GCP, subnet_ipv6_cidr should not be required even with insane mode
					resource.TestCheckResourceAttrSet(resourceName, "subnet_ipv6_cidr"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"gcloud_project_credentials_filepath",
					"vnet_and_resource_group_names",
				},
			},
		},
	})
}

func preAwsTransitGatewayIPv6Check(t *testing.T, msgCommon string) {
	requiredEnvVars := []string{
		"AWS_VPC_ID4",
		"AWS_SUBNET4",
		"AWS_REGION",
		"AWS_SUBNET_IPV6_CIDR",
	}
	for _, v := range requiredEnvVars {
		if os.Getenv(v) == "" {
			t.Fatalf("Env Var %s required %s", v, msgCommon)
		}
	}
}

func preAwsTransitGatewayIPv6HACheck(t *testing.T, msgCommon string) {
	requiredEnvVars := []string{
		"AWS_VPC_ID4",
		"AWS_SUBNET4",
		"AWS_REGION",
		"AWS_SUBNET_IPV6_CIDR",
		"AWS_HA_SUBNET",
		"AWS_HA_SUBNET_IPV6_CIDR",
	}
	for _, v := range requiredEnvVars {
		if os.Getenv(v) == "" {
			t.Fatalf("Env Var %s required %s", v, msgCommon)
		}
	}
}

func preAzureTransitGatewayIPv6Check(t *testing.T, msgCommon string) {
	requiredEnvVars := []string{
		"AZURE_VNET_ID",
		"AZURE_SUBNET",
		"AZURE_REGION",
		"AZURE_GW_SIZE",
		"AZURE_SUBNET_IPV6_CIDR",
	}
	for _, v := range requiredEnvVars {
		if os.Getenv(v) == "" {
			t.Fatalf("Env Var %s required %s", v, msgCommon)
		}
	}
}

func preGCPTransitGatewayIPv6Check(t *testing.T, msgCommon string) {
	requiredEnvVars := []string{
		"GCP_PROJECT_ID",
		"GCP_VPC_ID",
		"GCP_ZONE",
		"GCP_SUBNET",
		"GOOGLE_CREDENTIALS_FILEPATH",
	}
	for _, v := range requiredEnvVars {
		if os.Getenv(v) == "" {
			t.Fatalf("Env Var %s required %s", v, msgCommon)
		}
	}
}

func testAccTransitGatewayConfigAWSIPv6(rName string) string {
	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t2.micro"
	}
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_transit_gateway" "test_transit_gateway_ipv6" {
	cloud_type        = 1
	account_name      = aviatrix_account.test_acc_aws.account_name
	gw_name           = "tfg-aws-transit-ipv6-%[1]s"
	vpc_id            = "%[5]s"
	vpc_reg           = "%[6]s"
	gw_size           = "%[7]s"
	subnet            = "%[8]s"
	enable_ipv6       = true
	subnet_ipv6_cidr  = "%[9]s"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID4"), os.Getenv("AWS_REGION"), awsGwSize, os.Getenv("AWS_SUBNET4"),
		os.Getenv("AWS_SUBNET_IPV6_CIDR"))
}

func testAccTransitGatewayConfigAWSIPv6WithHA(rName string) string {
	awsGwSize := os.Getenv("AWS_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "t2.micro"
	}
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_transit_gateway" "test_transit_gateway_ipv6_ha" {
	cloud_type           = 1
	account_name         = aviatrix_account.test_acc_aws.account_name
	gw_name              = "tfg-aws-transit-ipv6-ha-%[1]s"
	vpc_id               = "%[5]s"
	vpc_reg              = "%[6]s"
	gw_size              = "%[7]s"
	subnet               = "%[8]s"
	enable_ipv6          = true
	subnet_ipv6_cidr     = "%[9]s"
	ha_subnet            = "%[10]s"
	ha_subnet_ipv6_cidr  = "%[11]s"
	ha_gw_size           = "%[7]s"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID4"), os.Getenv("AWS_REGION"), awsGwSize, os.Getenv("AWS_SUBNET4"),
		os.Getenv("AWS_SUBNET_IPV6_CIDR"), os.Getenv("AWS_HA_SUBNET"), os.Getenv("AWS_HA_SUBNET_IPV6_CIDR"))
}

func testAccTransitGatewayConfigAzureIPv6(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_azure" {
	account_name        = "tfa-azure-%s"
	cloud_type          = 8
	arm_subscription_id = "%s"
	arm_directory_id    = "%s"
	arm_application_id  = "%s"
	arm_application_key = "%s"
}
resource "aviatrix_transit_gateway" "test_transit_gateway_ipv6_azure" {
	cloud_type       = 8
	account_name     = aviatrix_account.test_acc_azure.account_name
	gw_name          = "tfg-azure-transit-ipv6-%[1]s"
	vpc_id           = "%[6]s"
	vpc_reg          = "%[7]s"
	gw_size          = "%[8]s"
	subnet           = "%[9]s"
	enable_ipv6      = true
	subnet_ipv6_cidr = "%[10]s"
}
	`, rName, os.Getenv("ARM_SUBSCRIPTION_ID"), os.Getenv("ARM_DIRECTORY_ID"),
		os.Getenv("ARM_APPLICATION_ID"), os.Getenv("ARM_APPLICATION_KEY"),
		os.Getenv("AZURE_VNET_ID"), os.Getenv("AZURE_REGION"),
		os.Getenv("AZURE_GW_SIZE"), os.Getenv("AZURE_SUBNET"), os.Getenv("AZURE_SUBNET_IPV6_CIDR"))
}

// TestAccAviatrixTransitGateway_ipv6WithInsaneMode tests IPv6 with Insane Mode enabled
func TestAccAviatrixTransitGateway_ipv6WithInsaneMode(t *testing.T) {
	var gateway goaviatrix.Gateway

	rName := acctest.RandString(5)
	resourceName := "aviatrix_transit_gateway.test_transit_gateway_ipv6_insane"

	msgCommon := ". Set SKIP_TRANSIT_GATEWAY_IPV6_INSANE_MODE to yes to skip Transit Gateway IPv6 Insane Mode tests"

	skipGwIPv6InsaneMode := os.Getenv("SKIP_TRANSIT_GATEWAY_IPV6_INSANE_MODE")
	if skipGwIPv6InsaneMode == "yes" {
		t.Skip("Skipping Transit Gateway IPv6 Insane Mode test as SKIP_TRANSIT_GATEWAY_IPV6_INSANE_MODE is set")
	}

	awsGwSize := os.Getenv("AWS_INSANE_MODE_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "c5.xlarge"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAwsTransitGatewayIPv6InsaneModeCheck(t, msgCommon)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTransitGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayConfigAWSIPv6WithInsaneMode(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayExists(resourceName, &gateway),
					resource.TestCheckResourceAttr(resourceName, "gw_name", fmt.Sprintf("tfg-aws-transit-ipv6-insane-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "gw_size", awsGwSize),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("tfa-aws-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", os.Getenv("AWS_VPC_ID4")),
					resource.TestCheckResourceAttr(resourceName, "vpc_reg", os.Getenv("AWS_REGION")),
					resource.TestCheckResourceAttr(resourceName, "enable_ipv6", "true"),
					resource.TestCheckResourceAttr(resourceName, "subnet_ipv6_cidr", os.Getenv("AWS_SUBNET_IPV6_CIDR")),
					resource.TestCheckResourceAttr(resourceName, "insane_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "insane_mode_az", os.Getenv("AWS_AVAILABILITY_ZONE")),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"gcloud_project_credentials_filepath",
					"vnet_and_resource_group_names",
				},
			},
		},
	})
}

func preAwsTransitGatewayIPv6InsaneModeCheck(t *testing.T, msgCommon string) {
	requiredEnvVars := []string{
		"AWS_VPC_ID4",
		"AWS_REGION",
		"AWS_SUBNET_IPV6_CIDR",
		"AWS_AVAILABILITY_ZONE",
		"AWS_INSANE_MODE_SUBNET",
	}
	for _, v := range requiredEnvVars {
		if os.Getenv(v) == "" {
			t.Fatalf("Env Var %s required %s", v, msgCommon)
		}
	}
}

func testAccTransitGatewayConfigAWSIPv6WithInsaneMode(rName string) string {
	awsGwSize := os.Getenv("AWS_INSANE_MODE_GW_SIZE")
	if awsGwSize == "" {
		awsGwSize = "c5.xlarge"
	}
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_aws" {
	account_name       = "tfa-aws-%s"
	cloud_type         = 1
	aws_account_number = "%s"
	aws_iam            = false
	aws_access_key     = "%s"
	aws_secret_key     = "%s"
}
resource "aviatrix_transit_gateway" "test_transit_gateway_ipv6_insane" {
	cloud_type        = 1
	account_name      = aviatrix_account.test_acc_aws.account_name
	gw_name           = "tfg-aws-transit-ipv6-insane-%[1]s"
	vpc_id            = "%[5]s"
	vpc_reg           = "%[6]s"
	gw_size           = "%[7]s"
	subnet            = "%[8]s"
	enable_ipv6       = true
	subnet_ipv6_cidr  = "%[9]s"
	insane_mode       = true
	insane_mode_az    = "%[10]s"
}
	`, rName, os.Getenv("AWS_ACCOUNT_NUMBER"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"),
		os.Getenv("AWS_VPC_ID4"), os.Getenv("AWS_REGION"), awsGwSize, os.Getenv("AWS_INSANE_MODE_SUBNET"),
		os.Getenv("AWS_SUBNET_IPV6_CIDR"), os.Getenv("AWS_AVAILABILITY_ZONE"))
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
					"type":  "MANAGEMENT",
					"index": 0,
				},
				map[string]interface{}{
					"name":  "eth1",
					"type":  "WAN",
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
			intfType:    "wan0",
			wanCount:    3,
			expected:    "eth0",
			expectedErr: nil,
		},
		{
			name:        "Valid WAN interface with index 1",
			intfType:    "wan1",
			wanCount:    3,
			expected:    "eth1",
			expectedErr: nil,
		},
		{
			name:        "Valid WAN interface with index 2",
			intfType:    "wan2",
			wanCount:    3,
			expected:    "eth3",
			expectedErr: nil,
		},
		{
			name:        "Valid MANAGEMENT interface with index 0",
			intfType:    "mgmt0",
			wanCount:    3,
			expected:    "eth2",
			expectedErr: nil,
		},
		{
			name:        "Valid MANAGEMENT interface with index 1",
			intfType:    "mgmt1",
			wanCount:    3,
			expected:    "eth4",
			expectedErr: nil,
		},
		{
			name:        "Invalid interface type",
			intfType:    "INVALID",
			wanCount:    3,
			expected:    "",
			expectedErr: errors.New("invalid logical interface name: INVALID"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := getInterfaceName(tt.intfType, tt.wanCount)

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
		cloudType   int
		expected    map[string][]goaviatrix.EipMap
		expectedErr error
	}{
		{
			name: "Valid EIP map with WAN and MANAGEMENT interfaces for AEP",
			eipMap: []interface{}{
				map[string]interface{}{
					"logical_ifname": "wan0",
					"private_ip":     "192.168.0.10",
					"public_ip":      "203.0.113.10",
				},
				map[string]interface{}{
					"logical_ifname": "mgmt0",
					"private_ip":     "192.168.1.10",
					"public_ip":      "203.0.113.11",
				},
			},
			wanCount:  3,
			cloudType: 262144,
			expected: map[string][]goaviatrix.EipMap{
				"eth0": {
					{PrivateIP: "192.168.0.10", PublicIP: "203.0.113.10"},
				},
				"eth2": {
					{PrivateIP: "192.168.1.10", PublicIP: "203.0.113.11"},
				},
			},
			expectedErr: nil,
		},
		{
			name: "Valid EIP map with WAN and MANAGEMENT interfaces for Megaport",
			eipMap: []interface{}{
				map[string]interface{}{
					"logical_ifname": "wan0",
					"private_ip":     "192.168.0.10",
					"public_ip":      "203.0.113.10",
				},
				map[string]interface{}{
					"logical_ifname": "mgmt0",
					"private_ip":     "192.168.1.10",
					"public_ip":      "203.0.113.11",
				},
			},
			wanCount:  3,
			cloudType: 1048576,
			expected: map[string][]goaviatrix.EipMap{
				"wan0": {
					{PrivateIP: "192.168.0.10", PublicIP: "203.0.113.10"},
				},
				"mgmt0": {
					{PrivateIP: "192.168.1.10", PublicIP: "203.0.113.11"},
				},
			},
			expectedErr: nil,
		},
		{
			name: "Invalid EIP map: missing logical interface name",
			eipMap: []interface{}{
				map[string]interface{}{
					"private_ip": "192.168.0.10",
					"public_ip":  "203.0.113.10",
				},
			},
			wanCount:    3,
			cloudType:   1048576,
			expected:    nil,
			expectedErr: errors.New("logical interface name must be a string"),
		},
		{
			name: "Invalid EIP map: invalid logical interface name",
			eipMap: []interface{}{
				map[string]interface{}{
					"logical_ifname": 123,
					"private_ip":     "192.168.0.10",
					"public_ip":      "203.0.113.10",
				},
			},
			wanCount:    3,
			cloudType:   1048576,
			expected:    nil,
			expectedErr: errors.New("logical interface name must be a string"),
		},
		{
			name:        "Empty EIP map",
			eipMap:      []interface{}{},
			wanCount:    3,
			cloudType:   1048576,
			expected:    map[string][]goaviatrix.EipMap{},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := getEipMapDetails(tt.eipMap, tt.wanCount, tt.cloudType)

			// Check for errors
			if err != nil && tt.expectedErr != nil {
				if err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
				}
			} else if (err != nil && tt.expectedErr == nil) || (err == nil && tt.expectedErr != nil) {
				t.Errorf("unexpected error: %v", err)
			}

			// Check result
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected result: %v, got: %v", tt.expected, result)
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
	expectedWANCount := 5
	if wanCount != expectedWANCount {
		t.Errorf("Expected %d WAN interfaces, got %d", expectedWANCount, wanCount)
	}
}

// test to get the interface details from the resource
func TestGetInterfaceDetails(t *testing.T) {
	// get the interface details
	cloudType := 1048576
	interfaceDetails, err := getInterfaceDetails(interfaces, cloudType)
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

func TestSetEipMapDetails(t *testing.T) {
	tests := []struct {
		name              string
		eipMap            map[string][]goaviatrix.EipMap
		ifNameTranslation map[string]string
		expectedResult    []map[string]interface{}
		expectedErr       string
	}{
		{
			name: "Valid input",
			eipMap: map[string][]goaviatrix.EipMap{
				"eth0": {
					{PrivateIP: "192.168.1.10", PublicIP: "34.123.45.67"},
				},
				"eth1": {
					{PrivateIP: "192.168.1.11", PublicIP: "34.123.45.68"},
					{PrivateIP: "192.168.1.12", PublicIP: "34.123.45.69"},
				},
			},
			ifNameTranslation: map[string]string{
				"eth0": "WAN.0",
				"eth1": "WAN.1",
			},
			expectedResult: []map[string]interface{}{
				{
					"logical_ifname": "wan0",
					"private_ip":     "192.168.1.10",
					"public_ip":      "34.123.45.67",
				},
				{
					"logical_ifname": "wan1",
					"private_ip":     "192.168.1.11",
					"public_ip":      "34.123.45.68",
				},
				{
					"logical_ifname": "wan1",
					"private_ip":     "192.168.1.12",
					"public_ip":      "34.123.45.69",
				},
			},
			expectedErr: "",
		},
		{
			name: "Error converting interface index",
			eipMap: map[string][]goaviatrix.EipMap{
				"eth0": {
					{PrivateIP: "192.168.1.10", PublicIP: "34.123.45.67"},
				},
			},
			ifNameTranslation: map[string]string{
				"eth0": "WAN.invalid_index",
			},
			expectedResult: nil,
			expectedErr:    "failed to convert interface index to integer",
		},
		{
			name:              "Empty EIP map",
			eipMap:            map[string][]goaviatrix.EipMap{},
			ifNameTranslation: map[string]string{},
			expectedResult:    []map[string]interface{}{},
			expectedErr:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := setEipMapDetails(tt.eipMap, tt.ifNameTranslation)

			if tt.expectedErr != "" {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestSetInterfaceMappingDetails(t *testing.T) {
	tests := []struct {
		name              string
		interfaceMapping  []goaviatrix.InterfaceMapping
		expectedResult    []map[string]interface{}
		expectedOrderFunc func([]map[string]interface{}) bool
	}{
		{
			name: "Valid interface mapping with multiple interfaces",
			interfaceMapping: []goaviatrix.InterfaceMapping{
				{Name: "eth0", Type: "WAN", Index: 0},
				{Name: "eth1", Type: "MANAGEMENT", Index: 1},
				{Name: "eth2", Type: "WAN", Index: 2},
			},
			expectedResult: []map[string]interface{}{
				{"name": "eth0", "type": "WAN", "index": 0},
				{"name": "eth1", "type": "MANAGEMENT", "index": 1},
				{"name": "eth2", "type": "WAN", "index": 2},
			},
			expectedOrderFunc: func(result []map[string]interface{}) bool {
				// Check the order based on "name" for sorting
				return result[0]["name"] == "eth0" && result[1]["name"] == "eth1" && result[2]["name"] == "eth2"
			},
		},
		{
			name:             "Empty interface mapping",
			interfaceMapping: []goaviatrix.InterfaceMapping{},
			expectedResult:   []map[string]interface{}{},
			expectedOrderFunc: func(result []map[string]interface{}) bool {
				return len(result) == 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := setInterfaceMappingDetails(tt.interfaceMapping)
			assert.Equal(t, tt.expectedResult, result)
			if tt.expectedOrderFunc != nil {
				assert.True(t, tt.expectedOrderFunc(result))
			}
		})
	}
}

func TestSetInterfaceDetails(t *testing.T) {
	// Define the interface order
	interfaceOrder := []string{"wan0", "wan1", "mgmt0", "wan2", "wan3"}

	// Define test cases
	tests := []struct {
		name       string
		interfaces []goaviatrix.EdgeTransitInterface
		expected   []map[string]interface{}
	}{
		{
			name: "Single WAN interface",
			interfaces: []goaviatrix.EdgeTransitInterface{
				{LogicalIfName: "wan0", PublicIp: "1.1.1.1", Dhcp: true, IpAddress: "10.0.0.1", GatewayIp: "10.0.0.254"},
			},
			expected: []map[string]interface{}{
				{
					"logical_ifname": "wan0",
					"public_ip":      "1.1.1.1",
					"dhcp":           true,
					"ip_address":     "10.0.0.1",
					"gateway_ip":     "10.0.0.254",
				},
			},
		},
		{
			name: "Multiple WAN and MANAGEMENT interfaces",
			interfaces: []goaviatrix.EdgeTransitInterface{
				{LogicalIfName: "wan2", GatewayIp: "192.168.1.1"},
				{LogicalIfName: "mgmt0", Dhcp: true},
				{LogicalIfName: "wan0", IpAddress: "10.0.0.2"},
				{LogicalIfName: "wan1", IpAddress: "10.0.0.3"},
			},
			expected: []map[string]interface{}{
				{"logical_ifname": "wan0", "ip_address": "10.0.0.2"},
				{"logical_ifname": "wan1", "ip_address": "10.0.0.3"},
				{"logical_ifname": "mgmt0", "dhcp": true},
				{"logical_ifname": "wan2", "gateway_ip": "192.168.1.1"},
			},
		},
		{
			name: "Custom interface with Secondary CIDRs",
			interfaces: []goaviatrix.EdgeTransitInterface{
				{
					LogicalIfName:  "mgmt0",
					SecondaryCIDRs: []string{"10.0.1.0/24", "10.0.2.0/24"},
				},
			},
			expected: []map[string]interface{}{
				{
					"logical_ifname":              "mgmt0",
					"secondary_private_cidr_list": []string{"10.0.1.0/24", "10.0.2.0/24"},
				},
			},
		},
		{
			name:       "Empty interface list",
			interfaces: []goaviatrix.EdgeTransitInterface{},
			expected:   []map[string]interface{}{},
		},
		{
			name: "Ignore empty SecondaryCIDRs",
			interfaces: []goaviatrix.EdgeTransitInterface{
				{
					LogicalIfName:  "wan0",
					SecondaryCIDRs: []string{"", "10.0.3.0/24", ""},
				},
			},
			expected: []map[string]interface{}{
				{
					"logical_ifname":              "wan0",
					"secondary_private_cidr_list": []string{"10.0.3.0/24"},
				},
			},
		},
		{
			name: "WAN interface with underlay CIDR",
			interfaces: []goaviatrix.EdgeTransitInterface{
				{LogicalIfName: "wan0", IpAddress: "192.168.1.10/24", GatewayIp: "169.254.100.1", UnderlayCidr: "169.254.100.2/30"},
			},
			expected: []map[string]interface{}{
				{
					"logical_ifname": "wan0",
					"ip_address":     "192.168.1.10/24",
					"gateway_ip":     "169.254.100.1", // Gateway IP within underlay_cidr subnet
					"underlay_cidr":  "169.254.100.2/30",
				},
			},
		},
		{
			name: "WAN interface with typical link-local underlay CIDR range",
			interfaces: []goaviatrix.EdgeTransitInterface{
				{LogicalIfName: "wan1", IpAddress: "10.0.1.10/24", GatewayIp: "169.254.100.1", UnderlayCidr: "169.254.100.2/28"},
			},
			expected: []map[string]interface{}{
				{
					"logical_ifname": "wan1",
					"ip_address":     "10.0.1.10/24",
					"gateway_ip":     "169.254.100.1", // Gateway IP within underlay_cidr subnet
					"underlay_cidr":  "169.254.100.2/28",
				},
			},
		},
		{
			name: "Multiple interfaces with and without underlay CIDR",
			interfaces: []goaviatrix.EdgeTransitInterface{
				{LogicalIfName: "wan0", IpAddress: "192.168.1.10/24", GatewayIp: "192.168.1.1"},
				{LogicalIfName: "wan1", IpAddress: "10.0.1.10/24", GatewayIp: "169.254.1.1", UnderlayCidr: "169.254.1.2/30"},
			},
			expected: []map[string]interface{}{
				{
					"logical_ifname": "wan0",
					"ip_address":     "192.168.1.10/24",
					"gateway_ip":     "192.168.1.1",
				},
				{
					"logical_ifname": "wan1",
					"ip_address":     "10.0.1.10/24",
					"gateway_ip":     "169.254.1.1", // Gateway IP within underlay_cidr subnet
					"underlay_cidr":  "169.254.1.2/30",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := setInterfaceDetails(tt.interfaces, interfaceOrder)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDeleteZtpFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ztp_test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Cleanup the temp directory after test
	gatewayName := "test-gateway"
	vpcID := "vpc-123456"
	ztpFileDownloadPath := tempDir
	fileName := filepath.Join(ztpFileDownloadPath, fmt.Sprintf("%s-%s-cloud-init.txt", gatewayName, vpcID))
	// Create a temporary file to simulate the ztp file
	if err := os.WriteFile(fileName, []byte("test content"), 0o644); err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	// Test Case 1: File exists and is successfully deleted
	// Ensure the file exists before calling the function
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		t.Fatalf("File does not exist before deletion: %v", fileName)
	}
	err = deleteZtpFile(gatewayName, vpcID, ztpFileDownloadPath)
	if err != nil {
		t.Errorf("deleteZtpFile returned an error: %v", err)
	}

	// Verify the file is deleted
	if _, err := os.Stat(fileName); err == nil || !os.IsNotExist(err) {
		t.Errorf("File was not deleted: %v", fileName)
	}

	// Test Case 2: File does not exist and no error should be returned
	// Try to delete the file again (it doesn't exist now)
	err = deleteZtpFile(gatewayName, vpcID, ztpFileDownloadPath)
	if err != nil {
		t.Errorf("deleteZtpFile returned an error when file does not exist: %v", err)
	}
}

func TestCreateBackupLinkConfig(t *testing.T) {
	tests := []struct {
		name                   string
		gwName                 string
		peerBackupLogicalNames []interface{}
		connectionType         string
		wanCount               int
		cloudType              int
		expected               string
		expectedErr            error
	}{
		{
			name:                   "Valid backup link config for AEP",
			gwName:                 "gw1",
			peerBackupLogicalNames: []interface{}{"wan0", "wan1"},
			connectionType:         "private",
			wanCount:               3,
			cloudType:              goaviatrix.EDGENEO,
			expected:               `[{"peer_gw_name":"gw1","peer_backup_port":"eth0,eth1","self_backup_port":"eth0,eth1","connection_type":"private"}]`,
			expectedErr:            nil,
		},
		{
			name:                   "Valid backup link config for Megaport",
			gwName:                 "gw2",
			peerBackupLogicalNames: []interface{}{"wan2", "wan3"},
			connectionType:         "public",
			wanCount:               4,
			cloudType:              goaviatrix.EDGEMEGAPORT,
			expected:               `[{"peer_gw_name":"gw2","connection_type":"public","peer_backup_logical_ifnames":["wan2","wan3"],"self_backup_logical_ifnames":["wan2","wan3"]}]`,
			expectedErr:            nil,
		},
		{
			name:                   "Invalid logical name in backup link config",
			gwName:                 "gw3",
			peerBackupLogicalNames: []interface{}{"wan0", "invalid_eth"},
			connectionType:         "private",
			wanCount:               3,
			cloudType:              goaviatrix.EDGENEO,
			expected:               "",
			expectedErr:            fmt.Errorf("failed to get the peer backup port name for logical name invalid_eth: invalid logical interface name: invalid_eth"),
		},
		{
			name:                   "Empty logical names in backup link config",
			gwName:                 "gw4",
			peerBackupLogicalNames: []interface{}{},
			connectionType:         "private",
			wanCount:               3,
			cloudType:              goaviatrix.EDGENEO,
			expected:               `[{"peer_gw_name":"gw4","connection_type":"private"}]`,
			expectedErr:            nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := createBackupLinkConfig(tt.gwName, tt.peerBackupLogicalNames, tt.connectionType, tt.wanCount, tt.cloudType)

			// Check for expected error
			if (err != nil || tt.expectedErr != nil) && (err == nil || err.Error() != tt.expectedErr.Error()) {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
			}

			// Check the result
			if result != tt.expected {
				t.Errorf("expected result: %s, got: %s", tt.expected, result)
			}
		})
	}
}

func TestCalculateInterfaceName(t *testing.T) {
	tests := []struct {
		name      string
		intfType  string
		intfIndex int
		wanCount  int
		expected  string
		expectErr bool
	}{
		{
			name:      "First WAN interface",
			intfType:  "WAN",
			intfIndex: 0,
			wanCount:  2,
			expected:  "eth0",
			expectErr: false,
		},
		{
			name:      "Second WAN interface",
			intfType:  "WAN",
			intfIndex: 1,
			wanCount:  2,
			expected:  "eth1",
			expectErr: false,
		},
		{
			name:      "Third WAN interface",
			intfType:  "WAN",
			intfIndex: 2,
			wanCount:  2,
			expected:  "eth3",
			expectErr: false,
		},
		{
			name:      "First MANAGEMENT interface",
			intfType:  "MANAGEMENT",
			intfIndex: 0,
			wanCount:  2,
			expected:  "eth2",
			expectErr: false,
		},
		{
			name:      "Invalid interface type",
			intfType:  "INVALID",
			intfIndex: 0,
			wanCount:  2,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculateInterfaceName(tt.intfType, tt.intfIndex, tt.wanCount)
			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %s, got %s", tt.expected, result)
				}
			}
		})
	}
}

func TestParseInterface(t *testing.T) {
	tests := []struct {
		name      string
		ifaceInfo map[string]interface{}
		wanCount  int
		cloudType int
		expected  goaviatrix.EdgeTransitInterface
		expectErr bool
	}{
		{
			name: "Valid WAN interface",
			ifaceInfo: map[string]interface{}{
				"logical_ifname":              "wan0",
				"gateway_ip":                  "192.168.1.1",
				"ip_address":                  "192.168.1.2",
				"public_ip":                   "203.0.113.1",
				"dhcp":                        true,
				"secondary_private_cidr_list": []interface{}{"10.0.0.0/16", "10.1.0.0/16"},
			},
			wanCount:  1,
			cloudType: goaviatrix.EDGEMEGAPORT,
			expected: goaviatrix.EdgeTransitInterface{
				GatewayIp:      "192.168.1.1",
				PublicIp:       "203.0.113.1",
				Dhcp:           true,
				IpAddress:      "192.168.1.2",
				SecondaryCIDRs: []string{"10.0.0.0/16", "10.1.0.0/16"},
				LogicalIfName:  "wan0",
			},
			expectErr: false,
		},
		{
			name: "WAN interface with underlay CIDR",
			ifaceInfo: map[string]interface{}{
				"logical_ifname": "wan1",
				"gateway_ip":     "169.254.100.1", // Gateway IP within underlay_cidr subnet
				"ip_address":     "192.168.2.10/24",
				"underlay_cidr":  "169.254.100.2/30", // underlay_cidr range: 169.254.100.0-169.254.100.3
			},
			wanCount:  2,
			cloudType: goaviatrix.EDGEMEGAPORT,
			expected: goaviatrix.EdgeTransitInterface{
				GatewayIp:     "169.254.100.1",
				IpAddress:     "192.168.2.10/24",
				UnderlayCidr:  "169.254.100.2/30",
				LogicalIfName: "wan1",
			},
			expectErr: false,
		},
		{
			name: "WAN interface with link-local underlay CIDR typical range",
			ifaceInfo: map[string]interface{}{
				"logical_ifname": "wan2",
				"gateway_ip":     "169.254.100.1", // Gateway IP within underlay_cidr subnet
				"ip_address":     "10.0.1.10/24",
				"underlay_cidr":  "169.254.100.2/28", // underlay_cidr range: 169.254.100.0-169.254.100.15
			},
			wanCount:  3,
			cloudType: goaviatrix.EDGEMEGAPORT,
			expected: goaviatrix.EdgeTransitInterface{
				GatewayIp:     "169.254.100.1",
				IpAddress:     "10.0.1.10/24",
				UnderlayCidr:  "169.254.100.2/28",
				LogicalIfName: "wan2",
			},
			expectErr: false,
		},
		{
			name: "WAN interface with point-to-point underlay CIDR",
			ifaceInfo: map[string]interface{}{
				"logical_ifname": "wan0",
				"gateway_ip":     "169.254.1.1", // Gateway IP within underlay_cidr subnet
				"ip_address":     "172.16.1.10/24",
				"underlay_cidr":  "169.254.1.2/30", // underlay_cidr range: 169.254.1.0-169.254.1.3
			},
			wanCount:  1,
			cloudType: goaviatrix.EDGEMEGAPORT,
			expected: goaviatrix.EdgeTransitInterface{
				GatewayIp:     "169.254.1.1",
				IpAddress:     "172.16.1.10/24",
				UnderlayCidr:  "169.254.1.2/30",
				LogicalIfName: "wan0",
			},
			expectErr: false,
		},
		{
			name: "Valid MANAGEMENT interface",
			ifaceInfo: map[string]interface{}{
				"logical_ifname": "mgmt0",
				"dhcp":           true,
			},
			wanCount:  1,
			cloudType: 0,
			expected: goaviatrix.EdgeTransitInterface{
				Dhcp: true,
				Name: "eth2",
				Type: "MANAGEMENT",
			},
			expectErr: false,
		},
		{
			name: "Invalid logical_ifname",
			ifaceInfo: map[string]interface{}{
				"logical_ifname": 12345, // Invalid type
			},
			wanCount:  1,
			cloudType: goaviatrix.EDGEMEGAPORT,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseInterface(tt.ifaceInfo, tt.wanCount, tt.cloudType)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestGetUserInterfaceOrder(t *testing.T) {
	tests := []struct {
		name        string
		interfaces  []interface{}
		expected    []string
		expectError bool
	}{
		{
			name: "Valid interface list",
			interfaces: []interface{}{
				map[string]interface{}{"logical_ifname": "wan0"},
				map[string]interface{}{"logical_ifname": "wan1"},
				map[string]interface{}{"logical_ifname": "mgmt0"},
			},
			expected:    []string{"wan0", "wan1", "mgmt0"},
			expectError: false,
		},
		{
			name:        "Empty interface list",
			interfaces:  []interface{}{},
			expected:    []string{},
			expectError: false,
		},
		{
			name: "Interface is not a map[string]interface{}",
			interfaces: []interface{}{
				"invalid_entry",
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "Missing logical_ifname key",
			interfaces: []interface{}{
				map[string]interface{}{"interface_type": "WAN"},
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "Logical_ifname is not a string",
			interfaces: []interface{}{
				map[string]interface{}{"logical_ifname": 12345},
			},
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := getUserInterfaceOrder(tt.interfaces)
			if tt.expectError {
				assert.Error(t, err, "Expected an error but got nil")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestSortInterfacesByCustomOrder(t *testing.T) {
	tests := []struct {
		name               string
		interfaces         []goaviatrix.EdgeTransitInterface
		userInterfaceOrder []string
		expected           []goaviatrix.EdgeTransitInterface
	}{
		{
			name: "Sort interfaces based on custom order",
			interfaces: []goaviatrix.EdgeTransitInterface{
				{LogicalIfName: "wan2"},
				{LogicalIfName: "wan0"},
				{LogicalIfName: "mgmt0"},
				{LogicalIfName: "wan1"},
			},
			userInterfaceOrder: []string{"wan0", "wan1", "mgmt0", "wan2"},
			expected: []goaviatrix.EdgeTransitInterface{
				{LogicalIfName: "wan0"},
				{LogicalIfName: "wan1"},
				{LogicalIfName: "mgmt0"},
				{LogicalIfName: "wan2"},
			},
		},
		{
			name: "Handles interfaces not in custom order",
			interfaces: []goaviatrix.EdgeTransitInterface{
				{LogicalIfName: "extra0"},
				{LogicalIfName: "wan1"},
				{LogicalIfName: "mgmt0"},
			},
			userInterfaceOrder: []string{"wan1", "mgmt0"},
			expected: []goaviatrix.EdgeTransitInterface{
				{LogicalIfName: "wan1"},
				{LogicalIfName: "mgmt0"},
				{LogicalIfName: "extra0"},
			},
		},
		{
			name: "Handles empty custom order",
			interfaces: []goaviatrix.EdgeTransitInterface{
				{LogicalIfName: "wan2"},
				{LogicalIfName: "wan0"},
			},
			userInterfaceOrder: []string{},
			expected: []goaviatrix.EdgeTransitInterface{
				{LogicalIfName: "wan2"},
				{LogicalIfName: "wan0"},
			},
		},
		{
			name:               "Handles empty interface list",
			interfaces:         []goaviatrix.EdgeTransitInterface{},
			userInterfaceOrder: []string{"wan0", "wan1"},
			expected:           []goaviatrix.EdgeTransitInterface{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sortInterfacesByCustomOrder(tt.interfaces, tt.userInterfaceOrder)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func testAccTransitGatewayConfigGCPIPv6(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_gcp" {
	account_name                        = "tfa-gcp-%s"
	cloud_type                          = 4
	gcloud_project_id                   = "%s"
	gcloud_project_credentials_filepath = "%s"
}
resource "aviatrix_transit_gateway" "test_transit_gateway_ipv6_gcp" {
	cloud_type       = 4
	account_name     = aviatrix_account.test_acc_gcp.account_name
	gw_name          = "tfg-gcp-ipv6-%[1]s"
	vpc_id           = "%[4]s"
	vpc_reg          = "%[5]s"
	gw_size          = "n1-standard-1"
	subnet           = "%[6]s"
	enable_ipv6      = true
	connected_transit = true
}
	`, rName, os.Getenv("GCP_PROJECT_ID"), os.Getenv("GOOGLE_CREDENTIALS_FILEPATH"),
		os.Getenv("GCP_VPC_ID"), os.Getenv("GCP_ZONE"), os.Getenv("GCP_SUBNET"))
}

func testAccTransitGatewayConfigGCPIPv6WithInsaneMode(rName string) string {
	return fmt.Sprintf(`
resource "aviatrix_account" "test_acc_gcp" {
	account_name                        = "tfa-gcp-%s"
	cloud_type                          = 4
	gcloud_project_id                   = "%s"
	gcloud_project_credentials_filepath = "%s"
}
resource "aviatrix_transit_gateway" "test_transit_gateway_ipv6_gcp_insane" {
	cloud_type        = 4
	account_name      = aviatrix_account.test_acc_gcp.account_name
	gw_name           = "tfg-gcp-ipv6-insane-%[1]s"
	vpc_id            = "%[4]s"
	vpc_reg           = "%[5]s"
	gw_size           = "n1-highmem-4"
	subnet            = "%[6]s"
	enable_ipv6       = true
	connected_transit = true
	insane_mode       = true
}
	`, rName, os.Getenv("GCP_PROJECT_ID"), os.Getenv("GOOGLE_CREDENTIALS_FILEPATH"),
		os.Getenv("GCP_VPC_ID"), os.Getenv("GCP_ZONE"), os.Getenv("GCP_SUBNET"))
}
