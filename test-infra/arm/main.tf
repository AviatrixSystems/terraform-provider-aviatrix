resource "azurerm_resource_group" "aviatrix" {
  name                 = "${var.name}-group"
  location             = "${var.azure_region}"
}
resource "azurerm_virtual_network" "aviatrix" {
  name                 = "${var.name}-vnet"
  resource_group_name  = "${azurerm_resource_group.aviatrix.name}"
  location             = "${azurerm_resource_group.aviatrix.location}"
  address_space        = ["${var.azure_vpc_cidr}"]
}
resource "azurerm_subnet" "aviatrix" {
  name                 = "aviatrix-subnet"
  resource_group_name  = "${azurerm_resource_group.aviatrix.name}"
  virtual_network_name = "${azurerm_virtual_network.aviatrix.name}"
  address_prefix       = "${var.azure_vpc_subnet}"
}
