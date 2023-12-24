package app

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"testing"
)

func TestGenerateRsaKeyPair(t *testing.T) {
	core := &Core{}

	privKey, pubKey, err := core.GenerateRsaKeyPair()
	if err != nil {
		t.Errorf("Error generating RSA key pair: %v", err)
	}

	if privKey == nil {
		t.Error("Private key is nil")
	}

	if pubKey == nil {
		t.Error("Public key is nil")
	}
}

func TestExportRsaPrivateKeyAsStr(t *testing.T) {
	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	expectedStr := generateExpectedStr(privKey)

	core := &Core{}

	result := core.ExportRsaPrivateKeyAsStr(privKey)

	if result != expectedStr {
		t.Errorf("ExportRsaPrivateKeyAsStr() = %s, expected %s", result, expectedStr)
	}
}

func generateExpectedStr(privKey *rsa.PrivateKey) string {
	return base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(privKey))
}

func TestIPV4ToIPV6(t *testing.T) {
	core := &Core{}
	ip := "192.168.0.1"
	expected := "::ffff:c0a8:0001"
	result, err := core.IPV4ToIPV6(ip)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if result != expected {
		t.Errorf("Expected %s, but got: %s", expected, result)
	}
}
