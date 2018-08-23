# BigBlueButton Plugin for Mattermost
BigBlueButton is an open source web conferencing system for online learning. Teams can create, join and manage their BigBlueButton meetings from inside Mattermost.

Jump to:

- [Installation and Setup](https://github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/blob/master/README.md#installation-and-setup)  
- [Usage](https://github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/blob/master/README.md#usage)
- [Contributing](https://github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/blob/master/README.md#contributing)

See this [YouTube video](https://youtu.be/gg7J9B4wGa4) for an overview of the Mattermost integration for BigBlueButton.

## Installation and Setup

 1. Go to: https://github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/releases
 2. Download `bigbluebutton.tar.gz` you do not need to extract the tar file once you download it.![enter image description here](https://raw.githubusercontent.com/blindsidenetworks/mattermost-plugin-bigbluebutton/master/docs_images/download_binary.png)
 3. Inside Mattermost, go to **System Console > Integrations > Custom Integrations**. Make sure the following are turned to true:
	- `Enable Incoming Webhooks`
	- `Enable Outgoing Webhooks`
	- `Enable Custom Slash Commands`
	- `Enable integrations to override usernames`
	- `Enable integrations to override profile picture icons`
 4. Next we must enable Plugins. Go to **System Console > Plugins > Configuration** and set `Enable Plugins` to true. ![enter image description here](https://raw.githubusercontent.com/blindsidenetworks/mattermost-plugin-bigbluebutton/master/docs_images/enableplugins.png)
 Depending on your Mattermost version, an additional step may be required to enable uploading plugins in your Mattermost **config.json** file:
	 - `vi /opt/mattermost/config/config.json`
	 - Under `PluginSettings`, make sure `Enable` and `Enable Uploads` are both set to `true`
	 - Restart your Mattermost with `sudo systemctl restart mattermost` assuming you used *systemd* for Mattermost 	services
 5. Go to **System Console > Plugins > Management** and upload your `bigbluebutton.tar.gz`. The BigBlueButton Plugin should appear under **Installed Plugins**.    ![
](https://raw.githubusercontent.com/blindsidenetworks/mattermost-plugin-bigbluebutton/master/docs_images/PluginManagement.png)
 6. Before activating the plugin, we must go configure the plugin settings. By default, you are given a BigBlueButton test server to try it out. See [BigBlueButton documentation](http://docs.bigbluebutton.org/install/install.html#Install_) to install your own BigBlueButton server. The secret key is to securely connect to BigBlueButton. To check your secret, in the command line, enter `bbb-conf --secret`.
Alternatively contact **Blindside Networks**, the company behind BigBlueButton, for  [Setup & Support, Custom Development, and Premium Hosting.](https://blindsidenetworks.com/services/)

	The **Site URL** is the site of your Mattermost without any paths. For example, if the location of your Mattermost Town Square is : `https://mysite.mattermost.com/core/channels/town-square`, enter: `https://mysite.mattermost.com`![
](https://raw.githubusercontent.com/blindsidenetworks/mattermost-plugin-bigbluebutton/master/docs_images/BBBsettingspage.png)

 7. Afterwards, go back to **System Console > Plugins > Management** and `Activate` the plugin. ![](https://raw.githubusercontent.com/blindsidenetworks/mattermost-plugin-bigbluebutton/master/docs_images/activate_plugin.png)


## Usage
#### Create a BigBlueButton meeting in any channel
![
](https://raw.githubusercontent.com/blindsidenetworks/mattermost-plugin-bigbluebutton/master/docs_images/createchannelheader.png)

#### Users can join BigBlueButton meetings through the post message
![enter image description here](https://raw.githubusercontent.com/blindsidenetworks/mattermost-plugin-bigbluebutton/master/docs_images/insideBBB.png)

#### Plugin provides live meeting details during and after the meeting has ended
![
](https://raw.githubusercontent.com/blindsidenetworks/mattermost-plugin-bigbluebutton/master/docs_images/recordingmanagment.png)

#### You can search for past BigBlueButton recordings
![
](https://raw.githubusercontent.com/blindsidenetworks/mattermost-plugin-bigbluebutton/master/docs_images/view_recordings.png)

#### Alternative way to start a BigBlueButton meeting is through clicking on a user's name and getting their profile popover
![
](https://raw.githubusercontent.com/blindsidenetworks/mattermost-plugin-bigbluebutton/master/docs_images/popover.png)

#### Slash command `/bbb` can also be used to start a meeting
![
](https://raw.githubusercontent.com/blindsidenetworks/mattermost-plugin-bigbluebutton/master/docs_images/slashcommand.png)

#### For any direct or group message, popup alert will open anywhere inside Mattermost to notify that someone has requested a meeting with you.
![
](https://raw.githubusercontent.com/blindsidenetworks/mattermost-plugin-bigbluebutton/master/docs_images/popup_modal.png)

## Contributing

Plugin is written in Golang for server side and Javascript for client side. Use `make build` to build the plugin.
The dependencies are managed with Glide for Go and NPM for javascript.

The plugin should be placed in a directory such as `~/go/src/github.com/blindsidenetworks/mattermost-plugin-bigbluebutton`

To download a local version: `mkdir -p ~/go/src/github.com/blindsidenetworks` and `git clone https://github.com/blindsidenetworks/mattermost-plugin-bigbluebutton.git`

Mattermost plugin development guides available here: https://developers.mattermost.com/extend/plugins/

BigBlueButton API available here: http://docs.bigbluebutton.org/dev/api.html
