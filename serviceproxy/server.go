package serviceproxy

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	defaultAddr = "127.0.0.1"
	servicePath = "/services"
)

// Registry contains the services registered with the proxy server
var registry Registry

func badRequest(w http.ResponseWriter, msg string) {
	fmt.Println(msg)
	http.Error(w, msg, http.StatusBadRequest)
}

func serviceHandler(w http.ResponseWriter, r *http.Request) {
	requestType := strings.ToLower(r.Method)

	if requestType != "get" {
		badRequest(w, "only supports gets")
	}

	registry.updateRegistry()
	members := registry.getServicesString()

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, members)
}

func Serve(port int) error {
	webAddress := defaultAddr + ":" + strconv.Itoa(port)

	fmt.Printf("Web server address: %s\n", webAddress)

	httpserver := &http.Server{Addr: webAddress}

	http.HandleFunc("/services", serviceHandler)

	httpserver.ListenAndServe()

	return nil
}
