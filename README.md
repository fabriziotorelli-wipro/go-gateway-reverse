# go-gateway-reverse

Go-Lang Simple Gateway Reverse Proxy

## Goals

Define a simple reverse proxy enabled by reference label and triggering single or load-balanced (simple balancing strategy) gateway to configured services.


## Architecture

Reverse Proxy, acquiring an URL defined as :
`<protocol>>//<host>:<service-operating-port>/label/rest/.../...?query=....&....`
The gateway search in the configuration of the service the resource and accorrding to site policy or balanced site policy to define the current reverse upstreaming chanell to :
`<real-protocol>>//<real-host>:<real-port>/api../rest/...?apiQuery...&query...`
or
`<real-protocol>>//<real-host>:<real-port>/rest/.../api..?apiQuery...&query...`

This allows to retain the original verb and method, challenging a channel to the real service with the service specified policies


## Configuration

Availale Configurations are :
* serverindex.json defining main information to open an indexing service available to the other service to require root meta-data
* config.json with multiple services on multiple ports available for the gateway protocol, defining each of them a datafile
* <datafiles>.json define streaming informations and overriding of service properties to allow multiple APIs broadcasted on a single Gateway port

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
