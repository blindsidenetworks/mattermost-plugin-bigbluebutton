import {getChannelByNameAndTeamName} from 'mattermost-redux/actions/channels';

export function getDirectChannel(teamname, username) {
  if (username) {
    var channel = getChannelByNameAndTeamName(teamname, username);
    return channel;
  }
  return {};
}
