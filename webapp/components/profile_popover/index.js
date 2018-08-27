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
import {getCurrentTeam} from 'mattermost-redux/selectors/entities/teams';
import {getCurrentUser} from 'mattermost-redux/selectors/entities/users';
import ProfilePopover from './profile_popover.jsx';

import {startMeeting} from '../../actions';

function mapStateToProps(state, ownProps) {
  const post = ownProps.post || {};
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
  return {
    actions: bindActionCreators({
      startMeeting
    }, dispatch)
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(ProfilePopover);
