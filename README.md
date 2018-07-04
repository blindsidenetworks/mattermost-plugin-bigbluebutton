# BigBlueButton Plugin for Mattermost
BigBlueButton plugin for Mattermost. Teams can create, join and manage their BigBlueButton meetings from inside Mattermost.

## Installation and Setup

 1. Go to: https://github.com/ypgao1/mattermost-plugin-bigbluebutton/releases
 2. Download `bigbluebutton.tar.gz` you do not need to extract the tar file once you download it.![enter image description here](https://raw.githubusercontent.com/ypgao1/mattermost-plugin-bigbluebutton/master/docs_images/download_binary.png?token=AQcJwAEZnlU-0YcwkuRX5CIiis4L7ENbks5bRiAmwA==)
 3. Inside Mattermost, go to **System Console > Integrations > Custom Integrations**. Make sure the following are turned to true:
	- `Enable Incoming Webhooks`
	- `Enable Outgoing Webhooks`
	- `Enable Custom Slash Commands`
	- `Enable integrations to override usernames`
	- `Enable integrations to override profile picture icons`
 4. Next we must enable Plugins. Go to **System Console > Plugins > Configuration** and set `Enable Plugins` to true. ![enter image description here](https://raw.githubusercontent.com/ypgao1/mattermost-plugin-bigbluebutton/master/docs_images/enableplugins.png?token=AQcJwEJmN0uiifTscnFiiU48DWrluxuqks5bRiYKwA==)
 Depending on your Mattermost version, an additional step may be required to enable uploading plugins in your Mattermost **config.json** file: 
	 - `vi /opt/mattermost/config/config.json`
	 - Under `PluginSettings`, make sure `Enable` and `Enable Uploads` are both set to `true`

## Developing

Plugin is written in Golang for server side and Javascript for Client side. Use `make build` to build the plugin (generate the tar file).
The dependencies are managed with Glide for Go and NPM for javascript.

Mattermost plugin development guides available here: https://developers.mattermost.com/extend/plugins/

BigBlueButton API available here: http://docs.bigbluebutton.org/dev/api.html
