package dns

import (
	"fmt"
	"github.com/MarlikAlmighty/mdns/internal/gen/models"
	"github.com/miekg/dns"
	"log"
	"net"
	"strings"
)

func (s *DNS) a(msg *dns.Msg, entry *models.DNSEntry, header dns.RR_Header) {

	if len(entry.Ipv4s) > 0 {
		for _, ipv4 := range entry.Ipv4s {
			msg.Answer = append(msg.Answer,
				&dns.A{
					Hdr: header,
					A:   net.ParseIP(ipv4),
				})
		}
	}
}

func (s *DNS) aaaa(msg *dns.Msg, entry *models.DNSEntry, header dns.RR_Header) {

	if len(entry.Ipv6s) > 0 {
		for _, ipv6 := range entry.Ipv6s {
			msg.Answer = append(msg.Answer,
				&dns.AAAA{
					Hdr:  header,
					AAAA: net.ParseIP(ipv6),
				})
		}
	}
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

		var (
			ipv4 string
			spf  []string
		)

		ipv4 = strings.Join(entry.Ipv4s, ",")
		spf = append(spf, fmt.Sprintf("v=spf1 ip4:%v include:_spf.%v a mx ~all", ipv4, outDot))
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
