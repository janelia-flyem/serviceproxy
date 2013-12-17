package register

import (
	"github.com/hashicorp/serf/command/agent"
	"github.com/mitchellh/cli"
	"github.com/hashicorp/serf/command"
        "fmt"
        "os"
        "strconv"
        "strings"
)

const (

    startingPort = 25001
    buffersize = 10000
   
    // ports for proxy service
    serfPort = 7946
    rpcPort = 7373
    proxyName = "proxy"
)

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
        serfPort  int
        rpcPort   int
        
        Debug     bool
        Blocking  bool
}

func NewAgent(name string, port int) (*SerfAgent) {
        agent := &SerfAgent{name: name, port: port, serfPort: -1, 
                        rpcPort: -1, Debug: false, Blocking: false}       
        hname, _ := os.Hostname()
        agent.haddr = hname
        agent.serfname = name+"#"+hname+":"+strconv.Itoa(port)

        if name == proxyName {
            agent.serfPort = serfPort 
            agent.rpcPort = rpcPort 
        } 
        
        return agent
}

func (s *SerfAgent) GetSerfPort() (int) {
        return s.serfPort;
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

        if !s.Blocking {
                go s.launchAgent(ac, dargs, writer)
        } else {
                s.launchAgent(ac, dargs, writer)
        }

        return nil
}

func (s *SerfAgent) UnRegisterService() (error) {
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
        dargs = append(dargs, "-rpc-addr=" + s.haddr + ":" + strconv.Itoa(s.rpcPort))
        ac.Run(dargs)
        
        return nil
}


func (s *SerfAgent) launchAgent(ac *agent.Command, dargs []string, writer *AgentWriter) {
        hasdefault := false
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


