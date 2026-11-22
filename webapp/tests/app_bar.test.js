import {afterEach, beforeAll, beforeEach, expect, test, vi} from 'vitest';
import React, {act} from 'react';
import {createRoot} from 'react-dom/client';
import {createPortal} from 'react-dom';
import * as Redux from 'redux';
import superagent from 'superagent';
import reducer from '../reducer';
import {CLOSE_ROOT_MODAL, OPEN_ROOT_MODAL, SET_PLUGIN_CONFIG} from '../action_types';

let Root;
let plugin;
let root;
let container;
let store;
let overlayProps;
let unsubscribe;
let actions;
const channel = {id: 'channel', name: 'town-square', display_name: 'Town Square', type: 'O'};

beforeAll(async () => {
	globalThis.IS_REACT_ACT_ENVIRONMENT = true;
	window.Redux = Redux;
	window.ReactRedux = {connect: () => (Component) => Component};
	const Section = ({children, id, style}) => React.createElement('section', {id, style}, children);
	window.ReactBootstrap = {
		Popover: Section,
		// Enforce the host overlay contract: a visible overlay needs a target,
		// and is portaled into its container rather than the root component.
		Overlay: (props) => {
			overlayProps = props;
			return props.show && props.target?.() ? createPortal(props.children, props.container) : null;
		},
		Modal: Section,
		OverlayTrigger: Section,
		Tooltip: Section,
	};
	window.registerPlugin = (id, instance) => { plugin = instance; };
	({default: Root} = await import('../components/root/root.jsx'));
	await import('../index');
});

beforeEach(() => {
	overlayProps = null;
	container = document.createElement('div');
	document.body.appendChild(container);
	root = createRoot(container);
	const entities = {
		general: {config: {SiteURL: 'https://chat.example'}}, // EnableAppBar is absent on current hosts.
		channels: {currentChannelId: 'channel', channels: {channel}},
	};
	store = Redux.createStore(Redux.combineReducers({
		entities: (state = entities) => state,
		'plugins-bigbluebutton': reducer,
	}));
	store.dispatch({type: SET_PLUGIN_CONFIG, data: {ALLOW_RECORDINGS: true}});
	actions = {
		startMeeting: vi.fn().mockResolvedValue({data: true}),
		showRecordings: vi.fn(),
		closePopover: () => store.dispatch({type: CLOSE_ROOT_MODAL}),
	};
});

afterEach(async () => {
	unsubscribe?.();
	unsubscribe = null;
	await act(async () => root.unmount());
	container.remove();
	document.querySelectorAll('[data-test-anchor]').forEach((element) => element.remove());
	vi.restoreAllMocks();
});

function anchor(id) {
	const element = document.createElement('button');
	element.id = id;
	element.dataset.testAnchor = 'true';
	document.body.appendChild(element);
	return element;
}

async function mount() {
	const render = () => {
		const state = store.getState();
		root.render(React.createElement(Root, {
			state, channel, channelId: channel.id,
			visible: state['plugins-bigbluebutton'].rootModalVisible,
			pluginConfig: state['plugins-bigbluebutton'].pluginConfig,
			theme: {centerChannelBg: '#ffffff', centerChannelColor: '#333333', linkColor: '#0055aa'},
			actions,
		}));
	};
	unsubscribe = store.subscribe(render);
	await act(async () => render());
}

test('the registered app-bar click opens an anchored menu without EnableAppBar', async () => {
	const registry = Object.fromEntries([
		'registerChannelHeaderButtonAction', 'registerAppBarComponent',
		'registerPopoverUserActionsComponent', 'registerRootComponent',
		'registerReducer', 'registerWebSocketEventHandler',
	].map((name) => [name, vi.fn()]));
	const response = {send: vi.fn().mockReturnThis(), set: vi.fn().mockReturnThis(), type: vi.fn().mockReturnThis(), accept: vi.fn().mockResolvedValue({body: {ALLOW_RECORDINGS: true}})};
	vi.spyOn(superagent, 'post').mockReturnValue(response);
	const dispatch = (action) => typeof action === 'function' ? action(dispatch, store.getState) : store.dispatch(action);
	await plugin.initialize(registry, {getState: store.getState, dispatch});
	const button = anchor('app-bar-icon-bigbluebutton');
	button.onclick = registry.registerAppBarComponent.mock.calls[0][1];
	await mount();
	expect(document.getElementById('bbbPopover')).toBeNull();
	await act(async () => button.click());
	expect(store.getState()['plugins-bigbluebutton'].rootModalVisible).toBe(true);
	expect(overlayProps.target()).toBe(button);
	expect(overlayProps.placement).toBe('left');
	expect(overlayProps.container).toBe(document.body);
	const menu = document.getElementById('bbbPopover');
	expect(menu.textContent).toContain('Start New Meeting');
	expect(menu.style.marginLeft).toBe('');
	expect(menu.style.marginTop).toBe('');
	const start = Array.from(menu.querySelectorAll('button')).find((item) => item.textContent.includes('Allow Recording'));
	await act(async () => start.click());
	expect(actions.startMeeting).toHaveBeenCalledExactlyOnceWith('channel', true, '', 'Town Square');
	expect(document.getElementById('bbbPopover')).toBeNull();
});

test('the header button remains a supported overlay target when the app-bar button is absent', async () => {
	const button = anchor('bbbChannelHeaderButton');
	store.dispatch({type: OPEN_ROOT_MODAL});
	await mount();
	expect(overlayProps.target()).toBe(button);
	expect(overlayProps.placement).toBe('bottom');
	expect(document.getElementById('bbbPopover')).not.toBeNull();
	await act(async () => overlayProps.onHide());
	expect(document.getElementById('bbbPopover')).toBeNull();
});

test('prefers the app-bar anchor when both controls are present', async () => {
	anchor('bbbChannelHeaderButton');
	const button = anchor('app-bar-icon-bigbluebutton');
	store.dispatch({type: OPEN_ROOT_MODAL});
	await mount();
	expect(overlayProps.target()).toBe(button);
	expect(overlayProps.placement).toBe('left');
});

test('does not show a popover at the page origin when no target exists', async () => {
	store.dispatch({type: OPEN_ROOT_MODAL});
	await mount();
	expect(overlayProps.show).toBe(false);
	expect(document.getElementById('bbbPopover')).toBeNull();
});

test('the app-bar action starts an unrecorded meeting when recording options are disabled', async () => {
	anchor('app-bar-icon-bigbluebutton');
	store.dispatch({type: SET_PLUGIN_CONFIG, data: {ALLOW_RECORDINGS: false}});
	await mount();
	await act(async () => store.dispatch({type: OPEN_ROOT_MODAL}));
	expect(actions.startMeeting).toHaveBeenCalledExactlyOnceWith('channel', false, '', 'Town Square');
	expect(store.getState()['plugins-bigbluebutton'].rootModalVisible).toBe(false);
});
