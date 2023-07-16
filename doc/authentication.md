# Authentication

The `netatmo-exporter` needs an access-token to access the NetAtmo API to retrieve the data on behalf of the user. There's currently only one authentication method available, the [authorization code grant](https://dev.netatmo.com/apidocumentation/oauth#authorization-code). The "client credentials grant" used for username/password authentication in previous versions of the exporter has been deprecated in October 2022 and removed in July 2023.

Unfortunately, the new grant needs user interaction to get the initial access-token. Once the initial authentication is done, the access-token can subsequently be automatically renewed by the exporter.

## Application Registration

The first step in getting your netatmo-exporter authenticated with the NetAtmo API is to create an "application" in the NetAtmo developer console. Visit the [NetAtmo Developer Console] to create an application. The application name and description are mostly relevant to you.

Once your application has been created, take note of the **Client ID** and **Client Secret** as these are the two configuration options needed by the netatmo-exporter.

They can either be provided as command-line arguments (`--client-id` and `--client-secret`) or using environment variables (`NETATMO_CLIENT_ID` and `NETATMO_CLIENT_SECRET`), both variants are equivalent and can also be mixed, depending on your preference of providing configuration and your environment.

## Configuring the Token-File

To avoid needing user interaction every time the netatmo-exporter is restarted, it requires a file to store the access-token. This "token file" is read on startup and written on shutdown, so it needs to be in a location both read- and writable by the netatmo-exporter.

The location of the token-file can be configured using the `--token-file` command-line argument or the `NETATMO_EXPORTER_TOKEN_FILE` environment variable.

More details about the structure of the token-file can be found [here](token-file.md).

## Getting the Access Token

Now that both the application in the developer console and the token-file path are configured, we need to get an initial access-token for the exporter to use.

There are currently two ways of getting the access-token for configuring the netatmo-exporter. Again, both lead to the same result, pick the method that is most suitable for you:

- Authenticate using the [NetAtmo Developer Console] and manually enter the token into the configuration
- Use the integrated web-interface of the netatmo-exporter

### Using the Developer Console

1. Open the [NetAtmo Developer Console] and click on the button for your created application.
2. Scroll down a bit until you reach the section titled "Token Generator".
3. Select the `read_station` scope and click on the "Generate Token" button.
4. You will be redirected to an authorization page from NetAtmo. Click "Yes, I accept".
5. You will return to the previous page with a new section which contains an "Access Token" and a "Refresh Token".

You now have the access-token and refresh-token needed for the netatmo-exporter. There are two ways of getting this information into the exporter:

- Open the web-interface of the netatmo-exporter and enter the **refresh-token** into it.

  The exporter has a simple web-interface when you navigate to it (for example at `http://localhost:9210` if running locally). Paste the **refresh-token** into the textfield and click the update button to submit the token to the exporter.

- Create a token-file and let the exporter read it

  See the section above and the [token file document](token-file.md) for more detail.

### Using the Integrated Web-Interface

This method requires the exporter to be reachable by the user using a web-browser. The URL needed for that depends on the environment the exporter is placed in. If the exporter is running on the same machine that the user is using it can be as simple as opening `http://localhost:9210`.

Because the user is automatically redirected back to the exporter during the authentication flow, the exporter needs to be configured with the URL it is reachable from. This can be done using the `--external-url` command-line parameter or the `NETATMO_EXPORTER_EXTERNAL_URL` environment variable. For example if the exporter is running on a computer with the IP-address `192.168.1.10` and using the default port, then the configuration would be:

```plain
--external-url http://192.168.1.10:9210
```

Keep in mind that this URL does not need to be reachable _from the internet_, but just for the user authenticating the exporter.

Once the exporter is configured using the client-id, client-secret, token-file and external-url, you should be able to visit the URL. In the interface shown to you, click the "authorize here" link. This should redirect you to the NetAtmo website and ask for confirmation.

Once the confirmation is given, you will be redirected to the exporter and end up at the same page you started. It should now show you as authenticated. If this redirect does not work properly, check the `--external-url` configuration.

[NetAtmo Developer Console]: https://dev.netatmo.com/apps/
