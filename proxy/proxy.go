package proxy

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
        "github.com/janelia-flyem/serviceproxy/register"
        "strconv"
)

type ServiceProxy struct {
	Port  int
	Debug bool
}

func (proxy *ServiceProxy) Run(srcroot string) error {
        // create agent and launch (no join node is specified)
        serfagent := register.NewAgent("proxy", proxy.Port)
        serfagent.Debug = proxy.Debug
        serfagent.RegisterService("")

	hname, _ := os.Hostname()
	addrs, _ := net.LookupHost(hname)
	serfaddr := addrs[1] + ":" + strconv.Itoa(serfagent.GetSerfPort())

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
	return Serve(proxy.Port, srcroot)
}
