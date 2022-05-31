import {dismissOpenMeeting} from '../../actions';
import {getPluginState} from '../../selectors';

const {bindActionCreators} = window.Redux;
const {connect} = window.ReactRedux;
import OpenMeeting from './open_meeting.jsx';
import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/common';
import utils from '../../utils/utils';

function mapDispatchTooProps(dispatch) {
	return {
		actions: bindActionCreators({
			dismissOpenMeeting,
		}, dispatch)
	};
}

function mapStateToProps(state) {
	let siteURL = utils.cleanSiteURL(state.entities.general.config.SiteURL);

	return {
		siteURL,
		channelID: getCurrentChannelId(state),
		meeting: getPluginState(state).meeting || {}
	};
}

export default connect(mapStateToProps, mapDispatchTooProps)(OpenMeeting);
