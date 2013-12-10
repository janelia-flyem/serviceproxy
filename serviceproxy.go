package main

import (
    "flag"
    "fmt"
    "os"
    "github.com/stephenplaza/serviceproxy/serviceproxy"
)

const defaultPort = 15333

var (
    showHelp = flag.Bool("help", false, "")
    portNum = flag.Int("port", defaultPort, "")
    debugSerf = flag.Bool("debug", false, "")
)

const helpMessage = `
The service proxy tool registers different services through the
serf package. 
 
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

    proxy := serviceproxy.ServiceProxy{Port: *portNum, Debug: *debugSerf}
    proxy.Run()
}


