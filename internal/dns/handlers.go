package dns

import (
	"context"
	"fmt"
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

	log.Printf("[REQ]: %v %v\n", msg.Question[0].Name, msg.Question[0].Qtype)

	if entry.Domain == "" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var err error
		if msg, err = s.Lookup(ctx, r, s.NameServers); err != nil {
			log.Printf("[LOOKUP ERROR]: %v\n", err)
			return
		}
	}

	msg.Authoritative = true

	header := dns.RR_Header{
		Name:   msg.Question[0].Name,
		Rrtype: msg.Question[0].Qtype,
		Class:  dns.ClassINET,
		Ttl:    60,
	}

	switch r.Question[0].Qtype {

	case dns.TypeA:

		msg.Answer = append(msg.Answer,
			&dns.A{
				Hdr: header,
				A:   net.ParseIP(entry.IPV4),
			})

		if len(entry.Ips) > 0 {
			for _, sock := range entry.Ips {
				msg.Answer = append(msg.Answer,
					&dns.A{
						Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
						A:   net.ParseIP(sock),
					})
			}
		}

	case dns.TypeSOA:

		msg.Answer = append(msg.Answer,
			&dns.SOA{
				Hdr:     header,
				Ns:      "ns1." + domain,
				Mbox:    "admin." + domain,
				Serial:  258342863,
				Refresh: 300,
				Retry:   300,
				Expire:  1200,
				Minttl:  60,
			})

	case dns.TypeNS:

		msg.Answer = append(msg.Answer,
			&dns.NS{
				Hdr: header,
				Ns:  "ns1." + domain,
			},
			&dns.NS{
				Hdr: header,
				Ns:  "ns2." + domain,
			})

	case dns.TypeAAAA:

		msg.Answer = append(msg.Answer,
			&dns.AAAA{
				Hdr:  header,
				AAAA: net.ParseIP(entry.IPV6),
			})

	case dns.TypePTR:

		ipAddress := net.ParseIP(entry.IPV4)
		reverseIpAddress := s.reverseIP(ipAddress) + ".in-addr.arpa."

		msg.Answer = append(msg.Answer,
			&dns.PTR{
				Hdr: dns.RR_Header{
					Name:   reverseIpAddress,
					Rrtype: dns.TypePTR,
					Class:  dns.ClassINET,
					Ttl:    60},
				Ptr: domain,
			})

	case dns.TypeSPF:

		spf := []string{"v=spf1 a mx ptr all"}

		msg.Answer = append(msg.Answer,
			&dns.SPF{
				Hdr: header,
				Txt: spf,
			})

	case dns.TypeTXT:

		var (
			dkim []string
			adsp []string
		)

		outDot := strings.TrimSuffix(domain, ".")
		spf := []string{"v=spf1 a mx ptr all"}
		spf2 := []string{"v=spf1 include:_spf." + outDot + " ~all"}
		dmarc := []string{"v=DMARC1; pct=1; p=reject; adkim=s; aspf=s"}
		dkim = append(dkim, fmt.Sprintf("v=DKIM1; k=rsa; t=s; p=%s", entry.Dkim))
		adsp = append(adsp, "dkim=all")

		msg.Answer = append(msg.Answer,

			&dns.SPF{
				Hdr: header,
				Txt: spf,
			},
			&dns.TXT{
				Hdr: header,
				Txt: spf2,
			},
			&dns.TXT{
				Hdr: dns.RR_Header{
					Name:   "_dmarc." + domain,
					Rrtype: dns.TypeTXT,
					Class:  dns.ClassINET,
					Ttl:    60},
				Txt: dmarc,
			},
			&dns.TXT{
				Hdr: dns.RR_Header{
					Name:   "mail._domainkey." + domain,
					Rrtype: dns.TypeTXT,
					Class:  dns.ClassINET, Ttl: 60},
				Txt: dkim,
			},
			&dns.TXT{
				Hdr: dns.RR_Header{
					Name:   "_adsp._domainkey." + domain,
					Rrtype: dns.TypeTXT,
					Class:  dns.ClassINET, Ttl: 60},
				Txt: adsp,
			},
			&dns.TXT{
				Hdr: dns.RR_Header{
					Name:   "_acme-challenge." + domain,
					Rrtype: dns.TypeTXT,
					Class:  dns.ClassINET,
					Ttl:    60},
				Txt: entry.Acme,
			})

	case dns.TypeMX:

		msg.Answer = append(msg.Answer, &dns.MX{
			Hdr:        header,
			Preference: 10,
			Mx:         "mail." + domain,
		})
		msg.Answer = append(msg.Answer, &dns.MX{
			Hdr:        header,
			Preference: 10,
			Mx:         "smtp." + domain,
		})
	}

	if len(msg.Answer) > 0 {
		log.Printf("[RESP]: %v\n", msg.Answer)
	}

	w.WriteMsg(msg)
}
