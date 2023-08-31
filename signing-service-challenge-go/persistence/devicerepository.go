package persistence

import "github.com/fiskaly/coding-challenges/signing-service-challenge/domain"

type SignatureDeviceRepository interface {
	SaveSignatureDevice(device *domain.SignatureDevice) error
	GetSignatureDevice(id string) (*domain.SignatureDevice, error)
	ListSignatureDevices() ([]*domain.SignatureDevice, error)
}
