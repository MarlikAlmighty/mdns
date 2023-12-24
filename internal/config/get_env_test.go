package config

import (
	"os"
	"reflect"
	"testing"
)

func TestConfiguration_GetEnv(t *testing.T) {

	if err := os.Setenv("HTTP_PORT", "8081"); err != nil {
		t.Errorf("Error: %v", err)
	}

	if err := os.Setenv("DNS_TCP_PORT", "1053"); err != nil {
		t.Errorf("Error: %v", err)
	}

	if err := os.Setenv("DNS_UDP_PORT", "1053"); err != nil {
		t.Errorf("Error: %v", err)
	}

	if err := os.Setenv("NAME_SERVERS", "1.1.1.1,1.0.0.1,8.8.8.8,8.8.4.4"); err != nil {
		t.Errorf("Error: %v", err)
	}

	cfg := New()
	if err := cfg.GetEnv(); err != nil {
		t.Errorf("Error: %v", err)
	}

	tests := []struct {
		name    string
		fields  Configuration
		wantErr bool
	}{
		{"get_env", *cfg, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cnf := &Configuration{
				HTTPPort:    tt.fields.HTTPPort,
				DnsTcpPort:  tt.fields.DnsTcpPort,
				DnsUdpPort:  tt.fields.DnsUdpPort,
				NameServers: tt.fields.NameServers,
			}
			if err := cnf.GetEnv(); (err != nil) != tt.wantErr {
				t.Errorf("GetEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNew(t *testing.T) {
	cfg := New()
	tests := []struct {
		name string
		want *Configuration
	}{
		{"new_config", cfg},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
