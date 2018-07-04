import {PostTypes} from 'mattermost-redux/action_types';
import {getCurrentTeamId} from 'mattermost-redux/selectors/entities/teams';
import {searchPosts} from 'mattermost-redux/actions/search';
import {batchActions} from 'redux-batched-actions';
import {ActionTypes, RHSStates} from '../utils/constants.jsx';
import {getCurrentChannel} from 'mattermost-redux/selectors/entities/channels';

import Client from '../client';


export function startMeeting(channelId, description = '', topic = '', meetingId = 0) {
    return async (dispatch, getState) => {
        try {
            await Client.startMeeting(getState().entities.users.currentUserId, channelId,topic,description);
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
                message: 'BigBlueButton did not successfully start a meeting',
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
export function getJoinURL(channelId, meetingid,creatorid) {
    return async (dispatch, getState) => {
      let url;
      var curUserId = getState().entities.users.currentUserId;
      var ismod = "FALSE"
      if (curUserId == creatorid){
        ismod = "TRUE"
      }
        try {
            url = await Client.getJoinURL(curUserId, meetingid,ismod);
            return {data:{joinurl : url}};
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
                    }},
                channelId
            });
            return {error};
        }
    };
}
export function isMeetingRunning(meetingid){
  return async (dispatch, getState) => {
    let response;
      try {
          response = await Client.isMeetingRunning(meetingid);
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
            url = await Client.endMeeting(getState().entities.users.currentUserId, meetingid);
            return {data:{joinurl : url}}; // don't know if its right to put return inside try
        } catch (error) {
            return {error};
        }
    };
}
export function getAttendees(channelId, meetingid) {
    return async (dispatch, getState) => {

        try {
            var resp = await Client.getAttendees(meetingid);
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
                    }},
                channelId
            });
            return {error};
        }
    };
}


export function publishRecordings(channelId, recordid,publish,meetingId) {
    return async (dispatch, getState) => {

        try {
            var resp = await Client.publishRecordings(recordid,publish,meetingId);
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
                message: 'Cant publish/unpublish',
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
                    }},
                channelId
            });
            return {error};
        }
    };
}

export function deleteRecordings(channelId, recordid,meetingId) {
    return async (dispatch, getState) => {

        try {
            var resp = await Client.deleteRecordings(recordid,meetingId);
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
                message: JSON.stringify(error),
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
                    }},
                channelId
            });
            return {error};
        }
    };
}

export function showRecordings() {
    return (dispatch, getState) => {

        const channel = getCurrentChannel(getState());

        const terms = "in:"+channel.name + " #recording" // change this to recordings

      //  trackEvent('api', 'api_posts_search_mention');

        dispatch(performSearch(terms, false));
        dispatch(batchActions([
            {
                type: ActionTypes.UPDATE_RHS_SEARCH_TERMS,
                terms,
            },
            {
                type: ActionTypes.UPDATE_RHS_STATE,
                state: RHSStates.MENTION,
            },
        ]));
    };
}
export function performSearch(terms, isMentionSearch) {
    return (dispatch, getState) => {
        const teamId = getCurrentTeamId(getState());

        return dispatch(searchPosts(teamId, terms, isMentionSearch));
    };
}
