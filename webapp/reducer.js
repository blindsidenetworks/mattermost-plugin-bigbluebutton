import {combineReducers} from 'redux';
import {
	CLOSE_ROOT_MODAL,
	DISMISS_INCOMING_CALL,
	INCOMING_CALL,
	OPEN_ROOT_MODAL,
	SET_PLUGIN_CONFIG,
	STATUS_CHANGE
} from './action_types';

const enabled = (state = false, action) => {
	switch (action.type) {
	case STATUS_CHANGE:
		return action.data;

	default:
		return state;
	}
};

const rootModalVisible = (state = false, action) => {
	switch (action.type) {
	case OPEN_ROOT_MODAL:
		return true;
	case CLOSE_ROOT_MODAL:
		return false;
	default:
		return state;
	}
};

const pluginConfig = (state = {}, action) => {
	switch (action.type) {
	case SET_PLUGIN_CONFIG:
		return action.data;
	default:
		return state;
	}
};

const incomingCall = (state = {}, action) => {
	switch (action.type) {
	case INCOMING_CALL:
		return {
			...action.data,
			dismissed: false,
		};
	case DISMISS_INCOMING_CALL:
		return {
			dismissed: true,
		};
	default:
		return state;
	}
};

export default combineReducers({
	enabled,
	rootModalVisible,
	pluginConfig,
	incomingCall,
});
