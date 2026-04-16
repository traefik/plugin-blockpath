# Block Path

[![Build Status](https://github.com/traefik/plugin-blockpath/workflows/Main/badge.svg?branch=master)](https://github.com/traefik/plugin-blockpath/actions)

Block Path is a middleware plugin for [Traefik](https://github.com/traefik/traefik) which blocks requests
with a configurable status code (default `403 Forbidden`) when the requested HTTP path matches one of the
configured [regular expressions](https://github.com/google/re2/wiki/Syntax).

## Configuration

## Static

```toml
[pilot]
    token="xxx"

[experimental.plugins.blockpath]
    modulename = "github.com/traefik/plugin-blockpath"
    version = "v0.2.1"
```

## Dynamic

To configure the `Block Path` plugin you should create a [middleware](https://docs.traefik.io/middlewares/overview/) in 
your dynamic configuration as explained [here](https://docs.traefik.io/middlewares/overview/). The following example creates
and uses the `blockpath` middleware plugin to block all HTTP requests with a path starting with `/foo`. 

```toml
[http.routers]
  [http.routers.my-router]
    rule = "Host(`localhost`)"
    middlewares = ["block-foo"]
    service = "my-service"

# Block all paths starting with /foo
[http.middlewares]
  [http.middlewares.block-foo.plugin.blockpath]
    regex = ["^/foo(.*)"]
    code = 404 # optional, responds with code 404 instead of 403

[http.services]
  [http.services.my-service]
    [http.services.my-service.loadBalancer]
      [[http.services.my-service.loadBalancer.servers]]
        url = "http://127.0.0.1"
```
