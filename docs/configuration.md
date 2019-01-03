Configuration path:
- **Windows**: $USER_HOME/.artemis/config.json
- **Unix**: /etc/artemis/config.json

Default configuration:
```json
{
	"broker": {
	  "host": "amqp://guest:guest@localhost",
	  "port": "5672",
	  "exchangeKey": "artemis",
	  "broadcastRoute": "broadcast.all"
	},
	"logging": {
	  "displayTimeStamp": true,
	  "debug": true
	},
	"heartbeatInterval": 1500,
	"electionTimeout": 2000,
	"clusterSize": 3
}
```
