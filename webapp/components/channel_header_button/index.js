const {connect} = window['react-redux'];
const {bindActionCreators} = window.redux;

import {startMeeting, showRecordings} from '../../actions';
import {getChannelsInCurrentTeam, getDirectChannels, getSortedUnreadChannelIds, makeGetChannel} from 'mattermost-redux/selectors/entities/channels';

import ChannelHeaderButton from './channel_header_button.jsx';

function mapStateToProps(state, ownProps) {
  let channelId = state.entities.channels.currentChannelId;
  const channel = state.entities.channels.channels[channelId] || {};
  const userId = state.entities.users.currentUserId;
  if (channel.name === `${userId}__${userId}`) {
    channelId = '';
  }
  let teamId = state.entities.teams.currentTeamId;

  return {
    state,
    channelId,
    channel: channel,
    channelName: channel.name,
    directChannels: getDirectChannels(state),
    teamId,
    ...ownProps
  };
}

function mapDispatchToProps(dispatch) {
  return {
    actions: bindActionCreators({
      startMeeting,
      showRecordings
    }, dispatch)
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(ChannelHeaderButton);
