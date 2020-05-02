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

import {PostTypes} from 'mattermost-redux/action_types';
import {batchActions} from 'redux-batched-actions';
import {getCurrentChannel} from 'mattermost-redux/selectors/entities/channels';
import {getCurrentTeamId} from 'mattermost-redux/selectors/entities/teams';
import {searchPosts} from 'mattermost-redux/actions/search';

import {ActionTypes, RHSStates} from '../utils/constants.jsx';
import PluginId from '../plugin_id';
import {STATUS_CHANGE, OPEN_ROOT_MODAL, CLOSE_ROOT_MODAL} from '../action_types';
import {GetClient} from "../client";

export const openRootModal = () => (dispatch) => {
    dispatch({
        type: OPEN_ROOT_MODAL,
    });
};

export const closeRootModal = () => (dispatch) => {
    dispatch({
        type: CLOSE_ROOT_MODAL,
    });
};

export const mainMenuAction = openRootModal;
export const channelHeaderButtonAction = openRootModal;

export function startMeeting(channelId, description = '', topic = '', meetingId = 0) {
  return async (dispatch, getState) => {
    try {
      await GetClient.startMeeting(getState().entities.users.currentUserId, channelId, topic, description);
    } catch (error) {
      var message_text = 'BigBlueButton did not successfully start a meeting';
      if (error.status == 422 ) { // SiteURL is not set
         message_text =error.response.text;
      }
      const post = {
        id: 'bbbPlugin' + Date.now(),
        create_at: Date.now(),
        update_at: 0,
        edit_at: 0,
        delete_at: 0,
        is_pinned: false,
        user_id: getState().entities.users.currentUserId,
        channel_id: channelId,
        root_id: '',
        parent_id: '',
        original_id: '',
        message: message_text,
        type: 'system_ephemeral',
        props: {},
        hashtags: '',
        pending_post_id: ''
      };

      dispatch({
        type: PostTypes.RECEIVED_POSTS,
        data: {
          order: [],
          posts: {
            [post.id]: post
          }
        },
        channelId
      });

      return {error};
    }

    return {data: true};
  };
}

// Get join url:
export function getJoinURL(channelId, meetingid, creatorid) {
  return async (dispatch, getState) => {
    let url;
    var curUserId = getState().entities.users.currentUserId;
    var ismod = "FALSE"
    if (curUserId == creatorid) {
      ismod = "TRUE"
    }
    try {
      url = await GetClient().getJoinURL(curUserId, meetingid, ismod);
      return {
        data: {
          joinurl: url
        }
      };
    } catch (error) {
      const post = {
        id: 'bbbPlugin' + Date.now(),
        create_at: Date.now(),
        update_at: 0,
        edit_at: 0,
        delete_at: 0,
        is_pinned: false,
        user_id: getState().entities.users.currentUserId,
        channel_id: channelId,
        root_id: '',
        parent_id: '',
        original_id: '',
        message: 'Cant get a join url',
        type: 'system_ephemeral',
        props: {},
        hashtags: '',
        pending_post_id: ''
      };
      dispatch({
        type: PostTypes.RECEIVED_POSTS,
        data: {
          order: [],
          posts: {
            [post.id]: post
          }
        },
        channelId
      });
      return {error};
    }
  };
}
export function isMeetingRunning(meetingid) {
  return async (dispatch, getState) => {
    let response;
    try {
      response = await GetClient().isMeetingRunning(meetingid);
      return response
    } catch (error) {
      return {error};
    }
  };
}

export function endMeeting(channelId, meetingid) {
  return async (dispatch, getState) => {
    let url;
    try {
      url = await GetClient().endMeeting(getState().entities.users.currentUserId, meetingid);
      return {
        data: {
          joinurl: url
        }
      };
    } catch (error) {
      return {error};
    }
  };
}
export function getAttendees(channelId, meetingid) {
  return async (dispatch, getState) => {

    try {
      var resp = await GetClient().getAttendees(meetingid);
      return resp;
    } catch (error) {
      const post = {
        id: 'bbbPlugin' + Date.now(),
        create_at: Date.now(),
        update_at: 0,
        edit_at: 0,
        delete_at: 0,
        is_pinned: false,
        user_id: getState().entities.users.currentUserId,
        channel_id: channelId,
        root_id: '',
        parent_id: '',
        original_id: '',
        message: 'Cant get attendees info',
        type: 'system_ephemeral',
        props: {},
        hashtags: '',
        pending_post_id: ''
      };
      dispatch({
        type: PostTypes.RECEIVED_POSTS,
        data: {
          order: [],
          posts: {
            [post.id]: post
          }
        },
        channelId
      });
      return {error};
    }
  };
}

export function publishRecordings(channelId, recordid, publish, meetingId) {
  console.log(recordid + " " + publish)
  return async (dispatch, getState) => {

    try {
      var resp = await GetClient().publishRecordings(recordid, publish, meetingId);
      return resp;
    } catch (error) {
      const post = {
        id: 'bbbPlugin' + Date.now(),
        create_at: Date.now(),
        update_at: 0,
        edit_at: 0,
        delete_at: 0,
        is_pinned: false,
        user_id: getState().entities.users.currentUserId,
        channel_id: channelId,
        root_id: '',
        parent_id: '',
        original_id: '',
        message: error.response.text ,
        type: 'system_ephemeral',
        props: {},
        hashtags: '',
        pending_post_id: ''
      };
      dispatch({
        type: PostTypes.RECEIVED_POSTS,
        data: {
          order: [],
          posts: {
            [post.id]: post
          }
        },
        channelId
      });
      return {error};
    }
  };
}

export function deleteRecordings(channelId, recordid, meetingId) {
  return async (dispatch, getState) => {

    try {
      var resp = await GetClient().deleteRecordings(recordid, meetingId);
      return resp;
    } catch (error) {
      const post = {
        id: 'bbbPlugin' + Date.now(),
        create_at: Date.now(),
        update_at: 0,
        edit_at: 0,
        delete_at: 0,
        is_pinned: false,
        user_id: getState().entities.users.currentUserId,
        channel_id: channelId,
        root_id: '',
        parent_id: '',
        original_id: '',
        message: error.response.text,
        type: 'system_ephemeral',
        props: {},
        hashtags: '',
        pending_post_id: ''
      };
      dispatch({
        type: PostTypes.RECEIVED_POSTS,
        data: {
          order: [],
          posts: {
            [post.id]: post
          }
        },
        channelId
      });
      return {error};
    }
  };
}

export function showRecordings() {
  return(dispatch, getState) => {

    const channel = getCurrentChannel(getState());

    const terms = "in:" + channel.name + " #recording"

    dispatch(performSearch(terms, false));
    dispatch(batchActions([
      {
        type: ActionTypes.UPDATE_RHS_SEARCH_TERMS,
        terms
      }, {
        type: ActionTypes.UPDATE_RHS_STATE,
        state: RHSStates.MENTION
      }
    ]));
  };
}
export function performSearch(terms, isMentionSearch) {
  return(dispatch, getState) => {
    const teamId = getCurrentTeamId(getState());

    return dispatch(searchPosts(teamId, terms, isMentionSearch));
  };
}
