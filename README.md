# Revox

Revox is a project which aims to implement a reverse proxy. Written in Go.
The package contains the proxy server, services and minimal unit tests.

The solution can be deployed using either docker-compose or helm.


## Features
- Load Balancing via two methods:
  - Random 
  - Round Robin
- Multiple upstream services
- Healthchecks
- Retries
- HTTP Cache
- YAML configuration

## Build

To build the docker image, `make build`.

## Docker-compose
To run the docker-copose setup, `make start`.
To test, 
```
curl http://localhost:8080/get -H "Host: httpbin.org"
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