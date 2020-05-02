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

import ChannelHeaderButton from './components/channel_header_button';
import ProfilePopover from './components/profile_popover';
import PostTypebbb from './components/post_type_bbb';
import Root from './components/root';
import PluginId from './plugin_id';

import {channelHeaderButtonAction} from './actions';
import reducer from './reducer';
import {initClient} from "./client";

class PluginClass {
  initialize(registry, store) {
    window.store = store;

    const siteURL = store.getState().entities.general.config.SiteURL;
    initClient(siteURL);

    registry.registerPostTypeComponent('custom_bbb', PostTypebbb);
    registry.registerChannelHeaderButtonAction(
      <ChannelHeaderButton/>,() => store.dispatch(channelHeaderButtonAction()), 'BigBlueButton');
    registry.registerPopoverUserActionsComponent(ProfilePopover);
    registry.registerRootComponent(Root);
    registry.registerReducer(reducer);

  }
}

global.window.registerPlugin(PluginId, new PluginClass());
