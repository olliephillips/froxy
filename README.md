# Froxy

## Fencer proxy with websocket relay and webhooks

For many apps the Fencer REST API is fine for use in a client/server manner. However, in other types of applications, external systems may benefit by being able to "act" on the knowledge (whether a user is inside or outside one or more geofences) and the client application itself is only required to determine position.

Fencer.io currently does not work well ( but may offer support in future) for this second style of "remote tracking" application.

Froxy is a proxy server for the Fencer API. It supports both websockets and webhooks for push notifications to external applications.

### Data flow

1. Client application (Web/App) makes request to Froxy rather than Fencer API
2. Froxy queries a geofence via Fencer API with client user's lat/lng coordinates and Fencer returns response to Froxy
3. Applications connected by Websocket are appraised of change events (inside, outside geofence)
4. Configured webhooks are triggered with payload

### Usage

Froxy is configured with a single `config.toml` file

[work in progress]