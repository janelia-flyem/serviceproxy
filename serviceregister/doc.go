/*
serviceregister is an executable that wraps the serf tool to register a given service
with the given serviceproxy.

Usage: serviceregister <service name> <service port> <service proxy registry address>

This command is blocking and should be run as a background process.  It must also be
run on the machine the service is running on.  Ideally, this process should be tied
to the lifetime of the service.  Otherwise, the service can be unregistered by killing
the process.
*/
package main
