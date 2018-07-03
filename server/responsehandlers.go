package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	bbbAPI "github.com/mattermost/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/api"
	"github.com/mattermost/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/dataStructs"
	"github.com/mattermost/mattermost-server/model"
)

type RequestCreateMeetingJSON struct {
	User_id    string `json:"user_id"`
	Channel_id string `json:"channel_id"`
	Topic      string `json:"title"`
	Desc       string `json:"description"`
}

func (p *Plugin) handleCreateMeeting(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	var request RequestCreateMeetingJSON
	json.Unmarshal([]byte(body), &request)

	meetingpointer := new(dataStructs.MeetingRoom) //this is actually a pointer but golang deallocate pointers for you
	if request.Topic == "" {

		p.PopulateMeeting(meetingpointer, nil, request.Desc)
	} else {
		p.PopulateMeeting(meetingpointer, []string{"create", request.Topic}, request.Desc)
	}

	p.createStartMeetingPost(request.User_id, request.Channel_id, meetingpointer)

	//this doesnt really do anything here
	//sends it to the client side startmeeting
	myresp := RequestCreateMeetingJSON{
		User_id:    request.User_id,
		Channel_id: request.Channel_id,
	}

	userJson, err2 := json.Marshal(myresp)

	//TODO: Use key value score instead
	p.Meetings = append(p.Meetings, *meetingpointer)

	if err2 != nil {
		panic(err2)
	}

	w.Write(userJson)

}

//TODO: reconfigure this
type ButtonRequestJSON struct {
	User_id   string `json:"user_id"`
	MeetingId string `json:"meeting_id"`
	IsMod     string `json:"is_mod"`
}
type ButtonResponseJSON struct {
	Url string `json:"url"`
}

func (p *Plugin) handleJoinMeeting(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var request ButtonRequestJSON
	json.Unmarshal([]byte(body), &request)
	meetingID := request.MeetingId
	meetingpointer := p.FindMeeting(meetingID)

	if meetingpointer == nil {
		myresp := ButtonResponseJSON{
			Url: "error",
		}
		userJson, err2 := json.Marshal(myresp)
		if err2 != nil {
			panic(err2)
		}
		w.Write(userJson)
		return
	} else {
		if !meetingpointer.Created {
			bbbAPI.CreateMeeting(meetingpointer)
			meetingpointer.Created = true
			var fullMeetingInfo dataStructs.GetMeetingInfoResponse
			bbbAPI.GetMeetingInfo(meetingID, meetingpointer.ModeratorPW_, &fullMeetingInfo)
			meetingpointer.InternalMeetingId = fullMeetingInfo.InternalMeetingID
			meetingpointer.CreatedAt = time.Now().Unix()

			//TEMP
			p.LoopingMeetings = append(p.LoopingMeetings, *meetingpointer)
		}

		user, _ := p.api.GetUser(request.User_id)
		username := user.Username
		meetingpointer.AttendeeNames = append(meetingpointer.AttendeeNames, username)
		var participant = dataStructs.Participants{}
		participant.FullName_ = username
		participant.MeetingID_ = meetingID
		// if request.IsMod == "TRUE" {
		// 	participant.Password_ = meetingpointer.ModeratorPW_
		// } else {
		// 	participant.Password_ = meetingpointer.AttendeePW_
		// }
		participant.Password_ = meetingpointer.ModeratorPW_ //make everyone in channel a mod
		joinURL := bbbAPI.GetJoinURL(&participant)
		// var text string
		// text = "[Join here!](" + joinURL + ")" + "   for testing purposes: " + meetingpointer.InternalMeetingId
		myresp := ButtonResponseJSON{
			Url: joinURL,
		}

		userJson, err2 := json.Marshal(myresp)
		if err2 != nil {
			panic(err2)
		}

		postid := meetingpointer.PostId
		if postid == "" {
			panic("no post id found")
		}
		//TODO: update kvp

		post, err := p.api.GetPost(postid)
		if err != nil {
			http.Error(w, err.Error(), err.StatusCode)
			return
		}
		Length, attendantsarray := GetAttendees(meetingID, meetingpointer.ModeratorPW_)
		attendantsarray = append(attendantsarray, username)
		post.Props["user_count"] = Length + 1
		post.Props["attendees"] = strings.Join(attendantsarray, ",")

		if _, err := p.api.UpdatePost(post); err != nil {
			http.Error(w, err.Error(), err.StatusCode)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(userJson)
	}
}
func (p *Plugin) handleAnotherCallbackTry(w http.ResponseWriter, r *http.Request, path string) {

	meetingid := path[22:] //22 length of "/meetingendedcallback?"

	meetingpointer := p.FindMeeting(meetingid)
	if meetingpointer == nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	postid := meetingpointer.PostId
	if postid == "" {
		panic("no post id found")
	}
	post, err := p.api.GetPost(postid)
	if err != nil {
		http.Error(w, err.Error(), err.StatusCode)
		return
	}
	if meetingpointer.EndedAt == 0 {
		meetingpointer.EndedAt = time.Now().Unix()
	}
	// p.api.DeletePost(postid)
	//post.Message = "Meeting has ended."
	post.Props["meeting_status"] = "ENDED"
	post.Props["attendents"] = strings.Join(meetingpointer.AttendeeNames, ",")
	timediff := meetingpointer.EndedAt - meetingpointer.CreatedAt
	durationstring := FormatSeconds(timediff)
	post.Props["duration"] = durationstring

	p.api.UpdatePost(post)

	w.WriteHeader(http.StatusOK)
}

func (p *Plugin) handleEndMeeting(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var request ButtonRequestJSON
	json.Unmarshal([]byte(body), &request)
	meetingID := request.MeetingId //TODO: do stuff with meetingID
	meetingpointer := p.FindMeeting(meetingID)

	user, _ := p.api.GetUser(request.User_id)
	username := user.Username

	if meetingpointer == nil {
		myresp := model.PostActionIntegrationResponse{
			EphemeralText: "meeting has already ended",
		}
		userJson, err2 := json.Marshal(myresp)
		if err2 != nil {
			panic(err2)
		}
		w.Write(userJson)
		return
	} else {
		bbbAPI.EndMeeting(meetingpointer.MeetingID_, meetingpointer.ModeratorPW_)

		if meetingpointer.EndedAt == 0 {
			meetingpointer.EndedAt = time.Now().Unix()
		}

		//p.DeleteMeeting(meetingpointer.MeetingID_) //only delete it from our memory when webhook calls back
		text := "meeting deleted"
		myresp := model.PostActionIntegrationResponse{
			EphemeralText: text,
		}

		userJson, err2 := json.Marshal(myresp)
		if err2 != nil {
			panic(err2)
		}
		postid := meetingpointer.PostId
		if postid == "" {
			panic("no post id found")
		}
		post, err := p.api.GetPost(postid)
		if err != nil {
			http.Error(w, err.Error(), err.StatusCode)
			return
		}

		//post.Message = "Meeting has ended."
		post.Props["meeting_status"] = "ENDED"
		post.Props["attendents"] = strings.Join(meetingpointer.AttendeeNames, ",")
		post.Props["ended_by"] = username
		timediff := meetingpointer.EndedAt - meetingpointer.CreatedAt
		if meetingpointer.CreatedAt == 0 {
			timediff = 0
		}
		durationstring := FormatSeconds(timediff)
		post.Props["duration"] = durationstring

		if _, err := p.api.UpdatePost(post); err != nil {
			http.Error(w, err.Error(), err.StatusCode)
			return
		}
		w.Write(userJson)
	}
}

type isRunningRequestJSON struct {
	MeetingId string `json:"meeting_id"`
}
type isRunningResponseJSON struct {
	IsRunning bool `json:"running"`
}

func (p *Plugin) handleIsMeetingRunning(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var request ButtonRequestJSON
	json.Unmarshal([]byte(body), &request)
	meetingID := request.MeetingId //TODO: do stuff with meetingID

	resp := bbbAPI.IsMeetingRunning(meetingID)

	myresp := isRunningResponseJSON{
		IsRunning: resp,
	}
	userJson, err2 := json.Marshal(myresp)
	if err2 != nil {
		panic(err2)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(userJson)

}

type WebHookRequestJSON struct {
	Header struct {
		Timestamp   string `json:"timestamp"`
		Name        string `json:"name"`
		CurrentTime string `json:"current_time"`
		Version     string `json:"version"`
	} `json:"header"`
	Payload struct {
		MeetingId string `json:"meeting_id"`
	} `json:"payload"`
}

type WebHookResponseEncoded struct {
	Payload map[string]interface{} `form:"payload"`
}

func (p *Plugin) handleWebhookMeetingEnded(w http.ResponseWriter, r *http.Request) {

	out := ""
	r.ParseForm()
	for key, value := range r.Form {
		out += fmt.Sprintf("%s = %s\n", key, value)
	}
	events := (r.FormValue("event"))

	internal_meetingid := events[strings.Index(events, "\""+"meeting_id"+"\"")+14:]
	internal_meetingid = internal_meetingid[:strings.IndexByte(internal_meetingid, '"')]

	meetingpointer := p.FindMeetingfromInternal(internal_meetingid)

	if meetingpointer == nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	//
	postid := meetingpointer.PostId
	if postid == "" {
		panic("no post id found")
	}
	post, err := p.api.GetPost(postid)
	if err != nil {
		http.Error(w, err.Error(), err.StatusCode)
		return
	}
	if meetingpointer.EndedAt == 0 {
		meetingpointer.EndedAt = time.Now().Unix()
	}
	// p.api.DeletePost(postid)
	//post.Message = "Meeting has ended."
	post.Props["meeting_status"] = "ENDED"
	post.Props["attendents"] = strings.Join(meetingpointer.AttendeeNames, ",")
	timediff := meetingpointer.EndedAt - meetingpointer.CreatedAt
	durationstring := FormatSeconds(timediff)
	post.Props["duration"] = durationstring

	if _, err := p.api.UpdatePost(post); err != nil {
		http.Error(w, err.Error(), err.StatusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type MyCustomClaims struct {
	MeetingID string `json:"meeting_id"`
	RecordID  string `json:"record_id"`
	jwt.StandardClaims
}

func (p *Plugin) handleRecordingReady(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	parameters := (r.FormValue("signed_parameters"))
	token, _ := jwt.ParseWithClaims(parameters, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil
	})
	claims, _ := token.Claims.(*MyCustomClaims)
	meetingid := claims.MeetingID
	recordid := claims.RecordID

	recordingsresponse, _ := bbbAPI.GetRecordings(meetingid, recordid, "")
	if recordingsresponse.ReturnCode != "SUCCESS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	meetingpointer := p.FindMeeting(meetingid)

	if meetingpointer == nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	//
	postid := meetingpointer.PostId
	if postid == "" {
		panic("no post id found")
	}
	post, err := p.api.GetPost(postid)
	if err != nil {
		http.Error(w, err.Error(), err.StatusCode)
		return
	}
	// p.api.DeletePost(postid)
	post.Message = "#BigBlueButton #" + meetingpointer.Name_ + " #" + recordid + " #recording" + " recordings"
	post.Props["recording_status"] = "COMPLETE"
	post.Props["is_published"] = "true"
	post.Props["record_id"] = recordid
	post.Props["recording_url"] = recordingsresponse.Recordings.Recording[0].Playback.Format[0].Url
	//	post.Props["images"] = recordingsresponse.Recordings.Recording[0].Playback.Format[0].Images
	post.Props["images"] = strings.Join(recordingsresponse.Recordings.Recording[0].Playback.Format[0].Images, ",")
	// post.Props["image1"] = recordingsresponse.Recordings.Recording[0].Playback.Format[0].Images[0]
	// post.Props["image2"] = recordingsresponse.Recordings.Recording[0].Playback.Format[0].Images[1]
	// post.Props["image3"] = recordingsresponse.Recordings.Recording[0].Playback.Format[0].Images[2]
	// ImageArr := []string{Image1, Image2, Image3}
	//
	// post.Props["images"] = ImageArr

	if _, err := p.api.UpdatePost(post); err != nil {
		http.Error(w, err.Error(), err.StatusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

type AttendeesRequestJSON struct {
	MeetingId string `json:"meeting_id"`
}
type AttendeesResponseJSON struct {
	Num       int      `json:"num"`
	Attendees []string `json:"attendees"`
}

func (p *Plugin) handleGetAttendeesInfo(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var request AttendeesRequestJSON
	json.Unmarshal([]byte(body), &request)
	meetingID := request.MeetingId //TODO: do stuff with meetingID
	meetingpointer := p.FindMeeting(meetingID)
	if meetingpointer == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	postid := meetingpointer.PostId
	if postid == "" {
		panic("no post id found")
	}
	post, err := p.api.GetPost(postid)
	if err != nil {
		http.Error(w, err.Error(), err.StatusCode)
		return
	}
	Length, Array := GetAttendees(meetingID, meetingpointer.ModeratorPW_)
	post.Props["user_count"] = Length
	post.Props["attendees"] = strings.Join(Array, ",")

	if _, err := p.api.UpdatePost(post); err != nil {
		http.Error(w, err.Error(), err.StatusCode)
		return
	}
	myresp := AttendeesResponseJSON{
		Num:       Length,
		Attendees: Array,
	}
	userJson, err2 := json.Marshal(myresp)
	if err2 != nil {
		panic(err2)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(userJson)
}

type RecordingsRequestJSON struct {
	ChannelId string `json:"channel_id"`
}
type RecordingsResponseJSON struct {
	Recordings []SingleRecording `json:"recordings"`
}
type SingleRecording struct {
	RecordingUrl string `json:"recordingurl"`
	Title        string `json:"title"`
}

type PublishRecordingsRequestJSON struct {
	RecordId  string `json:"record_id"`
	Publish   string `json:"publish"` //string  true or false
	MeetingId string `json:"meeting_id"`
}

func (p *Plugin) handlePublishRecordings(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	//what if we just serve through xml?

	var request PublishRecordingsRequestJSON
	json.Unmarshal([]byte(body), &request)
	recordid := request.RecordId
	publish := request.Publish

	publishrecordingsresponse := bbbAPI.PublishRecordings(recordid, publish)

	if publishrecordingsresponse.ReturnCode != "SUCCESS" {
		http.Error(w, "Forbidden", http.StatusForbidden) //prob should change this error status
		return
	}

	meetingID := request.MeetingId
	meetingpointer := p.FindMeeting(meetingID)
	if meetingpointer == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	postid := meetingpointer.PostId
	if postid == "" {
		panic("no post id found")
	}
	post, err := p.api.GetPost(postid)
	if err != nil {
		http.Error(w, err.Error(), err.StatusCode)
		return
	}

	post.Props["is_published"] = publish

	if _, err := p.api.UpdatePost(post); err != nil {
		http.Error(w, err.Error(), err.StatusCode)
		return
	}
	//update post props with new recording  status
	w.WriteHeader(http.StatusOK)
}

type DeleteRecordingsRequestJSON struct {
	RecordId  string `json:"record_id"`
	MeetingId string `json:"meeting_id"`
}

func (p *Plugin) handleDeleteRecordings(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	//what if we just serve through xml?

	var request DeleteRecordingsRequestJSON
	json.Unmarshal([]byte(body), &request)
	recordid := request.RecordId

	deleterecordingsresponse := bbbAPI.DeleteRecordings(recordid)

	if deleterecordingsresponse.ReturnCode != "SUCCESS" {
		http.Error(w, "Forbidden", http.StatusForbidden) //prob should change this error status
		return
	}

	meetingID := request.MeetingId
	meetingpointer := p.FindMeeting(meetingID)
	if meetingpointer == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	postid := meetingpointer.PostId
	if postid == "" {
		panic("no post id found")
	}
	post, err := p.api.GetPost(postid)
	if err != nil {
		http.Error(w, err.Error(), err.StatusCode)
		return
	}

	post.Props["is_deleted"] = "true" //careful when setting true like this
	post.Props["record_status"] = "Recording Deleted"
	if _, err := p.api.UpdatePost(post); err != nil {
		http.Error(w, err.Error(), err.StatusCode)
		return
	}

	//update post props with new recording  status
	w.WriteHeader(http.StatusOK)
}

func (p *Plugin) Loopthroughrecordings() {

	for i := 0; i < len(p.LoopingMeetings); i++ {
		Meeting := p.LoopingMeetings[i]
		if Meeting.LoopCount > 144 {
			p.LoopingMeetings = append(p.LoopingMeetings[:i], p.LoopingMeetings[i+1:]...)
			i--
			continue
		}

		recordingsresponse, _ := bbbAPI.GetRecordings(Meeting.MeetingID_, "", "")
		if recordingsresponse.ReturnCode == "SUCCESS" {
			if len(recordingsresponse.Recordings.Recording) > 0 {
				postid := Meeting.PostId
				if postid != "" {
					post, _ := p.api.GetPost(postid)
					post.Message = "#BigBlueButton #" + Meeting.Name_ + " #" + Meeting.MeetingID_ + " #recording" + " recordings"
					post.Props["recording_status"] = "COMPLETE"
					post.Props["is_published"] = "true"
					post.Props["record_id"] = recordingsresponse.Recordings.Recording[0].RecordID
					post.Props["recording_url"] = recordingsresponse.Recordings.Recording[0].Playback.Format[0].Url
					post.Props["images"] = strings.Join(recordingsresponse.Recordings.Recording[0].Playback.Format[0].Images, ",")
					// post.Props["image1"] = recordingsresponse.Recordings.Recording[0].Playback.Format[0].Images[0]
					// post.Props["image2"] = recordingsresponse.Recordings.Recording[0].Playback.Format[0].Images[1]
					// post.Props["image3"] = recordingsresponse.Recordings.Recording[0].Playback.Format[0].Images[2]

					if _, err := p.api.UpdatePost(post); err != nil {
						p.LoopingMeetings = append(p.LoopingMeetings[:i], p.LoopingMeetings[i+1:]...)
						i--
					}

				}
			}
		}
	}
}
