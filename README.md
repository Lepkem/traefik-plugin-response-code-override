# traefik-plugin-response-code-override
This traefik plugin allows you to take control over the response code returned along with removing response headers and respnse body.

# Configuration
Traefik rules will require that plugin has the configuration following the spec of:

```
    RemoveBody: true
    HeadersToRemove:
        - "Referrer-Policy"
        - "Alt-Svc"
    Overrides:
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