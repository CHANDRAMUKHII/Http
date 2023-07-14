package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	http.HandleFunc("/", proxyRequest)
	fmt.Print("Listening on port 8001")
	log.Fatal(http.ListenAndServe(":8001", nil))

}

func proxyRequest(w http.ResponseWriter, r *http.Request) {
	targetURL := "localhost:8002"

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   targetURL,
	})

	r.URL.Host = targetURL
	r.URL.Scheme = "http"
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = targetURL

	proxy.ServeHTTP(w, r)
}
