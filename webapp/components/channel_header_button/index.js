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

const {connect} = window.ReactRedux;
import {getTheme} from 'mattermost-redux/selectors/entities/preferences';
import ChannelHeaderButton from './channel_header_button.jsx';
import {isRootModalVisible} from '../../selectors';

function mapStateToProps(state) {
	let channelId = state.entities.channels.currentChannelId;
	const channel = state.entities.channels.channels[channelId] || {};
	const userId = state.entities.users.currentUserId;
	if (channel.name === `${userId}__${userId}`) {
		channelId = '';
	}

	return {
		channelId,
		theme: getTheme(state),
		visible: isRootModalVisible(state),
	};
}

export default connect(mapStateToProps)(ChannelHeaderButton);
