/*
proxy contains the main functionality for the serviceproxy web server and
maintains registry that contains the different client services.

The server handles requests that list the registry contents and to retrieve
the location of a service.  It also redirects service requests automatically
to a node (random selection for now).
*/
package proxy
