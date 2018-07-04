const {connect} = window['react-redux'];
const {bindActionCreators} = window.redux;
import {getCurrentTeam} from 'mattermost-redux/selectors/entities/teams';
import {getCurrentUser} from 'mattermost-redux/selectors/entities/users';
import {getLastPostPerChannel} from 'mattermost-redux/selectors/entities/posts';

import {getSortedDirectChannelWithUnreadsIds} from 'mattermost-redux/selectors/entities/channels';
import {getJoinURL } from '../../actions';

import Root from './root.jsx';

//Root component is used for creating a popup notifying user about a
// BigBlueButton meeting started from a direct message

function mapStateToProps(state, ownProps) {
  const post = ownProps.post || {};
    /* Provide values for any custom props or override any existing props here */
    let team = getCurrentTeam(state) || {};
    let teamname = team.name;
    let cur_user = getCurrentUser(state) || {};
      const keepChannelIdAsUnread = state.views.channel.keepChannelIdAsUnread;
    return {
      cur_user,
      teamname,
      state,
      lastpostperchannel: getLastPostPerChannel(state),
      unreadChannelIds: getSortedDirectChannelWithUnreadsIds(state, keepChannelIdAsUnread),
        ...ownProps,
    };
}

function mapDispatchToProps(dispatch) {
    /* Provide actions here if needed */
    return {
        actions: bindActionCreators({getJoinURL,
        }, dispatch)
    };
}

export default connect(mapStateToProps, mapDispatchToProps)(Root);
