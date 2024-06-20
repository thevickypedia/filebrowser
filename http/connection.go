package http

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

var sessionInfo = make(map[string]string)

func getUserAgent(r *http.Request) string {
	userAgent := r.UserAgent()
	if userAgent != "" {
		return userAgent
	}
	return r.Header.Get("user-agent")
}

func getHostAddress(r *http.Request) string {
	hostAddress := r.Header.Get("host")
	if hostAddress != "" {
		return hostAddress
	}
	return r.RemoteAddr
}

// Logs the connection information
func logConnection(r *http.Request) {
	if path, exists := sessionInfo[r.Host]; exists {
		if path != r.URL.Path && // streaming content
			r.URL.Path != "/api/renew" && // page refresh
			!strings.HasPrefix(r.URL.Path, "/files") && // duplicate, since path is /files
			!strings.HasPrefix(r.URL.Path, "/thumb") && // thumbnails for images
			!strings.HasPrefix(r.URL.Path, "/big") && // actual images redundant to path
			!strings.HasPrefix(r.URL.Path, "assets") { // internal calls to fetch JS
			log.Printf("%s %s", r.Method, r.URL.Path)
			sessionInfo[r.Host] = r.URL.Path
		}
	} else {
		sessionInfo[r.Host] = r.URL.Path
		logStatment := fmt.Sprintf("Connection received from %s client-host: %s, host-header: %s",
			strings.Split(r.Host, ":")[0], r.Host, getHostAddress(r))
		xFwdHost := r.Header.Get("x-forwarded-host")
		if xFwdHost != "" {
			logStatment += fmt.Sprintf(", x-fwd-host: %s", xFwdHost)
		}
		log.Print(logStatment)
		log.Printf("user-agent: %s", getUserAgent(r))
	}
}
