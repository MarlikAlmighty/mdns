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

	name := msg.Question[0].Name
	domain := strings.ToLower(name)
	entry := s.Resolver.Get(domain)

	header := dns.RR_Header{
		Name:   msg.Question[0].Name,
		Rrtype: msg.Question[0].Qtype,
		Class:  dns.ClassINET,
		Ttl:    60,
	}

	log.Printf("[REQ]: %v %v\n", msg.Question[0].Name, dns.TypeToString[msg.Question[0].Qtype])

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

	if entry.Domain == "" {

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var err error
		if msg, err = s.Lookup(ctx, r, s.NameServers); err != nil {
			log.Printf("[LOOKUP ERROR]: %v\n", err)
			msg.SetRcodeFormatError(r)
		}
	}

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
	case dns.TypeSPF:
		s.spf(msg)
	case dns.TypeMX:
		s.mx(msg)
	}

	if len(msg.Answer) > 0 {
		for _, answer := range msg.Answer {
			log.Printf("[RESP]: %v\n", answer.String())
		}
	}

	w.WriteMsg(msg)
	return
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

	var dkim []string
	dkim = append(dkim, fmt.Sprintf("v=DKIM1; k=rsa; t=s; p=%s", entry.DkimPublicKey))
	dmarc := []string{"v=DMARC1; p=reject; sp=reject; pct=0; adkim=r; aspf=r"}

	msg.Answer = append(msg.Answer,

		&dns.TXT{
			Hdr: dns.RR_Header{
				Name:   "_dmarc." + entry.Domain,
				Rrtype: dns.TypeTXT,
				Class:  dns.ClassINET,
				Ttl:    60,
			},
			Txt: dmarc,
		},
		&dns.TXT{
			Hdr: dns.RR_Header{
				Name:   "mail._domainkey." + entry.Domain,
				Rrtype: dns.TypeTXT,
				Class:  dns.ClassINET,
				Ttl:    60,
			},
			Txt: dkim,
		},
		&dns.TXT{
			Hdr: dns.RR_Header{
				Name:   "_acme-challenge." + entry.Domain,
				Rrtype: dns.TypeTXT,
				Class:  dns.ClassINET,
				Ttl:    60,
			},
			Txt: entry.Acme,
		})
}

func (s *DNS) soa(msg *dns.Msg, entry *models.DNSEntry) {

	msg.Answer = append(msg.Answer,
		&dns.SOA{
			Hdr: dns.RR_Header{
				Name:   entry.Domain,
				Rrtype: dns.TypeSOA,
				Class:  dns.ClassINET,
				Ttl:    2399,
			},
			Ns:      "ns1." + entry.Domain,
			Mbox:    "admin." + entry.Domain,
			Serial:  258342863,
			Refresh: 900,
			Retry:   900,
			Expire:  1800,
			Minttl:  86399,
		})
}

func (s *DNS) ns(msg *dns.Msg, entry *models.DNSEntry) {

	msg.Answer = append(msg.Answer,
		&dns.NS{
			Hdr: dns.RR_Header{
				Name:   entry.Domain,
				Rrtype: dns.TypeNS,
				Class:  dns.ClassINET,
				Ttl:    86399,
			},
			Ns: "ns1." + entry.Domain,
		},
		&dns.NS{
			Hdr: dns.RR_Header{
				Name:   entry.Domain,
				Rrtype: dns.TypeNS,
				Class:  dns.ClassINET,
				Ttl:    86399,
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
				Name:   reverseIpAddress,
				Rrtype: dns.TypePTR,
				Class:  dns.ClassINET,
				Ttl:    86399,
			},
			Ptr: entry.Domain,
		})
}

func (s *DNS) spf(msg *dns.Msg) {

	outDot := strings.TrimSuffix(msg.Question[0].Name, ".")
	spf := []string{"v=spf1 include:_spf." + outDot + " a mx ptr ~all"}

	msg.Answer = append(msg.Answer,
		&dns.SPF{
			Hdr: dns.RR_Header{
				Name:   msg.Question[0].Name,
				Rrtype: dns.TypeSPF,
				Class:  dns.ClassINET,
				Ttl:    86399,
			},
			Txt: spf,
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
