package persistence

import (
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

type InMemoryPersistence struct {
	devices map[string]*domain.SignatureDevice
	mutex   sync.RWMutex
}

func NewInMemoryPersistence() *InMemoryPersistence {
	return &InMemoryPersistence{
		devices: make(map[string]*domain.SignatureDevice),
	}
}

func (p *InMemoryPersistence) SaveSignatureDevice(device *domain.SignatureDevice) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.devices[device.Id] = device
	return nil
}

func (p *InMemoryPersistence) GetSignatureDevice(id string) (*domain.SignatureDevice, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	device, ok := p.devices[id]
	if !ok {
		return nil, domain.ErrDeviceNotFound
	}
	return device, nil
}

func (p *InMemoryPersistence) ListSignatureDevices() ([]*domain.SignatureDevice, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	var devices []*domain.SignatureDevice
	for _, device := range p.devices {
		devices = append(devices, device)
	}
	return devices, nil
}
