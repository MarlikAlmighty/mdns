package data

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/MarlikAlmighty/mdns/internal/gen/models"
	"golang.org/x/crypto/acme"
	"strings"
	"sync"
	"time"
)

const acmeURL = "https://acme-staging-v02.api.letsencrypt.org/directory"

// Resolver for assertion
type Resolver interface {
	Set(domain string, md *models.DNSEntry)
	Get(domain string) *models.DNSEntry
	Delete(domain string)
	GetMap() map[string]models.DNSEntry
	FetchCert(domain, ipv4 string) (*models.DNSEntry, error)
	GenerateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error)
	ExportRsaPrivateKeyAsStr(privKey *rsa.PrivateKey) string
	ExportRsaPublicKeyAsStr(pubKey *rsa.PublicKey) (string, error)
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

/*
Start block for fetch HTTP certs
*/

// FetchCert fetch cert from ca
func (r *ResolvedData) FetchCert(domain, ipv4 string) (*models.DNSEntry, error) {
	outDot := strings.TrimSuffix(domain, ".")
	// check domain
	identifiers := acme.DomainIDs(strings.Fields(outDot)...)
	if len(identifiers) == 0 {
		return nil, errors.New("at least one domain is required")
	}
	// context with cancel
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	// create a new private key for account.
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("ecdsa.GenerateKey for a cert: %v", err)
	}
	// new client
	var cl *acme.Client
	if cl, err = r.newClient(ctx, k, acmeURL); err != nil {
		return nil, fmt.Errorf("register: %v", err)
	}
	// new order
	var o *acme.Order
	if o, err = cl.AuthorizeOrder(ctx, identifiers); err != nil {
		return nil, fmt.Errorf("authorizeOrder: %v", err)
	}
	// auth urls
	var challenge *acme.Challenge
	if challenge, err = r.authURL(ctx, o, cl); err != nil {
		return nil, err
	}
	// get token for dns record "_acme-challenge." + domain
	var token string
	if token, err = cl.DNS01ChallengeRecord(challenge.Token); err != nil {
		return nil, fmt.Errorf("dns01ChallengeRecord: %v", err)
	}
	var md models.DNSEntry
	md.Domain = domain
	md.IPV4 = ipv4
	md.Acme = []string{token}
	r.Set(domain, &md)
	// accept informs the server that the client accepts one of its challenges
	if _, err = cl.Accept(ctx, challenge); err != nil {
		return nil, fmt.Errorf("accept(%q): %v", challenge.URI, err)
	}
	// polls an authorization at the given URL
	if _, err = cl.WaitAuthorization(ctx, challenge.URI); err != nil {
		return nil, fmt.Errorf("waitAutorization %v", err)
	}
	var urls []string
	urls = append(urls, challenge.URI)
	// polls an order from the given URL
	if _, err = cl.WaitOrder(ctx, o.URI); err != nil {
		return nil, fmt.Errorf("waitOrder(%q): %v", o.URI, err)
	}
	// create a new private key for certs
	var newKey *ecdsa.PrivateKey
	if newKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader); err != nil {
		return nil, fmt.Errorf("ecdsa.GenerateKey for a cert: %v", err)
	}
	// create csr
	var csr []byte
	if csr, err = r.newCSR(identifiers, newKey); err != nil {
		return nil, err
	}
	// submits the CSR (Certificate Signing Request) to a CA at the specified URL
	var der [][]byte
	if der, _, err = cl.CreateOrderCert(ctx, o.FinalizeURL, csr, true); err != nil {
		return nil, fmt.Errorf("createOrderCert: %v", err)
	}
	// check cert
	var cert []byte
	if cert, err = r.checkCert(der, identifiers); err != nil {
		return nil, fmt.Errorf("invalid cert: %v", err)
	}
	// RevokeAuthorization relinquishes an existing authorization identified by the given URL
	for _, v := range urls {
		if err = cl.RevokeAuthorization(ctx, v); err != nil {
			return nil, fmt.Errorf("revokAuthorization(%q): %v", v, err)
		}
	}
	// DeactivateReg permanently disables an existing account associated with key
	if err = cl.DeactivateReg(ctx); err != nil {
		return nil, fmt.Errorf("deactivateReg: %v", err)
	}
	// RevokeCert revokes a previously issued certificate cert, provided in DER format
	if err = cl.RevokeCert(ctx, newKey, der[0], acme.CRLReasonCessationOfOperation); err != nil {
		return nil, fmt.Errorf("revokeCert: %v", err)
	}
	// convert public cert []byte to []string
	publicCert := r.exportPublicCert(cert)
	// convert private cert []byte to []string
	var privateKey string
	if privateKey, err = r.exportEDSAtPrivateKey(newKey); err != nil {
		return nil, err
	}
	md.Acme = []string{""}
	md.HTTPPublicKey = publicCert
	md.HTTPPrivateKey = privateKey
	r.Set(domain, &md)
	return &md, nil
}

// checkCert verifies the public certificate and returns it
func (r *ResolvedData) checkCert(derChain [][]byte, id []acme.AuthzID) ([]byte, error) {
	if len(derChain) == 0 {
		return nil, errors.New("cert chain is zero bytes")
	}
	var publicCert []byte
	for i, b := range derChain {
		crt, err := x509.ParseCertificate(b)
		if err != nil {
			return nil, fmt.Errorf("%d: ParseCertificate: %v", i, err)
		}
		if i > 0 {
			continue
		}
		publicCert = b
		for _, v := range id {
			if err = crt.VerifyHostname(v.Value); err != nil {
				return nil, err
			}
		}
	}
	return publicCert, nil
}

// newCSR creates a signing certificate
func (r *ResolvedData) newCSR(identifiers []acme.AuthzID, k *ecdsa.PrivateKey) ([]byte, error) {
	var cr x509.CertificateRequest
	for _, id := range identifiers {
		switch id.Type {
		case "dns":
			cr.DNSNames = append(cr.DNSNames, id.Value)
		default:
			return nil, fmt.Errorf("newCSR: unknown identifier type %q", id.Type)
		}
	}
	csr, err := x509.CreateCertificateRequest(rand.Reader, &cr, k)
	if err != nil {
		return nil, fmt.Errorf("newCSR: x509.CreateCertificateRequest: %v", err)
	}
	return csr, nil
}

// authURL represents authorizations to complete before a certificate
func (r *ResolvedData) authURL(ctx context.Context, o *acme.Order, client *acme.Client) (*acme.Challenge, error) {
	var (
		z         *acme.Authorization
		challenge *acme.Challenge
		err       error
	)
	for _, u := range o.AuthzURLs {
		if z, err = client.GetAuthorization(ctx, u); err != nil {
			return nil, fmt.Errorf("getAuthorization(%q): %v", u, err)
		}
		if z.Status != acme.StatusPending {
			return nil, fmt.Errorf("authz status is %q; skipping", z.Status)
		}
		for _, c := range z.Challenges {
			if c.Type == "dns-01" {
				challenge = c
				break
			}
		}
		if challenge == nil {
			return nil, fmt.Errorf("challenge type %q wasn't offered for authz %s", "dns-01", z.URI)
		}
	}
	return challenge, nil
}

// newClient create acme client
func (r *ResolvedData) newClient(ctx context.Context, k *ecdsa.PrivateKey, directoryURL string) (*acme.Client, error) {
	cl := &acme.Client{Key: k, DirectoryURL: directoryURL}
	a := &acme.Account{Contact: strings.Fields("")}
	if _, err := cl.Register(ctx, a, acme.AcceptTOS); err != nil {
		return &acme.Client{}, fmt.Errorf("register: %v", err)
	}
	return cl, nil
}

// exportEDSAPrivateKey creates a private pem block
func (r *ResolvedData) exportEDSAtPrivateKey(k *ecdsa.PrivateKey) (string, error) {
	key, err := x509.MarshalECPrivateKey(k)
	if err != nil {
		return "", err
	}
	block := pem.Block{Type: "ECDSA PRIVATE KEY", Bytes: key}
	ky := pem.EncodeToMemory(&block)
	return string(ky), nil
}

// exportPublicCert creates a public pem block
func (r *ResolvedData) exportPublicCert(b []byte) string {
	p := &pem.Block{Type: "CERTIFICATE", Bytes: b}
	c := pem.EncodeToMemory(p)
	return string(c)
}

/*
Start block for generate RSA keys
*/

func (r *ResolvedData) GenerateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	return privKey, &privKey.PublicKey, nil
}

func (r *ResolvedData) ExportRsaPrivateKeyAsStr(privKey *rsa.PrivateKey) string {
	return base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(privKey))
}

func (r *ResolvedData) ExportRsaPublicKeyAsStr(pubKey *rsa.PublicKey) (string, error) {
	pu, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(pu), nil
}
