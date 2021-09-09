import {combineReducers} from 'redux';
import {STATUS_CHANGE, OPEN_ROOT_MODAL, CLOSE_ROOT_MODAL, SET_PLUGIN_CONFIG} from './action_types';

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
}

export default combineReducers({
    enabled,
    rootModalVisible,
    pluginConfig,
});
