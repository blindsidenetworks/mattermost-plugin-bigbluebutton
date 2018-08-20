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

// class PluginClass {
//   initialize(registerComponents, store) {
//     window.store = store;
//
//     registerComponents({
//       ChannelHeaderButton,
//       ProfilePopover,
//       Root
//     }, {custom_bbb: PostTypebbb});
//   }
// }

class PluginClass {
  initialize(registry, store) {
    window.store = store;

    registerComponents({
      ChannelHeaderButton,
      ProfilePopover,
      Root
    }, {custom_bbb: PostTypebbb});
  }
}

// global.window.plugins['bigbluebutton'] = new PluginClass();
window.registerPlugin('bigbluebutton', new MyPlugin());
