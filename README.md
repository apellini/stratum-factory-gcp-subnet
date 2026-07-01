# Factory Module: `stratum-factory-gcp-subnet`

Provisions a GCP regional subnetwork for the STRATUM platform.

## Purpose

Creates a `google_compute_subnetwork` within an existing VPC. The Wrapper passes the VPC `self_link` as the `network` input. The subnet's `self_link` is passed downstream to GCE instance modules.

## Usage

```hcl
module "dev_subnet" {
  source = "git::https://github.com/apellini/stratum-factory-gcp-subnet.git?ref=v0.1.0"

  environment              = "dev"
  project_id               = "stratum-dev-sandbox"
  region                   = "europe-west1"
  name_prefix              = "stratum-dev"
  network                  = module.dev_vpc.vpc_self_link
  subnet_cidr              = "10.100.0.0/24"
  private_ip_google_access = true
  tags                     = { environment = "dev", managed_by = "opentofu" }
}
```

## Inputs

| Name | Type | Required | Validation | Description |
|------|------|----------|------------|-------------|
| `environment` | `string` | yes | one of `dev`, `stage`, `main` | Deployment environment |
| `project_id` | `string` | yes | non-empty, no whitespace | GCP project ID |
| `region` | `string` | yes | valid GCP region format | GCP region for the subnet |
| `name_prefix` | `string` | yes | 3-24 chars, lowercase alphanumeric/hyphens, starts with letter | Resource name prefix |
| `network` | `string` | yes | valid GCP network self-link URI | VPC self-link (from vpc module output) |
| `subnet_cidr` | `string` | yes | valid IPv4 CIDR | IP range for the subnet |
| `private_ip_google_access` | `bool` | optional (default: `true`) | — | Enable Private Google Access |
| `tags` | `map(string)` | optional | all keys/values non-empty | Labels applied to all resources |

## Outputs

| Name | Type | Description |
|------|------|-------------|
| `subnet_id` | `string` | Unique identifier of the subnetwork |
| `subnet_name` | `string` | Name of the subnetwork |
| `subnet_self_link` | `string` | Self-link URI — pass to GCE instance modules |
| `subnet_cidr` | `string` | IP CIDR range of the subnetwork |

## Factory rules applied

- **Stateless** — no local state, no remote state reads
- **Input-only** — all configuration via `variables.tf`; nothing hardcoded
- **Strict validation** — every variable has a `validation` block
- **No secrets** — no sensitive values in code or outputs
- **Documented for humans and RAG** — this README + inline comments

## Release

```hcl
source = "git::https://github.com/apellini/stratum-factory-gcp-subnet.git?ref=v0.1.0"
```
