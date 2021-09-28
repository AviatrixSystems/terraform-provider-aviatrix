output "vnet" {
   value = azurerm_subnet.aviatrix.virtual_network_name
} 
output "group" {
   value = azurerm_subnet.aviatrix.resource_group_name
}
output "subnet" {
   value = azurerm_subnet.aviatrix.address_prefix
}
output "guid" {
   value = azurerm_virtual_network.aviatrix.guid
}