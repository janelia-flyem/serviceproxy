# Service Proxy Tool [![Picture](https://raw.github.com/janelia-flyem/janelia-flyem.github.com/master/images/gray_janelia_logo.png)](http://janelia.org/)
[![Build Status](https://drone.io/github.com/janelia-flyem/serviceproxy/status.png)](https://drone.io/github.com/janelia-flyem/serviceproxy/latest)

This go package provides a server that registers different services that
connect to it using the [serf](https://github.com/hashicorp/serf)
package.  The server provides the locations of these services and
can also act as the proxy between clients and the
servers by routing service requests (currently serviceproxy will only
redirect service calls to the appropriate address rather than act as
the intermediary).

The tool has been tested on linux but should work flexibly in other
environments.  Because the serf protocol can work on external
networks, services can be managed across the internet.

##Installation and Basic Usage
This package includes the main executable for launching the service
proxy ('serviceproxy').  The package also includes an
example client service 'adderexample'
and a service registration executable 'serviceregister'.

To install serviceproxy:

    % go get github.com/janelia-flyem/serviceproxy

To install adderexample and serviceregister

    % cd $GOPATH/src/github.com/janelia-flyem/serviceproxy
    % go get ./...
    
Before launching serviceproxy, GOPATH should be set
to the path that roots the downloaded go packages.  This is
necessary for now since there are some static file resources
needed by the web server.

To launch the server:

    % serviceproxy [-port WEBPORT (default 15333)]

This will start a web server at the given port on the current
machine.
The command will return the registry address (ADDR) (a different
port from the web server) used by clients to register their service.

The rest interface specification is in [RAML](http://raml.org) format.
To view the interface, navigate to
"127.0.0.1:WEBPORT/interface".  To view the RAML interface, we use
the [api-console](https://github.com/mulesoft/api-console) javascript-based viewer.
(The RAML text can be found at proxy/interface.go.)

The adderexample service can be launched by:

    % adderexample ADDR
    
Using serf node discovery, the serviceproxy is now aware of the
adder service.  The location of this service and its interface are
exposed by the serviceproxy web server.

## Designing a Client Service

The adderexample provides a sample service written in Go.  While the example
is very simple, it indicates best practices one should follow
when designing a client service.

At the minimum, a service only needs
register itself to the serviceproxy as defined in the next section.  To better
exploit features in the serviceproxy, the service should implement a REST interface
at the registered address.
Beyond these bare requirements, we recommend the following:

* Define a REST interface using RAML.  This RAML should be accessible via the /interface/interface.raml URI.
* Non-binary request and response data should use JSON format when possible
* Specific JSON mime-types should be defined using (JSON schema)[http://json-schema.org/] for request and response data
* The name of the json schema should also define the mimetype name, so that "<schema>.json" should be of mime type "<schema>+json"
* Potentially reusable JSON schema should be saved in some global repository that is CORS-enabled and included in a RAML specification using !include syntax.  Janelia maintains a JSON schema directory at: http://janelia-flyem.github.io/schema (not currently CORS enabled)
* Calls to services should be non-blocking if they require non-trivial computation and immediately return a callback URL(s) (unless otherwise noted in the interface), which will indicate the service status or provide the result(s)
* JSON schema validators should be used to validate JSON data

The adderexample follows all of these suggestions (except that no global JSON schema is used).


## Registering Client Service

To make service discoverable, it must launch a serf agent that
contains the name of the service and its location on the network.

A service can be registered by calling the included go utility "registerservice".
To register a service called "foo" on port "15555" to the
the registry address (ADDR), call the following on the machine the service is running on:
    
    % registerservice foo 15555 ADDR

This process will run indefinitely and could be run as a background daemon.
To unregister, this process must be killed.

We also provide a utility in Go to register a service.  First,
import the register package:

    import "github.com/janelia-flyem/serviceproxy/register"

To register a service call:
    
    serfagent := register.NewAgent("foo", 15555)
    serfagent.RegisterService(ADDR)

The service will be unregisterd by terminating the parent
program or by calling:

    serfagent.UnregisterService()

In theory, service registration can also be done by creating
a serf agent explicitly on the command using the following
specific formatting
(where IP is the IP address of the service):

    % serf agent -node=foo#IP:15555 -port=ADDR


##TODO

* Better decouple serf service agent wrapper in current implementation
* Support redirection protocols that favor certain nodes as appropriate (maybe each service can provide information)
* Enable a system of versioning; the same service with different versions should be distinguished
* Support wrapping serf to check services for busy states
* Have serviceproxy act as an intermediary between users and the services (beyond redirection)
* Allow for some basic orchestration between services
* Enable data encryption and authentication for services in non-internal networks
* Provide caching mechanisms for quicker service access
