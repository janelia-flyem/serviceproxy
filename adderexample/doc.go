/*
adderexample is a simple REST service that adds two numbers together.  It registers
itself with the serviceproxy tool.

Usage: adderexample <registryaddress> [-port default: 23230]

The RAML interface can be retrieved by navigating to http://127.0.0.1:23230/interface/interface.raml.
To view a nicely formatted version of the interface, use the service proxy with the URI:
/services/adder/interface.
*/
package main
