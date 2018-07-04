package api

import (
	"github.com/ypgao1/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/dataStructs"
	"github.com/ypgao1/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/helpers"
	"log"
	"net/url"
	"strconv"
)

var BASE_URL string
var SALT string

func SetAPI(url string, salt string) {
	BASE_URL = url
	SALT = salt
}

//Creates A BigBlueButton meeting
// note: a BigBlueButton meeting will terminate 1 minute after its creation
// if there are no attendees currently present in the meeting

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

	checksum := helpers.GetChecksum("create" + createParam + SALT)

	response := helpers.HttpGet(BASE_URL + "create?" + createParam + "&checksum=" +
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

// we send in a Participant struct and get back a joinurl that participant can go to
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

	checksum := helpers.GetChecksum("join" + joinParam + SALT)
	joinUrl := BASE_URL + "join?" + joinParam + "&checksum=" + checksum
	participants.JoinURL = joinUrl

	return joinUrl
}

//only returns true when someone has joined the meeting
func IsMeetingRunning(meetingID string) bool {
	checksum := helpers.GetChecksum("isMeetingRunning" + "meetingID=" + meetingID + SALT)
	getURL := BASE_URL + "isMeetingRunning?" + "meetingID=" + meetingID + "&checksum=" + checksum
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

func EndMeeting(meeting_ID string, mod_PW string) string {
	meetingID := "meetingID=" + url.QueryEscape(meeting_ID)
	modPW := "&password=" + url.QueryEscape(mod_PW)
	param := meetingID + modPW
	checksum := helpers.GetChecksum("end" + param + SALT)

	getURL := BASE_URL + "end?" + param + "&checksum=" + checksum

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

//pass in meeting id, moderator password and address of a response structure,
// able to see new response info without having to get passed back the structure
func GetMeetingInfo(meeting_ID string, mod_PW string, responseXML *dataStructs.GetMeetingInfoResponse) string {
	meetingID := "meetingID=" + url.QueryEscape(meeting_ID)
	modPW := "&password=" + url.QueryEscape(mod_PW)
	param := meetingID + modPW
	checksum := helpers.GetChecksum("getMeetingInfo" + param + SALT)

	getURL := BASE_URL + "getMeetingInfo?" + param + "&checksum=" + checksum
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

//Gets all meetings and the details by returning a struct
func GetMeetings() dataStructs.GetMeetingsResponse {
	checksum := helpers.GetChecksum("getMeetings" + SALT)
	getURL := BASE_URL + "getMeetings?" + "&checksum=" + checksum
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
	checksum := helpers.GetChecksum("getRecordings" + param + SALT)
	getURL := BASE_URL + "getRecordings?" + param + "&checksum=" + checksum
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

func PublishRecordings(recordid string, publish string) dataStructs.PublishRecordingsResponse {
	recordID := "recordID=" + url.QueryEscape(recordid)
	Publish := "&publish=" + url.QueryEscape(publish)

	param := recordID + Publish
	checksum := helpers.GetChecksum("publishRecordings" + param + SALT)

	getURL := BASE_URL + "publishRecordings?" + param + "&checksum=" + checksum
	//log.Println(getURL)
	response := helpers.HttpGet(getURL)
	var XMLResp dataStructs.PublishRecordingsResponse

	helpers.ReadXML(response, &XMLResp)

	return XMLResp
}

func DeleteRecordings(recordid string) dataStructs.DeleteRecordingsResponse {
	recordID := "recordID=" + url.QueryEscape(recordid)
	param := recordID
	checksum := helpers.GetChecksum("deleteRecordings" + param + SALT)

	getURL := BASE_URL + "deleteRecordings?" + param + "&checksum=" + checksum
	//log.Println(getURL)
	response := helpers.HttpGet(getURL)
	var XMLResp dataStructs.DeleteRecordingsResponse

	helpers.ReadXML(response, &XMLResp)

	return XMLResp
}
