# Service Proxy Tool

This go package provides a server that registers different services that
connect to it using the serf package.  The server provides the locations
of these services and can also act as the proxy between clients and the
servers by routing service requests

##Installing and Using
To install:

    % go get github.com/janelia-flyem/serviceproxy

To launch the server:

    % serviceproxy

To make service discoverable, they must launch a serf-agent that joins
the proxy servers network.  The node name should be <servicename>.<portnum>.

##TODO

* Factor out registry functionality
* Add interface for getting addresses (use json format)
* Add interface for getting nodes
* Add interface for showing interface
* Add interface for passing requests through to servers
* Add register package to allow clients to easily register services
* Make an example client service
* Add documentation, comments, testing
* Support service state or other status?

