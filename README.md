# Service Proxy Tool

This go package provides a server that registers different services that
connect to it using the [serf](https://github.com/hashicorp/serf)
package.  The server provides the locations
of these services and can also act as the proxy between clients and the
servers by routing service requests (currently serviceproxy will only
redirect service calls to the appropriate address rather than act as
the intermediary).

##Installing and Usage
To install:

    % go get github.com/janelia-flyem/serviceproxy

To launch the server:

    % serviceproxy [-port WEBPORT]

This will start a web server at the given port on the current
machine.  The command will return the registry address that should
be used for clients to register their service.  The address is the same
as the web address hosting the proxy but the port is different.

To see the supported REST interface, navigate to
"127.0.0.1:WEBPORT/interface".  This will return an interface
specification in [RAML](http://raml.org) format.  

## Making Services Discoverable

To make service discoverable, it must launch a serf agent that
contains the name of the service and its location on the network.

A service can be registered by calling the included go utility "registerservice" (TBD).
To register a service called "foo" on port "15555" to the
the registry address (ADDR), call the following on the machine the service is running on:
    
    % registerservice foo 15555 ADDR

If the -addr option is specified, one can specify the node address that will be running
the service.

We also provide a utility in Go to register a service (TBD).  First,
import the register package:

    import "github.com/janelia-flyem/serviceproxy/register"

To register a service call:

    register.RegisterService("foo", 15555, ADDR)

In theory, service registration can also be done directly by creating
a serf agent on the command using the following particular formatting
(where IP is the IP address of the service):

    % serf agent -node=foo#IP:15555 -port=ADDR

Although serviceproxy is intended to be most relevant
for services that contain a REST interface, the serviceproxy
will discover services independent of their communication protocol.

We recommend that each service contain a URI /interface, which
returns a definition of the REST interface, in RAML.
If JSON is used for a service's interface, the service should
define appropriate JSON schemas.

##TODO

* Add register package to allow clients to easily register services
* Make an example client service
* Add documentation, comments, testing
* Setup proxy server to communicate between client and service rather than just redirect
* Support caching mechanisms

