package main

import (
    "flag"
    "fmt"
    "github.com/janelia-flyem/serviceproxy/register"
    "os"
    "strconv"
)

var (
	showHelp  = flag.Bool("help", false, "")
	removeAgent  = flag.Bool("remove", false, "")
)


const helpMessage = `
Register a service on a given port to the registry.

Usage: serviceregister <service name> <service port> <registry address>
    -remove     (flag)      Remove service from registry
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

        serfagent := register.NewAgent(flag.Arg(0), port)
 
        if *removeAgent {
               fmt.Println("unregister")
                fmt.Println(flag.Arg(0))
                fmt.Println(flag.Arg(1))
                fmt.Println(flag.Arg(2))
               serfagent.UnRegisterService() 
               os.Exit(0)
        } else {
                serfagent.Blocking = true
                serfagent.RegisterService(flag.Arg(2)) 
        }

        
}
