package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

type RSASigner struct {
	PrivateKey *rsa.PrivateKey
}

func NewRSASigner(privateKey *rsa.PrivateKey) RSASigner {
	return RSASigner{PrivateKey: privateKey}
}

func (signer RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hash := sha256.Sum256(dataToBeSigned)
	signature, err := rsa.SignPKCS1v15(rand.Reader, signer.PrivateKey, crypto.SHA256, hash[:])
	if err != nil {
		return nil, err
	}
	return signature, nil
}

type ECCSigner struct {
	PrivateKey *ecdsa.PrivateKey
}

func NewECCSigner(privateKey *ecdsa.PrivateKey) ECCSigner {
	return ECCSigner{PrivateKey: privateKey}
}

func (signer ECCSigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hash := sha256.Sum256(dataToBeSigned)
	signature, err := ecdsa.SignASN1(rand.Reader, signer.PrivateKey, hash[:])
	if err != nil {
		return nil, err
	}

	return signature, nil
}
