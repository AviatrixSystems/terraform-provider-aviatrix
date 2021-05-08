resource random_integer csr_vpc_int {
  count = 2
  min   = 1
  max   = 126
}

resource aws_vpc csr_vpc {
  cidr_block = join(".", [
    random_integer.csr_vpc_int[0].result,
    random_integer.csr_vpc_int[1].result,
    "0.0/16"])

  tags = {
    Name    = "csr-vpc"
    Purpose = "Terraform Acceptance"
  }
}

resource aws_subnet csr_vpc_subnet_1 {
  vpc_id                  = aws_vpc.csr_vpc.id
  cidr_block              = join(".", [
    random_integer.csr_vpc_int[0].result,
    random_integer.csr_vpc_int[1].result,
    "10.0/24"])
  availability_zone       = format("%s%s", var.aws_region, "b")
  map_public_ip_on_launch = true

  tags = {
    Name    = "csr-vpc-subnet-${random_integer.csr_vpc_int[0].result}"
    Purpose = "Terraform Acceptance"
  }
}

resource aws_network_interface csr_aws_netw_interface_1 {
  subnet_id       = aws_subnet.csr_vpc_subnet_1.id
  security_groups = [
    aws_security_group.csr_sec_group.id]

  tags = {
    Name    = "csr-aws-netw-interface-${random_integer.csr_vpc_int[0].result}"
    Purpose = "Terraform Acceptance"
  }
}

resource aws_eip csr_eip_1 {
  tags = {
    Name    = "csr-eip-${random_integer.csr_vpc_int[0].result}"
    Purpose = "Terraform Acceptance"
  }

  lifecycle {
    ignore_changes = [
      tags]
  }
}

resource aws_eip_association csr_eip_association_1 {
  allocation_id        = aws_eip.csr_eip_1.id
  network_interface_id = aws_network_interface.csr_aws_netw_interface_1.id
}
