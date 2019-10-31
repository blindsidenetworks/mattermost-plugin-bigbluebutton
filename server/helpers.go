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

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	bbbAPI "github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/api"
	"github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/dataStructs"
	"github.com/mattermost/mattermost-server/model"
	"github.com/segmentio/ksuid"
	"net/url"
	"strings"
)

func (p *Plugin) PopulateMeeting(m *dataStructs.MeetingRoom, details []string, description string) error {

	if len(details) == 2 {
		m.Name_ = details[1]
	} else {
		m.Name_ = "Big Blue Button Meeting"
	}

	siteconfig := p.API.GetConfig()

	var callbackURL string
	if siteconfig.ServiceSettings.SiteURL != nil {
		callbackURL = *siteconfig.ServiceSettings.SiteURL
	} else {
		return errors.New("SiteURL not set")
	}
	if !strings.HasPrefix(callbackURL, "http") {
		callbackURL = "http://" + callbackURL
	}

	m.MeetingID_ = GenerateRandomID()
	m.AttendeePW_ = "ap"
	m.ModeratorPW_ = "mp"
	m.Record = "true"
	m.AllowStartStopRecording = true
	m.AutoStartRecording = false
	m.Meta = description
	var RedirectUrl *url.URL
	RedirectUrl, _ = url.Parse(callbackURL)
	RedirectUrl.Path += "/plugins/bigbluebutton/redirect"
	StringRedirectUrl := RedirectUrl.String()
	m.LogoutURL = StringRedirectUrl
	m.LoopCount = 0
	m.ValidToken = GenerateRandomID()

	var recordingcallbackurl string
	var Url *url.URL
	Url, _ = url.Parse(callbackURL)
	Url.Path += "/plugins/bigbluebutton/recordingready"
	recordingcallbackurl = Url.String()
	m.Meta_bn_recording_ready_url = recordingcallbackurl

	var UrlEnd *url.URL
	UrlEnd, _ = url.Parse(callbackURL)
	UrlEnd.Path += "/plugins/bigbluebutton/meetingendedcallback?" + m.MeetingID_ + "&" + m.ValidToken
	Endmeetingcallback := UrlEnd.String()
	m.Meta_endcallbackurl = Endmeetingcallback
	return nil
}

func (p *Plugin) LoadMeetingsFromStore() {
	byted, _ := p.API.KVGet("all_meetings")
	json.Unmarshal(byted, &p.Meetings)

	recordingsBytes, _ := p.API.KVGet("recording_queue")
	json.Unmarshal(recordingsBytes, &p.MeetingsWaitingforRecordings)

}

func (p *Plugin) SaveMeetingToStore() {
	byted, _ := json.Marshal(p.Meetings)
	p.API.KVSet("all_meetings", byted)

	recordingBytes, _ := json.Marshal(p.MeetingsWaitingforRecordings)
	p.API.KVSet("recording_queue", recordingBytes)

}

// Returns a meeting pointer so we'll be able to manipulate its content from outside the array.
func (p *Plugin) FindMeeting(meetingId string) *dataStructs.MeetingRoom {
	for i := range p.Meetings {
		if p.Meetings[i].MeetingID_ == meetingId {
			return &(p.Meetings[i])
		}
	}
	return nil
}

func (p *Plugin) FindMeetingfromInternal(meetingId string) *dataStructs.MeetingRoom {
	for i := range p.Meetings {
		if p.Meetings[i].InternalMeetingId == meetingId {
			return &(p.Meetings[i])
		}
	}
	return nil
}

func (p *Plugin) createStartMeetingPost(userId string, channelId string, m *dataStructs.MeetingRoom) {

	config := p.config()
	// If config page is not set uh oh.
	if err := config.IsValid(); err != nil {
		return
	}

	textPost := &model.Post{UserId: userId, ChannelId: channelId,
		Message: "#BigBlueButton #" + m.Name_ + " #ID" + m.MeetingID_, Type: "custom_bbb"}

	textPost.Props = model.StringInterface{
		"from_webhook":      "true",
		"override_username": "BigBlueButton",
		"override_icon_url": "https://pbs.twimg.com/profile_images/467451035837923328/JxPpOTL6_400x400.jpeg",
		"meeting_id":        m.MeetingID_,
		"meeting_status":    "STARTED",
		"meeting_personal":  false,
		"meeting_topic":     m.Name_, // Fill in this meeting topic.
		"meeting_desc":      m.Meta,
		"user_count":        0,
	}

	postpointer, _ := p.API.CreatePost(textPost)
	m.PostId = postpointer.Id
}

func (p *Plugin) DeleteMeeting(meetingId string) {
	var index int
	for i := range p.Meetings {
		if p.Meetings[i].MeetingID_ == meetingId {
			index = i
			break
		}
	}
	p.Meetings = append(p.Meetings[:index], p.Meetings[index+1:]...)
}

func GetAttendees(meetingId string, modPw string) (int, []string) {
	var meetinginfo dataStructs.GetMeetingInfoResponse

	if _, err := bbbAPI.GetMeetingInfo(meetingId, modPw, &meetinginfo); err != nil {
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
