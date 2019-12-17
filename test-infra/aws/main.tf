resource "aws_vpc" "vpc" {
  cidr_block         = var.aws_vpc_cidr
  enable_classiclink = "false"
  tags = {
    Name = "aviatrix-vpc"
  }
}
resource "aws_subnet" "vpc-public" {
  vpc_id                  = aws_vpc.vpc.id
  cidr_block              = var.aws_vpc_subnet
  map_public_ip_on_launch = "true"
  tags = {
    Name = "aviatrix-public"
  }
}
resource "aws_internet_gateway" "vpc-gw" {
  vpc_id = aws_vpc.vpc.id
  tags = {
    Name = "aviatrix-igw"
  }
}
resource "aws_route_table" "vpc-route" {
  vpc_id = aws_vpc.vpc.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.vpc-gw.id
  }
  tags = {
    Name = "aviatrix-route"
  }
  lifecycle {
    ignore_changes = [route]
  }
}
resource "aws_route_table_association" "vpc-ra" {
  subnet_id          = aws_subnet.vpc-public.id
  route_table_id     = aws_route_table.vpc-route.id
  depends_on         = [
    aws_subnet.vpc-public,
    aws_route_table.vpc-route,
    aws_internet_gateway.vpc-gw,
    aws_vpc.vpc,
  ]
}


