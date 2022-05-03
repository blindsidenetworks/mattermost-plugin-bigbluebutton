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

const MATTERMOST_PRODUCT_PREFIXES = [
	'boards',
	'playbooks',
];

export default class OpenMeeting extends React.Component {
	// TODO if we need to check for userID match as well
	static propTypes = {
		siteURL: PropTypes.string.isRequired,
		channelID: PropTypes.string.isRequired,
		meeting: PropTypes.shape({
			channelID: PropTypes.string.isRequired,
			joinURL: PropTypes.string.isRequired,
		}),
		actions: PropTypes.shape({
			dismissOpenMeeting: PropTypes.func.isRequired,
		}),
	}

	constructor(props) {
		super(props);
	}

	isInProduct() {
		const currentURL = window.location.href;
		return MATTERMOST_PRODUCT_PREFIXES
			.find(
				(productPrefix) => currentURL.startsWith(this.props.siteURL + '/' + productPrefix)
			) || false;
	}

	componentDidUpdate(prevProps) {
		if (prevProps.meeting.joinURL === this.props.meeting.joinURL) {
			return;
		}

		const isTargetChannel = this.props.channelID === this.props.meeting.channelID;
		const isInProduct = this.isInProduct();

		if (!isTargetChannel || isInProduct) {
			this.props.actions.dismissOpenMeeting();
			return;
		}

		window.open(this.props.meeting.joinURL, '_blank');
		this.props.actions.dismissOpenMeeting();
	}

	render() {
		return null;
	}
}
