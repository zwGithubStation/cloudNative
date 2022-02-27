package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"
)

func index(w http.ResponseWriter, r *http.Request) {
	os.Setenv("VERSION", "1.0.0")
	version := os.Getenv("VERSION")
	w.Header().Set("VERSION", version)

	for k, v := range r.Header {
		for _, vv := range v {
			w.Header().Set(k, vv)
		}
	}

	clientip := getCurrentIP(r)
	log.Printf("Success! Response code: %d", 200)   
	log.Printf("Success! clientip: %s", clientip)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "working")
}

func getCurrentIP(r *http.Request) string {     
	ip := r.Header.Get("X-Real-IP")   
	if ip == "" {     
		ip = strings.Split(r.RemoteAddr, ":")[0]   
	}   
	return ip
}

func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" { 
		return ip   
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" { 
		return ip   
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}


func main() {
	mux := http.NewServeMux()   
	mux.HandleFunc("/debug/pprof/", pprof.Index)   
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)   
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)   
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)   
	mux.HandleFunc("/", index)   
	mux.HandleFunc("/healthz", healthz)   
	if err := http.ListenAndServe(":8080", mux); err != nil {      
		log.Fatalf("start http server failed, error: %s\n", err.Error())   
	}
}

