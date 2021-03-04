# Block Path With Whitelist

[![Build Status](https://github.com/traefik/plugin-blockpath/workflows/Main/badge.svg?branch=master)](https://github.com/traefik/plugin-blockpath/actions)

Block Path With Whitelist is a middleware plugin for [Traefik](https://github.com/traefik/traefik) which sends an HTTP `403 Forbidden` 
response when the requested HTTP path matches one the configured [regular expressions](https://github.com/google/re2/wiki/Syntax). Users of the plugin can also choose a Whitelist regex array to override the blocklist if required.

## Configuration

## Static

```toml
[pilot]
    token="xxx"

[experimental.plugins.blockpath]
    modulename = "github.com/AdrianFletcher/plugin-blockpath-with-whitelist"
    version = "v0.2.1"
```

## Dynamic

To configure the `Block Path` plugin you should create a [middleware](https://docs.traefik.io/middlewares/overview/) in 
your dynamic configuration as explained [here](https://docs.traefik.io/middlewares/overview/). The following example creates
and uses the `blockpath` middleware plugin to block all HTTP requests with a path starting with `/foo`, but allow a path starting with `/foobar`.

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
    regexwhitelist = ["^/foobar(.*)"]

[http.services]
  [http.services.my-service]
    [http.services.my-service.loadBalancer]
      [[http.services.my-service.loadBalancer.servers]]
        url = "http://127.0.0.1"
```
