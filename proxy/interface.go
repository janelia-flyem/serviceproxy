package proxy

const ramlInterface = `#%%RAML 0.8
title: "Janelia Service Handler"
version: v1
baseUri: /
/services:
  get:
    description: "Retrieves a list of registered services"
    responses:
      200:
        body:
          application/json:
            schema: |
              { "$schema": "http://json-schema.org/schema#",
                "title" : "List of services registered",
                "type" : "object",
                "properties" : {
                  "services" : {
                    "description": "List of services",
                    "type" : "array",
                    "uniqueItems" : true,
                    "items" : {
                      "type" : "string"
                    }
                   }
                 },
                 "required" : ["services"] 
              }
  /{service}/node:
    get:
      description: "Finds a node for the requested service"
      responses:
        200:
          body:
            application/json:
              schema: |
                { "$schema": "http://json-schema.org/schema#",
                  "title" : "Retrieves a node that hosts the given service",
                  "type" : "object",
                  "properties" : {
                    "service-location" : {
                    "description" : "Web address or null if an address does not exist",
                    "type" : "string"
                    }
                  },
                  "required" : ["service-location"] 
                }
  /{service}/interface:
    get:
      description: "Shows the RAML interface of the given service. The service must retrieve a RAML for /interface"
/nodes:
  get:
    description: "Show addresses of registered compute nodes"
    responses:
      200:
        body:
          application/json:
            schema: |
              { "$schema" : "http://json-schema.org/schema#",
                "title" : "Shows all the nodes visible to the proxy server",
                "type" : "object",
                "properties" : {
                  "nodes": {
                    "description": "List of different compute node addresses",
                    "type" : "array",
                    "uniqueItems" : true,
                    "items" : {
                      "type" : "string"
                    }
                  }
                }
              }
/interface:
  get:
    description: "Shows the RAML interface of the service handler."
  /raw:
    get:
      description: "Retrieves the raw RAML interface for the service handler"
      responses:
        200:
          body:
            application/raml+yaml:
    /exec/{service}/{request}:
      description: "Redirects any request to the given service"
`            
 
