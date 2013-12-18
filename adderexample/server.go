package main

import (
	"encoding/json"
	"fmt"
	"github.com/sigu-399/gojsonschema"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	// Contain URI location for interface
	interfacePath = "/interface/"
)

// String representing interface for adder example
const ramlInterface = `#%%RAML 0.8
title: Adder Service
/:
  post:
    description: "Call service to add two numbers"
    body:
      application/json:
        schema: |
          { "$schema": "http://json-schema.org/schema#",
            "title": "Provide numbers to be added together",
            "type": "object",
            "properties": {
              "num1" : { "type" : "integer" },
              "num2" : { "type" : "integer" }
            },
            "required" : ["num1", "num2"]
          }
    responses:
      200:
        body:
          application/json:
            schema: |
              { "$schema" : "http://json-schema.org/schema#",
                "title" : "Provides callback link for results",
                "type" : "object",
                "properties" : {
                  "result-callback" : {
                    "description" : "URL for results",
                    "type" : "string"
                  }
                },
                "required" : ["result-callback"]
              }
/jobs/{id}:
  get:
    description: "Get the result from a particular job"
    responses:
      200:
        body:
          application/json:
            schema: |
              { "$schema" : "http://json-schema.org/schema#",
                "title" : "Results from adder service",
                "type" : "object",
                "properties" : {
                  "result" : {
                    "type" : "integer"
                  }
                },
                "required" : ["result"]
              }         
/interface/interface.raml:
  get:
    description: "Get the interface for the adder service"
    responses:
      200:
        body:
          application/raml+yaml:
`

// String representing the JSON schema for the service call
const serviceSchema = `
{ "$schema": "http://json-schema.org/schema#",
  "title": "Provide numbers to be added together",
  "type": "object",
  "properties": {
    "num1" : { "type" : "integer" },
    "num2" : { "type" : "integer" }
  },
  "required" : ["num1", "num2"]
}
`

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

// webAddress is the http address for the server
var webAddress string

// badRequest is a halper for printing an http error message
func badRequest(w http.ResponseWriter, msg string) {
	fmt.Println(msg)
	http.Error(w, msg, http.StatusBadRequest)
}

// InterfaceHandler returns the RAML interface for any request at
// the /interface URI.
func interfaceHandler(w http.ResponseWriter, r *http.Request) {
	// allow resources to be accessed via ajax
	w.Header().Set("Content-Type", "application/raml+yaml")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, ramlInterface)
}

// AddRequest is the datastructure to be populated with the JSON input
type AddRequest struct {
	num1 int
	num2 int
	name string
}

// JobResults contains all the results from service requests in memory
type JobResults struct {
	Results map[string]interface{}
}

var jobResults JobResults

// randomHex computes a random hash for storing service results
func randomHex() (randomStr string) {
	randomStr = ""
	for i := 0; i < 8; i++ {
		val := rand.Intn(16)
		randomStr += strconv.FormatInt(int64(val), 16)
	}
	return
}

// addService is the actual 'service' for computing addition
func addService(addRequest AddRequest) {
	// wait 10 seconds till result is computed
	time.Sleep(10 * time.Second)
	result := addRequest.num1 + addRequest.num2
	jobResults.Results[addRequest.name] = result
}

// serviceHandler handlers post request to "/"
func serviceHandler(w http.ResponseWriter, r *http.Request) {
	pathlist, requestType, err := parseURI(r, "/")
	if err != nil || len(pathlist) != 0 {
		badRequest(w, "Error: incorrectly formatted request")
		return
	}
	if requestType != "post" {
		badRequest(w, "only supports posts")
		return
	}

	// read json
	decoder := json.NewDecoder(r.Body)
	var json_data map[string]interface{}
	err = decoder.Decode(&json_data)

	// convert schema to json data
	var schema_data interface{}
	json.Unmarshal([]byte(serviceSchema), &schema_data)

	// validate json schema
	schema, err := gojsonschema.NewJsonSchemaDocument(schema_data)
	validationResult := schema.Validate(json_data)
	if !validationResult.IsValid() {
		badRequest(w, "JSON did not pass validation")
		return
	}

	var addRequest AddRequest
	addRequest.num1 = int(json_data["num1"].(float64))
	addRequest.num2 = int(json_data["num2"].(float64))
	if err != nil {
		badRequest(w, "JSON not formatted properly")
		return
	}

	jobid := randomHex()
	var empty interface{}
	jobResults.Results[jobid] = empty
	addRequest.name = jobid

	w.Header().Set("Content-Type", "application/json")
	jsondata, _ := json.Marshal(map[string]string{
		"result-callback": "http://" + webAddress + "/jobs/" + jobid,
	})
	fmt.Fprintf(w, string(jsondata))

	// non-blocking call to actual service
	go addService(addRequest)
}

// jobHandler gets requests for job status
func jobHandler(w http.ResponseWriter, r *http.Request) {
	pathlist, requestType, err := parseURI(r, "/jobs/")
	if err != nil || len(pathlist) != 1 {
		badRequest(w, "Error: incorrectly formatted request")
		return
	}
	if requestType != "get" {
		badRequest(w, "only supports gets")
		return
	}
	result, ok := jobResults.Results[(pathlist[0])]

	if !ok {
		badRequest(w, "job does not exist")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsondata, _ := json.Marshal(map[string]interface{}{
		"result": result,
	})
	fmt.Fprintf(w, string(jsondata))
}

// Serve is the main server function call that creates http server and handlers
func Serve(port int) {
	jobResults.Results = make(map[string]interface{})

	hname, _ := os.Hostname()
	webAddress = hname + ":" + strconv.Itoa(port)

	fmt.Printf("Web server address: %s\n", webAddress)
	fmt.Printf("Running...\n")

	httpserver := &http.Server{Addr: webAddress}

	// serve out static json schema and raml (allow access)
	http.HandleFunc(interfacePath, interfaceHandler)

	// serve out static json schema and raml (allow access)
	http.HandleFunc("/", serviceHandler)

	http.HandleFunc("/jobs/", jobHandler)

	// exit server if user presses Ctrl-C
	go func() {
		sigch := make(chan os.Signal)
		signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
		<-sigch
		fmt.Println("Exiting...")
		os.Exit(0)
	}()

	httpserver.ListenAndServe()
}
