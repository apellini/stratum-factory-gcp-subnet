# Factory Module: `stratum-factory-gcp-subnet`

Provisions a GCP regional subnetwork for the STRATUM platform, with optional secondary IP ranges.

## Purpose

Creates a `google_compute_subnetwork` within an existing VPC. Secondary IP ranges (e.g. for GKE
pods and services) can be added via the `secondary_ip_ranges` input. When `secondary_ip_ranges`
is empty (the default), no secondary ranges are configured.

## Usage

```hcl
module "dev_subnet" {
  source = "git::https://github.com/apellini/stratum-factory-gcp-subnet.git?ref=v0.2.0"

  environment              = "dev"
  project_id               = "stratum-dev-sandbox"
  region                   = "europe-west1"
  name_prefix              = "stratum-dev"
  network                  = module.dev_vpc.vpc_self_link
  subnet_cidr              = "10.100.0.0/24"
  private_ip_google_access = true
  tags                     = { environment = "dev", managed_by = "opentofu" }

  # Optional: secondary IP ranges (e.g. for GKE)
  secondary_ip_ranges = [
    { range_name = "gke-pods",     ip_cidr_range = "10.101.0.0/16" },
    { range_name = "gke-services", ip_cidr_range = "10.102.0.0/20" },
  ]
}
```

## Inputs

| Name | Type | Required | Validation | Description |
|------|------|----------|------------|-------------|
| `environment` | `string` | yes | one of `dev`, `stage`, `main` | Deployment environment |
| `project_id` | `string` | yes | non-empty, no whitespace | GCP project ID |
| `region` | `string` | yes | valid GCP region format | GCP region for the subnet |
| `name_prefix` | `string` | yes | 3–24 chars, lowercase alphanumeric/hyphens, starts with letter | Resource name prefix |
| `network` | `string` | yes | valid GCP network self-link URI | VPC self-link (from vpc module output) |
| `subnet_cidr` | `string` | yes | valid IPv4 CIDR | Primary IP range for the subnet |
| `private_ip_google_access` | `bool` | optional (default `true`) | — | Enable Private Google Access |
| `secondary_ip_ranges` | `list(object)` | optional (default `[]`) | see below | Secondary IP ranges |
| `tags` | `map(string)` | optional (default `{}`) | all keys/values non-empty | Labels applied to all resources |

### `secondary_ip_ranges` object attributes

| Attribute | Type | Validation | Description |
|-----------|------|------------|-------------|
| `range_name` | `string` | lowercase letters/digits/hyphens, starts and ends with letter or digit | Name for the secondary range |
| `ip_cidr_range` | `string` | valid IPv4 CIDR | CIDR block for this secondary range |

## Outputs

| Name | Type | Description |
|------|------|-------------|
| `subnet_id` | `string` | Unique identifier of the subnetwork |
| `subnet_name` | `string` | Name of the subnetwork |
| `subnet_self_link` | `string` | Self-link URI — pass to GCE instance modules |
| `subnet_cidr` | `string` | Primary IP CIDR range of the subnetwork |
| `secondary_ip_range_names` | `list(string)` | Names of secondary IP ranges configured |

## Factory rules applied

- **Stateless** — no local state, no remote state reads
- **Input-only** — all configuration via `variables.tf`; nothing hardcoded
- **Strict validation** — every variable has a `validation` block
- **No secrets** — no sensitive values in code or outputs
- **Documented for humans and RAG** — this README + inline comments

## Release

```hcl
source = "git::https://github.com/apellini/stratum-factory-gcp-subnet.git?ref=v0.2.0"
```
