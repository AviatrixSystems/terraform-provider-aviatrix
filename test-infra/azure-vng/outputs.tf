output "resource_group" {
  value = azurerm_resource_group.test.name
}

output "vnet" {
  value = azurerm_virtual_network.test.name
}

output "subnet" {
  value = azurerm_subnet.test1.address_prefixes[0]
}

output "vng" {
  value = azurerm_virtual_network_gateway.test.name
}
