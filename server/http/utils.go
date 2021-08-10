package httpServer

import (
	"buildings_info/logging"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"strings"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func serveOrigin(allowed []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err == nil {
				return
			}
			stack := debug.Stack()
			_, _ = os.Stderr.Write(stack)
		}()

		origin := r.Header.Get("Origin")
		if len(origin) > 0 {
			if !isOriginAllowed(origin, allowed) {
				logging.Logging.Error("origin %v is not allowed", origin)
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("origin not allowed"))
				return
			}
			header := w.Header()
			header.Add(acao, origin)
			header.Add(acam, allowMethods)
			header.Add(acah, allowHeaders)
			header.Add(acac, allowCredentials)
			header.Add(acma, maxAge)
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func isOriginAllowed(origin string, allowed []string) bool {
	if len(allowed) == 0 {
		return true
	}
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	h, _, _ := net.SplitHostPort(u.Host)
	if h == "" {
		h = u.Host
	}
	h = strings.ToLower(h)
	for _, it := range allowed {
		if h == it || strings.HasSuffix(h, "."+it) || it == "*" {
			return true
		}
	}
	return false
}

const (
	acao = "Access-Control-Allow-Origin"
	acam = "Access-Control-Allow-Methods"
	acah = "Access-Control-Allow-Headers"
	acac = "Access-Control-Allow-Credentials"
	acma = "Access-Control-Max-Age"

	allowHeaders     = "Origin, Authorization, Accept, Accept-Encoding, Cache-Control, Content-Type, Content-Length, X-SupplierId, X-UUID"
	allowMethods     = "GET, POST, PUT, DELETE, OPTIONS"
	allowCredentials = "true"
	maxAge           = "3600"
)
