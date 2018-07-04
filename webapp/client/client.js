import request from 'superagent';
//superagent helps make post request

//client.js is used to communicate with out backend server
export default class Client {
  constructor() {
    this.url = '/plugins/bigbluebutton';
  }

  startMeeting = async (userid, channelid, topic, description) => {
    return this.doPost(`${this.url}/create`, {
      user_id: userid,
      channel_id: channelid,
      title: topic,
      description: description
    });
  }

  getJoinURL = async (userid, meetingid, ismod) => {
    var body = await this.doPost(`${this.url}/joinmeeting`, {
      user_id: userid,
      meeting_id: meetingid,
      is_mod: ismod
    });
    return body;
  }
  endMeeting = async (userid, meetingid) => {
    var body = await this.doPost(`${this.url}/endmeeting`, {
      user_id: userid,
      meeting_id: meetingid
    });
    return body;
  }
  isMeetingRunning = async (meetingid) => {
    var body = await this.doPost(`${this.url}/ismeetingrunning`, {meeting_id: meetingid});
    return body;
  }
  getAttendees = async (meetingid) => {
    var body = await this.doPost(`${this.url}/getattendees`, {meeting_id: meetingid});
    return body;
  }

  publishRecordings = async (recordid, publish, meetingId) => {
    return await this.doPost(`${this.url}/publishrecordings`, {
      record_id: recordid,
      publish: publish,
      meeting_id: meetingId
    });
  }

  deleteRecordings = async (recordid, meetingId) => {
    return await this.doPost(`${this.url}/deleterecordings`, {
      record_id: recordid,
      meeting_id: meetingId
    });
  }

  doPost = async (url, body, headers = {}) => {
    headers['X-Requested-With'] = 'XMLHttpRequest';

    try {
      const response = await request.post(url).send(body).set(headers).type('application/json').accept('application/json');

      return response.body;
    } catch (err) {
      console.log(err);
      throw err;
    }
  }
}
