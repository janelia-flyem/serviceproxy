package main

import (
//        "encoding/json"
	"net"
	"net/http"
	"fmt"
	"syscall"
	"os"
        "os/signal"
	"strconv"
//	"strings"
)

const (
        staticPath = "/static/"
        interfacePath = "/interface/"
)

var srcPATH string

func InterfaceHandler(w http.ResponseWriter, r *http.Request) {
        // allow resources to be accessed via ajax
        w.Header().Set("Content-Type", "application/raml+yaml")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        http.ServeFile(w, r, srcPATH + "interface/interface.raml")
}

func Serve(port int, srcroot string) {
        srcPATH = srcroot + "/src/github.com/janelia-flyem/serviceproxy/adderexample/"

	hname, _ := os.Hostname()
	addrs, _ := net.LookupHost(hname)
        
	webAddress := addrs[1] + ":" + strconv.Itoa(port)	
        
        fmt.Printf("Web server address: %s\n", webAddress)
	fmt.Printf("Running...\n")

	httpserver := &http.Server{Addr: webAddress}
        
        // serve out static json schema and raml (allow access)
        http.HandleFunc(interfacePath, InterfaceHandler)

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




