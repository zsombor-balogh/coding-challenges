package domain

import (
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
)

var (
	ErrDeviceNotFound       = fmt.Errorf("signature device not found")
	ErrUnsupportedAlgorithm = fmt.Errorf("unsupported algorithm")
)

type SignatureDevice struct {
	Id               string
	Algorithm        string
	Label            string
	SignatureCounter int
	LastSignature    string

	signerLock sync.Mutex
	signer     crypto.Signer
}

func NewSignatureDevice(id, algorithm, label string) (*SignatureDevice, error) {
	var signer crypto.Signer
	switch algorithm {
	case "RSA":
		rsaGenerator := crypto.RSAGenerator{}
		keyPair, err := rsaGenerator.Generate()
		if err != nil {
			return nil, err
		}
		signer = crypto.NewRSASigner(keyPair.Private)
	case "ECC":
		eccGenerator := crypto.ECCGenerator{}
		keyPair, err := eccGenerator.Generate()
		if err != nil {
			return nil, err
		}
		signer = crypto.NewECCSigner(keyPair.Private)
	default:
		return nil, ErrUnsupportedAlgorithm
	}

	return &SignatureDevice{
		Id:        id,
		Algorithm: algorithm,
		Label:     label,
		signer:    signer,
	}, nil
}

func (d *SignatureDevice) SignTransaction(dataToBeSigned string) (string, string, error) {
	d.signerLock.Lock()
	defer d.signerLock.Unlock()

	securedDataToBeSigned := fmt.Sprintf("%d_%s_%s", d.SignatureCounter, dataToBeSigned, d.LastSignature)
	signature, err := d.signer.Sign([]byte(securedDataToBeSigned))
	if err != nil {
		return "", "", err
	}

	var lastSignature string
	if d.SignatureCounter == 0 {
		lastSignature = base64.StdEncoding.EncodeToString([]byte(d.Id))
	} else {
		lastSignature = d.LastSignature
	}

	signedData := fmt.Sprintf("%d_%s_%s", d.SignatureCounter, dataToBeSigned, lastSignature)

	d.SignatureCounter++
	d.LastSignature = base64.StdEncoding.EncodeToString(signature)

	return base64.StdEncoding.EncodeToString(signature), signedData, nil
}
