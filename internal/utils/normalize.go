package utils

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

func NormalizeTarget(input string) (string, bool, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", false, fmt.Errorf("empty target")
	}

	raw := input
	if !strings.Contains(raw, "://") {
		raw = "http://" + raw
	}

	u, err := url.Parse(raw)
	if err != nil {
		// در صورتی که url.Parse شکست خورد، سعی می‌کنیم مستقیماً آن را به عنوان host/ipv4/ipv6 بررسی کنیم
		// مثلاً ورودی‌هایی مثل "127.0.0.1:8080" ممکن است parse کامل نداشته باشند.
		// در این حالت تلاش می‌کنیم بخش قبل از ":" را جدا و بررسی کنیم.
		// اما اگر parse نرمال شد، ادامه می‌دهیم.
		// بازگرداندن خطا تا caller بداند.
		return "", false, fmt.Errorf("invalid target '%s': %w", input, err)
	}

	host := u.Hostname()

	if host == "" {

		h := raw

		if idx := strings.Index(h, "://"); idx != -1 {
			h = h[idx+3:]

		}

		if idx := strings.Index(h, ":/"); idx != -1 {
			h = h[:idx]

		}
		if idx := strings.LastIndex(h, "@"); idx != -1 {
			h = h[idx+1:]
		}
		if idx := strings.LastIndex(h, ":"); idx != -1 {
			h = h[:idx]
		}
		host = h
	}

	host = strings.TrimPrefix(host, "[")
	host = strings.TrimSuffix(host, "]")

	if net.ParseIP(host) != nil {
		return host, true, nil
	}

	if strings.ContainsAny(host, "/\\") || host == "" {
		return "", false, fmt.Errorf("invalid hostname after normalization: '%s'", host)
	}
	return host, false, nil
}
