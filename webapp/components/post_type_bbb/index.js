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
const {bindActionCreators} = window.Redux;

import {getBool} from 'mattermost-redux/selectors/entities/preferences';
import {displayUsernameForUser} from '../../utils/user_utils';
import {
  getJoinURL,
  endMeeting,
  getAttendees,
  publishRecordings,
  deleteRecordings,
  isMeetingRunning
} from '../../actions';
import {getCurrentUserId} from 'mattermost-redux/selectors/entities/users';
import PostTypebbb from './post_type_bbb.jsx';

//custom post for users to join meetings, end meetings, view recordings, etc

function mapStateToProps(state, ownProps) {
  const post = ownProps.post || {};
  const user = state.entities.users.profiles[post.user_id] || {};
  let channelId = state.entities.channels.currentChannelId;
  const channel = state.entities.channels.channels[channelId]
  const userid = getCurrentUserId(state) || {}; 
  return {
    channelId,
    channel,
    state,
    ...ownProps,
    currentUserId: userid,
    creatorId: user.id,
    username: user.username,
    creatorName: displayUsernameForUser(user, state.entities.general.config),
    useMilitaryTime: getBool(state, 'display_settings', 'use_military_time', false)
  };
}

function mapDispatchToProps(dispatch) {
  return {
    actions: bindActionCreators({
      getJoinURL,
      endMeeting,
      getAttendees,
      publishRecordings,
      deleteRecordings,
      isMeetingRunning
    }, dispatch)
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(PostTypebbb);
