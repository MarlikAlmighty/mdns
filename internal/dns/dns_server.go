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
func (s *DNS) Run() error {

	tcpHandler := dns.NewServeMux()
	tcpHandler.HandleFunc(".", s.Handler)
	udpHandler := dns.NewServeMux()
	udpHandler.HandleFunc(".", s.Handler)
	s.TcpServer.Addr = s.Config.DnsHost + s.Config.DnsTcpPort
	s.TcpServer.Net = "tcp"
	s.TcpServer.Handler = tcpHandler
	s.UdpServer.Addr = s.Config.DnsHost + s.Config.DnsUdpPort
	s.UdpServer.Net = "udp"
	s.UdpServer.Handler = udpHandler

	var errs []string

	log.Printf("Serving mdns on tcp %v \n", s.TcpServer.Addr)

	go func() {
		if err := s.TcpServer.ListenAndServe(); err != nil {
			errs = append(errs, err.Error())
		}
	}()

	log.Printf("Serving mdns on udp %v \n", s.UdpServer.Addr)

	go func() {
		if err := s.UdpServer.ListenAndServe(); err != nil {
			errs = append(errs, err.Error())
		}
	}()

	if len(errs) > 0 {
		return fmt.Errorf("errs: %s", strings.Join(errs, ","))
	}

	return nil
}

func (s *DNS) Lookup(ctx context.Context, req *dns.Msg, nameServers []string) (*dns.Msg, error) {

	var (
		r   *dns.Msg
		err error
	)

	answer := make(chan *dns.Msg, 1)

	for _, v := range nameServers {

		go func(v string, answer chan *dns.Msg) {

			r, _, err = s.Client.Exchange(req, v)
			if err != nil {
				return
			}

			if r != nil && r.Rcode == dns.RcodeServerFailure {
				answer <- r.SetRcodeFormatError(r)
				return
			}

			if r.Rcode == dns.RcodeSuccess {
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

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
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
