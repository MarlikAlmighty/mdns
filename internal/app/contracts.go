package app

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/MarlikAlmighty/mdns/internal/gen/models"
	"net"
	"strings"
)

type (
	// App application methods
	App interface {
		ReverseIP(ip net.IP) string
		IPV4ToIPV6(ip string) (string, error)
	}
	// Resolver methods
	Resolver interface {
		Set(domain string, md *models.DNSEntry)
		Get(domain string) *models.DNSEntry
	}
	// Store methods
	Store interface {
	}
	// Config methods
	Config interface {
	}
)

// ReverseIP reverse ipv4 address for ptr
func (core *Core) ReverseIP(ip net.IP) string {
	if ip.To4() != nil {
		addressSlice := strings.Split(ip.String(), ".")
		var reverseSlice []string
		for i := range addressSlice {
			octet := addressSlice[len(addressSlice)-1-i]
			reverseSlice = append(reverseSlice, octet)
		}
		return strings.Join(reverseSlice, ".")
	}
	return ""
}

// IPV4ToIPV6 convert ipv4 addres to ipv6
func (core *Core) IPV4ToIPV6(ip string) (string, error) {
	a := net.ParseIP(ip)
	if a == nil {
		return "", errors.New("invalid ip")
	}
	dst := make([]byte, hex.EncodedLen(len(a)))
	hex.Encode(dst, a)
	return fmt.Sprintf("0:0:0:0:0:%s:%s:%s \n", dst[20:24], dst[24:28], dst[28:]), nil
}
