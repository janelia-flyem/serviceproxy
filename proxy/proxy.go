package proxy

import (
	"fmt"
	"github.com/janelia-flyem/serviceproxy/register"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

// ServiceProxy contains information on proxy server
type ServiceProxy struct {
	Port  int
	Debug bool
}

// Contains the RPC address for the serf agent attached to the proxy server
var rpcAddr string

// init sets up the default RPC address for the serf agent
func init() {
	hname, _ := os.Hostname()
	rpcAddr = hname + ":7373"
}

// Run creates the serf agent and launches the http server
func (proxy *ServiceProxy) Run() error {
	// create agent and launch (no join node is specified)
	serfagent := register.NewAgent("proxy", proxy.Port)
	serfagent.Debug = proxy.Debug
	serfagent.RegisterService("")

	hname, _ := os.Hostname()
	serfaddr := hname + ":" + strconv.Itoa(serfagent.GetSerfPort())

	// address for clients to register (port does not need to be specified
	// by the client if using the Go register interface)
	fmt.Printf("Registry address: %s\n", serfaddr)

	// exit server if user presses Ctrl-C
	go func() {
		sigch := make(chan os.Signal)
		signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
		<-sigch
		fmt.Println("Exiting...")
		os.Exit(0)
	}()

	// create web server
	return Serve(proxy.Port)
}
