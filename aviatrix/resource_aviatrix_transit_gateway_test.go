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
						resource.TestCheckResourceAttr(resourceNameAEP, "enable_jumbo_frame", "true"),
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
	enable_jumbo_frame = true
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
	expectedWANCount := 4
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
				{LogicalIfName: "wan0", IpAddress: "10.0.0.2"},
				{LogicalIfName: "wan1", IpAddress: "10.0.0.3"},
				{LogicalIfName: "wan2", GatewayIp: "192.168.1.1"},
				{LogicalIfName: "mgmt0", Dhcp: true},
			},
			expected: []map[string]interface{}{
				{"logical_ifname": "wan0", "ip_address": "10.0.0.2"},
				{"logical_ifname": "wan1", "ip_address": "10.0.0.3"},
				{"logical_ifname": "wan2", "gateway_ip": "192.168.1.1"},
				{"logical_ifname": "mgmt0", "dhcp": true},
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := setInterfaceDetails(tt.interfaces)
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
