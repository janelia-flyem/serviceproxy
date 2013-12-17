package main

import (
	"flag"
	"fmt"
	"github.com/janelia-flyem/serviceproxy/proxy"
	"os"
	"strings"
)

// Arbitrary port chosen as the default
const defaultPort = 15333

var (
	// Prints help message
	showHelp = flag.Bool("help", false, "")

	// Specify port address for web server
	portNum = flag.Int("port", defaultPort, "")

	// Show serf agent debug output
	debugSerf = flag.Bool("debug", false, "")
)

const helpMessage = `
The service proxy tool registers different services through the
serf package. (GOPATH must be set to the current src path.) 
 
Usage: serviceproxy
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

	// GOPATH needed to access static resources.  (To be removed)
	srcroot := os.Getenv("GOPATH")
	if srcroot == "" || strings.Contains(srcroot, ":") {
		fmt.Printf("GOPATH must be set to current src path\n")
		os.Exit(0)
	}

	proxy := proxy.ServiceProxy{Port: *portNum, Debug: *debugSerf}
	// launch tool to create proxy serf agent and web server
	proxy.Run(srcroot)
}
