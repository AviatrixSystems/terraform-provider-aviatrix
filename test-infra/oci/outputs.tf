output "vpc_id" {
   value = oci_core_vcn.vpc.display_name
} 
output "subnet" {
   value = oci_core_subnet.subnet.cidr_block
} 
