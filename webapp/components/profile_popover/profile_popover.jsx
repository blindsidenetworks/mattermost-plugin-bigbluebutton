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

import PropTypes from 'prop-types';
import * as ChannelActions from 'mattermost-redux/actions/channels';
import {Client4} from 'mattermost-redux/client';

export default class ProfilePopover extends React.PureComponent {
	static propTypes = {
		user: PropTypes.object.isRequired,
		cur_user: PropTypes.object.isRequired,
		state: PropTypes.object.isRequired,
		actions: PropTypes.shape({startMeeting: PropTypes.func.isRequired}).isRequired,
		teamname: PropTypes.string.isRequired,
		siteURL: PropTypes.string.isRequired,
	};

	constructor(props) {
		super(props);
	}

	handleDirectMessage = async (e, href) => {
		e.preventDefault();
		const dispatch = window.store.dispatch;
		Client4.setUrl(this.props.siteURL);
		const result = await ChannelActions.createDirectChannel(this.props.user.id, this.props.cur_user.id)(dispatch, this.props.state);
		await this.props.actions.startMeeting(result.data.id, '', this.props.cur_user.username + ' ' + this.props.user.username);
		window.location = href;
	};

	render() {
		const user = this.props.user;
		const myteam = this.props.teamname;
		const url = `${this.props.siteURL}/${myteam}/messages/@${user.username}`;

		if (this.props.user.id === this.props.cur_user.id) {
			return null;
		}

		return (
			<div>
				<hr style={{
					margin: '10px -15px 10px'
				}}/>
				<div>
					<a onClick={(e) => this.handleDirectMessage(e, url)}>
						<i className="fa fa-video-camera"/>{'  Start BigBlueButton Meeting'}
					</a>
					<br/>
				</div>
			</div>
		);
	}
}
