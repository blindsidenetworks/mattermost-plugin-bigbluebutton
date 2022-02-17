import {dismissIncomingCall, getJoinURL} from '../../actions';
import {getPluginState} from '../../selectors';
import IncomingCallPopup from './incoming_call.jsx';

const {bindActionCreators} = window.Redux;
const {connect} = window.ReactRedux;
import {getTheme} from 'mattermost-redux/selectors/entities/preferences';
import {getTeam, getCurrentTeamId} from 'mattermost-redux/selectors/entities/teams';

function mapDispatchTooProps(dispatch) {
	return {
		actions: bindActionCreators({
			dismissIncomingCall,
			getJoinURL,
		}, dispatch)
	};
}

function mapStateToProps(state, ownProps) {
	const pluginState = getPluginState(state);
	return {
		theme: getTheme(state),
		siteURL: state.entities.general.config.SiteURL,
		currentTeam: getTeam(state, getCurrentTeamId(state)),
		incomingCall: getPluginState(state).incomingCall,
		...ownProps,
	};
}

export default connect(mapStateToProps, mapDispatchTooProps)(IncomingCallPopup);
