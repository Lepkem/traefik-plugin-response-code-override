# traefik-plugin-response-code-override
This traefik plugin allows you to take control over the response code returned along with removing response headers and respnse body.

# Adding & configuring

## Static config
To add a new plugin to a Traefik instance, you must modify that instance's static configuration. The code to be added is provided by the Traefik Pilot UI when you choose Install the Plugin.

> Check the most recent release version - in the example version v0.0.3 is used

```
entryPoints:
  web:
    address: :80

pilot:
    token: xxxxxxxxx

experimental:
  plugins:
    rco:
      moduleName: github.com/Lepkem/traefik-plugin-response-code-override
      version: v0.0.3

```

## Dynamic config
Plugin needs to be configured by adding a dynamic configuration

```
http:
  middlewares:
    my-rco:
      plugin:
        rco:
          RemoveBody: false
          HeadersToRemove:
            - "Referrer-Policy"
            - "Alt-Svc"
            - "X-Be"
            - "H1"
          Overrides:
            '200': 403
            '403': 404
```

The configuration options are following:
* RemoveBody: controls if the plugin will remove response body
* HeadersToRemove: are extra headers that you would like to remove from the response to the client
    > applies to all response code overrrides !
* Overrides: is a list of origin response mapped to desired client response

# Local development
* Download Traefik from
    ```
    curl -L https://github.com/traefik/traefik/releases/download/v2.3.2/traefik_v2.3.2_darwin_amd64.tar.gz -output traefik.tar.gz && tar xzf traefik.tar.gz
    ```
* Download and run docker container
    ```
    docker run -d -p 5000:80 --name whoami containous/whoami
    ```

From this point onward you can develop the plugin. Make sure that the paths will be changed

# Contributors
* Lepkem
* Rafpe
