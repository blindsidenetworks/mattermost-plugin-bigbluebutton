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
	"github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/helpers"
	"github.com/mattermost/mattermost-server/v5/utils"
	"net/url"
	"strconv"
	"strings"

	bbbAPI "github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/api"
	"github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/dataStructs"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/segmentio/ksuid"
)

const (
	// KV Store key prefixes
	prefixMeeting     = "m_"
	prefixMeetingList = "m_list_"
)

func (p *Plugin) PopulateMeeting(m *dataStructs.MeetingRoom, details []string, description string, UserId string, channelId string) error {
	if len(details) == 2 {
		m.Name_ = details[1]
	} else {
		m.Name_ = "Big Blue Button Meeting"
	}

	siteConfig := p.API.GetConfig()

	var callbackURL string
	if siteConfig.ServiceSettings.SiteURL == nil {
		return errors.New("SiteURL not set")
	}

	callbackURL = *siteConfig.ServiceSettings.SiteURL
	if !strings.HasPrefix(callbackURL, "http") {
		callbackURL = "http://" + callbackURL
	}

	m.MeetingID_ = GenerateRandomID()
	m.AttendeePW_ = "ap"
	m.ModeratorPW_ = "mp"
	m.Record = strconv.FormatBool(p.config().AllowRecordings)
	m.AllowStartStopRecording = p.config().AllowRecordings
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
	m.Meta_bbb_recording_ready_url = recordingcallbackurl

	m.Meta_channelid = channelId

	var UrlEnd *url.URL
	UrlEnd, _ = url.Parse(callbackURL)
	UrlEnd.Path += "/plugins/bigbluebutton/meetingendedcallback?" + m.MeetingID_ + "&" + m.ValidToken
	Endmeetingcallback := UrlEnd.String()
	m.Meta_endcallbackurl = Endmeetingcallback

	m.Meta_bbb_origin = "Mattermost"
	m.Meta_bbb_origin_version = helpers.PluginVersion
	if siteConfig.ServiceSettings.SiteURL != nil {
		m.Meta_bbb_origin_server_name = utils.GetHostnameFromSiteURL(*siteConfig.ServiceSettings.SiteURL)
	} else {
		return errors.New("SiteURL not set")
	}

	user, err := p.API.GetUser(UserId)
	if err != nil {
		return errors.New("Error resolving UserId")
	}
	m.Meta_dc_creator = user.Email

	return nil
}

// Returns a meeting pointer so we'll be able to manipulate its content from outside the array.
func (p *Plugin) FindMeeting(meetingId string) *dataStructs.MeetingRoom {
	meeting, _ := p.GetMeeting(meetingId)
	return meeting
}

func (p *Plugin) createStartMeetingPost(userId string, channelId string, m *dataStructs.MeetingRoom) {
	config := p.config()
	// If config page is not set uh oh.
	if err := config.IsValid(); err != nil {
		return
	}

	post := &model.Post{
		UserId:    userId,
		ChannelId: channelId,
	}

	user, appErr := p.API.GetUser(userId)
	if appErr != nil {
		p.API.LogError("Failed to fetch user id: " + userId)
		return
	}

	post.AddProp("created_by", user.Id)

	attachments := []*model.SlackAttachment{
		{
			Text: "Meeting created by @" + user.Username,
			Fields: []*model.SlackAttachmentField{
				{
					Title: "Attendees",
					Value: "*There are no attendees in this session*",
					Short: false,
				},
				{
					Title: "",
					Value: "",
					Short: false,
				},
			},
			Actions: []*model.PostAction{
				{
					Id:    "bigBlueButtonJoinMeeting",
					Type:  "button",
					Name:  "Join Meeting",
					Style: "primary",
					Integration: &model.PostActionIntegration{
						URL: "/plugins/bigbluebutton/joinmeeting",
						Context: map[string]interface{}{
							"meetingId": m.MeetingID_,
						},
					},
				},
				{
					Id:    "bigBlueButtonEndMeeting",
					Type:  "button",
					Name:  "End Meeting",
					Style: "danger",
					Integration: &model.PostActionIntegration{
						URL: "/plugins/bigbluebutton/endmeeting",
						Context: map[string]interface{}{
							"meetingId": m.MeetingID_,
						},
					},
				},
			},
		},
	}

	model.ParseSlackAttachment(post, attachments)

	post.AddProp("from_webhook", true)
	post.AddProp("override_username", "BigBlueButton")
	post.AddProp("override_icon_url", "strings.TrimSuffix(*p.API.GetConfig().ServiceSettings.SiteURL, \"/\") + \"/plugins/bigbluebutton/bbb.png\",")
	post.AddProp("meeting_id", m.MeetingID_)
	post.AddProp("meeting_status", "started")
	post.AddProp("meeting_personal", false)
	post.AddProp("meeting_topic", m.Name_)
	post.AddProp("meeting_desc", m.Meta)
	post.AddProp("user_count", 0)

	postpointer, _ := p.API.CreatePost(post)
	m.PostId = postpointer.Id
}

func (p *Plugin) DeleteMeeting(meetingId string) {
	if appErr := p.API.KVDelete(prefixMeeting + meetingId); appErr != nil {
		p.API.LogError(fmt.Sprintf("Unable to delete meeting from KV store. Meeting ID: {%s}, error: {%s}", meetingId, appErr.Error()))
	}
}

func (p *Plugin) SaveMeeting(meeting *dataStructs.MeetingRoom) error {
	data, err := json.Marshal(meeting)
	if err != nil {
		p.API.LogError(fmt.Sprintf("Unable to marshal meeting for storing in KV store.Meeting ID: {%s}, error: {%s}", meeting.MeetingID_, err.Error()))
		return err
	}

	appErr := p.API.KVSet(prefixMeeting+meeting.MeetingID_, data)
	if appErr != nil {
		p.API.LogError(fmt.Sprintf("Unable to save meeting in KV store. Meeting ID: {%s}, error: {%s}", meeting.MeetingID_, appErr.Error()))
		return err
	}

	if err := p.addToMeetingList(meeting.MeetingID_); err != nil {
		return err
	}

	return nil
}

func (p *Plugin) addToMeetingList(meetingID string) error {
	meetings, err := p.GetMeetingList()
	if err != nil {
		return err
	}

	meetings = append(meetings, prefixMeeting+meetingID)
	data, err := json.Marshal(meetings)
	if err != nil {
		p.API.LogError(fmt.Sprintf("Unable to marshal meeting list. Error: {%s}", err.Error()))
		return err
	}

	if appErr := p.API.KVSet(prefixMeetingList, data); appErr != nil {
		p.API.LogError(fmt.Sprintf("Unable to save updated meeting list in KV store. Error: {%s}", appErr.Error()))
		return err
	}

	return nil
}

func (p *Plugin) GetMeetingList() ([]string, error) {
	var meetings *[]string
	data, appErr := p.API.KVGet(prefixMeetingList)
	if appErr != nil {
		p.API.LogError(fmt.Sprintf("Unable to fetch meeting list. Error: {%s}", appErr.Error()))
		return nil, appErr
	}

	// This handles the case of no data present in KV store.
	// Happens on fresh installation.
	if len(data) == 0 {
		data = []byte("[]")
	}

	if err := json.Unmarshal(data, &meetings); err != nil {
		p.API.LogError(fmt.Sprintf("Unable to deserialize meetings data. Error: {%s}", err.Error()))
		return nil, errors.New(err.Error())
	}

	return *meetings, nil
}

func (p *Plugin) GetMeeting(meetingId string) (*dataStructs.MeetingRoom, error) {
	var meeting *dataStructs.MeetingRoom

	data, appErr := p.API.KVGet(prefixMeeting + meetingId)
	if appErr != nil {
		p.API.LogError(fmt.Sprintf("Unable to fetch meeting from KV store. Error: {%s}", appErr.Error()))
		return nil, appErr
	}

	_ = json.Unmarshal(data, &meeting)
	return meeting, nil
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
	hours := seconds / 3600
	seconds = seconds - 3600*hours

	minutes := seconds / 60
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
