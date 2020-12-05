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
	"github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/server/mattermost"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	bbbAPI "github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/api"
	"github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/dataStructs"
	"github.com/mattermost/mattermost-server/model"
)

type RequestCreateMeetingJSON struct {
	UserId    string `json:"user_id"`
	ChannelId string `json:"channel_id"`
	Topic     string `json:"title"`
	Desc      string `json:"description"`
}

//Create meeting doesn't call the BBB api to start a meeting
//Only populates the meeting with details. Meeting is started when first person joins
func (p *Plugin) handleCreateMeeting(w http.ResponseWriter, r *http.Request) {

	// reads in information to create a meeting from client inside
	// whats being read in is the stuff in RequestCreateMeetingJSON
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	var request RequestCreateMeetingJSON
	json.Unmarshal(body, &request)

	meetingpointer := new(dataStructs.MeetingRoom)
	var err error
	if request.Topic == "" {
		err = p.PopulateMeeting(meetingpointer, nil, request.Desc, request.ChannelId)
	} else {
		err = p.PopulateMeeting(meetingpointer, []string{"create", request.Topic}, request.Desc, request.ChannelId)
	}

	if err != nil {
		http.Error(w, "Please provide a 'Site URL' in Settings > General > Configuration.", http.StatusUnprocessableEntity)
		return
	}

	//creates the start meeting post
	p.createStartMeetingPost(request.UserId, request.ChannelId, meetingpointer)

	// add our newly created meeting to our array of meetings
	p.Meetings = append(p.Meetings, *meetingpointer)

	w.WriteHeader(http.StatusOK)
}

type ButtonRequestJSON struct {
	UserId    string `json:"user_id"`
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
	json.Unmarshal(body, &request)
	meetingID := request.MeetingId
	meetingpointer := p.FindMeeting(meetingID)

	if meetingpointer == nil {
		myresp := ButtonResponseJSON{
			Url: "error",
		}
		userJson, _ := json.Marshal(myresp)
		w.Write(userJson)
		return
	} else {
		//check if meeting has actually been created and can be joined
		if !meetingpointer.Created {
			bbbAPI.CreateMeeting(meetingpointer)
			meetingpointer.Created = true
			var fullMeetingInfo dataStructs.GetMeetingInfoResponse
			bbbAPI.GetMeetingInfo(meetingID, meetingpointer.ModeratorPW_, &fullMeetingInfo) // this is used to get the InternalMeetingID
			meetingpointer.InternalMeetingId = fullMeetingInfo.InternalMeetingID
			meetingpointer.CreatedAt = time.Now().Unix()
		}

		user, _ := p.API.GetUser(request.UserId)
		username := user.Username

		//golang doesnt have sets so have to iterate through array to check if meeting participant is already in meeeting
		if !IsItemInArray(username, meetingpointer.AttendeeNames) {
			meetingpointer.AttendeeNames = append(meetingpointer.AttendeeNames, username)
		}

		var participant = dataStructs.Participants{} //set participant as an empty struct of type Participants
		participant.FullName_ = user.GetFullName()
		if len(participant.FullName_) == 0 {
			participant.FullName_ = user.Username
		}

		participant.MeetingID_ = meetingID

		post, appErr := p.API.GetPost(meetingpointer.PostId)
		if appErr != nil {
			http.Error(w, appErr.Error(), appErr.StatusCode)
			return
		}
		config := p.config()
		if (config.AdminOnly) {
			participant.Password_ = meetingpointer.AttendeePW_
			if(post.UserId == request.UserId ) {
				participant.Password_ = meetingpointer.ModeratorPW_ // the creator of a room is always moderator
			} else {
				for _, role := range user.GetRoles() {
					if role == "SYSTEM_ADMIN" || role == "TEAM_ADMIN" {
						participant.Password_ = meetingpointer.ModeratorPW_
						break
					}
				}
			}
		} else {
			participant.Password_ = meetingpointer.ModeratorPW_ //make everyone in channel a mod
		}
		joinURL, err := bbbAPI.GetJoinURL(&participant)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		myresp := ButtonResponseJSON{
			Url: joinURL,
		}
		userJson, _ := json.Marshal(myresp)

		Length, attendantsarray := GetAttendees(meetingID, meetingpointer.ModeratorPW_)
		// we immediately add our current attendee thats trying to join the meeting
		// to avoid the delay
		attendantsarray = append(attendantsarray, username)
		post.Props["user_count"] = Length + 1
		post.Props["attendees"] = strings.Join(attendantsarray, ",")

		if _, err := p.API.UpdatePost(post); err != nil {
			http.Error(w, err.Error(), err.StatusCode)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(userJson)
	}
}

//this method is responsible for updating meeting has ended inside mattermost when
// we end our meeting from inside BigBlueButton
func (p *Plugin) handleImmediateEndMeetingCallback(w http.ResponseWriter, r *http.Request, path string) {

	startpoint := len("/meetingendedcallback?")
	endpoint := strings.Index(path, "&")
	meetingid := path[startpoint:endpoint]
	validation := path[endpoint+1:]
	meetingpointer := p.FindMeeting(meetingid)
	if meetingpointer == nil || meetingpointer.ValidToken != validation {
		w.WriteHeader(http.StatusOK)
		return
	}
	post, err := p.API.GetPost(meetingpointer.PostId)
	if err != nil {
		http.Error(w, err.Error(), err.StatusCode)
		return
	}
	if meetingpointer.EndedAt == 0 {
		meetingpointer.EndedAt = time.Now().Unix()
	}
	p.MeetingsWaitingforRecordings = append(p.MeetingsWaitingforRecordings, *meetingpointer)
	post.Props["meeting_status"] = "ENDED"
	post.Props["attendents"] = strings.Join(meetingpointer.AttendeeNames, ",")
	timediff := meetingpointer.EndedAt - meetingpointer.CreatedAt
	durationstring := FormatSeconds(timediff)
	post.Props["duration"] = durationstring

	p.API.UpdatePost(post)

	w.WriteHeader(http.StatusOK)
}

//when user clicks endmeeting button inside Mattermost
func (p *Plugin) handleEndMeeting(w http.ResponseWriter, r *http.Request) {

	//for debugging
	mattermost.API.LogInfo("Processing End Meeting Request")

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var request ButtonRequestJSON
	json.Unmarshal(body, &request)
	meetingID := request.MeetingId
	meetingpointer := p.FindMeeting(meetingID)

	user, _ := p.API.GetUser(request.UserId)
	username := user.Username

	if meetingpointer == nil {
		myresp := model.PostActionIntegrationResponse{
			EphemeralText: "meeting has already ended",
		}
		userJson, _ := json.Marshal(myresp)
		w.Write(userJson)
		return
	} else {
		bbbAPI.EndMeeting(meetingpointer.MeetingID_, meetingpointer.ModeratorPW_)
		//for debugging
		mattermost.API.LogInfo("Meeting Ended")

		if meetingpointer.EndedAt == 0 {
			meetingpointer.EndedAt = time.Now().Unix()
		}
		p.MeetingsWaitingforRecordings = append(p.MeetingsWaitingforRecordings, *meetingpointer)

		post, err := p.API.GetPost(meetingpointer.PostId)
		if err != nil {
			http.Error(w, err.Error(), err.StatusCode)
			return
		}

		post.Props["meeting_status"] = "ENDED"
		post.Props["attendents"] = strings.Join(meetingpointer.AttendeeNames, ",")
		post.Props["ended_by"] = username
		timediff := meetingpointer.EndedAt - meetingpointer.CreatedAt
		if meetingpointer.CreatedAt == 0 {
			timediff = 0
		}
		durationstring := FormatSeconds(timediff)
		post.Props["duration"] = durationstring

		if _, err := p.API.UpdatePost(post); err != nil {
			http.Error(w, err.Error(), err.StatusCode)
			return
		}
		w.WriteHeader(http.StatusOK)
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
	json.Unmarshal(body, &request)
	meetingID := request.MeetingId

	resp, _ := bbbAPI.IsMeetingRunning(meetingID)

	myresp := isRunningResponseJSON{
		IsRunning: resp,
	}
	userJson, _ := json.Marshal(myresp)

	w.Header().Set("Content-Type", "application/json")
	w.Write(userJson)

}

// type WebHookRequestJSON struct {
// 	Header struct {
// 		Timestamp   string `json:"timestamp"`
// 		Name        string `json:"name"`
// 		CurrentTime string `json:"current_time"`
// 		Version     string `json:"version"`
// 	} `json:"header"`
// 	Payload struct {
// 		MeetingId string `json:"meeting_id"`
// 	} `json:"payload"`
// }
//
// type WebHookResponseEncoded struct {
// 	Payload map[string]interface{} `form:"payload"`
// }

//webhook to send additional information about meeting that had ended
//has a 4-5 minute delay which is why handleImmediateEndMeetingCallback() is
//used instead for updating end meeting post on Mattermost. Keeping it here in
//case bbbserver is not up to date with the immediate meeting ended callback feature
// func (p *Plugin) handleWebhookMeetingEnded(w http.ResponseWriter, r *http.Request) {
//
// 	out := ""
// 	r.ParseForm()
// 	for key, value := range r.Form {
// 		out += fmt.Sprintf("%s = %s\n", key, value)
// 	}
// 	events := (r.FormValue("event"))
//
// 	internal_meetingid := events[strings.Index(events, "\""+"meeting_id"+"\"")+14:]
// 	internal_meetingid = internal_meetingid[:strings.IndexByte(internal_meetingid, '"')]
//
// 	meetingpointer := p.FindMeetingfromInternal(internal_meetingid)
//
// 	if meetingpointer == nil {
// 		w.WriteHeader(http.StatusOK)
// 		return
// 	}
//
// 	postid := meetingpointer.PostId
// 	if postid == "" {
// 		panic("no post id found")
// 	}
// 	post, err := p.API.GetPost(postid)
// 	if err != nil {
// 		http.Error(w, err.Error(), err.StatusCode)
// 		return
// 	}
// 	if meetingpointer.EndedAt == 0 {
// 		meetingpointer.EndedAt = time.Now().Unix()
// 	}
// 	p.MeetingsWaitingforRecordings = append(p.MeetingsWaitingforRecordings, *meetingpointer)
// 	post.Props["meeting_status"] = "ENDED"
// 	post.Props["attendents"] = strings.Join(meetingpointer.AttendeeNames, ",")
// 	timediff := meetingpointer.EndedAt - meetingpointer.CreatedAt
// 	durationstring := FormatSeconds(timediff)
// 	post.Props["duration"] = durationstring
//
// 	if _, err := p.API.UpdatePost(post); err != nil {
// 		http.Error(w, err.Error(), err.StatusCode)
// 		return
// 	}
//
// 	w.WriteHeader(http.StatusOK)
// }

// type MyCustomClaims struct {
// 	MeetingID string `json:"meeting_id"`
// 	RecordID  string `json:"record_id"`
// 	jwt.StandardClaims
// }

func (p *Plugin) handleRecordingReady(w http.ResponseWriter, r *http.Request) {
	// p.API.LogDebug("handleRecordingReady reached")
	// r.ParseForm()
	// parameters := (r.FormValue("signed_parameters"))
	// token, _ := jwt.ParseWithClaims(parameters, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
	// 	return []byte("AllYourBase"), nil
	// })
	// claims, _ := token.Claims.(*MyCustomClaims)
	// meetingid := claims.MeetingID
	// recordid := claims.RecordID
	// p.API.LogDebug(meetingid + " " + recordid)
	// recordingsresponse, _ := bbbAPI.GetRecordings(meetingid, recordid, "")
	// if recordingsresponse.ReturnCode != "SUCCESS" {
	// 	w.WriteHeader(http.StatusOK)
	// 	return
	// }
	//
	// meetingpointer := p.FindMeeting(meetingid)
	//
	// if meetingpointer == nil {
	// 	w.WriteHeader(http.StatusOK)
	// 	return
	// }
	//
	// postid := meetingpointer.PostId
	// if postid == "" {
	// 	panic("no post id found")
	// }
	// post, err := p.API.GetPost(postid)
	// if err != nil {
	// 	http.Error(w, err.Error(), err.StatusCode)
	// 	return
	// }
	//
	// post.Message = "#BigBlueButton #" + meetingpointer.Name_ + " #" + recordid + " #recording" + " recordings"
	// post.Props["recording_status"] = "COMPLETE"
	// post.Props["is_published"] = "true"
	// post.Props["record_id"] = recordid
	// post.Props["recording_url"] = recordingsresponse.Recordings.Recording[0].Playback.Format[0].Url
	//
	// post.Props["images"] = strings.Join(recordingsresponse.Recordings.Recording[0].Playback.Format[0].Images, ",")
	// if _, err := p.API.UpdatePost(post); err != nil {
	// 	http.Error(w, err.Error(), err.StatusCode)
	// 	return
	// }

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
	json.Unmarshal(body, &request)
	meetingID := request.MeetingId
	meetingpointer := p.FindMeeting(meetingID)
	if meetingpointer == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	postid := meetingpointer.PostId
	if postid == "" {
		w.WriteHeader(http.StatusOK)
		return
	}
	post, err := p.API.GetPost(postid)
	if err != nil {
		http.Error(w, err.Error(), err.StatusCode)
		return
	}
	Length, Array := GetAttendees(meetingID, meetingpointer.ModeratorPW_)
	post.Props["user_count"] = Length
	post.Props["attendees"] = strings.Join(Array, ",")

	if _, err := p.API.UpdatePost(post); err != nil {
		http.Error(w, err.Error(), err.StatusCode)
		return
	}
	myresp := AttendeesResponseJSON{
		Num:       Length,
		Attendees: Array,
	}
	userJson, _ := json.Marshal(myresp)

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

	var request PublishRecordingsRequestJSON
	json.Unmarshal(body, &request)
	recordid := request.RecordId
	publish := request.Publish

	meetingpointer := p.FindMeeting(request.MeetingId)
	if meetingpointer == nil {
		http.Error(w, "Error: Cannot find the meeting_id for the recording, MeetingID#"+request.MeetingId, http.StatusForbidden)
		return
	}

	if _, err := bbbAPI.PublishRecordings(recordid, publish); err != nil {
		http.Error(w, "Error: Recording not found", http.StatusForbidden)
		return
	}

	post, appErr := p.API.GetPost(meetingpointer.PostId)
	if appErr != nil {
		http.Error(w, "Error: cannot find the post message for this recording \n"+appErr.Error(), appErr.StatusCode)
		return
	}

	post.Props["is_published"] = publish

	if _, err := p.API.UpdatePost(post); err != nil {
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

	var request DeleteRecordingsRequestJSON
	json.Unmarshal(body, &request)
	recordid := request.RecordId

	if _, err := bbbAPI.DeleteRecordings(recordid); err != nil {
		http.Error(w, "Error: Recording not found", http.StatusForbidden)
		return
	}

	meetingID := request.MeetingId
	meetingpointer := p.FindMeeting(meetingID)
	if meetingpointer == nil {
		http.Error(w, "Error: Cannot find the meeting_id for the recording", http.StatusForbidden)
		return
	}

	post, appErr := p.API.GetPost(meetingpointer.PostId)
	if appErr != nil {
		http.Error(w, "Error: cannot find the post message for this recording \n"+appErr.Error(), appErr.StatusCode)
		return
	}

	post.Props["is_deleted"] = "true"
	post.Props["record_status"] = "Recording Deleted"
	if _, err := p.API.UpdatePost(post); err != nil {
		http.Error(w, "Error: could not update post \n"+err.Error(), err.StatusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (p *Plugin) Loopthroughrecordings() {

	for i := 0; i < len(p.MeetingsWaitingforRecordings); i++ {
		Meeting := p.MeetingsWaitingforRecordings[i]
		// TODO Harshil Sharma: explore better alternative of waiting for specific count of re-tries
		// instead of duration of re-tries.
		if Meeting.LoopCount > 144 {
			p.MeetingsWaitingforRecordings = append(p.MeetingsWaitingforRecordings[:i], p.MeetingsWaitingforRecordings[i+1:]...)
			i--
			continue
		}
		Meeting.LoopCount++

		recordingsresponse, _, _ := bbbAPI.GetRecordings(Meeting.MeetingID_, "", "")
		if recordingsresponse.ReturnCode == "SUCCESS" {
			if len(recordingsresponse.Recordings.Recording) > 0 {
				postid := Meeting.PostId
				if postid != "" {
					post, _ := p.API.GetPost(postid)
					post.Message = "#BigBlueButton #" + Meeting.Name_ + " #" + Meeting.MeetingID_ + " #recording" + " recordings"
					post.Props["recording_status"] = "COMPLETE"
					post.Props["is_published"] = "true"
					post.Props["record_id"] = recordingsresponse.Recordings.Recording[0].RecordID
					post.Props["recording_url"] = recordingsresponse.Recordings.Recording[0].Playback.Format[0].Url
					post.Props["images"] = strings.Join(recordingsresponse.Recordings.Recording[0].Playback.Format[0].Images, ",")

					if _, err := p.API.UpdatePost(post); err == nil {
						p.MeetingsWaitingforRecordings = append(p.MeetingsWaitingforRecordings[:i], p.MeetingsWaitingforRecordings[i+1:]...)
						i--
					}
				}
			}
		}
	}
}
