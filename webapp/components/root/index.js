const {connect} = window['react-redux'];
const {bindActionCreators} = window.redux;
import {getCurrentTeam} from 'mattermost-redux/selectors/entities/teams';
import {getCurrentUser} from 'mattermost-redux/selectors/entities/users';
import {getLastPostPerChannel} from 'mattermost-redux/selectors/entities/posts';
import {Client4} from 'mattermost-redux/client';


import {
    getSortedPublicChannelWithUnreadsIds,
    getSortedPrivateChannelWithUnreadsIds,
    getSortedFavoriteChannelWithUnreadsIds,
    getSortedDirectChannelWithUnreadsIds,
    getCurrentChannel,
    getUnreads,
    getSortedUnreadChannelIds,
    getSortedDirectChannelIds,
    getSortedFavoriteChannelIds,
    getSortedPublicChannelIds,
    getSortedPrivateChannelIds,
} from 'mattermost-redux/selectors/entities/channels';
import {getJoinURL } from '../../actions';

import Root from './root.jsx';





function mapStateToProps(state, ownProps) {
  const post = ownProps.post || {};
    /* Provide values for any custom props or override any existing props here */
    let team = getCurrentTeam(state) || {};
    let teamname = team.name;
    let cur_user = getCurrentUser(state) || {};
      const keepChannelIdAsUnread = state.views.channel.keepChannelIdAsUnread;
    //console.log("getting all posts:" +state.entities.posts.posts)
  //  console.log("what does get unread print? " + JSON.stringify(getUnreads(state)));
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
