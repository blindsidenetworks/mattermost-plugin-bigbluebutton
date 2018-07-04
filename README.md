# mattermost-plugin-bigbluebutton
BigBlueButton plugin for Mattermost. Teams can create, join and manage their BigBlueButton meetings from inside Mattermost.

## Installation

Get the tar file from one of the releases. Install the plugin in System Console > Plugins > Management

## Developing

Plugin is written in Golang for server side and Javascript for Client side. Use `make build` to build the plugin (generate the tar file).
The dependencies are managed with Glide for Go and NPM for javascript.

Mattermost plugin development guides available here: https://developers.mattermost.com/extend/plugins/

BigBlueButton API available here: http://docs.bigbluebutton.org/dev/api.html
