package main

import (
	"fmt"
	"os"
	"testing"
)

func main() {
	// Simple demonstration of the checks.
	spf := CheckSPF("v=spf1 ip4:192.0.2.1 -all", "192.0.2.1")
	fmt.Printf("SPF check: pass=%v, reason=%q\n", spf.Pass, spf.Reason)

	dkimHeaders := map[string]string{
		"From":    "alice@example.com",
		"Subject": "Hello",
	}
	dkim := CheckDKIM("v=1; a=rsa-sha256; d=example.com; s=selector; h=From:Subject", dkimHeaders)
	fmt.Printf("DKIM check: pass=%v, reason=%q\n", dkim.Pass, dkim.Reason)

	tlsa := CheckTLSA("v=0 1 1 5F:3A:...:AB", "5F:3A:...:AB")
	fmt.Printf("TLSA check: pass=%v, reason=%q\n", tlsa.Pass, tlsa.Reason)

	tls := CheckTLS("secure.mail.example.com", 587)
	fmt.Printf("TLS check: pass=%v, reason=%q\n", tls.Pass, tls.Reason)

	if !spf.Pass || !dkim.Pass || !tlsa.Pass || !tls.Pass {
		os.Exit(1)
	}
}

// Unit tests
func TestCheckSPF(t *testing.T) {
	tests := []struct {
		name   string
		record string
		ip     string
		want   bool
	}{
		{"valid", "v=spf1 ip4:192.0.2.1 -all", "192.0.2.1", true},
		{"invalid ip", "v=spf1 ip4:192.0.2.1 -all", "192.0.2.2", false},
		{"missing v", "ip4:192.0.2.1 -all", "192.0.2.1", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckSPF(tt.record, tt.ip)
			if got.Pass != tt.want {
				t.Errorf("CheckSPF() = %v, want %v", got.Pass, tt.want)
			}
		})
	}
}

func TestCheckDKIM(t *testing.T) {
	headers := map[string]string{
		"From":    "alice@example.com",
		"Subject": "Hello",
	}
	tests := []struct {
		name      string
		signature string
		wantPass  bool
	}{
		{"valid", "v=1; a=rsa-sha256; d=example.com; s=selector; h=From:Subject", true},
		{"missing d", "v=1; a=rsa-sha256; s=selector; h=From:Subject", false},
		{"missing header", "v=1; a=rsa-sha256; d=example.com; s=selector; h=From:Subject:To", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckDKIM(tt.signature, headers)
			if got.Pass != tt.wantPass {
				t.Errorf("CheckDKIM() = %v, want %v", got.Pass, tt.wantPass)
			}
		})
	}
}

func TestCheckTLSA(t *testing.T) {
	tests := []struct {
		name        string
		record      string
		fingerprint string
		wantPass    bool
	}{
		{"valid", "v=0 1 1 5F:3A:...:AB", "5F:3A:...:AB", true},
		{"missing v", "1 1 5F:3A:...:AB", "5F:3A:...:AB", false},
		{"mismatch", "v=0 1 1 5F:3A:...:AB", "00:11:22", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckTLSA(tt.record, tt.fingerprint)
			if got.Pass != tt.wantPass {
				t.Errorf("CheckTLSA() = %v, want %v", got.Pass, tt.wantPass)
			}
		})
	}
}

func TestCheckTLS(t *testing.T) {
	tests := []struct {
		name   string
		host   string
		port   int
		want   bool
	}{
		{"secure host", "secure.mail.example.com", 587, true},
		{"insecure host", "mail.example.com", 587, false},
		{"secure port", "mail.example.com", 443, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckTLS(tt.host, tt.port)
			if got.Pass != tt.want {
				t.Errorf("CheckTLS() = %v, want %v", got.Pass, tt.want)
			}
		})
	}
}

// Benchmarks
func BenchmarkCheckSPF(b *testing.B) {
	record := "v=spf1 ip4:192.0.2.1 -all"
	ip := "192.0.2.1"
	for i := 0; i < b.N; i++ {
		CheckSPF(record, ip)
	}
}

func BenchmarkCheckDKIM(b *testing.B) {
	signature := "v=1; a=rsa-sha256; d=example.com; s=selector; h=From:Subject"
	headers := map[string]string{
		"From":    "alice@example.com",
		"Subject": "Hello",
	}
	for i := 0; i < b.N; i++ {
		CheckDKIM(signature, headers)
	}
}

func BenchmarkCheckTLSA(b *testing.B) {
	record := "v=0 1 1 5F:3A:...:AB"
	fingerprint := "5F:3A:...:AB"
	for i := 0; i < b.N; i++ {
		CheckTLSA(record, fingerprint)
	}
}

func BenchmarkCheckTLS(b *testing.B) {
	host := "secure.mail.example.com"
	port := 587
	for i := 0; i < b.N; i++ {
		CheckTLS(host, port)
	}
}
