const {connect} = window['react-redux'];
const {bindActionCreators} = window.redux;
import {getCurrentTeam} from 'mattermost-redux/selectors/entities/teams';
import {getCurrentUser} from 'mattermost-redux/selectors/entities/users';
import ProfilePopover from './profile_popover.jsx';


import {startMeeting} from '../../actions';


function mapStateToProps(state, ownProps) {
  const post = ownProps.post || {};
    /* Provide values for any custom props or override any existing props here */
    let team = getCurrentTeam(state) || {};
    let teamname = team.name;
    let cur_user = getCurrentUser(state) || {};
    return {
      state,
      cur_user,
      teamname,
        ...ownProps
    };
}

function mapDispatchToProps(dispatch) {
    /* Provide actions here if needed */
    return {
        actions: bindActionCreators({
          startMeeting
        }, dispatch)
    };
}

export default connect(mapStateToProps, mapDispatchToProps)(ProfilePopover);
