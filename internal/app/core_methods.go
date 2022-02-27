package app

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/MarlikAlmighty/mdns/internal/gen/models"
	"net"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/miekg/dns"

	"github.com/MarlikAlmighty/mdns/internal/config"
)

type (
	// App application methods
	App interface {
		Handler() dns.HandlerFunc
		CertsFromFile() error
		GenerateCerts() error
		IPV4ToIPV6(ip string) (string, error)
		ReverseIP(ip net.IP) string
	}
	// Resolver methods
	Resolver interface {
		Set(domain string, md *models.DNSEntry)
		Get(domain string) *models.DNSEntry
	}
)

// CertsFromFile get certs from file
func (core *Core) CertsFromFile() ([]byte, []byte, error) {
	crt := path.Join(core.Config.(*config.Configuration).CertDir,
		core.Config.(*config.Configuration).Domain) + ".crt"
	key := path.Join(core.Config.(*config.Configuration).CertDir,
		core.Config.(*config.Configuration).Domain) + ".key"
	certPEMBlock, err := os.ReadFile(crt)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	keyPEMBlock, err := os.ReadFile(key)
	if err != nil {
		return []byte{}, []byte{}, err
	}
	return certPEMBlock, keyPEMBlock, nil
}

// GenerateCerts generating certs on host
func (core *Core) GenerateCerts() error {
	if core.Config.(*config.Configuration).CertDir == "" {
		return errors.New("folder for certs not defined")
	}
	if core.Config.(*config.Configuration).Domain == "" {
		return errors.New("domain not defined")
	}
	a := net.ParseIP(core.Config.(*config.Configuration).IPV4)
	if a == nil || a.IsLoopback() || a.IsPrivate() {
		return errors.New("invalid ipv4 address")
	}
	if _, err := os.Stat(core.Config.(*config.Configuration).CertDir); err != nil {
		if os.IsNotExist(err) {
			key := path.Join(core.Config.(*config.Configuration).CertDir,
				core.Config.(*config.Configuration).Domain) + ".key"
			crt := path.Join(core.Config.(*config.Configuration).CertDir,
				core.Config.(*config.Configuration).Domain) + ".crt"
			if err := os.MkdirAll(core.Config.(*config.Configuration).CertDir, 0744); err != nil {
				return errors.New("make dir: " + err.Error())
			}
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, "/usr/bin/openssl",
				"req", "-x509", "-nodes", "-days", "365",
				"-newkey", "rsa:2048", "-keyout", key,
				"-out", crt, "-subj",
				"/C=US/ST=Oregon/L=Portland/O=Marlik Almighty/OU=Org/CN="+core.Config.(*config.Configuration).IPV4)
			if err := cmd.Run(); err != nil {
				return errors.New("cmd run: " + err.Error())
			}
		}
	}
	return nil
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
