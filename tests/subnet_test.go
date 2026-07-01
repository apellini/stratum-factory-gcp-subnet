// tests/subnet_test.go — Terratest unit tests for stratum-factory-gcp-subnet.
//
// Run locally:
//
//	cd tests && go mod tidy && go test -v -timeout 10m ./...
//
// Tests are credential-free:
//   - Positive tests use `tofu validate` (static analysis, no API calls).
//   - Negative tests use InitAndPlanE with invalid var values (validation fires before provider auth).
package test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const validNetwork = "https://www.googleapis.com/compute/v1/projects/stratum-dev-sandbox/global/networks/stratum-dev-vpc"

func moduleDir(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs("..")
	require.NoError(t, err)
	return root
}

func tofuOptions(t *testing.T, vars map[string]interface{}) *terraform.Options {
	t.Helper()
	return &terraform.Options{
		TerraformDir:    moduleDir(t),
		TerraformBinary: "tofu",
		Vars:            vars,
		NoColor:         true,
	}
}

func printReport(t *testing.T, rows [][]string) {
	t.Helper()
	header := fmt.Sprintf("%-60s %-10s %s", "Test", "Result", "Detail")
	sep := strings.Repeat("─", 100)
	t.Log("\n" + sep)
	t.Log("  STRATUM-FACTORY — GCP Subnet Module Test Report")
	t.Log(sep)
	t.Log(header)
	t.Log(sep)
	for _, row := range rows {
		t.Logf("  %-60s %-10s %s", row[0], row[1], row[2])
	}
	t.Log(sep)
	if summaryFile := os.Getenv("GITHUB_STEP_SUMMARY"); summaryFile != "" {
		f, err := os.OpenFile(summaryFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return
		}
		defer f.Close()
		fmt.Fprintln(f, "## GCP Subnet Module Test Report")
		fmt.Fprintln(f, "| Test | Result | Detail |")
		fmt.Fprintln(f, "|------|--------|--------|")
		for _, row := range rows {
			fmt.Fprintf(f, "| %s | %s | %s |\n", row[0], row[1], row[2])
		}
	}
}

// TestSubnetValidate verifies that tofu validate succeeds for valid inputs.
func TestSubnetValidate(t *testing.T) {
	t.Parallel()
	opts := tofuOptions(t, map[string]interface{}{
		"environment": "dev",
		"project_id":  "stratum-dev-sandbox",
		"region":      "europe-west1",
		"name_prefix": "stratum-dev",
		"network":     validNetwork,
		"subnet_cidr": "10.100.0.0/24",
	})
	terraform.Init(t, opts)
	_, err := terraform.RunTerraformCommandE(t, opts, "validate")
	result, detail := "✅ PASS", "validate completed"
	if err != nil {
		result, detail = "❌ FAIL", err.Error()
	}
	printReport(t, [][]string{{"SubnetValidate", result, detail}})
	require.NoError(t, err, "tofu validate must pass for valid inputs")
}

// TestSubnetValidateWithSecondaryRange verifies validate succeeds with a secondary IP range.
func TestSubnetValidateWithSecondaryRange(t *testing.T) {
	t.Parallel()
	opts := tofuOptions(t, map[string]interface{}{
		"environment": "dev",
		"project_id":  "stratum-dev-sandbox",
		"region":      "europe-west1",
		"name_prefix": "stratum-dev",
		"network":     validNetwork,
		"subnet_cidr": "10.100.0.0/24",
		"secondary_ip_ranges": []interface{}{
			map[string]interface{}{
				"range_name":    "gke-pods",
				"ip_cidr_range": "10.101.0.0/16",
			},
		},
	})
	terraform.Init(t, opts)
	_, err := terraform.RunTerraformCommandE(t, opts, "validate")
	result, detail := "✅ PASS", "validate completed with secondary IP range"
	if err != nil {
		result, detail = "❌ FAIL", err.Error()
	}
	printReport(t, [][]string{{"SubnetValidateWithSecondaryRange", result, detail}})
	require.NoError(t, err, "tofu validate must pass with a valid secondary IP range")
}

// TestSubnetRejectsInvalidEnvironment verifies that an invalid environment is rejected.
func TestSubnetRejectsInvalidEnvironment(t *testing.T) {
	t.Parallel()
	opts := tofuOptions(t, map[string]interface{}{
		"environment": "production",
		"project_id":  "stratum-dev-sandbox",
		"region":      "europe-west1",
		"name_prefix": "stratum-dev",
		"network":     validNetwork,
		"subnet_cidr": "10.100.0.0/24",
	})
	_, err := terraform.InitAndPlanE(t, opts)
	result, detail := "✅ PASS", "plan correctly rejected invalid environment"
	if err == nil {
		result, detail = "❌ FAIL", "plan should have failed"
	}
	printReport(t, [][]string{{"SubnetRejectsInvalidEnvironment", result, detail}})
	assert.Error(t, err, "must reject environment='production'")
}

// TestSubnetRejectsInvalidRegion verifies that a malformed region is rejected.
func TestSubnetRejectsInvalidRegion(t *testing.T) {
	t.Parallel()
	opts := tofuOptions(t, map[string]interface{}{
		"environment": "dev",
		"project_id":  "stratum-dev-sandbox",
		"region":      "us east 1",
		"name_prefix": "stratum-dev",
		"network":     validNetwork,
		"subnet_cidr": "10.100.0.0/24",
	})
	_, err := terraform.InitAndPlanE(t, opts)
	result, detail := "✅ PASS", "plan correctly rejected invalid region"
	if err == nil {
		result, detail = "❌ FAIL", "plan should have failed"
	}
	printReport(t, [][]string{{"SubnetRejectsInvalidRegion", result, detail}})
	assert.Error(t, err, "must reject region with spaces")
}

// TestSubnetRejectsInvalidSubnetCidr verifies that a malformed CIDR is rejected.
func TestSubnetRejectsInvalidSubnetCidr(t *testing.T) {
	t.Parallel()
	opts := tofuOptions(t, map[string]interface{}{
		"environment": "dev",
		"project_id":  "stratum-dev-sandbox",
		"region":      "europe-west1",
		"name_prefix": "stratum-dev",
		"network":     validNetwork,
		"subnet_cidr": "not-a-cidr",
	})
	_, err := terraform.InitAndPlanE(t, opts)
	result, detail := "✅ PASS", "plan correctly rejected invalid CIDR"
	if err == nil {
		result, detail = "❌ FAIL", "plan should have failed"
	}
	printReport(t, [][]string{{"SubnetRejectsInvalidSubnetCidr", result, detail}})
	assert.Error(t, err, "must reject subnet_cidr='not-a-cidr'")
}

// TestSubnetRejectsInvalidNamePrefix verifies that a name_prefix starting with a digit is rejected.
func TestSubnetRejectsInvalidNamePrefix(t *testing.T) {
	t.Parallel()
	opts := tofuOptions(t, map[string]interface{}{
		"environment": "dev",
		"project_id":  "stratum-dev-sandbox",
		"region":      "europe-west1",
		"name_prefix": "1bad",
		"network":     validNetwork,
		"subnet_cidr": "10.100.0.0/24",
	})
	_, err := terraform.InitAndPlanE(t, opts)
	result, detail := "✅ PASS", "plan correctly rejected name_prefix starting with digit"
	if err == nil {
		result, detail = "❌ FAIL", "plan should have failed"
	}
	printReport(t, [][]string{{"SubnetRejectsInvalidNamePrefix", result, detail}})
	assert.Error(t, err, "must reject name_prefix='1bad'")
}

// TestSubnetRejectsInvalidSecondaryRangeName verifies that an invalid range_name is rejected.
func TestSubnetRejectsInvalidSecondaryRangeName(t *testing.T) {
	t.Parallel()
	opts := tofuOptions(t, map[string]interface{}{
		"environment": "dev",
		"project_id":  "stratum-dev-sandbox",
		"region":      "europe-west1",
		"name_prefix": "stratum-dev",
		"network":     validNetwork,
		"subnet_cidr": "10.100.0.0/24",
		"secondary_ip_ranges": []interface{}{
			map[string]interface{}{
				"range_name":    "1bad-name",
				"ip_cidr_range": "10.101.0.0/16",
			},
		},
	})
	_, err := terraform.InitAndPlanE(t, opts)
	result, detail := "✅ PASS", "plan correctly rejected invalid secondary range name"
	if err == nil {
		result, detail = "❌ FAIL", "plan should have failed"
	}
	printReport(t, [][]string{{"SubnetRejectsInvalidSecondaryRangeName", result, detail}})
	assert.Error(t, err, "must reject range_name starting with digit")
}

// TestSubnetRejectsInvalidSecondaryCidr verifies that an invalid secondary ip_cidr_range is rejected.
func TestSubnetRejectsInvalidSecondaryCidr(t *testing.T) {
	t.Parallel()
	opts := tofuOptions(t, map[string]interface{}{
		"environment": "dev",
		"project_id":  "stratum-dev-sandbox",
		"region":      "europe-west1",
		"name_prefix": "stratum-dev",
		"network":     validNetwork,
		"subnet_cidr": "10.100.0.0/24",
		"secondary_ip_ranges": []interface{}{
			map[string]interface{}{
				"range_name":    "gke-pods",
				"ip_cidr_range": "not-a-cidr",
			},
		},
	})
	_, err := terraform.InitAndPlanE(t, opts)
	result, detail := "✅ PASS", "plan correctly rejected invalid secondary CIDR"
	if err == nil {
		result, detail = "❌ FAIL", "plan should have failed"
	}
	printReport(t, [][]string{{"SubnetRejectsInvalidSecondaryCidr", result, detail}})
	assert.Error(t, err, "must reject invalid ip_cidr_range for secondary range")
}

// TestNoTerraformBinary verifies that no .tf file references the `terraform` binary.
func TestNoTerraformBinary(t *testing.T) {
	t.Parallel()
	tfFiles, _ := filepath.Glob("../*.tf")
	issues := []string{}
	for _, f := range tfFiles {
		content, err := os.ReadFile(f)
		require.NoError(t, err)
		lines := strings.Split(string(content), "\n")
		for i, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.Contains(trimmed, "terraform") &&
				!strings.HasPrefix(trimmed, "#") &&
				!strings.HasPrefix(trimmed, "//") &&
				trimmed != "terraform {" &&
				!strings.HasPrefix(trimmed, "required_version") &&
				!strings.HasPrefix(trimmed, "required_providers") &&
				!strings.HasPrefix(trimmed, "backend") &&
				!strings.Contains(trimmed, "TerraformBinary") {
				if strings.Contains(trimmed, "\"terraform\"") || strings.Contains(trimmed, "`terraform`") {
					issues = append(issues, fmt.Sprintf("%s:%d: %s", f, i+1, trimmed))
				}
			}
		}
	}
	result, detail := "✅ PASS", "no terraform binary references found"
	if len(issues) > 0 {
		result, detail = "❌ FAIL", fmt.Sprintf("found: %v", issues)
	}
	printReport(t, [][]string{{"NoTerraformBinary", result, detail}})
	assert.Empty(t, issues, "no .tf file should reference the 'terraform' binary")
}
