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

package dataStructs

import "errors"

//the following structs are the types we create to interact with the API
// ie participants, meetingRooms, recordings

type Recording struct {
	MeetingID string
	RecordID  string
	State     string
	Meta      string
	Publish   string
}

type Participants struct {
	IsAdmin_     int
	FullName_    string
	MeetingID_   string
	Password_    string
	CreateTime   string
	UserID       string
	WebVoiceConf string
	ConfigToken  string
	AvatarURL    string
	Redirect     string
	ClientURL    string
	JoinURL      string
}

func (p *Participants) IsValid() error {
	if p.FullName_ == "" {
		return errors.New("full name cannot be empty")
	}

	if p.MeetingID_ == "" {
		return errors.New("meeting ID cannot be empty")
	}

	if p.Password_ == "" {
		return errors.New("password cannot be empty")
	}

	return nil
}

type MeetingRoom struct {
	Name_                   string
	MeetingID_              string
	InternalMeetingId       string
	AttendeePW_             string
	ModeratorPW_            string
	Welcome                 string
	DialNumber              string
	VoiceBridge             string
	WebVoice                string
	LogoutURL               string
	Record                  string
	Duration                int
	Meta                    string
	ModeratorOnlyMessage    string
	AutoStartRecording      bool
	AllowStartStopRecording bool
	Created                 bool
	PostId                  string
	CreatedAt               int64
	EndedAt                 int64
	AttendeeNames           []string
	LoopCount               int
	ValidToken              string

	Meta_bn_recording_ready_url string //this needs to be properly url encoded
	Meta_channelid              string
	Meta_endcallbackurl         string

	CreateMeetingResponse CreateMeetingResponse
	MeetingInfo           GetMeetingInfoResponse
}

type WebHook struct {
	HookID      string
	CallBackURL string
	MeetingId   string

	WebhookResponse CreateWebhookResponse
}
