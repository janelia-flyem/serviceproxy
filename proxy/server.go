package proxy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Contains paths used by the server
const (
	servicePath       = "/services/"
	nodesPath         = "/nodes/"
	execPath          = "/exec/"
	interfacePath     = "/interface/"
	staticPath        = "/static/"
	interfaceFile     = "interface/interface.raml"
	interfaceFilePath = "/interface/raw"
)

// ramlHTML is the interface for the proxy server
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

// srcPATH is the location of the static files from GOPATH.
// TBD: put static files in memory to avoid this hack
var srcPATH string

// Registry contains the services registered with the proxy server
var registry Registry

// init sets the srcPATH variable to the GOLANG
func init() {
	// GOPATH needed to access static resources.  (To be removed)
	srcroot := os.Getenv("GOPATH")
	if srcroot == "" || strings.Contains(srcroot, ":") {
		fmt.Printf("GOPATH must be set to current src path\n")
		os.Exit(0)
	}

	srcPATH = srcroot + "/src/github.com/janelia-flyem/serviceproxy/"
}

// badRequest is a halper for printing an http error message
func badRequest(w http.ResponseWriter, msg string) {
	fmt.Println(msg)
	http.Error(w, msg, http.StatusBadRequest)
}

// parseURI is a utility function for retrieving parts of the URI
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

// interfaceHandler returns the raml interface
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

	// return the raw file
	if len(pathlist) == 1 && pathlist[0] == "raw" {
		w.Header().Set("Content-Type", "application/raml+yaml")
		fmt.Fprintf(w, ramlInterface)
	} else { // return a javascript view of the interface
		w.Header().Set("Content-Type", "text/html")
		interfaceHTML := strings.Replace(ramlHTML, "ADDRESS", interfaceFilePath, 1)
		fmt.Fprintf(w, interfaceHTML)
	}
}

// execHandler redirects call to the appropriate service
func execHandler(w http.ResponseWriter, r *http.Request) {
	pathlist, _, err := parseURI(r, execPath)

	if err != nil || len(pathlist) == 0 {
		badRequest(w, "error handling URI")
		return
	}

	// determine node that implements service
	registry.UpdateRegistry()
	addr, err := registry.GetServiceAddr(pathlist[0])

	if err != nil {
		badRequest(w, "error in processing: "+pathlist[0])
	} else {
		url := "http://" + addr + "/"
		if len(pathlist) > 1 {
			url += strings.Join(pathlist[1:], "/")
		}

		http.Redirect(w, r, url, http.StatusFound)
	}
}

// nodesHandler returns a list of nodes that are visible to the proxy server
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

	registry.UpdateRegistry()
	nodes := registry.GetActiveNodes()

	data := make(map[string]interface{})
	data["nodes"] = nodes
	jsonStr, _ := json.Marshal(data)
	fmt.Fprintf(w, string(jsonStr))
}

// serviceHandler retrieves information about the services
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

	// return member services
	registry.UpdateRegistry()
	if len(pathlist) == 0 {
		members := registry.GetServicesSlice()

		w.Header().Set("Content-Type", "application/json")

		data := make(map[string]interface{})
		var services_json []interface{}

		for _, service := range members {
			service_map := map[string]string{service: servicePath + service}
			services_json = append(services_json, service_map)
		}
		data["services"] = services_json
		jsonStr, _ := json.Marshal(data)
		fmt.Fprintf(w, string(jsonStr))
	} else {
		addr, err := registry.GetServiceAddr(pathlist[0])

		// retrieve service interface
		if pathlist[1] == "interface" {
			if err != nil {
				badRequest(w, "Service "+pathlist[0]+" not found")
			} else {
				// ASSUME interface defined at client
				addr = "http://" + addr + "/interface/interface.raml"
				w.Header().Set("Content-Type", "text/html")
				interfaceHTML := strings.Replace(ramlHTML, "ADDRESS", addr, 1)
				fmt.Fprintf(w, interfaceHTML)
			}
		} else if pathlist[1] == "node" { // retrieve random node implementing interface
			w.Header().Set("Content-Type", "application/json")
			var data map[string]interface{}
			if err != nil {
				data = map[string]interface{}{"service-location": nil}
			} else {
				addr = "http://" + addr
				data = map[string]interface{}{"service-location": addr}
			}
			jsonStr, _ := json.Marshal(data)
			fmt.Fprintf(w, string(jsonStr))
		} else {
			badRequest(w, "Bad request for service: "+pathlist[0])
		}
	}
}

// SourceHandler serves static files
func SourceHandler(w http.ResponseWriter, r *http.Request) {
	// allow resources to be accessed via ajax
	w.Header().Set("Access-Control-Allow-Origin", "*")
	http.ServeFile(w, r, srcPATH+r.URL.Path[1:])
}

// serveHTML loads an app that shows all the active services
func serveHTML(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, htmldata) 
}

// Serve creates http server and sets handlers
func Serve(port int) error {
        hname, _ := os.Hostname()
	webAddress := hname + ":" + strconv.Itoa(port)

	fmt.Printf("Web server address: %s\n", webAddress)
	fmt.Printf("Running...\n")

	httpserver := &http.Server{Addr: webAddress}

	http.HandleFunc(servicePath, serviceHandler)
	http.HandleFunc(nodesPath, nodesHandler)
	http.HandleFunc(execPath, execHandler)
	http.HandleFunc(interfacePath, interfaceHandler)

        http.HandleFunc("/", serveHTML)

	// should be only called internally?

	// serve out static javascript pages for api handler
	http.HandleFunc(staticPath, SourceHandler)

	httpserver.ListenAndServe()

	return nil
}
