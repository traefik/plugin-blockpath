# Block Path

This [Traefik](https://github.com/containous/traefik) plugin is as middleware which sends an HTTP `403 Forbidden` response 
when the HTTP request path matches one the configured [expressions](https://github.com/google/re2/wiki/Syntax).

## Configuration

To configure this plugin you should add its configuration to the Traefik dynamic configuration as explained [here](https://docs.traefik.io/getting-started/configuration-overview/#the-dynamic-configuration).
The following snippet shows how to configure this plugin with the File provider in TOML and YAML: 

```toml
# Block all paths starting with /foo
[http.middlewares]
  [http.middlewares.block-foo.blockPath]
    regex = ["^/foo(.*)"]
```

```yaml
# Block all paths containing bar
http:
  middlewares:
    block-bar:
      plugin:
        blockpath:
          regex: 
            - "bar"
```
