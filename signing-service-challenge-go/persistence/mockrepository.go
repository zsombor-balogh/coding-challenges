package persistence

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

type MockRepository struct {
	Devices map[string]*domain.SignatureDevice
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		Devices: make(map[string]*domain.SignatureDevice),
	}
}

func (r *MockRepository) SaveSignatureDevice(device *domain.SignatureDevice) error {
	r.Devices[device.Id] = device
	return nil
}

func (r *MockRepository) GetSignatureDevice(deviceId string) (*domain.SignatureDevice, error) {
	if device, ok := r.Devices[deviceId]; ok {
		return device, nil
	}
	return nil, domain.ErrDeviceNotFound
}

func (r *MockRepository) ListSignatureDevices() ([]*domain.SignatureDevice, error) {
	devices := make([]*domain.SignatureDevice, 0, len(r.Devices))
	for _, device := range r.Devices {
		devices = append(devices, device)
	}
	return devices, nil
}
