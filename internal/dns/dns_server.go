package dns

import (
	"context"
	"fmt"
	"github.com/MarlikAlmighty/mdns/internal/data"
	"github.com/miekg/dns"
	"log"
	"net"
	"strings"
	"time"
)

// DNS wrapper over dns server
type DNS struct {
	TcpServer   *dns.Server
	UdpServer   *dns.Server
	Client      *dns.Client
	Resolver    data.Resolver
	NameServers []string
	IPV4        string
}

// New simple constructor
func New(nameServers []string, d data.Resolver) *DNS {

	t, u := &dns.Server{}, &dns.Server{}

	c := &dns.Client{
		Net:          "tcp",
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second * 3,
	}

	return &DNS{
		TcpServer:   t,
		UdpServer:   u,
		Client:      c,
		Resolver:    d,
		NameServers: nameServers,
	}
}

// Run start dns server
func (s *DNS) Run() error {

	tcpHandler := dns.NewServeMux()
	tcpHandler.HandleFunc(".", s.Handler)
	udpHandler := dns.NewServeMux()
	udpHandler.HandleFunc(".", s.Handler)
	s.TcpServer.Addr = s.IPV4 + ":53"
	s.TcpServer.Net = "tcp"
	s.TcpServer.Handler = tcpHandler
	s.UdpServer.Addr = s.IPV4 + ":53"
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

func (s *DNS) Lookup(req *dns.Msg, nameServers []string) (*dns.Msg, error) {

	var (
		r   *dns.Msg
		err error
	)

	for _, v := range nameServers {
		r, _, err = s.Client.Exchange(req, v)
		if err != nil {
			return nil, err
		}
		if r != nil && r.Rcode == dns.RcodeServerFailure {
			continue
		} else if r.Rcode == dns.RcodeSuccess {
			break
		}
	}

	return r, nil
}

// Close stop dns server
func (s *DNS) Close() error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
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
