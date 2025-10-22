package scanner

import (
	"fmt"
	"net"
	"strings"

	"github.com/ToomyGO/pasunes_scan/internal/utils"
)

type ScanResult struct {
	Input    string
	IsIP     bool
	Resolved string
}

func ResolvedTarget(target string, verbose bool) (*ScanResult, error) {
	normalized, isIP, err := utils.NormalizeTarget(target)
	if err != nil {
		return nil, err
	}
	result := &ScanResult{
		Input: target,
		IsIP:  isIP,
	}

	if isIP {
		if verbose {
			fmt.Printf("🔍 Input '%s' normalized to IP: %s. Doing reverse lookup...\n", target, normalized)
		}
		names, err := net.LookupAddr(normalized)
		if err != nil || len(names) == 0 {
			result.Resolved = normalized
			if verbose {
				fmt.Printf("⚠️  Reverse lookup failed: %v\n", err)
			}
			return result, nil
		}
		result.Resolved = strings.TrimSuffix(names[0], ".")

		return result, nil
	}

	if verbose {
		fmt.Printf("🔍 Input '%s' normalized to hostname: %s. Doing forward lookup...\n", target, normalized)
	}
	ips, err := net.LookupIP(normalized)

	if err != nil || len(ips) == 0 {
		return nil, fmt.Errorf("no IP found for domain %s: %w", normalized, err)
	}

	result.Resolved = ips[0].String()
	return result, nil
	

}
