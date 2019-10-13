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

package api

import (
	"github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/dataStructs"
	"github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/helpers"
	"log"
	"net/url"
	"strconv"
)

//url of the BigBlueButton server
var BaseUrl string

//Secret of the BigBlueButton server
var salt string

//Sets the BaseUrl and salt
func SetAPI(url string, saltParam string) {
	BaseUrl = url
	salt = saltParam
}

//CreateMeeting creates A BigBlueButton meeting
// note: a BigBlueButton meeting will terminate 1 minute after its creation
// if there are no attendees currently present in the meeting
//
// see http://docs.bigbluebutton.org/dev/api.html for API documentation
func CreateMeeting(meetingRoom *dataStructs.MeetingRoom) string {
	if meetingRoom.Name_ == "" || meetingRoom.MeetingID_ == "" ||
		meetingRoom.AttendeePW_ == "" || meetingRoom.ModeratorPW_ == "" {
		log.Println("ERROR: PARAM ERROR.")
		return "ERROR: PARAM ERROR."
	}

	name := "name=" + url.QueryEscape(meetingRoom.Name_)
	meetingID := "&meetingID=" + url.QueryEscape(meetingRoom.MeetingID_)
	attendeePW := "&attendeePW=" + url.QueryEscape(meetingRoom.AttendeePW_)
	moderatorPW := "&moderatorPW=" + url.QueryEscape(meetingRoom.ModeratorPW_)
	welcome := "&welcome=" + url.QueryEscape(meetingRoom.Welcome)
	dialNumber := "&dialNumber=" + url.QueryEscape(meetingRoom.DialNumber)
	logoutURL := "&logoutURL=" + url.QueryEscape(meetingRoom.LogoutURL)
	record := "&record=" + url.QueryEscape(meetingRoom.Record)
	duration := "&duration=" + url.QueryEscape(strconv.Itoa(meetingRoom.Duration))
	allowStartStopRecording := "&allowStartStopRecording=" +
		url.QueryEscape(strconv.FormatBool(meetingRoom.AllowStartStopRecording))
	moderatorOnlyMessage := "&moderatorOnlyMessage=" +
		url.QueryEscape(meetingRoom.ModeratorOnlyMessage)
	meta_bn_recording_ready_url := "&meta_bn-recording-ready-url=" +
		url.QueryEscape(meetingRoom.Meta_bn_recording_ready_url)
	meta_channelid := "&meta_channelid=" +
		url.QueryEscape(meetingRoom.Meta_channelid)
	meta_endcallback := "&meta_endcallbackurl=" +
		url.QueryEscape(meetingRoom.Meta_endcallbackurl)
	voiceBridge := "&voiceBridge=" + url.QueryEscape(meetingRoom.VoiceBridge)

	createParam := name + meetingID + attendeePW + moderatorPW + welcome + dialNumber +
		voiceBridge + logoutURL + record + duration + moderatorOnlyMessage + meta_bn_recording_ready_url + meta_channelid +
		meta_endcallback + allowStartStopRecording

	checksum := helpers.GetChecksum("create" + createParam + salt)

	response := helpers.HttpGet(BaseUrl + "create?" + createParam + "&checksum=" +
		checksum)

	if "ERROR" == response {
		log.Println("ERROR: HTTP ERROR.")
		return "ERROR: HTTP ERROR."
	}
	err := helpers.ReadXML(response, &meetingRoom.CreateMeetingResponse)

	if nil != err {
		log.Println("XML PARSE ERROR: " + err.Error())
		return "ERROR: XML PARSE ERROR."
	}

	if "SUCCESS" == meetingRoom.CreateMeetingResponse.Returncode {
		log.Println("SUCCESS CREATE MEETINGROOM. MEETING ID: " +
			meetingRoom.CreateMeetingResponse.MeetingID)
		return meetingRoom.CreateMeetingResponse.MeetingID
	} else {
		log.Println("CREATE MEETINGROOM FAILD: " + response)
		return "FAILED"
	}
	return "ERROR: UNKNOWN."
}

// GetJoinURL: we send in a Participant struct and get back a joinurl that participant can go to
func GetJoinURL(participants *(dataStructs.Participants)) string {
	if "" == participants.FullName_ || "" == participants.MeetingID_ ||
		"" == participants.Password_ {
		return "ERROR: PARAM ERROR."
	}

	fullName := "fullName=" + url.QueryEscape(participants.FullName_)
	meetingID := "&meetingID=" + url.QueryEscape(participants.MeetingID_)
	password := "&password=" + url.QueryEscape(participants.Password_)

	var createTime string
	var userID string
	var configToken string
	var avatarURL string
	var redirect string
	var clientURL string

	if "" != participants.CreateTime {
		createTime = "&createTime=" + url.QueryEscape(participants.CreateTime)
	}

	if "" != participants.UserID {
		userID = "&userID=" + url.QueryEscape(participants.UserID)
	}

	if "" != participants.ConfigToken {
		configToken = "&configToken=" + url.QueryEscape(participants.ConfigToken)
	}

	if "" != participants.AvatarURL {
		avatarURL = "&avatarURL=" + url.QueryEscape(participants.AvatarURL)
	}

	if "" != participants.ClientURL {
		redirect = "&redirect=true"
		clientURL = "&clientURL=" + url.QueryEscape(participants.ClientURL)
	}
	joinviahtml := "&joinViaHtml5=true"

	joinParam := fullName + meetingID + password + createTime + userID +
		configToken + avatarURL + redirect + clientURL + joinviahtml

	checksum := helpers.GetChecksum("join" + joinParam + salt)
	joinUrl := BaseUrl + "join?" + joinParam + "&checksum=" + checksum
	participants.JoinURL = joinUrl

	return joinUrl
}

//IsMeetingRunning: only returns true when someone has joined the meeting
func IsMeetingRunning(meetingID string) bool {
	checksum := helpers.GetChecksum("isMeetingRunning" + "meetingID=" + meetingID + salt)
	getURL := BaseUrl + "isMeetingRunning?" + "meetingID=" + meetingID + "&checksum=" + checksum
	response := helpers.HttpGet(getURL)
	if "ERROR" == response {
		log.Println("ERROR: HTTP ERROR.")
		return false
	}
	var XMLResp dataStructs.IsMeetingRunningResponse
	err := helpers.ReadXML(response, &XMLResp)
	if nil != err {
		return false
	}

	return XMLResp.Running
}

//EndMeeting ends a BBB meeting
func EndMeeting(meeting_ID string, mod_PW string) string {
	meetingID := "meetingID=" + url.QueryEscape(meeting_ID)
	modPW := "&password=" + url.QueryEscape(mod_PW)
	param := meetingID + modPW
	checksum := helpers.GetChecksum("end" + param + salt)

	getURL := BaseUrl + "end?" + param + "&checksum=" + checksum

	response := helpers.HttpGet(getURL)

	if "ERROR" == response {
		log.Println("ERROR: HTTP ERROR.")
		return "Could not end meeting " + meeting_ID
	}
	var XMLResp dataStructs.EndResponse

	err := helpers.ReadXML(response, &XMLResp)
	if nil != err {
		return "Could not end meeting " + meeting_ID
	}

	if "SUCCESS" == XMLResp.ReturnCode {

		return "Successfully ended meeting " + meeting_ID
	} else {
		return "Could not end meeting " + meeting_ID
	}

}

//GetMeetingInfo: pass in meeting id, moderator password and address of a response structure,
// able to see new response info without having to get passed back the structure
func GetMeetingInfo(meeting_ID string, mod_PW string, responseXML *dataStructs.GetMeetingInfoResponse) string {
	meetingID := "meetingID=" + url.QueryEscape(meeting_ID)
	modPW := "&password=" + url.QueryEscape(mod_PW)
	param := meetingID + modPW
	checksum := helpers.GetChecksum("getMeetingInfo" + param + salt)

	getURL := BaseUrl + "getMeetingInfo?" + param + "&checksum=" + checksum
	response := helpers.HttpGet(getURL)

	if "ERROR" == response {
		log.Println("ERROR: HTTP ERROR.")
		return "FAILED"
	}

	err := helpers.ReadXML(response, responseXML)
	if nil != err {
		return "FAILED"
	}

	if "SUCCESS" == responseXML.ReturnCode {
		println("Successfully got meeting info")
		return "Successfully got meeting info" + meeting_ID
	} else {
		println("Could not get meeting info ")
		return "FAILED"
	}

}

//GetMeetings: Gets all meetings and the details by returning a struct
func GetMeetings() dataStructs.GetMeetingsResponse {
	checksum := helpers.GetChecksum("getMeetings" + salt)
	getURL := BaseUrl + "getMeetings?" + "&checksum=" + checksum
	response := helpers.HttpGet(getURL)

	if "ERROR" == response {
		log.Println("ERROR: HTTP ERROR.")
	}
	var XMLResp dataStructs.GetMeetingsResponse

	helpers.ReadXML(response, &XMLResp)

	if "SUCCESS" == XMLResp.ReturnCode {
		println("Successfully got meetings info")

	} else {
		println("Could not get meetings info ")
	}
	return XMLResp

}

//GetRecordings gets a recording for a BBB meeting
func GetRecordings(meeting_id string, record_id string, metachannelid string) (dataStructs.GetRecordingsResponse, string) {

	meetingID := "meetingID=" + url.QueryEscape(meeting_id)
	recordid := "&recordID=" + url.QueryEscape(record_id)
	var param string
	if metachannelid != "" {
		meta_channelid := "meta_channelid=" +
			url.QueryEscape(metachannelid)
		param = meta_channelid
	} else if meeting_id != "" && record_id != "" {
		param = meetingID + recordid
	} else if meeting_id != "" {
		param = meetingID
	}
	checksum := helpers.GetChecksum("getRecordings" + param + salt)
	getURL := BaseUrl + "getRecordings?" + param + "&checksum=" + checksum
	response := helpers.HttpGet(getURL)

	if "ERROR" == response {
		log.Println("ERROR: HTTP ERROR.")
	}
	var XMLResp dataStructs.GetRecordingsResponse

	err := helpers.ReadXML(response, &XMLResp)
	if nil != err {

	}
	if "SUCCESS" == XMLResp.ReturnCode {
		println("Successfully got recordings info")

	} else {
		println("Could not get recordings info ")
	}
	return XMLResp, response
}

//PublishRecordings
func PublishRecordings(recordid string, publish string) dataStructs.PublishRecordingsResponse {
	recordID := "recordID=" + url.QueryEscape(recordid)
	Publish := "&publish=" + url.QueryEscape(publish)

	param := recordID + Publish
	checksum := helpers.GetChecksum("publishRecordings" + param + salt)

	getURL := BaseUrl + "publishRecordings?" + param + "&checksum=" + checksum
	//log.Println(getURL)
	response := helpers.HttpGet(getURL)
	var XMLResp dataStructs.PublishRecordingsResponse

	helpers.ReadXML(response, &XMLResp)

	return XMLResp
}

//DeleteRecordings
func DeleteRecordings(recordid string) dataStructs.DeleteRecordingsResponse {
	recordID := "recordID=" + url.QueryEscape(recordid)
	param := recordID
	checksum := helpers.GetChecksum("deleteRecordings" + param + salt)

	getURL := BaseUrl + "deleteRecordings?" + param + "&checksum=" + checksum
	//log.Println(getURL)
	response := helpers.HttpGet(getURL)
	var XMLResp dataStructs.DeleteRecordingsResponse

	helpers.ReadXML(response, &XMLResp)

	return XMLResp
}
