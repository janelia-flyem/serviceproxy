/*
register provides functionality that wraps serf to allow service registration.  The main
ojbect is the SerfAgent.  A SerfAgent is created with a service name and port number (the address
is inferred from the current machine).  The SerfAgent will create a serf node encoded in a
manner that allows its discovery by serviceproxy.  It also allows one to unregister the agent.

Registering a service is run as a separate Go process by default.  It will run as long as the program
is alive and the service is still registered.  Unregistering a service is
a blocking call that will execute quickly.
*/
package register
