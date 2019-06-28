resource "google_compute_network" "vpc" {
  auto_create_subnetworks = false
  name                 = "aviatrix-network"
  project              = "${var.gcp_project_id}"
  provider             = "google"
}
resource "google_compute_subnetwork" "subnet" {
  name                 = "aviatrix-subnet"
  ip_cidr_range        = "${var.gcp_vpc_cidr}"
  region               = "${var.gcp_region}"
  network              = "${google_compute_network.vpc.self_link}"
  provider             = "google"
}
