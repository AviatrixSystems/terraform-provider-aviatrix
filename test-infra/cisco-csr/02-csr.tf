data aws_ami csr_ami {
  most_recent = true
  owners      = [
    "aws-marketplace"]
  # Canonical
  filter {
    name   = "product-code"
    values = [
      "5tiyrfb5tasxk9gmnab39b843"]
    # aws ec2 describe-images --region us-east-2 --filters "Name=product-code,Values=5tiyrfb5tasxk9gmnab39b843"
  }
}

resource "tls_private_key" "key_pair_material" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_key_pair" "csr_key_pair" {
  key_name = "csr-kp"
  public_key = tls_private_key.key_pair_material.public_key_openssh
}

locals {
  key_file_path = "/tmp/csr-kp.pem"
}

resource "null_resource" "key_pair_file" {
  provisioner "local-exec" {
    command = "echo \"${tls_private_key.key_pair_material.private_key_pem}\" > ${local.key_file_path}"
  }
}

resource aws_instance csr_instance_1 {
  ami                     = data.aws_ami.csr_ami.id
  disable_api_termination = false
  instance_type           = "t2.medium"
  key_name                = aws_key_pair.csr_key_pair.key_name

  network_interface {
    network_interface_id = aws_network_interface.csr_aws_netw_interface_1.id
    device_index         = 0
  }

  tags = {
    Name    = "csr-instance-1"
    Purpose = "Terraform Acceptance"
  }
}
