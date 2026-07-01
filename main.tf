# stratum-factory-gcp-subnet/main.tf
#
# FACTORY MODULE — GCP regional subnetwork.
#
# STATELESS & INPUT-ONLY:
#   - All configuration via variables.tf.
#   - No hardcoded values, no remote state reads, no secrets.
#   - State owned by the calling Wrapper.
#
# Consumed by the Wrapper as:
#   source = "git::https://github.com/apellini/stratum-factory-gcp-subnet.git?ref=v<semver>"

# ── Subnetwork ────────────────────────────────────────────────────────────────
resource "google_compute_subnetwork" "subnet" {
  project                  = var.project_id
  name                     = "${var.name_prefix}-subnet"
  region                   = var.region
  network                  = var.network
  ip_cidr_range            = var.subnet_cidr
  private_ip_google_access = var.private_ip_google_access
  labels                   = var.tags

  description = "STRATUM ${var.environment} subnet — managed by OpenTofu (stratum-factory-gcp-subnet)"

  dynamic "secondary_ip_range" {
    for_each = var.secondary_ip_ranges
    content {
      range_name    = secondary_ip_range.value.range_name
      ip_cidr_range = secondary_ip_range.value.ip_cidr_range
    }
  }
}
