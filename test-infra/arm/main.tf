resource "azurerm_resource_group" "aviatrix" {
  name     = "${var.name}-group"
  location = var.azure_region
}
resource "azurerm_virtual_network" "aviatrix" {
  name                 = "${var.name}-vnet"
  resource_group_name  = azurerm_resource_group.aviatrix.name
  location             = azurerm_resource_group.aviatrix.location
  address_space        = ["${var.azure_vpc_cidr}"]
}

resource "azurerm_subnet" "aviatrix" {
  name                 = "aviatrix-subnet"
  resource_group_name  = azurerm_resource_group.aviatrix.name
  virtual_network_name = azurerm_virtual_network.aviatrix.name
  address_prefix       = var.azure_vpc_subnet
}

resource "azurerm_route_table" "aviatrix" {
  name                          = "acceptanceTestSecurityGroup1"
  location                      = azurerm_resource_group.aviatrix.location
  resource_group_name           = azurerm_resource_group.aviatrix.name
}

resource "azurerm_subnet" "aviatrix-private" {
  name                 = "aviatrix-private-subnet"
  resource_group_name  = azurerm_resource_group.aviatrix.name
  virtual_network_name = azurerm_virtual_network.aviatrix.name
  address_prefix       = var.azure_address_prefix
}

resource "azurerm_subnet_route_table_association" "aviatrix" {
  subnet_id      = azurerm_subnet.aviatrix-private.id
  route_table_id = azurerm_route_table.aviatrix.id
}