# go-gateway-reverse

Go-Lang Simple Gateway Reverse Proxy

## Goals

Define a simple reverse proxy enabled by reference label and triggering single or load-balanced (simple balancing strategy) gateway to configured services.


## Pre-requisites

To compile and run this project you have to check availability of following software:
* [Go](https://golang.org/dl/) (tested with version 1.8)
* Any program (curl, wget) or Browser plugin (REST Easy, etc..) to test token http calls in HEAD space 

## Architecture

Reverse Proxy, acquiring an URL defined as :
`<protocol>://<host>:<service-operating-port>/label/rest/.../...?query=....&....`
The gateway search in the configuration of the service the resource and accorrding to site policy or balanced site policy to define the current reverse upstreaming chanell to :
`<real-protocol>://<real-host>:<real-port>/api../rest/...?apiQuery...&query...`
or
`<real-protocol>://<real-host>:<real-port>/rest/.../api..?apiQuery...&query...`

This allows to retain the original verb and method, challenging a channel to the real service with the service specified policies

<p align="center"><img width="640" height="480" src="/images/arch1.png" /></p>

The single GateWay Port Server manages identifiers that mask the access to the real service and open streaming channels
to the target server/port as required in the original call :

<p align="center"><img width="640" height="480" src="/images/arch2.png" /></p>

## Configuration

Availale Configurations are :
* serverindex.json defining main information to open an indexing service available to the other service to require root meta-data
* config.json with multiple services on multiple ports available for the gateway protocol, defining each of them a datafile
* `<datafiles>.json` define streaming information and overriding of default Gateway Service Ports, to allow multiple APIs provisioning on a single Gateway Port

## Index Server

Server Index file configure a service JSON http server that expose following end-point
* /poweroff : PowerOff all gateway Port Servers and the Gateway Application will exit
* /{n} : return the JSON output of service list, exposed in the Gateway Port Server at the {n} position in the Gateway Port servers Index, used by GateWay Port Servers, when a root call is required
* /error : service for balancing and shaping errors, required by GateWay Port Servers

*IMPORTANT :*
_This server requires a Token Protection information, defined in own configuration_

Configuration descriptor (`indexservice.json`) :
* "enabled": Status of service for Index Server
* "ipaddress": Host Name/IP Address used by Index Server or "" (for any address)
* "serviceaddress": Address where GateWay Ports should recover information or "localhost"
* "port": Port Number used by Index Server (integer)
* "protocol": Protocol used for connecting the Index Server by GateWay Ports
* "usetokenprotection": Flag defining the user to check in the request HEADs the `X-GATEWAY-TOKEN` tag
* "securitytoken": Security Token recovered in Head Tag `X-GATEWAY-TOKEN`
* "usetls": Enable/Disable SSL/TLS configuration for the GateWay Index Server
* "tlsx509certificatefile": X509 Certificate full qualified file path
* "tlsx509certificatekeyfile": Certificate Server Key full qualified file path

Example :
```
{
  "enabled": true,
  "ipaddress": "",
  "serviceaddress": "localhost",
  "port": 10098,
  "protocol": "http",
  "securitytoken": "J1qK1c18UUGJFAzz9xnH56584l4",
  "usetls": true,
  "tlsx509certificatefile": "./data/server.pem",
  "tlsx509certificatekeyfile": "./data/server.key"
}
```

## Port Servers

GateWay Port Server List file configure a set of ports that consumes services defined in a specific data file.

Configuration descriptor for any of the ports (`config.json`) :
* "ipaddress": Host Name/IP Address used by Gateway Port Server or "" (for any address)
* "port": Port Number used by Gateway Port Server (integer)
* "apiurl": API URL Base Address
* "concatenate": Flag defining if server should concatenate API and Call path and Query parameters 
* "beforeapi": Concatenate call path and Query parameters before the API URL data
* "servicefile": Full qualified path of JSON file containing Gateway Port Service data
* "protocol": Protocol used for connecting the GateWay Port Services
* "user": Authentication User Name / Code (not yet implemented)
* "password": Authentication User Password / Token (not yet implemented)
* "usetokenprotection": Flag defining the user to check in the request HEADs the `X-GATEWAY-TOKEN` tag
* "securitytoken": Security Token recovered in Head Tag `X-GATEWAY-TOKEN`
* "usetls": Enable/Disable SSL/TLS configuration for the GateWay Port Server
* "tlsx509certificatefile": X509 Certificate full qualified file path
* "tlsx509certificatekeyfile": Certificate Server Key full qualified file path

Example :
```
[
  {
    "ipaddress": "",
    "port":10099,
    "apiurl": "/api/json?pretty=true",
    "concatenate": true,
    "beforeapi": true,
    "servicefile": "./data/data.json",
    "protocol": "http",
    "user": "",
    "password": "",
    "usetokenprotection": false,
    "securitytoken": "",
    "usetls": false,
    "tlsx509certificatefile": "",
    "tlsx509certificatekeyfile": ""


  },
  {
    "ipaddress": "",
    "port":10100,
    "apiurl": "/api/json?pretty=true",
    "concatenate": true,
    "beforeapi": true,
    "servicefile": "./data/data2.json",
    "protocol": "http",
    "user": "",
    "password": "",
    "usetokenprotection": true,
    "securitytoken": "J1qK1c18UUGJFAzz9xnH56584l4",
    "usetls": false,
    "tlsx509certificatefile": "",
    "tlsx509certificatekeyfile": ""
  }
]
```

## Port Servers Data Files

GateWay Port Server Service file contains information about upstreaming and reverse proxy rules, shading the real addresses.

Configuration descriptor for any of the port services (`<data-file>.json` contained in the `config.json` single service `servicefile` JSON element) :
* "site" : Desired Label for masking the call to the server (overlapping of the label causes Load Balancing)
* "address" : Host Name/IP Address used by the reverse proxy engine to connect the real server...
* "port" : Port used by the reverse proxy engine to connect the real server...
* "protocol": URL protocol used in the merging of the real server proxying
* "scheme" : URL schema used in the URL element, merging of the real server proxying
* "type" : Informative data
* "override": Flag that define if following information overrides the GateWay Port Configuration items
* "apiuri": API URL Base Address
* "concatenatepath": Concatenate call path and Query parameters before the API URL data
* "concatenatebeforeapi": Concatenate call path and Query parameters before the API URL data

Example :
```
[
  {
    "site" : "Jenkins1",
    "address" : "10.10.243.50",
    "port" : 8080,
    "protocol": "http",
    "scheme" : "http",
    "type" : "json",
    "override": false,
    "apiuri": "",
    "concatenatepath": false,
    "concatenatebeforeapi": false
  },
  {
    "site" : "Jenkins2",
    "address" : "10.10.243.53",
    "port" : 8080,
    "protocol": "http",
    "scheme" : "http",
    "type" : "json",
    "override": true,
    "apiuri": "/api/json?pretty=true",
    "concatenatepath": true,
    "concatenatebeforeapi": true
  }
]
```


## Checkout and test this repository

Go in you `GOPATH\src` folder and type :
```
 git clone https://github.com/fabriziotorelli-wipro/go-gateway-reverse.git gateway

```

Project GO package folder name is `gateway`.

## Build

It's present a make file that returns an help on the call :

```
make
```
provided help returns :
```
make [all|test|build|exe|run|clean|install]
all: test build exe run
test: run unit test
build: build the module
exe: make executable for the module
clean: clean module C objects
run: exec the module code
install: install the module in go libs
```

Alternatively you can call following commands :
 * `go build .` to build the project
 * `go test` to run unit and integration test on the project
 * `go run main.go` to execute the project
 * `go build --buildmode exe .` to create an executable command

 
## Further test 

You can access information on GateWay Token protected ports using following command :

* POST:
```
curl -i -H Accept:application/json -H X-GATEWAY-TOKEN:<YOUR-TOKEN-HERE> -X POST http://<HOST>:<PORT>/<MASKED-SERVICE> -H Content-Type: application/json -d ''
```

* GET:
```
curl -i -H Accept:application/json -H X-GATEWAY-TOKEN:<YOUR-TOKEN-HERE> -X GET http://<HOST>:<PORT>/<MASKED-SERVICE>
```

## TLS Cerificate test

##### Generate private key (.key)

```sh
# Key considerations for algorithm "RSA" ≥ 2048-bit
openssl genrsa -out server.key 2048
    
# Key considerations for algorithm "ECDSA" ≥ secp384r1
# List ECDSA the supported curves (openssl ecparam -list_curves)
openssl ecparam -genkey -name secp384r1 -out server.key
```

##### Generation of self-signed(x509) public key (PEM-encodings `.pem`|`.crt`) based on the private (`.key`)

```sh
openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650
```

Now you can assign enable TLS mode for the server, using X509 SSL Certificate and Server Key, configuring the relative information on the PORT, then you can call the SSL channel.

Here one example of call :

```sh
curl -k https://<gw-port-address>:<gw-port-number>/ -v –key /path/to/server.key –cert /path/to/server.key https://<gw-port-address>:<gw-port-number>/<gw-port-end-point>
```

## Execution

*Coming soon ...*

## License

Copyright (c) 2016-2017 [BuildIt, Inc.](http://buildit.digital)

Licensed under the [MIT](/LICENSE) License (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[https://opensource.org/licenses/MIT](https://opensource.org/licenses/MIT)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
