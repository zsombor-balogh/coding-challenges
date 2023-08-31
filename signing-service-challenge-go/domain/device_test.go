package domain

import (
	"errors"
	"testing"
)

func TestNewSignatureDeviceRSA(t *testing.T) {
	deviceId := "test-device-rsa"
	algorithm := "RSA"
	label := "Test Device RSA"

	device, err := NewSignatureDevice(deviceId, algorithm, label)
	if err != nil {
		t.Errorf("Error creating signature device: %v", err)
	}

	if device.Id != deviceId {
		t.Errorf("Device Id doesn't match")
	}
	if device.Algorithm != algorithm {
		t.Errorf("Algorithm doesn't match")
	}
	if device.Label != label {
		t.Errorf("Label doesn't match")
	}
	if device.SignatureCounter != 0 {
		t.Errorf("Initial signature counter should be zero")
	}
}

func TestNewSignatureDeviceECC(t *testing.T) {
	deviceId := "test-device-ecc"
	algorithm := "ECC"
	label := "Test Device ECC"

	device, err := NewSignatureDevice(deviceId, algorithm, label)
	if err != nil {
		t.Errorf("Error creating signature device: %v", err)
	}

	if device.Id != deviceId {
		t.Errorf("Device Id doesn't match")
	}
	if device.Algorithm != algorithm {
		t.Errorf("Algorithm doesn't match")
	}
	if device.Label != label {
		t.Errorf("Label doesn't match")
	}
	if device.SignatureCounter != 0 {
		t.Errorf("Initial signature counter should be zero")
	}
}

func TestNewSignatureDeviceUnsupported(t *testing.T) {
	deviceId := "test-device-unsupported"
	algorithm := "UNSUPPORTED"
	label := "Test Device Unsupported"

	_, err := NewSignatureDevice(deviceId, algorithm, label)
	if err == nil || !errors.Is(err, ErrUnsupportedAlgorithm) {
		t.Errorf("Expected unsupported algorithm error")
	}
}

func TestSignatureDeviceSignTransaction(t *testing.T) {
	deviceId := "test-device"
	algorithm := "RSA"
	label := "Test Device"

	device, _ := NewSignatureDevice(deviceId, algorithm, label)

	signature, signedData, err := device.SignTransaction("data-to-be-signed")
	if err != nil {
		t.Errorf("Error signing transaction: %v", err)
	}

	if len(signature) == 0 {
		t.Errorf("Signature should not be empty")
	}

	if len(signedData) == 0 {
		t.Errorf("Signed data should not be empty")
	}
	if device.SignatureCounter != 1 {
		t.Errorf("Signature counter should increment")
	}
}
