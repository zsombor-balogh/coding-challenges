package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"testing"
)

func TestRSASigner(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal("Failed to generate RSA key:", err)
	}

	signer := NewRSASigner(privateKey)

	testData := []byte("Test-Data-RSA")
	hash := sha256.Sum256(testData)

	signature, err := signer.Sign(testData)
	if err != nil {
		t.Fatal("RSA signing failed:", err)
	}

	// Verify the signature
	err = rsa.VerifyPKCS1v15(&privateKey.PublicKey, crypto.SHA256, hash[:], signature)
	if err != nil {
		t.Fatal("RSA signature verification failed:", err)
	}
}

func TestECCSigner(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal("Failed to generate ECC key:", err)
	}

	signer := NewECCSigner(privateKey)

	testData := []byte("Test-Data-ECC")
	signature, err := signer.Sign(testData)
	if err != nil {
		t.Fatal("ECC signing failed:", err)
	}

	hash := sha256.Sum256(testData)

	valid := ecdsa.VerifyASN1(&signer.PrivateKey.PublicKey, hash[:], signature)

	if !valid {
		t.Fatal("ECC signature verification failed")
	}
}
