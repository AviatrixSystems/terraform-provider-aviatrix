output "vpc" {
   value = aws_subnet.vpc-public.vpc_id
} 
output "subnet" {
   value = aws_subnet.vpc-public.cidr_block
}
output "subnet_id" {
   value = aws_subnet.vpc-public.id
}
