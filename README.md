# Revox

Revox is a project which aims to implement a reverse proxy. The package contains the proxy server, services, configuration file, helm chart and minimal unit tests.

The solution can be deployed using either docker-compose or helm.
It supports HTTP/1.1 and messages encoded in JSON.

The proxy server listens for HTTP requests and forwards them to downstream services if the host header matches. In the eventuality of a match, the proxy server will implement the load balancing strategy specified in the configuration file and if no explicit strategy is declared, it will default to round robin.

*Written in Go.*

## Features
- Load Balancing via two methods:
  - Random 
  - Round Robin
- Multiple downstream services
- Healthchecks
- Retries
- HTTP Cache
- YAML configuration

## Getting Started

### Build

To build the docker image, `make build`.

## Docker-compose
To run the docker-copose setup, `make start`.
To test, 
```
curl http://localhost:8080/get -H "Host: my-service.my-company.com"
curl http://localhost:8080/get -H "Host: httpbin.org"
curl http://localhost:8080/cache/5 -H "Host: my-service.my-company.com"

```

## Helm installation.

`kubectl create ns demo-revox`

`kubectl config set-context --current --namespace demo-revox`

`make install-demo` will install the chart using the demo values (deploys 2 httpbins, and a proxy configuration for those)

In order to test with the local k8s cluster (docker for mac):
```
k port-forward svc/revox 8080&

curl http://localhost:8080/get -H "Host: httpbin.org"
Handling connection for 8080
{
  "args": {},
  "headers": {
    "Accept": "*/*",
    "Accept-Encoding": "gzip",
    "Host": "httpbin.org",
    "User-Agent": "curl/7.77.0",
    "X-Amzn-Trace-Id": "Root=1-626fbda7-2de2976f7d97c0370eda402c"
  },
  "origin": "127.0.0.1, 82.76.227.130",
  "url": "https://httpbin.org/get"
}
```


## TODO:

Improve unit testing