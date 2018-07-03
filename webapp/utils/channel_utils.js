import {getChannelByNameAndTeamName} from 'mattermost-redux/actions/channels';

export function getDirectChannel(teamname, username){
  console.log("teamname: " +teamname)
  console.log("username: " + username)
  if (username){
    channel = getChannelByNameAndTeamName(teamname,username);
    return channel;
  }return {};
}
