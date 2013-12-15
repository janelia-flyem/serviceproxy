package main

import (
//        "encoding/json"
	"net/http"
	"fmt"
	"syscall"
	"os"
        "os/signal"
	"strconv"
//	"strings"
)

const (
	defaultAddr = "localhost"
        staticPath = "/static/"
        interfacePath = "/interface/"
)

func InterfaceHandler(w http.ResponseWriter, r *http.Request) {
        // allow resources to be accessed via ajax
        w.Header().Set("Content-Type", "application/raml+yaml")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
        w.Header().Set("Access-Control-Allow-Methods", "GET")        
        http.ServeFile(w, r, "interface/interface.raml")
}

func SourceHandler(w http.ResponseWriter, r *http.Request) {
        // allow resources to be accessed via ajax
        w.Header().Set("Access-Control-Allow-Origin", "*")
        http.ServeFile(w, r, r.URL.Path[1:])
}

func Serve(port int) {
	webAddress := defaultAddr + ":" + strconv.Itoa(port)	

        fmt.Printf("Web server address: %s\n", webAddress)
	fmt.Printf("Running...\n")

	httpserver := &http.Server{Addr: webAddress}
        
        // serve out static json schema and raml (allow access)
        http.HandleFunc(staticPath, SourceHandler)
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




