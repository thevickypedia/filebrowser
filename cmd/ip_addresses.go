package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
)

// GetLocalIP returns the first non-loopback IPv4 address of the host machine.
func GetLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				// Check if it's not a loopback address and is IPv4
				if !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
					return ipNet.IP.String(), nil
				}
			}
		}
	}
	return "", fmt.Errorf("no suitable IP address found")
}

func GetPublicIP() string {
	// Regular expression to match IP address
	r := `^(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9])\.){3}(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9])$`
	ipRegex := regexp.MustCompile(r)

	// Define functions for handling different API responses
	opt1 := func(body []byte) string {
		return string(body)
	}
	opt2 := func(body []byte) string {
		var response struct {
			Origin string `json:"origin"`
		}
		err := json.Unmarshal(body, &response)
		if err != nil {
			log.Printf("Failed to get public IP from JSON payload.")
			return ""
		}
		return response.Origin
	}

	// Mapping of URLs to handling functions
	mapping := map[string]func([]byte) string{
		"https://checkip.amazonaws.com/": opt1,
		"https://api.ipify.org/":         opt1,
		"https://ipinfo.io/ip/":          opt1,
		"https://v4.ident.me/":           opt1,
		"https://httpbin.org/ip":         opt2,
		"https://myip.dnsomatic.com/":    opt1,
	}

	// Iterate over each URL and try to fetch public IP
	for url, handler := range mapping {
		resp, err := http.Get(url) //nolint:gosec
		if err != nil {
			continue
		}
		defer resp.Body.Close() //nolint:gocritic
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		ip := handler(body)
		if ipRegex.MatchString(ip) {
			return ip
		}
	}
	return ""
}

func existsAlready(addr string, addrArray []string) bool {
	for _, eachAddr := range addrArray {
		if addr == eachAddr {
			return true
		}
	}
	return false
}
