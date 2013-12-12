# Service Proxy Tool

This go package provides a server that registers different services that
connect to it using the [serf](https://github.com/hashicorp/serf)
package.  The server provides the locations
of these services and can also act as the proxy between clients and the
servers by routing service requests.

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

We provide a utility in Go to register a service (TBD).  First,
import the register package:

    import "github.com/janelia-flyem/serviceproxy/register"

To register a service called "foo" on port "15555" to the
the registry address (ADDR), call:

    register.RegisterService("foo", 15555, ADDR)

Non-Go services running on address ip address (IP) can easily
register with the serviceproxy by creating a serf
agent on the command-line:

    % serf agent -node=foo#IP:15555 -port=ADDR

Although serviceproxy is intended to be most relevant
for services that contain a REST interface, the serviceproxy
will discover services independent of their communication protocol.

We recommend that each service contain a URI /interface, which
returns a definition of the REST interface, in RAML.
If JSON is used for a service's interface, the service should
define appropriate JSON schemas.

##TODO

* Add interface for showing interface
* Add interface for passing requests through to servers
* Add register package to allow clients to easily register services
* Make an example client service
* Add documentation, comments, testing
* Support service state or other status?

