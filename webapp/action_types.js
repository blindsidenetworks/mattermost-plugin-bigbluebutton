import PluginId from './plugin_id';

// Namespace your actions to avoid collisions.
export const STATUS_CHANGE = PluginId + '_status_change';
export const OPEN_ROOT_MODAL = PluginId + '_open_root_modal';
export const CLOSE_ROOT_MODAL = PluginId + '_close_root_modal';
export const SET_PLUGIN_CONFIG = PluginId + '_set_plugin_config';
export const INCOMING_CALL = PluginId + '_incoming_call';
export const DISMISS_INCOMING_CALL = PluginId + '_dismiss_incoming_call';
export const OPEN_MEETING = PluginId + '_open_meeting';
export const DISMISS_OPEN_MEETING = PluginId + '_dismiss_open_meeting';
