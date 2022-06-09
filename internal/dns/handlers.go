package dns

import (
	"context"
	"fmt"
	"github.com/MarlikAlmighty/mdns/internal/gen/models"
	"github.com/miekg/dns"
	"log"
	"net"
	"strings"
	"time"
)

// Handler serve dns requests
func (s *DNS) Handler(w dns.ResponseWriter, r *dns.Msg) {

	defer w.Close()

	msg := &dns.Msg{}
	msg.SetReply(r)

	log.Printf("[REQ]: %v %v\n", msg.Question[0].Name, dns.TypeToString[msg.Question[0].Qtype])

	header := dns.RR_Header{
		Name:   msg.Question[0].Name,
		Rrtype: msg.Question[0].Qtype,
		Class:  dns.ClassINET,
		Ttl:    60,
	}

	// to lower case
	domain := strings.ToLower(msg.Question[0].Name)

	// find domain in map
	entry := s.Resolver.Get(domain)

	// find sub domain in map
	if entry.Domain == "" {
		mp := s.Resolver.GetMap()
		for k, v := range mp {
			if ok := dns.IsSubDomain(k, domain); ok {
				domain = k
				entry = &v
				break
			}
		}
	}

	// if domain or sub domain find
	if entry.Domain != "" {
		switch r.Question[0].Qtype {
		case dns.TypeA:
			s.a(msg, entry, header)
		case dns.TypeAAAA:
			s.aaaa(msg, entry, header)
		case dns.TypeCAA:
			s.caa(msg, header)
		case dns.TypeTXT:
			s.txt(msg, entry)
		case dns.TypeSOA:
			s.soa(msg, entry)
		case dns.TypeNS:
			s.ns(msg, entry)
		case dns.TypePTR:
			s.ptr(msg, entry)
		case dns.TypeMX:
			s.mx(msg)
		default:
			s.soa(msg, entry)
		}

	} else {

		// if not found domain on server, doing look up request in internet
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var err error
		if msg, err = s.Lookup(ctx, r, s.Config.NameServers); err != nil {
			log.Printf("[LOOKUP ERROR]: %v\n", err)
			msg.SetRcodeFormatError(r)
		}
	}

	if len(msg.Answer) > 0 {
		for _, answer := range msg.Answer {
			log.Printf("[RESP]: %v\n", answer.String())
		}
	}

	w.WriteMsg(msg)
}

func (s *DNS) a(msg *dns.Msg, entry *models.DNSEntry, header dns.RR_Header) {
	msg.Answer = append(msg.Answer,
		&dns.A{
			Hdr: header,
			A:   net.ParseIP(entry.IPV4),
		})
	if len(entry.Ips) > 0 {
		for _, sock := range entry.Ips {
			msg.Answer = append(msg.Answer,
				&dns.A{
					Hdr: header,
					A:   net.ParseIP(sock),
				})
		}
	}
}

func (s *DNS) aaaa(msg *dns.Msg, entry *models.DNSEntry, header dns.RR_Header) {
	msg.Answer = append(msg.Answer,
		&dns.AAAA{
			Hdr:  header,
			AAAA: net.ParseIP(entry.IPV6),
		})
}

func (s *DNS) caa(msg *dns.Msg, header dns.RR_Header) {
	msg.Answer = append(msg.Answer,
		&dns.CAA{
			Hdr:   header,
			Flag:  0,
			Tag:   "issue",
			Value: "letsencrypt.org",
		})
}

func (s *DNS) txt(msg *dns.Msg, entry *models.DNSEntry) {

	outDot := strings.TrimSuffix(msg.Question[0].Name, ".")
	outDot = strings.ToLower(outDot)
	serviceSegments := dns.SplitDomainName(outDot)
	outDotEntry := strings.TrimSuffix(entry.Domain, ".")

	switch serviceSegments[0] {

	case "mail":

		cert := fmt.Sprintf("v=DKIM1; k=rsa; t=s; p=%s", entry.DkimPublicKey)
		if len(cert) > dns.MinMsgSize {
			log.Printf("[ERR]: cert is %v over then dns.MinMsgSize \n", len(cert))
			return
		}

		if len(cert) == 0 {
			log.Printf("[ERR]: cert is %v zero length \n", len(cert))
			return
		}

		var dkim []string
		dkim = append(dkim, cert[:254])
		dkim = append(dkim, cert[254:])

		msg.Answer = append(msg.Answer,
			&dns.TXT{
				Hdr: dns.RR_Header{
					Name:   "mail._domainkey." + entry.Domain,
					Rrtype: dns.TypeTXT,
					Class:  dns.ClassINET,
					Ttl:    60,
				},
				Txt: dkim,
			})

	case "_dmarc":
		msg.Answer = append(msg.Answer,
			&dns.TXT{
				Hdr: dns.RR_Header{
					Name:   "_dmarc." + entry.Domain,
					Rrtype: dns.TypeTXT,
					Class:  dns.ClassINET,
					Ttl:    60,
				},
				Txt: []string{"DMARC1; p=reject; sp=reject; adkim=s; aspf=s; rua=mailto:admin@" + outDotEntry},
			})

	case "_acme-challenge":
		msg.Answer = append(msg.Answer,
			&dns.TXT{
				Hdr: dns.RR_Header{
					Name:   "_acme-challenge." + entry.Domain,
					Rrtype: dns.TypeTXT,
					Class:  dns.ClassINET,
					Ttl:    60,
				},
				Txt: entry.Acme,
			})

	default:
		spf := []string{"v=spf1 include:_spf." + outDot + " a mx ptr ~all"}
		msg.Answer = append(msg.Answer,
			&dns.TXT{
				Hdr: dns.RR_Header{
					Name:   entry.Domain,
					Rrtype: dns.TypeTXT,
					Class:  dns.ClassINET,
					Ttl:    60,
				},
				Txt: spf,
			})
	}
}

func (s *DNS) soa(msg *dns.Msg, entry *models.DNSEntry) {
	msg.Answer = append(msg.Answer,
		&dns.SOA{
			Hdr: dns.RR_Header{
				Name:   entry.Domain,
				Rrtype: dns.TypeSOA,
				Class:  dns.ClassINET,
				Ttl:    3600,
			},
			Ns:      "ns1." + entry.Domain,
			Mbox:    "admin." + entry.Domain,
			Serial:  2279726185,
			Refresh: 900,
			Retry:   900,
			Expire:  1800,
			Minttl:  3600,
		})
}

func (s *DNS) ns(msg *dns.Msg, entry *models.DNSEntry) {
	msg.Answer = append(msg.Answer,
		&dns.NS{
			Hdr: dns.RR_Header{
				Name:   entry.Domain,
				Rrtype: dns.TypeNS,
				Class:  dns.ClassINET,
				Ttl:    86400,
			},
			Ns: "ns1." + entry.Domain,
		},
		&dns.NS{
			Hdr: dns.RR_Header{
				Name:   entry.Domain,
				Rrtype: dns.TypeNS,
				Class:  dns.ClassINET,
				Ttl:    86400,
			},
			Ns: "ns2." + entry.Domain,
		})
}

func (s *DNS) ptr(msg *dns.Msg, entry *models.DNSEntry) {
	ipAddress := net.ParseIP(entry.IPV4)
	reverseIpAddress := s.reverseIP(ipAddress) + ".in-addr.arpa."
	msg.Answer = append(msg.Answer,
		&dns.PTR{
			Hdr: dns.RR_Header{
				Name:   entry.Domain,
				Rrtype: dns.TypePTR,
				Class:  dns.ClassINET,
				Ttl:    86399,
			},
			Ptr: reverseIpAddress,
		})
}

func (s *DNS) mx(msg *dns.Msg) {
	msg.Answer = append(msg.Answer,
		&dns.MX{
			Hdr: dns.RR_Header{
				Name:   msg.Question[0].Name,
				Rrtype: dns.TypeMX,
				Class:  dns.ClassINET,
				Ttl:    86399,
			},
			Preference: 10,
			Mx:         "mail." + msg.Question[0].Name,
		})
}
