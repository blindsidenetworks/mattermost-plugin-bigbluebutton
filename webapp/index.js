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

import React from 'react';
import ChannelHeaderButton from './components/channel_header_button';
import ProfilePopover from './components/profile_popover';
import Root from './components/root';
import PluginId from './plugin_id';

import {channelHeaderButtonAction} from './actions';
import reducer from './reducer';
import {GetClient, initClient} from './client';
import {INCOMING_CALL, SET_PLUGIN_CONFIG} from './action_types';
import IncomingCallPopup from './components/incoming_call';

class PluginClass {
	async initialize(registry, store) {
		window.store = store;
		initClient(() => store.getState().entities.general.config.SiteURL);

		registry.registerChannelHeaderButtonAction(
			<ChannelHeaderButton/>,
			() => store.dispatch(channelHeaderButtonAction()),
			'BigBlueButton',
		);
		registry.registerPopoverUserActionsComponent(ProfilePopover);
		registry.registerRootComponent(Root);
		registry.registerRootComponent(IncomingCallPopup);
		registry.registerReducer(reducer);

		registry.registerWebSocketEventHandler(
			`custom_${PluginId}_config_update`,
			(payload) => {
				store.dispatch({
					type: SET_PLUGIN_CONFIG,
					data: payload.data.config,
				});
			}
		);

		registry.registerWebSocketEventHandler(
			`custom_${PluginId}_incoming_call`,
			(payload) => {
				store.dispatch({
					type: INCOMING_CALL,
					data: payload.data,
				});
			}
		);

		await this.setPluginConfig(store);
	}

	async setPluginConfig(store) {
		const pluginConfig = await GetClient().getPluginConfig();
		store.dispatch({
			type: SET_PLUGIN_CONFIG,
			data: pluginConfig,
		});
	}
}

window.registerPlugin(PluginId, new PluginClass());
