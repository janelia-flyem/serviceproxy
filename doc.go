/*
serviceproxy is a server that is a registry of different client services.
It implements service discovery by wrapping the serf tool.  Because
serf can work on external networks, services can be managed
across the internet.

Before launching serviceproxy, GOPATH should be set
to the path that roots the downloaded go packages.  This is
necessary for now since there are some static file resources
needed by the web server.

To launch the server:

    % serviceproxy [-port WEBPORT (default 15333)]

This will start a web server at the given port on the current machine.
The command will return the registry address (ADDR) (a different
port from the web server) used by clients to register their service.

The rest interface specification is in RAML (http://raml.org) format.
To view the interface, navigate to
"127.0.0.1:WEBPORT/interface".
*/
package main
