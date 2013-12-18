package register

import (
	"fmt"
	"github.com/hashicorp/serf/command"
	"github.com/hashicorp/serf/command/agent"
	"github.com/mitchellh/cli"
	"os"
	"strconv"
	"strings"
)

const (
	// First port checked for availability
        startingPort = 25001

        // Size of the write buffer used with serf
	buffersize   = 100000

	// Default serf port for proxy service
	serfPort  = 7946

        // Default serf rpc port for proxy service
	rpcPort   = 7373

        // Name of the proxy service
	proxyName = "proxy"
)

// AgentWriter implements the Writer interface and contains serf output
type AgentWriter struct {
	// Contains the last several bytes of serf output
        bytes []byte

        // If true it will dump output to console
	debug bool
}

// Write implements Writer interface and dumps output to buffer
func (w *AgentWriter) Write(p []byte) (n int, err error) {
	if w.debug {
		fmt.Println(string(p))
	}

        // clears buffer when buffer size is reached
        // TODO: make buffer a FIFO queue instead
	if len(w.bytes) > buffersize {
		w.bytes = make([]byte, 0)
	}

	n = len(p)
	w.bytes = append(w.bytes, p...)
	err = nil

	return n, err
}

// GetString returns the string value of the written bytes
func (w *AgentWriter) GetString() string {
	return string(w.bytes[:])
}

// SerfAgent contains information neeeded to call serf facilities
type SerfAgent struct {
	name     string
	port     int
	serfname string
	haddr    string
	serfPort int
	rpcPort  int

        // Print debug output
	Debug    bool

        // False by default, if true, a new process is not created for the agent
	Blocking bool
}

// NewAgent creates an agent using specific naming conventions
func NewAgent(name string, port int) *SerfAgent {
	agent := &SerfAgent{name: name, port: port, serfPort: -1,
		rpcPort: -1, Debug: false, Blocking: false}
	hname, _ := os.Hostname()
	agent.haddr = hname
	agent.serfname = name + "#" + hname + ":" + strconv.Itoa(port)

        // specialized check for proxy service registration: use default ports
	if name == proxyName {
		agent.serfPort = serfPort
		agent.rpcPort = rpcPort
	}

	return agent
}

// GetSerfPort returns the port number that the agent is (will be) registered at
func (s *SerfAgent) GetSerfPort() int {
	return s.serfPort
}

// RegisterService calls serf and launches an agent.
// It join the serf network at the node specified by the registry.  If
// registry is an empty string, the agent does not join other nodes.
func (s *SerfAgent) RegisterService(registry string) error {
	// place serf output in a debug buffer
	writer := &AgentWriter{debug: s.Debug}

	// create agent command
	ui := &cli.BasicUi{Writer: writer}
	ac := &agent.Command{Ui: ui, ShutdownCh: make(chan struct{})}
	var dargs []string

	// format proper node name (service + address + port)
	dargs = append(dargs, "-node="+s.serfname)

	if registry != "" {
		dargs = append(dargs, "-join="+registry)
	}

	if !s.Blocking {
		go s.launchAgent(ac, dargs, writer)
	} else {
		s.launchAgent(ac, dargs, writer)
	}

	return nil
}

// UnRegisterService stops a previously launches serf agent.
// This function is a blocking call and is generally only callable
// if SerfAgent was initially registered.
func (s *SerfAgent) UnRegisterService() error {
	writer := &AgentWriter{debug: s.Debug}

	// create agent command
	ui := &cli.BasicUi{Writer: writer}
	ac := &command.LeaveCommand{Ui: ui}

	// add the kill address and run
	var dargs []string

	// use current rpc port or find rpc port
	if s.rpcPort == -1 {
		return fmt.Errorf("Cannot unregister without an RPC address\n")
	}
	dargs = append(dargs, "-rpc-addr="+s.haddr+":"+strconv.Itoa(s.rpcPort))
	ac.Run(dargs)

        // re-initialize ports if not a proxy service	
        if s.name != proxyName {
                s.serfPort = -1
                s.rpcPort = -1
        }

	return nil
}

// launchAgent actually calls the serf interface.
// By default, it is called as a go process.  Except for a proxy service,
// it will keep launching the agent until a valid port combination is tried.
func (s *SerfAgent) launchAgent(ac *agent.Command, dargs []string, writer *AgentWriter) {
	hasdefault := false
        
        // should always be true if it is not a proxy
	if s.serfPort != -1 || s.rpcPort != -1 {
		hasdefault = true
	}
	start := startingPort

	for {
		var args []string
		for _, darg := range dargs {
			args = append(args, darg)
		}

		// pick 2 ports for the serf agent
		if !hasdefault {
			s.serfPort = start
			s.rpcPort = start + 1
			start += 2
		}

		// add port options
		serfaddr := s.haddr + ":" + strconv.Itoa(s.serfPort)
		rpcaddr := s.haddr + ":" + strconv.Itoa(s.rpcPort)
		args = append(args, "-bind="+serfaddr)
		args = append(args, "-rpc-addr="+rpcaddr)

		// launch blocking call to serf agent
		ac.Run(args)

		// check final output for errors
		output := writer.GetString()
		if strings.Contains(output, "Failed to start") || strings.Contains(output, "Error starting") {
			if hasdefault {
				panic("Agent failed to start at specified ports")
			}
		} else {
			break
		}
	}
}
