package persistence

import (
	"errors"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

func TestInMemoryPersistence(t *testing.T) {
	persistence := NewInMemoryPersistence()

	deviceID := "test-device"
	device := &domain.SignatureDevice{
		Id:        deviceID,
		Algorithm: "RSA",
		Label:     "Test Device",
	}

	// Save a device
	err := persistence.SaveSignatureDevice(device)
	if err != nil {
		t.Errorf("Error saving device: %v", err)
	}

	// Get the saved device
	savedDevice, err := persistence.GetSignatureDevice(deviceID)
	if err != nil {
		t.Errorf("Error getting device: %v", err)
	}

	if savedDevice.Id != device.Id || savedDevice.Algorithm != device.Algorithm || savedDevice.Label != device.Label {
		t.Errorf("Saved device does not match expected values")
	}

	// List devices
	devices, err := persistence.ListSignatureDevices()
	if err != nil {
		t.Errorf("Error listing devices: %v", err)
	}

	if len(devices) != 1 || devices[0].Id != device.Id {
		t.Errorf("Listed devices do not match expected values")
	}

	// Get a non-existing device
	nonExistingDeviceID := "non-existing-device"
	_, err = persistence.GetSignatureDevice(nonExistingDeviceID)
	if err == nil || !errors.Is(err, domain.ErrDeviceNotFound) {
		t.Errorf("Expected device not found error")
	}
}
