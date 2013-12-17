package main

import (
	"github.com/janelia-flyem/serviceproxy/proxy"
	"github.com/janelia-flyem/serviceproxy/register"
	"os"
	"testing"
	"time"
)

func TestProxyRegister(t *testing.T) {
	// register proxy
	serfagent := register.NewAgent("proxy", 15333)
	err := serfagent.RegisterService("")

	// make sure service is up
	time.Sleep(1 * time.Second)

	if err != nil {
		t.Errorf("Error when registering")
		return
	}

	port := serfagent.GetSerfPort()
	if port != 7946 {
		t.Errorf("Proxy not registered at port 7946")
		return
	}

	// ?! check that service is alive

	err = serfagent.UnRegisterService()
	if err != nil {
		t.Errorf("Failed to unregister proxy")
		return
	}
}

func TestServiceRegister(t *testing.T) {
	// register proxy
	serfagent := register.NewAgent("adder", 23230)
	err := serfagent.RegisterService("")

	// make sure service is up
	time.Sleep(1 * time.Second)

	if err != nil {
		t.Errorf("Error when registering")
		return
	}

	port := serfagent.GetSerfPort()
	if port <= 25000 {
		t.Errorf("Service port not found in correct range")
		return
	}

	// ?! check that service is alive

	err = serfagent.UnRegisterService()
	if err != nil {
		t.Errorf("Failed to unregister service")
		return
	}
}

func TestServiceMemberIdentification(t *testing.T) {
	// register proxy
	serfagent := register.NewAgent("proxy", 15333)
	serfagent.RegisterService("")

	// make sure service is up
	time.Sleep(1 * time.Second)

	// register dummy adder service
	serfagent2 := register.NewAgent("adder", 23230)
	hname, _ := os.Hostname()
	serfagent2.RegisterService(hname + ":7946")

	// make sure service is up
	time.Sleep(10 * time.Second)

	var registry proxy.Registry
	registry.UpdateRegistry()
	members := registry.GetServicesSlice()

	if len(members) != 1 {
		t.Errorf("Number service members %d (should be 1)", len(members))
		return
	}

	if members[0] != "adder" {
		t.Errorf("adder should be the service member")
		return
	}
}
