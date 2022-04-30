import PluginId from './plugin_id';

export const getPluginState = (state) => state['plugins-' + PluginId] || {};

export const isRootModalVisible = (state) => getPluginState(state).rootModalVisible || false;
