package serviceproxy

import (
        "encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	defaultAddr = "localhost"
	servicePath = "/services/"
)

// Registry contains the services registered with the proxy server
var registry Registry

func badRequest(w http.ResponseWriter, msg string) {
	fmt.Println(msg)
	http.Error(w, msg, http.StatusBadRequest)
}

func parseURI(r *http.Request, prefix string) ([]string, string, error) {
	requestType := strings.ToLower(r.Method)
        prefix = strings.Trim(prefix, "/")
        path := strings.Trim(r.URL.Path, "/")
        prefix_list := strings.Split(prefix, "/")
        url_list := strings.Split(path, "/")
        var path_list []string   

        if len(prefix_list) > len(url_list) {
                return path_list, requestType, fmt.Errorf("Incorrectly formatted URI")
        }

        for i, val := range prefix_list {
                if val != url_list[i] {
                        return path_list, requestType, fmt.Errorf("Incorrectly formatted URI")
                }
        }

        if len(prefix_list) < len(url_list) {
                path_list = url_list[len(prefix_list):]
        }

        return path_list, requestType, nil 
}

func serviceHandler(w http.ResponseWriter, r *http.Request) {

        pathlist, requestType, err := parseURI(r, servicePath)
        
        if err != nil {
                badRequest(w, "error handling URI")
                return
        }
	if requestType != "get" {
		badRequest(w, "only supports gets")
                return
	}
        
        if len(pathlist) > 1 {
                badRequest(w, "incorrectly formatted URL")
                return
        }

        if len(pathlist) == 0 {
                registry.updateRegistry()
                members := registry.getServicesSlice()

                w.Header().Set("Content-Type", "text/html")

                for _, service := range members {
                        href := "<a href=\"" + servicePath + service + "\">"
                        fmt.Fprintf(w, href + service + "</a>")
                        fmt.Fprintf(w, "<br>") 
                }
        } else {
                registry.updateRegistry()
                addr, err := registry.getServiceAddr(pathlist[0])
    
                w.Header().Set("Content-Type", "application/json")
                var data  map[string]string
                if err != nil {
                        data = map[string]string{ pathlist[0] : "null" }
                } else {
                        data = map[string]string{ pathlist[0] : addr }
                }
                jsonStr, _ := json.Marshal(data)
                fmt.Fprintf(w, string(jsonStr))
        }
}

func Serve(port int) error {
	webAddress := defaultAddr + ":" + strconv.Itoa(port)

	fmt.Printf("Web server address: %s\n", webAddress)
	fmt.Printf("Running...\n")

	httpserver := &http.Server{Addr: webAddress}

	http.HandleFunc(servicePath, serviceHandler)

	httpserver.ListenAndServe()

	return nil
}
