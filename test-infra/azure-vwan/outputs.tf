output "azure_resource_group" {
   value = azurerm_resource_group.test.name
}

output "azure_hub_name" {
   value = azurerm_virtual_hub.test.name
}

output "azure_vpn_gateway_asn" {
   value = azurerm_vpn_gateway.test.bgp_settings.asn
}
