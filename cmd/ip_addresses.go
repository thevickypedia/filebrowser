package cmd

import (
	"encoding/json"
	"io/ioutil"
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
			switch v := addr.(type) {
			case *net.IPNet:
				// Check if it's not a loopback address and is IPv4
				if !v.IP.IsLoopback() && v.IP.To4() != nil {
					return v.IP.String(), nil
				}
			}
		}
	}
	return "", nil
}

func GetPublicIP() string {
	// Regular expression to match IP address
	ipRegex := regexp.MustCompile(`^(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9])\.){3}(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9])$`)

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
		resp, err := http.Get(url)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
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
