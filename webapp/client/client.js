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

import {Client4} from 'mattermost-redux/client';

import request from 'superagent';
//superagent helps make post request

//client.js is used to communicate with out backend server
export default class Client {
	constructor(siteUrlFunc) {
		this.siteUrlFunc = siteUrlFunc;
	}

	get baseUrl() {
		return `${this.siteUrlFunc()}/plugins/bigbluebutton`;
	}

	startMeeting = async (userid, channelid, topic, description, allowRecording) => {
		return this.doPost(`${this.baseUrl}/create`, {
			user_id: userid,
			channel_id: channelid,
			title: topic,
			description: description,
			allow_recording: allowRecording,
		});
	}

	getJoinURL = async (userid, meetingid, ismod) => {
		var body = await this.doPost(`${this.baseUrl}/joininvite`, {
			user_id: userid,
			meeting_id: meetingid,
			is_mod: ismod
		});
		return body;
	}
	endMeeting = async (userid, meetingid) => {
		var body = await this.doPost(`${this.baseUrl}/endmeeting`, {
			user_id: userid,
			meeting_id: meetingid
		});
		return body;
	}
	isMeetingRunning = async (meetingid) => {
		var body = await this.doPost(`${this.baseUrl}/ismeetingrunning`, {meeting_id: meetingid});
		return body;
	}
	getAttendees = async (meetingid) => {
		var body = await this.doPost(`${this.baseUrl}/getattendees`, {meeting_id: meetingid});
		return body;
	}

	publishRecordings = async (recordid, publish, meetingId) => {
		return await this.doPost(`${this.baseUrl}/publishrecordings`, {
			record_id: recordid,
			publish: publish,
			meeting_id: meetingId
		});
	}

	deleteRecordings = async (recordid, meetingId) => {
		return await this.doPost(`${this.baseUrl}/deleterecordings`, {
			record_id: recordid,
			meeting_id: meetingId
		});
	}

	getPluginConfig = async () => {
		return await this.doPost(`${this.baseUrl}/config`);
	}

	doPost = async (url, body, headers = {}) => {
		const options = Client4.getOptions({
			method: 'post',
			headers,
		});

		try {
			const response = await request.post(url).send(body).set(options.headers).type('application/json').accept('application/json');

			return response.body;
		} catch (err) {
			console.log(err);
			throw err;
		}
	}
}
