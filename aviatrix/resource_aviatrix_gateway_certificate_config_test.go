package aviatrix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAviatrixGatewayCertificateConfig_basic(t *testing.T) {
	if os.Getenv("SKIP_GATEWAY_CERTIFICATE_CONFIG") == "yes" {
		t.Skip("Skipping Branch Router test as SKIP_GATEWAY_CERTIFICATE_CONFIG is set")
	}

	resourceName := "aviatrix_gateway_certificate_config.test_gateway_certificate_config"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGatewayCertificateConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayCertificateConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayCertificateConfigExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ca_certificate", "ca_private_key"},
			},
		},
	})
}

func testAccGatewayCertificateConfigBasic() string {
	return `
resource "aviatrix_gateway_certificate_config" "test_gateway_certificate_config" {
	ca_private_key = <<EOT
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAlnBS/Ei6vci53ReH8dLlW6aALupuMEcekN7vf7nWGygyr3i3
WG83syQ270yhDuI6ic9F4qrb6qP9lWhfBlMKrZQAUfhG+0a+GfMy4k+mY20Fvfcz
gY+EwNUxlGbSW/WQoDO8Mo07pBkMcolH+mhg0WvQKMw00iMAyAm0zIjHJluyVDh4
1TFO1KjZIdihhH5UrRyIruBaU4nOZEFmFBgtH1LwLZtRLKVdGkW9vdqk8M41/xCC
6+hBhQv9pzD2ziZ9sQ/ZOu5Aemd92LKPBtW9TM3Wqw0zoqTWdeIKSCgNyVO0a5BN
t4W3GG7D/IO2bM9x7kgi06Sj1yRXRmtdqUaxhQIDAQABAoIBABy5K48Rz93mkl49
XO52Juad3sGWbx12psZgWngXomKjBTJtqQwQiCEDHB4KkoHF/glL8vr5Rm4Bi6xY
NjR97H8B0CHWdq8JbciEn8WIEavQKBWCOmpVXbJ3wjSkgSufslJ0Lk20m5uUUMUZ
ow2TmlDB9gekHb21gzOubr/SqbuP4YNcLjr1EPKV/FwsOrXtWRfr5/oTfLEyeg+V
hD1JEb0Pg2JH8T8kWICouZmzJpZ26CtP6ssUujUsR9Ll2zfihXNVu2nPghfVWQVS
zEz6DBuLBBGgre9iGb19OgLY8q2i1h5pS1s3bj+NqeHAe4tJIKIWsd10E919pBIk
aSwzHQ0CgYEAyCmsvPfTOi++scMf5rBlzK+zJxFgzkH/YnZeNlle39CsEFGq+Q4k
dLwPgY5lDOvsfz4ccn6pUsVVuhjiIvIRDoTQFvDl25rCH47kPbg19z7lNdKWd7jT
3na2qluUOLW+LshDX9fEIZ6AD9AjjF8W++CHachCfKiWnynl68vght8CgYEAwGeu
ybugACbUM+MzQmz7j61GW6l/fZHat4Clgcsa3oAIuCxldajqRwKNV9EBqzeb9QaN
W/XbyKLlZp5VabRmLZCN79b0PoGrh4T0r3Z4eIDyEVW82uL2PfQaGOkPNUgsRFQd
5juElkEJ3TR9hDhNR5EsRnj7YeMCdzNIL1dOiBsCgYBHpGr4Y/eNjwNBGubzKdX1
8jk8VYMBsCuZcWZ9K3XCxCyh1qlMZVx1D864/195xYOrc265KE6wmoL5jeh6u4uR
V8YnP+f1tymeJAXbdXCTY0alAg/rIBNtP65XwVmHBr3gfrtmtZK0ucd6YXQnzO0s
EAxHt27csXbf1x49TIa7nwKBgA1PnwIJ3vnjtc6ZK1SvIYBBJpr00QMo8g678bdX
C9bU2MWVHLfVJdAf1xN6PDlSMZH0EBeKnNf2nRRKY/JaLq98TrWHE4K3idxoUF1L
Pu5nTVfxrI0gIpUKrDuI9CplgdqAT0k8WOHkQgBxvzVEh+QpoPyHJi0RfXwtqhLM
YXx9AoGBAJmxo7/4LKKjLOWv0qNm2qetVOQqeMy5e/J7T+60DP3Ctqtv+pwLj1aS
hAHjaMPxNRUr6M7OFsmVZ8PO/7Ud/Llf89rmdqOmAqrJls+XNf/TREmloYh88pXO
3BmygRRtb3yPCFw++hXtEzWG7z7oDyyyiT6JXZINhOCfvh4Y9icw
-----END RSA PRIVATE KEY-----

EOT

	ca_certificate = <<EOT
-----BEGIN CERTIFICATE-----
MIIDKjCCAhICCQDgRt4UZsINwjANBgkqhkiG9w0BAQsFADBXMQswCQYDVQQGEwJV
UzELMAkGA1UECAwCQ0ExEjAQBgNVBAcMCUN1cGVydGlubzEOMAwGA1UECgwFQXBw
bGUxFzAVBgNVBAMMDmN5cnVzamF2YW4uY29tMB4XDTIxMDIwMjE4NDU1M1oXDTI2
MDIwMTE4NDU1M1owVzELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMRIwEAYDVQQH
DAlDdXBlcnRpbm8xDjAMBgNVBAoMBUFwcGxlMRcwFQYDVQQDDA5jeXJ1c2phdmFu
LmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAJZwUvxIur3Iud0X
h/HS5VumgC7qbjBHHpDe73+51hsoMq94t1hvN7MkNu9MoQ7iOonPReKq2+qj/ZVo
XwZTCq2UAFH4RvtGvhnzMuJPpmNtBb33M4GPhMDVMZRm0lv1kKAzvDKNO6QZDHKJ
R/poYNFr0CjMNNIjAMgJtMyIxyZbslQ4eNUxTtSo2SHYoYR+VK0ciK7gWlOJzmRB
ZhQYLR9S8C2bUSylXRpFvb3apPDONf8QguvoQYUL/acw9s4mfbEP2TruQHpnfdiy
jwbVvUzN1qsNM6Kk1nXiCkgoDclTtGuQTbeFtxhuw/yDtmzPce5IItOko9ckV0Zr
XalGsYUCAwEAATANBgkqhkiG9w0BAQsFAAOCAQEARTwMmBLw0R+D1HRPIipf+mx1
udKDadcjypkQpqAnTkzXngnc/+tZi7vB0EhiU3ODAUvc3dv8BGBx7XFhR03jh31+
xWQZSAY8zVzfzwkYPlgYL7+L8TW+WT+rfkoaFF+xFCzSpCD6dKPTpzMNCHwshqua
Tz0kEeJ6d2ZuXICGNyl0gMxnULapJjW4sDbMNeK9bl3cJPF9BsfT1nDIlNGiJ6vr
KQZ0NERnAlm69cJiIvOx1xYKW9pw+sQHDJPouIAjFoH+eDvhZLqIrE7aZKUNFrrK
JSbEHyZUdM+2F2xWRk1oBHmatyKdzi6vtOC0Yix+oa2iDUJfNGdsI3xKaF+oKA==
-----END CERTIFICATE-----

EOT

}
`
}

func testAccCheckGatewayCertificateConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("gateway_certificate_config Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no gateway_certificate_config ID is set")
		}

		client := testAccProvider.Meta().(*goaviatrix.Client)

		status, err := client.GetGatewayCertificateStatus(context.Background())
		if err != nil {
			return err
		}
		if status != "enabled" {
			return fmt.Errorf("gateway_certificate_config not found")
		}

		return nil
	}
}

func testAccCheckGatewayCertificateConfigDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*goaviatrix.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aviatrix_gateway_certificate_config" {
			continue
		}

		status, _ := client.GetGatewayCertificateStatus(context.Background())
		if status != "disabled" {
			return fmt.Errorf("gateway_certificate_config still exists")
		}
	}

	return nil
}
