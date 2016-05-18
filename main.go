package main

import (
	"log"
	"net/http"
	"net"
	"net/http/httputil"
	"time"
)

func main() {
	http.HandleFunc("/", ReverseProxy())
	println("ready")
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func ReverseProxy() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		transport := &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: func(network, addr string) (net.Conn, error) {
				return getConnection(req)
			},
			TLSHandshakeTimeout: 10 * time.Second,
		}
		(&httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = "http"
				req.URL.Host = "/latest" + req.RequestURI
			},
			Transport: transport,
		}).ServeHTTP(w, req)
	}
}

func getConnection(req *http.Request) (net.Conn, error) {
	return net.Dial("tcp", "rancher-metadata:80")
}
