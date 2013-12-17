package proxy

import (
	"fmt"
	"github.com/hashicorp/serf/command"
	"github.com/mitchellh/cli"
	"math/rand"
	"strings"
)

// Service contains all the addresses for a current service.
// The Service current supports random retrieval of available services.
type Service struct {
	// name of service
	name string

	// ip and port that host service
	Addresses []string
}

func NewService(name string) *Service {
	return &Service{name: name}
}

func (s *Service) addAddress(addr string) error {
	s.Addresses = append(s.Addresses, addr)
	return nil
}

func (s *Service) getAddress() (address string, err error) {
	err = nil
	address = ""
	if len(s.Addresses) > 0 {
		address = s.Addresses[rand.Intn(len(s.Addresses))]
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

func (w *MembersWriter) GetString() string {
	return string(w.bytes[:])
}

type Registry struct {
	services map[string]*Service
}

func (r *Registry) UpdateRegistry() error {
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
		serviceport := strings.Split(fields[0], "#")

		// there should be no hash marks int the name
		service_name := serviceport[0]
		if service_name == "proxy" {
			continue
		}

		if len(serviceport) != 2 {
			fmt.Errorf("service name incorrectly formatted: %s ", serviceport)
			continue
		}
		complete_address_name := serviceport[1]
//		service_address := strings.Split(complete_address_name, ":")

		if len(fields) != 3 {
			fmt.Errorf("incorrect number of fields for service")
			continue
		}
//		address_fields := strings.Split(fields[1], ":")
//		serf_address := address_fields[0]

//		if serf_address != service_address[0] {
//			fmt.Errorf("Service address does not match serf agent address: %s\n", service_name)
//			continue
//		}

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

func (r *Registry) GetActiveNodes() []string {
        var nodes []string
        unique_nodes := make(map[string]bool)

        for _, service := range r.services {
                for _, val := range service.Addresses {
                        addr := strings.Split(val, ":")[0]
                        unique_nodes[addr] = true
                }
        }

        for node, _ := range unique_nodes {
                nodes = append(nodes, node)
        }

        return nodes
}

func (r *Registry) GetServicesSlice() []string {
	var services []string
	for key, _ := range r.services {
	        services = append(services, key)	
	}

	return services
}

func (r* Registry) GetServiceAddr(service string) (string, error) {
        var err error
        _, ok := r.services[service]
        addr := ""
        if ok {
                serviceInfo := r.services[service]
                addr, err = serviceInfo.getAddress()
        } else {
                err = fmt.Errorf("Service not in registry: " + service)
        }

        return addr, err
}


