# Golang OAuth2 Browser Loopback Example

This app launches a browser for the user to authenticate via
an Auth server, then listens on `LISTEN_ADDRESS` for
a response.

## Configuration

The app will look for these environment variables.

```shell
LISTEN_ADDRESS=:9999
REDIRECT_URL=http://127.0.0.1:9999/oauth/callback

## Get these from your Auth provider
CLIENT_ID=
AUTH_URL=
TOKEN_URL=
RESOURCE_URL=
```

## Running

```shell
go run .
```

## Related Links
