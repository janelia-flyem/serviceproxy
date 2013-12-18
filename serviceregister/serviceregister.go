package main

import (
	"flag"
	"fmt"
	"github.com/janelia-flyem/serviceproxy/register"
	"os"
	"strconv"
)

var (
	showHelp = flag.Bool("help", false, "")
)

const helpMessage = `
Register a service on a given port to the registry.

Usage: serviceregister <service name> <service port> <registry address>
  -h, -help     (flag)          Show help message
`

func main() {
	flag.BoolVar(showHelp, "h", false, "Show help message")
	flag.Parse()

	if *showHelp {
		fmt.Printf(helpMessage)
		os.Exit(0)
	}

	if flag.NArg() != 3 {
		fmt.Printf("Need 3 arguments")
		fmt.Printf(helpMessage)
		os.Exit(0)
	}
	port, err := strconv.Atoi(flag.Arg(1))

	if err != nil {
		fmt.Printf("Not a valid port")
		os.Exit(0)
	}

	// creates new serf agent with a node name comprised of the service name and port
	serfagent := register.NewAgent(flag.Arg(0), port)

	// service call should be blocking -- will run until killed
	serfagent.Blocking = true

	// service registraters to service registry address
	serfagent.RegisterService(flag.Arg(2))
}
