package serviceproxy

import (
    "fmt"
    "os"
    "github.com/hashicorp/serf/command/agent"
    "github.com/mitchellh/cli"  
    "os/signal"
    "io/ioutil"
    "syscall"
    "strconv"
    "net"
)

const defaultPort = 7946

type ServiceProxy struct {
    Port int
    Debug bool
}

func (proxy *ServiceProxy) Run() (error) {
    // create a default serf agent at the default port 7373
    writer := ioutil.Discard
    if proxy.Debug {
        writer = os.Stdout
    }
    ui := &cli.BasicUi{Writer: writer} 
    ac := &agent.Command{Ui: ui, ShutdownCh: make(chan struct{}),}
    var dargs []string
    dargs = append(dargs, "-node=proxy."+strconv.Itoa(proxy.Port))
    
    hname, _ := os.Hostname()
    addrs, _ := net.LookupHost(hname)

    serfaddr := addrs[1] + ":" + strconv.Itoa(defaultPort)
    dargs = append(dargs, "-bind="+serfaddr)
    go ac.Run(dargs)
   
    // address for clients to register (port does not need to be specified
    // by the client if using the Go register interface)  
    fmt.Printf("Registry address: %s\n", serfaddr)    

    // exit server if user presses Ctrl-C 
    go func () {
        sigch := make(chan os.Signal)
        signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
        <-sigch
        fmt.Println("Exiting...")
        os.Exit(0)
    } ()

    // ?! create web server 
    return Serve(proxy.Port)    
}
