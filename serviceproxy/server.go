package serviceproxy

import (
    "fmt"
    "strconv"
    "net/http"
    "github.com/mitchellh/cli"  
    "github.com/hashicorp/serf/command"
    "strings"
    "math/rand"
)

const (
    defaultAddr = "127.0.0.1"
    servicePath = "/services"
)


// Service contains all the addresses for a current service.
// The Service current supports random retrieval of available services.
type Service struct {
    // name of service
    name string

    // ip and port that host service
    addresses []string
}

func NewService(name string) *Service {
    return &Service{name: name}
}

func (s *Service) addAddress(addr string) (error) {
    s.addresses = append(s.addresses, addr)
    return nil
}

func (s* Service) getAddress() (address string, err error) {
    err = nil
    address = ""
    if len(s.addresses) > 0 {
        address = s.addresses[rand.Intn(len(s.addresses))]
    } else {
        err = fmt.Errorf("Address does not exist for service: %s", s.name)
    }

    return address, err
}

type MembersWriter struct {
    bytes []byte
}

func (w *MembersWriter) Write(p []byte) (n int, err error) {
    n = len(p)
    w.bytes = append(w.bytes, p...)
    err = nil

    return n, err
}

func (w *MembersWriter) GetString() (string) {
    return string(w.bytes[:])
}

type Registry struct {
    services map[string]*Service
}

var registry Registry

func badRequest(w http.ResponseWriter, msg string) {
    fmt.Println(msg)
    http.Error(w, msg, http.StatusBadRequest)
}


func (r *Registry) updateRegistry() (error) {
    // retrieve members that are alive
    writer := new(MembersWriter)
    ui := &cli.BasicUi{Writer: writer} 
    mc := &command.MembersCommand{Ui: ui}
    var dargs []string
    dargs = append(dargs, "-status=alive")
    mc.Run(dargs) 
 
    mem_str := writer.GetString() 
    mems := strings.Split(strings.Trim(mem_str, "\n"), "\n")

    r.services = make(map[string]*Service)    
    for _, member := range mems {
        fields := strings.Fields(member)
        serviceport := strings.Split(fields[0], ".")
        
        // there should be no periods in the name
        service_name := serviceport[0]
        if service_name == "proxy" {
            continue
        }
        
        if len(serviceport) == 1 {
            fmt.Errorf("service name incorrectly formatted: %s ", serviceport)
            continue
        }
        port_name := serviceport[len(serviceport)-1]        

        if len(fields) != 3 {
            fmt.Errorf("incorrect number of fields for service")
            continue
        }
        address_fields := strings.Split(fields[1], ":")
        address_name := address_fields[0] 
        complete_address_name := address_name + ":" + port_name

        _, ok := r.services[service_name] 
        var service *Service
        if ok {
            service = r.services[service_name]       
        } else {
            service = NewService(service_name)       
            r.services[service_name] = service
        }      

        service.addAddress(complete_address_name)
    }

    return nil
}

func (r *Registry) getServicesString() (string) {
    service_str := ""
    for key, _ := range r.services {
        service_str = service_str + key + "\n"
    }

    return service_str
}


func serviceHandler(w http.ResponseWriter, r *http.Request) {
    requestType := strings.ToLower(r.Method)

    if requestType != "get" {
        badRequest(w, "only supports gets")
    }

    registry.updateRegistry()
    members := registry.getServicesString()
 
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprintf(w, members)
}

func Serve(port int) (error) {
    webAddress := defaultAddr + ":" + strconv.Itoa(port)

    fmt.Printf("Web server address: %s\n", webAddress)

    httpserver := &http.Server{Addr: webAddress}

    http.HandleFunc("/services", serviceHandler)
    

    httpserver.ListenAndServe()


    return nil
}



