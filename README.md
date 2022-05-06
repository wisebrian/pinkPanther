# Revox

Revox is a project which aims to implement a reverse proxy. The package contains the proxy server, services, configuration file, helm chart and minimal unit tests.

The solution can be deployed using either docker-compose or helm.
It supports HTTP/1.1 and messages encoded in JSON.

The proxy server listens for HTTP requests and forwards them to downstream services if the host header matches. In the eventuality of a match, the proxy server will implement the load balancing strategy specified in the configuration file and if no explicit strategy is declared, it will default to round robin.

*Written in Go. Had [godon](https://github.com/bmf-san/godon) as a starting point.*

## Design decisions

I tried to make the project as modular as possible in order to facilitate testing and configuration.

## Features
- Load Balancing via two methods:
  - Random 
  - Round Robin
- Multiple downstream services
- Healthchecks
- Retries
- HTTP Cache
- Configuration file




## Getting Started

First clone the repository on your machine.

```
git clone https://github.com/wisebrian/revox.git`
```

CD into revox and run `make go-install` and `make go-build`.



### Build

To build the docker image, `make build`.



## Docker-compose

Start revox by executing `make start` which will run the docker-copose setup. Please check if the three containers are up and running, 1 proxy (`revox`) & 2 services (`httpbin`).

Other useful commands would be:

```
make restart                              ///      rebuild and recreate the containers
make stop                                 ///      stop the containers
make remove                               ///      remove containers, networks, etc.
make state                                ///      lists and displays containers state
make logs                                 ///      prints log stream
```


To test HTTP requests:
```
curl http://localhost:8080/get -H "Host: my-service.my-company.com"           /// If you want to test
curl http://localhost:8080/get -H "Host: httpbin.org"                         /// different Host headers

curl http://localhost:8080/cache/5 -H "Host: my-service.my-company.com"       /// If you want to test cache hits

```
to simulate an upstream which exposes Cache-Control response headers, we can call the /cache/ttl endpoint of httpbin. Revox will try to cache the HTTP response based on RFC 7234, using `github.com/lox/httpcache`.
To make sure it is working you can open up a different terminal window and
watch the logs by using the command listed earlier in the section and adding -i (include HTTP response headers) to the curl command. This will print out if the request was a cache hit or miss which can also be verified in the logs because no requests will be forwarded by the proxy.



## Helm installation


Create a namespace for revox using: `kubectl create ns demo-revox`

Set context to revox by running: `kubectl config set-context --current --namespace demo-revox`
In order to install the chart using the demo values found in the `values-demo.yaml` file, run
`make install-demo` (deploys 2 httpbins, and a proxy configuration for those).

To test with the local k8s cluster (Kubernetes feature in Docker for Mac):
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


## Limitations


- Timeouts per host not implemented yet
- O(n) time complexity



## Improvements / TODOs :


- Improve unit testing
- Support for rate limiting
- Implement timeouts per host
- Healthchecks
- O(1) time complexity