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
		// following condition prevents long videos from spamming the logs
		if path != r.URL.Path {
			log.Printf("%s %s", strings.ToUpper(r.Method), path)
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
