package app

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"

	"github.com/MarlikAlmighty/mdns/internal/gen/models"
)

type (
	// App application methods
	App interface {
		IPV4ToIPV6(ip string) (string, error)
	}
	Config interface {
	}
	// Resolver methods
	Resolver interface {
		Set(domain string, md *models.DNSEntry)
		Get(domain string) *models.DNSEntry
		Delete(domain string)
		GetMap() map[string]models.DNSEntry
		FetchCert(domain, ipv4 string) (*models.DNSEntry, error)
	}
)

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
