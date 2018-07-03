package dataStructs

//the following structs are the types we create to interact with the API
// ie participants, meetingRooms, recordings


type Recording struct {
  MeetingID     string
  RecordID      string
  State         string
  Meta          string
  Publish       string
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

  Meta_bn_recording_ready_url string //this needs to be properly url encoded
  Meta_channelid          string
  Meta_endcallbackurl     string

  CreateMeetingResponse CreateMeetingResponse
	MeetingInfo           GetMeetingInfoResponse

}

type WebHook struct {
  HookID      string
  CallBackURL string
  MeetingId   string

  WebhookResponse CreateWebhookResponse
}
