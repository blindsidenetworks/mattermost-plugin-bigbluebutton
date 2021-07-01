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
	"fmt"
	"github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/server/mattermost"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	bbbAPI "github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/api"
	"github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/dataStructs"
	"github.com/mattermost/mattermost-server/v5/model"
)

type RequestCreateMeetingJSON struct {
	UserId    string `json:"user_id"`
	ChannelId string `json:"channel_id"`
	Topic     string `json:"title"`
	Desc      string `json:"description"`
}

type ButtonRequestJSON struct {
	UserId    string `json:"user_id"`
	MeetingId string `json:"meeting_id"`
	IsMod     string `json:"is_mod"`
}

type ButtonResponseJSON struct {
	Url string `json:"url"`
}

type AttendeesRequestJSON struct {
	MeetingId string `json:"meeting_id"`
}

type AttendeesResponseJSON struct {
	Num       int      `json:"num"`
	Attendees []string `json:"attendees"`
}

type PublishRecordingsRequestJSON struct {
	RecordId  string `json:"record_id"`
	Publish   string `json:"publish"` //string  true or false
	MeetingId string `json:"meeting_id"`
}

func (p *Plugin) Loopthroughrecordings() {
	meetingsWaitingforRecordings, err := p.GetRecordingWaitingList()
	if err != nil {
		return
	}

	for _, meetingID := range meetingsWaitingforRecordings {
		Meeting, err := p.GetMeetingWaitingForRecording(meetingID)
		if err != nil {
			continue
		}

		// TODO Harshil Sharma: explore better alternative of waiting for specific count of re-tries
		// instead of duration of re-tries.
		if Meeting.LoopCount > 144 {
			_ = p.RemoveMeetingWaitingForRecording(Meeting.MeetingID_)
			continue
		}

		recordingsresponse, _, _ := bbbAPI.GetRecordings(Meeting.MeetingID_, "", "")
		if recordingsresponse.ReturnCode == "SUCCESS" {
			if len(recordingsresponse.Recordings.Recording) > 0 {
				recordings, err := json.Marshal(recordingsresponse.Recordings.Recording)
				if err != nil {
					p.API.LogError(err.Error())
				} else {
					p.API.LogInfo(string(recordings))
				}

				postid := Meeting.PostId
				if postid != "" {
					post, _ := p.API.GetPost(postid)
					// TODO
					post.Message = "#BigBlueButton #" + Meeting.Name_ + " #" + Meeting.MeetingID_ + " #recording" + " recordings"
					post.AddProp("recording_status", "COMPLETE")
					post.AddProp("is_published", "true")

					attachments := post.Attachments()
					recordingsAdded := false

					for _, playback := range recordingsresponse.Recordings.Recording[0].Playback.Format {
						switch playback.Type {
						case "presentation":
							post.AddProp("record_id", recordingsresponse.Recordings.Recording[0].RecordID)
							post.AddProp("recording_url", playback.Url)
							post.AddProp("images", strings.Join(playback.Images, ", "))

							attachments[0].Fields = append(attachments[0].Fields, &model.SlackAttachmentField{
								Title: "Recordings",
								Value: fmt.Sprintf("[Click to view recordings](%s)", playback.Url),
								Short: true,
							})

							recordingsAdded = true
						case "notes":
							post.AddProp("notes", true)
							post.AddProp("notes_url", playback.Url)

							attachments[0].Fields = append(attachments[0].Fields, &model.SlackAttachmentField{
								Title: "Notes",
								Value: fmt.Sprintf("[Click to view notes](%s)", playback.Url),
							})

							recordingsAdded = true
						}
					}

					if recordingsAdded {
						meetingAttachment := attachments[0]
						attachments = []*model.SlackAttachment{
							meetingAttachment,
							{
								Actions: []*model.PostAction{
									{
										Id:    "toggle_recording_visibility",
										Type:  "button",
										Name:  "Make Recording Invisible",
										Style: "secondary",
										Integration: &model.PostActionIntegration{
											URL: "/plugins/bigbluebutton/publishrecordings",
											Context: map[string]interface{}{
												"publish":    "false",
												"meeting_id": meetingID,
												"record_id":  recordingsresponse.Recordings.Recording[0].RecordID,
											},
										},
									},
									{
										Id:    "delete_recordings",
										Type:  "button",
										Name:  "Delete Recordings",
										Style: "danger",
										Integration: &model.PostActionIntegration{
											URL: "/plugins/bigbluebutton/deleterecordings",
											Context: map[string]interface{}{
												"meeting_id": meetingID,
												"record_id":  recordingsresponse.Recordings.Recording[0].RecordID,
											},
										},
									},
								},
							},
						}
					}

					model.ParseSlackAttachment(post, attachments)

					if _, err := p.API.UpdatePost(post); err == nil {
						_ = p.RemoveMeetingWaitingForRecording(Meeting.MeetingID_)
					}
				}
			}
		}
	}
}

type DeleteRecordingsRequestJSON struct {
	RecordId  string `json:"record_id"`
	MeetingId string `json:"meeting_id"`
}

//Create meeting doesn't call the BBB api to start a meeting
//Only populates the meeting with details. Meeting is started when first person joins
func (p *Plugin) handleCreateMeeting(w http.ResponseWriter, r *http.Request) {

	// reads in information to create a meeting from client inside
	// whats being read in is the stuff in RequestCreateMeetingJSON
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	var request RequestCreateMeetingJSON
	_ = json.Unmarshal(body, &request)

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
	if err := p.SaveMeeting(meetingpointer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (p *Plugin) handleJoinMeeting(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var request *model.PostActionIntegrationRequest
	if err := json.Unmarshal(body, &request); err != nil {
		p.API.LogError("Error occured unmarshaling join meeting request body. Error: " + err.Error())
		return
	}

	//var request ButtonRequestJSON
	//_ = json.Unmarshal(body, &request)
	//meetingID := request.MeetingId

	meetingID := request.Context["meetingId"].(string)
	meetingpointer := p.FindMeeting(meetingID)

	if meetingpointer == nil {
		myresp := ButtonResponseJSON{
			Url: "error",
		}
		userJson, _ := json.Marshal(myresp)
		_, _ = w.Write(userJson)
		return
	} else {
		//check if meeting has actually been created and can be joined
		if !meetingpointer.Created {
			if _, err := bbbAPI.CreateMeeting(meetingpointer); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			meetingpointer.Created = true
			var fullMeetingInfo dataStructs.GetMeetingInfoResponse

			// this is used to get the InternalMeetingID
			if _, err := bbbAPI.GetMeetingInfo(meetingID, meetingpointer.ModeratorPW_, &fullMeetingInfo); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			meetingpointer.InternalMeetingId = fullMeetingInfo.InternalMeetingID
			meetingpointer.CreatedAt = time.Now().Unix()
		}

		user, _ := p.API.GetUser(request.UserId)
		username := user.Username

		//golang doesnt have sets so have to iterate through array to check if meeting participant is already in meeeting
		if !IsItemInArray(username, meetingpointer.AttendeeNames) {
			meetingpointer.AttendeeNames = append(meetingpointer.AttendeeNames, username)
		}

		if err := p.SaveMeeting(meetingpointer); err != nil {
			p.API.LogError("Error occurred updating meeting info in handleJoinMeeting. Error: " + err.Error())
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
		if config.AdminOnly {
			participant.Password_ = meetingpointer.AttendeePW_
			if post.UserId == request.UserId {
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
		post.AddProp("user_count", Length+1)
		post.AddProp("attendees", strings.Join(attendantsarray, ","))

		slackAttachments := post.Attachments()
		slackAttachments[0].Fields[0].Title = fmt.Sprintf("Attendees (%d)", Length+1)
		slackAttachments[0].Fields[0].Value = "@" + strings.Join(attendantsarray, ", @")

		model.ParseSlackAttachment(post, slackAttachments)

		if _, err := p.API.UpdatePost(post); err != nil {
			p.API.LogError("Error occurred updating meeting post during user joining it. Error: " + err.Error())
			http.Error(w, err.Error(), err.StatusCode)
			return
		}

		p.API.SendEphemeralPost(request.UserId, &model.Post{
			ChannelId: request.ChannelId,
			Type:      model.POST_EPHEMERAL,
			Message:   fmt.Sprintf("Join the BBB meeting [here](%s)", joinURL),
		})

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(userJson)
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
	if err := p.AddMeetingWaitingForRecording(meetingpointer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	post.Props["meeting_status"] = "ENDED"
	post.Props["attendents"] = strings.Join(meetingpointer.AttendeeNames, ",")
	timediff := meetingpointer.EndedAt - meetingpointer.CreatedAt
	durationstring := FormatSeconds(timediff)
	post.Props["duration"] = durationstring

	attachments := []*model.SlackAttachment{
		{
			Title: "**Meeting Ended**",
			Fields: []*model.SlackAttachmentField{
				{
					Title: "Date Started At",
					Value: time.Unix(meetingpointer.CreatedAt, 0).Format("Jan _2 at 3:04 PM"),
					Short: true,
				},
				{
					Title: "Duration",
					Value: FormatSeconds(timediff),
					Short: false,
				},
				{
					Title: "Attendees",
					Value: "@" + strings.Join(meetingpointer.AttendeeNames, ", @"),
					Short: false,
				},
			},
		},
	}

	model.ParseSlackAttachment(post, attachments)

	if _, err := p.API.UpdatePost(post); err != nil {
		p.API.LogError(fmt.Sprintf("Unable to update post. Error: {%s}", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

//when user clicks endmeeting button inside Mattermost
func (p *Plugin) handleEndMeeting(w http.ResponseWriter, r *http.Request) {

	p.API.LogInfo("00000000000000000000000000000000000000000000000000000000000000000")
	p.API.LogInfo("00000000000000000000000000000000000000000000000000000000000000000")
	p.API.LogInfo("00000000000000000000000000000000000000000000000000000000000000000")
	p.API.LogInfo("00000000000000000000000000000000000000000000000000000000000000000")

	//for debugging
	mattermost.API.LogInfo("Processing End Meeting Request")

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var request ButtonRequestJSON
	_ = json.Unmarshal(body, &request)
	meetingID := request.MeetingId
	meetingpointer := p.FindMeeting(meetingID)

	user, _ := p.API.GetUser(request.UserId)
	username := user.Username

	if meetingpointer == nil {
		myresp := model.PostActionIntegrationResponse{
			EphemeralText: "meeting has already ended",
		}
		userJson, _ := json.Marshal(myresp)
		_, _ = w.Write(userJson)
		return
	} else {
		if _, err := bbbAPI.EndMeeting(meetingpointer.MeetingID_, meetingpointer.ModeratorPW_); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			p.API.LogError(err.Error())
			return
		}

		//for debugging
		mattermost.API.LogInfo("Meeting Ended")

		if meetingpointer.EndedAt == 0 {
			meetingpointer.EndedAt = time.Now().Unix()
		}
		if err := p.AddMeetingWaitingForRecording(meetingpointer); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			p.API.LogError(err.Error())
			return
		}

		post, err := p.API.GetPost(meetingpointer.PostId)
		if err != nil {
			http.Error(w, err.Error(), err.StatusCode)
			p.API.LogError(err.Error())
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

		//attachments := post.Attachments()
		//attachments[0].Text = ""
		//attachments[0].Fields[0].Title = "Meeting Ended"
		//attachments[0].Fields[0].Value = ""
		//attachments[0].Actions = []*model.PostAction{}
		//
		//attachments[0].Fields = append(
		//	attachments[0].Fields,
		//	&model.SlackAttachmentField{
		//		Title: "Date Started At",
		//		Value: time.Unix(meetingpointer.CreatedAt, 0).Format("Jan _2 at 3:04 PM"),
		//		Short: true,
		//	},
		//	&model.SlackAttachmentField{
		//		Title: "Duration",
		//		Value: FormatSeconds(timediff),
		//		Short: true,
		//	},
		//	&model.SlackAttachmentField{
		//		Title: "Attendees",
		//		Value: "@" + strings.Join(meetingpointer.AttendeeNames, ", @"),
		//		Short: false,
		//	},
		//)

		attachments := []*model.SlackAttachment{
			{
				Fields: []*model.SlackAttachmentField{
					{
						Title: "Meeting Ended",
						Short: false,
					},
					{
						Title: "Date Started At",
						Value: time.Unix(meetingpointer.CreatedAt, 0).Format("Jan _2 at 3:04 PM"),
						Short: true,
					},
					{
						Title: "Duration",
						Value: FormatSeconds(timediff),
						Short: false,
					},
					{
						Title: "Attendees",
						Value: "@" + strings.Join(meetingpointer.AttendeeNames, ", @"),
						Short: false,
					},
				},
			},
		}

		model.ParseSlackAttachment(post, attachments)

		p.API.LogInfo("&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&")
		p.API.LogInfo(fmt.Sprintf("%d", len(post.Attachments()[0].Actions)))
		p.API.LogInfo("&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&")

		if _, err := p.API.UpdatePost(post); err != nil {
			http.Error(w, err.Error(), err.StatusCode)
			p.API.LogError(err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (p *Plugin) handleIsMeetingRunning(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var request ButtonRequestJSON
	_ = json.Unmarshal(body, &request)
	meetingID := request.MeetingId

	resp, _ := bbbAPI.IsMeetingRunning(meetingID)

	myresp := isRunningResponseJSON{
		IsRunning: resp,
	}
	userJson, _ := json.Marshal(myresp)

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(userJson)

}

func (p *Plugin) handleRecordingReady(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (p *Plugin) handleGetAttendeesInfo(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var request AttendeesRequestJSON
	_ = json.Unmarshal(body, &request)
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
	_, _ = w.Write(userJson)
}

func (p *Plugin) handlePublishRecordings(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var request *model.PostActionIntegrationRequest
	if err := json.Unmarshal(body, &request); err != nil {
		p.API.LogError("Error occurred unmarshalling publish recording API request payload. Error: " + err.Error())
		http.Error(w, "Error: couldn't unmarshal request payload", http.StatusInternalServerError)
		return
	}

	recordid := request.Context["record_id"].(string)
	publish := request.Context["publish"].(string)
	meetingID := request.Context["meeting_id"].(string)

	meetingpointer := p.FindMeeting(meetingID)
	if meetingpointer == nil {
		http.Error(w, "Error: Cannot find the meeting_id for the recording, MeetingID#"+meetingID, http.StatusForbidden)
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

	post.AddProp("is_published", publish)

	originalAttachments := post.Attachments()
	newAttachments := []*model.SlackAttachment{
		{},
	}

	for i := range originalAttachments[0].Fields {
		field := originalAttachments[0].Fields[i]

		if field.Title != "Notes" && field.Title != "Recordings" {
			newAttachments[0].Fields = append(newAttachments[0].Fields, field)
		}
	}

	newAttachments = append(newAttachments, originalAttachments[1])
	newAttachments[1].Actions[0].Name = "Make Recording Visible"
	newAttachments[1].Actions[0].Integration.Context["publish"] = "true"

	model.ParseSlackAttachment(post, newAttachments)

	if _, err := p.API.UpdatePost(post); err != nil {
		http.Error(w, err.Error(), err.StatusCode)
		return
	}
	//update post props with new recording  status
	w.WriteHeader(http.StatusOK)
}

func (p *Plugin) handleDeleteRecordings(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var request DeleteRecordingsRequestJSON
	_ = json.Unmarshal(body, &request)
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

type isRunningResponseJSON struct {
	IsRunning bool `json:"running"`
}
