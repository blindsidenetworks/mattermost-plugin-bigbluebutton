import {dismissOpenMeeting} from '../../actions';
import {getPluginState} from '../../selectors';

const {bindActionCreators} = window.Redux;
const {connect} = window.ReactRedux;
import OpenMeeting from './open_meeting.jsx';
import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/common';

function mapDispatchTooProps(dispatch) {
	return {
		actions: bindActionCreators({
			dismissOpenMeeting,
		}, dispatch)
	};
}

function mapStateToProps(state) {
	let siteURL = state.entities.general.config.SiteURL;
	siteURL = siteURL.endsWith('/') ? siteURL.substring(0, siteURL.length - 1) : siteURL;
	siteURL = siteURL.trim();

	return {
		siteURL,
		channelID: getCurrentChannelId(state),
		meeting: getPluginState(state).meeting || {}
	};
}

export default connect(mapStateToProps, mapDispatchTooProps)(OpenMeeting);
