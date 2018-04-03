# Froxy

## Fencer proxy with websocket relay and webhooks

For many apps the Fencer REST API is fine for use in a client/server manner. However, in other types of applications, external systems may benefit by being able to "act" on the knowledge (whether a user is inside or outside one or more geofences) and the client application itself is only required to determine position.

Fencer.io currently does not work well for this second style of "remote tracking" application (but may offer support in future).

Froxy is a proxy server for the Fencer API. It supports both websockets and webhooks for push notifications to external applications.

### Data flow

1. Client application (Web/App) makes request to Froxy rather than Fencer API
2. Froxy queries a geofence via Fencer API with client user's lat/lng coordinates and Fencer returns response to Froxy
3. Applications connected by Websocket are appraised of change events (inside, outside geofence)
4. Configured webhooks are triggered with payload

### Usage

Froxy is configured with a single `config.toml` file

Example `config.toml`

```toml
## Example Froxy configuration file

# Fencer.io API key
apikey 		= 	"34xx59-xxx-xxx-xxx-696xxx4010b9"

# Example geofence with websockets enabled and a single webhook to IFTTT
# Maker Webhooks service
[[geofence]]
	alias 		= 	"Home" 
	accesskey 	= 	"3096eb87-xxxx-xxxx-xxxx-5dfxxxx25273"
	websocket 	=	true
	webhooks 	= 	[
						[
							"https://maker.ifttt.com/trigger/hook/with/key/cLnxxxxxxxq1UpCW",
							"{ \"value1\" : \"{client_id}\", \"value2\" : \"{inside}\", \"value3\" : \"{lng_pos}\"}",
							""
						]
					]

# Example geofence with websockets enabled
[[geofence]]
	alias 		= 	"Work"
	accesskey 	= 	"9cfxxa37-da4a-4edd-xxxxx-xxxxx8f0"
	websocket 	= 	true

```

[work in progress]

## Notes
There's no authentication in Froxy (yet)
No TLS (yet)