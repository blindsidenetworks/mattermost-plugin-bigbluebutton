import request from 'superagent';
//superagent helps make post request

export default class Client {
    constructor() {
        this.url = '/plugins/bigbluebutton';
    }

    startMeeting = async (userid, channelid,topic,description) => {
        return this.doPost(`${this.url}/create`, {user_id: userid, channel_id: channelid, title : topic, description: description});
    }

    getJoinURL = async (userid, meetingid, ismod) => {
      var body = await this.doPost(`${this.url}/joinmeeting`, {user_id: userid, meeting_id : meetingid, is_mod: ismod});
      //console.log("from client " + JSON.stringify(body));
      return body;
    }
    endMeeting = async (userid, meetingid) => {
      var body = await this.doPost(`${this.url}/endmeeting`, {user_id: userid, meeting_id : meetingid});
      //console.log("from client " + JSON.stringify(body));
      return body;
    }
    isMeetingRunning = async (meetingid) =>{
      var body = await this.doPost(`${this.url}/ismeetingrunning`, {meeting_id : meetingid});
      //console.log("from client " + JSON.stringify(body));
      return body;
    }
    getAttendees = async (meetingid) =>{
      var body = await this.doPost(`${this.url}/getattendees`, {meeting_id: meetingid});
      //console.log(body)
      return body;
    }
    getRecordingsByChannel = async (channelid) =>{


      var body = await this.doPostXML(`${this.url}/getrecordingsbychannel`, {channel_id:channelid});

      console.log(body);
    }

    publishRecordings = async (recordid,publish,meetingId) => {
      return await this.doPost(`${this.url}/publishrecordings`, {record_id: recordid, publish : publish, meeting_id:meetingId});
    }

    deleteRecordings = async (recordid,meetingId) => {
      return await this.doPost(`${this.url}/deleterecordings`, {record_id: recordid,meeting_id:meetingId});
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
    doPostXML = async (url, body, headers = {}) => {
        headers['X-Requested-With'] = 'XMLHttpRequest';

        try {
            const response = await request.post(url).send(body).set(headers).type('application/json').accept('application/xml');

            return response.body;
        } catch (err) {
          console.log(err);
            throw err;
        }
    }
}
