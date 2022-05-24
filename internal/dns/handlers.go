package dns

import (
	"github.com/miekg/dns"
	"log"
	"net"
	"strings"
)

// Handler serve dns requests
func (s *DNS) Handler(w dns.ResponseWriter, r *dns.Msg) {

	defer func() {
		if err := w.Close(); err != nil {
			log.Println(err)
		}
	}()

	msg := &dns.Msg{}
	msg.SetReply(r)
	src := msg.Question[0].Name
	domain := strings.ToLower(src)
	entry := s.Resolver.Get(domain)

	log.Printf("REQ: %v\n", domain)

	if entry.Domain == "" {
		var err error
		if msg, err = s.Lookup(r, s.NameServers); err != nil {
			log.Printf("ERROR: %v\n", err.Error())
			r.SetRcode(r, dns.RcodeServerFailure)
		}
	} else {

		msg.Authoritative = true

		switch r.Question[0].Qtype {

		case dns.TypeA:

			/*
				for _, sock := range entry.Ips {
					msg.Answer = append(msg.Answer,
						&dns.A{
							Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
							A:   net.ParseIP(sock),
						})
				}
			*/

			msg.Answer = append(msg.Answer,
				&dns.A{
					Hdr: dns.RR_Header{
						Name:   domain,
						Rrtype: dns.TypeA,
						Class:  dns.ClassINET,
						Ttl:    60},
					A: net.ParseIP(entry.IPV4),
				})

		case dns.TypeSOA:

			msg.Answer = append(msg.Answer, &dns.SOA{
				Hdr: dns.RR_Header{
					Name:   domain,
					Rrtype: dns.TypeSOA,
					Class:  dns.ClassINET,
					Ttl:    60,
				},
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
					Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: 60},
					Ns:  "ns1." + domain,
				},
				&dns.NS{
					Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: 60},
					Ns:  "ns2." + domain,
				})

		case dns.TypeAAAA:

			msg.Answer = append(msg.Answer, &dns.AAAA{
				Hdr:  dns.RR_Header{Name: domain, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: 60},
				AAAA: net.ParseIP(entry.IPV6),
			})

		case dns.TypePTR:

			ipAddress := net.ParseIP(entry.IPV4)
			reverseIpAddress := s.reverseIP(ipAddress) + ".in-addr.arpa."

			msg.Answer = append(msg.Answer,
				&dns.PTR{
					Hdr: dns.RR_Header{Name: reverseIpAddress, Rrtype: dns.TypePTR, Class: dns.ClassINET, Ttl: 60},
					Ptr: domain,
				})

		case dns.TypeSPF:

			spf := []string{"v=spf1 a mx ptr all"}

			msg.Answer = append(msg.Answer,
				&dns.SPF{
					Hdr: dns.RR_Header{
						Name:   domain,
						Rrtype: dns.TypeSPF,
						Class:  dns.ClassINET,
						Ttl:    60,
					},
					Txt: spf,
				})

		case dns.TypeTXT:

			//spf := []string{"v=spf1 a mx ptr all"}

			//dmarc := []string{"v=DMARC1; pct=1; p=reject; adkim=s; aspf=s"}

			// "v=DKIM1;k=rsa;p=..."

			msg.Answer = append(msg.Answer,
				/*
					&dns.SPF{
						Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeSPF, Class: dns.ClassINET, Ttl: 60},
						Txt: spf,
					},
					&dns.TXT{
						Hdr: dns.RR_Header{Name: "_dmarc." + domain, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 60},
						Txt: dmarc,
					},
					&dns.TXT{
						Hdr: dns.RR_Header{Name: "_domainkey." + domain, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 60},
						Txt: entry.Dkim,
					},
				*/
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
				Hdr: dns.RR_Header{
					Name:   domain,
					Rrtype: dns.TypeMX,
					Class:  dns.ClassINET,
					Ttl:    60,
				},
				Preference: 10,
				Mx:         "mail." + domain,
			})
		}
	}

	for k, v := range msg.Answer {
		log.Printf("RESP: %v %v\n", k, v)
	}

	if err := w.WriteMsg(msg); err != nil {
		log.Printf("WRITE MSG: %v\n", err.Error())
	}

	return
}
