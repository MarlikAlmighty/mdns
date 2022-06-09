package data

import (
	"github.com/MarlikAlmighty/mdns/internal/gen/models"
	"sync"
)

// Resolver for assertion
type Resolver interface {
	Set(domain string, md *models.DNSEntry)
	Get(domain string) *models.DNSEntry
	Delete(domain string)
	GetMap() map[string]models.DNSEntry
}

// ResolvedData saved records of dns
type ResolvedData struct {
	Records map[string]models.DNSEntry
	mux     sync.Mutex
}

// New simple constructor
func New() *ResolvedData {
	return &ResolvedData{
		Records: make(map[string]models.DNSEntry),
	}
}

// Set add data to map
func (r *ResolvedData) Set(domain string, md *models.DNSEntry) {
	r.mux.Lock()
	r.Records[domain] = *md
	r.mux.Unlock()
}

// Get fetch data from map by value
func (r *ResolvedData) Get(domain string) *models.DNSEntry {
	r.mux.Lock()
	md := r.Records[domain]
	r.mux.Unlock()
	return &md
}

// Delete record from map
func (r *ResolvedData) Delete(domain string) {
	r.mux.Lock()
	delete(r.Records, domain)
	r.mux.Unlock()
}

// GetMap get all map
func (r *ResolvedData) GetMap() map[string]models.DNSEntry {
	r.mux.Lock()
	mp := r.Records
	r.mux.Unlock()
	return mp
}
