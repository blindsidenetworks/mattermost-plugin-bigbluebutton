package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/mattermost/mattermost-server/model"
	"github.com/segmentio/ksuid"
	bbbAPI "github.com/ypgao1/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/api"
	"github.com/ypgao1/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/dataStructs"
	BBBwh "github.com/ypgao1/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/webhook"
)

const key = "key"

func (p *Plugin) PopulateMeeting(m *dataStructs.MeetingRoom, details []string, description string) {

	config := p.config()

	if len(details) == 2 {
		m.Name_ = details[1]
	} else {
		m.Name_ = "Big Blue Button Meeting"
	}
	m.MeetingID_ = GenerateRandomID()
	m.AttendeePW_ = "ap"
	m.ModeratorPW_ = "mp"
	m.Record = "true"
	m.AllowStartStopRecording = true
	m.AutoStartRecording = false
	m.Meta = description
	m.LogoutURL = "javascript:window.close();"
	m.LoopCount = 0

	var recordingcallbackurl string
	var Url *url.URL

	BaseUrl := config.CallBack_URL
	if !strings.HasPrefix(BaseUrl, "http") {
		BaseUrl = "http://" + BaseUrl
	}

	Url, _ = url.Parse(BaseUrl)
	Url.Path += "/plugins/bigbluebutton/recordingready"
	recordingcallbackurl = Url.String()
	m.Meta_bn_recording_ready_url = recordingcallbackurl

	var UrlEnd *url.URL
	UrlEnd, _ = url.Parse(BaseUrl)
	UrlEnd.Path += "/plugins/bigbluebutton/meetingendedcallback?" + m.MeetingID_
	Endmeetingcallback := UrlEnd.String()
	m.Meta_endcallbackurl = Endmeetingcallback

}

func (p *Plugin) LoadMeetingsFromStore() {
	byted, _ := p.api.KeyValueStore().Get("all_meetings")
	json.Unmarshal(byted, &p.Meetings)

	bytedRecordings, _ := p.api.KeyValueStore().Get("recording_meetings")
	json.Unmarshal(bytedRecordings, &p.MeetingsWaitingforRecordings)

	bytedLive, _ := p.api.KeyValueStore().Get("live_meetings")
	json.Unmarshal(bytedLive, &p.ActiveMeetings)

}
func (p *Plugin) SaveMeetingToStore() {
	byted, _ := json.Marshal(p.Meetings)
	p.api.KeyValueStore().Set("all_meetings", byted)

	bytedRecordings, _ := json.Marshal(p.MeetingsWaitingforRecordings)
	p.api.KeyValueStore().Set("recording_meetings", bytedRecordings)

	bytedLive, _ := json.Marshal(p.ActiveMeetings)
	p.api.KeyValueStore().Set("live_meetings", bytedLive)
}

func (p *Plugin) FindMeeting(meeting_id string) *dataStructs.MeetingRoom {
	for i := range p.Meetings {
		if p.Meetings[i].MeetingID_ == meeting_id {
			return &(p.Meetings[i])
		}
	}
	return nil
}
func (p *Plugin) FindMeetingfromInternal(meeting_id string) *dataStructs.MeetingRoom {
	for i := range p.Meetings {
		if p.Meetings[i].InternalMeetingId == meeting_id {
			return &(p.Meetings[i])
		}
	}
	return nil
}

func (p *Plugin) createEndMeetingWebook(callback_url string, meeting_id string) string {
	webhook := new(dataStructs.WebHook)
	webhook.CallBackURL = callback_url
	if meeting_id != "" {
		webhook.MeetingId = meeting_id
	}
	BBBwh.CreateHook(webhook)
	p.webhooks = append(p.webhooks, webhook)

	return webhook.HookID
}

func (p *Plugin) createStartMeetingPost(user_id string, channel_id string, m *dataStructs.MeetingRoom) {

	config := p.config()
	//if config page is not set uh oh
	if err := config.IsValid(); err != nil {
		return
	}

	textPost := &model.Post{UserId: user_id, ChannelId: channel_id,
		Message: "#BigBlueButton #" + m.Name_ + " #ID" + m.MeetingID_, Type: "custom_bbb"} //RootId: args.RootId, ParentId: args.ParentId,

	textPost.Props = model.StringInterface{
		"from_webhook":      "true",
		"override_username": "BigBlueButton",
		"override_icon_url": "https://pbs.twimg.com/profile_images/467451035837923328/JxPpOTL6_400x400.jpeg",
		"meeting_id":        m.MeetingID_,
		"meeting_link":      "www.google.com",
		"meeting_status":    "STARTED",
		"meeting_personal":  false,
		"meeting_topic":     m.Name_, //fill in this meeting topic
		"meeting_desc":      m.Meta,
		"user_count":        0,
	}

	postpointer, _ := p.api.CreatePost(textPost)
	m.PostId = postpointer.Id
}

func (p *Plugin) UpdatePostButtons(postid string, message string) {
	post, err := p.api.GetPost(postid)
	if err != nil {
		return
	}
	post.Message = message
	post.Props = nil

	p.api.UpdatePost(post)
	return

}

func (p *Plugin) DeleteMeeting(meeting_id string) {
	var index int
	for i := range p.Meetings {
		if p.Meetings[i].MeetingID_ == meeting_id {
			index = i
			break
		}
	}
	p.Meetings = append(p.Meetings[:index], p.Meetings[index+1:]...)
}
func (p *Plugin) DeleteActiveMeeting(meeting_id string) {
	var index int
	for i := range p.ActiveMeetings {
		if p.ActiveMeetings[i].MeetingID_ == meeting_id {
			index = i
			break
		}
	}
	p.ActiveMeetings = append(p.ActiveMeetings[:index], p.ActiveMeetings[index+1:]...)
}


// check if an attendee is a MODERATOR, useful feature decided not to include
// func (p *Plugin) IsModerator(meeting_id string, modpw string, fullname string) bool {
// 	var meetinginfo dataStructs.GetMeetingInfoResponse
// 	resp := bbbAPI.GetMeetingInfo(meeting_id, modpw, &meetinginfo)
// 	if resp == "FAILED" {
// 		return true
// 	}
// 	attendeesArray := meetinginfo.Attendees.Attendees
// 	for i := 0; i < len(attendeesArray); i++ {
// 		if attendeesArray[i].FullName == fullname {
// 			if attendeesArray[i].Role == "VIEWER" {
// 				return false
// 			} else if attendeesArray[i].Role == "MODERATOR" {
// 				return true
// 			}
// 		}
// 	}
// 	return false
// }

func GetAttendees(meeting_id string, modpw string) (int, []string) {
	var meetinginfo dataStructs.GetMeetingInfoResponse
	resp := bbbAPI.GetMeetingInfo(meeting_id, modpw, &meetinginfo)
	if resp == "FAILED" {
		return 0, []string{}
	}
	attendeesArray := meetinginfo.Attendees.Attendees

	NumAttendees := len(attendeesArray)
	var Fullnames []string
	for i := 0; i < NumAttendees; i++ {
		Fullnames = append(Fullnames, attendeesArray[i].FullName)
	}
	return NumAttendees, Fullnames
}
func FormatSeconds(seconds int64) string {
	var hours int64
	hours = seconds / 3600
	seconds = seconds - 3600*hours

	var minutes int64
	minutes = seconds / 60
	seconds = seconds - 60*minutes

	if hours != 0 {
		return fmt.Sprintf("%d hours, %d minutes and %d seconds", hours, minutes, seconds)
	} else if minutes != 0 {
		return fmt.Sprintf("%d minutes and %d seconds", minutes, seconds)
	}
	return fmt.Sprintf("%d seconds", seconds)
}

func GenerateRandomID() string {
	id := ksuid.New()
	return id.String()
}

func IsItemInArray(name string, array []string) bool {
	for _, word := range array {
		if word == name {
			return true
		}
	}
	return false
}
