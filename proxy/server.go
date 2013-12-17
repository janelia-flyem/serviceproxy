package proxy

import (
        "encoding/json"
	"fmt"
        "net"
        "os"
	"net/http"
	"strconv"
	"strings"
)

const (
	servicePath = "/services/"
	nodesPath = "/nodes/"
	execPath = "/exec/"
	interfacePath = "/interface/"
	staticPath = "/static/"
        interfaceFile = "interface/interface.raml"
        interfaceFilePath = "/interface/raw"
)

const ramlHTML = `
<html>
<head>
  <link rel="stylesheet" href="/static/api-console/dist/styles/app.css" type="text/css" />
</head>
<body ng-app="ramlConsoleApp" ng-cloak id="raml-console-unembedded">
  <script src="/static/api-console/dist/scripts/vendor.js"></script>
  <script src="/static/api-console/dist/scripts/app.js"></script>
  <raml-console src="ADDRESS"/> 
</body>
</html>
`

var srcPATH string

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

        if err != nil {
                badRequest(w, "error handling URI")
                return
        }
        if requestType != "get" {
                badRequest(w, "only supports gets")
                return
        }
   
        if len(pathlist) == 1 && pathlist[0] == "raw" {
                w.Header().Set("Content-Type", "application/raml+yaml")
                fmt.Fprintf(w, ramlInterface) 
        } else {
                w.Header().Set("Content-Type", "text/html")
                interfaceHTML := strings.Replace(ramlHTML, "ADDRESS", interfaceFilePath, 1)
                fmt.Fprintf(w, interfaceHTML) 
        } 
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
        
        if len(pathlist) > 2 || len(pathlist) == 1 {
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

                if pathlist[1] == "interface" {
                        if err != nil {
                                badRequest(w, "Service " + pathlist[0] + " not found")
                        } else {
                                // ASSUME interface defined at client
                                addr = "http://" + addr + "/interface/interface.raml" 
                                w.Header().Set("Content-Type", "text/html")
                                interfaceHTML := strings.Replace(ramlHTML, "ADDRESS", addr, 1)
                                fmt.Fprintf(w, interfaceHTML)
                        }
                } else if pathlist[1] == "node" { 
                        w.Header().Set("Content-Type", "application/json")
                        var data  map[string]interface{}
                        if err != nil {
                                data = map[string]interface{}{ "service-location" : nil }
                        } else {
                                addr = "http://" + addr
                                data = map[string]interface{}{ "service-location" : addr }
                        }
                        jsonStr, _ := json.Marshal(data)
                        fmt.Fprintf(w, string(jsonStr))
                } else {
                        badRequest(w, "Bad request for service: " + pathlist[0])
                }
        }
}

func SourceHandler(w http.ResponseWriter, r *http.Request) {
        // allow resources to be accessed via ajax
        w.Header().Set("Access-Control-Allow-Origin", "*")
        http.ServeFile(w, r, srcPATH + r.URL.Path[1:])
}

func Serve(port int, srcroot string) error {
        srcPATH = srcroot + "/src/github.com/janelia-flyem/serviceproxy/"
	
        hname, _ := os.Hostname()
	addrs, _ := net.LookupHost(hname)
        webAddress := addrs[1] + ":" + strconv.Itoa(port)

	fmt.Printf("Web server address: %s\n", webAddress)
	fmt.Printf("Running...\n")

	httpserver := &http.Server{Addr: webAddress}

	http.HandleFunc(servicePath, serviceHandler)
	http.HandleFunc(nodesPath, nodesHandler)
	http.HandleFunc(execPath, execHandler)
	http.HandleFunc(interfacePath, interfaceHandler)
        
        // should be only called internally?
        
        // serve out static javascript pages for api handler
        http.HandleFunc(staticPath, SourceHandler)

	httpserver.ListenAndServe()

	return nil
}
