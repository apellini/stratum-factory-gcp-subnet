# stratum-factory-gcp-subnet/outputs.tf
#
# FACTORY MODULE outputs — documented for humans and RAG.

output "subnet_id" {
  description = <<-EOT
    Unique identifier of the subnetwork.
    Type: string.
    Example: "projects/stratum-dev-sandbox/regions/europe-west1/subnetworks/stratum-dev-subnet"
  EOT
  value       = google_compute_subnetwork.subnet.id
}

output "subnet_name" {
  description = <<-EOT
    Name of the subnetwork.
    Type: string. Example: "stratum-dev-subnet"
  EOT
  value       = google_compute_subnetwork.subnet.name
}

output "subnet_self_link" {
  description = <<-EOT
    Self-link URI of the subnetwork.
    Type: string.
    Example: "https://www.googleapis.com/compute/v1/projects/stratum-dev-sandbox/regions/europe-west1/subnetworks/stratum-dev-subnet"
    Pass to GCE instance Factory modules via the Wrapper.
  EOT
  value       = google_compute_subnetwork.subnet.self_link
}

output "subnet_cidr" {
  description = <<-EOT
    IP CIDR range of the subnetwork.
    Type: string. Example: "10.100.0.0/24"
  EOT
  value       = google_compute_subnetwork.subnet.ip_cidr_range
}

output "secondary_ip_range_names" {
  description = <<-EOT
    List of secondary IP range names configured on the subnetwork.
    Type: list(string).
    Example: ["gke-pods", "gke-services"]
    Empty list when secondary_ip_ranges = [].
  EOT
  value = [for r in var.secondary_ip_ranges : r.range_name]
}
