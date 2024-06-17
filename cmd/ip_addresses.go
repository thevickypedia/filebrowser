package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	"time"

	"github.com/thevickypedia/filebrowser/v2/settings"
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

func refreshAllowedOrigins(server *settings.Server) {
	if server.AllowPrivateIP {
		privateIP, err := GetLocalIP()
		if err != nil {
			log.Printf("Warning: Failed to get local IP: %v", err)
		} else {
			privateIPString := fmt.Sprintf("%s:%s", privateIP, server.Port)
			if !existsAlready(privateIPString, server.AllowedOrigins) {
				log.Printf("Adding local IP address [%s] to allowed origins", privateIPString)
				server.AllowedOrigins = append(server.AllowedOrigins, privateIPString)
			}
		}
	}
	if server.AllowPublicIP {
		publicIP := GetPublicIP()
		if publicIP == "" {
			log.Printf("Warning: Failed to get public IP")
		} else {
			publicIPString := fmt.Sprintf("%s:%s", publicIP, server.Port)
			if !existsAlready(publicIPString, server.AllowedOrigins) {
				log.Printf("Adding public IP address [%s] to allowed origins", publicIPString)
				server.AllowedOrigins = append(server.AllowedOrigins, publicIPString)
			}
		}
	}
}

// Function to start the background task with arguments
func startBackgroundTask(server *settings.Server) chan bool {
	log.Printf("Allowed origins will be refreshed every %d seconds", server.RefreshAllowedOrigins)
	ticker := time.NewTicker(time.Duration(server.RefreshAllowedOrigins) * time.Second)
	// Channel to signal the goroutine to stop
	done := make(chan bool)

	// Use a goroutine to handle the ticker
	go func() {
		for {
			select {
			case t := <-ticker.C:
				log.Printf("Refreshing allowed origins at %s", t.Format("2006-01-02 15:04:05 MST"))
				refreshAllowedOrigins(server)
			case <-done:
				// Stop the ticker and exit the goroutine
				ticker.Stop()
				log.Print("Stopping Goroutine to refresh origins")
				return
			}
		}
	}()

	return done
}
