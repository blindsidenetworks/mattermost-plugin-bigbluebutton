import {beforeEach, expect, test, vi} from 'vitest';
import React from 'react';
import {renderToStaticMarkup} from 'react-dom/server';
import * as Redux from 'redux';
import entitiesReducer from 'mattermost-redux/reducers/entities';
import {Client4} from 'mattermost-redux/client';
import Client from '../client/client';
import superagent from 'superagent';
import reducer from '../reducer';
import {DISMISS_INCOMING_CALL, INCOMING_CALL, OPEN_ROOT_MODAL, SET_PLUGIN_CONFIG} from '../action_types';

const connected = [];

beforeEach(() => {
	connected.length = 0;
	window.Redux = Redux;
	window.ReactRedux = {
		connect: (mapState) => (Component) => {
			connected.push(mapState);
			return Component;
		},
	};
	// Mattermost supplies these UI components. Keep them outside the plugin bundle.
	window.ReactBootstrap = Object.fromEntries(['Popover', 'Overlay', 'Modal', 'OverlayTrigger', 'Tooltip'].map((name) => [name, ({children}) => React.createElement('div', {}, children)]));
});

test('plugin initializes and its containers accept the current Mattermost Redux state', async () => {
	vi.resetModules();
	window.registerPlugin = vi.fn();
	await import('../index');
	const [id, plugin] = window.registerPlugin.mock.calls[0];
	expect(id).toBe('bigbluebutton');
	const registry = Object.fromEntries([
		'registerChannelHeaderButtonAction', 'registerAppBarComponent',
		'registerPopoverUserActionsComponent', 'registerRootComponent',
		'registerReducer', 'registerWebSocketEventHandler',
	].map((name) => [name, vi.fn()]));
	const state = {entities: entitiesReducer(undefined, {type: 'INIT'}), views: {channel: {}}};
	state.entities.general.config = {SiteURL: 'https://chat.example', EnableAppBar: 'true'};
	state.entities.channels.currentChannelId = 'channel';
	state.entities.channels.channels.channel = {id: 'channel', name: 'town-square', type: 'O', display_name: 'Town Square'};
	const dispatch = vi.fn();
	const response = {send: vi.fn().mockReturnThis(), set: vi.fn().mockReturnThis(), type: vi.fn().mockReturnThis(), accept: vi.fn().mockResolvedValue({body: {ALLOW_RECORDINGS: true}})};
	const request = vi.spyOn(superagent, 'post').mockReturnValue(response);
	await plugin.initialize(registry, {getState: () => state, dispatch});
	expect(registry.registerRootComponent).toHaveBeenCalledTimes(3);
	expect(registry.registerWebSocketEventHandler).toHaveBeenCalledTimes(3);
	for (const mapState of connected) {
		expect(() => mapState(state, {})).not.toThrow();
	}
	const handler = registry.registerWebSocketEventHandler.mock.calls.find(([name]) => name.endsWith('_incoming_call'))[1];
	handler({data: {meetingId: 'meeting', fromUserID: 'alice'}});
	const action = dispatch.mock.calls.at(-1)[0];
	expect(reducer(undefined, action).incomingCall).toMatchObject({meetingId: 'meeting', dismissed: false});
	request.mockRestore();
});

test('meeting API posts JSON with the current client headers and propagates errors', async () => {
	const client = new Client(() => 'https://chat.example');
	const options = vi.spyOn(Client4, 'getOptions').mockReturnValue({headers: {'X-Requested-With': 'XMLHttpRequest'}});
	// Superagent uses browser XMLHttpRequest in jsdom. Stub only the transport chain.
	const {default: request} = await import('superagent');
	const chain = {send: vi.fn().mockReturnThis(), set: vi.fn().mockReturnThis(), type: vi.fn().mockReturnThis(), accept: vi.fn().mockResolvedValue({body: {meetingId: 'meeting'}})};
	const post = vi.spyOn(request, 'post').mockReturnValue(chain);
	await expect(client.startMeeting('alice', 'channel', 'Topic', 'Description', true)).resolves.toEqual({meetingId: 'meeting'});
	expect(post).toHaveBeenCalledWith('https://chat.example/plugins/bigbluebutton/create');
	expect(chain.send).toHaveBeenCalledWith({user_id: 'alice', channel_id: 'channel', title: 'Topic', description: 'Description', allow_recording: true});
	expect(chain.set).toHaveBeenCalledWith({'X-Requested-With': 'XMLHttpRequest'});
	const error = new Error('offline');
	chain.accept.mockRejectedValueOnce(error);
	const log = vi.spyOn(console, 'log').mockImplementation(() => {});
	await expect(client.getPluginConfig()).rejects.toBe(error);
	log.mockRestore();
	post.mockRestore();
	options.mockRestore();
});

test('Redux 5 preserves meeting state and React renders meeting controls', async () => {
	const store = Redux.createStore(reducer);
	store.dispatch({type: SET_PLUGIN_CONFIG, data: {ALLOW_RECORDINGS: true}});
	store.dispatch({type: OPEN_ROOT_MODAL});
	store.dispatch({type: INCOMING_CALL, data: {meetingId: 'meeting'}});
	store.dispatch({type: DISMISS_INCOMING_CALL});
	expect(store.getState()).toMatchObject({rootModalVisible: true, pluginConfig: {ALLOW_RECORDINGS: true}, incomingCall: {dismissed: true}});
	const {default: Root} = await import('../components/root/root.jsx');
	const channel = {id: 'channel', name: 'town-square', display_name: 'Town Square', type: 'O'};
	const html = renderToStaticMarkup(React.createElement(Root, {
		state: {entities: {channels: {channels: {channel}}}},
		channel, channelId: 'channel', visible: true, appBarEnabled: false,
		pluginConfig: {ALLOW_RECORDINGS: true},
		theme: {centerChannelBg: '#ffffff', centerChannelColor: '#333333', linkColor: '#0055aa'},
		actions: {startMeeting: vi.fn(), showRecordings: vi.fn(), closePopover: vi.fn()},
	}));
	expect(html).toContain('Start New Meeting');
	expect(html).toContain('Recording Disabled');
	expect(html).toContain('View Recordings');
});

test('meeting websocket changes open a join URL only for the current channel', async () => {
	const {default: OpenMeeting} = await import('../components/open_meeting/open_meeting.jsx');
	const dismissOpenMeeting = vi.fn();
	const open = vi.spyOn(window, 'open').mockImplementation(() => {});
	const component = new OpenMeeting({
		siteURL: 'https://chat.example', channelID: 'channel',
		meeting: {channelID: 'channel', joinURL: 'https://bbb.example/join'},
		actions: {dismissOpenMeeting},
	});
	component.componentDidUpdate({meeting: {}});
	expect(open).toHaveBeenCalledWith('https://bbb.example/join', '_blank');
	expect(dismissOpenMeeting).toHaveBeenCalledOnce();
	open.mockClear();
	component.props = {...component.props, channelID: 'other-channel'};
	component.componentDidUpdate({meeting: {}});
	expect(open).not.toHaveBeenCalled();
	open.mockRestore();
});
