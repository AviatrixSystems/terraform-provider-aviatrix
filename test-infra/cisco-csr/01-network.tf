resource aws_internet_gateway csr_vpc_igw {
  vpc_id = aws_vpc.csr_vpc.id
  tags   = {
    Name    = "csr-igw-${random_integer.csr_vpc_int[0].result}"
    Purpose = "Terraform Regression"
  }
}

resource aws_route_table csr_vpc_rtb {
  vpc_id = aws_vpc.csr_vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.csr_vpc_igw.id
  }
  tags = {
    Name    = "csr-vpc-rtb-${random_integer.csr_vpc_int[0].result}"
    Purpose = "Terraform Regression"
  }
}

resource aws_route_table_association csr_vpc_subnet_rtb_1 {
  subnet_id      = aws_subnet.csr_vpc_subnet_1.id
  route_table_id = aws_route_table.csr_vpc_rtb.id
}

resource aws_security_group csr_sec_group {
  name        = "csr-sec-group-${random_integer.csr_vpc_int[0].result}"
  description = "Aviatrix - Controller Security Group"
  vpc_id      = aws_vpc.csr_vpc.id

  tags = {
    Name    = "csr-aws-sec-group-${random_integer.csr_vpc_int[0].result}"
    Purpose = "Terraform Acceptance"
  }
}

resource aws_security_group_rule egress_rule {
  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = [
    "0.0.0.0/0"]
  security_group_id = aws_security_group.csr_sec_group.id
}

resource aws_security_group_rule ingress_rule {
  type              = "ingress"
  from_port         = 22
  to_port           = 22
  protocol          = "tcp"
  cidr_blocks       = [
    "0.0.0.0/0"]
  security_group_id = aws_security_group.csr_sec_group.id
}
