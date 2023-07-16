# Token File

The exporter uses a "token file" to save the authentication information it receives from NetAtmo. It is a simple JSON file with usually three attributes:

```json
{
  "access_token": "a long string",
  "refresh_token": "another long string",
  "expiry": "2023-07-16T20:32:06.400559267+02:00"
}
```

## Attributes

All of the three attributes are necessary for the exporter to work correctly, but they have different purposes:

- `access_token` this is the "key" that is used to communicate with the NetAtmo API and fetch the data available for the user. It is only valid for a limited time after which the API will return a 403 error when the access-token is used.
- `expiry` this is the time when the `access_token` will expire. The exporter needs to know this, so that it can get a new access-token in time ("refresh" it).
- `refresh_token` this "key" is used when the exporter wants to renew the `access_token`. It can not be used to retrieve the data, only to get a new access-token.

## Startup

When starting the exporter it will try to load the file specified with `--token-file`. If it does not exist, it will just start up without any authentication and wait for the user to initiate authentication.

If the token-file is available, it is read by the exporter. If all three attributes are available and the token is still valid, the exporter will immediately start working properly.

The exporter will issue warnings if the token-file is missing attributes during startup. If the `expiry` is missing a new short expiry will be set, so that the refresh happens as early as possible.

If the `refresh_token` is missing during the startup of the exporter, it will issue a warning that it can not automatically refresh the token. It will continue to work normally until the `access_token` expires after which the user needs to initiate a new authentication. The exporter can not automatically recover from this case.

**Note:** Due to the facts that the `access_token` can be regenerated using the `refresh_token` and that the exporter will automatically set an early `expiry`, it is technically possible to start the exporter with a token-file that only contains a `refresh_token`. If the refresh token is valid, it will immediately renew the token and have a proper `access_token` and `expiry` afterward.

## Shutdown

When the exporter has a valid token in memory when shutting down, it will try to save the token to the path specified using `--token-file`. It will emit an error if this is not successful, but will not try again.
