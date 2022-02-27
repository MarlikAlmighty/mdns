package data

import (
	"github.com/MarlikAlmighty/mdns/internal/app"
	"github.com/MarlikAlmighty/mdns/internal/config"
	"github.com/MarlikAlmighty/mdns/internal/gen/models"
	"github.com/go-openapi/strfmt"
	"sync"
)

// TODO Change me
var _ app.Resolver = &ResolvedData{}

// Resolver for assertion
type Resolver interface {
	Set(domain string, md *models.DNSEntry)
	Get(domain string) *models.DNSEntry
}

// ResolvedData saved records of dns
type ResolvedData struct {
	Records map[string]*models.DNSEntry
	mux     sync.Mutex
}

// NewResolvedData simple constructor
func NewResolvedData(c *config.Configuration) *ResolvedData {
	mp := make(map[string]*models.DNSEntry)
	md := &models.DNSEntry{
		Dkim:   c.PrivateKey,
		Domain: c.Domain,
		IPV4:   strfmt.IPv4(c.IPV4),
		IPV6:   strfmt.IPv6(c.IPV6),
	}
	mp[c.Domain] = md
	return &ResolvedData{
		Records: mp,
	}
}

func (r *ResolvedData) Set(domain string, md *models.DNSEntry) {
	r.mux.Lock()
	r.Records[domain] = md
	r.mux.Unlock()
}

func (r *ResolvedData) Get(domain string) *models.DNSEntry {
	r.mux.Lock()
	md := r.Records[domain]
	r.mux.Unlock()
	return md
}
