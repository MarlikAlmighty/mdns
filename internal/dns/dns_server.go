package dns

import (
	"context"
	"fmt"
	"github.com/MarlikAlmighty/mdns/internal/config"
	"github.com/MarlikAlmighty/mdns/internal/data"
	"github.com/miekg/dns"
	"log"
	"net"
	"strings"
	"time"
)

// DNS wrapper over dns server
type DNS struct {
	TcpServer *dns.Server
	UdpServer *dns.Server
	Client    *dns.Client
	Resolver  *data.ResolvedData
	Config    *config.Configuration
}

// New simple constructor
func New(d *data.ResolvedData, cnf *config.Configuration) *DNS {
	t, u := &dns.Server{}, &dns.Server{}
	c := &dns.Client{
		Net: "udp",
	}
	return &DNS{
		TcpServer: t,
		UdpServer: u,
		Client:    c,
		Resolver:  d,
		Config:    cnf,
	}
}

// Run start dns server
func (s *DNS) Run() {

	tcpHandler := dns.NewServeMux()
	tcpHandler.HandleFunc(".", s.Handler)
	udpHandler := dns.NewServeMux()
	udpHandler.HandleFunc(".", s.Handler)

	s.TcpServer.Addr = "0.0.0.0" + ":" + s.Config.DnsTcpPort
	s.TcpServer.Net = "tcp4"
	s.TcpServer.Handler = tcpHandler

	s.UdpServer.Addr = "0.0.0.0" + ":" + s.Config.DnsUdpPort
	s.UdpServer.Net = "udp4"
	s.UdpServer.Handler = udpHandler

	go func() {
		log.Printf("Serving mdns on tcp %v \n", s.TcpServer.Addr)
		if err := s.TcpServer.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		log.Printf("Serving mdns on udp %v \n", s.UdpServer.Addr)
		if err := s.UdpServer.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
}

// Handler serve dns requests
func (s *DNS) Handler(w dns.ResponseWriter, r *dns.Msg) {

	defer func() {
		if err := w.Close(); err != nil {
			log.Printf("%v\n", err)
		}
	}()

	msg := &dns.Msg{}
	msg.SetReply(r)

	// for logging ***************************

	host, _, err := net.SplitHostPort(w.RemoteAddr().String())
	if err != nil {
		log.Printf("%v\n", w.RemoteAddr().String())
		return
	}

	// Logging request
	log.Printf("[REQ]: from: %v %v %v\n",
		host,
		msg.Question[0].Name,
		dns.TypeToString[msg.Question[0].Qtype])

	// *******************************************

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
			s.mx(msg, entry)
		default:
			s.soa(msg, entry)
		}

	} else {

		/*
			if not found domain on server, doing look up request in internet,
			but before check if request is not from local net
		*/
		n := net.ParseIP(host)
		if n.IsPrivate() || n.IsLoopback() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if msg, err = s.Lookup(ctx, r, s.Config.NameServers); err != nil {
				log.Printf("[ERR]: %v\n", err)
				return
			}
		} else {
			log.Printf("[ERR]: deny request from %v\n",
				n.String())
			return
		}
	}

	if err = w.WriteMsg(msg); err != nil {
		log.Printf("[ERR]: write msg %v\n", err)
	}
}

func (s *DNS) Lookup(ctx context.Context, req *dns.Msg, nameServers []string) (*dns.Msg, error) {

	var (
		r   *dns.Msg
		err error
	)

	answer := make(chan *dns.Msg, 1)

	for _, v := range nameServers {
		go func(v string, answer chan *dns.Msg) {
			r, _, err = s.Client.Exchange(req, v+":53")
			if err != nil {
				log.Printf("[ERR]: client exchange, host: %v err: %v\n", v, err)
				return
			}
			if r != nil {
				answer <- r
			}
		}(v, answer)
	}

	select {
	case a := <-answer:
		return a, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Close stop dns server
func (s *DNS) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	var errs []string
	log.Printf("Stopped serving tcp on %v \n", s.TcpServer.Addr)
	if err := s.TcpServer.ShutdownContext(ctx); err != nil {
		errs = append(errs, err.Error())
	}
	log.Printf("Stopped serving udp on %v \n", s.UdpServer.Addr)
	if err := s.UdpServer.ShutdownContext(ctx); err != nil {
		errs = append(errs, err.Error())
	}
	if len(errs) > 0 {
		return fmt.Errorf("errs: %s", strings.Join(errs, ","))
	}
	return nil
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
