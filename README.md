# Froxy

## Fencer proxy with websocket relay and webhooks

Froxy is a proxy server for the Fencer API. It fulfills the remote application requirement by providing support for both websockets and webhooks both of which enable communication with supporting external apps.

### Background

For many Client/Server style Apps, the "Geofencing As A Service" REST API model offered by Fencer.io is exactly what is required. 

In some other types of application, where external systems can benefit by being able to "act" on knowledge of whether a user is inside/outside one or more geofences, it's currently not ideal. In this "remote" application scenario, the Client application itself is only required to determine geofence inclusion - but the remote application needs to know about it.

Fencer.io currently doesn't support this second style of App, where remote applications need to know about remote clients. 

But Froxy does.

### Data flow

1. Client application (Web/Native) makes request to Froxy rather than Fencer API
2. Froxy queries a geofence via Fencer API with client user's lat/lng coordinates and Fencer returns response to Froxy
3. Applications connected by Websocket are appraised of change events (inside, outside geofence)
4. Configured webhooks are triggered with payload

## Usage

Froxy is configured with a single `config.toml` file

Example `config.toml`

```toml
## Example Froxy configuration file

# Fencer.io API key
apikey = "34xx59-xxx-xxx-xxx-696xxx4010b9"

# Example geofence with websockets enabled and a single webhook to IFTTT
# Maker Webhooks service
[[geofence]]
  alias = "Home"
  accesskey = "3096eb87-xxxx-xxxx-xxxx-5dfxxxx25273" # Fencer geofence access key
  websocket = true
  webhooks = [
              [
                "https://maker.ifttt.com/trigger/hook/with/key/cLnxxxxxxxq1UpCW",
                "{ \"value1\" : \"{client_id}\", \"value2\" : \"{inside}\", \"value3\" : \"{lng_pos}\"}",
                ""
              ]
             ]

# Example geofence with only websockets enabled
[[geofence]]
  alias = "Work"
  accesskey = "9cfxxa37-da4a-4edd-xxxxx-xxxxx8f0" # Fencer geofence access key
  websocket = true

```
### REST (Client)

By default Froxy listens on port 9000.

A typical HTTP request from a web or native client app would look like the below. 

```html
http://hostdomain:9000/client/309xxx7-fcxx5-4xxb2-b1xxf-5dfxxxx5273
```

The long hash on the end of the URI is the access key of the geofence to be queried.

Client requests can contain up to three request headers. 

```
Lat-Pos   : Latitude of the client      (required)
Lng-Pos   : Longitude of the client     (required)
Client-ID : Identifier for user/client  (optional)
```

### Websockets

A single websocket connection is supported for each geofence. A websocket connection is established by making the following request. Events fire on change in inside/outside status.

```html
ws://hostdomain:9000/ws/309xxx7-fcxx5-4xxb2-b1xxf-5dfxxxx5273
```

A JSON object is provided with each event:

```json
{
  "client_id": "ollie",
  "geofence_alias": "Work",
  "lat_pos": "55.345239",
  "lng_pos": "-2.639349",
  "inside": true
}
```

### Webhooks

[work in progress]

## Notes
- There's no authentication in Froxy (yet)
- No TLS (yet)