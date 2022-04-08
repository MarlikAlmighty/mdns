package dns

import (
	"github.com/MarlikAlmighty/mdns/internal/data"
	"github.com/miekg/dns"
	"log"
	"net"
	"strings"
)

// DNS wrapper over dns server
type DNS struct {
	Server   *dns.Server
	Resolver data.Resolver
}

// New simple constructor
func New(port string, data data.Resolver) *DNS {
	return &DNS{
		Server: &dns.Server{
			Addr:      ":" + port,
			Net:       "udp",
			ReusePort: true,
		},
		Resolver: data,
	}
}

// Run start dns server
func (s *DNS) Run() error {
	dns.HandleFunc(".", s.Handler())
	log.Printf("Serving mdns on %v \n", s.Server.Addr)
	if err := s.Server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

// Close stop dns server
func (s *DNS) Close() error {
	log.Printf("Stopped serving mdns on %v \n", s.Server.Addr)
	if err := s.Server.Shutdown(); err != nil {
		return err
	}
	return nil
}

// Handler serve dns requests
func (s *DNS) Handler() dns.HandlerFunc {
	return func(w dns.ResponseWriter, r *dns.Msg) {

		msg := &dns.Msg{}

		msg.SetReply(r)

		msg.Authoritative = true

		domain := msg.Question[0].Name

		entry := s.Resolver.Get(domain)

		switch r.Question[0].Qtype {

		case dns.TypeSOA:

			msg.Answer = append(msg.Answer, &dns.SOA{
				Hdr:     dns.RR_Header{Name: domain, Rrtype: dns.TypeSOA, Class: dns.ClassINET, Ttl: 60},
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

		case dns.TypeA:

			msg.Answer = append(msg.Answer,
				&dns.A{
					Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
					A:   net.ParseIP(entry.IPV4),
				},
				&dns.A{Hdr: dns.RR_Header{Name: "ns1." + domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
					A: net.ParseIP(entry.IPV4)},

				&dns.A{Hdr: dns.RR_Header{Name: "ns2." + domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
					A: net.ParseIP(entry.IPV4),
				},
				&dns.A{Hdr: dns.RR_Header{Name: "mail." + domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
					A: net.ParseIP(entry.IPV4),
				})

			for _, sock := range entry.Ips {
				msg.Answer = append(msg.Answer,
					&dns.A{
						Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
						A:   net.ParseIP(sock),
					})
			}

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
					Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeSPF, Class: dns.ClassINET, Ttl: 60},
					Txt: spf,
				})

		case dns.TypeTXT:

			spf := []string{"v=spf1 a mx ptr all"}

			dmarc := []string{"v=DMARC1; pct=1; p=reject; adkim=s; aspf=s"}

			// "v=DKIM1;k=rsa;p=..."

			msg.Answer = append(msg.Answer,
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
				&dns.TXT{
					Hdr: dns.RR_Header{Name: "_acme-challenge." + domain, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 60},
					Txt: entry.Acme,
				})

		case dns.TypeMX:

			msg.Answer = append(msg.Answer, &dns.MX{
				Hdr:        dns.RR_Header{Name: domain, Rrtype: dns.TypeMX, Class: dns.ClassINET, Ttl: 60},
				Preference: 10,
				Mx:         "mail." + domain,
			})
		}

		if err := w.WriteMsg(msg); err != nil {
			log.Println(err)
		}
	}
}

// reverseIP reverse ipv4 address for ptr
func (s *DNS) reverseIP(ip net.IP) string {
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
