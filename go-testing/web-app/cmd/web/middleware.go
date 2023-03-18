package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

type contextKey string

const contextUserKey contextKey = "user_ip"

func (app *application) ipFromContext(ctx context.Context) string {
	return ctx.Value(contextUserKey).(string)
}

func (app *application) addIPToContext(next http.Handler) http.Handler {
	//return inline function
	return http.HandlerFunc(func(write http.ResponseWriter, read *http.Request) {
		var ctx = context.Background()
		// get the ip (as accurately as possible)
		ip, err := getIP(read)
		if err != nil {
			ip, _, _ = net.SplitHostPort(read.RemoteAddr)
			if len(ip) == 0 {
				ip = "unknown"
			}
			ctx = context.WithValue(read.Context(), contextUserKey, ip)
		} else {
			ctx = context.WithValue(read.Context(), contextUserKey, ip)
		}
		next.ServeHTTP(write, read.WithContext(ctx))
	})
}

func getIP(read *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(read.RemoteAddr)
	if err != nil {
		return "unknown", err
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return "", fmt.Errorf("userip: %q is not IP:port", read.RemoteAddr)
	}

	forward := read.Header.Get("X-Forwarded-For")
	if len(forward) > 0 {
		ip = forward
	}

	if len(ip) == 0 {
		ip = "forward"
	}

	return ip, nil
}
