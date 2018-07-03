// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

import ChannelHeaderButton from './components/channel_header_button';
import ProfilePopover from './components/profile_popover';
import PostTypebbb from './components/post_type_bbb';
import Root from './components/root';


class PluginClass {
    initialize(registerComponents, store) {
        window.store = store;

        registerComponents({ChannelHeaderButton,ProfilePopover,Root}, {custom_bbb: PostTypebbb});
    }
}

global.window.plugins['bigbluebutton'] = new PluginClass();
