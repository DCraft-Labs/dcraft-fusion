package tenancy

import "testing"

func TestDeploymentClassification(t *testing.T) {
	if !IsCustomerOwned(DeploymentCustomerVPC) || !IsCustomerOwned(DeploymentOnPrem) {
		t.Fatal("customer-owned deployment profiles were not classified correctly")
	}
	if IsCustomerOwned(DeploymentPooledSaaS) {
		t.Fatal("pooled SaaS should not be customer-owned")
	}
	if !IsDedicated(DeploymentDedicatedDatabase) || !IsDedicated(DeploymentDedicatedStack) || !IsDedicated(DeploymentCustomerVPC) {
		t.Fatal("dedicated deployment profiles were not classified correctly")
	}
	if IsDedicated(DeploymentPooledSaaS) {
		t.Fatal("pooled SaaS should not be dedicated")
	}
}
