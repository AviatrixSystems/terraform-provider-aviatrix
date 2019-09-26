resource "oci_core_vcn" "vpc" {
  display_name   = "oracle-network"
  cidr_block     = var.oci_vpc_cidr1
  compartment_id = var.oci_compartment_id
}
resource "oci_core_subnet" "subnet" {
  cidr_block     = var.oci_vpc_cidr1
  compartment_id = var.oci_compartment_id
  vcn_id         = oci_core_vcn.vpc.id
}
