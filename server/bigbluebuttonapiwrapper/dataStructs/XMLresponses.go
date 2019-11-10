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

//the following structs are the possible responses we can receive through
//api responses. Using func ReadXML we'll be able to turn our XML file to
// golang style structs

type CreateMeetingResponse struct {
	Returncode           string `xml:"returncode"`
	MeetingID            string `xml:"meetingID"`
	CreateTime           string `xml:"createTime"`
	AttendeePW           string `xml:"attendeePW"`
	ModeratorPW          string `xml:"moderatorPW"`
	HasBeenForciblyEnded string `xml:"hasBeenForciblyEnded"`
	MessageKey           string `xml:"messageKey"`
	Message              string `xml:"message"`
}

type IsMeetingRunningResponse struct {
	ReturnCode string `xml:"returncode"`
	Running    bool   `xml:"running"`
}

type EndResponse struct {
	ReturnCode string `xml:"returncode"`
	MessageKey string `xml:"messageKey"`
	Message    string `xml:"message"`
}

type GetMeetingsResponse struct {
	ReturnCode string      `xml:"returncode"`
	Meetings   allMeetings `xml:"meetings"`
}

type allMeetings struct {
	MeetingInfo []GetMeetingInfoResponse `xml:"meeting"`
}

type GetMeetingInfoResponse struct {
	ReturnCode            string    `xml:"returncode"`
	MeetingName           string    `xml:"meetingName"`
	MeetingID             string    `xml:"meetingID"`
	InternalMeetingID     string    `xml:"internalMeetingID"`
	CreateTime            string    `xml:"createTime"`
	CreateDate            string    `xml:"createDate"`
	VoiceBridge           string    `xml:"voiceBridge"`
	DialNumber            string    `xml:"dialNumber"`
	AttendeePW            string    `xml:"attendeePW"`
	ModeratorPW           string    `xml:"moderatorPW"`
	Running               bool      `xml:"running"`
	Duration              int       `xml:"duration"`
	HasUserJoined         bool      `xml:"hasUserJoined"`
	Recording             bool      `xml:"recording"`
	HasBeenForciblyEnded  bool      `xml:"hasBeenForciblyEnded"`
	StartTime             string    `xml:"startTime"`
	EndTime               string    `xml:"endTime"`
	ParticipantCount      int       `xml:"participantCount"`
	ListenerCount         int       `xml:"listenerCount"`
	VoiceParticipantCount int       `xml:"voiceParticipantCount"`
	VideoCount            int       `xml:"videoCount"`
	MaxUsers              int       `xml:"maxUsers"`
	ModeratorCount        int       `xml:"moderatorCount"`
	Attendees             attendees `xml:"attendees"`
	Metadata              string    `xml:"metadata"`
	MessageKey            string    `xml:"messageKey"`
	Message               string    `xml:"message"`
	//untested
	BreakoutRooms breakoutRooms `xml:"breakoutRooms"`
}

type breakoutRooms struct {
	BreakoutRooms []string `xml:"breakout"`
}

type attendees struct {
	Attendees []attendee `xml:"attendee"`
}

type attendee struct {
	UserID          string `xml:"userID"`
	FullName        string `xml:"fullName"`
	Role            string `xml:"role"`
	IsPresenter     bool   `xml:"isPresenter"`
	IsListeningOnly bool   `xml:"isListeningOnly"`
	HasJoinedVoice  bool   `xml:"hasJoinedVoice"`
	HasVideo        bool   `xml:"hasVideo"`
	Customdata      string `xml:"customdata"`
}

type GetRecordingsResponse struct {
	ReturnCode string     `xml:"returncode"`
	Recordings recordings `xml:"recordings"`
}

type recordings struct {
	Recording []recording `xml:"recording"`
}

type recording struct {
	RecordID     string   `xml:"recordID"`
	MeetingID    string   `xml:"meetingID"`
	Name         string   `xml:"name"`
	Published    string   `xml:"published"`
	State        string   `xml:"state"`
	StartTime    string   `xml:"startTime"`
	EndTime      string   `xml:"endTime"`
	Participants string   `xml:"participants"`
	MetaData     metadata `xml:"metadata"`
	Playback     struct {
		Format []struct {
			Type   string   `xml:"type"`
			Url    string   `xml:"url"`
			Length string   `xml:"length"`
			Images []string `xml:"preview>images>image"`
		} `xml:"format"`
	} `xml:"playback"`
}

type metadata struct {
	Title       string `xml:"title"`
	Subject     string `xml:"subject"`
	Description string `xml:"description"`
	Creator     string `xml:"creator"`
	Contributor string `xml:"contributor"`
	Language    string `xml:"language"`
}

type CreateWebhookResponse struct {
	Returncode string `xml:"returncode"`
	MessageKey string `xml:"messageKey"`
	Message    string `xml:"message"`
	HookID     string `xml:"hookID"`
}

type DestroyedWebhookResponse struct {
	Returncode string `xml:"returncode"`
	MessageKey string `xml:"messageKey"`
	Message    string `xml:"message"`
	Removed    string `xml:"removed"`
}

type PublishRecordingsResponse struct {
	ReturnCode string `xml:"returncode"`
	Published  string `xml:"published"`
}

type DeleteRecordingsResponse struct {
	ReturnCode string `xml:"returncode"`
	Deleted    string `xml:"deleted"`
}
