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
	"github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/server/mattermost"
	"github.com/pkg/errors"
	"net/url"
	"strconv"
)

//url of the BigBlueButton server
var BaseUrl string

//Secret of the BigBlueButton server
var secret string

//Sets the BaseUrl and secret
func SetAPI(url string, secretParam string) {
	BaseUrl = url
	secret = secretParam
}

//CreateMeeting creates A BigBlueButton meeting
// note: a BigBlueButton meeting will terminate 1 minute after its creation
// if there are no attendees currently present in the meeting
//
// see http://docs.bigbluebutton.org/dev/api.html for API documentation
func CreateMeeting(meetingRoom *dataStructs.MeetingRoom) (string, error) {
	if meetingRoom.Name_ == "" {
		return "", errors.New("meeting name cannot be empty")
	}

	if meetingRoom.MeetingID_ == "" {
		return "", errors.New("meeting ID cannot be empty")
	}

	if meetingRoom.AttendeePW_ == "" {
		return "", errors.New("attendee PW cannot be empty")
	}

	if meetingRoom.ModeratorPW_ == "" {
		return "", errors.New("moderator PW cannot be empty")
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
	metaBnRecordingReadyUrl := "&meta_bn-recording-ready-url=" +
		url.QueryEscape(meetingRoom.Meta_bn_recording_ready_url)
	metaChannelId := "&meta_channelid=" +
		url.QueryEscape(meetingRoom.Meta_channelid)
	metaEndcallback := "&meta_endcallbackurl=" +
		url.QueryEscape(meetingRoom.Meta_endcallbackurl)
	voiceBridge := "&voiceBridge=" + url.QueryEscape(meetingRoom.VoiceBridge)

	createParam := name + meetingID + attendeePW + moderatorPW + welcome + dialNumber +
		voiceBridge + logoutURL + record + duration + moderatorOnlyMessage + metaBnRecordingReadyUrl + metaChannelId +
		metaEndcallback + allowStartStopRecording

	checksum := helpers.GetChecksum("create" + createParam + secret)

	response, err := helpers.HttpGet(BaseUrl + "create?" + createParam + "&checksum=" + checksum)

	if err != nil {
		mattermost.API.LogError("ERROR: HTTP ERROR: " + response)
		return "", errors.New(response)
	}

	if err := helpers.ReadXML(response, &meetingRoom.CreateMeetingResponse); err != nil {
		return "", err
	}

	if "SUCCESS" == meetingRoom.CreateMeetingResponse.Returncode {
		mattermost.API.LogInfo("SUCCESS CREATE MEETINGROOM. MEETING ID: " +
			meetingRoom.CreateMeetingResponse.MeetingID)
		return meetingRoom.CreateMeetingResponse.MeetingID, nil
	}

	mattermost.API.LogError("CREATE MEETINGROOM FAILD: " + response)
	return "", errors.New(response)
}

// GetJoinURL: we send in a Participant struct and get back a joinurl that participant can go to
func GetJoinURL(participants *dataStructs.Participants) (string, error) {
	if err := participants.IsValid(); err != nil {
		return "", err
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
	joinViaHtml5 := "&joinViaHtml5=true"

	joinParam := fullName + meetingID + password + createTime + userID +
		configToken + avatarURL + redirect + clientURL + joinViaHtml5

	checksum := helpers.GetChecksum("join" + joinParam + secret)
	joinUrl := BaseUrl + "join?" + joinParam + "&checksum=" + checksum
	participants.JoinURL = joinUrl

	return joinUrl, nil
}

//IsMeetingRunning: only returns true when someone has joined the meeting
func IsMeetingRunning(meetingID string) (bool, error) {
	checksum := helpers.GetChecksum("isMeetingRunning" + "meetingID=" + meetingID + secret)
	getURL := BaseUrl + "isMeetingRunning?" + "meetingID=" + meetingID + "&checksum=" + checksum
	response, err := helpers.HttpGet(getURL)
	if err != nil {
		return false, err
	}
	var xmlResp dataStructs.IsMeetingRunningResponse
	err = helpers.ReadXML(response, &xmlResp)
	if nil != err {
		return false, err
	}

	return xmlResp.Running, nil
}

//EndMeeting ends a BBB meeting
func EndMeeting(meetingId string, modPw string) (string, error) {
	meetingID := "meetingID=" + url.QueryEscape(meetingId)
	modPW := "&password=" + url.QueryEscape(modPw)
	param := meetingID + modPW
	checksum := helpers.GetChecksum("end" + param + secret)

	getURL := BaseUrl + "end?" + param + "&checksum=" + checksum

	response, err := helpers.HttpGet(getURL)
	if err != nil {
		return "", err
	}
	var XMLResp dataStructs.EndResponse

	err = helpers.ReadXML(response, &XMLResp)
	if nil != err {
		return "", err
	}

	if "SUCCESS" == XMLResp.ReturnCode {
		return "Successfully ended meeting " + meetingId, nil
	}

	return "", errors.New("Could not end meeting " + meetingId)
}

//GetMeetingInfo: pass in meeting id, moderator password and address of a response structure,
// able to see new response info without having to get passed back the structure
func GetMeetingInfo(meetingId string, modPw string, responseXML *dataStructs.GetMeetingInfoResponse) (string, error) {
	meetingID := "meetingID=" + url.QueryEscape(meetingId)
	modPW := "&password=" + url.QueryEscape(modPw)
	param := meetingID + modPW
	checksum := helpers.GetChecksum("getMeetingInfo" + param + secret)

	getURL := BaseUrl + "getMeetingInfo?" + param + "&checksum=" + checksum
	response, err := helpers.HttpGet(getURL)
	if err != nil {
		return "", err
	}

	err = helpers.ReadXML(response, responseXML)
	if nil != err {
		return "", err
	}

	if "SUCCESS" == responseXML.ReturnCode {
		mattermost.API.LogInfo("Successfully got meeting info")
		return "Successfully got meeting info" + meetingId, nil
	} else {
		return "", errors.New("Could not get meeting info ")
	}
}

//GetMeetings: Gets all meetings and the details by returning a struct
func GetMeetings() (dataStructs.GetMeetingsResponse, error) {
	checksum := helpers.GetChecksum("getMeetings" + secret)
	getURL := BaseUrl + "getMeetings?" + "&checksum=" + checksum
	response, err := helpers.HttpGet(getURL)

	if err != nil {
		return dataStructs.GetMeetingsResponse{}, err
	}
	var XMLResp dataStructs.GetMeetingsResponse

	if err := helpers.ReadXML(response, &XMLResp); err != nil {
		return dataStructs.GetMeetingsResponse{}, err
	}

	if "SUCCESS" == XMLResp.ReturnCode {
		mattermost.API.LogInfo("Successfully got meetings info")
	} else {
		mattermost.API.LogError("Could not get meetings info ")
	}
	return XMLResp, nil
}

//GetRecordings gets a recording for a BBB meeting
func GetRecordings(meetingId string, recordId string, metachannelid string) (dataStructs.GetRecordingsResponse, string, error) {
	meetingID := "meetingID=" + url.QueryEscape(meetingId)
	recordid := "&recordID=" + url.QueryEscape(recordId)
	var param string
	if metachannelid != "" {
		metaChannelId := "meta_channelid=" +
			url.QueryEscape(metachannelid)
		param = metaChannelId
	} else if meetingId != "" && recordId != "" {
		param = meetingID + recordid
	} else if meetingId != "" {
		param = meetingID
	}
	checksum := helpers.GetChecksum("getRecordings" + param + secret)
	getURL := BaseUrl + "getRecordings?" + param + "&checksum=" + checksum
	response, err := helpers.HttpGet(getURL)

	if err != nil {
		return dataStructs.GetRecordingsResponse{}, "", err
	}
	var XMLResp dataStructs.GetRecordingsResponse

	if err := helpers.ReadXML(response, &XMLResp); nil != err {
		return dataStructs.GetRecordingsResponse{}, "", err
	}

	if "SUCCESS" == XMLResp.ReturnCode {
		mattermost.API.LogInfo("Successfully got recordings info")
	} else {
		return dataStructs.GetRecordingsResponse{}, "", errors.New("Could not get recordings info")
	}
	return XMLResp, response, nil
}

//PublishRecordings
func PublishRecordings(recordid string, publish string) (dataStructs.PublishRecordingsResponse, error) {
	recordID := "recordID=" + url.QueryEscape(recordid)
	Publish := "&publish=" + url.QueryEscape(publish)

	param := recordID + Publish
	checksum := helpers.GetChecksum("publishRecordings" + param + secret)

	getURL := BaseUrl + "publishRecordings?" + param + "&checksum=" + checksum
	response, err := helpers.HttpGet(getURL)
	if err != nil {
		return dataStructs.PublishRecordingsResponse{}, err
	}

	var XMLResp dataStructs.PublishRecordingsResponse
	helpers.ReadXML(response, &XMLResp)
	return XMLResp, nil
}

//DeleteRecordings
func DeleteRecordings(recordid string) (dataStructs.DeleteRecordingsResponse, error) {
	recordID := "recordID=" + url.QueryEscape(recordid)
	param := recordID
	checksum := helpers.GetChecksum("deleteRecordings" + param + secret)

	getURL := BaseUrl + "deleteRecordings?" + param + "&checksum=" + checksum
	response, err := helpers.HttpGet(getURL)
	if err != nil {
		return dataStructs.DeleteRecordingsResponse{}, err
	}

	var XMLResp dataStructs.DeleteRecordingsResponse
	helpers.ReadXML(response, &XMLResp)
	return XMLResp, nil
}
