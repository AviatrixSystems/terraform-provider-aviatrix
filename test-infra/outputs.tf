#output "aws_vpc_id" {
#  value = "${aws_vpc.vpc.*.id}"
#}
#output "aws_vpc-subnet" {
#  value = "${aws_subnet.vpc-public-ha.*.cidr_block}"
#}
#output "arm_vnet" {
#  value = "${azurerm_virtual_network.hashicorp.*.name}"
#}
#output "arm_resource_group" {
#  value = "${azurerm_virtual_network.hashicorp.*.resource_group_name}"
#}
#output "arm_vpc-subnet" {
#  value = "${azurerm_subnet.hashicorp.*.address_prefix}"
#}
#output "gcp_vpc_id" {
#  value = "${google_compute_network.vpc.*.id}"
#}
#output "gcp_vpc-subnet" {
#  value = "${google_compute_subnetwork.subnet.*.ip_cidr_range}"
#}
