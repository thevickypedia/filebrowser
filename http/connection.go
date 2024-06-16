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
	if userAgent == "" {
		userAgent = r.Header.Get("user-agent")
	}
	return userAgent
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
		logStatment := fmt.Sprintf("Connection received from client-host: %s", r.Host)
		hostHeader := r.Header.Get("host")
		xFwdHost := r.Header.Get("x-forwarded-host")
		if hostHeader != "" {
			logStatment += fmt.Sprintf(", host-header: %s", hostHeader)
		}
		if xFwdHost != "" {
			logStatment += fmt.Sprintf(", x-fwd-host: %s", xFwdHost)
		}
		log.Print(logStatment)
		userAgent := getUserAgent(r)
		if userAgent != "" {
			log.Printf("user-agent: %s", userAgent)
		}
	}
}
