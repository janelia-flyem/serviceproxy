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
	nodesPath = "/nodes/"
	execPath = "/exec/"
	interfacePath = "/interface/"
        interfaceFile = "interface/interface.raml"
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

func interfaceHandler(w http.ResponseWriter, r *http.Request) {
        pathlist, requestType, err := parseURI(r, interfacePath)

        if err != nil || len(pathlist) != 0 {
                badRequest(w, "error handling URI")
                return
        }
        if requestType != "get" {
                badRequest(w, "only supports gets")
                return
        }
    
        w.Header().Set("Content-Type", "application/raml+yaml")
        http.ServeFile(w, r, interfaceFile)
}

func execHandler(w http.ResponseWriter, r *http.Request) {
        pathlist, _, err := parseURI(r, execPath)

        if err != nil || len(pathlist) == 0 {
                badRequest(w, "error handling URI")
                return
        }

        registry.updateRegistry()
        addr, err := registry.getServiceAddr(pathlist[0])

        if err != nil {
                badRequest(w, "error in processing: " + pathlist[0])
        } else {
                url := "http://" + addr + "/"
                if len(pathlist) > 1 {
                        url += strings.Join(pathlist[1:], "/")
                }

                http.Redirect(w, r, url, http.StatusFound)
        }
}
   


func nodesHandler(w http.ResponseWriter, r *http.Request) {
        pathlist, requestType, err := parseURI(r, nodesPath)

        if err != nil || len(pathlist) != 0 {
                badRequest(w, "error handling URI")
                return
        }
        if requestType != "get" {
                badRequest(w, "only supports gets")
                return
        }
    
        w.Header().Set("Content-Type", "application/json")
        
        registry.updateRegistry()
        nodes := registry.getActiveNodes()
        
        data := make(map[string]interface{})
        data["nodes"] = nodes
        jsonStr, _ := json.Marshal(data)
        fmt.Fprintf(w, string(jsonStr))
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

        registry.updateRegistry()
        if len(pathlist) == 0 {
                members := registry.getServicesSlice()

                w.Header().Set("Content-Type", "application/json")

                data := make(map[string]interface{})
                var services_json []interface{}
                
                for _, service := range members {
                        service_map := map[string]string{ service : servicePath + service}
                        services_json = append(services_json, service_map)
                }
                data["services"] = services_json
                jsonStr, _ := json.Marshal(data)
                fmt.Fprintf(w, string(jsonStr))
        } else {
                addr, err := registry.getServiceAddr(pathlist[0])
    
                w.Header().Set("Content-Type", "application/json")
                var data  map[string]interface{}
                if err != nil {
                        data = map[string]interface{}{ pathlist[0] : nil }
                } else {
                        data = map[string]interface{}{ pathlist[0] : addr }
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
	http.HandleFunc(nodesPath, nodesHandler)
	http.HandleFunc(execPath, execHandler)
	http.HandleFunc(interfacePath, interfaceHandler)

	httpserver.ListenAndServe()

	return nil
}
