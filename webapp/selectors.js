import PluginId from './plugin_id';

export const getPluginState = (state) => state['plugins-' + PluginId] || {};

export const isEnabled = (state) => getPluginState(state).enabled;

export const isRootModalVisible = (state) => getPluginState(state).rootModalVisible;

export const getPluginConfig = (state) => getPluginState(state).pluginConfig;
