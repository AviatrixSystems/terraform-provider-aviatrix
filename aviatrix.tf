provider "aviatrix" {
  controller_ip = "1.2.3.4"
  username = "admin"
  password = "password"
}

resource "aviatrix_account" "test_account" {
  account_name = "myacc"
  account_password = "P@55w0rd"
  account_email = "support@aviatrix.com"
  cloud_type = 1
  aws_account_number = "123456789"
  aws_access_key = "ABCDEFGHIJKL"
  aws_secret_key = "ABCDEFGHIJKLabcdefghijkl"
}
