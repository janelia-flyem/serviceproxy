package register

import (
	"github.com/hashicorp/serf/command/agent"
	"github.com/mitchellh/cli"
	"github.com/hashicorp/serf/command"
        "fmt"
        "os"
        "net"
        "strconv"
        "strings"
)

const startingPort = 25001
const buffersize = 10000

type AgentWriter struct {
	bytes []byte
        debug bool
}

func (w *AgentWriter) Write(p []byte) (n int, err error) {
	if w.debug {
            fmt.Println(string(p))
        }
        
        if len(w.bytes) > buffersize {
                w.bytes = make([]byte, 0)
        }

        n = len(p)
	w.bytes = append(w.bytes, p...)
	err = nil

	return n, err
}

func (w *AgentWriter) GetString() string {
	return string(w.bytes[:])
}

type SerfAgent struct {
        name      string
        port      int
        serfname  string
        haddr     string
        
        SerfPort  int
        RPCPort   int
        Debug     bool
}

func NewAgent(name string, port int) (*SerfAgent) {
        agent := &SerfAgent{name: name, port: port, SerfPort: -1, 
                        RPCPort: -1, Debug: false}       
        hname, _ := os.Hostname()
        addrs, _ := net.LookupHost(hname)
        agent.haddr = addrs[1]
        agent.serfname = name+"#"+addrs[1]+":"+strconv.Itoa(port)

        return agent
}

func (s *SerfAgent) RegisterService(registry string) (error) {
        // place serf output in a debug buffer
	writer := &AgentWriter{debug: s.Debug}

        // create agent command
        ui := &cli.BasicUi{Writer: writer}
        ac := &agent.Command{Ui: ui, ShutdownCh: make(chan struct{})}
	var dargs []string

        // format proper node name (service + address + port)	
        dargs = append(dargs, "-node=" + s.serfname)

        if registry != "" {
                dargs = append(dargs, "-join=" + registry)
        }

	go s.launchAgent(ac, dargs, writer)

        return nil
}

func (s *SerfAgent) UnRegisterService() (error) {
	writer := &AgentWriter{debug: s.Debug}
        
        // create agent command
        ui := &cli.BasicUi{Writer: writer}
        ac := &command.LeaveCommand{Ui: ui}

        // add the kill address and run	
	var dargs []string
        dargs = append(dargs, "-rpc-addr=" + s.haddr + ":" + strconv.Itoa(s.RPCPort))
        ac.Run(dargs)
        
        return nil
}


func (s *SerfAgent) launchAgent(ac *agent.Command, dargs []string, writer *AgentWriter) {
        hasdefault := false
        if s.SerfPort != -1 || s.RPCPort != -1 {
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
                    s.SerfPort = start
                    s.RPCPort = start + 1
                    start += 2 
                }
        
                // add port options
                serfaddr := s.haddr + ":" + strconv.Itoa(s.SerfPort) 
                rpcaddr := s.haddr + ":" + strconv.Itoa(s.RPCPort) 
                args = append(args, "-bind="+serfaddr)
                args = append(args, "-rpc-addr="+rpcaddr)
                
                // launch blocking call to serf agent 
                ac.Run(args)

                // check final output for errors
                output := writer.GetString()
                fmt.Println(output)
                if strings.Contains(output, "Failed to start") || strings.Contains(output, "Error starting") {
                        if hasdefault {
                                panic("Agent failed to start at specified ports")
                        }
                } else {
                        break
                }
        }
}


