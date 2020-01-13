# Note
We have updated the Mattermost plugin. Only versions 2.1.0 and above will work with the default BigBlueButton server credentials.

Older versions will, however, continue to work when configured with your own BigBlueButton server.

# BigBlueButton Plugin for Mattermost
BigBlueButton is an open source web conferencing system for online learning. Teams can create, join and manage their BigBlueButton meetings from inside Mattermost.

Jump to:

- [Installation and Setup](#installation-and-setup)  
- [Usage](#usage)
- [Contributing](#contributing)

Want to see how the BigBlueButton integration with Mattermost works?  Checkout the video below.

[![Alt text](https://img.youtube.com/vi/gg7J9B4wGa4/0.jpg)](https://www.youtube.com/watch?v=gg7J9B4wGa4)

## Installation and Setup

 1. Go to: https://github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/releases
 2. Download the `tar.gz` file corresponding to your server platform. You do not need to extract the tar file once you download it.![enter image description here](docs_images/download_binary.png)
 3. Inside Mattermost, go to **System Console > Integrations > Integration Management**. Make sure the following are turned to true:
	- `Enable Incoming Webhooks`
	- `Enable Outgoing Webhooks`
	- `Enable Custom Slash Commands`
	- `Enable integrations to override usernames`
	- `Enable integrations to override profile picture icons`
 4. Next you must enable Plugins. Go to **System Console > Plugin Management > Configuration** and set `Enable Plugins` to true. ![enter image description here](docs_images/enableplugins.png)
 Depending on your Mattermost version, an additional step may be required to enable uploading plugins in your Mattermost **config.json** file:
	 - `vi /opt/mattermost/config/config.json`
	 - Under `PluginSettings`, make sure `Enable` and `Enable Uploads` are both set to `true`
	 - Restart your Mattermost with `sudo systemctl restart mattermost` assuming you used *systemd* for Mattermost services
 5. Go to **System Console > Plugins > Management** and upload your `bigbluebutton_<arch>_amd64.tar.gz`. The BigBlueButton Plugin should appear under **Installed Plugins**.    ![
](docs_images/PluginManagement.png)
 6. Before activating the plugin, you must configure the plugin settings. By default, you are given a BigBlueButton test server to try it out. However, you have options.  Like Mattermost, BigBlueButton is open source.  You are (more than) welcome to [setup your own BigBlueButton server](http://docs.bigbluebutton.org/install/install.html#Install_).  If you do, the command `sudo bbb-conf --secret` will print out the server's URL and secret key for configuration with Mattermost.  Alternatively, you can [contact](https://blindsidenetworks.com/contact/) Blindside Networks for [hosting options](https://blindsidenetworks.com/services/).

	![](docs_images/BBBsettingspage.png)

 7. Next, go back to **System Console > Plugins > Management** and `Activate` the plugin. ![](docs_images/activate_plugin.png)


## Usage

#### Create a BigBlueButton meeting in any channel

You can create a meeting that all channel participants can join.

![](docs_images/createchannelheader.png)

Clicking the **Join Meeting** button immediately loads the BigBlueButton HTML5 client.

![enter image description here](docs_images/insideBBB.png)

#### Plugin provides live meeting details during and after the meeting has ended

After the meeting ends, you see the **Date**, **Meeting Length**, and **Attendees**.

![](docs_images/recordingmanagment.png)

#### You can search for past BigBlueButton recordings

Using the drop-down menu you can easily search a channel for all past recordings.

![](docs_images/view_recordings.png)

#### Directly meeting with any user

You can click on any user's name and choose **Start BigBlueButton Meeting**.

![](docs_images/popover.png)

When you invite a user to a meeting, they will get a pop-up notification to **Join Meeting**.

![](docs_images/popup_modal.png)

You can type `/bbb` in any channel to create a meeting.  When 

![](docs_images/slashcommand.png)

## Setting up your own BigBlueButton server

Using the [bbb-install.sh](https://github.com/bigbluebutton/bbb-install) script you can setup your own BigBlueButton server in about 15 minutes.  If your interested in going through the steps in detail, see [BigBlueButton install guide](http://docs.bigbluebutton.org/install/install.html).

## Contributing

Plugin is written in Golang for server side and Javascript and React for client side. Use `make build` to build the plugin. You can also use `make quickbuild` following first build for faster builds.
The dependencies are managed with Glide for Go and NPM for javascript.

The plugin should be placed in a directory such as `~/go/src/github.com/blindsidenetworks/mattermost-plugin-bigbluebutton`

To download a local version: `mkdir -p ~/go/src/github.com/blindsidenetworks` and `git clone https://github.com/blindsidenetworks/mattermost-plugin-bigbluebutton.git`

Mattermost plugin development guides available here: https://developers.mattermost.com/extend/plugins/

BigBlueButton API available here: http://docs.bigbluebutton.org/dev/api.html
