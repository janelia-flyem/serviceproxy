package main

import (
        "flag"
        "fmt"
        "github.com/janelia-flyem/serviceproxy/register"
        "os"
        "strings"
)

const defaultPort = 23230

var (
	portNum   = flag.Int("port", defaultPort, "")
	showHelp  = flag.Bool("help", false, "")
)

const helpMessage = `
Launches example service that adds two numbers.

Usage: adderexample <registry address>
      -port     (number)        Port for HTTP server
  -h, -help     (flag)          Show help message
` 

func main() {
	flag.BoolVar(showHelp, "h", false, "Show help message")
	flag.Parse()

	if *showHelp {
		fmt.Printf(helpMessage)
		os.Exit(0)
	}

        // register service
        if flag.NArg() != 1 {
                fmt.Printf("Must provide registry address")
                fmt.Printf(helpMessage)
                os.Exit(0)
        }   

        srcroot := os.Getenv("GOPATH")
        if srcroot == "" || strings.Contains(srcroot, ":") {
                fmt.Printf("GOPATH must be set to current src path\n")
                os.Exit(0)
        }

        serfagent := register.NewAgent("adder", *portNum)
        serfagent.RegisterService(flag.Arg(0)) 
 
        Serve(*portNum, srcroot)
}
