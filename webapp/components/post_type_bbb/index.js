const {connect} = window['react-redux'];
const {bindActionCreators} = window.redux;


import {getBool} from 'mattermost-redux/selectors/entities/preferences';
import {displayUsernameForUser} from '../../utils/user_utils';
import {getJoinURL,endMeeting,getAttendees,publishRecordings,deleteRecordings, isMeetingRunning } from '../../actions';
import {getCurrentUserId} from 'mattermost-redux/selectors/entities/users';
import PostTypebbb from './post_type_bbb.jsx';

//custom post for users to join meetings, end meetings, view recordings, etc

function mapStateToProps(state, ownProps) {
    const post = ownProps.post || {};
    const user = state.entities.users.profiles[post.user_id] || {};
      let channelId = state.entities.channels.currentChannelId;
      const channel = state.entities.channels.channels[channelId]
    const userid = getCurrentUserId(state) || {}; //dont know if we should have this here
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
        actions: bindActionCreators({getJoinURL,endMeeting,getAttendees,publishRecordings,deleteRecordings,isMeetingRunning
        }, dispatch)
    };
}

export default connect(mapStateToProps, mapDispatchToProps)(PostTypebbb);
