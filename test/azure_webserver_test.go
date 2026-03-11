package test

import (
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

var subscriptionID string = "b7a56bf2-e1b1-422d-a2b2-b7652c051180"

func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: "../",
		Vars: map[string]interface{}{
			"labelPrefix": "mend0214",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	nicName := terraform.Output(t, terraformOptions, "nic_name")

	// Confirm VM exists
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))

	// Confirm NIC exists
	assert.True(t, azure.NetworkInterfaceExists(t, nicName, resourceGroupName, subscriptionID))

	// Confirm NIC is attached to the VM
	vm := azure.GetVirtualMachine(t, vmName, resourceGroupName, subscriptionID)

	assert.NotNil(t, vm.NetworkProfile)
	assert.NotNil(t, vm.NetworkProfile.NetworkInterfaces)
	assert.True(t, len(*vm.NetworkProfile.NetworkInterfaces) > 0)

	attachedNicID := *(*vm.NetworkProfile.NetworkInterfaces)[0].ID
	assert.True(t, strings.Contains(attachedNicID, nicName))

	// Confirm correct Ubuntu image
	assert.Equal(t, "Canonical", *vm.StorageProfile.ImageReference.Publisher)
	assert.Equal(t, "0001-com-ubuntu-server-jammy", *vm.StorageProfile.ImageReference.Offer)
	assert.Equal(t, "22_04-lts-gen2", *vm.StorageProfile.ImageReference.Sku)
}