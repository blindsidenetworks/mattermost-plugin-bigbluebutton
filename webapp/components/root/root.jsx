/*
Copyright 2018 Blindside Networks

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import React from 'react';
import PopoverListMembersItem from './popover_list_members_item.jsx';
import PropTypes from 'prop-types';
import {makeStyleFromTheme} from 'mattermost-redux/utils/theme_utils';
import {getChannel} from 'mattermost-redux/selectors/entities/channels';

// eslint-disable-next-line no-unused-vars
const {Tooltip, Popover, OverlayTrigger, Modal, Overlay} = window.ReactBootstrap;

export default class Root extends React.PureComponent {
	static propTypes = {
		cur_user: PropTypes.object.isRequired,
		state: PropTypes.object.isRequired,
		teamname: PropTypes.string.isRequired,
		lastpostperchannel: PropTypes.object.isRequired,
		unreadChannelIds: PropTypes.array.isRequired,
		theme: PropTypes.object.isRequired,
		channelName: PropTypes.string.isRequired,
		channel: PropTypes.object.isRequired,
		pluginConfig: PropTypes.object.isRequired,
		channelId: PropTypes.string,
		visible: PropTypes.bool,
		actions: PropTypes.shape({
			getJoinURL: PropTypes.func.isRequired,
			channelId: PropTypes.string.isRequired,
			directChannels: PropTypes.array.isRequired,
			teamId: PropTypes.string.isRequired,
			visible: PropTypes.bool.isRequired,
			actions: PropTypes.shape({
				startMeeting: PropTypes.func.isRequired,
				showRecordings: PropTypes.func.isRequired,
				closePopover: PropTypes.func.isRequired
			}).isRequired,
			startMeeting: PropTypes.func.isRequired,
			showRecordings: PropTypes.func.isRequired,
			closePopover: PropTypes.func.isRequired,
		}).isRequired
	};

	constructor(props) {
		super(props);

		this.state = {
			ignoredPosts: [],
			show: false,
			channelId: '',
			channelName: '',
			meetingId: '',
			profilePicUrl: '',
			channelURL: ''
		};
	}

	handleClose = () => {
		this.setState({show: false});
	};

	searchRecordings = () => {
		this.props.actions.showRecordings();
	};

	startMeeting = async (allowRecording) => {
		await this.props.actions.startMeeting(this.props.channelId, allowRecording, '', this.props.channel.display_name);
		this.close_the_popover();
	};

	close_the_popover = () => {
		this.props.actions.closePopover();
		this.setState({showPopover: false});
	};

	componentDidUpdate(prevProps) {
		if (!this.props.pluginConfig.ALLOW_RECORDINGS && !prevProps.visible && this.props.visible) {
			this.startMeeting(false);
		}
	}

	render() {
		const pos_width = (window.innerWidth - 400 + 'px');
		const style = getStyle(pos_width, this.props.theme);

		let popoverButton = this.props.pluginConfig.ALLOW_RECORDINGS ? (
			<div className="more-modal__button" style={style.viewRecordingBtn}>
				<a className="btn  btn-link" onClick={this.searchRecordings}>
					{'View Recordings'}
				</a>
			</div>
		) : null;

		style.popover['marginLeft'] = pos_width;
		style.popoverDM['marginLeft'] = pos_width;

		const channel = getChannel(this.props.state, this.props.channelId);
		let channelName = '';
		let ownChannel = false;
		if (channel == undefined) {
			ownChannel = true;
		} else {
			channelName = channel.display_name;
		}

		const directMessageListItem = (
			<PopoverListMembersItem
				ariaLabel={'Call ' + channelName}
				onItemClick={() => this.startMeeting(true)}
				icon={'BBBCAM'}
				text={
					<span>{'Call '}<strong>{channelName}</strong></span>
				}
				theme={this.props.theme}
			/>
		);

		const allowRecordingBtn = (
			<PopoverListMembersItem
				ariaLabel={'Create a BigBlueButton Meeting'}
				onItemClick={() => this.startMeeting(true)}
				icon={'ALLOW_RECORDING'}
				text={
					<React.Fragment>
						<span> {'Start New Meeting'}</span>
						<br/>
						<span> {'Allow Recording'}</span>
					</React.Fragment>
				}
				theme={this.props.theme}
			/>
		);

		const noRecordingBtn = (
			<PopoverListMembersItem
				ariaLabel={'Create a BigBlueButton Meeting'}
				onItemClick={() => this.startMeeting(false)}
				icon={'DONT_ALLOW_RECORDING'}
				text={
					<React.Fragment>
						<span> {'Start New Meeting'}</span>
						<br/>
						<span> {'Recording Disabled'}</span>
					</React.Fragment>
				}
				theme={this.props.theme}
			/>
		);

		const channelListItem = (
			<React.Fragment>
				{allowRecordingBtn}
				{noRecordingBtn}
			</React.Fragment>
		);

		if (!this.props.pluginConfig.ALLOW_RECORDINGS) {
			return null;
		}

		return (
			<div>
				{!ownChannel && //shows popup when not on own channel
				<Overlay rootClose={true} show={this.props.visible} onHide={this.close_the_popover} placement="bottom">
					<Popover
						id="bbbPopover"
						style={this.props.channel.type === 'D' ? style.popoverDM : style.popover}>
						<div style={this.props.channel.type === 'D' ? style.popoverBodyDM : style.popoverBody}>
							{
								this.props.channel.type === 'D' ? directMessageListItem : channelListItem
							}
						</div>
						{popoverButton}
					</Popover>
				</Overlay>
				}
			</div>);
	}
}

/* Define CSS styles here */
var getStyle = makeStyleFromTheme((theme) => {
	var x_pos = (window.innerWidth - 400 + 'px'); //shouldn't be set here as it doesn't rerender
	return {
		popover: {
			marginLeft: x_pos,
			marginTop: '50px',
			maxWidth: '250px',
			width: '250px',
			background: theme.centerChannelBg,
			borderRadius: '12px',
			overflow: 'hidden',
			paddingTop: '12px',
		},
		popoverDM: {
			marginLeft: x_pos,
			marginTop: '50px',
			maxWidth: '220px',
			height: '105px',
			width: '220px',
			background: theme.centerChannelBg
		},
		header: {
			background: '#FFFFFF',
			color: '#0059A5',
			borderStyle: 'none',
			height: '10px',
			minHeight: '28px'
		},
		body: {
			padding: '0px 0px 10px 0px'
		},
		bodyText: {
			textAlign: 'center',
			margin: '20px 0 0 0',
			fontSize: '17px',
			lineHeight: '19px'
		},
		meetingId: {
			marginTop: '55px'
		},
		backdrop: {
			position: 'absolute',
			display: 'flex',
			top: 0,
			left: 0,
			right: 0,
			bottom: 0,
			backgroundColor: 'rgba(0, 0, 0, 0.50)',
			zIndex: 2000,
			alignItems: 'center',
			justifyContent: 'center',
		},
		modal: {
			height: '250px',
			width: '400px',
			padding: '1em',
			color: theme.centerChannelColor,
			backgroundColor: theme.centerChannelBg,
		},
		iconStyle: {
			position: 'relative',
			top: '-1px'
		},

		popoverBody: {
			maxHeight: '305px',
			overflow: 'auto',
			position: 'relative',
			width: '298px',
			left: '-14px',
			top: '-9px',
		},

		popoverBodyDM: {
			maxHeight: '305px',
			overflow: 'auto',
			position: 'relative',
			width: '218px',
			left: '-14px',
			top: '-9px',
		},
		viewRecordingBtn: {
			borderTop: '1px solid',
			borderColor: theme.centerChannelColor,
			padding: '8px 0px 0px',
		}
	};
});
