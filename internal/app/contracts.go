package app

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"

	"github.com/MarlikAlmighty/mdns/internal/config"
	"github.com/MarlikAlmighty/mdns/internal/gen/models"
)

type (
	// App application methods
	App interface {
		IPV4ToIPV6(ip string) (string, error)
	}
	// Resolver methods
	Resolver interface {
		Set(domain string, md *models.DNSEntry)
		Get(domain string) *models.DNSEntry
		Delete(domain string)
		SetMap(mp map[string]models.DNSEntry)
		GetMap() map[string]models.DNSEntry
		FetchCert(cnf *config.Configuration) (*models.DNSEntry, error)
	}
)

// IPV4ToIPV6 convert ipv4 addres to ipv6
func (core *Core) IPV4ToIPV6(ip string) (string, error) {
	a := net.ParseIP(ip)
	if a == nil {
		return "", errors.New("invalid ip")
	}
	dst := make([]byte, hex.EncodedLen(len(a)))
	hex.Encode(dst, a)
	return fmt.Sprintf("0:0:0:0:0:%s:%s:%s", dst[20:24], dst[24:28], dst[28:]), nil
}
