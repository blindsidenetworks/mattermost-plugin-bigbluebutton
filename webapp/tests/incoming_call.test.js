import {afterEach, beforeAll, beforeEach, expect, test, vi} from 'vitest';
import React, {act} from 'react';
import {createRoot} from 'react-dom/client';
import {Client4} from 'mattermost-redux/client';

let IncomingCallPopup;
let root;
let container;
let props;
let getUser;

beforeAll(async () => {
	globalThis.IS_REACT_ACT_ENVIRONMENT = true;
	const Section = ({children}) => React.createElement('section', null, children);
	window.ReactBootstrap = {Modal: Object.assign(Section, {Header: Section, Body: Section, Footer: Section})};
	({default: IncomingCallPopup} = await import('../components/incoming_call/incoming_call.jsx'));
});

beforeEach(() => {
	container = document.createElement('div');
	document.body.appendChild(container);
	root = createRoot(container);
	props = {
		theme: {},
		siteURL: 'http://localhost:8065',
		currentTeam: {name: 'team'},
		incomingCall: {},
		actions: {dismissIncomingCall: vi.fn(), getJoinURL: vi.fn()},
	};
	getUser = vi.spyOn(Client4, 'getUser').mockResolvedValue({id: 'alice', username: 'alice'});
	vi.spyOn(Client4, 'setUrl').mockImplementation(() => {});
});

afterEach(async () => {
	if (root) {
		await act(async () => root.unmount());
	}
	container.remove();
	vi.restoreAllMocks();
});

async function render(call = props.incomingCall, changes = {}) {
	props = {...props, ...changes, incomingCall: call};
	await act(async () => root.render(React.createElement(IncomingCallPopup, props)));
}

function deferred() {
	let resolve;
	let reject;
	const promise = new Promise((res, rej) => {
		resolve = res;
		reject = rej;
	});
	return {promise, resolve, reject};
}

const aliceCall = {meetingId: 'meeting-a', fromUserID: 'alice'};
const bobCall = {meetingId: 'meeting-b', fromUserID: 'bob'};

test.each([
	null,
	{},
	{dismissed: false},
	{meetingId: 'meeting'},
	{meetingId: 'meeting', fromUserID: ''},
	{meetingId: 'meeting', fromUserID: '   '},
	{meetingId: 'meeting', fromUserID: 123},
	{fromUserID: 'alice'},
	{fromUserID: 'alice', meetingId: ''},
	{...aliceCall, dismissed: true},
])('does not request a user for an absent, malformed, or dismissed call: %j', async (call) => {
	await render(call);
	await render(call, {theme: {centerChannelBg: '#ffffff'}});
	expect(getUser).not.toHaveBeenCalled();
	expect(container.textContent).toBe('');
});

test('loads a call already present at mount and renders its caller', async () => {
	await render(aliceCall);
	expect(getUser).toHaveBeenCalledExactlyOnceWith('alice');
	expect(container.textContent).toContain('@alice');
	expect(container.querySelector('img').getAttribute('src')).toBe('http://localhost:8065/api/v4/users/alice/image');
});

test('unrelated updates and the profile state update do not duplicate pending requests', async () => {
	const pending = deferred();
	getUser.mockReturnValue(pending.promise);
	await render();
	await render(aliceCall);
	await render(aliceCall, {theme: {centerChannelBg: '#ffffff'}});
	await render(aliceCall, {currentTeam: {name: 'another-team'}});
	expect(getUser).toHaveBeenCalledExactlyOnceWith('alice');
	await act(async () => pending.resolve({id: 'alice', username: 'alice'}));
	expect(getUser).toHaveBeenCalledTimes(1);
	expect(container.textContent).toContain('@alice');
});

test('handles a rejected lookup without an unhandled promise or retry loop', async () => {
	getUser.mockRejectedValue(new Error('user unavailable'));
	await render(aliceCall);
	await render(aliceCall, {theme: {centerChannelBg: '#ffffff'}});
	expect(getUser).toHaveBeenCalledTimes(1);
	expect(props.actions.dismissIncomingCall).toHaveBeenCalledOnce();
	expect(container.textContent).toBe('');
	getUser.mockResolvedValue({id: 'bob', username: 'bob'});
	await render(bobCall);
	expect(container.textContent).toContain('@bob');
});

test('clears the previous profile while the next caller is loading', async () => {
	await render(aliceCall);
	const pending = deferred();
	getUser.mockReturnValue(pending.promise);
	await render(bobCall);
	expect(container.textContent).not.toContain('@alice');
	await act(async () => pending.resolve({id: 'bob', username: 'bob'}));
	expect(container.textContent).toContain('@bob');
});

test('clears a loaded profile when the call is reset without requesting an undefined user', async () => {
	await render(aliceCall);
	await render({});
	expect(container.textContent).toBe('');
	expect(getUser).toHaveBeenCalledExactlyOnceWith('alice');
});

test('ignores an older lookup that resolves after a newer call', async () => {
	const old = deferred();
	getUser.mockReturnValueOnce(old.promise).mockResolvedValueOnce({id: 'bob', username: 'bob'});
	await render(aliceCall);
	await render(bobCall);
	await act(async () => old.resolve({id: 'alice', username: 'alice'}));
	expect(container.textContent).toContain('@bob');
	expect(container.textContent).not.toContain('@alice');
});

test('an older failed lookup cannot dismiss the newer call', async () => {
	const old = deferred();
	getUser.mockReturnValueOnce(old.promise).mockResolvedValueOnce({id: 'bob', username: 'bob'});
	await render(aliceCall);
	await render(bobCall);
	await act(async () => old.reject(new Error('old request failed')));
	expect(props.actions.dismissIncomingCall).not.toHaveBeenCalled();
	expect(container.textContent).toContain('@bob');
});

test('dismissal invalidates a pending response', async () => {
	const pending = deferred();
	getUser.mockReturnValue(pending.promise);
	await render(aliceCall);
	await render({dismissed: true});
	await act(async () => pending.resolve({id: 'alice', username: 'alice'}));
	expect(container.textContent).toBe('');
	expect(getUser).toHaveBeenCalledTimes(1);
});

test.each(['resolve', 'reject'])('ignores a request that settles after unmount: %s', async (settlement) => {
	const pending = deferred();
	getUser.mockReturnValue(pending.promise);
	await render(aliceCall);
	await act(async () => root.unmount());
	root = null;
	await act(async () => {
		if (settlement === 'resolve') {
			pending.resolve({id: 'alice', username: 'alice'});
		} else {
			pending.reject(new Error('request failed after unmount'));
		}
	});
	expect(container.textContent).toBe('');
	expect(props.actions.dismissIncomingCall).not.toHaveBeenCalled();
});
