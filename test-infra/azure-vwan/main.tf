resource "azurerm_resource_group" "test" {
  name     = "cwan-rg"
  location = var.azure_region
}

resource "azurerm_virtual_wan" "test" {
  name                = "cwan-vwan"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_hub" "test" {
  name                = "cwan-vhub"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.1.0/24"
}

resource "azurerm_vpn_gateway" "test" {
  name                = "cwan-vpngateway"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_hub_id      = azurerm_virtual_hub.test.id
}
