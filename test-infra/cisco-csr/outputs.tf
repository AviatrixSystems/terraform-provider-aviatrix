output "DEVICE_PUBLIC_IP" {
  value = aws_eip.csr_eip_1.public_ip
}

output "DEVICE_KEY_FILE_PATH" {
  value = local.key_file_path
}
