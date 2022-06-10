package app

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/MarlikAlmighty/mdns/internal/gen/models"
	"net"
)

type (
	// App application methods
	App interface {
		GenerateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error)
		ExportRsaPrivateKeyAsStr(privKey *rsa.PrivateKey) string
		ExportRsaPublicKeyAsStr(pubKey *rsa.PublicKey) (string, error)
		IPV4ToIPV6(ip string) (string, error)
	}
	// Resolver methods
	Resolver interface {
		Set(domain string, md *models.DNSEntry)
		Get(domain string) *models.DNSEntry
		Delete(domain string)
		GetMap() map[string]models.DNSEntry
	}
	Config interface {
	}
)

// GenerateRsaKeyPair generate rsa pair
func (core *Core) GenerateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	return privKey, &privKey.PublicKey, nil
}

// ExportRsaPrivateKeyAsStr encode *rsa.PrivateKey to base64
func (core *Core) ExportRsaPrivateKeyAsStr(privKey *rsa.PrivateKey) string {
	return base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(privKey))
}

// ExportRsaPublicKeyAsStr encode *rsa.PublicKey to base64
func (core *Core) ExportRsaPublicKeyAsStr(pubKey *rsa.PublicKey) (string, error) {
	pu, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(pu), nil
}

// IPV4ToIPV6 convert ipv4 address to ipv6
func (core *Core) IPV4ToIPV6(ip string) (string, error) {
	a := net.ParseIP(ip)
	if a == nil {
		return "", errors.New("invalid ip")
	}
	dst := make([]byte, hex.EncodedLen(len(a)))
	hex.Encode(dst, a)
	return fmt.Sprintf("::%s:%s:%s", dst[20:24], dst[24:28], dst[28:]), nil
}
