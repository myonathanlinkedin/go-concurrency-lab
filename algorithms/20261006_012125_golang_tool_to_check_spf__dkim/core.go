package main

import (
	"strings"
)

// SPFResult represents the outcome of an SPF check.
type SPFResult struct {
	Pass   bool
	Reason string
}

// CheckSPF evaluates a raw SPF record against a given IP address.
// It performs a minimal validation: the record must contain "v=spf1" and
// an "ip4:" or "ip6:" mechanism that matches the supplied IP.
func CheckSPF(record, ip string) SPFResult {
	if !strings.Contains(record, "v=spf1") {
		return SPFResult{false, "missing v=spf1"}
	}
	if strings.Contains(record, "ip4:"+ip) || strings.Contains(record, "ip6:"+ip) {
		return SPFResult{true, ""}
	}
	return SPFResult{false, "ip not allowed"}
}

// DKIMResult represents the outcome of a DKIM signature verification.
type DKIMResult struct {
	Pass   bool
	Reason string
}

// CheckDKIM performs a lightweight verification of a DKIM signature string.
// It checks for the presence of "d=", "s=", and "h=" tags and ensures that
// all header fields listed in "h=" exist in the provided headers map.
func CheckDKIM(signature string, headers map[string]string) DKIMResult {
	if !strings.Contains(signature, "d=") || !strings.Contains(signature, "s=") {
		return DKIMResult{false, "missing d= or s="}
	}
	hIndex := strings.Index(signature, "h=")
	if hIndex == -1 {
		return DKIMResult{false, "missing h="}
	}
	hPart := signature[hIndex+2:]
	fields := strings.Split(hPart, ":")
	for _, f := range fields {
		if _, ok := headers[f]; !ok {
			return DKIMResult{false, "missing header field: " + f}
		}
	}
	return DKIMResult{true, ""}
}

// TLSAResult represents the outcome of a TLSA record validation.
type TLSAResult struct {
	Pass   bool
	Reason string
}

// CheckTLSA validates a TLSA record against a certificate fingerprint.
// It checks for the presence of "v=0" and that the record contains the fingerprint.
func CheckTLSA(record, fingerprint string) TLSAResult {
	if !strings.Contains(record, "v=0") {
		return TLSAResult{false, "missing v=0"}
	}
	if !strings.Contains(record, fingerprint) {
		return TLSAResult{false, "fingerprint mismatch"}
	}
	return TLSAResult{true, ""}
}

// TLSResult represents the outcome of a TLS connectivity check.
type TLSResult struct {
	Pass   bool
	Reason string
}

// CheckTLS performs a mock TLS connectivity test.
// It accepts a host and port and returns Pass if the host contains "secure"
// or the port is 443; otherwise it fails.
func CheckTLS(host string, port int) TLSResult {
	if strings.Contains(host, "secure") || port == 443 {
		return TLSResult{true, ""}
	}
	return TLSResult{false, "insecure host or port"}
}
