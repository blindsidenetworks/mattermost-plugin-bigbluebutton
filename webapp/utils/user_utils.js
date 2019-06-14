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
import {getFullName} from 'mattermost-redux/utils/user_utils';

export function displayUsernameForUser(user, config) {
  if (user) {
    const nameFormat = config.TeammateNameDisplay;
    let name = user.username;
    if (nameFormat === 'nickname_full_name' && user.nickname && user.nickname !== '') {
      name = user.nickname;
    } else if ((user.first_name || user.last_name) && (nameFormat === 'nickname_full_name' || nameFormat === 'full_name')) {
      name = getFullName(user);
    }

    return name;
  }

  return '';
}
